package airtable

import (
	"testing"
)

func TestGetBaseSchema(t *testing.T) {
	client := testClient(t)
	baseschema := client.GetBaseSchema("test")
	baseschema.client.baseURL = mockResponse("base_schema.json").URL

	result, err := baseschema.Do()
	if err != nil {
		t.Errorf("there should not be an err, but was: %v", err)
	}
	if len(result.Tables) != 2 {
		t.Errorf("there should be 2 tales, but was %v", len(result.Tables))
	}

	baseschema.client.baseURL = mockErrorResponse(400).URL
	_, err = baseschema.Do()
	if err == nil {
		t.Errorf("there should be an err, but was nil")
	}
}
