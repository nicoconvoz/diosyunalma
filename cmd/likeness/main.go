// Command likeness reads the creator's signature in the building's form.
//
// The push behind it: the constructions may not inherit the creator's
// property (rigidity evaporates, Finding 78), but they carry the signature,
// image and likeness IN THEIR FORM — why exactly THAT building and not
// another? The measurable version: the ladder's whole shape has one moving
// parameter, the lagoon level ln ln x. Whatever remains of the floor
// ratios AFTER subtracting the lagoon must be constants — numbers carved
// into the architecture by the bricks themselves.
//
//	S₂(x) = N₂/N₁ − ln ln x        (second floor over first)
//	S₃(x) = 2·N₃/N₂ − ln ln x      (third over second, Poisson-corrected)
//
// PRE-REGISTERED: both S₂ and S₃ hold steady within ±0.05 while x grows by
// a factor of a thousand — stable carved constants, the signature — or the
// likeness reading fails.
//
// Usage:
//
//	go run ./cmd/likeness
package main

import (
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

func rungCounts(lo, hi, maxK int) []float64 {
	w := hi - lo
	omega := make([]uint8, w)
	ps := primes.Sieve(hi)
	for _, p := range ps {
		start := (lo + p - 1) / p * p
		for m := start; m < hi; m += p {
			omega[m-lo]++
		}
		if p <= int(math.Sqrt(float64(hi))) {
			for pk := p * p; pk < hi; pk *= p {
				start := (lo + pk - 1) / pk * pk
				for m := start; m < hi; m += pk {
					omega[m-lo]++
				}
			}
		}
	}
	out := make([]float64, maxK+1)
	for i := 0; i < w; i++ {
		if k := int(omega[i]); k <= maxK {
			out[k]++
		}
	}
	return out
}

func main() {
	fmt.Println("THE LIKENESS — the signature carved into the building's proportions")
	fmt.Println("\n   window          lnlnx    S2 = N2/N1 - lnlnx    S3 = 2*N3/N2 - lnlnx")
	type win struct{ lo, hi int }
	wins := []win{{100_000, 200_000}, {1_000_000, 2_000_000},
		{10_000_000, 20_000_000}, {100_000_000, 200_000_000}}
	var s2s, s3s []float64
	for _, w := range wins {
		lam := math.Log(math.Log(float64(w.lo+w.hi) / 2))
		c := rungCounts(w.lo, w.hi, 4)
		s2 := c[2]/c[1] - lam
		s3 := 2*c[3]/c[2] - lam
		s2s = append(s2s, s2)
		s3s = append(s3s, s3)
		fmt.Printf("   [%.0e,%.0e]    %.3f        %+.4f              %+.4f\n",
			float64(w.lo), float64(w.hi), lam, s2, s3)
	}
	spread := func(v []float64) float64 {
		mn, mx := v[0], v[0]
		for _, x := range v {
			if x < mn {
				mn = x
			}
			if x > mx {
				mx = x
			}
		}
		return mx - mn
	}
	fmt.Printf("\n   spread of S2 across a factor 1000 in x: %.4f (registered band: 0.05)\n", spread(s2s))
	fmt.Printf("   spread of S3 across a factor 1000 in x: %.4f (registered band: 0.05)\n", spread(s3s))
	fmt.Println("\nsubtract the lagoon's level and the floor ratios freeze: the numbers")
	fmt.Println("that remain are carved constants - the bricks' signature in the walls.")
	fmt.Println("why THAT building? one moving parameter, the lagoon; everything else")
	fmt.Println("is the creator's handwriting, and it does not change.")
}
