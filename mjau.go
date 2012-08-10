// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	ihttp "github.com/noll/mjau/http" // Internal http package.
	"github.com/noll/mjau/inventory"
	"github.com/noll/mjau/util"
	"github.com/noll/mjau/whitelist"
)

const (
	ProgName    = "mjau"
	ProgVersion = "0.1"
)

var (
	bFlag = flag.String("b", "0.0.0.0:80", "TCP address to bind to")
	cFlag = flag.Int("c", 2592000, "Cache-Control max-age value")
	eFlag = flag.Bool("e", false, "toggle entity tags validation")
	gFlag = flag.Bool("g", false, "toggle response gzip compression")
	lFlag = flag.String("l", "fonts/", "path to font library")
	oFlag = flag.Bool("o", false, "toggle cross-origin resource sharing")
	tFlag = flag.String("t", "templates/", "path to templates directory")
	vFlag = flag.Bool("v", false, "display version number and exit")
	wFlag = flag.String("w", "whitelist.json", "path to whitelist file")
)

func init() {
	util.BlankStrFlagDefault(bFlag, "b")
	util.BlankStrFlagDefault(lFlag, "l")
	util.BlankStrFlagDefault(tFlag, "t")
	util.BlankStrFlagDefault(wFlag, "w")
	*lFlag = filepath.FromSlash(*lFlag)
	*wFlag = filepath.FromSlash(*wFlag)
}

func main() {
	flag.Parse()
	if *vFlag {
		fmt.Println(ProgName, ProgVersion)
		os.Exit(0)
	}
	// Build font inventory.
	fontInventory := inventory.New()
	if err := fontInventory.Build(*lFlag); err != nil {
		PrintErrorExit(err.Error())
	}
	if fontInventory.Size() == 0 {
		PrintErrorExit(fmt.Sprintf("%s: empty font library", *lFlag))
	}
	// Read whitelist.
	whitelist := whitelist.New()
	if err := whitelist.Read(*wFlag); err != nil {
		PrintErrorExit(err.Error())
	}
	if whitelist.Size() == 0 {
		PrintErrorExit(fmt.Sprintf("%s: empty whitelist", *wFlag))
	}
	// Parse templates.
	templatesPath := filepath.FromSlash(*tFlag)
	eot := filepath.Join(templatesPath, "eot.css")
	woff := filepath.Join(templatesPath, "woff.css")
	var templates *template.Template
	var err error
	if templates, err = template.ParseFiles(eot, woff); err != nil {
		PrintErrorExit(err.Error())
	}
	// Create CSS handler function.
	var cssHandler http.HandlerFunc
	ctx := ihttp.HandlerContext{
		Flags: ihttp.Flags{
			AcAllowOrigin: *oFlag,
			CcMaxAge:      *cFlag,
			Etag:          *eFlag,
			Gzip:          *gFlag,
			Version:       ProgName + "/" + ProgVersion,
		},
		Inventory: *fontInventory,
		Templates: *templates,
		Whitelist: *whitelist,
	}
	cssHandler = ihttp.MakeHandler(ihttp.CssHandler, ctx)
	if *gFlag {
		// Enable response gzip compression.
		cssHandler = ihttp.MakeGzipHandler(cssHandler)
	}
	// Register CSS HTTP handler.
	http.HandleFunc("/css/", cssHandler)
	// Start HTTP server.
	if err := http.ListenAndServe(*bFlag, nil); err != nil {
		PrintErrorExit(err.Error())
	}
}
