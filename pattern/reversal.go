package pattern

// ReversalPairs counts the pairs v < r within seq where r is v with its
// decimal digits reversed, both at most limit.
//
// This detector exists as a monument. The four pairs below 100 — 13/31,
// 17/71, 37/73, 79/97 — looked like structure and died against the control:
// random odd sets of the same density produce 2.6 pairs on average, and four
// is comfortably inside that noise (p = 0.195). The detector is kept so the
// death stays reproducible; a findings record that can only re-run its
// survivors is a sales brochure.
//
// Palindromic values pair with themselves and are not counted.
func ReversalPairs(seq []int, limit int) int {
	in := map[int]bool{}
	for _, v := range seq {
		if v <= limit {
			in[v] = true
		}
	}

	count := 0
	for v := range in {
		r := reverseDigits(v)
		if r != v && r <= limit && in[r] && v < r {
			count++
		}
	}

	return count
}

func reverseDigits(n int) int {
	r := 0
	for n > 0 {
		r = r*10 + n%10
		n /= 10
	}
	return r
}
