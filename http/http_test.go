// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package http

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/noll/mjau/font"
	"github.com/noll/mjau/inventory"
	"github.com/noll/mjau/test"
	"github.com/noll/mjau/util"
	"github.com/noll/mjau/whitelist"
)

var (
	amf = filepath.Join(fl, "Amaranth")      // Amaranth family path.
	fl  = filepath.FromSlash("../fonts")     // Font library path.
	tp  = filepath.FromSlash("../templates") // Templates path.
)

var QueriesCases = []struct {
	Family  string // Family form value.
	Format  font.Format
	Queries []*inventory.Query
}{
	// Case 1
	{
		Family:      "Amaranth",
		font.Format: font.EOT,
		Queries: []*inventory.Query{
			&inventory.Query{"Amaranth", "eot400normal"},
		},
	},
	// Case 2
	{
		Family:      "Amaranth|Open+Sans",
		font.Format: font.WOFF,
		Queries: []*inventory.Query{
			&inventory.Query{"Amaranth", "woff400normal"},
			&inventory.Query{"Open+Sans", "woff400normal"},
		},
	},
	// Case 3
	{
		Family:      "Amaranth:700italic|Open+Sans:800normal",
		font.Format: font.EOT,
		Queries: []*inventory.Query{
			&inventory.Query{"Amaranth", "eot700italic"},
			&inventory.Query{"Open+Sans", "eot800normal"},
		},
	},
	// Case 4
	{
		Family:      "Amaranth:400normal|Open+Sans:300normal,600italic",
		font.Format: font.WOFF,
		Queries: []*inventory.Query{
			&inventory.Query{"Amaranth", "woff400normal"},
			&inventory.Query{"Open+Sans", "woff300normal"},
			&inventory.Query{"Open+Sans", "woff600italic"},
		},
	},
	// Case 5
	{
		Family:      "Amaranth:400normal,700normal|Open+Sans:700normal",
		font.Format: font.EOT,
		Queries: []*inventory.Query{
			&inventory.Query{"Amaranth", "eot400normal"},
			&inventory.Query{"Amaranth", "eot700normal"},
			&inventory.Query{"Open+Sans", "eot700normal"},
		},
	},
	// Case 6
	{
		Family:      "Amaranth:400,700|Open+Sans:700",
		font.Format: font.EOT,
		Queries: []*inventory.Query{
			&inventory.Query{"Amaranth", "eot400normal"},
			&inventory.Query{"Amaranth", "eot700normal"},
			&inventory.Query{"Open+Sans", "eot700normal"},
		},
	},
	// Case 7
	{
		Family:      "|Amaranth",
		font.Format: font.WOFF,
	},
	// Case 8
	{
		Family:      ":Amaranth",
		font.Format: font.WOFF,
	},
	// Case 9
	{
		Family:      ",Amaranth",
		font.Format: font.WOFF,
	},
	// Case 10
	{
		Family:      "|Open+Sans:700,300italic",
		font.Format: font.WOFF,
	},
	// Case 11
	{
		Family:      ":Open+Sans:700,300italic",
		font.Format: font.WOFF,
	},
	// Case 12
	{
		Family:      ",Open+Sans:700,300italic",
		font.Format: font.WOFF,
	},
	// Case 13
	{
		Family:      "Open+Sans:700,300italic|",
		font.Format: font.WOFF,
		Queries: []*inventory.Query{
			&inventory.Query{"Open+Sans", "woff700normal"},
			&inventory.Query{"Open+Sans", "woff300italic"},
		},
	},
	// Case 14
	{
		Family:      "Open+Sans:700,300italic||Amaranth",
		font.Format: font.EOT,
		Queries: []*inventory.Query{
			&inventory.Query{"Open+Sans", "eot700normal"},
			&inventory.Query{"Open+Sans", "eot300italic"},
			&inventory.Query{"Amaranth", "eot400normal"},
		},
	},
	// Case 15
	{
		Family:      "Open+Sans:700,300italic,",
		font.Format: font.EOT,
		Queries: []*inventory.Query{
			&inventory.Query{"Open+Sans", "eot700normal"},
			&inventory.Query{"Open+Sans", "eot300italic"},
		},
	},
	// Case 16
	{
		Family:      "Open+Sans:700,300italic,|Amaranth",
		font.Format: font.WOFF,
		Queries: []*inventory.Query{
			&inventory.Query{"Open+Sans", "woff700normal"},
			&inventory.Query{"Open+Sans", "woff300italic"},
			&inventory.Query{"Amaranth", "woff400normal"},
		},
	},
	// Case 17
	{
		Family:      "Open+Sans:700,300italic,,",
		font.Format: font.WOFF,
		Queries: []*inventory.Query{
			&inventory.Query{"Open+Sans", "woff700normal"},
			&inventory.Query{"Open+Sans", "woff300italic"},
		},
	},
	// Case 18
	{
		Family:      "Open+Sans:700,300italic,,|Amaranth",
		font.Format: font.EOT,
		Queries: []*inventory.Query{
			&inventory.Query{"Open+Sans", "eot700normal"},
			&inventory.Query{"Open+Sans", "eot300italic"},
			&inventory.Query{"Amaranth", "eot400normal"},
		},
	},
	// Case 19
	{
		Family:      "Open+Sans:700,300italic,,400|Amaranth",
		font.Format: font.WOFF,
		Queries: []*inventory.Query{
			&inventory.Query{"Open+Sans", "woff700normal"},
			&inventory.Query{"Open+Sans", "woff300italic"},
			&inventory.Query{"Open+Sans", "woff400normal"},
			&inventory.Query{"Amaranth", "woff400normal"},
		},
	},
}

func TestFontFaceFromFont(t *testing.T) {
	ff := &FontFace{}
	f := font.Font{
		Family: "Amaranth",
		Format: font.WOFF,
		Path:   filepath.Join(amf, "amaranth-regular.woff"),
		Style:  "normal",
		Weight: 400,
	}
	err := ff.FromFont(f)
	test.VerifyFatal(t, 1, 0, true, nil == err)
	data, err := f.Contents()
	test.Verify(t, 2, 0, util.Base64(data), ff.Base64Data)
	test.Verify(t, 3, 0, f.Family, ff.Family)
	test.Verify(t, 4, 0, f.Format.String(), ff.Format)
	wMimeType := f.MimeType()
	test.VerifyFatal(t, 5, 0, false, "" == wMimeType)
	test.Verify(t, 6, 0, wMimeType, ff.MimeType)
	test.Verify(t, 7, 0, f.Style, ff.Style)
	test.Verify(t, 8, 0, f.Weight, ff.Weight)
}

func TestMakeHandler(t *testing.T) {
	var cases = []struct {
		Context HandlerContext
		Header  map[string]string
	}{
		// Case 1
		{
			Context: HandlerContext{
				Flags: Flags{
					AcAllowOrigin: false,
					Gzip:          true,
					Version:       "test/0.1",
				},
			},
			Header: map[string]string{
				"Server": "test/0.1",
				"Vary":   "Accept-Encoding",
			},
		},
		// Case 2
		{
			Context: HandlerContext{
				Flags: Flags{
					AcAllowOrigin: true,
					Gzip:          false,
					Version:       "test/0.1",
				},
			},
			Header: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Server":                      "test/0.1",
			},
		},
		// Case 3
		{
			Context: HandlerContext{
				Flags: Flags{
					AcAllowOrigin: false,
					Gzip:          true,
					Version:       "",
				},
			},
			Header: map[string]string{
				"Vary": "Accept-Encoding",
			},
		},
		// Case 4
		{
			Context: HandlerContext{
				Flags: Flags{},
			},
			Header: map[string]string{},
		},
	}

	for i, c := range cases {
		j := i + 1

		handler := MakeHandler(emptyHandler, c.Context)
		server := httptest.NewServer(handler)
		defer server.Close()

		client := http.Client{}
		resp, err := client.Get(server.URL)
		defer resp.Body.Close()
		test.VerifyFatal(t, 1, j, true, nil == err)

		wHeader := c.Header
		gHeader := resp.Header
		test.Verify(t, 2, j, wHeader["Server"], gHeader.Get("Server"))
		test.Verify(t, 3, j, wHeader["Vary"], gHeader.Get("Vary"))
		wAcAllowOrigin := wHeader["Access-Control-Allow-Origin"]
		gAcAllowOrigin := gHeader.Get("Access-Control-Allow-Origin")
		test.Verify(t, 4, j, wAcAllowOrigin, gAcAllowOrigin)
	}
}

func TestCssHandler(t *testing.T) {
	type Request struct {
		Method string
		URL    string
	}

	// Build an inventory.
	// Used in cases 6-7.
	inv := inventory.New()
	err := inv.Build(fl)
	test.VerifyFatal(t, 1, 0, true, err == nil)

	// Build an "allow all" whitelist.
	// Use this whitelist to allow all
	// referrers to fetch the resource.
	// Used in cases 3-7.
	aawl := whitelist.New()
	aawl.Domains = append(aawl.Domains, "")

	// Parse templates.
	// Used in cases 7-9.
	eot := filepath.Join(tp, "eot.css")
	woff := filepath.Join(tp, "woff.css")
	tmpl, err := template.ParseFiles(eot, woff)
	test.VerifyFatal(t, 2, 0, true, err == nil)

	// Execute template containing Amaranth Regular.
	// Used in case 7.
	ar := &font.Font{
		Family: "Amaranth",
		Format: font.WOFF,
		Path:   filepath.Join(amf, "amaranth-regular.woff"),
		Style:  "normal",
		Weight: 400,
	}
	var tmplData []*FontFace
	data, err := ar.Contents()
	test.VerifyFatal(t, 3, 0, true, err == nil)
	mimeType := ar.MimeType()
	test.VerifyFatal(t, 4, 0, false, "" == mimeType)
	arff := &FontFace{
		Base64Data: util.Base64(data),
		Family:     ar.Family,
		Format:     ar.Format.String(),
		MimeType:   mimeType,
		Style:      ar.Style,
		Weight:     ar.Weight,
	}
	buf := new(bytes.Buffer)
	tmplName := ar.Format.String() + ".css"
	tmplData = append(tmplData, arff)
	err = tmpl.ExecuteTemplate(buf, tmplName, tmplData)
	test.VerifyFatal(t, 5, 0, true, nil == err)
	arBody := buf.Bytes()

	// Generate the entity tag corresponding to Amaranth Regular.
	// Uncompressed response.
	// Used in case 7.
	hash := md5.New()
	modtime, err := ar.ModTime()
	test.VerifyFatal(t, 6, 0, true, nil == err)
	io.WriteString(hash, modtime.String())
	arEtag := fmt.Sprintf("%x", hash.Sum(nil))
	// Gzip compressed response.
	// Used in cases 8-9.
	arEtagGzip := arEtag + "+gzip"

	var cases = []struct {
		Body        []byte
		Context     HandlerContext
		Header      map[string]string
		IfNoneMatch string // If-None-Match client request header.
		Request     *Request
		StatusCode  int
	}{
		// Case 1
		{
			Header: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			Request: &Request{
				Method: "POST",
				URL:    "",
			},
			StatusCode: http.StatusNotImplemented,
		},
		// Case 2
		{
			Header: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			Request: &Request{
				Method: "GET",
				URL:    "",
			},
			StatusCode: http.StatusForbidden,
		},
		// Case 3
		{
			Context: HandlerContext{
				Whitelist: *aawl,
			},
			Header: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			Request: &Request{
				Method: "GET",
				URL:    "?family=",
			},
			StatusCode: http.StatusBadRequest,
		},
		// Case 4
		{
			Context: HandlerContext{
				Whitelist: *aawl,
			},
			Header: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			Request: &Request{
				Method: "GET",
				URL:    "?family=Amaranth&format=nonexistent",
			},
			StatusCode: http.StatusBadRequest,
		},
		// Case 5
		{
			Context: HandlerContext{
				Whitelist: *aawl,
			},
			Header: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			Request: &Request{
				Method: "GET",
				URL:    "?family=|Amaranth",
			},
			StatusCode: http.StatusBadRequest,
		},
		// Case 6
		{
			Context: HandlerContext{
				Flags: Flags{
					Etag: true,
				},
				Inventory: *inv,
				Whitelist: *aawl,
			},
			Header: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
				"Etag":         "",
			},
			Request: &Request{
				Method: "GET",
				URL:    "?family=Nonexistent",
			},
			StatusCode: http.StatusBadRequest,
		},
		// Case 7
		{
			Body: arBody,
			Context: HandlerContext{
				Flags: Flags{
					CcMaxAge: 2592000,
					Etag:     true,
				},
				Inventory: *inv,
				Templates: *tmpl,
				Whitelist: *aawl,
			},
			Header: map[string]string{
				"Cache-Control": "max-age=2592000",
				"Content-Type":  "text/css; charset=utf-8",
				"Etag":          arEtag,
			},
			Request: &Request{
				Method: "GET",
				URL:    "?family=Amaranth",
			},
			StatusCode: http.StatusOK,
		},
		// Case 8
		{
			Body: arBody,
			Context: HandlerContext{
				Flags: Flags{
					Etag: true,
					Gzip: true,
				},
				Inventory: *inv,
				Templates: *tmpl,
				Whitelist: *aawl,
			},
			Header: map[string]string{
				"Cache-Control": "max-age=0",
				"Content-Type":  "text/css; charset=utf-8",
				"Etag":          arEtagGzip,
			},
			Request: &Request{
				Method: "GET",
				URL:    "?family=Amaranth",
			},
			StatusCode: http.StatusOK,
		},
		// Case 9
		{
			Context: HandlerContext{
				Flags: Flags{
					Etag: true,
					Gzip: true,
				},
				Inventory: *inv,
				Templates: *tmpl,
				Whitelist: *aawl,
			},
			Header: map[string]string{
				"Cache-Control": "max-age=0",
				"Etag":          arEtagGzip,
			},
			IfNoneMatch: arEtagGzip,
			Request: &Request{
				Method: "GET",
				URL:    "?family=Amaranth",
			},
			StatusCode: http.StatusNotModified,
		},
	}

	for i, c := range cases {
		j := i + 1

		handler := MakeHandler(CssHandler, c.Context)
		server := httptest.NewServer(handler)
		defer server.Close()

		client := http.Client{}
		reqURL := server.URL + c.Request.URL
		req, err := http.NewRequest(c.Request.Method, reqURL, nil)
		test.VerifyFatal(t, 6, j, true, nil == err)

		if c.IfNoneMatch != "" {
			req.Header.Add("If-None-Match", c.IfNoneMatch)
		}

		resp, err := client.Do(req)
		defer resp.Body.Close()
		test.VerifyFatal(t, 7, j, true, nil == err)
		test.Verify(t, 8, j, c.StatusCode, resp.StatusCode)

		wContentType := c.Header["Content-Type"]
		gContentType := resp.Header.Get("Content-Type")
		test.Verify(t, 9, j, wContentType, gContentType)

		wCacheControl := c.Header["Cache-Control"]
		gCacheControl := resp.Header.Get("Cache-Control")
		test.Verify(t, 10, j, wCacheControl, gCacheControl)

		wEtag := c.Header["Etag"]
		gEtag := resp.Header.Get("Etag")
		test.Verify(t, 11, j, wEtag, gEtag)

		if (c.StatusCode == http.StatusOK) &&
			(resp.StatusCode == http.StatusOK) {
			gbody, err := ioutil.ReadAll(resp.Body)
			test.VerifyFatal(t, 13, j, true, nil == err)
			wbody := c.Body
			test.Verify(t, 14, j, true, bytes.Equal(wbody, gbody))
		}
	}
}

func TestQueries(t *testing.T) {
	for i, c := range QueriesCases {
		j := i + 1

		gqueries := Queries(c.Family, c.Format)
		gqueriesSize := len(gqueries)
		wqueriesSize := len(c.Queries)
		test.VerifyFatal(t, 1, j, wqueriesSize, gqueriesSize)

		for k, wquery := range c.Queries {
			gquery := gqueries[k]
			test.Verify(t, 2, j, wquery.RowKey, gquery.RowKey)
			test.Verify(t, 3, j, gquery.ColumnKey, gquery.ColumnKey)
		}
	}
}

func emptyHandler(w http.ResponseWriter, r *http.Request, ctx HandlerContext) {}
