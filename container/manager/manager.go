package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	log "common/clog"

	"fdocker/container/config"
	"fdocker/container/types"
)

func SaveContainerInfo(containerPID int, commandArray []string, containerName, id, volume, ip, nwName string) error {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	containerInfo := &types.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      config.RUNNING,
		Name:        containerName,
		Volume:      volume,
		Ip:          ip,
		NetworkName: nwName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Blog.Errorf("Record container info error %v", err)
		return err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(config.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Blog.Errorf("Mkdir error %s error %v", dirUrl, err)
		return err
	}
	fileName := dirUrl + "/" + config.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Blog.Errorf("Create file %s error %v", fileName, err)
		return err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Blog.Errorf("File write string error %v", err)
		return err
	}

	return nil
}

func UpdateContainerInfo(containerInfo *types.ContainerInfo) error {
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Blog.Errorf("Json marshal %s error %v", containerInfo.Name, err)
		return err
	}
	dirURL := fmt.Sprintf(config.DefaultInfoLocation, containerInfo.Name)
	configFilePath := dirURL + config.ConfigName
	if err := ioutil.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		log.Blog.Errorf("Write file %s error", configFilePath, err)
		return err
	}
	return nil
}

func DeleteContainerInfo(containerName string) {
	dirURL := fmt.Sprintf(config.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Blog.Errorf("Remove dir %s error %v", dirURL, err)
	}
}

func GetContainerInfoByName(containerName string) (*types.ContainerInfo, error) {
	dirURL := fmt.Sprintf(config.DefaultInfoLocation, containerName)
	configFilePath := dirURL + config.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Blog.Errorf("Read file %s error %v", configFilePath, err)
		return nil, err
	}
	var containerInfo types.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		log.Blog.Errorf("GetContainerInfoByName unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}
