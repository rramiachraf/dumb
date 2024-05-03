package handlers

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

func isValidExt(ext string) bool {
	valid := []string{"jpg", "jpeg", "png", "gif"}
	for _, c := range valid {
		if strings.ToLower(ext) == c {
			return true
		}
	}

	return false
}

func imageProxy(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		f := v["filename"]
		ext := v["ext"]

		if !isValidExt(ext) {
			w.WriteHeader(http.StatusBadRequest)
			views.ErrorPage(400, "something went wrong").Render(context.Background(), w)
			return
		}

		// first segment of URL resize the image to reduce bandwith usage.
		url := fmt.Sprintf("https://t2.genius.com/unsafe/300x300/https://images.genius.com/%s.%s", f, ext)

		res, err := utils.SendRequest(url)
		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorPage(500, "cannot reach Genius servers").Render(context.Background(), w)
			return
		}

		if res.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorPage(500, "something went wrong").Render(context.Background(), w)
			return
		}

		w.Header().Add("Content-type", mime.TypeByExtension("."+ext))
		w.Header().Add("Cache-Control", "max-age=1296000")
		if _, err = io.Copy(w, res.Body); err != nil {
			l.Error("unable to write image, %s", err.Error())
		}
	}
}
