package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/rramiachraf/dumb/utils"
)

func TestArtist(t *testing.T) {
	url := "/artists/Red-hot-chili-peppers"
	name := "Red Hot Chili Peppers"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)
	m := New(l, &assets{})

	m.ServeHTTP(rr, r)

	defer rr.Result().Body.Close()

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected %d, got %d\n", http.StatusOK, rr.Result().StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	artistName := doc.Find("#metadata-info > h1").First().Text()
	if artistName != name {
		t.Fatalf("expected %q, got %q\n", name, artistName)
	}
}

func TestArtistWithAllAlbums(t *testing.T) {
	url := "/artists/Red-hot-chili-peppers?allAlbums=true"
	name := "Red Hot Chili Peppers"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)
	m := New(l, &assets{})

	m.ServeHTTP(rr, r)

	defer rr.Result().Body.Close()

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected %d, got %d\n", http.StatusOK, rr.Result().StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	artistName := doc.Find("#metadata-info > h1").First().Text()
	if artistName != name {
		t.Fatalf("expected %q, got %q\n", name, artistName)
	}

	albumsCount := doc.Find("##artist-single-album").Length()
	if albumsCount < 30 || albumsCount > 200 {
		t.Fatalf("expected album count to be between 30 and 200, was %q\n", albumsCount)
	}
}
