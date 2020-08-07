// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"fmt"
	"net/url"
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
		grc.params.Add("fields", fieldName)
	}
	return grc
}

// SetOffset set records offset to nest request
func (grc *GetRecordsConfig) SetOffset(offset string) *GetRecordsConfig {
	grc.params.Set("offset", offset)
	return grc
}

// WithFilterFormula add filter to request
func (grc *GetRecordsConfig) WithFilterFormula(filterFormula string) *GetRecordsConfig {
	grc.params.Add("filterByFormula", filterFormula)
	return grc
}

// WithSort add sorting to request
func (grc *GetRecordsConfig) WithSort(sortQueries ...struct {
	fieldName string
	direction string
}) *GetRecordsConfig {
	for queryNum, sortQuery := range sortQueries {
		grc.params.Add(fmt.Sprintf("sort[%v][field]", queryNum), sortQuery.fieldName)
		grc.params.Add(fmt.Sprintf("sort[%v][direction]", queryNum), sortQuery.direction)
	}
	return grc
}

// FromView add view parameter to get records
func (grc *GetRecordsConfig) FromView(viewNameOrID string) *GetRecordsConfig {
	grc.params.Add("view", viewNameOrID)
	return grc
}

// InStringFormat add parameter to get records in string format
// it require timezone
// https://support.airtable.com/hc/en-us/articles/216141558-Supported-timezones-for-SET-TIMEZONE
// and user locale data
// https://support.airtable.com/hc/en-us/articles/220340268-Supported-locale-modifiers-for-SET-LOCALE
func (grc *GetRecordsConfig) InStringFormat(timeZone, userLocale string) *GetRecordsConfig {
	grc.params.Add("cellFormat", "string")
	grc.params.Add("timeZone", timeZone)
	grc.params.Add("userLocale", userLocale)
	return grc
}

// Do send the prepared get records request
func (grc *GetRecordsConfig) Do() (*Records, error) {
	return grc.table.GetRecordsWithParams(grc.params)
}
