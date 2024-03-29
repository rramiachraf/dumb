package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/views"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(mustHeaders)

	r.Handle("/", templ.Handler(views.HomePage()))
	r.HandleFunc("/{id}-lyrics", lyrics(logger)).Methods("GET")
	r.HandleFunc("/albums/{artist}/{albumName}", album(logger)).Methods("GET")
	r.HandleFunc("/images/{filename}.{ext}", imageProxy(logger)).Methods("GET")
	r.HandleFunc("/search", search(logger)).Methods("GET")
	r.HandleFunc("/{id}/{artist-song}/{verse}/annotations", annotations(logger)).Methods("GET")
	r.HandleFunc("/instances.json", instances(logger)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		views.ErrorPage(404, "page not found").Render(context.Background(), w)
	})

	return r
}
