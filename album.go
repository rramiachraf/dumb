package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type album struct {
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

func (a *album) parseAlbumData(doc *goquery.Document) {
	pageMetadata, exists := doc.Find("meta[itemprop='page_data']").Attr("content")
	if !exists {
		return
	}

	var albumMetadataFromPage albumMetadata
	json.Unmarshal([]byte(pageMetadata), &albumMetadataFromPage)

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
}

func (a *album) parse(doc *goquery.Document) {
	a.parseAlbumData(doc)
}

func albumHandler(w http.ResponseWriter, r *http.Request) {
	artist := mux.Vars(r)["artist"]
	albumName := mux.Vars(r)["albumName"]

	id := fmt.Sprintf("%s/%s", artist, albumName)

	if data, err := getCache(id); err == nil {
		render("album", w, data)
		return
	}

	url := fmt.Sprintf("https://genius.com/albums/%s/%s", artist, albumName)

	resp, err := sendRequest(url)
	if err != nil {
		logger.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "cannot reach genius servers",
		})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		render("error", w, map[string]string{
			"Status": "404",
			"Error":  "page not found",
		})
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "something went wrong",
		})
		return
	}

	cf := doc.Find(".cloudflare_content").Length()
	if cf > 0 {
		logger.Errorln("cloudflare got in the way")
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "damn cloudflare, issue #21 on GitHub",
		})
		return
	}

	var a album
	a.parse(doc)

	render("album", w, a)

	setCache(id, a)
}
