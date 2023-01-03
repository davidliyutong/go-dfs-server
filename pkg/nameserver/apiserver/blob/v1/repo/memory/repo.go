package memory

import (
	"bytes"
	"errors"
	v12 "go-dfs-server/pkg/dataserver/client/v1"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	repo2 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type blobRepo struct {
	volume   string
	servers  server.DataServerManager
	sessions server.SessionManager
}

var _ repo2.BlobRepo = &blobRepo{"", nil, nil}

func newBlobRepo(volume string, servers server.DataServerManager, sessions server.SessionManager) repo2.BlobRepo {
	return &blobRepo{
		volume:   volume,
		servers:  servers,
		sessions: sessions,
	}
}

func (r *blobRepo) getAbsPath(path string) (string, error) {
	return utils.JoinSubPathSafe(r.volume, path)
}

func (r *blobRepo) Open(path string, mode int) (v1.BlobMetaData, error) {
	filePath, err := r.getAbsPath(path)
	if err != nil {
		return v1.BlobMetaData{}, err
	}
	if utils.PathExists(filePath) {
		if utils.GetFileState(filePath) {
			var session server.Session
			session, err := r.sessions.Get(path)
			if err != nil {
				session, err = r.sessions.New(path, filePath, mode)
				if err != nil {
					return v1.BlobMetaData{}, err
				}
				err = session.LoadBlobMetaData()
				if err != nil {
					return v1.BlobMetaData{}, err
				}
			}
			err = session.Open()
			defer func(session server.Session) {
				err := session.Close()
				if err != nil {

				}
			}(session)
			if err != nil {
				return v1.BlobMetaData{}, err
			}

			session.SyncMutex().RLock()
			defer session.SyncMutex().RUnlock()
			session.ExtendTo(utils.GetChunkID(session.GetBlobMetaData().Size))

			return session.GetBlobMetaData(), err
		}
		return v1.BlobMetaData{}, errors.New("not a valid file")
	} else {

		switch mode {
		case os.O_RDONLY:
			return v1.BlobMetaData{}, errors.New("file does not exist")
		case os.O_WRONLY:
			fallthrough
		case os.O_RDWR:
			session, err := r.sessions.Get(path)
			if err == nil {
				// clean orphans
				err := r.SessionManager().Delete(path)
				if err != nil {
					return v1.BlobMetaData{}, err
				}
			}

			session, err = r.sessions.New(path, filePath, mode)
			if err != nil {
				return v1.BlobMetaData{}, err
			}
			err = session.Open()
			defer func(session server.Session) {
				err := session.Close()
				if err != nil {

				}
			}(session)
			if err != nil {
				return v1.BlobMetaData{}, err
			}
			// layout file on data servers
			clients := r.DataServerManager().GetAllClients()
			if len(clients) < 3 {
				return v1.BlobMetaData{}, session.SetErrorClose(errors.New("no enough data server available"))
			}
			clientErrors := make([]error, len(clients))

			wg := sync.WaitGroup{}
			wg.Add(len(clients) + 1)
			for idx, client := range clients {
				idx := idx
				client := client
				go func() {
					defer wg.Done()
					err := client.(v12.DataServerClient).BlobCreateFile(*session.Path())
					if err != nil {
						clientErrors[idx] = err
					} else {
						clientErrors[idx] = nil
					}

				}()
			}

			// create file meta data on name server
			var err1, err2 error
			go func() {
				defer wg.Done()
				err1 = os.Mkdir(*session.FilePath(), 0775)
				session.SetBlobMetaData(v1.NewBlobMetaData(v1.BlobFileTypeName, filepath.Base(*session.Path())))
				err2 = session.DumpBlobMetaData()
			}()
			wg.Wait()

			if utils.HasError(clientErrors) {
				for _, err := range clientErrors {
					_ = session.SetErrorClose(err)
				}
				return v1.BlobMetaData{}, session.SetErrorClose(errors.New("create file failed, dataserver error"))
			} else if err1 != nil || err2 != nil {
				_ = session.SetErrorClose(err1)
				_ = session.SetErrorClose(err2)
				return v1.BlobMetaData{}, session.SetErrorClose(errors.New("create file failed, local error"))
			}

			clientErrors = make([]error, server.NameServerNumOfReplicas)
			clients = utils.SelectRandomNFromArray(r.DataServerManager().GetAllClients(), server.NameServerNumOfReplicas)
			newDistribution := make([]string, len(clients))
			wg.Add(len(clients))
			for idx, client := range clients {
				idx := idx
				client := client
				go func() {
					defer wg.Done()
					err = client.(v12.DataServerClient).BlobCreateChunk(*session.Path(), int64(0))
					newDistribution[idx] = client.(v12.DataServerClient).GetUUID()
					clientErrors[idx] = err
				}()
			}
			wg.Wait()

			if utils.HasError(clientErrors) {
				_ = session.SetErrorClose(utils.GetFirstError(clientErrors))
			}
			session.ExtendTo(0)
			newMeta := session.GetBlobMetaData()
			checksum, _ := utils.GetBufferMD5(nil)
			newMeta.ChunkChecksums[0] = checksum
			newMeta.ChunkDistribution[0] = newDistribution
			session.SetBlobMetaData(newMeta)
			return session.GetBlobMetaData(), nil
		default:
			return v1.BlobMetaData{}, errors.New("invalid mode")
		}
	}
}

func (r *blobRepo) Sync(path string, src v1.BlobMetaData) (v1.BlobMetaData, error) {
	if !utils.IsValidBlobMetaData(src) {
		return v1.BlobMetaData{}, errors.New("invalid blob meta data")
	}

	var session server.Session

	session, err := r.sessions.Get(path)
	if err != nil {
		return src, errors.New("session not found")
	} else {

		err = session.Open()
		defer func(session server.Session) {
			err := session.Close()
			if err != nil {

			}
		}(session)
		if err != nil {
			return src, err
		}

		session.SyncMutex().Lock()
		defer session.SyncMutex().Unlock()

		dst := session.GetBlobMetaData()
		if utils.IsNewerThan(dst, src) {
			return dst, err
		}

		switch {
		case len(dst.ChunkChecksums) < len(src.ChunkChecksums):
			// expand the file
			dstNChunks := len(dst.ChunkChecksums)
			srcNChunks := len(src.ChunkChecksums)
			clientErrors := make([]error, server.NameServerNumOfReplicas)
			wg := new(sync.WaitGroup)
			for newChunkID := dstNChunks; newChunkID < srcNChunks; newChunkID++ {
				clients := utils.SelectRandomNFromArray(r.DataServerManager().GetAllClients(), server.NameServerNumOfReplicas)
				newDistribution := make([]string, len(clients))
				wg.Add(len(clients))
				for idx, client := range clients {
					idx := idx
					client := client
					go func() {
						defer wg.Done()
						err = client.(v12.DataServerClient).BlobCreateChunk(*session.Path(), int64(newChunkID))
						newDistribution[idx] = client.(v12.DataServerClient).GetUUID()
						clientErrors[idx] = err
					}()
				}
				wg.Wait()

				if utils.HasError(clientErrors) {
					_ = session.SetErrorClose(utils.GetFirstError(clientErrors))
				}
				src.ChunkDistribution[newChunkID] = newDistribution
			}
			session.ExtendTo(int64(srcNChunks))

		case len(dst.ChunkChecksums) > len(src.ChunkChecksums):
			// shrink the file
			dstNChunks := len(dst.ChunkChecksums)
			srcNChunks := len(src.ChunkChecksums)
			wg := new(sync.WaitGroup)

			for oldChunkID := srcNChunks + 1; oldChunkID <= dstNChunks; oldChunkID++ {
				clientUUIDs, err := session.GetChunkDistribution(int64(oldChunkID))
				if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
					//return v1.BlobMetaData{}, errors.New("current chunk is not present")
					continue
				}

				clients, err := r.DataServerManager().GetClients(clientUUIDs)
				if err != nil {
					//return v1.BlobMetaData{}, errors.New("cannot get related data server")
					_ = session.SetErrorClose(err)
					continue
				}

				clientErrors := make([]error, len(clients))

				wg.Add(len(clients))
				for idx, client := range clients {
					idx := idx
					client := client
					go func() {
						defer wg.Done()
						clientErrors[idx] = client.(v12.DataServerClient).BlobDeleteChunk(*session.Path(), int64(oldChunkID))
					}()
				}
				wg.Wait()

				if utils.HasError(clientErrors) {
					_ = session.SetErrorClose(utils.GetFirstError(clientErrors))
				}
			}
			session.TruncateTo(int64(srcNChunks))

		case len(dst.ChunkChecksums) == len(src.ChunkChecksums):
			// TODO: check if the file is modified src and dst should only have Version as difference
			src.Version -= 1
			break
		}
		session.SetBlobMetaData(src)
		err := session.DumpBlobMetaData()
		if err != nil {
			return src, err
		}
		return session.GetBlobMetaData(), nil
	}
}

func (r *blobRepo) Read(buffer io.Writer, path string, chunkID int64, chunkOffset int64, size int64) (int64, error) {
	session, err := r.sessions.Get(path)
	if err != nil {
		return 0, errors.New("session not found")
	}
	session.Add(1)
	defer session.Done()
	session.SyncMutex().RLock()
	defer session.SyncMutex().RUnlock()

	err = session.RLockChunk(chunkID)
	if err != nil {
		return 0, err
	}
	defer func(session server.Session, chunkID int64) {
		err := session.RUnlockChunk(chunkID)
		if err != nil {

		}
	}(session, chunkID)

	clientUUIDs, err := session.GetChunkDistribution(chunkID)
	if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
		return 0, session.SetErrorClose(errors.New("current chunk is not present"))
	}

	clients, err := r.DataServerManager().GetClients(clientUUIDs)
	if err != nil {
		return 0, session.SetErrorClose(err)
	}

	for _, client := range clients {
		reader, err := client.(v12.DataServerClient).BlobReadChunk(*session.Path(), chunkID, chunkOffset, size)
		if err != nil {
			continue
		}
		n, err := io.Copy(buffer, reader)
		err = reader.Close()
		if err != nil {
			return n, err
		}
		return n, nil
	}
	return 0, nil

}

func (r *blobRepo) Write(path string, chunkID int64, chunkOffset int64, size int64, version int64, data io.ReadCloser) (string, int64, error) {
	session, err := r.sessions.Get(path)
	if err != nil {
		return "", 0, errors.New("session not found")
	}
	session.Add(1)
	defer session.Done()
	session.SyncMutex().RLock()
	defer session.SyncMutex().RUnlock()

	err = session.LockChunk(chunkID)
	if err != nil {
		return "", 0, err
	}
	defer func(session server.Session, chunkID int64) {
		err := session.UnlockChunk(chunkID)
		if err != nil {

		}
	}(session, chunkID)

	clientUUIDs, err := session.GetChunkDistribution(chunkID)
	if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
		return "", 0, session.SetErrorClose(errors.New("current chunk is not present"))
	}

	clients, err := r.DataServerManager().GetClients(clientUUIDs)
	if err != nil {
		return "", 0, session.SetErrorClose(err)
	}

	MD5Strings := make([]string, len(clients))
	numWritten := make([]int64, len(clients))
	buf, err := io.ReadAll(data)
	if err != nil {
		return "", 0, err
	}

	for idx, client := range clients {

		MD5Strings[idx], numWritten[idx], err = client.(v12.DataServerClient).BlobWriteChunk(*session.Path(), chunkID, chunkOffset, size, version, bytes.NewBuffer(buf))
		if err != nil {
			return "", 0, session.SetErrorClose(err)
		}
	}
	if utils.IsSameString(MD5Strings) && utils.IsSameInt64(numWritten) {
		return MD5Strings[0], numWritten[0], nil
	} else {
		return "", 0, errors.New("not all servers have the same checksum")
	}

}

func (r *blobRepo) Rmdir(path string) error {
	directoryPath, err := r.getAbsPath(path)
	if err != nil {
		return err
	}

	if utils.IsSameDirectory(server.GlobalServerDesc.Opt.Volume, directoryPath) {
		return errors.New("cannot delete this directory")
	}

	_, err = os.Stat(directoryPath)
	if err == nil {
		err = os.RemoveAll(directoryPath)
		if err != nil {
			return errors.New("failed to remove directory")
		} else {
			clients := r.DataServerManager().GetAllClients()
			wg := new(sync.WaitGroup)
			wg.Add(len(clients))
			errs := make([]error, len(clients))
			for idx, client := range clients {
				client := client
				idx := idx
				go func() {
					wg.Done()
					errs[idx] = client.(v12.DataServerClient).BlobDeleteDirectory(path)
				}()
			}
			wg.Wait()
			if utils.HasError(errs) {
				return errors.New("failed to remove directory from some data servers")
			}

		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	}
	return err
}

func (r *blobRepo) Rm(path string, recursive bool) error {
	if recursive {
		return r.Rmdir(path)
	}
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
				clients := r.DataServerManager().GetAllClients()
				wg := new(sync.WaitGroup)
				wg.Add(len(clients))
				errs := make([]error, len(clients))
				for idx, client := range clients {
					client := client
					idx := idx
					go func() {
						wg.Done()
						errs[idx] = client.(v12.DataServerClient).BlobDeleteFile(path)
					}()
				}
				wg.Wait()

				if utils.HasError(errs) {
					return errors.New("failed to remove file from some data servers")
				} else {
					_ = r.SessionManager().Delete(path)
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

func (r *blobRepo) Mkdir(path string) error {
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
			clients := r.DataServerManager().GetAllClients()
			wg := new(sync.WaitGroup)
			wg.Add(len(clients))
			errs := make([]error, len(clients))
			for idx, client := range clients {
				client := client
				idx := idx
				go func() {
					wg.Done()
					errs[idx] = client.(v12.DataServerClient).BlobCreateDirectory(path)
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
		return errors.New("cannot create directory")
	}
}

func (r *blobRepo) Ls(path string) ([]v1.BlobMetaData, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	isFile := utils.GetFileState(filePath)
	if isFile {
		res := make([]v1.BlobMetaData, 1)
		res[0] = v1.BlobMetaData{
			BaseName: filepath.Base(filePath),
			Size:     server.NewSession(path, filePath, "", os.O_RDONLY).GetBlobMetaData().Size,
			Type:     v1.BlobFileTypeName,
		}
		return res, nil
	} else {
		lst, err := os.ReadDir(filePath)
		if err != nil {
			return nil, err
		}
		res := make([]v1.BlobMetaData, 0)
		for _, val := range lst {
			if !utils.GetFileState(filepath.Join(filePath, val.Name())) {
				res = append(res, v1.BlobMetaData{
					BaseName: val.Name(),
					Type:     v1.BlobDirTypeName,
				})
			} else {
				res = append(res, v1.BlobMetaData{
					BaseName: val.Name(),
					Size:     server.NewSession(path, filePath, "", os.O_RDONLY).GetBlobMetaData().Size,
					Type:     v1.BlobFileTypeName,
				})
			}
		}
		return res, nil
	}
}

func (r *blobRepo) SessionManager() server.SessionManager {
	return r.sessions
}

func (r *blobRepo) DataServerManager() server.DataServerManager {
	return r.servers
}

type repo struct {
	blobRepo repo2.BlobRepo
}

//var _ repo3.BlobRepo = (*repo)(nil)

var (
	r    repo
	once sync.Once
)

// Repo creates and returns the store client instance.
func Repo(volume string, servers server.DataServerManager, sessions server.SessionManager) (repo2.Repo, error) {
	once.Do(func() {
		r = repo{
			blobRepo: newBlobRepo(volume, servers, sessions),
		}
	})

	return r, nil
}

func (r repo) BlobRepo() repo2.BlobRepo {
	return r.blobRepo
}

// Close closes the repo.
func (r repo) Close() error {
	return r.Close()
}
