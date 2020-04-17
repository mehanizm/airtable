// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"errors"
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
