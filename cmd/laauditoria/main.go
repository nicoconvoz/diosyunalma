package main

// LA AUDITORIA INDEPENDIENTE - Fase XXI. A from-scratch SECOND IMPLEMENTATION
// of the Phase XX experiment, sharing NO code, tables or intermediate results
// with cmd/laformulamadre. Deliberately different algorithmic route:
//
//   - S(g) summed PER PRIME AND POWER: S = -(1/pi) sum_p sum_m sin(m g ln p)/(m p^{m/2})
//     (implementation A used a precomputed Lambda-weight table over prime powers)
//   - theta(t) with one extra asymptotic term (+31/80640 t^5), toggleable for
//     the numerical-sensitivity test
//   - model points found by GRID SCAN of upward half-integer crossings of the
//     counting function, then bisection (A marched point by point)
//   - real zeros re-found with a fresh Riemann-Siegel scan at step 0.01
//     (A used 0.02)
//   - fresh echo/binning/crossing/T4 code, and a DIFFERENT RNG (splitmix64;
//     A used xorshift64)
//
// Tests, in her order: (1) the four truncation rungs, all headline quantities,
// plus model-to-zero distances; (2) gamma-range windows; (3) fine numerical
// sensitivity; and the loose end 0.42 -> 0.39 diagnosed by window, not fixed.
// Frozen definitions (unfolding, zoom bins, periods, slope window) are kept
// identical BY DEFINITION - that is the experiment; only the code is new.

import (
	"fmt"
	"math"
)

// ---------- fresh RNG: splitmix64 -------------------------------------------

type rng struct{ s uint64 }

func (r *rng) u() float64 {
	r.s += 0x9E3779B97F4A7C15
	z := r.s
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	z ^= z >> 31
	return float64(z>>11) / float64(uint64(1)<<53)
}

// ---------- fresh smooth phase ----------------------------------------------

var thetaExtra = true

func fase(t float64) float64 { // Riemann-Siegel theta, one extra term
	v := t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t*t)
	if thetaExtra {
		v += 31 / (80640 * t * t * t * t * t)
	}
	return v
}

func nLiso(t float64) float64 { return fase(t)/math.Pi + 1 }

// ---------- fresh S: per prime, per power ------------------------------------

type termino struct{ w, l float64 } // weight 1/(m p^{m/2}), log = m ln p

func terminos(tope int) []termino {
	esP := make([]bool, tope+1)
	for i := 2; i <= tope; i++ {
		esP[i] = true
	}
	for i := 2; i*i <= tope; i++ {
		if esP[i] {
			for j := i * i; j <= tope; j += i {
				esP[j] = false
			}
		}
	}
	var ts []termino
	for p := 2; p <= tope; p++ {
		if !esP[p] {
			continue
		}
		lp := math.Log(float64(p))
		pm := float64(p)
		for m := 1; ; m++ {
			if pm > float64(tope) {
				break
			}
			ts = append(ts, termino{1 / (float64(m) * math.Sqrt(pm)), float64(m) * lp})
			pm *= float64(p)
		}
	}
	return ts
}

func fluct(g float64, ts []termino) float64 {
	s := 0.0
	for _, t := range ts {
		s += t.w * math.Sin(g*t.l)
	}
	return -s / math.Pi
}

// ---------- fresh model finder: grid scan of half-integer crossings ----------

func modelo(gLo, gHi, paso float64, ts []termino) []float64 {
	F := func(g float64) float64 { return nLiso(g) + fluct(g, ts) }
	var out []float64
	prev := F(gLo)
	k := math.Floor(prev + 0.5) // next half-integer level index
	for x := gLo + paso; x <= gHi; x += paso {
		cur := F(x)
		for cur >= k+0.5 { // upward crossing(s) of level k+1/2 in this cell
			lo, hi := x-paso, x
			for i := 0; i < 50; i++ {
				m := (lo + hi) / 2
				if F(m) < k+0.5 {
					lo = m
				} else {
					hi = m
				}
			}
			out = append(out, (lo+hi)/2)
			k++
		}
		prev = cur
		_ = prev
	}
	return out
}

// ---------- fresh real zeros: Riemann-Siegel, step 0.01 ----------------------

func zRS(t float64) float64 {
	th := fase(t)
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
	sg := 1.0
	if N%2 == 0 {
		sg = -1
	}
	return s + sg*math.Pow(2*math.Pi/t, 0.25)*c0
}

func cerosReales(gLo, gHi, paso float64) []float64 {
	var g []float64
	a, za := gLo, zRS(gLo)
	for b := gLo + paso; b <= gHi; b += paso {
		zb := zRS(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 50; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zRS(m) < 0) {
					hi = m
				} else {
					lo = m
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

// ---------- fresh mirror pipeline --------------------------------------------

const (
	bLo, bHi, bW = 0.50, 1.30, 0.05
)

var periodos []float64

func init() {
	for _, p := range []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} {
		periodos = append(periodos, math.Log(float64(p)))
	}
}

type parejaB struct{ base, gap float64 }

func parejas(g []float64) []parejaB {
	var ps []parejaB
	for i := 0; i+1 < len(g); i++ {
		ps = append(ps, parejaB{g[i], g[i+1] - g[i]})
	}
	return ps
}

func sNorm(p parejaB) float64 {
	return p.gap * math.Log(p.base/(2*math.Pi)) / (2 * math.Pi)
}

func curva(ps []parejaB, gaps []float64) []float64 {
	nb := int((bHi - bLo) / bW)
	sc := make([]float64, nb)
	cn := make([]int, nb)
	for i, p := range ps {
		gp := p.gap
		if gaps != nil {
			gp = gaps[i]
		}
		s := gp * math.Log(p.base/(2*math.Pi)) / (2 * math.Pi)
		if s < bLo || s >= bHi {
			continue
		}
		b := int((s - bLo) / bW)
		cn[b]++
		m := p.base + gp/2
		for _, T := range periodos {
			sc[b] += math.Cos(m * T)
		}
	}
	out := make([]float64, nb)
	for b := range out {
		if cn[b] > 0 {
			out[b] = -2 * sc[b] / float64(cn[b]*len(periodos))
		}
	}
	return out
}

func espejo21(g []float64, r *rng, nulos int) (E []float64, amp, cruce, pend float64) {
	ps := parejas(g)
	gaps := make([]float64, len(ps))
	for i, p := range ps {
		gaps[i] = p.gap
	}
	nb := int((bHi - bLo) / bW)
	nm := make([]float64, nb)
	for k := 0; k < nulos; k++ {
		br := append([]float64(nil), gaps...)
		for i := len(br) - 1; i > 0; i-- {
			j := int(r.u() * float64(i+1))
			br[i], br[j] = br[j], br[i]
		}
		c := curva(ps, br)
		for b := range c {
			nm[b] += c[b] / float64(nulos)
		}
	}
	real := curva(ps, nil)
	E = make([]float64, nb)
	for b := range E {
		E[b] = real[b] - nm[b]
		if math.Abs(E[b]) > amp {
			amp = math.Abs(E[b])
		}
	}
	for b := 0; b+1 < nb; b++ {
		if E[b] > 0 && E[b+1] < 0 {
			s1 := bLo + (float64(b)+0.5)*bW
			cruce = s1 + (-E[b]/(E[b+1]-E[b]))*bW
			break
		}
	}
	var sx, sy, sxx, sxy, n float64
	for b := 0; b < nb; b++ {
		sm := bLo + (float64(b)+0.5)*bW
		if math.Abs(sm-cruce) <= 0.15 {
			sx += sm
			sy += E[b]
			sxx += sm * sm
			sxy += sm * E[b]
			n++
		}
	}
	pend = (n*sxy - sx*sy) / (n*sxx - sx*sx)
	return
}

func t4De(g []float64) float64 { // T4 of the first zoom bin, fresh code
	ps := parejas(g)
	// bin 1: s in [0.50, 0.55)
	n := 0
	for _, p := range ps {
		s := sNorm(p)
		if s < 0.50 || s >= 0.55 {
			continue
		}
		n++
	}
	if n == 0 {
		return 0
	}
	t4 := 0.0
	for _, T := range periodos {
		var sc2, sg2, sa float64
		for _, p := range ps {
			sa += math.Sin(p.base * T)
		}
		sa /= float64(len(ps))
		for _, p := range ps {
			s := sNorm(p)
			if s < 0.50 || s >= 0.55 {
				continue
			}
			sc2 += math.Sin(p.base * T)
			sg2 += math.Sin(p.gap * T / 2)
		}
		t4 += 2 * ((sc2/float64(n) - sa) * (sg2 / float64(n))) / float64(len(periodos))
	}
	return t4
}

func desvS(g []float64) float64 {
	var v []float64
	for _, p := range parejas(g) {
		v = append(v, sNorm(p))
	}
	m := 0.0
	for _, x := range v {
		m += x
	}
	m /= float64(len(v))
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return math.Sqrt(s / float64(len(v)-1))
}

func distancias(mod, real []float64) (med, mx float64, pct int) {
	j, n := 0, 0
	for _, x := range mod {
		for j+1 < len(real) && math.Abs(real[j+1]-x) < math.Abs(real[j]-x) {
			j++
		}
		d := math.Abs(x - real[j])
		med += d
		if d > mx {
			mx = d
		}
		if d < 0.1 {
			pct++
		}
		n++
	}
	return med / float64(n), mx, 100 * pct / n
}

func residuo(a, b []float64) float64 {
	m := 0.0
	for i := range a {
		if d := math.Abs(a[i] - b[i]); d > m {
			m = d
		}
	}
	return m
}

func main() {
	fmt.Println("🔬🪞 FASE XXI — AUDITORÍA INDEPENDIENTE: segunda implementación desde cero")
	fmt.Println("   ruta B: S por primo y potencia · theta con término extra · buscador por")
	fmt.Println("   barrido de semienteros · ceros reales re-hallados a paso 0,01 · RNG splitmix64")
	fmt.Println("   Nada compartido con la implementación A. Definiciones congeladas idénticas.")

	real := cerosReales(30, 4000, 0.01)
	r := &rng{s: 20260828}
	Ereal, ampR, crR, pendR := espejo21(real, r, 80)
	fmt.Printf("\n   ceros reales (ruta B, paso 0,01): %d — A había hallado 3474\n", len(real))
	fmt.Printf("   REAL (ruta B): amplitud %.3f · cruce %.3f · pendiente %.3f · s desv %.2f · T4 %.3f\n",
		ampR, crR, pendR, desvS(real), t4De(real))
	fmt.Println("   REAL (ruta A):          0,348 ·       0,866 ·           −0,924 ·        0,42 ·    0,267")

	// ---- Prueba 1: the four rungs ------------------------------------------
	fmt.Println("\n§1 · PRUEBA 1 — los cuatro escalones, ruta B (A entre paréntesis)")
	fmt.Printf("   %-8s %9s %9s %10s %7s %7s %9s %10s %6s\n",
		"N", "amplitud", "cruce", "pendiente", "s desv", "T4", "resid", "dist med", "<0,1")
	refA := map[int][4]float64{97: {0.267, 0.847, -0.615, 0.215}, 997: {0.315, 0.921, -0.995, 0.205},
		9973: {0.353, 0.869, -1.011, 0.256}, 99991: {0.374, 0.861, -0.867, 0.292}}
	for _, N := range []int{97, 997, 9973, 99991} {
		ts := terminos(N)
		mod := modelo(30, 4000, 0.01, ts)
		E, amp, cr, pend := espejo21(mod, r, 60)
		dm, _, pct := distancias(mod, real)
		a := refA[N]
		fmt.Printf("   %-8d %9.3f %9.3f %10.3f %7.2f %7.3f %9.4f %10.4f %5d%%\n",
			N, amp, cr, pend, desvS(mod), t4De(mod), residuo(Ereal, E), dm, pct)
		fmt.Printf("   %-8s (%7.3f) (%7.3f) (%8.3f)         (%5.3f)\n", "  A:", a[0], a[1], a[2], a[3])
	}

	// ---- Prueba 2: gamma windows -------------------------------------------
	fmt.Println("\n§2 · PRUEBA 2 — el rango de γ (N = 9973), ventanas declaradas")
	ts := terminos(9973)
	for _, w := range [][2]float64{{30, 2000}, {2000, 4000}, {30, 4000}} {
		var mr []float64
		for _, x := range real {
			if x >= w[0] && x <= w[1] {
				mr = append(mr, x)
			}
		}
		mod := modelo(w[0], w[1], 0.01, ts)
		Er, _, crr, _ := espejo21(mr, r, 60)
		Em, amp, cr, _ := espejo21(mod, r, 60)
		dm, _, pct := distancias(mod, mr)
		fmt.Printf("   γ∈[%4.0f,%4.0f]: modelo amp %.3f cruce %.3f s-desv %.2f · real cruce %.3f s-desv %.2f · resid %.4f · dist %.4f (%d%%)\n",
			w[0], w[1], amp, cr, desvS(mod), crr, desvS(mr), residuo(Er, Em), dm, pct)
	}

	// ---- Prueba 3: numerical sensitivity ------------------------------------
	fmt.Println("\n§3 · PRUEBA 3 — sensibilidad numérica fina (N = 9973)")
	base := modelo(30, 4000, 0.01, ts)
	_, ampB, crB, _ := espejo21(base, r, 60)
	fina := modelo(30, 4000, 0.004, ts)
	_, ampF, crF, _ := espejo21(fina, r, 60)
	thetaExtra = false
	sinExtra := modelo(30, 4000, 0.01, ts)
	_, ampT, crT, _ := espejo21(sinExtra, r, 60)
	thetaExtra = true
	fmt.Printf("   paso 0,010            : amp %.3f · cruce %.3f · puntos %d\n", ampB, crB, len(base))
	fmt.Printf("   paso 0,004            : amp %.3f · cruce %.3f · puntos %d\n", ampF, crF, len(fina))
	fmt.Printf("   θ sin término extra   : amp %.3f · cruce %.3f · puntos %d\n", ampT, crT, len(sinExtra))

	// ---- loose end: 0.42 -> 0.39 diagnosed by window and truncation ----------
	fmt.Println("\n§4 · EL CABO SUELTO 0,42→0,39 — diagnóstico por ventana y truncación (sin corregir)")
	for _, N := range []int{97, 9973, 99991} {
		t2 := terminos(N)
		lo := modelo(30, 2000, 0.01, t2)
		hi := modelo(2000, 4000, 0.01, t2)
		fmt.Printf("   N=%-6d s-desv: baja %.3f · alta %.3f\n", N, desvS(lo), desvS(hi))
	}
	var rl, rh []float64
	for _, x := range real {
		if x <= 2000 {
			rl = append(rl, x)
		} else {
			rh = append(rh, x)
		}
	}
	fmt.Printf("   REAL     s-desv: baja %.3f · alta %.3f\n", desvS(rl), desvS(rh))
}
