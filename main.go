package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	c, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 24))
	if err != nil {
		logger.Fatalln("can't initialize caching")
	}
	cache = c

	r := mux.NewRouter()

	r.Use(securityHeaders)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { render("home", w, nil) })
	r.HandleFunc("/search", fullSearchHandler).Methods("GET")
	r.HandleFunc("/form_search", formSearchHandler).Methods("POST")
	r.HandleFunc("/{id}-lyrics", lyricsHandler)
	r.HandleFunc("/images/{filename}.{ext}", proxyHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		render("error", w, map[string]string{
			"Status": "404",
			"Error":  "page not found",
		})

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
