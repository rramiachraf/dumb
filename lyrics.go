package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
	"github.com/valyala/fastjson"
)

type song struct {
	Artist       string
	Title        string
	Image        string
	Lyrics       string
	Credits      map[string]string
	About        [2]string
	PrimaryColor string
}

func fixJSON(in []byte) []byte {
	var out = in
	replaceList := map[string]string{
		`{\"`:  `{"`,
		`\":`:  `":`,
		`:\"`:  `:"`,
		`\",`:  `",`,
		`,\"`:  `,"`,
		`\"}`:  `"}`,
		`[\"`:  `["`,
		`\"],`: `"],`,
		`\"]}`: `"]}`,
		`\\n`:  ``,
		`\'`:   `'`,
		`\\"`:  `"`,
	}

	for match, replacer := range replaceList {
		out = bytes.ReplaceAll(out, []byte(match), []byte(replacer))
	}

	return out
}

func (s *song) parse(urlPath string, preload []byte) {
	jsonData := fixJSON(preload)

	var parser fastjson.Parser

	v, err := parser.Parse(string(jsonData))
	if err != nil {
		logger.Errorf(`%s: %s\n`, urlPath, err)
	}

	v.Del("currentPage")
	v.Del("deviceType")
	v.Del("session")

	s.Lyrics = string(v.GetStringBytes("songPage", "lyricsData", "body", "html"))
	s.Credits = make(map[string]string)

	v.GetObject("entities", "songs").Visit(func(key []byte, v *fastjson.Value) {
		path := strings.ToLower(string(v.GetStringBytes("path")))
		if path == urlPath {
			s.Title = string(v.GetStringBytes("title"))
			s.Artist = string(v.GetStringBytes("artistNames"))
			s.About[0] = string(blackfriday.Run(v.GetStringBytes("description", "markdown")))
			s.About[1] = string(v.GetStringBytes("descriptionPreview"))
			if u, err := url.Parse(string(v.GetStringBytes("songArtImageUrl"))); err == nil {
				s.Image = fmt.Sprintf("/images%s", u.Path)
			}
			s.PrimaryColor = string(v.GetStringBytes("songArtPrimaryColor"))

			for _, v := range v.GetArray("customPerformances") {
				label := v.GetStringBytes("label")
				var artists []string
				for _, v := range v.GetArray("artists") {
					artists = append(artists, string(v.GetStringBytes("name")))
				}
				s.Credits[string(label)] = strings.Join(artists, ", ")
			}
		}
	})
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

	if resp.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		render("error", w, map[string]string{
			"Status": "404",
			"Error":  "page not found",
		})
		return
	}

	bodyHTML, err := io.ReadAll(resp.Body)

	rgx, err := regexp.Compile(`window\.__PRELOADED_STATE__ = JSON.parse\('(.*)'\);`)
	if err != nil {
		logger.Errorln(err)
	}

	preload := rgx.FindSubmatch(bodyHTML)[1]

	var s song
	s.parse(r.URL.RequestURI(), preload)

	render("lyrics", w, s)
	setCache(id, s)
}
