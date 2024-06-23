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
	Path     string
	Handler  func(*utils.Logger) http.HandlerFunc
	Method   string
	Template func() templ.Component
}

func New(logger *utils.Logger, staticFiles static) *mux.Router {
	r := mux.NewRouter()

	r.Use(utils.MustHeaders)
	r.Use(gorillaHandlers.CompressHandler)

	routes := []route{
		{Path: "/", Template: views.HomePage},
		{Path: "/robots.txt", Handler: robotsHandler},
		{Path: "/albums/{artist}/{albumName}", Handler: album},
		{Path: "/artists/{artist}", Handler: artist},
		{Path: "/a/{article}", Handler: article},
		{Path: "/images/{filename}.{ext}", Handler: imageProxy},
		{Path: "/search", Handler: search},
		{Path: "/{annotation-id}/{artist-song}/{verse}/annotations", Handler: annotations},
		{Path: "/instances.json", Handler: instances},
	}

	registerRoutes(r, routes, logger)

	r.PathPrefix("/static/").HandlerFunc(staticAssets(logger, staticFiles))
	r.PathPrefix("/{annotation-id}/{artist-song}-lyrics").HandlerFunc(lyrics(logger)).Methods("GET")
	r.PathPrefix("/{annotation-id}/{artist-song}").HandlerFunc(lyrics(logger)).Methods("GET")
	r.PathPrefix("/{annotation-id}").HandlerFunc(lyrics(logger)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		views.ErrorPage(404, "page not found").Render(context.Background(), w)
	})

	return r
}

func robotsHandler(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("User-agent: *\nDisallow: /\n")); err != nil {
			l.Error(err.Error())
		}
	}
}

func registerRoutes(router *mux.Router, routes []route, logger *utils.Logger) {
	for _, r := range routes {
		method := r.Method
		if method == "" {
			method = http.MethodGet
		}

		if r.Template != nil {
			router.Handle(r.Path, templ.Handler(r.Template())).Methods(method)
			continue
		}

		router.HandleFunc(r.Path, r.Handler(logger)).Methods(method)
	}
}
