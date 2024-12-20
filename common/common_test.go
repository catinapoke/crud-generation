package common

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	// Initialize test components
	g := Generator{}
	buf := &CodeBuffer{}

	// Test struct generation
	structData := StructData{
		Name: "Person",
		Fields: []FieldData{
			{
				Name: "ID",
				Type: "int",
				tags: []FieldTag{{Key: "json", Value: "id"}},
			},
			{
				Name: "Name",
				Type: "string",
				tags: []FieldTag{{Key: "json", Value: "name"}},
			},
		},
	}

	g.PackageHeader(buf, "common", "context")

	// Generate struct
	g.Struct(buf, structData)

	// Generate getter method
	g.Method(buf, "GetName", &FieldData{Name: "p", Type: "Person"}, []string{"string"}, []FieldData{}, func(b CodeWriter) {
		b.L("return p.Name")
	})

	buf.L("")

	g.Method(buf, "GetOthersName", nil, []string{"string"}, []FieldData{{Name: "p", Type: "Person"}}, func(b CodeWriter) {
		b.L("return p.Name")
	})

	err := g.WriteFile(buf, "person", "./")
	if err != nil {
		t.Fatal(err)
	}

	// add compile test
}
