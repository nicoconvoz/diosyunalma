// Command laley executes the captain's instruction literally.
//
// HIS ORDER: "no necesitamos analizar todo el infinito - necesitamos demostrar
// su FORMA nada mas. Que la ley sea pareja para todos los resultados. Lleva
// x=y <=> x-y=0 a la dimension 0 y armonizalo con la relacion 1/2. Quiero sacar
// una ley de ahi. Si probamos esa ley los ceros son totalmente explicables en
// todos los casos."
//
// The strategy is right, and this program does exactly what he says. Two things
// come out of it: a law that is exactly as even as he demanded, and an honest
// account of why that particular law cannot be the one that closes it - plus
// where the one that CAN is already sitting in this laboratory's own record.
//
// STEP 1 - HIS EQUALITY, CARRIED INTO DIMENSION 0. A pearl rho = beta + i*gamma
// goes to w = 1 - 1/rho. Sitting on the cable means beta = 1/2, and in the disk
// that means |w| = 1, that is |w|^2 = 1, that is
//
//	w * conj(w) - 1 = 0
//
// which is his x - y = 0 with x = w*conj(w) and y = 1. Expand it:
//
//	|w|^2 = (1 - 1/s)(1 - 1/sbar) = 1 - (s + sbar)/|s|^2 + 1/|s|^2
//	      = 1 - (2*beta - 1)/|s|^2
//
//	⟹  |w|^2 - 1 = -(2*beta - 1) / (beta^2 + gamma^2)
//
// THERE IS THE LAW, AND THERE IS HIS HALF, ALONE, IN THE NUMERATOR. The
// denominator is a sum of squares, so it is positive always, for every pearl,
// with no exceptions and no cases. The entire sign of the left side is decided
// by (2*beta - 1) - which is zero exactly when beta = 1/2.
//
//	beta > 1/2  →  |w| < 1  →  the pearl falls INSIDE the disk
//	beta = 1/2  →  |w| = 1  →  the pearl sits ON THE SKIN
//	beta < 1/2  →  |w| > 1  →  the pearl falls OUTSIDE
//
// That is exactly the "even law for all results" he asked for.
//
// STEP 2 - AND HERE IS WHY IT CANNOT BE THE ONE. It is an IDENTITY. It holds for
// every complex number in the plane, zero of zeta or not, prime or garbage. A
// thing that cannot fail cannot be proven, and proving it proves nothing. This
// program demonstrates that by running it on numbers that have nothing to do
// with zeta. It is the ninth time this laboratory meets that trap.
//
// The law does not answer the question. It TRANSLATES it: "are all pearls on the
// cable" becomes "is 2*beta - 1 always zero". Same question, new clothes - which
// is exactly what Finding 274 already said about a change of coordinate.
//
// STEP 3 - BUT THE LAW HE IS REACCHING FOR EXISTS, AND THIS LABORATORY ALREADY
// HAS IT. It is Li's criterion, and Finding 232 wrote it in dimension 0:
//
//	lambda_n = sum over pairs {rho, rhobar} of [ 2 - 2*Re(w^n) ]
//
// and RH is true if and only if lambda_n >= 0 for every n. That IS a law even
// for all results, and unlike the identity above IT CAN FAIL.
//
// Now harmonise THAT with his half, which is the step he asked for. Each pair's
// contribution splits exactly:
//
//	2 - 2*Re(w^n) = |1 - w^n|^2 + (1 - |w|^(2n))
//
// and the two halves of that sum say different things:
//
//	|1 - w^n|^2        is >= 0 ALWAYS, whatever the pearl does
//	1 - |w|^(2n)       is >= 0 exactly when |w| <= 1, that is when beta >= 1/2
//
// So the second term is the HALF's own term, and it is the only one that can go
// negative. Li's criterion in dimension 0 says: the unconditional part must
// always dominate the deficit that off-line pearls open. That is the law worth
// proving - and nobody has.
//
// PRE-REGISTERED PREDICTIONS, written before running:
//  1. The identity holds on all 38 measured pearls to ~1e-15.
//  2. It holds equally on numbers that are NOT zeros of zeta - which is what
//     makes it an identity and not a discovery.
//  3. The Li decomposition holds exactly, and for an off-line point the second
//     term goes negative while the first stays positive.
//
// Reproduce: go run ./cmd/laley
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66}
	fact := []float64{2, 24, 720, 40320, 3628800}
	poch := s
	for k := 1; k <= 5; k++ {
		if k > 1 {
			poch *= (s + complex(float64(2*k-3), 0)) * (s + complex(float64(2*k-2), 0))
		}
		sum += complex(B[k-1]/fact[k-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*k-1), 0))*lnN)
	}
	return sum
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

// w es el cambiaformas: lleva un punto del libro a la dimension 0.
func w(s complex128) complex128 { return 1 - 1/s }

// izquierda es |w|^2 - 1, calculado sin la ley.
func izquierda(s complex128) float64 {
	ww := w(s)
	return real(ww*cmplx.Conj(ww)) - 1
}

// derecha es -(2*beta - 1)/(beta^2 + gamma^2), la ley del capitan.
func derecha(s complex128) float64 {
	b, g := real(s), imag(s)
	return -(2*b - 1) / (b*b + g*g)
}

func main() {
	fmt.Println("⚖️  LA LEY — x = y ⟺ x − y = 0, llevado a la dimensión 0 y armonizado con el ½")
	fmt.Println("\n   Su orden: «no necesitamos analizar todo el infinito, necesitamos demostrar")
	fmt.Println("   su FORMA. Que la ley sea pareja para todos los resultados».")
	fmt.Println("\n   La estrategia es correcta, y acá está ejecutada al pie de la letra.")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · SU IGUALDAD, LLEVADA A LA DIMENSIÓN 0")
	fmt.Println("\n   Estar en el cable es β = ½. En el disco eso es |w| = 1, o sea |w|² = 1,")
	fmt.Println("   o sea **w·w̄ − 1 = 0** — que es su x − y = 0 con x = w·w̄ e y = 1.")
	fmt.Println("\n   Desarrollándolo:")
	fmt.Println("\n        |w|² = (1 − 1/s)(1 − 1/s̄) = 1 − (s + s̄)/|s|² + 1/|s|²")
	fmt.Println("             = 1 − (2β − 1)/|s|²")
	fmt.Println("\n   ⟹ ⚡ **|w|² − 1 = −(2β − 1) / (β² + γ²)**")
	fmt.Println("\n   **AHÍ ESTÁ LA LEY, Y AHÍ ESTÁ SU ½, SOLO, EN EL NUMERADOR.**")
	fmt.Println("   El denominador es una suma de cuadrados: **positivo siempre, para toda")
	fmt.Println("   perla, sin excepciones y sin casos.** Todo el signo lo decide (2β − 1),")
	fmt.Println("   que vale cero exactamente cuando β = ½.")

	fmt.Printf("\nbuscando las perlas…\n")
	ps := perlas(120)
	fmt.Printf("perlas encontradas: %d\n", len(ps))

	fmt.Println("\n   Verificada sobre nuestras propias perlas:")
	fmt.Println("\n        γ            |w|² − 1 (a mano)     la ley           diferencia")
	peor := 0.0
	for i, g := range ps {
		s := complex(0.5, g)
		iz, de := izquierda(s), derecha(s)
		d := math.Abs(iz - de)
		if d > peor {
			peor = d
		}
		if i < 5 {
			fmt.Printf("   %12.6f %20.3e %16.3e %14.2e\n", g, iz, de, d)
		}
	}
	fmt.Printf("\n        peor diferencia sobre las %d perlas ......... %.2e\n", len(ps), peor)

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · Y ES PAREJA PARA TODOS LOS RESULTADOS — los tres casos, sin excepción")
	fmt.Println("\n        β         |w|² − 1        dónde cae la perla")
	for _, b := range []float64{0.30, 0.40, 0.50, 0.60, 0.70, 0.8085} {
		s := complex(b, 85.699)
		iz := izquierda(s)
		var donde string
		switch {
		case iz > 1e-14:
			donde = "AFUERA del disco (β < ½)"
		case iz < -1e-14:
			donde = "ADENTRO del disco (β > ½)"
		default:
			donde = "⚡ EN LA PIEL (β = ½)"
		}
		fmt.Printf("   %8.4f %15.3e        %s\n", b, iz, donde)
	}
	fmt.Println("\n   ⟹ **Un solo renglón decide los tres casos, y el que decide es el ½.**")
	fmt.Println("   Eso es exactamente la ley pareja que pidió.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚠️ PERO NO SE PUEDE DEMOSTRAR, PORQUE ES UNA IDENTIDAD")
	fmt.Println("   Vale para CUALQUIER número complejo, sea cero de zeta o sea basura.")
	fmt.Println("   Probémoslo con números que no tienen nada que ver:")
	fmt.Println("\n        s                        |w|² − 1        la ley         diferencia")
	basura := []complex128{
		complex(3, 7), complex(-2, 0.5), complex(0.9, 1e-3),
		complex(100, -40), complex(0.5, 1),
	}
	peorB := 0.0
	for _, s := range basura {
		iz, de := izquierda(s), derecha(s)
		d := math.Abs(iz - de)
		if d > peorB {
			peorB = d
		}
		fmt.Printf("   %-22s %14.6f %14.6f %14.2e\n",
			fmt.Sprintf("%.3f%+.3fi", real(s), imag(s)), iz, de, d)
	}
	fmt.Printf("\n        peor diferencia sobre pura basura ........... %.2e\n", peorB)
	fmt.Println("\n   ⟹ **Cierra igual de perfecto con basura que con perlas.** Y una cosa que")
	fmt.Println("   no puede fallar no se puede demostrar: **demostrarla no demuestra nada**.")
	fmt.Println("   **NOVENA aparición de la trampa del 0.0e+00 en este laboratorio.**")
	fmt.Println("\n   📌 La ley no CONTESTA la pregunta: la TRADUCE. «¿están todas las perlas en")
	fmt.Println("   el cable?» se vuelve «¿es 2β − 1 siempre cero?». Misma pregunta, ropa nueva.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚡⚡ PERO LA LEY QUE USTED BUSCA EXISTE — Y YA LA TENEMOS")
	fmt.Println("   Es el criterio de Li, y F232 lo escribió en la dimensión 0:")
	fmt.Println("\n        λₙ = Σ sobre pares {ρ, ρ̄} de [ 2 − 2·Re(wⁿ) ]")
	fmt.Println("\n   y **RH es cierta si y sólo si λₙ ≥ 0 para todo n**. Ésa SÍ es una ley pareja")
	fmt.Println("   para todos los resultados — y a diferencia de la anterior, **PUEDE FALLAR**.")
	fmt.Println("\n   Ahora armonizala con su ½, que es el paso que pidió. Cada par se parte:")
	fmt.Println("\n        2 − 2·Re(wⁿ) = |1 − wⁿ|² + (1 − |w|²ⁿ)")
	fmt.Println("\n   Verificado numéricamente:")
	fmt.Println("\n        γ         n     2−2Re(wⁿ)    |1−wⁿ|²    1−|w|²ⁿ      diferencia")
	peorL := 0.0
	for _, g := range []float64{ps[0], ps[1], ps[2]} {
		for _, n := range []int{1, 3, 7} {
			ww := w(complex(0.5, g))
			wn := cmplx.Pow(ww, complex(float64(n), 0))
			izq := 2 - 2*real(wn)
			t1 := cmplx.Abs(1-wn) * cmplx.Abs(1-wn)
			t2 := 1 - math.Pow(cmplx.Abs(ww), float64(2*n))
			d := math.Abs(izq - (t1 + t2))
			if d > peorL {
				peorL = d
			}
			fmt.Printf("   %10.5f %5d %12.6f %10.6f %11.3e %14.2e\n", g, n, izq, t1, t2, d)
		}
	}
	fmt.Printf("\n        peor diferencia de la descomposición ........ %.2e\n", peorL)
	fmt.Println("\n   ⟹ **Y las dos mitades dicen cosas DISTINTAS:**")
	fmt.Println("\n        |1 − wⁿ|²      es ≥ 0 SIEMPRE, haga lo que haga la perla")
	fmt.Println("        1 − |w|²ⁿ      es ≥ 0 **exactamente cuando |w| ≤ 1**, o sea β ≥ ½")
	fmt.Println("\n   **El segundo término es el término DEL ½, y es el único que puede ponerse")
	fmt.Println("   negativo.** Sobre el cable vale cero clavado — mírelo en la columna.")

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · QUÉ LE HACE UNA PERLA SUELTA A ESA LEY")
	fmt.Println("   Tomemos el punto fuera de línea del collar hermano (β = 0.8085) y su espejo")
	fmt.Println("   obligatorio por la ecuación funcional (β = 0.1915), y miremos el término del ½:")
	fmt.Println("\n        β        |w|        1 − |w|²ⁿ (n=1)    1 − |w|²ⁿ (n=10)")
	for _, b := range []float64{0.8085, 0.5, 0.1915} {
		ww := w(complex(b, 85.699))
		m := cmplx.Abs(ww)
		fmt.Printf("   %8.4f %10.6f %18.6e %18.6e\n", b, m,
			1-math.Pow(m, 2), 1-math.Pow(m, 20))
	}
	fmt.Println("\n   ⟹ **La perla de arriba del cable aporta positivo, la de abajo aporta")
	fmt.Println("   NEGATIVO, y el déficit crece con n.** Por eso Li puede fallar: hace falta")
	fmt.Println("   que la parte incondicional |1−wⁿ|² tape ese agujero **para todo n, siempre**.")
	fmt.Println("   **Ésa es la ley que vale la pena demostrar. Y nadie la demostró.**")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("✅ **SU LEY EXISTE, SALE EXACTA, Y ES TAN PAREJA COMO LA PIDIÓ:**")
	fmt.Println("\n        |w|² − 1 = −(2β − 1) / (β² + γ²)")
	fmt.Printf("\n  Un solo renglón, sin casos, sin excepciones, verificado sobre %d perlas a\n", len(ps))
	fmt.Printf("  %.0e. Y el ½ aparece solo, en el numerador, decidiendo todo el signo.\n", peor)
	fmt.Println("\n⚠️ **PERO NO SE PUEDE DEMOSTRAR: ES UNA IDENTIDAD.** Cierra igual con perlas")
	fmt.Printf("  que con basura (%.0e). **Novena trampa del 0.0e+00.** No contesta la\n", peorB)
	fmt.Println("  pregunta: la traduce.")
	fmt.Println("\n⚡ **Y SIN EMBARGO SU ESTRATEGIA ES LA CORRECTA, Y LA LEY QUE BUSCA EXISTE:**")
	fmt.Println("  es el criterio de Li, λₙ ≥ 0 para todo n, que este laboratorio ya escribió en")
	fmt.Println("  la dimensión 0 en F232. Ésa **puede fallar**, y demostrarla **ES** demostrar RH.")
	fmt.Println("\n  Y armonizada con su ½ como usted pidió, se parte en dos mitades exactas:")
	fmt.Println("  una que siempre suma, y otra —la del ½— que es la única que puede restar.")
	fmt.Println("\n⚖️ Y el límite honesto: Li es de 1997, se conoce, y nadie pudo probar λₙ ≥ 0.")
	fmt.Println("  Usted no encontró la puerta cerrada: encontró CUÁL es la puerta. Todavía no.")

	escribirLamina(ps, peor, peorB, peorL)
}

func escribirLamina(ps []float64, peor, peorB, peorL float64) {
	var b strings.Builder
	W, H := 1600.0, 1080.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⚖️ LA LEY — x = y ⟺ x − y = 0, en la dimensión 0, armonizado con el ½</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">la ley existe y es exacta · pero es una identidad · y la que SÍ puede fallar ya estaba en el registro</text>
`, W, H, W, H, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="40" y="102" width="1520" height="200" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="800" y="136" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡ SU LEY, EJECUTADA AL PIE DE LA LETRA</text>
<text x="800" y="172" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">estar en el cable es β = ½ · en el disco eso es |w| = 1 · o sea w·w̄ − 1 = 0, que es su x − y = 0</text>
<text x="800" y="222" font-size="27" text-anchor="middle" font-family="monospace" fill="#ffd98a">|w|² − 1 = −(2β − 1) / (β² + γ²)</text>
<text x="800" y="258" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">y ahí está su ½, SOLO, en el numerador — decidiendo todo el signo</text>
<text x="800" y="286" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">el denominador es suma de cuadrados: positivo siempre, sin casos, sin excepciones · verificado sobre %d perlas a %.0e</text>
`, len(ps), peor)

	fmt.Fprintf(&b, `<rect x="40" y="322" width="740" height="230" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="410" y="356" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">PAREJA PARA TODOS LOS RESULTADOS</text>
<text x="80" y="398" font-size="16" font-family="monospace" fill="#7fb2ff">β &gt; ½   →  |w| &lt; 1  →  ADENTRO del disco</text>
<text x="80" y="432" font-size="16" font-family="monospace" fill="#7ee0c0">β = ½   →  |w| = 1  →  EN LA PIEL</text>
<text x="80" y="466" font-size="16" font-family="monospace" fill="#ff8fa0">β &lt; ½   →  |w| &gt; 1  →  AFUERA del disco</text>
<text x="410" y="510" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">un solo renglón decide los tres casos</text>
<text x="410" y="534" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">y el que decide es el ½</text>`)

	fmt.Fprintf(&b, `<rect x="820" y="322" width="740" height="230" rx="12" fill="#33221c" stroke="#c0392b"/>
<text x="1190" y="356" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚠️ PERO NO SE PUEDE DEMOSTRAR</text>
<text x="1190" y="390" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Es una IDENTIDAD: vale para cualquier número complejo,</text>
<text x="1190" y="414" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">sea cero de zeta o sea basura.</text>
<text x="1190" y="452" font-size="16" text-anchor="middle" font-family="monospace" fill="#ff8fa0">con perlas: %.0e   ·   con basura: %.0e</text>
<text x="1190" y="486" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Una cosa que no puede fallar no se puede demostrar.</text>
<text x="1190" y="512" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">No CONTESTA la pregunta: la TRADUCE.</text>
<text x="1190" y="538" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">NOVENA aparición de la trampa del 0.0e+00</text>
`, peor, peorB)

	fmt.Fprintf(&b, `<rect x="40" y="572" width="1520" height="300" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="800" y="606" font-size="19" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚡⚡ PERO LA LEY QUE BUSCA EXISTE — Y YA ESTABA EN EL REGISTRO</text>
<text x="800" y="640" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el criterio de Li, escrito en la dimensión 0 por nuestro F232:</text>
<text x="800" y="682" font-size="22" text-anchor="middle" font-family="monospace" fill="#ffd98a">λₙ = Σ [ 2 − 2·Re(wⁿ) ]        y      RH ⟺ λₙ ≥ 0 para todo n</text>
<text x="800" y="716" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">ésa SÍ es pareja para todos los resultados — y a diferencia de la otra, PUEDE FALLAR</text>
<text x="800" y="756" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">y armonizada con su ½, como usted pidió, se parte en dos mitades exactas (verificado a %.0e):</text>
<text x="800" y="796" font-size="21" text-anchor="middle" font-family="monospace" fill="#dce8f7">2 − 2·Re(wⁿ)  =  |1 − wⁿ|²  +  (1 − |w|²ⁿ)</text>
<text x="440" y="832" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">≥ 0 SIEMPRE, haga lo que haga la perla</text>
<text x="1130" y="832" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">≥ 0 sólo si β ≥ ½ — EL TÉRMINO DEL ½</text>
<text x="800" y="858" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">el segundo es el único que puede restar · sobre el cable vale cero clavado</text>
`, peorL)

	fmt.Fprintf(&b, `<rect x="40" y="892" width="1520" height="160" rx="12" fill="#1a1030" stroke="#5a4fa8"/>
<text x="800" y="926" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ LO QUE ESTO ES, Y LO QUE NO</text>
<text x="800" y="960" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Su ley sale exacta y es tan pareja como la pidió — pero es una identidad, y demostrarla no demuestra nada.</text>
<text x="800" y="986" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">La que sí hay que demostrar es λₙ ≥ 0: es de Li, es de 1997, se conoce, y nadie pudo.</text>
<text x="800" y="1022" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Usted no encontró la puerta cerrada: encontró CUÁL es la puerta. Y eso vale.</text>
</svg>
`)

	if err := os.WriteFile("la-ley.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: la-ley.svg")
}
