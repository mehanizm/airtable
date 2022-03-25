// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestHTTPClientError_Error(t *testing.T) {
	e := &HTTPClientError{
		StatusCode: 300,
		Err:        errors.New("error message"),
	}
	expected := "status 300, err: error message"
	if got := e.Error(); got != expected {
		t.Errorf("HTTPClientError.Error() = %v, want %v", got, expected)
	}
}

func Test_makeHTTPClientError(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		resp     *http.Response
		expected string
	}{
		{
			name: "400",
			url:  "google.com",
			resp: &http.Response{
				StatusCode: 400,
				Status:     "Bad Request",
				Body:       io.NopCloser(bytes.NewReader([]byte("body"))),
			},
			expected: `status 400, err: HTTP request failure on google.com:
400 Bad Request
The request encoding is invalid; the request can't be parsed as a valid JSON.

Body: body`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeHTTPClientError(tt.url, tt.resp); err.Error() != tt.expected {
				t.Errorf("makeHTTPClientError() error:\n%v\n\nexpected:\n%v", err, tt.expected)
			}
		})
	}
}
