package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

type route struct {
	Path    string
	Handler func(*utils.Logger) http.HandlerFunc
	Method  string
}

func New(logger *utils.Logger, staticFiles static) *mux.Router {
	r := mux.NewRouter()

	r.Use(utils.MustHeaders)
	r.Use(gorillaHandlers.CompressHandler)

	routes := []route{
		{Path: "/albums/{artist}/{albumName}", Handler: album},
		{Path: "/artists/{artist}", Handler: artist},
		{Path: "/images/{filename}.{ext}", Handler: imageProxy},
		{Path: "/search", Handler: search},
		{Path: "/{annotation-id}/{artist-song}/{verse}/annotations", Handler: annotations},
		{Path: "/instances.json", Handler: instances},
	}

	for _, rr := range routes {
		method := rr.Method
		if method == "" {
			method = http.MethodGet
		}

		r.HandleFunc(rr.Path, rr.Handler(logger)).Methods(method)
	}

	r.PathPrefix("/static/").HandlerFunc(staticAssets(logger, staticFiles))

	r.Handle("/", templ.Handler(views.HomePage()))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nDisallow: /\n"))
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
