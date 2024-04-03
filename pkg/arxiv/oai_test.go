package arxiv_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lehigh-university-libraries/papercut/pkg/arxiv"
)

func TestGetOaiRecord(t *testing.T) {
	// Mocking HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a sample OAI response XML
		xml := `<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
		<responseDate>2024-03-21T13:07:35Z</responseDate>
		<request verb="GetRecord" identifier="oai:arXiv.org:2210.04727" metadataPrefix="arXiv">http://export.arxiv.org/oai2</request>
		<GetRecord>
		  <record>
		    <header>
					<identifier>oai:arXiv.org:123.abc</identifier>
					<datestamp>2023-02-22</datestamp>
					<setSpec>123</setSpec>
				</header>
				<metadata>
					<arXiv xmlns="http://arxiv.org/OAI/arXiv/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://arxiv.org/OAI/arXiv/ http://arxiv.org/OAI/arXiv.xsd">
							<id>123.abc</id>
							<license>https://test-license.com</license>
							<authors>
								<author>
									<keyname>Author1LastName</keyname>
									<forenames>Author1FirstName</forenames>
								</author>
								<author>
									<keyname>Author2LastName</keyname>
									<forenames>Author2FirstName</forenames>
								</author>
							</authors>
						</arXiv>
					</metadata>
				</record>
			</GetRecord>
		</OAI-PMH>`
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, xml)
	}))
	defer ts.Close()

	result := arxiv.GetOaiRecord(ts.URL)

	// Verify the result
	expected := map[string]string{
		"field_rights":       "https://test-license.com",
		"field_linked_agent": "relators:cre:person:Author1LastName, Author1FirstName|relators:cre:person:Author2LastName, Author2FirstName",
	}
	for key, value := range expected {
		if result[key] != value {
			t.Errorf("Expected value for %s to be %s, got %s", key, value, result[key])
		}
	}
}
