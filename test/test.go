// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

// Package test implements general purpose testing routines.
package test

import "testing"

func Verify(t *testing.T, i, j int, want, got interface{}) {
	if got != want {
		t.Errorf("(%d:%d) => want `%v', got `%v'", i, j, want, got)
	}
}

func VerifyFatal(t *testing.T, i, j int, want, got interface{}) {
	if got != want {
		t.Fatalf("(%d:%d) => want `%v', got `%v'", i, j, want, got)
	}
}
