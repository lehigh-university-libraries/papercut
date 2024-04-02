package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lehigh-university-libraries/papercut/pkg/doi"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "papercut",
	Short: "Command line utility to help fetch papers from various sources.	",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getResult(d, url, line, acceptContentType string) []byte {
	var err error
	content := checkCachedFile(d)
	if content != nil {
		var a doi.Affiliation
		err = json.Unmarshal(content, &a)
		if err == nil || acceptContentType == "text/html" {
			return content
		}
		log.Println("Error unmarshalling cached file:", err)
	}

	apiURL := fmt.Sprintf("%s/%s", url, line)

	log.Printf("Accessing %s\n", apiURL)

	doiObject, err := doi.GetObject(apiURL, acceptContentType)
	if err != nil {
		log.Fatal(err)
	}
	writeCachedFile(d, string(doiObject))
	return doiObject
}

func checkCachedFile(d string) []byte {
	// see if we can just get the cached file
	if _, err := os.Stat(d); err == nil {
		content, err := os.ReadFile(d)
		if err != nil {
			log.Println("Error reading cached file:", err)
			return nil
		}
		return content
	}
	return nil
}

func writeCachedFile(f, c string) {
	cacheFile, err := os.Create(f)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer cacheFile.Close()

	_, err = cacheFile.WriteString(c)
	if err != nil {
		log.Println("Error caching DOI JSON:", err)
	}
}
