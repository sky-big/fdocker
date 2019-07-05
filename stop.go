package main

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/sky-big/fdocker/container/config"
	"github.com/sky-big/fdocker/container/manager"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

// stop command
var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		stopContainer(containerName)
		return nil
	},
}

func stopContainer(containerName string) {
	containerInfo, err := manager.GetContainerInfoByName(containerName)
	if err != nil {
		glog.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	pidInt, err := strconv.Atoi(containerInfo.Pid)
	if err != nil {
		glog.Errorf("Conver pid from string to int error %v", err)
		return
	}

	if err := syscall.Kill(pidInt, syscall.SIGKILL); err != nil {
		glog.Errorf("Stop container %s error %v", containerName, err)
	}

	containerInfo.Status = config.STOP
	containerInfo.Pid = " "
	err = manager.UpdateContainerInfo(containerInfo)
	if err != nil {
		glog.Errorf("Update container info %v error %v", containerInfo, err)
	}
}
