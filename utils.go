package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"text/template"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/caffix/cloudflare-roundtripper/cfrt"
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
	_, err := w.Write(data)
	if err != nil {
		logger.Errorln(err)
	}

}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; object-src 'none'"
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
	w.Header().Set("content-type", "text/html")
	t := template.New(n + ".tmpl").Funcs(template.FuncMap{"extractURL": extractURL})
	t, err := t.ParseFiles(getTemplates(n, "navbar", "footer")...)
	if err != nil {
		logger.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = t.Execute(w, data); err != nil {
		logger.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

const UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"

var client = &http.Client{
	Timeout: 20 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 15 * time.Second,
			DualStack: true,
		}).DialContext,
	},
}

func sendRequest(u string) (*http.Response, error) {
	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	client.Transport, err = cfrt.New(client.Transport)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: map[string][]string{
			"Accept-Language": {"en-US"},
			"User-Agent":      {UA},
		},
	}

	return client.Do(req)
}
