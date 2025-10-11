package main

import (
	"github.com/duwa-lang/duwa/src/cli/docs"
	"github.com/spf13/cobra"
)

var generateDocs = &cobra.Command{
	Use:   "docs",
	Short: "Generate Docs",
	Run: func(cmd *cobra.Command, args []string) {
		docs.GenerateDocs()
	},
}
