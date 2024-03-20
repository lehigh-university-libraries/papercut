package arxiv

import (
	"encoding/xml"
	"regexp"
	"strings"
	"time"
)

type Entry struct {
	ID              string     `xml:"id"`
	Updated         time.Time  `xml:"updated"`
	Published       time.Time  `xml:"published"`
	Title           string     `xml:"title"`
	Summary         string     `xml:"summary"`
	Authors         []Author   `xml:"author"`
	DOI             string     `xml:"http://arxiv.org/schemas/atom doi"`
	Links           []Link     `xml:"link"`
	Comment         string     `xml:"http://arxiv.org/schemas/atom comment"`
	JournalRef      string     `xml:"http://arxiv.org/schemas/atom journal_ref"`
	PrimaryCategory Category   `xml:"http://arxiv.org/schemas/atom primary_category"`
	Categories      []Category `xml:"category"`
	License         string     `xml:"license"`
	PDF             string
}

// Implement the UnmarshalXML interface to cleanup title fields
func (e *Entry) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type entryAlias Entry
	var alias entryAlias

	err := d.DecodeElement(&alias, &start)
	if err != nil {
		return err
	}

	*e = Entry(alias)
	e.Title = cleanString(e.Title)

	for _, link := range e.Links {
		if link.Title == "pdf" {
			e.PDF = link.Href
		}
	}

	return nil
}

// Helper function to clean string (remove new lines and trim whitespace)
func cleanString(s string) string {
	cs := strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
	cs = strings.ReplaceAll(cs, "\t", " ")

	// Replace double spaces with a single space
	re := regexp.MustCompile(`\s{2,}`)
	return re.ReplaceAllString(cs, " ")
}
