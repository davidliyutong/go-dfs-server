package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	v12 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"os"
	"path/filepath"
)

func upload(cli v1.NameServerClient, src string, dst string, force bool) error {
	log.Infof("%s -> %s, %v", src, dst, force)
	handle, err := cli.Open(dst, os.O_RDWR)
	if err != nil {
		log.Errorln(err)
		return err
	}
	input, err := os.OpenFile(src, os.O_RDONLY, 0664)
	if err != nil {
		log.Errorln(err)
		return err
	}

	var total = 0
	for {
		buf := make([]byte, v12.DefaultBlobChunkSize)
		read, err := input.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				log.Errorln(err)
			}
			break
		}
		written, err := handle.Write(buf[:read])
		if err != nil {
			log.Errorln(err)
			break
		}
		total += written
	}

	err = handle.Close()
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func uploadRecursive(cli v1.NameServerClient, src string, dst string, force bool) error {
	log.Infof("%s <-> %s, %v recursive", src, dst, force)
	lst, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, val := range lst {
		srcSubPath := filepath.Join(src, val.Name())
		dstSubPath := filepath.Join(dst, val.Name())

		if val.IsDir() {
			// Create remote directory
			isDir, _, err := cli.BlobLs(dstSubPath)
			if err != nil {
				err = cli.BlobMkdir(dstSubPath)
				if err != nil {
					return err
				}
				log.Info("creating: ", dstSubPath)
			} else {
				if !isDir {
					log.Warning("skipping existing file: ", dstSubPath)
				}
			}
			if err = uploadRecursive(cli, srcSubPath, dstSubPath, force); err != nil {
				return err
			}
		} else {
			isDir, _, err := cli.BlobLs(dstSubPath)
			if err == nil {
				if isDir {
					log.Warning("skipping existing directory: ", dstSubPath)
				} else {
					if !force {
						log.Warning("skipping existing file: ", dstSubPath)
					} else {
						log.Warning("replacing existing file: ", dstSubPath)
						err := cli.BlobRm(dstSubPath, false)
						if err != nil {
							return err
						}
						err = upload(cli, srcSubPath, dstSubPath, force)
					}
				}
			} else {
				err = upload(cli, srcSubPath, dstSubPath, force)
			}

		}
	}
	return nil
}

func execPut(cli v1.NameServerClient, srcPath string, dstPath string, srcIsDir bool, dstIsDir bool, dstExist bool, force bool) error {
	if srcIsDir {
		if !dstExist {
			err := cli.BlobMkdir(dstPath)
			if err != nil {
				return err
			}
		}
		return uploadRecursive(cli, srcPath, dstPath, force)
	} else {
		if dstIsDir {
			actualDstPath := filepath.Join(dstPath, filepath.Base(srcPath))
			return upload(cli, srcPath, actualDstPath, force)
		} else {
			return upload(cli, srcPath, dstPath, force)
		}
	}
}

func Put(cmd *cobra.Command, args []string) {
	opt := config.NewClientOpt()
	_, err := opt.Parse(cmd)
	if err != nil {
		log.Println("cannot find credential, run login first")
	} else {
		recursive, _ := cmd.Flags().GetBool("recursive")
		force, _ := cmd.Flags().GetBool("force")
		cli := v1.NewNameServerClient(opt.Token, opt.Hostname, opt.Port, opt.UseTLS)

		srcPath := args[0]
		info, err := os.Stat(srcPath)
		if err != nil {
			log.Error("source directory does not exist")
			return
		}
		if info.IsDir() && !recursive {
			log.Error("source is a directory, use -r to upload recursively")
			return
		}

		srcIsDir := info.IsDir()

		dstPath := args[1]
		dstParent, _ := filepath.Split(dstPath)

		var dstIsDir bool
		var dstExist bool
		isDir, _, err := cli.BlobLs(dstPath)

		if err == nil {
			dstExist = true
			dstIsDir = isDir
			if srcIsDir && !dstIsDir {
				log.Error("cannot upload directory to a file")
				return
			}
		} else {
			_, _, err := cli.BlobLs(dstParent)
			if err == nil {
				dstExist = false
				dstIsDir = srcIsDir
			} else {
				log.Error("destination directory does not exist")
			}
		}
		err = execPut(cli, srcPath, dstPath, srcIsDir, dstIsDir, dstExist, force)
		if err != nil {
			log.Errorln(err)
		}
		//log.Infof("srcIsDir: %v, dstIsDir: %v, dstExist: %v, recursive %v, force %v", srcIsDir, dstIsDir, dstExist, recursive, force)
	}
}
