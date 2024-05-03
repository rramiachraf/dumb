package handlers

import (
	"context"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/a-h/templ"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

func New(logger *utils.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(utils.MustHeaders)
	r.Use(gorillaHandlers.CompressHandler)

	r.Handle("/", templ.Handler(views.HomePage()))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nDisallow: /\n"))
	})
	r.HandleFunc("/albums/{artist}/{albumName}", album(logger)).Methods("GET")
	r.HandleFunc("/images/{filename}.{ext}", imageProxy(logger)).Methods("GET")
	r.HandleFunc("/search", search(logger)).Methods("GET")
	r.HandleFunc("/{annotation-id}/{artist-song}/{verse}/annotations", annotations(logger)).Methods("GET")
	r.HandleFunc("/instances.json", instances(logger)).Methods("GET")
	r.PathPrefix("/static/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.Replace(r.URL.Path, "/static", "static", 1)
		f, err := os.Open(url)
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

	})
	r.PathPrefix("/{annotation-id}/{artist-song}-lyrics").HandlerFunc(lyrics(logger)).Methods("GET")
	r.PathPrefix("/{annotation-id}/{artist-song}").HandlerFunc(lyrics(logger)).Methods("GET")
	r.PathPrefix("/{annotation-id}").HandlerFunc(lyrics(logger)).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		views.ErrorPage(404, "page not found").Render(context.Background(), w)
	})

	return r
}
