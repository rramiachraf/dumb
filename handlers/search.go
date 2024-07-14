package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rramiachraf/dumb/data"
	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

func search(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		url := fmt.Sprintf(`https://genius.com/api/search/multi?q=%s`, url.QueryEscape(query))

		res, err := utils.SendRequest(url)
		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderTemplate(w, views.ErrorPage(500, "cannot reach Genius servers"), l)
			return
		}

		defer res.Body.Close()

		var sRes data.SearchResponse

		d := json.NewDecoder(res.Body)
		if err = d.Decode(&sRes); err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderTemplate(w, views.ErrorPage(500, "something went wrong"), l)
		}

		results := data.SearchResults{Query: query, Sections: sRes.Response.Sections}

		utils.RenderTemplate(w, views.SearchPage(results), l)
	}

}
