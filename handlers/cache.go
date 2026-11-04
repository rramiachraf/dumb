package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rramiachraf/dumb/data"
)

type cachable interface {
	data.Album | data.Song | data.Annotation | data.Artist | data.Article | []byte
}

const cacheMaxSizeMB = 256

func cacheConfig(maxSizeMB int) bigcache.Config {
	config := bigcache.DefaultConfig(time.Hour * 24)
	config.HardMaxCacheSize = maxSizeMB
	// The default preallocates more than each shard's share of the cap.
	config.MaxEntriesInWindow = config.Shards * 10

	return config
}

var c, _ = bigcache.New(context.Background(), cacheConfig(cacheMaxSizeMB))

func setCache(key string, entry interface{}) error {
	data, err := json.Marshal(&entry)
	if err != nil {
		return err
	}

	return c.Set(key, data)
}

func getCache[v cachable](key string) (v, error) {
	var decoded v

	data, err := c.Get(key)
	if err != nil {
		return decoded, err
	}

	if err = json.Unmarshal(data, &decoded); err != nil {
		return decoded, err
	}

	return decoded, nil
}
