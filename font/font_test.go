// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package font

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/noll/mjau/test"
)

var (
	amf = filepath.Join(fl, "Amaranth")  // Amaranth family path.
	fl  = filepath.FromSlash("../fonts") // Font library path.
)

func TestFontContents(t *testing.T) {
	font := &Font{
		Family: "Amaranth",
		Format: WOFF,
		Path:   filepath.Join(amf, "amaranth-regular.woff"),
		Style:  "normal",
		Weight: 400,
	}
	gContents, err := font.Contents()
	test.VerifyFatal(t, 1, 0, true, nil == err)
	wContents, err := ioutil.ReadFile(font.Path)
	test.VerifyFatal(t, 2, 0, true, nil == err)
	test.Verify(t, 3, 0, 0, bytes.Compare(wContents, gContents))
}

func TestFontMimeType(t *testing.T) {
	font := &Font{
		Family: "Amaranth",
		Format: WOFF,
		Path:   filepath.Join(amf, "amaranth-regular.woff"),
		Style:  "normal",
		Weight: 400,
	}
	gMimeType := font.MimeType()
	wMimeType := "application/x-font-woff"
	test.Verify(t, 1, 0, wMimeType, gMimeType)
}

func TestFontModTime(t *testing.T) {
	font := &Font{
		Family: "Amaranth",
		Format: WOFF,
		Path:   filepath.Join(amf, "amaranth-regular.woff"),
		Style:  "normal",
		Weight: 400,
	}
	gModTime, err := font.ModTime()
	test.VerifyFatal(t, 1, 0, true, nil == err)
	fi, err := os.Stat(font.Path)
	test.VerifyFatal(t, 2, 0, true, nil == err)
	wModTime := fi.ModTime()
	test.Verify(t, 3, 0, true, gModTime.Equal(wModTime))
}

func TestFormatEqual(t *testing.T) {
	format := EOT
	test.Verify(t, 1, 0, true, format.Equal(EOT))
	test.Verify(t, 2, 0, false, format.Equal(WOFF))
}

func TestFormatFromString(t *testing.T) {
	format := NOF
	format.FromString("eot")
	test.Verify(t, 1, 0, EOT, format)
	format.FromString("WOFF")
	test.Verify(t, 2, 0, WOFF, format)
}

func TestFormatString(t *testing.T) {
	format := EOT
	test.Verify(t, 1, 0, "eot", format.String())
	format = NOF
	test.Verify(t, 2, 0, "", format.String())
}

func TestMetadataFonts(t *testing.T) {
	metadata := &Metadata{
		Family: "Amaranth",
		Subfamilies: []Subfamily{
			// Subfamily 1
			Subfamily{
				Basename: "amaranth-regular",
				Formats:  []string{"eot", "woff"},
				Style:    "normal",
				Weight:   400,
			},
			// Subfamily 2
			Subfamily{
				Basename: "amaranth-italic",
				Formats:  []string{"eot", "woff"},
				Style:    "italic",
				Weight:   400,
			},
			// Subfamily 3
			Subfamily{
				Basename: "amaranth-bold",
				Formats:  []string{"eot", "woff"},
				Style:    "normal",
				Weight:   700,
			},
			// Subfamily 4
			Subfamily{
				Basename: "amaranth-bolditalic",
				Formats:  []string{"eot", "woff"},
				Style:    "italic",
				Weight:   700,
			},
		},
		path: filepath.Join(amf, "metadata.json"),
	}
	gFonts := metadata.Fonts()
	wFonts := []*Font{
		// Font 1
		&Font{
			Family: "Amaranth",
			Format: EOT,
			Path:   filepath.Join(amf, "amaranth-regular.eot"),
			Style:  "normal",
			Weight: 400,
		},
		// Font 2
		&Font{
			Family: "Amaranth",
			Format: WOFF,
			Path:   filepath.Join(amf, "amaranth-regular.woff"),
			Style:  "normal",
			Weight: 400,
		},
		// Font 3
		&Font{
			Family: "Amaranth",
			Format: EOT,
			Path:   filepath.Join(amf, "amaranth-italic.eot"),
			Style:  "italic",
			Weight: 400,
		},
		// Font 4
		&Font{
			Family: "Amaranth",
			Format: WOFF,
			Path:   filepath.Join(amf, "amaranth-italic.woff"),
			Style:  "italic",
			Weight: 400,
		},
		// Font 5
		&Font{
			Family: "Amaranth",
			Format: EOT,
			Path:   filepath.Join(amf, "amaranth-bold.eot"),
			Style:  "normal",
			Weight: 700,
		},
		// Font 6
		&Font{
			Family: "Amaranth",
			Format: WOFF,
			Path:   filepath.Join(amf, "amaranth-bold.woff"),
			Style:  "normal",
			Weight: 700,
		},
		// Font 7
		&Font{
			Family: "Amaranth",
			Format: EOT,
			Path:   filepath.Join(amf, "amaranth-bolditalic.eot"),
			Style:  "italic",
			Weight: 700,
		},
		// Font 8
		&Font{
			Family: "Amaranth",
			Format: WOFF,
			Path:   filepath.Join(amf, "amaranth-bolditalic.woff"),
			Style:  "italic",
			Weight: 700,
		},
	}
	test.VerifyFatal(t, 1, 0, len(wFonts), len(gFonts))
	for i, wFont := range wFonts {
		j := i + 1
		gFont := gFonts[i]
		test.Verify(t, 2, j, wFont.Family, gFont.Family)
		test.Verify(t, 3, j, true, wFont.Format.Equal(gFont.Format))
		test.Verify(t, 4, j, wFont.Path, gFont.Path)
		test.Verify(t, 5, j, wFont.Style, gFont.Style)
		test.Verify(t, 6, j, wFont.Weight, gFont.Weight)
	}
}

func TestMetadataRead(t *testing.T) {
	gMetadata := new(Metadata)
	err := gMetadata.Read(filepath.Join(amf, "metadata.json"))
	test.VerifyFatal(t, 1, 0, true, nil == err)
	wMetadata := &Metadata{
		Family: "Amaranth",
		Subfamilies: []Subfamily{
			// Subfamily 1
			Subfamily{
				Basename: "amaranth-regular",
				Formats:  []string{"eot", "woff"},
				Style:    "normal",
				Weight:   400,
			},
			// Subfamily 2
			Subfamily{
				Basename: "amaranth-italic",
				Formats:  []string{"eot", "woff"},
				Style:    "italic",
				Weight:   400,
			},
			// Subfamily 3
			Subfamily{
				Basename: "amaranth-bold",
				Formats:  []string{"eot", "woff"},
				Style:    "normal",
				Weight:   700,
			},
			// Subfamily 4
			Subfamily{
				Basename: "amaranth-bolditalic",
				Formats:  []string{"eot", "woff"},
				Style:    "italic",
				Weight:   700,
			},
		},
		path: filepath.Join(amf, "metadata.json"),
	}
	test.Verify(t, 2, 0, wMetadata.Family, gMetadata.Family)

	gSubfamilies := gMetadata.Subfamilies
	wSubfamilies := wMetadata.Subfamilies
	test.VerifyFatal(t, 3, 0, len(wSubfamilies), len(gSubfamilies))
	for i, wSubfamily := range wMetadata.Subfamilies {
		j := i + 1

		gSubfamily := gMetadata.Subfamilies[i]
		test.Verify(t, 4, j, wSubfamily.Basename, gSubfamily.Basename)

		gFormats := gSubfamily.Formats
		wFormats := wSubfamily.Formats
		test.VerifyFatal(t, 5, j, len(wFormats), len(gFormats))
		for k, wFormat := range wFormats {
			gFormat := gFormats[k]
			test.Verify(t, 6, j, wFormat, gFormat)
		}

		test.Verify(t, 7, j, wSubfamily.Style, gSubfamily.Style)
		test.Verify(t, 8, j, wSubfamily.Weight, gSubfamily.Weight)
	}

	test.Verify(t, 9, 0, wMetadata.path, gMetadata.path)
}
