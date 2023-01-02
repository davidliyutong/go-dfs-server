package v1

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
)

const DefaultBlobChunkSize = 2 * 1024 * 1024 // 2MB
const BlobFileTypeName = "file"
const BlobDirTypeName = "dir"

type BlobMetaData struct {
	Type              string     `json:"type"`
	BaseName          string     `json:"base_name"`
	Size              int64      `json:"size"`
	Versions          []int64    `json:"versions"`
	ChunkChecksums    []string   `json:"chunk_checksums"`
	ChunkDistribution [][]string `json:"chunk_distribution"`
}

func (o *BlobMetaData) GetNumOfChunks() int64 {
	return int64(math.Ceil(float64(o.Size) / float64(DefaultBlobChunkSize)))
}

func (o *BlobMetaData) GetChunkDistribution(chunkID int64) ([]string, error) {
	if chunkID < 0 || chunkID >= o.GetNumOfChunks() {
		return nil, errors.New("wrong id")
	}
	if chunkID >= int64(len(o.ChunkDistribution)) {
		return nil, errors.New("blob corrupted")
	}
	return o.ChunkDistribution[chunkID], nil
}

func (o *BlobMetaData) ExtendTo(chunkID int64) {
	for {
		if chunkID >= int64(len(o.ChunkDistribution)) {
			o.ChunkDistribution = append(o.ChunkDistribution, make([]string, 0))
		} else if chunkID >= int64(len(o.ChunkChecksums)) {
			o.ChunkChecksums = append(o.ChunkChecksums, "")
		} else if chunkID >= int64(len(o.Versions)) {
			o.Versions = append(o.Versions, 0)
		} else {
			break
		}
	}
}

func (o *BlobMetaData) TruncateTo(chunkID int64) {
	o.ChunkDistribution = o.ChunkDistribution[:chunkID]
	o.ChunkChecksums = o.ChunkChecksums[:chunkID]
	o.Versions = o.Versions[:chunkID]
}

func (o *BlobMetaData) Dump(path string) error {
	filePtr, err := os.Create(path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(o)
	if err != nil {
		return filePtr.Close()
	}
	err = filePtr.Close()
	if err == nil {
		err = os.Chmod(path, 0775)
		return err
	}
	return err
}

func (o *BlobMetaData) Load(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return errors.New("cannot open metadata")
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	buffer, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(buffer, o)
	return err
}

func NewBlobMetaData(blobType string, baseName string) BlobMetaData {
	return BlobMetaData{
		Type:              blobType,
		BaseName:          baseName,
		Size:              0,
		Versions:          make([]int64, 0),
		ChunkChecksums:    make([]string, 0),
		ChunkDistribution: make([][]string, 0),
	}
}
