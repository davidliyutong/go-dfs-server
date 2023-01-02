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
	setErrorClose(err error) error
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
}

func (s *session) Truncate(size int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *session) pushChunk(chunkID int64) error {
	// no lock
	if !s.IsOpened() {
		return errors.New("session closed")
	}

	if _, ok := s.Chunks[chunkID]; !ok {
		return s.setErrorClose(errors.New("chunk not found"))
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
	wg.Add(len(clients) + 1)
	go func() {
		clientErrors := make([]error, 0)

		for _, client := range clients {
			if needCreateChunk {
				_ = client.(v12.DataServerClient).BlobCreateChunk(s.Chunks[chunkID].path, chunkID)
			}
			remoteMD5, err := client.(v12.DataServerClient).BlobWriteChunk(s.Chunks[chunkID].path, chunkID, s.Chunks[chunkID].version, bytes.NewBuffer(s.Chunks[chunkID].Bytes()))

			if localMD5 != remoteMD5 && remoteMD5 != "" {
				err = errors.New(fmt.Sprintf("checksum %s mismatch %s", localMD5, remoteMD5))
			}
			clientErrors = append(clientErrors, err)
			wg.Done()
			continue
		}

		if utils.HasError(clientErrors) {
			_ = s.setErrorClose(errors.New(fmt.Sprint("push: ", chunkID, ";", "errors: ", clientErrors)))
			wg.Done()
			return
		} else {
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
			s.Chunks[chunkID].SetPushed(true)

			log.Debugln("sync: ", chunkID, ";", "errors: ", clientErrors)
			wg.Done()
		}

	}()
	wg.Wait()

	return nil

}

func (s *session) pullChunk(chunkID int64) error {
	// no Lock
	if !s.IsOpened() {
		return errors.New("session closed")
	}

	chunkOffset := s.GetChunkOffset()
	_ = s.keepBuffer(chunkID)

	if chunkID >= int64(len(s.Blob.ChunkChecksums)) {
		return nil
	}
	clientUUIDs, err := s.Blob.GetChunkDistribution(chunkID)
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
				s.Chunks[chunkID].SetPushed(false)
				s.Blob.Versions[chunkID] = version

				buf, _ := io.ReadAll(reader)
				_, err := s.Chunks[chunkID].WriteInPlace(buf)
				if err != nil {
					return err
				}
				_, err = s.Chunks[chunkID].Seek(chunkOffset, io.SeekStart)
				return err
			}
		}

		return errors.New("cannot read chunk from any data servers")
	}
}

func (s *session) deleteChunk(chunkID int64) error {
	// no lock
	if !s.IsOpened() {
		return errors.New("session closed")
	}

	if chunkID >= int64(len(s.Blob.ChunkChecksums)) {
		return errors.New("invalid chunkID")
	}
	errChan := make(chan error, 3)
	wg := new(sync.WaitGroup)

	clientUUIDs, err := s.Blob.GetChunkDistribution(chunkID)
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

func (s *session) updateChunk(chunkID int64, chunkOffset int64, data []byte) error {
	// Chunk RLock
	if _, ok := s.Chunks[chunkID]; !ok {
		err := s.pullChunk(chunkID)
		if err != nil {
			return err
		}
	}

	s.ChunkMutex.Lock()
	defer s.ChunkMutex.Unlock()

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

func (s *session) Write(data []byte, n int64, wg *sync.WaitGroup, errChan chan error) error {
	// Session-Level Lock
	if !s.IsOpened() {
		return errors.New("session closed")
	}

	var numBytesWritten int64 = 0
	numBytesLeft := n

	s.SessionMutex.Lock()

	defer s.SessionMutex.Unlock()
	defer wg.Done()

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

		err := s.updateChunk(currChunkID, currChunkOffset, data[numBytesWritten:numBytesWritten+numBytesToWrite])
		if err != nil {
			errChan <- s.setErrorClose(err)
			return err
		}
		s.ChunkMutex.RLock()
		err = s.pushChunk(currChunkID)
		s.ChunkMutex.RUnlock()
		if err != nil {
			errChan <- s.setErrorClose(err)
			return err
		} else {
			numBytesLeft -= numBytesToWrite
			numBytesWritten += numBytesToWrite
			s.Offset += numBytesToWrite
			errChan <- nil
			s.freeBuffer()
		}
	}

	err := s.DumpBlobMetaData()
	if err != nil {
		errChan <- err
	}

	return nil
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
	// Chunk Lock
	s.ChunkMutex.Lock()
	go func() {
		defer s.ChunkMutex.Unlock()
		currChunkID := s.GetChunkID()
		for k, v := range s.Chunks {
			if math.Abs(float64(k)-float64(currChunkID)) > 1 && v.IsPushed() {
				delete(s.Chunks, k)
			}
		}
	}()
}

func (s *session) Open() error {
	// check path
	if s.Opened {
		return errors.New("session already opened")
	}

	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	filePath, err := utils.JoinSubPathSafe(GlobalServerDesc.Opt.Volume, s.Path)
	if err != nil {
		return s.setErrorClose(err)
	} else {
		s.FilePath = filePath
	}

	var pathExists = utils.PathExists(s.FilePath)
	var isValidFile = utils.GetFileState(s.FilePath)

	if pathExists && !isValidFile {
		return s.setErrorClose(errors.New("path is not a file"))
	}

	if isValidFile {
		err := s.LoadBlobMetaData()
		if err != nil {
			return s.setErrorClose(err)
		}
	}

	switch s.Mode {
	case os.O_RDONLY:
		if !isValidFile {
			return s.setErrorClose(errors.New("file does not exist or invalid"))
		}

	case os.O_RDWR:
		if !isValidFile {
			clients := BlobDataServerManger.GetAllClients()
			if len(clients) < 3 {
				return s.setErrorClose(errors.New("no enough data server available"))
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
			var err1, err2 error
			go func() {
				defer wg.Done()
				err1 = os.Mkdir(filePath, 0775)
				s.SetBlobMetaData(v1.NewBlobMetaData(v1.BlobFileTypeName, filepath.Base(s.Path)))
				err2 = s.DumpBlobMetaData()
			}()

			wg.Wait()

			if utils.HasError(clientErrors) {
				for _, err := range clientErrors {
					_ = s.setErrorClose(err)
				}
				return s.setErrorClose(errors.New("create file failed, dataserver error"))
			} else if err1 != nil || err2 != nil {
				_ = s.setErrorClose(err1)
				_ = s.setErrorClose(err2)
				return s.setErrorClose(errors.New("create file failed, local error"))
			}
		}

	default:
		return s.setErrorClose(errors.New("invalid mode"))

	}
	s.Opened = true
	if isValidFile {
		err := s.pullChunk(0)
		if err != nil {
			s.Opened = false
			return s.setErrorClose(errors.New("failed to pull read buffer"))
		}
	}

	return nil
}

func (s *session) IsOpened() bool {
	return s.Opened
}

func (s *session) Seek(offset int64, whence int) (int64, error) {
	// Session-level lock
	// Chunk-level lock
	if !s.IsOpened() {
		return -1, errors.New("session closed")
	}

	var newOffset int64
	if whence == io.SeekStart {
		newOffset = offset
	} else if whence == io.SeekCurrent {
		newOffset += offset
	} else if whence == io.SeekEnd {
		newOffset = s.Blob.Size - 1 - offset
	}

	s.SessionMutex.Lock()
	if newOffset < 0 || newOffset > s.Blob.Size {
		s.SessionMutex.Unlock()
		return 0, errors.New("invalid offset")
	} else {

		s.Offset = newOffset

		s.ChunkMutex.Lock()
		go func() {
			defer s.SessionMutex.Unlock()
			defer s.ChunkMutex.Unlock()
			currChunkID := s.GetChunkID()
			err := s.pullChunk(currChunkID)
			if err != nil {
				_ = s.setErrorClose(err)
				return
			}
		}()

		return s.Offset, nil

	}
}

func (s *session) Close() error {
	// Session-level lock
	if !s.IsOpened() {
		return errors.New("session closed")
	}
	err := s.Flush()
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()
	if err != nil {
		_ = s.setErrorClose(err)
		return err
	} else {
		_ = s.setErrorClose(nil)

		s.ChunkMutex = nil
		s.Chunks = nil
		return nil
	}
}

func (s *session) Flush() error {
	// Session-level lock
	// Chunk-level lock
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	if !s.IsOpened() {
		return errors.New("session closed")
	}

	switch s.Mode {
	case os.O_RDONLY:
		return nil
	case os.O_RDWR:
		s.ChunkMutex.RLock()
		for chunkID := range s.Chunks {
			err := s.pushChunk(chunkID)
			if err != nil && err.Error() != "chunk already pushed" {
				return s.setErrorClose(err)
			}
		}
		s.ChunkMutex.RUnlock()

		s.freeBuffer()

		maxChunkID := utils.GetChunkID(s.Blob.Size)
		for idx, _ := range s.Blob.ChunkChecksums {
			if int64(idx) > maxChunkID {
				err := s.deleteChunk(int64(idx))
				if err != nil {
					return s.setErrorClose(err)
				}
			}
		}

		s.Blob.TruncateTo(maxChunkID)
		err := s.DumpBlobMetaData()
		if err != nil {
			return s.setErrorClose(err)
		} else {
			return nil
		}

	default:
		return s.setErrorClose(errors.New("invalid mode"))
	}
}

func (s *session) fetchReadBuffer(chunkID int64, chunkOffset int64, buffer []byte) (int64, error) {
	// Chunk-level lock
	if _, ok := s.Chunks[chunkID]; !ok {
		err := s.pullChunk(chunkID)
		if err != nil {
			return 0, err
		}
	}

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

func (s *session) Read(buffer []byte, size int64) (int64, error) {
	// Session-Level lock
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

func (s *session) setErrorClose(err error) error {
	if err != nil {
		log.Errorf("session %v closed due to error: %v", s.GetID(), err)
	} else {
		log.Errorf("session %v closed", s.GetID())
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

func (s *session) GetBlobMetaData() *v1.BlobMetaData {
	return &s.Blob
}

func (s *session) SetBlobMetaData(blob v1.BlobMetaData) {
	s.Blob = blob
}

func (s *session) DumpBlobMetaData() error {
	return s.Blob.Dump(s.GetMetaFilePath())
}

func (s *session) LoadBlobMetaData() error {
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
	}
}
