package handlers

import (
	"bytes"
	"testing"
)

func TestCache(t *testing.T) {
	key := "testkey"
	value := []byte("testvalue")

	err := setCache(key, value)
	if err != nil {
		t.Fatalf("unable to set cache, %q\n", err)
	}

	v, err := getCache[[]byte](key)
	if err != nil {
		t.Fatalf("unable to get cache, %q\n", err)
	}

	if !bytes.Equal(v, value) {
		t.Fatalf("expected %q, got %q\n", value, v)
	}
}
