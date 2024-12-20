package common

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type FieldData struct {
	Name string
	Type string
	tags []FieldTag
}

func (f FieldData) String() string {
	builder := strings.Builder{}
	builder.WriteString(f.Name + " " + f.Type + " ")
	builder.WriteString("`")
	for _, tag := range f.tags {
		builder.WriteString(tag.String())
		builder.WriteString(" ")
	}
	builder.WriteString("`")
	return builder.String()
}

func (f FieldData) WithoutTags() string {
	return f.Name + " " + f.Type
}

type FieldTag struct {
	Key   string
	Value string
}

func (f FieldTag) String() string {
	return f.Key + `:"` + f.Value + `"`
}

type StructData struct {
	Name   string
	Fields []FieldData
}

type Generator struct {
}

func (Generator) Struct(g *CodeBuffer, items StructData) {
	g.P("type ", items.Name, " struct")
	WithCurls(g, func(b CodeWriter) {
		for _, item := range items.Fields {
			b.L(item.String())
		}
	})
}

func (Generator) PackageHeader(g *CodeBuffer, packageName string, imports ...string) {
	g.L("package ", packageName)
	g.L("")

	if len(imports) > 0 {
		g.L("import (")
		for _, item := range imports {
			g.L(g.Tab(), "\"", item, "\"")
		}
		g.L(")")
		g.L("")
	}
}

func (Generator) Method(g *CodeBuffer, name string, originType *FieldData, returnTypes []string, params []FieldData, bodyWriter func(b CodeWriter)) {
	g.P("func ")

	if originType != nil {
		g.P("(", originType.Name, " ", *&originType.Type, ") ")
	}

	g.P(name, "(")

	// write params

	// TODO: rewrite so it doesn't make line break if it's one param
	for _, item := range params {
		g.L(item.WithoutTags(), ",")
	}
	g.P(")")

	if len(returnTypes) > 0 {
		g.P("(", strings.Join(returnTypes, ", "), ")")
	}

	// write body
	WithCurls(g, bodyWriter)
}

func (Generator) WriteFile(g *CodeBuffer, name string, basePath string) error {
	data := g.buf.String()
	filename := strings.ToLower(name) + ".gen.go"

	// write to file
	err := os.WriteFile(basePath+filename, []byte(data), 0644)
	if err != nil {
		return errors.Wrap(err, "write generated file")
	}

	// go format and check
	cmd := exec.Command("go", "fmt", filename)
	cmd.Stdin = bytes.NewBufferString(data)
	out, err := cmd.Output()
	if err != nil {
		return errors.WithMessagef(err, "go fmt: %s", out)
	}

	// goimports
	cmd = exec.Command("goimports", "-w", filename)
	cmd.Stdin = bytes.NewBufferString(data)
	out, err = cmd.Output()
	if err != nil {
		return errors.WithMessagef(err, "goimports: %s", out)
	}

	return nil
}
