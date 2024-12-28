package create

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/lib"
)

func init() {
	randomNumber := strconv.Itoa(rand.IntN(9999999))
	defaultClusterName := "cluster-" + randomNumber

	createCluster.Flags().String("preset", "none", "Cluster presets (none, small, medium, large)")

	createCluster.Flags().String("cluster-name", defaultClusterName, "Cluster Name")

	createCluster.Flags().Int("worker-nodes", 3, "Amount of worker nodes")
	createCluster.Flags().Int("master-nodes", 1, "Amount of master nodes")

	createCluster.Flags().Int("worker-node-disk", 20, "Disk size of a worker node")
	createCluster.Flags().Int("master-node-disk", 30, "Disk size of a master node")

	createCluster.Flags().Int("worker-node-cpus", 1, "Core count for worker node")
	createCluster.Flags().Int("master-node-cpus", 2, "Core count for master node")

	createCluster.Flags().Int("worker-node-memory", 2048, "Memory for worker node")
	createCluster.Flags().Int("master-node-memory", 2048, "Memory count for master node")

	createCluster.Flags().String("worker-node-tags", fmt.Sprintf("k8s;worker;%s", defaultClusterName), "Tags for worker nodes")
	createCluster.Flags().String("master-node-tags", fmt.Sprintf("k8s;worker;%s", defaultClusterName), "Tags for master nodes")

	createCluster.Flags().String("ip-range", "", "IP range for the cluster")
}

var createCluster = &cobra.Command{
	Use:   "cluster",
	Short: "Create a cluster/group of VMs (optimized for Kubernetes)",
	Run: func(cmd *cobra.Command, args []string) {
		apiNode, _ := cmd.Flags().GetString("api-url")
		targetNode, _ := cmd.Flags().GetString("target-node")

		// Read generic flags
		nodeType, _ := cmd.Flags().GetString("type")
		nodeOSTemplate, _ := cmd.Flags().GetString("os-template")
		nodeISO, _ := cmd.Flags().GetString("iso")
		storageBackend, _ := cmd.Flags().GetString("storage-backend")
		templateBackend, _ := cmd.Flags().GetString("template-backend")
		networkBridge, _ := cmd.Flags().GetString("network-bridge")
		networkInterface, _ := cmd.Flags().GetString("network-interface")
		defaultGateway, _ := cmd.Flags().GetString("default-gateway")
		onBoot, _ := cmd.Flags().GetBool("start-on-boot")

		// Read cluster-specialized flags
		clusterName, _ := cmd.Flags().GetString("cluster-name")
		preset, _ := cmd.Flags().GetString("preset")
		workerNodes, _ := cmd.Flags().GetInt("worker-nodes")
		masterNodes, _ := cmd.Flags().GetInt("master-nodes")
		workerNodeDisk, _ := cmd.Flags().GetInt("worker-node-disk")
		masterNodeDisk, _ := cmd.Flags().GetInt("master-node-disk")
		workerNodeCores, _ := cmd.Flags().GetInt("worker-node-cpus")
		masterNodeCores, _ := cmd.Flags().GetInt("master-node-cpus")
		workerNodeMemory, _ := cmd.Flags().GetInt("worker-node-memory")
		masterNodeMemory, _ := cmd.Flags().GetInt("master-node-memory")
		workerNodeTags, _ := cmd.Flags().GetString("worker-node-tags")
		masterNodeTags, _ := cmd.Flags().GetString("master-node-tags")
		ipRange, _ := cmd.Flags().GetString("ip-range")
		sshKeys, _ := cmd.Flags().GetStringSlice("ssh-pubkey-file")

		// Print the values for debugging (optional)
		fmt.Printf("Cluster Name: %s\n", clusterName)
		fmt.Printf("Preset: %s\n", preset)
		fmt.Printf("Node Type: %s\n", nodeType)

		fmt.Printf("Worker Node Tags: %s\n", workerNodeTags)
		fmt.Printf("Master Node Tags: %s\n", masterNodeTags)
		fmt.Printf("Node OS Template: %s\n", nodeOSTemplate)
		fmt.Printf("Node ISO: %s\n", nodeISO)

		fmt.Printf("Storage Backend: %s\n", storageBackend)
		fmt.Printf("Template Backend: %s\n", templateBackend)
		fmt.Printf("Network Bridge: %s\n", networkBridge)
		fmt.Printf("Network Interface: %s\n", networkInterface)
		fmt.Printf("Default Gateway: %s\n", networkInterface)
		fmt.Printf("IP Range: %s\n", ipRange)
		fmt.Printf("Start On Boot: %t\n", onBoot)
		fmt.Printf("SSH Key files to read: %s\n", sshKeys)

		if preset == "none" {
			fmt.Printf("Worker Nodes: %d\n", workerNodes)
			fmt.Printf("Master Nodes: %d\n", masterNodes)
			fmt.Printf("Worker Node Disk: %dGB\n", workerNodeDisk)
			fmt.Printf("Master Node Disk: %dGB\n", masterNodeDisk)
			fmt.Printf("Worker Node Cores: %d\n", workerNodeCores)
			fmt.Printf("Master Node Cores: %d\n", masterNodeCores)
			fmt.Printf("Worker Node Memory: %d\n", workerNodeMemory)
			fmt.Printf("Master Node Memory: %d\n", masterNodeMemory)
		}

		workerNodeConfig := lib.MachineConfig{}
		masterNodeConfig := lib.MachineConfig{}

		clusterConfig := &lib.ClusterConfig{}

		if preset == "none" {

			workerNodeConfig = lib.MachineConfig{
				// Name:             "machine-" + fmt.Sprintf("%d", randint),
				// ID:               randint,
				MachineType:      nodeType,
				OsTemplate:       nodeOSTemplate,
				NetworkBridge:    networkBridge,
				NetworkInterface: networkInterface,
				// DefaultGateway:   "192.168.1.1",
				// IPAddress:        "192.168.1.150/24",
				StorageBackend:  storageBackend,
				TemplateBackend: templateBackend,
				DiskSize:        workerNodeDisk,
				SwapSize:        0,
				CPUCount:        workerNodeCores,
				Memory:          workerNodeMemory,
				OnBoot:          onBoot,
				ISO:             nodeISO,
				Tags:            workerNodeTags + fmt.Sprintf(";%s", clusterName),
				SshKeys:         sshKeys,
			}

			masterNodeConfig = lib.MachineConfig{
				// Name:             "machine-" + fmt.Sprintf("%d", randint),
				// ID:               randint,
				MachineType:      nodeType,
				OsTemplate:       nodeOSTemplate,
				NetworkBridge:    networkBridge,
				NetworkInterface: networkInterface,
				// DefaultGateway:   "192.168.1.1",
				// IPAddress:        "192.168.1.150/24",
				StorageBackend:  storageBackend,
				TemplateBackend: templateBackend,
				DiskSize:        masterNodeDisk,
				Memory:          masterNodeMemory,
				SwapSize:        0,
				CPUCount:        masterNodeCores,
				OnBoot:          onBoot,
				ISO:             nodeISO,
				Tags:            masterNodeTags + fmt.Sprintf(";%s", clusterName),
				SshKeys:         sshKeys,
			}

			clusterConfig = &lib.ClusterConfig{
				WorkerConfig:    workerNodeConfig,
				MasterConfig:    masterNodeConfig,
				WorkerNodeCount: workerNodes,
				MasterNodeCount: masterNodes,
				ClusterName:     clusterName,
				IPRange:         ipRange,
			}

		} else {

			clusterConfig = returnClusterConfigFromPreset(clusterName, preset)
			if clusterConfig == nil {
				log.Fatalln("cluster config is undefined")
			}

			clusterConfig.IPRange = ipRange

			// i hate this
			clusterConfig.WorkerConfig.MachineType = nodeType
			clusterConfig.WorkerConfig.OsTemplate = nodeOSTemplate
			clusterConfig.WorkerConfig.NetworkBridge = networkBridge
			clusterConfig.WorkerConfig.NetworkInterface = networkInterface
			clusterConfig.WorkerConfig.DefaultGateway = defaultGateway
			clusterConfig.WorkerConfig.StorageBackend = storageBackend
			clusterConfig.WorkerConfig.TemplateBackend = templateBackend
			clusterConfig.WorkerConfig.ISO = nodeISO
			clusterConfig.WorkerConfig.OnBoot = onBoot
			clusterConfig.WorkerConfig.SshKeys = sshKeys

			// what the fuck
			clusterConfig.MasterConfig.MachineType = nodeType
			clusterConfig.MasterConfig.OsTemplate = nodeOSTemplate
			clusterConfig.MasterConfig.NetworkBridge = networkBridge
			clusterConfig.MasterConfig.NetworkInterface = networkInterface
			clusterConfig.MasterConfig.DefaultGateway = defaultGateway
			clusterConfig.MasterConfig.StorageBackend = storageBackend
			clusterConfig.MasterConfig.TemplateBackend = templateBackend
			clusterConfig.MasterConfig.ISO = nodeISO
			clusterConfig.MasterConfig.OnBoot = onBoot
			clusterConfig.MasterConfig.SshKeys = sshKeys

			fmt.Printf("Worker Nodes: %d\n", clusterConfig.WorkerNodeCount)
			fmt.Printf("Master Nodes: %d\n", clusterConfig.MasterNodeCount)
			fmt.Printf("Worker Node Disk: %dGB\n", clusterConfig.WorkerConfig.DiskSize)
			fmt.Printf("Master Node Disk: %dGB\n", clusterConfig.MasterConfig.DiskSize)
			fmt.Printf("Worker Node Cores: %d\n", clusterConfig.WorkerConfig.CPUCount)
			fmt.Printf("Master Node Cores: %d\n", clusterConfig.MasterConfig.CPUCount)
			fmt.Printf("Worker Node Memory: %d\n", clusterConfig.WorkerConfig.Memory)
			fmt.Printf("Master Node Memory: %d\n", clusterConfig.MasterConfig.Memory)
		}

		_, err := lib.CreateCluster(apiNode, targetNode, *clusterConfig)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func returnClusterConfigFromPreset(clusterName, preset string) *lib.ClusterConfig {
	switch preset {
	case "small":
		fmt.Println("Small preset")
		clusterconf := &lib.ClusterConfig{
			MasterConfig:    SMALL_PRESET_MASTER,
			WorkerConfig:    SMALL_PRESET_WORKER,
			WorkerNodeCount: 3,
			MasterNodeCount: 1,
			ClusterName:     clusterName,
		}
		clusterconf.MasterConfig.Tags += fmt.Sprintf(";%s", clusterName)
		clusterconf.WorkerConfig.Tags += fmt.Sprintf(";%s", clusterName)

		return clusterconf

	case "balanced":
		fmt.Println("Balanced preset")
		clusterconf := &lib.ClusterConfig{
			MasterConfig:    SMALL_PRESET_MASTER,
			WorkerConfig:    MEDIUM_PRESET_WORKER,
			WorkerNodeCount: 5,
			MasterNodeCount: 3,
			ClusterName:     clusterName,
		}

		clusterconf.MasterConfig.Tags += fmt.Sprintf(";%s", clusterName)
		clusterconf.WorkerConfig.Tags += fmt.Sprintf(";%s", clusterName)

		return clusterconf

	case "balanced-storage":
		fmt.Println("Balanced storage-oriented preset")
		clusterconf := &lib.ClusterConfig{
			MasterConfig:    SMALL_STORAGE_PRESET_MASTER,
			WorkerConfig:    MEDIUM_STORAGE_PRESET_WORKER,
			WorkerNodeCount: 5,
			MasterNodeCount: 3,
			ClusterName:     clusterName,
		}

		clusterconf.MasterConfig.Tags += fmt.Sprintf(";%s", clusterName)
		clusterconf.WorkerConfig.Tags += fmt.Sprintf(";%s", clusterName)

		return clusterconf

	case "medium":
		fmt.Println("Medium preset")
		clusterconf := &lib.ClusterConfig{
			MasterConfig:    MEDIUM_PRESET_MASTER,
			WorkerConfig:    MEDIUM_PRESET_WORKER,
			WorkerNodeCount: 5,
			MasterNodeCount: 3,
			ClusterName:     clusterName,
		}

		clusterconf.MasterConfig.Tags += fmt.Sprintf(";%s", clusterName)
		clusterconf.WorkerConfig.Tags += fmt.Sprintf(";%s", clusterName)

		return clusterconf

	case "large":
		fmt.Println("Large preset")
		clusterconf := &lib.ClusterConfig{
			MasterConfig:    LARGE_PRESET_MASTER,
			WorkerConfig:    LARGE_PRESET_WORKER,
			WorkerNodeCount: 9,
			MasterNodeCount: 5,
			ClusterName:     clusterName,
		}
		clusterconf.MasterConfig.Tags += fmt.Sprintf(";%s", clusterName)
		clusterconf.WorkerConfig.Tags += fmt.Sprintf(";%s", clusterName)

		return clusterconf

	default:
		fmt.Println("no such preset")
		return nil
	}
}
