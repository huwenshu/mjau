// Copyright (c) 2012, Robert Dinu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license which can be found in the LICENSE file.

package http

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/noll/mjau/test"
)

func TestMakeGzipHandler(t *testing.T) {
	handler := MakeGzipHandler(helloHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	test.VerifyFatal(t, 1, 0, true, nil == err)

	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	test.VerifyFatal(t, 2, 0, true, nil == err)

	wContentEncoding := "gzip"
	gContentEncoding := resp.Header.Get("Content-Encoding")
	test.Verify(t, 3, 0, wContentEncoding, gContentEncoding)

	buf := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(buf)
	fmt.Fprint(gzipWriter, "Hej!")
	err = gzipWriter.Close()
	test.VerifyFatal(t, 4, 0, true, nil == err)
	wBody := buf.Bytes()
	gBody, err := ioutil.ReadAll(resp.Body)
	test.VerifyFatal(t, 5, 0, true, nil == err)
	test.Verify(t, 6, 0, true, bytes.Equal(wBody, gBody))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hej!")
}
