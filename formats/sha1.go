// formats/sha1.go
package formats

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/cydside/combina/combinacore"
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
