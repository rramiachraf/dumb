package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/rramiachraf/dumb/utils"
)

func TestArticle(t *testing.T) {
	url := "/a/genius-celebrates-hip-hops-50th-anniversary-with-a-look-back-at-the-music-thats-defined-this-site"
	title := "Genius Celebrates Hip-Hop’s 50th Anniversary With A Look Back At The Music That’s Defined This Site"
	subtitle := "The first post in a yearlong look at the genre’s storied history."

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

	articleTitle := doc.Find("#article-title").First().Text()
	if articleTitle != title {
		t.Fatalf("expected %q, got %q\n", title, articleTitle)
	}

	articleSubtitle := doc.Find("#article-subtitle").First().Text()
	if articleSubtitle != subtitle {
		t.Fatalf("expected %q, got %q\n", subtitle, articleSubtitle)
	}

	articleBody := doc.Find("#article-body").First().Text()
	if len(articleBody) == 0 {
		t.Fatal("missing article body\n")
	}
}
