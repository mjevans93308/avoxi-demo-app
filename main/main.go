package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mjevans93308/geolocate-ip-demo-app/api"
	"github.com/mjevans93308/geolocate-ip-demo-app/config"
)

func serverCmd() *cobra.Command {
	return &cobra.Command{
		Use: config.Serve_Command,
		PreRun: func(cmd *cobra.Command, args []string) {
			loadConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			a := api.App{}
			a.Initialize(false)
			a.Run(config.Address)

			return nil
		},
	}
}

func loadConfig() {
	viper.SetConfigType(config.ConfigFileType)
	viper.AddConfigPath(config.ConfigFileLocation)
	viper.SetConfigName(config.CobraCommandName)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func main() {
	cmd := &cobra.Command{
		Use:     config.CobraCommandName,
		Short:   "geolocate-ip Demo App",
		Version: "0.0.1",
	}

	cmd.AddCommand(serverCmd())

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
