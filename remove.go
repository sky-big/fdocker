package fdocker

import (
	"fmt"

	log "common/clog"
	"fdocker/cgroups"
	"fdocker/container/config"
	"fdocker/container/logs"
	"fdocker/container/manager"
	"fdocker/container/volume"
	"fdocker/network"

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
		log.Blog.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	if containerInfo.Status != config.STOP {
		log.Blog.Errorf("Couldn't remove running container")
		return
	}

	volume.DeleteWorkSpace(containerInfo.Volume)
	cgroupManager := cgroups.NewCgroupManager(containerInfo.Id)
	cgroupManager.Destroy()
	logs.DeleteLogFile(containerName)
	network.Init()
	network.Disconnect(containerInfo)
	manager.DeleteContainerInfo(containerName)
}
