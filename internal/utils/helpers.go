package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
