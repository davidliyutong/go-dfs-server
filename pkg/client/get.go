package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"io"
	"os"
	"path/filepath"
)

func download(cli v1.NameServerClient, src string, dst string, force bool) error {
	log.Infof("%s -> %s, %v", src, dst, force)
	handle, err := cli.Open(src, os.O_RDONLY)
	if err != nil {
		log.Errorln(err)
		return err
	}
	output, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Errorln(err)
		return err
	}
	_, _ = io.Copy(output, handle)

	err = handle.Close()
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func downloadRecursive(cli v1.NameServerClient, src string, dst string, force bool) error {
	log.Infof("%s <-> %s, %v recursive", src, dst, force)
	_, res, err := cli.BlobLs(src)
	if err != nil {
		return err
	}
	for _, val := range res {
		srcSubPath := filepath.Join(src, val.BaseName)
		dstSubPath := filepath.Join(dst, val.BaseName)

		if val.IsDir() {
			// Create local directory
			info, err := os.Stat(dstSubPath)
			if err != nil {
				err = os.Mkdir(dstSubPath, 0775)
				if err != nil {
					return err
				}
				log.Info("creating: ", dstSubPath)
			} else {
				if !info.IsDir() {
					log.Warning("skipping existing file: ", dstSubPath)
				}
			}
			if err = downloadRecursive(cli, srcSubPath, dstSubPath, force); err != nil {
				return err
			}
		} else {
			info, err := os.Stat(dstSubPath)
			if err == nil {
				if info.IsDir() {
					log.Warning("skipping existing directory: ", dstSubPath)
				} else {
					if !force {
						log.Warning("skipping existing file: ", dstSubPath)
					} else {
						log.Warning("replacing existing file: ", dstSubPath)
						err = os.Remove(dstSubPath)
						if err != nil {
							return err
						}
						err = download(cli, srcSubPath, dstSubPath, force)
					}
				}
			} else {
				err = download(cli, srcSubPath, dstSubPath, force)
			}

		}
	}
	return nil
}

func execGet(cli v1.NameServerClient, srcPath string, dstPath string, srcIsDir bool, dstIsDir bool, dstExist bool, force bool) error {
	if srcIsDir {
		if !dstExist {
			err := os.Mkdir(dstPath, 0775)
			if err != nil {
				return err
			}
		}
		return downloadRecursive(cli, srcPath, dstPath, force)
	} else {
		if dstIsDir {
			actualDstPath := filepath.Join(dstPath, filepath.Base(srcPath))
			return download(cli, srcPath, actualDstPath, force)
		} else {
			return download(cli, srcPath, dstPath, force)
		}
	}
}

func Get(cmd *cobra.Command, args []string) {
	opt := config.NewClientOpt()
	vipCfg, err := opt.Parse(cmd)
	if err != nil {
		log.Println("cannot find credential, run login first")
	} else {
		recursive, _ := cmd.Flags().GetBool("recursive")
		force, _ := cmd.Flags().GetBool("force")
		cli := v1.NewNameServerClient(opt.Token, opt.Hostname, opt.Port, opt.UseTLS)
		defer refreshToken(cli, vipCfg)

		srcPath := args[0]
		isDir, _, err := cli.BlobLs(srcPath)
		if err != nil {
			log.Error("source directory does not exist")
			return
		}
		if isDir && !recursive {
			log.Error("source is a directory, use -r to download recursively")
			return
		}

		srcIsDir := isDir

		dstPath := args[1]
		dstParent, _ := filepath.Split(dstPath)

		var dstIsDir bool
		var dstExist bool
		info, err := os.Stat(dstPath)
		if err == nil {
			dstExist = true
			dstIsDir = info.IsDir()
			if srcIsDir && !dstIsDir {
				log.Error("cannot download directory to a file")
				return
			}
		} else {
			_, err := os.Stat(dstParent)
			if err == nil {
				dstExist = false
				dstIsDir = srcIsDir
			} else {
				log.Error("destination directory does not exist")
			}
		}
		err = execGet(cli, srcPath, dstPath, srcIsDir, dstIsDir, dstExist, force)
		if err != nil {
			log.Errorln(err)
		}
		//log.Infof("srcIsDir: %v, dstIsDir: %v, dstExist: %v, recursive %v, force %v", srcIsDir, dstIsDir, dstExist, recursive, force)
	}
}
