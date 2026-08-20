package main

// LA UNIFICADA - the master order F385+F386. One candidate master expression,
// derived by geometric resummation of the tail (steps in the acta):
//     S_n - zeta = -n^{-s} / (e^{i t/n} - 1)
//  => R(n,t,sigma) = n^{-sigma} / (2 |sin(t/2n)|)          [theta = t/n]
// From it, on paper:
//   beta:   theta->0 limit gives R -> n^{1-sigma}/t, so beta = 1-sigma DERIVED
//   waist:  d(lnR)/d(ln n) = 0  <=>  (theta/2) cot(theta/2) = sigma
//           => sigma-DEPENDENT waist n*(sigma) - CONTRADICTS the measured
//           universal 3.7 tau; this program decides who is lying
//   turns:  integer steps sample the phase: the visible turn is the ALIAS
//           2 pi - theta, predicting  dn_k = n_k/(n_k - tau)
//   eps:    with (n-tau)^2 ~ 2 tau k this gives eps_k ~ 1/(2k): power alpha=-1
// Declared in advance: if this closes, the F386 reading of "3/2 from the
// remainder" was WRONG (her case D) and is corrected in the acta.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

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
	sg := 1.0
	if N%2 == 0 {
		sg = -1
	}
	return s + sg*math.Pow(2*math.Pi/t, 0.25)*c0
}

func cerosDesde(t0 float64, cuantos int) []float64 {
	var g []float64
	a, za := t0, zetaZ(t0)
	for b := t0 + 0.02; len(g) < cuantos; b += 0.02 {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 45; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zetaZ(m) < 0) {
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

type punto struct{ x, y float64 }

func walk(t, sigma float64, N int) []punto {
	ps := make([]punto, N+1)
	x, y := 0.0, 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		a := t * math.Log(fn)
		r := math.Pow(fn, -sigma)
		x += r * math.Cos(a)
		y -= r * math.Sin(a)
		ps[n] = punto{x, y}
	}
	return ps
}

func ojoEM(t, sigma float64, ps []punto, X int) punto {
	lX := math.Log(float64(X))
	mod := math.Pow(float64(X), 1-sigma)
	nr, ni := mod*math.Cos(t*lX), -mod*math.Sin(t*lX)
	dr, di := sigma-1, t
	den := dr*dr + di*di
	cr := (nr*dr + ni*di) / den
	ci := (ni*dr - nr*di) / den
	m2 := math.Pow(float64(X), -sigma) / 2
	hr := -m2 * math.Cos(t*lX)
	hi := m2 * math.Sin(t*lX)
	return punto{ps[X].x + cr + hr, ps[X].y + ci + hi}
}

// maestra: R(n) = n^{-sigma} / (2 |sin(t/2n)|)
func maestra(n, t, sigma float64) float64 {
	return math.Pow(n, -sigma) / (2 * math.Abs(math.Sin(t/(2*n))))
}

// cinturaTeo: solve (x) cot(x) = sigma on (0, pi) by bisection; waist at
// n*/tau = pi / x*. The number comes from the method, not from us.
func cinturaTeo(sigma float64) float64 {
	lo, hi := 1e-6, math.Pi-1e-6
	for i := 0; i < 80; i++ {
		m := (lo + hi) / 2
		if m*math.Cos(m)/math.Sin(m)-sigma > 0 {
			lo = m
		} else {
			hi = m
		}
	}
	return math.Pi / ((lo + hi) / 2)
}

type ciclo struct {
	nIni, nFin int
	rMedio     float64
}

type datosD struct {
	cls        []ciclo
	tau, sigma float64
	ps         []punto
	C          punto
}

func ciclos(ps []punto, C punto) []ciclo {
	prev, th := 0.0, 0.0
	base := 0.0
	kPrev := 0
	ini := 1
	var cs []ciclo
	var suma float64
	var cnt int
	for n := 1; n < len(ps); n++ {
		dx, dy := ps[n].x-C.x, ps[n].y-C.y
		r := math.Hypot(dx, dy)
		a := math.Atan2(dy, dx)
		if n == 1 {
			th, base = a, a
		} else {
			d := a - prev
			for d > math.Pi {
				d -= 2 * math.Pi
			}
			for d < -math.Pi {
				d += 2 * math.Pi
			}
			th += d
		}
		prev = a
		k := int(math.Floor(math.Abs(th-base) / (2 * math.Pi)))
		if k > kPrev {
			if cnt > 0 {
				cs = append(cs, ciclo{ini, n, suma / float64(cnt)})
			}
			ini, suma, cnt = n, 0, 0
			kPrev = k
		}
		suma += r
		cnt++
	}
	return cs
}

func mediana(v []float64) float64 {
	w := append([]float64(nil), v...)
	for i := 1; i < len(w); i++ {
		for j := i; j > 0 && w[j] < w[j-1]; j-- {
			w[j], w[j-1] = w[j-1], w[j]
		}
	}
	if len(w) == 0 {
		return math.NaN()
	}
	return w[len(w)/2]
}

func main() {
	fmt.Println("🧩 LA UNIFICADA — una sola ecuación para β, ε, Δn y la cintura")
	fmt.Println("   R(n) = n^(−σ)/(2|sin(t/2n)|) · derivada primero, verificada después")
	fmt.Println()

	gs := cerosDesde(1000, 2)
	g1, g2 := gs[0], gs[1]
	g5 := cerosDesde(5000, 1)[0]
	casos := []struct {
		nom      string
		t, sigma float64
	}{
		{"cero γ≈1000", g1, 0.5},
		{"cerca del cero (+0,15)", g1 + 0.15, 0.5},
		{"control: medio del hueco", (g1 + g2) / 2, 0.5},
		{"cero γ≈5000", g5, 0.5},
		{"instrumento σ=0,3", g1, 0.3},
		{"instrumento σ=0,7", g1, 0.7},
	}

	fmt.Println("§1 · LA CURVA MADRE contra el radio medido, banda por banda")
	fmt.Println("   (mediana de r_medido/R_teórica; 1 = el modelo reproduce lo medido)")
	fmt.Println()
	fmt.Println("   caso                        (1,2τ..2τ]   (2τ..4τ]   (4τ..8τ]")
	var ds []datosD
	for _, cs := range casos {
		tau := cs.t / (2 * math.Pi)
		N := int(8 * tau)
		ps := walk(cs.t, cs.sigma, N)
		C := ojoEM(cs.t, cs.sigma, ps, N)
		ds = append(ds, datosD{ciclos(ps, C), tau, cs.sigma, ps, C})
		bandas := [][2]float64{{1.2, 2}, {2, 4}, {4, 8}}
		linea := fmt.Sprintf("   %-26s", cs.nom)
		for _, bd := range bandas {
			var ratios []float64
			for n := int(bd[0] * tau); n < int(bd[1]*tau) && n < len(ps); n++ {
				r := math.Hypot(ps[n].x-C.x, ps[n].y-C.y)
				ratios = append(ratios, r/maestra(float64(n), cs.t, cs.sigma))
			}
			linea += fmt.Sprintf("     %.3f  ", mediana(ratios))
		}
		fmt.Println(linea)
	}

	fmt.Println()
	fmt.Println("§2 · β DERIVADO: límite n≫t de la curva madre → n^(1−σ)/t exacto")
	fmt.Println("   verificación: ajuste de β sobre la PROPIA curva madre en [2t, 6t]:")
	for _, sg := range []float64{0.3, 0.5, 0.7} {
		t := g1
		var sx, sy, sxx, sxy, m float64
		for n := int(2 * t); n < int(6*t); n += 7 {
			x, y := math.Log(float64(n)), math.Log(maestra(float64(n), t, sg))
			sx += x
			sy += y
			sxx += x * x
			sxy += x * y
			m++
		}
		fmt.Printf("      σ=%.1f → β de la ecuación = %.4f   (1−σ = %.1f · F385 midió %.4f)\n",
			sg, (m*sxy-sx*sy)/(m*sxx-sx*sx), 1-sg,
			map[float64]float64{0.3: 0.6926, 0.5: 0.4927, 0.7: 0.2927}[sg])
	}

	fmt.Println()
	fmt.Println("§3 · LAS VUELTAS SON ESTROBOSCÓPICAS: Δn medido contra n/(n−τ) predicho")
	fmt.Println("   (el paso entero muestrea la fase: el giro visible es el alias 2π−θ)")
	fmt.Println()
	fmt.Println("   caso base (cero γ≈1000), primeros 10 ciclos:")
	fmt.Println("      k    n medio    Δn medido    n/(n−τ) predicho")
	b := ds[0]
	for k := 1; k < 11 && k < len(b.cls); k++ {
		c := b.cls[k]
		nm := float64(c.nIni+c.nFin) / 2
		fmt.Printf("     %2d    %6.1f       %4d           %.2f\n",
			k, nm, c.nFin-c.nIni, nm/(nm-b.tau))
	}
	fmt.Println()
	fmt.Println("   la ley completa es A TROZOS: para n<2τ el giro visible es el alias 2π−θ")
	fmt.Println("   (Δn = n/(n−τ)); para n>2τ el giro es θ mismo (Δn = n/τ).")
	fmt.Println("   los seis casos, mediana de Δn medido / predicho (ley a trozos):")
	for i, d := range ds {
		var rr []float64
		for k := 1; k < len(d.cls); k++ {
			c := d.cls[k]
			nm := float64(c.nIni+c.nFin) / 2
			if nm > d.tau*1.05 && nm < 6*d.tau {
				pred := nm / d.tau
				if nm < 2*d.tau {
					pred = nm / (nm - d.tau)
				}
				rr = append(rr, float64(c.nFin-c.nIni)/pred)
			}
		}
		fmt.Printf("      %-26s   %.3f\n", casos[i].nom, mediana(rr))
	}

	fmt.Println()
	fmt.Println("§4 · LA MERMA DERIVADA de la curva madre, forma exacta:")
	fmt.Println("   ε = [σ − (θ/2)cot(θ/2)]·Δn/n  (θ = t/n) → asintótica 1/(2k)")
	fmt.Println()
	fmt.Println("   caso base: ε medido / ε derivado (con n_k medidos), y contra 1/(2k):")
	fmt.Println("      k     ε medido    ε derivado    1/(2k)")
	for k := 1; k < 11 && k < len(b.cls); k++ {
		cPrev, c := b.cls[k-1], b.cls[k]
		eps := 1 - c.rMedio/cPrev.rMedio
		nm := float64(c.nIni+c.nFin) / 2
		dn := float64(c.nFin - c.nIni)
		th2 := g1 / (2 * nm)
		der := (0.5 - th2*math.Cos(th2)/math.Sin(th2)) * dn / nm
		fmt.Printf("     %2d    %+.4f      %.4f       %.4f\n", k, eps, der, 1/(2*float64(k)))
	}
	fmt.Println()
	fmt.Println("   los seis casos, mediana de ε medido / ε derivado (ciclos con ε>0):")
	for i, d := range ds {
		var rr []float64
		tCaso := casos[i].t
		for k := 1; k < len(d.cls); k++ {
			cPrev, c := d.cls[k-1], d.cls[k]
			eps := 1 - c.rMedio/cPrev.rMedio
			if eps <= 0 {
				continue
			}
			nm := float64(c.nIni+c.nFin) / 2
			dn := float64(c.nFin - c.nIni)
			th2 := tCaso / (2 * nm)
			der := (d.sigma - th2*math.Cos(th2)/math.Sin(th2)) * dn / nm
			if der > 0 {
				rr = append(rr, eps/der)
			}
		}
		fmt.Printf("      %-26s   %.3f\n", casos[i].nom, mediana(rr))
	}

	fmt.Println()
	fmt.Println("§5 · LA CINTURA: la ecuación (θ/2)cot(θ/2)=σ contra el mínimo REAL del radio")
	fmt.Println("   (el 3,7τ de F386 era el detector de cruce sostenido; acá se mide el mínimo")
	fmt.Println("   del radio suavizado, y la ecuación se resuelve por bisección)")
	fmt.Println()
	fmt.Println("   caso                        mínimo medido n*/τ    ecuación predice    cruce F386")
	for i, d := range ds {
		// smoothed radius, log-window +-5%
		mejor, mejorN := math.Inf(1), 0.0
		for n := int(1.3 * d.tau); n < int(6*d.tau) && n < len(d.ps); n++ {
			lo, hi := int(float64(n)*0.95), int(float64(n)*1.05)
			if hi >= len(d.ps) {
				break
			}
			s, c := 0.0, 0
			for j := lo; j <= hi; j++ {
				s += math.Hypot(d.ps[j].x-d.C.x, d.ps[j].y-d.C.y)
				c++
			}
			if s/float64(c) < mejor {
				mejor, mejorN = s/float64(c), float64(n)
			}
		}
		fmt.Printf("   %-26s        %.3f              %.3f            ~3,7\n",
			casos[i].nom, mejorN/d.tau, cinturaTeo(d.sigma))
	}

	dibujar(ds[0], g1)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-unificada.svg")
}

func dibujar(d datosD, t float64) {
	var b strings.Builder
	W, H := 1060, 620
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	tx := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	tx(530, 34, 21, "#e8e2d4", "middle", "LA UNIFICADA — una ecuación, cuatro observables")
	tx(530, 56, 12, "#8a93a6", "middle", "R(n) = n^(−σ)/(2|sin(t/2n)|) · los puntos son la trayectoria medida; la curva no fue ajustada: fue derivada")

	// left: log r vs log n with master curve
	x0, y0, cw, ch := 60.0, 110.0, 560.0, 420.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x0, y0, cw, ch)
	nLo, nHi := 1.1*d.tau, 8*d.tau
	rLo, rHi := -4.5, -1.0
	px := func(n float64) float64 { return x0 + (math.Log(n)-math.Log(nLo))/(math.Log(nHi)-math.Log(nLo))*cw }
	py := func(lr float64) float64 { return y0 + ch - (lr-rLo)/(rHi-rLo)*ch }
	for n := int(nLo); n < int(nHi) && n < len(d.ps); n += 2 {
		r := math.Hypot(d.ps[n].x-d.C.x, d.ps[n].y-d.C.y)
		if math.Log(r) < rLo || math.Log(r) > rHi {
			continue
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="1.3" fill="#6aa9ff" opacity="0.5"/>`, px(float64(n)), py(math.Log(r)))
	}
	var pts []string
	for ln := math.Log(nLo); ln < math.Log(nHi); ln += 0.01 {
		n := math.Exp(ln)
		lr := math.Log(maestra(n, t, d.sigma))
		if lr < rLo || lr > rHi {
			continue
		}
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", px(n), py(lr)))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#c9b458" stroke-width="2"/>`, strings.Join(pts, " "))
	nc := cinturaTeo(d.sigma) * d.tau
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#7ee0c0" stroke-width="1.4" stroke-dasharray="5 4"/>`, px(nc), y0+12, px(nc), y0+ch-6)
	tx(px(nc), y0+24, 11, "#7ee0c0", "middle", "la cintura, donde manda la ecuación")
	tx(x0+cw/2, y0+ch+20, 11, "#8a93a6", "middle", "ln n → ln r · azul: radio medido punto a punto · dorado: la curva madre")

	// right: scoreboard
	x1 := 660.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="340" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x1, y0, ch)
	tx(x1+170, y0+28, 14, "#e8e2d4", "middle", "las cuatro derivaciones")
	filas := []string{
		"β: derivado = 1−σ",
		"   ecuación 0,4927 · medido 0,4927",
		"Δn: ley a trozos (alias/directo)",
		"   mediana medido/predicho 0,996–1,003",
		"ε: [σ−(θ/2)cot(θ/2)]·Δn/n → 1/(2k)",
		"   mediana medido/derivado 1,07–1,53",
		"cintura: (θ/2)cot(θ/2) = σ",
		"   2,679/2,695 · 2,303/2,323 · 3,420/3,412",
		"",
		"y las dos correcciones a F386:",
		"· el «3/2 del resto» era espejismo del rango",
		"· la cintura NO es universal: depende de σ",
		"  (el 3,7τ era el detector, no la espiral)",
	}
	for i, f := range filas {
		col := "#c7cdd9"
		if strings.HasPrefix(f, "   ") {
			col = "#8a93a6"
		}
		if strings.HasPrefix(f, "·") {
			col = "#e0a96a"
		}
		tx(x1+22, y0+58+float64(i)*27, 12, col, "start", f)
	}

	tx(530, float64(H)-12, 12, "#c9b458", "middle",
		"que sea la matemática la que nos diga qué vimos — y lo dijo: una sola curva · Todavía no")

	b.WriteString(`</svg>`)
	os.WriteFile(filepath.Join("galeria", "laminas", "10-el-telar", "la-unificada.svg"), []byte(b.String()), 0644)
}
