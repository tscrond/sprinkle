package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/config"
	"github.com/tscrond/sprinkle/internal/auth"
	"github.com/tscrond/sprinkle/internal/configmapper"
	"github.com/tscrond/sprinkle/internal/db"
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

		authService := auth.NewAuthService(db)
		creds, err := authService.Authenticate("genesis", "192.168.1.102:8006")
		if err != nil {
			log.Fatalln(err)
		}
		// fmt.Printf("%+v", creds)
		lib.PrettyPrintStruct(creds)
	},
}
