package v1

func (b blobService) Rm(path string, recursive bool) error {
	err := b.repo.BlobRepo().Rm(path, recursive)
	if err != nil {
		// TODO: handle error
	}
	return err
}
