package main

import (
	"fmt"

	"github.com/mehanizm/airtable"
)

const (
	airtableAPIKey    = "xxx"
	airtableDBName    = "xxx"
	airtableTableName = "xxx"
)

func main() {
	airtableClient := airtable.NewClient(airtableAPIKey)
	airtableTable := airtableClient.GetTable(airtableDBName, airtableTableName)

	offset := ""

	for {
		records, err := airtableTable.GetRecords().
			WithFilterFormula("NOT({SomeBoolColumn})").
			ReturnFields("Column1", "Column2", "Column3", "Column4").
			MaxRecords(100).
			PageSize(10).
			WithOffset(offset).
			Do()
		if err != nil {
			panic(err)
		}

		for recordNum, record := range records.Records {
			fmt.Println("====iteration====")
			fmt.Println(recordNum, record)
		}

		offset = records.Offset
		if offset == "" {
			break
		}
	}

}
