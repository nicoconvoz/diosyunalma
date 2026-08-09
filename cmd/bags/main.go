// Command bags2 is the repaired test. The first design compared residues
// within one wheel's lanes — and was broken by the interlock: the OTHER
// wheel's compatibility bans different small gap VALUES for different
// residues (residue 1 can never jump 4; residue 7 can never jump 2), which
// shifts every marginal split mechanically. First-order determinism wearing a
// second-order costume, caught because the z came out absurd.
//
// The honest null: ONE shared bag of gap values, filtered by each residue's
// legal menu and renormalised. w(d) is estimated from the pooled counts
// divided by the number of residues whose menu allows d. Whatever deviation
// survives THIS null is the genuine Lemke Oliver–Soundararajan content: the
// bag itself differing by standing residue beyond legality.
package main

import (
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

var residues = []int{1, 2, 4, 7, 8, 11, 13, 14}

const dMax = 150

func compatible(r, d int) bool {
	s := r + d
	return s%3 != 0 && s%5 != 0
}

func main() {
	walk := primes.From(primes.Sieve(100_000_000), 7)

	idx := map[int]int{}
	for i, r := range residues {
		idx[r] = i
	}

	var counts [8][dMax + 1]float64
	var used [8]float64
	dropped := 0
	for i := 0; i+1 < len(walk); i++ {
		d := walk[i+1] - walk[i]
		if d > dMax {
			dropped++
			continue
		}
		r := idx[walk[i]%15]
		counts[r][d]++
		used[r]++
	}
	fmt.Printf("transitions used: %.0f (dropped %d with gap > %d)\n",
		sum(used[:]), dropped, dMax)

	// The shared bag: pooled weight per gap value, divided by how many
	// residues may legally draw it.
	var w [dMax + 1]float64
	for d := 2; d <= dMax; d += 2 {
		var pool float64
		menus := 0
		for _, r := range residues {
			pool += counts[idx[r]][d]
			if compatible(r, d) {
				menus++
			}
		}
		if menus > 0 {
			w[d] = pool / float64(menus)
		}
	}

	fmt.Println("\nchi-square per residue against the shared-bag-with-menu null:")
	fmt.Printf("%-10s %-12s %-8s %s\n", "residue", "chi2", "bins", "biggest single deviation")
	totalChi, totalDof := 0.0, 0
	for _, r := range residues {
		ri := idx[r]
		var z float64
		for d := 2; d <= dMax; d += 2 {
			if compatible(r, d) {
				z += w[d]
			}
		}
		chi, bins := 0.0, 0
		worstD, worst := 0, 0.0
		for d := 2; d <= dMax; d += 2 {
			if !compatible(r, d) {
				continue
			}
			e := w[d] / z * used[ri]
			if e < 5 {
				continue
			}
			dev := (counts[ri][d] - e) * (counts[ri][d] - e) / e
			chi += dev
			bins++
			if dev > worst {
				worst, worstD = dev, d
			}
		}
		totalChi += chi
		totalDof += bins - 1
		obs := counts[ri][worstD]
		e := w[worstD] / z * used[ri]
		fmt.Printf("%-10d %-12.1f %-8d gap %d: obs %.0f vs exp %.0f (%+.1f%%)\n",
			r, chi, bins, worstD, obs, e, 100*(obs-e)/e)
	}
	totalDof -= 0 // w is fitted from the pooled data; dof is approximate
	zScore := (totalChi - float64(totalDof)) / math.Sqrt(2*float64(totalDof))
	fmt.Printf("\nTOTAL chi2 = %.1f on ~%d dof    z = %+.1f\n", totalChi, totalDof, zScore)

	if zScore > 5 {
		fmt.Println("\nVERDICT: the bag depends on the standing residue BEYOND legality —")
		fmt.Println("genuine second-order structure. This is the real LO-S content.")
	} else {
		fmt.Println("\nVERDICT: one shared bag plus legal menus explains everything —")
		fmt.Println("no second-order structure at 5 sigma.")
	}
}

func sum(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s
}
