package handlers

import (
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
			utils.RenderTemplate(w, views.LyricsPage(s), l)
			return
		}

		url := fmt.Sprintf("https://genius.com/%s", id)
		resp, err := utils.SendRequest(url)
		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderTemplate(w, views.ErrorPage(500, "cannot reach Genius servers"), l)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			utils.RenderTemplate(w, views.ErrorPage(404, "page not found"), l)
			return
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderTemplate(w, views.ErrorPage(500, "something went wrong"), l)
			return
		}

		cf := doc.Find(".cloudflare_content").Length()
		if cf > 0 {
			l.Error("cloudflare got in the way")
			utils.RenderTemplate(w, views.ErrorPage(500, "cloudflare is detected"), l)
			return
		}

		var s data.Song
		if err := s.Parse(doc); err != nil {
			l.Error(err.Error())
		}

		utils.RenderTemplate(w, views.LyricsPage(s), l)

		if err = setCache(id, s); err != nil {
			l.Error(err.Error())
		}
	}
}
