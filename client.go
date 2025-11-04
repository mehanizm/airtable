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
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	airtableBaseURL                 = "https://api.airtable.com/v0"
	airtableUploadAttachmentBaseURL = "https://content.airtable.com/v0"
	rateLimit                       = 4
)

var (
	// backoff settings for retrying 503 responses
	backoffMaxRetries = 3
	backoffBaseDelay  = 100 * time.Millisecond
	backoffMaxJitter  = 100 * time.Millisecond
	backoffMaxDelay   = 2 * time.Second
)

// Client client for airtable api.
type Client struct {
	client                  *http.Client
	rateLimiter             *rate.Limiter
	baseURL                 string
	uploadAttachmentBaseURL string
	apiKey                  string
	// backoff parameters (per-client)
	maxRetries      int
	backoffBase     time.Duration
	backoffJitter   time.Duration
	backoffMaxDelay time.Duration
}

// NewClient airtable client constructor
// your API KEY you can get on your account page
// https://airtable.com/account
func NewClient(apiKey string) *Client {
	return &Client{
		client:                  http.DefaultClient,
		rateLimiter:             rate.NewLimiter(rate.Limit(rateLimit), 1),
		apiKey:                  apiKey,
		baseURL:                 airtableBaseURL,
		uploadAttachmentBaseURL: airtableUploadAttachmentBaseURL,
		maxRetries:              backoffMaxRetries,
		backoffBase:             backoffBaseDelay,
		backoffJitter:           backoffMaxJitter,
		backoffMaxDelay:         backoffMaxDelay,
	}
}

// SetBackoffRetries sets how many retries will be attempted on 503 responses.
func (at *Client) SetBackoffRetries(n int) {
	at.maxRetries = n
}

// SetBackoffBaseDelay sets the base delay used for exponential backoff.
func (at *Client) SetBackoffBaseDelay(d time.Duration) {
	at.backoffBase = d
}

// SetBackoffMaxJitter sets the maximum jitter added to backoff sleeps.
func (at *Client) SetBackoffMaxJitter(d time.Duration) {
	at.backoffJitter = d
}

// SetBackoffMaxDelay sets the maximum delay (cap) for exponential backoff.
// If set to 0, backoff delay will be uncapped
func (at *Client) SetBackoffMaxDelay(d time.Duration) {
	at.backoffMaxDelay = d
}

// Set custom http client for custom usage
func (at *Client) SetCustomClient(client *http.Client) {
	at.client = client
}

// SetRateLimit rate limit setter for custom usage
// Airtable limit is 5 requests per second (we use 4)
// https://airtable.com/{yourDatabaseID}/api/docs#curl/ratelimits
func (at *Client) SetRateLimit(customRateLimit int) {
	at.rateLimiter = rate.NewLimiter(rate.Limit(customRateLimit), 1)
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

func (at *Client) rateLimit(ctx context.Context) error {
	return at.rateLimiter.Wait(ctx)
}

func (at *Client) get(ctx context.Context, db, table, recordID string, params url.Values, target any) error {
	err := at.rateLimit(ctx)
	if err != nil {
		return err
	}

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
	err := at.rateLimit(ctx)
	if err != nil {
		return err
	}

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
	err := at.rateLimit(ctx)
	if err != nil {
		return err
	}

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
	err := at.rateLimit(ctx)
	if err != nil {
		return err
	}

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
	err := at.rateLimit(ctx)
	if err != nil {
		return err
	}

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
	err := at.rateLimit(ctx)
	if err != nil {
		return err
	}

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

	// Ensure request body can be replayed for retries. If GetBody is not
	// provided, read the body and set GetBody so we can recreate Body for
	// each attempt.
	if req.Body != nil && req.GetBody == nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("reading request body before retrying: %w", err)
		}
		// reset original Body
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}
	}

	// Attempt requests with retries on 503 and 429. On non-retryable
	// non-2xx responses we return immediately. On retryable responses we
	// retry with exponential backoff + jitter, capped by backoffMaxDelay.
	for attempt := 0; attempt <= at.maxRetries; attempt++ {
		creq := req.Clone(req.Context())
		if req.GetBody != nil {
			rc, err := req.GetBody()
			if err != nil {
				return fmt.Errorf("failed to get request body for retry: %w", err)
			}
			creq.Body = rc
		}

		resp, err := at.client.Do(creq)
		if err != nil {
			return fmt.Errorf("HTTP request failure on %s: %w", url, err)
		}
		// If it's a retryable status (503 or 429), handle retries.
		if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusTooManyRequests {
			if attempt == at.maxRetries {
				// final attempt, return error; do not close body so the
				// error helper can read it.
				return makeHTTPClientError(url, resp)
			}

			// Determine sleep duration: prefer Retry-After header if present.
			var sleep time.Duration
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				// try parse as integer seconds
				if secs, err := strconv.Atoi(strings.TrimSpace(ra)); err == nil {
					sleep = time.Duration(secs) * time.Second
				} else if t, err := http.ParseTime(ra); err == nil {
					sleep = time.Until(t)
					if sleep < 0 {
						sleep = 0
					}
				}
			}

			// If no Retry-After, use exponential backoff
			if sleep == 0 {
				sleep = at.backoffBase * time.Duration(1<<attempt)
			}

			// clip to max delay
			if at.backoffMaxDelay > 0 && sleep > at.backoffMaxDelay {
				sleep = at.backoffMaxDelay
			}

			// drain and close so idle connection can be reused
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			// add jitter
			var jitter time.Duration
			if at.backoffJitter > 0 {
				jitter = time.Duration(rand.Int63n(int64(at.backoffJitter)))
			}

			time.Sleep(sleep + jitter)

			// try again
			continue
		}

		// Non-2xx responses (other than 503 handled above)
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return makeHTTPClientError(url, resp)
		}

		// Success path: read the body, close, and unmarshal
		b, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("HTTP Read error on response for %s: %w", url, err)
		}

		err = json.Unmarshal(b, response)
		if err != nil {
			return fmt.Errorf("JSON decode failed on %s:\n%s\nerror: %w", url, string(b), err)
		}

		return nil
	}

	// Should never reach here
	return errors.New("request retry loop exhausted")
}
