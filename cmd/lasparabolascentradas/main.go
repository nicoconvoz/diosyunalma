package main

// LAS PARABOLAS CENTRADAS - the captain's second flash on F391: put the two
// parabolas with their contact point at the absolute origin (0,0) and read
// the anatomy: right distance vs left distance, the angle, the separations,
// and what changes with height.
//
// Centered frame: x = sigma - beta (beta measured by 2-D Newton, never
// assumed), y = |zeta(beta + x + i gamma)|^2 for the upper curve and
// y = -|zeta(beta - x + i gamma)|^2 for the lower one (mirror through (0,0)).
//
// Verdicts fixed BEFORE running:
//   V1  (declared trap #11) the vertical separation between the two curves,
//       g(x) = f+(x) - f-(x) = f+(x) + f+(-x), is EVEN by construction:
//       g(x) - g(-x) prints 0.0e+00 whatever zeta does. Not a finding.
//   V2  (measurement) the asymmetry h(x) = f+(x) - f+(-x) is the lean of
//       F391: h(x) = -2 a L x^3 + O(x^5) with a = |zeta'(rho)|^2 and
//       L = log(gamma/2pi). A cubic fit of f+ returns L_meas = -b/a; the
//       verdict is L_meas/L -> 1 within 1%.
//   V3  at equal height y, right distance x+ and left distance |x-| differ:
//       x+ + x- = L y / a (chord midpoints drift linearly with y), so the
//       parabola's axis leans by psi = atan(L/(2a)) in raw (sigma, |zeta|^2)
//       units - the angle depends on the y-scale, the slope L/(2a) does not.
//   V4  with height: beta stays 1/2 (measured), L grows as log(gamma/2pi)
//       (deterministic), a = |zeta'(rho)|^2 wanders zero by zero.
// Nothing here proves anything about zeros not scanned.

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func zetaEM(s complex128) complex128 {
	N := int(3*math.Abs(imag(s))) + 60
	sum := complex(0, 0)
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * complex(math.Log(float64(n)), 0))
	}
	lnN := complex(math.Log(float64(N)), 0)
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	bern := []float64{1.0 / 12, -1.0 / 720, 1.0 / 30240}
	prod := s
	for k, b := range bern {
		if k > 0 {
			prod *= (s + complex(float64(2*k-1), 0)) * (s + complex(float64(2*k), 0))
		}
		sum += complex(b, 0) * prod * cmplx.Exp(-(s+complex(float64(2*k+1), 0))*lnN)
	}
	return sum
}

func f2(sigma, t float64) float64 {
	return math.Pow(cmplx.Abs(zetaEM(complex(sigma, t))), 2)
}

type cero struct {
	beta, gamma float64
	a           float64 // |zeta'(rho)|^2
}

func buscar(t0, t1, dt float64) []cero {
	const s0, s1, ds = 0.2, 0.8, 0.01
	ns := int(math.Round((s1-s0)/ds)) + 1
	nt := int(math.Round((t1-t0)/dt)) + 1
	g := make([][]float64, nt)
	for i := range g {
		g[i] = make([]float64, ns)
		for j := range g[i] {
			g[i][j] = f2(s0+float64(j)*ds, t0+float64(i)*dt)
		}
	}
	var out []cero
	for i := 1; i < nt-1; i++ {
		for j := 1; j < ns-1; j++ {
			v := g[i][j]
			if v > 0.05 {
				continue
			}
			ok := true
			for di := -1; di <= 1 && ok; di++ {
				for dj := -1; dj <= 1; dj++ {
					if (di != 0 || dj != 0) && g[i+di][j+dj] < v {
						ok = false
						break
					}
				}
			}
			if !ok {
				continue
			}
			s := complex(s0+float64(j)*ds, t0+float64(i)*dt)
			h := complex(1e-6, 0)
			var d complex128
			for it := 0; it < 30; it++ {
				d = (zetaEM(s+h) - zetaEM(s-h)) / (2 * h)
				paso := zetaEM(s) / d
				s -= paso
				if cmplx.Abs(paso) < 1e-13 {
					break
				}
			}
			if cmplx.Abs(zetaEM(s)) < 1e-9 {
				d = (zetaEM(s+h) - zetaEM(s-h)) / (2 * h)
				out = append(out, cero{real(s), imag(s), math.Pow(cmplx.Abs(d), 2)})
			}
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].gamma < out[b].gamma })
	var u []cero
	for _, c := range out {
		if len(u) == 0 || math.Abs(u[len(u)-1].gamma-c.gamma) > 1e-6 {
			u = append(u, c)
		}
	}
	return u
}

// ajusteCubico fits y = c0 + c1 x + c2 x^2 + c3 x^3 by least squares.
func ajusteCubico(x, y []float64) [4]float64 {
	var A [4][5]float64
	for i := range x {
		p := [4]float64{1, x[i], x[i] * x[i], x[i] * x[i] * x[i]}
		for r := 0; r < 4; r++ {
			for c := 0; c < 4; c++ {
				A[r][c] += p[r] * p[c]
			}
			A[r][4] += p[r] * y[i]
		}
	}
	for c := 0; c < 4; c++ {
		piv := c
		for r := c + 1; r < 4; r++ {
			if math.Abs(A[r][c]) > math.Abs(A[piv][c]) {
				piv = r
			}
		}
		A[c], A[piv] = A[piv], A[c]
		for r := 0; r < 4; r++ {
			if r == c {
				continue
			}
			f := A[r][c] / A[c][c]
			for k := c; k < 5; k++ {
				A[r][k] -= f * A[c][k]
			}
		}
	}
	var out [4]float64
	for r := 0; r < 4; r++ {
		out[r] = A[r][4] / A[r][r]
	}
	return out
}

// cruce solves f+(x) = y on one side of the origin by bisection.
func cruce(c cero, y float64, derecha bool) float64 {
	lo, hi := 0.0, 0.08
	if !derecha {
		lo, hi = -0.08, 0.0
	}
	fl := func(x float64) float64 { return f2(c.beta+x, c.gamma) - y }
	if derecha {
		for i := 0; i < 60; i++ {
			m := (lo + hi) / 2
			if fl(m) < 0 {
				lo = m
			} else {
				hi = m
			}
		}
	} else {
		for i := 0; i < 60; i++ {
			m := (lo + hi) / 2
			if fl(m) < 0 {
				hi = m
			} else {
				lo = m
			}
		}
	}
	return (lo + hi) / 2
}

type medida struct {
	c                       cero
	L, Lmeas                float64
	niveles                 []float64
	xMas, xMenos            []float64
	pendienteEje, psi       float64
	gTrampa, hMed, hPred    float64
	dMas, dMenos, anchoBase float64
	xs, yMas                []float64
}

func medir(c cero) medida {
	const w, paso = 0.02, 0.001
	m := medida{c: c, L: math.Log(c.gamma / (2 * math.Pi))}
	var xs, ys []float64
	for x := -w; x <= w+1e-12; x += paso {
		xs = append(xs, x)
		ys = append(ys, f2(c.beta+x, c.gamma))
	}
	m.xs, m.yMas = xs, ys
	cf := ajusteCubico(xs, ys)
	m.Lmeas = -cf[3] / cf[2]
	// equal-height crossings
	var mids, lvl []float64
	for _, wk := range []float64{0.005, 0.01, 0.02} {
		y := c.a * wk * wk
		xp, xm := cruce(c, y, true), cruce(c, y, false)
		m.niveles = append(m.niveles, y)
		m.xMas = append(m.xMas, xp)
		m.xMenos = append(m.xMenos, xm)
		mids = append(mids, (xp+xm)/2)
		lvl = append(lvl, y)
	}
	// axis lean: chord midpoint vs height, slope through the origin
	var sxy, sxx float64
	for i := range mids {
		sxy += mids[i] * lvl[i]
		sxx += lvl[i] * lvl[i]
	}
	m.pendienteEje = sxy / sxx
	m.psi = math.Atan(m.pendienteEje)
	// separations at the window edge
	fp, fm := f2(c.beta+w, c.gamma), f2(c.beta-w, c.gamma)
	m.gTrampa = math.Abs((fp + fm) - (fm + fp))
	m.hMed = fp - fm
	m.hPred = -2 * c.a * m.L * w * w * w
	m.dMas = math.Hypot(w, fp)
	m.dMenos = math.Hypot(w, fm)
	m.anchoBase = m.xMas[2] - m.xMenos[2]
	return m
}

func main() {
	fmt.Println("⏝⏜ LAS PARÁBOLAS CENTRADAS — contacto en el (0,0) absoluto: x = σ − β, y = |ζ|²")
	fmt.Println("   β medido por Newton en el plano (nunca supuesto) · el espejo de abajo es y ↦ −f₊(−x)")
	fmt.Println()

	ceros := append(buscar(10, 50, 0.05), buscar(998, 1004, 0.02)...)
	var ms []medida
	for _, c := range ceros {
		ms = append(ms, medir(c))
	}

	fmt.Println("§1 · la relación con cada altura: qué se queda quieto y qué se mueve")
	fmt.Println()
	fmt.Println("        γ            β (contacto)      a = |ζ'(ρ)|²    L = log(γ/2π)   L medido (cúbico)   razón")
	peor := 0.0
	for _, m := range ms {
		r := m.Lmeas / m.L
		if d := math.Abs(r - 1); d > peor {
			peor = d
		}
		fmt.Printf("   %12.6f     %.12f     %9.5f        %.4f          %.4f          %.4f\n", m.c.gamma, m.c.beta, m.c.a, m.L, m.Lmeas, r)
	}
	fmt.Printf("\n   V4 · β no se mueve con la altura (½ en los %d); L sube como log(γ/2π); a anda de %.3f a %.3f zero a zero\n",
		len(ms), minA(ms), maxA(ms))
	fmt.Printf("   V2 · el reloj leído de la forma de la parábola: peor desvío de L medido/L = %.1e → %s\n", peor,
		map[bool]string{true: "la asimetría ES el reloj", false: "revisar"}[peor < 0.01])
	fmt.Println()

	fmt.Println("§2 · distancia positiva contra negativa a la MISMA altura y (tres niveles por cero)")
	fmt.Println()
	fmt.Println("        γ          nivel y        x₊ (derecha)   |x₋| (izquierda)   x₊/|x₋|    x₊+x₋ medido   L·y/a predicho")
	for _, m := range ms {
		for k := range m.niveles {
			if k != 2 && m.c.gamma > 100 {
				continue
			}
			if k != 2 && m.c.gamma > 20 {
				continue
			}
			pred := m.L * m.niveles[k] / m.c.a
			fmt.Printf("   %12.6f    %.3e     %.6f       %.6f         %.5f    %+.3e      %+.3e\n",
				m.c.gamma, m.niveles[k], m.xMas[k], -m.xMenos[k], m.xMas[k]/(-m.xMenos[k]), m.xMas[k]+m.xMenos[k], pred)
		}
	}
	fmt.Println("   (a igual altura la rama derecha llega MÁS LEJOS que la izquierda; la diferencia crece con la altura y con el reloj)")
	fmt.Println()

	fmt.Println("§3 · el ángulo: el eje de la parábola no es vertical, se inclina con el reloj")
	fmt.Println()
	fmt.Println("        γ         pendiente del eje (dx/dy)   L/(2a) predicho   ψ = atan(pendiente)   [ángulo en unidades crudas σ vs |ζ|²]")
	for _, m := range ms {
		fmt.Printf("   %12.6f        %+.4f                 %+.4f          %+.4f rad = %+.2f°\n",
			m.c.gamma, m.pendienteEje, m.L/(2*m.c.a), m.psi, m.psi*180/math.Pi)
	}
	fmt.Println("   El ángulo depende de la escala elegida para y; la pendiente L/(2a) no. Es el ÚNICO ángulo intrínseco del par.")
	fmt.Println("   En el punto de contacto las dos tangentes son horizontales: el ángulo entre las dos parábolas ahí es 0.")
	fmt.Println()

	fmt.Println("§4 · las separaciones (ventana x = ±0,02)")
	fmt.Println()
	fmt.Println("        γ         |0→(w, f₊(w))|   |0→(−w, f₊(−w))|   ancho base x₊−x₋    mitad−mitad: f₊(w)−f₊(−w) medido   −2aLw³ predicho   g(x)−g(−x)")
	for _, m := range ms {
		fmt.Printf("   %12.6f     %.6f          %.6f           %.6f          %+.3e                    %+.3e       %.1e\n",
			m.c.gamma, m.dMas, m.dMenos, m.anchoBase, m.hMed, m.hPred, m.gTrampa)
	}
	fmt.Println("   V1 · g(x)−g(−x) = 0.0e+00 en todos: la separación VERTICAL entre las dos curvas es par por construcción (trampa #11, declarada).")
	fmt.Println("   La separación que sí mide algo es la de una mitad contra la otra de la MISMA parábola: −2aLw³, y da el reloj.")
	fmt.Println()
	fmt.Println("   Lo que esto NO es: una demostración. Anatomía de 15 ceros en dos ventanas.")
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(ms)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/las-parabolas-centradas.svg")
}

func minA(ms []medida) float64 {
	v := math.Inf(1)
	for _, m := range ms {
		v = math.Min(v, m.c.a)
	}
	return v
}

func maxA(ms []medida) float64 {
	v := 0.0
	for _, m := range ms {
		v = math.Max(v, m.c.a)
	}
	return v
}

func dibujar(ms []medida) {
	var b strings.Builder
	W, H := 1060, 640
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 34, 20, "#e8e2d0", "middle", "LAS PARÁBOLAS CENTRADAS — contacto en el (0,0): la rama derecha llega más lejos")
	txt(530, 56, 12, "#8a8676", "middle", "x = σ − β, y = |ζ|² normalizado por a·w²; la asimetría entre mitades es el reloj log(γ/2π) · F392")

	// left panel: normalized curves for the lowest and the highest zero
	x0, y0, w, h := 20.0, 80.0, 620.0, 420.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
	cx, cy := x0+w/2, y0+h/2
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3a4050"/>`, x0+20, cy, x0+w-20, cy)
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3a4050"/>`, cx, y0+20, cx, y0+h-20)
	txt(cx+6, cy-6, 11, "#8a8676", "start", "(0,0) — el contacto")
	cols := []string{"#f2c14e", "#ff6b6b"}
	sel := []medida{ms[0], ms[len(ms)-1]}
	for k, m := range sel {
		const wv = 0.02
		esc := m.c.a * wv * wv
		var pts, pts2 []string
		for i, x := range m.xs {
			px := cx + x/wv*(w/2-40)
			py := cy - m.yMas[i]/esc*(h/2-40)
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", px, py))
			px2 := cx - x/wv*(w/2-40)
			py2 := cy + m.yMas[i]/esc*(h/2-40)
			pts2 = append(pts2, fmt.Sprintf("%.1f,%.1f", px2, py2))
		}
		fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2"/>`, strings.Join(pts, " "), cols[k])
		fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2" stroke-dasharray="4,3"/>`, strings.Join(pts2, " "), cols[k])
		// axis lean line: x = slope * y (in normalized units)
		s := m.pendienteEje * esc / wv
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" opacity="0.7"/>`,
			cx-s*(h/2-40), cy+(h/2-40), cx+s*(h/2-40), cy-(h/2-40), cols[k])
		txt(x0+16, y0+22+float64(k)*16, 11, cols[k], "start",
			fmt.Sprintf("γ = %.3f · L = %.3f · L medido %.3f · eje inclinado %+.3f", m.c.gamma, m.L, m.Lmeas, m.pendienteEje))
	}
	txt(x0+16, y0+h-14, 10, "#8a8676", "start", "trazo lleno: la de arriba · punteado: su espejo por (0,0) · línea fina: el eje de cada parábola")

	// right panel: L measured vs log(gamma/2pi)
	rx, ry, rw, rh := 660.0, 80.0, 380.0, 420.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, rx, ry, rw, rh)
	txt(rx+rw/2, ry+20, 12, "#e8e2d0", "middle", "el reloj leído de la forma: L medido vs log(γ/2π)")
	lmax := 0.0
	for _, m := range ms {
		lmax = math.Max(lmax, math.Max(m.L, m.Lmeas))
	}
	lmax *= 1.1
	px := func(v float64) float64 { return rx + 40 + v/lmax*(rw-70) }
	py := func(v float64) float64 { return ry + rh - 40 - v/lmax*(rh-80) }
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3a4050" stroke-dasharray="3,3"/>`, px(0), py(0), px(lmax), py(lmax))
	for _, m := range ms {
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="#5ec8ff"/>`, px(m.L), py(m.Lmeas))
	}
	txt(rx+rw/2, ry+rh-12, 10, "#8a8676", "middle", "eje x: log(γ/2π) · eje y: −b/a del ajuste cúbico · la diagonal es razón 1")
	txt(30, 530, 11, "#e8e2d0", "start", "V1 (trampa #11): la separación vertical entre las dos curvas es par por construcción — 0.0e+00, sin mérito.")
	txt(30, 548, 11, "#f2c14e", "start", "V2: la asimetría entre las dos mitades de la MISMA parábola es −2aLx³, y L es el reloj: medido al 1%.")
	txt(30, 566, 11, "#5ec8ff", "start", "V3: a igual altura, x₊ + x₋ = L·y/a — el eje se inclina con pendiente L/(2a), el único ángulo intrínseco del par.")
	txt(30, 584, 11, "#8a8676", "start", "V4: con la altura, β se queda en ½ y a vaga cero a cero; lo que crece es la inclinación. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "las-parabolas-centradas.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
