package utils

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestSendRequestReusesProxyConnections(t *testing.T) {
	var conns atomic.Int32

	proxy := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	proxy.Config.ConnState = func(_ net.Conn, s http.ConnState) {
		if s == http.StateNew {
			conns.Add(1)
		}
	}
	proxy.Start()
	defer proxy.Close()

	t.Setenv("PROXY", proxy.URL)

	const requests = 5
	for range requests {
		res, err := SendRequest("http://upstream.invalid/")
		if err != nil {
			t.Fatal(err)
		}

		_, _ = io.Copy(io.Discard, res.Body)
		if err := res.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}

	if n := conns.Load(); n != 1 {
		t.Fatalf("expected %d requests to share 1 proxy connection, got %d connections", requests, n)
	}
}
