package database

import (
	"fmt"
	"strings"

	"github.com/zitadel/logging"
)

type TextFilter[T text] interface {
	Condition
}

type TextCondition[T text] struct {
	Filter[textCompare, T]
}

func NewTextEqual[T text](t T) *TextCondition[T] {
	return newTextFilter(textEqual, t)
}

func NewTextUnequal[T text](t T) *TextCondition[T] {
	return newTextFilter(textUnequal, t)
}

func NewTextEqualInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textEqualInsensitive, t)
}

func NewTextUnequalInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textUnequalInsensitive, t)
}

func NewTextStartsWith[T text](t T) *TextCondition[T] {
	return newTextFilter(textStartsWith, t)
}

func NewTextStartsWithInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textStartsWithInsensitive, t)
}

func NewTextEndsWith[T text](t T) *TextCondition[T] {
	return newTextFilter(textEndsWith, t)
}

func NewTextEndsWithInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textEndsWithInsensitive, t)
}

func NewTextContains[T text](t T) *TextCondition[T] {
	return newTextFilter(textContains, t)
}

func NewTextContainsInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textContainsInsensitive, t)
}

func newTextFilter[T text](comp textCompare, t T) *TextCondition[T] {
	return &TextCondition[T]{
		Filter: Filter[textCompare, T]{
			comp:  comp,
			value: t,
		},
	}
}

func (f TextCondition[T]) Write(stmt *Statement, columnName string) {
	if f.comp.isInsensitive() {
		f.writeCaseInsensitive(stmt, columnName)
		return
	}
	f.Filter.Write(stmt, columnName)
}

func (f *TextCondition[T]) writeCaseInsensitive(stmt *Statement, columnName string) {
	stmt.Builder.WriteString("LOWER(")
	stmt.Builder.WriteString(columnName)
	stmt.Builder.WriteString(") ")
	stmt.Builder.WriteString(f.comp.String())
	f.writeArg(stmt)
}

func (f *TextCondition[T]) writeArg(stmt *Statement) {
	// TODO: condition must know if it's args are named parameters or not
	// var v any = f.value
	// workaround for placeholder
	// if placeholder, ok := v.(placeholder); ok {
	// 	stmt.Builder.WriteString(" LOWER(")
	// 	stmt.WriteArg(placeholder)
	// 	stmt.Builder.WriteString(")")
	// }
	stmt.WriteArg(strings.ToLower(fmt.Sprint(f.value)))
}

type textCompare uint8

const (
	textEqual textCompare = iota
	textUnequal
	textEqualInsensitive
	textUnequalInsensitive
	textStartsWith
	textStartsWithInsensitive
	textEndsWith
	textEndsWithInsensitive
	textContains
	textContainsInsensitive
)

func (c textCompare) String() string {
	switch c {
	case textEqual, textEqualInsensitive:
		return "="
	case textUnequal, textUnequalInsensitive:
		return "<>"
	case textStartsWith, textStartsWithInsensitive, textEndsWith, textEndsWithInsensitive, textContains, textContainsInsensitive:
		return "LIKE"
	default:
		logging.WithFields("compare", c).Panic("comparison type not implemented")
		return ""
	}
}

func (c textCompare) isInsensitive() bool {
	return c == textEqualInsensitive ||
		c == textStartsWithInsensitive ||
		c == textEndsWithInsensitive ||
		c == textContainsInsensitive
}

type text interface {
	~string
	// TODO: condition must know if it's args are named parameters or not
	// ~string | placeholder
}
