package v1

import (
	"errors"
	v1 "go-dfs-server/pkg/dataserver/client/v1"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"os"
	"sync"
)

func (b blobService) Mkdir(path string) error {
	directoryPath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(directoryPath)
	if err == nil {
		return errors.New("file or directory exists")
	} else if os.IsNotExist(err) {
		err := os.Mkdir(directoryPath, 0775)
		if err != nil {
			return err
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
					errs[idx] = client.(v1.DataServerClient).BlobCreateDirectory(path)
				}()
			}
			wg.Wait()
			if utils.HasError(errs) {
				return errors.New("failed to create directory at some data servers")
			} else {
				return nil
			}
		}
	} else {
		return errors.New("cannot create direcotry")
	}

}
