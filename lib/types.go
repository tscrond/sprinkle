package lib

type MachineConfig struct {
	Name             string
	ID               int
	MachineType      string
	OsTemplate       string
	NetworkBridge    string
	NetworkInterface string
	DefaultGateway   string
	IPAddress        string
	StorageBackend   string
	TemplateBackend  string
	DiskSize         int
	SwapSize         int
	CPUCount         int
	Memory           int
	OnBoot           bool
	ISO              string
	Tags             string
	SshKeys          []string
}

type ClusterConfig struct {
	WorkerConfig    MachineConfig
	MasterConfig    MachineConfig
	MasterNodeCount int
	WorkerNodeCount int
	ClusterName     string
	IPRange         string // like 192.168.1.140-150/24
}
