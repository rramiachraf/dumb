package handlers

import (
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rramiachraf/dumb/utils"
)

func TestStaticAssets(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/static/style.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)
	m := New(l, &assets{})

	m.ServeHTTP(rr, r)

	contentType := rr.Header().Get("content-type")
	expectedContentType := mime.TypeByExtension(".css")

	if contentType != expectedContentType {
		t.Fatalf("expected %q, got %q", expectedContentType, contentType)
	}

	if rr.Code != 200 {
		t.Fatalf("expected %d, got %d", 200, rr.Code)
	}
}
