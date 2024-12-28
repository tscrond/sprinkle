package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sprinkle",
	Short: "Sprinkle - Proxmox resource provisioning as a CLI",
	Long:  `Create Promox VMs/LXC containers using Proxmox API`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("api-url", "", "Proxmox API address (example: proxmox.example.com)")

	rootCmd.AddCommand(pveLogin)
	rootCmd.AddCommand(testConn)

	rootCmd.AddCommand(createResource)
}
