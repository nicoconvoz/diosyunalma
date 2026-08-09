package primes

// ProgressionPairs counts the values p for which p and p+d are both prime and
// p+d stays inside the set.
//
// Like ProgressionTriples it imposes no consecutiveness: other primes may sit
// between the two. That is what makes it useful. Dividing the count of
// consecutive pairs with gap d by this number gives the fraction of prime pairs
// at distance d that happen to have nothing between them.
//
// With the triple analogue, the two survival rates separate the arithmetic part
// of a measurement from the cost of demanding adjacency:
//
//	observed / independent  =  [ T(d)·n / P₂(d)² ]  ×  [ s₃(d) / s₂(d)² ]
//	                             arithmetic            consecutiveness
func ProgressionPairs(s Set, d int) int {
	if d < 1 {
		return 0
	}

	count := 0
	for p := 2; p+d < len(s); p++ {
		if s[p] && s[p+d] {
			count++
		}
	}

	return count
}
