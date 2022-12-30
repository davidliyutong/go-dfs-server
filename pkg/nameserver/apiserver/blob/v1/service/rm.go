package v1

import (
	"errors"
	v1 "go-dfs-server/pkg/dataserver/client/v1"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"os"
	"sync"
)

func (b blobService) Rm(path string) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			err = os.RemoveAll(filePath)
			if err != nil {
				return errors.New("failed to remove directory")
			} else {
				clients := b.repo.BlobRepo().DataServerManager().GetAllClients()
				wg := new(sync.WaitGroup)
				wg.Add(len(clients))
				errs := make([]error, len(clients))
				for idx, client := range clients {
					client := client
					idx := idx
					go func() {
						wg.Done()
						errs[idx] = client.(v1.DataServerClient).BlobDeleteFile(path)
					}()
				}
				wg.Wait()

				if utils.HasError(errs) {
					return errors.New("failed to remove file from some data servers")
				} else {
					return nil
				}
			}
		} else {
			return errors.New("not a file")
		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	} else {
		return err
	}
}
