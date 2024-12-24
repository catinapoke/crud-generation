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
	g.Method(buf, MethodData{
		Name:        "GetName",
		OriginType:  &FieldData{Name: "p", Type: "Person"},
		ReturnTypes: []string{"string"},
		Params:      []FieldData{},
	}, func(b CodeWriter) {
		b.L("return p.Name")
	})

	buf.L("")

	g.Method(buf, MethodData{
		Name:        "GetOthersName",
		OriginType:  nil,
		ReturnTypes: []string{"string"},
		Params:      []FieldData{{Name: "p", Type: "Person"}},
	}, func(b CodeWriter) {
		b.L("return p.Name")
	})

	g.Method(buf, MethodData{
		Name:        "Empty",
		OriginType:  nil,
		ReturnTypes: []string{"int"},
		Params:      []FieldData{{Name: "_", Type: "int"}, {Name: "_", Type: "Person"}},
	}, func(b CodeWriter) {
		b.L("return 1")
	})

	err := g.WriteFile(buf, "person", "./")
	if err != nil {
		t.Fatal(err)
	}

	// TODO: add compile test
}
