package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rramiachraf/dumb/utils"
)

func TestInstancesList(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/instances.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	l := utils.NewLogger(os.Stdout)

	m := New(l, &assets{})
	m.ServeHTTP(rr, r)

	c := rr.Result().Header.Get("content-type")
	if c != ContentTypeJSON {
		t.Fatalf("expected %q, got %q", ContentTypeJSON, c)
	}

	defer rr.Result().Body.Close()

	d := json.NewDecoder(rr.Result().Body)
	instances := []map[string]any{}
	if err := d.Decode(&instances); err != nil {
		t.Fatalf("unable to decode json from response, %q\n", err)
	}

	if _, exists := instances[0]["clearnet"]; !exists {
		t.Fatal("unable to get clearnet value from instances list")
	}
}
