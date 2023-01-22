package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"

	"github.com/allegro/bigcache/v3"
)

var cache *bigcache.BigCache

func setCache(key string, entry interface{}) error {
	data, err := json.Marshal(&entry)
	if err != nil {
		return err
	}

	return cache.Set(key, data)
}

func getCache(key string) (interface{}, error) {
	data, err := cache.Get(key)
	if err != nil {
		return nil, err
	}

	var decoded interface{}

	if err = json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

func write(w http.ResponseWriter, status int, data []byte) {
	w.WriteHeader(status)
	w.Write(data)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' images.genius.com; object-src 'none'"
		w.Header().Add("content-security-policy", csp)
		w.Header().Add("referrer-policy", "no-referrer")
		w.Header().Add("x-content-type-options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func getTemplates(templates ...string) []string {
	var pths []string
	for _, t := range templates {
		tmpl := path.Join("views", fmt.Sprintf("%s.tmpl", t))
		pths = append(pths, tmpl)
	}
	return pths
}

func render(n string, w http.ResponseWriter, data interface{}) {
	t, err := template.ParseFiles(getTemplates(n, "navbar", "footer")...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = t.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
