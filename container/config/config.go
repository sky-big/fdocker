package config

const (
	RUNNING             string = "running"
	STOP                string = "stopped"
	Exit                string = "exited"
	DefaultInfoLocation string = "/var/run/fdocker/%s/"
	InitCommandFile     string = "init"
	ConfigName          string = "config.json"
	ContainerLogFile    string = "stdout.log"
	ContainerErrFile    string = "stderr.log"
	RootUrl             string = "/var"
	Runtime             string = "runtime"
	MntUrl              string = "/root/mnt/%s"
	WriteLayerUrl       string = "/root/writeLayer/%s"
)
