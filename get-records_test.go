// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"context"
	"testing"
)

func TestGetRecordsConfig_Do(t *testing.T) {
	table := testTable()
	table.client.baseURL = mockResponse("get_records_with_filter.json").URL
	sortQuery1 := struct {
		FieldName string
		Direction string
	}{"Field1", "desc"}
	sortQuery2 := struct {
		FieldName string
		Direction string
	}{"Field2", "asc"}

	records, err := table.GetRecords().
		FromView("view_1").
		WithFilterFormula("AND({Field1}='value_1',NOT({Field2}='value_2'))").
		WithSort(sortQuery1, sortQuery2).
		WithFieldsById().
		ReturnFields("Field1", "Field2").
		InStringFormat("Europe/Moscow", "ru").
		MaxRecords(100).
		PageSize(10).
		WithOffset("hhh").
		Do()
	if err != nil {
		t.Errorf("there should not be an err, but was: %v", err)
	}
	if len(records.Records) != 3 {
		t.Errorf("there should be 3 records, but was %v", len(records.Records))
	}

	table.client.baseURL = mockErrorResponse(400).URL
	_, err = table.GetRecords().Do()
	if err == nil {
		t.Errorf("there should be an err, but was nil")
	}
}
func TestGetRecordsConfig_IterAllRecords(t *testing.T) {
	expectedRecordsIDs := []string{
		"recnTq6CsvFM6vX2m",
		"recr3qAQbM7juKa4o",
		"recr3qAQbM7juKa4a",
		"recnFR6CsvFM6vX2m",
		"recrLdAQbM7juKa4o",
		"recr64AQbM7juKa4a",
	}

	table := testTable()
	table.client.baseURL = mockResponse("get_records_with_filter.json").URL

	recordsIter := table.GetRecords().
		MaxRecords(100).
		PageSize(3). // mocked response for each page is 3 records
		IterAllRecords(context.Background())

	recordsNum := 0
	for rec, err := range recordsIter {
		if err != nil {
			t.Errorf("there should not be an err, but was: %v", err)
			continue
		}
		if rec == nil {
			t.Errorf("there should not be a nil record, but was")
			continue
		}
		if len(expectedRecordsIDs) > recordsNum {
			if rec.ID != expectedRecordsIDs[recordsNum] {
				t.Errorf("for record num %d, expected record ID %s, but got %s", recordsNum, expectedRecordsIDs[recordsNum], rec.ID)
			}
		}
		recordsNum++
		if recordsNum == 3 {
			table.client.baseURL = mockResponse("get_records_stage_2.json").URL
		}
	}
	if recordsNum != len(expectedRecordsIDs) {
		t.Errorf("there should be %d records, but was %d", len(expectedRecordsIDs), recordsNum)
	}

}
