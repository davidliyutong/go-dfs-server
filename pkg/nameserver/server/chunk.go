package server

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"
	"sync"
	"time"
)

type ChunkBuffer struct {
	path    string
	chunkID int64
	version int64
	pushed  bool
	buffer  *ByteBuffer
	rwmutex *sync.RWMutex
	trmutex *sync.RWMutex
}

func (mf *ChunkBuffer) IsPushed() bool {
	mf.trmutex.RLock()
	defer mf.trmutex.RUnlock()
	return mf.pushed
}

func (mf *ChunkBuffer) SetPushed(flag bool) {
	mf.trmutex.Lock()
	defer mf.trmutex.Unlock()
	mf.pushed = flag
}

func (mf *ChunkBuffer) Version() int64 {
	mf.rwmutex.RLock()
	defer mf.rwmutex.RUnlock()
	return mf.version
}

func (mf *ChunkBuffer) SetVersion(version int64) {
	mf.rwmutex.Lock()
	defer mf.rwmutex.Unlock()
	mf.version = version
}

func (mf *ChunkBuffer) Bytes() []byte {
	mf.rwmutex.RLock()
	defer mf.rwmutex.RUnlock()
	return mf.buffer.buffer
}

func NewChunkBuffer(path string, chunkID int64, version int64, buffer []byte) *ChunkBuffer {
	return &ChunkBuffer{
		path:    path,
		chunkID: chunkID,
		version: version,
		pushed:  false,
		buffer:  MakeByteBuffer(buffer),
		rwmutex: new(sync.RWMutex),
		trmutex: new(sync.RWMutex),
	}
}

func (mf *ChunkBuffer) Stat() (fs.FileInfo, error) {
	return mf, nil
}

func (mf *ChunkBuffer) Read(buffer []byte) (int, error) {
	mf.rwmutex.RLock()
	defer mf.rwmutex.RUnlock()
	return mf.buffer.Read(buffer)
}

func (mf *ChunkBuffer) Close() error {
	return nil
}

func (mf *ChunkBuffer) Write(buffer []byte) (n int, err error) {
	mf.rwmutex.Lock()
	mf.trmutex.Lock()
	defer mf.rwmutex.Unlock()
	defer mf.trmutex.Unlock()
	mf.pushed = false
	mf.version += 1
	return mf.buffer.Write(buffer)
}

func (mf *ChunkBuffer) WriteInPlace(buffer []byte) (n int, err error) {
	mf.rwmutex.Lock()
	mf.trmutex.Lock()
	defer mf.rwmutex.Unlock()
	defer mf.trmutex.Unlock()
	mf.pushed = false
	mf.version += 1
	if mf.buffer.index+len(buffer) <= len(mf.buffer.buffer) {
		n = copy(mf.buffer.buffer[mf.buffer.index:mf.buffer.index+len(buffer)], buffer)
		mf.buffer.index += n
		if len(buffer) != n {
			return len(mf.buffer.buffer), errors.New("failed to write in place")
		} else {
			return len(mf.buffer.buffer), nil

		}
	} else {
		mf.buffer.index, err = mf.buffer.Write(buffer)
		return mf.buffer.index, err
	}
}

func (mf *ChunkBuffer) Seek(offset int64, whence int) (int64, error) {
	mf.rwmutex.Lock()
	defer mf.rwmutex.Unlock()
	return mf.buffer.Seek(offset, whence)
}

func (mf *ChunkBuffer) Position() int {
	mf.rwmutex.Lock()
	defer mf.rwmutex.Unlock()
	return mf.buffer.Position()
}

func (mf *ChunkBuffer) Path() string       { return mf.path }
func (mf *ChunkBuffer) Name() string       { return mf.path + "." + strconv.Itoa(int(mf.chunkID)) }
func (mf *ChunkBuffer) Size() int64        { return int64(mf.buffer.Len()) }
func (mf *ChunkBuffer) Mode() os.FileMode  { return 0666 }
func (mf *ChunkBuffer) ModTime() time.Time { return time.Time{} }
func (mf *ChunkBuffer) IsDir() bool        { return false }
func (mf *ChunkBuffer) Sys() interface{}   { return nil }

type ByteBuffer struct {
	buffer []byte
	index  int
}

func MakeByteBuffer(buffer []byte) *ByteBuffer {
	return &ByteBuffer{
		buffer: buffer,
		index:  0,
	}
}

func (bb *ByteBuffer) Reset() {
	bb.index = 0
}

func (bb *ByteBuffer) Len() int {
	return len(bb.buffer)
}

func (bb *ByteBuffer) Position() int {
	return bb.index
}

func (bb *ByteBuffer) Bytes() []byte {
	return bb.buffer
}

func (bb *ByteBuffer) Read(buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, nil
	}

	if bb.index >= bb.Len() {
		return 0, io.EOF
	}

	// copy 会判断 buffer 的大小
	last := copy(buffer, bb.buffer[bb.index:])
	bb.index += last
	return last, nil
}

func (bb *ByteBuffer) Write(buffer []byte) (int, error) {
	bb.buffer = append(bb.buffer[:bb.index], buffer...)
	return len(buffer), nil
}

func (bb *ByteBuffer) Seek(offset int64, whence int) (int64, error) {
	var newIndex int
	switch whence {
	default:
	case io.SeekStart:
		newIndex = int(offset)
	case io.SeekCurrent:
		newIndex += int(offset)
	case io.SeekEnd:
		newIndex = bb.Len() - 1 - int(offset)
	}

	if newIndex < 0 || newIndex > bb.Len() {
		return int64(bb.index), errors.New("invalid offset")
	} else {
		bb.index = newIndex
		return int64(bb.index), nil
	}
}
