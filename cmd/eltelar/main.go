// Command eltelar answers the auditor's Phase II question (Auditoria/35, §15):
//
//	can we build, from data admissible under R6, a structure whose effective
//	geometry grows like ln(E/2pi) and whose spectrum echoes at k*log p, without
//	introducing the gamma_n at any point of the definition?
//
// It splits that question into two halves that can be measured separately, and
// measures both against the same protocol the auditor fixed in her §4:
//
//	§2  THE BREATHING BOX ALONE. Levels from the smooth law only, whose effective
//	    log-length is exactly ln(E/2pi). It gets the counting right by
//	    construction. The question is what else it gets.
//	§3  THE PRIMES ALONE. Levels from the smooth law PLUS the oscillating term of
//	    the explicit formula, truncated to prime powers - no zero is ever read.
//	    Then: how close to the true gamma_n, and how does that depend on how many
//	    primes we are allowed?
//
// Every candidate carries its input list (§1, the R6 audit) and every spectrum
// is judged by the same four layers the auditor separated: counting,
// correlations, identity, arithmetic.
//
// Reproduce: go run ./cmd/eltelar
package main

import (
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// THE ADMISSIBLE DATA - everything here comes from the functional equation and
// the primes. No zero is read anywhere in this file except in the section that
// explicitly measures the distance to them.
// ---------------------------------------------------------------------------

// theta is the Riemann-Siegel phase: it comes from the functional equation, not
// from the zeros.
func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t*t)
}

// suave is the smooth counting law N(T) = theta(T)/pi + 1. Its derivative is
// (1/2pi)*ln(T/2pi): the density of a box whose log-length is ln(T/2pi). THIS
// is the breathing box, and it is admissible: it descends from the functional
// equation, not from where the zeros sit.
func suave(T float64) float64 { return theta(T)/math.Pi + 1 }

// largoCaja is the effective log-length the smooth law implies at height T.
func largoCaja(T float64) float64 { return math.Log(T / (2 * math.Pi)) }

// mangoldt returns the prime powers n = p^k <= tope with their weights
// Lambda(n) = log p. This is the whole arithmetic input of the candidate.
func mangoldt(tope int) (ns []float64, lam []float64) {
	criba := make([]bool, tope+1)
	for p := 2; p <= tope; p++ {
		if criba[p] {
			continue
		}
		for m := 2 * p; m <= tope; m += p {
			criba[m] = true
		}
		lp := math.Log(float64(p))
		for pk := p; pk <= tope; {
			ns = append(ns, float64(pk))
			lam = append(lam, lp)
			if pk > tope/p {
				break
			}
			pk *= p
		}
	}
	idx := make([]int, len(ns))
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool { return ns[idx[a]] < ns[idx[b]] })
	on, ol := make([]float64, len(ns)), make([]float64, len(ns))
	for i, j := range idx {
		on[i], ol[i] = ns[j], lam[j]
	}
	return on, ol
}

// sPrimos is the oscillating term of the counting function, heard through the
// primes only:  S(t) = -(1/pi) * sum_{n=p^k} Lambda(n) sin(t log n)/(sqrt(n) log n).
// It is the arithmetic half of Riemann-von Mangoldt, and it reads no zero.
func sPrimos(t float64, ns, lam []float64) float64 {
	s := 0.0
	for i, n := range ns {
		ln := math.Log(n)
		s += lam[i] * math.Sin(t*ln) / (math.Sqrt(n) * ln)
	}
	return -s / math.Pi
}

// ---------------------------------------------------------------------------
// THE TRUE SEA - used ONLY to measure how close a candidate gets. No candidate
// is allowed to look at this.
// ---------------------------------------------------------------------------

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
// BUILDING A SPECTRUM FROM A COUNTING FUNCTION
// ---------------------------------------------------------------------------

// nivelesDe solves cuenta(T) = j - 1/2 for consecutive j: the level ladder a
// counting function implies. Bisection, so a non-monotonic count is handled by
// taking the first crossing.
func nivelesDe(cuenta func(float64) float64, t0, t1 float64) []float64 {
	var out []float64
	n0 := cuenta(t0)
	j := math.Ceil(n0 + 0.5)
	for {
		obj := j - 0.5
		lo, hi := t0, t1
		if cuenta(hi) < obj {
			break
		}
		for i := 0; i < 80; i++ {
			m := (lo + hi) / 2
			if cuenta(m) < obj {
				lo = m
			} else {
				hi = m
			}
		}
		out = append(out, (lo+hi)/2)
		j++
	}
	return out
}

// ---------------------------------------------------------------------------
// THE PROTOCOL (the auditor's §4 and §5): four layers, one verdict each
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

// desplegar turns raw levels into unfolded spacings with the smooth law.
func desplegar(niv []float64) []float64 {
	sp := make([]float64, 0, len(niv))
	for i := 0; i+1 < len(niv); i++ {
		sp = append(sp, suave(niv[i+1])-suave(niv[i]))
	}
	return sp
}

// sigma2 is the number variance: the two-level statistic that Phase I showed is
// NOT blind where the one-spacing exam is.
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

// eco is the auditor's bat: E(T) = (2/M) sum_n w_n cos(gamma_n T), raised cosine.
func eco(niv []float64, T float64) float64 {
	M := len(niv)
	s := 0.0
	for i, g := range niv {
		w := 0.5 - 0.5*math.Cos(2*math.Pi*(float64(i)+0.5)/float64(M))
		s += w * math.Cos(g*T)
	}
	return 2 * s / float64(M)
}

func periodosAritmeticos(tope float64) []float64 {
	ns, _ := mangoldt(int(math.Exp(tope)) + 2)
	var ps []float64
	for _, n := range ns {
		if l := math.Log(n); l <= tope && l > 0.3 {
			ps = append(ps, l)
		}
	}
	return ps
}

// murcielago scores a spectrum: mean |E| on the given periods divided by mean
// |E| off them.
func murcielago(niv, per []float64, tope float64) (dentro, fuera, razon float64) {
	const ancho = 0.012
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
	return dentro / float64(nIn), fuera / float64(nOut), (dentro / float64(nIn)) / (fuera / float64(nOut))
}

type dado struct{ s uint64 }

func (d *dado) u() float64 {
	d.s ^= d.s << 13
	d.s ^= d.s >> 7
	d.s ^= d.s << 17
	return float64(d.s>>11) / float64(uint64(1)<<53)
}

// periodosAlAzar draws a control set of non-arithmetic periods.
func periodosAlAzar(d *dado, n int, per []float64, tope float64) []float64 {
	var out []float64
	for len(out) < n {
		T := 0.4 + d.u()*(tope-0.45)
		ok := true
		for _, q := range per {
			if math.Abs(T-q) < 0.05 {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, T)
		}
	}
	sort.Float64s(out)
	return out
}
