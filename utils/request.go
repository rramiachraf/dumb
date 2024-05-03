package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func MustHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; object-src 'none'"
		w.Header().Add("content-security-policy", csp)
		w.Header().Add("referrer-policy", "no-referrer")
		w.Header().Add("x-content-type-options", "nosniff")
		w.Header().Add("content-type", "text/html")
		next.ServeHTTP(w, r)
	})
}

const UA = "Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0"

func SendRequest(u string) (*http.Response, error) {
	proxy := os.Getenv("PROXY")
	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("unable to parse proxy url: %s", err.Error())
		}

		client.Timeout = time.Minute * 1
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

	}

	headers := http.Header{}
	headers.Set("user-agent", UA)

	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: headers,
	}

	return client.Do(req)
}
