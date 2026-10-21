package data

import (
	"fmt"
	"net/url"
)

func ExtractImageURL(image string) string {
	u, err := url.Parse(image)
	if err != nil {
		return ""
	}

	// The image proxy can't fetch from assets.genius.com.
	if u.Host == "assets.genius.com" && u.Path == "/images/default_cover_image.png" {
		return "/static/default_cover_image.png"
	}

	return fmt.Sprintf("/images%s", u.Path)
}
