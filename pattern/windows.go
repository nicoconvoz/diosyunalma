package pattern

// Windows returns a copy of every palindromic window of length k that
// satisfies keep, in order of appearance.
//
// PalindromesWith answers how many; this answers which. The distinction became
// necessary when the geometric law broke at high k (Finding 23): explaining the
// break requires knowing whether the surviving windows are scattered
// coincidences or a few rigid constellations repeating — and that is a question
// about the windows themselves, not their count.
//
// Each returned window is an independent copy, safe to mutate or tally.
func Windows(gaps []int, k int, keep func(window []int) bool) [][]int {
	if keep == nil || k < 2 || k > len(gaps) {
		return nil
	}

	out := [][]int{}
	for start := 0; start+k <= len(gaps); start++ {
		window := gaps[start : start+k]
		if isMirrored(window) && keep(window) {
			out = append(out, append([]int(nil), window...))
		}
	}

	return out
}
