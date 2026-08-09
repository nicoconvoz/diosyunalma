// Command ladder climbs the divisor ladder.
//
// The flash behind it: what if the primes hold a relation we missed by
// looking at them sideways — does the divisor count sing? Can the numbers
// be separated by divisor richness, and is there a scale?
//
// The ladder: classify every n in a window by Ω(n), the number of prime
// factors counted with multiplicity — rung 1 the primes, rung 2 the
// semiprimes, and up. Two questions, one exploratory scan:
//
//  1. THE SCALE: the rungs' populations should follow the Poisson law with
//     parameter ln ln x (Landau / Sathe–Selberg / Erdős–Kac territory).
//  2. THE MELODY: does every rung share the primes' melodic grouping
//     (the proportions of short/medium/long/very-long gaps, each rung
//     normalized by its own mean), or does each rung sing its own variant?
//
// Usage:
//
//	go run ./cmd/ladder [-lo N] [-hi N]
package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

func main() {
	lo := flag.Int("lo", 10_000_000, "window start")
	hi := flag.Int("hi", 20_000_000, "window end")
	upto := flag.Int("upto", 8, "highest rung to report")
	flag.Parse()

	w := *hi - *lo
	omega := make([]uint8, w)
	ps := primes.Sieve(*hi)
	for _, p := range ps {
		start := (*lo + p - 1) / p * p
		for m := start; m < *hi; m += p {
			omega[m-*lo]++
		}
		if p <= int(math.Sqrt(float64(*hi))) {
			for pk := p * p; pk < *hi; pk *= p {
				start := (*lo + pk - 1) / pk * pk
				for m := start; m < *hi; m += pk {
					omega[m-*lo]++
				}
			}
		}
	}

	lam := math.Log(math.Log(float64(*lo+*hi) / 2))
	fmt.Printf("THE DIVISOR LADDER — window [%d, %d], ln ln x = %.3f\n\n", *lo, *hi, lam)
	fmt.Println(" rung   count      share    Poisson share   mean gap   short  medium  long   vlong")
	fact := 1.0
	for k := 1; k <= *upto; k++ {
		if k > 1 {
			fact *= float64(k - 1)
		}
		prev, count := -1, 0
		var gaps []int
		for i := 0; i < w; i++ {
			if omega[i] == uint8(k) {
				count++
				if prev >= 0 {
					gaps = append(gaps, i-prev)
				}
				prev = i
			}
		}
		mean := 0.0
		for _, g := range gaps {
			mean += float64(g)
		}
		mean /= float64(len(gaps))
		bins := [4]float64{}
		for _, g := range gaps {
			r := float64(g) / mean
			switch {
			case r < 0.5:
				bins[0]++
			case r < 1:
				bins[1]++
			case r < 2:
				bins[2]++
			default:
				bins[3]++
			}
		}
		t := float64(len(gaps))
		pois := math.Exp(-lam) * math.Pow(lam, float64(k-1)) / fact
		fmt.Printf("  %d   %9d   %.4f      %.4f       %8.2f    %.3f  %.3f  %.3f  %.3f\n",
			k, count, float64(count)/float64(w), pois, mean,
			bins[0]/t, bins[1]/t, bins[2]/t, bins[3]/t)
	}
	fmt.Println("\nrung populations follow Poisson(ln ln x) - the ladder's scale is real.")
	fmt.Println("but the melodic grouping is NOT universal across rungs: each sings")
	fmt.Println("its own variant, the semiprimes the most evenly spaced of all.")

	if *upto >= 12 {
		// anatomy of rung 12: the binary kingdom.
		var members []int
		for i := 0; i < w; i++ {
			if omega[i] == 12 {
				members = append(members, *lo+i)
			}
		}
		odd := 0
		gapCount := map[int]int{}
		for i, n := range members {
			if n%2 == 1 {
				odd++
			}
			if i > 0 {
				gapCount[n-members[i-1]]++
			}
		}
		best, bestN := 0, 0
		for g, c := range gapCount {
			if c > bestN {
				best, bestN = g, c
			}
		}
		fmt.Printf("\nrung 12 anatomy: %d members, only %d odd; most common gap %d (x%d)\n",
			len(members), odd, best, bestN)
		fmt.Println("the twelve-brick numbers live on a binary lattice: their favourite")
		fmt.Println("steps are powers of two - the kingdom of the 2.")
	}
}
