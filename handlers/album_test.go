package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

func TestAlbum(t *testing.T) {
	url := "/albums/Daft-punk/Random-access-memories"
	title := "Give Life Back to Music"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := logrus.New()
	m := New(l)

	m.ServeHTTP(rr, r)

	defer rr.Result().Body.Close()

	doc, err := goquery.NewDocumentFromReader(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	docTitle := doc.Find("#album-tracklist > a > p").First().Text()

	if docTitle != title {
		t.Fatalf("expected %q, got %q\n", title, docTitle)
	}
}
