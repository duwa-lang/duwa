package main

import (
	"github.com/sevenreup/duwa/src/cli/std"
	"github.com/spf13/cobra"
)

var generateStdTypes = &cobra.Command{
	Use:   "std-types",
	Short: "Generate Std Types",
	Run: func(cmd *cobra.Command, args []string) {
		std.GenerateStdTypes()
	},
}
