// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"context"
	"net/url"
)

// Records base type of airtable records.
type Records struct {
	Records []*Record `json:"records"`
	Offset  string    `json:"offset,omitempty"`

	// The Airtable API will perform best-effort automatic data conversion
	// from string values if the typecast parameter is passed in.
	// Automatic conversion is disabled by default to ensure data integrity,
	// but it may be helpful for integrating with 3rd party data sources.
	Typecast bool `json:"typecast,omitempty"`
}

// Table represents table object.
type Table struct {
	client    *Client
	dbName    string
	tableName string
}

// GetTable return table object.
func (c *Client) GetTable(dbName, tableName string) *Table {
	return &Table{
		client:    c,
		dbName:    dbName,
		tableName: tableName,
	}
}

// GetRecordsWithParams get records with url values params
// https://airtable.com/{yourDatabaseID}/api/docs#curl/table:{yourTableName}:list
func (t *Table) GetRecordsWithParams(params url.Values) (*Records, error) {
	return t.GetRecordsWithParamsContext(context.Background(), params)
}

// GetRecordsWithParamsContext get records with url values params
// with custom context
func (t *Table) GetRecordsWithParamsContext(ctx context.Context, params url.Values) (*Records, error) {
	records := new(Records)

	err := t.client.get(ctx, t.dbName, t.tableName, "", params, records)
	if err != nil {
		return nil, err
	}

	for _, record := range records.Records {
		record.client = t.client
		record.table = t
	}

	return records, nil
}

// AddRecords method to add lines to table (up to 10 in one request)
// https://airtable.com/{yourDatabaseID}/api/docs#curl/table:{yourTableName}:create
func (t *Table) AddRecords(records *Records) (*Records, error) {
	return t.AddRecordsContext(context.Background(), records)
}

// AddRecordsContext method to add lines to table (up to 10 in one request)
// with custom context
func (t *Table) AddRecordsContext(ctx context.Context, records *Records) (*Records, error) {
	result := new(Records)

	err := t.client.post(ctx, t.dbName, t.tableName, records, result)
	if err != nil {
		return nil, err
	}

	for _, record := range result.Records {
		record.client = t.client
		record.table = t
	}

	return result, err
}

// UpdateRecords full update records.
func (t *Table) UpdateRecords(records *Records) (*Records, error) {
	return t.UpdateRecordsContext(context.Background(), records)
}

// UpdateRecordsContext full update records with custom context.
func (t *Table) UpdateRecordsContext(ctx context.Context, records *Records) (*Records, error) {
	response := new(Records)

	err := t.client.put(ctx, t.dbName, t.tableName, records, response)
	if err != nil {
		return nil, err
	}

	for _, record := range response.Records {
		record.client = t.client
		record.table = t
	}

	return response, nil
}

// UpdateRecordsPartial partial update records.
func (t *Table) UpdateRecordsPartial(records *Records) (*Records, error) {
	return t.UpdateRecordsPartialContext(context.Background(), records)
}

// UpdateRecordsPartialContext partial update records with custom context.
func (t *Table) UpdateRecordsPartialContext(ctx context.Context, records *Records) (*Records, error) {
	response := new(Records)

	err := t.client.patch(ctx, t.dbName, t.tableName, records, response)
	if err != nil {
		return nil, err
	}

	for _, record := range response.Records {
		record.client = t.client
		record.table = t
	}

	return response, nil
}

// DeleteRecords delete records by recordID
// up to 10 ids in one request.
func (t *Table) DeleteRecords(recordIDs []string) (*Records, error) {
	return t.DeleteRecordsContext(context.Background(), recordIDs)
}

// DeleteRecordsContext delete records by recordID
// with custom context
func (t *Table) DeleteRecordsContext(ctx context.Context, recordIDs []string) (*Records, error) {
	response := new(Records)

	err := t.client.delete(ctx, t.dbName, t.tableName, recordIDs, response)
	if err != nil {
		return nil, err
	}

	for _, record := range response.Records {
		record.client = t.client
		record.table = t
	}

	return response, nil
}
