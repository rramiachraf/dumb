package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type song struct {
	Artist string
	Title  string
	Image  string
	Lyrics string
}

func (s *song) parseLyrics(doc *goquery.Document) {
	doc.Find("[data-lyrics-container='true']").Each(func(i int, ss *goquery.Selection) {
		h, err := ss.Html()
		if err != nil {
			log.Println(err)
		}

		s.Lyrics += h
	})
}

func (s *song) parseMetadata(doc *goquery.Document) {
	artist := doc.Find("a[class*='Artist']").First().Text()
	title := doc.Find("h1[class*='Title']").First().Text()
	image, exists := doc.Find("meta[property='og:image']").Attr("content")
	if exists {
		s.Image = image
	}

	s.Title = title
	s.Artist = artist
}

func (s *song) parse(doc *goquery.Document) {
	s.parseLyrics(doc)
	s.parseMetadata(doc)
}

func lyricsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if data, err := getCache(id); err == nil {
		render(w, data)
		return
	}

	url := fmt.Sprintf("https://genius.com/%s-lyrics", id)
	resp, err := http.Get(url)
	if err != nil {
		write(w, http.StatusInternalServerError, []byte("can't reach genius servers"))
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		write(w, http.StatusNotFound, []byte("Not found"))
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		write(w, http.StatusInternalServerError, []byte("something went wrong"))
		return
	}

	var s song
	s.parse(doc)

	w.Header().Set("content-type", "text/html")
	t, err := template.ParseFiles(path.Join("views", "lyrics.tmpl"))
	if err != nil {
		write(w, http.StatusInternalServerError, []byte("something went wrong"))
		return
	}

	t.Execute(w, s)
	setCache(id, s)
}
