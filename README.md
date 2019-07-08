# Introduction

fdocker is a simple container runc implementation in go

# QuickStart

1. make sure operation system have aufs filesystem

```
make init

make run
```

# Usage

```
root@sky_big:~/gopath/src/github.com/sky-big/fdocker# fdocker
NAME:
   fdocker - fdocker is a simple container runtime for function invoke.

USAGE:
   fdocker [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     init     Init container process run user's process in container. Do not call it outside
     run      Create a container with namespace and cgroups limit ie: fdocker run -ti [image] [command]
     ps       list all the containers
     logs     print logs of a container
     exec     exec a command into container
     stop     stop a container
     rm       remove unused containers
     network  container network commands
     inspec   inspec a container into
     mem      get a container mem info
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version


root@sky_big:~/gopath/src/github.com/sky-big/fdocker# fdocker run -h
NAME:
   fdocker run - Create a container with namespace and cgroups limit ie: fdocker run -ti [image] [command]

USAGE:
   fdocker run [command options] [arguments...]

OPTIONS:
   --ti              enable tty
   -d                detach container
   -m value          memory limit
   --cpushare value  cpushare limit
   --cpuset value    cpuset limit
   --name value      container name
   -v value          volume
   -e value          set environment
   --net value       container network
   -p value          port mapping
   -u value          user command owner
   --images value    image store path
```

# Links

1. https://github.com/xianlubird/mydocker
