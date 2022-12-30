package v1

func (b blobService) Close(sessionID string) error {
	return b.repo.BlobRepo().Close(sessionID)
}
