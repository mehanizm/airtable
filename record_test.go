// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"errors"
	"reflect"
	"testing"
)

func TestRecord_GetRecord(t *testing.T) {
	table := testTable()
	table.client.baseURL = mockResponse("get_record.json").URL
	record, err := table.GetRecord("recnTq6CsvFM6vX2m")
	if err != nil {
		t.Error("must be no error")
	}
	expected := &Record{
		client:      table.client,
		table:       table,
		ID:          "recnTq6CsvFM6vX2m",
		CreatedTime: "2020-04-10T11:30:57.000Z",
		Fields: map[string]any{
			"Field1": "Field1",
			"Field2": true,
			"Field3": "2020-04-06T06:00:00.000Z",
		},
	}
	if !reflect.DeepEqual(record, expected) {
		t.Errorf("expected: %#v\nbut got: %#v\n", expected, record)
	}
	table.client.baseURL = mockErrorResponse(404).URL
	_, err = table.GetRecord("recnTq6CsvFM6vX2m")
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}

func TestRecord_DeleteRecord(t *testing.T) {
	record := testRecord(t)
	record.client.baseURL = mockResponse("delete_record.json").URL
	res, err := record.DeleteRecord()
	if err != nil {
		t.Error("must be no error")
	}
	if !res.Deleted {
		t.Errorf("expected that record will be deleted, but was: %#v", record.Deleted)
	}
	record.client.baseURL = mockErrorResponse(404).URL
	_, err = record.DeleteRecord()
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}

func TestRecord_UpdateRecordPartial(t *testing.T) {
	record := testRecord(t)
	record.client.baseURL = mockResponse("get_records_with_filter.json").URL
	res, err := record.UpdateRecordPartial(map[string]any{"Field_2": true})
	if err != nil {
		t.Error("must be no error")
	}
	resBool, ok := res.Fields["Field2"].(bool)
	if !ok {
		t.Errorf("Field2 should be bool type, but was %#v\n\nFull resp: %#v", res.Fields["Field2"], res)
	}
	if !resBool {
		t.Errorf("expected that Field_2 will be true, but was: %#v", res.Fields["Field2"].(bool))
	}
	record.client.baseURL = mockErrorResponse(404).URL
	_, err = record.UpdateRecordPartial(map[string]any{})
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}

func testRecord(t *testing.T) *Record {
	table := testTable()
	table.client.baseURL = mockResponse("get_record.json").URL
	record, err := table.GetRecord("recordID")
	if err != nil {
		t.Error("must be no error")
	}
	return record
}
