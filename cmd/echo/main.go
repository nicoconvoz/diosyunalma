// Command echo is the harmonic scalpel for the chest's teeth.
//
// Finding 54 left a sharpened question: what generates the repetition
// anomaly's signature {6, 12, 36, 48} with 12 the lone positive? The musical
// hypothesis, registered here: the anomaly is an ECHO — the melody of
// blocking primes inside the first corridor repeating, note for note, at the
// same offsets one bar later inside the second.
//
// For every triple p, p+d, p+2d (all prime), and every even offset a inside
// a corridor, compare the joint rate of "prime at p+a AND prime at p+d+a"
// (the same note in both bars) against the product of the marginal rates.
//
// PRE-REGISTERED: the echo's sign must match the tooth's sign — positive
// for d = 12, negative for d = 36 and 48, near zero for the flat gaps
// 18, 24, 42 — or the musical reading is dead.
//
// Usage:
//
//	go run ./cmd/echo [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/nicoconvoz/diosyunalma/primes"
)

func main() {
	limit := flag.Int("limit", 1_000_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	ps := primes.Sieve(*limit)
	bits := make([]byte, *limit/8+1)
	for _, p := range ps {
		bits[p/8] |= 1 << (p % 8)
	}
	isPrime := func(n int) bool { return bits[n/8]&(1<<(n%8)) != 0 }

	fmt.Printf("THE ECHO — does the corridor's riff repeat one bar later? (x <= %d)\n", *limit)
	fmt.Println("\n   d   triples   echo z (Stouffer)   loudest same-note offsets      tooth (F54)")

	teeth := map[int]string{6: "-0.7%", 12: "+2.2%", 18: "flat", 24: "flat",
		30: "shrinking", 36: "-12.5%", 42: "flat", 48: "-9.1%"}

	for _, d := range []int{6, 12, 18, 24, 30, 36, 42, 48} {
		nA := d/2 - 1
		if nA < 1 {
			continue
		}
		nX := make([]float64, d)
		nY := make([]float64, d)
		nXY := make([]float64, d)
		var n float64
		for _, p := range ps {
			if p+2*d > *limit {
				break
			}
			if p < 5 || !isPrime(p+d) || !isPrime(p+2*d) {
				continue
			}
			n++
			for a := 2; a < d; a += 2 {
				x := isPrime(p + a)
				y := isPrime(p + d + a)
				if x {
					nX[a]++
				}
				if y {
					nY[a]++
				}
				if x && y {
					nXY[a]++
				}
			}
		}
		zSum, zN := 0.0, 0.0
		type loud struct {
			a int
			z float64
		}
		best := []loud{}
		for a := 2; a < d; a += 2 {
			exp := nX[a] * nY[a] / n
			if exp < 25 {
				continue
			}
			z := (nXY[a] - exp) / math.Sqrt(exp)
			zSum += z
			zN++
			best = append(best, loud{a, z})
		}
		stouffer := zSum / math.Sqrt(zN)
		// pick the two loudest offsets by |z|.
		b1, b2 := loud{}, loud{}
		for _, l := range best {
			if math.Abs(l.z) > math.Abs(b1.z) {
				b2 = b1
				b1 = l
			} else if math.Abs(l.z) > math.Abs(b2.z) {
				b2 = l
			}
		}
		fmt.Printf("  %2d  %8.0f      %+6.1f            a=%2d (z%+.1f)  a=%2d (z%+.1f)    %s\n",
			d, n, stouffer, b1.a, b1.z, b2.a, b2.z, teeth[d])
	}
	fmt.Println("\nthe echo is the same note struck in both bars: prime at p+a and at p+d+a.")
	fmt.Println("if its sign follows the teeth, the anomaly is melody remembering itself.")
}
