package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/lehigh-university-libraries/papercut/internal/utils"
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
			queries := []string{}

			dl, err := cmd.Flags().GetString("directory-listing")
			if err != nil {
				log.Fatal(err)
			}
			if dl != "" {
				if query != "" || ids != "" {
					log.Fatal("query or ids can not be used with the directory listing option.")
				}
				queries, err = utils.FetchEmails(dl)
				if err != nil {
					log.Fatal(err)
				}
			}

			if query != "" {
				queries = append(queries, query)
			}

			if len(queries) == 0 && ids == "" {
				log.Fatal("query or ids required.")
			}

			wr := csv.NewWriter(os.Stdout)

			// header
			wr.Write([]string{"id", "published", "updated", "title", "doi", "pdf", "query"})

			for _, query := range queries {

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

				result, err := arxiv.GetResults(apiURL)
				if err != nil {
					log.Fatal(err)
				}
				for true {
					for _, e := range result.Entries {
						wr.Write([]string{
							e.ID,
							e.Published.String(),
							e.Updated.String(),
							e.Title,
							e.DOI,
							e.PDF,
							query,
						})
						wr.Flush()
					}

					log.Println("Pausing between requests. arXiv requests a three second delay between API requests...")
					time.Sleep(3 * time.Second)
					next := result.StartIndex + result.ItemsPerPage
					if result.TotalResults > next {
						params.Set("start", strconv.Itoa(next))
						apiURL := fmt.Sprintf("%s?%s", url, params.Encode())
						log.Printf("Accessing %s\n", apiURL)
						result, err = arxiv.GetResults(apiURL)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						break
					}
				}
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
	arxivCmd.Flags().String("directory-listing", "", "URL to a web page listing faculty")
}
