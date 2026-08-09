package pattern

// LagEquality counts the positions where two gaps separated by lag are equal.
//
// It answers a sharper question than Palindromes: not "is this window
// symmetric" but "at what distance do distances repeat". Running it across a
// range of lags isolates where a correlation lives, which a symmetry count
// alone cannot show — a palindrome of length k folds several lags together, so
// an excess there could come from any of them.
//
// A lag below 1, or at least the length of gaps, yields 0.
func LagEquality(gaps []int, lag int) int {
	if lag < 1 || lag >= len(gaps) {
		return 0
	}

	count := 0
	for i := 0; i+lag < len(gaps); i++ {
		if gaps[i] == gaps[i+lag] {
			count++
		}
	}

	return count
}
