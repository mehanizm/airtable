// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"errors"
	"testing"
)

func TestTable_DeleteRecords(t *testing.T) {
	table := testTable(t)
	table.client.baseURL = mockResponse("delete_records.json").URL
	records, err := table.DeleteRecords([]string{"recnTq6CsvFM6vX2m", "recr3qAQbM7juKa4o"})
	if err != nil {
		t.Error("must be no error")
	}
	for _, record := range records.Records {
		if !record.Deleted {
			t.Errorf("expected that record will be deleted, but was: %#v", record.Deleted)
		}
	}
	table.client.baseURL = mockErrorResponse(404).URL
	_, err = table.DeleteRecords([]string{})
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}

func TestTable_AddRecords(t *testing.T) {
	table := testTable(t)
	table.client.baseURL = mockResponse("get_records_with_filter.json").URL
	toSend := new(Records)
	records, err := table.AddRecords(toSend)
	if err != nil {
		t.Error("must be no error")
	}
	if len(records.Records) != 3 {
		t.Errorf("should be 3 records in result, but was: %v", len(records.Records))
	}
	table.client.baseURL = mockErrorResponse(404).URL
	_, err = table.AddRecords(toSend)
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}

func TestTable_UpdateRecords(t *testing.T) {
	table := testTable(t)
	table.client.baseURL = mockResponse("get_records_with_filter.json").URL
	toSend := new(Records)
	records, err := table.UpdateRecords(toSend)
	if err != nil {
		t.Error("must be no error")
	}
	if len(records.Records) != 3 {
		t.Errorf("should be 3 records in result, but was: %v", len(records.Records))
	}
	table.client.baseURL = mockErrorResponse(404).URL
	_, err = table.UpdateRecords(toSend)
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}

func TestTable_UpdateRecordsPartial(t *testing.T) {
	table := testTable(t)
	table.client.baseURL = mockResponse("get_records_with_filter.json").URL
	toSend := new(Records)
	records, err := table.UpdateRecordsPartial(toSend)
	if err != nil {
		t.Error("must be no error")
	}
	if len(records.Records) != 3 {
		t.Errorf("should be 3 records in result, but was: %v", len(records.Records))
	}
	table.client.baseURL = mockErrorResponse(404).URL
	_, err = table.UpdateRecordsPartial(toSend)
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}

func testTable(t *testing.T) *Table {
	client := testClient(t)
	return client.GetTable("dbName", "tableName")
}
