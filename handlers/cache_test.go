package handlers

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/allegro/bigcache/v3"
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

func TestCacheIsBounded(t *testing.T) {
	const maxSizeMB = 8

	bc, err := bigcache.New(context.Background(), cacheConfig(maxSizeMB))
	if err != nil {
		t.Fatal(err)
	}

	entry := bytes.Repeat([]byte("x"), 1024)
	for i := range 4 * maxSizeMB * 1024 {
		if err := bc.Set(fmt.Sprintf("/song-%d-lyrics", i), entry); err != nil {
			t.Fatal(err)
		}
	}

	if got, limit := bc.Capacity(), maxSizeMB<<20; got > limit {
		t.Fatalf("expected cache capacity <= %d bytes, got %d", limit, got)
	}
}

func TestCacheStoresLargeEntries(t *testing.T) {
	key := "large-entry"
	value := bytes.Repeat([]byte("x"), 100*1024)

	if err := setCache(key, value); err != nil {
		t.Fatalf("unable to set cache, %q\n", err)
	}

	if _, err := getCache[[]byte](key); err != nil {
		t.Fatalf("unable to get cache, %q\n", err)
	}
}
