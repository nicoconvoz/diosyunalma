// Package primes provides the primitives used to explore the distribution of
// prime numbers: generation, gap analysis, and density comparisons.
package primes

// Sieve returns every prime less than or equal to limit, in ascending order.
//
// It implements the sieve of Eratosthenes, which runs in O(n log log n) time
// and allocates one byte per candidate. Limits below 2 yield an empty slice.
func Sieve(limit int) []int {
	if limit < 2 {
		return []int{}
	}

	composite := make([]bool, limit+1)

	// Every multiple of p below p*p already carries a smaller prime factor and
	// has been marked, so marking can start at p*p. Once p*p passes the limit
	// there is nothing left to mark.
	for p := 2; p*p <= limit; p++ {
		if composite[p] {
			continue
		}
		for multiple := p * p; multiple <= limit; multiple += p {
			composite[multiple] = true
		}
	}

	primes := make([]int, 0)
	for n := 2; n <= limit; n++ {
		if !composite[n] {
			primes = append(primes, n)
		}
	}

	return primes
}
