// Command starship is the spacecraft: the flagship's certified hull plus
// a fourth tier where deep blocks are not rowed term by term — they are
// slung once through the Fresnel gearbox and enter the light bucket as
// single super-terms.
//
// The physics that permits it, checked twice: within a hunting window the
// internal shape of a block changes by less than 1e-8 radians — ALL the
// t-dependence rides on the carrier e^{-i t ln(kc)}. One fold serves the
// whole window. The block's amplitude is flat to ~1e-8 (blocks are slivers
// of k), its quadratic phase model is exact to 0.003 rad by the choice of
// block length L = (0.009/t)^(1/3) * k, and the carrier phase is kept on
// the same double-double chain that carves the facets.
//
// The sky opens at t ~ 1e21: below that, blocks are shorter than the fold's
// own overhead and the ship simply rows (it degenerates exactly into the
// flagship). The flight test is an A/B duel at t = 2.22e21: pure hull vs
// folded hull over the same virgin window — same zeros, same sphere, or
// the spacecraft stays in the hangar.
//
// Usage:
//
//	go run ./cmd/starship            # certification gates (pure water)
//	go run ./cmd/starship -tunnel    # wind tunnel: fold error/cost vs edge
//	go run ./cmd/starship -flight    # the A/B flight test at t = 2.22e21
package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ---------- double-double core (copied verbatim from the flagship) ----------

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
	piDD  = parse("3.14159265358979323846264338327950288419716939937511")
)

var twoPi26 = dd{twoPi.hi * 67108864, twoPi.lo * 67108864}

func mod2pi(x dd) float64 {
	n := math.Floor(x.hi / twoPi.hi)
	n1 := math.Floor(n / 67108864)
	n2 := n - n1*67108864
	x = add(x, mulF(twoPi26, -n1))
	x = add(x, mulF(twoPi, -n2))
	r := x.hi + x.lo
	// beyond the certified ceiling (t ~ 4e24) the counters overflow 53
	// bits and the residue can be far off: math.Mod as the safety net so
	// the cleanup can NEVER spin forever. Within the certified regime
	// this is identical to the one-or-two subtractions it replaces.
	if r >= twoPi.hi || r < 0 {
		r = math.Mod(r, twoPi.hi)
		if r < 0 {
			r += twoPi.hi
		}
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

// ---------- the colossal absorber (F109): FFT deposit, flat cost ----------

var colosal bool

// foldWOverride, when positive, pins the fold tier's segment count. The
// quad approximation realizes its (in-budget) error differently under
// different block boundaries, so certification against a pre-surgery
// photograph needs the pre-surgery blocking: -foldw 1.
var foldWOverride int

const msp = 12

// fft: in-place iterative radix-2, forward.
func fft(re, im []float64) {
	n := len(re)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit
		if i < j {
			re[i], re[j] = re[j], re[i]
			im[i], im[j] = im[j], im[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		ang := -2 * math.Pi / float64(length)
		wr, wi := math.Cos(ang), math.Sin(ang)
		for i := 0; i < n; i += length {
			cwr, cwi := 1.0, 0.0
			for j := 0; j < length/2; j++ {
				ur, ui := re[i+j], im[i+j]
				vr := re[i+j+length/2]*cwr - im[i+j+length/2]*cwi
				vi := re[i+j+length/2]*cwi + im[i+j+length/2]*cwr
				re[i+j], im[i+j] = ur+vr, ui+vi
				re[i+j+length/2], im[i+j+length/2] = ur-vr, ui-vi
				cwr, cwi = cwr*wr-cwi*wi, cwr*wi+cwi*wr
			}
		}
	}
}

func engineName() string {
	if colosal {
		return "colosal"
	}
	return "classic"
}

// ---------- the Fresnel gearbox (from cmd/fingers, edge adjustable) ----------

var edge float64 = 48

func fresnel(x float64) (float64, float64) {
	negx := x < 0
	if negx {
		x = -x
	}
	var c, s float64
	if x <= 3.2 {
		t := math.Pi / 2 * x * x
		baseC := 1.0
		baseS := t
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
	if negx {
		return -c, -s
	}
	return c, s
}

func frac(x float64) float64 {
	return x - math.Floor(x)
}

func quadDirect(a, b float64, L int64) (float64, float64) {
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

// quad computes S(a,b,L) = sum_{j<L} e^{2 pi i (a j + b j^2)}.
func quad(a, b float64, L int64, depth int) (float64, float64) {
	a, b = frac(a), frac(b)
	if b > 0.5 {
		b -= 1
	}
	if b < 0 {
		r, i := quad(frac(-a), -b, L, depth)
		return r, -i
	}
	if b > 0.25 {
		return quad(frac(a+0.5), b-0.5, L, depth)
	}
	if L <= 256 || depth > 24 {
		return quadDirect(a, b, L)
	}
	if b*float64(L) < 1e-9 {
		if a < 1e-12 || a > 1-1e-12 {
			return float64(L), 0
		}
		ph := math.Pi * a * float64(L-1)
		amp := math.Sin(math.Pi*a*float64(L)) / math.Sin(math.Pi*a)
		return amp * math.Cos(ph), amp * math.Sin(ph)
	}
	x0, x1 := -0.5, float64(L)-0.5
	d0, d1 := a+2*b*x0, a+2*b*x1
	mLoAll := int64(math.Ceil(d0 - edge))
	mHiAll := int64(math.Floor(d1 + edge))
	mLoIn := int64(math.Ceil(d0 + edge))
	mHiIn := int64(math.Floor(d1 - edge))
	inv2b := 1 / (2 * b)
	sqb2 := 2 * math.Sqrt(b)
	var sr, si float64
	for m := mLoAll; m <= mHiAll; m++ {
		if mHiIn >= mLoIn && m >= mLoIn && m <= mHiIn {
			continue
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

// ---------- the hull (flagship machinery + the fold tier) ----------

type facet struct {
	k0                   int64
	n                    int32
	phi0, d1, d2, d3, d4 float64
	cA, cB, cC, cD       dd
	lnk0, invk0, r0      float64
}

type ship struct {
	t0                float64
	small             []float32
	kSmall            int64
	facets            []facet
	nTop              int64
	thAnchor, thPrime float64
	h, dt0            float64
	fr, fi            []float64
	c0term            float64
	// the fold tier: blocks from foldFrom to nTop are slung, not rowed.
	foldFrom int64
	c        float64 // block length fraction: L = c * k
	lnB      dd      // chain handoff: ln(foldFrom)
	phB      float64 // chain handoff: t0*ln(foldFrom) mod 2pi
	nBlocks  int64
	foldSecs float64
}

// blockFrac is the quadratic-model tolerance: L = (0.009/t)^(1/3) * k
// keeps the neglected cubic phase under 0.003 rad at the block's end.
//
// gruesoFactor is the captain's short-number equivalence (F136): scale
// the block length by N for a COARSE first pass - N^3 times the phase
// tolerance (N=3: ~0.08 rad, Z to a few percent), N times fewer block
// evaluations. Reconnaissance sails fast and approximate; only where
// something shines does the fleet DECOMPRESS: re-sail that little
// window with the full, strict numbers.
var gruesoFactor float64 = 1

func blockFrac(t0 float64) float64 {
	return gruesoFactor * math.Cbrt(9.0e-3/t0)
}

func build(t0 float64, nTop, foldFrom int64) *ship {
	ub := math.Pow(0.1/t0, 0.2)
	kSmall := int64(16/ub) + 1
	if kSmall > nTop {
		kSmall = nTop + 1
	}
	sp := &ship{t0: t0, kSmall: kSmall, nTop: nTop, foldFrom: foldFrom}
	sp.c = blockFrac(t0)
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
	fTop := nTop
	if foldFrom <= nTop {
		fTop = foldFrom - 1
	}
	for k0 := kSmall; k0 <= fTop; {
		b := int64(ub * float64(k0))
		if b < 16 {
			b = 16
		}
		if k0+b-1 > fTop {
			b = fTop - k0 + 1
		}
		invK := ddInv(float64(k0))
		a := mulF(invK, t0)
		k2 := mul(invK, invK)
		b2 := mulF(k2, -t0/2)
		c := mulF(mul(k2, invK), t0/3)
		d := mulF(mul(k2, k2), -t0/4)
		sp.facets = append(sp.facets, facet{
			k0: k0, n: int32(b),
			phi0: ph,
			d1:   mod2pi(add(add(a, b2), add(c, d))),
			d2:   mod2pi(add(mulF(b2, 2), add(mulF(c, 6), mulF(d, 14)))),
			d3:   mod2pi(add(mulF(c, 6), mulF(d, 36))),
			d4:   mod2pi(mulF(d, 24)),
			cA:   a, cB: b2, cC: c, cD: d,
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
	// chain handoff to the fold tier: lnPrev = ln(foldFrom),
	// ph = t0*ln(foldFrom) mod 2pi — seamless, no bridge, no gap.
	sp.lnB, sp.phB = lnPrev, ph
	return sp
}

var refreshEvery int64 = 4096

func wrap2pi(x float64) float64 {
	x = math.Mod(x, twoPi.hi)
	if x < 0 {
		x += twoPi.hi
	}
	return x
}

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

// ---------- the DeLorean's memory: checkpointed collection ----------
// The light bucket is tiny (a few KB), so the whole voyage state can be
// photographed every few seconds of work: the facet tier saves after each
// of 64 shifts, the fold tier saves its position k plus the exact dd
// chain. Stop the ship anywhere; it resumes from the last photograph.

var memory bool

type ckptFacets struct {
	T0, Span float64
	S, Shift int
	Engine   string
	Fr, Fi   []float64
}

type ckptFold struct {
	T0, Span   float64
	S          int
	Engine     string
	Edge       float64
	K0         int64
	LnHi, LnLo float64
	Ph         float64
	NBlocks    int64
	Fr, Fi     []float64
}

// engMatch: old photographs carry no Engine field and were classic.
func engMatch(e string) bool {
	if e == "" {
		e = "classic"
	}
	return e == engineName()
}

func ckptPath(kind string, t0 float64) string {
	// coarse runs photograph apart: their blocking must never poison an
	// exact resume, nor overwrite an exact tile (F136).
	if gruesoFactor > 1 {
		return fmt.Sprintf("ckpt/%s-%.6g-grueso.gob", kind, t0)
	}
	return fmt.Sprintf("ckpt/%s-%.6g.gob", kind, t0)
}

func saveCkpt(path string, v interface{}) {
	os.MkdirAll(filepath.Dir(path), 0755)
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return
	}
	if gob.NewEncoder(f).Encode(v) != nil {
		f.Close()
		return
	}
	f.Sync()
	f.Close()
	os.Remove(path)
	os.Rename(tmp, path)
}

func loadCkpt(path string, v interface{}) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(v) == nil
}

func (sp *ship) collect(dtStart, span float64) {
	lnTop := math.Log(float64(sp.nTop))
	sp.h = 2 * math.Pi / lnTop / 3
	guard := 8 * sp.h
	sp.dt0 = dtStart - guard
	S := int((span+2*guard)/sp.h) + 2
	workers := runtime.NumCPU() - 1
	if workers < 1 {
		workers = 1
	}

	// colossal geometry (F109): fine grid, Gaussian kernel, one final FFT.
	M := 1
	for M < 4*S {
		M <<= 1
	}
	MG := M + 2*msp + 1
	dxg := 2 * math.Pi / float64(M)
	taug := 4 * math.Pi * float64(msp) / (float64(M) * float64(M) * 3)
	c1g := 1 / (4 * taug)
	c2g := dxg / (2 * taug)
	e3g := make([]float64, 2*msp+1)
	for l := -msp; l <= msp; l++ {
		e3g[l+msp] = math.Exp(-float64(l*l) * dxg * dxg / (4 * taug))
	}
	gridLen := S
	if colosal {
		gridLen = MG
	}

	// master facet grid, possibly restored from the ship's memory.
	frM := make([]float64, gridLen)
	fiM := make([]float64, gridLen)
	const nShifts = 64
	startShift := 0
	if memory {
		var ck ckptFacets
		if loadCkpt(ckptPath("facets", sp.t0), &ck) &&
			ck.T0 == sp.t0 && ck.Span == span && ck.S == gridLen && engMatch(ck.Engine) {
			copy(frM, ck.Fr)
			copy(fiM, ck.Fi)
			startShift = ck.Shift + 1
			fmt.Printf("  memory: facet tier resumes at shift %d/%d\n", startShift, nShifts)
		}
	}

	// THE FOLD TIER, split in halves (the captain's bisection applied to
	// the ENGINE itself): the fold used to row with one oarsman while the
	// rest of the crew idled after the facet shifts. Now the k-range is
	// cut into geometric segments - halves of halves - each with its own
	// dd chain, its own grid, its own memory photograph, folding in
	// parallel. The hours-long fold tail becomes minutes, on BOTH engines.
	frF := make([]float64, gridLen)
	fiF := make([]float64, gridLen)
	foldDone := make(chan struct{})
	go func() {
		defer close(foldDone)
		sp.nBlocks = 0
		if sp.foldFrom > sp.nTop {
			return
		}
		t0 := sp.t0
		foldStart := time.Now()
		W := workers
		if foldWOverride > 0 {
			W = foldWOverride
		}
		if W > 8 {
			W = 8
		}
		if W < 1 {
			W = 1
		}
		kb := make([]int64, W+1)
		kb[0] = sp.foldFrom
		ratio := math.Pow(float64(sp.nTop)/float64(sp.foldFrom), 1.0/float64(W))
		for i := 1; i < W; i++ {
			kb[i] = int64(float64(kb[i-1]) * ratio)
			if kb[i] <= kb[i-1] {
				kb[i] = kb[i-1] + 1
			}
		}
		kb[W] = sp.nTop + 1
		frSeg := make([][]float64, W)
		fiSeg := make([][]float64, W)
		counts := make([]int64, W)
		var wgF sync.WaitGroup
		for i := 0; i < W; i++ {
			frSeg[i] = make([]float64, gridLen)
			fiSeg[i] = make([]float64, gridLen)
			wgF.Add(1)
			go func(i int) {
				defer wgF.Done()
				segLo, segHi := kb[i], kb[i+1]-1
				if segHi > sp.nTop {
					segHi = sp.nTop
				}
				var lnPrev dd
				var ph float64
				if i == 0 {
					lnPrev, ph = sp.lnB, sp.phB
				} else {
					lnPrev = ddLnF(float64(segLo))
					ph = mod2pi(mulF(lnPrev, t0))
				}
				fr := frSeg[i]
				fi := fiSeg[i]
				k0 := segLo
				ckName := ckptPath(fmt.Sprintf("fold%d", i), sp.t0)
				if memory {
					var ck ckptFold
					if loadCkpt(ckName, &ck) &&
						ck.T0 == sp.t0 && ck.Span == span && ck.S == gridLen &&
						ck.Edge == edge && engMatch(ck.Engine) &&
						ck.K0 >= segLo && ck.K0 <= segHi+1 {
						copy(fr, ck.Fr)
						copy(fi, ck.Fi)
						k0 = ck.K0
						lnPrev = dd{ck.LnHi, ck.LnLo}
						ph = ck.Ph
						counts[i] = ck.NBlocks
						fmt.Printf("  memory: fold segment %d resumes at k = %.4g\n", i, float64(k0))
					}
				}
				save := func() {
					saveCkpt(ckName, ckptFold{
						T0: sp.t0, Span: span, S: gridLen, Engine: engineName(), Edge: edge,
						K0: k0, LnHi: lnPrev.hi, LnLo: lnPrev.lo, Ph: ph,
						NBlocks: counts[i], Fr: fr, Fi: fi,
					})
				}
				var sinceSave int64
				for k0 <= segHi {
					k0f := float64(k0)
					L := int64(sp.c * k0f)
					if L < 64 {
						L = 64
					}
					if k0+L-1 > segHi {
						L = segHi - k0 + 1
					}
					invK := ddInv(k0f)
					a := mod2pi(mulF(invK, t0)) / twoPi.hi
					b := -t0 / (4 * math.Pi * k0f * k0f)
					sr, si := quad(a, b, L, 0)
					si = -si
					Lc := float64(L) / 2
					lnR := ln1pSeries(mulF(invK, Lc))
					psi := mod2pi(mulF(lnR, t0))
					cp, sq := math.Cos(psi), math.Sin(psi)
					rc := 1 / math.Sqrt(k0f+Lc)
					gr := rc * (cp*sr - sq*si)
					gi := rc * (cp*si + sq*sr)
					phc := ph + psi
					if phc >= twoPi.hi {
						phc -= twoPi.hi
					}
					lnkc := lnPrev.hi + lnPrev.lo + lnR.hi + lnR.lo
					ang := phc + sp.dt0*lnkc
					ca, sa := math.Cos(-ang), math.Sin(-ang)
					zr := gr*ca - gi*sa
					zi := gr*sa + gi*ca
					if colosal {
						xig := sp.h * lnkc
						m0 := int(xig/dxg + 0.5)
						xp := xig - float64(m0)*dxg
						fk := math.Exp(-xp*xp*c1g - float64(msp)*xp*c2g)
						es := math.Exp(xp * c2g)
						for l := 0; l <= 2*msp; l++ {
							w := fk * e3g[l]
							fr[m0+l] += w * zr
							fi[m0+l] += w * zi
							fk *= es
						}
					} else {
						wr, wi := math.Cos(-sp.h*lnkc), math.Sin(-sp.h*lnkc)
						for s := 0; s < S; s++ {
							fr[s] += zr
							fi[s] += zi
							zr, zi = zr*wr-zi*wi, zr*wi+zi*wr
						}
					}
					st := ln1pSeries(mulF(invK, float64(L)))
					lnPrev = add(lnPrev, st)
					ph += mod2pi(mulF(st, t0))
					if ph >= twoPi.hi {
						ph -= twoPi.hi
					}
					k0 += L
					counts[i]++
					sinceSave++
					if memory && sinceSave >= 1000000 {
						save()
						sinceSave = 0
					}
				}
				if memory {
					save()
				}
			}(i)
		}
		wgF.Wait()
		for i := 0; i < W; i++ {
			for s := 0; s < gridLen; s++ {
				frF[s] += frSeg[i][s]
				fiF[s] += fiSeg[i][s]
			}
			sp.nBlocks += counts[i]
		}
		sp.foldSecs = time.Since(foldStart).Seconds()
	}()

	// facet tier in checkpointed shifts.
	nf := len(sp.facets)
	frs := make([][]float64, workers)
	fis := make([][]float64, workers)
	for w := 0; w < workers; w++ {
		frs[w] = make([]float64, gridLen)
		fis[w] = make([]float64, gridLen)
	}
	runShift := func(sLo, sHi int) {
		var wg sync.WaitGroup
		chunk := (sHi - sLo + workers - 1) / workers
		if chunk < 1 {
			chunk = 1
		}
		for w := 0; w < workers; w++ {
			lo, hi := sLo+w*chunk, sLo+(w+1)*chunk
			if hi > sHi {
				hi = sHi
			}
			if lo >= hi {
				continue
			}
			wg.Add(1)
			go func(lo, hi, w int) {
				defer wg.Done()
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
					if colosal {
						xig := sp.h * lnLin
						m0 := int(xig/dxg + 0.5)
						xp := xig - float64(m0)*dxg
						fk := math.Exp(-xp*xp*c1g - float64(msp)*xp*c2g)
						es := math.Exp(xp * c2g)
						for l := 0; l <= 2*msp; l++ {
							w := fk * e3g[l]
							fr[m0+l] += w * cr
							fi[m0+l] += w * ci
							fk *= es
						}
					} else {
						wr, wi := wr0, wi0
						for s := 0; s < S; s++ {
							fr[s] += cr
							fi[s] += ci
							cr, ci = cr*wr-ci*wi, cr*wi+ci*wr
						}
					}
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
			}(lo, hi, w)
		}
		wg.Wait()
	}
	for shift := startShift; shift < nShifts; shift++ {
		sLo, sHi := shift*nf/nShifts, (shift+1)*nf/nShifts
		for w := 0; w < workers; w++ {
			for s := 0; s < gridLen; s++ {
				frs[w][s] = 0
				fis[w][s] = 0
			}
		}
		runShift(sLo, sHi)
		for w := 0; w < workers; w++ {
			for s := 0; s < gridLen; s++ {
				frM[s] += frs[w][s]
				fiM[s] += fis[w][s]
			}
		}
		if memory {
			saveCkpt(ckptPath("facets", sp.t0), ckptFacets{
				T0: sp.t0, Span: span, S: gridLen, Engine: engineName(),
				Shift: shift, Fr: frM, Fi: fiM,
			})
		}
	}

	// small tier + k=1 on the main goroutine (fast, always fresh).
	fr0 := make([]float64, gridLen)
	fi0 := make([]float64, gridLen)
	deposit0 := func(lnk, cr, ci float64) {
		if colosal {
			xig := sp.h * lnk
			m0 := int(xig/dxg + 0.5)
			xp := xig - float64(m0)*dxg
			fk := math.Exp(-xp*xp*c1g - float64(msp)*xp*c2g)
			es := math.Exp(xp * c2g)
			for l := 0; l <= 2*msp; l++ {
				w := fk * e3g[l]
				fr0[m0+l] += w * cr
				fi0[m0+l] += w * ci
				fk *= es
			}
		} else {
			wr, wi := math.Cos(-sp.h*lnk), math.Sin(-sp.h*lnk)
			for s := 0; s < S; s++ {
				fr0[s] += cr
				fi0[s] += ci
				cr, ci = cr*wr-ci*wi, cr*wi+ci*wr
			}
		}
	}
	deposit0(0, 1, 0) // k = 1: frequency zero, amplitude one
	for k := int64(2); k < sp.kSmall && k <= sp.nTop; k++ {
		lnk := math.Log(float64(k))
		amp := 1 / math.Sqrt(float64(k))
		ang := float64(sp.small[k]) + sp.dt0*lnk
		deposit0(lnk, amp*math.Cos(-ang), -amp*math.Sin(ang))
	}

	<-foldDone
	if colosal {
		// the colossal re-propulsion: fold the guards, one FFT, deconvolve.
		rr := make([]float64, M)
		ri := make([]float64, M)
		for m := 0; m < MG; m++ {
			idx := ((m-msp)%M + M) % M
			rr[idx] += frM[m] + fr0[m] + frF[m]
			ri[idx] += fiM[m] + fi0[m] + fiF[m]
		}
		fft(rr, ri)
		sp.fr = make([]float64, S)
		sp.fi = make([]float64, S)
		normg := dxg / (2 * math.Sqrt(math.Pi*taug))
		for s := 0; s < S; s++ {
			c := normg * math.Exp(float64(s)*float64(s)*taug)
			sp.fr[s] = rr[s] * c
			sp.fi[s] = ri[s] * c
		}
	} else {
		sp.fr, sp.fi = fr0, fi0
		for s := 0; s < S; s++ {
			sp.fr[s] += frM[s] + frF[s]
			sp.fi[s] += fiM[s] + fiF[s]
		}
	}
	// the voyage is complete: burn the photographs so no stale memory can
	// ever contaminate a future anchorage.
	if memory {
		os.Remove(ckptPath("facets", sp.t0))
		os.Remove(ckptPath("fold", sp.t0))
		for i := 0; i < 16; i++ {
			os.Remove(ckptPath(fmt.Sprintf("fold%d", i), sp.t0))
		}
	}
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

// lightArchive is the notebook's compression made mathematics: the whole
// collected window — billions of rowed and folded terms — fits in a few
// kilobytes of band-limited light. Archived, it grants infinite zoom
// forever: re-scan at any resolution without ever sailing again.
type lightArchive struct {
	T0, Span, H, Dt0, C0, ThA, ThP float64
	NTop                           int64
	Fr, Fi                         []float64
	// Bow is the compass's window shift (0 on tiles from before the
	// bow existed - gob leaves missing fields zero, so old photographs
	// stay correct). Every reader of the tile must scan [Bow, Bow+Span].
	Bow float64
}

func archivePath(t0 float64) string {
	if gruesoFactor > 1 {
		return fmt.Sprintf("luz/luz-%.6g-grueso.gob", t0)
	}
	return fmt.Sprintf("luz/luz-%.6g.gob", t0)
}

// rescan hunts the already-recorded light: the zoom, at any step, for free.
func (sp *ship) rescan(dtStart, span, spacing float64) ([]float64, float64) {
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
				z := lo
				if math.Abs(zhi) < math.Abs(zlo) {
					z = hi
				}
				// the gravitational lens (the supermassive magnifier):
				// near a zero the light bends linearly, so each Newton
				// step through the lens SQUARES the precision - the
				// closer you approach, the more it amplifies.
				for it := 0; it < 2; it++ {
					hd := spacing * 1e-4
					d := (sp.zAt(z+hd) - sp.zAt(z-hd)) / (2 * hd)
					if d != 0 {
						z -= sp.zAt(z) / d
					}
				}
				found = append(found, z)
			}
			prevD, prevZ = d, zd
		}
		return found
	}
	zeros := scan(spacing / 60)
	exact := ((end-dtStart)*sp.thPrime + (end*end-dtStart*dtStart)/(4*sp.t0)) / math.Pi
	if float64(len(zeros)) < exact-0.5 {
		zeros = scan(spacing / 600)
	}
	// the turbulence glass (F105 orchestrated with F101): hidden pairs
	// live in turbulent water, never in crystal. A suspiciously narrow
	// pass raises the glass as a precaution, before any count deficit.
	if len(zeros) >= 2 {
		minG := math.Inf(1)
		for i := 0; i+1 < len(zeros); i++ {
			if g := zeros[i+1] - zeros[i]; g < minG {
				minG = g
			}
		}
		if minG < 0.25*spacing {
			zeros = scan(spacing / 600)
		}
	}
	return zeros, exact
}

func (sp *ship) sweep(dtStart, span, spacing float64) ([]float64, float64) {
	sp.collect(dtStart, span)
	return sp.rescan(dtStart, span, spacing)
}

// tolS is the barometric tolerance (F102): the sea's number variance
// SATURATES (measured ~0.5 at gamma~1e4, Berry saturation), growing only
// as ln ln t. The sphere's honest gauge is 2.5 sigma of that pressure,
// calibrated on cmd/barometro's measurement.
func tolS(t float64) float64 {
	return 2.5 * math.Sqrt((math.Log(math.Log(t))+2.72)/(math.Pi*math.Pi))
}

// sPred forecasts the predictable part of the tide S(t0+dt) from the
// loudest prime voices — the pacemaker of Finding 100: the rhythm of the
// rhythm is the primes, so a large share of the sphere's delta can be
// predicted BEFORE sailing. S(t) ~ -(1/pi) sum sin(t ln p)/sqrt(p).
func sPred(t0, dt float64) float64 {
	s := 0.0
	for _, p := range []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37,
		41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101} {
		ph := mod2pi(mulF(ddLnF(p), t0)) + dt*math.Log(p)
		s -= math.Sin(ph) / math.Sqrt(p)
	}
	return s / math.Pi
}

// ---------- EL LIBRO DE COORDENADAS (the captain's split number) ----------
// One float64 cannot write "10^21 plus a third of a spacing" - but TWO
// variables can, exactly: the anchor T and the window offset U. Double
// the size, double the truth (the cure, made permanent, for F122's
// rounded-coordinate lesson). Every treasure the laboratory finds is a
// row in a matrix that GROWS with time - the coordinate book - and each
// row carries a REDUCED string form: a short street name to talk with,
// while the matrix keeps the full-precision pair.

type coordRow struct {
	T, U       float64 // the split number: anchor + offset (double size)
	Kind, Addr string  // what it is; base-twelve address when the compass sang it
}

const coordBook = "luz/coordenadas.gob"

// coordString is the reduced form: short enough to speak, unique enough
// to find the row again.
func coordString(c coordRow) string {
	return fmt.Sprintf("%s:%.6g@u%+.3f", c.Kind, c.T, c.U)
}

// saveCoord appends one row to the growing matrix.
func saveCoord(kind string, t, u float64, addr string) {
	var book []coordRow
	loadCkpt(coordBook, &book)
	book = append(book, coordRow{T: t, U: u, Kind: kind, Addr: addr})
	saveCkpt(coordBook, book)
}

// zVerdad evaluates Z(t0+u) by the direct term-by-term dd chain - the
// supreme judge's arithmetic, sharing nothing with the engines but the
// dd core. Streams the chain (no table): any depth, O(nTop) per point.
func zVerdad(t0, u float64) float64 {
	nTop := int64(math.Sqrt(t0 / twoPi.hi))
	lnArg := add(ddLnF(t0), neg(ln2pi))
	thP := (lnArg.hi + lnArg.lo) / 2
	th := mod2pi(mulF(lnArg, t0/2)) - mod2pi(dd{t0 / 2, 0}) - math.Pi/8 + 1/(48*t0)
	th += thP*u + u*u/(4*t0)
	th = wrap2pi(th)
	lnPrev := dd{0, 0}
	ph := 0.0
	needed := math.Log(t0 * 1e4)
	var z float64
	for k := int64(1); k <= nTop; k++ {
		if k > 1 {
			var delta dd
			if k <= 64 {
				lk := ddLnF(float64(k))
				delta = add(lk, neg(lnPrev))
				lnPrev = lk
			} else {
				uu := ddInv(float64(k - 1))
				m := int(needed/math.Log(float64(k-1))) + 2
				if m > 16 {
					m = 16
				}
				s := mulF(ddInv(float64(m)), 1)
				for j := m - 1; j >= 1; j-- {
					s = mul(neg(s), uu)
					s = add(s, mulF(ddInv(float64(j)), 1))
				}
				delta = mul(s, uu)
				lnPrev = add(lnPrev, delta)
			}
			ph += mod2pi(mulF(delta, t0))
			if ph >= twoPi.hi {
				ph -= twoPi.hi
			}
		}
		lnk := lnPrev.hi + lnPrev.lo
		z += math.Cos(th-ph-u*lnk) / math.Sqrt(float64(k))
	}
	return 2 * z
}

// compassTwelve is the captain's fractal compass: cut the shout's window
// into twelve sectors, take the treasure's twelfth, cut THAT into twelve
// - infinite redundancy - until the sector is finer than a thousandth of
// a spacing. Returns the treasure's address in base twelve and the BOW:
// the window-start offset that seats the treasure mid-frame, so it can
// never be clipped at an edge (F119 addendum 3 folded into the aim).
// calmMode hunts the stillest point (land); otherwise the wildest (storm).
func compassTwelve(t0 float64, calmMode bool) (string, float64) {
	spc := twoPi.hi / math.Log(t0/twoPi.hi)
	span := 5 * spc
	lo, hi := 0.0, span
	addr := ""
	for hi-lo > spc/1000 {
		bestJ, bestV := 0, math.Inf(-1)
		if calmMode {
			bestV = math.Inf(1)
		}
		w := (hi - lo) / 12
		for j := 0; j < 12; j++ {
			u := lo + w*(float64(j)+0.5)
			v := math.Abs(sPred(t0, u))
			if (!calmMode && v > bestV) || (calmMode && v < bestV) {
				bestV, bestJ = v, j
			}
		}
		if addr != "" {
			addr += "."
		}
		addr += fmt.Sprintf("%d", bestJ+1)
		lo, hi = lo+w*float64(bestJ), lo+w*float64(bestJ+1)
	}
	return addr, (lo+hi)/2 - span/2
}

// launch builds a ship at t0; fold=true opens the fourth tier wherever the
// sky is deep enough (blocks of at least 256 terms), fold=false rows pure.
func launch(t0 float64, fold bool) *ship {
	nTop := int64(math.Sqrt(t0 / twoPi.hi))
	foldFrom := nTop + 1
	if fold {
		kB := int64(math.Ceil(256 / blockFrac(t0)))
		if kB <= nTop {
			foldFrom = kB
		}
	}
	return build(t0, nTop, foldFrom)
}

// hunt sails a window of the given width at t0. bow shifts the window
// start (the compass's centered aim): at these heights a sub-spacing
// shift is far below float64's resolution in t itself, so the bow lives
// where the fine precision lives - in the window offset.
func hunt(t0, spacings float64, fold bool, bow float64) ([]float64, float64, *ship, float64) {
	start := time.Now()
	spacing := twoPi.hi / math.Log(t0/twoPi.hi)
	sp := launch(t0, fold)
	zeros, exact := sp.sweep(bow, spacings*spacing, spacing)
	return zeros, exact, sp, time.Since(start).Minutes()
}

func gate(name string, t0, spacings float64, expect []float64, tol float64) bool {
	zeros, exact, sp, mins := hunt(t0, spacings, true, 0)
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
	fmt.Printf("  gate %-22s found %d/%d  worst dev %.5f (tol %.3f)  sphere %.2f  %d facets  %d blocks  %.1f min  [%s]\n",
		name, len(zeros), len(expect), worst, tol, exact, len(sp.facets), sp.nBlocks, mins, verdict)
	return ok
}

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

// tunnel measures the fold's error and cost against direct summation on
// blocks with the exact parameters the flight will meet at t = 2.22e21.
func tunnel() {
	fmt.Println("\nWIND TUNNEL - fold error and cost vs Fresnel edge width")
	t0 := 2.22e21
	nTop := int64(math.Sqrt(t0 / twoPi.hi))
	c := blockFrac(t0)
	kB := int64(math.Ceil(256 / c))
	fmt.Printf("  block ensemble at t=%.3g: k in [%.3g, %.3g], L = %.0f..%.0f\n",
		t0, float64(kB), float64(nTop), 256.0, c*float64(nTop))

	rng := rand.New(rand.NewSource(2026))
	const trials = 40
	type blk struct {
		a, b float64
		L    int64
	}
	blocks := make([]blk, trials)
	for i := range blocks {
		u := rng.Float64()
		k0 := float64(kB) * math.Pow(float64(nTop)/float64(kB), u)
		blocks[i] = blk{
			a: rng.Float64(),
			b: -t0 / (4 * math.Pi * k0 * k0),
			L: int64(c * k0),
		}
	}
	fmt.Println("    edge   worst abs err   mean abs err   mean cost/block")
	for _, e := range []float64{350, 160, 80, 48, 24, 12} {
		edge = e
		worst, mean := 0.0, 0.0
		st := time.Now()
		for _, bl := range blocks {
			fr, fi := quad(bl.a, bl.b, bl.L, 0)
			dr, di := quadDirect(frac(bl.a), bl.b-math.Floor(bl.b), bl.L)
			err := math.Hypot(fr-dr, fi-di)
			mean += err
			if err > worst {
				worst = err
			}
		}
		el := time.Since(st).Seconds() / trials
		mean /= trials
		fmt.Printf("    %4.0f   %13.2e   %12.2e   %10.1f us\n", e, worst, mean, el*1e6)
	}
	nb := math.Log(float64(nTop)/float64(kB)) / c
	rc := 1 / math.Sqrt(float64(kB))
	fmt.Printf("\n  flight will fold %.3g blocks; aggregate error ~ sqrt(blocks) * rc * mean\n", nb)
	fmt.Printf("  (rc = %.1e, sqrt(blocks) = %.0f -> multiply mean abs err by %.1e)\n",
		rc, math.Sqrt(nb), rc*math.Sqrt(nb))
	fmt.Println("  pick the narrowest edge whose projected aggregate stays under ~1e-3.")
}

// flight is the A/B duel at t = 2.22e21: pure hull vs folded hull over the
// same virgin window. Same zeros, same sphere, or the craft stays grounded.
func flight() {
	t0 := 2.22e21
	spacings := 5.0
	spacing := twoPi.hi / math.Log(t0/twoPi.hi)
	span := spacings * spacing
	logLine("")
	logLine("## Vuelo de prueba de la nave espacial — " + time.Now().Format("2006-01-02 15:04"))
	logLine(fmt.Sprintf("- anclaje t = %.6g, ventana %.0f espaciamientos, borde Fresnel %.0f", t0, spacings, edge))

	fmt.Printf("\nFLIGHT TEST at t = %.3g (window %.0f spacings, edge %.0f)\n", t0, spacings, edge)
	fmt.Println("  leg 1: the pure hull rows every term...")
	st := time.Now()
	spPure := launch(t0, false)
	zP, exP := spPure.sweep(0, span, spacing)
	tPure := time.Since(st).Minutes()
	fmt.Printf("    pure:   %d zeros, sphere %.2f, %.1f min\n", len(zP), exP, tPure)

	fmt.Println("  leg 2: the spacecraft slings the deep blocks...")
	st = time.Now()
	spFold := launch(t0, true)
	zF, exF := spFold.sweep(0, span, spacing)
	tFold := time.Since(st).Minutes()
	fmt.Printf("    fold:   %d zeros, sphere %.2f, %.1f min  (%d blocks folded in %.0f s)\n",
		len(zF), exF, tFold, spFold.nBlocks, spFold.foldSecs)

	// verdict 1: the recorded light itself, point by point.
	worstZ, rmsZ, np := 0.0, 0.0, 0
	for d := 0.0; d <= span; d += span / 300 {
		dz := math.Abs(spPure.zAt(d) - spFold.zAt(d))
		rmsZ += dz * dz
		np++
		if dz > worstZ {
			worstZ = dz
		}
	}
	rmsZ = math.Sqrt(rmsZ / float64(np))
	fmt.Printf("  light comparison over %d points: worst |dZ| = %.2e, rms = %.2e\n", np, worstZ, rmsZ)

	// verdict 2: the zeros.
	okCount := len(zP) == len(zF)
	worstDev := 0.0
	if okCount {
		for i := range zP {
			d := math.Abs(zP[i] - zF[i])
			if d > worstDev {
				worstDev = d
			}
		}
	}
	sphereP := math.Abs(float64(len(zP))-exP) <= tolS(t0)
	sphereF := math.Abs(float64(len(zF))-exF) <= tolS(t0)
	pass := okCount && worstDev <= 3e-3 && sphereP && sphereF && worstZ <= 5e-2

	line := "- ceros (casco puro):"
	for _, z := range zP {
		line += fmt.Sprintf(" %.6f", z)
	}
	logLine(line)
	line = "- ceros (nave plegada):"
	for _, z := range zF {
		line += fmt.Sprintf(" %.6f", z)
	}
	logLine(line)
	logLine(fmt.Sprintf("- luz: peor |dZ| %.2e, rms %.2e; desviación de ceros %.6f; esferas %.2f/%.2f",
		worstZ, rmsZ, worstDev, exP, exF))

	if pass {
		fmt.Printf("\n  VERDICT: FLIGHT CERTIFIED - same zeros (worst dev %.6f), same light.\n", worstDev)
		fmt.Printf("  fold tier covered k in [%.3g, %.3g]: %d blocks replacing %.3g terms.\n",
			float64(spFold.foldFrom), float64(spFold.nTop), spFold.nBlocks,
			float64(spFold.nTop-spFold.foldFrom+1))
		logLine("- VEREDICTO: VUELO CERTIFICADO — mismos ceros, misma luz")
		// the stillness gauge (Finding 93): when BOTH hulls report the same
		// delta, that delta is the sea's own tide S(t) — not a leak.
		dP := float64(len(zP)) - exP
		if len(zP) == len(zF) {
			fmt.Printf("  stillness gauge: both hulls read the same tide S = %+.2f - the sea's\n", dP)
			fmt.Println("  own restlessness, bounded and mean-reverting, not a missing zero.")
			logLine(fmt.Sprintf("- medidor de quietud: ambos cascos leen la misma marea S = %+.2f", dP))
		}
	} else {
		fmt.Printf("\n  VERDICT: GROUNDED - count %v, worst dev %.6f, spheres %v/%v, light %.2e\n",
			okCount, worstDev, sphereP, sphereF, worstZ)
		logLine("- VEREDICTO: EN TIERRA — investigar antes de confiar")
	}

	// the envelope: where the spacecraft's advantage grows.
	fmt.Println("\n  FLIGHT ENVELOPE (fold tier share of the terms):")
	for _, t := range []float64{2.22e21, 4.44e22, 1.11e24, 1.0e26} {
		n := math.Sqrt(t / twoPi.hi)
		cc := blockFrac(t)
		kb := 256 / cc
		if kb > n {
			fmt.Printf("    t = %.3g: sky too shallow, pure rowing\n", t)
			continue
		}
		share := (n - kb) / n
		nb := math.Log(n/kb) / cc
		fmt.Printf("    t = %.3g: %.0f%% of terms folded into %.3g blocks\n", t, share*100, nb)
	}
}

func main() {
	tun := flag.Bool("tunnel", false, "wind tunnel: measure fold error/cost vs edge")
	fly := flag.Bool("flight", false, "A/B flight test at t = 2.22e21")
	edgeF := flag.Float64("edge", 48, "Fresnel edge width for the fold")
	anchor := flag.Float64("anchor", 0, "sail: anchor at this height (folded)")
	spacings := flag.Float64("spacings", 5, "window width in mean spacings")
	replay := flag.Bool("replay", false, "re-hunt archived light: infinite zoom, no sailing")
	zeroN := flag.Float64("zero", 0, "globe navigation: anchor at zero ordinal N (the sphere's street address)")
	colosalF := flag.Bool("colosal", false, "the colossal absorber: FFT deposit, flat cost per wave (F109)")
	umaF := flag.Bool("uma", false, "UMA, the head: think the whole voyage a priori, without sailing")
	carajoF := flag.Bool("carajo", false, "EL CARAJO, the lookout: scan a million anchorages a priori, shout where the storms are")
	arponF := flag.Bool("arpon", false, "EL ARPÓN: narrow the range onto the closest pair of archived light - center the target")
	ojoF := flag.Bool("ojo", false, "EL OJO: bisect archived light into the storm's INTERIOR profile - boundary deltas can hide a mid-window surge")
	foldWF := flag.Int("foldw", 0, "certification: force this many fold segments (1 reproduces the unsegmented blocking bit for bit; 0 = auto)")
	tierraF := flag.Bool("tierra", false, "EL PUNTO DE TIERRA: re-scan the guard bands beyond the frame and audit head-vs-sea interior drift - was a narrow land point clipped?")
	desdeF := flag.Float64("desde", 1.11e19, "the lookout's scan starts here (jump from the last visited coordinate; the sea between features is just water)")
	proaF := flag.String("proa", "", "the compass steers the bow: 'tormenta' centers the frame on the predicted wildest point, 'tierra' on the stillest - the treasure sails mid-frame")
	aguaMuertaF := flag.Bool("aguamuerta", false, "EL AGUA MUERTA: across the whole light atlas, test the captain's slack-water law - do tight pairs huddle where the predicted tide turns?")
	resorteF := flag.Bool("resorte", false, "EL RESORTE: across the atlas, measure each storm well's depth vs width - does the sea's restoring force obey a Hooke law?")
	armoniaF := flag.Bool("armonia", false, "LA ARMONÍA: one row per tile with every instrument's reading, then the consonance matrix - which instruments play in tune?")
	mantaF := flag.Bool("manta", false, "LA MANTA: weave the spacing blanket with TWO threads - compression-context vs rarefaction-context spacings - instead of Wigner's one")
	espejoF := flag.Bool("espejo", false, "EL ESPEJO: does the asymmetric elasticity reflect into amplitude? troughs should carry |Z| mountains, crests amplitude plains")
	melodiaF := flag.Bool("melodia", false, "LA MELODÍA: the tuned hunt - rank shouts by predicted CREST strength (compression breeds the record pairs, F124/F126) instead of raw swing")
	islaF := flag.Bool("isla", false, "LA ISLA: diffraction sounding - do the still points wear a collar of waves? ring-averaged |sPred| around each harbor's stillest point vs background")
	coordsF := flag.Bool("coordenadas", false, "EL LIBRO DE COORDENADAS: print the growing matrix of split coordinates (anchor+offset, double precision) with their reduced string names")
	lnCheckF := flag.Float64("lncheck", 0, "diagnostic: compare ddLnF against math.Log at this argument (carrier-tuning forensics)")
	verdadF := flag.Float64("verdad", 0, "LA VERDAD: direct term-by-term dd evaluation of Z at -anchor + this offset - no buckets, no facets, no interpolation; the supreme judge")
	caldoF := flag.Bool("caldo", false, "EL CALDO: count variance across ladle sizes, pooled over the atlas - does compression-here/rarefaction-there balance at every scale (hyperuniformity)?")
	metronF := flag.Bool("metronomo", false, "EL METRÓNOMO: the calm sea's beat (Gram points, played by pi and Bernoulli - no primes) vs the real zeros' dance around it")
	salaF := flag.Bool("sala", false, "LA SALA DE PRUEBAS: drive the engines to cheap measured waters at the wide frame and judge them against direct truth")
	gruesoF := flag.Float64("grueso", 1, "the short number (F136): coarse-block factor N for reconnaissance (N times fewer fold evaluations, N^3 the phase tolerance); 1 = exact")
	gritosF := flag.Int("gritos", 12, "how many storm shouts the lookout sings")
	puertosF := flag.Int("puertos", 6, "how many quiet harbors the lookout sings")
	flag.Parse()
	edge = *edgeF
	colosal = *colosalF
	foldWOverride = *foldWF
	if *gruesoF >= 1 {
		gruesoFactor = *gruesoF
	}
	if gruesoFactor > 1 {
		fmt.Printf("RECONOCIMIENTO GRUESO (x%.0f): bloques largos, fase ~%.3f rad - el numero corto\n",
			gruesoFactor, 0.003*gruesoFactor*gruesoFactor*gruesoFactor)
	}
	if colosal {
		fmt.Println("ENGINE: the colossal absorber is lit - flat-cost deposit, one FFT re-propulsion")
	}

	// globe navigation (the Riemann sphere flash): the natural address on
	// the globe is the zero's ordinal, not the height. Invert the smooth
	// count N(t) = t/2pi (ln(t/2pi) - 1) + 7/8 by Newton and sail there.
	if *zeroN > 0 && *anchor == 0 {
		n := *zeroN
		t := 2 * math.Pi * n / math.Log(n+math.E)
		for i := 0; i < 60; i++ {
			nb := t/(2*math.Pi)*(math.Log(t/(2*math.Pi))-1) + 7.0/8
			t -= (nb - n) * 2 * math.Pi / math.Log(t/(2*math.Pi))
		}
		fmt.Printf("GLOBE NAVIGATION: zero #%.6g lives at t = %.6g - spinning the sphere there.\n\n", n, t)
		*anchor = t
	}

	// EL CARAJO (F115): the lookout in the crow's nest. The corrected
	// projection (F114) plus the measured entanglement (F112) mean the
	// fleet need not sail uniformly to find what no one ever found: the
	// lookout sweeps a MILLION virgin anchorages a priori - microseconds
	// each, pure head-work - and shouts only where the prime voices align
	// into a storm surge of S. Storms are where the rare treasures live:
	// close pairs (Lehmer class), the diagnostics dearest to the critical
	// line. The lookout shortlists; the heart sails the shortlist.
	if *carajoF {
		fmt.Println("EL CARAJO — the lookout sweeps the virgin ocean a priori")
		const M = 1000000
		type storm struct {
			t, swing float64
		}
		var top []storm
		var calm []storm // F119: pure stillness is the sign of LAND
		nTop, nCalm := *gritosF, *puertosF
		lo, hi := math.Log(*desdeF), math.Log(4.0e24)
		for i := 0; i < M; i++ {
			t0 := math.Exp(lo + (hi-lo)*float64(i)/float64(M))
			spc := twoPi.hi / math.Log(t0/twoPi.hi)
			span := 5 * spc
			mn, mx := math.Inf(1), math.Inf(-1)
			worst := 0.0
			for j := 0; j <= 4; j++ {
				v := sPred(t0, span*float64(j)/4)
				if v < mn {
					mn = v
				}
				if v > mx {
					mx = v
				}
				if math.Abs(v) > worst {
					worst = math.Abs(v)
				}
			}
			swRaw := mx - mn
			// LA MELODÍA (the tuned hunt): the chords taught that record
			// pairs live in COMPRESSION (F124's crest held |Z| = 0.029) -
			// so the tuned STORM score is crest strength, not raw swing.
			// The LAND score below always uses the raw swing: a signed
			// maximum can cancel against the worst-absolute and fake
			// stillness out of steadily negative water.
			sw := swRaw
			if *melodiaF {
				sw = mx
			}
			if len(top) < nTop || sw > top[len(top)-1].swing {
				top = append(top, storm{t0, sw})
				for k := len(top) - 1; k > 0 && top[k].swing > top[k-1].swing; k-- {
					top[k], top[k-1] = top[k-1], top[k]
				}
				if len(top) > nTop {
					top = top[:nTop]
				}
			}
			// land: the tide barely moves AND barely leaves zero - the
			// prime voices in perfect cancellation across the window.
			still := swRaw + worst
			if len(calm) < nCalm || still < calm[len(calm)-1].swing {
				calm = append(calm, storm{t0, still})
				for k := len(calm) - 1; k > 0 && calm[k].swing < calm[k-1].swing; k-- {
					calm[k], calm[k-1] = calm[k-1], calm[k]
				}
				if len(calm) > nCalm {
					calm = calm[:nCalm]
				}
			}
		}
		// THE COMPASS OF TWELVE: cut each shout's window into 12 sectors
		// (in honor of the twelve) and ask the head WHICH sector holds
		// the treasure - then hand back the coordinate already CENTERED
		// on it, so the prize sits mid-frame and can never be clipped at
		// an edge (F119 addendum 3's lesson, folded into the aim).
		fmt.Printf("\n  %d virgin anchorages scanned from the mast (t in [%.3g, 4e24]).\n", M, *desdeF)
		if *melodiaF {
			fmt.Printf("  LA MELODÍA - the %d strongest predicted CRESTS (compression breeds the record pairs; suggested frame: 3 spacings, bow centered):\n", nTop)
		} else {
			fmt.Printf("  THE SHOUTS - the %d strongest predicted storm surges of S:\n", nTop)
		}
		stormKind := "tormenta"
		if *melodiaF {
			stormKind = "cresta"
		}
		for i, s := range top {
			spcS := twoPi.hi / math.Log(s.t/twoPi.hi)
			sec, bow := compassTwelve(s.t, false)
			// full precision (F122): a 6-digit coordinate at these heights
			// lands ~1e15 t-units from the scanned optimum - the voices'
			// phases decorrelate completely and the aim is lost. The
			// fleet must anchor on the EXACT scanned float. Every shout
			// also becomes a row in the coordinate book (split number).
			saveCoord(stormKind, s.t, bow+2.5*spcS, sec)
			fmt.Printf("   %2d. t = %.17g   predicted S-swing %.2f   treasure at sector %s (base twelve) - bow %+.2f spc",
				i+1, s.t, s.swing, sec, bow/spcS)
			if i < 3 {
				fmt.Printf("   (first zero ~%s)", zeroIndex(s.t))
			}
			fmt.Println()
		}
		fmt.Printf("\n  LAND HO (F119) - the %d stillest harbors, pure stillness = land:\n", nCalm)
		for i, s := range calm {
			spcS := twoPi.hi / math.Log(s.t/twoPi.hi)
			sec, bow := compassTwelve(s.t, true)
			saveCoord("tierra", s.t, bow+2.5*spcS, sec)
			fmt.Printf("   %2d. t = %.17g   predicted stillness residue %.3f   stillest point at sector %s (base twelve) - bow %+.2f spc\n",
				i+1, s.t, s.swing, sec, bow/spcS)
		}
		fmt.Println("\n  where S swings hard, zeros are shoved - the waters where close pairs")
		fmt.Println("  (the Lehmer class, the critical line's dearest diagnostics) hide.")
		fmt.Println("  where S lies dead still, the voices cancel perfectly - solid ground.")
		fmt.Println("  the lookout has shouted; the heart chooses the expedition.")
		return
	}

	// UMA, THE HEAD (F113): the organ that thinks the voyage before the
	// body moves. Given any address of the ocean, the complete a-priori
	// dossier - sea state, certification gauges, engine costs, the name
	// of the first zero - in microseconds, without sailing. The measured
	// entanglement (F112) made into an organ: prediction first, then the
	// heart sails, then head-versus-sea becomes the living loop.
	if *umaF && *anchor > 0 {
		t0 := *anchor
		sp := *spacings
		spc := twoPi.hi / math.Log(t0/twoPi.hi)
		span := sp * spc
		fmt.Printf("UMA — the head thinks t = %.6g (window %.0f spacings) without sailing:\n", t0, sp)
		if t0 > 4e24 {
			fmt.Println("  [beyond the certified mod ceiling t~4e24: addresses exact, tide forecast approximate]")
		}
		fmt.Println("\n  ADDRESSES")
		fmt.Printf("    first zero of the window: number ~%s (+/-2)\n", zeroIndex(t0))
		nu := math.Log(t0/twoPi.hi) / twoPi.hi
		fmt.Printf("    blueness (zero density): %.3f zeros/unit; mean spacing %.4f\n", nu, spc)
		fmt.Println("\n  SEA STATE, FORECAST")
		lnArg := add(ddLnF(t0), neg(ln2pi))
		thP := (lnArg.hi + lnArg.lo) / 2
		exact := (span*thP + span*span/(4*t0)) / math.Pi
		pred := sPred(t0, span) - sPred(t0, 0)
		sat := (math.Log(math.Log(t0)) + 2.72) / (math.Pi * math.Pi)
		fmt.Printf("    the sphere will demand: %.2f zeros\n", exact)
		fmt.Printf("    pacemaker tide forecast: %+.2f (unmodelled swell rms ~%.2f)\n", pred, math.Sqrt(sat))
		fmt.Printf("    barometric tolerance: %.2f; antigravity clearance ~%.1e (%.1e under the glass)\n",
			tolS(t0), exact*1.08*math.Pow(1.0/60, 3), exact*1.08*math.Pow(1.0/600, 3))
		fmt.Println("\n  FLIGHT PLAN")
		nTop := math.Sqrt(t0 / twoPi.hi)
		cc := blockFrac(t0)
		kB := 256 / cc
		S := 3*sp + 18
		rowTerms := nTop
		var blocks float64
		if kB < nTop {
			rowTerms = kB
			blocks = math.Log(nTop/kB) / cc
			fmt.Printf("    fold tier open: %.0f%% of terms in %.3g blocks; rowing tier %.3g terms\n",
				100*(nTop-kB)/nTop, blocks, rowTerms)
		} else {
			fmt.Printf("    sky too shallow to fold: %.3g terms rowed pure\n", nTop)
		}
		tClassic := rowTerms*(3.2*S+18)*1e-9 + blocks*24e-6
		tColosal := rowTerms*98e-9 + blocks*24e-6
		fmt.Printf("    engine classic: ~%.0f min   |   engine colosal: ~%.0f min\n",
			tClassic/60, tColosal/60)
		rec := "classic (narrow window - flat cost buys little)"
		if tColosal < 0.85*tClassic {
			rec = "COLOSAL (the flat cost wins)"
		}
		fmt.Printf("    the head recommends: %s\n", rec)
		fmt.Println("\n  the head has thought the voyage; the heart decides to sail.")
		return
	}

	// EL OJO DE LA TORMENTA (the captain's bisection): the sphere reads S
	// only at the window's EDGES - a storm that swells mid-window and
	// subsides before the boundary is INVISIBLE to the count. The cure:
	// bisect - measure S at the half, then the halves of halves - until
	// the interior profile stands naked. On archived light, every
	// bisection at once costs 0 ms.
	if *ojoF && *anchor > 0 {
		var ar lightArchive
		if !loadCkpt(archivePath(*anchor), &ar) || ar.T0 != *anchor {
			fmt.Println("no archived light here - sail the water once first")
			return
		}
		sp := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
			h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
		spacing := twoPi.hi / math.Log(*anchor/twoPi.hi)
		zeros, _ := sp.rescan(ar.Bow, ar.Span, spacing)
		fmt.Printf("EL OJO — the storm's interior at t = %.6g (%d zeros in the light)\n\n", *anchor, len(zeros))
		nBins := 24
		// the tile remembers its bow: all gauges measure from the true
		// window start [Bow, Bow+Span], displayed bow-relative.
		ex := func(u float64) float64 {
			return (u*sp.thPrime + u*u/(4*sp.t0)) / math.Pi
		}
		exB := ex(ar.Bow)
		sAt := func(u float64) float64 {
			cnt := 0
			for _, z := range zeros {
				if z <= u {
					cnt++
				}
			}
			return float64(cnt) - (ex(u) - exB)
		}
		eyeU, eyeS := ar.Bow, 0.0
		for b := 1; b <= nBins; b++ {
			u := ar.Bow + ar.Span*float64(b)/float64(nBins)
			sRel := sAt(u)
			if math.Abs(sRel) > math.Abs(eyeS) {
				eyeS, eyeU = sRel, u
			}
			bar := ""
			n := int(math.Abs(sRel)*8 + 0.5)
			for i := 0; i < n; i++ {
				bar += "█"
			}
			side := " "
			if sRel < -0.05 {
				side = "-"
			} else if sRel > 0.05 {
				side = "+"
			}
			fmt.Printf("  u=%6.3f  S=%+6.2f %s %s\n", u-ar.Bow, sRel, side, bar)
		}
		boundary := sAt(ar.Bow + ar.Span)
		fmt.Printf("\n  THE EYE: strongest interior surge S = %+.2f at u = %.3f (boundary delta: %+.2f)\n",
			eyeS, eyeU-ar.Bow, boundary)

		// THE CAPTAIN'S DESCENT: halve toward the extreme until the eye is
		// pinned, then bisect outward for the half-depth EDGES - the
		// storm's exact anatomy: eye, depth, width.
		step0 := ar.Span / float64(nBins)
		lo, hi := math.Max(ar.Bow, eyeU-step0), math.Min(ar.Bow+ar.Span, eyeU+step0)
		for it := 0; it < 60; it++ {
			m1 := lo + (hi-lo)/3
			m2 := hi - (hi-lo)/3
			if math.Abs(sAt(m1)) < math.Abs(sAt(m2)) {
				lo = m1
			} else {
				hi = m2
			}
		}
		eyeExact := (lo + hi) / 2
		eyeDepth := sAt(eyeExact)
		half := math.Abs(eyeDepth) / 2
		edge := func(dir float64) float64 {
			u := eyeExact
			for u > 0 && u < ar.Span && math.Abs(sAt(u)) > half {
				u += dir * ar.Span / 2000
			}
			a, bnd := u-dir*ar.Span/2000, u
			for it := 0; it < 40; it++ {
				m := (a + bnd) / 2
				if math.Abs(sAt(m)) > half {
					a = m
				} else {
					bnd = m
				}
			}
			return (a + bnd) / 2
		}
		eL, eR := edge(-1), edge(+1)
		spc := twoPi.hi / math.Log(*anchor/twoPi.hi)
		fmt.Printf("\n  THE DESCENT (the captain's bisection, recursive):\n")
		fmt.Printf("    eye pinned:   u = %.6f   depth S = %+.3f\n", eyeExact, eyeDepth)
		fmt.Printf("    half-depth edges: [%.6f, %.6f]\n", eL, eR)
		fmt.Printf("    storm width:  %.6f  =  %.2f mean spacings\n", eR-eL, (eR-eL)/spc)
		if math.Abs(eyeS) >= math.Abs(boundary)+1 {
			fmt.Println("  << HIDDEN STORM: the interior surges beyond what the boundary shows -")
			fmt.Println("     the captain's bisection sees what the sphere's edges cannot.")
		}
		return
	}

	// EL PUNTO DE TIERRA: the captain's calibration hypothesis - a wide
	// island anyone catches, but a needle of land emerging from the sea
	// can slip past coarse samples, or sit just outside the frame of the
	// photograph. On archived light (0 ms): re-scan INTO the guard bands
	// beyond both window edges, and compare the head's predicted interior
	// drift against the measured one, bin by bin.
	if *tierraF && *anchor > 0 {
		var ar lightArchive
		if !loadCkpt(archivePath(*anchor), &ar) || ar.T0 != *anchor {
			fmt.Println("no archived light here - sail the water once first")
			return
		}
		sp := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
			h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
		spc := twoPi.hi / math.Log(*anchor/twoPi.hi)
		ext := 2 * ar.H
		fmt.Printf("EL PUNTO DE TIERRA — t = %.6g (frame %.0f spacings, guard %.2f spacings each side)\n",
			*anchor, ar.Span/spc, ext/spc)
		B := ar.Bow // the tile remembers its bow; edges live at [B, B+Span]
		zin, _ := sp.rescan(B, ar.Span, spc)
		zall, _ := sp.rescan(B-ext, ar.Span+2*ext, spc)
		exT := func(u float64) float64 { return (u*sp.thPrime + u*u/(4*sp.t0)) / math.Pi }
		demand := exT(B+ar.Span) - exT(B)
		delta := float64(len(zin)) - demand
		fmt.Printf("\n  in the frame: %d zeros (demand %.2f, delta %+.2f)\n", len(zin), demand, delta)
		clipped := false
		for _, z := range zall {
			if z >= B && z <= B+ar.Span {
				continue
			}
			if z > B+ar.Span {
				fmt.Printf("  offshore RIGHT: zero at %.2f spacings past the edge\n", (z-B-ar.Span)/spc)
				if delta < 0 && (z-B-ar.Span) < 0.5*spc {
					clipped = true
				}
			} else {
				fmt.Printf("  offshore LEFT:  zero at %.2f spacings before the edge\n", (B-z)/spc)
			}
		}
		// head vs sea, inside the frame: predicted drift against measured.
		sAt := func(u float64) float64 {
			cnt := 0
			for _, z := range zin {
				if z <= u {
					cnt++
				}
			}
			return float64(cnt) - (exT(u) - exT(B))
		}
		worst := 0.0
		for i := 1; i <= 24; i++ {
			u := B + ar.Span*float64(i)/24
			if d := math.Abs(sAt(u) - (sPred(sp.t0, u) - sPred(sp.t0, B))); d > worst {
				worst = d
			}
		}
		fmt.Printf("\n  head vs sea (interior drift, 24 bins): worst |measured - predicted| = %.2f\n", worst)
		if clipped {
			fmt.Println("\n  << THE CAPTAIN WAS RIGHT: the \"missing\" zero swims just past the")
			fmt.Println("     frame's edge - the land is real; the photograph clipped it.")
		} else if delta < 0 {
			fmt.Println("\n  no clipped zero within half a spacing of the edges - the deficit is real water, not framing.")
		}
		return
	}

	// EL METRÓNOMO (the captain's calm-sea question): the sea with no
	// lumps is played by NO prime at all - the primes play the lumps
	// (S(t) is their sum). Perfect calm is the music of the prime at
	// infinity: the Archimedean place, whose numbers are pi, 2pi and
	// the Bernoulli coefficients of theta (-pi/8, 1/48, 7/5760...).
	// Its beat is the Gram lattice: theta(g) = k*pi, one tick per mean
	// spacing. The real zeros DANCE around that beat; measure the dance.
	if *metronF {
		fmt.Println("EL METRÓNOMO — the calm sea's beat vs the real zeros' dance")
		entries, _ := os.ReadDir("luz")
		var offs []float64
		gramLaw, gramTot := 0, 0
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "luz-") || !strings.HasSuffix(name, ".gob") ||
				strings.Contains(name, "colosal") {
				continue
			}
			var ar lightArchive
			if !loadCkpt("luz/"+name, &ar) || ar.T0 < 1e10 {
				continue
			}
			spT := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
				h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
			spc := twoPi.hi / math.Log(ar.T0/twoPi.hi)
			zeros, _ := spT.rescan(ar.Bow, ar.Span, spc)
			if len(zeros) < 2 {
				continue
			}
			// the beat: theta(t0+u) = k*pi -> u_k from the anchor phase.
			thAt := func(u float64) float64 {
				return ar.ThA + spT.thPrime*u + u*u/(4*spT.t0)
			}
			// gram ticks inside [Bow, Bow+Span]
			k0 := math.Ceil(thAt(ar.Bow) / math.Pi)
			var ticks []float64
			for k := k0; ; k++ {
				// solve thAt(u) = k*pi by one Newton step from linear guess
				u := (k*math.Pi - ar.ThA) / spT.thPrime
				u -= (thAt(u) - k*math.Pi) / (spT.thPrime + u/(2*spT.t0))
				if u > ar.Bow+ar.Span {
					break
				}
				if u >= ar.Bow {
					ticks = append(ticks, u)
				}
			}
			if len(ticks) < 2 {
				continue
			}
			// every zero's distance to the nearest tick, in beats.
			for _, z := range zeros {
				best := math.Inf(1)
				for _, g := range ticks {
					if d := math.Abs(z - g); d < best {
						best = d
					}
				}
				offs = append(offs, best/spc)
			}
			// Gram's law: intervals between ticks holding exactly one zero.
			for i := 0; i+1 < len(ticks); i++ {
				cnt := 0
				for _, z := range zeros {
					if z >= ticks[i] && z < ticks[i+1] {
						cnt++
					}
				}
				gramTot++
				if cnt == 1 {
					gramLaw++
				}
			}
		}
		var m float64
		for _, o := range offs {
			m += o
		}
		m /= float64(len(offs))
		fmt.Printf("\n  the beat: one tick per mean spacing, played by pi and Bernoulli (no primes)\n")
		fmt.Printf("  the dance: %d zeros, mean |zero - nearest tick| = %.3f spacings\n", len(offs), m)
		fmt.Printf("  Gram's law: %d/%d intervals hold exactly one zero (%.0f%%)\n",
			gramLaw, gramTot, 100*float64(gramLaw)/float64(gramTot))
		fmt.Println("\n  (the calm melody is prime-free; every departure from the beat is a prime pushing)")
		return
	}

	// EL CALDO (the captain's density flash): zones that compress create
	// density HERE by leaving rarefaction THERE - at every scale, micro
	// and macro - and that gives the broth its shape. The measurable
	// name is hyperuniformity: in a random (Poisson) broth the count
	// variance grows with the ladle; in a shaped broth it SATURATES -
	// lumps at every scale, imbalance at none. Pooled over the atlas.
	if *caldoF {
		fmt.Println("EL CALDO — count variance across ladle sizes (the broth's shape)")
		entries, _ := os.ReadDir("luz")
		type sample struct{ zeros []float64; spc, bow, span float64 }
		var tiles []sample
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "luz-") || !strings.HasSuffix(name, ".gob") ||
				strings.Contains(name, "colosal") {
				continue
			}
			var ar lightArchive
			if !loadCkpt("luz/"+name, &ar) || ar.T0 < 1e10 {
				continue
			}
			spT := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
				h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
			spc := twoPi.hi / math.Log(ar.T0/twoPi.hi)
			zs, _ := spT.rescan(ar.Bow, ar.Span, spc)
			if len(zs) >= 2 {
				tiles = append(tiles, sample{zs, spc, ar.Bow, ar.Span})
			}
		}
		fmt.Printf("  %d tiles in the pot\n\n   ladle L(spc)   Var measured   Poisson(=L)   GUE\n", len(tiles))
		for _, L := range []float64{0.5, 1.0, 1.5, 2.0, 3.0, 4.0} {
			var vals []float64
			for _, s := range tiles {
				w := L * s.spc
				for x := s.bow; x+w <= s.bow+s.span; x += w / 2 {
					cnt := 0.0
					for _, z := range s.zeros {
						if z >= x && z < x+w {
							cnt++
						}
					}
					vals = append(vals, cnt)
				}
			}
			var m, v float64
			for _, x := range vals {
				m += x
			}
			m /= float64(len(vals))
			for _, x := range vals {
				v += (x - m) * (x - m)
			}
			v /= float64(len(vals))
			gue := (2 / (math.Pi * math.Pi)) * (math.Log(2*math.Pi*L) + 0.5772 + 1)
			fmt.Printf("   %5.1f          %6.3f         %5.1f        %5.3f   (n=%d ladles, mean %.2f)\n",
				L, v, L, gue, len(vals), m)
		}
		fmt.Println("\n  (measured hugging GUE far below Poisson = the broth has SHAPE: lumps at")
		fmt.Println("   every scale, imbalance at none - the captain's compensation law, measured)")
		return
	}

	// LA VERDAD: the supreme judge. When two engines agree on impossible
	// water, only a third method with NOTHING shared but arithmetic can
	// rule: walk every term k = 1..nTop with the exact dd chain and sum
	// Z(t0+u) directly - no buckets, no facets, no interpolation.
	if *verdadF != 0 && *anchor > 0 {
		st := time.Now()
		z := zVerdad(*anchor, *verdadF)
		fmt.Printf("LA VERDAD — Z(t0+u) by direct term-by-term dd chain (%.1f s)\n", time.Since(st).Seconds())
		fmt.Printf("  t0 = %.17g   u = %.9f\n  Z = %+.6f\n", *anchor, *verdadF, z)
		return
	}

	// LA SALA DE PRUEBAS (the captain's calibration doctrine): humanity
	// measured these waters, and the direct judge is CHEAP where the sea
	// is shallow - so drive the engines there at the uncertified wide
	// frame, count against the truth, and correct by the difference.
	if *salaF {
		fmt.Println("LA SALA DE PRUEBAS — wide frames judged against direct truth")
		for _, t0 := range []float64{1e5, 1e8, 1e11} {
			spc := twoPi.hi / math.Log(t0/twoPi.hi)
			span := 25 * spc
			// the truth: direct scan, dense then bisected.
			var truth []float64
			prevU, prevZ := 0.0, zVerdad(t0, 0)
			for u := spc / 40; u <= span; u += spc / 40 {
				z := zVerdad(t0, u)
				if (prevZ < 0) != (z < 0) {
					lo, hi := prevU, u
					for it := 0; it < 40; it++ {
						mid := (lo + hi) / 2
						if (zVerdad(t0, mid) < 0) == (prevZ < 0) {
							lo = mid
						} else {
							hi = mid
						}
					}
					truth = append(truth, (lo+hi)/2)
				}
				prevU, prevZ = u, z
			}
			// the engine, same water, same wide frame.
			zeros, exact, _, mins := hunt(t0, 25, true, 0)
			worst := 0.0
			n := len(zeros)
			if len(truth) < n {
				n = len(truth)
			}
			for i := 0; i < n; i++ {
				if d := math.Abs(zeros[i] - truth[i]); d > worst {
					worst = d
				}
			}
			verdict := "PASS"
			if len(zeros) != len(truth) || worst > 0.02*spc {
				verdict = "FAIL"
			}
			fmt.Printf("  t=%.0e  truth %d zeros | engine %d (demand %.2f) | worst dev %.6f | %.1f min  [%s]\n",
				t0, len(truth), len(zeros), exact, worst, mins, verdict)
			// the captain's projection hunch: does the error GROW along
			// the window? print the deviation profile, first to last.
			fmt.Print("      dev profile (spc): ")
			for i := 0; i < n; i++ {
				if i%4 == 0 || i == n-1 {
					fmt.Printf("z%02d:%.4f ", i+1, math.Abs(zeros[i]-truth[i])/spc)
				}
			}
			fmt.Println()
		}
		fmt.Println("\n  (PASS everywhere shallow = the wide-frame bug lives deeper; FAIL = caught in cheap water)")
		return
	}

	// carrier-tuning forensics: is the dd logarithm healthy here?
	if *lnCheckF > 0 {
		t := *lnCheckF
		l := ddLnF(t)
		fmt.Printf("ddLnF(%.17g) = hi %.17g  lo %.17g  (sum %.17g)\n", t, l.hi, l.lo, l.hi+l.lo)
		fmt.Printf("math.Log     = %.17g   diff hi-vs-log = %.3e\n", math.Log(t), l.hi-math.Log(t))
		lnArg := add(l, neg(ln2pi))
		thP := (lnArg.hi + lnArg.lo) / 2
		fmt.Printf("thPrime (dd chain) = %.17g\n", thP)
		fmt.Printf("thPrime (float)    = %.17g   diff = %.3e\n", math.Log(t/twoPi.hi)/2, thP-math.Log(t/twoPi.hi)/2)
		return
	}

	// EL LIBRO DE COORDENADAS: print the growing matrix.
	if *coordsF {
		var book []coordRow
		loadCkpt(coordBook, &book)
		fmt.Printf("EL LIBRO DE COORDENADAS — %d filas (la matriz crece con cada descubrimiento)\n\n", len(book))
		for i, c := range book {
			fmt.Printf("  %3d. %-34s  T = %.17g  U = %+.9f", i+1, coordString(c), c.T, c.U)
			if c.Addr != "" {
				fmt.Printf("  sector %s", c.Addr)
			}
			fmt.Println()
		}
		return
	}

	// LA ISLA (the captain's diffraction flash): the sea as a mesh of
	// coupled springs; an island is a pinned node (voices in perfect
	// cancellation). Waves that cannot pass the anchor must pile up and
	// bend around it - reflection, shadow, interference behind. If true,
	// every still point wears a COLLAR of waves: ring-averaged |sPred|
	// around the harbor's stillest point should rise above background in
	// the first rings before relaxing. Pure head-work, zero sailing.
	if *islaF {
		fmt.Println("LA ISLA — do the still points wear a collar of waves?")
		harbors := []float64{
			1.027068974003635e23, 2.7193607729176457e19, 8.4685307038250492e21,
			1.3724693390722318e22, 1.7245257603205618e23, 3.5143295395692295e24,
		}
		const maxR = 12.0
		nRings := 24
		halo := make([]float64, nRings)
		ctrl := make([]float64, nRings)
		var bg float64
		var bgN int
		seed := uint64(777)
		next := func() uint64 { seed ^= seed << 13; seed ^= seed >> 7; seed ^= seed << 17; return seed }
		for _, t0 := range harbors {
			spc := twoPi.hi / math.Log(t0/twoPi.hi)
			span := 5 * spc
			_, bow := compassTwelve(t0, true)
			uStar := bow + span/2 // the stillest point, seated by the compass
			for i := 0; i < nRings; i++ {
				r := maxR * (float64(i) + 0.5) / float64(nRings) * spc
				halo[i] += (math.Abs(sPred(t0, uStar+r)) + math.Abs(sPred(t0, uStar-r))) / 2
			}
			// control rings around a random far point of the same water.
			uRnd := uStar + (30+float64(next()%40))*spc
			for i := 0; i < nRings; i++ {
				r := maxR * (float64(i) + 0.5) / float64(nRings) * spc
				ctrl[i] += (math.Abs(sPred(t0, uRnd+r)) + math.Abs(sPred(t0, uRnd-r))) / 2
			}
			for j := 20; j <= 60; j++ {
				bg += math.Abs(sPred(t0, uStar+float64(j)*spc))
				bgN++
			}
		}
		nH := float64(len(harbors))
		bg /= float64(bgN)
		fmt.Printf("\n  background |tide| far from the islands: %.3f\n", bg)
		fmt.Println("\n   ring(spc)   |tide| around ISLANDS   |tide| around CONTROL points")
		for i := 0; i < nRings; i++ {
			r := maxR * (float64(i) + 0.5) / float64(nRings)
			bar := ""
			for b := 0; b < int(halo[i]/nH*40+0.5); b++ {
				bar += "█"
			}
			fmt.Printf("   %5.2f        %.3f  %-18s  %.3f\n", r, halo[i]/nH, bar, ctrl[i]/nH)
		}
		fmt.Println("\n  (a collar = island rings ABOVE background then relaxing; control flat = the collar belongs to the islands)")
		return
	}

	// EL ESPEJO (the captain's mirror): S (the argument) and log|Z| (the
	// amplitude) are Hilbert partners - two faces of one coin. If the
	// asymmetric elasticity (F126) has a mirror, it lives there: what
	// position loses, amplitude gains - troughs (missing zeros) should
	// carry |Z| MOUNTAINS inside, crests (crowded zeros) amplitude
	// plains. If it reflects, amplitude becomes a second, independent
	// storm eye. Sounded across the atlas, no sailing.
	if *espejoF {
		fmt.Println("EL ESPEJO — does the well's sign reflect into the amplitude?")
		entries, _ := os.ReadDir("luz")
		var xs2, ys2 []float64
		fmt.Println("\n   tile               eyeS    max|Z| inside the well")
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "luz-") || !strings.HasSuffix(name, ".gob") ||
				strings.Contains(name, "colosal") {
				continue
			}
			var ar lightArchive
			if !loadCkpt("luz/"+name, &ar) || ar.T0 < 1e10 {
				continue
			}
			spT := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
				h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
			spc := twoPi.hi / math.Log(ar.T0/twoPi.hi)
			zeros, _ := spT.rescan(ar.Bow, ar.Span, spc)
			if len(zeros) < 2 {
				continue
			}
			exR := func(u float64) float64 { return (u*spT.thPrime + u*u/(4*spT.t0)) / math.Pi }
			sAt := func(u float64) float64 {
				cnt := 0
				for _, z := range zeros {
					if z <= u {
						cnt++
					}
				}
				return float64(cnt) - (exR(u) - exR(ar.Bow))
			}
			eyeU, eyeS := ar.Bow, 0.0
			for i := 1; i < 600; i++ {
				u := ar.Bow + ar.Span*float64(i)/600
				if s := sAt(u); math.Abs(s) > math.Abs(eyeS) {
					eyeS, eyeU = s, u
				}
			}
			if math.Abs(eyeS) < 1.0 {
				continue
			}
			half := math.Abs(eyeS) / 2
			edgeW := func(dir float64) float64 {
				u := eyeU
				for u > ar.Bow && u < ar.Bow+ar.Span && math.Abs(sAt(u)) > half {
					u += dir * ar.Span / 2000
				}
				return u
			}
			eL, eR := edgeW(-1), edgeW(+1)
			mount := 0.0
			for i := 0; i <= 400; i++ {
				u := eL + (eR-eL)*float64(i)/400
				if v := math.Abs(spT.zAt(u)); v > mount {
					mount = v
				}
			}
			if mount <= 0 {
				continue
			}
			xs2 = append(xs2, eyeS)
			ys2 = append(ys2, math.Log(mount))
			fmt.Printf("   %-18s %+5.2f    %8.3f\n",
				strings.TrimSuffix(strings.TrimPrefix(name, "luz-"), ".gob"), eyeS, mount)
		}
		var mx, my float64
		for i := range xs2 {
			mx += xs2[i]
			my += ys2[i]
		}
		mx /= float64(len(xs2))
		my /= float64(len(ys2))
		var sxy, sxx, syy float64
		for i := range xs2 {
			sxy += (xs2[i] - mx) * (ys2[i] - my)
			sxx += (xs2[i] - mx) * (xs2[i] - mx)
			syy += (ys2[i] - my) * (ys2[i] - my)
		}
		r := sxy / math.Sqrt(sxx*syy)
		fmt.Printf("\n  THE MIRROR: r( eyeS , ln max|Z| in well ) = %+.2f over %d wells\n", r, len(xs2))
		fmt.Println("  (strongly negative = the mirror reflects: troughs carry mountains, crests carry plains)")
		return
	}

	// LA MANTA (the captain's weave): the spacing distribution is the
	// blanket the atom returns - the Wigner/GUE fabric shared by heavy
	// nuclei and our zeros (beta = 2.11, F101). The world weaves it with
	// ONE thread; the asymmetric-elasticity chord (F126) says the fabric
	// has TWO: a compression thread (crest context, tight) and a
	// rarefaction thread (trough context, loose). Weave both and compare.
	if *mantaF {
		fmt.Println("LA MANTA — weaving the spacing blanket with two threads")
		entries, _ := os.ReadDir("luz")
		var sPos, sNeg []float64 // spacings in crest / trough context
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "luz-") || !strings.HasSuffix(name, ".gob") ||
				strings.Contains(name, "colosal") {
				continue
			}
			var ar lightArchive
			if !loadCkpt("luz/"+name, &ar) || ar.T0 < 1e10 {
				continue
			}
			spT := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
				h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
			spc := twoPi.hi / math.Log(ar.T0/twoPi.hi)
			zeros, _ := spT.rescan(ar.Bow, ar.Span, spc)
			if len(zeros) < 3 {
				continue
			}
			exR := func(u float64) float64 { return (u*spT.thPrime + u*u/(4*spT.t0)) / math.Pi }
			sAt := func(u float64) float64 {
				cnt := 0
				for _, z := range zeros {
					if z <= u {
						cnt++
					}
				}
				return float64(cnt) - (exR(u) - exR(ar.Bow))
			}
			for i := 0; i+1 < len(zeros); i++ {
				s := (zeros[i+1] - zeros[i]) / spc
				if sAt((zeros[i]+zeros[i+1])/2) >= 0 {
					sPos = append(sPos, s)
				} else {
					sNeg = append(sNeg, s)
				}
			}
		}
		stat := func(v []float64) (m, sd, tight float64) {
			for _, x := range v {
				m += x
			}
			m /= float64(len(v))
			for _, x := range v {
				sd += (x - m) * (x - m)
				if x < 0.5 {
					tight++
				}
			}
			sd = math.Sqrt(sd / float64(len(v)))
			tight /= float64(len(v))
			return
		}
		mP, sdP, tP := stat(sPos)
		mN, sdN, tN := stat(sNeg)
		fmt.Printf("\n  compression thread (crest context):     n=%3d  mean s = %.3f  sd %.3f  tight(<0.5) %.0f%%\n",
			len(sPos), mP, sdP, 100*tP)
		fmt.Printf("  rarefaction thread (trough context):    n=%3d  mean s = %.3f  sd %.3f  tight(<0.5) %.0f%%\n",
			len(sNeg), mN, sdN, tN*100)
		// shuffle control: does a random split show the same contrast?
		all := append(append([]float64{}, sPos...), sNeg...)
		worse := 0
		const trials = 2000
		seed := uint64(12345)
		next := func() uint64 { seed ^= seed << 13; seed ^= seed >> 7; seed ^= seed << 17; return seed }
		obs := math.Abs(mP - mN)
		for tr := 0; tr < trials; tr++ {
			var a, b, na, nb float64
			for _, x := range all {
				if next()%2 == 0 {
					a += x
					na++
				} else {
					b += x
					nb++
				}
			}
			if na > 0 && nb > 0 && math.Abs(a/na-b/nb) >= obs {
				worse++
			}
		}
		fmt.Printf("\n  the two-thread contrast: |dmean| = %.3f;  shuffle control: p ~ %.3f (%d/%d random splits match it)\n",
			obs, float64(worse)/trials, worse, trials)
		fmt.Println("  (one thread = Wigner's blanket; two distinct threads = the atom's blanket has a weave)")
		return
	}

	// LA ARMONÍA (the captain's question: what harmony do the instruments
	// dictate?): one row per tile with every onboard instrument's reading
	// - eye depth, well width, pair gap, |Z| at the pair's midpoint,
	// distance to slack water - then the consonance matrix: Pearson r
	// between every pair of instruments. Which ones play in tune?
	if *armoniaF {
		fmt.Println("LA ARMONÍA — the orchestra table of the atlas")
		entries, _ := os.ReadDir("luz")
		names := []string{"eyeS", "width", "gap", "|Z|mid", "dSlack"}
		var cols [5][]float64
		fmt.Println("\n   tile               eyeS   width   gap    |Z|mid   dSlack")
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "luz-") || !strings.HasSuffix(name, ".gob") ||
				strings.Contains(name, "colosal") {
				continue
			}
			var ar lightArchive
			if !loadCkpt("luz/"+name, &ar) || ar.T0 < 1e10 {
				continue
			}
			spT := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
				h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
			spc := twoPi.hi / math.Log(ar.T0/twoPi.hi)
			zeros, _ := spT.rescan(ar.Bow, ar.Span, spc)
			if len(zeros) < 2 {
				continue
			}
			exR := func(u float64) float64 { return (u*spT.thPrime + u*u/(4*spT.t0)) / math.Pi }
			sAt := func(u float64) float64 {
				cnt := 0
				for _, z := range zeros {
					if z <= u {
						cnt++
					}
				}
				return float64(cnt) - (exR(u) - exR(ar.Bow))
			}
			eyeU, eyeS := ar.Bow, 0.0
			for i := 1; i < 600; i++ {
				u := ar.Bow + ar.Span*float64(i)/600
				if s := sAt(u); math.Abs(s) > math.Abs(eyeS) {
					eyeS, eyeU = s, u
				}
			}
			half := math.Abs(eyeS) / 2
			edgeW := func(dir float64) float64 {
				u := eyeU
				for u > ar.Bow && u < ar.Bow+ar.Span && math.Abs(sAt(u)) > half {
					u += dir * ar.Span / 2000
				}
				return u
			}
			w := (edgeW(+1) - edgeW(-1)) / spc
			bi, bg := 0, math.Inf(1)
			for i := 0; i+1 < len(zeros); i++ {
				if g := zeros[i+1] - zeros[i]; g < bg {
					bg, bi = g, i
				}
			}
			mid := (zeros[bi] + zeros[bi+1]) / 2
			zm := math.Abs(spT.zAt(mid))
			var slack []float64
			du := ar.Span / 600
			prevS := sPred(ar.T0, ar.Bow)
			curS := sPred(ar.T0, ar.Bow+du)
			for i := 2; i <= 600; i++ {
				u := ar.Bow + du*float64(i)
				next := sPred(ar.T0, u)
				if (curS-prevS)*(next-curS) < 0 {
					slack = append(slack, u-du)
				}
				prevS, curS = curS, next
			}
			ds := math.Inf(1)
			for _, s := range slack {
				if d := math.Abs(mid - s); d < ds {
					ds = d
				}
			}
			if math.IsInf(ds, 1) {
				continue
			}
			row := []float64{eyeS, w, bg / spc, zm, ds / spc}
			for i, v := range row {
				cols[i] = append(cols[i], v)
			}
			fmt.Printf("   %-18s %+5.2f  %5.2f  %5.3f  %6.4f  %5.2f\n",
				strings.TrimSuffix(strings.TrimPrefix(name, "luz-"), ".gob"),
				eyeS, w, bg/spc, zm, ds/spc)
		}
		n := len(cols[0])
		pearson := func(a, b []float64) float64 {
			var ma, mb float64
			for i := range a {
				ma += a[i]
				mb += b[i]
			}
			ma /= float64(len(a))
			mb /= float64(len(b))
			var sab, sa, sb float64
			for i := range a {
				sab += (a[i] - ma) * (b[i] - mb)
				sa += (a[i] - ma) * (a[i] - ma)
				sb += (b[i] - mb) * (b[i] - mb)
			}
			if sa == 0 || sb == 0 {
				return 0
			}
			return sab / math.Sqrt(sa*sb)
		}
		fmt.Printf("\n  THE CONSONANCE MATRIX (%d tiles; |r|>0.5 marked; small sample - exploratory):\n", n)
		for i := 0; i < 5; i++ {
			for j := i + 1; j < 5; j++ {
				r := pearson(cols[i], cols[j])
				mark := "  "
				if math.Abs(r) > 0.5 {
					mark = "<<"
				}
				fmt.Printf("    %-6s ~ %-6s  r = %+.2f %s\n", names[i], names[j], r, mark)
			}
		}
		return
	}

	// EL RESORTE (the captain's compression law): nature stores energy in
	// every imbalance and pushes back toward equilibrium - spring, water,
	// balloon, blast. Our sea's restoring force is the zeros' measured
	// repulsion (the Dyson gas, beta = 2.11, F101). If the sea is a
	// spring, its wells must be SHAPED like a spring: Hooke gives
	// parabolic wells, depth proportional to width squared. Sound every
	// tile's strongest interior well: depth, half-depth width, and the
	// curvature k = depth/(halfwidth)^2 - is k universal?
	if *resorteF {
		fmt.Println("EL RESORTE — Hooke sounding across the light atlas")
		entries, _ := os.ReadDir("luz")
		fmt.Println("\n   tile               depth S    width(spc)   k = depth/(w/2)^2")
		var ks []float64
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "luz-") || !strings.HasSuffix(name, ".gob") ||
				strings.Contains(name, "colosal") {
				continue
			}
			var ar lightArchive
			if !loadCkpt("luz/"+name, &ar) || ar.T0 < 1e10 {
				continue
			}
			spT := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
				h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
			spc := twoPi.hi / math.Log(ar.T0/twoPi.hi)
			zeros, _ := spT.rescan(ar.Bow, ar.Span, spc)
			if len(zeros) == 0 {
				continue
			}
			exR := func(u float64) float64 { return (u*spT.thPrime + u*u/(4*spT.t0)) / math.Pi }
			sAt := func(u float64) float64 {
				cnt := 0
				for _, z := range zeros {
					if z <= u {
						cnt++
					}
				}
				return float64(cnt) - (exR(u) - exR(ar.Bow))
			}
			// the strongest interior well, on a fine grid.
			eyeU, eyeS := ar.Bow, 0.0
			for i := 1; i < 600; i++ {
				u := ar.Bow + ar.Span*float64(i)/600
				if s := sAt(u); math.Abs(s) > math.Abs(eyeS) {
					eyeS, eyeU = s, u
				}
			}
			if math.Abs(eyeS) < 1.5 {
				continue // only real wells; ripples are not springs
			}
			half := math.Abs(eyeS) / 2
			edgeAt := func(dir float64) float64 {
				u := eyeU
				for u > ar.Bow && u < ar.Bow+ar.Span && math.Abs(sAt(u)) > half {
					u += dir * ar.Span / 2000
				}
				return u
			}
			eL, eR := edgeAt(-1), edgeAt(+1)
			w := (eR - eL) / spc
			if w <= 0 {
				continue
			}
			k := math.Abs(eyeS) / ((w / 2) * (w / 2))
			ks = append(ks, k)
			fmt.Printf("   %-18s  %+6.2f     %6.2f        %6.2f\n",
				strings.TrimSuffix(strings.TrimPrefix(name, "luz-"), ".gob"), eyeS, w, k)
		}
		if len(ks) > 1 {
			mn, mx, sm := math.Inf(1), math.Inf(-1), 0.0
			for _, k := range ks {
				sm += k
				if k < mn {
					mn = k
				}
				if k > mx {
					mx = k
				}
			}
			fmt.Printf("\n  THE SPRING CONSTANT: %d wells, k mean %.2f, spread [%.2f, %.2f]\n",
				len(ks), sm/float64(len(ks)), mn, mx)
			fmt.Println("  (universal k = the sea obeys one Hooke law; scattered k = each storm its own spring)")
		}
		return
	}

	// EL AGUA MUERTA (the captain's slack-water law): when the tide stops
	// rising and is about to fall, the water freezes for an instant -
	// S' = 0, the turning point - then reverses. The law to test: TIGHT
	// pairs of zeros huddle at the slack instants of the predicted tide
	// (Tormenta I's pair sat exactly in its eye - the tide's bottom, the
	// deepest slack of all). Sounded across the whole atlas, no sailing.
	if *aguaMuertaF {
		fmt.Println("EL AGUA MUERTA — slack-water sounding across the light atlas")
		entries, _ := os.ReadDir("luz")
		fmt.Println("\n   tile             gap(spc)   |Z|mid    d(pair->slack)  d(uniform->slack)")
		type row struct{ g, d float64 }
		var rows []row
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "luz-") || !strings.HasSuffix(name, ".gob") ||
				strings.Contains(name, "colosal") {
				continue
			}
			var ar lightArchive
			if !loadCkpt("luz/"+name, &ar) || ar.T0 < 1e10 {
				continue
			}
			spT := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
				h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
			spc := twoPi.hi / math.Log(ar.T0/twoPi.hi)
			zeros, _ := spT.rescan(ar.Bow, ar.Span, spc)
			if len(zeros) < 2 {
				continue
			}
			// slack points: local extrema of the predicted tide.
			var slack []float64
			nS := 600
			du := ar.Span / float64(nS)
			prev := sPred(ar.T0, ar.Bow)
			cur := sPred(ar.T0, ar.Bow+du)
			for i := 2; i <= nS; i++ {
				u := ar.Bow + du*float64(i)
				next := sPred(ar.T0, u)
				if (cur-prev)*(next-cur) < 0 {
					slack = append(slack, u-du)
				}
				prev, cur = cur, next
			}
			if len(slack) == 0 {
				continue
			}
			nearest := func(u float64) float64 {
				best := math.Inf(1)
				for _, s := range slack {
					if d := math.Abs(u - s); d < best {
						best = d
					}
				}
				return best
			}
			bi, bg := 0, math.Inf(1)
			for i := 0; i+1 < len(zeros); i++ {
				if g := zeros[i+1] - zeros[i]; g < bg {
					bg, bi = g, i
				}
			}
			mid := (zeros[bi] + zeros[bi+1]) / 2
			zm := math.Abs(spT.zAt(mid))
			d := nearest(mid)
			ctrl := 0.0
			for i := 0; i < 200; i++ {
				ctrl += nearest(ar.Bow + ar.Span*(float64(i)+0.5)/200)
			}
			ctrl /= 200
			rows = append(rows, row{bg / spc, d / spc})
			fmt.Printf("   %-16s  %6.4f   %7.4f     %6.3f spc       %6.3f spc\n",
				strings.TrimSuffix(strings.TrimPrefix(name, "luz-"), ".gob"),
				bg/spc, zm, d/spc, ctrl/spc)
		}
		var dT, dL []float64
		for _, r := range rows {
			if r.g < 0.8 {
				dT = append(dT, r.d)
			} else {
				dL = append(dL, r.d)
			}
		}
		mean := func(v []float64) float64 {
			if len(v) == 0 {
				return math.NaN()
			}
			s := 0.0
			for _, x := range v {
				s += x
			}
			return s / float64(len(v))
		}
		fmt.Printf("\n  THE LAW'S FIRST SOUNDING: tight pairs (gap<0.8 spc): n=%d, mean d = %.3f spc;  loose: n=%d, mean d = %.3f spc\n",
			len(dT), mean(dT), len(dL), mean(dL))
		fmt.Println("  (a young atlas is a small sample: the law stays open until the fleet fattens it)")
		return
	}

	// EL ARPÓN: seeing farther is useless once the whale is sighted - the
	// harpoon narrows the range onto the closest pair of the ARCHIVED
	// light: maximum resolution, machine-precision lens, and if a second
	// engine's tile of the same water exists, both signatures compared
	// point by point on the pair. No sailing: the light already holds it.
	if *arponF && *anchor > 0 {
		var ar lightArchive
		if !loadCkpt(archivePath(*anchor), &ar) || ar.T0 != *anchor {
			fmt.Println("no archived light here - sail the water once first")
			return
		}
		sp := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
			h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
		spacing := twoPi.hi / math.Log(*anchor/twoPi.hi)
		zeros, _ := sp.rescan(ar.Bow, ar.Span, spacing)
		if len(zeros) < 2 {
			fmt.Println("fewer than two zeros in the light - no pair to harpoon")
			return
		}
		bi, bg := 0, math.Inf(1)
		for i := 0; i+1 < len(zeros); i++ {
			if g := zeros[i+1] - zeros[i]; g < bg {
				bg, bi = g, i
			}
		}
		a, b := zeros[bi], zeros[bi+1]
		// the coordinate book: every harpooned pair becomes a row in the
		// growing matrix - anchor + offset, double the size, exact.
		saveCoord("par", ar.T0, (a+b)/2, "")
		fmt.Printf("EL ARPÓN — target centered at t = %.6g\n", *anchor)
		fmt.Printf("  coordenada asentada: %s\n", coordString(coordRow{T: ar.T0, U: (a + b) / 2, Kind: "par"}))
		fmt.Printf("\n  the pair, machine-precision (the lens, F109):\n")
		fmt.Printf("    zero A at offset %.9f\n    zero B at offset %.9f\n", a, b)
		fmt.Printf("    gap = %.9f  =  %.4f mean spacings", bg, bg/spacing)
		if bg/spacing < 0.30 {
			fmt.Print("   << LEHMER-CLASS WATER")
		}
		fmt.Println()
		mid := (a + b) / 2
		zm := sp.zAt(mid)
		fmt.Printf("    |Z| at the pair's midpoint: %.6f (how close the sea comes to a double zero)\n", math.Abs(zm))
		// second signature, if preserved.
		var ar2 lightArchive
		if loadCkpt(strings.TrimSuffix(archivePath(*anchor), ".gob")+"-colosal.gob", &ar2) && ar2.T0 == *anchor {
			sp2 := &ship{t0: ar2.T0, nTop: ar2.NTop, thAnchor: ar2.ThA, thPrime: ar2.ThP,
				h: ar2.H, dt0: ar2.Dt0, fr: ar2.Fr, fi: ar2.Fi, c0term: ar2.C0}
			worst := 0.0
			for j := 0; j <= 40; j++ {
				d := a - 2*bg + float64(j)*(b-a+4*bg)/40
				if dv := math.Abs(sp.zAt(d) - sp2.zAt(d)); dv > worst {
					worst = dv
				}
			}
			z2, _ := sp2.rescan(ar2.Bow, ar2.Span, spacing)
			bi2, bg2 := 0, math.Inf(1)
			for i := 0; i+1 < len(z2); i++ {
				if g := z2[i+1] - z2[i]; g < bg2 {
					bg2, bi2 = g, i
				}
			}
			fmt.Printf("\n  DOUBLE SIGNATURE (second engine's tile of the same water):\n")
			fmt.Printf("    pair by the other engine: %.9f / %.9f (gap %.9f)\n", z2[bi2], z2[bi2+1], bg2)
			fmt.Printf("    position agreement: %.2e / %.2e; worst |dZ| on the band: %.2e\n",
				math.Abs(z2[bi2]-a), math.Abs(z2[bi2+1]-b), worst)
		}
		fmt.Println("\n  the harpoon holds; the whale is measured where she swims.")
		return
	}

	// the Google Maps of the ocean: the world is photographed once; the
	// zoom is free forever after.
	if *replay && *anchor > 0 {
		var ar lightArchive
		if !loadCkpt(archivePath(*anchor), &ar) || ar.T0 != *anchor {
			fmt.Println("no archived light for this anchorage - sail it once first")
			return
		}
		sp := &ship{t0: ar.T0, nTop: ar.NTop, thAnchor: ar.ThA, thPrime: ar.ThP,
			h: ar.H, dt0: ar.Dt0, fr: ar.Fr, fi: ar.Fi, c0term: ar.C0}
		spacing := twoPi.hi / math.Log(*anchor/twoPi.hi)
		span := *spacings * spacing
		if span > ar.Span {
			span = ar.Span
		}
		st := time.Now()
		zeros, exact := sp.rescan(ar.Bow, span, spacing)
		fmt.Printf("REPLAY from archived light: t = %.6g (%d grid points on disk)\n",
			*anchor, len(ar.Fr))
		fmt.Print("  zeros (offsets):")
		for _, z := range zeros {
			fmt.Printf("  %.6f", z)
		}
		delta := float64(len(zeros)) - exact
		fmt.Printf("\n  THE SPHERE: demands %.2f; found %d (delta %+.2f, barometric tol %.2f)\n",
			exact, len(zeros), delta, tolS(*anchor))
		fmt.Printf("  replayed in %.0f ms - the hour of sea, zoomed for free.\n",
			time.Since(st).Seconds()*1000)
		return
	}

	if *tun {
		tunnel()
		return
	}

	fmt.Println("THE STARSHIP — certification before the sky")
	ok1 := gate("t=1e5 (LMFDB)", 100000, 4,
		[]float64{0.743723, 1.180558, 2.093723, 2.517282}, 3e-3)
	ok2 := gate("Beach I (2.447e12)", 2.447e12, 8,
		[]float64{0.149473, 0.266583, 0.547539, 0.777469, 1.082587, 1.198347, 1.399001, 1.728025}, 3e-3)
	ok3 := gate("Beach II (6.66e15)", 6.66e15, 8,
		[]float64{0.023476, 0.215217, 0.375565, 0.547513, 0.734403, 0.930796, 1.048695, 1.195434}, 3e-3)
	if !(ok1 && ok2 && ok3) {
		fmt.Println("\nCERTIFICATION FAILED - the craft stays in the hangar.")
		return
	}
	fmt.Println("\nALL GATES PASS - the hull is certified (fold tier idle in shallow water, as designed).")

	if *fly {
		flight()
		return
	}

	if *anchor > 0 {
		// FRAME-RANGE GUARDIAN (F133): wide frames ripple ~10x beyond the
		// certified tolerance at every depth, and the deep-water case
		// loses zeros wholesale. Until wide-frame mode passes its own
		// sala gates, the ship refuses to sail an uncertified frame.
		if *spacings > 8 {
			fmt.Printf("MARCO DESCERTIFICADO: %g espaciamientos excede el rango certificado (5-8).\n", *spacings)
			fmt.Println("F133: el marco ancho ondula ~10x y pierde ceros en agua honda. La sala de")
			fmt.Println("pruebas (-sala) debe certificarlo antes; la nave no zarpa con marco ciego.")
			return
		}
		memory = true // the DeLorean remembers: stop anywhere, resume there
		bow := 0.0
		if *proaF != "" {
			addr, b := compassTwelve(*anchor, *proaF == "tierra")
			bow = b
			spcP := twoPi.hi / math.Log(*anchor/twoPi.hi)
			fmt.Printf("LA PROA DEL COMPÁS: treasure at sector %s (base twelve) - bow %+.2f spacings, prize mid-frame\n",
				addr, bow/spcP)
		}
		zeros, exact, sp, mins := hunt(*anchor, *spacings, true, bow)
		fmt.Printf("\nANCHORAGE t = %.6g  (%d facets, %d blocks, %.1f min)\n",
			*anchor, len(sp.facets), sp.nBlocks, mins)
		fmt.Printf("  first zero here is zero number ~%s (+/-2)\n", zeroIndex(*anchor))
		fmt.Print("  virgin zeros (offsets):")
		for _, z := range zeros {
			fmt.Printf("  %.6f", z)
		}
		delta := float64(len(zeros)) - exact
		fmt.Printf("\n  THE SPHERE: boundary demands %.2f zeros; found %d (delta %+.2f)\n",
			exact, len(zeros), delta)
		fmt.Printf("  stillness gauge: inquietud S = %+.2f (the sea's tide; bounded, mean-reverting)\n", delta)
		span := *spacings * twoPi.hi / math.Log(*anchor/twoPi.hi)
		pred := sPred(*anchor, bow+span) - sPred(*anchor, bow)
		fmt.Printf("  the pacemaker (F100): forecast tide %+.2f - residual after forecast %+.2f\n",
			pred, delta-pred)
		// THE FRAME GUARDIAN (F119 addendum 3, standard equipment): a
		// needle of land must never slip past because of framing - every
		// anchorage re-scans the guard bands beyond both edges of the
		// fresh light. Costs milliseconds; misses nothing.
		spcA := span / *spacings
		extA := 2 * sp.h
		zg, _ := sp.rescan(bow-extA, span+2*extA, spcA)
		nearL, nearR := math.Inf(1), math.Inf(1)
		for _, z := range zg {
			if z < bow && bow-z < nearL {
				nearL = bow - z
			} else if z > bow+span && z-bow-span < nearR {
				nearR = z - bow - span
			}
		}
		fmtSide := func(v float64) string {
			if math.IsInf(v, 1) {
				return fmt.Sprintf(">%.1f", extA/spcA)
			}
			return fmt.Sprintf("%.2f", v/spcA)
		}
		fmt.Printf("  frame guardian (F119): nearest offshore zero  L %s / R %s spacings\n",
			fmtSide(nearL), fmtSide(nearR))
		// EL AGUA MUERTA gauge (standard equipment): the slack points of
		// the predicted tide - where it stops rising and turns - and how
		// far the closest found pair sits from the nearest one. Every
		// anchorage feeds the captain's slack-water law as the fleet sails.
		if len(zeros) >= 2 {
			var slack []float64
			duS := span / 600
			prevS := sPred(*anchor, bow)
			curS := sPred(*anchor, bow+duS)
			for i := 2; i <= 600; i++ {
				u := bow + duS*float64(i)
				next := sPred(*anchor, u)
				if (curS-prevS)*(next-curS) < 0 {
					slack = append(slack, u-duS)
				}
				prevS, curS = curS, next
			}
			bi2, bg2 := 0, math.Inf(1)
			for i := 0; i+1 < len(zeros); i++ {
				if g := zeros[i+1] - zeros[i]; g < bg2 {
					bg2, bi2 = g, i
				}
			}
			if len(slack) > 0 {
				midP := (zeros[bi2] + zeros[bi2+1]) / 2
				best := math.Inf(1)
				for _, s := range slack {
					if d := math.Abs(midP - s); d < best {
						best = d
					}
				}
				fmt.Printf("  agua muerta: %d turning points of the predicted tide; closest pair (gap %.2f spc) sits %.2f spc from the nearest\n",
					len(slack), bg2/spcA, best/spcA)
			}
		}
		if delta <= -0.5 && math.Min(nearL, nearR) < 0.5*spcA {
			fmt.Println("  << FRAME CLIP: the deficit's zero swims just past the edge - the count")
			fmt.Println("     closes with a breathing frame; the land was framed out, not absent.")
			logLine(fmt.Sprintf("- GUARDIÁN DEL MARCO en t=%.6g: déficit %+.2f pero cero a %s esp del borde — el marco recortó tierra",
				*anchor, delta, fmtSide(math.Min(nearL, nearR))))
		}
		// the antigravity clearance (F101): harmonic repulsion makes hidden
		// pairs cubically rare - each beach carries its certificate.
		clear60 := exact * 1.08 * math.Pow(1.0/60, 3)
		clear600 := exact * 1.08 * math.Pow(1.0/600, 3)
		fmt.Printf("  antigravity clearance (F101): P(hidden pair under the sweep) ~ %.1e (%.1e under the glass)\n",
			clear60, clear600)
		// EL OJO on board (F118): every anchorage X-rays its own interior -
		// a mid-window surge is invisible to the boundary delta.
		nB := 24
		eyeS, eyeU := 0.0, bow
		exQ := func(u float64) float64 { return (u*sp.thPrime + u*u/(4*sp.t0)) / math.Pi }
		for bq := 1; bq <= nB; bq++ {
			u := bow + span*float64(bq)/float64(nB)
			cnt := 0
			for _, z := range zeros {
				if z <= u {
					cnt++
				}
			}
			sRel := float64(cnt) - (exQ(u) - exQ(bow))
			if math.Abs(sRel) > math.Abs(eyeS) {
				eyeS, eyeU = sRel, u
			}
		}
		fmt.Printf("  el ojo (F118): interior surge max S = %+.2f at u = %.3f\n", eyeS, eyeU-bow)
		logLine(fmt.Sprintf("- el ojo (F118): marejada interior máxima S = %+.2f en u = %.3f", eyeS, eyeU))
		// EL RESORTE (F125, standard equipment): the strongest well's
		// anatomy - depth, half-depth width, curvature - so every
		// anchorage feeds the compression ledger (twin wells, sawtooth
		// asymmetry) as the fleet sails.
		if math.Abs(eyeS) >= 1.5 {
			half := math.Abs(eyeS) / 2
			edgeW := func(dir float64) float64 {
				u := eyeU
				for u > bow && u < bow+span {
					cnt := 0
					for _, z := range zeros {
						if z <= u {
							cnt++
						}
					}
					if math.Abs(float64(cnt)-(exQ(u)-exQ(bow))) <= half {
						break
					}
					u += dir * span / 2000
				}
				return u
			}
			eLw, eRw := edgeW(-1), edgeW(+1)
			wSpc := (eRw - eLw) / spcA
			if wSpc > 0 {
				fmt.Printf("  el resorte (F125): well S = %+.2f, width %.2f spc, k = %.2f\n",
					eyeS, wSpc, math.Abs(eyeS)/((wSpc/2)*(wSpc/2)))
			}
			// EL AGUA MUDA (F127, standard equipment): the amplitude
			// mountain inside the well. A deep deficit running hushed -
			// Z grazing the axis without crossing - is the water the
			// Hypothesis lives in; shout it and log it.
			mount := 0.0
			for i := 0; i <= 400; i++ {
				u := eLw + (eRw-eLw)*float64(i)/400
				if v := math.Abs(sp.zAt(u)); v > mount {
					mount = v
				}
			}
			fmt.Printf("  el agua muda (F127): max|Z| inside the well = %.3f\n", mount)
			if eyeS <= -2 && mount < 1.0 {
				fmt.Println("  << AGUA MUDA PROFUNDA: a deep deficit running silent - near-miss")
				fmt.Println("     audit recommended; this is the water the Hypothesis lives in.")
				logLine(fmt.Sprintf("- AGUA MUDA PROFUNDA en t=%.6g: pozo %+.2f con max|Z|=%.3f — auditoría de casi-cruces recomendada",
					*anchor, eyeS, mount))
			}
		}
		// EL CIRCUITO (F125 addendum): constant voltage, growing current.
		fmt.Printf("  el circuito (F125): V = %+.2f (bounded at every height), I = %.3f zeros/unit (grows with depth), R ~ %.3f\n",
			delta, 1/spcA, math.Abs(delta)*spcA)

		// the crystal gauge (F105): crystallinity of the window - crystal
		// water is certified calm; turbulent water got the automatic glass.
		if len(zeros) >= 3 {
			spacingA := twoPi.hi / math.Log(*anchor/twoPi.hi)
			mean, minG := 0.0, math.Inf(1)
			for i := 0; i+1 < len(zeros); i++ {
				g := zeros[i+1] - zeros[i]
				mean += g
				if g < minG {
					minG = g
				}
			}
			mean /= float64(len(zeros) - 1)
			vv := 0.0
			for i := 0; i+1 < len(zeros); i++ {
				g := zeros[i+1] - zeros[i]
				vv += (g - mean) * (g - mean)
			}
			cvv := math.Sqrt(vv/float64(len(zeros)-1)) / mean
			state := "normal"
			if cvv <= 0.15 {
				state = "CRISTAL - agua en red, sin escondites"
			} else if cvv >= 0.5 {
				state = "TURBULENTA - la lupa ya pasó de oficio"
			}
			fmt.Printf("  crystal gauge (F105): CV %.0f%% (%s); narrowest pass %.1f sweep steps\n",
				100*cvv, state, minG/(spacingA/60))
			logLine(fmt.Sprintf("- medidor de cristal (F105): CV %.0f%% (%s); pase más angosto %.1f pasos de barrido",
				100*cvv, state, minG/(spacingA/60)))
		}
		// the light archive: the hour of sea, kept forever in kilobytes.
		saveCkpt(archivePath(*anchor), lightArchive{
			T0: sp.t0, Span: span, H: sp.h, Dt0: sp.dt0, C0: sp.c0term,
			ThA: sp.thAnchor, ThP: sp.thPrime, NTop: sp.nTop, Fr: sp.fr, Fi: sp.fi,
			Bow: bow,
		})
		fmt.Printf("  light archived: luz/luz-%.6g.gob (%d points) - infinite replay from now on\n",
			*anchor, len(sp.fr))
		// the black box: every anchorage survives any closed panel.
		logLine("")
		logLine(fmt.Sprintf("### Nave espacial fondeada en t = %.6g — %s",
			*anchor, time.Now().Format("2006-01-02 15:04")))
		logLine(fmt.Sprintf("- primer cero de la ventana: cero número ~%s (±2)", zeroIndex(*anchor)))
		line := "- ceros vírgenes (offsets):"
		for _, z := range zeros {
			line += fmt.Sprintf(" %.6f", z)
		}
		logLine(line)
		verdict := "esfera OK"
		if math.Abs(delta) > tolS(*anchor) {
			verdict = "ESFERA ROTA: investigar antes de confiar"
		}
		logLine(fmt.Sprintf("- esfera: exige %.2f, hallados %d (delta %+.2f) — %s; marea S = %+.2f",
			exact, len(zeros), delta, verdict, delta))
		logLine(fmt.Sprintf("- marcapasos (F100): marea prevista %+.2f, residuo tras pronóstico %+.2f",
			pred, delta-pred))
		logLine(fmt.Sprintf("- despeje antigravitatorio (F101): P(par oculto) ~ %.1e bajo el barrido, %.1e bajo la lupa",
			clear60, clear600))
		logLine(fmt.Sprintf("- instrumento: %d facetas + %d bloques plegados (borde %.0f); %.1f min",
			len(sp.facets), sp.nBlocks, edge, mins))
	}
}
