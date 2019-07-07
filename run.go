package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/sky-big/fdocker/cgroups"
	"github.com/sky-big/fdocker/cgroups/subsystems"
	"github.com/sky-big/fdocker/container"
	"github.com/sky-big/fdocker/container/common"
	"github.com/sky-big/fdocker/container/config"
	"github.com/sky-big/fdocker/container/logs"
	"github.com/sky-big/fdocker/container/manager"
	"github.com/sky-big/fdocker/container/types"
	fvolume "github.com/sky-big/fdocker/container/volume"
	"github.com/sky-big/fdocker/network"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

// run command
var RunCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit ie: fdocker run -ti [image] [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "set environment",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping",
		},
		cli.StringFlag{
			Name:  "u",
			Usage: "user command owner",
		},
		cli.StringFlag{
			Name:  "images",
			Usage: "image store path",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		//get image name
		imageName := cmdArray[0]
		cmdArray = cmdArray[1:]

		createTty := context.Bool("ti")
		detach := context.Bool("d")

		if createTty && detach {
			return fmt.Errorf("ti and d paramter can not both provided")
		}
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}
		glog.Infof("createTty %v", createTty)
		containerName := context.String("name")
		volume := context.String("v")
		network := context.String("net")
		user := context.String("u")

		envSlice := context.StringSlice("e")
		portmapping := context.StringSlice("p")

		// image store path
		imageStorePath := context.String("images")
		if imageStorePath != "" {
			config.ImageStorePath = imageStorePath
		}

		Run(createTty, cmdArray, resConf, containerName, volume, imageName, envSlice, network, user, portmapping)
		return nil
	},
}

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, containerName, volume, imageName string,
	envSlice []string, nw, user string, portmapping []string) {
	containerID := common.RandStringBytes(10)
	if containerName == "" {
		containerName = containerID
	}

	glog.Infof("parent %s starting", containerName)

	parent := container.NewParentProcess(tty, containerName, volume, imageName, user, envSlice)
	if parent == nil {
		glog.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		glog.Error(err)
	}

	// use containerID as cgroup name
	cgroupManager := cgroups.NewCgroupManager(containerID)
	cgroupManager.Set(res)
	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		glog.Error(err)
	}

	ipStr := ""
	if nw != "" {
		// config container network
		network.Init()
		containerInfo := &types.ContainerInfo{
			Id:          containerID,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: portmapping,
		}
		if ip, err := network.Connect(nw, containerInfo); err != nil {
			glog.Errorf("Error Connect Network %v", err)
			return
		} else {
			ipStr = ip
		}
	}

	// record container info
	err := manager.SaveContainerInfo(parent.Process.Pid, comArray, containerName, containerID, volume, ipStr, nw)
	if err != nil {
		glog.Errorf("Record container info error %v", err)
		return
	}

	glog.Infof("parent %s store meta data success", containerName)

	// record init process run command
	err = saveInitCommand(containerName, comArray)
	if err != nil {
		glog.Warningf("parent process save init command error : %v", err)
		return
	}

	if tty {
		parent.Wait()
		manager.DeleteContainerInfo(containerName)
		fvolume.DeleteWorkSpace(volume, containerName)
		cgroupManager.Destroy()
		logs.DeleteLogFile(containerName)
	}
}

func saveInitCommand(containerName string, comArray []string) error {
	command := strings.Join(comArray, " ")
	glog.Infof("command all is %s", command)

	filePath := fmt.Sprintf(config.DefaultInfoLocation, containerName) + config.InitCommandFile
	filePathBack := filePath + ".bk"
	f, err := os.Create(filePathBack)
	if err != nil {
		glog.Warningf("save init command create file error : %v", err)
		return err
	}

	_, err = io.WriteString(f, command)
	if err != nil {
		glog.Warningf("save init command write error : %v", err)
		return err
	}

	if err := os.Rename(filePathBack, filePath); err != nil {
		glog.Warningf("change filename(%s) to filename(%s) err: %s", filePathBack, filePath, err.Error())
		return err
	}

	return nil
}
