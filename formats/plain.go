// formats/plain.go
package formats

import (
	"github.com/cydside/combina/combinacore"
)

type PlainFormatter struct {
	BaseFormatter
}

func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{
		BaseFormatter: BaseFormatter{name: "plain"},
	}
}

func (f *PlainFormatter) Format(indices []int, gen combinacore.Generator) (string, error) {
	return gen.BuildFullString(indices), nil
}
