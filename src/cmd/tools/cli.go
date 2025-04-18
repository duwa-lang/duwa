package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	_ "github.com/tliron/commonlog/simple"
)

var rootCmd = &cobra.Command{
	Use:   "duwa-intl",
	Short: "duwa-intl",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(1)
	},
}

var generate = &cobra.Command{
	Use:   "gen",
	Short: "Generate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(lspCommand)
	rootCmd.AddCommand(generate)
	generate.AddCommand(generateDocs)
	generate.AddCommand(generateStdTypes)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
