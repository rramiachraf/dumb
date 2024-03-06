package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

func TestLyrics(t *testing.T) {
	url := "/The-silver-seas-catch-yer-own-train-lyrics"
	title := "The Silver Seas"
	artist := "Catch Yer Own Train"

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

	docTitle := doc.Find("#metadata-info > h2").Text()
	docArtist := doc.Find("#metadata-info > h1").Text()

	if docTitle != title {
		t.Fatalf("expected %q, got %q\n", title, docTitle)
	}

	if docArtist != artist {
		t.Fatalf("expected %q, got %q\n", artist, docArtist)
	}
}
