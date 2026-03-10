// combinacore/generator.go
package combinacore

import (
	"fmt"
	"strings"
)

// GeneratorFactory crea generatori in base alla configurazione
type GeneratorFactory struct{}

// CreateGenerator crea un generatore appropriato in base alla configurazione
func (f *GeneratorFactory) CreateGenerator(config Config) (Generator, error) {
	// Valida la configurazione
	if len(config.Elements) == 0 {
		return nil, fmt.Errorf("no elements provided")
	}

	if config.Length < 1 {
		return nil, fmt.Errorf("length must be >= 1")
	}

	if config.ProgressiveStart < 1 || config.ProgressiveStart > config.Length {
		return nil, fmt.Errorf("progressive start must be between 1 and length")
	}

	// Verifica vincoli per combinazioni senza ripetizione
	if (config.CombinationSimple || config.PermutationSimple) &&
		config.Length > len(config.Elements) {
		return nil, fmt.Errorf("length cannot exceed number of elements for combinations without repetition")
	}

	// Crea il generatore appropriato
	switch {
	case config.PermutationSimple:
		return &PermutationSimpleGenerator{BaseGenerator: BaseGenerator{config: config}}, nil
	case config.CombinationSimple:
		return &CombinationSimpleGenerator{BaseGenerator: BaseGenerator{config: config}}, nil
	case config.CombinationRepeat:
		return &CombinationRepeatGenerator{BaseGenerator: BaseGenerator{config: config}}, nil
	default: // PermutationRepeat
		return &PermutationRepeatGenerator{BaseGenerator: BaseGenerator{config: config}}, nil
	}
}

// BaseGenerator fornisce funzionalità comuni a tutti i generatori
type BaseGenerator struct {
	config Config
}

func (g *BaseGenerator) GetElement(idx int) string {
	if idx < 0 || idx >= len(g.config.Elements) {
		return "?"
	}
	return g.config.Elements[idx]
}

func (g *BaseGenerator) GetElementsLen() int {
	return len(g.config.Elements)
}

func (g *BaseGenerator) GetConfig() Config {
	return g.config
}

// BuildString costruisce una stringa dagli indici usando OutputSeparator
func (g *BaseGenerator) BuildString(indices []int) string {
	if g.config.Mode == ModeChars {
		// Modalità caratteri: concatenazione diretta (ignora separatori)
		var result strings.Builder
		for _, idx := range indices {
			if idx >= 0 && idx < len(g.config.Elements) {
				result.WriteString(g.config.Elements[idx])
			}
		}
		return result.String()
	}

	// Modalità parole: join con OutputSeparator
	parts := make([]string, len(indices))
	for i, idx := range indices {
		if idx >= 0 && idx < len(g.config.Elements) {
			parts[i] = g.config.Elements[idx]
		} else {
			parts[i] = "?"
		}
	}
	return strings.Join(parts, g.config.OutputSeparator)
}

// BuildFullString costruisce la stringa completa con prefisso e suffisso
func (g *BaseGenerator) BuildFullString(indices []int) string {
	var result strings.Builder

	// Aggiungi prefisso
	result.WriteString(g.config.Prefix)

	// Aggiungi la combinazione
	result.WriteString(g.BuildString(indices))

	// Aggiungi suffisso
	result.WriteString(g.config.Suffix)

	return result.String()
}
