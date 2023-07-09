// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"context"
	"net/url"
)

// Record base time of airtable record fields.
type Record struct {
	client      *Client
	table       *Table
	ID          string         `json:"id,omitempty"`
	Fields      map[string]any `json:"fields"`
	CreatedTime string         `json:"createdTime,omitempty"`
	Deleted     bool           `json:"deleted,omitempty"`

	// The Airtable API will perform best-effort automatic data conversion
	// from string values if the typecast parameter is passed in.
	// Automatic conversion is disabled by default to ensure data integrity,
	// but it may be helpful for integrating with 3rd party data sources.
	Typecast bool `json:"typecast,omitempty"`
}

// GetRecord get record from table
// https://airtable.com/{yourDatabaseID}/api/docs#curl/table:{yourTableName}:retrieve
func (t *Table) GetRecord(recordID string) (*Record, error) {
	return t.GetRecordContext(context.Background(), recordID)
}

// GetRecordContext get record from table
// with custom context
func (t *Table) GetRecordContext(ctx context.Context, recordID string) (*Record, error) {
	result := new(Record)

	err := t.client.get(ctx, t.dbName, t.tableName, recordID, url.Values{}, result)
	if err != nil {
		return nil, err
	}

	result.client = t.client
	result.table = t

	return result, nil
}

// UpdateRecordPartial updates partial info on record.
func (r *Record) UpdateRecordPartial(changedFields map[string]any) (*Record, error) {
	return r.UpdateRecordPartialContext(context.Background(), changedFields)
}

// UpdateRecordPartialContext updates partial info on record
// with custom context
func (r *Record) UpdateRecordPartialContext(ctx context.Context, changedFields map[string]any) (*Record, error) {
	data := &Records{
		Records: []*Record{
			{
				ID:     r.ID,
				Fields: changedFields,
			},
		},
	}
	response := new(Records)

	err := r.client.patch(ctx, r.table.dbName, r.table.tableName, data, response)
	if err != nil {
		return nil, err
	}

	result := response.Records[0]

	result.client = r.client
	result.table = r.table

	return result, nil
}

// DeleteRecord delete one record.
func (r *Record) DeleteRecord() (*Record, error) {
	return r.DeleteRecordContext(context.Background())
}

// DeleteRecordContext delete one record
// with custom context
func (r *Record) DeleteRecordContext(ctx context.Context) (*Record, error) {
	response := new(Records)

	err := r.client.delete(ctx, r.table.dbName, r.table.tableName, []string{r.ID}, response)
	if err != nil {
		return nil, err
	}

	result := response.Records[0]
	result.client = r.client
	result.table = r.table

	return result, nil
}
