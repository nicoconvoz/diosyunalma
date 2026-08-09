package pattern

// GilbreathRows repeatedly takes absolute differences of seq and reports how
// many of the first maxRows rows begin with 1, and the row where that first
// fails (-1 if it never does within the run).
//
// Gilbreath's conjecture (1958, open) says every row of this triangle built
// from the primes begins with 1. Shuffled-gap decoys survive only rarely, so
// the ordering of the gaps carries information the distribution alone does
// not — but the shuffle also destroys the small-gaps-first tendency, so the
// control is imperfect and the measurement stays an open finding, not a
// confirmed one.
func GilbreathRows(seq []int, maxRows int) (rowsOK int, brokeAt int) {
	current := seq
	for row := 1; row <= maxRows && len(current) > 1; row++ {
		next := make([]int, len(current)-1)
		for i := range next {
			d := current[i+1] - current[i]
			if d < 0 {
				d = -d
			}
			next[i] = d
		}
		if next[0] != 1 {
			return rowsOK, row
		}
		rowsOK = row
		current = next
	}
	return rowsOK, -1
}
