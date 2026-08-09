// Command wheels couples the two measured wheels: the mod-3 lane walk of
// Finding 13 and the golden character of Finding 37.
//
// THE QUESTION. Each wheel repels on its own: P(stay-3) ≈ 0.43, P(stay-g) ≈
// 0.45. If the wheels turn independently, the joint stay probability is their
// product. Any deviation is COUPLING — the two repulsions talking to each
// other through the shared gap.
//
// Also, the 256 bombita: the z of Finding 37 is not a constant of nature. It
// grows as the square root of the sample. Measured here at two limits to bury
// that with evidence.
package main

import (
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

func golden(p int) bool { r := p % 5; return r == 1 || r == 4 }

func measure(walk []int) (s3, s5, joint float64, n int) {
	var c3, c5, cj int
	for i := 0; i+1 < len(walk); i++ {
		a3 := walk[i]%3 == walk[i+1]%3
		a5 := golden(walk[i]) == golden(walk[i+1])
		if a3 {
			c3++
		}
		if a5 {
			c5++
		}
		if a3 && a5 {
			cj++
		}
	}
	n = len(walk) - 1
	return float64(c3) / float64(n), float64(c5) / float64(n), float64(cj) / float64(n), n
}

func main() {
	walk := primes.From(primes.Sieve(100_000_000), 7)

	// The bombita first: z grows with sqrt(n), it is not a constant.
	fmt.Println("THE 256 BOMBITA, DEFUSED — z versus sample size")
	for _, cut := range []int{10_000_000, 100_000_000} {
		end := 0
		for end < len(walk) && walk[end] <= cut {
			end++
		}
		_, s5, _, n := measure(walk[:end])
		z := (s5 - 0.5) / math.Sqrt(0.25/float64(n))
		fmt.Printf("  N=%.0e: P(stay-golden) = %.5f   z = %+.1f\n", float64(cut), s5, z)
	}
	fmt.Println("  same physics, different z: the number measures data, not nature.")

	// The double wheel.
	s3, s5, joint, n := measure(walk)
	product := s3 * s5

	fmt.Println("\nTHE DOUBLE WHEEL — do the two repulsions turn independently?")
	fmt.Printf("  P(stay mod-3)             = %.5f\n", s3)
	fmt.Printf("  P(stay golden)            = %.5f\n", s5)
	fmt.Printf("  product (if independent)  = %.5f\n", product)
	fmt.Printf("  P(both stay), measured    = %.5f\n", joint)

	// Coupling: the phi coefficient of the two stay indicators.
	num := joint - product
	den := math.Sqrt(s3 * (1 - s3) * s5 * (1 - s5))
	corr := num / den
	z := corr * math.Sqrt(float64(n))
	fmt.Printf("\n  coupling (correlation of the two stays) = %+.5f    z = %+.1f\n", corr, z)

	// Where does the coupling live? The 2x2 table of (stay3, stay5) against
	// the independent prediction.
	fmt.Println("\n  the full 2x2, measured / predicted-if-independent:")
	p11, p10 := joint, s3-joint
	p01, p00 := s5-joint, 1-s3-s5+joint
	q11, q10 := s3*s5, s3*(1-s5)
	q01, q00 := (1-s3)*s5, (1-s3)*(1-s5)
	fmt.Printf("    both stay      %.5f / %.5f  = %.4f\n", p11, q11, p11/q11)
	fmt.Printf("    only mod-3     %.5f / %.5f  = %.4f\n", p10, q10, p10/q10)
	fmt.Printf("    only golden    %.5f / %.5f  = %.4f\n", p01, q01, p01/q01)
	fmt.Printf("    both switch    %.5f / %.5f  = %.4f\n", p00, q00, p00/q00)

	if math.Abs(z) > 5 {
		fmt.Println("\n  VERDICT: the wheels are COUPLED - the repulsions talk to each other.")
	} else {
		fmt.Println("\n  VERDICT: independent at 5 sigma - the wheels turn alone.")
	}
}
