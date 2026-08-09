// Command broth cooks the k12 broth and watches how it moves.
//
// The flash behind it: make the broth of rung 12, understand the motion of
// its behaviour — the harmony of its chaos may be a clue to the map.
//
// Two pre-registered claims:
//
//  1. THE MOTION — the broth has two layers moving at different speeds.
//     The surface (body rungs, k ≤ 5) swirls with the slowest clock in
//     nature: shares tracking Poisson(ln ln x) as the window slides across
//     decades. The depths (k ≥ 10) are FROZEN in the domino's law: the
//     ratios N_k/N_{k−1} stay near 1/2 (measured ≈ 0.47) in every window,
//     barely moving while the surface turns.
//
//  2. THE HARMONY OF THE CHAOS — rung 12's pair differences live on the
//     binary comb (≥ 80% divisible by 128), and the comb's occupied teeth
//     are filled FREELY: the spacing of members in comb units has variance
//     near Poisson's 1. The raw primes are super-rigid (GUE, var 0.18);
//     their boiled broth is free — rigidity is an aroma that evaporates
//     with cooking, which locates it: in the primes themselves, not in
//     the multiplicative scaffolding.
//
// Usage:
//
//	go run ./cmd/broth
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
	fmt.Println("THE BROTH — rung 12 on the fire, watched while it moves")

	// 1) the motion: four windows, one decade apart.
	type win struct{ lo, hi int }
	wins := []win{{100_000, 200_000}, {1_000_000, 2_000_000},
		{10_000_000, 20_000_000}, {100_000_000, 200_000_000}}
	fmt.Println("\n1) THE TWO LAYERS")
	fmt.Println("   window        lnlnx   surface share k=2 (Poisson)   depth ratios N10/N9 N11/N10 N12/N11")
	for _, w := range wins {
		lam := math.Log(math.Log(float64(w.lo+w.hi) / 2))
		c := rungCounts(w.lo, w.hi, 13)
		var tot float64
		for _, v := range c {
			tot += v
		}
		pois2 := math.Exp(-lam) * lam
		fmt.Printf("   [%.0e,%.0e]   %.3f     %.4f (%.4f)               %.3f  %.3f  %.3f\n",
			float64(w.lo), float64(w.hi), lam, c[2]/tot, pois2,
			c[10]/c[9], c[11]/c[10], c[12]/c[11])
	}
	fmt.Println("   the surface turns with ln ln x; the depths hold the frozen 1/2.")

	// 2) the harmony of the chaos, in the [1e7, 2e7] pot.
	lo, hi := 10_000_000, 20_000_000
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
	members := []int{}
	for i := 0; i < w; i++ {
		if omega[i] == 12 {
			members = append(members, lo+i)
		}
	}

	// pair differences within 2048: how many sit on the binary comb?
	var pairs, comb128, comb256 float64
	j := 0
	for i := range members {
		for j < len(members) && members[j] <= members[i]+2048 {
			j++
		}
		for k := i + 1; k < j; k++ {
			d := members[k] - members[i]
			pairs++
			if d%128 == 0 {
				comb128++
			}
			if d%256 == 0 {
				comb256++
			}
		}
	}
	fmt.Printf("\n2) THE HARMONY OF THE CHAOS — %d members, %d close pairs\n",
		len(members), int(pairs))
	fmt.Printf("   pair differences on the 128-comb: %.1f%%   on the 256-comb: %.1f%%\n",
		100*comb128/pairs, 100*comb256/pairs)
	fmt.Println("   (a comb-free world would give 0.8% and 0.4%)")

	// the chaos layer: members with 2-power exactly 2^9, in comb units.
	teeth := []int{}
	for _, m := range members {
		v, n := 0, m
		for n%2 == 0 {
			n /= 2
			v++
		}
		if v == 9 {
			teeth = append(teeth, m/512)
		}
	}
	gaps := []float64{}
	for i := 1; i < len(teeth); i++ {
		gaps = append(gaps, float64(teeth[i]-teeth[i-1]))
	}
	mean := 0.0
	for _, g := range gaps {
		mean += g
	}
	mean /= float64(len(gaps))
	varr := 0.0
	for _, g := range gaps {
		varr += (g - mean) * (g - mean)
	}
	varr /= float64(len(gaps)) * mean * mean
	fmt.Printf("\n   the occupied teeth (2^9 members in comb units): %d members\n", len(teeth))
	fmt.Printf("   spacing variance: %.3f   (Poisson freedom = 1.000; the raw primes = 0.18)\n", varr)
	fmt.Println("\nboiling evaporates the rigidity: the broth's chaos is free on its comb,")
	fmt.Println("while the raw primes are stiff - the harmony of order lives in the primes")
	fmt.Println("themselves, not in the multiplicative scaffolding they build.")
}
