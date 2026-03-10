// formats/sha1.go
package formats

import (
	"combina/combinacore"
	"crypto/sha1"
	"encoding/hex"
)

type SHA1Formatter struct {
	BaseFormatter
}

func NewSHA1Formatter() *SHA1Formatter {
	return &SHA1Formatter{
		BaseFormatter: BaseFormatter{name: "sha1"},
	}
}

func (f *SHA1Formatter) Format(indices []int, gen combinacore.Generator) (string, error) {
	str := gen.BuildFullString(indices)
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:]), nil
}
