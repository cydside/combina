// formats/format.go
package formats

import (
	"fmt"

	"github.com/cydside/combina/combinacore"
)

// Formatter è l'interfaccia per i formattatori
type Formatter interface {
	Format(indices []int, gen combinacore.Generator) (string, error)
	Name() string
}

// BaseFormatter fornisce funzionalità comuni
type BaseFormatter struct {
	name string
}

func (f *BaseFormatter) Name() string {
	return f.name
}

// FormatterFactory crea formattatori
type FormatterFactory struct{}

// CreateFormatter crea un formattatore in base al tipo
func (f *FormatterFactory) CreateFormatter(formatType string) (Formatter, error) {
	switch formatType {
	case "plain", "":
		return NewPlainFormatter(), nil
	case "md5":
		return NewMD5Formatter(), nil
	case "sha1":
		return NewSHA1Formatter(), nil
	case "hex":
		return NewHexFormatter(), nil
	default:
		return nil, fmt.Errorf("unknown format type: %s", formatType)
	}
}
