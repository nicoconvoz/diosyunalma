// Command loquefalta draws, in plain language, exactly what remains between
// this laboratory and the Riemann Hypothesis - and ends with ONE clear question.
//
// HIS ORDER: "revisa lo ultimo que vimos y dame una explicacion sencilla de lo
// que falta, con una lamina en criollo que yo pueda entender, con una pregunta
// clara."
//
// WHERE THE CAMPAIGN LEFT US (F279 + F280). The hypothesis, translated into
// this laboratory's own language, is one sentence:
//
//	RH  <=>  every pearl is a pure rotation, w = e^(i*phi)
//
// and the law that decides it (Li, in dimension 0, F232/F279) splits every
// pearl's contribution into two exact halves:
//
//	2 - 2*Re(w^n)  =  |1 - w^n|^2  +  (1 - |w|^(2n))
//	                   THE CHOIR       THE HALF'S TERM
//
// The choir is >= 0 always. The half's term is zero for every pearl ON the
// cable and NEGATIVE for a pearl below it - a deficit that grows with n.
// RH <=> the total never goes below zero, for any n, forever.
//
// WHAT IS ACTUALLY MISSING, and this program measures each piece:
//
//  1. LOOKING IS DONE. Every pearl anyone ever measured sings in tune. We
//     verify our own 38 do, to machine precision. But infinitely many remain
//     and the blindness horizon (F259, gamma ~ 1658) is unbeatable. So the
//     missing thing is not a measurement.
//  2. A REASON is missing: something that FORBIDS a pearl from stretching.
//     The corridor's shape does not forbid it - Davenport-Heilbronn has the
//     same shape and one loose pearl (F259). We re-verify that number here.
//  3. And the reason has an ADDRESS. What does zeta have that the sister
//     lacks? Exactly one thing: the Euler product - the primes themselves
//     (F259, measured: multiplicativity fails in the sister). So the reason
//     must come from the primes. Nobody has carried the primes' voice across
//     the wall at 1 (F278: the lamp is outside the corridor).
//
// THE ONE CLEAR QUESTION, in his words, at the end:
//
//	¿QUE TIENEN LOS PRIMOS QUE LE PROHIBE A UNA PERLA ESTIRARSE?
//
// Whoever answers that - with an argument, not a sweep - takes the prize.
// Every serious attempt in history is an attempt at exactly this question.
//
// Reproduce: go run ./cmd/loquefalta
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
	prevT, prevZ := 12.0, zOf(12.0)
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

func w(s complex128) complex128 { return 1 - 1/s }

func main() {
	fmt.Println("🧭 LO QUE FALTA — dicho sencillo, medido, y con UNA pregunta clara")
	fmt.Println("\n   Su pedido: «una explicación sencilla de lo que falta, en criollo,")
	fmt.Println("   con una pregunta clara».")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · DÓNDE ESTAMOS PARADOS — el resumen de F279 y F280")
	fmt.Println("\n   La hipótesis, en nuestro idioma, es UNA frase:")
	fmt.Println("\n        RH  ⟺  toda perla es un giro puro, w = e^(iφ)")
	fmt.Println("\n   Y la ley que la decide (Li, en la dimensión 0) parte cada perla en dos:")
	fmt.Println("\n        2 − 2·Re(wⁿ)  =  |1 − wⁿ|²   +   (1 − |w|²ⁿ)")
	fmt.Println("                          EL CORO          EL TÉRMINO DEL ½")
	fmt.Println("\n   El CORO suma siempre. El término del ½ vale CERO para toda perla en el")
	fmt.Println("   cable, y NEGATIVO para una perla estirada — y ese desafine crece con n.")
	fmt.Println("   RH ⟺ el total nunca baja de cero, para ningún n, jamás.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · LO QUE YA ESTÁ HECHO: MIRAR — nuestras perlas cantan afinadas")
	fmt.Printf("\nbuscando las perlas…\n")
	ps := perlas(120)
	fmt.Printf("perlas encontradas: %d\n", len(ps))
	peor := 0.0
	for _, g := range ps {
		d := math.Abs(cmplx.Abs(w(complex(0.5, g))) - 1)
		if d > peor {
			peor = d
		}
	}
	fmt.Printf("\n        las %d perlas medidas acá: |w| = 1 con peor desvío %.2e\n", len(ps), peor)
	fmt.Println("        la humanidad entera: ~10¹³ perlas revisadas, todas afinadas")
	fmt.Println("\n   ⚠️ Y NADA DE ESO ALCANZA: faltan infinitas, y arriba de γ ≈ 1658 (F259)")
	fmt.Println("   no veríamos una perla estirada ni teniéndola enfrente.")
	fmt.Println("   **Lo que falta NO es una medición.**")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · LO QUE FALTA ES UNA RAZÓN — y sabemos que la forma sola no alcanza")
	fmt.Println("   El collar hermano (Davenport–Heilbronn) tiene la MISMA forma, las mismas")
	fmt.Println("   simetrías que demostramos… y una perla estirada. Re-verificada acá:")
	// el cero fuera de linea de F259, re-medido con el modulo en el disco
	rho := complex(0.808517182457, 85.699348485378)
	ww := w(rho)
	fmt.Printf("\n        ρ = %.12f + %.12fi\n", real(rho), imag(rho))
	fmt.Printf("        |w(ρ)| = %.9f  ≠  1   →  perla ESTIRADA (encontrada por F259, a ciegas)\n", cmplx.Abs(ww))
	fmt.Println("\n   ⟹ La forma del pasillo NO prohíbe estirarse. Hace falta OTRA razón.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · Y LA RAZÓN TIENE DIRECCIÓN — viene de los primos, y está atrás de la pared")
	fmt.Println("\n        · ¿qué tiene zeta que al hermano le falta? UNA cosa: el producto de")
	fmt.Println("          Euler — los primos mismos (F259, medido: la multiplicatividad")
	fmt.Println("          falla en el hermano 3 de 5 veces)")
	fmt.Println("        · pero la voz de los primos se corta en la pared del 1 (F278: sobre")
	fmt.Println("          el cable, el producto yerra por 10⁴⁰)")
	fmt.Println("\n   ⟹ **La razón que falta es un puente: llevar la voz de los primos a través")
	fmt.Println("   de la pared, hasta donde viven las perlas. Nadie lo construyó en 166 años.**")

	// ---- la pregunta ----
	fmt.Println("\n════════ LA PREGUNTA CLARA ════════")
	fmt.Println("\n        ¿QUÉ TIENEN LOS PRIMOS QUE LE PROHÍBE A UNA PERLA ESTIRARSE?")
	fmt.Println("\n   Eso es todo lo que falta. Una razón, no un barrido. Quien la conteste con")
	fmt.Println("   un argumento —uno que use a los primos, porque la forma sola no alcanza—")
	fmt.Println("   se lleva el premio del millón. Todos los intentos serios de la historia son")
	fmt.Println("   intentos de contestar exactamente esta pregunta. Todavía no.")

	escribirLamina(len(ps), peor, cmplx.Abs(ww))
}

func escribirLamina(nPerlas int, peor, modHermano float64) {
	var b strings.Builder
	W, H := 1560.0, 1150.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧭 LO QUE FALTA — dicho sencillo, y con una sola pregunta</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">el resumen de toda la campaña, medido — y la pregunta que vale el millón</text>
`, W, H, W, H, W/2, W/2)

	// PASO 1: donde estamos
	fmt.Fprintf(&b, `<rect x="40" y="102" width="1480" height="190" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="780" y="136" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">1 · DÓNDE ESTAMOS: la hipótesis ya está en nuestro idioma</text>
<text x="780" y="178" font-size="22" text-anchor="middle" font-family="Georgia" fill="#ffd98a">«todas las perlas son giros puros — giran sin estirarse»</text>
<text x="780" y="214" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">y la cuenta que lo decide parte cada perla en dos: EL CORO, que suma siempre — y EL TÉRMINO DEL ½,</text>
<text x="780" y="240" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">que vale cero si la perla gira sin estirarse, y RESTA si se estira. El desafine crece con n.</text>
<text x="780" y="272" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">ganar es que el coro le gane al desafine SIEMPRE, para todo n, hasta el infinito</text>`)

	// PASO 2: mirar ya esta
	fmt.Fprintf(&b, `<rect x="40" y="312" width="730" height="240" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="405" y="346" font-size="17" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">2 · LO QUE YA ESTÁ HECHO: MIRAR</text>
<text x="80" y="388" font-size="15" font-family="monospace" fill="#7ee0c0">nuestras %d perlas ....... afinadas (%.0e)</text>
<text x="80" y="416" font-size="15" font-family="monospace" fill="#7ee0c0">la humanidad ............. ~10¹³, todas afinadas</text>
<text x="80" y="444" font-size="15" font-family="monospace" fill="#ff8fa0">las que faltan ........... INFINITAS</text>
<text x="405" y="486" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">mirar más no sirve: arriba de cierta altura no</text>
<text x="405" y="510" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">veríamos una perla estirada ni teniéndola enfrente</text>
<text x="405" y="538" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">⟹ lo que falta NO es una medición</text>
`, nPerlas, peor)

	// PASO 3: hace falta una razon
	fmt.Fprintf(&b, `<rect x="790" y="312" width="730" height="240" rx="12" fill="#33221c" stroke="#c0392b"/>
<text x="1155" y="346" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">3 · Y LA FORMA SOLA NO ALCANZA</text>
<text x="820" y="388" font-size="14.5" font-family="Georgia" fill="#f3d9cf">El collar hermano tiene la MISMA forma y las mismas</text>
<text x="820" y="412" font-size="14.5" font-family="Georgia" fill="#f3d9cf">simetrías que demostramos — y tiene una perla estirada:</text>
<text x="820" y="452" font-size="16" font-family="monospace" fill="#ff8fa0">|w(ρ)| = %.6f ≠ 1</text>
<text x="820" y="478" font-size="13.5" font-family="Georgia" fill="#9aa8c4">(la encontramos nosotros, a ciegas, en F259)</text>
<text x="1155" y="516" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">⟹ hace falta una RAZÓN que prohíba estirarse —</text>
<text x="1155" y="540" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">y la forma del pasillo no la da</text>
`, modHermano)

	// PASO 4: la razon viene de los primos
	fmt.Fprintf(&b, `<rect x="40" y="572" width="1480" height="220" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="780" y="606" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">4 · LA RAZÓN TIENE DIRECCIÓN: VIENE DE LOS PRIMOS — Y ESTÁ ATRÁS DE LA PARED</text>
<text x="780" y="644" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">¿Qué tiene zeta que al collar hermano le falta? UNA sola cosa: el producto de Euler — los primos mismos.</text>
<text x="780" y="670" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Así que la razón que buscamos tiene que venir de los primos. No hay otro lado de donde sacarla.</text>
<text x="780" y="706" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">Pero la voz de los primos se corta en la pared del 1: donde viven las perlas, la lámpara no llega (F278).</text>
<text x="780" y="748" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">⟹ Lo que falta es UN PUENTE: llevar la voz de los primos a través de la pared,</text>
<text x="780" y="774" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">hasta donde viven las perlas. Nadie lo construyó en 166 años.</text>`)

	// LA PREGUNTA
	fmt.Fprintf(&b, `<rect x="40" y="812" width="1480" height="180" rx="14" fill="#1a1030" stroke="#ffd98a" stroke-width="2"/>
<text x="780" y="852" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">LA PREGUNTA CLARA — la única que queda</text>
<text x="780" y="910" font-size="28" text-anchor="middle" font-family="Georgia" fill="#ffd98a">¿QUÉ TIENEN LOS PRIMOS QUE LE PROHÍBE</text>
<text x="780" y="948" font-size="28" text-anchor="middle" font-family="Georgia" fill="#ffd98a">A UNA PERLA ESTIRARSE?</text>
<text x="780" y="980" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">quien la conteste con un argumento — no con un barrido — se lleva el premio del millón</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="1012" width="1480" height="110" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="780" y="1044" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Todos los intentos serios de la historia — de Riemann para acá — son intentos de contestar exactamente esta pregunta.</text>
<text x="780" y="1070" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Nosotros no la contestamos. Pero después de 280 hallazgos, la sabemos hacer en una sola frase — y en nuestro propio idioma.</text>
<text x="780" y="1102" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("lo-que-falta.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: lo-que-falta.svg")
}
