package arxiv

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
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

func GetResults(url string) (Feed, error) {
	var result Feed

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error requesting XML data:", err)
		return result, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result, err
	}

	err = xml.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return result, err
	}

	return result, nil
}
