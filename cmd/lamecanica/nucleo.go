// Command porosidad tests the captain's Phase VIII flash (Auditoria/41):
//
//	"y si en este sistema hay cierta porosidad: mayor o menor porosidad,
//	 materiales mas blandos y otros mas duros?"
//
// Two points separated by the same distance would then NOT interact the same
// way: the interaction would depend on the local material, not only on |i-j|.
//
// THE TRANSLATION (her §7 demands we define it). Each site i carries a local
// permeability h_i > 0, and the kernel is
//
//	C(i,j) = sqrt(h_i * h_j) * f(|i-j|)
//
// which keeps the matrix symmetric and makes "porous" mean exactly: a site whose
// influence passes easily in and out. Hard = small h, soft = large h. Porosity
// and hardness are NOT assumed to be two variables: they are one, read in two
// directions, which is what her §7 asked us to decide.
//
// THE THREE ARMS (her §6), and the third is the whole point:
//
//	A - HOMOGENEOUS: h_i = 1. The Phase VII long-tail baseline.
//	B - SHUFFLED: the SAME multiset of h values, randomly permuted across sites.
//	    Same statistics, same spread, same everything - except the arrangement.
//	C - ARITHMETIC: h_i in its true arithmetic order.
//
// B is the sharp control: if C beats B, then what matters is the ARRANGEMENT,
// which is the only thing that could be arithmetic. If C and B tie, the gain is
// texture and has nothing to do with the primes.
//
// R6: the permeability is declared before any spectrum is computed, from
// Lambda(p) and the mode's own prime only. No gamma_n anywhere except as a ruler.
//
// Reproduce: go run ./cmd/porosidad
package main

// nucleo.go - the numerical core, inherited UNCHANGED from Phase VIII (cmd/porosidad)
// so that the identity check can reproduce its table bit for bit. The single
// addition is the bond-sign field in espectro, which is force-neutral by sg^2 = 1.

import (
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// The medium
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

type modo struct {
	w, amp float64
	p, k   int
}

func medio(topeP, n int, t0 float64) []modo {
	var ms []modo
	for _, p := range primos(topeP) {
		lp := math.Log(float64(p))
		amp := lp / math.Sqrt(float64(p))
		for k := 1; ; k++ {
			w := 2 * math.Pi * float64(k) / lp
			if w > t0*40 {
				break
			}
			if w >= t0 {
				ms = append(ms, modo{w, amp, p, k})
			}
		}
	}
	sort.Slice(ms, func(a, b int) bool { return ms[a].w < ms[b].w })
	if len(ms) > n {
		ms = ms[:n]
	}
	return ms
}

// ---------------------------------------------------------------------------
// The permeability field: the captain's porosity, made mathematical
// ---------------------------------------------------------------------------

// permeabilidad returns the local h_i, normalised to mean 1 so that every arm
// carries the same total material and only its DISTRIBUTION differs.
//
// The arithmetic rule, declared before any spectrum: a mode that belongs to a
// SMALL prime is a long, slack wave in the medium and lets influence through;
// one that belongs to a large prime is short and stiff. So h_i = 1/log(p_i).
// Nothing here looks at a zero.
func permeabilidad(ms []modo) []float64 {
	h := make([]float64, len(ms))
	s := 0.0
	for i, m := range ms {
		h[i] = 1 / math.Log(float64(m.p))
		s += h[i]
	}
	f := float64(len(ms)) / s
	for i := range h {
		h[i] *= f
	}
	return h
}

// permeabilidadDivisores is a second arithmetic rule, deliberately UNCORRELATED
// with the mode's amplitude: the hardness of a site is read off the divisor
// count of its harmonic index. It exists so the sheet does not rest on a single
// arithmetic choice.
func permeabilidadDivisores(ms []modo) []float64 {
	h := make([]float64, len(ms))
	s := 0.0
	for i, m := range ms {
		d := 0
		for q := 1; q*q <= m.k; q++ {
			if m.k%q == 0 {
				d += 2
				if q*q == m.k {
					d--
				}
			}
		}
		h[i] = float64(d)
		s += h[i]
	}
	f := float64(len(ms)) / s
	for i := range h {
		h[i] *= f
	}
	return h
}

// ordenadaNoAritmetica is the control this sheet needed and did not have: a
// field that is just as ORDERED as the arithmetic one - a smooth ramp, and a
// smooth wave - but carries no arithmetic whatsoever. If these match the
// arithmetic arms, then what the medium is responding to is the ORDER of the
// field and NOT the primes, and the arithmetic claim dies right here.
func ordenadaRampa(n int, disp float64) []float64 {
	h := make([]float64, n)
	for i := range h {
		h[i] = 1 + disp*(float64(i)/float64(n-1)-0.5)
	}
	return normalizar1(h)
}

func ordenadaOnda(n int, disp float64) []float64 {
	h := make([]float64, n)
	for i := range h {
		h[i] = 1 + disp*math.Sin(6*math.Pi*float64(i)/float64(n-1))
	}
	return normalizar1(h)
}

// normalizar1 forces mean 1 and clamps positivity, so every arm carries the
// same total material.
func normalizar1(h []float64) []float64 {
	for i := range h {
		if h[i] < 1e-6 {
			h[i] = 1e-6
		}
	}
	s := 0.0
	for _, x := range h {
		s += x
	}
	f := float64(len(h)) / s
	for i := range h {
		h[i] *= f
	}
	return h
}

// homogenea is arm A: no heterogeneity at all.
func homogenea(n int) []float64 {
	h := make([]float64, n)
	for i := range h {
		h[i] = 1
	}
	return h
}

// mezclar is arm B: the SAME values, permuted. Same statistics, no arrangement.
func mezclar(h []float64, d *dado) []float64 {
	out := append([]float64(nil), h...)
	for i := len(out) - 1; i > 0; i-- {
		j := int(d.u() * float64(i+1))
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// bloques is her §8 localisation medium: an explicit hard zone, soft zone and
// middling zone, with the same mean permeability as every other arm.
func bloques(n int, contraste float64) []float64 {
	h := make([]float64, n)
	for i := range h {
		switch {
		case i < n/3:
			h[i] = contraste // soft
		case i < 2*n/3:
			h[i] = 1
		default:
			h[i] = 1 / contraste // hard
		}
	}
	s := 0.0
	for _, x := range h {
		s += x
	}
	f := float64(n) / s
	for i := range h {
		h[i] *= f
	}
	return h
}

type dado struct{ s uint64 }

func (d *dado) u() float64 {
	d.s ^= d.s << 13
	d.s ^= d.s >> 7
	d.s ^= d.s << 17
	return float64(d.s>>11) / float64(uint64(1)<<53)
}

// ---------------------------------------------------------------------------
// The kernels
// ---------------------------------------------------------------------------

// colaLarga is the Phase VII survivor: a power-law tail, no node.
func colaLarga(s float64) func(int) float64 {
	return func(k int) float64 {
		if k == 0 {
			return 0
		}
		return 1 / math.Pow(float64(k), s)
	}
}

// conNodo is the captain's hourglass: the same tail with a node at k0.
func conNodo(k0, s float64) func(int) float64 {
	return func(k int) float64 {
		if k == 0 {
			return 0
		}
		fk := float64(k)
		return (1 - fk/k0) / math.Pow(fk, s)
	}
}

// ---------------------------------------------------------------------------
// The spectrum, with eigenvectors so localisation can be measured
// ---------------------------------------------------------------------------

type espectroRes struct {
	niveles []float64
	prMedio float64 // mean participation ratio / N: 1 = extended, ~0 = localised
	vivos   int
}

// espectro takes an optional BOND SIGN field sg(i,j) in {+1,-1}. Because sg^2 = 1 the
// total force sum_{i!=j} H_ij^2 is unchanged, so the F normalisation cannot move: only
// the MATERIAL changes. H stays real symmetric, so R1 (self-adjointness) is intact.
func espectro(ms []modo, h []float64, c func(int) float64, kmax int, fuerzaObj float64, sg func(int, int) float64) espectroRes {
	n := len(ms)
	H := make([][]float64, n)
	V := make([][]float64, n)
	for i := range H {
		H[i] = make([]float64, n)
		V[i] = make([]float64, n)
		V[i][i] = 1
	}
	// first pass: build the perturbation and measure its own force
	P := make([][]float64, n)
	for i := range P {
		P[i] = make([]float64, n)
	}
	f2 := 0.0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n && j-i <= kmax; j++ {
			x := math.Sqrt(h[i]*h[j]) * c(j-i) * ms[i].amp * ms[j].amp
			if sg != nil {
				x *= sg(i, j)
			}
			P[i][j], P[j][i] = x, x
			f2 += 2 * x * x
		}
	}
	esc := 1.0
	if f := math.Sqrt(f2); f > 0 {
		esc = fuerzaObj / f
	}
	for i := range ms {
		H[i][i] = ms[i].w
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			H[i][j] += esc * P[i][j]
			H[j][i] += esc * P[j][i]
		}
	}
	jacobiV(H, V, 20)
	// pair each eigenvalue with its vector's participation ratio
	type par struct {
		e, pr float64
	}
	pares := make([]par, n)
	for c2 := 0; c2 < n; c2++ {
		s4 := 0.0
		for r := 0; r < n; r++ {
			v2 := V[r][c2] * V[r][c2]
			s4 += v2 * v2
		}
		pr := 0.0
		if s4 > 0 {
			pr = 1 / s4 / float64(n)
		}
		pares[c2] = par{H[c2][c2], pr}
	}
	sort.Slice(pares, func(a, b int) bool { return pares[a].e < pares[b].e })
	lo, hi := ms[0].w, ms[len(ms)-1].w
	var niv []float64
	prSum, prN := 0.0, 0
	for _, p := range pares {
		if p.e >= lo && p.e <= hi {
			niv = append(niv, p.e)
			prSum += p.pr
			prN++
		}
	}
	pr := 0.0
	if prN > 0 {
		pr = prSum / float64(prN)
	}
	return espectroRes{niv, pr, len(niv)}
}

// jacobiV diagonalises a in place and accumulates the rotations in v, so the
// eigenvectors are available and localisation can be measured.
func jacobiV(a, v [][]float64, barridos int) {
	n := len(a)
	for sw := 0; sw < barridos; sw++ {
		off := 0.0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				off += a[i][j] * a[i][j]
			}
		}
		if off < 1e-20 {
			return
		}
		for p := 0; p < n-1; p++ {
			for q := p + 1; q < n; q++ {
				if math.Abs(a[p][q]) < 1e-16 {
					continue
				}
				th := (a[q][q] - a[p][p]) / (2 * a[p][q])
				t := math.Copysign(1, th) / (math.Abs(th) + math.Sqrt(th*th+1))
				cc := 1 / math.Sqrt(t*t+1)
				s := t * cc
				for k := 0; k < n; k++ {
					akp, akq := a[k][p], a[k][q]
					a[k][p] = cc*akp - s*akq
					a[k][q] = s*akp + cc*akq
				}
				for k := 0; k < n; k++ {
					apk, aqk := a[p][k], a[q][k]
					a[p][k] = cc*apk - s*aqk
					a[q][k] = s*apk + cc*aqk
				}
				for k := 0; k < n; k++ {
					vkp, vkq := v[k][p], v[k][q]
					v[k][p] = cc*vkp - s*vkq
					v[k][q] = s*vkp + cc*vkq
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// The bench
// ---------------------------------------------------------------------------

func media(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func desvio(v []float64) float64 {
	if len(v) < 2 {
		return 0
	}
	m := media(v)
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return math.Sqrt(s / float64(len(v)-1))
}

func desplegarPropio(niv []float64) []float64 {
	if len(niv) < 3 {
		return nil
	}
	d := float64(len(niv)-1) / (niv[len(niv)-1] - niv[0])
	sp := make([]float64, 0, len(niv))
	for i := 0; i+1 < len(niv); i++ {
		sp = append(sp, (niv[i+1]-niv[i])*d)
	}
	return sp
}

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
	vv := 0.0
	for _, y := range ns {
		vv += (y - m) * (y - m)
	}
	return vv / float64(len(ns))
}

func repulsion(sp []float64) (frac, minimo float64) {
	minimo = math.Inf(1)
	n := 0
	for _, s := range sp {
		if s < minimo {
			minimo = s
		}
		if s < 0.1 {
			n++
		}
	}
	return float64(n) / float64(len(sp)), minimo
}

// obs is the full observable set, with the discipline Phase VII forced on us:
// nothing is reported unless enough levels stayed inside the band.
type obs struct {
	s5, s10, s20, alfa, frac, min, pr float64
	vivos                             int
	valido                            bool
}

func medir(r espectroRes) obs {
	sp := desplegarPropio(r.niveles)
	if len(sp) < 60 {
		return obs{vivos: r.vivos, pr: r.prMedio}
	}
	s5, s10, s20 := sigma2(sp, 5, 400), sigma2(sp, 10, 400), sigma2(sp, 20, 400)
	fr, mn := repulsion(sp)
	return obs{s5, s10, s20, math.Log(s20/s5) / math.Log(4), fr, mn, r.prMedio, r.vivos, true}
}
