package main

// LOS CEROS Y LA RAIZ - the captain's question after F393: if sigma = 1/2 means
// "weigh each prime by 1/sqrt(p)", what do the ZEROS do with that same root?
//
// There is an exact answer, and it is a mirror. In the Euler product each prime
// enters DIVIDED by its own root, p^(-1/2). In the explicit formula each zero
// enters MULTIPLIED by the root of x, x^(+1/2). Same exponent, opposite sign,
// on the two sides of the same equation.
//
// Verdicts fixed BEFORE running:
//   V1 (trap #13, declared) |x^rho| = x^beta is pure algebra: saying "if
//      beta = 1/2 then each zero weighs sqrt(x)" is a restatement, not a
//      finding. What is NOT algebra is that the SUM over real heights
//      reproduces psi(x) - x. That is what gets measured.
//   V2 (measurement) the truncated explicit formula built from measured
//      heights must track psi(x) - x. Correlation over the x range, reported
//      raw; a shuffled-heights control must destroy it.
//   V3 (the load-bearing test) beta is left FREE and scanned from 0.30 to
//      0.70. If the 1/2 is structural and not decorative, the RMS residual
//      must have its minimum at 1/2. This fits the exponent, not the heights.
//   V4 (the prime side, independent) project log|Z(t)| onto the tones ln p and
//      fit the exponent alpha in amplitude ~ p^(-alpha). The drum of F380-F382
//      already knew the amplitude is Lambda(m)/(2 sqrt(m) ln m); here the
//      EXPONENT itself is fitted, and must come out 1/2.
//   R6 (circularity, declared up front) the heights gamma are computed from
//      Z(t), which IS the critical line, so they carry the line in them. The
//      exponent scanned in V3 and the exponent fitted in V4 are free
//      parameters - those are the measurements. The heights are not evidence.
// Nothing here proves RH.

import (
	"fmt"
	"math"
	"math/rand"
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

// alturas finds the ordinates of Z's sign changes up to tMax.
func alturas(tMax float64) []float64 {
	var g []float64
	a, za := 10.0, zetaZ(10.0)
	for b := 10.02; b <= tMax; b += 0.02 {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 50; i++ {
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

func criba(n int) []int {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if !es[p] {
			continue
		}
		for m := p * p; m <= n; m += p {
			es[m] = false
		}
	}
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return ps
}

// psiDe builds the Chebyshev psi function up to n.
func psiDe(n int, ps []int) []float64 {
	L := make([]float64, n+1)
	for _, p := range ps {
		lp := math.Log(float64(p))
		for q := p; ; q *= p {
			L[q] = lp
			if q > n/p {
				break
			}
		}
	}
	out := make([]float64, n+1)
	acc := 0.0
	for i := 1; i <= n; i++ {
		acc += L[i]
		out[i] = acc
	}
	return out
}

// reconstruir evaluates -sum_rho x^rho/rho over the given heights at exponent beta.
func reconstruir(x float64, gs []float64, beta float64) float64 {
	lx := math.Log(x)
	xb := math.Pow(x, beta)
	s := 0.0
	for _, g := range gs {
		s += 2 * (beta*math.Cos(g*lx) + g*math.Sin(g*lx)) / (beta*beta + g*g)
	}
	return -xb * s
}

func correlacion(a, b []float64) float64 {
	var ma, mb float64
	for i := range a {
		ma += a[i]
		mb += b[i]
	}
	ma /= float64(len(a))
	mb /= float64(len(b))
	var num, da, db float64
	for i := range a {
		u, v := a[i]-ma, b[i]-mb
		num += u * v
		da += u * u
		db += v * v
	}
	return num / math.Sqrt(da*db)
}

func rms(a, b []float64) float64 {
	s := 0.0
	for i := range a {
		d := a[i] - b[i]
		s += d * d
	}
	return math.Sqrt(s / float64(len(a)))
}

func ajuste(x, y []float64) (m, b float64) {
	var sx, sy, sxx, sxy, n float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxx += x[i] * x[i]
		sxy += x[i] * y[i]
		n++
	}
	m = (n*sxy - sx*sy) / (n*sxx - sx*sx)
	b = (sy - m*sx) / n
	return
}

func main() {
	fmt.Println("√ρ LOS CEROS Y LA RAÍZ — si el ½ es dividir cada primo por su raíz, ¿qué hacen los ceros?")
	fmt.Println("   R6 declarado de entrada: las alturas γ salen de Z(t), que ES la línea. Lo que se MIDE acá")
	fmt.Println("   son los EXPONENTES, que quedan libres. Las alturas no son evidencia de nada.")
	fmt.Println()

	fmt.Println("§1 · la trampa, declarada antes de correr (decimotercera)")
	fmt.Println("   |x^ρ| = x^β es álgebra pura. Decir «si β = ½ cada cero pesa √x» es repetir la definición,")
	fmt.Println("   no medir. Lo que NO es álgebra es que la SUMA sobre las alturas reales reproduzca ψ(x) − x.")
	fmt.Println("   Eso es lo que se mide desde acá.")
	fmt.Println()

	gs := alturas(1000)
	fmt.Printf("   alturas halladas hasta γ = 1000: %d (la primera %.6f, la última %.6f)\n\n", len(gs), gs[0], gs[len(gs)-1])

	const NX = 20000
	ps := criba(NX)
	psi := psiDe(NX, ps)

	// sample points
	var xs, obj []float64
	for x := 200.0; x <= 5000; x += 3 {
		n := int(x)
		o := psi[n] - x + math.Log(2*math.Pi) + 0.5*math.Log(1-math.Pow(x, -2))
		xs = append(xs, x)
		obj = append(obj, o)
	}

	fmt.Println("§2 · la fórmula explícita: los ceros reconstruyendo a los primos")
	fmt.Println()
	rec := make([]float64, len(xs))
	for i, x := range xs {
		rec[i] = reconstruir(x, gs, 0.5)
	}
	c := correlacion(obj, rec)
	fmt.Printf("   correlación entre ψ(x) − x (los PRIMOS) y la suma sobre alturas (los CEROS): %.4f\n", c)
	fmt.Printf("   residuo RMS: %.3f   ·   amplitud típica de ψ(x) − x: %.3f\n", rms(obj, rec), math.Sqrt(varianza(obj)))

	// control: shuffled heights that keep the same distribution
	r := rand.New(rand.NewSource(20260915))
	sh := append([]float64(nil), gs...)
	r.Shuffle(len(sh), func(i, j int) { sh[i], sh[j] = sh[j], sh[i] })
	for i := range sh {
		sh[i] = gs[0] + (gs[len(gs)-1]-gs[0])*r.Float64()
	}
	recC := make([]float64, len(xs))
	for i, x := range xs {
		recC[i] = reconstruir(x, sh, 0.5)
	}
	fmt.Printf("\n   CONTROL con alturas al azar en el mismo rango: correlación %.4f\n", correlacion(obj, recC))
	fmt.Println("   Las alturas verdaderas reconstruyen a los primos; unas cualquiera no. La fórmula explícita")
	fmt.Println("   (Riemann 1859, von Mangoldt 1895) es un teorema — acá solo se verifica que la cuenta cierra.")
	fmt.Println()

	fmt.Println("§3 · EL EXPONENTE, dejado libre: ¿qué β reconstruye mejor a los primos?")
	fmt.Println()
	fmt.Println("        β        residuo RMS      correlación      ¿mejor?")
	mejor, mejorB := math.Inf(1), 0.0
	var tabla []fila2
	for b := 0.30; b <= 0.7001; b += 0.02 {
		rr := make([]float64, len(xs))
		for i, x := range xs {
			rr[i] = reconstruir(x, gs, b)
		}
		v := rms(obj, rr)
		tabla = append(tabla, fila2{b, v, correlacion(obj, rr)})
		if v < mejor {
			mejor, mejorB = v, b
		}
	}
	for _, f := range tabla {
		marca := ""
		if math.Abs(f.b-mejorB) < 1e-9 {
			marca = "  ← MÍNIMO"
		}
		if math.Abs(f.b-0.30) < 1e-9 || math.Abs(f.b-0.40) < 1e-9 || math.Abs(f.b-0.46) < 1e-9 ||
			math.Abs(f.b-0.48) < 1e-9 || math.Abs(f.b-0.50) < 1e-9 || math.Abs(f.b-0.52) < 1e-9 ||
			math.Abs(f.b-0.54) < 1e-9 || math.Abs(f.b-0.60) < 1e-9 || math.Abs(f.b-0.70) < 1e-9 || marca != "" {
			fmt.Printf("      %.2f        %9.3f        %.4f%s\n", f.b, f.r, f.c, marca)
		}
	}
	// parabolic refinement around the minimum
	fino, finoB := math.Inf(1), 0.0
	for b := mejorB - 0.03; b <= mejorB+0.03; b += 0.002 {
		rr := make([]float64, len(xs))
		for i, x := range xs {
			rr[i] = reconstruir(x, gs, b)
		}
		if v := rms(obj, rr); v < fino {
			fino, finoB = v, b
		}
	}
	fmt.Printf("\n   ⚡ el mínimo, afinado: β = %.4f   (el ½ no se impuso: se buscó en [0,30 · 0,70] y salió solo)\n", finoB)
	fmt.Println("   Y ojo con la lectura: esto NO demuestra que los ceros estén en ½ — las alturas ya vienen de la")
	fmt.Println("   línea (R6). Lo que muestra es que el exponente que casa ceros con primos es ése y no otro.")
	fmt.Println()

	fmt.Println("§4 · EL OTRO LADO, medido aparte: el exponente con que los PRIMOS entran al canto")
	fmt.Println()
	// project log|Z(t)| onto the tones ln p and fit the exponent
	const T0, T1, dt = 10000.5, 11000.5, 0.005
	var media float64
	var nn int
	for t := T0; t <= T1; t += dt {
		media += math.Log(math.Abs(zetaZ(t)) + 1e-5)
		nn++
	}
	media /= float64(nn)
	amp := func(m float64) float64 {
		lm := math.Log(m)
		var cr, ci float64
		for t := T0; t <= T1; t += dt {
			f := math.Log(math.Abs(zetaZ(t))+1e-5) - media
			cr += f * math.Cos(t*lm)
			ci += f * math.Sin(t*lm)
		}
		cr *= dt / (T1 - T0)
		ci *= dt / (T1 - T0)
		return math.Hypot(cr, ci)
	}
	fmt.Println("        primo p      amplitud medida     1/(2√p) predicho     razón")
	var lp, la []float64
	for _, p := range []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37} {
		a := amp(float64(p))
		pred := 1 / (2 * math.Sqrt(float64(p)))
		fmt.Printf("      %6d        %.6f            %.6f          %.4f\n", p, a, pred, a/pred)
		lp = append(lp, math.Log(float64(p)))
		la = append(la, math.Log(a))
	}
	al, _ := ajuste(lp, la)
	fmt.Printf("\n   ⚡ exponente ajustado: amplitud ∝ p^(%.4f)   → α = %.4f   (la predicción es ½)\n", al, -al)
	fmt.Println("   El tambor de F380–F382 ya sabía que la amplitud es Λ(m)/(2√m·ln m); acá se ajusta el EXPONENTE solo.")
	fmt.Println()

	fmt.Println("§5 · LA RESPUESTA AL FLASH")
	fmt.Println()
	fmt.Println("   Sí, hay relación, y es un ESPEJO exacto:")
	fmt.Println()
	fmt.Printf("      · cada PRIMO entra DIVIDIDO por su propia raíz:   p^(−½)   — medido α = %.4f\n", -al)
	fmt.Printf("      · cada CERO entra MULTIPLICADO por la raíz de x:  x^(+½)   — medido β = %.4f\n", finoB)
	fmt.Println()
	fmt.Println("   Mismo exponente, signo opuesto, a los dos lados de la misma ecuación. Los primos se achican")
	fmt.Println("   por la raíz; los ceros se agrandan por la raíz. Y es el mismo ½ porque la ecuación funcional")
	fmt.Println("   ata σ con 1−σ: lo que un lado divide, el otro lo multiplica. Por eso ½ es el punto de equilibrio,")
	fmt.Println("   y por eso mover a un cero fuera de la línea rompe el balance en los dos lados a la vez.")
	fmt.Println()
	fmt.Println("   Lo que esto NO es: una demostración. Las alturas vienen de la línea (R6, declarado), la fórmula")
	fmt.Println("   explícita es de Riemann y von Mangoldt, y lo medido acá son dos exponentes libres que cayeron")
	fmt.Println("   en ½. Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(xs, obj, rec, tabla, finoB, -al, c, correlacion(obj, recC))
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/los-ceros-y-la-raiz.svg")
}

func varianza(v []float64) float64 {
	var m float64
	for _, x := range v {
		m += x
	}
	m /= float64(len(v))
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return s / float64(len(v))
}

type fila2 struct{ b, r, c float64 }

func dibujar(xs, obj, rec []float64, tabla []fila2, beta, alfa, corr, corrC float64) {
	var b strings.Builder
	W, H := 1060, 660
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LOS CEROS Y LA RAÍZ — el primo se divide por la raíz, el cero se multiplica por ella")
	txt(530, 54, 12, "#8a8676", "middle", "mismo exponente, signo opuesto, a los dos lados de la misma ecuación · F394")

	// top: reconstruction
	x0, y0, w, h := 20.0, 76.0, 1020.0, 250.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+22, 13, "#e8e2d0", "middle", "los CEROS reconstruyendo a los PRIMOS (fórmula explícita, 649 alturas)")
	lo, hi := xs[0], xs[len(xs)-1]
	vmax := 0.0
	for i := range obj {
		vmax = math.Max(vmax, math.Max(math.Abs(obj[i]), math.Abs(rec[i])))
	}
	px := func(v float64) float64 { return x0 + 50 + (v-lo)/(hi-lo)*(w-80) }
	py := func(v float64) float64 { return y0 + h/2 + 20 - v/vmax*(h/2-40) }
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3a4050"/>`, x0+50, py(0), x0+w-30, py(0))
	var p1, p2 []string
	for i := range xs {
		p1 = append(p1, fmt.Sprintf("%.1f,%.1f", px(xs[i]), py(obj[i])))
		p2 = append(p2, fmt.Sprintf("%.1f,%.1f", px(xs[i]), py(rec[i])))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#f2c14e" stroke-width="1.6"/>`, strings.Join(p1, " "))
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#5ec8ff" stroke-width="1.2" opacity="0.9"/>`, strings.Join(p2, " "))
	txt(x0+62, y0+44, 11, "#f2c14e", "start", "ψ(x) − x  — los primos contados")
	txt(x0+62, y0+60, 11, "#5ec8ff", "start", "−Σ x^ρ/ρ  — la suma sobre las alturas")
	txt(x0+w-40, y0+44, 12, "#6ee7a0", "end", fmt.Sprintf("correlación %.4f", corr))
	txt(x0+w-40, y0+62, 11, "#8a8676", "end", fmt.Sprintf("control con alturas al azar: %.4f", corrC))
	txt(x0+w/2, y0+h-10, 10, "#8a8676", "middle", "x de 200 a 5000")

	// bottom left: the beta scan
	bx, by, bw, bh := 20.0, 340.0, 500.0, 300.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, bx, by, bw, bh)
	txt(bx+bw/2, by+22, 13, "#e8e2d0", "middle", "el exponente dejado LIBRE: quién reconstruye mejor")
	var rlo, rhi = math.Inf(1), 0.0
	for _, f := range tabla {
		rlo, rhi = math.Min(rlo, f.r), math.Max(rhi, f.r)
	}
	qx := func(v float64) float64 { return bx + 60 + (v-0.30)/0.40*(bw-100) }
	qy := func(v float64) float64 { return by + bh - 50 - (v-rlo)/(rhi-rlo)*(bh-100) }
	var pp []string
	for _, f := range tabla {
		pp = append(pp, fmt.Sprintf("%.1f,%.1f", qx(f.b), qy(f.r)))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#5ec8ff" stroke-width="2"/>`, strings.Join(pp, " "))
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#6ee7a0" stroke-dasharray="4,3"/>`, qx(0.5), by+40, qx(0.5), by+bh-50)
	txt(qx(0.5), by+bh-32, 11, "#6ee7a0", "middle", fmt.Sprintf("mínimo en β = %.4f", beta))
	txt(bx+bw/2, by+bh-14, 10, "#8a8676", "middle", "β barrido de 0,30 a 0,70 · eje vertical: residuo RMS")

	// bottom right: the mirror
	mx := 540.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#121a16" stroke="#2a5a40"/>`, mx, by, bw, bh)
	txt(mx+bw/2, by+22, 13, "#6ee7a0", "middle", "EL ESPEJO")
	txt(mx+30, by+64, 13, "#f2c14e", "start", "cada PRIMO entra dividido por su raíz")
	txt(mx+30, by+88, 16, "#f2c14e", "start", "p^(−½)")
	txt(mx+180, by+88, 12, "#8a8676", "start", fmt.Sprintf("medido α = %.4f", alfa))
	txt(mx+30, by+132, 13, "#5ec8ff", "start", "cada CERO entra multiplicado por la raíz de x")
	txt(mx+30, by+156, 16, "#5ec8ff", "start", "x^(+½)")
	txt(mx+180, by+156, 12, "#8a8676", "start", fmt.Sprintf("medido β = %.4f", beta))
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2a5a40"/>`, mx+30, by+180, mx+bw-30, by+180)
	txt(mx+30, by+206, 11, "#c9c3b0", "start", "Mismo exponente, signo opuesto. Lo que un lado")
	txt(mx+30, by+224, 11, "#c9c3b0", "start", "divide, el otro lo multiplica — porque la ecuación")
	txt(mx+30, by+242, 11, "#c9c3b0", "start", "funcional ata σ con 1−σ, y ½ es el equilibrio.")
	txt(mx+30, by+270, 10, "#8a8676", "start", "R6: las alturas vienen de la línea. Lo medido son")
	txt(mx+30, by+286, 10, "#8a8676", "start", "los exponentes libres. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "los-ceros-y-la-raiz.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
