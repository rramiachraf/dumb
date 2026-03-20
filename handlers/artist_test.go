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

	expectedFirstAlbumName := "Road Trippin’ through Time"
	actualFirstAlbumName := doc.Find("#artist-single-album > p").First().Text()

	if expectedFirstAlbumName != actualFirstAlbumName {
		t.Fatalf("expected %q, got %q\n", expectedFirstAlbumName, actualFirstAlbumName)
	}

	expectedFirstAlbumReleaseDate := "May 5, 2011"
	actualFirstAlbumReleaseDate := doc.Find("#artist-single-album > small").First().Text()

	if expectedFirstAlbumReleaseDate != actualFirstAlbumReleaseDate {
		t.Fatalf("expected %q, got %q\n", expectedFirstAlbumReleaseDate, expectedFirstAlbumReleaseDate)
	}
}
