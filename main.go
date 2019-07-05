package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

const (
	usage = `fdocker is a simple container runtime for function invoke.`
)

func main() {
	app := cli.NewApp()
	app.Name = "fdocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		fdocker.InitCommand,
		fdocker.RunCommand,
		fdocker.ListCommand,
		fdocker.LogCommand,
		fdocker.ExecCommand,
		fdocker.StopCommand,
		fdocker.RemoveCommand,
		fdocker.NetworkCommand,
		fdocker.InspecCommand,
		fdocker.GetMemCommand,
	}

	app.Before = func(context *cli.Context) error {
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		glog.Warningf("fdocker run error : ", err)
	}
}
