package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type Song struct {
	Artist  string
	Title   string
	Image   string
	Lyrics  string
	Credits map[string]string
	About   [2]string
	Album   struct {
		URL   string
		Name  string
		Image string
	}
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
				URL   string `json:"url"`
				Name  string `json:"name"`
				Image string `json:"cover_art_url"`
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

func (s *Song) parseLyrics(doc *goquery.Document) {
	doc.Find("[data-lyrics-container='true']").Each(func(i int, ss *goquery.Selection) {
		h, err := ss.Html()
		if err != nil {
			logrus.Errorln("unable to parse lyrics", err)
		}
		s.Lyrics += h
	})
}

func (s *Song) parseSongData(doc *goquery.Document) {
	attr, exists := doc.Find("meta[property='twitter:app:url:iphone']").Attr("content")
	if exists {
		songID := strings.Replace(attr, "genius://songs/", "", 1)

		u := fmt.Sprintf("https://genius.com/api/songs/%s?text_format=plain", songID)

		res, err := sendRequest(u)
		if err != nil {
			logrus.Errorln(err)
		}

		defer res.Body.Close()

		var data songResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&data)
		if err != nil {
			logrus.Errorln(err)
		}

		songData := data.Response.Song
		s.Artist = songData.ArtistNames
		s.Image = songData.Image
		s.Title = songData.Title
		s.About[0] = songData.Description.Plain
		s.About[1] = truncateText(songData.Description.Plain)
		s.Credits = make(map[string]string)
		s.Album.Name = songData.Album.Name
		s.Album.URL = strings.Replace(songData.Album.URL, "https://genius.com", "", -1)
		s.Album.Image = ExtractImageURL(songData.Album.Image)

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

func (s *Song) Parse(doc *goquery.Document) {
	s.parseLyrics(doc)
	s.parseSongData(doc)
}
