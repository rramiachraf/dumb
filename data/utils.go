package data

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/caffix/cloudflare-roundtripper/cfrt"
)

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
