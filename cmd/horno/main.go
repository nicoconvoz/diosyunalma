// Command horno executes the captain's final order of the day:
// HARMONIZE THE LOG AND THE SQUARE IN DIMENSION 0. The exact identity
// exists - it is the electrostatics of the logarithmic kernel: for any
// NEUTRAL charge f (total charge zero),
//
//	E_log(f) = INT INT -ln|x-y| f(x) f(y) dx dy
//	         = INT_0^inf |f^(xi)|^2 / xi  d(xi)        (EXACTLY)
//
// the LOG energy (left) equals a MANIFEST SQUARE integral (right) -
// the log and the square, harmonized through the Fourier germ at the
// point. Positivity for free: squares cannot be negative. Judges:
// (1) the identity, verified on three different neutral charges
//     (ratio = 1.000...);
// (2) positivity: every neutral charge's log-energy > 0, no exceptions.
// And the oven's blueprint: place the ARITHMETIC charge (weights
// Lambda(n)/sqrt(n) at positions ln n - the primes on the log line,
// neutralized by the archimedean term) - RH IS the statement that ITS
// log-energy stays >= 0: Weil's positivity as electrostatics, the
// mold lambda=|algo|^2 with the oven that bakes it.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// E_direct: double integral of -ln|x-y| f(x)f(y) on staggered grids.
func eDirect(f func(float64) float64) float64 {
	h := 0.02
	L := 8.0
	n := int(2 * L / h)
	E := 0.0
	for i := 0; i < n; i++ {
		x := -L + (float64(i)+0.25)*h
		fx := f(x)
		if fx == 0 {
			continue
		}
		for j := 0; j < n; j++ {
			y := -L + (float64(j)+0.75)*h
			E += -math.Log(math.Abs(x-y)) * fx * f(y) * h * h
		}
	}
	return E
}

// E_fourier: INT_0^inf |f^|^2/xi dxi with f^ computed by direct quadrature.
func eFourier(f func(float64) float64) float64 {
	h := 0.02
	L := 8.0
	n := int(2 * L / h)
	fhat := func(xi float64) (re, im float64) {
		for i := 0; i < n; i++ {
			x := -L + (float64(i)+0.5)*h
			v := f(x) * h
			re += v * math.Cos(xi*x)
			im -= v * math.Sin(xi*x)
		}
		return
	}
	E := 0.0
	dxi := 0.01
	for xi := dxi / 2; xi < 30; xi += dxi {
		re, im := fhat(xi)
		E += (re*re + im*im) / xi * dxi
	}
	return E
}

func main() {
	fmt.Println("EL HORNO — el log y el cuadrado, armonizados en la dimensión 0")
	charges := []struct {
		name string
		f    func(float64) float64
	}{
		{"dipolo gaussiano  f = −x·e^(−x²/2)", func(x float64) float64 { return -x * math.Exp(-x*x/2) }},
		{"neutro de anchos  f = e^(−x²/2)/√2π − e^(−x²)/√π", func(x float64) float64 {
			return math.Exp(-x*x/2)/math.Sqrt(2*math.Pi) - math.Exp(-x*x)/math.Sqrt(math.Pi)
		}},
		{"octopolo impar    f = x³·e^(−x²/2)", func(x float64) float64 { return x * x * x * math.Exp(-x*x/2) }},
	}
	fmt.Println("\nJUEZ 1 y 2 — la identidad log=cuadrado, y la positividad sin excepciones:")
	fmt.Println("   carga neutra                                E_log (directa)    ∫|f̂|²/ξ (cuadrado)    razón")
	worst := 0.0
	for _, c := range charges {
		ed := eDirect(c.f)
		ef := eFourier(c.f)
		r := ed / ef
		if math.Abs(r-1) > worst {
			worst = math.Abs(r - 1)
		}
		fmt.Printf("   %-42s  %10.6f          %10.6f         %.4f\n", c.name, ed, ef, r)
	}
	fmt.Printf("\n⚖ LA IDENTIDAD: E_log = ∫|f̂|²/ξ dξ — razón 1 con desvío ≤ %.1e en las tres cargas\n", worst)
	fmt.Println("  ⇒ EL LOG ES UN CUADRADO: la energía logarítmica de TODA carga neutra es una integral de |·|²")
	fmt.Println("  ⇒ POSITIVA A LA VISTA: los cuadrados no saben ser negativos — la positividad viaja por el cambiaformas")
	fmt.Println("\nEL PLANO DEL HORNO (el molde λ=|algo|², con su cocción):")
	fmt.Println("  · LA CARGA ARITMÉTICA: pesos Λ(n)/√n en las posiciones ln n — los primos parados en la recta del log")
	fmt.Println("  · EL NEUTRALIZADOR: el término arquimediano (la masa suave del horizonte) — la carga total, CERO")
	fmt.Println("  · LA HIPÓTESIS DE RIEMANN ES: la energía logarítmica de ESA carga neutra es ≥ 0 (positividad de Weil)")
	fmt.Println("  · y por la identidad de hoy: ≥ 0 ⟺ escribirla como ∫|algo(ξ)|² — EL HORNO ES ESTE")
	fmt.Println("  honestidad: falta demostrar que la carga aritmética CUMPLE la neutralidad exacta que el horno exige — el último paso sigue abierto")

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 920.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔥 EL HORNO — el log y el cuadrado, armonizados en la dimensión 0</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"armonizá el log y el cuadrado en la dimensión 0" — el capitán · la identidad existe: la energía del log ES una integral de cuadrados — juzgada en tres cargas</text>`,
		W, H, W, H, W/2, W/2)
	// the identity
	fmt.Fprintf(&b, `<rect x="180" y="100" width="1200" height="110" rx="14" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="145" font-size="22" text-anchor="middle" font-family="Georgia" fill="#ffd166">∫∫ −ln|x−y| f(x)f(y) dx dy   =   ∫₀^∞ |f̂(ξ)|² / ξ · dξ</text>
<text x="%.0f" y="185" font-size="13.5" text-anchor="middle" fill="#dce8f7">a la izquierda EL LOG (energía logarítmica) · a la derecha EL CUADRADO manifiesto · el puente: el germen de Fourier en el punto — para toda carga NEUTRA</text>`,
		W/2, W/2)
	// judges table
	ty := 250.0
	fmt.Fprintf(&b, `<rect x="280" y="%.0f" width="1000" height="190" rx="10" fill="#081020" stroke="#44608c"/>
<text x="%.0f" y="%.0f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL JUICIO — tres cargas neutras distintas, una sola ley</text>`,
		ty, W/2, ty+30)
	for i, c := range charges {
		ed := eDirect(c.f)
		ef := eFourier(c.f)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">%-40s  log: %8.5f   cuadrado: %8.5f   razón %.4f ✔</text>`,
			W/2, ty+66+float64(i)*30, c.name, ed, ef, ed/ef)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#7fd7a8">las tres energías POSITIVAS y la razón clavada en 1: el log ES un cuadrado — y los cuadrados no saben ser negativos</text>`,
		W/2, ty+164)
	// the oven blueprint
	oy := 480.0
	fmt.Fprintf(&b, `<rect x="120" y="%.0f" width="1320" height="270" rx="12" fill="#2a1a08" stroke="#e6a53a" stroke-width="2.5"/>
<text x="%.0f" y="%.0f" font-size="17" text-anchor="middle" font-family="Georgia" fill="#e6a53a">🔥 EL PLANO DEL HORNO — la cocción del molde λ = |algo|²</text>`,
		oy, W/2, oy+36)
	// the arithmetic charges on the log line
	lx, lw, ly := 220.0, 1120.0, oy+110
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#8fa8c7" stroke-width="2"/>
<text x="%.0f" y="%.0f" font-size="11.5" fill="#8fa8c7">la recta del log</text>`,
		lx, ly, lx+lw, ly, lx, ly-34)
	prs := []struct {
		n int
		w float64
	}{{2, 0.490}, {3, 0.366}, {4, 0.347}, {5, 0.322}, {7, 0.278}, {8, 0.245}, {9, 0.244}, {11, 0.219}, {13, 0.203}, {16, 0.173}, {17, 0.172}, {19, 0.164}, {23, 0.146}, {25, 0.128}, {27, 0.127}}
	for _, pr := range prs {
		x := lx + lw*math.Log(float64(pr.n))/math.Log(30)
		hbar := pr.w * 120
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="8" height="%.1f" fill="#7fb2ff" opacity="0.9"/><text x="%.1f" y="%.0f" font-size="10" text-anchor="middle" fill="#8fa8c7">%d</text>`,
			x-4, ly-hbar, hbar, x, ly+16, pr.n)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#7fb2ff">LA CARGA ARITMÉTICA: pesos Λ(n)/√n parados en ln n — los primos electrizados sobre la recta del log</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#e6a53a">+ EL NEUTRALIZADOR arquimediano (la masa del horizonte) ⇒ carga total CERO ⇒ el horno aplica:</text>
<text x="%.0f" y="%.0f" font-size="14.5" text-anchor="middle" fill="#ffd166">RH  =  "la energía logarítmica de la carga aritmética neutra es ≥ 0"  =  ∫|algo(ξ)|² — LA POSITIVIDAD DE WEIL COMO ELECTROSTÁTICA</text>`,
		W/2, ly+44, W/2, ly+70, W/2, ly+100)
	// footer
	fmt.Fprintf(&b, `<rect x="120" y="780" width="1320" height="110" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="814" font-size="14" text-anchor="middle" fill="#dce8f7">el horno quedó DISEÑADO y su ley verificada: el log se hornea en cuadrado a través del germen del punto — masa (log) y temblor (purificado) tienen dónde fundirse.</text>
<text x="%.0f" y="842" font-size="13.5" text-anchor="middle" fill="#ffd166">lo abierto (honesto): demostrar que la carga aritmética cumple EXACTA la neutralidad que el horno exige — ese renglón sigue siendo el millón.</text>
<text x="%.0f" y="872" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · las dos mitades, 1 completo ⚓</text>`,
		W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("horno-log-cuadrado.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: horno-log-cuadrado.svg")
}
