package v1

func (b blobService) Rmdir(path string) error {
	err := b.repo.BlobRepo().Rm(path, true)
	if err != nil {
		// TODO: handle error
	}
	return err
}
