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
)

func main() {
	c, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 2))
	if err != nil {
		fatal("can't initialize caching")
	}
	cache = c

	r := mux.NewRouter()

	r.Use(securityHeaders)

	r.HandleFunc("/{id}-lyrics", lyricsHandler)
	r.HandleFunc("/images/{filename}.jpg", proxyHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := &http.Server{
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	if port == 0 {
		port = 5555
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fatal(err)
	}

	info(fmt.Sprintf("server is listening on port %d", port))

	fatal(server.Serve(l))
}
