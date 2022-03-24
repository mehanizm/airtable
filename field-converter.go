package airtable

import (
	"errors"
	"time"
)

const (
	dateTimeFormat = "2006-01-02T15:04:05.000Z"
)

var ErrNotDateTime = errors.New("field is not date time")

func ToDateTime(field interface{}) (time.Time, error) {
	fS, err := field.(string)
	if !err {
		return time.Time{}, ErrNotDateTime
	}
	return time.Parse(dateTimeFormat, fS)
}

func FromDateTime(t time.Time) interface{} {
	return t.Format(dateTimeFormat)
}
