// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

// Error reports a generic error.
type Error string

func (e Error) Error() string {
	return ProgName + ": " + string(e)
}

// PrintErrorExit prints the given error message to standard error
// and exits the program signaling abnormal termination.
func PrintErrorExit(message string) {
	fmt.Fprintln(os.Stderr, Error(message))
	os.Exit(1)
}
