// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"fmt"
	"net/url"
	"strconv"
)

// GetRecordsConfig helper type to use in
// step by step get records
type GetRecordsConfig struct {
	table  *Table
	params url.Values
}

// GetRecords prepare step to get records
func (t *Table) GetRecords() *GetRecordsConfig {
	return &GetRecordsConfig{
		table:  t,
		params: url.Values{},
	}
}

// ReturnFields set returning field names
func (grc *GetRecordsConfig) ReturnFields(fieldNames ...string) *GetRecordsConfig {
	for _, fieldName := range fieldNames {
		grc.params.Add("fields[]", fieldName)
	}
	return grc
}

// WithFilterFormula add filter to request
func (grc *GetRecordsConfig) WithFilterFormula(filterFormula string) *GetRecordsConfig {
	grc.params.Set("filterByFormula", filterFormula)
	return grc
}

// WithSort add sorting to request
func (grc *GetRecordsConfig) WithSort(sortQueries ...struct {
	fieldName string
	direction string
}) *GetRecordsConfig {
	for queryNum, sortQuery := range sortQueries {
		grc.params.Set(fmt.Sprintf("sort[%v][field]", queryNum), sortQuery.fieldName)
		grc.params.Set(fmt.Sprintf("sort[%v][direction]", queryNum), sortQuery.direction)
	}
	return grc
}

// FromView add view parameter to get records
func (grc *GetRecordsConfig) FromView(viewNameOrID string) *GetRecordsConfig {
	grc.params.Set("view", viewNameOrID)
	return grc
}

// The maximum total number of records that will be returned in your requests.
// If this value is larger than pageSize (which is 100 by default),
// you may have to load multiple pages to reach this total.
// See the Pagination section below for more.
func (grc *GetRecordsConfig) MaxRecords(maxRecords int) *GetRecordsConfig {
	grc.params.Set("maxRecords", strconv.Itoa(maxRecords))
	return grc
}

// The number of records returned in each request.
// Must be less than or equal to 100. Default is 100.
// See the Pagination section below for more.
func (grc *GetRecordsConfig) PageSize(pageSize int) *GetRecordsConfig {
	grc.params.Set("pageSize", strconv.Itoa(pageSize))
	return grc
}

// Pagination
// The server returns one page of records at a time.
// Each page will contain pageSize records, which is 100 by default.

// If there are more records, the response will contain an offset.
// To fetch the next page of records, include offset in the next request's parameters.

// Pagination will stop when you've reached the end of your table.
// If the maxRecords parameter is passed, pagination will stop once you've reached this maximum.
func (grc *GetRecordsConfig) WithOffset(offset string) *GetRecordsConfig {
	grc.params.Set("offset", offset)
	return grc
}

// InStringFormat add parameter to get records in string format
// it require timezone
// https://support.airtable.com/hc/en-us/articles/216141558-Supported-timezones-for-SET-TIMEZONE
// and user locale data
// https://support.airtable.com/hc/en-us/articles/220340268-Supported-locale-modifiers-for-SET-LOCALE
func (grc *GetRecordsConfig) InStringFormat(timeZone, userLocale string) *GetRecordsConfig {
	grc.params.Set("cellFormat", "string")
	grc.params.Set("timeZone", timeZone)
	grc.params.Set("userLocale", userLocale)
	return grc
}

// Do send the prepared get records request
func (grc *GetRecordsConfig) Do() (*Records, error) {
	return grc.table.GetRecordsWithParams(grc.params)
}
