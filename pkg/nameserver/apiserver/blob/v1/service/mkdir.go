package v1

func (b blobService) Mkdir(path string) error {
	err := b.repo.BlobRepo().Mkdir(path)
	if err != nil {
		// TODO: handle error
	}
	return err
}
