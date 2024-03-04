package handlers

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

func MustHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; object-src 'none'"
		w.Header().Add("content-security-policy", csp)
		w.Header().Add("referrer-policy", "no-referrer")
		w.Header().Add("x-content-type-options", "nosniff")
		next.ServeHTTP(w, r)
	})
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
