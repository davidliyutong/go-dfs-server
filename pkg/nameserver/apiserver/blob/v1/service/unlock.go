package v1

func (b blobService) Unlock(sessionID string) error {
	err := b.repo.BlobRepo().Unlock(sessionID)
	return err
}
