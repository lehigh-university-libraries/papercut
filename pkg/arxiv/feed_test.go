package arxiv_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lehigh-university-libraries/papercut/pkg/arxiv"
)

func TestGetResults(t *testing.T) {
	// Mocking HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a sample Atom feed XML
		xml := `<feed xmlns="http://www.w3.org/2005/Atom" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/">
			<link/>
			<title>Test Feed</title>
			<id>urn:uuid:60a76c80-74b1-11eb-9439-0242ac130002</id>
			<updated>2023-10-30T10:00:00Z</updated>
			<opensearch:totalResults>2</opensearch:totalResults>
			<opensearch:startIndex>1</opensearch:startIndex>
			<opensearch:itemsPerPage>10</opensearch:itemsPerPage>
			<entry></entry>
			<entry></entry>
		</feed>`
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, xml)
	}))
	defer ts.Close()

	// Test case
	url := ts.URL
	feed, err := arxiv.GetResults(url)
	if err != nil {
		t.Errorf("GetResults(%s) returned error: %v", url, err)
	}

	// Verify the result
	if feed.Title != "Test Feed" {
		t.Errorf("Expected feed title 'Test Feed', got '%s'", feed.Title)
	}
	if feed.TotalResults != 2 {
		t.Errorf("Expected totalResults to be 2, got %d", feed.TotalResults)
	}
}
