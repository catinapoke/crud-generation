package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/catinapoke/crud-generation/common"
	"github.com/pkg/errors"
)

func NewDatabaseTag(rowname string) common.FieldTag {
	return common.FieldTag{Key: "db", Value: rowname}
}

type EnitityRow struct {
	common.FieldData
	IsPrimaryKey bool
}

type DatabaseEntity struct {
	Name         string
	DatabaseName string
	Items        []EnitityRow
}

func (items DatabaseEntity) GetPrimaryKeys() []EnitityRow {
	primary := []EnitityRow{}
	for _, item := range items.Items {
		if item.IsPrimaryKey {
			primary = append(primary, item)
		}
	}
	return primary
}

func (Generator) SelectQuery(g *CodeBuffer, items DatabaseEntity) {
	g.L("`select ")
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

	primary := items.GetPrimaryKeys()

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
}

func (gen Generator) Generate(items DatabaseEntity) string {
	g := &CodeBuffer{}

	// generate header
	gen.PackageHeader(g, "context")

	// generate struct
	gen.Struct(g, items)

	// generate query
	g.L("const(")
	g.P(g.Tab(), "queryGet", items.Name, " = ")
	gen.SelectQuery(g, items)
	g.L(")")

	// generate get method
	g.P("func Get", items.Name, "(ctx context.Context, db Querier, ")
	primary := items.GetPrimaryKeys()
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

	return g.buf.String()
}

func (g *Generator) GenerateInFile(items DatabaseEntity, basePath string) error {
	data := g.Generate(items)

	// write to file
	err := os.WriteFile(basePath+strings.ToLower(items.Name)+".gen.go", []byte(data), 0644)
	if err != nil {
		return errors.Wrap(err, "write generated file")
	}

	// go format and check
	cmd := exec.Command("go", "fmt", strings.ToLower(items.Name)+".gen.go")
	cmd.Stdin = bytes.NewBufferString(data)
	out, err := cmd.Output()
	if err != nil {
		return errors.WithMessagef(err, "go fmt: %s", out)
	}

	//goimports
	cmd = exec.Command("goimports", "-w", strings.ToLower(items.Name)+".gen.go")
	cmd.Stdin = bytes.NewBufferString(data)
	out, err = cmd.Output()
	if err != nil {
		return errors.WithMessagef(err, "goimports: %s", out)
	}

	return nil
}
