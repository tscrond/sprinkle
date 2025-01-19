package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/config"
	"github.com/tscrond/sprinkle/internal/auth"
	"github.com/tscrond/sprinkle/internal/configmapper"
	"github.com/tscrond/sprinkle/internal/db"
)

func init() {

}

var ApplyConfig = &cobra.Command{
	Use:   "apply",
	Short: "Apply a piece of configuration to fit your setup",
	Run: func(cmd *cobra.Command, args []string) {

		dbPath, _ := cmd.Flags().GetString("db-path")
		configFile, _ := cmd.Flags().GetString("config")
		db, err := db.NewResourceRepository(dbPath)
		if err != nil {
			panic(err)
		}

		configManager := config.NewConfigManager()
		config, err := configManager.LoadConfigFromYAML(configFile)
		if err != nil {
			panic(err)
		}

		config = configmapper.PropagateDefaults(config)
		authConfigs := configmapper.MapConfigToAuthConfig(config)

		authService := auth.NewAuthService(db, authConfigs)
		fmt.Println(authService)
	},
}
