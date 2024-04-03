package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lehigh-university-libraries/papercut/internal/utils"
	"github.com/lehigh-university-libraries/papercut/pkg/doi"
	"github.com/lehigh-university-libraries/papercut/pkg/romeo"
	"github.com/spf13/cobra"
)

var (
	// used for flags.
	filePath     string
	downloadPdfs bool
	doiCmd       = &cobra.Command{
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
				"field_part_detail",
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
				doiStr := strings.TrimSpace(scanner.Text())
				doiObject, err := doi.GetDoi(doiStr, url)
				if err != nil {
					log.Println(err)
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
				fieldRights := ""
				for _, i := range doiObject.ISSN {
					identifiers = append(identifiers, fmt.Sprintf(`{"attr0":"issn","value":"%s"}`, i))
					if fieldRights == "" {
						fieldRights = romeo.FindIssnLicense(i)
					}
				}

				partDetail := []string{}
				if doiObject.Volume != "" {
					partDetail = append(partDetail, fmt.Sprintf(`{"type": "volume", "number": "%s"}`, doiObject.Volume))
				}
				if doiObject.Issue != "" {
					partDetail = append(partDetail, fmt.Sprintf(`{"type": "volume", "number": "%s"}`, doiObject.Issue))
				}

				relatedItem := []string{}
				if doiObject.ContainerTitle != "" {
					relatedItem = append(relatedItem, fmt.Sprintf(`{"title": "%s"}`, doiObject.ContainerTitle))
				}
				extent := ""
				if doiObject.Page != "" {
					extent = fmt.Sprintf(`{"attr0": "page", "number": "%s"}`, doiObject.Page)
				}

				pdf := ""
				if downloadPdfs {
					pdf = doiObject.DownloadPdf()
				}

				fullTitle := ""
				if len(doiObject.Title) > 255 {
					fullTitle = doiObject.Title
				}
				err = wr.Write([]string{
					doiStr,
					doi.JoinDate(doiObject.Issued),
					utils.TrimToMaxLen(doiObject.Title, 255),
					fullTitle,
					doiObject.Abstract,
					"Digital Document",
					strings.Join(linkedAgent, "|"),
					strings.Join(identifiers, "|"),
					strings.Join(partDetail, "|"),
					strings.Join(relatedItem, "|"),
					extent,
					doiObject.Language,
					fieldRights,
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
	doiCmd.Flags().BoolVarP(&downloadPdfs, "download-pdfs", "d", true, "whether to download the PDFs")
}
