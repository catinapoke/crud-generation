package main

import (
	"testing"
)

func TestGenerateSelect(t *testing.T) {
	g := Generator{}

	items := InputStruct{
		Name:         "Example",
		DatabaseName: "example",
		Items: []InputItems{
			{Name: "id", DatabaseName: "id", Type: "int", IsPrimaryKey: true},
			{Name: "item", DatabaseName: "item", Type: "string", IsPrimaryKey: false},
		},
	}

	err := g.GenerateInFile(items, "./")
	if err != nil {
		t.Fatal(err)
	}
}
