package volume

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "common/clog"
	"fdocker/container/common"
	"fdocker/container/config"
)

//Create a AUFS filesystem as container root workspace
func NewWorkSpace(volume string) {
	CreateReadOnlyLayer()
	//	CreateWriteLayer(containerName)
	//	CreateMountPoint(containerName, imageName)
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs)
			log.Blog.Infof("NewWorkSpace volume urls %q", volumeURLs)
		} else {
			log.Blog.Infof("Volume parameter input is not correct.")
		}
	}
}

//Decompression tar image
func CreateReadOnlyLayer() error {
	unTarFolderUrl := filepath.Join(config.RootUrl, config.Runtime)
	exist, err := common.PathExists(unTarFolderUrl)
	if err != nil {
		log.Blog.Infof("Fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
		return err
	}
	if !exist {
		log.Blog.Errorf("runtime %s not exist", unTarFolderUrl)
		return err
	}
	return nil
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(config.WriteLayerUrl, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Blog.Infof("Mkdir write layer dir %s error. %v", writeURL, err)
	}
}

func MountVolume(volumeURLs []string) error {
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		log.Blog.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
	}
	containerUrl := volumeURLs[1]
	mntURL := filepath.Join(config.RootUrl, config.Runtime)
	containerVolumeURL := filepath.Join(mntURL, containerUrl)
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		log.Blog.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}
	//	dirs := "dirs=" + parentUrl
	//	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL).CombinedOutput()
	//	if err != nil {
	//		log.Blog.Errorf("Mount volume failed. %v", err)
	//		return err
	//	}

	_, err := exec.Command("mount", "--bind", parentUrl, containerVolumeURL).CombinedOutput()
	if err != nil {
		log.Blog.Errorf("Mount volume(%s, %s) failed. %v", parentUrl, containerVolumeURL, err)
		return err
	}
	return nil
}

func CreateMountPoint(containerName, imageName string) error {
	mntUrl := fmt.Sprintf(config.MntUrl, containerName)
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Blog.Errorf("Mkdir mountpoint dir %s error. %v", mntUrl, err)
		return err
	}
	tmpWriteLayer := fmt.Sprintf(config.WriteLayerUrl, containerName)
	tmpImageLocation := config.RootUrl + "/" + imageName
	mntURL := fmt.Sprintf(config.MntUrl, containerName)
	dirs := "dirs=" + tmpWriteLayer + ":" + tmpImageLocation
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL).CombinedOutput()
	if err != nil {
		log.Blog.Errorf("Run command for creating mount point failed %v", err)
		return err
	}
	return nil
}

//Delete the AUFS filesystem while container exit
func DeleteWorkSpace(volume string) {
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteVolume(volumeURLs)
		}
	}
	//	DeleteMountPoint(containerName)
	//	DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) error {
	mntURL := fmt.Sprintf(config.MntUrl, containerName)
	_, err := exec.Command("umount", mntURL).CombinedOutput()
	if err != nil {
		log.Blog.Errorf("Unmount %s error %v", mntURL, err)
		return err
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Blog.Errorf("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}

func DeleteVolume(volumeURLs []string) error {
	containerUrl := volumeURLs[1]
	mntURL := filepath.Join(config.RootUrl, config.Runtime)
	containerVolumeURL := filepath.Join(mntURL, containerUrl)

	if _, err := exec.Command("umount", containerVolumeURL).CombinedOutput(); err != nil {
		log.Blog.Errorf("Umount volume %s failed. %v", containerUrl, err)
		return err
	}
	return nil
}

func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(config.WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		log.Blog.Infof("Remove writeLayer dir %s error %v", writeURL, err)
	}
}
