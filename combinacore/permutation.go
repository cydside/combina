// combinacore/permutation.go
package combinacore

import (
	"math"
)

// PermutationRepeatGenerator genera permutazioni con ripetizione
type PermutationRepeatGenerator struct {
	BaseGenerator
}

func (g *PermutationRepeatGenerator) Generate() (<-chan []int, error) {
	ch := make(chan []int, 1000)

	go func() {
		defer close(ch)

		length := g.config.Length
		indices := make([]int, length)
		elementCount := len(g.config.Elements)

		// Invia la prima combinazione
		first := make([]int, length)
		copy(first, indices)
		ch <- first

		for {
			// Genera la prossima combinazione
			j := length - 1
			for j >= 0 && indices[j] == elementCount-1 {
				j--
			}

			if j < 0 {
				break
			}

			indices[j]++
			for k := j + 1; k < length; k++ {
				indices[k] = 0
			}

			// Invia copia degli indici
			next := make([]int, length)
			copy(next, indices)
			ch <- next
		}
	}()

	return ch, nil
}

func (g *PermutationRepeatGenerator) TotalCombinations() int64 {
	return int64(math.Pow(float64(len(g.config.Elements)), float64(g.config.Length)))
}

// PermutationSimpleGenerator genera permutazioni senza ripetizione
type PermutationSimpleGenerator struct {
	BaseGenerator
}

func (g *PermutationSimpleGenerator) Generate() (<-chan []int, error) {
	ch := make(chan []int, 1000)

	go func() {
		defer close(ch)

		length := g.config.Length
		n := len(g.config.Elements)

		// Se length > n, non ci sono permutazioni
		if length > n {
			return
		}

		// Algoritmo di Heap per permutazioni
		indices := make([]int, length)
		for i := 0; i < length; i++ {
			indices[i] = i
		}

		c := make([]int, length)
		for i := range c {
			c[i] = 0
		}

		// Invia prima permutazione
		first := make([]int, length)
		copy(first, indices)
		ch <- first

		i := 0
		for i < length {
			if c[i] < i {
				if i%2 == 0 {
					indices[0], indices[i] = indices[i], indices[0]
				} else {
					indices[c[i]], indices[i] = indices[i], indices[c[i]]
				}

				next := make([]int, length)
				copy(next, indices)
				ch <- next

				c[i]++
				i = 0
			} else {
				c[i] = 0
				i++
			}
		}
	}()

	return ch, nil
}

func (g *PermutationSimpleGenerator) TotalCombinations() int64 {
	n := int64(len(g.config.Elements))
	k := int64(g.config.Length)

	if k > n {
		return 0
	}

	total := int64(1)
	for i := int64(0); i < k; i++ {
		total *= n - i
	}
	return total
}
