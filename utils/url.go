package utils

import (
	"net/url"
	"strings"
)

func TrimURL(u string) string {
	uu, err := url.Parse(u)
	if err != nil {
		return ""
	}

	if strings.HasPrefix(uu.Path, "/") {
		return uu.Path
	}

	return "/" + uu.Path
}
