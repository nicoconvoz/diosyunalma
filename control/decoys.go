package control

import (
	"math"
	"math/rand"
	"sort"
)

// OddDecoy returns a sorted set of count values: the number 2 followed by
// distinct random odd numbers up to limit.
//
// It reproduces the two most trivial structural facts about the primes —
// exactly one even member, everything else odd — and nothing more. Detectors
// whose findings survive against this decoy have found something beyond
// those facts; the emirp pairs of Finding 5 did not.
func OddDecoy(limit, count int, rng *rand.Rand) []int {
	if count < 1 {
		return []int{}
	}

	seen := map[int]bool{2: true}
	for len(seen) < count {
		seen[3+2*rng.Intn((limit-1)/2)] = true
	}

	out := make([]int, 0, count)
	for v := range seen {
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}

// CramerDecoy returns 2 followed by every n up to limit that a coin with
// probability 1/ln n admits — Cramér's model of the primes as a random set
// with the right density.
//
// It matches the primes' density and nothing arithmetic. It is the standard
// null model of this project's spectral experiments: any structure present
// here belongs to density alone.
func CramerDecoy(limit int, rng *rand.Rand) []int {
	out := []int{2}
	for n := 3; n <= limit; n++ {
		if rng.Float64() < 1/math.Log(float64(n)) {
			out = append(out, n)
		}
	}
	return out
}
