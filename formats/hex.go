// formats/hex.go
package formats

import (
	"combina/combinacore"
	"encoding/hex"
)

type HexFormatter struct {
	BaseFormatter
}

func NewHexFormatter() *HexFormatter {
	return &HexFormatter{
		BaseFormatter: BaseFormatter{name: "hex"},
	}
}

func (f *HexFormatter) Format(indices []int, gen combinacore.Generator) (string, error) {
	str := gen.BuildFullString(indices)
	return hex.EncodeToString([]byte(str)), nil
}
