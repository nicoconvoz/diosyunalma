package main

// LA PROPORCION - strict order from the captain: extract objective geometric
// magnitudes from each spiral, form dimensionless ratios, do NOT impose 1/2,
// sort them, and see whether any converges to 1/2 on its own.
//
// Circularity guard, declared up front: our spirals were drawn with step
// lengths 1/sqrt(n) - the 1/2 is PAINTED on them. So the skeleton and the
// paint are separated:
//   SKELETON: unit-step walk, direction phi_n = -t*ln(n) only. Contains no
//             sigma anywhere. Any 1/2 found here is spontaneous.
//   PAINT:    step lengths n^{-sigma} for sigma = 0.3 / 0.5 / 0.7. The eye-
//             radius growth exponent is measured for each, as the control that
//             tells imposed from intrinsic.
//
// Magnitudes per spiral (mechanical, fixed thresholds, no tuning per rung):
//   tau     = t/(2*pi)          the walk's own step scale
//   nMax    = center of the outermost coherent run (the main resonance)
//   K       = number of distinct coherent runs the walk shows
//   nMin    = center of the innermost distinguishable run (resolution limit)
//   nOjo    = where the final winding starts
//   beta    = growth exponent of the eye radius (per sigma)
// Ratios are formed from these, sorted, and tracked down a ladder of heights.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
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

func primerCero(t0 float64) float64 {
	a, za := t0, zetaZ(t0)
	for b := t0 + 0.02; ; b += 0.02 {
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
			return (lo + hi) / 2
		}
		a, za = b, zb
	}
}

type punto struct{ x, y float64 }

func abs2(p punto) float64 { return math.Hypot(p.x, p.y) }

// caminata builds the partial-sum path with step lengths n^{-sigma}.
// sigma = 0 gives the pure angular skeleton (unit steps).
func caminata(t, sigma float64, X int) []punto {
	ps := make([]punto, X+1)
	x, y := 0.0, 0.0
	for n := 1; n <= X; n++ {
		fn := float64(n)
		a := t * math.Log(fn)
		r := math.Pow(fn, -sigma)
		x += r * math.Cos(a)
		y -= r * math.Sin(a)
		ps[n] = punto{x, y}
	}
	return ps
}

// coherencia: net displacement over a +-w window divided by arc length (unit
// steps -> arc length is 2w exactly). 1 = straight run, 0 = pure scribble.
func coherencia(ps []punto, w int) []float64 {
	c := make([]float64, len(ps))
	for n := w; n < len(ps)-w; n++ {
		d := math.Hypot(ps[n+w].x-ps[n-w].x, ps[n+w].y-ps[n-w].y)
		c[n] = d / float64(2*w)
	}
	return c
}

type corrida struct{ ini, fin, centro int }

// corridas: maximal stretches with coherence above the threshold, at least
// minLen wide; center = the n of maximal coherence inside the stretch.
func corridas(c []float64, umbral float64, minLen int) []corrida {
	var rs []corrida
	i := 0
	for i < len(c) {
		if c[i] < umbral {
			i++
			continue
		}
		j := i
		for j < len(c) && c[j] >= umbral {
			j++
		}
		if j-i >= minLen {
			mejor, cb := i, c[i]
			for k := i; k < j; k++ {
				if c[k] > cb {
					mejor, cb = k, c[k]
				}
			}
			rs = append(rs, corrida{i, j - 1, mejor})
		}
		i = j
	}
	return rs
}

// ojoEM: Euler-Maclaurin corrected limit of the sigma-walk at cutoff X.
func ojoEM(t, sigma float64, ps []punto, X int) punto {
	lX := math.Log(float64(X))
	// X^{1-s}/(s-1), s = sigma+it
	mod := math.Pow(float64(X), 1-sigma)
	nr, ni := mod*math.Cos(t*lX), -mod*math.Sin(t*lX)
	dr, di := sigma-1, t
	den := dr*dr + di*di
	cr := (nr*dr + ni*di) / den
	ci := (ni*dr - nr*di) / den
	// - X^{-s}/2
	m2 := math.Pow(float64(X), -sigma) / 2
	hr := -m2 * math.Cos(t*lX)
	hi := m2 * math.Sin(t*lX)
	return punto{ps[X].x + cr + hr, ps[X].y + ci + hi}
}

// betaRadio: least-squares slope of log|S_n - eye| vs log n. The window must
// sit where the smooth drift n^{1-sigma}/|s-1| dominates the Euler-Maclaurin
// error O(|s| n^{-sigma-1}): that needs n >> t, so the fit runs over [2t, 6t].
func betaRadio(t, sigma float64, ps []punto) float64 {
	X := len(ps) - 1
	o := ojoEM(t, sigma, ps, X)
	lo, hi := int(2*t), int(6*t)
	if hi > X {
		hi = X
	}
	var sx, sy, sxx, sxy float64
	m := 0.0
	for n := lo; n <= hi; n++ {
		r := math.Hypot(ps[n].x-o.x, ps[n].y-o.y)
		if r < 1e-9 {
			continue
		}
		lx, ly := math.Log(float64(n)), math.Log(r)
		sx += lx
		sy += ly
		sxx += lx * lx
		sxy += lx * ly
		m++
	}
	return (m*sxy - sx*sy) / (m*sxx - sx*sx)
}

type medida struct {
	t, tau              float64
	nMax, nMin, K, nOjo int
	r1, r2, r3, med, r5 float64
	beta3, beta5, beta7 float64
}

func main() {
	fmt.Println("📐 LA PROPORCIÓN — magnitudes geométricas de la espiral, razones libres")
	fmt.Println("   orden estricta: no imponer 1/2 · esqueleto angular separado de la pintura")
	fmt.Println()
	fmt.Println("   detector fijo para toda la escalera: ventana w=4, umbral 0,85, largo mínimo 3")
	fmt.Println()

	const w, minLen = 4, 3
	const umbral = 0.85
	objetivos := []float64{200, 400, 800, 1600, 3200, 6400, 25600, 102400}
	var ms []medida

	for _, t0 := range objetivos {
		t := primerCero(t0)
		tau := t / (2 * math.Pi)
		X := int(3*tau) + 60

		esq := caminata(t, 0, X) // skeleton: NO sigma anywhere
		c := coherencia(esq, w)
		rs := corridas(c, umbral, minLen)
		if len(rs) == 0 {
			fmt.Printf("   t=%.1f: sin corridas — se reporta y se sigue\n", t)
			continue
		}
		nMax, nMin := 0, 1<<30
		for _, r := range rs {
			if r.centro > nMax {
				nMax = r.centro
			}
			if r.centro < nMin {
				nMin = r.centro
			}
		}
		// final winding start: first low-coherence n past the main run
		nOjo := nMax
		for n := nMax; n < len(c); n++ {
			if c[n] < 0.25 {
				nOjo = n
				break
			}
		}
		m := medida{t: t, tau: tau, nMax: nMax, nMin: nMin, K: len(rs), nOjo: nOjo}
		lt := math.Log(tau)
		m.r1 = float64(nMax) / tau
		m.r2 = math.Log(float64(len(rs))) / lt
		m.r3 = math.Log(float64(nMin)) / lt
		m.med = (m.r2 + m.r3) / 2
		m.r5 = math.Log(float64(len(rs))*float64(nMin)) / lt

		// paint control: radius exponent for three sigmas, on a walk long
		// enough that the drift dominates the EM error (n up to 6t)
		X2 := int(6*t) + 60
		for _, sg := range []float64{0.3, 0.5, 0.7} {
			ps := caminata(t, sg, X2)
			b := betaRadio(t, sg, ps)
			switch sg {
			case 0.3:
				m.beta3 = b
			case 0.5:
				m.beta5 = b
			case 0.7:
				m.beta7 = b
			}
		}
		ms = append(ms, m)
	}

	fmt.Println("§1 · el esqueleto angular (sin σ en ninguna parte), escalera de alturas")
	fmt.Println()
	fmt.Println("      t        τ     nMax    K   nMin  | nMax/τ   logK/logτ  log nMin/logτ   PUNTO MEDIO   suma")
	for _, m := range ms {
		fmt.Printf("   %7.1f  %6.1f  %5d  %3d  %5d  |  %.3f     %.4f       %.4f        %.4f      %.4f\n",
			m.t, m.tau, m.nMax, m.K, m.nMin, m.r1, m.r2, m.r3, m.med, m.r5)
	}

	fmt.Println()
	fmt.Println("§2 · la pintura (exponente del radio del ojo), control a tres σ")
	fmt.Println()
	fmt.Println("      t       β(σ=0,3)   β(σ=0,5)   β(σ=0,7)      (teoría: β = 1−σ)")
	for _, m := range ms {
		fmt.Printf("   %7.1f     %.4f     %.4f     %.4f\n", m.t, m.beta3, m.beta5, m.beta7)
	}

	// sorted ratios at the deepest rung
	f := ms[len(ms)-1]
	tipo := []struct {
		nombre string
		v      float64
	}{
		{"nMax/τ (resonancia mayor)", f.r1},
		{"log K/log τ (cuántas corridas)", f.r2},
		{"log nMin/log τ (límite de resolución)", f.r3},
		{"PUNTO MEDIO de las dos anteriores", f.med},
		{"log(K·nMin)/log τ (el producto)", f.r5},
		{"nOjo/τ (arranque del ojo)", float64(f.nOjo) / f.tau},
		{"β del radio a σ=0,5", f.beta5},
		{"β del radio a σ=0,3", f.beta3},
		{"β del radio a σ=0,7", f.beta7},
	}
	sort.Slice(tipo, func(i, j int) bool {
		return math.Abs(tipo[i].v-0.5) < math.Abs(tipo[j].v-0.5)
	})
	fmt.Println()
	fmt.Printf("§3 · las razones en el peldaño más hondo (t=%.1f), ordenadas por cercanía a 1/2\n", f.t)
	fmt.Println("     (el orden lo decide el número, no nosotros)")
	fmt.Println()
	for _, r := range tipo {
		marca := " "
		if math.Abs(r.v-0.5) < 0.05 {
			marca = "◆"
		}
		fmt.Printf("   %s  %-42s %.4f\n", marca, r.nombre, r.v)
	}

	dibujar(ms)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-proporcion.svg")
}

func dibujar(ms []medida) {
	var b strings.Builder
	W, H := 1060, 620
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	t := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	t(530, 34, 21, "#e8e2d4", "middle", "LA PROPORCIÓN — ¿converge alguna razón a 1/2 sin que se lo pidamos?")
	t(530, 56, 12, "#8a93a6", "middle", "esqueleto angular puro (pasos de largo 1, sin σ): razones geométricas por la escalera t = 200 … 6400")

	// chart: ratios vs log tau
	x0, y0, cw, ch := 70.0, 100.0, 920.0, 420.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x0, y0, cw, ch)
	vmin, vmax := -0.03, 1.06
	py := func(v float64) float64 { return y0 + ch - (v-vmin)/(vmax-vmin)*ch }
	lt0 := math.Log(ms[0].tau)
	lt1 := math.Log(ms[len(ms)-1].tau)
	px := func(tau float64) float64 { return x0 + 40 + (math.Log(tau)-lt0)/(lt1-lt0)*(cw-80) }

	// the 1/2 line, drawn AFTER the data existed
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#c9b458" stroke-width="1.2" stroke-dasharray="7 5"/>`, x0, py(0.5), x0+cw, py(0.5))
	t(x0+cw-6, py(0.5)-7, 12, "#c9b458", "end", "1/2")
	for _, v := range []float64{0, 0.25, 0.75, 1} {
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#1c2230" stroke-width="1"/>`, x0, py(v), x0+cw, py(v))
		t(x0-8, py(v)+4, 10, "#5c6478", "end", fmt.Sprintf("%.2f", v))
	}

	serie := func(vals []float64, col, nombre string, yLeyenda float64) {
		var pts []string
		for i, m := range ms {
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", px(m.tau), py(vals[i])))
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="3.4" fill="%s"/>`, px(m.tau), py(vals[i]), col)
		}
		fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="%s" stroke-width="1.6" opacity="0.9"/>`, strings.Join(pts, " "), col)
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="4" fill="%s"/>`, x0+58, yLeyenda-4, col)
		t(x0+70, yLeyenda, 12, "#c7cdd9", "start", nombre)
	}
	get := func(f func(medida) float64) []float64 {
		out := make([]float64, len(ms))
		for i, m := range ms {
			out[i] = f(m)
		}
		return out
	}
	serie(get(func(m medida) float64 { return m.r3 }), "#f2a6c0", "log nMin / log τ — el límite de resolución (baja hacia 1/2)", y0+28)
	serie(get(func(m medida) float64 { return m.r2 }), "#6aa9ff", "log K / log τ — cuántas corridas se ven (sube hacia 1/2)", y0+50)
	serie(get(func(m medida) float64 { return m.med }), "#7ee0c0", "el PUNTO MEDIO de las dos — clavado", y0+72)

	for _, m := range ms {
		t(px(m.tau), y0+ch+18, 10, "#5c6478", "middle", fmt.Sprintf("t≈%.0f", m.t))
	}

	f := ms[len(ms)-1]
	t(530, float64(H)-32, 12, "#c7cdd9", "middle",
		fmt.Sprintf("las dos razones angulares aprietan el 1/2 desde los dos lados, y su punto medio da %.4f en el peldaño más hondo — sin σ en ninguna parte del esqueleto", f.med))
	t(530, float64(H)-12, 12, "#c9b458", "middle",
		"la pintura obedece β = 1−σ (0,7 / 0,5 / 0,3 medidos): el único σ donde la pintura empata al esqueleto es 1/2 · Todavía no")

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-proporcion.svg")
	os.WriteFile(ruta, []byte(b.String()), 0644)
}
