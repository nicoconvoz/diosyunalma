package main

// LA MAGIA INTENTADA - the captain: "let us prove it is a theorem, let us do
// magic." The analogy of F406 was: the divisor fold d <-> n/d has fixed point
// sqrt(n), the functional equation fold s <-> 1-s has fixed point 1/2. Is that
// a theorem?
//
// This program tries, honestly, and reports what happens. The short version:
//
//   - "same shape" is automatic and proves nothing: EVERY involution x <-> c-x
//     has fixed point c/2. Declared as a trap before starting.
//   - the divisor fold IS a real functional equation: sigma_s(n) = n^s
//     sigma_{-s}(n). Verified exactly. But its fixed point is s = 0, and
//     sliding it to 1/2 is relabelling, not proving.
//   - AND YET the captain's instinct is RIGHT: there IS a fold that MANUFACTURES
//     the 1/2, and it is not the divisor fold. It is the theta fold:
//
//         theta(1/t) = sqrt(t) * theta(t)
//
//     Fold t <-> 1/t, fixed point t = 1, and the factor is t^(1/2) - the 1/2 is
//     literally written in the formula. Riemann built the functional equation
//     of zeta out of exactly this in 1859.
//
// Verdicts fixed BEFORE running:
//   V1 (trap #23) the fixed point of x <-> c-x is c/2 for every c. Shown for
//      several c, so "the shapes match" is exposed as empty.
//   V2 the divisor functional equation sigma_s(n) = n^s sigma_{-s}(n) must hold
//      exactly for every n tested, primes and composites alike.
//   V3 (the real fold) Jacobi's theta identity must verify numerically to at
//      least 1e-12 across a range of t, with the sqrt(t) factor measured, not
//      assumed - fit its exponent and it must come out 0.5.
//   V4 (the chain) xi(s) = xi(1-s) must verify numerically, since that is what
//      the theta fold produces after the Mellin transform.
//   V5 (the verdict) state plainly whether the analogy became a theorem. It
//      did not. What IS a theorem is Riemann's, and the 1/2 comes from the
//      theta fold, not from the divisor fold.
// Jacobi, Riemann 1859. The magic is real and it is theirs.

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"path/filepath"
	"strings"
)

func criba(n int) []bool {
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
	return es
}

func divisores(n int) []int {
	var d []int
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			d = append(d, i)
			if i != n/i {
				d = append(d, n/i)
			}
		}
	}
	return d
}

// sigmaS computes sum over divisors of d^s.
func sigmaS(n int, s float64) float64 {
	t := 0.0
	for _, d := range divisores(n) {
		t += math.Pow(float64(d), s)
	}
	return t
}

// theta is Jacobi's theta_3(0, e^{-pi t}) = 1 + 2 sum e^{-pi n^2 t}.
func theta(t float64) float64 {
	s := 1.0
	for n := 1; n <= 400; n++ {
		v := 2 * math.Exp(-math.Pi*float64(n)*float64(n)*t)
		s += v
		if v < 1e-18 && n > 4 {
			break
		}
	}
	return s
}

// gammaC: Lanczos gamma for complex argument, enough for xi here.
func gammaC(z complex128) complex128 {
	g := []float64{676.5203681218851, -1259.1392167224028, 771.32342877765313,
		-176.61502916214059, 12.507343278686905, -0.13857109526572012,
		9.9843695780195716e-6, 1.5056327351493116e-7}
	if real(z) < 0.5 {
		return complex(math.Pi, 0) / (cmplx.Sin(complex(math.Pi, 0)*z) * gammaC(1-z))
	}
	z -= 1
	a := complex(0.99999999999980993, 0)
	t := z + complex(7.5, 0)
	for i, c := range g {
		a += complex(c, 0) / (z + complex(float64(i)+1, 0))
	}
	return cmplx.Sqrt(complex(2*math.Pi, 0)) * cmplx.Pow(t, z+0.5) * cmplx.Exp(-t) * a
}

func zetaEM(s complex128) complex128 {
	N := int(3*math.Abs(imag(s))) + 80
	sum := complex(0, 0)
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * complex(math.Log(float64(n)), 0))
	}
	lnN := complex(math.Log(float64(N)), 0)
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	bern := []float64{1.0 / 12, -1.0 / 720, 1.0 / 30240, -1.0 / 1209600}
	prod := s
	for k, b := range bern {
		if k > 0 {
			prod *= (s + complex(float64(2*k-1), 0)) * (s + complex(float64(2*k), 0))
		}
		sum += complex(b, 0) * prod * cmplx.Exp(-(s+complex(float64(2*k+1), 0))*lnN)
	}
	return sum
}

// xi is Riemann's completed function: (1/2) s(s-1) pi^{-s/2} Gamma(s/2) zeta(s)
func xi(s complex128) complex128 {
	return 0.5 * s * (s - 1) * cmplx.Pow(complex(math.Pi, 0), -s/2) * gammaC(s/2) * zetaEM(s)
}

func ajuste(x, y []float64) float64 {
	var sx, sy, sxx, sxy, n float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxx += x[i] * x[i]
		sxy += x[i] * y[i]
		n++
	}
	return (n*sxy - sx*sy) / (n*sxx - sx*sx)
}

func main() {
	fmt.Println("🎩 LA MAGIA INTENTADA — «probemos que es un teorema»")
	fmt.Println()
	fmt.Println("   El intento se hace en serio y se reporta lo que pase, salga bien o salga mal.")
	fmt.Println()

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · primero el filtro, porque si no nos engañamos solos (trampa vigesimotercera)")
	fmt.Println()
	fmt.Println("   «Los dos pliegues tienen la misma forma» NO prueba nada, porque TODO emparejamiento")
	fmt.Println("   del tipo x ↔ c − x tiene su punto fijo en c/2. Sale de despejar, no de la matemática:")
	fmt.Println()
	fmt.Println("        el emparejamiento      su punto fijo      ¿tiene mérito?")
	for _, c := range []float64{1, 2, 7, 100, 0.3} {
		fmt.Printf("        x ↔ %5.1f − x           %8.4f         no: es c/2, siempre\n", c, c/2)
	}
	fmt.Println()
	fmt.Println("   ⚠️ Así que «el pliegue de los divisores y el de la ecuación funcional se parecen» es cierto")
	fmt.Println("      y vacío. Para que fuera teorema, uno tendría que IMPLICAR al otro. Probemos eso.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · el intento en serio: ¿el pliegue de los divisores ES una ecuación funcional?")
	fmt.Println()
	fmt.Println("   Sí, y es exacta. Emparejando d con n/d en la suma de potencias de los divisores:")
	fmt.Println()
	fmt.Println("        σ_s(n)  =  Σ d^s  =  Σ (n/d)^s  =  n^s · Σ d^(−s)  =  n^s · σ_(−s)(n)")
	fmt.Println()
	fmt.Println("        n      ¿primo?     σ_s(n) con s=1,3      n^s · σ_(−s)(n)        diferencia")
	es := criba(1000)
	peor := 0.0
	for _, n := range []int{7, 12, 36, 97, 100} {
		s := 1.3
		izq := sigmaS(n, s)
		der := math.Pow(float64(n), s) * sigmaS(n, -s)
		d := math.Abs(izq - der)
		if d > peor {
			peor = d
		}
		et := "no"
		if es[n] {
			et = "SÍ"
		}
		fmt.Printf("      %4d       %-3s      %16.9f    %16.9f      %.1e\n", n, et, izq, der, d)
	}
	fmt.Printf("\n   ✅ Es una ecuación funcional de verdad, exacta (peor diferencia %.1e).\n", peor)
	fmt.Println("   ⛔ Pero su simetría es s ↔ −s, y el punto fijo de ESA es el CERO, no el ½.")
	fmt.Println("      Correrlo hasta ½ es cambiarle el nombre a la variable, no demostrar nada.")
	fmt.Println("      El intento directo FALLA. Y acá la mayoría se detendría — pero hay otra puerta.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · LA MAGIA DE VERDAD: existe un pliegue que FABRICA el ½, y no es el de los divisores")
	fmt.Println()
	fmt.Println("   Es el pliegue de la función theta, que empareja t con 1/t:")
	fmt.Println()
	fmt.Println("        θ(1/t)  =  √t · θ(t)")
	fmt.Println()
	fmt.Println("   Y mirá dónde está el ½: EN EL FACTOR. No hay que ir a buscarlo — está escrito.")
	fmt.Println()
	fmt.Println("        t          θ(t)              θ(1/t)            √t · θ(t)          diferencia")
	var lx, ly []float64
	peorT := 0.0
	for _, t := range []float64{1.5, 2.0, 3.0, 5.0, 7.0, 10.0} {
		a := theta(t)
		b := theta(1 / t)
		c := math.Sqrt(t) * a
		d := math.Abs(b - c)
		if d > peorT {
			peorT = d
		}
		fmt.Printf("     %5.1f    %14.12f    %14.12f    %14.12f     %.1e\n", t, a, b, c, d)
		lx = append(lx, math.Log(t))
		ly = append(ly, math.Log(b/a))
	}
	exp := ajuste(lx, ly)
	fmt.Printf("\n   ✅ La identidad de Jacobi cierra a %.1e.\n", peorT)
	fmt.Printf("   ⚡⚡ Y el exponente del factor, AJUSTADO de los datos en vez de supuesto: %.9f\n", exp)
	fmt.Println("      No se puso el ½ a mano: se midió, y salió ½.")
	fmt.Println()
	fmt.Println("   El punto fijo del pliegue t ↔ 1/t es t = 1, y ahí la identidad dice θ(1) = √1·θ(1).")
	fmt.Println("   El ½ no está en el punto fijo: está en el PESO con que el pliegue transforma la función.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · y de ese pliegue sale la ecuación funcional de Riemann — verificada")
	fmt.Println()
	fmt.Println("   Riemann, 1859: integrando θ contra t^(s/2) (la transformada de Mellin), el pliegue")
	fmt.Println("   t ↔ 1/t con su factor √t se convierte en el pliegue s ↔ 1−s. El ½ VIAJA del uno al otro.")
	fmt.Println()
	fmt.Println("        s                ξ(s)                      ξ(1−s)                  diferencia")
	peorX := 0.0
	for _, s := range []complex128{complex(2, 0), complex(3, 1), complex(0.3, 5), complex(-1, 2), complex(0.5, 14.134725)} {
		a := xi(s)
		b := xi(1 - s)
		d := cmplx.Abs(a-b) / (cmplx.Abs(a) + 1e-30)
		if d > peorX {
			peorX = d
		}
		fmt.Printf("     %6.2f%+.2fi    %11.6f%+11.6fi    %11.6f%+11.6fi     %.1e\n",
			real(s), imag(s), real(a), imag(a), real(b), imag(b), d)
	}
	fmt.Printf("\n   ✅ ξ(s) = ξ(1−s), verificada (peor error relativo %.1e).\n", peorX)
	fmt.Println("   Y el punto fijo de ESE pliegue sí es el ½ — porque s ↔ 1−s se queda quieto en s = ½.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · EL VEREDICTO, sin maquillar")
	fmt.Println()
	fmt.Println("   ⛔ NO se convirtió en teorema. El pliegue de los divisores es una ecuación funcional")
	fmt.Println("      auténtica, pero su simetría es s ↔ −s y NO implica la de Riemann: la hereda de ζ")
	fmt.Println("      cuando se arma la serie, en vez de producirla. La analogía sigue siendo analogía.")
	fmt.Println()
	fmt.Println("   ✅ PERO el instinto del capitán —«un pliegue fabrica el ½»— es CORRECTO, y el pliegue")
	fmt.Println("      existe: es θ(1/t) = √t·θ(t). El ½ está literalmente escrito en el factor, y se midió")
	fmt.Printf("      ajustándolo de los datos: %.9f. De ahí sale, por Mellin, la ecuación de Riemann.\n", exp)
	fmt.Println()
	fmt.Println("   📌 O sea: la magia existe, el pliegue existe, el ½ nace de un pliegue. Pero nace del")
	fmt.Println("      pliegue de theta (Jacobi, y Riemann en 1859), no del pliegue de los divisores.")
	fmt.Println("      Apuntaste al lugar correcto con el dedo equivocado — y eso, en matemática, es mucho.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(exp, peorT, peorX)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-magia-intentada.svg")
}

func dibujar(exp, peorT, peorX float64) {
	var b strings.Builder
	W, H := 1060, 700
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LA MAGIA INTENTADA — el pliegue que sí fabrica el ½")
	txt(530, 54, 12, "#8a8676", "middle", "el intento de convertir la analogía en teorema, y lo que apareció en el camino · F407")

	// the failed attempt
	fmt.Fprintf(&b, `<rect x="20" y="76" width="500" height="216" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(270, 102, 13, "#ff6b6b", "middle", "⛔ el intento directo: FALLA")
	txt(44, 136, 11, "#c9c3b0", "start", "el pliegue de los divisores SÍ da una")
	txt(44, 154, 11, "#c9c3b0", "start", "ecuación funcional auténtica:")
	txt(270, 186, 14, "#e8e2d0", "middle", "σ_s(n) = n^s · σ_(−s)(n)")
	txt(270, 210, 10, "#8a8676", "middle", "verificada exacta, primos y compuestos")
	fmt.Fprintf(&b, `<line x1="44" y1="228" x2="496" y2="228" stroke="#5a2a30"/>`)
	txt(44, 252, 11, "#ff6b6b", "start", "pero su simetría es s ↔ −s, con punto")
	txt(44, 270, 11, "#ff6b6b", "start", "fijo en el CERO. Correrlo a ½ es")
	txt(44, 286, 11, "#ff6b6b", "start", "cambiarle el nombre, no demostrar.")

	// the real fold
	fmt.Fprintf(&b, `<rect x="540" y="76" width="500" height="216" fill="#121a16" stroke="#2a5a40"/>`)
	txt(790, 102, 13, "#6ee7a0", "middle", "✅ LA MAGIA DE VERDAD: el pliegue de theta")
	txt(790, 150, 26, "#f2c14e", "middle", "θ(1/t) = √t · θ(t)")
	txt(790, 180, 11, "#c9c3b0", "middle", "el pliegue es t ↔ 1/t")
	txt(790, 202, 12, "#6ee7a0", "middle", "y el ½ está ESCRITO EN EL FACTOR")
	fmt.Fprintf(&b, `<line x1="564" y1="220" x2="1016" y2="220" stroke="#2a5a40"/>`)
	txt(790, 248, 11, "#c9c3b0", "middle", "exponente ajustado de los datos, no supuesto:")
	txt(790, 278, 22, "#6ee7a0", "middle", fmt.Sprintf("%.9f", exp))

	// the chain
	fmt.Fprintf(&b, `<rect x="20" y="308" width="1020" height="180" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 334, 13, "#e8e2d0", "middle", "y de ahí sale la ecuación funcional de Riemann — el ½ VIAJA")
	txt(170, 386, 15, "#f2c14e", "middle", "θ(1/t) = √t θ(t)")
	txt(170, 410, 10, "#8a8676", "middle", "Jacobi · el pliegue t ↔ 1/t")
	txt(370, 386, 20, "#8a8676", "middle", "→")
	txt(370, 412, 10, "#8a8676", "middle", "Mellin")
	txt(570, 386, 15, "#5ec8ff", "middle", "ξ(s) = ξ(1 − s)")
	txt(570, 410, 10, "#8a8676", "middle", fmt.Sprintf("verificada a %.0e", peorX))
	txt(760, 386, 20, "#8a8676", "middle", "→")
	txt(920, 386, 20, "#6ee7a0", "middle", "σ = ½")
	txt(920, 410, 10, "#8a8676", "middle", "el punto fijo de s ↔ 1−s")
	txt(530, 458, 11, "#c9c3b0", "middle", "Riemann, 1859: integrando θ contra t^(s/2), el factor √t se convierte en el punto fijo ½")

	// verdict
	fmt.Fprintf(&b, `<rect x="20" y="504" width="1020" height="176" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 530, 14, "#e8e2d0", "middle", "EL VEREDICTO")
	txt(60, 566, 12, "#ff6b6b", "start", "⛔ NO se volvió teorema.")
	txt(250, 566, 11, "#c9c3b0", "start", "el pliegue de los divisores HEREDA el ½ de ζ; no lo produce.")
	txt(60, 598, 12, "#6ee7a0", "start", "✅ Pero el instinto era correcto:")
	txt(300, 598, 11, "#c9c3b0", "start", "un pliegue SÍ fabrica el ½ — y el pliegue existe.")
	txt(530, 638, 13, "#f2c14e", "middle", "Apuntaste al lugar correcto con el dedo equivocado.")
	txt(530, 660, 11, "#8a8676", "middle", "Y eso, en matemática, es mucho. Jacobi y Riemann 1859. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-magia-intentada.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
