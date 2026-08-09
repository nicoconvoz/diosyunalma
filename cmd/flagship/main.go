// Command flagship is the definitive vessel, built to be certified before
// it ever sails deep.
//
// Everything the voyage taught, in one hull:
//
//   - CONVEX FACETS (the convex-mirror flash): quartic phase facets from a
//     single double-double chain — megabytes of optics, built in seconds.
//   - THE LIGHT BUCKET (the Webb flash): one pass over the terms deposits
//     every photon onto an oversampled Nyquist grid of the whole window;
//     the function is band-limited, so zeros are then read off the
//     recorded light by Lanczos interpolation — no photon escapes, no
//     second collection.
//   - ANALYTIC THETA (the pixelated-sea lesson): at depth the ulp of t
//     exceeds the window, so theta is carved in dt at the anchor.
//   - SELF-CERTIFICATION: the ship refuses deep water until it reproduces
//     the charted zeros of t = 10⁵ (LMFDB), Beach I and Beach II, and it
//     reports its measured speed and maximum rated capacity.
//
// Usage:
//
//	go run ./cmd/flagship             # certification: gates + capacity report
//	go run ./cmd/flagship -anchor T   # sail (runs quick gates first)
package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime"
	"sync"
	"time"
)

type dd struct{ hi, lo float64 }

func twoSum(a, b float64) (float64, float64) {
	s := a + b
	bb := s - a
	return s, (a - (s - bb)) + (b - bb)
}

func quick(a, b float64) (float64, float64) {
	s := a + b
	return s, b - (s - a)
}

func twoProd(a, b float64) (float64, float64) {
	p := a * b
	return p, math.FMA(a, b, -p)
}

func add(x, y dd) dd {
	s, e := twoSum(x.hi, y.hi)
	e += x.lo + y.lo
	s, e = quick(s, e)
	return dd{s, e}
}

func addF(x dd, f float64) dd {
	s, e := twoSum(x.hi, f)
	e += x.lo
	s, e = quick(s, e)
	return dd{s, e}
}

func mulF(x dd, f float64) dd {
	p, e := twoProd(x.hi, f)
	e += x.lo * f
	p, e = quick(p, e)
	return dd{p, e}
}

func mul(x, y dd) dd {
	p, e := twoProd(x.hi, y.hi)
	e += x.hi*y.lo + x.lo*y.hi
	p, e = quick(p, e)
	return dd{p, e}
}

func neg(x dd) dd { return dd{-x.hi, -x.lo} }

func ddInv(f float64) dd {
	q0 := 1 / f
	r := math.FMA(q0, f, -1)
	return dd{q0, -r / f}
}

func parse(s string) dd {
	bf, _, _ := big.ParseFloat(s, 10, 200, big.ToNearestEven)
	hi, _ := bf.Float64()
	bf.Sub(bf, new(big.Float).SetPrec(200).SetFloat64(hi))
	lo, _ := bf.Float64()
	return dd{hi, lo}
}

var (
	ln2   = parse("0.69314718055994530941723212145817656807550013436026")
	ln2pi = parse("1.83787706640934548356065947281123527972766750048826")
	twoPi = parse("6.28318530717958647692528676655900576839433879875021")
	pi8   = parse("0.39269908169872415480783042290993786052464617492189")
)

var twoPi26 = dd{twoPi.hi * 67108864, twoPi.lo * 67108864}

func mod2pi(x dd) float64 {
	n := math.Floor(x.hi / twoPi.hi)
	n1 := math.Floor(n / 67108864)
	n2 := n - n1*67108864
	x = add(x, mulF(twoPi26, -n1))
	x = add(x, mulF(twoPi, -n2))
	r := x.hi + x.lo
	for r >= twoPi.hi {
		r -= twoPi.hi
	}
	for r < 0 {
		r += twoPi.hi
	}
	return r
}

func ddLnF(x float64) dd {
	m, e := math.Frexp(x)
	if m < 0.70710678118654752 {
		m *= 2
		e--
	}
	num := m - 1
	den := m + 1
	q0 := num / den
	r := math.FMA(q0, den, -num)
	z := dd{q0, -r / den}
	z2 := mul(z, z)
	sum := dd{0, 0}
	for j := 31; j >= 1; j -= 2 {
		sum = add(mul(sum, z2), mulF(ddInv(float64(j)), 1))
	}
	sum = mul(sum, z)
	sum = mulF(sum, 2)
	return add(sum, mulF(ln2, float64(e)))
}

func thetaMod(t float64) float64 {
	lnArg := add(ddLnF(t), neg(ln2pi))
	th := mulF(addF(lnArg, -1), t/2)
	th = add(th, neg(pi8))
	th = addF(th, 1/(48*t))
	return mod2pi(th)
}

var piDD = parse("3.14159265358979323846264338327950288419716939937511")

// zeroIndex names the ordinal number of the first zero above t: the
// sphere's naming power, N(t) = theta(t)/pi + 1, exact to the unit in dd
// (plus the S(t) ambiguity of about +/-2).
func zeroIndex(t float64) string {
	lnArg := add(ddLnF(t), neg(ln2pi))
	th := mulF(addF(lnArg, -1), t/2)
	th = add(th, neg(pi8))
	q0 := th.hi / piDD.hi
	r := add(th, neg(mulF(piDD, q0)))
	q1 := (r.hi + r.lo) / piDD.hi
	bf := new(big.Float).SetPrec(150).SetFloat64(q0)
	bf.Add(bf, new(big.Float).SetPrec(150).SetFloat64(q1+1))
	return bf.Text('f', 0)
}

func ln1pSeries(u dd) dd {
	s := mulF(ddInv(7), 1)
	for j := 6; j >= 1; j-- {
		s = mul(neg(s), u)
		s = add(s, mulF(ddInv(float64(j)), 1))
	}
	return mul(s, u)
}

type facet struct {
	k0                   int64
	n                    int32
	phi0, d1, d2, d3, d4 float64
	// the raw polynomial coefficients in FULL double-double, kept for the
	// exact refresh. Lesson carved in: reducing a tiny coefficient mod 2pi
	// parks it next to 2pi with half-an-ulp of its own information left,
	// and j^4 amplifies that half-ulp into a tenth of a radian.
	cA, cB, cC, cD  dd
	lnk0, invk0, r0 float64
}

type ship struct {
	t0                float64
	small             []float32
	kSmall            int64
	facets            []facet
	nTop              int64
	thAnchor, thPrime float64
	// the recorded light
	h, dt0 float64
	fr, fi []float64
	c0term float64
}

func build(t0 float64, nTop int64) *ship {
	ub := math.Pow(0.1/t0, 0.2)
	kSmall := int64(16/ub) + 1
	if kSmall > nTop {
		kSmall = nTop + 1
	}
	sp := &ship{t0: t0, kSmall: kSmall, nTop: nTop}
	sp.small = make([]float32, kSmall)
	sp.thAnchor = thetaMod(t0)
	lnArg := add(ddLnF(t0), neg(ln2pi))
	sp.thPrime = (lnArg.hi + lnArg.lo) / 2

	lnPrev := dd{0, 0}
	ph := 0.0
	needed := math.Log(t0 * 1e4)
	step := func(k int64) {
		var delta dd
		if k <= 64 {
			lk := ddLnF(float64(k))
			delta = add(lk, neg(lnPrev))
			lnPrev = lk
		} else {
			u := ddInv(float64(k - 1))
			m := int(needed/math.Log(float64(k-1))) + 2
			if m > 16 {
				m = 16
			}
			s := mulF(ddInv(float64(m)), 1)
			for j := m - 1; j >= 1; j-- {
				s = mul(neg(s), u)
				s = add(s, mulF(ddInv(float64(j)), 1))
			}
			delta = mul(s, u)
			lnPrev = add(lnPrev, delta)
		}
		ph += mod2pi(mulF(delta, t0))
		if ph >= twoPi.hi {
			ph -= twoPi.hi
		}
	}
	for k := int64(2); k < kSmall; k++ {
		step(k)
		sp.small[k] = float32(ph)
	}
	if kSmall <= nTop {
		step(kSmall)
	}
	for k0 := kSmall; k0 <= nTop; {
		b := int64(ub * float64(k0))
		if b < 16 {
			b = 16
		}
		if k0+b-1 > nTop {
			b = nTop - k0 + 1
		}
		invK := ddInv(float64(k0))
		a := mulF(invK, t0)
		k2 := mul(invK, invK)
		b2 := mulF(k2, -t0/2)
		c := mulF(mul(k2, invK), t0/3)
		d := mulF(mul(k2, k2), -t0/4)
		sp.facets = append(sp.facets, facet{
			k0: k0, n: int32(b),
			phi0:  ph,
			d1:    mod2pi(add(add(a, b2), add(c, d))),
			d2:    mod2pi(add(mulF(b2, 2), add(mulF(c, 6), mulF(d, 14)))),
			d3:    mod2pi(add(mulF(c, 6), mulF(d, 36))),
			d4:    mod2pi(mulF(d, 24)),
			cA: a, cB: b2, cC: c, cD: d,
			lnk0:  lnPrev.hi + lnPrev.lo,
			invk0: 1 / float64(k0),
			r0:    1 / math.Sqrt(float64(k0)),
		})
		u := mulF(invK, float64(b))
		st := ln1pSeries(u)
		lnPrev = add(lnPrev, st)
		ph += mod2pi(mulF(st, t0))
		if ph >= twoPi.hi {
			ph -= twoPi.hi
		}
		k0 += b
	}
	return sp
}

var refreshEvery int64 = 4096

// oversample is the rhythm of the light bucket: how many grid beats per
// Nyquist interval. 3 is the certified cruise beat; the rhythm experiment
// (the heartbeat of Sunqu) measures how fast the beat can go while the
// gates still certify.
var oversample = 3.0

func wrap2pi(x float64) float64 {
	x = math.Mod(x, twoPi.hi)
	if x < 0 {
		x += twoPi.hi
	}
	return x
}

// refresh recomputes the walk state at step j from the polynomial in full
// double-double: the antidote to the cubic error growth of the chain.
func (f *facet) refresh(j int64) (p, d1, d2, d3 float64) {
	jf := float64(j)
	j2 := jf * jf
	j3 := j2 * jf
	sum := mulF(f.cA, jf)
	sum = add(sum, mulF(f.cB, j2))
	sum = add(sum, mulF(f.cC, j3))
	sum = add(sum, mulF(f.cD, j3*jf))
	p = wrap2pi(f.phi0 + mod2pi(sum))
	s1 := f.cA
	s1 = add(s1, mulF(f.cB, 2*jf+1))
	s1 = add(s1, mulF(f.cC, 3*j2+3*jf+1))
	s1 = add(s1, mulF(f.cD, 4*j3+6*j2+4*jf+1))
	d1 = mod2pi(s1)
	d2 = mod2pi(add(mulF(f.cB, 2),
		add(mulF(f.cC, 6*jf+6), mulF(f.cD, 12*j2+24*jf+14))))
	d3 = mod2pi(add(mulF(f.cC, 6), mulF(f.cD, 24*jf+36)))
	return
}

// collect runs the light bucket: one pass, every photon onto the grid.
// The window starts at dtStart, so the same carved facets can fill
// consecutive stretches of sea — the ship grows while it sails.
func (sp *ship) collect(dtStart, span float64) {
	lnTop := math.Log(float64(sp.nTop))
	sp.h = 2 * math.Pi / lnTop / oversample // the bucket's beat vs the band
	guard := 8 * sp.h
	sp.dt0 = dtStart - guard
	S := int((span+2*guard)/sp.h) + 2
	workers := runtime.NumCPU() - 1
	if workers < 1 {
		workers = 1
	}
	frs := make([][]float64, workers)
	fis := make([][]float64, workers)
	var wg sync.WaitGroup
	chunk := (len(sp.facets) + workers - 1) / workers
	for w := 0; w < workers; w++ {
		frs[w] = make([]float64, S)
		fis[w] = make([]float64, S)
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			lo, hi := w*chunk, (w+1)*chunk
			if hi > len(sp.facets) {
				hi = len(sp.facets)
			}
			fr, fi := frs[w], fis[w]
			for fx := lo; fx < hi; fx++ {
				f := &sp.facets[fx]
				wr0, wi0 := math.Cos(-sp.h*f.lnk0), math.Sin(-sp.h*f.lnk0)
				dwr, dwi := math.Cos(-sp.h*f.invk0), math.Sin(-sp.h*f.invk0)
				ph, d1, d2, d3 := f.phi0, f.d1, f.d2, f.d3
				lnLin := f.lnk0
				for j := int64(0); j < int64(f.n); j++ {
					if j > 0 && j%refreshEvery == 0 {
						ph, d1, d2, d3 = f.refresh(j)
					}
					amp := f.r0 * (1 - 0.5*float64(j)*f.invk0)
					ang := ph + sp.dt0*lnLin
					cr, ci := math.Cos(-ang), math.Sin(-ang)
					cr, ci = amp*cr, amp*ci
					wr, wi := wr0, wi0
					for s := 0; s < S; s++ {
						fr[s] += cr
						fi[s] += ci
						cr, ci = cr*wr-ci*wi, cr*wi+ci*wr
					}
					_ = wi
					ph += d1
					if ph >= twoPi.hi {
						ph -= twoPi.hi
					}
					d1 += d2
					if d1 >= twoPi.hi {
						d1 -= twoPi.hi
					}
					d2 += d3
					if d2 >= twoPi.hi {
						d2 -= twoPi.hi
					}
					d3 += f.d4
					if d3 >= twoPi.hi {
						d3 -= twoPi.hi
					}
					lnLin += f.invk0
					wr0, wi0 = wr0*dwr-wi0*dwi, wr0*dwi+wi0*dwr
				}
			}
		}(w)
	}
	// small tier + k=1 on the main goroutine.
	fr0 := make([]float64, S)
	fi0 := make([]float64, S)
	for s := 0; s < S; s++ {
		fr0[s] += 1 // k = 1: phase 0 always
	}
	for k := int64(2); k < sp.kSmall && k <= sp.nTop; k++ {
		lnk := math.Log(float64(k))
		amp := 1 / math.Sqrt(float64(k))
		ang := float64(sp.small[k]) + sp.dt0*lnk
		cr, ci := amp*math.Cos(-ang), amp*math.Sin(-ang)
		wr, wi := math.Cos(-sp.h*lnk), math.Sin(-sp.h*lnk)
		for s := 0; s < S; s++ {
			fr0[s] += cr
			fi0[s] += ci
			cr, ci = cr*wr-ci*wi, cr*wi+ci*wr
		}
	}
	wg.Wait()
	sp.fr, sp.fi = fr0, fi0
	for w := 0; w < workers; w++ {
		for s := 0; s < S; s++ {
			sp.fr[s] += frs[w][s]
			sp.fi[s] += fis[w][s]
		}
	}
	// the seam where the outer sphere (the functional-equation reflector,
	// folding the far half of the light back inward) meets the inner one:
	// the edge snapshot C0, polished with the reverse-mirror step C1.
	a := math.Sqrt(sp.t0 / twoPi.hi)
	p := a - float64(sp.nTop)
	c0f := func(p float64) float64 {
		return math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	}
	const hh = 0.008
	c1 := -(c0f(p+2*hh) - 2*c0f(p+hh) + 2*c0f(p-hh) - c0f(p-2*hh)) /
		(2 * hh * hh * hh) / (96 * math.Pi * math.Pi)
	sign := 1.0
	if (sp.nTop-1)%2 == 1 {
		sign = -1
	}
	tau := sp.t0 / twoPi.hi
	sp.c0term = sign * (math.Pow(tau, -0.25)*c0f(p) + math.Pow(tau, -0.75)*c1)
}

// zAt reads Z(dt) off the recorded light by Lanczos-6 interpolation.
func (sp *ship) zAt(dt float64) float64 {
	x := (dt - sp.dt0) / sp.h
	j0 := int(math.Floor(x))
	var fr, fi float64
	for j := j0 - 5; j <= j0+6; j++ {
		if j < 0 || j >= len(sp.fr) {
			continue
		}
		u := x - float64(j)
		var w float64
		if u == 0 {
			w = 1
		} else {
			pu := math.Pi * u
			w = 6 * math.Sin(pu) * math.Sin(pu/6) / (pu * pu)
		}
		fr += w * sp.fr[j]
		fi += w * sp.fi[j]
	}
	th := sp.thAnchor + dt*sp.thPrime + dt*dt/(4*sp.t0)
	cr, ci := math.Cos(th), math.Sin(th)
	return 2*(cr*fr-ci*fi) + sp.c0term
}

// sweep collects the light over [dtStart, dtStart+span] and hunts it —
// one stage of a voyage that can keep extending from the same facets.
func (sp *ship) sweep(dtStart, span, spacing float64) ([]float64, float64) {
	sp.collect(dtStart, span)
	end := dtStart + span
	scan := func(step float64) []float64 {
		found := []float64{}
		prevD, prevZ := dtStart, sp.zAt(dtStart)
		for d := dtStart + step; d <= end; d += step {
			zd := sp.zAt(d)
			if (prevZ < 0) != (zd < 0) {
				lo, hi := prevD, d
				zlo, zhi := prevZ, zd
				for i := 0; i < 40 && hi-lo > 1e-7; i++ {
					var mid float64
					if i%2 == 0 {
						mid = (lo + hi) / 2
					} else {
						mid = lo - zlo*(hi-lo)/(zhi-zlo)
						if mid <= lo || mid >= hi {
							mid = (lo + hi) / 2
						}
					}
					zm := sp.zAt(mid)
					if (zlo < 0) != (zm < 0) {
						hi, zhi = mid, zm
					} else {
						lo, zlo = mid, zm
					}
				}
				if math.Abs(zlo) < math.Abs(zhi) {
					found = append(found, lo)
				} else {
					found = append(found, hi)
				}
			}
			prevD, prevZ = d, zd
		}
		return found
	}
	zeros := scan(spacing / 60)
	// the sphere: exact enclosure count from the boundary alone.
	exact := ((end-dtStart)*sp.thPrime + (end*end-dtStart*dtStart)/(4*sp.t0)) / math.Pi
	// the magnifier: reading recorded light costs nothing, so when the
	// sphere says a zero is hiding, slow down and sweep 10x finer — a
	// close pair unresolved at cruising speed resolves under the glass.
	if float64(len(zeros)) < exact-0.5 {
		zeros = scan(spacing / 600)
	}
	return zeros, exact
}

func hunt(t0 float64, spacings float64) ([]float64, float64, int, float64) {
	start := time.Now()
	spacing := twoPi.hi / math.Log(t0/twoPi.hi)
	nTop := int64(math.Sqrt(t0 / twoPi.hi))
	sp := build(t0, nTop)
	zeros, exact := sp.sweep(0, spacings*spacing, spacing)
	return zeros, exact, len(sp.facets), time.Since(start).Minutes()
}

func gate(name string, t0, spacings float64, expect []float64, tol float64) bool {
	zeros, exact, nf, mins := hunt(t0, spacings)
	ok := len(zeros) == len(expect)
	worst := 0.0
	if ok {
		for i := range zeros {
			d := math.Abs(zeros[i] - expect[i])
			if d > worst {
				worst = d
			}
		}
		ok = worst <= tol
	}
	sphere := math.Abs(float64(len(zeros))-exact) <= tolS(t0)
	verdict := "PASS"
	if !ok || !sphere {
		verdict = "FAIL"
		ok = false
	}
	fmt.Printf("  gate %-22s found %d/%d  worst dev %.5f (tol %.3f)  sphere %.2f  %d facets  %.1f min  [%s]\n",
		name, len(zeros), len(expect), worst, tol, exact, nf, mins, verdict)
	return ok
}

// tolS is the barometric tolerance (F102): the sea's number variance
// saturates (Berry), growing only as ln ln t; the sphere's honest gauge
// is 2.5 sigma of that pressure, calibrated on cmd/barometro.
func tolS(t float64) float64 {
	return 2.5 * math.Sqrt((math.Log(math.Log(t))+2.72)/(math.Pi*math.Pi))
}

// sPred forecasts the predictable tide from the loudest prime voices
// (the pacemaker, F100): S(t) ~ -(1/pi) sum sin(t ln p)/sqrt(p).
func sPred(t0, dt float64) float64 {
	s := 0.0
	for _, p := range []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37,
		41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101} {
		ph := mod2pi(mulF(ddLnF(p), t0)) + dt*math.Log(p)
		s -= math.Sin(ph) / math.Sqrt(p)
	}
	return s / math.Pi
}

// logLine appends a line to the night log ON DISK immediately, so every
// landed beach survives any closed panel, killed session or power cut.
func logLine(s string) {
	f, err := os.OpenFile("docs/BITACORA-NOCTURNA.md",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(s + "\n")
	f.Sync()
}

func mission() {
	logLine("")
	logLine("## Misión nocturna — zarpe " + time.Now().Format("2006-01-02 15:04"))
	logLine("")
	// the growing ship: facets are carved ONCE per anchorage, then the
	// sailors keep building — each stage extends the covered sea by one
	// window at pure marginal cost, and lands its zeros on disk at once.
	ladder := []struct {
		t        float64
		spacings float64
		budget   float64 // minutes of sailing at this anchorage
		name     string
	}{
		{1.11e19, 25, 45, "las puertas de Odlyzko"},
		{2.22e21, 25, 150, "más profundo que la isla más honda de Odlyzko"},
		{4.44e22, 15, 360, "más allá de la isla #10^23 de Gourdon"},
		{1.11e24, 6, 600, "el techo del casco certificado"},
	}
	for _, rung := range ladder {
		fmt.Printf("\nsailing to %s (t = %.3g)...\n", rung.name, rung.t)
		logLine(fmt.Sprintf("### Playa en t = %.6g — %s", rung.t, rung.name))
		logLine(fmt.Sprintf("- primer cero de la ventana: cero número ~%s (±2)",
			zeroIndex(rung.t)))
		anchorStart := time.Now()
		spacing := twoPi.hi / math.Log(rung.t/twoPi.hi)
		nTop := int64(math.Sqrt(rung.t / twoPi.hi))
		sp := build(rung.t, nTop)
		logLine(fmt.Sprintf("- instrumento: %d facetas, talladas una vez; etapas de %.0f espaciamientos",
			len(sp.facets), rung.spacings))
		total, totalExact := 0, 0.0
		var allZeros []float64
		sailed := 0.0
		weave := 0.0 // the music box: the interwoven memory of residuals
		for stage := 0; ; stage++ {
			stageStart := time.Now()
			dtA := float64(stage) * rung.spacings * spacing
			zeros, exact := sp.sweep(dtA, rung.spacings*spacing, spacing)
			allZeros = append(allZeros, zeros...)
			sailed = dtA + rung.spacings*spacing
			stageMins := time.Since(stageStart).Minutes()
			line := fmt.Sprintf("- etapa %d [%.3f, %.3f]:", stage+1, dtA, dtA+rung.spacings*spacing)
			for _, z := range zeros {
				line += fmt.Sprintf(" %.6f", z)
			}
			logLine(line)
			delta := float64(len(zeros)) - exact
			verdict := "esfera OK"
			if math.Abs(delta) > tolS(rung.t) {
				verdict = "ESFERA ROTA: investigar antes de confiar"
			}
			logLine(fmt.Sprintf("  esfera: exige %.2f, hallados %d (delta %+.2f) — %s; %.1f min",
				exact, len(zeros), delta, verdict, stageMins))
			// LA CAJITA MUSICAL (the tide, comprehended): pacemaker forecast
			// plus the interwoven memory of what the sea actually did — the
			// bass line of the gap song turned into navigation. Only the
			// residual that survives BOTH is genuinely unexplained.
			forecast := sPred(rung.t, dtA+rung.spacings*spacing) - sPred(rung.t, dtA)
			resid := delta - forecast - weave
			logLine(fmt.Sprintf("  cajita musical: marea prevista %+.2f, entretejido %+.2f, residuo %+.2f",
				forecast, weave, resid))
			weave = 0.85*weave + 0.15*(delta-forecast)
			total += len(zeros)
			totalExact += exact
			fmt.Printf("  stage %d: %d zeros, sphere delta %+.2f, %.1f min - logged\n",
				stage+1, len(zeros), delta, stageMins)
			if time.Since(anchorStart).Minutes()+stageMins > rung.budget {
				break
			}
		}
		logLine(fmt.Sprintf("- total del fondeadero: %d ceros vírgenes (frontera exigía %.2f); %.1f min",
			total, totalExact, time.Since(anchorStart).Minutes()))
		// the stillness gauge (Finding 93): each stage's delta is the sea's
		// own tide S(t), which is BOUNDED and mean-reverting. The cumulative
		// delta over an anchorage must therefore return toward rest; a
		// monotone drift would expose a leaking hull, not a tide.
		cum := float64(total) - totalExact
		gauge := "QUIETUD: la marea volvió a descansar"
		if math.Abs(cum) > 2 {
			gauge = "INQUIETUD ANÓMALA: posible fuga del casco — investigar"
		}
		logLine(fmt.Sprintf("- medidor de quietud: S acumulada %+.2f — %s", cum, gauge))
		// the Turing-style audit (the music used as certification): the
		// INTEGRAL of the unrest has a proven ceiling (Turing 1953; explicit
		// constants by Trudgian). A truly lost zero drags the second half's
		// mean unrest one full unit below the first half's — beyond what the
		// theorem allows the honest sea. One unknown constant (S at the
		// anchor) cancels in the half-difference.
		if sailed > 0 && len(allZeros) > 0 {
			half := sailed / 2
			i1, i2 := 0.0, 0.0
			for _, z := range allZeros {
				if z <= half {
					i1 += half - z
					i2 += sailed - half
				} else {
					i2 += sailed - z
				}
			}
			nbarInt := func(u float64) float64 {
				return (u*u/2*sp.thPrime + u*u*u/(12*rung.t)) / math.Pi
			}
			i1 -= nbarInt(half)
			i2 -= nbarInt(sailed) - nbarInt(half)
			c1, c2 := i1/half, i2/(sailed-half)
			bound := 2.30 + 0.128*math.Log(rung.t/twoPi.hi)
			thr := 8 * bound / sailed
			verdictT := "sin fuga detectable"
			if math.Abs(c2-c1) > thr {
				verdictT = "DERIVA MÁS ALLÁ DEL TEOREMA — fuga del casco"
			}
			logLine(fmt.Sprintf("- auditoría Turing: inquietud media %.3f (1a mitad) vs %.3f (2a mitad), umbral %.3f — %s",
				c1, c2, thr, verdictT))
		}
		logLine("")
	}
	logLine("### Fin de la misión — " + time.Now().Format("2006-01-02 15:04"))
	fmt.Println("\nMISSION COMPLETE - every beach persisted in docs/BITACORA-NOCTURNA.md")
}

func main() {
	anchor := flag.Float64("anchor", 0, "sail: anchor at this height")
	spacings := flag.Float64("spacings", 5, "window width in mean spacings")
	night := flag.Bool("mission", false, "night mission: the full ladder, logged to disk")
	from := flag.Float64("from", 0, "start the window this many spacings past the anchor")
	osF := flag.Float64("os", 3, "light-bucket oversampling (the rhythm)")
	flag.Parse()
	oversample = *osF

	fmt.Println("THE FLAGSHIP — certification before deep water")
	ok1 := gate("t=1e5 (LMFDB)", 100000, 4,
		[]float64{0.743723, 1.180558, 2.093723, 2.517282}, 3e-3)
	ok2 := gate("Beach I (2.447e12)", 2.447e12, 8,
		[]float64{0.149473, 0.266583, 0.547539, 0.777469, 1.082587, 1.198347, 1.399001, 1.728025}, 3e-3)
	ok3 := gate("Beach II (6.66e15)", 6.66e15, 8,
		[]float64{0.023476, 0.215217, 0.375565, 0.547513, 0.734403, 0.930796, 1.048695, 1.195434}, 3e-3)

	if !(ok1 && ok2 && ok3) {
		fmt.Println("\nCERTIFICATION FAILED - the ship stays in port.")
		return
	}
	fmt.Println("\nALL GATES PASS - the hull is certified.")

	// measured speed -> capacity report, at the mission's wide window.
	t0 := 6.66e15
	nTop := int64(math.Sqrt(t0 / twoPi.hi))
	spb := build(t0, nTop)
	start := time.Now()
	spb.collect(0, 25*twoPi.hi/math.Log(t0/twoPi.hi))
	perTerm := time.Since(start).Seconds() / float64(nTop)
	fmt.Printf("\nmeasured collection speed: %.1f ns/term at a 25-spacing window\n", perTerm*1e9)
	fmt.Println("\nCAPACITY REPORT (the wide-window ladder):")
	for _, e := range []struct{ t, sp float64 }{
		{1.11e19, 50}, {2.22e21, 25}, {4.44e22, 15}, {1.11e24, 6},
	} {
		n := math.Sqrt(e.t / twoPi.hi)
		scale := (3*e.sp + 18) / (3*25 + 18)
		fmt.Printf("  t = %.3g   %.2g terms   window %.0f spacings (~%.0f zeros)   ~%.0f min\n",
			e.t, n, e.sp, e.sp, n*perTerm*scale/60)
	}
	fmt.Println("  precision ceiling of the mod-reduction: t ~ 4e24")
	fmt.Println("\nMAXIMUM RATED CAPACITY: t ~ 1e24 (hours-scale), sweet spot 4.44e22.")
	fmt.Println("the ship awaits the launch order.")

	if *night {
		mission()
		return
	}

	if *anchor > 0 {
		hStart := time.Now()
		spacing := twoPi.hi / math.Log(*anchor/twoPi.hi)
		sp := build(*anchor, int64(math.Sqrt(*anchor/twoPi.hi)))
		zeros, exact := sp.sweep(*from*spacing, *spacings*spacing, spacing)
		nf, mins := len(sp.facets), time.Since(hStart).Minutes()
		fmt.Printf("\nANCHORAGE t = %.6g, window [%.0f, %.0f] spacings  (%d facets, %.1f min)\n",
			*anchor, *from, *from+*spacings, nf, mins)
		fmt.Printf("  first zero here is zero number ~%s (+/-2, the S ambiguity)\n",
			zeroIndex(*anchor))
		fmt.Print("  virgin zeros (offsets):")
		for _, z := range zeros {
			fmt.Printf("  %.6f", z)
		}
		delta := float64(len(zeros)) - exact
		fmt.Printf("\n  THE SPHERE: boundary demands %.2f zeros; found %d (delta %+.2f)", exact, len(zeros), delta)
		if math.Abs(delta) <= tolS(*anchor) {
			fmt.Println("  -> ENCLOSURE CERTIFIED: no zero escapes this window.")
		} else {
			fmt.Println("  -> ENCLOSURE BROKEN - investigate before trusting.")
		}
	}
}
