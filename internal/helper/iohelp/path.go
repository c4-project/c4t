// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// ExtlessFile gets the file part of slash-path fpath without its extension.
func ExtlessFile(fpath string) string {
	_, file := path.Split(fpath)
	ext := path.Ext(file)
	return strings.TrimSuffix(file, ext)
}

// ExpandMany applies homedir.Expand to every (file)path in paths.
func ExpandMany(paths []string) ([]string, error) {
	var err error
	xpaths := make([]string, len(paths))
	for i, p := range paths {
		xpaths[i], err = homedir.Expand(p)
		if err != nil {
			return nil, err
		}
	}
	return xpaths, nil
}
