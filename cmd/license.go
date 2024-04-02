package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/lehigh-university-libraries/papercut/internal/utils"
	"github.com/lehigh-university-libraries/papercut/pkg/doi"
	"github.com/lehigh-university-libraries/papercut/pkg/romeo"
	"github.com/spf13/cobra"
)

var (
	// used for flags.
	licenseFilePath string
	romeoApiKey     = os.Getenv("SHERPA_ROMEO_API_KEY")
	licenseCmd      = &cobra.Command{
		Use:   "license",
		Short: "Get license for a DOI",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(licenseFilePath)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()

			// Create a scanner to read the file line by line
			scanner := bufio.NewScanner(file)
			url, err := cmd.Flags().GetString("url")
			if err != nil {
				log.Fatal(err)
			}
			wr := csv.NewWriter(os.Stdout)

			// CSV header
			err = wr.Write([]string{
				"id",
				"field_rights",
			})
			if err != nil {
				log.Fatalf("Unable to write to CSV: %v", err)
			}
			for scanner.Scan() {
				var doiObject doi.Article
				line := strings.TrimSpace(scanner.Text())
				dirPath := filepath.Join("dois", line)
				dirPath, err = utils.MkTmpDir(dirPath)
				if err != nil {
					log.Printf("Unable to create cached file directory: %v", err)
					continue
				}

				d := filepath.Join(dirPath, "doi.json")
				result := getResult(d, url, line, "application/json")
				err = json.Unmarshal(result, &doiObject)
				if err != nil {
					log.Printf("Could not unmarshal JSON for %s: %v", line, err)
					continue
				}

				fieldRights := ""
				for _, i := range doiObject.ISSN {
					d, err = utils.MkTmpDir("issns")
					if err != nil {
						continue
					}
					d = filepath.Join(d, i)
					publicationId := checkCachedFile(d)
					id := string(publicationId)
					if publicationId == nil {
						id = romeo.GetIdFromIssn(i)
						if id != "" {
							writeCachedFile(d, id)
						}
					}
					if id == "" {
						log.Println("Could not find publication ID for ISSN", i)
						continue
					}
					filter := fmt.Sprintf("[[\"id\",\"equals\",\"%s\"]]", id)
					romeUrl := fmt.Sprintf("https://v2.sherpa.ac.uk/cgi/retrieve?item-type=publication&format=Json&limit=10&offset=0&order=-id&filter=%s&api-key=%s", neturl.QueryEscape(filter), romeoApiKey)
					d, _ = utils.MkTmpDir(filepath.Join("issns", "ids"))
					d = filepath.Join(d, id)
					publication := checkCachedFile(d)
					if publication == nil {
						publication = romeo.GetPublication(romeUrl)
						if publication != nil {
							writeCachedFile(d, string(publication))
						}
					}
					if publication == nil {
						log.Println("Could not find publication info for", i)
						continue
					}
					var r romeo.Response
					err = json.Unmarshal(publication, &r)
					if err != nil {
						log.Printf("Unable to read publication: %v", err)
						continue
					}

					fieldRights = r.GetLicense()
					if fieldRights != "" {
						break
					}
				}

				err = wr.Write([]string{
					line,
					fieldRights,
				})
				if err != nil {
					log.Fatalf("Unable to write to CSV: %v", err)
				}
				wr.Flush()
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error scanning file:", err)
				return
			}
		},
	}
)

func init() {
	getCmd.AddCommand(licenseCmd)

	licenseCmd.Flags().StringP("url", "u", "https://dx.doi.org", "The DOI API url")
	licenseCmd.Flags().StringVarP(&licenseFilePath, "file", "f", "", "path to file containing one DOI per line")
}
