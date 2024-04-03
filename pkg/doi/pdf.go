package doi

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lehigh-university-libraries/papercut/internal/utils"
)

func (d *Article) DownloadPdf() string {
	pdfUrl := ""
	pdf := ""
	for _, l := range d.Link {
		if l.ContentType == "application/pdf" || strings.Contains(strings.ToLower(l.URL), "pdf") {
			pdfUrl = l.URL
		}
	}
	if pdfUrl == "" {
		dirPath, err := utils.MkTmpDir(filepath.Join("dois", d.DOI))
		if err != nil {
			log.Fatal("Unable to write to tmp filesystem")
		}
		dir := filepath.Join(dirPath, "doi.html")
		log.Println(dir)
		result := utils.GetResult(dir, d.URL, "text/html")
		pattern := `<meta name="citation_pdf_url" content="([^"]+)".*>`
		re := regexp.MustCompile(pattern)
		matches := re.FindAllSubmatch(result, -1)
		var pdfURLs []string
		for _, match := range matches {
			if len(match) >= 2 {
				pdfURLs = append(pdfURLs, string(match[1]))
			}
		}
		for _, url := range pdfURLs {
			pdfUrl = url
			break
		}
	}

	if pdfUrl != "" {
		hash := md5.Sum([]byte(d.DOI))
		hashStr := hex.EncodeToString(hash[:])

		pdf = fmt.Sprintf("papers/dois/%s.pdf", hashStr)
		err := utils.DownloadPdf(pdfUrl, pdf)
		if err != nil {
			err = os.Remove(pdf)
			if err != nil {
				log.Println("Error deleting file:", err)
				return ""
			}
			pdf = pdfUrl
		}
	}

	return pdf
}
