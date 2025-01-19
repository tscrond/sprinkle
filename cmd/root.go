package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sprinkle",
	Short: "Sprinkle - Proxmox resource provisioning as a CLI",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	workdir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current working dir: ", err)
	}

	defaultDBPath := fmt.Sprintf("%s/.storage/sprinkle.db", workdir)

	rootCmd.PersistentFlags().String("config", workdir+"/"+"sprinkle-config.yaml", "Infrastructure config file path")
	rootCmd.PersistentFlags().String("api-url", "", "Proxmox API address (example: proxmox.example.com)")
	rootCmd.PersistentFlags().String("target-node", "", "Target PVE node name")
	rootCmd.PersistentFlags().String("db-path", defaultDBPath, "Path to DB storing the state")

	rootCmd.AddCommand(ApplyConfig)
}
