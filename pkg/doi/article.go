package doi

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Affiliation struct {
	Name string `json:"name"`
}

type Author struct {
	Given       string        `json:"given"`
	Family      string        `json:"family"`
	Sequence    string        `json:"sequence"`
	Affiliation []Affiliation `json:"affiliation"`
}

type ContentDomain struct {
	Domain               []string `json:"domain"`
	CrossmarkRestriction bool     `json:"crossmark-restriction"`
}

type Link struct {
	URL                 string `json:"URL"`
	ContentType         string `json:"content-type"`
	ContentVersion      string `json:"content-version"`
	IntendedApplication string `json:"intended-application"`
}

type Resource struct {
	Primary struct {
		URL string `json:"URL"`
	} `json:"primary"`
}

type JournalIssue struct {
	Issue           string    `json:"issue"`
	PublishedOnline DateParts `json:"published-online"`
	PublishedPrint  DateParts `json:"published-print"`
}

type DateParts struct {
	Dates [][]int `json:"date-parts"`
}

// Define the main struct
type Article struct {
	ReferenceCount      int           `json:"reference-count"`
	Publisher           string        `json:"publisher"`
	Issue               string        `json:"issue"`
	ContentDomain       ContentDomain `json:"content-domain"`
	Abstract            string        `json:"abstract"`
	DOI                 string        `json:"DOI"`
	Type                string        `json:"type"`
	Created             DateParts     `json:"created"`
	Page                string        `json:"page"`
	Source              string        `json:"source"`
	IsReferencedByCount int           `json:"is-referenced-by-count"`
	PublishedPrint      DateParts     `json:"published-print"`
	Title               string        `json:"title"`
	Prefix              string        `json:"prefix"`
	Volume              string        `json:"volume"`
	Authors             []Author      `json:"author"`
	Member              string        `json:"member"`
	PublishedOnline     DateParts     `json:"published-online"`
	ContainerTitle      string        `json:"container-title"`
	OriginalTitle       []string      `json:"original-title"`
	Language            string        `json:"language"`
	Link                []Link        `json:"link"`
	Deposited           DateParts     `json:"deposited"`
	Score               int           `json:"score"`
	Resource            Resource      `json:"resource"`
	Subtitle            []string      `json:"subtitle"`
	ShortTitle          []string      `json:"short-title"`
	Issued              DateParts     `json:"issued"`
	ReferencesCount     int           `json:"references-count"`
	JournalIssue        JournalIssue  `json:"journal-issue"`
	URL                 string        `json:"URL"`
	ISSN                []string      `json:"ISSN"`
	Subject             []string      `json:"subject"`
	ContainerTitleShort string        `json:"container-title-short"`
	PublishedDate       DateParts     `json:"published"`
	// do not have a need for relation ATM
	//	Relation            interface{}   `json:"relation"`
}

func GetObject(url, acceptContentType string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", acceptContentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func JoinDate(d DateParts) string {
	l := len(d.Dates[0])

	// Convert each integer to a string
	strNumbers := make([]string, l)
	for i, num := range d.Dates[0] {
		strNumbers[i] = strconv.Itoa(num)
	}

	dateString := strings.Join(strNumbers, "-")

	// if it's just a year, return it
	if l == 1 {
		return dateString
	}

	// else we need to convert to EDTF dates
	sourcePattern := "2006-1-2"
	targetPattern := "2006-01-02"
	if l == 2 {
		sourcePattern = "2006-1"
		targetPattern = "2006-01"
	}
	parsedTime, err := time.Parse(sourcePattern, dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return dateString
	}
	return parsedTime.Format(targetPattern)
}
