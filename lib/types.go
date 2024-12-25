package lib

type MachineConfig struct {
	Name             string
	ID               int
	OsTemplate       string
	NetworkBridge    string
	NetworkInterface string
	DefaultGateway   string
	IPAddress        string
	StorageBackend   string
	DiskSize         int
	SwapSize         int
	CPUCount         int
	OnBoot           bool
	ISO              string
}

/*

	// Request parameters
	vmid := "601"
	ostemplate := "local:vztmpl/debian-8.0-standard_8.0-1_amd64.tar.gz"
	ipAddress := "192.168.1.100/24" // Example IP address with CIDR notation
	bridge := "vmbr0"
	defaultGateway := "192.168.1.1"
	net0 := fmt.Sprintf("name=myct0,bridge=%s,ip=%s,gw=%s", bridge, ipAddress, defaultGateway)


	// Construct the URL
	apiURL := fmt.Sprintf("https://%s:8006/api2/json/nodes/%s/lxc", apiNode, targetNode)


*/
