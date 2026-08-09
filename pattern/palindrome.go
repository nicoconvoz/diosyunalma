// Package pattern holds detectors: functions that count how often a candidate
// structure appears in a sequence of distances.
//
// A detector on its own proves nothing. Its output only becomes evidence once
// the control package has scored it against decoy sequences.
package pattern

// Palindromes counts the windows of k consecutive gaps that read identically
// forwards and backwards.
//
// Windows overlap: every starting offset is counted separately. A k below 2, or
// larger than the input, yields 0.
//
// The measurement is direction-blind by construction, which is the point: it
// asks whether the distances between primes carry a symmetry that survives
// being read from either end.
func Palindromes(gaps []int, k int) int {
	if k < 2 || k > len(gaps) {
		return 0
	}

	count := 0
	for start := 0; start+k <= len(gaps); start++ {
		if isMirrored(gaps[start : start+k]) {
			count++
		}
	}

	return count
}

// isMirrored reports whether window reads the same in both directions. Only
// half the comparisons are needed; a lone middle element always matches itself.
func isMirrored(window []int) bool {
	for i := 0; i < len(window)/2; i++ {
		if window[i] != window[len(window)-1-i] {
			return false
		}
	}
	return true
}
