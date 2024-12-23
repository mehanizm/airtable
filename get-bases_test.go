package airtable

import (
	"testing"
)

func TestGetBases_Do(t *testing.T) {
	client := testClient()
	bases := client.GetBases()
	bases.client.baseURL = mockResponse("get_bases.json").URL

	result, err := bases.WithOffset("0").Do()
	if err != nil {
		t.Errorf("there should not be an err, but was: %v", err)
	}
	if len(result.Bases) != 2 {
		t.Errorf("there should be 2 bases, but was %v", len(result.Bases))
	}

	bases.client.baseURL = mockErrorResponse(400).URL
	_, err = bases.Do()
	if err == nil {
		t.Errorf("there should be an err, but was nil")
	}
}
