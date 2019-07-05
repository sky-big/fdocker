package fdocker

import (
	"encoding/json"
	"fmt"
	"os"

	log "common/clog"
	"fdocker/container/manager"

	"github.com/urfave/cli"
)

// inspec command
var InspecCommand = cli.Command{
	Name:  "inspec",
	Usage: "inspec a container into",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		InspecContainer(containerName)
		return nil
	},
}

func InspecContainer(containerName string) {
	containerInfo, err := manager.GetContainerInfoByName(containerName)
	if err != nil {
		log.Blog.Errorf("Get container %s error %v", containerName, err)
		return
	}

	content, err := json.Marshal(containerInfo)
	if err != nil {
		log.Blog.Errorf("Marshal container info error %v", err)
	}
	fmt.Fprint(os.Stdout, string(content))
}
