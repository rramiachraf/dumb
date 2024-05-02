package handlers

import (
	"context"
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

		res, err := sendRequest(url)
		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorPage(500, "cannot reach Genius servers").Render(context.Background(), w)
			return
		}

		defer res.Body.Close()

		var sRes data.SearchResponse

		d := json.NewDecoder(res.Body)
		if err = d.Decode(&sRes); err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorPage(500, "something went wrong").Render(context.Background(), w)
		}

		results := data.SearchResults{Query: query, Sections: sRes.Response.Sections}

		views.SearchPage(results).Render(context.Background(), w)
	}

}
