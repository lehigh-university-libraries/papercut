package romeo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Publication struct {
}

func GetIdFromIssn(i string) string {
	url := fmt.Sprintf("https://v2.sherpa.ac.uk//cgi/romeosearch?publication_title-auto=%s", i)
	// Create a custom HTTP client with redirection disabled
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevent automatic redirection
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Check if the response status code is 301 (Moved Permanently)
	if resp.StatusCode == http.StatusFound {
		location := strings.Split(resp.Header.Get("Location"), "/")
		return location[len(location)-1]
	}
	log.Println(resp.StatusCode)

	return ""
}

func GetPublication(url string) []byte {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/pdf")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error downloading PDF:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("Error: HTTP status %d\n", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return body
}
