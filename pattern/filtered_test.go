package pattern

import "testing"

func TestPalindromesWithKeepsOnlyTheMatchingWindows(t *testing.T) {
	// [6 4 6 2 6] — the k=3 palindromes are [6 4 6] at 0 and [6 2 6] at 2.
	gaps := []int{6, 4, 6, 2, 6}

	all := func([]int) bool { return true }
	none := func([]int) bool { return false }
	centreIsTwo := func(w []int) bool { return w[len(w)/2] == 2 }

	tests := []struct {
		name string
		keep func([]int) bool
		want int
	}{
		{name: "keeping everything matches the plain count", keep: all, want: 2},
		{name: "keeping nothing yields zero", keep: none, want: 0},
		{name: "filtering on the centre gap", keep: centreIsTwo, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PalindromesWith(gaps, 3, tt.keep); got != tt.want {
				t.Errorf("PalindromesWith(%v, 3, %s) = %d, want %d", gaps, tt.name, got, tt.want)
			}
		})
	}
}

// TestPalindromesWithMatchesPalindromesWhenUnfiltered pins the two counters
// together: an unconditional filter must reproduce the plain detector exactly,
// or one of them is scanning differently.
func TestPalindromesWithMatchesPalindromesWhenUnfiltered(t *testing.T) {
	inputs := [][]int{
		{1, 2, 2, 4, 2, 4, 2, 4, 6},
		{6, 6, 6, 6, 6},
		{2, 4, 6, 4, 2, 8, 2, 4, 6, 10, 2, 2},
		{},
	}

	keep := func([]int) bool { return true }
	for _, gaps := range inputs {
		for k := 2; k <= 5; k++ {
			with := PalindromesWith(gaps, k, keep)
			plain := Palindromes(gaps, k)
			if with != plain {
				t.Errorf("%v k=%d: PalindromesWith = %d but Palindromes = %d",
					gaps, k, with, plain)
			}
		}
	}
}

// TestPalindromesWithSplitsIntoTheLawsTwoBranches exercises the split the odd
// excess is investigated through. The law allows an odd palindrome either
// because its centre gap is not a multiple of 3, or because no gap in it is —
// and the two branches must together account for every one of them.
func TestPalindromesWithSplitsIntoTheLawsTwoBranches(t *testing.T) {
	gaps := []int{6, 4, 6, 6, 6, 6, 2, 6, 12, 6, 12}

	centreFree := func(w []int) bool { return w[len(w)/2]%3 != 0 }
	allDivisible := func(w []int) bool {
		for _, g := range w {
			if g%3 != 0 {
				return false
			}
		}
		return true
	}

	for _, k := range []int{3, 5} {
		a := PalindromesWith(gaps, k, centreFree)
		b := PalindromesWith(gaps, k, allDivisible)
		if total := Palindromes(gaps, k); a+b != total {
			t.Errorf("k=%d: branches %d + %d = %d, but there are %d palindromes",
				k, a, b, a+b, total)
		}
	}
}

func TestPalindromesWithHandlesDegenerateInput(t *testing.T) {
	keep := func([]int) bool { return true }

	if got := PalindromesWith([]int{6, 6}, 1, keep); got != 0 {
		t.Errorf("k below two = %d, want 0", got)
	}
	if got := PalindromesWith([]int{6, 6}, 9, keep); got != 0 {
		t.Errorf("k beyond the input = %d, want 0", got)
	}
	if got := PalindromesWith([]int{6, 6}, 2, nil); got != 0 {
		t.Errorf("nil filter = %d, want 0", got)
	}
}
