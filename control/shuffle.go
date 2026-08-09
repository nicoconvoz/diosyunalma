// Package control turns a measurement into evidence.
//
// A detector that reports "I found 2994 of these" says nothing on its own. The
// question is always: how many would a sequence with no such structure have
// produced? This package builds those decoy sequences and scores the real
// measurement against them.
package control

import "math/rand"

// ShuffleGaps returns a permutation of gaps: the same distances in a random
// order. The input is left untouched.
//
// This is the sharp control for any claim about the ARRANGEMENT of primes. The
// decoy keeps the exact multiset of real gaps, so the gap distribution, the
// mean, the variance and the jumping champion are all identical. The only thing
// destroyed is the order. Anything that survives the shuffle was a property of
// the distribution; anything that collapses was a property of the arrangement.
func ShuffleGaps(gaps []int, rng *rand.Rand) []int {
	out := make([]int, len(gaps))
	copy(out, gaps)

	rng.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})

	return out
}
