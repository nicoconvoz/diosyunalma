package main

// LAS DOS PARABOLAS - the captain's flash: put one parabola above and one
// below, touching at a single point that reaches 0. Does the contact point
// MEASURE the 1/2, or does it carry it built in?
//
// Construction, with nothing assumed about where zeros live:
//   f+(sigma) =  |zeta(sigma + i t)|^2        (the parabola above)
//   f-(sigma) = -|zeta((1-sigma) + i t)|^2    (the parabola below, the mirror)
// Both are computed DIRECTLY by Euler-Maclaurin at their own argument; the
// functional equation is never invoked numerically.
//
// Heights are found BLIND: a 2-D scan of |zeta|^2 over the whole rectangle
// 0.2 <= sigma <= 0.8, refined by 2-D Newton on zeta(s) = 0. The real part
// of each zero is an OUTPUT, never an input.
//
// Verdicts fixed BEFORE running:
//   V1  (declared trap) the midpoint of the two vertices is 1/2 at EVERY
//       height, zero or not: f-(sigma) = -f+(1-sigma) is a point reflection
//       through (1/2, 0) by construction. If it prints 0.0e+00 it is the
//       tenth appearance of the trap, not a finding. The controls must show it.
//   V2  (measurement) at a zero, EACH vertex on its own lands on the zero's
//       real part; whether that is 1/2 is what the data says, not the code.
//   V3  the two parabolas touch (same vertex AND both reach 0) only at zeros;
//       at control heights the minima stay positive and there is no contact.
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

// zetaEM evaluates zeta(s) by Euler-Maclaurin with three Bernoulli terms.
func zetaEM(s complex128) complex128 {
	N := int(3*math.Abs(imag(s))) + 60
	sum := complex(0, 0)
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * complex(math.Log(float64(n)), 0))
	}
	fN := complex(float64(N), 0)
	lnN := complex(math.Log(float64(N)), 0)
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	// B2/2! s N^{-s-1} + B4/4! s(s+1)(s+2) N^{-s-3} + B6/6! s...(s+4) N^{-s-5}
	bern := []float64{1.0 / 12, -1.0 / 720, 1.0 / 30240}
	prod := s
	for k, b := range bern {
		if k > 0 {
			prod *= (s + complex(float64(2*k-1), 0)) * (s + complex(float64(2*k), 0))
		}
		sum += complex(b, 0) * prod * cmplx.Exp(-(s+complex(float64(2*k+1), 0))*lnN)
	}
	_ = fN
	return sum
}

func f2(sigma, t float64) float64 {
	return math.Pow(cmplx.Abs(zetaEM(complex(sigma, t))), 2)
}

type cero struct{ beta, gamma float64 }

// buscar scans the rectangle blind and refines each candidate by 2-D Newton.
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
			for it := 0; it < 30; it++ {
				h := complex(1e-6, 0)
				d := (zetaEM(s+h) - zetaEM(s-h)) / (2 * h)
				paso := zetaEM(s) / d
				s -= paso
				if cmplx.Abs(paso) < 1e-13 {
					break
				}
			}
			if cmplx.Abs(zetaEM(s)) < 1e-9 {
				out = append(out, cero{real(s), imag(s)})
			}
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].gamma < out[b].gamma })
	// merge duplicates converging to the same zero
	var u []cero
	for _, c := range out {
		if len(u) == 0 || math.Abs(u[len(u)-1].gamma-c.gamma) > 1e-6 {
			u = append(u, c)
		}
	}
	return u
}

// ajusteParabola fits y = a x^2 + b x + c by least squares; returns vertex and its value.
func ajusteParabola(x, y []float64) (vx, vy, a float64) {
	var s0, s1, s2, s3, s4, t0, t1, t2 float64
	for i := range x {
		xi, yi := x[i], y[i]
		s0++
		s1 += xi
		s2 += xi * xi
		s3 += xi * xi * xi
		s4 += xi * xi * xi * xi
		t0 += yi
		t1 += xi * yi
		t2 += xi * xi * yi
	}
	// normal equations (3x3) solved by Cramer
	m := [3][3]float64{{s4, s3, s2}, {s3, s2, s1}, {s2, s1, s0}}
	r := [3]float64{t2, t1, t0}
	det := func(m [3][3]float64) float64 {
		return m[0][0]*(m[1][1]*m[2][2]-m[1][2]*m[2][1]) -
			m[0][1]*(m[1][0]*m[2][2]-m[1][2]*m[2][0]) +
			m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0])
	}
	D := det(m)
	col := func(k int) float64 {
		mm := m
		for i := 0; i < 3; i++ {
			mm[i][k] = r[i]
		}
		return det(mm) / D
	}
	a, b, c := col(0), col(1), col(2)
	vx = -b / (2 * a)
	vy = c - b*b/(4*a)
	return
}

type fila struct {
	nombre string
	t      float64
	esCero bool
	beta   float64 // Newton real part (only for zeros)
	// fits on the SAME sigma window for both parabolas
	vMas, vMenos, mMas, mMenos float64
	// fit of the lower parabola on the MIRRORED window (1-sigma)
	vMenosEspejo float64
	brecha       float64
	identidad    float64
	// lean of the upper parabola: measured vertex - beta, and the value
	// predicted from the clock rate L = log(gamma/2pi) for this window
	inclinacion, inclinacionPred, inclinacionFina float64
	xs, yMas, yMenos                              []float64
}

// momento returns sum(x^4)/sum(x^2) for the symmetric grid used by the fit.
func momento(ancho, paso float64) float64 {
	var s2, s4 float64
	for x := -ancho; x <= ancho+1e-12; x += paso {
		s2 += x * x
		s4 += x * x * x * x
	}
	return s4 / s2
}

func ajustarVentana(t, centro, ancho, paso float64, abajo bool) (xs, ys []float64) {
	for x := centro - ancho; x <= centro+ancho+1e-12; x += paso {
		xs = append(xs, x)
		if abajo {
			ys = append(ys, -f2(1-x, t))
		} else {
			ys = append(ys, f2(x, t))
		}
	}
	return
}

func medir(nombre string, t, centro float64, esCero bool, beta float64) fila {
	const ancho, paso = 0.02, 0.002
	xs, yp := ajustarVentana(t, centro, ancho, paso, false)
	_, ym := ajustarVentana(t, centro, ancho, paso, true)
	identidad := 0.0
	for i, x := range xs {
		if d := math.Abs(ym[i] + f2(1-x, t)); d > identidad {
			identidad = d
		}
	}
	vp, mp, _ := ajusteParabola(xs, yp)
	vm, mm, _ := ajusteParabola(xs, ym)
	xe, ye := ajustarVentana(t, 1-centro, ancho, paso, true)
	ve, _, _ := ajusteParabola(xe, ye)
	sc := (vp + ve) / 2
	brecha := f2(sc, t) + f2(1-sc, t)
	f := fila{nombre: nombre, t: t, esCero: esCero, beta: beta,
		vMas: vp, vMenos: vm, mMas: mp, mMenos: mm, vMenosEspejo: ve,
		brecha: brecha, identidad: identidad, xs: xs, yMas: yp, yMenos: ym}
	if esCero {
		L := math.Log(t / (2 * math.Pi))
		f.inclinacion = vp - beta
		f.inclinacionPred = L * momento(ancho, paso) / 2
		xf, yf := ajustarVentana(t, centro, ancho/4, paso/4, false)
		vf, _, _ := ajusteParabola(xf, yf)
		f.inclinacionFina = vf - beta
	}
	return f
}

func main() {
	fmt.Println("⏝⏜ LAS DOS PARÁBOLAS — la de arriba |ζ(σ+it)|², la de abajo −|ζ(1−σ+it)|²")
	fmt.Println("   ceros buscados a ciegas en el rectángulo 0,2 ≤ σ ≤ 0,8; el ½ nunca entra")
	fmt.Println()

	fmt.Println("§0 · calibración del evaluador")
	z2 := zetaEM(complex(2, 0))
	fmt.Printf("   ζ(2) = %.12f   (π²/6 = %.12f, desvío %.1e)\n", real(z2), math.Pi*math.Pi/6, math.Abs(real(z2)-math.Pi*math.Pi/6))
	zm := zetaEM(complex(-1, 0))
	fmt.Printf("   ζ(−1) = %.12f  (−1/12 = %.12f, desvío %.1e)\n", real(zm), -1.0/12, math.Abs(real(zm)+1.0/12))
	fmt.Println()

	fmt.Println("§1 · los ceros, hallados a ciegas (mínimos de |ζ|² en 2-D + Newton en el plano)")
	ceros := buscar(10, 50, 0.05)
	hondos := buscar(998, 1004, 0.02)
	todos := append(append([]cero{}, ceros...), hondos...)
	fmt.Println()
	fmt.Println("        γ (altura)        β = Re ρ (Newton, salida)     |ζ(ρ)|")
	for _, c := range todos {
		fmt.Printf("   %14.9f     %.12f          %.1e\n", c.gamma, c.beta, cmplx.Abs(zetaEM(complex(c.beta, c.gamma))))
	}
	fmt.Println()

	var filas []fila
	for i, c := range todos {
		filas = append(filas, medir(fmt.Sprintf("cero γ=%.3f", c.gamma), c.gamma, c.beta, true, c.beta))
		if i+1 < len(ceros) {
			tc := (c.gamma + ceros[i+1].gamma) / 2
			// control: locate the sigma of minimum |zeta|^2 at that height, blind
			best, bs := math.Inf(1), 0.5
			for x := 0.2; x <= 0.8; x += 0.005 {
				if v := f2(x, tc); v < best {
					best, bs = v, x
				}
			}
			filas = append(filas, medir(fmt.Sprintf("control t=%.3f", tc), tc, bs, false, math.NaN()))
		}
	}

	fmt.Println("§2 · las dos parábolas, ajustadas por mínimos cuadrados (vértice = salida)")
	fmt.Println()
	fmt.Println("   caso                  vértice ARRIBA   vértice ABAJO    mín arriba   mín abajo   | ABAJO ventana espejo   punto medio   brecha de contacto")
	for _, f := range filas {
		fmt.Printf("   %-20s   %.6f         %.6f         %+.2e    %+.2e   |     %.6f           %.6f      %.2e\n",
			f.nombre, f.vMas, f.vMenos, f.mMas, f.mMenos, f.vMenosEspejo, (f.vMas+f.vMenosEspejo)/2, f.brecha)
	}
	fmt.Println("   (misma ventana para las dos: en los controles el vértice de abajo cae fuera de su ventana y el ajuste es basura;")
	fmt.Println("    ventana espejo: cada parábola ajustada alrededor de su propio vértice)")
	fmt.Println()

	fmt.Println("§3 · veredictos pre-registrados")
	peorId := 0.0
	peorMedio, peorMedioMisma := 0.0, 0.0
	for _, f := range filas {
		if f.identidad > peorId {
			peorId = f.identidad
		}
		if d := math.Abs((f.vMas+f.vMenosEspejo)/2 - 0.5); d > peorMedio {
			peorMedio = d
		}
		if d := math.Abs((f.vMas+f.vMenos)/2 - 0.5); d > peorMedioMisma && !f.esCero {
			peorMedioMisma = d
		}
	}
	fmt.Printf("   V1 · punto medio de los vértices con ventanas espejo, peor desvío de ½ en TODOS los casos: %.1e\n", peorMedio)
	fmt.Printf("        con la MISMA ventana, en los controles se va hasta %.2f: el ½ vive en el espejo, no en zeta\n", peorMedioMisma)
	fmt.Printf("        f−(σ) + f+(1−σ), peor valor: %.1e  → %s\n", peorId,
		map[bool]string{true: "TRAMPA DEL 0.0e+00 CONFIRMADA: el ½ del punto medio es de la construcción", false: "la identidad no cerró: revisar"}[peorId < 1e-12])
	fmt.Println("        (corrección al pre-registro: la frase «a toda altura» olvidó decir «con ventanas espejadas»; sin espejar, no)")
	peorV, peorC := 0.0, math.Inf(1)
	for _, f := range filas {
		if f.esCero {
			if d := math.Abs(f.vMas - 0.5); d > peorV {
				peorV = d
			}
		} else if d := math.Abs(f.vMas - 0.5); d < peorC {
			peorC = d
		}
	}
	fmt.Printf("   V2 · en los ceros, cada vértice POR SEPARADO: peor desvío de ½ = %.1e (medido, sin espejo)\n", peorV)
	fmt.Printf("        en los controles, el vértice de arriba más cercano a ½ dista %.3f: el ½ es de los ceros, no de la altura\n", peorC)
	peorBrechaCero, menorBrechaCtl := 0.0, math.Inf(1)
	for _, f := range filas {
		if f.esCero {
			if f.brecha > peorBrechaCero {
				peorBrechaCero = f.brecha
			}
		} else if f.brecha < menorBrechaCtl {
			menorBrechaCtl = f.brecha
		}
	}
	fmt.Printf("   V3 · contacto: brecha máxima en ceros %.1e · brecha mínima en controles %.2e → %s\n",
		peorBrechaCero, menorBrechaCtl, map[bool]string{true: "se tocan SOLO en los ceros", false: "revisar"}[peorBrechaCero < 1e-8 && menorBrechaCtl > 1e-3])
	fmt.Println()

	fmt.Println("§4 · por qué el vértice de arriba no cae clavado: la parábola está inclinada por el reloj")
	fmt.Println("   |ζ(σ+iγ)|² = |ζ'|²·[(σ−β)² − L·(σ−β)³ + …] con L = −Re(ζ''/ζ')(ρ) = 2θ'(γ) ≈ log(γ/2π)")
	fmt.Println("   (Cauchy–Riemann + ζ = e^(−iθ)Z sobre la línea); el ajuste cuadrático de un cúbico corre el")
	fmt.Println("   vértice en L·(Σx⁴/Σx²)/2 para la grilla usada. Predicho antes de mirar la columna medida:")
	fmt.Println()
	fmt.Println("        γ            L=log(γ/2π)    inclinación PREDICHA   MEDIDA (vértice−β)   razón   ventana/4 → medida")
	peorRazon := 0.0
	for _, f := range filas {
		if !f.esCero {
			continue
		}
		L := math.Log(f.t / (2 * math.Pi))
		r := f.inclinacion / f.inclinacionPred
		if d := math.Abs(r - 1); d > peorRazon {
			peorRazon = d
		}
		fmt.Printf("   %12.6f     %.4f          %+.3e            %+.3e        %.4f    %+.2e\n",
			f.t, L, f.inclinacionPred, f.inclinacion, r, f.inclinacionFina)
	}
	fmt.Printf("\n   peor desvío de la razón medida/predicha respecto de 1: %.1e", peorRazon)
	fmt.Println(map[bool]string{true: "  → la inclinación ES el reloj", false: "  → la derivación no cierra: revisar"}[peorRazon < 0.05])
	fmt.Println("   Con la ventana a un cuarto, la inclinación cae ~16× (va como el cuadrado del ancho): el vértice")
	fmt.Println("   converge a β. Y en el par espejado las dos inclinaciones son opuestas: se cancelan en el punto medio.")
	fmt.Println()
	fmt.Println("   Lo que esto NO es: una demostración. Mide β para los ceros hallados en dos ventanas")
	fmt.Println("   y nada más. Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(filas)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/las-dos-parabolas.svg")
}

func dibujar(filas []fila) {
	var b strings.Builder
	W, H := 1060, 640
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 34, 20, "#e8e2d0", "middle", "LAS DOS PARÁBOLAS — arriba |ζ(σ+it)|², abajo −|ζ(1−σ+it)|²")
	txt(530, 56, 12, "#8a8676", "middle", "ceros hallados a ciegas en 0,2 ≤ σ ≤ 0,8 · el ½ nunca entra al cálculo · F391")

	panel := func(f fila, x0, y0, w, h float64, titulo string) {
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
		txt(x0+w/2, y0+20, 13, "#e8e2d0", "middle", titulo)
		ymax := 0.0
		for i := range f.xs {
			if v := math.Abs(f.yMas[i]); v > ymax {
				ymax = v
			}
			if v := math.Abs(f.yMenos[i]); v > ymax {
				ymax = v
			}
		}
		if ymax == 0 {
			ymax = 1
		}
		xmin, xmax := f.xs[0], f.xs[len(f.xs)-1]
		px := func(x float64) float64 { return x0 + 30 + (x-xmin)/(xmax-xmin)*(w-60) }
		py := func(y float64) float64 { return y0 + h/2 - y/ymax*(h/2-40) }
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3a4050"/>`, x0+30, py(0), x0+w-30, py(0))
		for i := range f.xs {
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.5" fill="#f2c14e"/>`, px(f.xs[i]), py(f.yMas[i]))
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.5" fill="#5ec8ff"/>`, px(f.xs[i]), py(f.yMenos[i]))
		}
		sc := (f.vMas + f.vMenosEspejo) / 2
		if sc >= xmin && sc <= xmax {
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ff6b6b" stroke-dasharray="3,3"/>`, px(sc), y0+30, px(sc), y0+h-30)
		}
		txt(x0+w/2, y0+h-12, 11, "#ff6b6b", "middle", fmt.Sprintf("punto medio con ventanas espejo: %.6f (tautología)", sc))
		txt(px(f.vMas), py(f.mMas)-8, 11, "#f2c14e", "middle", fmt.Sprintf("vértice %.5f", f.vMas))
		txt(x0+w-8, y0+36, 10, "#8a8676", "end", fmt.Sprintf("brecha %.1e", f.brecha))
		txt(x0+8, y0+36, 10, "#8a8676", "start", fmt.Sprintf("σ de %.2f a %.2f", xmin, xmax))
	}
	var z, c fila
	for _, f := range filas {
		if f.esCero && z.xs == nil {
			z = f
		}
		if !f.esCero && c.xs == nil {
			c = f
		}
	}
	panel(z, 20, 80, 500, 330, "EN UN CERO: se tocan en (β, 0) — "+z.nombre)
	panel(c, 540, 80, 500, 330, "CONTROL, sin cero: no se tocan — "+c.nombre)

	txt(30, 440, 12, "#e8e2d0", "start", "vértice de ARRIBA, medido por separado (sin espejo):")
	y := 460.0
	for _, f := range filas {
		if !f.esCero {
			continue
		}
		txt(30, y, 11, "#c9c3b0", "start", fmt.Sprintf("%-18s  β(Newton) = %.9f   vértice = %.6f", f.nombre, f.beta, f.vMas))
		y += 14
		if y > 600 {
			break
		}
	}
	txt(560, 440, 12, "#e8e2d0", "start", "lo que es de la construcción y lo que es de zeta:")
	txt(560, 462, 11, "#ff6b6b", "start", "punto medio de los vértices = ½ SIEMPRE, con o sin cero (trampa declarada)")
	txt(560, 480, 11, "#f2c14e", "start", "vértice individual = ½ SOLO en los ceros (medido, no impuesto)")
	txt(560, 498, 11, "#5ec8ff", "start", "contacto (brecha 0) SOLO en los ceros; en los controles la brecha es finita")
	txt(560, 530, 11, "#8a8676", "start", "Mide β para los ceros de dos ventanas. No demuestra nada")
	txt(560, 546, 11, "#8a8676", "start", "sobre los que no miró. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "las-dos-parabolas.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
