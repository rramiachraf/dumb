package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/data"
	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

func artist(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artistName := mux.Vars(r)["artist"]
		allAlbumsStr := r.URL.Query().Get("allAlbums")
		allAlbums := false
		if allAlbumsStr != "" {
			v, err := strconv.ParseBool(allAlbumsStr)
			if err != nil {
				http.Error(w, "invalid value for allAlbums, use true or false", http.StatusBadRequest)
				return
			}
			allAlbums = v
		}

		cacheKey := fmt.Sprintf("artist:%s-%t", artistName, allAlbums)

		if a, err := getCache[data.Artist](cacheKey); err == nil {
			utils.RenderTemplate(w, views.ArtistPage(a), l)
			return
		}

		url := fmt.Sprintf("https://genius.com/artists/%s", artistName)

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

		var a data.Artist
		if err = a.Parse(doc); err != nil {
			l.Error(err.Error())
		}

		if allAlbums {
			err = a.GetAllAlbums()
			if err != nil {
				l.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				utils.RenderTemplate(w, views.ErrorPage(500, "something went wrong fetching all albums"), l)
				return
			}
			a.AlbumsComplete = true
		}

		utils.RenderTemplate(w, views.ArtistPage(a), l)

		if err = setCache(cacheKey, a); err != nil {
			l.Error(err.Error())
		}
	}
}
