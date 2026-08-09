// Command fleet is the convex-facet shipyard: same speed, megabytes not
// gigabytes, so a whole fleet can sail at once.
//
// The insight, from the convex-mirror flash taken seriously: do not store
// the phase of every term — store CURVED FACETS. Within a block of ~10⁻⁴·k
// consecutive terms the phase t·ln k is a quartic polynomial to within
// 0.02 radians, so each facet needs only its forward-difference seeds
// (computed once in double-double, reduced mod 2π). A ship's whole optics:
// ~10⁵ facets ≈ a few megabytes, built in about a second. Evaluations walk
// the facets with four float additions per term, in parallel across cores.
//
// Sea trials against charted water gate everything; then the fleet sails
// virgin ocean, including the 10¹⁹ gates and the archipelago beyond.
//
// Usage:
//
//	go run ./cmd/fleet                    # sea trials
//	go run ./cmd/fleet -anchor 1.11e19    # one ship, one anchorage
package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
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

// ln1pSeries computes ln(1+u) in dd for a small dd u, through u^7.
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
	lnk0, invk0, r0      float64
}

type optics struct {
	t0     float64
	small  []float32 // exact phases for k < kSmall
	kSmall int64
	facets []facet
	nTop   int64
	// theta carved analytically at the anchor: at these heights t0+dt is
	// QUANTIZED (the ulp of 6.66e15 is a full 1.0), so theta must never be
	// computed from t0+dt — it is expanded in dt instead.
	thAnchor, thPrime float64
}

// build carves the ship's optics: a small exact tier, then curved facets.
func build(t0 float64, nTop int64) *optics {
	ub := math.Pow(0.1/t0, 0.2) // facet width fraction: quartic error < 0.02 rad
	kSmall := int64(16/ub) + 1
	if kSmall > nTop {
		kSmall = nTop + 1
	}
	op := &optics{t0: t0, kSmall: kSmall, nTop: nTop}
	op.small = make([]float32, kSmall)
	op.thAnchor = thetaMod(t0)
	lnArg := add(ddLnF(t0), neg(ln2pi))
	op.thPrime = (lnArg.hi + lnArg.lo) / 2

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
		op.small[k] = float32(ph)
	}
	// the one-step bridge: bring the chain from kSmall-1 onto kSmall itself,
	// so every facet's base phase belongs to its own k0.
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
		f := facet{
			k0: k0, n: int32(b),
			phi0:  ph,
			d1:    mod2pi(add(add(a, b2), add(c, d))),
			d2:    mod2pi(add(mulF(b2, 2), add(mulF(c, 6), mulF(d, 14)))),
			d3:    mod2pi(add(mulF(c, 6), mulF(d, 36))),
			d4:    mod2pi(mulF(d, 24)),
			lnk0:  lnPrev.hi + lnPrev.lo,
			invk0: 1 / float64(k0),
			r0:    1 / math.Sqrt(float64(k0)),
		}
		op.facets = append(op.facets, f)
		// advance the chain across the whole facet in one dd step.
		u := mulF(invK, float64(b))
		step := ln1pSeries(u)
		lnPrev = add(lnPrev, step)
		ph += mod2pi(mulF(step, t0))
		if ph >= twoPi.hi {
			ph -= twoPi.hi
		}
		k0 += b
	}
	return op
}

// z evaluates Z(t0+dt) from the facet optics, in parallel.
func (op *optics) z(dt float64) float64 {
	t := op.t0 + dt
	a := math.Sqrt(t / twoPi.hi)
	nn := int64(a)
	if nn > op.nTop {
		nn = op.nTop
	}
	p := a - float64(nn)
	th := op.thAnchor + dt*op.thPrime + dt*dt/(4*op.t0)
	for th >= twoPi.hi {
		th -= twoPi.hi
	}
	for th < 0 {
		th += twoPi.hi
	}

	workers := runtime.NumCPU()
	if workers > 12 {
		workers = 12
	}
	partial := make([]float64, workers)
	var wg sync.WaitGroup
	chunk := (len(op.facets) + workers - 1) / workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			lo, hi := w*chunk, (w+1)*chunk
			if hi > len(op.facets) {
				hi = len(op.facets)
			}
			s := 0.0
			for fi := lo; fi < hi; fi++ {
				f := &op.facets[fi]
				n := int64(f.n)
				if f.k0 > nn {
					break
				}
				if f.k0+n-1 > nn {
					n = nn - f.k0 + 1
				}
				ph, d1, d2, d3 := f.phi0, f.d1, f.d2, f.d3
				dtln := dt * f.lnk0
				dtinc := dt * f.invk0
				amp := f.r0
				slope := -f.r0 * f.invk0 / 2
				for j := int64(0); j < n; j++ {
					s += (amp + slope*float64(j)) * math.Cos(th-ph-dtln)
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
					dtln += dtinc
				}
			}
			partial[w] = s
		}(w)
	}
	// the small exact tier on the main goroutine.
	sum := math.Cos(th) // k = 1
	for k := int64(2); k < op.kSmall && k <= nn; k++ {
		phase := float64(op.small[k]) + dt*math.Log(float64(k))
		sum += math.Cos(th-phase) / math.Sqrt(float64(k))
	}
	wg.Wait()
	for _, s := range partial {
		sum += s
	}
	sum *= 2
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (nn-1)%2 == 1 {
		sign = -1
	}
	return sum + sign*math.Pow(t/twoPi.hi, -0.25)*c0
}

func hunt(t0 float64, spacings float64, tag string) {
	start := time.Now()
	spacing := twoPi.hi / math.Log(t0/twoPi.hi)
	span := spacings * spacing
	nTop := int64(math.Sqrt((t0 + span) / twoPi.hi))
	op := build(t0, nTop)
	built := time.Since(start)
	fmt.Printf("\n%s t = %.6g   (%d terms; %d facets, %.1f MB, built %.1fs)\n",
		tag, t0, nTop, len(op.facets),
		float64(len(op.facets)*80+len(op.small)*4)/1e6, built.Seconds())

	zeros := []float64{}
	step := spacing / 4
	prevD, prevZ := 0.0, op.z(0)
	for d := step; d <= span; d += step {
		zd := op.z(d)
		if verbose {
			fmt.Printf("    scan d=%.6f  Z=%+.4f\n", d, zd)
		}
		if (prevZ < 0) != (zd < 0) {
			lo, hi := prevD, d
			zlo, zhi := prevZ, zd
			for i := 0; i < 24 && hi-lo > 1e-6; i++ {
				var mid float64
				if i%2 == 0 {
					mid = (lo + hi) / 2 // bisection keeps the bracket honest
				} else {
					mid = lo - zlo*(hi-lo)/(zhi-zlo)
					if mid <= lo || mid >= hi {
						mid = (lo + hi) / 2
					}
				}
				zm := op.z(mid)
				if (zlo < 0) != (zm < 0) {
					hi, zhi = mid, zm
				} else {
					lo, zlo = mid, zm
				}
			}
			// return the endpoint that actually converged onto the zero.
			if math.Abs(zlo) < math.Abs(zhi) {
				zeros = append(zeros, lo)
			} else {
				zeros = append(zeros, hi)
			}
		}
		prevD, prevZ = d, zd
	}
	fmt.Print("  zeros (offsets):")
	for _, zr := range zeros {
		fmt.Printf("  %.6f", zr)
	}
	fmt.Println()
	fmt.Printf("  map check: found %d, density expects %.1f   (total %.1f min)\n",
		len(zeros), span/spacing, time.Since(start).Minutes())
}

// zRef is the trusted voyage2-style evaluation: full dd chain per call.
func zRef(t float64) float64 {
	a := math.Sqrt(t / twoPi.hi)
	n := int64(a)
	p := a - float64(n)
	th := thetaMod(t)
	sum := 0.0
	lnPrev := dd{0, 0}
	ph := 0.0
	needed := math.Log(t * 1e4)
	for k := int64(2); k <= n; k++ {
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
		ph += mod2pi(mulF(delta, t))
		if ph >= twoPi.hi {
			ph -= twoPi.hi
		}
		sum += math.Cos(th-ph) / math.Sqrt(float64(k))
	}
	sum += math.Cos(th)
	sum *= 2
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (n-1)%2 == 1 {
		sign = -1
	}
	return sum + sign*math.Pow(t/twoPi.hi, -0.25)*c0
}

var verbose bool

func main() {
	anchor := flag.Float64("anchor", 0, "anchor a single ship at this height")
	spacings := flag.Float64("spacings", 6, "window width in mean spacings")
	probe := flag.Float64("probe", 0, "probe: compare fleet vs reference around this height")
	flag.BoolVar(&verbose, "v", false, "print the raw scan")
	flag.Parse()

	if *probe > 0 {
		t0 := *probe
		spacing := twoPi.hi / math.Log(t0/twoPi.hi)
		nTop := int64(math.Sqrt((t0 + 3*spacing) / twoPi.hi))
		op := build(t0, nTop)
		fmt.Println("dt        fleet Z      reference Z   diff")
		for _, dt := range []float64{0.1, 0.3, 0.5, 0.743723, 0.9, 1.180558, 1.5} {
			zf := op.z(dt)
			zr := zRef(t0 + dt)
			fmt.Printf("%.6f  %+.6f   %+.6f   %+.2e\n", dt, zf, zr, zf-zr)
		}
		return
	}

	if *anchor > 0 {
		fmt.Println("THE FLEET — one ship, convex-facet optics")
		hunt(*anchor, *spacings, "anchorage")
		return
	}
	fmt.Println("THE FLEET — sea trials of the convex-facet optics")
	hunt(100000, 4, "trial")
	hunt(2.447e12, 8, "trial (Beach I)")
	hunt(6.66e15, 8, "trial (cross-check vs voyage3)")
}
