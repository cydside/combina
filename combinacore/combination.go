// combinacore/combination.go
package combinacore

// CombinationRepeatGenerator genera combinazioni con ripetizione
type CombinationRepeatGenerator struct {
	BaseGenerator
}

func (g *CombinationRepeatGenerator) Generate() (<-chan []int, error) {
	ch := make(chan []int, 1000)

	go func() {
		defer close(ch)

		length := g.config.Length
		indices := make([]int, length)
		elementCount := len(g.config.Elements)

		// Invia prima combinazione
		first := make([]int, length)
		copy(first, indices)
		ch <- first

		for {
			j := length - 1
			for j >= 0 && indices[j] == elementCount-1 {
				j--
			}

			if j < 0 {
				break
			}

			indices[j]++
			for k := j + 1; k < length; k++ {
				indices[k] = indices[j]
			}

			next := make([]int, length)
			copy(next, indices)
			ch <- next
		}
	}()

	return ch, nil
}

func (g *CombinationRepeatGenerator) TotalCombinations() int64 {
	n := int64(len(g.config.Elements))
	k := int64(g.config.Length)

	// C(n+k-1, k)
	total := int64(1)
	for i := int64(1); i <= k; i++ {
		total = total * (n + k - i) / i
	}
	return total
}

// CombinationSimpleGenerator genera combinazioni senza ripetizione
type CombinationSimpleGenerator struct {
	BaseGenerator
}

func (g *CombinationSimpleGenerator) Generate() (<-chan []int, error) {
	ch := make(chan []int, 1000)

	go func() {
		defer close(ch)

		length := g.config.Length
		elementCount := len(g.config.Elements)

		// Se length > elementCount, non ci sono combinazioni
		if length > elementCount {
			close(ch)
			return
		}

		indices := make([]int, length)
		for i := 0; i < length; i++ {
			indices[i] = i
		}

		// Invia prima combinazione
		first := make([]int, length)
		copy(first, indices)
		ch <- first

		for {
			j := length - 1
			for j >= 0 && indices[j] >= elementCount-length+j {
				j--
			}

			if j < 0 {
				break
			}

			indices[j]++
			for k := j + 1; k < length; k++ {
				indices[k] = indices[k-1] + 1
			}

			next := make([]int, length)
			copy(next, indices)
			ch <- next
		}
	}()

	return ch, nil
}

func (g *CombinationSimpleGenerator) TotalCombinations() int64 {
	n := int64(len(g.config.Elements))
	k := int64(g.config.Length)

	if k > n {
		return 0
	}

	if k > n-k {
		k = n - k
	}

	total := int64(1)
	for i := int64(1); i <= k; i++ {
		total = total * (n - k + i) / i
	}
	return total
}
