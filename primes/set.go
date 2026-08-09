package primes

// Set is a membership lookup for primality up to a fixed limit.
//
// Sieve returns the primes as a list, which is what most measurements want.
// Some questions instead ask "is this particular number prime", millions of
// times, and a list forces a search. Set answers in constant time.
type Set []bool

// NewSet builds a primality lookup covering 0 through limit inclusive.
func NewSet(limit int) Set {
	if limit < 0 {
		return Set{}
	}

	prime := make(Set, limit+1)
	for n := 2; n <= limit; n++ {
		prime[n] = true
	}
	for p := 2; p*p <= limit; p++ {
		if !prime[p] {
			continue
		}
		for m := p * p; m <= limit; m += p {
			prime[m] = false
		}
	}

	return prime
}

// Has reports whether n is prime. Values outside the set's range are false.
func (s Set) Has(n int) bool {
	if n < 0 || n >= len(s) {
		return false
	}
	return s[n]
}

// ProgressionTriples counts the values p for which p, p+d and p+2d are all
// prime and p+2d stays inside the set.
//
// This is the quantity the Hardy–Littlewood k-tuple conjecture predicts. It
// makes no demand that the three primes be CONSECUTIVE — other primes may sit
// inside either interval.
//
// Comparing this against the count of two equal consecutive gaps isolates the
// cost of consecutiveness on its own, which is the only part of the measured
// deficit the singular series does not reach.
func ProgressionTriples(s Set, d int) int {
	if d < 1 {
		return 0
	}

	count := 0
	for p := 2; p+2*d < len(s); p++ {
		if s[p] && s[p+d] && s[p+2*d] {
			count++
		}
	}

	return count
}
