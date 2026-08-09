// Command fingers builds the Fresnel gearbox: incomplete quadratic sums
// folded recursively — the piece that closes the gap in the painting.
//
// The crankshaft (Finding 88) folds COMPLETE Gauss sums exactly. The
// fingers fold the real thing: S(a,b,L) = Σ_{j<L} e^{2πi(aj+bj²)} with
// arbitrary real parameters. Mechanism: Poisson summation over the
// half-integer window. Interior resonances carry the complete Fresnel
// factor (1+i) and their phases form ANOTHER quadratic sum of length
// ~2bL — the recursion, mirror within mirror — while ~700 edge
// resonances get exact Fresnel integrals. Each level at least halves the
// length.
//
// Version 1 targets zero-hunting grade (absolute error ~1e-3 on sums of
// magnitude ~√L); the bench measures the truth. Full precision and the
// assembly into the flagship are the registered final stage.
//
// Usage:
//
//	go run ./cmd/fingers
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// fresnel returns F(x) = C(x) + i·S(x) = ∫_0^x e^{iπu²/2} du.
func fresnel(x float64) (float64, float64) {
	neg := x < 0
	if neg {
		x = -x
	}
	var c, s float64
	if x <= 3.2 {
		t := math.Pi / 2 * x * x
		baseC := 1.0 // (-1)^n t^{2n} / (2n)!
		baseS := t   // (-1)^n t^{2n+1} / (2n+1)!
		for n := 0; n < 80; n++ {
			c += baseC * x / float64(4*n+1)
			s += baseS * x / float64(4*n+3)
			baseC *= -t * t / float64((2*n+1)*(2*n+2))
			baseS *= -t * t / float64((2*n+2)*(2*n+3))
			if math.Abs(baseC) < 1e-18 && math.Abs(baseS) < 1e-18 {
				break
			}
		}
	} else {
		px := math.Pi * x
		px2 := math.Pi * x * x
		f := 1/px - 3/(px*px2*px2)*math.Pi
		g := 1/(px*px2) - 15/(px*px2*px2*px2)*math.Pi
		sn, cs := math.Sin(px2/2), math.Cos(px2/2)
		c = 0.5 + f*sn - g*cs
		s = 0.5 - f*cs - g*sn
	}
	if neg {
		return -c, -s
	}
	return c, s
}

func frac(x float64) float64 {
	x -= math.Floor(x)
	return x
}

var maxDepth int

// quad computes S(a,b,L) = Σ_{j=0}^{L-1} e^{2πi(aj+bj²)}.
func quad(a, b float64, L int64, depth int) (float64, float64) {
	if depth > maxDepth {
		maxDepth = depth
	}
	a, b = frac(a), frac(b)
	if b > 0.5 {
		b -= 1
	}
	if b < 0 {
		r, i := quad(frac(-a), -b, L, depth)
		return r, -i
	}
	if b > 0.25 {
		// j² ≡ j (mod 2): shift half a turn into the linear term.
		return quad(frac(a+0.5), b-0.5, L, depth)
	}
	if L <= 256 {
		var sr, si float64
		ph, dph := 0.0, 0.0
		for j := int64(0); j < L; j++ {
			sr += math.Cos(2 * math.Pi * ph)
			si += math.Sin(2 * math.Pi * ph)
			dph = frac(a + b*float64(2*j+1))
			ph = frac(ph + dph)
		}
		return sr, si
	}
	if b*float64(L) < 1e-9 {
		// effectively geometric.
		if a < 1e-12 || a > 1-1e-12 {
			return float64(L), 0
		}
		ph := math.Pi * a * float64(L-1)
		amp := math.Sin(math.Pi*a*float64(L)) / math.Sin(math.Pi*a)
		return amp * math.Cos(ph), amp * math.Sin(ph)
	}

	// Poisson over the window [-1/2, L-1/2].
	const edge = 350
	x0, x1 := -0.5, float64(L)-0.5
	d0, d1 := a+2*b*x0, a+2*b*x1
	mLoAll := int64(math.Ceil(d0 - edge))
	mHiAll := int64(math.Floor(d1 + edge))
	mLoIn := int64(math.Ceil(d0 + edge))
	mHiIn := int64(math.Floor(d1 - edge))

	inv2b := 1 / (2 * b)
	sqb2 := 2 * math.Sqrt(b) // the Fresnel change of variable: v = 2*sqrt(b)*(x-xs)
	var sr, si float64

	// edge resonances: exact Fresnel integrals, one single pass.
	for m := mLoAll; m <= mHiAll; m++ {
		if mHiIn >= mLoIn && m >= mLoIn && m <= mHiIn {
			continue // handled by the folded interior
		}
		xs := (float64(m) - a) * inv2b
		v0 := (x0 - xs) * sqb2
		v1 := (x1 - xs) * sqb2
		c1, s1 := fresnel(v1)
		c0, s0 := fresnel(v0)
		fr, fi := c1-c0, s1-s0
		phm := -math.Pi * (float64(m) - a) * (float64(m) - a) * inv2b
		cr, ci := math.Cos(phm), math.Sin(phm)
		sr += (cr*fr - ci*fi) / sqb2
		si += (cr*fi + ci*fr) / sqb2
	}

	// interior resonances: complete Fresnel (1+i) times the folded sum.
	if mHiIn >= mLoIn {
		nL := mHiIn - mLoIn + 1
		alpha := -(float64(mLoIn) - a) * inv2b
		beta := -inv2b / 2
		phi0 := -math.Pi * (float64(mLoIn) - a) * (float64(mLoIn) - a) * inv2b
		gr, gi := quad(frac(alpha), frac(beta), nL, depth+1)
		c0, s0 := math.Cos(phi0), math.Sin(phi0)
		tr := gr*c0 - gi*s0
		ti := gr*s0 + gi*c0
		sr += (tr - ti) / sqb2
		si += (tr + ti) / sqb2
	}
	return sr, si
}

// direct brute-forces S(a,b,L) with incremental exact-fraction phases.
func direct(a, b float64, L int64) (float64, float64) {
	var sr, si float64
	ph := 0.0
	for j := int64(0); j < L; j++ {
		sr += math.Cos(2 * math.Pi * ph)
		si += math.Sin(2 * math.Pi * ph)
		ph = frac(ph + frac(a+b*float64(2*j+1)))
	}
	return sr, si
}

func main() {
	fmt.Println("THE FINGERS — incomplete quadratic sums through the Fresnel gearbox")

	rng := rand.New(rand.NewSource(2026))
	const trials = 60
	worstAbs, worstRel := 0.0, 0.0
	deepest := 0
	var totDirect, totFold time.Duration
	for i := 0; i < trials; i++ {
		a := rng.Float64()
		b := rng.Float64()
		L := int64(rng.Intn(200000) + 20000)
		t0 := time.Now()
		dr, di := direct(a, b, L)
		totDirect += time.Since(t0)
		t0 = time.Now()
		maxDepth = 0
		fr, fi := quad(a, b, L, 0)
		totFold += time.Since(t0)
		if maxDepth > deepest {
			deepest = maxDepth
		}
		abs := math.Hypot(fr-dr, fi-di)
		rel := abs / math.Sqrt(float64(L))
		if abs > worstAbs {
			worstAbs = abs
		}
		if rel > worstRel {
			worstRel = rel
		}
	}
	fmt.Printf("\n  bench: %d random incomplete sums, L up to 220000, full (a,b) range\n", trials)
	fmt.Printf("  worst absolute error: %.2e   worst relative (vs sqrt L): %.2e\n",
		worstAbs, worstRel)
	fmt.Printf("  deepest mirror recursion: %d levels\n", deepest)
	fmt.Printf("  time: direct %.2fs vs folded %.3fs (speedup ~%.0fx)\n",
		totDirect.Seconds(), totFold.Seconds(),
		totDirect.Seconds()/math.Max(totFold.Seconds(), 1e-9))
	fmt.Println("\nif the errors sit at zero-hunting grade, the fingers move: the")
	fmt.Println("assembly into the flagship is the registered final stage.")
}
