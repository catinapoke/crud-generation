package common

import (
	"bytes"
	"fmt"
)

type CodeWriter interface {
	P(args ...interface{})
	L(args ...interface{})
}

type CodeBuffer struct {
	buf bytes.Buffer
}

func (h *CodeBuffer) P(args ...interface{}) {
	fmt.Fprint(&h.buf, args...)
}

func (h *CodeBuffer) L(args ...interface{}) {
	h.P(args...)
	h.P("\n")
}

func (h *CodeBuffer) Tab() string {
	return "\t"
}

type FormatOption func(writer CodeWriter, args ...interface{})

func TabWriter(writer CodeWriter, _ ...interface{}) {
	writer.P("\t")
}

type FormattedHelper struct {
	base      CodeWriter
	formatter []FormatOption
}

func NewFormattedHelper(h CodeWriter, formatter ...FormatOption) *FormattedHelper {
	return &FormattedHelper{
		base:      h,
		formatter: formatter,
	}
}

func (h *FormattedHelper) P(args ...interface{}) {
	if len(h.formatter) > 0 {
		for i := len(h.formatter) - 1; i > -1; i-- {
			h.formatter[i](h.base, args...)
		}
	}

	h.base.P(args...)
}

func (h *FormattedHelper) L(args ...interface{}) {
	if len(h.formatter) > 0 {
		for i := len(h.formatter) - 1; i > -1; i-- {
			h.formatter[i](h.base, args...)
		}
	}

	h.base.L(args...)
}

func WithCurls(h CodeWriter, f func(h CodeWriter)) {
	h.L("{")
	f(NewFormattedHelper(h, TabWriter))
	h.L("}")
}

// TODO: remake WithInComment, so it can be used with Generator
func WithInComment(h CodeWriter, f func(h CodeWriter)) {
	h.L("/*")
	f(h)
	h.L("*/")
}
