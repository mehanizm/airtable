package airtable

import (
	"net/url"
)

// GetBasesConfig helper type to use in.
// step by step get bases.
type GetBasesConfig struct {
	client *Client
	params url.Values
}

// GetBases prepare step to get bases.
func (c *Client) GetBases() *GetBasesConfig {
	return &GetBasesConfig{
		client: c,
		params: url.Values{},
	}
}

// Pagination
// The server returns one page of bases at a time.

// If there are more records, the response will contain an offset.
// To fetch the next page of records, include offset in the next request's parameters.

// WithOffset Pagination will stop when you've reached the end of your bases.
func (gbc *GetBasesConfig) WithOffset(offset string) *GetBasesConfig {
	gbc.params.Set("offset", offset)
	return gbc
}

// Do send the prepared get records request.
func (gbc *GetBasesConfig) Do() (*Bases, error) {
	return gbc.client.GetBasesWithParams(gbc.params)
}
