package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// arxivCmd represents the arxiv command
var arxivCmd = &cobra.Command{
	Use:   "arxiv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Thank you to arXiv for use of its open access interoperability.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("arxiv called")
	},
}

func init() {
	searchCmd.AddCommand(arxivCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// arxivCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// arxivCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
