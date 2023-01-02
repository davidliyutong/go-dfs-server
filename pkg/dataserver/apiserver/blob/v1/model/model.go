package v1

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type BlobMetaData struct {
	Path           string           `json:"path"`
	Versions       map[int64]int64  `json:"versions"`
	ChunkChecksums map[int64]string `json:"chunk_checksums"`
}

func NewBlobMetaData(path string) BlobMetaData {
	return BlobMetaData{
		Path:           path,
		Versions:       make(map[int64]int64),
		ChunkChecksums: make(map[int64]string),
	}
}

func (o *BlobMetaData) Load() error {
	jsonFile, err := os.Open(o.Path)
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
	if o.ChunkChecksums == nil {
		o.ChunkChecksums = make(map[int64]string)
	}
	return err
}

func (o *BlobMetaData) Dump() error {
	filePtr, err := os.Create(o.Path)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(o)
	if err != nil {
		return filePtr.Close()
	}
	err = filePtr.Close()
	return err
}
