package lib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func CreateMachine(machineType, apiNode, targetNode string, config MachineConfig) (int, error) {

	var apiURL string
	if machineType == "lxc" {
		// Construct the URL
		apiURL = fmt.Sprintf("https://%s/api2/json/nodes/%s/lxc", apiNode, targetNode)
	} else if machineType == "vm" {
		// Construct the URL
		apiURL = fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu", apiNode, targetNode)
	} else {
		log.Fatalln("No such machine type defined: ", machineType)
	}

	// // Create form data
	data := url.Values{}
	data.Set("vmid", strconv.Itoa(config.ID)) // VM or container ID
	data.Set("storage", config.StorageBackend)

	// data.Set("rootfs", strconv.Itoa(config.DiskSize))     // Disk size (in GB)
	data.Set("storage", config.StorageBackend)       // Storage backend (example: local-lvm, local, ceph etc.)
	data.Set("cores", strconv.Itoa(config.CPUCount)) // CPU cores

	if machineType == "lxc" {
		// Parameters specific to LXC
		data.Set("hostname", config.Name)
		data.Set("ostemplate", fmt.Sprintf("%s:vztmpl/%s", config.StorageBackend, config.OsTemplate))                                                          // OS Template
		data.Set("rootfs", fmt.Sprintf("%s:%d", config.StorageBackend, config.DiskSize))                                                                       // Disk size and storage
		data.Set("net0", fmt.Sprintf("name=%s,bridge=%s,ip=%s,gw=%s", config.NetworkInterface, config.NetworkBridge, config.IPAddress, config.DefaultGateway)) // Network
		data.Set("swap", strconv.Itoa(config.SwapSize))                                                                                                        // Swap size in MB
	} else if machineType == "vm" {
		// Parameters specific to VM
		data.Set("name", config.Name)                                                             // VM Name
		data.Set("ide0", fmt.Sprintf("%s:%.2f", config.StorageBackend, float64(config.DiskSize))) // Disk size
		data.Set("ide2", fmt.Sprintf("%s:iso/%s,media=cdrom", config.StorageBackend, config.ISO)) // CD/DVD drive with ISO
		data.Set("net0", fmt.Sprintf("model=virtio,bridge=%s", config.NetworkBridge))             // Network model and bridge
	}

	data.Set("onboot", strconv.Itoa(btoi(config.OnBoot))) // Start on boot (1 for true, 0 for false)

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v\n", err)
	}

	client, req, err := ConfigureAuth(&http.Client{}, req)
	if err != nil {
		return -1, err
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
		return -1, nil
	}

	fmt.Println(string(body))

	// Print the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed: %s\n", resp.Status)
		return resp.StatusCode, nil
	}

	return resp.StatusCode, nil
}

// btoi is a helper function to convert bool to int (1 for true, 0 for false)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
