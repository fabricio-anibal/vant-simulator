package commands

import (
	"vantsimulator/internal/processor"

	"github.com/spf13/cobra"
)

var simCmd = &cobra.Command{
	Use:   "sim",
	Short: "Inicia o simulador de VANT",
	Run: func(cmd *cobra.Command, args []string) {
		process()
	},
}

func process() {
	processor.Process()

}
