package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search various sources for articles.",
	Long: `Search various sources for articles.

A subcommand is required in order to search a specific database.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("search called without a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
