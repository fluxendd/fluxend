package cmd

import (
	"fluxton/configs"
	"fluxton/utils"
	"github.com/spf13/cobra"
)

// aboutCmd Prints information about the application
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Prints information about the application",
	Run: func(cmd *cobra.Command, args []string) {
		aboutFluxton()
	},
}

func aboutFluxton() {
	utils.DumpJSON(configs.AboutFluxton)
}
