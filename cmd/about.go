package cmd

import (
	"fluxton/configs"
	"fluxton/pkg"
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
	pkg.DumpJSON(configs.AboutFluxton)
}
