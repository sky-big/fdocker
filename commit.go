package fdocker

import (
	"fmt"
	"os/exec"

	log "common/clog"
	"fdocker/container/config"

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
	mntURL := fmt.Sprintf(config.MntUrl, containerName)
	mntURL += "/"

	imageTar := config.RootUrl + "/" + imageName + ".tar"

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Blog.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
