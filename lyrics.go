package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type song struct {
	Artist      string
	Title       string
	Image       string
	Lyrics      string
	Credits     map[string]string
	About       [2]string
	Album       string
	LinkToAlbum string
}

type songResponse struct {
	Response struct {
		Song struct {
			ArtistNames string `json:"artist_names"`
			Image       string `json:"song_art_image_thumbnail_url"`
			Title       string
			Description struct {
				Plain string
			}
			Album struct {
				Url  string `json:"url"`
				Name string `json:"name"`
			}
			CustomPerformances []customPerformance `json:"custom_performances"`
		}
	}
}

type customPerformance struct {
	Label   string
	Artists []struct {
		Name string
	}
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

func (s *song) parseSongData(doc *goquery.Document) {
	attr, exists := doc.Find("meta[property='twitter:app:url:iphone']").Attr("content")
	if exists {
		songID := strings.Replace(attr, "genius://songs/", "", 1)

		u := fmt.Sprintf("https://genius.com/api/songs/%s?text_format=plain", songID)

		res, err := sendRequest(u)
		if err != nil {
			logger.Errorln(err)
		}

		defer res.Body.Close()

		var data songResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&data)
		if err != nil {
			logger.Errorln(err)
		}

		songData := data.Response.Song
		s.Artist = songData.ArtistNames
		s.Image = songData.Image
		s.Title = songData.Title
		s.About[0] = songData.Description.Plain
		s.About[1] = truncateText(songData.Description.Plain)
		s.Credits = make(map[string]string)
		s.Album = songData.Album.Name
		s.LinkToAlbum = strings.Replace(songData.Album.Url, "https://genius.com", "", -1)

		for _, perf := range songData.CustomPerformances {
			var artists []string
			for _, artist := range perf.Artists {
				artists = append(artists, artist.Name)
			}
			s.Credits[perf.Label] = strings.Join(artists, ", ")
		}
	}
}

func truncateText(text string) string {
	textArr := strings.Split(text, "")

	if len(textArr) > 250 {
		return strings.Join(textArr[0:250], "") + "..."
	}

	return text
}

func (s *song) parse(doc *goquery.Document) {
	s.parseLyrics(doc)
	s.parseSongData(doc)
}

func lyricsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if data, err := getCache(id); err == nil {
		render("lyrics", w, data)
		return
	}

	url := fmt.Sprintf("https://genius.com/%s-lyrics", id)
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

	var s song
	s.parse(doc)

	render("lyrics", w, s)
	err = setCache(id, s)
	if err != nil {
		logger.Errorln(err)
	}

}
