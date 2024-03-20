package arxiv

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// OAIResponse represents the XML structure of the OAI response
type OAIResponse struct {
	XMLName xml.Name `xml:"OAI-PMH"`
	Record  Record   `xml:"GetRecord>record"`
}

// Record represents the XML structure of the record element
type Record struct {
	XMLName xml.Name `xml:"record"`
	ID      string   `xml:"metadata>arXiv>id"`
	License string   `xml:"metadata>arXiv>license"`
	Authors Authors  `xml:"metadata>arXiv>authors"`
}

// Authors represents the XML structure of the authors element
type Authors struct {
	Authors []OaiAuthor `xml:"author"`
}

// Author represents the XML structure of the author element
type OaiAuthor struct {
	KeyName  string `xml:"keyname"`
	ForeName string `xml:"forenames"`
}

func GetOaiRecord(id string) map[string]string {
	values := map[string]string{
		"field_rights": "https://arxiv.org/licenses/nonexclusive-distrib/1.0/license.html",
	}
	log.Println("Fetching", id)
	url := fmt.Sprintf("https://export.arxiv.org/oai2?verb=GetRecord&identifier=oai:arXiv.org:%s&metadataPrefix=arXiv", id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	var oaiResponse OAIResponse
	err = xml.Unmarshal(body, &oaiResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	if oaiResponse.Record.License != "" {
		values["field_rights"] = oaiResponse.Record.License
	}

	var authors []string
	for _, author := range oaiResponse.Record.Authors.Authors {
		authors = append(authors, fmt.Sprintf("%s, %s", author.KeyName, author.ForeName))
	}
	values["field_linked_agent"] = strings.Join(authors, "|")

	return values
}
