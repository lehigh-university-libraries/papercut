package arxiv

type Author struct {
	Name        string `xml:"name"`
	Affiliation string `xml:"http://arxiv.org/schemas/atom affiliation"`
}
