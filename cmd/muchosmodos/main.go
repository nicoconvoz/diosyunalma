// Command muchosmodos runs the Phase V experiment the auditor asked for
// (Auditoria/38), which is the prediction our own Phase IV act left written:
//
//	H = D + sum_{a=1..K} g_a |v_a><v_a|
//
// and measure how Sigma^2(L) moves as K grows. Phase IV showed that ONE common
// mode makes real short-range repulsion (minimum spacing x270) but saturates,
// because a rank-one coupling interlaces: exactly one root between consecutive
// levels, so no level can be pushed past its immediate neighbour.
//
// The structural reason to expect K to matter is the same theorem read forward:
// for a rank-K positive perturbation the eigenvalues satisfy w_i <= E_i <=
// w_{i+K}. THE RANK IS LITERALLY HOW FAR A LEVEL CAN BE PUSHED. So the question
// "how many channels does the fluid need" is the question "how far must the
// waves feel each other".
//
// THE CHANNELS. They must be arithmetic and not fitted. We take the Dirichlet
// characters modulo a prime q: the group is cyclic of order q-1 = K, so channel
// a weights the excitation of prime p by cos(2*pi*a*ind(p)/(q-1)), where ind is
// the discrete logarithm to a primitive root. That is how arithmetic itself
// organises the primes into families, it is R6-clean, and a = 0 reproduces
// Phase IV's single common mode exactly.
//
// CONTROL (the auditor's §10). The same K channels with DISJOINT support - each
// channel touching only the primes of one residue class - which is K independent
// rank-one problems and no cross-talk at all. Same degrees of freedom, same
// total coupling: whatever improvement survives that comparison belongs to the
// interaction and not to the count of parameters.
//
// R6: the only arithmetic inputs are Lambda(n) and the residues mod q. The
// gamma_n are read only as a ruler, after every spectrum exists.
//
// Reproduce: go run ./cmd/muchosmodos
package main

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
	w   float64 // frequency 2*pi*k/log p
	amp float64 // Lambda(p)/sqrt(p) = log p * p^(-1/2)
	p   int
}

// medio builds the uncoupled ladders and keeps the first n modes above t0.
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
// The channels: Dirichlet characters mod a prime q, so the group is cyclic
// ---------------------------------------------------------------------------

// raizPrimitiva returns a primitive root mod the prime q.
func raizPrimitiva(q int) int {
	if q == 2 {
		return 1
	}
	var fac []int
	m := q - 1
	for d := 2; d*d <= m; d++ {
		if m%d == 0 {
			fac = append(fac, d)
			for m%d == 0 {
				m /= d
			}
		}
	}
	if m > 1 {
		fac = append(fac, m)
	}
	for g := 2; g < q; g++ {
		ok := true
		for _, f := range fac {
			if potMod(g, (q-1)/f, q) == 1 {
				ok = false
				break
			}
		}
		if ok {
			return g
		}
	}
	return 1
}

func potMod(b, e, m int) int {
	r := 1
	b %= m
	for e > 0 {
		if e&1 == 1 {
			r = r * b % m
		}
		b = b * b % m
		e >>= 1
	}
	return r
}

// indices returns the discrete logarithm of every residue mod q.
func indices(q int) map[int]int {
	g := raizPrimitiva(q)
	ind := map[int]int{}
	x := 1
	for e := 0; e < q-1; e++ {
		ind[x] = e
		x = x * g % q
	}
	return ind
}

// canales builds the K coupling vectors. cruzados = true gives the characters
// (every channel touches every prime); false gives the disjoint control.
func canales(ms []modo, q int, cruzados bool) [][]float64 {
	K := q - 1
	if K < 1 {
		K = 1
	}
	ind := indices(q)
	vs := make([][]float64, K)
	for a := range vs {
		vs[a] = make([]float64, len(ms))
	}
	for i, m := range ms {
		r := m.p % q
		e, ok := ind[r]
		if !ok { // p divides q: it belongs to no class, leave it uncoupled
			continue
		}
		if cruzados {
			for a := 0; a < K; a++ {
				vs[a][i] = m.amp * math.Cos(2*math.Pi*float64(a)*float64(e)/float64(K))
			}
		} else {
			vs[e%K][i] = m.amp // disjoint support: one class, one channel
		}
	}
	return vs
}

// rangoEfectivo measures how many channels are actually independent, by the
// eigenvalues of the K x K Gram matrix. The cosine family has repeats, and the
// sheet must report the rank it really has rather than the one it asked for.
func rangoEfectivo(vs [][]float64) int {
	K := len(vs)
	G := make([][]float64, K)
	for a := range G {
		G[a] = make([]float64, K)
		for b := range G[a] {
			s := 0.0
			for i := range vs[a] {
				s += vs[a][i] * vs[b][i]
			}
			G[a][b] = s
		}
	}
	ev := jacobi(G, 60)
	mx := 0.0
	for _, e := range ev {
		if math.Abs(e) > mx {
			mx = math.Abs(e)
		}
	}
	n := 0
	for _, e := range ev {
		if math.Abs(e) > 1e-9*mx {
			n++
		}
	}
	return n
}

// ---------------------------------------------------------------------------
// The spectrum: dense symmetric eigenvalues by cyclic Jacobi
// ---------------------------------------------------------------------------

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
				c := 1 / math.Sqrt(t*t+1)
				s := t * c
				for k := 0; k < n; k++ {
					akp, akq := a[k][p], a[k][q]
					a[k][p] = c*akp - s*akq
					a[k][q] = s*akp + c*akq
				}
				for k := 0; k < n; k++ {
					apk, aqk := a[p][k], a[q][k]
					a[p][k] = c*apk - s*aqk
					a[q][k] = s*apk + c*aqk
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

// normalizar rescales the channels so the TOTAL perturbation strength is the
// same for every K. Without this, raising K also raises the total coupling and
// the experiment confuses "more channels" with "more force" - which is exactly
// what the auditor's control asks us not to do.
func normalizar(vs [][]float64, objetivo float64) {
	tot := 0.0
	for _, v := range vs {
		for _, x := range v {
			tot += x * x
		}
	}
	if tot <= 0 {
		return
	}
	f := math.Sqrt(objetivo / tot)
	for _, v := range vs {
		for i := range v {
			v[i] *= f
		}
	}
}

// fuerzaTotal is sum_a |v_a|^2, the trace of the perturbation at g = 1.
func fuerzaTotal(vs [][]float64) float64 {
	t := 0.0
	for _, v := range vs {
		for _, x := range v {
			t += x * x
		}
	}
	return t
}

// espectro diagonalises D + g * sum_a |v_a><v_a|.
func espectro(ms []modo, vs [][]float64, g float64) []float64 {
	n := len(ms)
	H := make([][]float64, n)
	for i := range H {
		H[i] = make([]float64, n)
	}
	for i := range ms {
		H[i][i] = ms[i].w
	}
	for _, v := range vs {
		for i := 0; i < n; i++ {
			if v[i] == 0 {
				continue
			}
			for j := i; j < n; j++ {
				if v[j] == 0 {
					continue
				}
				x := g * v[i] * v[j]
				H[i][j] += x
				if i != j {
					H[j][i] += x
				}
			}
		}
	}
	ev := jacobi(H, 22)
	// A positive rank-K perturbation pushes K "doorway" states OUT of the band:
	// they are the collective modes, and they are not part of the medium's
	// response. Keeping them would inflate the range and make every normalised
	// spacing meaningless - which is exactly the trap this sheet fell into on
	// its first run. Only the levels inside the original band are kept.
	lo, hi := ms[0].w, ms[len(ms)-1].w
	dentro := ev[:0]
	for _, e := range ev {
		if e >= lo && e <= hi {
			dentro = append(dentro, e)
		}
	}
	return dentro
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

// desplegarPropio unfolds by the candidate's own mean density.
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
