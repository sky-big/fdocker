package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sky-big/fdocker/container/config"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

// log command
var LogCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "out",
			Usage: "stdout log",
		},
		cli.BoolFlag{
			Name:  "err",
			Usage: "stderr log",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Please input your container name")
		}
		containerName := context.Args().Get(0)
		stdoutlog := context.Bool("out")
		stderrlog := context.Bool("err")
		logContainer(containerName, stdoutlog, stderrlog)
		return nil
	},
}

func logContainer(containerName string, stdout, stderr bool) {
	dirURL := fmt.Sprintf(config.DefaultInfoLocation, containerName)
	logFileLocation := ""
	if stderr {
		logFileLocation = dirURL + config.ContainerErrFile
	}
	if stdout {
		logFileLocation = dirURL + config.ContainerLogFile
	}

	file, err := os.OpenFile(logFileLocation, os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		glog.Errorf("Log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		glog.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}
	fmt.Fprint(os.Stdout, string(content))

	// clean file
	file1, _ := os.OpenFile(logFileLocation, os.O_RDWR|os.O_TRUNC, 0666)
	defer file1.Close()
}
