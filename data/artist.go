package data

import (
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/rramiachraf/dumb/utils"
)

type ArtistPreview struct {
	Name string
	URL  string
}

type Artist struct {
	Name        string
	Description string
	Albums      []AlbumPreview
	Image       string
}

type artistMetadata struct {
	Artist struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description_preview"`
		Image       string `json:"image_url"`
	}
	Albums []struct {
		Id    int    `json:"id"`
		Image string `json:"cover_art_thumbnail_url"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"artist_albums"`
}

func (a *Artist) parseArtistData(doc *goquery.Document) error {
	pageMetadata, exists := doc.Find("meta[itemprop='page_data']").Attr("content")
	if !exists {
		return nil
	}

	var artistMetadataFromPage artistMetadata
	if err := json.Unmarshal([]byte(pageMetadata), &artistMetadataFromPage); err != nil {
		return err
	}

	a.Name = artistMetadataFromPage.Artist.Name
	a.Description = artistMetadataFromPage.Artist.Description
	a.Image = artistMetadataFromPage.Artist.Image

	for _, album := range artistMetadataFromPage.Albums {
		a.Albums = append(a.Albums, AlbumPreview{
			Name:  album.Name,
			Image: album.Image,
			URL:   utils.TrimURL(album.URL),
		})
	}

	return nil
}

func (a *Artist) Parse(doc *goquery.Document) error {
	return a.parseArtistData(doc)
}
