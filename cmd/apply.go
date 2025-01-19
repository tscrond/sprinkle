package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/config"
	"github.com/tscrond/sprinkle/internal/configmapper"
	"github.com/tscrond/sprinkle/pkg/lib"
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
		config, err := configManager.LoadConfigFromYAML(configFile)
		if err != nil {
			panic(err)
		}

		config = configmapper.PropagateDefaults(config)
		fmt.Println("")
		lib.PrettyPrintStruct(config)

		// dbConfig, err := configmapper.ConvertConfigToDBModel(config)
		// if err != nil {
		// 	panic(err)
		// }
		// lib.PrettyPrintStruct(dbConfig)
	},
}
