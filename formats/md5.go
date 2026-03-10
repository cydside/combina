// formats/md5.go
package formats

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/cydside/combina/combinacore"
)

type MD5Formatter struct {
	BaseFormatter
}

func NewMD5Formatter() *MD5Formatter {
	return &MD5Formatter{
		BaseFormatter: BaseFormatter{name: "md5"},
	}
}

func (f *MD5Formatter) Format(indices []int, gen combinacore.Generator) (string, error) {
	str := gen.BuildFullString(indices)
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:]), nil
}
