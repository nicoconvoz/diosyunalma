// Command hiladofino runs the two Phase VII campaigns (Auditoria/40), in the
// order her §16 demands: "primero hilar fino, despues cambiar el hilo".
//
//	CAMPAIGN A - refine the A/B phase boundary around (30,32), where Phase VI
//	found alpha = 0.313, and ask whether the transition is gradual or abrupt,
//	and whether it survives changes of size, window and discretisation.
//
//	CAMPAIGN B - the hourglass kernel born from the captain's drawing:
//	H_ij = c(|i-j|) with a NODE at k0 where c(k0) = 0 and the sign reversing
//	beyond it. Her §8 forbids assuming k0 = 5: it is a parameter.
//
// AND HER §6, WHICH IS THE POINT OF THE WHOLE THING: a sign change does NOT
// prove there is useful cancellation. So the kernel is compared, at EQUAL total
// force, against the very same |c(k)| with every sign POSITIVE. Whatever the
// node buys has to show up in that difference and nowhere else.
//
// R6: the only arithmetic input is Lambda(p)/sqrt(p). No parameter is chosen by
// looking at the gamma_n; the zeros' value is quoted once, at the end, as a ruler.
//
// Reproduce: go run ./cmd/hiladofino
package main

import (
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// The medium: one excitation per prime, its ladder, its Weil weight
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
	p      int
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
				ms = append(ms, modo{w, amp, p})
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
// The two kernel families
// ---------------------------------------------------------------------------

// nucleoAB is Phase VI's kernel: amplitude A, exponential reach B, one sign.
func nucleoAB(A, B float64) func(k int) float64 {
	return func(k int) float64 {
		if k == 0 {
			return 0
		}
		return A * math.Exp(-float64(k)/B)
	}
}

// nucleoReloj is the captain's hourglass: a power-law tail with a NODE at k0,
// positive before it and negative after. c(k0) = 0 exactly.
//
//	c(k) = A * (1 - k/k0) / k^s
func nucleoReloj(A, k0, s float64) func(k int) float64 {
	return func(k int) float64 {
		if k == 0 {
			return 0
		}
		fk := float64(k)
		return A * (1 - fk/k0) / math.Pow(fk, s)
	}
}

// sinNodo is the same kernel with every sign made POSITIVE - her §15 control.
// If the node buys anything, it must show up against THIS and nothing else.
func sinNodo(c func(int) float64) func(int) float64 {
	return func(k int) float64 { return math.Abs(c(k)) }
}

// sumaYnorma reports the two quantities her §6 demands: the plain sum of the
// kernel (does it cancel globally?) and its effective norm (how much force is
// really there?).
func sumaYnorma(c func(int) float64, kmax int) (suma, norma float64) {
	for k := 1; k <= kmax; k++ {
		v := c(k)
		suma += v
		norma += v * v
	}
	return suma, math.Sqrt(norma)
}

// ---------------------------------------------------------------------------
// The spectrum
// ---------------------------------------------------------------------------

// fuerza is the Frobenius norm of the off-diagonal perturbation: the quantity
// held equal when two kernels are compared.
func fuerza(ms []modo, c func(int) float64, kmax int) float64 {
	s := 0.0
	for i := range ms {
		for j := range ms {
			if i == j {
				continue
			}
			k := j - i
			if k < 0 {
				k = -k
			}
			if k > kmax {
				continue
			}
			x := c(k) * ms[i].amp * ms[j].amp
			s += x * x
		}
	}
	return math.Sqrt(s)
}

func espectro(ms []modo, c func(int) float64, kmax int, escala float64) []float64 {
	n := len(ms)
	H := make([][]float64, n)
	for i := range H {
		H[i] = make([]float64, n)
	}
	for i := range ms {
		H[i][i] = ms[i].w
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n && j-i <= kmax; j++ {
			x := escala * c(j-i) * ms[i].amp * ms[j].amp
			H[i][j] += x
			H[j][i] += x
		}
	}
	ev := jacobi(H, 20)
	lo, hi := ms[0].w, ms[len(ms)-1].w
	dentro := ev[:0]
	for _, e := range ev {
		if e >= lo && e <= hi {
			dentro = append(dentro, e)
		}
	}
	return dentro
}

func jacobi(a [][]float64, barridos int) []float64 {
	n := len(a)
	for sw := 0; sw < barridos; sw++ {
		off := 0.0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				off += a[i][j] * a[i][j]
			}
		}
		if off < 1e-20 {
			break
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
			}
		}
	}
	ev := make([]float64, n)
	for i := range ev {
		ev[i] = a[i][i]
	}
	sort.Float64s(ev)
	return ev
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
	v := 0.0
	for _, y := range ns {
		v += (y - m) * (y - m)
	}
	return v / float64(len(ns))
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

// medir returns the full observable set of a spectrum in one call.
type obs struct {
	s5, s10, s20, s40, alfa, frac, min float64
	vivos                              int  // how many levels stayed inside the band
	valido                             bool // enough of them to mean anything
}

// medir reports the observables AND how many levels survived inside the band.
// That count is not a detail: a strong coupling expels most states, and Sigma^2
// measured on the handful that remain is not a property of the medium. Phase VI
// did not report it, and this sheet must.
func medir(niv []float64) obs {
	sp := desplegarPropio(niv)
	if len(sp) < 60 {
		return obs{vivos: len(niv), valido: false}
	}
	s5, s10 := sigma2(sp, 5, 400), sigma2(sp, 10, 400)
	s20, s40 := sigma2(sp, 20, 400), sigma2(sp, 40, 400)
	fr, mn := repulsion(sp)
	return obs{s5, s10, s20, s40, math.Log(s20/s5) / math.Log(4), fr, mn, len(niv), true}
}
