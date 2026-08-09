package pattern

// PalindromesWith counts the palindromic windows of length k that also satisfy
// keep.
//
// It exists so a population can be split without reimplementing the scan. The
// parity law allows an odd-length palindrome by either of two routes — its
// centre gap is not a multiple of 3, or no gap in it is — and asking which
// route carries the odd excess needs the same detector applied twice under
// different filters.
//
// keep receives the window itself and must not modify it. A nil filter, a k
// below 2, or a k beyond the input all yield 0.
func PalindromesWith(gaps []int, k int, keep func(window []int) bool) int {
	if keep == nil || k < 2 || k > len(gaps) {
		return 0
	}

	count := 0
	for start := 0; start+k <= len(gaps); start++ {
		window := gaps[start : start+k]
		if isMirrored(window) && keep(window) {
			count++
		}
	}

	return count
}
