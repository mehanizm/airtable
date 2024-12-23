package airtable

import "context"

type Attachment struct {
	ContentType string `json:"contentType"`
	File        string `json:"file"`
	FileName    string `json:"filename"`
}

type FieldAttachments struct {
	client *Client
	table  *Table
	// Parent record's ID of the attached file.
	Id          string `json:"id"`
	CreatedTime string `json:"createdTime"`
	// Mapped array of attachments attached to this field.
	//
	// Key is the ID of the record's attachment field and the value is the list of attached files.
	Attachments map[string][]FieldAttachmentDetails `json:"fields"`
}

type FieldAttachmentDetails struct {
	// Airtable's attachment ID
	Id       string `json:"id"`
	URL      string `json:"url"`
	FileName string `json:"filename"`
	// In bytes
	Size int `json:"size"`
	// Content-Type value
	Type string `json:"type"`
}

func (t *Table) UploadAttachment(recordID string, attachmentFieldIdOrName string, data Attachment) (*FieldAttachments, error) {
	return t.UploadAttachmentContext(context.Background(), recordID, attachmentFieldIdOrName, data)
}
func (t *Table) UploadAttachmentContext(ctx context.Context, recordID string, attachmentFieldIdOrName string, data Attachment) (*FieldAttachments, error) {
	result := new(FieldAttachments)

	err := t.client.postAttachment(ctx, t.dbName, recordID, attachmentFieldIdOrName, data, result)
	if err != nil {
		return nil, err
	}

	result.client = t.client
	result.table = t

	return result, nil
}
