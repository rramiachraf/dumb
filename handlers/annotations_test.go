package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestAnnotations(t *testing.T) {
	url := "/61590/Black-star-respiration/The-new-moon-rode-high-in-the-crown-of-the-metropolis/annotations"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := logrus.New()
	m := New(l)

	m.ServeHTTP(rr, r)

	defer rr.Result().Body.Close()

	decoder := json.NewDecoder(rr.Result().Body)
	annotation := map[string]string{}

	if err := decoder.Decode(&annotation); err != nil {
		t.Fatal(err)
	}

	if _, exists := annotation["html"]; !exists {
		t.Fatalf("html field not found on annotation\n")
	}
}
