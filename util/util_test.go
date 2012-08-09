// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package util

import (
	"path/filepath"
	"testing"

	"github.com/noll/mjau/test"
)

var (
	td = filepath.FromSlash("./test") // Test directory.

	ed = filepath.Join(td, "test.dir")  // Existent directory.
	ef = filepath.Join(td, "test.file") // Existent file.
	nd = filepath.Join(td, "ne.dir")    // Nonexistent directory.
	nf = filepath.Join(td, "ne.file")   // Nonexistent file.
)

func TestBase64(t *testing.T) {
	value := []byte("test")
	test.Verify(t, 1, 0, "dGVzdA==", Base64(value))
}

func TestExists(t *testing.T) {
	// Files
	test.Verify(t, 1, 0, true, Exists(ef))
	test.Verify(t, 2, 0, false, Exists(nf))
	// Directories
	test.Verify(t, 3, 0, true, Exists(ed))
	test.Verify(t, 4, 0, false, Exists(nd))
}

func TestIsDir(t *testing.T) {
	// Files
	test.Verify(t, 1, 0, false, IsDir(ef))
	test.Verify(t, 2, 0, false, IsDir(nf))
	// Directories
	test.Verify(t, 3, 0, true, IsDir(ed))
	test.Verify(t, 4, 0, false, IsDir(nd))
}

func TestReadJson(t *testing.T) {
	type TestJson struct {
		Test string
	}

	testJson := &TestJson{}
	err := ReadJson(filepath.Join(td, "test.json"), testJson)
	test.VerifyFatal(t, 1, 0, true, nil == err)
	test.Verify(t, 2, 0, "json", testJson.Test)
}
