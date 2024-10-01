// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	airtableBaseURL                 = "https://api.airtable.com/v0"
	airtableUploadAttachmentBaseURL = "https://content.airtable.com/v0/"
	rateLimit                       = 4
)

// Client client for airtable api.
type Client struct {
	client                  *http.Client
	rateLimiter             <-chan time.Time
	baseURL                 string
	uploadAttachmentBaseURL string
	apiKey                  string
}

// NewClient airtable client constructor
// your API KEY you can get on your account page
// https://airtable.com/account
func NewClient(apiKey string) *Client {
	return &Client{
		client:                  http.DefaultClient,
		rateLimiter:             time.Tick(time.Second / time.Duration(rateLimit)),
		apiKey:                  apiKey,
		baseURL:                 airtableBaseURL,
		uploadAttachmentBaseURL: airtableUploadAttachmentBaseURL,
	}
}

// Set custom http client for custom usage
func (at *Client) SetCustomClient(client *http.Client) {
	at.client = client
}

// SetRateLimit rate limit setter for custom usage
// Airtable limit is 5 requests per second (we use 4)
// https://airtable.com/{yourDatabaseID}/api/docs#curl/ratelimits
func (at *Client) SetRateLimit(customRateLimit int) {
	at.rateLimiter = time.Tick(time.Second / time.Duration(customRateLimit))
}

func (at *Client) SetBaseURL(baseURL string) error {
	url, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse baseURL: %s", err)
	}

	if url.Scheme == "" {
		return fmt.Errorf("scheme of http or https must be specified")
	}

	if url.Scheme != "https" && url.Scheme != "http" {
		return fmt.Errorf("http or https baseURL must be used")
	}

	at.baseURL = url.String()

	return nil
}

func (at *Client) rateLimit() {
	<-at.rateLimiter
}

func (at *Client) get(ctx context.Context, db, table, recordID string, params url.Values, target any) error {
	at.rateLimit()

	url := fmt.Sprintf("%s/%s/%s", at.baseURL, db, table)
	if recordID != "" {
		url += fmt.Sprintf("/%s", recordID)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.apiKey))

	req.URL.RawQuery = params.Encode()

	err = at.do(req, target)
	if err != nil {
		return err
	}

	return nil
}

func (at *Client) post(ctx context.Context, db, table string, data, response any) error {
	at.rateLimit()

	url := fmt.Sprintf("%s/%s/%s", at.baseURL, db, table)

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.apiKey))

	return at.do(req, response)
}

func (at *Client) postAttachment(ctx context.Context, db, recordID string, attachmentFieldIdOrName string, data Attachment, response any) error {
	at.rateLimit()

	url := fmt.Sprintf("%s/%s/%s/%s/uploadAttachment", at.uploadAttachmentBaseURL, db, recordID, attachmentFieldIdOrName)

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.apiKey))

	return at.do(req, response)
}

func (at *Client) delete(ctx context.Context, db, table string, recordIDs []string, target any) error {
	at.rateLimit()

	rawURL := fmt.Sprintf("%s/%s/%s", at.baseURL, db, table)
	params := url.Values{}

	for _, recordID := range recordIDs {
		params.Add("records[]", recordID)
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", rawURL, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.apiKey))

	req.URL.RawQuery = params.Encode()

	err = at.do(req, target)
	if err != nil {
		return err
	}

	return nil
}

func (at *Client) patch(ctx context.Context, db, table, data, response any) error {
	at.rateLimit()

	url := fmt.Sprintf("%s/%s/%s", at.baseURL, db, table)

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.apiKey))

	return at.do(req, response)
}

func (at *Client) put(ctx context.Context, db, table, data, response any) error {
	at.rateLimit()

	url := fmt.Sprintf("%s/%s/%s", at.baseURL, db, table)

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.apiKey))

	return at.do(req, response)
}

func (at *Client) do(req *http.Request, response any) error {
	if req == nil {
		return errors.New("nil request")
	}

	url := req.URL.RequestURI()

	resp, err := at.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failure on %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return makeHTTPClientError(url, resp)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP Read error on response for %s: %w", url, err)
	}

	err = json.Unmarshal(b, response)
	if err != nil {
		return fmt.Errorf("JSON decode failed on %s:\n%s\nerror: %w", url, string(b), err)
	}

	return nil
}
