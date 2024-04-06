package data

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Album struct {
	Artist string
	Name   string
	Image  string
	About  [2]string

	Tracks []Track
}

type Track struct {
	Title string
	Url   string
}

type albumMetadata struct {
	Album struct {
		Id          int    `json:"id"`
		Image       string `json:"cover_art_thumbnail_url"`
		Name        string `json:"name"`
		Description string `json:"description_preview"`
		Artist      `json:"artist"`
	}
	AlbumAppearances []AlbumAppearances `json:"album_appearances"`
}

type AlbumAppearances struct {
	Id          int `json:"id"`
	TrackNumber int `json:"track_number"`
	Song        struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	}
}

type Artist struct {
	Name string `json:"name"`
}

func (a *Album) parseAlbumData(doc *goquery.Document) error {
	pageMetadata, exists := doc.Find("meta[itemprop='page_data']").Attr("content")
	if !exists {
		return nil
	}

	var albumMetadataFromPage albumMetadata
	if err := json.Unmarshal([]byte(pageMetadata), &albumMetadataFromPage); err != nil {
		return err
	}

	albumData := albumMetadataFromPage.Album
	a.Artist = albumData.Artist.Name
	a.Name = albumData.Name
	a.Image = albumData.Image
	a.About[0] = albumData.Description
	a.About[1] = truncateText(albumData.Description)

	for _, track := range albumMetadataFromPage.AlbumAppearances {
		url := strings.Replace(track.Song.Url, "https://genius.com", "", -1)
		a.Tracks = append(a.Tracks, Track{Title: track.Song.Title, Url: url})
	}

	return nil
}

func (a *Album) Parse(doc *goquery.Document) error {
	return a.parseAlbumData(doc)
}
