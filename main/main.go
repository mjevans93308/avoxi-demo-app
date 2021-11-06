package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mjevans93308/avoxi-demo-app/api"
	"github.com/mjevans93308/avoxi-demo-app/config"
)

func serverCmd() *cobra.Command {
	return &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {

			a := api.App{}
			a.Initialize()
			a.Run(config.ADDRESS)

			return nil

		},
	}
}

func main() {
	cmd := &cobra.Command{
		Use:     "avoxi-demo-app",
		Short:   "Avoxi Demo App",
		Version: "0.1",
	}

	cmd.AddCommand(serverCmd())

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
