package primes

// Gaps returns the distances between consecutive entries of seq.
//
// For a prime sequence the result is the gap sequence g[i] = p[i+1] - p[i].
// Sequences shorter than two entries yield an empty slice.
//
// Only the first gap is odd: 2 is the sole even prime, so every prime past it
// is odd and every distance between them is even. That single arithmetic fact
// drives most of the structure the pattern package measures.
func Gaps(seq []int) []int {
	if len(seq) < 2 {
		return []int{}
	}

	out := make([]int, len(seq)-1)
	for i := range out {
		out[i] = seq[i+1] - seq[i]
	}

	return out
}
