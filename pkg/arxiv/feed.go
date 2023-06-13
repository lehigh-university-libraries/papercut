package arxiv

import (
	"encoding/xml"
)

type Feed struct {
	XMLName      xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Link         Link     `xml:"link"`
	Title        string   `xml:"title"`
	ID           string   `xml:"id"`
	Updated      string   `xml:"updated"`
	TotalResults int      `xml:"http://a9.com/-/spec/opensearch/1.1/ totalResults"`
	StartIndex   int      `xml:"http://a9.com/-/spec/opensearch/1.1/ startIndex"`
	ItemsPerPage int      `xml:"http://a9.com/-/spec/opensearch/1.1/ itemsPerPage"`
	Entries      []Entry  `xml:"entry"`
}
