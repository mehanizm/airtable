// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import "net/url"

// Record base time of airtable record fields.
type Record struct {
	client      *Client
	table       *Table
	ID          string                 `json:"id,omitempty"`
	Fields      map[string]interface{} `json:"fields"`
	CreatedTime string                 `json:"createdTime,omitempty"`
	Deleted     bool                   `json:"deleted,omitempty"`

	// The Airtable API will perform best-effort automatic data conversion
	// from string values if the typecast parameter is passed in.
	// Automatic conversion is disabled by default to ensure data integrity,
	// but it may be helpful for integrating with 3rd party data sources.
	Typecast bool `json:"typecast,omitempty"`
}

// GetRecord get record from table
// https://airtable.com/{yourDatabaseID}/api/docs#curl/table:{yourTableName}:retrieve
func (t *Table) GetRecord(recordID string) (*Record, error) {
	result := new(Record)

	err := t.client.get(t.dbName, t.tableName, recordID, url.Values{}, result)
	if err != nil {
		return nil, err
	}

	result.client = t.client
	result.table = t

	return result, nil
}

// UpdateRecordPartial updates partial info on record.
func (r *Record) UpdateRecordPartial(changedFields map[string]interface{}) (*Record, error) {
	data := &Records{
		Records: []*Record{
			{
				ID:     r.ID,
				Fields: changedFields,
			},
		},
	}
	response := new(Records)

	err := r.client.patch(r.table.dbName, r.table.tableName, data, response)
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
	response := new(Records)

	err := r.client.delete(r.table.dbName, r.table.tableName, []string{r.ID}, response)
	if err != nil {
		return nil, err
	}

	result := response.Records[0]
	result.client = r.client
	result.table = r.table

	return result, nil
}
