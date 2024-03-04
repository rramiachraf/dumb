package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/handlers"
	"github.com/rramiachraf/dumb/views"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	r := mux.NewRouter()

	r.Use(handlers.MustHeaders)

	r.Handle("/", templ.Handler(views.HomePage()))
	r.HandleFunc("/{id}-lyrics", handlers.Lyrics(logger)).Methods("GET")
	r.HandleFunc("/albums/{artist}/{albumName}", handlers.Album(logger)).Methods("GET")
	r.HandleFunc("/images/{filename}.{ext}", handlers.ImageProxy(logger)).Methods("GET")
	r.HandleFunc("/search", handlers.Search(logger)).Methods("GET")
	r.HandleFunc("/{id}/{artist-song}/{verse}/annotations", handlers.Annotations(logger)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		views.ErrorPage(404, "page not found").Render(context.Background(), w)
	})

	server := &http.Server{
		Handler:      r,
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	if port == 0 {
		port = 5555
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Infof("server is listening on port %d\n", port)

	logger.Fatalln(server.Serve(l))
}
