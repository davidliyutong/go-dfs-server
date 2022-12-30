package v1

func (b blobService) Lock(sessionID string) error {
	err := b.repo.BlobRepo().Lock(sessionID)
	return err
}
