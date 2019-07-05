package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/sky-big/fdocker/container/config"
	"github.com/sky-big/fdocker/container/manager"
	"github.com/sky-big/fdocker/container/types"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

// list command
var ListCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the containers",
	Action: func(context *cli.Context) error {
		ListContainers()
		return nil
	},
}

func ListContainers() {
	dirURL := fmt.Sprintf(config.DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		glog.Errorf("Read dir %s error %v", dirURL, err)
		return
	}

	var containers []*types.ContainerInfo
	for _, file := range files {
		if file.Name() == "network" {
			continue
		}
		tmpContainer, err := manager.GetContainerInfoByName(file.Name())
		if err != nil {
			glog.Errorf("Get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}
	if err := w.Flush(); err != nil {
		glog.Errorf("Flush error %v", err)
		return
	}
}
