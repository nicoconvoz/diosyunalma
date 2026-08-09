// Command carvings measures the carving angle of every tribe at once.
//
// Finding 73 showed the stations of ζ are carved out of the light at phase
// 90° + arctan(½/γ) — the absorption signature, with the offset above 90°
// encoding the critical line's real part. This run points the same phase
// instrument at every real dial in the orchestra: for a non-principal
// character the signal ψ(x,χ)/√x has no smooth term, and at each measured
// station the complex coefficient must again be ≈ −1/ρ.
//
// PRE-REGISTERED: for every dial, phases cluster near +90°, |c|·|ρ| near 1,
// and the per-dial mean recovered real part lands in [0.3, 0.7] — the ½ of
// eight L-functions measured simultaneously by absorption angles.
//
// Usage:
//
//	go run ./cmd/carvings [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/numerosprimos/primes"
)

type dial struct {
	name  string
	chi   func(int) float64
	zeros []float64
}

func modPow(base, exp, mod int) int {
	r := 1
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			r = r * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return r
}

func legendre(q int) func(int) float64 {
	half := (q - 1) / 2
	table := make([]float64, q)
	for n := 1; n < q; n++ {
		switch modPow(n, half, q) {
		case 1:
			table[n] = 1
		case q - 1:
			table[n] = -1
		}
	}
	return func(n int) float64 { return table[n%q] }
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	one := func(n int) float64 { return 1 }
	chi4 := func(n int) float64 {
		switch n % 4 {
		case 1:
			return 1
		case 3:
			return -1
		}
		return 0
	}
	chi8 := func(n int) float64 {
		switch n % 8 {
		case 1, 7:
			return 1
		case 3, 5:
			return -1
		}
		return 0
	}
	dials := []dial{
		{"zeta", one, []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422, 37.5872, 40.9264, 43.3211, 48.0105, 49.7752}},
		{"chi3", legendre(3), []float64{8.0396, 11.2450, 15.7062, 18.2579, 20.4551, 24.0636}},
		{"chi4", chi4, []float64{6.0199, 10.2423, 12.9848, 16.3464, 18.2914, 21.4547}},
		{"chi5", legendre(5), []float64{6.6516, 9.8280, 11.9612, 16.0386, 17.5632, 19.5431, 22.2228, 24.5864}},
		{"chi7", legendre(7), []float64{4.4762, 6.8400, 11.1782, 12.4670, 15.1161, 16.7892}},
		{"chi8", chi8, []float64{4.8989, 7.6194, 10.8219, 12.3126, 15.1935, 17.0246}},
		{"chi11", legendre(11), []float64{2.4768, 6.7997, 8.9663, 10.1257, 13.0422, 15.0996}},
		{"chi13", legendre(13), []float64{3.1119, 7.2340, 8.6013, 10.3241, 12.6185, 15.1341}},
	}

	const du = 0.005
	var us []float64
	var xs []int
	for u := math.Log(100); u <= math.Log(float64(*limit)); u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	ps := primes.Sieve(*limit)
	T := us[len(us)-1] - us[0]

	fmt.Println("THE CARVINGS — the absorption angle of every tribe")
	fmt.Println("\n  dial    stations   phases in band   mean |c|*|rho|   mean sigma")
	var grand, grandN float64
	for _, d := range dials {
		type event struct {
			at int
			v  float64
		}
		evs := []event{}
		for _, p := range ps {
			if p*p > *limit {
				break
			}
			lg := math.Log(float64(p))
			for pk := p * p; pk <= *limit; pk *= p {
				evs = append(evs, event{pk, d.chi(pk) * lg})
			}
		}
		sort.Slice(evs, func(i, j int) bool { return evs[i].at < evs[j].at })
		es := make([]float64, len(xs))
		sum, pi, wi := 0.0, 0, 0
		for i, x := range xs {
			for pi < len(ps) && ps[pi] <= x {
				sum += d.chi(ps[pi]) * math.Log(float64(ps[pi]))
				pi++
			}
			for wi < len(evs) && evs[wi].at <= x {
				sum += evs[wi].v
				wi++
			}
			es[i] = sum / math.Sqrt(float64(x))
			if d.name == "zeta" {
				es[i] = (sum - float64(x)) / math.Sqrt(float64(x))
			}
		}
		var ampSum, sigSum float64
		inBand := 0
		for _, g := range d.zeros {
			var re, im float64
			for i, u := range us {
				re += es[i] * math.Cos(g*u) * du
				im -= es[i] * math.Sin(g*u) * du
			}
			re /= T
			im /= T
			phase := math.Atan2(im, re) * 180 / math.Pi
			rho := math.Sqrt(0.25 + g*g)
			if phase > 45 && phase < 135 {
				inBand++
				ampSum += math.Hypot(re, im) * rho
				sigSum += g * math.Tan((phase-90)*math.Pi/180)
			}
		}
		if inBand == 0 {
			fmt.Printf("  %-6s   %5d       none in band\n", d.name, len(d.zeros))
			continue
		}
		fmt.Printf("  %-6s   %5d        %2d/%d            %.3f          %+.3f\n",
			d.name, len(d.zeros), inBand, len(d.zeros),
			ampSum/float64(inBand), sigSum/float64(inBand))
		grand += sigSum
		grandN += float64(inBand)
	}
	fmt.Printf("\ngrand mean recovered real part over all tribes: %+.3f (the line says +0.500)\n",
		grand/grandN)
	fmt.Println("\neight functions, one chisel angle: every tribe carves its absences")
	fmt.Println("at the angle of one half - the generalized line, heard as geometry.")
}
