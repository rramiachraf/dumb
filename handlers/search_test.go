package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/rramiachraf/dumb/utils"
)

func TestSearch(t *testing.T) {
	url := "/search?q=it+aint+hard+to+tell"
	artist := "Nas"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)
	m := New(l)

	m.ServeHTTP(rr, r)

	defer rr.Result().Body.Close()

	doc, err := goquery.NewDocumentFromReader(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	docArtist := doc.Find("#search-item > div > span").First().Text()
	if docArtist != artist {
		t.Fatalf("expected %q, got %q\n", artist, docArtist)
	}
}
