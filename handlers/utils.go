package handlers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	fhttp "github.com/Danny-Dasilva/fhttp"
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

const UA = "Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0"
const JA3 = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53,0-23-65281-10-11-16-5-34-51-43-13-45-28-65037-41,29-23-24-25-256-257,0"

func sendRequest(u string) (*fhttp.Response, error) {
	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	client := &fhttp.Client{
		Transport: cycletls.NewTransport(JA3, UA),
		Timeout:   20 * time.Second,
	}

	req := &fhttp.Request{
		Method: http.MethodGet,
		URL:    url,
	}

	return client.Do(req)
}
