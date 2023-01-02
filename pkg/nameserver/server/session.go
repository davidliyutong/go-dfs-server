package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru/v2"
	log "github.com/sirupsen/logrus"
	v12 "go-dfs-server/pkg/dataserver/client/v1"
	"go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/utils"
	"io"
	"math"
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
	Write(data []byte, n int64, wg *sync.WaitGroup, errChan chan error) error
	Truncate(size int64) error
	SetErrorClose(err error) error
	GetTime() time.Time
	GetID() *string
	GetErrors() *[]error
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
	Path         string
	FilePath     string
	ID           string
	Mode         int
	Time         time.Time
	Offset       int64
	Opened       bool
	Error        []error
	Blob         v1.BlobMetaData
	SessionMutex *sync.RWMutex // Controls access to the session
	ChunkMutex   *sync.RWMutex // Controls access to the Buffer
	Chunks       map[int64]*ChunkBuffer
	ChunkLRU     *lru.ARCCache[int64, int64]
	BlobMutex    *sync.RWMutex
	EventGroup   *sync.WaitGroup
}

func (s *session) Truncate(size int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *session) pushChunk(chunkID int64) error {
	// no lock
	if _, ok := s.Chunks[chunkID]; !ok {
		return s.SetErrorClose(errors.New("chunk not found"))
	}

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
			return err
		}
	}

	if s.Chunks[chunkID].IsPushed() {
		return nil
	}

	localMD5, _ := utils.GetBufferMD5(s.Chunks[chunkID].Bytes())

	wg := new(sync.WaitGroup)
	wg.Add(len(clients))
	clientErrors := make([]error, len(clients))

	for idx, client := range clients {
		client := client
		idx := idx
		go func() {
			defer wg.Done()
			if needCreateChunk {
				_ = client.(v12.DataServerClient).BlobCreateChunk(s.Chunks[chunkID].path, chunkID)
			}
			remoteMD5, err := client.(v12.DataServerClient).BlobWriteChunk(s.Chunks[chunkID].path, chunkID, s.Chunks[chunkID].version, bytes.NewBuffer(s.Chunks[chunkID].Bytes()))

			if localMD5 != remoteMD5 && remoteMD5 != "" {
				err = errors.New(fmt.Sprintf("checksum %s mismatch %s", localMD5, remoteMD5))
			}
			clientErrors[idx] = err
		}()
	}

	wg.Wait()

	if utils.HasError(clientErrors) {
		return s.SetErrorClose(errors.New(fmt.Sprint("push: ", chunkID, ";", "errors: ", clientErrors)))
	} else {
		s.BlobMutex.Lock()
		s.Blob.ExtendTo(chunkID)
		s.Blob.ChunkChecksums[chunkID] = localMD5
		s.Blob.Versions[chunkID] = s.Chunks[chunkID].Version()

		if needCreateChunk {
			newClientUUIDs := make([]string, len(clients))
			for idx, client := range clients {
				newClientUUIDs[idx] = client.(v12.DataServerClient).GetUUID()
			}
			s.Blob.ChunkDistribution[chunkID] = newClientUUIDs
		}

		s.Blob.Size = utils.MaxInt64(s.Blob.Size, chunkID*v1.DefaultBlobChunkSize+int64(s.Chunks[chunkID].Position()))
		s.BlobMutex.Unlock()
		s.Chunks[chunkID].SetPushed(true)

		log.Debugln("sync: ", chunkID, ";", "errors: ", clientErrors)
	}

	return nil

}

func (s *session) pullChunk(chunkID int64) error {
	_ = s.keepBuffer(chunkID)

	if chunkID >= int64(len(s.Blob.ChunkChecksums)) {
		return nil
	}
	clientUUIDs, err := s.GetChunkDistribution(chunkID)
	if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
		return errors.New("current chunk is not present")
	} else {
		clients, err := BlobDataServerManger.GetClients(clientUUIDs)
		if err != nil {
			return errors.New("cannot get related data server")
		}
		for _, client := range clients {
			version, _, err := client.(v12.DataServerClient).BlobReadChunkMeta(s.Path, chunkID)
			if err != nil {
				// TODO: handle error
				continue
			}
			reader, err := client.(v12.DataServerClient).BlobReadChunk(s.Path, chunkID)
			if err != nil {
				// TODO: handle error
				continue
			} else {
				s.Chunks[chunkID].SetVersion(version)
				s.Chunks[chunkID].SetPushed(true)
				s.BlobMutex.Lock()
				s.Blob.Versions[chunkID] = version
				s.BlobMutex.Unlock()

				buf, _ := io.ReadAll(reader)
				_, err := s.Chunks[chunkID].WriteInPlace(buf)
				return err
			}
		}
		return errors.New("cannot read chunk from any data servers")
	}
}

func (s *session) deleteChunk(chunkID int64) error {
	// no lock

	if chunkID >= int64(len(s.Blob.ChunkChecksums)) {
		return errors.New("invalid chunkID")
	}
	errChan := make(chan error, 3)
	wg := new(sync.WaitGroup)

	clientUUIDs, err := s.GetChunkDistribution(chunkID)
	if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
		return errors.New("current chunk is not present")
	} else {
		clients, err := BlobDataServerManger.GetClients(clientUUIDs)
		if err != nil {
			return errors.New("cannot get related data server")
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
		wg.Wait()
		errs := utils.ReceiveErrors(errChan)
		if utils.HasError(errs) {
			return utils.GetFirstError(errs)
		} else {
			return nil
		}
	}
}

func (s *session) fetchChunk(chunkID int64, chunkOffset int64, buffer []byte) (int64, error) {
	// no lock
	_, err := s.Chunks[chunkID].Seek(chunkOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	n, err := s.Chunks[chunkID].Read(buffer)
	if err != nil {
		return 0, err
	} else {
		return int64(n), nil
	}
}

func (s *session) updateChunk(chunkID int64, chunkOffset int64, data []byte) error {
	// no lock
	_, err := s.Chunks[chunkID].Seek(chunkOffset, io.SeekStart)
	if err != nil {
		return err
	}

	if s.Chunks[chunkID].Position()+len(data) > v1.DefaultBlobChunkSize {
		return errors.New("buffer overflow")
	} else {
		_, err := s.Chunks[chunkID].WriteInPlace(data)
		return err
	}
}

func (s *session) keepBuffer(chunkID int64) error {
	// no lock
	if _, ok := s.Chunks[chunkID]; ok {
		return nil
	} else {
		s.Chunks[chunkID] = NewChunkBuffer(s.Path, s.GetChunkID(), 0, make([]byte, 0))
		return nil
	}
}

func (s *session) freeBuffer() {
	// no lock
	go func() {
		if len(s.Chunks) <= 4 {
			return
		}
		currChunkID := s.GetChunkID()
		for k, v := range s.Chunks {
			if math.Abs(float64(k)-float64(currChunkID)) > 1 && v.IsPushed() {
				delete(s.Chunks, k)
			}
		}
	}()
}

func (s *session) Write(data []byte, n int64, wg *sync.WaitGroup, errChan chan error) error {
	// SessionMutex: RLock

	s.SessionMutex.RLock()
	defer s.SessionMutex.RUnlock()
	if !s.IsOpened() {
		return errors.New("session closed")
	}
	s.EventGroup.Add(1)
	defer s.EventGroup.Done()

	defer wg.Done()

	var numBytesWritten int64 = 0
	numBytesLeft := n

	for numBytesLeft > 0 {
		// update offset
		s.BlobMutex.Lock()
		currChunkID := s.GetChunkID()
		currChunkOffset := s.GetChunkOffset()
		currOffset := *s.GetOffset()
		boundary := (currChunkID + 1) * v1.DefaultBlobChunkSize
		var numBytesToWrite int64
		if currOffset+numBytesLeft >= boundary {
			numBytesToWrite = boundary - currOffset
		} else {
			numBytesToWrite = numBytesLeft
		}
		numBytesLeft -= numBytesToWrite
		numBytesWritten += numBytesToWrite
		s.Offset += numBytesToWrite
		s.BlobMutex.Unlock()

		s.ChunkMutex.Lock()
		if _, ok := s.Chunks[currChunkID]; !ok {
			err := s.pullChunk(currChunkID)
			if err != nil {
				errChan <- s.SetErrorClose(err)
				return err
			}
		}
		s.ChunkMutex.Unlock()

		err := s.updateChunk(currChunkID, currChunkOffset, data[numBytesWritten-numBytesToWrite:numBytesWritten])
		if err != nil {
			errChan <- s.SetErrorClose(err)
			return err
		}

		s.ChunkMutex.RLock()
		err = s.pushChunk(currChunkID)
		s.ChunkMutex.RUnlock()

		if err != nil {
			errChan <- s.SetErrorClose(err)
			return err
		} else {
			errChan <- nil
			s.ChunkMutex.Lock()
			s.freeBuffer()
			s.ChunkMutex.Unlock()
		}
	}

	err := s.DumpBlobMetaData()
	if err != nil {
		errChan <- err
	}

	return nil
}

func (s *session) createFile() error {
	// layout file on data servers
	clients := BlobDataServerManger.GetAllClients()
	if len(clients) < 3 {
		return s.SetErrorClose(errors.New("no enough data server available"))
	}
	clientErrors := make([]error, len(clients))
	wg := sync.WaitGroup{}
	wg.Add(len(clients) + 1)
	for idx, client := range clients {
		idx := idx
		client := client
		go func() {
			defer wg.Done()
			err := client.(v12.DataServerClient).BlobCreateFile(s.Path)
			if err != nil {
				clientErrors[idx] = err
			} else {
				clientErrors[idx] = nil
			}
		}()
	}

	// create file meta data on name server
	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = os.Mkdir(s.FilePath, 0775)
		s.SetBlobMetaData(v1.NewBlobMetaData(v1.BlobFileTypeName, filepath.Base(s.Path)))
		err2 = s.DumpBlobMetaData()
	}()

	wg.Wait()
	if utils.HasError(clientErrors) {
		for _, err := range clientErrors {
			_ = s.SetErrorClose(err)
		}
		return s.SetErrorClose(errors.New("create file failed, dataserver error"))
	} else if err1 != nil || err2 != nil {
		_ = s.SetErrorClose(err1)
		_ = s.SetErrorClose(err2)
		return s.SetErrorClose(errors.New("create file failed, local error"))
	}
	return nil
}

func (s *session) Open() error {
	// check opened
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()
	if s.Opened {
		return errors.New("session already opened")
	}
	if utils.HasError(s.Error) {
		return errors.New("session has unresolved error")
	}

	// prevent session to be closed
	s.EventGroup.Add(1)
	defer s.EventGroup.Done()

	// check path
	filePath, err := utils.JoinSubPathSafe(GlobalServerDesc.Opt.Volume, s.Path)
	if err != nil {
		return s.SetErrorClose(err)
	} else {
		s.FilePath = filePath
	}

	// check metadata and type of path
	var pathExists = utils.PathExists(s.FilePath)
	var isValidFile = utils.GetFileState(s.FilePath)

	if pathExists && !isValidFile {
		return s.SetErrorClose(errors.New("path is not a file"))
	}

	if isValidFile {
		err := s.LoadBlobMetaData() // this will lock BlobMutex
		if err != nil {
			return s.SetErrorClose(err)
		}
	}

	// open session according to mode
	switch s.Mode {
	case os.O_RDONLY:
		if !isValidFile {
			return s.SetErrorClose(errors.New("file does not exist or invalid"))
		}

	case os.O_RDWR:
		if !isValidFile {
			err = s.createFile() // this will close session on failure
			if err != nil {
				return err
			}
		}

	default:
		return s.SetErrorClose(errors.New("invalid mode"))

	}
	s.Opened = true
	return nil
}

func (s *session) IsOpened() bool {
	return s.Opened
}

func (s *session) Seek(offset int64, whence int) (int64, error) {
	// SessionMutex: Rlock
	// ChunkMutex: Lock

	s.SessionMutex.RLock()
	defer s.SessionMutex.RUnlock()
	if !s.IsOpened() {
		return -1, errors.New("session closed")
	}
	s.EventGroup.Add(1)
	defer s.EventGroup.Done()

	var newOffset int64
	if whence == io.SeekStart {
		newOffset = offset
	} else if whence == io.SeekCurrent {
		newOffset += offset
	} else if whence == io.SeekEnd {
		newOffset = s.Blob.Size - 1 - offset
	}

	s.BlobMutex.Lock()
	if newOffset < 0 || newOffset > s.Blob.Size {
		s.BlobMutex.Unlock()
		return 0, errors.New("invalid offset")
	} else {
		s.Offset = newOffset
		s.BlobMutex.Unlock()
		s.ChunkMutex.Lock()
		go func() {
			defer s.ChunkMutex.Unlock()
			currChunkID := s.GetChunkID()
			err := s.pullChunk(currChunkID)
			if err != nil {
				_ = s.SetErrorClose(err)
				return
			}
			s.freeBuffer()
		}()

		return s.Offset, nil
	}
}

func (s *session) flush() error {
	return nil
}
func (s *session) Close() error {
	// Session-level lock
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()
	if !s.IsOpened() {
		return errors.New("session closed")
	}

	err := s.flush()
	if err != nil {
		_ = s.SetErrorClose(err)
		return err
	} else {
		s.Opened = false
		_ = s.SetErrorClose(nil)
		s.Chunks = make(map[int64]*ChunkBuffer)
		return nil
	}
}

func (s *session) Flush() error {
	// SessionMutex: Lock
	// ChunkMutex: Lock

	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()
	if !s.IsOpened() {
		return errors.New("session closed")
	}

	s.EventGroup.Wait()

	switch s.Mode {
	case os.O_RDONLY:
		return nil
	case os.O_RDWR:
		s.ChunkMutex.Lock()
		defer s.ChunkMutex.Unlock()

		for chunkID := range s.Chunks {
			err := s.pushChunk(chunkID)
			if err != nil && err.Error() != "chunk already pushed" {
				return s.SetErrorClose(err)
			}
		}

		s.freeBuffer()

		maxChunkID := utils.GetChunkID(s.Blob.Size)
		for idx, _ := range s.Blob.ChunkChecksums {
			if int64(idx) > maxChunkID {
				err := s.deleteChunk(int64(idx))
				if err != nil {
					return s.SetErrorClose(err)
				}
			}
		}

		s.TruncateTo(maxChunkID + 1)
		err := s.DumpBlobMetaData()
		if err != nil {
			return s.SetErrorClose(err)
		} else {
			return nil
		}

	default:
		return s.SetErrorClose(errors.New("invalid mode"))
	}
}

func (s *session) Read(buffer []byte, size int64) (int64, error) {
	// SessionMutex: RLock

	s.SessionMutex.RLock()
	defer s.SessionMutex.RUnlock()
	if !s.IsOpened() {
		return -1, errors.New("session closed")
	}
	s.EventGroup.Add(1)
	defer s.EventGroup.Done()

	var numBytesRead int64 = 0
	numBytesLeft := utils.MinInt64(int64(len(buffer)), size)
	s.ChunkMutex.Lock()
	defer s.ChunkMutex.Unlock()

	for numBytesLeft > 0 {
		s.BlobMutex.Lock()
		currChunkID := s.GetChunkID()
		currChunkOffset := s.GetChunkOffset()
		currOffset := *s.GetOffset()
		boundary := utils.MinInt64(s.Blob.Size, (currChunkID+1)*v1.DefaultBlobChunkSize)
		var numBytesToRead int64
		if currOffset+numBytesLeft >= boundary {
			numBytesToRead = boundary - currOffset
		} else {
			numBytesToRead = numBytesLeft
		}
		numBytesRead += numBytesToRead
		numBytesLeft -= numBytesToRead
		s.Offset += numBytesToRead
		s.BlobMutex.Unlock()

		if numBytesToRead <= 0 {
			return numBytesRead, io.EOF
		}

		if _, ok := s.Chunks[currChunkID]; !ok {
			err := s.pullChunk(currChunkID)
			if err != nil {
				return numBytesRead, s.SetErrorClose(err)
			}
		}
		_, err := s.fetchChunk(currChunkID, currChunkOffset, buffer[numBytesRead-numBytesToRead:numBytesRead+numBytesLeft])
		if err != nil {
			return numBytesRead, s.SetErrorClose(err)
		}

		s.freeBuffer()

	}
	return numBytesRead, nil
}

func (s *session) SetErrorClose(err error) error {
	if err != nil {
		log.Errorf("session %v closed due to error: %v", s.GetID(), err)
	} else {
		log.Debugf("session %v closed", s.GetID())
	}
	s.Error = append(s.Error, err)
	s.Opened = false
	return err
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

func (s *session) GetErrors() *[]error {
	return &s.Error
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

func (s *session) GetNumOfChunks() int64 {
	s.BlobMutex.RLock()
	defer s.BlobMutex.RUnlock()
	return s.Blob.GetNumOfChunks()
}

func (s *session) GetChunkDistribution(chunkID int64) ([]string, error) {
	s.BlobMutex.RLock()
	defer s.BlobMutex.RUnlock()
	return s.Blob.GetChunkDistribution(chunkID)
}

func (s *session) ExtendTo(chunkID int64) {
	s.BlobMutex.Lock()
	defer s.BlobMutex.Unlock()
	s.Blob.ExtendTo(chunkID)
}

func (s *session) TruncateTo(chunkID int64) {
	s.BlobMutex.Lock()
	defer s.BlobMutex.Unlock()
	s.Blob.TruncateTo(chunkID)
}

func (s *session) GetBlobMetaData() *v1.BlobMetaData {
	s.BlobMutex.RLock()
	defer s.BlobMutex.RUnlock()
	return &s.Blob
}

func (s *session) SetBlobMetaData(blob v1.BlobMetaData) {
	s.BlobMutex.Lock()
	defer s.BlobMutex.Unlock()
	s.Blob = blob
}

func (s *session) DumpBlobMetaData() error {
	s.BlobMutex.Lock()
	defer s.BlobMutex.Unlock()
	return s.Blob.Dump(s.GetMetaFilePath())
}

func (s *session) LoadBlobMetaData() error {
	s.BlobMutex.RLock()
	defer s.BlobMutex.RUnlock()
	return s.Blob.Load(s.GetMetaFilePath())
}

func NewSession(path string, filePath string, id string, mode int) Session {
	newCache, _ := lru.NewARC[int64, int64](32)
	return &session{
		Path:         path,
		FilePath:     filePath,
		ID:           id,
		Mode:         mode,
		Time:         time.Now(),
		Offset:       0,
		Opened:       false,
		Error:        make([]error, 0),
		Blob:         v1.BlobMetaData{},
		SessionMutex: new(sync.RWMutex),
		ChunkMutex:   new(sync.RWMutex),
		Chunks:       make(map[int64]*ChunkBuffer),
		ChunkLRU:     newCache,
		BlobMutex:    new(sync.RWMutex),
		EventGroup:   new(sync.WaitGroup),
	}
}
