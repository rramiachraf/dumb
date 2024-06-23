package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CleanBody(body string) string {
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(body)); err == nil {
		doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
			src, exists := s.Attr("src")
			if exists {
				html := fmt.Sprintf(`<a id="iframed-link" href="%s">Link</a>`, src)
				s.ReplaceWithHtml(html)
			}
		})

		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, exists := s.Attr("src")
			if exists {
				re := regexp.MustCompile(`(?i)https:\/\/images\.(rapgenius|genius)\.com\/(images\/)?`)
				pSrc := re.ReplaceAllString(src, "/images/")
				s.SetAttr("src", pSrc)
			}
		})

		if source, err := doc.Html(); err == nil {
			body = source
		}
	}

	re := regexp.MustCompile(`https?:\/\/[a-z]*.?genius.com`)
	return re.ReplaceAllString(body, "")
}
