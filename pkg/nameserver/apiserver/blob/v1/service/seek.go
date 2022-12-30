package v1

func (b blobService) Seek(sessionID string, offset int64, whence int) (int64, error) {
	session, err := b.repo.BlobRepo().SessionManager().Get(sessionID)
	if err != nil {
		return -1, err
	}
	seek, err := session.Seek(offset, whence)
	if err != nil {
		return *session.GetOffset(), err
	} else {
		return seek, nil
	}
}
