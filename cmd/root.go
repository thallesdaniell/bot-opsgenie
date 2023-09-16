package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "opslink",
	Short:   "Um utilitário de linha de comando para opslink",
	Long:    `Você pode usar este utilitário para interagir com o opslink.`,
	Example: "opslink nextOncall",
	Version: "1.0.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
