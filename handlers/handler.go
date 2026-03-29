package handlers

import (
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

type ResponseType struct {
	asApi bool
}

func New(logger *utils.Logger, staticFiles static) *mux.Router {
	r := mux.NewRouter()

	r.Use(utils.MustHeaders)
	r.Use(gorillaHandlers.CompressHandler)

	rH := ResponseType{asApi: false}
	rJ := ResponseType{asApi: true}

	routes := []route{
		{Path: "/", Template: views.HomePage},
		{Path: "/robots.txt", Handler: robotsHandler},
		{Path: "/albums/{artist}/{albumName}", Handler: rH.albums},
		{Path: "/artists/{artist}", Handler:  rH.artist},
		{Path: "/a/{article}", Handler: rH.article},
		{Path: "/images/{filename}.{ext}", Handler: imageProxy},
		{Path: "/search", Handler:  rH.search},
		{Path: "/{annotation-id}/{artist-song}/{verse}/annotations", Handler:  rH.annotations},
		{Path: "/{annotation-id}/{artist-song}/annotations", Handler:  rH.annotations},
		{Path: "/instances.json", Handler: instances},

		{Path: "api/v1/albums/{artist}/{albumName}", Handler: rJ.albums},
		{Path: "api/v1/artists/{artist}", Handler: rJ.artist},
		{Path: "api/v1/annotations/{annotation-id}", Handler: rJ.annotations},
		{Path: "api/v1/search", Handler: rJ.search},
	}

	registerRoutes(r, routes, logger)

	r.PathPrefix("/static/").HandlerFunc(staticAssets(logger, staticFiles))
	r.PathPrefix("/{annotation-id}/{artist-song}-lyrics").HandlerFunc(rH.lyrics(logger)).Methods("GET")
	r.PathPrefix("/{annotation-id}/{artist-song}").HandlerFunc(rH.lyrics(logger)).Methods("GET")
	r.PathPrefix("/{annotation-id}").HandlerFunc(rH.lyrics(logger)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		utils.RenderTemplate(w, views.ErrorPage(404, "page not found"), logger)
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
