package airtable

import (
	"errors"
	"reflect"
	"testing"
)

func TestTable_UploadAttachment(t *testing.T) {
	table := testTable()
	table.client.uploadAttachmentBaseURL = mockResponse("upload_attachment.json").URL
	attachment := Attachment{
		ContentType: "image/png",
		File:        "base64encodedstring",
		FileName:    "test.png",
	}
	fieldAttachments, err := table.UploadAttachment("recnTq6CsvFM6vX2m", "Attachments", attachment)
	if err != nil {
		t.Errorf("must be no error, but: %v", err)
	}
	expected := &FieldAttachments{
		client:      table.client,
		table:       table,
		Id:          "recnTq6CsvFM6vX2m",
		CreatedTime: "2020-04-10T11:30:57.000Z",
		Attachments: map[string][]FieldAttachmentDetails{
			"Attachments": {
				{
					Id:       "att1",
					URL:      "https://example.com/test.png",
					FileName: "test.png",
					Size:     12345,
					Type:     "image/png",
				},
			},
		},
	}
	if !reflect.DeepEqual(fieldAttachments, expected) {
		t.Errorf("expected: %#v\nbut got: %#v\n", expected, fieldAttachments)
	}
	table.client.baseURL = mockErrorResponse(404).URL
	_, err = table.UploadAttachment("recnTq6CsvFM6vX2m", "Attachments", attachment)
	var e *HTTPClientError
	if errors.Is(err, e) {
		t.Errorf("should be an http error, but was not: %v", err)
	}
}
