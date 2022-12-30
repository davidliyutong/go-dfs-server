package server

import (
	"bytes"
	json "encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	v12 "go-dfs-server/pkg/dataserver/client/v1"
	"go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Session interface {
	Open() error
	IsOpened() bool
	Seek(offset int64, whence int) (int64, error)
	Close() error
	Flush() error
	Read(buffer []byte, size int64) (int64, error)
	Write(data []byte, n int64) error
	Truncate(size int64) error
	GetTime() time.Time
	GetID() *string
	GetChunkID() int64
	GetMode() *int
	GetPath() *string
	GetFilePath() *string
	GetOffset() *int64
	GetMetaFilePath() string
	GetBlobMetaData() *v1.BlobMetaData
	SetBlobMetaData(blob v1.BlobMetaData)
	DumpBlobMetaData() error
	LoadBlobMetaData() error
}

type session struct {
	Path             string
	FilePath         string
	ID               string
	Mode             int
	Version          int64
	Time             time.Time
	Offset           int64
	Opened           bool
	Blob             v1.BlobMetaData
	SessionMutex     *sync.RWMutex           // Controls access to the session
	ChunkMutex       *sync.RWMutex           // Controls access to the chunk
	ChunkBufferMutex map[int64]*sync.RWMutex // Controls access to the chunk buffer
	ChunkBuffer      map[int64]*ChunkBuffer
	TransferMutex    *sync.RWMutex // Controls access to the chunk transfer
}

func (s *session) Truncate(size int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *session) pushChunk(chunkID int64) (*sync.WaitGroup, chan error, error) {
	if !s.IsOpened() {
		return nil, nil, errors.New("session closed")
	}

	s.ChunkMutex.RLock()
	defer s.ChunkMutex.RUnlock()

	wg := new(sync.WaitGroup)

	if _, ok := s.ChunkBufferMutex[chunkID]; !ok {
		//s.Opened = false
		return nil, nil, errors.New("chunk not found")
	}

	s.ChunkBufferMutex[chunkID].RLock()
	defer s.ChunkBufferMutex[chunkID].RUnlock()

	if s.ChunkBuffer[chunkID].pushed {
		//s.Opened = false
		return wg, nil, errors.New("chunk already pushed")
	}
	localMD5, _ := utils.GetBufferMD5(s.ChunkBuffer[chunkID].buffer.Bytes())

	var clients []interface{}
	var needCreateChunk = false
	oldClientUUIDs, err := s.Blob.GetChunkDistribution(chunkID)
	if err != nil || len(oldClientUUIDs) == 0 || oldClientUUIDs == nil {
		clients = utils.SelectRandomNFromArray(BlobDataServerManger.GetAllClients(), NameServerNumOfReplicas)
		needCreateChunk = true
	} else {
		clients, err = BlobDataServerManger.GetClients(oldClientUUIDs)
		if err != nil {
			log.Warningln(err.Error())
		}
	}

	errChan := make(chan error, len(clients)*2)
	wg.Add(len(clients) + 1)
	go func() {
		clientErrors := make([]error, 0)
		for _, client := range clients {
			if needCreateChunk {
				_ = client.(v12.DataServerClient).BlobCreateChunk(s.ChunkBuffer[chunkID].path, chunkID)
			}

			remoteMD5, err := client.(v12.DataServerClient).BlobWriteChunk(s.ChunkBuffer[chunkID].path, chunkID, s.ChunkBuffer[chunkID].version, bytes.NewBuffer(s.ChunkBuffer[chunkID].buffer.Bytes()))
			if err != nil {
				s.Opened = false
				errChan <- err
				clientErrors = append(clientErrors, err)
				wg.Done()
				continue
			}

			if localMD5 != remoteMD5 {
				s.Opened = false
				err = errors.New(fmt.Sprintf("checksum %s mismatch %s", localMD5, remoteMD5))
			} else {
				err = nil
			}
			errChan <- err
			clientErrors = append(clientErrors, err)
			wg.Done()
			continue
		}

		if utils.HasError(clientErrors) {
			log.Warningln("push: ", chunkID, ";", "errors: ", clientErrors)
		}

		s.Blob.ExtendTo(chunkID)
		s.Blob.ChunkChecksums[chunkID] = localMD5
		s.Blob.Versions[chunkID] = s.ChunkBuffer[chunkID].version
		s.ChunkBuffer[chunkID].pushed = true
		if needCreateChunk {
			newClientUUIDs := make([]string, len(clients))
			for idx, client := range clients {
				newClientUUIDs[idx] = client.(v12.DataServerClient).GetUUID()
			}
			s.Blob.ChunkDistribution[chunkID] = newClientUUIDs
		}
		s.Blob.Size = utils.MaxInt64(s.Blob.Size, chunkID*v1.DefaultBlobChunkSize+int64(s.ChunkBuffer[chunkID].buffer.Position()))

		log.Debugln("sync: ", chunkID, ";", "errors: ", clientErrors)
		errChan <- s.DumpBlobMetaData()
		wg.Done()
		close(errChan)
	}()

	return wg, errChan, nil

}

func (s *session) pullChunk(chunkID int64) (*sync.WaitGroup, chan error, error) {
	if !s.IsOpened() {
		return nil, nil, errors.New("session closed")
	}

	chunkOffset := s.GetChunkOffset()
	_ = s.keepBuffer(chunkID)

	errChan := make(chan error, 10)
	wg := new(sync.WaitGroup)
	if chunkID >= int64(len(s.Blob.ChunkChecksums)) {
		close(errChan)
		return wg, errChan, nil
	}
	clientUUIDs, err := s.Blob.GetChunkDistribution(chunkID)
	if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
		close(errChan)
		return nil, nil, errors.New("current chunk is not present")
	} else {
		clients, err := BlobDataServerManger.GetClients(clientUUIDs)
		if err != nil {
			close(errChan)
			return nil, nil, errors.New("cannot get related data server")
		}
		for _, client := range clients {
			version, _, err := client.(v12.DataServerClient).BlobReadChunkMeta(s.Path, chunkID)
			if err != nil {
				continue
			}
			reader, err := client.(v12.DataServerClient).BlobReadChunk(s.Path, chunkID)
			if err != nil {
				continue
			} else {
				wg.Add(1)
				go func() {
					s.ChunkBufferMutex[chunkID].Lock()
					s.ChunkBuffer[chunkID].version = version
					s.Blob.Versions[chunkID] = version
					defer wg.Done()
					defer s.ChunkBufferMutex[chunkID].Unlock()
					defer close(errChan)

					buf, _ := io.ReadAll(reader)
					_, err := s.ChunkBuffer[chunkID].WriteInPlace(buf)
					if err != nil {
						s.Opened = false
						errChan <- err
						return
					}
					_, err = s.ChunkBuffer[chunkID].buffer.Seek(chunkOffset, io.SeekStart)
					if err != nil {
						s.Opened = false
						errChan <- err
						return
					}
				}()
				return wg, errChan, nil
			}
		}
		close(errChan)
		return nil, nil, errors.New("cannot read chunk")
	}
}

func (s *session) deleteChunk(chunkID int64) (*sync.WaitGroup, chan error, error) {
	if !s.IsOpened() {
		return nil, nil, errors.New("session closed")
	}

	if chunkID >= int64(len(s.Blob.ChunkChecksums)) {
		return nil, nil, errors.New("invalid chunkID")
	}
	errChan := make(chan error, 1)
	wg := new(sync.WaitGroup)

	clientUUIDs, err := s.Blob.GetChunkDistribution(chunkID)
	if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
		return nil, nil, errors.New("current chunk is not present")
	} else {
		clients, err := BlobDataServerManger.GetClients(clientUUIDs)
		if err != nil {
			return nil, nil, errors.New("cannot get related data server")
		}
		wg.Add(len(clients))
		go func() {
			defer close(errChan)
			for _, client := range clients {
				err := client.(v12.DataServerClient).BlobDeleteChunk(s.Path, chunkID)
				errChan <- err
				wg.Done()
			}
		}()

		return wg, errChan, nil
	}
}

func (s *session) freeBuffer() {
	currChunkID := s.GetChunkID()
	for k, v := range s.ChunkBuffer {
		if k != currChunkID && v.IsPushed() {
			s.dropBuffer(k)
		}
	}
}

func (s *session) updateWriteBuffer(chunkID int64, chunkOffset int64, data []byte) error {
	if _, ok := s.ChunkBuffer[chunkID]; !ok {
		wg, _, err := s.pullChunk(chunkID)
		if err != nil {
			return err
		}
		wg.Wait()
	}

	s.ChunkBufferMutex[chunkID].Lock()
	defer s.ChunkBufferMutex[chunkID].Unlock()

	_, err := s.ChunkBuffer[chunkID].buffer.Seek(chunkOffset, io.SeekStart)
	if err != nil {
		return err
	}

	if s.ChunkBuffer[chunkID].buffer.Position()+len(data) > v1.DefaultBlobChunkSize {
		return errors.New("buffer overflow")
	} else {
		_, err := s.ChunkBuffer[chunkID].WriteInPlace(data)
		return err
	}
}

func (s *session) Write(data []byte, n int64) error {
	if !s.IsOpened() {
		return errors.New("session closed")
	}

	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()
	var numBytesWritten int64 = 0
	numBytesLeft := n
	for numBytesLeft > 0 {
		currChunkID := s.GetChunkID()
		currChunkOffset := s.GetChunkOffset()
		boundary := (currChunkID + 1) * v1.DefaultBlobChunkSize
		var numBytesToWrite int64
		if s.Offset+numBytesLeft >= boundary {
			numBytesToWrite = boundary - s.Offset
		} else {
			numBytesToWrite = numBytesLeft
		}

		err := s.updateWriteBuffer(currChunkID, currChunkOffset, data[numBytesWritten:numBytesWritten+numBytesToWrite])
		if err != nil {
			return err
		}
		wg, _, err := s.pushChunk(currChunkID)
		if err != nil {
			return err
		}
		wg.Wait()
		//s.freeBuffer()
		numBytesLeft -= numBytesToWrite
		numBytesWritten += numBytesToWrite

		s.Offset += numBytesToWrite
	}
	return nil
}

func (s *session) initBuffer() {
	if s.ChunkBuffer == nil {
		s.ChunkBuffer = make(map[int64]*ChunkBuffer)
	}
	if s.ChunkBufferMutex == nil {
		s.ChunkBufferMutex = make(map[int64]*sync.RWMutex)
	}
	if s.ChunkMutex == nil {
		s.ChunkMutex = new(sync.RWMutex)
	}
}

func (s *session) keepBuffer(chunkID int64) error {
	if _, ok := s.ChunkBuffer[chunkID]; ok {
		return nil
	} else {
		s.ChunkMutex.Lock()
		s.ChunkBuffer[chunkID] = NewChunkBuffer(s.Path, s.GetChunkID(), 0, make([]byte, 0))
		s.ChunkBufferMutex[chunkID] = new(sync.RWMutex)
		s.ChunkMutex.Unlock()
		return nil
	}
}

func (s *session) dropBuffer(chunkID int64) {
	if _, ok := s.ChunkBufferMutex[chunkID]; ok {
		s.ChunkMutex.Lock()

		s.ChunkBufferMutex[chunkID].Lock()
		delete(s.ChunkBuffer, chunkID)
		s.ChunkBufferMutex[chunkID].Unlock()
		delete(s.ChunkBufferMutex, chunkID)

		s.ChunkMutex.Unlock()
	}
}

func (s *session) Open() error {
	// check path
	filePath, err := utils.JoinSubPathSafe(GlobalServerDesc.Opt.Volume, s.Path)
	if err != nil {
		return err
	} else {
		s.FilePath = filePath
	}

	var pathExists = utils.PathExists(s.FilePath)
	var isValidFile = utils.GetFileState(s.FilePath)

	if pathExists && !isValidFile {
		return errors.New("path is not a file")
	}

	if isValidFile {
		err := s.LoadBlobMetaData()
		if err != nil {
			return err
		}
	}

	switch s.Mode {
	case os.O_RDONLY:
		if !isValidFile {
			return errors.New("file does not exist or invalid")
		}

	case os.O_RDWR:
		if !isValidFile {
			clients := BlobDataServerManger.GetAllClients()
			if len(clients) < 3 {
				return errors.New("no enough data server available")
			}
			clientErrors := make([]error, len(clients))
			wg := sync.WaitGroup{}
			wg.Add(len(clients))
			for idx, client := range clients {
				idx := idx
				client := client
				go func() {
					err := client.(v12.DataServerClient).BlobCreateFile(s.Path)
					if err != nil && err.Error() != "file or directory exists" {
						clientErrors[idx] = err
					} else {
						clientErrors[idx] = nil
					}
					wg.Done()
				}()
			}
			wg.Wait()
			err1 := os.Mkdir(filePath, 0775)
			s.SetBlobMetaData(v1.NewBlobMetaData(v1.BlobFileTypeName, filepath.Base(s.Path)))
			err2 := s.DumpBlobMetaData()

			if utils.HasError(clientErrors) || err1 != nil || err2 != nil {
				wg.Add(len(clients))
				for _, client := range clients {
					client := client
					go func() { _ = client.(v12.DataServerClient).BlobDeleteFile(s.Path); wg.Done() }()
				}
				wg.Wait()
				return errors.New("create file failed")
			} else {
				return nil
			}
		}

	default:
		return errors.New("invalid mode")

	}
	s.Opened = true
	if isValidFile {
		_, _, err = s.pullChunk(0)
		if err != nil {
			s.Opened = false
			return errors.New("failed to pull read buffer")
		}
	}

	return nil
}

func (s *session) IsOpened() bool {
	return s.Opened
}

func (s *session) Seek(offset int64, whence int) (int64, error) {
	if !s.IsOpened() {
		return -1, errors.New("session closed")
	}
	err := s.Flush()
	if err != nil {
		return -1, err
	}
	var newOffset int64
	if whence == io.SeekStart {
		newOffset = offset
	} else if whence == io.SeekCurrent {
		newOffset += offset
	} else if whence == io.SeekEnd {
		newOffset = s.Blob.Size + offset
	}
	if newOffset < 0 || newOffset > s.Blob.Size {
		return 0, errors.New("invalid offset")
	} else {
		s.Offset = newOffset
		currChunkID := s.GetChunkID()
		_, _, err := s.pullChunk(currChunkID)
		if err != nil {
			return -1, err
		}
		return s.Offset, nil
	}
}

func (s *session) Close() error {
	if !s.IsOpened() {
		return errors.New("session closed")
	}
	err := s.Flush()
	if err != nil {
		return err
	}

	s.SessionMutex.Lock()
	s.Opened = false
	s.SessionMutex.Unlock()

	s.SessionMutex = nil
	s.ChunkMutex = nil
	s.ChunkBufferMutex = nil
	s.ChunkBuffer = nil

	return nil
}

func (s *session) Flush() error {
	if !s.IsOpened() {
		return errors.New("session closed")
	}
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	switch s.Mode {
	case os.O_RDONLY:
		return nil
	case os.O_RDWR:
		for chunkID := range s.ChunkBuffer {
			wg, _, err := s.pushChunk(chunkID)
			wg.Wait()
			if err != nil && err.Error() != "chunk already pushed" {
				return err
			}
		}
		s.freeBuffer()
		maxChunkID := utils.GetChunkID(s.Blob.Size)
		for idx, _ := range s.Blob.ChunkChecksums {
			if int64(idx) > maxChunkID {
				wg, _, err := s.deleteChunk(int64(idx))
				wg.Wait()
				if err != nil {
					continue
				}
			}
		}
		s.Blob.TruncateTo(maxChunkID)

		return s.DumpBlobMetaData()
	default:
		return errors.New("invalid mode")
	}
}

func (s *session) fetchReadBuffer(chunkID int64, chunkOffset int64, buffer []byte) (int64, error) {
	if _, ok := s.ChunkBuffer[chunkID]; !ok {
		wg, _, err := s.pullChunk(chunkID)
		if err != nil {
			return 0, err
		}
		wg.Wait()
	}

	s.ChunkBufferMutex[chunkID].RLock()
	defer s.ChunkBufferMutex[chunkID].RUnlock()

	_, err := s.ChunkBuffer[chunkID].buffer.Seek(chunkOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	n, err := s.ChunkBuffer[chunkID].Read(buffer)
	if err != nil {
		return 0, err
	} else {
		return int64(n), nil
	}
}

func (s *session) Read(buffer []byte, size int64) (int64, error) {
	if !s.IsOpened() {
		return 0, errors.New("session closed")
	}

	s.SessionMutex.RLock()
	defer s.SessionMutex.RUnlock()
	var numBytesRead int64 = 0
	numBytesLeft := utils.MinInt64(int64(len(buffer)), size)
	for numBytesLeft > 0 {
		currChunkID := s.GetChunkID()
		currChunkOffset := s.GetChunkOffset()

		n, err := s.fetchReadBuffer(currChunkID, currChunkOffset, buffer[numBytesRead:numBytesRead+numBytesLeft])
		if err != nil {
			return numBytesRead, err
		}
		numBytesRead += n
		numBytesLeft -= n
		s.Offset += n
	}
	return numBytesRead, nil
}

func (s *session) GetTime() time.Time {
	return s.Time
}

func (s *session) GetID() *string {
	return &s.ID
}

func (s *session) GetChunkID() int64 {
	return s.Offset / v1.DefaultBlobChunkSize
}

func (s *session) GetChunkOffset() int64 {
	return s.Offset % v1.DefaultBlobChunkSize
}

func (s *session) GetMode() *int {
	return &s.Mode
}

func (s *session) GetPath() *string {
	return &s.Path
}

func (s *session) GetFilePath() *string {
	return &s.FilePath
}

func (s *session) GetOffset() *int64 {
	return &s.Offset
}

func (s *session) GetMetaFilePath() string {
	return filepath.Join(s.FilePath, "meta.json")
}

func (s *session) GetBlobMetaData() *v1.BlobMetaData {
	return &s.Blob
}

func (s *session) SetBlobMetaData(blob v1.BlobMetaData) {
	s.Blob = blob
}

func (s *session) DumpBlobMetaData() error {
	filePtr, err := os.Create(s.GetMetaFilePath())
	if err != nil {
		return err
	}
	defer func(filePtr *os.File) {
		err := filePtr.Close()
		if err != nil {

		}
	}(filePtr)
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(s.Blob)
	return err
}

func (s *session) LoadBlobMetaData() error {
	jsonFile, err := os.Open(s.GetMetaFilePath())
	if err != nil {
		return errors.New("cannot open metadata")
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	buffer, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(buffer, &s.Blob)
	return err
}

func NewSession(path string, filePath string, id string, mode int) Session {
	return &session{
		Path:             path,
		FilePath:         filePath,
		ID:               id,
		Mode:             mode,
		Time:             time.Now(),
		Offset:           0,
		Opened:           true,
		Blob:             v1.BlobMetaData{},
		SessionMutex:     new(sync.RWMutex),
		ChunkMutex:       new(sync.RWMutex),
		ChunkBufferMutex: make(map[int64]*sync.RWMutex),
		ChunkBuffer:      make(map[int64]*ChunkBuffer),
		TransferMutex:    new(sync.RWMutex),
	}
}
