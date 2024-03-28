package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lehigh-university-libraries/papercut/internal/utils"
	"github.com/lehigh-university-libraries/papercut/pkg/doi"
	"github.com/spf13/cobra"
)

var (
	// used for flags.
	filePath string

	doiCmd = &cobra.Command{
		Use:   "doi",
		Short: "Get DOI metadata and PDF",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(filePath)
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
				"field_edtf_date_issued",
				"title",
				"field_full_title",
				"field_abstract",
				"field_linked_agent",
				"field_identifier",
				"field_language",
				"field_rights",
				"field_subject",
				"file",
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
				if _, err := os.Stat(d); err == nil {
					content, err := os.ReadFile(d)
					if err != nil {
						fmt.Println("Error reading file:", err)
						return
					}
					err = json.Unmarshal(content, &doiObject)
					if err != nil {
						log.Printf("Unable to unmarshal cached file: %v", err)
						continue
					}
				} else {
					apiURL := fmt.Sprintf("%s/%s", url, line)

					log.Printf("Accessing %s\n", apiURL)

					doiObject, err = doi.GetResults(apiURL)
					if err != nil {
						log.Fatal(err)
					}
					if err != nil {
						log.Printf("Unable to create directory to cache DOI: %v", err)
						continue
					}
					jsonData, err := json.Marshal(doiObject)
					if err != nil {
						continue
					}

					cacheFile, err := os.Create(d)
					if err != nil {
						fmt.Println("Error creating file:", err)
						return
					}
					defer cacheFile.Close()

					_, err = cacheFile.WriteString(string(jsonData))
					if err != nil {
						fmt.Println("Error caching DOI JSON:", err)
					}
					time.Sleep(500 * time.Microsecond)
				}

				var linkedAgent []string
				for _, author := range doiObject.Authors {
					linkedAgent = append(linkedAgent, fmt.Sprintf("relators:aut:person:%s, %s", author.Family, author.Given))
				}
				if doiObject.Publisher != "" {
					linkedAgent = append(linkedAgent, fmt.Sprintf("relators:pbl:corporate_body:%s", doiObject.Publisher))
				}
				identifiers := []string{
					fmt.Sprintf(`{"attr0":"doi","value":"%s"}`, doiObject.DOI),
				}
				for _, i := range doiObject.ISSN {
					identifiers = append(identifiers, fmt.Sprintf(`{"attr0":"issn","value":"%s"}`, i))
				}
				err = wr.Write([]string{
					line,
					doi.JoinDate(doiObject.Issued),
					utils.TrimToMaxLen(doiObject.Title, 255),
					doiObject.Title,
					doiObject.Abstract,
					strings.Join(linkedAgent, "|"),
					strings.Join(identifiers, "|"),
					doiObject.Language,
					"",
					strings.Join(doiObject.Subject, "|"),
					"file",
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
			/*
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
			*/
		},
	}
)

func init() {
	getCmd.AddCommand(doiCmd)

	doiCmd.Flags().StringP("url", "u", "https://dx.doi.org", "The DOI API url")
	doiCmd.Flags().StringVarP(&filePath, "file", "f", "", "path to file containing one DOI per line")
}
