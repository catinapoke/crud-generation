package main

import (
	"os"
	"testing"
)

func TestCurls(t *testing.T) {
	g := &CodeBuffer{}

	items := DatabaseEntity{
		Name:         "Example",
		DatabaseName: "example",
		Items: []EnitityRow{
			{Name: "id", DatabaseName: "id", Type: "int", IsPrimaryKey: true},
			{Name: "item", DatabaseName: "item", Type: "string", IsPrimaryKey: false},
		},
	}

	g.L("package main")
	WithInComment(g, func(h CodeWriter) {
		g.P("type ", items.Name, " struct")
		WithCurls(g, func(h CodeWriter) {
			for _, item := range items.Items {
				h.L(item.Name, " ", item.Type, " `db:\"", item.DatabaseName, "\"`")
			}
		})
	})

	code := g.buf.String()

	if code == "" {
		t.FailNow()
	}

	data := g.buf.String()
	err := os.WriteFile("curls.gen.go", []byte(data), 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(code)
}
