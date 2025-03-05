package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "fluxton",
	Short: "Fluxton CLI for managing the BaaS platform",
	Long:  `Fluxton CLI allows you to start the server, run seeders, and inspect routes.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.AddCommand(seedCmd)
	RootCmd.AddCommand(routesCmd)
}
