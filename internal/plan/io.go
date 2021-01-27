// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/c4-project/c4t/internal/helper/errhelp"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// WriteFlag is the type of writing flags for use in
type WriteFlag uint8

const (
	// Ext is the file extension, if any, that should be used for plan files.
	Ext = ".json"
	// ExtCompress is the file extension, if any, that should be used for compressed plan files.
	ExtCompress = Ext + ".gz"

	// WriteNone is the absence of write flags.
	WriteNone WriteFlag = 0
	// WriteHuman should be passed to plan writers to request human-readable indentation.
	WriteHuman WriteFlag = 1 << iota
	// WriteCompress should be passed to plan writers to request compression.
	WriteCompress
)

// Read reads uncompressed plan information from r into p.
func Read(r io.Reader, p *Plan) error {
	return json.NewDecoder(r).Decode(p)
}

// ReadCompressed reads compressed plan information from r into p.
func ReadCompressed(r io.Reader, p *Plan) error {
	rc, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	rerr := Read(rc, p)
	cerr := rc.Close()
	return errhelp.FirstError(rerr, cerr)
}

func hasGzMagic(r io.ReadSeeker) (bool, error) {
	var bs [2]byte
	n, err := io.ReadFull(r, bs[:])
	if n != 2 {
		err = fmt.Errorf("wanted 2 bytes, only got %d", n)
	}
	if err != nil {
		return false, err
	}
	isGz := bs[0] == 0x1F && bs[1] == 0x8B
	isGz = isGz || (bs[1] == 0x1F && bs[0] == 0x8B)
	_, err = r.Seek(-2, io.SeekCurrent)
	return isGz, err
}

// ReadMagic is ReadCompressed if r starts with the gzip magic number, and Read otherwise.
// As we seek backwards after reading the magic number, r must be a ReadSeeker.
func ReadMagic(r io.ReadSeeker, p *Plan) error {
	isGz, err := hasGzMagic(r)
	if err != nil {
		return err
	}
	if isGz {
		return ReadCompressed(r, p)
	}
	return Read(r, p)
}

// ReadFile reads plan information from the file named by path into p.
func ReadFile(path string, p *Plan) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening plan file %q: %w", path, err)
	}
	perr := ReadMagic(r, p)
	cerr := r.Close()
	return errhelp.FirstError(perr, cerr)
}

// wrapForCompress wraps w according to whether flags requires it to compress.
// The result will be a gzip writer that can be closed if so, and a nop-closer over w otherwise.
func wrapForCompress(w io.Writer, flags WriteFlag) io.WriteCloser {
	if flags&WriteCompress != 0 {
		return gzip.NewWriter(w)
	}
	return iohelp.NopWriteCloser{Writer: w}
}

// Write dumps plan p to w.
func (p *Plan) Write(w io.Writer, flags WriteFlag) error {
	wc := wrapForCompress(w, flags)
	enc := json.NewEncoder(wc)
	if flags&WriteHuman != 0 {
		enc.SetIndent("", "\t")
	}
	err := enc.Encode(p)
	// We need to close in case we're compressing; else, the footer won't be written.
	cerr := wc.Close()
	return errhelp.FirstError(err, cerr)
}

// WriteFile dumps plan p to the file named by path.
func (p *Plan) WriteFile(path string, flags WriteFlag) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating plan file: %w", err)
	}
	err = p.Write(f, flags)
	cerr := f.Close()
	return errhelp.FirstError(err, cerr)
}
