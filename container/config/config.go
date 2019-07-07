package config

const (
	RUNNING                      string = "running"
	STOP                         string = "stopped"
	Exit                         string = "exited"
	DefaultInfoLocation          string = "/var/run/fdocker/%s/"
	InitCommandFile              string = "init"
	ConfigName                   string = "config.json"
	ContainerLogFile             string = "stdout.log"
	ContainerErrFile             string = "stderr.log"
	DefaultNetworkLocation       string = "/var/run/fdocker/network/network/"
	IpamDefaultAllocatorLocation string = "/var/run/fdocker/network/ipam/subnet.json"
	Root                         string = "/var/run/fdocker/root/"
	MntPath                      string = Root + "mnt/"
	WritePath                    string = Root + "writelayer/"
)

var (
	ImageStorePath string = Root
)
