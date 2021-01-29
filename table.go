// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"net/url"
)

// Records base type of airtable records
type Records struct {
	Records []*Record `json:"records"`
	Offset  string    `json:"offset,omitempty"`
}

// Table represents table object
type Table struct {
	client    *Client
	dbName    string
	tableName string
}

// GetTable return table object
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
	records := new(Records)
	err := t.client.get(t.dbName, t.tableName, "", params, records)
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
	result := new(Records)
	err := t.client.post(t.dbName, t.tableName, records, result)
	if err != nil {
		return nil, err
	}
	for _, record := range result.Records {
		record.client = t.client
		record.table = t
	}
	return result, err
}

// UpdateRecords full update records
func (t *Table) UpdateRecords(records *Records) (*Records, error) {
	response := new(Records)
	err := t.client.post(t.dbName, t.tableName, records, response)
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
func (t *Table) DeleteRecords(recordIDs []string) (*Records, error) {
	response := new(Records)
	err := t.client.delete(t.dbName, t.tableName, recordIDs, response)
	if err != nil {
		return nil, err
	}
	for _, record := range response.Records {
		record.client = t.client
		record.table = t
	}
	return response, nil
}
