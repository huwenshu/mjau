// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

// Package http implements routines for serving @font-face CSS files over HTTP.
package http

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/noll/mjau/font"
	"github.com/noll/mjau/inventory"
	"github.com/noll/mjau/util"
	"github.com/noll/mjau/whitelist"
)

// Flags holds the state of the HTTP flags.
type Flags struct {
	AcAllowOrigin bool   // Cross-origin resource sharing toggle.
	CcMaxAge      uint64 // Cache-Control max-age value.
	Etag          bool   // Entity tags validation toggle.
	Gzip          bool   // Response gzip compression toggle.
	Version       string // Server version string.
}

// FontFace represents a single @font-face CSS rule.
type FontFace struct {
	Base64Data string
	Family     string
	Format     string
	MimeType   string
	Style      string
	Weight     int
}

// HandlerContext represents a context for a handler.
type HandlerContext struct {
	Flags     Flags
	Inventory inventory.Inventory
	Templates template.Template
	Whitelist whitelist.Whitelist
}

type HandlerFunc func(http.ResponseWriter, *http.Request, HandlerContext)

// FromFont initializes ff from the given font. If it is called on an already
// initialized font face, changes the font face according to the given font.
// Returns an error if the initialization fails.
func (ff *FontFace) FromFont(f font.Font) error {
	data, err := f.Contents()
	if err != nil {
		return err
	}
	mimeType := f.MimeType()
	if mimeType == "" {
		// Should not happen.
		return fmt.Errorf("can't determine font MIME type")
	}
	ff.Base64Data = util.Base64(data)
	ff.Family = f.Family
	ff.Format = f.Format.String()
	ff.MimeType = mimeType
	ff.Style = f.Style
	ff.Weight = f.Weight
	return nil
}

func CssHandler(w http.ResponseWriter, r *http.Request, ctx HandlerContext) {
	if r.Method != "GET" {
		// TODO: Add logging.
		NotImplemented(w, r)
		return
	}
	// Allow only whitelisted referers to fetch the resource.
	if !ctx.Whitelist.Contains(r.Referer()) {
		// TODO: Add logging.
		Forbidden(w, r)
		return
	}
	family := r.FormValue("family")
	if family == "" {
		// TODO: Add logging.
		BadRequest(w, r)
		return
	}
	sFormat := r.FormValue("format")
	if sFormat == "" {
		// Font format not specified,
		// default to WOFF.
		sFormat = "woff"
	}
	format := font.NOF
	format.FromString(sFormat)
	if format == font.NOF {
		// TODO: Add logging.
		BadRequest(w, r)
		return
	}
	queries := Queries(family, format)
	if len(queries) == 0 {
		// TODO: Add logging.
		BadRequest(w, r)
		return
	}
	if ctx.Flags.Etag && Etag(w, r, queries, ctx) {
		return
	}
	var templateData []*FontFace
	var templateName string
	for _, query := range queries {
		fnt := ctx.Inventory.Query(*query)
		if fnt == nil {
			// TODO: Add logging.
			BadRequest(w, r)
			return
		}
		fontFace := new(FontFace)
		err := fontFace.FromFont(*fnt)
		if err != nil {
			// TODO: Add logging.
			InternalServerError(w, r)
			return
		}
		templateData = append(templateData, fontFace)
		templateName = fnt.Format.String() + ".css.tmpl"
	}
	buf := new(bytes.Buffer)
	err := ctx.Templates.ExecuteTemplate(buf, templateName, templateData)
	if err != nil {
		// TODO: Add logging.
		InternalServerError(w, r)
		return
	}
	maxAge := strconv.FormatUint(ctx.Flags.CcMaxAge, 10)
	w.Header().Set("Cache-Control", "max-age="+maxAge)
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	io.Copy(w, buf)
}

// Etag generates and validates entity tags.
// Returns true if the resource has not been modified.
func Etag(w http.ResponseWriter, r *http.Request, queries []*inventory.Query,
	ctx HandlerContext) bool {
	var failed bool
	hash := md5.New()
	for _, query := range queries {
		fnt := ctx.Inventory.Query(*query)
		if fnt == nil {
			// TODO: Log error.
			failed = true
			break
		}
		modtime, err := fnt.ModTime()
		if err != nil {
			// TODO: Log error.
			failed = true
			break
		}
		io.WriteString(hash, modtime.String())
	}
	if !failed {
		etag := fmt.Sprintf("%x", hash.Sum(nil))
		// Add "+gzip" suffix to entity tag if the response
		// is going to be gzip compressed.
		if ctx.Flags.Gzip {
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				etag = etag + "+gzip"
			}
		}
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			maxAge := strconv.FormatUint(ctx.Flags.CcMaxAge, 10)
			w.Header().Set("Cache-Control", "max-age="+maxAge)
			NotModified(w, r)
			return true
		}
	} else {
		// Something went wrong and the entity tag cannot be
		// validated, make sure the entity tag received from
		// the client is resent.
		etag := r.Header.Get("If-None-Match")
		if etag != "" {
			w.Header().Set("ETag", etag)
		}
	}
	return false
}

func MakeHandler(fn HandlerFunc, ctx HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.Flags.AcAllowOrigin {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if ctx.Flags.Version != "" {
			w.Header().Set("Server", ctx.Flags.Version)
		}
		// Set the "Vary: Accept-Encoding" HTTP response
		// header if gzip compression is enabled.
		if ctx.Flags.Gzip {
			w.Header().Set("Vary", "Accept-Encoding")
		}
		fn(w, r, ctx)
	}
}

// Queries builds and returns a slice of pointers to inventory queries from the
// given family form value and font format format.
func Queries(family string, format font.Format) []*inventory.Query {
	var queries []*inventory.Query
	if strings.HasPrefix(family, "|") || strings.HasPrefix(family, ":") ||
		strings.HasPrefix(family, ",") {
		return queries
	}
	families := strings.Split(family, "|")
	for _, f := range families {
		// Contains font family name at index 0
		// and specified styles at index 1.
		familyStyles := strings.Split(f, ":")
		switch len(familyStyles) {
		case 1:
			q := new(inventory.Query)
			q.RowKey = familyStyles[0]
			// Weight and style are not specified, default
			// weight to 400 and style to normal.
			q.ColumnKey = "400normal"
			if q.RowKey == "" || q.ColumnKey == "" {
				continue
			}
			q.ColumnKey = format.String() + q.ColumnKey
			queries = append(queries, q)
		case 2:
			styles := strings.Split(familyStyles[1], ",")
			for _, s := range styles {
				q := new(inventory.Query)
				q.RowKey = familyStyles[0]
				if _, err := strconv.Atoi(s); err == nil {
					// Only weight is specified,
					// default style to normal.
					q.ColumnKey = s + "normal"
				} else {
					q.ColumnKey = s
				}
				if q.RowKey == "" || q.ColumnKey == "" {
					continue
				}
				q.ColumnKey = format.String() + q.ColumnKey
				queries = append(queries, q)
			}
		}
	}
	return queries
}
