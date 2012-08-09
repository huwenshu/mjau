// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package whitelist

import (
	"path/filepath"
	"testing"

	"github.com/noll/mjau/test"
)

var td = filepath.FromSlash("./test") // Test directory.

func TestWhitelistContains(t *testing.T) {
	whitelist := New()
	whitelist.Domains = append(whitelist.Domains, "http://one/")
	test.Verify(t, 1, 0, true, whitelist.Contains("http://one/"))
	test.Verify(t, 2, 0, true, whitelist.Contains("http://one/two"))
	test.Verify(t, 3, 0, false, whitelist.Contains("http://one"))
}

func TestWhitelistRead(t *testing.T) {
	gWhitelist := New()
	err := gWhitelist.Read(filepath.Join(td, "whitelist.json"))
	test.VerifyFatal(t, 1, 0, true, nil == err)
	wWhitelist := Whitelist{
		Domains: []string{
			"http://one/",
			"http://two/",
			"http://three/",
		},
	}
	test.VerifyFatal(t, 2, 0, wWhitelist.Size(), gWhitelist.Size())

	for i, wdomain := range wWhitelist.Domains {
		j := i + 1
		gdomain := gWhitelist.Domains[i]
		test.Verify(t, 3, j, wdomain, gdomain)
	}
}

func TestWhitelistSize(t *testing.T) {
	whitelist := New()
	test.Verify(t, 1, 0, 0, whitelist.Size())

	whitelist.Domains = append(whitelist.Domains, "one")
	test.Verify(t, 2, 0, 1, whitelist.Size())

	whitelist.Domains = append(whitelist.Domains, "two")
	test.Verify(t, 3, 0, 2, whitelist.Size())
}
