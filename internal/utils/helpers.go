package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func FetchEmails(url string) ([]string, error) {
	queries := []string{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error requesting directory listing:", err)
		return queries, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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
