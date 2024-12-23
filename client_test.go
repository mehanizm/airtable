// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"testing"
)

func testClient() *Client {
	c := NewClient("apiKey")
	c.SetRateLimit(1000)
	return c
}

func TestClient_do(t *testing.T) {
	c := testClient()
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

func TestClient_SetBaseURL(t *testing.T) {
	testCases := []struct {
		description   string
		url           string
		expectedError string
	}{
		{
			description: "accepts a valid URL",
			url:         "http://localhost:3000",
		},
		{
			description: "accepts a valid HTTPS URL",
			url:         "https://example.com",
		},
		{
			description:   "rejects non http/https scheme url",
			url:           "ftp://example.com",
			expectedError: "http or https baseURL must be used",
		},
		{
			description:   "rejects url without scheme",
			url:           "example.com",
			expectedError: "scheme of http or https must be specified",
		},
	}

	c := testClient()

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			err := c.SetBaseURL(testCase.url)
			if err != nil {
				if testCase.expectedError != "" {
					if !strings.Contains(err.Error(), testCase.expectedError) {
						t.Fatalf("unexpected error: %s", err)
					}
				} else {
					t.Fatalf("unexpected error: %s", err)
				}
			}
		})
	}
}
