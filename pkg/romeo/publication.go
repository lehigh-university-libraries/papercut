package romeo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lehigh-university-libraries/papercut/internal/utils"
)

type Response struct {
	Publications []Publication `json:"items"`
}

type Publication struct {
	PublisherPolicies []PublisherPolicy `json:"publisher_policy"`
}

type PublisherPolicy struct {
	Uri                  string       `json:"uri"`
	OpenAccessProhibited string       `json:"open_access_prohibited"`
	PermittedOa          []OpenAccess `json:"permitted_oa"`
}

type OpenAccess struct {
	ArticleVersion []string  `json:"article_version"`
	Conditions     []string  `json:"conditions"`
	Embargo        Embargo   `json:"embargo,omitempty"`
	License        []License `json:"license,omitempty"`
	Location       Location  `json:"location"`
	AdditonalFee   string    `json:"additional_oa_fee"`
}

type Location struct {
	Locations []string `json:"location"`
}

type Embargo struct {
	Amount int    `json:"amount,omitempty"`
	Units  string `json:"units,omitempty"`
}

type License struct {
	Value   string `json:"license"`
	Version string `json:"version"`
}

func GetIdFromIssn(i string) string {
	url := fmt.Sprintf("https://v2.sherpa.ac.uk//cgi/romeosearch?publication_title-auto=%s", i)
	// Create a custom HTTP client with redirection disabled
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevent automatic redirection
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Check if the response status code is 301 (Moved Permanently)
	if resp.StatusCode == http.StatusFound {
		location := strings.Split(resp.Header.Get("Location"), "/")
		return location[len(location)-1]
	}
	log.Println(resp.StatusCode)

	return ""
}

func GetPublication(url string) []byte {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/pdf")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error downloading PDF:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("Error: HTTP status %d\n", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return body
}

func (r *Response) GetLicense() string {
	license := ""
	for _, p := range r.Publications {
		for _, policy := range p.PublisherPolicies {
			license = policy.Uri
			for _, oa := range policy.PermittedOa {
				if utils.StrInSlice("published", oa.ArticleVersion) {
					if !utils.StrInSlice("any_website", oa.Location.Locations) && !utils.StrInSlice("non_commercial_website", oa.Location.Locations) && !utils.StrInSlice("institutional_repository", oa.Location.Locations) && !utils.StrInSlice("non_commercial_repository", oa.Location.Locations) {
						continue
					}

					if oa.Embargo.Amount == 0 {
						for _, l := range oa.License {
							uri := l.Uri()
							if uri != "" {
								return uri
							}
						}
					}
				}
			}
		}
	}

	return license
}

func (l License) Uri() string {
	c := strings.Split(l.Value, "_")
	if c[0] == "cc" {

		if l.Version == "" {
			l.Version = "4.0"
		}
		uri := strings.Join(c[1:], "-")
		if uri == "public-domain" {
			uri = "publicdomain/"
		} else {
			uri = fmt.Sprintf("licenses/%s/%s/", uri, l.Version)
		}
		return fmt.Sprintf("https://creativecommons.org/%s", uri)
	}

	return ""
}
