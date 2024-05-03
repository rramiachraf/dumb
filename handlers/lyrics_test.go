package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/rramiachraf/dumb/utils"
)

func TestLyrics(t *testing.T) {
	urls := []string{"/The-silver-seas-catch-yer-own-train-lyrics",
		"/1784308/The-silver-seas-catch-yer-own-train",
		"/1784308/The-silver-seas-catch-yer-own-train-lyrics",
		"/1784308/The-silver-seas-catch-yer-own-train/Baby-you-and-i-are-not-the-same-you-say-you-like-sun-i-like-the-rain",
		"/1784308/The-silver-seas-catch-yer-own-train-lyrics/Baby-you-and-i-are-not-the-same-you-say-you-like-sun-i-like-the-rain",
		"/1784308"}
	for _, url := range urls {
		t.Run(url, func(t *testing.T) { testLyrics(t, url) })
	}
}

func testLyrics(t *testing.T, url string) {
	title := "The Silver Seas"
	artist := "Catch Yer Own Train"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)
	m := New(l, &assets{})

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
