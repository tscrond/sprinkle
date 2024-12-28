package auth

import "github.com/spf13/cobra"

func init() {
	Auth.AddCommand(pveLogin)
	Auth.AddCommand(testConn)
}

var Auth = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Proxmox VE",
	Run:   func(cmd *cobra.Command, args []string) {},
}
