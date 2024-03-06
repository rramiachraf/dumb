package handlers

/*
import (
	"encoding/json"
	"testing"
)

func TestSendRequest(t *testing.T) {
	res, err := sendRequest("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	type fingerprint struct {
		JA3 string `json:"ja3"`
	}

	decoder := json.NewDecoder(res.Body)
	var fg fingerprint

	if err := decoder.Decode(&fg); err != nil {
		t.Fatal(err)
	}

	if fg.JA3 != JA3 {
		t.Fatalf("expected %q, got %q\n", JA3, fg.JA3)
	}
}
*/
