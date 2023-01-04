package memory

import (
	"bytes"
	v12 "go-dfs-server/pkg/dataserver/client/v1"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	repo2 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/status"
	"go-dfs-server/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
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
			session.ExtendToID(utils.GetChunkID(session.GetBlobMetaData().Size))

			return session.GetBlobMetaData(), err
		}
		return v1.BlobMetaData{}, status.ErrFileInValid
	} else {

		switch mode {
		case os.O_RDONLY:
			return v1.BlobMetaData{}, status.ErrFileNotExist
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
				return v1.BlobMetaData{}, session.SetErrorClose(status.ErrDataServerInsufficient)
			}
			clientErrors := make([]error, len(clients))
			createStatus := make([]string, len(clients))

			wg := new(sync.WaitGroup)
			wg.Add(len(clients) + 1)
			var numSuccessCreation int64 = 0
			for idx, client := range clients {
				idx := idx
				client := client
				go func() {
					defer wg.Done()
					err := client.(v12.DataServerClient).BlobCreateFile(*session.Path())
					if err != nil {
						clientErrors[idx] = err
						createStatus[idx] = ""
					} else {
						clientErrors[idx] = nil
						createStatus[idx] = client.(v12.DataServerClient).GetUUID()
						atomic.AddInt64(&numSuccessCreation, 1)
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
				if numSuccessCreation >= server.NameServerNumOfReplicas {
					_ = session.SetError(status.ErrDataServerOfflineSome)
				} else {
					for _, err := range clientErrors {
						_ = session.SetErrorClose(err)
					}
					return v1.BlobMetaData{}, session.SetErrorClose(status.ErrCreateErrorRemote)
				}
			} else if err1 != nil || err2 != nil {
				_ = session.SetErrorClose(err1)
				_ = session.SetErrorClose(err2)
				return v1.BlobMetaData{}, session.SetErrorClose(status.ErrCreateErrorLocal)
			}

			filePresence := utils.FilterEmptyString(createStatus)
			_ = session.SetFilePresence(filePresence)

			presentClients, err := r.DataServerManager().GetClients(filePresence)
			if err != nil {
				return v1.BlobMetaData{}, session.SetErrorClose(err)
			}

			clients = utils.SelectRandomNFromArray(presentClients, server.NameServerNumOfReplicas)
			clientErrors = make([]error, len(clients))
			newDistribution := make([]string, len(clients))
			newChecksums := make([]string, len(clients))
			var numFailed int64 = 0
			wg.Add(len(clients))
			for idx, client := range clients {
				idx := idx
				client := client
				go func() {
					defer wg.Done()
					err = client.(v12.DataServerClient).BlobCreateChunk(*session.Path(), int64(0))
					newDistribution[idx] = client.(v12.DataServerClient).GetUUID()
					newChecksums[idx], _ = utils.GetBufferMD5(nil)
					clientErrors[idx] = err
					if err != nil {
						atomic.AddInt64(&numFailed, 1)
					}
				}()
			}
			wg.Wait()

			if utils.HasError(clientErrors) {
				if numFailed <= 1 {
					_ = session.SetError(utils.GetFirstError(clientErrors))
				} else {
					_ = session.SetErrorClose(utils.GetFirstError(clientErrors))
				}
			}
			session.ExtendToID(0)
			newMeta := session.GetBlobMetaData()
			newMeta.ChunkChecksums[0] = newChecksums
			newMeta.ChunkDistribution[0] = newDistribution
			session.SetBlobMetaData(newMeta)

			err = session.DumpBlobMetaData()
			if err != nil {
				return v1.BlobMetaData{}, session.SetErrorClose(err)
			}
			return session.GetBlobMetaData(), nil
		default:
			return v1.BlobMetaData{}, status.ErrIOModeInvalid
		}
	}
}

func (r *blobRepo) Sync(path string, src v1.BlobMetaData) (v1.BlobMetaData, error) {
	if !utils.IsValidBlobMetaData(src) {
		return v1.BlobMetaData{}, status.ErrBlobCorrupted
	}

	var session server.Session

	session, err := r.sessions.Get(path)
	if err != nil {
		return src, status.ErrSessionNotFound
	} else {
		session.Add(1)
		defer session.Done()

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

			filePresence, _ := session.GetFilePresence()
			presentClients, err := r.DataServerManager().GetClients(filePresence)
			if err != nil {
				return dst, session.SetErrorClose(err)
			}

			for newChunkID := dstNChunks; newChunkID < srcNChunks; newChunkID++ {
				clients := utils.SelectRandomNFromArray(presentClients, server.NameServerNumOfReplicas)
				newDistribution := make([]string, len(clients))
				newChecksums := make([]string, len(clients))
				var numFailed int64 = 0
				wg.Add(len(clients))
				for idx, client := range clients {
					idx := idx
					client := client
					go func() {
						defer wg.Done()
						err = client.(v12.DataServerClient).BlobCreateChunk(*session.Path(), int64(newChunkID))
						newDistribution[idx] = client.(v12.DataServerClient).GetUUID()
						newChecksums[idx], _ = utils.GetBufferMD5(nil)
						clientErrors[idx] = err
						if err != nil {
							atomic.AddInt64(&numFailed, 1)
						}
					}()
				}
				wg.Wait()

				if utils.HasError(clientErrors) {
					if numFailed <= 1 {
						_ = session.SetError(utils.GetFirstError(clientErrors))
					} else {
						_ = session.SetErrorClose(utils.GetFirstError(clientErrors))
					}
				}
				src.ChunkDistribution[newChunkID] = newDistribution
				src.ChunkChecksums[newChunkID] = newChecksums
			}
			session.ExtendToID(int64(srcNChunks))

		case len(dst.ChunkChecksums) > len(src.ChunkChecksums):
			// shrink the file
			dstNChunks := len(dst.ChunkChecksums)
			srcNChunks := len(src.ChunkChecksums)
			wg := new(sync.WaitGroup)

			for oldChunkID := srcNChunks + 1; oldChunkID <= dstNChunks; oldChunkID++ {
				clientUUIDs, err := session.GetChunkDistribution(int64(oldChunkID))
				if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
					_ = session.SetError(err)
					continue
				}

				clients, err := r.DataServerManager().GetClients(clientUUIDs)
				if err != nil {
					_ = session.SetError(err)
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
			session.TruncateToID(int64(srcNChunks))

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
		return 0, status.ErrSessionNotFound
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
		return 0, session.SetErrorClose(status.ErrChunkCurrentNotPresent)
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

func (r *blobRepo) Write(path string, chunkID int64, chunkOffset int64, size int64, version int64, data io.ReadCloser) ([]string, int64, error) {
	session, err := r.sessions.Get(path)
	if err != nil {
		return nil, 0, status.ErrSessionNotFound
	}
	session.Add(1)
	defer session.Done()
	session.SyncMutex().RLock()
	defer session.SyncMutex().RUnlock()

	err = session.LockChunk(chunkID)
	if err != nil {
		return nil, 0, err
	}
	defer func(session server.Session, chunkID int64) {
		err := session.UnlockChunk(chunkID)
		if err != nil {

		}
	}(session, chunkID)

	clientUUIDs, err := session.GetChunkDistribution(chunkID)
	if err != nil || len(clientUUIDs) == 0 || clientUUIDs == nil {
		return nil, 0, session.SetErrorClose(status.ErrChunkCurrentNotPresent)
	}

	clients, err := r.DataServerManager().GetClients(clientUUIDs)
	if err != nil {
		return nil, 0, session.SetErrorClose(err)
	}

	MD5Strings := make([]string, len(clients))
	numWritten := make([]int64, len(clients))
	buf, err := io.ReadAll(data)
	if err != nil {
		return nil, 0, err
	}

	var numFailed int64 = 0
	wg := new(sync.WaitGroup)
	wg.Add(len(clients))
	for idx, client := range clients {
		idx := idx
		client := client
		go func() {
			defer wg.Done()
			MD5Strings[idx], numWritten[idx], err = client.(v12.DataServerClient).BlobWriteChunk(*session.Path(), chunkID, chunkOffset, size, version, bytes.NewBuffer(buf))
			if err != nil {
				atomic.AddInt64(&numFailed, 1)
			}
		}()
	}
	wg.Wait()

	if numFailed == 0 && utils.IsSameString(MD5Strings) && utils.IsSameInt64(numWritten) {
		return MD5Strings, numWritten[0], nil
	} else {
		if numFailed <= 1 {
			_ = session.SetError(status.ErrDataServerChecksumMismatch)
			for idx, v := range MD5Strings {
				if v != "" {
					return MD5Strings, numWritten[idx], nil
				}
			}
			return MD5Strings, numWritten[0], nil

		} else {
			return nil, 0, session.SetErrorClose(status.ErrDataServerChecksumMismatch)
		}
	}

}

func (r *blobRepo) Rmdir(path string) error {
	directoryPath, err := r.getAbsPath(path)
	if err != nil {
		return err
	}

	if utils.IsSameDirectory(server.GlobalServerDesc.Opt.Volume, directoryPath) {
		return status.ErrDirectoryCannotDelete
	}

	_, err = os.Stat(directoryPath)
	if err == nil {
		err = os.RemoveAll(directoryPath)
		if err != nil {
			return status.ErrDirectoryCannotRemoveLocal
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
				return status.ErrDirectoryCannotRemoveRemote
			}

		}
	} else if os.IsNotExist(err) {
		return status.ErrFileOrDirectoryNotExist
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
				return status.ErrDirectoryCannotRemoveLocal
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
					return status.ErrDataServerCannotRemoveSome
				} else {
					_ = r.SessionManager().Delete(path)
					return nil
				}
			}
		} else {
			return status.ErrFileNotFile
		}
	} else if os.IsNotExist(err) {
		return status.ErrFileOrDirectoryNotExist
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
		return status.ErrFileExists
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
				return status.ErrDataServerCannotCreateSome
			} else {
				return nil
			}
		}
	} else {
		return status.ErrDirectoryCannotCreate
	}
}

func (r *blobRepo) Ls(path string) (bool, []v1.BlobMetaData, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return false, nil, err
	}
	_, err = os.Stat(filePath)
	if err != nil {
		return false, nil, err
	}
	isFile := utils.GetFileState(filePath)
	if isFile {
		s := server.NewSession(path, filePath, "", os.O_RDONLY)
		_ = s.LoadBlobMetaData()
		res := make([]v1.BlobMetaData, 1)
		res[0] = v1.BlobMetaData{
			BaseName: filepath.Base(filePath),
			Size:     s.GetBlobMetaData().Size,
			Type:     v1.BlobFileTypeName,
		}
		return false, res, nil
	} else {
		lst, err := os.ReadDir(filePath)
		if err != nil {
			return true, nil, err
		}
		res := make([]v1.BlobMetaData, 0)
		for _, val := range lst {
			if !utils.GetFileState(filepath.Join(filePath, val.Name())) {
				res = append(res, v1.BlobMetaData{
					BaseName: val.Name(),
					Type:     v1.BlobDirTypeName,
				})
			} else {
				s := server.NewSession(path, filepath.Join(filePath, val.Name()), "", os.O_RDONLY)
				_ = s.LoadBlobMetaData()
				res = append(res, v1.BlobMetaData{
					BaseName: val.Name(),
					Size:     s.GetBlobMetaData().Size,
					Type:     v1.BlobFileTypeName,
				})
			}
		}
		return true, res, nil
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
