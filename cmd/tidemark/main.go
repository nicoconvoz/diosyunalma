// Command tidemark tests Finding 79's candidate at deep water.
//
// The carved constant S₂ = N₂/N₁ − ln ln x extrapolated to ≈ 0.273 — a
// registered candidate identity with Mertens' M = 0.26149. This run
// measures S₂ and S₃ in windows near 10⁹ and 10¹⁰ by windowed
// factorization: every n in the window is divided by all primes up to
// √hi with multiplicity; whatever remains is either 1 or a single large
// prime. Ω exact, no full sieve.
//
// PRE-REGISTERED (from the four-window fit S₂(x) = S∞ − c/ln x with
// S∞ = 0.273, c = 1.75): S₂ lands in [0.17, 0.21] near 10⁹ and in
// [0.18, 0.22] near 10¹⁰, still climbing; the six-point refit keeps
// S∞ inside [0.23, 0.30] — Mertens' tide-mark — or the candidate dies.
//
// Usage:
//
//	go run ./cmd/tidemark
package main

import (
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

const width = 20_000_000

func rungs(lo int) (n1, n2, n3 float64) {
	hi := lo + width
	rem := make([]int64, width)
	omega := make([]uint8, width)
	for i := range rem {
		rem[i] = int64(lo + i)
	}
	base := primes.Sieve(int(math.Sqrt(float64(hi))) + 1)
	for _, p := range base {
		pp := int64(p)
		start := (lo + p - 1) / p * p
		for m := start; m < hi; m += p {
			i := m - lo
			for rem[i]%pp == 0 {
				rem[i] /= pp
				omega[i]++
			}
		}
	}
	for i := range rem {
		if rem[i] > 1 {
			omega[i]++
		}
		switch omega[i] {
		case 1:
			n1++
		case 2:
			n2++
		case 3:
			n3++
		}
	}
	return
}

func main() {
	fmt.Println("THE TIDEMARK — the carved constant tested at deep water")
	fmt.Println("\n   window        lnlnx    S2       S3      registered band for S2")

	type pt struct{ lnx, s2 float64 }
	pts := []pt{
		{math.Log(150_000), 0.1258}, {math.Log(1_500_000), 0.1457},
		{math.Log(15_000_000), 0.1654}, {math.Log(150_000_000), 0.1800},
	}
	bands := map[int][2]float64{1_000_000_000: {0.17, 0.21}, 10_000_000_000: {0.18, 0.22}}
	for _, lo := range []int{1_000_000_000, 10_000_000_000} {
		n1, n2, n3 := rungs(lo)
		mid := float64(lo) + width/2
		lam := math.Log(math.Log(mid))
		s2 := n2/n1 - lam
		s3 := 2*n3/n2 - lam
		b := bands[lo]
		verdict := "MISS"
		if s2 >= b[0] && s2 <= b[1] {
			verdict = "inside"
		}
		fmt.Printf("   [%.0e]      %.3f   %+.4f  %+.4f   [%.2f, %.2f] %s\n",
			float64(lo), lam, s2, s3, b[0], b[1], verdict)
		pts = append(pts, pt{math.Log(mid), s2})
	}

	// six-point refit of S2 = Sinf - c/lnx.
	var sx, sy, sxx, sxy float64
	for _, p := range pts {
		x := 1 / p.lnx
		sx += x
		sy += p.s2
		sxx += x * x
		sxy += x * p.s2
	}
	n := float64(len(pts))
	c := -(n*sxy - sx*sy) / (n*sxx - sx*sx)
	sinf := (sy + c*sx) / n
	fmt.Printf("\nsix-point refit: S2(x) = %.4f - %.2f/ln x\n", sinf, c)
	fmt.Printf("the tide-mark candidate: S_inf = %.4f   (Mertens' M = 0.26149)\n", sinf)
	if sinf > 0.23 && sinf < 0.30 {
		fmt.Println("\nVERDICT: the candidate lives - the lagoon's constant is carved in the wall.")
	} else {
		fmt.Println("\nVERDICT: the candidate dies - the carved constant is not the tide-mark.")
	}
}
