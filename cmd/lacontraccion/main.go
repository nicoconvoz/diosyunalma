package main

// LA CONTRACCION - Auditoria 58, strict order: do NOT look for 1/2. Measure
// the spiral. Decisions fixed BEFORE looking (her §2-§3):
//   trajectory: painted walk S_n = sum n^{-sigma-it} at height t, n = 1..N
//   center C:   the Euler-Maclaurin eye at cutoff N (fixed criterion)
//   cycles:     three separate definitions, all run:
//                 A1 full revolutions of the unwrapped angle around C
//                 A2 successive radial maxima
//                 A3 successive radial minima
//   size r_k:   mean radial distance over the cycle (max and min also stored)
//   ratios:     q_k = r_k / r_{k-1}
// Cases (her §7 and §11): real zero, near-zero (+0.15), mid-gap control,
// a second distant zero, and two instrument controls (sigma = 0.3 / 0.7).
// Robustness: N in {4tau, 8tau, 16tau}; scale control x3.
// The numbers 0.5 and 17/50 appear ONLY in the final blind section.
// Raw data written to docs/atomo/CONTRACCION-CRUDOS.csv.

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

// trayectoria: radial distance to C and unwrapped angle, n = 1..N.
func trayectoria(ps []punto, C punto) (r, th []float64) {
	r = make([]float64, len(ps))
	th = make([]float64, len(ps))
	prev := 0.0
	for n := 1; n < len(ps); n++ {
		dx, dy := ps[n].x-C.x, ps[n].y-C.y
		r[n] = math.Hypot(dx, dy)
		a := math.Atan2(dy, dx)
		if n == 1 {
			th[n] = a
		} else {
			d := a - prev
			for d > math.Pi {
				d -= 2 * math.Pi
			}
			for d < -math.Pi {
				d += 2 * math.Pi
			}
			th[n] = th[n-1] + d
		}
		prev = a
	}
	return
}

// ciclosAngulo: boundaries each time the ABSOLUTE winding advances 2*pi.
func ciclosAngulo(th []float64) []int {
	var bs []int
	base := th[1]
	prev := 0
	for n := 2; n < len(th); n++ {
		k := int(math.Floor(math.Abs(th[n]-base) / (2 * math.Pi)))
		if k > prev {
			bs = append(bs, n)
			prev = k
		}
	}
	return bs
}

// ciclosPicos: boundaries at strict local maxima (modo=+1) or minima (-1) of r,
// window +-2.
func ciclosPicos(r []float64, modo float64) []int {
	var bs []int
	for n := 3; n < len(r)-2; n++ {
		es := true
		for _, d := range []int{-2, -1, 1, 2} {
			if modo*(r[n]-r[n+d]) <= 0 {
				es = false
				break
			}
		}
		if es {
			bs = append(bs, n)
		}
	}
	return bs
}

func tamanos(r []float64, bs []int) (mean, mx, mn []float64) {
	for i := 0; i+1 < len(bs); i++ {
		a, b := bs[i], bs[i+1]
		s, hi, lo := 0.0, 0.0, math.Inf(1)
		for n := a; n < b; n++ {
			s += r[n]
			hi = math.Max(hi, r[n])
			lo = math.Min(lo, r[n])
		}
		mean = append(mean, s/float64(b-a))
		mx = append(mx, hi)
		mn = append(mn, lo)
	}
	return
}

func razones(rs []float64) []float64 {
	var q []float64
	for i := 1; i < len(rs); i++ {
		q = append(q, rs[i]/rs[i-1])
	}
	return q
}

func mediana(v []float64) float64 {
	if len(v) == 0 {
		return math.NaN()
	}
	w := append([]float64(nil), v...)
	sort.Float64s(w)
	return w[len(w)/2]
}

// betaExterior: slope of ln r vs ln n where the smooth drift DOMINATES the
// Euler-Maclaurin error - that needs n >> t (the F376 lesson: fitting over
// [2tau, 8tau] measures the error term, not the law). A separate long walk is
// built and the fit runs over [2t, 6t].
func betaExterior(t, sigma float64) float64 {
	N := int(6 * t)
	ps := walk(t, sigma, N)
	C := ojoEM(t, sigma, ps, N)
	lo, hi := int(2*t), N
	var sx, sy, sxx, sxy, m float64
	for n := lo; n < hi; n++ {
		r := math.Hypot(ps[n].x-C.x, ps[n].y-C.y)
		if r < 1e-12 {
			continue
		}
		x, y := math.Log(float64(n)), math.Log(r)
		sx += x
		sy += y
		sxx += x * x
		sxy += x * y
		m++
	}
	return (m*sxy - sx*sy) / (m*sxx - sx*sx)
}

type caso struct {
	nom      string
	t, sigma float64
	tau      float64
	qAng     []float64
	tamAng   []float64
	medAng   float64
	medMax   float64
	medMin   float64
	beta     float64
	nCiclos  int
	rangoQ   [2]float64
}

func main() {
	fmt.Println("🌀 LA CONTRACCIÓN — medir la espiral, no buscar nada (Auditoría 58)")
	fmt.Println("   centro = ojo EM (fijo) · tres definiciones de vuelta · tamaño = radio medio")
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
		{"instrumento σ=0,3 (misma altura)", g1, 0.3},
		{"instrumento σ=0,7 (misma altura)", g1, 0.7},
	}

	var csv strings.Builder
	csv.WriteString("caso,t,sigma,criterio,ciclo,r_k,r_prev,q_k\n")
	var res []caso

	for _, cs := range casos {
		tau := cs.t / (2 * math.Pi)
		N := int(8 * tau)
		ps := walk(cs.t, cs.sigma, N)
		C := ojoEM(cs.t, cs.sigma, ps, N)
		r, th := trayectoria(ps, C)

		bA := ciclosAngulo(th)
		mA, _, _ := tamanos(r, bA)
		qA := razones(mA)
		bMx := ciclosPicos(r, +1)
		mMx, _, _ := tamanos(r, bMx)
		qMx := razones(mMx)
		bMn := ciclosPicos(r, -1)
		mMn, _, _ := tamanos(r, bMn)
		qMn := razones(mMn)

		for i, q := range qA {
			csv.WriteString(fmt.Sprintf("%s,%.4f,%.1f,angular,%d,%.6f,%.6f,%.6f\n",
				cs.nom, cs.t, cs.sigma, i+1, mA[i+1], mA[i], q))
		}
		lo, hi := math.Inf(1), math.Inf(-1)
		for _, q := range qA {
			lo, hi = math.Min(lo, q), math.Max(hi, q)
		}
		res = append(res, caso{cs.nom, cs.t, cs.sigma, tau, qA, mA,
			mediana(qA), mediana(qMx), mediana(qMn),
			betaExterior(cs.t, cs.sigma), len(qA), [2]float64{lo, hi}})
	}

	fmt.Println("§1 · LA TABLA OBLIGATORIA (criterio angular, primeros 6 ciclos por caso)")
	fmt.Println("   (las columnas de comparación llegan recién en §5, como manda la orden)")
	fmt.Println()
	fmt.Println("   caso                              ciclo    r_k        r_{k-1}     q_k")
	for _, c := range res {
		tope := 6
		if len(c.qAng) < tope {
			tope = len(c.qAng)
		}
		for i := 0; i < tope; i++ {
			fmt.Printf("   %-32s   %2d    %.6f   %.6f   %.4f\n", c.nom, i+1, c.tamAng[i+1], c.tamAng[i], c.qAng[i])
		}
	}

	fmt.Println()
	fmt.Println("§2 · resumen por caso y por criterio (mediana de q_k)")
	fmt.Println()
	fmt.Println("   caso                               ciclos   q̃ angular   q̃ máximos   q̃ mínimos   rango q (angular)")
	for _, c := range res {
		fmt.Printf("   %-32s   %4d     %.4f      %.4f      %.4f      [%.3f, %.3f]\n",
			c.nom, c.nCiclos, c.medAng, c.medMax, c.medMin, c.rangoQ[0], c.rangoQ[1])
	}

	fmt.Println()
	fmt.Println("§3 · robustez con N (cero γ≈1000, criterio angular)")
	tau := g1 / (2 * math.Pi)
	for _, mult := range []float64{4, 8, 16} {
		N := int(mult * tau)
		ps := walk(g1, 0.5, N)
		C := ojoEM(g1, 0.5, ps, N)
		r, th := trayectoria(ps, C)
		mA, _, _ := tamanos(r, ciclosAngulo(th))
		q := razones(mA)
		fmt.Printf("   N = %2.0fτ: %3d ciclos · q̃ = %.4f\n", mult, len(q), mediana(q))
	}
	// scale control x3
	{
		N := int(8 * tau)
		ps := walk(g1, 0.5, N)
		for i := range ps {
			ps[i].x *= 3
			ps[i].y *= 3
		}
		C := ojoEM(g1, 0.5, walk(g1, 0.5, N), N)
		C.x *= 3
		C.y *= 3
		r, th := trayectoria(ps, C)
		mA, _, _ := tamanos(r, ciclosAngulo(th))
		fmt.Printf("   escala ×3: q̃ = %.4f (control de su §10: la razón no cambia)\n", mediana(razones(mA)))
	}

	fmt.Println()
	fmt.Println("§4 · el exponente del enrollado exterior (ajuste en [2t, 6t], donde la")
	fmt.Println("   deriva domina al error — la lección de F376, aplicada)")
	for _, c := range res {
		fmt.Printf("   %-32s   β = %.4f   ⟹ razón por vuelta e^(2πβ/t) = %.6f\n",
			c.nom, c.beta, math.Exp(2*math.Pi*c.beta/c.t))
	}

	fmt.Println()
	fmt.Println("§5 · RECIÉN AHORA — la comparación ciega (0,5 y 17/50 = 0,34)")
	fmt.Println()
	fmt.Println("   caso                               q̃ angular   |q̃−0,5|   |q̃−0,34|")
	for _, c := range res {
		fmt.Printf("   %-32s    %.4f     %.4f    %.4f\n", c.nom, c.medAng, math.Abs(c.medAng-0.5), math.Abs(c.medAng-0.34))
	}
	fmt.Println()
	fmt.Println("   y el lugar donde un medio SÍ vive en esta espiral, medido en §4:")
	for _, c := range res {
		fmt.Printf("   %-32s   β = %.4f   (1−σ = %.1f)\n", c.nom, c.beta, 1-c.sigma)
	}

	os.WriteFile(filepath.Join("docs", "atomo", "CONTRACCION-CRUDOS.csv"), []byte(csv.String()), 0644)
	fmt.Println("\n📄 datos crudos: docs/atomo/CONTRACCION-CRUDOS.csv")

	dibujar(res)
	fmt.Println("🖼️  lámina escrita: galeria/laminas/10-el-telar/la-contraccion.svg")
}

func dibujar(res []caso) {
	var b strings.Builder
	W, H := 1060, 620
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	t := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	t(530, 34, 21, "#e8e2d4", "middle", "LA CONTRACCIÓN — las razones por vuelta, medidas a ciegas")
	t(530, 56, 12, "#8a93a6", "middle", "q_k = r_k/r_{k−1}, criterio angular · sin ninguna línea de referencia dibujada")

	x0, y0, cw, ch := 70.0, 100.0, 920.0, 430.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x0, y0, cw, ch)
	vmin, vmax := 0.0, 2.0
	py := func(v float64) float64 { return y0 + ch - (v-vmin)/(vmax-vmin)*ch }
	for _, g := range []float64{0.25, 0.75, 1.0, 1.25, 1.75} {
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#1c2230"/>`, x0, py(g), x0+cw, py(g))
		t(x0-8, py(g)+4, 10, "#5c6478", "end", fmt.Sprintf("%.2f", g))
	}
	cols := []string{"#7ee0c0", "#6aa9ff", "#f2a6c0", "#c9b458", "#f0a05a", "#e06a6a"}
	for ci, c := range res {
		if ci >= len(cols) {
			break
		}
		nq := len(c.qAng)
		if nq == 0 {
			continue
		}
		for i, q := range c.qAng {
			x := x0 + 30 + (float64(i)+0.5)/float64(nq)*(cw-60)
			qq := math.Max(vmin, math.Min(q, vmax))
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.6" fill="%s" opacity="0.75"/>`, x, py(qq), cols[ci])
		}
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="4" fill="%s"/>`, x0+50, y0+22+float64(ci)*20-4, cols[ci])
		t(x0+62, y0+22+float64(ci)*20, 11, "#c7cdd9", "start", fmt.Sprintf("%s · mediana %.3f", c.nom, c.medAng))
	}
	t(x0+cw/2, y0+ch+22, 12, "#8a93a6", "middle", "ciclo k (normalizado por caso) → q_k · los seis casos superpuestos")

	t(530, float64(H)-32, 12, "#c7cdd9", "middle",
		"la comparación con 0,5 y 0,34 vive solo en el §5 de la salida — la lámina muestra los datos desnudos")
	t(530, float64(H)-12, 12, "#c9b458", "middle",
		"objetivo real: descubrir la ley, no demostrar el flash · Todavía no")

	b.WriteString(`</svg>`)
	os.WriteFile(filepath.Join("galeria", "laminas", "10-el-telar", "la-contraccion.svg"), []byte(b.String()), 0644)
}
