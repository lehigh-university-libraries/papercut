package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the search command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get articles.",
	Long: `Fetch PDFs and/or metadata for articles.

A subcommand is required in order to fetch the article from a specific source.`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
