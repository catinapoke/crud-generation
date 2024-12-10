package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type InputItems struct {
	Name         string
	DatabaseName string
	Type         string
	IsPrimaryKey bool
}

type InputStruct struct {
	Name         string
	DatabaseName string
	Items        []InputItems
}

func TestGenerateSelect(t *testing.T) {
	g := &Helper{}

	items := InputStruct{
		Name:         "Example",
		DatabaseName: "example",
		Items: []InputItems{
			{Name: "id", DatabaseName: "id", Type: "int", IsPrimaryKey: true},
			{Name: "item", DatabaseName: "item", Type: "string", IsPrimaryKey: false},
		},
	}

	// generate header
	g.L("package main")
	g.L("")
	g.L("import (")
	g.L(g.Tab(), `"context"`)
	g.L(")")
	g.L("")

	// generate struct
	g.L("type ", items.Name, " struct {")
	for _, item := range items.Items {
		g.L("    ", item.Name, " ", item.Type, " `db:\"", item.DatabaseName, "\"`")
	}
	g.L("}")

	// generate query
	g.L("const(")
	g.P(g.Tab(), "queryGet", items.Name, " = `select ")
	for i, item := range items.Items {
		g.P(item.DatabaseName)
		if i != len(items.Items)-1 {
			g.P(", ")
		}
	}
	g.L("")

	g.P(g.Tab(), "from ", items.DatabaseName)
	g.L("")
	g.P(g.Tab(), "where ")

	primary := []InputItems{}
	for _, item := range items.Items {
		if item.IsPrimaryKey {
			primary = append(primary, item)
		}
	}

	count := 1
	for i, item := range primary {
		g.P(item.DatabaseName, "=$", count)
		count++

		if i != len(primary)-1 {
			g.P(" AND ")
		}
	}
	g.L("")
	g.L("`")
	g.L(")")

	// generate get method
	g.P("func Get", items.Name, "(ctx context.Context, db Querier, ")
	for i, item := range primary {
		g.P(item.Name, " ", item.Type)
		if i != len(primary)-1 {
			g.P(", ")
		}
	}
	g.L(") (", items.Name, ", error) {")
	g.L(g.Tab(), "var item ", items.Name)
	g.P(g.Tab(), "if err := db.QueryRowContext(ctx, queryGet", items.Name, ", ")
	for _, item := range primary {
		g.P(item.Name, ", ")
	}
	g.L(").Scan(")
	for _, item := range items.Items {
		g.L(g.Tab(), g.Tab(), "&item.", item.Name, ", ")
	}
	g.L(g.Tab(), "); err != nil {")
	g.L(g.Tab(), g.Tab(), "return ", items.Name, "{}, err")
	g.L(g.Tab(), "}")
	g.L(g.Tab(), "return item, nil")
	g.L("}")

	// write to file
	data := g.buf.String()
	err := os.WriteFile(strings.ToLower(items.Name)+".gen.go", []byte(data), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// go format and check
	cmd := exec.Command("go", "fmt", strings.ToLower(items.Name)+".gen.go")
	cmd.Stdin = bytes.NewBufferString(data)
	out, err := cmd.Output()
	if err != nil {
		t.Error(err, out)
	}

	//goimports
	cmd = exec.Command("goimports", "-w", strings.ToLower(items.Name)+".gen.go")
	cmd.Stdin = bytes.NewBufferString(data)
	out, err = cmd.Output()
	if err != nil {
		t.Error(err, out)
	}

	t.Log(string(data))
}
