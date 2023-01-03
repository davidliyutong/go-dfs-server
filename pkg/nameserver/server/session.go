package server

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"path/filepath"
	"sync"
	"time"
)

type Session interface {
	Open() error
	IsOpened() bool
	Close() error
	SetErrorClose(err error) error

	GetTime() time.Time
	ID() *string
	Error() *[]error
	Path() *string
	FilePath() *string
	SyncMutex() *sync.RWMutex

	ExtendTo(chunkID int64)
	TruncateTo(chunkID int64)
	GetMetaFilePath() string
	GetChunkDistribution(chunkID int64) ([]string, error)

	LockChunk(chunkID int64) error
	RLockChunk(chunkID int64) error
	UnlockChunk(chunkID int64) error
	RUnlockChunk(chunkID int64) error

	GetBlobMetaData() v1.BlobMetaData
	SetBlobMetaData(blob v1.BlobMetaData)
	DumpBlobMetaData() error
	LoadBlobMetaData() error
	Delete()
	Wait()
	Add(n int)
	Done()
}

type session struct {
	path     string
	filePath string
	id       string
	time     time.Time
	opened   bool
	deleted  bool
	error    []error

	Blob v1.BlobMetaData

	sessionMutex *sync.RWMutex // protect from unwanted delete to the session | Open | Close | Delete
	syncMutex    *sync.RWMutex // avoid unwanted sync
	errorMutex   *sync.RWMutex // protect errors

	metaMutex  *sync.RWMutex  // protect blob and chunkMutex
	chunkMutex []sync.RWMutex // controls access to the chunks, protected by metaMutex

	eventGroup *sync.WaitGroup
}

func (s *session) Add(delta int) {
	s.time = time.Now()
	s.eventGroup.Add(delta)
}

func (s *session) Done() {
	s.eventGroup.Done()
}

func (s *session) Wait() {
	s.eventGroup.Wait()
}

func (s *session) SyncMutex() *sync.RWMutex {
	return s.syncMutex
}

func (s *session) LockChunk(chunkID int64) error {
	if chunkID < 0 || chunkID >= int64(len(s.chunkMutex)) {
		return errors.New("chunkID out of range")
	}
	s.chunkMutex[chunkID].Lock()
	return nil
}

func (s *session) RLockChunk(chunkID int64) error {
	if chunkID < 0 || chunkID >= int64(len(s.chunkMutex)) {
		return errors.New("chunkID out of range")
	}
	s.chunkMutex[chunkID].RLock()
	return nil
}

func (s *session) UnlockChunk(chunkID int64) error {
	if chunkID < 0 || chunkID >= int64(len(s.chunkMutex)) {
		return errors.New("chunkID out of range")
	}
	s.chunkMutex[chunkID].Unlock()
	return nil
}

func (s *session) RUnlockChunk(chunkID int64) error {
	if chunkID < 0 || chunkID >= int64(len(s.chunkMutex)) {
		return errors.New("chunkID out of range")
	}
	s.chunkMutex[chunkID].RUnlock()
	return nil
}

func (s *session) Delete() {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	s.deleted = true
}

func (s *session) Open() error {
	if s.deleted {
		return errors.New("session is deleted")
	}
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	s.opened = true
	s.time = time.Now()
	return nil
}

func (s *session) IsOpened() bool {
	return s.opened
}

func (s *session) Close() error {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	s.opened = false
	s.time = time.Now()
	return nil
}

func (s *session) SetErrorClose(err error) error {
	if s.deleted {
		return errors.New("session is deleted")
	}
	if err != nil {
		log.Errorf("session %v closed due to error: %v", s.Path(), err)
	} else {
		log.Debugf("session %v closed", s.Path())
	}
	s.errorMutex.Lock()
	defer s.errorMutex.Unlock()
	s.error = append(s.error, err)
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	s.deleted = true // mark as deleted to prevent further access
	return err
}

func (s *session) GetTime() time.Time {
	return s.time
}

func (s *session) ID() *string {
	return &s.id
}

func (s *session) Error() *[]error {
	return &s.error
}

func (s *session) Path() *string {
	return &s.path
}

func (s *session) FilePath() *string {
	return &s.filePath
}

func (s *session) GetMetaFilePath() string {
	return filepath.Join(s.filePath, "meta.json")
}

func (s *session) GetNumOfChunks() int64 {
	s.metaMutex.RLock()
	defer s.metaMutex.RUnlock()
	return s.Blob.GetNumOfChunks()
}

func (s *session) GetChunkDistribution(chunkID int64) ([]string, error) {
	s.metaMutex.RLock()
	defer s.metaMutex.RUnlock()
	return s.Blob.GetChunkDistribution(chunkID)
}

func (s *session) ExtendTo(chunkID int64) {
	s.metaMutex.Lock()
	defer s.metaMutex.Unlock()
	s.Blob.ExtendTo(chunkID)
	if chunkID >= int64(len(s.chunkMutex)) {
		s.chunkMutex = append(s.chunkMutex, make([]sync.RWMutex, 1+chunkID-int64(len(s.chunkMutex)))...)
	}
}

func (s *session) TruncateTo(chunkID int64) {
	s.metaMutex.Lock()
	defer s.metaMutex.Unlock()
	s.Blob.TruncateTo(chunkID)
	s.chunkMutex = s.chunkMutex[:chunkID]
}

func (s *session) GetBlobMetaData() v1.BlobMetaData {
	s.metaMutex.RLock()
	defer s.metaMutex.RUnlock()
	return s.Blob
}

func (s *session) SetBlobMetaData(blob v1.BlobMetaData) {
	s.metaMutex.Lock()
	defer s.metaMutex.Unlock()
	s.Blob = blob
}

func (s *session) DumpBlobMetaData() error {
	s.metaMutex.Lock()
	defer s.metaMutex.Unlock()
	return s.Blob.Dump(s.GetMetaFilePath())
}

func (s *session) LoadBlobMetaData() error {
	s.metaMutex.RLock()
	defer s.metaMutex.RUnlock()
	return s.Blob.Load(s.GetMetaFilePath())
}

func NewSession(path string, filePath string, id string, mode int) Session {
	return &session{
		path:     path,
		filePath: filePath,
		id:       id,
		time:     time.Now(),
		opened:   false,
		deleted:  false,
		error:    make([]error, 0),

		Blob: v1.BlobMetaData{},

		sessionMutex: new(sync.RWMutex),
		errorMutex:   new(sync.RWMutex),
		syncMutex:    new(sync.RWMutex),

		chunkMutex: make([]sync.RWMutex, 0),
		metaMutex:  new(sync.RWMutex),

		eventGroup: new(sync.WaitGroup),
	}
}
