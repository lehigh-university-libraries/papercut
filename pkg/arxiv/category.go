package arxiv

type Category struct {
	Term   string `xml:"term,attr"`
	Scheme string `xml:"scheme,attr"`
}
