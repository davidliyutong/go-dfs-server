package v1

import (
	"errors"
	"math"
)

const DefaultBlobChunkSize = 2 * 1024 * 1024 // 2MB

type BlobStruct struct {
	Size              int64
	ChunkChecksums    []string
	ChunkDistribution [][]string
}

func (o *BlobStruct) GetNumOfChunks() int64 {
	return int64(math.Ceil(float64(o.Size) / float64(DefaultBlobChunkSize)))
}

func (o *BlobStruct) GetChunkDistribution(chunkID int64) ([]string, error) {
	if chunkID < 0 || chunkID >= o.GetNumOfChunks() {
		return nil, errors.New("wrong id")
	}
	if chunkID >= int64(len(o.ChunkDistribution)) {
		return nil, errors.New("blob corrupted")
	}
	return o.ChunkDistribution[chunkID], nil
}
