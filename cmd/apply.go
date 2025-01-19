package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/config"
)

func init() {

}

var ApplyConfig = &cobra.Command{
	Use:   "apply",
	Short: "Apply a piece of configuration to fit your setup",
	Run: func(cmd *cobra.Command, args []string) {

		dbPath, _ := cmd.Flags().GetString("db-path")
		configFile, _ := cmd.Flags().GetString("config")
		fmt.Println(dbPath)

		configManager := config.NewConfigManager()
		configManager.LoadConfigFromYAML(configFile)
		config := configManager.Config

		
		fmt.Printf("%+v", config)

	},
}
