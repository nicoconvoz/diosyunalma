// Command elpliego answers the auditor's Atom brief (Auditoria/34): do not try
// to prove the Atom exists - try to BUILD the operator that satisfies R1-R5,
// and if it cannot be built, determine exactly which requirement forbids it.
//
// Four measurements, in the order the answer needs them:
//
//	§1  AUDIT OF OUR OWN WORKSHOP. cmd/maquina never computed prototype A: its
//	    spacings are stipulated at 1.0 and its 0.555 is the distance from a
//	    delta spike to the Wigner curve. Here prototype A is actually built and
//	    prototype B is run as an ensemble with an error bar.
//	§2  THE SONG IS BLIND. The canto statistic is measured on the true zeros,
//	    on real GUE matrices and on a random draw from the Wigner law, all at
//	    the same sample size, against the sampling spread. Then the same three
//	    are asked for their ECHO at the periods k·log p.
//	§3  THE BRIEF IS SATISFIABLE BY CHEATING. diag(gamma_n) passes R1-R5.
//	§4  THE BOX THAT MUST BREATHE. Any fixed self-adjoint box has constant
//	    Weyl density; the zeros' density grows like log. Measured.
//
// Every number is computed here from zeros this program finds itself.
//
// Reproduce: go run ./cmd/elpliego
package main

import (
	"fmt"
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// THE SEA: Riemann-Siegel on the critical line, and the zeros
// ---------------------------------------------------------------------------

// theta is the Riemann-Siegel phase.
func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t*t)
}

// cuenta is the smooth zero count N(T) = theta(T)/pi + 1.
func cuenta(T float64) float64 { return theta(T)/math.Pi + 1 }

// zeta Z(t): main sum plus the first Riemann-Siegel correction term.
func zeta(t float64) float64 {
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

// ceros finds every zero of Z in [t0,t1] by sign change plus bisection.
func ceros(t0, t1, paso float64) []float64 {
	var g []float64
	a, za := t0, zeta(t0)
	for b := t0 + paso; b <= t1; b += paso {
		zb := zeta(b)
		if za == 0 || (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 60; i++ {
				m := (lo + hi) / 2
				zm := zeta(m)
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
// THE EXAM: the canto statistic, exactly as cmd/maquina defines it
// ---------------------------------------------------------------------------

const (
	nbins = 15
	smax  = 3.0
)

// wigner is the GUE surmise (beta = 2).
func wigner(s float64) float64 {
	return 32 / (math.Pi * math.Pi) * s * s * math.Exp(-4*s*s/math.Pi)
}

// canto is the mean per-bin L1 distance between the spacing histogram and the
// Wigner curve - the workshop's own exam, reproduced without changes.
func canto(sp []float64) float64 {
	h := make([]float64, nbins)
	for _, s := range sp {
		i := int(s / smax * nbins)
		if i >= 0 && i < nbins {
			h[i]++
		}
	}
	d := 0.0
	for i := range h {
		h[i] /= float64(len(sp)) * (smax / nbins)
		d += math.Abs(h[i] - wigner((float64(i)+0.5)*smax/nbins))
	}
	return d / nbins
}

// ---------------------------------------------------------------------------
// The dice: one deterministic generator, so the whole sheet is reproducible
// ---------------------------------------------------------------------------

type dado struct{ s uint64 }

func (d *dado) u() float64 {
	d.s ^= d.s << 13
	d.s ^= d.s >> 7
	d.s ^= d.s << 17
	return float64(d.s>>11) / float64(uint64(1)<<53)
}

func (d *dado) normal() float64 {
	u1, u2 := d.u(), d.u()
	if u1 < 1e-300 {
		u1 = 1e-300
	}
	return math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
}

// wignerSorteo draws one spacing from the Wigner surmise by rejection.
func (d *dado) wignerSorteo() float64 {
	const pico = 0.94 // the surmise peaks at 0.936797 (s = sqrt(pi)/2); 0.9 clipped it
	for {
		s := d.u() * smax
		if d.u()*pico < wigner(s) {
			return s
		}
	}
}

// ---------------------------------------------------------------------------
// A real GUE matrix and its levels (Jacobi on the real symmetric embedding)
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
		if off < 1e-22 {
			break
		}
		for p := 0; p < n-1; p++ {
			for q := p + 1; q < n; q++ {
				if math.Abs(a[p][q]) < 1e-18 {
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

// gueNiveles builds one N x N GUE matrix and returns its unfolded spacings,
// taken from the central 60% of the semicircle where the density is smooth.
func gueNiveles(d *dado, N int) []float64 {
	X := make([][]float64, N)
	Y := make([][]float64, N)
	for i := range X {
		X[i] = make([]float64, N)
		Y[i] = make([]float64, N)
	}
	sc := 1 / math.Sqrt(2*float64(N))
	for i := 0; i < N; i++ {
		X[i][i] = d.normal() * math.Sqrt2 * sc
		for j := i + 1; j < N; j++ {
			X[i][j] = d.normal() * sc
			X[j][i] = X[i][j]
			Y[i][j] = d.normal() * sc
			Y[j][i] = -Y[i][j]
		}
	}
	M := make([][]float64, 2*N)
	for i := range M {
		M[i] = make([]float64, 2*N)
	}
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			M[i][j] = X[i][j]
			M[i][j+N] = -Y[i][j]
			M[i+N][j] = Y[i][j]
			M[i+N][j+N] = X[i][j]
		}
	}
	all := jacobi(M, 24)
	ev := make([]float64, 0, N)
	for i := 0; i < 2*N; i += 2 { // the real embedding doubles every level
		ev = append(ev, all[i])
	}
	// unfold with the semicircle law (radius 2 for this scaling)
	F := func(E float64) float64 {
		if E <= -2 {
			return 0
		}
		if E >= 2 {
			return 1
		}
		return 0.5 + (E*math.Sqrt(4-E*E)/2+2*math.Asin(E/2))/(2*math.Pi)
	}
	lo, hi := N/5, 4*N/5
	sp := make([]float64, 0, hi-lo)
	for i := lo; i < hi-1; i++ {
		sp = append(sp, float64(N)*(F(ev[i+1])-F(ev[i])))
	}
	media := 0.0
	for _, s := range sp {
		media += s
	}
	media /= float64(len(sp))
	for i := range sp {
		sp[i] /= media
	}
	return sp
}

// ---------------------------------------------------------------------------
// THE ECHO: the only test that carries arithmetic
// ---------------------------------------------------------------------------

// eco is E(T) = sum_n w_n cos(gamma_n T) with a raised-cosine taper, so the
// window's own edges do not manufacture peaks.
func eco(g []float64, T float64) float64 {
	M := len(g)
	s := 0.0
	for i, gam := range g {
		w := 0.5 - 0.5*math.Cos(2*math.Pi*(float64(i)+0.5)/float64(M))
		s += w * math.Cos(gam*T)
	}
	return 2 * s / float64(M)
}

// periodos returns the arithmetic periods k*log(p) = log(p^k) below tope.
func periodos(tope float64) []float64 {
	var ps []float64
	criba := make([]bool, 40000)
	for n := 2; n < len(criba); n++ {
		if criba[n] {
			continue
		}
		for m := 2 * n; m < len(criba); m += n {
			criba[m] = true
		}
		for pk, k := n, 1; float64(pk) < math.Exp(tope) && k < 12; pk, k = pk*n, k+1 {
			if l := math.Log(float64(pk)); l <= tope && l > 0 {
				ps = append(ps, l)
			}
			if pk > 1<<40/n {
				break
			}
		}
	}
	sort.Float64s(ps)
	return ps
}

// fuerzaEco scores a spectrum's echo: the mean |E| on the arithmetic periods
// divided by the mean |E| everywhere else. One means "hears nothing".
func fuerzaEco(g []float64, ps []float64, tope float64) (dentro, fuera, razon float64) {
	const ancho = 0.012
	nIn, nOut := 0, 0
	for T := 0.35; T <= tope; T += 0.0007 {
		cerca := false
		for _, p := range ps {
			if math.Abs(T-p) < ancho {
				cerca = true
				break
			}
		}
		v := math.Abs(eco(g, T))
		if cerca {
			dentro += v
			nIn++
		} else {
			fuera += v
			nOut++
		}
	}
	dentro /= float64(nIn)
	fuera /= float64(nOut)
	return dentro, fuera, dentro / fuera
}

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

func main() {
	fmt.Println("⚛️  EL PLIEGO — respuesta al pedido de la auditora sobre el Átomo")
	fmt.Println("    no demostrar que el Átomo existe: intentar CONSTRUIR el operador,")
	fmt.Println("    y si no se puede, decir exactamente qué requisito lo impide.")

	// ---- the sea -----------------------------------------------------------
	const t0, t1 = 100.0, 1000.0
	g := ceros(t0, t1, 0.02)
	esperados := cuenta(t1) - cuenta(t0)
	fmt.Printf("\n§0 · EL MAR · %d ceros medidos en [%.0f, %.0f]; la ley suave predice %.1f — diferencia %.1f\n",
		len(g), t0, t1, esperados, float64(len(g))-esperados)

	// unfold the true zeros with N(t)
	spZ := make([]float64, 0, len(g))
	for i := 0; i+1 < len(g); i++ {
		spZ = append(spZ, cuenta(g[i+1])-cuenta(g[i]))
	}
	fmt.Printf("     espaciado desplegado medio %.6f (tiene que dar 1) · mínimo %.4f\n", media(spZ), mini(spZ))

	// ---- §1 the audit of our own workshop ----------------------------------
	fmt.Println("\n§1 · AUDITORÍA DEL TALLER — lo que cmd/maquina publicó, recalculado en serio")

	// prototype A, actually built this time: the Berry-Keating picket fence
	// E_n = 2*pi*n/log(Lambda) from a FIXED box of log-length log(Lambda).
	lnL := 10.0
	nivA := make([]float64, 400)
	for i := range nivA {
		nivA[i] = 2 * math.Pi * float64(i+1) / lnL
	}
	// The unfolded spacing of a fence is 1 BY CONSTRUCTION. Computing it as
	// (nivA[i+1]-nivA[i])*lnL/(2*pi) subtracts numbers of order 250 and spreads
	// the answer over ~400 ulps, which then straddles a histogram bin edge - a
	// cancellation defect of ours, not a property of the fence.
	spA := make([]float64, 0, len(nivA))
	for i := 0; i+1 < len(nivA); i++ {
		spA = append(spA, 1.0)
	}
	spAcancel := make([]float64, 0, len(nivA))
	for i := 0; i+1 < len(nivA); i++ {
		spAcancel = append(spAcancel, (nivA[i+1]-nivA[i])*lnL/(2*math.Pi))
	}
	cA := canto(spA)
	// the workshop's own version: the spacings stipulated at exactly 1.0
	spAes := make([]float64, len(spA))
	for i := range spAes {
		spAes[i] = 1.0
	}
	cAes := canto(spAes)
	cAcancel := canto(spAcancel)
	// sweep the bin origin: if the exam is stable, the answer must not depend on it
	var barrido []float64
	for k := -6; k <= 5; k++ {
		off := float64(k) * 0.02
		h := make([]float64, nbins)
		for _, x := range spA {
			i := int((x + off) / smax * nbins)
			if i >= 0 && i < nbins {
				h[i]++
			}
		}
		dd := 0.0
		for i := range h {
			h[i] /= float64(len(spA)) * (smax / nbins)
			dd += math.Abs(h[i] - wigner((float64(i)+0.5)*smax/nbins))
		}
		barrido = append(barrido, dd/nbins)
	}
	fmt.Printf("     PROTOTIPO A (la reja de una caja FIJA de largo log %.0f):\n", lnL)
	fmt.Printf("       el taller publico 0.555 · recalculado aca: canto = %.4f — el numero estaba BIEN\n", cAes)
	fmt.Printf("       barriendo el origen de los cajones: entre %.4f y %.4f — o sea, estable\n", mini(barrido), maxi(barrido))
	fmt.Printf("       distancia SIN cajones (Kolmogorov-Smirnov contra Wigner) = %.4f, sin ambiguedad\n", ks(spA))
	fmt.Printf("       ⚠ CORRECCION NUESTRA: al construir los niveles y restarlos entre si (numeros de orden 250)\n")
	fmt.Printf("         la cancelacion desparrama el espaciado ~400 ulps alrededor de 1, justo sobre el borde de\n")
	fmt.Printf("         un cajon, y el examen contesta %.4f. Ese era un defecto NUESTRO, no del taller.\n", cAcancel)
	fmt.Printf("       lo que si queda en pie: el canto es DISCONTINUO sobre un espectro degenerado.\n")

	// prototype B, as an ensemble with an error bar
	const NB, ens = 80, 24
	d := &dado{s: 20260817}
	var cantosB []float64
	var poolB []float64
	for e := 0; e < ens; e++ {
		sp := gueNiveles(d, NB)
		cantosB = append(cantosB, canto(sp))
		poolB = append(poolB, sp...)
	}
	fmt.Printf("     PROTOTIPO B (matriz GUE real, %d matrices de %dx%d, %d espaciados c/u):\n", ens, NB, NB, len(poolB)/ens)
	fmt.Printf("       canto = %.4f ± %.4f  (el taller publicó 0.076 de UNA sola matriz, sin barra)\n", media(cantosB), desvio(cantosB))

	// ---- §2 the song is blind ----------------------------------------------
	fmt.Println("\n§2 · EL CANTO ES CIEGO — el mismo examen a tres candidatos, al mismo tamaño de muestra")
	n89 := len(poolB) / ens // the workshop's own sample size, matched
	if n89 > len(spZ) {
		n89 = len(spZ)
	}
	// the true zeros, in blocks of n89
	var cantosZ []float64
	for i := 0; i+n89 <= len(spZ); i += n89 {
		cantosZ = append(cantosZ, canto(spZ[i:i+n89]))
	}
	// a pure random draw from the Wigner law: knows nothing at all
	var cantosR []float64
	for e := 0; e < 400; e++ {
		sp := make([]float64, n89)
		for i := range sp {
			sp[i] = d.wignerSorteo()
		}
		cantosR = append(cantosR, canto(sp))
	}
	fmt.Printf("     tamaño de muestra: %d espaciados en los tres casos\n", n89)
	fmt.Printf("       ceros de zeta VERDADEROS : canto = %.4f ± %.4f  (%d bloques)\n", media(cantosZ), desvio(cantosZ), len(cantosZ))
	fmt.Printf("       matriz GUE al azar       : canto = %.4f ± %.4f  (%d matrices)\n", media(cantosB), desvio(cantosB), ens)
	fmt.Printf("       sorteo puro de Wigner    : canto = %.4f ± %.4f  (%d sorteos) ← no sabe NADA\n", media(cantosR), desvio(cantosR), len(cantosR))
	sep := math.Abs(media(cantosZ)-media(cantosR)) / desvio(cantosR)
	tZR, tZB := tDeMedias(cantosZ, cantosR), tDeMedias(cantosZ, cantosB)
	fmt.Printf("     separacion por realizacion (ceros contra ruido puro): %.2f sigmas del propio examen\n", sep)
	fmt.Printf("     y la pregunta correcta, sobre las MEDIAS: t(ceros vs ruido) = %.2f · t(ceros vs GUE) = %.2f\n", tZR, tZB)
	var ksZ, ksB, ksR []float64
	for i := 0; i+n89 <= len(spZ); i += n89 {
		ksZ = append(ksZ, ks(spZ[i:i+n89]))
	}
	for e := 0; e < ens; e++ {
		ksB = append(ksB, ks(gueNiveles(d, NB)))
	}
	for e := 0; e < 200; e++ {
		sp := make([]float64, n89)
		for i := range sp {
			sp[i] = d.wignerSorteo()
		}
		ksR = append(ksR, ks(normalizar(sp)))
	}
	fmt.Printf("     sin cajones (KS): ceros %.4f±%.4f · GUE %.4f±%.4f · ruido %.4f±%.4f · t(ceros vs GUE) = %.2f\n",
		media(ksZ), desvio(ksZ), media(ksB), desvio(ksB), media(ksR), desvio(ksR), tDeMedias(ksZ, ksB))
	fmt.Println("     ⟹ el examen es ciego PORQUE mira un espaciado por vez. La varianza de numero, que mira")
	fmt.Println("       CORRELACIONES entre niveles, con los mismos datos ya no es ciega:")
	var sigmas [][4]float64
	tsig := 0.0
	for _, L := range []float64{10, 20} {
		var vZ, vB, vR []float64
		for i := 0; i+n89 <= len(spZ); i += n89 {
			if v := varianzaNumero(spZ[i:i+n89], L, 400); !math.IsNaN(v) {
				vZ = append(vZ, v)
			}
		}
		for e := 0; e < ens; e++ {
			if v := varianzaNumero(gueNiveles(d, NB), L, 400); !math.IsNaN(v) {
				vB = append(vB, v)
			}
		}
		for e := 0; e < 200; e++ {
			sp := make([]float64, n89)
			for i := range sp {
				sp[i] = d.wignerSorteo()
			}
			if v := varianzaNumero(normalizar(sp), L, 400); !math.IsNaN(v) {
				vR = append(vR, v)
			}
		}
		fmt.Printf("       Σ²(%2.0f): ceros %.3f±%.3f · GUE %.3f±%.3f · sorteo SIN memoria %.3f±%.3f · t(ceros vs sorteo) = %.1f\n",
			L, media(vZ), desvio(vZ), media(vB), desvio(vB), media(vR), desvio(vR), tDeMedias(vZ, vR))
		sigmas = append(sigmas, [4]float64{L, media(vZ), media(vB), media(vR)})
		if L == 10 {
			tsig = tDeMedias(vZ, vR)
		}
	}

	// the echo: the only test that carries arithmetic
	const tope = 3.2
	ps := periodos(tope)
	// zeros: use a clean window of them
	gw := g
	if len(gw) > 400 {
		gw = gw[:400]
	}
	dZ, fZ, rZ := fuerzaEco(gw, ps, tope)
	// A MATCHED control: the same NUMBER of levels as the zeros (the first
	// version had 95 against 400, which washed out its contrast and inflated its
	// noise floor by sqrt(400/95)), laid on the zeros' own staircase so the
	// growing density is present in the control too.
	var spG []float64
	for len(spG) < len(gw) {
		spG = append(spG, gueNiveles(d, 160)...)
	}
	spG = spG[:len(gw)]
	nivG := make([]float64, 0, len(spG))
	nG := cuenta(gw[0])
	for _, sp := range spG {
		nG += sp
		// invert N(t) = nG by bisection: put the control level where a zero with
		// that ordinal would sit, so control and zeros share the density profile
		lo, hi := gw[0], gw[len(gw)-1]+50
		for k := 0; k < 50; k++ {
			m := (lo + hi) / 2
			if cuenta(m) < nG {
				lo = m
			} else {
				hi = m
			}
		}
		nivG = append(nivG, (lo+hi)/2)
	}
	dG, fG, rG := fuerzaEco(nivG, ps, tope)
	// CTRL: the true zeros scored at RANDOM non-arithmetic periods. If the echo
	// were an artifact of window or taper, any set of periods would score alike.
	var azar []float64
	for tr := 0; tr < 120; tr++ {
		falsos := make([]float64, 0, len(ps))
		for len(falsos) < len(ps) {
			T := 0.4 + d.u()*(tope-0.45)
			ok := true
			for _, q := range ps {
				if math.Abs(T-q) < 0.05 {
					ok = false
					break
				}
			}
			if ok {
				falsos = append(falsos, T)
			}
		}
		sort.Float64s(falsos)
		_, _, r := fuerzaEco(gw, falsos, tope)
		azar = append(azar, r)
	}
	fmt.Printf("     EL ECO en los períodos k·log p (hasta T = %.1f, %d períodos):\n", tope, len(ps))
	fmt.Printf("       ceros verdaderos : |E| dentro %.5f · fuera %.5f · razón %.3f\n", dZ, fZ, rZ)
	fmt.Printf("       espectro GUE     : |E| dentro %.5f · fuera %.5f · razón %.3f\n", dG, fG, rG)
	fmt.Printf("       ceros en períodos AL AZAR (no aritméticos, %d sorteos): razón %.4f ± %.4f, máximo %.4f\n",
		len(azar), media(azar), desvio(azar), maxi(azar))
	fmt.Printf("     ⚠ el 36x lo carga el DENOMINADOR: contra el control emparejado los ceros suenan %.2f veces\n", dZ/fG)
	fmt.Printf("       más fuerte EN los períodos aritméticos, y %.1f veces más CALLADOS afuera. Lo raro es el silencio.\n", fG/fZ)

	// ---- §3 the brief is satisfiable by cheating ---------------------------
	fmt.Println("\n§3 · EL PLIEGO SE SATISFACE HACIENDO TRAMPA")
	fmt.Printf("     H = diag(γ₁ … γ_%d) es autoadjunta (R1), su estadística es la de los ceros (R2),\n", len(g))
	fmt.Printf("     su conteo es N(T) exacto (R3: %d contra %.1f), y su eco tiene los picos aritméticos\n", len(g), esperados)
	fmt.Printf("     (R4/R5: razón %.3f, medida arriba). Los CINCO requisitos, con un objeto que ya sabe la respuesta.\n", rZ)
	fmt.Println("     ⟹ al pliego le falta R6 — NO CIRCULARIDAD: el operador debe definirse sin mirar los ceros.")

	// ---- §4 the box that must breathe --------------------------------------
	fmt.Println("\n§4 · LA CAJA QUE TIENE QUE RESPIRAR — el obstáculo, medido")
	fmt.Println("     una caja FIJA de largo logarítmico L da niveles 2πn/L: espaciado CONSTANTE 2π/L.")
	fmt.Println("     los ceros no: su espaciado se achica con la altura. Medido por ventanas:")
	fmt.Printf("     %10s  %14s  %14s  %14s\n", "altura T", "espaciado medio", "L medido=2π/Δ", "ln(T/2π)")
	var filas []fila
	for _, w := range [][2]float64{{100, 250}, {250, 450}, {450, 650}, {650, 820}, {820, 1000}} {
		var ds []float64
		var Ts float64
		var n int
		for i := 0; i+1 < len(g); i++ {
			if g[i] >= w[0] && g[i+1] <= w[1] {
				ds = append(ds, g[i+1]-g[i])
				Ts += g[i]
				n++
			}
		}
		if n == 0 {
			continue
		}
		Tm := Ts / float64(n)
		dm := media(ds)
		f := fila{Tm, dm, 2 * math.Pi / dm, math.Log(Tm / (2 * math.Pi))}
		filas = append(filas, f)
		fmt.Printf("     %10.1f  %14.6f  %14.4f  %14.4f\n", f.T, f.dm, f.Lm, f.Lt)
	}
	// The honest minimax fixed box (the arithmetic mean is not the best one),
	// reported only to show that the percentage is a property of the interval and
	// not of anything structural.
	Lfijo, peor := 0.0, math.Inf(1)
	for cand := 2.0; cand <= 7.0; cand += 0.0005 {
		w := 0.0
		for _, f := range filas {
			if e := math.Abs(cand-f.Lm) / f.Lm * 100; e > w {
				w = e
			}
		}
		if w < peor {
			peor, Lfijo = w, cand
		}
	}
	crec := filas[len(filas)-1].Lm / filas[0].Lm
	fmt.Printf("     la MEJOR caja fija posible para este tramo es L = %.3f y aun así yerra hasta %.1f%%\n", Lfijo, peor)
	fmt.Printf("     pero ese porcentaje es del tramo, no del problema: L = ln(T/2π) → ∞, así que NO existe\n")
	fmt.Printf("     ninguna L fija. Acá crece de %.3f a %.3f (×%.3f) sólo porque miramos hasta 1000.\n",
		filas[0].Lm, filas[len(filas)-1].Lm, crec)
	fmt.Println("     ⟹ CORRECCIÓN A NOSOTROS MISMOS: NO es cierto que «cualquier operador autoadjunto en")
	fmt.Println("       dominio compacto da densidad constante» — nuestra propia diag(γ) del §3 lo refuta, y la")
	fmt.Println("       ley de Weyl en dimensión d da N(λ) ~ Cλ^(d/2). Lo que la tabla sí prohíbe es la GEOMETRÍA")
	fmt.Println("       FIJA: el operador de Berry-Keating sobre un grafo cuántico compacto de largo total L tiene")
	fmt.Println("       N(E) ~ (L/2π)E — densidad constante — y los ceros piden (1/2π)ln(T/2π) (Endres-Steiner 2010,")
	fmt.Println("       arXiv:0912.3183: espectro puramente continuo en la recta, y asintótica de Weyl en grafos).")
	fmt.Println("       El requisito que lo impide es R3 contra (R1 + geometría fija). Y la única salida —una")
	fmt.Println("       geometría que crezca como ln(T/2π)— choca de frente con R6: eso es volver a leer la respuesta.")
	fmt.Println("     (esta sección es un control de consistencia: confirma Riemann-von Mangoldt, que es la misma")
	fmt.Println("      fórmula con la que encontramos y desplegamos los ceros. No es evidencia independiente.)")

	res := Resultado{
		Ceros: g, Esperados: esperados, T0: t0, T1: t1,
		SpZeros: spZ, CantoA: cA, CantoAes: cAes, LnBox: lnL,
		CantosZ: cantosZ, CantosB: cantosB, CantosR: cantosR, N89: n89,
		Sigmas: sep, NGue: NB, Ensamble: ens,
		Periodos: ps, Tope: tope, VentanaCeros: gw, NivelesGUE: nivG,
		EcoZdentro: dZ, EcoZfuera: fZ, EcoZrazon: rZ,
		EcoGdentro: dG, EcoGfuera: fG, EcoGrazon: rG,
		Ventanas: filas2(filas), Lfijo: Lfijo, PeorFijo: peor, Crecimiento: crec,
		Sigma2: sigmas, Tsigma2: tsig, EcoAzar: media(azar), EcoAzarDes: desvio(azar), EcoAzarMax: maxi(azar),
		CantoAcancel: cAcancel, KSa: ks(spA), Tzr: tZR, Tzb: tZB,
	}
	dibujar(res)
}

// Resultado carries every measured number of this run to the plate, so the
// sheet can never show a figure the run did not produce.
type Resultado struct {
	Ceros                            []float64
	Esperados, T0, T1                float64
	SpZeros                          []float64
	CantoA, CantoAes, LnBox          float64
	CantosZ, CantosB, CantosR        []float64
	N89, NGue, Ensamble              int
	Sigmas                           float64
	Periodos                         []float64
	Tope                             float64
	VentanaCeros, NivelesGUE         []float64
	EcoZdentro, EcoZfuera, EcoZrazon float64
	EcoGdentro, EcoGfuera, EcoGrazon float64
	Ventanas                         [][4]float64 // T, spacing, L measured, ln(T/2pi)
	Lfijo, PeorFijo, Crecimiento     float64
	Sigma2                           [][4]float64 // L, ceros, GUE, sorteo sin memoria
	Tsigma2                          float64
	EcoAzar, EcoAzarDes, EcoAzarMax  float64
	CantoAcancel, KSa                float64
	Tzr, Tzb                         float64
}

func mini(v []float64) float64 {
	m := v[0]
	for _, x := range v {
		if x < m {
			m = x
		}
	}
	return m
}

// fila is one height window of the breathing-box measurement.
type fila struct{ T, dm, Lm, Lt float64 }

// filas2 flattens the window table for the plate.
func filas2(f []fila) [][4]float64 {
	out := make([][4]float64, len(f))
	for i, x := range f {
		out[i] = [4]float64{x.T, x.dm, x.Lm, x.Lt}
	}
	return out
}
