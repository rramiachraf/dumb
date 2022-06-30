package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func fatal(err any) {
	log.Fatalf("[ERR] %s\n", err)
}

func info(s string) {
	log.Printf("[INFO] %s\n", s)
}

func write(w http.ResponseWriter, status int, data []byte) {
	w.WriteHeader(status)
	w.Write(data)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{id}-lyrics", lyricsHandler)
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
