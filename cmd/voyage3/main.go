// Command voyage3 installs the convex elements in the mirror.
//
// The design (thought before built, as ordered): the phase curve t·ln n is
// curved, and the true convex mirror — quadratic-phase Gauss-sum facets,
// self-similar under the fold — is Hiary's t^(1/3) road to 10³²⁺: the
// grand refit, registered. Tonight's safe convex elements:
//
//  1. WIDE-FIELD: the expensive double-double phase chain runs ONCE per
//     anchorage and is stored as a float32 phase field; every evaluation
//     of the hunt reuses it. (Odlyzko–Schönhage's philosophy, laptop-size.)
//  2. FEWER REFLECTIONS: coarse scan at a quarter-spacing plus
//     secant-bisection refinement — four times fewer evaluations.
//
// New hull rating ~3·10¹⁸ tonight (the 10¹⁹ gates need a larger phase
// field, registered). The fleet: sea trials in charted water, then virgin
// anchorages at 6.66·10¹⁵, 4.44·10¹⁷, 1.11·10¹⁸ and flagship 3.33·10¹⁸.
//
// Usage:
//
//	go run ./cmd/voyage3
package main

import (
	"fmt"
	"math"
	"math/big"
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

// chainPass runs the dd phase chain once at anchor t0, storing the phase
// field t0*ln k mod 2pi as float32 - the wide-field convex element.
func chainPass(t0 float64, n int) []float32 {
	field := make([]float32, n+1)
	lnPrev := dd{0, 0}
	ph := 0.0
	needed := math.Log(t0 * 1e4)
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
		ph += mod2pi(mulF(delta, t0))
		if ph >= twoPi.hi {
			ph -= twoPi.hi
		}
		field[k] = float32(ph)
	}
	return field
}

// zCached evaluates Z(t0+dt) from the stored phase field.
func zCached(t0, dt float64, field []float32, n int) float64 {
	t := t0 + dt
	a := math.Sqrt(t / twoPi.hi)
	nn := int(a)
	if nn > n {
		nn = n
	}
	p := a - float64(nn)
	th := thetaMod(t)
	sum := math.Cos(th) // k=1
	for k := 2; k <= nn; k++ {
		phase := float64(field[k]) + dt*math.Log(float64(k))
		sum += math.Cos(th-phase) / math.Sqrt(float64(k))
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
	nAtTop := int(math.Sqrt((t0 + span) / twoPi.hi))
	fmt.Printf("\n%s t = %.6g   (%d dd terms; wide-field pass...)\n", tag, t0, nAtTop)
	field := chainPass(t0, nAtTop)
	eval := func(dt float64) float64 { return zCached(t0, dt, field, nAtTop) }

	// coarse scan + secant-bisection refinement: the fewer-reflections element.
	zeros := []float64{}
	step := spacing / 4
	prevD, prevZ := 0.0, eval(0)
	for d := step; d <= span; d += step {
		zd := eval(d)
		if (prevZ < 0) != (zd < 0) {
			lo, hi := prevD, d
			zlo, zhi := prevZ, zd
			for i := 0; i < 12 && hi-lo > 1e-6; i++ {
				mid := lo - zlo*(hi-lo)/(zhi-zlo) // secant
				if mid <= lo || mid >= hi {
					mid = (lo + hi) / 2
				}
				zm := eval(mid)
				if (zlo < 0) != (zm < 0) {
					hi, zhi = mid, zm
				} else {
					lo, zlo = mid, zm
				}
			}
			zeros = append(zeros, (lo+hi)/2)
		}
		prevD, prevZ = d, zd
	}
	fmt.Print("  zeros (offsets):")
	for _, zr := range zeros {
		fmt.Printf("  %.6f", zr)
	}
	fmt.Println()
	expected := span / spacing
	fmt.Printf("  map check: found %d, density expects %.1f   (%.1f min)\n",
		len(zeros), expected, time.Since(start).Minutes())
}

func main() {
	fmt.Println("THE CONVEX MIRROR — wide-field phase cache + fewer reflections")

	fmt.Println("\nSEA TRIALS (charted water):")
	hunt(100000, 4, "trial")
	hunt(2.447e12, 8, "trial (Finding 87's beach)")

	fmt.Println("\nVIRGIN ANCHORAGES:")
	hunt(6.66e15, 8, "anchorage")
	hunt(4.44e17, 8, "anchorage")
	hunt(1.11e18, 6, "anchorage")
	hunt(3.33e18, 5, "FLAGSHIP")

	fmt.Println("\ndeep water reached with a mirror that runs its expensive optics")
	fmt.Println("once and reflects the whole field from memory. the 10^19 gates of")
	fmt.Println("Odlyzko's island await a bigger phase field - registered.")
}
