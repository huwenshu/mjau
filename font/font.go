// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

// Package font implements routines for manipulating font files.
package font

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/noll/mjau/util"
)

const (
	NOF Format = iota // No format.
	// Supported font file formats.
	EOT
	WOFF
)

// Font represents a single font file.
type Font struct {
	Family string
	Format Format
	Path   string
	Style  string
	Weight int
}

// Format represents the format of a font file.
type Format int

// Metadata represents the metadata of a font family.
type Metadata struct {
	Family      string
	Subfamilies []Subfamily
	path        string
}

// Subfamily represents the metadata of a font subfamily.
type Subfamily struct {
	Basename string
	Formats  []string
	Style    string
	Weight   int
}

// Contents reads and returns the contents of the font file.
// Returns an error if the font file cannot be read.
func (f *Font) Contents() ([]byte, error) {
	b, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}
	return b, err
}

// MimeType returns the MIME type associated with the font file.
func (f *Font) MimeType() string {
	switch f.Format {
	case EOT:
		return "application/vnd.ms-fontobject"
	case WOFF:
		return "application/x-font-woff"
	}
	// Should not happen.
	return ""
}

// ModTime returns the modification time of the font.
// If there is an error, it will be of type *os.PathError.
func (f *Font) ModTime() (time.Time, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return time.Time{}, err
	}
	return fi.ModTime(), nil
}

// Equal reports whether the value pointed to by f and
// the v value are the one and same font format.
func (f *Format) Equal(v Format) bool {
	return *f == v
}

// FontFormat determines f to point to the font format corresponding to the
// given font format string representation. If the given font format string
// representation is not valid, f remains unchanged.
func (f *Format) FromString(format string) {
	switch strings.ToLower(format) {
	case "eot":
		*f = EOT
	case "woff":
		*f = WOFF
	}
}

// String returns the string representation of the font format.
// Calling the method on NOF returns the empty string.
func (f *Format) String() string {
	switch *f {
	case EOT:
		return "eot"
	case WOFF:
		return "woff"
	}
	return ""
}

// Fonts returns all the valid fonts corresponding to the metadata.
// Must be used only after reading the metadata from a JSON-encoded file.
func (m *Metadata) Fonts() (fonts []*Font) {
	for _, s := range m.Subfamilies {
		for _, f := range s.Formats {
			format := NOF
			format.FromString(f)
			if format == NOF {
				// Unsupported font file format, skip it.
				// TODO: Add logging.
				continue
			}
			dir := filepath.Dir(m.path)
			filename := s.Basename + "." + f
			path := filepath.Join(dir, filename)
			font := &Font{
				Family: m.Family,
				Format: format,
				Path:   path,
				Style:  s.Style,
				Weight: s.Weight,
			}
			fonts = append(fonts, font)
		}
	}
	return
}

// Read reads and parses the JSON-encoded contents of the named metadata file.
// Returns an error if the named file cannot be read or correctly parsed.
func (m *Metadata) Read(name string) error {
	if err := util.ReadJson(name, &m); err != nil {
		return err
	}
	m.path = name
	return nil
}
