package cmd

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
				"field_model",
				"field_linked_agent",
				"field_identifier",
				"field_related_item",
				"field_extent",
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
				result := getResult(d, url, line, "application/json")
				err = json.Unmarshal(result, &doiObject)
				if err != nil {
					log.Printf("Could not unmarshal JSON for %s: %v", line, err)
					continue
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

				relatedItem := []string{}
				if doiObject.Volume != "" {
					relatedItem = append(relatedItem, fmt.Sprintf(`{"type": "volume", "number": "%s"}`, doiObject.Volume))
				}
				if doiObject.Issue != "" {
					relatedItem = append(relatedItem, fmt.Sprintf(`{"type": "volume", "number": "%s"}`, doiObject.Issue))
				}

				extent := ""
				if doiObject.Page != "" {
					extent = fmt.Sprintf(`{"attr0": "page", "number": "%s"}`, doiObject.Page)
				}

				var pdfUrl string
				pdf := ""
				for _, l := range doiObject.Link {
					if l.ContentType == "application/pdf" || strings.Contains(strings.ToLower(l.URL), "pdf") {
						pdfUrl = l.URL
					}
				}
				if pdfUrl == "" {
					d = filepath.Join(dirPath, "doi.html")
					result = getResult(d, url, line, "text/html")
					pattern := `<meta\s+name="citation_pdf_url"\s+content="([^"]+)"\s*>`
					re := regexp.MustCompile(pattern)
					matches := re.FindAllSubmatch(result, -1)
					var pdfURLs []string
					for _, match := range matches {
						if len(match) >= 2 {
							pdfURLs = append(pdfURLs, string(match[1]))
						}
					}
					for _, url := range pdfURLs {
						pdfUrl = url
						break
					}
				}
				if pdfUrl != "" {
					hash := md5.Sum([]byte(line))
					hashStr := hex.EncodeToString(hash[:])

					pdf = fmt.Sprintf("papers/dois/%s.pdf", hashStr)
					err = utils.DownloadPdf(pdfUrl, pdf)
					if err != nil {
						err = os.Remove(pdf)
						if err != nil {
							log.Println("Error deleting file:", err)
						}
						pdf = pdfUrl
					}
				}

				fullTitle := ""
				if len(doiObject.Title) > 255 {
					fullTitle = doiObject.Title
				}
				err = wr.Write([]string{
					line,
					doi.JoinDate(doiObject.Issued),
					utils.TrimToMaxLen(doiObject.Title, 255),
					fullTitle,
					doiObject.Abstract,
					"Digital Document",
					strings.Join(linkedAgent, "|"),
					strings.Join(identifiers, "|"),
					strings.Join(relatedItem, "|"),
					extent,
					doiObject.Language,
					"",
					strings.Join(doiObject.Subject, "|"),
					pdf,
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
	getCmd.AddCommand(doiCmd)

	doiCmd.Flags().StringP("url", "u", "https://dx.doi.org", "The DOI API url")
	doiCmd.Flags().StringVarP(&filePath, "file", "f", "", "path to file containing one DOI per line")
}

func getResult(d, url, line, acceptContentType string) []byte {
	var err error

	// see if we can just get the cached file
	if _, err := os.Stat(d); err == nil {
		content, err := os.ReadFile(d)
		if err != nil {
			fmt.Println("Error reading cached file:", err)
		} else {
			var a doi.Affiliation
			err = json.Unmarshal(content, &a)
			if err == nil || acceptContentType == "text/html" {
				return content
			}
			log.Println("Error unmarshalling cached file:", err)
		}
	}

	apiURL := fmt.Sprintf("%s/%s", url, line)

	log.Printf("Accessing %s\n", apiURL)

	doiObject, err := doi.GetObject(apiURL, acceptContentType)
	if err != nil {
		log.Fatal(err)
	}
	cacheFile, err := os.Create(d)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}
	defer cacheFile.Close()

	_, err = cacheFile.WriteString(string(doiObject))
	if err != nil {
		fmt.Println("Error caching DOI JSON:", err)
	}

	return doiObject
}
