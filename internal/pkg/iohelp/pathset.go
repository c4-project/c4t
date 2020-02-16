package iohelp

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

// ErrPathsetNil is a standard error for when things that expect a pathset don't get one.
var ErrPathsetNil = errors.New("pathset nil")

// Mkdirs tries to make each directory in the directory set d.
func Mkdirs(ps ...string) error {
	for _, dir := range ps {
		logrus.Debugf("mkdir %s\n", dir)
		if err := os.MkdirAll(dir, 0744); err != nil {
			return err
		}
	}
	return nil
}
