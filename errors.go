// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPClientError custom error to handle with response status.
type HTTPClientError struct {
	StatusCode int
	Err        error
}

func (e *HTTPClientError) Error() string {
	return fmt.Sprintf("status %d, err: %v", e.StatusCode, e.Err)
}

func makeHTTPClientError(url string, resp *http.Response) error {
	var resError error

	respStatusText := "Unknown status text"
	switch resp.StatusCode {
	case 400:
		respStatusText = "The request encoding is invalid; the request can't be parsed as a valid JSON."
	case 401:
		respStatusText = "Accessinga protected resource without authorization or with invalid credentials."
	case 402:
		respStatusText = "The account associated with the API key making requests hits a quota that can be increased by upgrading the Airtable account plan."
	case 403:
		respStatusText = "Accessing a protected resource with API credentials that don't have access to that resource."
	case 404:
		respStatusText = "Route or resource is not found. This error is returned when the request hits an undefined route, or if the resource doesn't exist (e.g. has been deleted)."
	case 413:
		respStatusText = "Too Large The request exceeded the maximum allowed payload size. You shouldn't encounter this under normal use."
	case 422:
		respStatusText = "The request data is invalid. This includes most of the base-specific validations. You will receive a detailed error message and code pointing to the exact issue."
	case 500:
		respStatusText = "Error The server encountered an unexpected condition."
	case 502:
		respStatusText = "Airtable's servers are restarting or an unexpected outage is in progress. You should generally not receive this error, and requests are safe to retry."
	case 503:
		respStatusText = "The server could not process your request in time. The server could be temporarily unavailable, or it could have timed out processing your request. You should retry the request with backoffs."
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resError = fmt.Errorf("HTTP request failure on %s:\n%d %s\n%s\n\nCannot parse body with err: %w",
			url, resp.StatusCode, resp.Status, respStatusText, err)
	} else {
		resError = fmt.Errorf("HTTP request failure on %s:\n%d %s\n%s\n\nBody: %v",
			url, resp.StatusCode, resp.Status, respStatusText, string(body))
	}

	return &HTTPClientError{
		StatusCode: resp.StatusCode,
		Err:        resError,
	}
}
