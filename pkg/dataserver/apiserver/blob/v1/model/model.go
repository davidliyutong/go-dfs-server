package v1

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type BlobMetaData struct {
	Path    string
	Content map[int64]string
}

func NewBlobMetaData(path string) BlobMetaData {
	return BlobMetaData{
		Path:    path,
		Content: make(map[int64]string),
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
	err = json.Unmarshal(buffer, &o.Content)
	if o.Content == nil {
		o.Content = make(map[int64]string)
	}
	return err
}

func (o *BlobMetaData) Dump() error {
	filePtr, err := os.Create(o.Path)
	if err != nil {
		return err
	}
	defer func(filePtr *os.File) {
		err := filePtr.Close()
		if err != nil {

		}
	}(filePtr)
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(o.Content)
	return err
}
