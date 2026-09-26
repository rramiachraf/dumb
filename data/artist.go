package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/rramiachraf/dumb/utils"
)

type ArtistPreview struct {
	Name string
	URL  string
}

type Artist struct {
	Id             int
	Name           string
	Description    string
	Albums         []AlbumPreview
	AlbumsComplete bool
	Image          string
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

	a.Id = artistMetadataFromPage.Artist.Id
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

type artistAlbumsMetadata struct {
	Response struct {
		Albums []struct {
			Id    int    `json:"id"`
			Image string `json:"cover_art_thumbnail_url"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		}
		NextPage int `json:"next_page"`
	}
}

func (a *Artist) GetAllAlbums() error {
	nextPage := 1
	a.Albums = []AlbumPreview{}

	for nextPage > 0 {
		url := fmt.Sprintf("https://genius.com/api/artists/%d/albums?page=%d", a.Id, nextPage)

		resp, err := utils.SendRequest(url)
		if err != nil {
			return fmt.Errorf("Failed reqquest for url %s: %w", url, err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("http error: status code %s for url %s", resp.Status, url)
		}

		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Failed reading response from for %s: %w", url, err)
		}

		var artistAlbumsMetadata artistAlbumsMetadata
		if err := json.Unmarshal(resBody, &artistAlbumsMetadata); err != nil {
			return fmt.Errorf("Failed unmarshalling response for url %s: %w", url, err)
		}

		for _, album := range artistAlbumsMetadata.Response.Albums {
			a.Albums = append(a.Albums,
				AlbumPreview{
					Name:  album.Name,
					Image: album.Image,
					URL:   utils.TrimURL(album.URL),
				})
		}

		nextPage = artistAlbumsMetadata.Response.NextPage
	}

	return nil
}
