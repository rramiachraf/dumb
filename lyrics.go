package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type song struct {
	Artist  string
	Title   string
	Image   string
	Lyrics  string
	Credits map[string]string
	About   [2]string
}

func (s *song) parseLyrics(doc *goquery.Document) {
	doc.Find("[data-lyrics-container='true']").Each(func(i int, ss *goquery.Selection) {
		h, err := ss.Html()
		if err != nil {
			logger.Errorln("unable to parse lyrics", err)
		}
		s.Lyrics += h
	})
}

func (s *song) parseMetadata(doc *goquery.Document) {
	artist := doc.Find("a[class*='Artist']").First().Text()
	title := doc.Find("h1[class*='Title']").First().Text()
	image, exists := doc.Find("meta[property='og:image']").Attr("content")
	if exists {
		s.Image = extractURL(image)
	}

	s.Title = title
	s.Artist = artist
}

func (s *song) parseCredits(doc *goquery.Document) {
	credits := make(map[string]string)

	doc.Find("[class*='SongInfo__Credit']").Each(func(i int, ss *goquery.Selection) {
		key := ss.Children().First().Text()
		value := ss.Children().Last().Text()
		credits[key] = value
	})

	s.Credits = credits
}

func (s *song) parseAbout(doc *goquery.Document) {
	s.About[0] = doc.Find("[class*='SongDescription__Content']").Text()
	summary := strings.Split(s.About[0], "")

	if len(summary) > 250 {
		s.About[1] = strings.Join(summary[0:250], "") + "..."
	}
}

func (s *song) parse(doc *goquery.Document) {
	s.parseLyrics(doc)
	s.parseMetadata(doc)
	s.parseCredits(doc)
	s.parseAbout(doc)
}

func lyricsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if data, err := getCache(id); err == nil {
		render("lyrics", w, data)
		return
	}

	url := fmt.Sprintf("https://genius.com/%s-lyrics", id)
	resp, err := http.Get(url)
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

	var s song
	s.parse(doc)

	render("lyrics", w, s)
	setCache(id, s)
}
