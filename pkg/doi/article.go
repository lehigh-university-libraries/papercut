package doi

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lehigh-university-libraries/papercut/internal/utils"
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

func GetDoi(d, url string) (Article, error) {
	var a Article
	var err error
	dirPath := filepath.Join("dois", d)
	dirPath, err = utils.MkTmpDir(dirPath)
	if err != nil {
		return Article{}, fmt.Errorf("unable to create cached file directory: %v", err)
	}

	dir := filepath.Join(dirPath, "doi.json")
	u := fmt.Sprintf("%s/%s", url, d)
	result := utils.GetResult(dir, u, "application/json")
	if result == nil {
		return Article{}, fmt.Errorf("could not find DOI %s", d)
	}

	err = json.Unmarshal(result, &a)
	if err != nil {
		return Article{}, fmt.Errorf("could not unmarshal JSON for %s: %v", d, err)
	}
	return a, nil
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
		return "invalid date"
	}

	return parsedTime.Format(targetPattern)
}
