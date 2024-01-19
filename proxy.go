package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func extractURL(image string) string {
	u, err := url.Parse(image)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("/images%s", u.Path)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	f := v["filename"]
	ext := v["ext"]

	if !isValidExt(ext) {
		w.WriteHeader(http.StatusBadRequest)
		render("error", w, map[string]string{
			"Status": "400",
			"Error":  "Something went wrong",
		})
		return
	}

	// first segment of URL resize the image to reduce bandwith usage.
	url := fmt.Sprintf("https://t2.genius.com/unsafe/300x300/https://images.genius.com/%s.%s", f, ext)

	res, err := sendRequest(url)
	if err != nil {
		logger.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "cannot reach genius servers",
		})
		return
	}

	if res.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "something went wrong",
		})

		return
	}

	w.Header().Add("Content-type", fmt.Sprintf("image/%s", ext))
	_, err = io.Copy(w, res.Body)
	if err != nil {
		logger.Errorln(err)
	}

}
