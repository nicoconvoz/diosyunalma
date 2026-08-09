package primes

import "testing"

// TestProgressionPairsCountsAllPrimePairs measures the 2-tuple analogue of
// ProgressionTriples: p and p+d both prime, with no requirement that they be
// consecutive. Dividing the consecutive-gap count by this isolates how often a
// prime pair happens to have nothing between it.
func TestProgressionPairsCountsAllPrimePairs(t *testing.T) {
	s := NewSet(20) // primes: 2 3 5 7 11 13 17 19

	tests := []struct {
		name string
		d    int
		want int
	}{
		{name: "d=1 finds only 2,3", d: 1, want: 1},
		{name: "d=2 finds the twin pairs", d: 2, want: 4},
		{name: "d=4 finds 3-7, 7-11, 13-17", d: 4, want: 3},
		{name: "d=6 finds 5-11, 7-13, 11-17, 13-19", d: 6, want: 4},
		{name: "d beyond the range finds nothing", d: 50, want: 0},
		{name: "d below one is not a step", d: 0, want: 0},
		{name: "negative d is not a step", d: -4, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProgressionPairs(s, tt.d); got != tt.want {
				t.Errorf("ProgressionPairs(d=%d) = %d, want %d", tt.d, got, tt.want)
			}
		})
	}
}

// TestProgressionPairsBoundsTheGapCount is the containment the decomposition
// rests on: every consecutive pair with gap d is also an all-prime pair, so the
// unconstrained count can never be smaller.
func TestProgressionPairsBoundsTheGapCount(t *testing.T) {
	const limit = 100_000
	s := NewSet(limit)
	gaps := Gaps(From(Sieve(limit), 5))

	for _, d := range []int{2, 4, 6, 12, 18, 24} {
		consecutive := 0
		for _, g := range gaps {
			if g == d {
				consecutive++
			}
		}
		all := ProgressionPairs(s, d)

		if all < consecutive {
			t.Errorf("d=%d: %d all-prime pairs but %d consecutive; the unconstrained count cannot be smaller",
				d, all, consecutive)
		}
	}
}

// TestProgressionPairsExceedsTriples records the ordering the survival ratios
// depend on: demanding a third term in progression can only remove candidates.
func TestProgressionPairsExceedsTriples(t *testing.T) {
	s := NewSet(100_000)

	for _, d := range []int{2, 6, 12, 30} {
		pairs := ProgressionPairs(s, d)
		triples := ProgressionTriples(s, d)

		if triples > pairs {
			t.Errorf("d=%d: %d triples exceeds %d pairs", d, triples, pairs)
		}
	}
}
