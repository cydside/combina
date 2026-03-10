// combinacore/core.go
package combinacore

// Mode rappresenta la modalità di generazione
type Mode int

const (
	ModeChars Mode = iota
	ModeWords
)

// Config contiene la configurazione per il generatore
type Config struct {
	Mode              Mode
	CombinationSimple bool
	CombinationRepeat bool
	PermutationSimple bool
	PermutationRepeat bool
	Length            int
	ProgressiveStart  int
	Elements          []string
	InputSeparator    string // Come separare la frase in input
	OutputSeparator   string // Come unire le parole nell'output
	Prefix            string
	Suffix            string
}

// Generator è l'interfaccia che tutti i generatori devono implementare
type Generator interface {
	Generate() (<-chan []int, error)
	GetElement(idx int) string
	GetElementsLen() int
	GetConfig() Config
	TotalCombinations() int64
	BuildString(indices []int) string
	BuildFullString(indices []int) string
}

// Formatter è l'interfaccia per formattare i risultati
type Formatter interface {
	Format(indices []int, gen Generator) (string, error)
	Name() string
}

// NewConfig crea una nuova configurazione con valori predefiniti
func NewConfig() Config {
	return Config{
		Mode:              ModeChars,
		PermutationRepeat: true,
		Length:            3,
		ProgressiveStart:  1,
		InputSeparator:    " ", // Default: spazio per input
		OutputSeparator:   " ", // Default: spazio per output
		Prefix:            "",
		Suffix:            "",
	}
}
