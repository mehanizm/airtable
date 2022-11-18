package airtable

import (
	"net/url"
)

// Base type of airtable base.
type Base struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PermissionLevel string `json:"permissionLevel"`
}

// Base type of airtable bases.
type Bases struct {
	Bases  []*Base `json:"bases"`
	Offset string  `json:"offset,omitempty"`
}

type Field struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type View struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type TableSchema struct {
	ID             string   `json:"id"`
	PrimaryFieldID string   `json:"primaryFieldId"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Fields         []*Field `json:"fields"`
	Views          []*View  `json:"views"`
}

type Tables struct {
	Tables []*TableSchema `json:"tables"`
}

// GetBasesWithParams get bases with url values params
// https://airtable.com/developers/web/api/list-bases
func (at *Client) GetBasesWithParams(params url.Values) (*Bases, error) {
	bases := new(Bases)

	err := at.get("meta", "bases", "", params, bases)
	if err != nil {
		return nil, err
	}

	return bases, nil
}

// Table represents table object.
type BaseConfig struct {
	client *Client
	dbId   string
	params url.Values
}

// GetBase return Base object.
func (c *Client) GetBase(dbId string) *BaseConfig {
	return &BaseConfig{
		client: c,
		dbId:   dbId,
	}
}

// Do send the prepared
func (b *BaseConfig) Do() (*Tables, error) {
	return b.GetTablesWithParams()
}

// GetTablesWithParams get tables from a base with url values params
// https://airtable.com/developers/web/api/get-base-schema
func (b *BaseConfig) GetTablesWithParams() (*Tables, error) {
	tables := new(Tables)

	err := b.client.get("meta/bases", b.dbId, "tables", nil, tables)
	if err != nil {
		return nil, err
	}

	return tables, nil
}
