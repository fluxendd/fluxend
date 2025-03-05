package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// aboutCmd Prints information about the application
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Prints information about the application",
	Run: func(cmd *cobra.Command, args []string) {
		aboutApplication()
	},
}

// TODO: implement later
func aboutApplication() {
	// container := InitializeContainer()

	log.Println("Print information about the application")
}
