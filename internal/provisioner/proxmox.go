package provisioner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/tscrond/sprinkle/internal/db"
	"github.com/tscrond/sprinkle/internal/pveclient"
	"github.com/tscrond/sprinkle/pkg/lib"
	"golang.org/x/crypto/ssh"
)

type ProxmoxProvisioner struct {
	PveClient *pveclient.PVEClient
}

func NewProxmoxProvisioner(creds *db.Credentials) *ProxmoxProvisioner {
	return &ProxmoxProvisioner{
		PveClient: &pveclient.PVEClient{
			Credentials: creds,
			Client:      &http.Client{},
		},
	}
}

func (p *ProxmoxProvisioner) CreateMachine(apiUrl, targetNode string, machineConfig *db.MachineConfig) error {
	fullApiURL := ""
	if machineConfig.Type == "lxc" {
		// Construct the URL
		fullApiURL = fmt.Sprintf("/api2/json/nodes/%s/lxc", targetNode)
	} else if machineConfig.Type == "vm" {
		// Construct the URL
		fullApiURL = fmt.Sprintf("/api2/json/nodes/%s/qemu", targetNode)
	} else {
		log.Fatalln("No such machine type defined: ", machineConfig.Type)
	}
	// // Create form data
	data := url.Values{}

	// If tags are present, add them to the request
	if machineConfig.Tags != "" {
		data.Set("tags", machineConfig.Tags)
	}

	var result error
	if machineConfig.Type == "lxc" {
		// Parameters specific to LXC
		data = p.SetLXCParams(machineConfig, data)
		fmt.Println("net0", fmt.Sprintf("name=%s,bridge=%s,ip=%s,gw=%s", machineConfig.NetworkInterface, machineConfig.NetworkBridge, machineConfig.IPAddress, machineConfig.DefaultGateway))
		result = p.makeMachineCreationRequest(fullApiURL, data)

	} else if machineConfig.Type == "vm" {
		// Parameters specific to VM
		fmt.Println("using cloudinit? ", machineConfig.UsingCloudInit)
		if machineConfig.UsingCloudInit {
			result = p.CreateVmUsingCloudInit(machineConfig)
		} else {
			data = p.SetVMParams(machineConfig, data)
			fmt.Println("net0", fmt.Sprintf("model=virtio,bridge=%s", machineConfig.NetworkBridge))
			result = p.makeMachineCreationRequest(fullApiURL, data)
		}
	}

	fmt.Println(fullApiURL)

	return result
}

func (p *ProxmoxProvisioner) makeMachineCreationRequest(fullApiURL string, data url.Values) error {
	pveClient := pveclient.NewPVEClient(p.PveClient.Credentials, &http.Client{})

	resp, err := pveClient.NewRequest("POST", fullApiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to make HTTP request: %v", err)
		return err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
		return err
	}
	fmt.Printf("Completed VM creation request with code %d: %s\n", statusCode, string(body))

	// fmt.Println("FULL API URL:", fullApiURL)

	return err
}

func (p *ProxmoxProvisioner) CreateVmUsingCloudInit(machineCfg *db.MachineConfig) error {
	user := lib.TrimLastSuffixAfter(p.PveClient.Credentials.Username, "@")
	password := p.PveClient.Credentials.Password
	address := lib.TrimSuffixAfter(p.PveClient.Credentials.ApiUrl, ":")

	fmt.Println("user: ", user)
	fmt.Println("password: ", password)
	fmt.Println("address: ", address)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password), // Use password authentication
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Ignore host key verification (not recommended for production)
	}

	client, err := ssh.Dial("tcp", address+":22", config)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	script := CLOUDINIT_SCRIPT

	sshkeys := func() []string {
		allkeys := []string{}
		for _, k := range machineCfg.SSHPublicKeys {
			allkeys = append(allkeys, *k.Key)
		}
		return allkeys
	}()

	cmd := fmt.Sprintf("bash -s -- \"%s\" \"%d\" \"%s\" \"%s\" \"%s\" \"%d\" \"%d\" \"%dG\" \"%s\" \"%s\" \"%s\" \"%s\" <<'EOF'\n%s\nEOF",
		machineCfg.IPAddress,
		machineCfg.VmId,
		machineCfg.Name,
		strings.Join(sshkeys, ","),
		machineCfg.DefaultGateway,
		machineCfg.Memory,
		machineCfg.CPUs,
		machineCfg.DiskSize,
		machineCfg.StorageBackend,
		machineCfg.ISO,
		machineCfg.NetworkBridge,
		machineCfg.Tags,
		script)

	// fmt.Println(cmd)

	// os.Exit(0)

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Printf("failed to run script: %v, output: \n%s\n", err, string(output))
		return errors.New("script_failed")
	}

	fmt.Println("dupa")

	fmt.Println("[i] Script output: ", string(output))

	return nil
}

func (p *ProxmoxProvisioner) DestroyMachine(apiUrl, targetNode string, machineConfig *db.MachineConfig) error {
	return nil
}

func (p *ProxmoxProvisioner) ApplyNewState(targetHost string, state []db.HostConfig) error {

	var hostConfig *db.HostConfig
	for _, hostConf := range state {
		if hostConf.TargetNode == targetHost {
			hostConfig = &hostConf
		}
	}
	if hostConfig == nil {
		log.Println("error: didnt find any host")
		return errors.New("host_is_nil")
	}

	machinesToCreate, machinesToDestroy, err := p.GetMachinesToCreateAndDestroy(hostConfig)
	if err != nil {
		return err
	}

	for _, m := range machinesToCreate {
		fmt.Printf("Creating machine %+v\n", m.VmId)
		err := p.CreateMachine(hostConfig.ApiURL, targetHost, &m)
		err = nil
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, m := range machinesToDestroy {
		fmt.Printf("Destroying machine %+v\n", m.VmId)
		err := p.DestroyMachine(hostConfig.ApiURL, targetHost, &m)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (p *ProxmoxProvisioner) SetLXCParams(machineConfig *db.MachineConfig, data url.Values) url.Values {
	data.Set("onboot", strconv.Itoa(lib.Btoi(machineConfig.StartOnBoot))) // Start on boot (1 for true, 0 for false)

	data.Set("vmid", strconv.Itoa(machineConfig.VmId))     // VM or container ID
	data.Set("storage", machineConfig.StorageBackend)      // Storage backend (example: local-lvm, local, ceph etc.)
	data.Set("cores", strconv.Itoa(machineConfig.CPUs))    // CPU cores
	data.Set("memory", strconv.Itoa(machineConfig.Memory)) // Memory size in MB

	data.Set("hostname", machineConfig.Name)
	data.Set("ostemplate", fmt.Sprintf("%s:vztmpl/%s", machineConfig.TemplateBackend, machineConfig.OsTemplate))                                                                       // OS Template
	data.Set("rootfs", fmt.Sprintf("%s:%d", machineConfig.StorageBackend, machineConfig.DiskSize))                                                                                     // Disk size and storage
	data.Set("net0", fmt.Sprintf("name=%s,bridge=%s,ip=%s,gw=%s", machineConfig.NetworkInterface, machineConfig.NetworkBridge, machineConfig.IPAddress, machineConfig.DefaultGateway)) // Network
	data.Set("swap", strconv.Itoa(machineConfig.SwapSize))                                                                                                                             // Swap size in MB
	// data.Set("ssh-public-keys", strings.Join(sshKeysFromFiles, "\n"))

	return data
}

func (p *ProxmoxProvisioner) SetVMParams(machineConfig *db.MachineConfig, data url.Values) url.Values { // VM Name
	data.Set("onboot", strconv.Itoa(lib.Btoi(machineConfig.StartOnBoot))) // Start on boot (1 for true, 0 for false)
	// Parameters specific to VM
	fmt.Println("net0", fmt.Sprintf("model=virtio,bridge=%s", machineConfig.NetworkBridge))
	data.Set("name", machineConfig.Name)                                                                     // VM Name
	data.Set("ide0", fmt.Sprintf("%s:%.2f", machineConfig.StorageBackend, float64(machineConfig.DiskSize)))  // Disk size
	data.Set("ide2", fmt.Sprintf("%s:iso/%s,media=cdrom", machineConfig.TemplateBackend, machineConfig.ISO)) // CD/DVD drive with ISO
	data.Set("net0", fmt.Sprintf("model=virtio,bridge=%s", machineConfig.NetworkBridge))                     // Network model and bridge
	return data
}

func (p *ProxmoxProvisioner) ConfigureCloudInit(machineConfig *db.MachineConfig) url.Values {

	data := url.Values{}
	data.Set("ciuser", fmt.Sprintf("guwno"))
	data.Set("cipassword", fmt.Sprintf("guwno"))
	data.Set("ipconfig0", fmt.Sprintf("ip=%s,gw=%s", machineConfig.IPAddress, machineConfig.DefaultGateway))
	data.Set("sshkeys", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIC48nEoa2rRazXTxZ4anL+6CL2bGXTo6w6XcDpmcd3pE tomasz.skrond@boar.network")

	return data
}

func (p *ProxmoxProvisioner) GetMachinesToCreateAndDestroy(hostConfig *db.HostConfig) ([]db.MachineConfig, []db.MachineConfig, error) {

	// Objective: Find machines managed by YAML to check if they exist
	// get current machines in the node
	pveClient := pveclient.NewPVEClient(p.PveClient.Credentials, &http.Client{})

	path := fmt.Sprintf("/api2/json/nodes/%s/qemu", p.PveClient.Credentials.TargetNode)

	resp, err := pveClient.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// serialize all machine IDs from node to temporary struct
	type data struct {
		VmId int `json:"vmid"`
	}

	var machinesDecoded struct {
		Data []data `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&machinesDecoded); err != nil {
		return nil, nil, err
	}

	machinesPresentOnHost := []int{}

	// retrieve relevant data to int array (vm id array)
	for _, vmdata := range machinesDecoded.Data {
		machinesPresentOnHost = append(machinesPresentOnHost, vmdata.VmId)
	}

	// check which machines from config match desired state
	// Condition 1 (C1): if some machine exists on the node AND in DB but does not exist in YAML, remove it and sync DB with YAML
	// Condition 2 (C2): if some machine exists in YAML but does not exist in node/DB, create it in the node/DB
	// DISCLAIMER: in this context, the hostConfig is the YAML config but in the FORMAT OF DB MODEL, so the state is based on YAML INPUT
	machinesToCreate := []db.MachineConfig{}
	machinesToDestroy := []db.MachineConfig{}

	hostMachines := hostConfig.Machines
	for _, m := range hostMachines {
		// C1: if does not exist in YAML, then remove
		if slices.Contains(machinesPresentOnHost, m.VmId) {
			machinesToDestroy = append(machinesToDestroy, m)
		} else { // C2: if exists in YAML, create
			machinesToCreate = append(machinesToCreate, m)
		}
	}

	return machinesToCreate, machinesToDestroy, nil
}
