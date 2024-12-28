package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CreateMachine(apiNode, targetNode string, config MachineConfig) (int, error) {

	var apiURL string
	if config.MachineType == "lxc" {
		// Construct the URL
		apiURL = fmt.Sprintf("https://%s/api2/json/nodes/%s/lxc", apiNode, targetNode)
	} else if config.MachineType == "vm" {
		// Construct the URL
		apiURL = fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu", apiNode, targetNode)
	} else {
		log.Fatalln("No such machine type defined: ", config.MachineType)
	}

	// // Create form data
	data := url.Values{}
	data.Set("vmid", strconv.Itoa(config.ID)) // VM or container ID
	data.Set("storage", config.StorageBackend)

	data.Set("storage", config.StorageBackend)       // Storage backend (example: local-lvm, local, ceph etc.)
	data.Set("cores", strconv.Itoa(config.CPUCount)) // CPU cores

	if config.MachineType == "lxc" {
		// Parameters specific to LXC
		data.Set("hostname", config.Name)
		data.Set("ostemplate", fmt.Sprintf("%s:vztmpl/%s", config.TemplateBackend, config.OsTemplate))                                                         // OS Template
		data.Set("rootfs", fmt.Sprintf("%s:%d", config.StorageBackend, config.DiskSize))                                                                       // Disk size and storage
		data.Set("net0", fmt.Sprintf("name=%s,bridge=%s,ip=%s,gw=%s", config.NetworkInterface, config.NetworkBridge, config.IPAddress, config.DefaultGateway)) // Network
		data.Set("swap", strconv.Itoa(config.SwapSize))                                                                                                        // Swap size in MB
	} else if config.MachineType == "vm" {
		// Parameters specific to VM
		data.Set("name", config.Name)                                                              // VM Name
		data.Set("ide0", fmt.Sprintf("%s:%.2f", config.StorageBackend, float64(config.DiskSize)))  // Disk size
		data.Set("ide2", fmt.Sprintf("%s:iso/%s,media=cdrom", config.TemplateBackend, config.ISO)) // CD/DVD drive with ISO
		data.Set("net0", fmt.Sprintf("model=virtio,bridge=%s", config.NetworkBridge))              // Network model and bridge
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

func CreateCluster(apiNode, targetNode string, clusterConfig ClusterConfig) (int, error) {

	ipRangeNumbers, cidr, baseIP, err := parseIPRange(clusterConfig.IPRange)
	if err != nil {
		fmt.Println("Error parsing IP range:", err)
		return -1, err
	}

	baseWorkerName := clusterConfig.WorkerConfig.Name
	baseMasterName := clusterConfig.MasterConfig.Name

	globalIndex := 0
	for i := range clusterConfig.WorkerNodeCount {

		nodeID := rand.IntN(9999999)
		clusterConfig.WorkerConfig.Name = baseWorkerName + fmt.Sprintf("-%d", i)
		clusterConfig.WorkerConfig.ID = assignIDToNode(apiNode, nodeID)
		clusterConfig.WorkerConfig.IPAddress = baseIP + fmt.Sprintf("%d", ipRangeNumbers[globalIndex]) + cidr
		if clusterConfig.WorkerConfig.ID == -1 {
			return -1, errors.New("logic_error")
		}
		fmt.Println(clusterConfig.WorkerConfig)

		if result, err := CreateMachine(apiNode, targetNode, clusterConfig.WorkerConfig); err != nil {
			fmt.Println("Error: ", err)
			return result, err
		}
		globalIndex += 1
	}

	for i := range clusterConfig.MasterNodeCount {
		nodeID := rand.IntN(9999999)
		clusterConfig.MasterConfig.Name = baseMasterName + fmt.Sprintf("-%d", i)
		clusterConfig.MasterConfig.ID = assignIDToNode(apiNode, nodeID)
		clusterConfig.MasterConfig.IPAddress = baseIP + fmt.Sprintf("%d", ipRangeNumbers[globalIndex]) + cidr
		if clusterConfig.MasterConfig.ID == -1 {
			return -1, errors.New("logic_error")
		}
		fmt.Println(clusterConfig.MasterConfig)

		if result, err := CreateMachine(apiNode, targetNode, clusterConfig.MasterConfig); err != nil {
			fmt.Println("Error: ", err)
			return result, err
		}
		globalIndex += 1
	}

	return 0, nil
}

func parseIPRange(ipRange string) ([]int, string, string, error) {
	// Split the IP range into the range part and the CIDR part
	parts := strings.Split(ipRange, "/")
	if len(parts) != 2 {
		return nil, "", "", fmt.Errorf("invalid IP range format: %s", ipRange)
	}

	rangePart := parts[0]
	cidr := "/" + parts[1]

	// Extract the numeric range
	ipParts := strings.Split(rangePart, ".")
	if len(ipParts) != 4 {
		return nil, "", "", fmt.Errorf("invalid IP format: %s", rangePart)
	}

	// Extract the base IP (e.g., 192.168.1.)
	baseIP := fmt.Sprintf("%s.%s.%s.", ipParts[0], ipParts[1], ipParts[2])

	// Extract the range (e.g., 140-150) from the last octet
	rangeTokens := strings.Split(ipParts[3], "-")
	if len(rangeTokens) != 2 {
		return nil, "", "", fmt.Errorf("invalid range in IP: %s", ipParts[3])
	}

	start, err := strconv.Atoi(rangeTokens[0])
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid start of range: %s", rangeTokens[0])
	}

	end, err := strconv.Atoi(rangeTokens[1])
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid end of range: %s", rangeTokens[1])
	}

	// Generate the array of numbers
	var ipRangeNumbers []int
	for i := start; i <= end; i++ {
		ipRangeNumbers = append(ipRangeNumbers, i)
	}

	return ipRangeNumbers, cidr, baseIP, nil
}
