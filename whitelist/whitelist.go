// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

// Package whitelist implements a JSON-driven domain names whitelist.
package whitelist

import (
	"fmt"
	"strings"

	"github.com/noll/mjau/util"
)

// Whitelist is the representation of a JSON-encoded whitelist file.
type Whitelist struct {
	Domains []string
}

// Contains reports whether the given domain name is present in the whitelist.
// Returns true if the given domain name has as a prefix one of the domain
// names in the whitelist.
func (w *Whitelist) Contains(domain string) (y bool) {
	for _, d := range w.Domains {
		if strings.HasPrefix(domain, d) {
			y = true
		}
	}
	return
}

// Read reads and parses the JSON-encoded contents of the named file and stores
// the result in the whitelist.
// Returns an error if the named file cannot be read or correctly parsed.
func (w *Whitelist) Read(name string) error {
	if util.IsDir(name) {
		return fmt.Errorf("%s: is a directory", name)
	}
	if err := util.ReadJson(name, &w); err != nil {
		return err
	}
	return nil
}

// Size returns the number of domain names in the whitelist.
func (w *Whitelist) Size() int {
	return len(w.Domains)
}

// New creates and returns a new (empty) whitelist.
func New() *Whitelist {
	return &Whitelist{}
}
