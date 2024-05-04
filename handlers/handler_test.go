package handlers

import (
	"io/fs"
	"os"
	"path"
)

type assets struct{}

func (assets) Open(p string) (fs.File, error) {
	return os.Open(path.Join("../", p))
}
