package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/data"
	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

func lyrics(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// prefer artist-song over annotation-id for cache key when available
		id := mux.Vars(r)["artist-song"]
		if id == "" {
			id = mux.Vars(r)["annotation-id"]
		} else {
			id = id + "-lyrics"
		}

		if s, err := getCache[data.Song](id); err == nil {
			views.LyricsPage(s).Render(context.Background(), w)
			return
		}

		url := fmt.Sprintf("https://genius.com/%s", id)
		resp, err := sendRequest(url)
		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorPage(500, "cannot reach Genius servers").Render(context.Background(), w)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			views.ErrorPage(404, "page not found").Render(context.Background(), w)
			return
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorPage(500, "something went wrong").Render(context.Background(), w)
			return
		}

		cf := doc.Find(".cloudflare_content").Length()
		if cf > 0 {
			l.Error("cloudflare got in the way")
			views.ErrorPage(500, "TODO: fix Cloudflare #21").Render(context.Background(), w)
			return
		}

		var s data.Song
		if err := s.Parse(doc); err != nil {
			l.Error(err.Error())
		}

		views.LyricsPage(s).Render(context.Background(), w)
		if err = setCache(id, s); err != nil {
			l.Error(err.Error())
		}
	}
}
