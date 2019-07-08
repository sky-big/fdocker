package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/sky-big/fdocker/container/manager"
	_ "github.com/sky-big/fdocker/container/nsenter"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const ENV_EXEC_PID = "fdocker_pid"
const ENV_EXEC_CMD = "fdocker_cmd"

// exec command
var ExecCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) error {
		//This is for callback
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", os.Getgid())
			return nil
		}

		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name or command")
		}
		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		ExecContainer(containerName, commandArray)
		return nil
	},
}

func ExecContainer(containerName string, comArray []string) {
	containerInfo, err := manager.GetContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}

	cmdStr := strings.Join(comArray, " ")
	log.Infof("container pid %s", containerInfo.Pid)
	log.Infof("command %s", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, containerInfo.Pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)
	containerEnvs := getEnvsByPid(containerInfo.Pid)
	cmd.Env = append(os.Environ(), containerEnvs...)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error %v", containerName, err)
	}
}

func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file %s error %v", path, err)
		return nil
	}
	//env split by \u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}
