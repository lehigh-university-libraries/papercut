package cmd

import (
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search various sources for articles.",
	Long: `Search various sources for articles.

A subcommand is required in order to search a specific database.`,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
