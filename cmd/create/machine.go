package create

import (
	"fmt"
	"math/rand/v2"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/lib"
)

var randint = rand.IntN(99999)

var DEFAULT_MACHINE_CONFIG = lib.MachineConfig{
	Name:             "machine-" + fmt.Sprintf("%d", randint),
	ID:               randint,
	MachineType:      "lxc",
	OsTemplate:       "debian-11-standard_11.7-1_amd64.tar.zst",
	NetworkBridge:    "vmbr0",
	NetworkInterface: "eth0",
	DefaultGateway:   "192.168.1.1",
	IPAddress:        "192.168.1.150/24",
	StorageBackend:   "local-lvm",
	TemplateBackend:  "local",
	DiskSize:         30,
	SwapSize:         0,
	CPUCount:         2,
	OnBoot:           false,
	ISO:              "ubuntu-22.04.3-live-server-amd64.iso",
	Tags:             "",
}

func init() {
	createMachine.Flags().String("type", DEFAULT_MACHINE_CONFIG.MachineType, "Determine type of the machine (lxc or vm)")

	createMachine.Flags().String("tags", "", "Tags for the machine, if more tags needed, enter with semicolon delimiter (for example: \"tag1;tag2\")")

	createMachine.Flags().Int("id", DEFAULT_MACHINE_CONFIG.ID, "LXC container/VM ID")                     //cluster-predefined (random)
	createMachine.Flags().Int("disk-size", DEFAULT_MACHINE_CONFIG.DiskSize, "Disk size for container/VM") //
	createMachine.Flags().Int("swap-size", DEFAULT_MACHINE_CONFIG.SwapSize, "Swap size for container/VM")
	createMachine.Flags().Int("cpus", DEFAULT_MACHINE_CONFIG.CPUCount, "CPU Cores Count")

}

var createMachine = &cobra.Command{
	Use:   "machine",
	Short: "Create a new LXC container/Virtual Machine",
	Run: func(cmd *cobra.Command, args []string) {

		apiNode, _ := cmd.Flags().GetString("api-url")
		targetNode, _ := cmd.Flags().GetString("target-node")
		machineType, _ := cmd.Flags().GetString("type")

		id, _ := cmd.Flags().GetInt("id")
		osTemplate, _ := cmd.Flags().GetString("os-template")
		networkBridge, _ := cmd.Flags().GetString("network-bridge")
		networkInterface, _ := cmd.Flags().GetString("network-interface")
		defaultGateway, _ := cmd.Flags().GetString("default-gateway")
		ipAddress, _ := cmd.Flags().GetString("ip-address")
		diskSize, _ := cmd.Flags().GetInt("disk-size")
		swapSize, _ := cmd.Flags().GetInt("swap-size")
		startOnBoot, _ := cmd.Flags().GetBool("start-on-boot")
		cpuCount, _ := cmd.Flags().GetInt("cpus")
		storageBackend, _ := cmd.Flags().GetString("storage-backend")
		machineName, _ := cmd.Flags().GetString("name")
		iso, _ := cmd.Flags().GetString("iso")
		tags, _ := cmd.Flags().GetString("tags")
		templateBackend, _ := cmd.Flags().GetString("template-backend")

		// Debug print for flag values
		fmt.Printf("Machine ID: %d\n", id)
		fmt.Printf("OS Template: %s\n", osTemplate)
		fmt.Printf("Network Bridge: %s\n", networkBridge)
		fmt.Printf("Network Interface: %s\n", networkInterface)
		fmt.Printf("Default Gateway: %s\n", defaultGateway)
		fmt.Printf("IP Address: %s\n", ipAddress)
		fmt.Printf("Disk Size: %d\n", diskSize)
		fmt.Printf("Swap Size: %d\n", swapSize)
		fmt.Printf("Start on Boot: %t\n", startOnBoot)
		fmt.Printf("CPU Count: %d\n", cpuCount)
		fmt.Printf("Storage Backend: %s\n", storageBackend)
		fmt.Printf("Template Backend: %s\n", templateBackend)
		fmt.Printf("ISO: %s\n", iso)
		fmt.Printf("Tags: %s\n", tags)

		machineConfig := lib.MachineConfig{
			ID:               id,
			MachineType:      machineType,
			OsTemplate:       osTemplate,
			TemplateBackend:  templateBackend,
			NetworkBridge:    networkBridge,
			NetworkInterface: networkInterface,
			DefaultGateway:   defaultGateway,
			IPAddress:        ipAddress,
			DiskSize:         diskSize,
			OnBoot:           startOnBoot,
			CPUCount:         cpuCount,
			StorageBackend:   storageBackend,
			Name:             machineName,
			ISO:              iso,
			Tags:             tags,
		}

		result, err := lib.CreateMachine(apiNode, targetNode, machineConfig)
		if err != nil {
			fmt.Println("Errors: ", err)
		}

		fmt.Println("Result: ", result)
	},
}
