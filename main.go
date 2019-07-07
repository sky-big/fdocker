package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

const (
	usage = `fdocker is a simple container runtime for function invoke.`
)

func main() {
	flag.Parse()

	app := cli.NewApp()
	app.Name = "fdocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		InitCommand,
		RunCommand,
		ListCommand,
		LogCommand,
		ExecCommand,
		StopCommand,
		RemoveCommand,
		NetworkCommand,
		InspecCommand,
		GetMemCommand,
	}

	app.Before = func(context *cli.Context) error {
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		glog.Warningf("fdocker run error : ", err)
	}
}
