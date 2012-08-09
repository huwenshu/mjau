// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

// Package util implements general purpose utility routines.
package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// Base64 returns the contents of b as a base64-encoded string.
func Base64(b []byte) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(b)
	encoder.Close()
	return buf.String()
}

// BlankStrFlagDefault sets value pointed to by v, if blank, to the default
// value of the named string command-line flag. If the named flag does not
// exist, exits the program signaling abnormal termination.
func BlankStrFlagDefault(v *string, name string) {
	if *v == "" {
		if flag := flag.Lookup(name); flag == nil {
			// Should not happen.
			os.Exit(1)
		} else {
			*v = flag.DefValue
		}
	}
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsDir reports whether d is a directory.
func IsDir(d string) (y bool) {
	if fi, err := os.Stat(d); err == nil {
		if fi.IsDir() {
			y = true
		}
	}
	return
}

// ReadJson reads and parses the JSON-encoded contents of the named file and
// stores the result in the value pointed to by v.
// Returns an error if the named file cannot be read or correctly parsed.
func ReadJson(name string, v interface{}) error {
	if b, err := ioutil.ReadFile(name); err != nil {
		return err
	} else {
		if err := json.Unmarshal(b, &v); err != nil {
			return fmt.Errorf("parse %s: %s", name, err)
		}
	}
	return nil
}
