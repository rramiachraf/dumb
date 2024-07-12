package data

import (
	"encoding/json"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rramiachraf/dumb/utils"
)

type Article struct {
	Title       string
	Subtitle    string
	HTML        string
	Authors     []Author
	PublishedAt time.Time
	Image       string
}

type Author struct {
	Name  string
	Role  string `json:"human_readable_role_for_display"`
	About string `json:"about_me_summary"`
}

type articleResponse struct {
	Article struct {
		Title    string
		Subtitle string `json:"dek"`
		Authors  []Author
		Body     struct {
			HTML string
		}
		PublishedAt int64  `json:"published_at"`
		Image       string `json:"preview_image"`
	}
}

func (a *Article) parseArticleData(doc *goquery.Document) error {
	pageMetadata, exists := doc.Find("meta[itemprop='page_data']").Attr("content")
	if !exists {
		return nil
	}

	var articleData articleResponse
	if err := json.Unmarshal([]byte(pageMetadata), &articleData); err != nil {
		return err
	}
	data := articleData.Article

	a.Title = data.Title
	a.Subtitle = data.Subtitle

	a.HTML = utils.CleanBody(data.Body.HTML)
	a.Authors = data.Authors
	a.PublishedAt = time.Unix(data.PublishedAt, 0)
	a.Image = ExtractImageURL(data.Image)

	return nil
}

func (a *Article) Parse(doc *goquery.Document) error {
	return a.parseArticleData(doc)
}
