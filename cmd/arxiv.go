package cmd

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/lehigh-university-libraries/papercut/pkg/arxiv"
	"github.com/spf13/cobra"
)

var (
	// used for flags.
	start   int
	results int
	ids     string
	query   string

	arxivCmd = &cobra.Command{
		Use:   "arxiv",
		Short: "Search arXiv for articles",
		Long: `Search arXiv for articles.

Thank you to arXiv for use of its open access interoperability.`,
		Run: func(cmd *cobra.Command, args []string) {
			if query == "" && ids == "" {
				log.Fatal("query or ids required.")
			}

			params := url.Values{}
			if query != "" {
				params.Set("search_query", query)
			}
			if ids != "" {
				params.Set("id_list", ids)
			}

			params.Set("start", strconv.Itoa(start))
			params.Set("max_results", strconv.Itoa(results))

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				log.Fatal(err)
			}
			apiURL := fmt.Sprintf("%s?%s", url, params.Encode())

			log.Printf("Accessing %s\n", apiURL)

			resp, err := http.Get(apiURL)

			if err != nil {
				fmt.Println("Error requesting XML data:", err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			var result arxiv.Feed
			err = xml.Unmarshal(body, &result)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			wr := csv.NewWriter(os.Stdout)

			// header
			wr.Write([]string{"id", "published", "updated", "title", "doi", "pdf"})

			for _, e := range result.Entries {
				wr.Write([]string{
					e.ID,
					e.Published.String(),
					e.Updated.String(),
					e.Title,
					e.DOI,
					e.PDF,
				})
				wr.Flush()
			}

		},
	}
)

func init() {
	searchCmd.AddCommand(arxivCmd)

	arxivCmd.Flags().StringP("url", "u", "https://export.arxiv.org/api/query", "The arXiv API url")
	arxivCmd.Flags().StringVarP(&query, "query", "q", "", "The arXiv API url")
	arxivCmd.Flags().StringVarP(&ids, "ids", "i", "", "The arXiv API url")
	arxivCmd.Flags().IntVarP(&start, "start", "s", 0, "The offset")
	arxivCmd.Flags().IntVarP(&results, "results", "r", 10, "The number of results to return in a response")
}
