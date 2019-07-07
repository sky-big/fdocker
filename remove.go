package main

import (
	"fmt"

	"github.com/sky-big/fdocker/cgroups"
	"github.com/sky-big/fdocker/container/config"
	"github.com/sky-big/fdocker/container/logs"
	"github.com/sky-big/fdocker/container/manager"
	"github.com/sky-big/fdocker/container/volume"
	"github.com/sky-big/fdocker/network"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

// remove command
var RemoveCommand = cli.Command{
	Name:  "rm",
	Usage: "remove unused containers",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		removeContainer(containerName)
		return nil
	},
}

func removeContainer(containerName string) {
	containerInfo, err := manager.GetContainerInfoByName(containerName)
	if err != nil {
		glog.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	if containerInfo.Status != config.STOP {
		glog.Errorf("Couldn't remove running container")
		return
	}

	volume.DeleteWorkSpace(containerInfo.Volume, containerName)
	cgroupManager := cgroups.NewCgroupManager(containerInfo.Id)
	cgroupManager.Destroy()
	logs.DeleteLogFile(containerName)
	network.Init()
	network.Disconnect(containerInfo)
	manager.DeleteContainerInfo(containerName)
}
