package primes

import "math"

// segmentSize is the window the segmented sieve marks at a time. Eight
// megabytes keeps the whole segment in cache-friendly territory while the
// base primes are reused across every segment.
const segmentSize = 1 << 23

// PsiSegmented returns Chebyshev's psi at each of xs, which must be ascending,
// without ever materialising a sieve of the full range.
//
// PsiAt allocates one byte per candidate, which caps it near 10^9 on ordinary
// hardware. This version sieves the range in windows of segmentSize using only
// the base primes up to √limit, so memory stays at megabytes while the limit
// climbs another decade. The telescope: resolution in the zero hunt grows with
// the log of the range, and this buys the next decade of it.
func PsiSegmented(limit int, xs []int) []float64 {
	if len(xs) == 0 {
		return []float64{}
	}

	base := Sieve(int(math.Sqrt(float64(limit))) + 1)

	// Prime powers p^k, k >= 2, contribute ln p at p^k. There are only about
	// √limit of them, so they are collected up front and merged by value.
	type event struct {
		at     int
		weight float64
	}
	powers := []event{}
	for _, p := range base {
		w := math.Log(float64(p))
		for pk := p * p; pk <= limit && pk > 0; pk *= p {
			powers = append(powers, event{pk, w})
		}
	}
	sortEvents := func(v []event) {
		for i := 1; i < len(v); i++ {
			for j := i; j > 0 && v[j].at < v[j-1].at; j-- {
				v[j], v[j-1] = v[j-1], v[j]
			}
		}
	}
	sortEvents(powers)

	out := make([]float64, len(xs))
	sum := 0.0
	xi, pi := 0, 0

	composite := make([]bool, segmentSize)
	for lo := 2; lo <= limit && xi < len(xs); lo += segmentSize {
		hi := lo + segmentSize - 1
		if hi > limit {
			hi = limit
		}
		span := hi - lo + 1
		for i := 0; i < span; i++ {
			composite[i] = false
		}
		for _, p := range base {
			start := (lo + p - 1) / p * p
			if s := p * p; s > start {
				start = s
			}
			for m := start; m <= hi; m += p {
				composite[m-lo] = true
			}
		}

		for n := lo; n <= hi; n++ {
			for xi < len(xs) && xs[xi] < n {
				out[xi] = sum
				xi++
			}
			if xi == len(xs) {
				break
			}
			if !composite[n-lo] {
				sum += math.Log(float64(n))
			}
			for pi < len(powers) && powers[pi].at == n {
				sum += powers[pi].weight
				pi++
			}
		}
	}
	for ; xi < len(xs); xi++ {
		out[xi] = sum
	}

	return out
}
