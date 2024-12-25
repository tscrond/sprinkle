package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/lib"
)

func init() {
	createMachine.Flags().String("target-node", "", "Target PVE node name")
	createMachine.Flags().String("type", "lxc", "Determine type of the machine (lxc or vm)")

	createMachine.Flags().Bool("start-on-boot", false, "Start Machine on Boot")

	createMachine.Flags().Int("id", 200, "LXC container/VM ID")
	createMachine.Flags().Int("disk-size", 30, "Disk size for container/VM")
	createMachine.Flags().Int("swap-size", 0, "Swap size for container/VM")
	createMachine.Flags().Int("cpus", 2, "CPU Cores Count")

	createMachine.Flags().String("os-template", "debian-11-standard_11.7-1_amd64.tar.zst", "Name for the OS template (located in vztmpl in one of the storage systems) (e.g., template.tar.gz)")
	createMachine.Flags().String("network-bridge", "vmbr0", "Network bridge to use (default: vmbr0)")
	createMachine.Flags().String("network-interface", "eth0", "Network interface name (default: eth0)")
	createMachine.Flags().String("default-gateway", "192.168.1.1", "Default gateway for the container (default: 192.168.1.1)")
	createMachine.Flags().String("ip-address", "192.168.1.150/24", "Static IP address for the container (default: 192.168.1.150)")
	createMachine.Flags().String("storage-backend", "local-lvm", "Storage backend to use (example: local, local-lvm, ceph etc.)")
	createMachine.Flags().String("name", "default", "Name for the machine")
	createMachine.Flags().String("iso", "ubuntu-22.04.3-live-server-amd64.iso", "ISO for VM")

}

var createMachine = &cobra.Command{
	Use:   "create",
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
		fmt.Printf("ISO: %s\n", iso)

		machineConfig := lib.MachineConfig{
			ID:               id,
			OsTemplate:       osTemplate,
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
		}

		result, err := lib.CreateMachine(machineType, apiNode, targetNode, machineConfig)
		if err != nil {
			fmt.Println("Errors: ", err)
		}

		fmt.Println("Result: ", result)
	},
}
