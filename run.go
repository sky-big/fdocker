package fdocker

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	log "common/clog"
	"fdocker/cgroups"
	"fdocker/cgroups/subsystems"
	"fdocker/container"
	"fdocker/container/common"
	"fdocker/container/config"
	"fdocker/container/logs"
	"fdocker/container/manager"
	"fdocker/container/types"
	fvolume "fdocker/container/volume"
	"fdocker/network"

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
		log.Blog.Infof("createTty %v", createTty)
		containerName := context.String("name")
		volume := context.String("v")
		network := context.String("net")
		user := context.String("u")

		envSlice := context.StringSlice("e")
		portmapping := context.StringSlice("p")

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

	log.Blog.Infof("parent %s starting", containerName)

	parent := container.NewParentProcess(tty, containerName, volume, imageName, user, envSlice)
	if parent == nil {
		log.Blog.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Blog.Error(err)
	}

	// use containerID as cgroup name
	cgroupManager := cgroups.NewCgroupManager(containerID)
	cgroupManager.Set(res)
	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		log.Blog.Error(err)
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
			log.Blog.Errorf("Error Connect Network %v", err)
			return
		} else {
			ipStr = ip
		}
	}

	// record container info
	err := manager.SaveContainerInfo(parent.Process.Pid, comArray, containerName, containerID, volume, ipStr, nw)
	if err != nil {
		log.Blog.Errorf("Record container info error %v", err)
		return
	}

	log.Blog.Infof("parent %s store meta data success", containerName)

	err = saveInitCommand(containerName, comArray)
	if err != nil {
		log.Blog.Warningf("parent process save init command error : %v", err)
		return
	}

	if tty {
		parent.Wait()
		manager.DeleteContainerInfo(containerName)
		fvolume.DeleteWorkSpace(volume)
		cgroupManager.Destroy()
		logs.DeleteLogFile(containerName)
	}
}

func saveInitCommand(containerName string, comArray []string) error {
	command := strings.Join(comArray, " ")
	log.Blog.Infof("command all is %s", command)

	filePath := fmt.Sprintf(config.DefaultInfoLocation, containerName) + config.InitCommandFile
	filePathBack := filePath + ".bk"
	f, err := os.Create(filePathBack)
	if err != nil {
		log.Blog.Warningf("save init command create file error : %v", err)
		return err
	}

	_, err = io.WriteString(f, command)
	if err != nil {
		log.Blog.Warningf("save init command write error : %v", err)
		return err
	}

	if err := os.Rename(filePathBack, filePath); err != nil {
		log.Blog.Warningf("change filename(%s) to filename(%s) err: %s", filePathBack, filePath, err.Error())
		return err
	}

	return nil
}
