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
	OnBoot           bool
	ISO              string
	Tags             string
}

type ClusterConfig struct {
	WorkerConfig    MachineConfig
	MasterConfig    MachineConfig
	MasterNodeCount int
	WorkerNodeCount int
	ClusterName     string
}
