// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"errors"
	"time"
)

const (
	dateTimeFormat = "2006-01-02T15:04:05.000Z"
)

var ErrNotDateTime = errors.New("field is not date time")

func ToDateTime(field any) (time.Time, error) {
	fS, err := field.(string)
	if !err {
		return time.Time{}, ErrNotDateTime
	}
	return time.Parse(dateTimeFormat, fS)
}

func FromDateTime(t time.Time) any {
	return t.Format(dateTimeFormat)
}
