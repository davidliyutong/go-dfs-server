package v1

func (b blobService) Flush(sessionID string) error {
	return b.repo.BlobRepo().Flush(sessionID)
}
