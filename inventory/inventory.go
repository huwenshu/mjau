// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

// Package inventory implements routines for font indexing and querying.
package inventory

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/noll/mjau/font"
	"github.com/noll/mjau/util"

	"github.com/noll/samling/table"
)

// Inventory represents a table for storing fonts.
type Inventory struct {
	*table.Table
}

// Query represents an inventory query.
type Query struct {
	RowKey    string
	ColumnKey string
}

// Build builds the inventory using the JSON-encoded metadata files from the
// first level subdirectories of the named directory. Returns an error if the
// named directory is not a directory, or if it cannot be read.
func (i *Inventory) Build(name string) error {
	if !util.IsDir(name) {
		return fmt.Errorf("%s: not a directory", name)
	}
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		mjson := filepath.Join(name, entry.Name(), "metadata.json")
		if !(entry.IsDir() && util.Exists(mjson)) {
			// No metadata file, skip entry.
			// TODO: Add logging.
			continue
		}
		metadata := new(font.Metadata)
		if err := metadata.Read(mjson); err != nil {
			// Invalid metadata file, skip font family.
			// TODO: Add logging.
			continue
		}
		fonts := metadata.Fonts()
		for _, font := range fonts {
			format := font.Format.String()
			weight := strconv.Itoa(font.Weight)
			columnKey := format + weight + font.Style
			i.Put(font.Family, columnKey, font)
		}
	}
	return nil
}

// Query queries the inventory and returns the font which conforms to the
// given query, or nil if there is no such font in the inventory.
func (i *Inventory) Query(query Query) *font.Font {
	if f := i.Get(query.RowKey, query.ColumnKey); f != nil {
		return f.(*font.Font)
	}
	return nil
}

// New creates and returns a new (empty) inventory.
func New() *Inventory {
	return &Inventory{table.New()}
}
