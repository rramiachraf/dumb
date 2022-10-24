package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func isValidExt(ext string) bool {
	valid := []string{"jpg", "jpeg", "png", "gif"}
	for _, c := range valid {
		if strings.ToLower(ext) == c {
			return true
		}
	}

	return false
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	f := v["filename"]
	ext := v["ext"]

	if !isValidExt(ext) {
		write(w, http.StatusBadRequest, []byte("not an image :/"))
		return
	}

	url := fmt.Sprintf("https://images.genius.com/%s.%s", f, ext)

	res, err := http.Get(url)
	if err != nil {
		write(w, http.StatusInternalServerError, []byte("can't reach genius servers"))
		return
	}

	if res.StatusCode != http.StatusOK {
		write(w, res.StatusCode, []byte{})
		return
	}

	w.Header().Add("Content-type", fmt.Sprintf("image/%s", ext))
	io.Copy(w, res.Body)
}
