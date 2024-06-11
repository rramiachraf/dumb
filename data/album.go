package data

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rramiachraf/dumb/utils"
)

type AlbumPreview struct {
	Name  string
	Image string
	URL   string
}

type Album struct {
	AlbumPreview
	Artist ArtistPreview
	About  string

	Tracks []Track
}

type Track struct {
	Title  string
	Url    string
	Number int
}

type albumMetadata struct {
	Album struct {
		Id                    int    `json:"id"`
		Image                 string `json:"cover_art_thumbnail_url"`
		Name                  string `json:"name"`
		Description           string `json:"description_preview"`
		artistPreviewMetadata `json:"artist"`
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

type artistPreviewMetadata struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
	a.Artist = ArtistPreview{
		Name: albumData.artistPreviewMetadata.Name,
		URL:  utils.TrimURL(albumData.artistPreviewMetadata.URL),
	}
	a.Name = albumData.Name
	a.Image = albumData.Image
	a.About = albumData.Description

	for _, track := range albumMetadataFromPage.AlbumAppearances {
		url := strings.Replace(track.Song.Url, "https://genius.com", "", -1)
		a.Tracks = append(a.Tracks, Track{Title: track.Song.Title, Url: url, Number: track.TrackNumber})
	}

	return nil
}

func (a *Album) Parse(doc *goquery.Document) error {
	return a.parseAlbumData(doc)
}
