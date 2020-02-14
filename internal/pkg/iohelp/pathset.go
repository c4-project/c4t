package iohelp

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

// ErrPathsetNil is a standard error for when things that expect a pathset don't get one.
var ErrPathsetNil = errors.New("pathset nil")

// Pathset represents any type that contains a list of directories.
type Pathset interface {
	// Dirs gets the list of directories mentioned in this directory set.
	Dirs() []string
}

// Mkdirs tries to make each directory in the directory set d.
func Mkdirs(d Pathset) error {
	for _, dir := range d.Dirs() {
		logrus.Debugf("mkdir %s\n", dir)
		if err := os.MkdirAll(dir, 0744); err != nil {
			return err
		}
	}
	return nil
}
