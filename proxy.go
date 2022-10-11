package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	f := mux.Vars(r)["filename"]
	url := fmt.Sprintf("https://images.genius.com/%s.jpg", f)

	res, err := http.Get(url)
	if err != nil {
		write(w, http.StatusInternalServerError, []byte("can't reach genius genius servers"))
		return
	}

	if res.StatusCode != http.StatusOK {
		write(w, res.StatusCode, []byte{})
		return
	}

	w.Header().Add("Content-type", "image/jpeg")
	io.Copy(w, res.Body)
}
