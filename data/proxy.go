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

	return fmt.Sprintf("/images%s", u.Path)
}

