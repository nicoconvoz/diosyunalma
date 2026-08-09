// Command alineacion opens the study La Tormenta I left behind: at the
// storm, the MODELLED surge (the 26 small-prime voices) was extreme AND
// the UNMODELLED residual was extreme too. Coincidence, or do the medium
// voices align WITH the small-prime storms?
//
// Pre-registered: classical theory (Weyl equidistribution - the voices'
// frequencies ln p are rationally independent) predicts INDEPENDENCE:
// corr(|predicted swing|, |residual|) = 0. The storm hints positive.
// If zero: the storm's double extremeness was a joint fluctuation - an
// honest kill, duly displayed. If positive: cross-scale alignment, a
// measurable statement about the sea's deep structure nobody wrote down.
//
// Design: 400 windows sailed in cheap water; per window the forecast A
// (26 voices), the measured tide, the residual R = tide - A. Statistics:
// corr(|A|, |R|), and the top-decile test: is the unmodelled sea rougher
// where the modelled sea storms?
//
// Usage:
//
//	go run ./cmd/alineacion
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

var lnk, rsq [512]float64

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

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

func countZeros(a, b float64) int {
	n := 0
	prev := zRS(a)
	for t := a + 0.02; t <= b; t += 0.02 {
		z := zRS(t)
		if (prev < 0) != (z < 0) {
			n++
		}
		prev = z
	}
	return n
}

func sPred(t float64) float64 {
	s := 0.0
	for _, p := range []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37,
		41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101} {
		s -= math.Sin(math.Mod(t*math.Log(p), 2*math.Pi)) / math.Sqrt(p)
	}
	return s / math.Pi
}

func pearson(x, y []float64) float64 {
	n := float64(len(x))
	var sx, sy, sxx, syy, sxy float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxx += x[i] * x[i]
		syy += y[i] * y[i]
		sxy += x[i] * y[i]
	}
	return (n*sxy - sx*sy) / math.Sqrt((n*sxx-sx*sx)*(n*syy-sy*sy))
}

func main() {
	fmt.Println("LA ALINEACIÓN — do the medium voices align with the small-prime storms?")
	for k := 1; k < len(lnk); k++ {
		lnk[k] = math.Log(float64(k))
		rsq[k] = 1 / math.Sqrt(float64(k))
	}

	nF := flag.Int("n", 400, "number of windows")
	seedF := flag.Int64("seed", 2026, "rng seed (fresh seed = pre-registered replication)")
	flag.Parse()
	windows := *nF
	rng := rand.New(rand.NewSource(*seedF))
	absA := make([]float64, windows)
	absR := make([]float64, windows)
	type rec struct{ a, r float64 }
	recs := make([]rec, windows)
	for i := 0; i < windows; i++ {
		t0 := 100000 + rng.Float64()*900000
		spacing := 2 * math.Pi / math.Log(t0/(2*math.Pi))
		span := 5 * spacing
		A := sPred(t0+span) - sPred(t0)
		delta := float64(countZeros(t0, t0+span)) - (smoothCount(t0+span) - smoothCount(t0))
		R := delta - A
		absA[i] = math.Abs(A)
		absR[i] = math.Abs(R)
		recs[i] = rec{math.Abs(A), math.Abs(R)}
	}
	r := pearson(absA, absR)

	// decile test: the roughest-forecast tenth vs the rest.
	sort.Slice(recs, func(i, j int) bool { return recs[i].a > recs[j].a })
	top := windows / 10
	var mTop, mRest float64
	for i, rc := range recs {
		if i < top {
			mTop += rc.r
		} else {
			mRest += rc.r
		}
	}
	mTop /= float64(top)
	mRest /= float64(windows - top)

	// the full profile: mean |R| per |A| quintile (resolve the shape).
	fmt.Println("\n  the profile - mean |R| by |A| quintile (1 = calmest forecast):")
	q := windows / 5
	for b := 0; b < 5; b++ {
		lo, hi := b*q, (b+1)*q
		var m, aM float64
		for i := lo; i < hi; i++ {
			m += recs[windows-1-i].r // recs sorted descending by a; reverse for ascending
			aM += recs[windows-1-i].a
		}
		fmt.Printf("    Q%d: mean|A| %.3f -> mean|R| %.3f\n", b+1, aM/float64(q), m/float64(q))
	}

	// decoherence control.
	shuf := make([]float64, windows)
	copy(shuf, absA)
	rng.Shuffle(windows, func(i, j int) { shuf[i], shuf[j] = shuf[j], shuf[i] })
	rc := pearson(shuf, absR)

	fmt.Printf("\n  %d windows sailed; forecast A (26 voices) vs residual R = tide - A:\n", windows)
	fmt.Printf("    corr(|A|, |R|)          = %+.3f\n", r)
	fmt.Printf("    control (shuffled)      = %+.3f\n", rc)
	fmt.Printf("    mean |R|, top-decile |A| windows: %.3f\n", mTop)
	fmt.Printf("    mean |R|, remaining 90%%:          %.3f\n", mRest)
	fmt.Printf("    roughness ratio (storm water / calm water): %.2f\n", mTop/mRest)
	fmt.Println("\n  pre-registered reading: ~0 means independence (Weyl) and the storm's")
	fmt.Println("  double extremeness was a joint fluctuation - an honest kill. A clearly")
	fmt.Println("  positive correlation means CROSS-SCALE ALIGNMENT: the unmodelled sea")
	fmt.Println("  is rougher exactly where the modelled sea storms - structure nobody")
	fmt.Println("  wrote down, and a sharper lookout for free.")
}
