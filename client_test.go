// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"errors"
	"log"
	"net/http"
	"testing"
)

func testClient(t *testing.T) *Client {
	c := NewClient("apiKey")
	c.SetRateLimit(1000)
	return c
}

func TestClient_do(t *testing.T) {
	c := testClient(t)
	url := mockErrorResponse(404).URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = c.do(req, nil)
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
	err = c.do(nil, nil)
	if err == nil {
		t.Errorf("there should be an error, but was nil")
	}
}
