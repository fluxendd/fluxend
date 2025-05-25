package commands

import (
	"fluxend/internal/config"
	"fluxend/pkg"
	"github.com/spf13/cobra"
)

// aboutCmd Prints information about the application
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Prints information about the application",
	Run: func(cmd *cobra.Command, args []string) {
		aboutFluxend()
	},
}

func aboutFluxend() {
	pkg.DumpJSON(config.AboutFluxend)
}
