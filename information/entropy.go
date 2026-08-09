// Package information measures how many bits a sequence costs to describe.
//
// It exists to give every finding in this project a common currency. A result
// that says "primes avoid repeating gaps" and a result that says "the residue
// chain remembers three steps" cannot be compared as they stand. Expressed as
// bits recovered per prime, they can — and a mechanism that recovers none can
// be discarded on the spot.
package information

import (
	"math"
	"strconv"
	"strings"
)

// Shannon returns the entropy in bits of a distribution given as raw counts.
//
// Zero counts contribute nothing. An empty or all-zero distribution is 0: with
// no observations there is no uncertainty to report.
func Shannon(counts []int) float64 {
	total := 0
	for _, c := range counts {
		if c > 0 {
			total += c
		}
	}
	if total == 0 {
		return 0
	}

	h := 0.0
	for _, c := range counts {
		if c <= 0 {
			continue
		}
		p := float64(c) / float64(total)
		h -= p * math.Log2(p)
	}

	return h
}

// ConditionalEntropy returns H(X_n | the previous `order` symbols) in bits —
// how much uncertainty about the next symbol survives after seeing that much
// history. Order 0 is the plain marginal entropy.
//
// A word of warning about the number it returns. Conditional entropy estimated
// from a finite sample is biased DOWNWARDS: as the history grows, each distinct
// history is seen fewer times, and rare histories look more deterministic than
// they are. The apparent structure is partly the estimator running out of data.
//
// The fix is the one this project uses everywhere else: run the same
// measurement on a shuffled control, where the true conditional entropy equals
// the marginal by construction. Whatever drop the control shows is the bias,
// and the real drop is what exceeds it.
func ConditionalEntropy(seq []int, order int) float64 {
	if order < 0 || len(seq) == 0 || order >= len(seq) {
		return 0
	}

	if order == 0 {
		return Shannon(tally(seq))
	}

	// Group the observed symbols by the history that preceded them.
	byHistory := map[string]map[int]int{}
	weight := map[string]int{}
	for i := order; i < len(seq); i++ {
		key := historyKey(seq[i-order : i])
		if byHistory[key] == nil {
			byHistory[key] = map[int]int{}
		}
		byHistory[key][seq[i]]++
		weight[key]++
	}

	total := len(seq) - order
	h := 0.0
	for key, dist := range byHistory {
		counts := make([]int, 0, len(dist))
		for _, c := range dist {
			counts = append(counts, c)
		}
		h += float64(weight[key]) / float64(total) * Shannon(counts)
	}

	return h
}

// MutualInformation returns how many bits the previous `order` symbols reveal
// about the next one: the marginal entropy less what survives conditioning.
//
// This is the value of the memory, not merely evidence that memory exists.
func MutualInformation(seq []int, order int) float64 {
	return ConditionalEntropy(seq, 0) - ConditionalEntropy(seq, order)
}

func tally(seq []int) []int {
	index := map[int]int{}
	counts := []int{}
	for _, v := range seq {
		pos, seen := index[v]
		if !seen {
			pos = len(counts)
			index[v] = pos
			counts = append(counts, 0)
		}
		counts[pos]++
	}
	return counts
}

func historyKey(window []int) string {
	parts := make([]string, len(window))
	for i, v := range window {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}
