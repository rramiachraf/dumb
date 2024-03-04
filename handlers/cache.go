package handlers

import (
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rramiachraf/dumb/data"
)

var c, _ = bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 24))

func setCache(key string, entry interface{}) error {
	data, err := json.Marshal(&entry)
	if err != nil {
		return err
	}

	return c.Set(key, data)
}

func getCache[v data.Album | data.Song | data.Annotation](key string) (v, error) {
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
