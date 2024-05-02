package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

const ContentTypeJSON = "application/json"

// TODO: move this to utils, so it can be used by other handlers.
func sendError(err error, status int, msg string, l *utils.Logger, w http.ResponseWriter) {
	l.Error(err.Error())
	w.WriteHeader(status)
	if err := views.ErrorPage(status, msg).Render(context.Background(), w); err != nil {
		l.Error(err.Error())
	}
}

func instances(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if instances, err := getCache[[]byte]("instances"); err == nil {
			w.Header().Set("content-type", ContentTypeJSON)
			_, err = w.Write(instances)
			if err == nil {
				return
			}
		}

		res, err := sendRequest("https://raw.githubusercontent.com/rramiachraf/dumb/main/instances.json")
		if err != nil {
			sendError(err, http.StatusInternalServerError, "something went wrong", l, w)
			return
		}

		defer res.Body.Close()

		instances, err := io.ReadAll(res.Body)
		if err != nil {
			sendError(err, http.StatusInternalServerError, "something went wrong", l, w)
			return
		}

		w.Header().Set("content-type", ContentTypeJSON)
		if _, err = w.Write(instances); err != nil {
			l.Error(err.Error())
		}

		if err = setCache("instances", instances); err != nil {
			l.Error(err.Error())
		}
	}
}
