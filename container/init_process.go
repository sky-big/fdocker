package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/sky-big/fdocker/container/common"
	"github.com/sky-big/fdocker/container/config"
	userOp "github.com/sky-big/fdocker/container/user"

	log "github.com/Sirupsen/logrus"
)

const (
	CheckInterval = 10
)

func RunContainerInitProcess(containerName, user string) error {
	log.Infof("containername : %s, user : %s", containerName, user)
	cmdArray, err := readUserCommand(containerName)
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error %v, cmdArray is %v", err, cmdArray)
	}

	log.Infof("Run container cmd : %v", cmdArray)

	// pivot root file system
	setUpMount()

	// set user
	if user != "" {
		if err := userOp.SetUser(user); err != nil {
			log.Errorf("Set User error %v", err)
			return err
		}
	}

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func readUserCommand(containerName string) ([]string, error) {
	for {
		time.Sleep(CheckInterval * time.Microsecond)
		filePath := fmt.Sprintf(config.DefaultInfoLocation, containerName) + config.InitCommandFile
		if common.CheckFileIsExist(filePath) {
			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Warningf("init process read command error : %v", err)
				return make([]string, 0), err
			}

			return strings.Split(string(b), " "), nil
		}
	}
	return make([]string, 0), nil
}

/**
Init 挂载点
*/
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)
	pivotRoot(pwd)

	//mount proc
	//defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	//syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	//syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotRoot := ".pivot_root_" + common.RandStringBytes(5)
	pivotDir := filepath.Join(root, pivotRoot)
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", pivotRoot)
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
