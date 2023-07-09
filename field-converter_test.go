// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"reflect"
	"testing"
	"time"
)

func TestToDateTime(t *testing.T) {
	tests := []struct {
		name    string
		field   any
		want    time.Time
		wantErr bool
	}{
		{"not string", any(1), time.Time{}, true},
		{"string not time", any("hello"), time.Time{}, true},
		{"string time", any("2022-03-24T11:12:13.000Z"), time.Date(2022, 0o3, 24, 11, 12, 13, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDateTime(tt.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromDateTime(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want any
	}{
		{"positive", time.Date(2022, 0o3, 24, 11, 12, 13, 1, time.UTC), any("2022-03-24T11:12:13.000Z")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromDateTime(tt.t); got.(string) != tt.want.(string) {
				t.Errorf("FromDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
