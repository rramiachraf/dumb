package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Response struct {
		Sections sections
	}
}

type result struct {
	ArtistNames string `json:"artist_names"`
	Title       string
	Path        string
	Thumbnail   string `json:"song_art_image_thumbnail_url"`
}

type hits []struct {
	Result result
}

type sections []struct {
	Type string
	Hits hits
}

type renderVars struct {
	Query    string
	Sections sections
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	url := fmt.Sprintf(`https://genius.com/api/search/multi?q=%s`, query)

	res, err := http.Get(url)
	if err != nil {
		logger.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "cannot reach genius servers",
		})
	}

	defer res.Body.Close()

	var data response

	d := json.NewDecoder(res.Body)
	d.Decode(&data)

	vars := renderVars{query, data.Response.Sections}

	render("search", w, vars)
}
