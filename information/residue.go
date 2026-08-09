package information

import "math"

// ResidueEntropy measures how evenly values spread across residue classes
// modulo m, as Shannon entropy normalised by log2(m): 1 means perfectly even,
// lower means concentration into fewer classes.
//
// Values up to the modulus are skipped — for a prime sequence they are the
// structural exceptions that generate the modular wheel, not samples of it.
//
// This is the sweep that discovered the wheel without being told about it:
// asked which modulus makes the primes look least random, it returns the
// primorials, while decoys sit flat at 1 for every modulus.
func ResidueEntropy(values []int, modulus int) float64 {
	if modulus < 2 {
		return 0
	}

	counts := make([]int, modulus)
	seen := false
	for _, v := range values {
		if v <= modulus {
			continue
		}
		counts[v%modulus]++
		seen = true
	}
	if !seen {
		return 0
	}

	return Shannon(counts) / math.Log2(float64(modulus))
}
