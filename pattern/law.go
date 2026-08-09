package pattern

// LawViolations counts the palindromic windows of length k that break the
// parity law.
//
// The law: a palindromic gap window holds an odd number of gaps that are not
// multiples of 3, or exactly zero of them. An even non-zero count cannot occur.
//
// It follows from the residue walk. Primes past 3 live only in residues 1 and 2
// modulo 3, so a gap either keeps the residue or flips it, and with two states
// available the flips must alternate 1, 2, 1, 2. Requiring that alternating
// chain to read the same in both directions forces the number of flips to be
// odd.
//
// Over primes above 3 this returns 0 for every k tested so far. The window
// (2, 2) from 3, 5, 7 is the sole exception in the natural numbers, and it
// vanishes once 3 is excluded — the prime that generates the law is the one
// place the law cannot reach.
func LawViolations(gaps []int, k int) int {
	if k < 2 || k > len(gaps) {
		return 0
	}

	count := 0
	for start := 0; start+k <= len(gaps); start++ {
		window := gaps[start : start+k]
		if !isMirrored(window) {
			continue
		}
		if flips := nonMultiplesOfThree(window); flips != 0 && flips%2 == 0 {
			count++
		}
	}

	return count
}

func nonMultiplesOfThree(window []int) int {
	n := 0
	for _, g := range window {
		if g%3 != 0 {
			n++
		}
	}
	return n
}

// SingularBoost returns the Hardy–Littlewood correction for a repeated gap of
// size d, relative to a model that treats the two gaps as independent.
//
// For three primes in arithmetic progression p, p+d, p+2d the singular series
// asks how many residue classes modulo each prime q the set {0, d, 2d}
// occupies. When q divides d it occupies one class and the constraint vanishes;
// otherwise it occupies three and one of them is 0, thinning the density. That
// gives a boost of (q-1)/(q-3) per prime q > 3 dividing d.
//
// The measured quantity is a ratio against the squared marginal gap rate, and
// that marginal already carries the 2-tuple series with boost (q-1)/(q-2).
// Dividing it out leaves
//
//	(q-2)² / ((q-1)(q-3))
//
// per distinct prime q > 3 dividing d — about 1.125 for 5, 1.042 for 7, and
// rapidly less for larger primes.
//
// The factors 2 and 3 carry their own terms in the series and are excluded
// here, so a gap that is a pure product of 2s and 3s returns exactly 1.
func SingularBoost(d int) float64 {
	if d < 0 {
		d = -d
	}
	for d%2 == 0 && d > 0 {
		d /= 2
	}
	for d%3 == 0 && d > 0 {
		d /= 3
	}

	boost := 1.0
	apply := func(q int) {
		f := float64(q)
		boost *= (f - 2) * (f - 2) / ((f - 1) * (f - 3))
	}

	for q := 5; q*q <= d; q += 2 {
		if d%q == 0 {
			apply(q)
			for d%q == 0 {
				d /= q
			}
		}
	}
	if d > 1 {
		apply(d)
	}

	return boost
}
