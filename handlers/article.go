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

func article(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articleSlug := mux.Vars(r)["article"]

		if a, err := getCache[data.Article](articleSlug); err == nil {
			utils.RenderTemplate(w, views.ArticlePage(a), l)
			return
		}

		url := fmt.Sprintf("https://genius.com/a/%s", articleSlug)

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

		var a data.Article
		if err = a.Parse(doc); err != nil {
			l.Error(err.Error())
		}

		utils.RenderTemplate(w, views.ArticlePage(a), l)

		if err = setCache(articleSlug, a); err != nil {
			l.Error(err.Error())
		}
	}
}
