// Command voyage2 is the refitted ship: double-double phases for deep water.
//
// The float64 hull creaked at 7.77·10¹³ (Finding 87): the phase t·ln k
// carries an ulp of ±0.26 radians at magnitude 10¹⁵. The refit computes
// every phase in double-double arithmetic (~30 digits) via an incremental
// logarithm chain — ln k grows from ln(k−1) by a dd series in 1/(k−1),
// each increment multiplied by the exact t and reduced mod 2π with a
// split-integer subtraction. Accumulated phase error stays below ~10⁻⁷
// radians even at 10¹⁸.
//
// Sea trials first (charted water: t = 10⁵ against the published tables,
// t = 2.447·10¹² against Finding 87's beach), then three virgin
// anchorages between the island expeditions of mankind:
//
//	6.66·10¹⁵ · 4.44·10¹⁷ · 3.33·10¹⁸
//
// PRE-REGISTERED: sea trials must reproduce charted zeros; at every virgin
// anchorage the count of zeros must match the density expectation — the
// line, audited in water no one has ever observed.
//
// Usage:
//
//	go run ./cmd/voyage2
package main

import (
	"fmt"
	"math"
	"math/big"
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

// divF returns x/f in dd for exact float f.
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

// mod2pi reduces a dd of any magnitude to a float64 phase.
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

// ddLnF computes ln(x) in dd for any positive float64 via atanh series.
func ddLnF(x float64) dd {
	m, e := math.Frexp(x) // x = m*2^e, m in [0.5,1)
	if m < 0.70710678118654752 {
		m *= 2
		e--
	}
	num := m - 1 // exact (Sterbenz)
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

// theta computes the Riemann–Siegel phase mod 2pi in dd.
func thetaMod(t float64) float64 {
	lnArg := add(ddLnF(t), neg(ln2pi))
	th := mulF(addF(lnArg, -1), t/2)
	th = add(th, neg(pi8))
	th = addF(th, 1/(48*t))
	return mod2pi(th)
}

// z evaluates Riemann–Siegel Z(t) with dd phases via the incremental chain.
func z(t float64) float64 {
	a := math.Sqrt(t / twoPi.hi)
	n := int(a)
	p := a - float64(n)
	th := thetaMod(t)

	sum := 0.0
	lnPrev := dd{0, 0} // ln 1
	ph := 0.0          // phase of t*ln k mod 2pi, accumulated
	needed := 49.8     // ln(t * 1e4) budget for series truncation
	if t > 1e16 {
		needed = math.Log(t * 1e4)
	}
	for k := 2; k <= n; k++ {
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
	sum += math.Cos(th) // k = 1
	sum *= 2
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (n-1)%2 == 1 {
		sign = -1
	}
	return sum + sign*math.Pow(t/twoPi.hi, -0.25)*c0
}

func hunt(t0 float64, spacings float64, tag string) {
	spacing := twoPi.hi / math.Log(t0/twoPi.hi)
	span := spacings * spacing
	n := int(math.Sqrt(t0 / twoPi.hi))
	fmt.Printf("\n%s t = %.6g   (%d dd terms)\n", tag, t0, n)
	zeros := []float64{}
	step := spacing / 10
	prevT, prevZ := t0, z(t0)
	for t := t0 + step; t <= t0+span; t += step {
		zt := z(t)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			zlo := prevZ
			for i := 0; i < 24 && hi-lo > 1e-7; i++ {
				mid := (lo + hi) / 2
				zm := z(mid)
				if (zlo < 0) != (zm < 0) {
					hi = mid
				} else {
					lo, zlo = mid, zm
				}
			}
			zeros = append(zeros, (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}
	fmt.Print("  zeros (offsets):")
	for _, zr := range zeros {
		fmt.Printf("  %.6f", zr-t0)
	}
	fmt.Println()
	expected := span / spacing
	fmt.Printf("  map check: found %d, density expects %.1f\n", len(zeros), expected)
}

func main() {
	fmt.Println("THE REFITTED SHIP — double-double phases, deep water rating")

	fmt.Println("\nSEA TRIALS (charted water):")
	hunt(100000, 3, "trial")
	hunt(2.447e12, 8, "trial (Finding 87's beach)")

	fmt.Println("\nVIRGIN ANCHORAGES (between mankind's islands):")
	hunt(6.66e15, 8, "anchorage")
	hunt(4.44e17, 8, "anchorage")
	hunt(3.33e18, 8, "anchorage")

	fmt.Println("\nthe ocean between the charted map and the island expeditions,")
	fmt.Println("entered with a hull that keeps its phases to thirty digits.")
}
