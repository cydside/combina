// cmd/combina/main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"combina/combinacore"
	"combina/formats"
)

func main() {
	// Parsing argomenti
	var (
		// Modalità caratteri
		alphaL    = flag.Bool("a", false, "lowercase letters")
		alphaU    = flag.Bool("A", false, "uppercase letters")
		numeric   = flag.Bool("n", false, "numbers")
		special   = flag.Bool("s", false, "special chars")
		userChars = flag.String("user", "", "custom charset")

		// Modalità parole
		phrase          = flag.String("phrase", "", "phrase to process")
		inputSeparator  = flag.String("input-separator", " ", "separator for splitting input phrase")
		outputSeparator = flag.String("output-separator", " ", "separator for joining words in output")

		// Prefisso e suffisso
		prefix = flag.String("prefix", "", "prefix to add before each combination")
		suffix = flag.String("suffix", "", "suffix to add after each combination")

		// Tipo combinazione
		simpleComb = flag.Bool("c", false, "combination without repetition")
		repeatComb = flag.Bool("m", false, "combination with repetition")
		simplePerm = flag.Bool("d", false, "permutation without repetition")
		repeatPerm = flag.Bool("r", false, "permutation with repetition")

		// Parametri
		length = flag.Int("k", 3, "password length")
		start  = flag.Int("p", 1, "progressive start length")

		// Formattazione
		format = flag.String("format", "plain", "output format (plain, md5, sha1, hex)")

		// Concorrenza
		workers = flag.Int("workers", 0, "number of workers (0 = auto)")
		verbose = flag.Bool("verbose", false, "show progress")

		// Altro
		help    = flag.Bool("help", false, "show help")
		version = flag.Bool("version", false, "show version")
	)

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *version {
		printVersion()
		return
	}

	// Determina modalità e crea elementi
	var mode combinacore.Mode
	var elements []string

	if *phrase != "" {
		mode = combinacore.ModeWords
		elements = splitPhrase(*phrase, *inputSeparator)
		fmt.Fprintf(os.Stderr, "Debug: Input phrase split with separator %q: %v\n",
			*inputSeparator, elements)
	} else {
		mode = combinacore.ModeChars
		elements = buildCharset(*alphaL, *alphaU, *numeric, *special, *userChars)
	}

	if len(elements) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no elements selected\n")
		os.Exit(1)
	}

	// Crea configurazione
	config := combinacore.NewConfig()
	config.Mode = mode
	config.CombinationSimple = *simpleComb
	config.CombinationRepeat = *repeatComb
	config.PermutationSimple = *simplePerm
	config.PermutationRepeat = *repeatPerm
	config.Length = *length
	config.ProgressiveStart = *start
	config.Elements = elements
	config.InputSeparator = *inputSeparator
	config.OutputSeparator = *outputSeparator
	config.Prefix = *prefix
	config.Suffix = *suffix

	// Se nessun tipo specificato, usa permutazione con ripetizione
	if !config.CombinationSimple && !config.CombinationRepeat &&
		!config.PermutationSimple && !config.PermutationRepeat {
		config.PermutationRepeat = true
	}

	// Crea generatore
	factory := &combinacore.GeneratorFactory{}
	gen, err := factory.CreateGenerator(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Crea formatter
	formatFactory := &formats.FormatterFactory{}
	formatter, err := formatFactory.CreateFormatter(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Imposta workers
	numWorkers := *workers
	if numWorkers == 0 {
		numWorkers = runtime.NumCPU()
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Generating %d combinations using %d workers\n",
			gen.TotalCombinations(), numWorkers)
		fmt.Fprintf(os.Stderr, "Elements: %v\n", elements)
		fmt.Fprintf(os.Stderr, "Output separator: %q\n", *outputSeparator)
		if *prefix != "" {
			fmt.Fprintf(os.Stderr, "Prefix: %q\n", *prefix)
		}
		if *suffix != "" {
			fmt.Fprintf(os.Stderr, "Suffix: %q\n", *suffix)
		}
	}

	// Genera risultati
	if numWorkers > 1 {
		err = generateConcurrent(gen, formatter, numWorkers, *verbose)
	} else {
		err = generateSequential(gen, formatter, *verbose)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func generateSequential(gen combinacore.Generator, formatter formats.Formatter, verbose bool) error {
	ch, err := gen.Generate()
	if err != nil {
		return err
	}

	count := int64(0)
	for indices := range ch {
		result, err := formatter.Format(indices, gen)
		if err != nil {
			return err
		}
		fmt.Println(result)
		count++

		if verbose && count%10000 == 0 {
			fmt.Fprintf(os.Stderr, "\rGenerated %d combinations", count)
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "\nDone! Total: %d\n", count)
	}

	return nil
}

func generateConcurrent(gen combinacore.Generator, formatter formats.Formatter, workers int, verbose bool) error {
	ch, err := gen.Generate()
	if err != nil {
		return err
	}

	results := make(chan string, 10000)
	done := make(chan bool)
	var wg sync.WaitGroup

	// Worker pool
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for indices := range ch {
				result, err := formatter.Format(indices, gen)
				if err == nil {
					results <- result
				}
			}
		}(w)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collector
	go func() {
		count := int64(0)
		for result := range results {
			fmt.Println(result)
			count++

			if verbose && count%10000 == 0 {
				fmt.Fprintf(os.Stderr, "\rGenerated %d combinations", count)
			}
		}
		done <- true
	}()

	// Attendi completamento
	<-done
	return nil
}

func buildCharset(alphaL, alphaU, numeric, special bool, user string) []string {
	var elements []string

	if alphaL {
		elements = append(elements, strings.Split("abcdefghijklmnopqrstuvwxyz", "")...)
	}
	if alphaU {
		elements = append(elements, strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")...)
	}
	if numeric {
		elements = append(elements, strings.Split("0123456789", "")...)
	}
	if special {
		elements = append(elements, strings.Split("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", "")...)
	}
	if user != "" {
		// Tratta la stringa utente come caratteri individuali
		for _, r := range user {
			elements = append(elements, string(r))
		}
	}

	return elements
}

func splitPhrase(phrase, separator string) []string {
	if separator == " " {
		// Per spazio, usa strings.Fields che gestisce spazi multipli
		return strings.Fields(phrase)
	}
	// Per separatore personalizzato, split semplice
	parts := strings.Split(phrase, separator)
	// Rimuovi spazi vuoti
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func printHelp() {
	fmt.Println("combina - Combinatorial generator")
	fmt.Println("\nUsage: combina [options]")
	fmt.Println("\nCharacter set options:")
	fmt.Println("  -a            lowercase letters")
	fmt.Println("  -A            uppercase letters")
	fmt.Println("  -n            numbers")
	fmt.Println("  -s            special characters")
	fmt.Println("  -user string  custom charset")
	fmt.Println("\nPhrase options:")
	fmt.Println("  -phrase string        phrase to process")
	fmt.Println("  -input-separator s    separator for splitting input (default \" \")")
	fmt.Println("  -output-separator s   separator for joining output (default \" \")")
	fmt.Println("\nPrefix/Suffix options:")
	fmt.Println("  -prefix string  prefix to add before each combination")
	fmt.Println("  -suffix string  suffix to add after each combination")
	fmt.Println("\nCombination type:")
	fmt.Println("  -c            combination without repetition")
	fmt.Println("  -m            combination with repetition")
	fmt.Println("  -d            permutation without repetition")
	fmt.Println("  -r            permutation with repetition (default)")
	fmt.Println("\nParameters:")
	fmt.Println("  -k int        length (default 3)")
	fmt.Println("  -p int        progressive start (default 1)")
	fmt.Println("\nOutput format:")
	fmt.Println("  -format f     plain, md5, sha1, hex (default plain)")
	fmt.Println("\nPerformance:")
	fmt.Println("  -workers n    number of workers (0 = auto)")
	fmt.Println("  -verbose      show progress")
	fmt.Println("\nOther:")
	fmt.Println("  -help         show this help")
	fmt.Println("  -version      show version")
	fmt.Println("\nExamples:")
	fmt.Println("  combina -phrase \"rosso,verde,blu\" -input-separator \",\" -output-separator \" \" -k 2 -r")
	fmt.Println("  combina -phrase \"gatto cane topo\" -k 2 -r -prefix \"[\" -suffix \"]\"")
}

func printVersion() {
	fmt.Println("combina version 1.2.0")
}
