// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package inventory

import (
	"path/filepath"
	"testing"

	"github.com/noll/mjau/font"
	"github.com/noll/mjau/test"
)

var (
	fl = filepath.FromSlash("../fonts")     // Font library path.
	tp = filepath.FromSlash("../templates") // Templates path.

	amf = filepath.Join(fl, "Amaranth")  // Amaranth family path.	
	osf = filepath.Join(fl, "Open Sans") // Open Sans family path.
)

var InventoryBuildTestCases = []struct {
	Query Query
	Font  *font.Font
}{
	// Case 1
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "eot400normal",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.EOT,
			Path:   filepath.Join(amf, "amaranth-regular.eot"),
			Style:  "normal",
			Weight: 400,
		},
	},
	// Case 2
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "woff400normal",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.WOFF,
			Path:   filepath.Join(amf, "amaranth-regular.woff"),
			Style:  "normal",
			Weight: 400,
		},
	},
	// Case 3
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "eot400italic",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.EOT,
			Path:   filepath.Join(amf, "amaranth-italic.eot"),
			Style:  "italic",
			Weight: 400,
		},
	},
	// Case 4
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "woff400italic",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.WOFF,
			Path:   filepath.Join(amf, "amaranth-italic.woff"),
			Style:  "italic",
			Weight: 400,
		},
	},
	// Case 5
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "eot700normal",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.EOT,
			Path:   filepath.Join(amf, "amaranth-bold.eot"),
			Style:  "normal",
			Weight: 700,
		},
	},
	// Case 6
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "woff700normal",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.WOFF,
			Path:   filepath.Join(amf, "amaranth-bold.woff"),
			Style:  "normal",
			Weight: 700,
		},
	},
	// Case 7
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "eot700italic",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.EOT,
			Path:   filepath.Join(amf, "amaranth-bolditalic.eot"),
			Style:  "italic",
			Weight: 700,
		},
	},
	// Case 8
	{
		Query{
			RowKey:    "Amaranth",
			ColumnKey: "woff700italic",
		},
		&font.Font{
			Family: "Amaranth",
			Format: font.WOFF,
			Path:   filepath.Join(amf, "amaranth-bolditalic.woff"),
			Style:  "italic",
			Weight: 700,
		},
	},
	// Case 9
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot300normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-light.eot"),
			Style:  "normal",
			Weight: 300,
		},
	},
	// Case 10
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff300normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-light.woff"),
			Style:  "normal",
			Weight: 300,
		},
	},
	// Case 11
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot300italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-lightitalic.eot"),
			Style:  "italic",
			Weight: 300,
		},
	},
	// Case 12
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff300italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-lightitalic.woff"),
			Style:  "italic",
			Weight: 300,
		},
	},
	// Case 13
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot400normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-regular.eot"),
			Style:  "normal",
			Weight: 400,
		},
	},
	// Case 14
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff400normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-regular.woff"),
			Style:  "normal",
			Weight: 400,
		},
	},
	// Case 15
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot400italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-italic.eot"),
			Style:  "italic",
			Weight: 400,
		},
	},
	// Case 16
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff400italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-italic.woff"),
			Style:  "italic",
			Weight: 400,
		},
	},
	// Case 17
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot600normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-semibold.eot"),
			Style:  "normal",
			Weight: 600,
		},
	},
	// Case 18
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff600normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-semibold.woff"),
			Style:  "normal",
			Weight: 600,
		},
	},
	// Case 19
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot600italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-semibolditalic.eot"),
			Style:  "italic",
			Weight: 600,
		},
	},
	// Case 20
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff600italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-semibolditalic.woff"),
			Style:  "italic",
			Weight: 600,
		},
	},
	// Case 21
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot700normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-bold.eot"),
			Style:  "normal",
			Weight: 700,
		},
	},
	// Case 22
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff700normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-bold.woff"),
			Style:  "normal",
			Weight: 700,
		},
	},
	// Case 23
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot700italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-bolditalic.eot"),
			Style:  "italic",
			Weight: 700,
		},
	},
	// Case 24
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff700italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-bolditalic.woff"),
			Style:  "italic",
			Weight: 700,
		},
	},
	// Case 25
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot800normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-extrabold.eot"),
			Style:  "normal",
			Weight: 800,
		},
	},
	// Case 26
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff800normal",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-extrabold.woff"),
			Style:  "normal",
			Weight: 800,
		},
	},
	// Case 27
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "eot800italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.EOT,
			Path:   filepath.Join(osf, "os-extrabolditalic.eot"),
			Style:  "italic",
			Weight: 800,
		},
	},

	// Case 28
	{
		Query{
			RowKey:    "Open Sans",
			ColumnKey: "woff800italic",
		},
		&font.Font{
			Family: "Open Sans",
			Format: font.WOFF,
			Path:   filepath.Join(osf, "os-extrabolditalic.woff"),
			Style:  "italic",
			Weight: 800,
		},
	},
}

func TestInventoryBuild(t *testing.T) {
	inventory := New()
	err := inventory.Build(fl)
	test.VerifyFatal(t, 1, 0, true, nil == err)
	wSize := len(InventoryBuildTestCases)
	test.VerifyFatal(t, 2, 0, wSize, inventory.Len())

	for i, c := range InventoryBuildTestCases {
		j := i + 1
		wFont := c.Font
		gFont := inventory.Query(c.Query)
		test.VerifyFatal(t, 3, j, false, nil == gFont)
		test.Verify(t, 4, j, wFont.Family, gFont.Family)
		test.Verify(t, 5, j, true, wFont.Format.Equal(gFont.Format))
		test.Verify(t, 6, j, wFont.Style, gFont.Style)
		test.Verify(t, 7, j, wFont.Weight, gFont.Weight)
		test.Verify(t, 8, j, wFont.Path, gFont.Path)
	}
}

func TestQueryQuery(t *testing.T) {
	inventory := New()
	query := Query{
		RowKey:    "Row",
		ColumnKey: "Column",
	}
	gFont := inventory.Query(query)
	test.VerifyFatal(t, 1, 0, true, nil == gFont)

	wFont := &font.Font{
		Family: "Amaranth",
		Format: font.EOT,
		Path:   filepath.Join(amf, "amaranth-regular.eot"),
		Style:  "normal",
		Weight: 400,
	}
	inventory.Put("Amaranth", "woff400normal", wFont)
	query = Query{
		RowKey:    "Amaranth",
		ColumnKey: "woff400normal",
	}
	gFont = inventory.Query(query)
	test.VerifyFatal(t, 2, 0, false, nil == gFont)
	test.Verify(t, 3, 0, wFont.Family, gFont.Family)
	test.Verify(t, 4, 0, true, wFont.Format.Equal(gFont.Format))
	test.Verify(t, 5, 0, wFont.Style, gFont.Style)
	test.Verify(t, 6, 0, wFont.Weight, gFont.Weight)
	test.Verify(t, 7, 0, wFont.Path, gFont.Path)
}
