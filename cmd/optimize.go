package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// optimizeCmd Flush all caches and optimize the application
var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Flush all caches and optimize the application",
	Run: func(cmd *cobra.Command, args []string) {
		optimize()
	},
}

// TODO: implement later
func optimize() {
	// container := InitializeContainer()

	log.Println("Flush all caches and optimize the application")
}
