package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lehigh-university-libraries/papercut/pkg/doi"
	"github.com/lehigh-university-libraries/papercut/pkg/romeo"
	"github.com/spf13/cobra"
)

var (
	// used for flags.
	licenseFilePath string
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
				doiStr := strings.TrimSpace(scanner.Text())
				doiObject, err := doi.GetDoi(doiStr, url)
				if err != nil {
					log.Println(err)
					continue
				}

				fieldRights := ""
				for _, i := range doiObject.ISSN {
					fieldRights = romeo.FindIssnLicense(i)
					if fieldRights != "" {
						break
					}
				}

				err = wr.Write([]string{
					doiStr,
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
