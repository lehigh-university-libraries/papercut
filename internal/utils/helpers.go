package utils

import (
	"fmt"
	"io"
	"net/http"
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
	// Check if the string length exceeds the maximum length
	if utf8.RuneCountInString(s) > maxLen {
		// Convert the string to a slice of runes
		runes := []rune(s)

		// Truncate the slice to the maximum length
		runes = runes[:maxLen]

		// Convert the slice of runes back to a string
		return string(runes)
	}
	return s
}
