package handlers

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

type static interface {
	Open(string) (fs.File, error)
}

func staticAssets(logger *utils.Logger, embededFiles static) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := strings.Replace(r.URL.Path, "/static", "static", 1)
		f, err := embededFiles.Open(url)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			views.ErrorPage(http.StatusNotFound, "page not found")
			return
		}

		defer f.Close()

		mimeType := mime.TypeByExtension(path.Ext(r.URL.Path))
		w.Header().Set("content-type", mimeType)

		if _, err := io.Copy(w, f); err != nil {
			logger.Error(err.Error())
		}
	}
}
