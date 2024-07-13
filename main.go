package main

import (
<<<<<<< HEAD
	"context"
=======
	"embed"
>>>>>>> upstream/main
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rramiachraf/dumb/handlers"
	"github.com/rramiachraf/dumb/utils"
)

//go:embed static
var staticFiles embed.FS

func main() {
<<<<<<< HEAD
	ctx := context.Background()
	c, err := bigcache.New(ctx, bigcache.DefaultConfig(time.Hour*24))
	if err != nil {
		logger.Fatalln("can't initialize caching")
	}
	cache = c

	r := mux.NewRouter()

	r.Use(securityHeaders)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { render("home", w, nil) })
	r.HandleFunc("/search", searchHandler).Methods("GET")
	r.HandleFunc("/{id}-lyrics", lyricsHandler)
	r.HandleFunc("/{id}/{artist-song}/{verse}/annotations", annotationsHandler)
	r.HandleFunc("/images/{filename}.{ext}", proxyHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/albums/{artist}/{albumName}", albumHandler)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		render("error", w, map[string]string{
			"Status": "404",
			"Error":  "page not found",
		})

	})
=======
	logger := utils.NewLogger(os.Stdout)
>>>>>>> upstream/main

	server := &http.Server{
		Handler:      handlers.New(logger, staticFiles),
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	PROXY_ENV := os.Getenv("PROXY")
	if PROXY_ENV != "" {
		if _, err := url.ParseRequestURI(PROXY_ENV); err != nil {
			logger.Fatal("invalid proxy")
		}

		logger.Info("using a custom proxy for requests")
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	if port == 0 {
		port = 5555
		logger.Info("using default port %d", port)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("server is listening on port %d", port)

	if err := server.Serve(l); err != nil {
		logger.Fatal(err.Error())
	}
}
