// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPClientError custom error to handle with response status
type HTTPClientError struct {
	StatusCode int
	Err        error
}

func (e *HTTPClientError) Error() string {
	return fmt.Sprintf("status %d, err: %v", e.StatusCode, e.Err)
}

func makeHTTPClientError(url string, resp *http.Response) error {
	var resError error
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resError = fmt.Errorf("HTTP request failure on %s with status %d\nCannot parse body with: %w", url, resp.StatusCode, err)
	} else {
		resError = fmt.Errorf("HTTP request failure on %s with status %d\nBody: %v", url, resp.StatusCode, string(body))
	}
	return &HTTPClientError{
		StatusCode: resp.StatusCode,
		Err:        resError,
	}
}
