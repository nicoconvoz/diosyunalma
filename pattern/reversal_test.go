package pattern

import "testing"

// TestReversalPairsOnThePrimesBelow100 pins the number that started the emirp
// investigation: exactly four pairs — 13/31, 17/71, 37/73, 79/97.
func TestReversalPairsOnThePrimesBelow100(t *testing.T) {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47,
		53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

	if got := ReversalPairs(primes, 100); got != 4 {
		t.Errorf("ReversalPairs(primes<100) = %d, want 4", got)
	}
}

// TestReversalPairsExcludesPalindromes records the convention: 11 reversed is
// 11, and a number paired with itself is not a pair.
func TestReversalPairsExcludesPalindromes(t *testing.T) {
	if got := ReversalPairs([]int{11, 101, 131}, 200); got != 0 {
		t.Errorf("palindromic values = %d pairs, want 0", got)
	}
}

// TestReversalPairsRespectsTheLimit checks that a partner beyond the limit
// does not count: 13's partner 31 is outside limit 30.
func TestReversalPairsRespectsTheLimit(t *testing.T) {
	if got := ReversalPairs([]int{13, 31}, 30); got != 0 {
		t.Errorf("partner beyond limit = %d, want 0", got)
	}
	if got := ReversalPairs([]int{13, 31}, 31); got != 1 {
		t.Errorf("partner at limit = %d, want 1", got)
	}
}

func TestReversalPairsHandlesDegenerateInput(t *testing.T) {
	if got := ReversalPairs(nil, 100); got != 0 {
		t.Errorf("empty = %d, want 0", got)
	}
}
