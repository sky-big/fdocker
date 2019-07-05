package logs

import (
	"fmt"
	"os"

	log "common/clog"
	"fdocker/container/config"
)

func NewLogFile(containerName, file string) *os.File {
	dirURL := fmt.Sprintf(config.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirURL, 0622); err != nil {
		log.Blog.Errorf("NewLogFile mkdir %s error %v", dirURL, err)
		return nil
	}
	stdLogFilePath := dirURL + file
	stdLogFile, err := os.OpenFile(stdLogFilePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil && os.IsNotExist(err) {
		stdLogFile, err = os.Create(stdLogFilePath)
		if err != nil {
			log.Blog.Errorf("NewLogFile create file %s error %v", stdLogFilePath, err)
			return nil
		}
	}
	stdLogFile.WriteString("")
	return stdLogFile
}

func DeleteLogFile(containerName string) {
	dirURL := fmt.Sprintf(config.DefaultInfoLocation, containerName)
	os.RemoveAll(dirURL)
}
