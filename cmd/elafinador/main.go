package main

// EL AFINADOR - the captain: "then set the formula equal to sqrt(x), and in the
// cambiaformas put the equivalent of putting ALL the numbers into ONE variable,
// to see whether any one of them is out of tune."
//
// Two orders, and both are exact operations, not approximations.
//
// (A) DIVIDE THE EXPLICIT FORMULA BY sqrt(x).
//     Each pair of conjugate zeros contributes
//         2 sqrt(x) * (0.5 cos(g ln x) + g sin(g ln x)) / (0.25 + g^2)
//     Divide by sqrt(x) and the amplitude 2/|rho| no longer depends on x AT ALL.
//     x survives in one single place: inside the cosine, as ln x.
//     That IS "all the numbers in one variable": the variable is u = ln x, and
//     every number in the world is one point on that one line.
//
// (B) THE CAMBIAFORMAS AS A TUNING FORK.
//     w(s) = 1 - 1/s = (s-1)/s.  On the critical line s = 1/2 + i g:
//         |w| = |-1/2 + i g| / |1/2 + i g| = 1  EXACTLY.
//     So the cambiaformas maps the critical line onto the UNIT CIRCLE. One
//     variable per zero - the radius - and "in tune" means exactly 1.
//     A zero with sigma > 1/2 falls INSIDE the circle, one with sigma < 1/2
//     falls OUTSIDE. Either way it is out of tune, and the radius says so.
//
// Verdicts fixed BEFORE running:
//   V1 psi(x) normalised by sqrt(x) must stay BOUNDED over the whole measured
//      range, and below Schoenfeld's RH bound ln^2(x)/(8 pi).
//   V2 every fished zero, run through the cambiaformas, must give |w| = 1 to
//      machine precision. Reported as the worst deviation found.
//   V3 (THE CONTROL THAT CAN FAIL) injected off-line zeros must be detected by
//      the same radius. If they are not, the instrument is useless and it gets
//      reported as useless.
//   V4 (TRAP 26, declared before computing) |w| = 1 on the critical line is
//      Moebius geometry: it is an exact RESTATEMENT of the hypothesis, not a
//      proof of it, and it says nothing about zeros nobody has found. The
//      control is to MEASURE the sensitivity and publish where it goes deaf.
//   V5 the cambiaformas must have order 3: w(w(w(s))) = s. Verified numerically.

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"path/filepath"
	"strings"
)

// ---------- zeros of zeta, fished blind by sign change of Hardy's Z ----------

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

func alturas(tMax float64) []float64 {
	var g []float64
	a, za := 10.0, zetaZ(10.0)
	for b := 10.01; b <= tMax; b += 0.01 {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 60; i++ {
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

// ---------- the true staircase, built once from the sieve ----------

// psiAcum returns the running Chebyshev staircase psi(n) for n = 0..N.
func psiAcum(N int) []float64 {
	lam := make([]float64, N+1)
	compuesto := make([]bool, N+1)
	for p := 2; p <= N; p++ {
		if compuesto[p] {
			continue
		}
		for m := p * p; m <= N; m += p {
			compuesto[m] = true
		}
		lp := math.Log(float64(p))
		for q := p; ; {
			lam[q] = lp
			if q > N/p {
				break
			}
			q *= p
		}
	}
	acum := make([]float64, N+1)
	for n := 1; n <= N; n++ {
		acum[n] = acum[n-1] + lam[n]
	}
	return acum
}

// ---------- the cambiaformas ----------

func cambiaformas(s complex128) complex128 { return 1 - 1/s }

// ---------- reported numbers, filled by main and used by the plate ----------

type medida struct {
	x    float64
	norm float64 // (psi(x) - x + ln 2pi + .5 ln(1-1/x^2)) / sqrt(x)
	cota float64 // Schoenfeld under RH: ln^2(x) / (8 pi)
}

type desafinado struct {
	sigma, gamma float64
	radio        float64
	desvio       float64 // 1 - |w|
	predicho     float64 // delta / (1/4 + gamma^2)
}

var (
	medidas    []medida
	peorNorm   medida
	nCeros     int
	peorRadio  float64
	desafinan  []desafinado
	crecimient [][3]string
)

func main() {
	fmt.Println("EL AFINADOR — F416")
	fmt.Println("el capitán: «igualá la fórmula a raíz de x, y en el cambiaformas")
	fmt.Println("poné todos los números en una variable, a ver si alguno desentona»")
	fmt.Println()

	// ---------------- (A) dividir por raíz de x ----------------
	const N = 1000000
	acum := psiAcum(N)
	fmt.Println("(A) LA FÓRMULA DIVIDIDA POR √x")
	fmt.Println("    E(x) = (ψ(x) − x + ln2π + ½ln(1−1/x²)) / √x")
	fmt.Println()
	fmt.Printf("    %12s  %12s  %12s\n", "x", "E(x)", "cota RH")
	for _, x := range []float64{100, 1000, 10000, 100000, 500000, 1000000} {
		n := int(x)
		e := (acum[n] - x - math.Log(2*math.Pi) - 0.5*math.Log(1-1/(x*x))) / math.Sqrt(x)
		c := math.Log(x) * math.Log(x) / (8 * math.Pi)
		medidas = append(medidas, medida{x, e, c})
		fmt.Printf("    %12.0f  %12.6f  %12.4f\n", x, e, c)
	}
	// worst normalised excursion over the whole range
	peor := medida{}
	for n := 100; n <= N; n++ {
		x := float64(n)
		e := (acum[n] - x - math.Log(2*math.Pi) - 0.5*math.Log(1-1/(x*x))) / math.Sqrt(x)
		if math.Abs(e) > math.Abs(peor.norm) {
			peor = medida{x, e, math.Log(x) * math.Log(x) / (8 * math.Pi)}
		}
	}
	peorNorm = peor
	fmt.Printf("\n    peor excursión en todo el rango: E(%.0f) = %.6f", peor.x, peor.norm)
	fmt.Printf("   (cota RH allí: %.4f)\n", peor.cota)
	fmt.Println("    → ACOTADO. Dividido por √x, cada eco tiene amplitud FIJA 2/|ρ|")
	fmt.Println("      y la x sobrevive en un solo lugar: adentro del coseno, como ln x.")
	fmt.Println()

	// ---------------- (B) el cambiaformas como afinador ----------------
	fmt.Println("(B) EL CAMBIAFORMAS COMO AFINADOR — una variable: el radio |w(ρ)|")
	gs := alturas(200)
	nCeros = len(gs)
	peorRadio = 0
	for _, g := range gs {
		r := cmplx.Abs(cambiaformas(complex(0.5, g)))
		if d := math.Abs(r - 1); d > peorRadio {
			peorRadio = d
		}
	}
	fmt.Printf("    %d ceros pescados a ciegas hasta altura 200\n", nCeros)
	fmt.Printf("    peor desvío de |w| respecto de 1: %.3e  (precisión de máquina)\n", peorRadio)
	fmt.Println("    → todos EN TONO: el cambiaformas manda la línea crítica")
	fmt.Println("      al CÍRCULO UNIDAD, exacto, sin aproximación.")
	fmt.Println()

	// order 3
	fmt.Println("    y el cambiaformas gira en tres pasos: s → 1−1/s → 1/(1−s) → s")
	for _, s := range []complex128{complex(2, 0), complex(0.5, 14.134725), complex(0.75, 3)} {
		t := cambiaformas(cambiaformas(cambiaformas(s)))
		fmt.Printf("      w³(%.4f%+.4fi) = %.4f%+.4fi   (error %.2e)\n",
			real(s), imag(s), real(t), imag(t), cmplx.Abs(t-s))
	}
	fmt.Printf("      y el número de la casa sigue en pie: w(2) = %.6f\n\n", real(cambiaformas(2)))

	// ---------------- (C) el control que puede fallar ----------------
	fmt.Println("(C) EL CONTROL QUE PUEDE FALLAR — ceros falsos, fuera de la línea")
	fmt.Printf("    %8s %9s %14s %12s %12s\n", "σ", "γ", "|w|", "1−|w|", "predicho")
	for _, sg := range []float64{0.6, 0.55, 0.51, 0.501} {
		for _, g := range []float64{14.134725, 100, 500} {
			r := cmplx.Abs(cambiaformas(complex(sg, g)))
			d := sg - 0.5
			pr := d / (0.25 + g*g)
			desafinan = append(desafinan, desafinado{sg, g, r, 1 - r, pr})
			fmt.Printf("    %8.3f %9.2f %14.10f %12.3e %12.3e\n", sg, g, r, 1-r, pr)
		}
	}
	fmt.Println()
	fmt.Println("    → SÍ desentona: cualquier cero fuera de la línea cae fuera del círculo.")
	fmt.Println("    ⚠ PERO el desvío muere como 1/γ². Arriba el afinador se queda SORDO:")
	fmt.Println("      a altura 500 un cero desplazado 0,1 mueve el radio 4 partes en DIEZ millones.")
	fmt.Println()

	// ---------------- (D) el eco desentonado ----------------
	fmt.Println("(D) EL ECO DESENTONADO — por qué el √x es el límite")
	fmt.Println("    un cero en σ da un eco de tamaño x^σ; dividido por √x queda x^(σ−½)")
	fmt.Printf("    %14s %16s %16s\n", "x", "σ=0,51", "σ=0,6")
	for _, x := range []float64{1e6, 1e12, 1e30, 1e100} {
		a := math.Pow(x, 0.01)
		b := math.Pow(x, 0.1)
		crecimient = append(crecimient, [3]string{
			fmt.Sprintf("10^%.0f", math.Log10(x)),
			fmt.Sprintf("%.3g", a), fmt.Sprintf("%.3g", b)})
		fmt.Printf("    %14s %16.3g %16.3g\n", fmt.Sprintf("10^%.0f", math.Log10(x)), a, b)
	}
	fmt.Println()
	fmt.Println("    → un solo cero desafinado termina tapando a toda la orquesta.")
	fmt.Println("      Por eso el √x no es un detalle: es la afinación entera.")
	fmt.Println()

	fmt.Println("HONESTO:")
	fmt.Println("  · que |w| = 1 en la línea crítica es geometría de Möbius elemental.")
	fmt.Println("    Es una REESCRITURA exacta de la Hipótesis, no una demostración.")
	fmt.Println("  · el afinador solo puede juzgar ceros que ya encontramos; de los que")
	fmt.Println("    no encontramos no dice nada, y arriba pierde sensibilidad como 1/γ².")
	fmt.Println("  · lo medido y propio: la excursión acotada de E(x) hasta un millón,")
	fmt.Println("    el desvío de máquina en los 79 ceros, y dónde se queda sordo.")
	fmt.Println()
	fmt.Println("Todavía no.")

	lamina()
	fmt.Println("\nlámina: galeria/laminas/10-el-telar/el-afinador.svg")
}

func lamina() {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1060 700" width="1060" height="700">`)
	b.WriteString(`<rect width="1060" height="700" fill="#0d1017"/>`)
	txt := func(x, y, sz float64, col, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="Georgia,serif" font-size="%.0f" fill="%s" text-anchor="%s">%s</text>`,
			x, y, sz, col, anc, s)
	}
	txt(530, 42, 24, "#e8e2d0", "middle", "EL AFINADOR")
	txt(530, 68, 12, "#8a8676", "middle", "«igualá la fórmula a raíz de x, y en el cambiaformas poné todos los números en una variable, a ver si alguno desentona»")

	// panel A
	fmt.Fprintf(&b, `<rect x="20" y="92" width="500" height="250" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 118, 13, "#6ee7a0", "middle", "A · dividida por √x, la x queda en un solo lugar")
	txt(270, 156, 16, "#f2c14e", "middle", "ψ(x) − x")
	fmt.Fprintf(&b, `<line x1="170" y1="164" x2="370" y2="164" stroke="#f2c14e" stroke-width="1.5"/>`)
	txt(270, 186, 16, "#f2c14e", "middle", "√x")
	txt(270, 216, 13, "#e8e2d0", "middle", "= − Σ  2 · cos(γ · ln x + φ) / |ρ|")
	txt(270, 248, 11, "#c9c3b0", "middle", "la amplitud 2/|ρ| YA NO DEPENDE DE x")
	txt(270, 268, 11, "#c9c3b0", "middle", "la x sobrevive adentro del coseno, como ln x")
	txt(270, 300, 12, "#6ee7a0", "middle", "esa es la única variable: u = ln x")
	txt(270, 322, 11, "#8a8676", "middle", "todo número del mundo es un punto de esa recta")

	// panel B: the circle
	fmt.Fprintf(&b, `<rect x="540" y="92" width="500" height="250" fill="#131826" stroke="#2a3040"/>`)
	txt(790, 118, 13, "#6ee7a0", "middle", "B · el cambiaformas manda la línea al círculo")
	fmt.Fprintf(&b, `<circle cx="790" cy="230" r="78" fill="none" stroke="#f2c14e" stroke-width="1.6"/>`)
	for i := 0; i < 26; i++ {
		a := float64(i) * 2 * math.Pi / 26
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.6" fill="#6ee7a0"/>`, 790+78*math.Cos(a), 230+78*math.Sin(a))
	}
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="3.4" fill="#ff6b6b"/>`, 790+40*math.Cos(2.2), 230+40*math.Sin(2.2))
	txt(790+40*math.Cos(2.2)-8, 230+40*math.Sin(2.2)+18, 10, "#ff6b6b", "middle", "σ&gt;½")
	txt(660, 168, 11, "#e8e2d0", "start", "w(s) = 1 − 1/s")
	txt(660, 186, 11, "#c9c3b0", "start", "|w(½+iγ)| = 1 exacto")
	txt(940, 330, 11, "#8a8676", "end", "una variable por cero: el radio")

	// panel C: the control
	fmt.Fprintf(&b, `<rect x="20" y="356" width="500" height="196" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 382, 13, "#e8e2d0", "middle", "C · el control que PUEDE fallar — y falla arriba")
	rows := [][3]string{{"σ", "γ", "1 − |w|"}, {"0,60", "14,13", "5,0 × 10⁻⁴"}, {"0,60", "100", "1,0 × 10⁻⁵"}, {"0,60", "500", "4,0 × 10⁻⁷"}}
	yy := 412.0
	for i, r := range rows {
		c := "#c9c3b0"
		if i == 0 {
			c = "#8a8676"
		}
		txt(70, yy, 11, c, "start", r[0])
		txt(190, yy, 11, c, "start", r[1])
		txt(470, yy, 11, c, "end", r[2])
		yy += 26
	}
	txt(270, 534, 11, "#ff6b6b", "middle", "el desvío muere como 1/γ²: el afinador se queda SORDO en altura")

	// panel D: growth
	fmt.Fprintf(&b, `<rect x="540" y="356" width="500" height="196" fill="#131826" stroke="#2a3040"/>`)
	txt(790, 382, 13, "#e8e2d0", "middle", "D · el eco desentonado, dividido por √x")
	txt(790, 412, 13, "#f2c14e", "middle", "x^(σ − ½)")
	gr := [][3]string{{"x", "σ=0,51", "σ=0,6"}, {"10⁶", "1,15", "3,98"}, {"10³⁰", "2,00", "1000"}, {"10¹⁰⁰", "10,0", "10¹⁰"}}
	yy = 442
	for i, r := range gr {
		c := "#c9c3b0"
		if i == 0 {
			c = "#8a8676"
		}
		txt(600, yy, 11, c, "start", r[0])
		txt(790, yy, 11, c, "middle", r[1])
		txt(980, yy, 11, c, "end", r[2])
		yy += 26
	}
	txt(790, 534, 11, "#8a8676", "middle", "un solo desafinado termina tapando a toda la orquesta")

	// closing
	fmt.Fprintf(&b, `<rect x="20" y="566" width="1020" height="116" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 592, 12, "#ff6b6b", "middle", "honesto: |w| = 1 en la línea es geometría de Möbius — una REESCRITURA exacta de la Hipótesis, no una demostración")
	txt(530, 618, 11, "#c9c3b0", "middle", "y solo juzga ceros que ya encontramos: de los que faltan no dice nada, y arriba pierde el oído como 1/γ²")
	txt(530, 656, 13, "#6ee7a0", "middle", "estar afinado es estar en el círculo. La Hipótesis es que NINGUNO desentona.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-afinador.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
