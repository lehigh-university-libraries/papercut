package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"unicode/utf8"
)

func FetchEmails(url string) ([]string, error) {
	queries := []string{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error requesting directory listing:", err)
		return queries, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return queries, err
	}

	// Regular expression pattern to match email addresses
	emailRegex := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`
	r := regexp.MustCompile(emailRegex)
	emails := r.FindAllString(string(body), -1)

	return emails, nil
}

func TrimToMaxLen(s string, maxLen int) string {
	if utf8.RuneCountInString(s) > maxLen {
		runes := []rune(s)
		runes = runes[:maxLen]
		return string(runes)
	}

	return s
}

func MkTmpDir(d string) (string, error) {
	tmpDir := os.TempDir()
	dirPath := filepath.Join(tmpDir, d)
	_, err := os.Stat(dirPath)
	if err == nil {
		return dirPath, nil
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return "", err
		}
	}

	return dirPath, nil
}

func DownloadPdf(url, filePath string) error {
	downloadDirectory := filepath.Dir(filePath)
	if err := os.MkdirAll(downloadDirectory, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("OK")
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating request:", err)
			return err
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/pdf")
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Cache-Control", "no-cache")

		response, err := client.Do(req)
		if err != nil {
			log.Println("Error downloading PDF:", err)
			return err
		}
		defer response.Body.Close()

		if response.StatusCode > 299 {
			log.Printf("Error: HTTP status %d\n", response.StatusCode)
			return fmt.Errorf("%s returned a non-200 status code: %d", url, response.StatusCode)
		}
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Println("Error copying PDF content to file:", err)
			return err
		}
	}

	time.Sleep(500 * time.Microsecond)

	return nil
}

func StrInSlice(s string, sl []string) bool {
	for _, a := range sl {
		if a == s {
			return true
		}
	}
	return false
}
