package volume

import (
	"os"
	"os/exec"
	"strings"

	"github.com/sky-big/fdocker/container/common"
	"github.com/sky-big/fdocker/container/config"

	// "github.com/golang/glog"
	glog "github.com/Sirupsen/logrus"
)

// Create a AUFS filesystem as container root workspace
func NewWorkSpace(volume, imageName, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName, imageName)
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs, containerName)
			glog.Infof("NewWorkSpace volume urls %q", volumeURLs)
		} else {
			glog.Infof("Volume parameter input is not correct.")
		}
	}
}

// Decompression tar image
func CreateReadOnlyLayer(imageName string) error {
	unTarFolderUrl := config.Root + imageName + "/"
	imageUrl := config.ImageStorePath + imageName + ".tar"
	exist, err := common.PathExists(unTarFolderUrl)
	if err != nil {
		glog.Infof("Fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
		return err
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			glog.Errorf("Mkdir %s error %v", unTarFolderUrl, err)
			return err
		}

		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			glog.Errorf("Untar dir %s error %v", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}

func CreateWriteLayer(containerName string) {
	writeURL := config.WritePath + containerName
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		glog.Infof("Mkdir write layer dir %s error. %v", writeURL, err)
	}
}

func MountVolume(volumeURLs []string, containerName string) error {
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		glog.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
	}
	containerUrl := volumeURLs[1]
	mntURL := config.MntPath + containerName
	containerVolumeURL := mntURL + "/" + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		glog.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}
	dirs := "dirs=" + parentUrl
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL).CombinedOutput()
	if err != nil {
		glog.Errorf("Mount volume failed. %v", err)
		return err
	}
	return nil
}

func CreateMountPoint(containerName, imageName string) error {
	mntUrl := config.MntPath + containerName
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		glog.Errorf("Mkdir mountpoint dir %s error. %v", mntUrl, err)
		return err
	}
	tmpWriteLayer := config.WritePath + containerName
	tmpImageLocation := config.Root + imageName
	mntURL := config.MntPath + containerName
	dirs := "dirs=" + tmpWriteLayer + ":" + tmpImageLocation
	glog.Infof("-------------:%v, %v", dirs, mntURL)
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL).CombinedOutput()
	if err != nil {
		glog.Errorf("Run command for creating mount point failed %v", err)
		return err
	}
	return nil
}

// Delete the AUFS filesystem while container exit
func DeleteWorkSpace(volume, containerName string) {
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteVolume(volumeURLs, containerName)
		}
	}
	DeleteMountPoint(containerName)
	DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) error {
	mntURL := config.MntPath + containerName
	_, err := exec.Command("umount", mntURL).CombinedOutput()
	if err != nil {
		glog.Errorf("Unmount %s error %v", mntURL, err)
		return err
	}
	if err := os.RemoveAll(mntURL); err != nil {
		glog.Errorf("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}

func DeleteVolume(volumeURLs []string, containerName string) error {
	mntURL := config.MntPath + containerName
	containerUrl := mntURL + "/" + volumeURLs[1]
	if _, err := exec.Command("umount", containerUrl).CombinedOutput(); err != nil {
		glog.Errorf("Umount volume %s failed. %v", containerUrl, err)
		return err
	}
	return nil
}

func DeleteWriteLayer(containerName string) {
	writeURL := config.WritePath + containerName
	if err := os.RemoveAll(writeURL); err != nil {
		glog.Infof("Remove writeLayer dir %s error %v", writeURL, err)
	}
}
