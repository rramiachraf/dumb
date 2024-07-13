package handlers

import (
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rramiachraf/dumb/utils"
)

func TestImageProxy(t *testing.T) {
	imgURL := "/images/3852401ae6c6d6a51aafe814d67199f0.1000x1000x1.jpg"

	r, err := http.NewRequest(http.MethodGet, imgURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)
	m := New(l, &assets{})

	m.ServeHTTP(rr, r)

	cc := rr.Result().Header.Get("cache-control")
	maxAge := "max-age=1296000"
	ct := rr.Result().Header.Get("content-type")
	mimeType := mime.TypeByExtension(".jpg")

	if cc != maxAge {
		t.Fatalf("expected %q, got %q\n", maxAge, cc)
	}

	if ct != mimeType {
		t.Fatalf("expected %q, got %q\n", mimeType, ct)
	}
}
