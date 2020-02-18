package iohelp

import (
	"path"
	"strings"
)

// ExtlessFile gets the file part of fpath without its extension.
func ExtlessFile(fpath string) string {
	_, file := path.Split(fpath)
	ext := path.Ext(file)
	return strings.TrimSuffix(file, ext)
}
