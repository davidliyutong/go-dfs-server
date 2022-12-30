package v1

import (
	"errors"
	v1 "go-dfs-server/pkg/nameserver/server"
)

func (b blobService) Open(path string, mode int) (v1.Session, error) {

	sessionID, err := b.repo.BlobRepo().Open(path, mode)
	if err != nil {
		return nil, err
	} else {
		session, err := b.repo.BlobRepo().SessionManager().Get(sessionID)
		if err != nil {
			return nil, err
		} else {
			if session.IsOpened() {
				return session, nil
			} else {
				return nil, errors.New("failed to open file")
			}

		}
	}

}
