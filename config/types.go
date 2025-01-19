package config

type HostConfigYAML struct {
	Hosts map[string]struct {
		ApiUrl     string `mapstructure:"api-url"`
		TargetNode string `mapstructure:"target-node"`
		LXCs       struct {
			Default  MachineConfigYAML   `mapstructure:"default"`
			Machines []MachineConfigYAML `mapstructure:"machines"`
		} `mapstructure:"lxc"`
		VMs struct {
			Default  MachineConfigYAML   `mapstructure:"default"`
			Machines []MachineConfigYAML `mapstructure:"machines"`
		} `mapstructure:"vm"`
	} `mapstructure:"hosts"`
}

type MachineConfigYAML struct {
	Name             string   `mapstructure:"name"`
	VmId             int      `mapstructure:"vmid"`
	OsTemplate       string   `mapstructure:"os-template,omitempty"`
	ISO              string   `mapstructure:"iso,omitempty"`
	IPAddress        string   `mapstructure:"ip-address"`
	CPUs             int      `mapstructure:"cpus"`
	Memory           int      `mapstructure:"memory"`
	DiskSize         int      `mapstructure:"disk-size"`
	SwapSize         int      `mapstructure:"swap-size"`
	Tags             string   `mapstructure:"tags"`
	StartOnBoot      bool     `mapstructure:"start-on-boot,omitempty"`
	StorageBackend   string   `mapstructure:"storage-backend,omitempty"`
	TemplateBackend  string   `mapstructure:"template-backend,omitempty"`
	NetworkBridge    string   `mapstructure:"network-bridge,omitempty"`
	NetworkInterface string   `mapstructure:"network-interface,omitempty"`
	DefaultGateway   string   `mapstructure:"default-gateway,omitempty"`
	SshPublicKeys    []SSHKey `mapstructure:"ssh-public-keys"`
}

type SSHKey struct {
	Path string `mapstructure:"path,omitempty"`
	Key  string `mapstructure:"key,omitempty"`
}
