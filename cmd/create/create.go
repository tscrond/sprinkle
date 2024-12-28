package create

import "github.com/spf13/cobra"

func init() {

	CreateResource.PersistentFlags().String("target-node", "", "Target PVE node name")
	CreateResource.PersistentFlags().String("type", "lxc", "Node type (lxc or vm)")

	CreateResource.PersistentFlags().String("os-template", "debian-11-standard_11.7-1_amd64.tar.zst", "Name for the OS template located in vztmpl in one of the storage systems (example: debian-11-standard_11.7-1_amd64.tar.zst, template.tar.gz)")
	CreateResource.PersistentFlags().String("network-bridge", "vmbr0", "Network bridge to use")
	CreateResource.PersistentFlags().String("network-interface", "eth0", "Network interface name")
	CreateResource.PersistentFlags().String("default-gateway", "192.168.1.1", "Default gateway for the container/VM")
	CreateResource.PersistentFlags().String("storage-backend", "local-lvm", "Storage backend to use for machine disks (example: local, local-lvm, ceph etc.)")
	CreateResource.PersistentFlags().String("template-backend", "local", "Storage backend to use for ISO/OS Templates (example: local, local-lvm, ceph etc.)")
	CreateResource.PersistentFlags().String("iso", "ubuntu-22.04.3-live-server-amd64.iso", "ISO for VM")
	CreateResource.PersistentFlags().Bool("start-on-boot", false, "Start Machine on Boot")

	CreateResource.PersistentFlags().String("config-file", "", "Specify config file to read configuration from")

	CreateResource.AddCommand(createCluster)
	CreateResource.AddCommand(createMachine)
}

var CreateResource = &cobra.Command{
	Use:   "create",
	Short: "Create a resource (cluster or machine)",
	Run:   func(cmd *cobra.Command, args []string) {},
}
