package main

import (
	"fmt"
	"os/exec"

	"github.com/sky-big/fdocker/container/config"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

// commit command
var CommitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into image.",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name and image name")
		}
		containerName := context.Args().Get(0)
		imageName := context.Args().Get(1)
		commitContainer(containerName, imageName)
		return nil
	},
}

func commitContainer(containerName, imageName string) {
	mntURL := config.MntPath + containerName
	mntURL += "/"

	imageTar := config.Root + imageName + ".tar"

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		glog.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
