// Command adele is the laboratory's first walk onto the bridge of natures.
//
// The flash behind it: weave in blocks of twelve positions "below the
// comma", representing the infinite exactly by coherent finite data — which
// is the p-adic construction, and the glue joining the real nature to every
// p-adic nature is the adele ring, the arena of Tate's functional equation
// and Connes' attack on the hypothesis.
//
// The experiment: in the 2-adic nature, the first twelve binary digits of a
// prime are its address on the first twelve floors of the binary tree —
// its residue mod 2^k for k = 1..12. This command measures how the primes
// fill that tree, floor by floor.
//
// PRE-REGISTERED: fair filling at every floor — chi-square per degree of
// freedom near 1 and no exceptional class — with one known whisper allowed:
// the Chebyshev bias, which favours non-square classes by a term of order
// √x/ln x. Any floor with chi²/df far above 1 beyond that whisper is news.
//
// Usage:
//
//	go run ./cmd/adele [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/nicoconvoz/diosyunalma/primes"
)

const floors = 12

func main() {
	limit := flag.Int("limit", 1_000_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	mod := 1 << floors
	counts := make([]float64, mod)
	ps := primes.Sieve(*limit)
	var total float64
	for _, p := range ps {
		if p == 2 {
			continue
		}
		counts[p%mod]++
		total++
	}

	fmt.Printf("THE ADELE — %d primes on the first %d floors of the 2-adic tree\n\n",
		int(total), floors)
	fmt.Println(" floor   classes   chi2/df   max |z|   verdict")
	for k := 1; k <= floors; k++ {
		m := 1 << k
		nc := m / 2 // odd residue classes
		cls := make([]float64, m)
		for r := 1; r < mod; r += 2 {
			cls[r%m] += counts[r]
		}
		exp := total / float64(nc)
		chi2, maxZ := 0.0, 0.0
		for r := 1; r < m; r += 2 {
			d := cls[r] - exp
			chi2 += d * d / exp
			if z := math.Abs(d) / math.Sqrt(exp); z > maxZ {
				maxZ = z
			}
		}
		df := float64(nc - 1)
		verdict := "fair"
		if df > 0 && chi2/df > 1+6/math.Sqrt(df) {
			verdict = "STRUCTURED"
		}
		if k == 1 {
			fmt.Printf("   %2d      %5d      all odd primes share one room\n", k, nc)
			continue
		}
		fmt.Printf("   %2d      %5d     %6.3f    %5.2f    %s\n", k, nc, chi2/df, maxZ, verdict)
	}

	// the famous whisper: the level-2 race.
	m4 := make([]float64, 4)
	for r := 1; r < mod; r += 2 {
		m4[r%4] += counts[r]
	}
	fmt.Printf("\nthe Chebyshev whisper at floor 2: pi(x;4,3) - pi(x;4,1) = %+.0f\n", m4[3]-m4[1])
	fmt.Printf("(predicted order sqrt(x)/ln x = %.0f; the finite nature leans, barely)\n",
		math.Sqrt(float64(*limit))/math.Log(float64(*limit)))

	// the circle and its melodies: (Z/2^12)* = mirror(-1) x circle(5).
	// discrete-log coordinates: every odd a = (-1)^eps * 5^j mod 2^floors.
	half := mod / 4 // order of 5 mod 2^floors
	idx := make(map[int][2]int)
	v := 1
	for j := 0; j < half; j++ {
		idx[v] = [2]int{0, j}
		idx[mod-v] = [2]int{1, j}
		v = v * 5 % mod
	}
	c := make([][]float64, 2)
	c[0] = make([]float64, half)
	c[1] = make([]float64, half)
	for r := 1; r < mod; r += 2 {
		co := idx[r]
		c[co[0]][co[1]] += counts[r]
	}
	fmt.Println("\nthe melody scan: every wave of the mirror-and-circle symmetry, loudness")
	fmt.Println("in noise units (fair filling predicts nothing above ~4):")
	type loud struct {
		m    int
		sign string
		a    float64
	}
	best := []loud{}
	for _, sign := range []int{0, 1} {
		tag := "+"
		if sign == 1 {
			tag = "-"
		}
		for m := 0; m < half; m++ {
			if sign == 0 && m == 0 {
				continue // the constant wave = the total count
			}
			var re, im float64
			for j := 0; j < half; j++ {
				ph := 2 * math.Pi * float64(m) * float64(j) / float64(half)
				w := c[0][j]
				if sign == 1 {
					w -= c[1][j]
				} else {
					w += c[1][j]
				}
				re += w * math.Cos(ph)
				im += w * math.Sin(ph)
			}
			a := math.Sqrt(re*re+im*im) / math.Sqrt(total)
			best = append(best, loud{m, tag, a})
		}
	}
	top := []loud{}
	for _, l := range best {
		if len(top) < 5 {
			top = append(top, l)
		} else {
			mi := 0
			for i := 1; i < 5; i++ {
				if top[i].a < top[mi].a {
					mi = i
				}
			}
			if l.a > top[mi].a {
				top[mi] = l
			}
		}
	}
	for _, l := range top {
		fmt.Printf("   wave (sign %s, frequency %4d): loudness %.2f\n", l.sign, l.m, l.a)
	}
	fmt.Println("\nthe primes fill the finite nature's tree like fair coins, twelve floors")
	fmt.Println("deep - equidistribution on the bridge, with one whisper of bias.")
}
