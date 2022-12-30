package v1

import (
	"errors"
	v1 "go-dfs-server/pkg/dataserver/client/v1"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"os"
	"path/filepath"
	"sync"
)

func (b blobService) Rmdir(path string) error {
	directoryPath := filepath.Join(server.GlobalServerDesc.Opt.Volume, path)
	_, err := os.Stat(directoryPath)
	if err == nil {
		err = os.RemoveAll(directoryPath)
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
					errs[idx] = client.(v1.DataServerClient).BlobDeleteDirectory(path)
				}()
			}
			wg.Wait()
			if utils.HasError(errs) {
				return errors.New("failed to remove directory from some data servers")
			} else {
				return nil
			}

		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	}
	return err
}
