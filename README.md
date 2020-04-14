Golang Airtable API
================

[![GoDoc](https://godoc.org/github.com/mehanizm/airtable?status.svg)](https://pkg.go.dev/github.com/mehanizm/airtable)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-92%25-brightgreen.svg?longCache=true&style=flat)</a>
<a href='https://goreportcard.com/report/github.com/mehanizm/airtable' target='_blank'>![goreportcard](https://goreportcard.com/badge/github.com/mehanizm/airtable)</a>

A simple #golang package to access the [Airtable API](https://airtable.com/api).

Table of contents
===
- [Golang Airtable API](#golang-airtable-api)
- [Table of contents](#table-of-contents)
  - [Installation](#installation)
  - [Basic usage](#basic-usage)
    - [Initialize client](#initialize-client)
    - [Get table](#get-table)
    - [List records](#list-records)
    - [Add records](#add-records)
    - [Get record by ID](#get-record-by-id)
    - [Update records](#update-records)
    - [Delete record](#delete-record)
  - [Special thanks](#special-thanks)
  

## Installation

The Golang Airtable API has been tested compatible with Go 1.13 on up.

```
go get github.com/mehanizm/airtable
```

## Basic usage

### Initialize client

You should get `your_api_token` in the airtable [account page](https://airtable.com/account)
```Go
client := airtable.NewClient("your_api_token")
```

### Get table

To get the `your_database_ID` you should go to [main API page](https://airtable.com/api) and select the database.

```Go
table := client.GetTable("your_database_ID", "your_table_name")
```

### List records

To get records from the table you can use something like this

```Go
records, err := table.GetRecords().
	FromView("view_1").
	WithFilterFormula("AND({Field1}='value_1',NOT({Field2}='value_2'))").
	WithSort(sortQuery1, sortQuery2).
	ReturnFields("Field1", "Field2").
	InStringFormat("Europe/Moscow", "ru").
	Do()
if err != nil {
	// Handle error
}
```

### Add records

```Go
recordsToSend := &airtable.Records{
    Records: []*airtable.Record{
        {
            Fields: map[string]interface{
                "Field1": "value1",
                "Field2": true,
            },
        },
    },
}
recievedRecords, err := table.AddRecords(recordsToSend)
if err != nil {
	// Handle error
}
```

### Get record by ID

```Go
record, err := table.GetRecord("recordID")
if err != nil {
	// Handle error
}
```

### Update records

To partial update one record

```Go
res, err := record.UpdateRecordPartial(map[string]interface{}{"Field_2": false})
if err != nil {
	// Handle error
}
```

To full update records

```Go
toUpdateRecords := &airtable.Records{
    Records: []*airtable.Record{
        {
            Fields: map[string]interface{
                "Field1": "value1",
                "Field2": true,
            },
        },
        {
            Fields: map[string]interface{
                "Field1": "value1",
                "Field2": true,
            },
        },
    },
}
updatedRecords, err := table.UpdateRecords(toUpdateRecords)
if err != nil {
	// Handle error
}
```

### Delete record

```Go
res, err := record.DeleteRecord()
if err != nil {
	// Handle error
}
```

## Special thanks

Inspired by [Go Trello API](github.com/adlio/trello)