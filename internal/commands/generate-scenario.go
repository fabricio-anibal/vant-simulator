package commands

import (
	"strconv"
	"vantsimulator/internal/processor"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Gera casos para de VANT na pasta /data",
	Run: func(cmd *cobra.Command, args []string) {
		processGenerator(args)
	},
}

func processGenerator(args []string) {
	qtdVants, _ := strconv.Atoi(args[0])
	qtdFiles, _ := strconv.Atoi(args[1])
	processor.ProcessGenerator(qtdVants, qtdFiles)
}