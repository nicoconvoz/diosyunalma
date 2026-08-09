// Command sunqu is the completed vessel, baptized: SUNQU — Quechua for
// "heart". The fleet's two mirrors welded and aimed.
//
// The order was precise: not to traverse the wormhole but to SET
// COORDINATES on it, BEND it into a lens, and TUNE it to the next
// harmony — to arrive exactly at a new prime. Sunqu does the three
// motions in sequence:
//
//	1. THE VOYAGE (cargo): hunt every zero up to gamma_max = pi * reach
//	   with the Riemann–Siegel mirror — tens of thousands of zeros,
//	   because the compass sees exactly as far as the voyage sails deep
//	   (the reach law of Finding 95).
//	2. SET COORDINATES & BEND (the lens): aim at x0 and apodize the
//	   aperture — Gaussian weights over the zeros bend the exotic lens
//	   so the image rings less and the peaks stand clean.
//	3. TUNE THE NEXT HARMONY: the first strong alignment of the orchestra
//	   past x0 is the landing point. No sieve, no division, no factoring
//	   in the detection — arithmetic is allowed only afterwards, as the
//	   independent verification that the landing point is truly prime.
//
// Usage:
//
//	go run ./cmd/sunqu              # fly to coordinate 15000
//	go run ./cmd/sunqu -x 30000     # fly anywhere within the hunt budget
package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
)

var lnk, rsq [160]float64

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

// zRS is the Riemann–Siegel mirror with the C0 edge snapshot.
func zRS(t float64) float64 {
	tau := t / (2 * math.Pi)
	a := math.Sqrt(tau)
	N := int(a)
	th := theta(t)
	var s float64
	for k := 1; k <= N; k++ {
		s += math.Cos(th-t*lnk[k]) * rsq[k]
	}
	p := a - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (N-1)%2 == 1 {
		sign = -1
	}
	return 2*s + sign*math.Pow(tau, -0.25)*c0
}

func smoothCount(t float64) float64 {
	return t/(2*math.Pi)*(math.Log(t/(2*math.Pi))-1) + 7.0/8
}

func huntZeros(tMax float64) []float64 {
	zeros := []float64{}
	prevT, prevZ := 14.0, zRS(14.0)
	for t := 14.05; t <= tMax; t += 0.05 {
		zt := zRS(t)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for i := 0; i < 25 && hi-lo > 1e-6; i++ {
				mid := (lo + hi) / 2
				if (zRS(mid) < 0) == (prevZ < 0) {
					lo = mid
				} else {
					hi = mid
				}
			}
			zeros = append(zeros, (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}
	return zeros
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			return false
		}
	}
	return true
}

// pratt emits the genealogical passport (F111): the rigorous PROOF of
// primality that IS the ancestor tree. p is proven prime by a witness g
// with g^(p-1) = 1 and g^((p-1)/q) != 1 for every prime parent q of p-1,
// each parent certified recursively - all the way down to Adam (2).
func pratt(p int64, indent int) {
	pad := strings.Repeat("  ", indent)
	if p == 2 {
		fmt.Printf("%s2 — Adán, el origen; primo por definición\n", pad)
		return
	}
	// parents: distinct prime factors of p-1.
	var parents []int64
	m := p - 1
	for q := int64(2); q*q <= m; q++ {
		if m%q == 0 {
			parents = append(parents, q)
			for m%q == 0 {
				m /= q
			}
		}
	}
	if m > 1 {
		parents = append(parents, m)
	}
	// witness: g with g^(p-1)=1 and g^((p-1)/q) != 1 for all parents.
	P := big.NewInt(p)
	Pm1 := big.NewInt(p - 1)
	one := big.NewInt(1)
	var g int64
	for cand := int64(2); cand < p; cand++ {
		G := big.NewInt(cand)
		if new(big.Int).Exp(G, Pm1, P).Cmp(one) != 0 {
			continue
		}
		ok := true
		for _, q := range parents {
			e := new(big.Int).Div(Pm1, big.NewInt(q))
			if new(big.Int).Exp(G, e, P).Cmp(one) == 0 {
				ok = false
				break
			}
		}
		if ok {
			g = cand
			break
		}
	}
	fmt.Printf("%s%d — testigo g=%d; padres de %d:", pad, p, g, p-1)
	for _, q := range parents {
		fmt.Printf(" %d", q)
	}
	fmt.Println()
	for _, q := range parents {
		pratt(q, indent+1)
	}
}

func main() {
	x0f := flag.Float64("x", 15000, "coordinate to aim the lens at")
	flag.Parse()
	x0 := *x0f

	fmt.Println("SUNQU — the heart of the fleet, aimed")
	for k := 1; k < len(lnk); k++ {
		lnk[k] = math.Log(float64(k))
		rsq[k] = 1 / math.Sqrt(float64(k))
	}

	// 1. the voyage: the cargo of zeros that gives the lens its reach.
	const window = 80.0
	gMax := 1.05 * math.Pi * (x0 + window)
	if int(math.Sqrt(gMax/(2*math.Pi))) >= len(lnk) {
		fmt.Println("  coordinate beyond this hull's hunt budget - raise the table size")
		return
	}
	fmt.Printf("\n  1. THE VOYAGE - hunting all zeros up to gamma = %.0f...\n", gMax)
	st := time.Now()
	gammas := huntZeros(gMax)
	expect := smoothCount(gMax) - smoothCount(14)
	fmt.Printf("     cargo: %d zeros in %.1f s (smooth count expects %.0f: %.2f%% complete)\n",
		len(gammas), time.Since(st).Seconds(), expect,
		100*float64(len(gammas))/expect)
	fmt.Printf("     gamma_1 = %.6f (known 14.134725)\n", gammas[0])

	// 2. set coordinates and bend: the apodized exotic lens.
	fmt.Printf("\n  2. COORDINATES SET: x0 = %.0f - bending the lens (Gaussian aperture)...\n", x0)
	w := make([]float64, len(gammas))
	for i, g := range gammas {
		u := g / gMax
		w[i] = math.Exp(-2 * u * u)
	}
	lens := func(x float64) float64 {
		lx := math.Log(x)
		var s float64
		for i, g := range gammas {
			s += w[i] * math.Cos(g*lx)
		}
		return -2 * s / math.Sqrt(x)
	}

	// 3. tune to the next harmony past x0.
	fmt.Println("\n  3. TUNING - listening for the next harmony past the coordinate...")
	st = time.Now()
	const dx = 0.05
	type pk struct{ x, d float64 }
	var best []pk
	maxD := 0.0
	prev := lens(x0 + 0.5)
	cur := lens(x0 + 0.5 + dx)
	for x := x0 + 0.5 + 2*dx; x <= x0+window; x += dx {
		next := lens(x)
		if cur > prev && cur >= next {
			best = append(best, pk{x - dx, cur})
			if cur > maxD {
				maxD = cur
			}
		}
		prev, cur = cur, next
	}
	landing := 0
	for _, p := range best {
		if p.d >= 0.6*maxD {
			landing = int(math.Round(p.x))
			break
		}
	}
	fmt.Printf("     harmony found in %.1f s - SUNQU LANDS AT n = %d\n",
		time.Since(st).Seconds(), landing)

	// independent verification - arithmetic allowed only now.
	verdict := "COMPOSITE - the lens mis-aimed"
	if isPrime(landing) {
		verdict = "PRIME - confirmed by independent arithmetic"
	}
	truth := int(x0) + 1
	for !isPrime(truth) {
		truth++
	}
	fmt.Printf("\n  VERIFICATION: %d is %s\n", landing, verdict)
	fmt.Printf("  the true next prime after %.0f is %d", x0, truth)
	if landing == truth {
		fmt.Println(" - EXACT LANDING: the heart aimed true.")
	} else {
		fmt.Println(" - the landing missed; honesty on display.")
	}
	fmt.Println("\n  no sieve, no division, no factoring touched the detection: the prime")
	fmt.Println("  was found by pointing the folded lens and listening for stillness.")

	if isPrime(landing) {
		fmt.Println("\n  EL PASAPORTE GENEALÓGICO (F111) - primality PROVEN by descent from Adam:")
		pratt(int64(landing), 2)
	}
}
