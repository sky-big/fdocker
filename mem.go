package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sky-big/fdocker/container/manager"

	log "github.com/Sirupsen/logrus"
	"github.com/containerd/cgroups"
	"github.com/urfave/cli"
)

var GetMemCommand = cli.Command{
	Name:  "mem",
	Usage: "get a container mem info",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		GetMemInfo(containerName)
		return nil
	},
}

func GetMemInfo(containerName string) {
	containerInfo, err := manager.GetContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}

	path := "/" + containerInfo.Id
	control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(path))
	if err != nil {
		log.Errorf("GetMemInfo error:%v", err)
		return
	}

	metrics, err := control.Stat()
	if err != nil {
		log.Errorf("GetMemInfo error:%v", err)
		return
	}

	memInfo, err := json.Marshal(metrics.Memory.Usage)
	if err != nil {
		log.Errorf("GetMemInfo error:%v", err)
		return
	}

	fmt.Fprint(os.Stdout, string(memInfo))
}
