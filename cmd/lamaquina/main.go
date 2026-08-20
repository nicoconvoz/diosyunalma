// Command lamaquina answers the auditor's Phase III route (Auditoria/36): stop
// looking for a list that imitates the zeros, look for a MACHINE that does not
// need to know them. Her §5 fixes the first step - "de primo a orbita": say what
// object a p^k orbit is, why its period is log(p^k), and DERIVE its weight and
// its sign instead of imposing them.
//
// This sheet does exactly that, and the answer is mostly negative - which is
// the useful kind:
//
//	§2  THE CONCATENATION OBSTRUCTION (new, with its exact hypotheses). If the
//	    primes are closed orbits of a system where orbits can be concatenated -
//	    any graph, any flow with a recurrent point - then the length spectrum is
//	    closed under addition, so it contains log n for EVERY integer n. The Weil
//	    explicit formula gives weight exactly ZERO to every n that is not a prime
//	    power. Measured: how many lengths a concatenating system is forced to
//	    carry that the zeros forbid.
//	§3  THE DISJOINT LADDERS. The only escape from §2 is orbits that never meet.
//	    Built and measured: density, correlations, echo, identity.
//	§4  THE WEIGHT, DERIVED. A circle gives weight log p with no damping. Asking
//	    for the p^(-m/2) of the explicit formula FORCES an unstable orbit with
//	    instability exponent equal to its own length - Lyapunov exponent exactly
//	    1, which is the xp flow. That is a derivation, not a choice.
//	§5  THE SIGN, and what is left over after the amplitudes are matched.
//
// R6: the only arithmetic input is Lambda(n). No zero is read except as a ruler
// in the section that explicitly measures distance to them.
//
// Reproduce: go run ./cmd/lamaquina
package main

import (
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// The admissible input: the primes and nothing else
// ---------------------------------------------------------------------------

func primos(tope int) []int {
	criba := make([]bool, tope+1)
	var ps []int
	for p := 2; p <= tope; p++ {
		if criba[p] {
			continue
		}
		ps = append(ps, p)
		for m := 2 * p; m <= tope; m += p {
			criba[m] = true
		}
	}
	return ps
}

// lambda is von Mangoldt: log p if n = p^k, else 0. This is the ENTIRE
// arithmetic content the explicit formula gives to the period log n.
func lambda(n int) float64 {
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			m := n
			for m%p == 0 {
				m /= p
			}
			if m == 1 {
				return math.Log(float64(p))
			}
			return 0
		}
	}
	if n > 1 {
		return math.Log(float64(n))
	}
	return 0
}

// ---------------------------------------------------------------------------
// The ruler: the true zeros. No candidate is allowed to look at these.
// ---------------------------------------------------------------------------

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t*t)
}

func suave(T float64) float64 { return theta(T)/math.Pi + 1 }

func zetaZ(t float64) float64 {
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if N%2 == 0 {
		sign = -1
	}
	return s + sign*math.Pow(2*math.Pi/t, 0.25)*c0
}

func cerosVerdaderos(t0, t1, paso float64) []float64 {
	var g []float64
	a, za := t0, zetaZ(t0)
	for b := t0 + paso; b <= t1; b += paso {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 60; i++ {
				m := (lo + hi) / 2
				zm := zetaZ(m)
				if (zlo < 0) != (zm < 0) {
					hi = m
				} else {
					lo, zlo = m, zm
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

// ---------------------------------------------------------------------------
// The measuring bench (independent of the Phase II one)
// ---------------------------------------------------------------------------

func media(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func desvio(v []float64) float64 {
	m := media(v)
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return math.Sqrt(s / float64(len(v)-1))
}

// sigma2 is the number variance of an unfolded ladder.
func sigma2(sp []float64, L float64, pasos int) float64 {
	x := make([]float64, len(sp)+1)
	for i, s := range sp {
		x[i+1] = x[i] + s
	}
	total := x[len(x)-1]
	if total <= L*1.2 {
		return math.NaN()
	}
	var ns []float64
	paso := (total - L) / float64(pasos)
	for k := 0; k < pasos; k++ {
		a := float64(k) * paso
		lo := sort.SearchFloat64s(x, a)
		hi := sort.SearchFloat64s(x, a+L)
		ns = append(ns, float64(hi-lo))
	}
	m := media(ns)
	v := 0.0
	for _, y := range ns {
		v += (y - m) * (y - m)
	}
	return v / float64(len(ns))
}

func eco(niv []float64, T float64) float64 {
	M := len(niv)
	s := 0.0
	for i, g := range niv {
		w := 0.5 - 0.5*math.Cos(2*math.Pi*(float64(i)+0.5)/float64(M))
		s += w * math.Cos(g*T)
	}
	return 2 * s / float64(M)
}

func murcielago(niv, per []float64, tope float64) float64 {
	const ancho = 0.012
	var dentro, fuera float64
	nIn, nOut := 0, 0
	for T := 0.35; T <= tope; T += 0.0007 {
		cerca := false
		for _, p := range per {
			if math.Abs(T-p) < ancho {
				cerca = true
				break
			}
		}
		v := math.Abs(eco(niv, T))
		if cerca {
			dentro += v
			nIn++
		} else {
			fuera += v
			nOut++
		}
	}
	return (dentro / float64(nIn)) / (fuera / float64(nOut))
}
