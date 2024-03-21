package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
			} else {
				emails, err := cmd.Flags().GetString("emails")
				if err != nil {
					log.Fatal(err)
				}
				if emails != "" {
					if query != "" || ids != "" {
						log.Fatal("query or ids can not be used with the email option.")
					}
					slice := strings.Split(emails, "\n")
					for _, s := range slice {
						s = strings.TrimSpace(s)
						if s != "" && strings.Contains(s, "@") {
							queries = append(queries, s)
						}
					}
				}
			}
			if query != "" {
				queries = append(queries, query)
			}
			if ids != "" {
				s := strings.Split(",", ids)
				id_list := []string{}
				for _, id := range s {
					id_list = append(id_list, id)
					if len(id_list) == results {
						queries = append(queries, strings.Join(id_list, ","))
						id_list = []string{}
					}
				}
				if len(id_list) > 0 {
					queries = append(queries, strings.Join(id_list, ","))
				}
			}
			if len(queries) == 0 && ids == "" {
				log.Fatal("query or ids required.")
			}

			wr := csv.NewWriter(os.Stdout)

			// header
			wr.Write([]string{
				"id",
				"field_edtf_date_issued",
				"title",
				"field_full_title",
				"field_abstract",
				"field_linked_agent",
				"field_identifier",
				"field_related_item",
				"field_rights",
				"field_subject",
				"file",
				"arXiv search query",
			})
			categoryNames := arxiv.GetCategoryLabels()
			for _, query := range queries {

				params := url.Values{}
				if ids != "" {
					params.Set("id_list", ids)

				} else {
					params.Set("search_query", query)
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
				pattern := `/abs/([0-9a-z\-]+(\/|\.)\d+)(?:v\d+)?$`
				re := regexp.MustCompile(pattern)

				for true {
					for _, e := range result.Entries {

						log.Println("Pausing between requests. arXiv requests a three second delay between API requests...")
						time.Sleep(3 * time.Second)

						matches := re.FindStringSubmatch(e.ID)
						if len(matches) <= 1 {
							log.Fatal(e.ID)
						}
						oai := arxiv.GetOaiRecord(matches[1])
						if e.JournalRef != "" {
							e.JournalRef = fmt.Sprintf(`{"title": "%s"}`, e.JournalRef)
						}
						var categories = []string{}
						for _, c := range e.Categories {
							term := strings.Split(c.Term, ".")
							group := "Physics"
							switch term[0] {
							case "cs":
								group = "Computer Science"
							case "econ":
								group = "Economics"
							case "eess":
								group = "Electrical Engineering and Systems Science"
							case "math":
								group = "Mathematics"
							case "astro-ph":
								group = "Physics--Astrophysics"
							case "cond-mat":
								group = "Physics--Condensed Matter"
							case "nlin":
								group = "Physics--Nonliner Sciences"
							case "q-bio":
								group = "Quantitative Biology"
							case "q-fin":
								group = "Quantitative Finance"
							case "stat":
								group = "Statistics"
							}
							if categoryName, ok := categoryNames[c.Term]; ok {
								c.Term = fmt.Sprintf("%s--%s", group, categoryName)
							}
							categories = append(categories, c.Term)
						}
						var identifiers = []string{
							fmt.Sprintf(`{"attr0":"arxiv","value":"%s"}`, e.ID),
						}
						if e.DOI != "" {
							doi := fmt.Sprintf(`{"attr0":"doi","value":"%s"}`, e.DOI)
							identifiers = append(identifiers, doi)
						}
						wr.Write([]string{
							e.ID,
							strings.Split(e.Published.String(), " ")[0],
							utils.TrimToMaxLen(e.Title, 255),
							e.Title,
							e.Summary,
							oai["field_linked_agent"],
							strings.Join(identifiers, "|"),
							e.JournalRef,
							oai["field_rights"],
							strings.Join(categories, "|"),
							e.PDF,
							query,
						})
						wr.Flush()

						if e.PDF != "" {
							downloadDirectory := "papers"
							if err := os.MkdirAll(downloadDirectory, 0755); err != nil {
								fmt.Println("Error creating directory:", err)
								return
							}

							_, filename := filepath.Split(e.PDF)
							// Ensure the filename has a .pdf extension
							if !strings.HasSuffix(filename, ".pdf") {
								filename = fmt.Sprintf("%s.pdf", filename)
							}
							filePath := filepath.Join(downloadDirectory, filename)

							if _, err := os.Stat(filePath); os.IsNotExist(err) {

								file, err := os.Create(filePath)
								if err != nil {
									fmt.Println("Error creating file:", err)
									return
								}
								defer file.Close()

								response, err := http.Get(e.PDF)
								if err != nil {
									fmt.Println("Error downloading PDF:", err)
									return
								}
								defer response.Body.Close()

								if response.StatusCode != http.StatusOK {
									fmt.Printf("Error: HTTP status %d\n", response.StatusCode)
									return
								}
								_, err = io.Copy(file, response.Body)
								if err != nil {
									fmt.Println("Error copying PDF content to file:", err)
									return
								}
							}
						}
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
	arxivCmd.Flags().StringVarP(&query, "query", "q", "", "The arXiv API search query to perform")
	arxivCmd.Flags().StringVarP(&ids, "ids", "i", "", "A comma separated list of arXiv IDs")
	arxivCmd.Flags().IntVarP(&start, "start", "s", 0, "The offset")
	arxivCmd.Flags().IntVarP(&results, "results", "r", 10, "The number of results to return in a response")
	arxivCmd.Flags().String("directory-listing", "", "URL to a web page listing faculty email addresses")
	arxivCmd.Flags().String("emails", "", "List of emails to search for")
}
