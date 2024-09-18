package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rramiachraf/dumb/data"
	"github.com/rramiachraf/dumb/utils"
)

func TestAnnotations(t *testing.T) {
	url := "/943841/Black-star-respiration/Shinin-like-who-on-top-of-this/annotations"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)
	m := New(l, &assets{})

	m.ServeHTTP(rr, r)

	defer rr.Result().Body.Close()

	decoder := json.NewDecoder(rr.Result().Body)
	var annotation data.Annotation

	if err := decoder.Decode(&annotation); err != nil {
		t.Fatal(err)
	}

	if annotation.State != "accepted" {
		t.Fatalf("expected state to be %q, got %q\n", "accepted", annotation.State)
	}
}
