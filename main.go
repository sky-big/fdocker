package main

import (
	"os"

	"common/clog"
	"fdocker"

	"github.com/urfave/cli"
)

const (
	usage = `fdocker is a simple container runtime for function invoke.`

	LogPath  = "/var/run/function"
	LogLevel = "DEBUG"
)

func main() {
	clog.LogInit(clog.LogConfig{
		LogDir:   LogPath,
		ToStderr: false,
		Level:    LogLevel,
	})

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
		clog.Flush()
	}
}
