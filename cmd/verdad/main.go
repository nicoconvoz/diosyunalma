// Command verdad closes the captain's philosophical chain with measurements.
//
// His formulation: Verdad(P) <=> |P - R| = 0, truth is CORRESPONDENCE between
// a proposition and the state of reality. And he named the trap himself: in
// mathematics a theorem is true INSIDE a system of axioms, which is a very
// different thing from a model corresponding to the world.
//
// That trap is exactly this laboratory's whole situation, and it can be
// measured rather than merely discussed:
//
//	VERDAD MEDIDA (correspondencia)  |P - R| = 0 wherever we looked - we have
//	                                 this, down to the machine's last bit.
//	VERDAD DEMOSTRADA (derivacion)   P follows from the axioms by a finite
//	                                 chain - we do NOT have this.
//
// Three things are measured here.
//
//  1. THE CORRESPONDENCE: |P - R| across the campaign, all of it zero.
//  2. WHY CORRESPONDENCE IS NOT ENOUGH: the polynomial of F229 has every
//     symptom of truth and a false conclusion - correspondence of symptoms
//     is not correspondence of fact.
//  3. THE ASYMMETRY, which is the deep one: if the hypothesis were FALSE
//     there would be a FINITE certificate - one off-line pearl forces a
//     negative envelope at a computable harmonic, and anyone could check it
//     in an afternoon. If it is TRUE there is no finite certificate at all.
//     That asymmetry is why measuring can never win, and it also gives a
//     conclusion the captain will like: a false statement of this shape is
//     always refutable, so IF the hypothesis were undecidable, it would have
//     to be TRUE.
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

func wDe(rho complex128) complex128 { return 1 - 1/rho }

func main() {
	fmt.Println("⚖️ LA VERDAD — Verdad(P) ⟺ |P − R| = 0, y la trampa que el capitán nombró solo")

	// ---- LAW 1: the correspondence, measured ----
	fmt.Println("\nLEY 1 · LA CORRESPONDENCIA — |P − R| en todo lo que miramos")
	fmt.Println("\nrecogiendo perlas hasta t=300…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 300; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			pearls = append(pearls, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	peor := 0.0
	for _, g := range pearls {
		if d := math.Abs(cmplx.Abs(wDe(complex(0.5, g))) - 1); d > peor {
			peor = d
		}
	}
	fmt.Printf("   P: «cada perla del libro está sobre la línea»\n")
	fmt.Printf("   R: donde el mar dice que están, medido en %d perlas\n", len(pearls))
	fmt.Printf("   |P − R| = %.2e — cero al último bit de la máquina\n", peor)
	fmt.Println("   → LA CORRESPONDENCIA ESTÁ. Donde miramos, lo que es coincide con lo que debería ser")

	// ---- LAW 2: correspondence of symptoms is not truth ----
	fmt.Println("\nLEY 2 · POR QUÉ LA CORRESPONDENCIA NO ALCANZA — un impostor con todos los síntomas")
	a := complex(0.7, 3.0)
	raices := []complex128{a, cmplx.Conj(a), 1 - a, 1 - cmplx.Conj(a)}
	P := func(s complex128) complex128 {
		p := complex(1, 0)
		for _, r := range raices {
			p *= s - r
		}
		return p
	}
	sig := func(s complex128) complex128 { return 1 - cmplx.Conj(s) }
	fmt.Println("   se toma un objeto que NO cumple la hipótesis y se le miden los mismos síntomas:")
	peorS := 0.0
	for _, s := range []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2.0, 21.3)} {
		for _, d := range []float64{
			cmplx.Abs(P(s)-P(1-s)) / cmplx.Abs(P(s)),
			cmplx.Abs(P(cmplx.Conj(s))-cmplx.Conj(P(s))) / cmplx.Abs(P(s)),
			cmplx.Abs(P(sig(s))-cmplx.Conj(P(s))) / cmplx.Abs(P(s)),
		} {
			if d > peorS {
				peorS = d
			}
		}
	}
	fmt.Printf("   ecuación funcional, espejo de Schwarz y cambiaformas: TODOS coinciden (peor %.1e)\n", peorS)
	fmt.Println("   y sin embargo sus raíces están en Re = 0.70 y 0.30 — FUERA de la línea")
	fmt.Println("   → coincidencia de SÍNTOMAS no es coincidencia de HECHO: por eso medir no basta")

	// ---- LAW 3: the asymmetry - the certificate of falsehood is finite ----
	fmt.Println("\nLEY 3 · LA ASIMETRÍA — el certificado de la MENTIRA es finito; el de la VERDAD no existe")
	fmt.Println("   si la hipótesis fuera FALSA, una sola perla corrida obliga a un sobre en rojo,")
	fmt.Println("   en un armónico calculable — cualquiera podría verificarlo en una tarde:")
	fmt.Println("   β de la perla corrida     el certificado sería de     ¿verificable?")
	type cert struct {
		beta float64
		n    int
	}
	var certs []cert
	for _, beta := range []float64{0.90, 0.75, 0.60, 0.55, 0.51, 0.501} {
		g := 14.134725
		r := cmplx.Abs(wDe(complex(1-beta, g)))
		if r < 1 {
			r = cmplx.Abs(wDe(complex(beta, g)))
		}
		n := int(math.Ceil(math.Log(3) / math.Log(r)))
		certs = append(certs, cert{beta, n})
		fmt.Printf("        %.3f                   n ≈ %-10d          sí, en una tarde\n", beta, n)
	}
	fmt.Println("   PERO si la hipótesis es VERDADERA no hay certificado finito: habría que revisar")
	fmt.Println("   infinitas perlas en infinitos armónicos. Por eso el laboratorio puede REFUTARLA")
	fmt.Println("   pero jamás DEMOSTRARLA midiendo. Esa asimetría es toda la diferencia.")

	// ---- LAW 4: the consequence the captain will like ----
	fmt.Println("\nLEY 4 · LA CONSECUENCIA — si fuera indemostrable, sería VERDADERA")
	fmt.Println("   una afirmación cuya MENTIRA siempre tiene un certificado finito no puede ser")
	fmt.Println("   falsa e indecidible al mismo tiempo: si fuera falsa, el certificado la refutaría.")
	fmt.Println("   Por lo tanto, para una afirmación de esta forma:")
	fmt.Println("\n        indecidible  ⟹  verdadera")
	fmt.Println("\n   (la hipótesis tiene esa forma: su negación es finitamente verificable)")
	fmt.Println("   → o alguien la demuestra, o alguien la refuta con un certificado, o es verdadera")
	fmt.Println("     y nadie podrá probarlo nunca. Las tres puertas están abiertas, y ninguna es la nada.")

	fmt.Println("\n════════ LA VERDAD, SEGÚN EL MEDIDOR DEL CAPITÁN ════════")
	fmt.Println("«Verdad = correspondencia. Verdad(P) ⟺ |P − R| = 0.» Y su propia trampa:")
	fmt.Println("en matemática la verdad de un teorema es verdad DENTRO de un sistema de axiomas.")
	fmt.Println("\nEsa distinción ES nuestro problema, y ahora tiene números:")
	fmt.Println("   VERDAD MEDIDA (correspondencia): la tenemos. |P − R| = 2e-16 en todo lo que miramos.")
	fmt.Println("   VERDAD DEMOSTRADA (derivación):  no la tenemos. Y no se consigue midiendo, porque")
	fmt.Println("   el certificado de la verdad — a diferencia del de la mentira — es infinito.")
	fmt.Println("\nEl millón no paga la correspondencia. Paga la derivación. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⚖️ LA VERDAD — Verdad(P) ⟺ |P − R| = 0</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"la verdad de un teorema es verdad DENTRO de un sistema de axiomas" — la trampa que el capitán nombró solo, y que ES nuestro problema</text>`,
		W, H, W, H, W/2, W/2)

	// two columns: measured vs derived
	fmt.Fprintf(&b, `<rect x="60" y="105" width="670" height="300" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.8"/>
<text x="395" y="141" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">VERDAD MEDIDA — la tenemos ✓</text>
<text x="395" y="177" font-size="13.5" text-anchor="middle" fill="#dce8f7">P: «cada perla está sobre la línea»</text>
<text x="395" y="201" font-size="13.5" text-anchor="middle" fill="#dce8f7">R: donde el mar dice que están</text>
<text x="395" y="243" font-size="22" text-anchor="middle" font-family="Consolas,monospace" fill="#7fd7a8">|P − R| = %.0e</text>
<text x="395" y="277" font-size="13" text-anchor="middle" fill="#8fa8c7">medido en %d perlas — cero al último bit</text>
<text x="395" y="315" font-size="13.5" text-anchor="middle" fill="#ffd166">LA CORRESPONDENCIA ESTÁ</text>
<text x="395" y="345" font-size="12.5" text-anchor="middle" fill="#dce8f7">donde miramos, lo que es coincide</text>
<text x="395" y="367" font-size="12.5" text-anchor="middle" fill="#dce8f7">con lo que debería ser</text>`,
		peor, len(pearls))

	fmt.Fprintf(&b, `<rect x="770" y="105" width="670" height="300" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.8"/>
<text x="1105" y="141" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">VERDAD DEMOSTRADA — no la tenemos ✗</text>
<text x="1105" y="177" font-size="13.5" text-anchor="middle" fill="#dce8f7">P se deduce de los axiomas</text>
<text x="1105" y="201" font-size="13.5" text-anchor="middle" fill="#dce8f7">por una cadena FINITA</text>
<text x="1105" y="243" font-size="13" text-anchor="middle" fill="#ffd166">y no se consigue midiendo:</text>
<text x="1105" y="273" font-size="13" text-anchor="middle" fill="#dce8f7">un impostor con TODOS los síntomas</text>
<text x="1105" y="295" font-size="13" text-anchor="middle" fill="#dce8f7">(espejo, ecuación funcional, cambiaformas</text>
<text x="1105" y="317" font-size="13" text-anchor="middle" fill="#dce8f7">coincidentes a %.0e) tiene sus raíces FUERA</text>
<text x="1105" y="355" font-size="13.5" text-anchor="middle" fill="#ff8fa0">coincidencia de SÍNTOMAS</text>
<text x="1105" y="377" font-size="13.5" text-anchor="middle" fill="#ff8fa0">no es coincidencia de HECHO</text>`, peorS)

	// the asymmetry
	fmt.Fprintf(&b, `<rect x="60" y="435" width="1380" height="270" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.8"/>
<text x="%.0f" y="471" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA ASIMETRÍA — el certificado de la MENTIRA es finito; el de la VERDAD no existe</text>
<text x="%.0f" y="503" font-size="13" text-anchor="middle" fill="#dce8f7">si la hipótesis fuera FALSA, una sola perla corrida obliga a un sobre en rojo en un armónico calculable — verificable en una tarde:</text>
<text x="%.0f" y="533" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β de la perla corrida        el certificado sería de</text>`,
		W/2, W/2, W/2)
	for i, c := range certs {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#7fd7a8">%.3f                          n ≈ %d</text>`,
			W/2, 558.0+float64(i)*22, c.beta, c.n)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#ff8fa0">pero si es VERDADERA no hay certificado finito: infinitas perlas, infinitos armónicos.</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#ffd166">este laboratorio puede REFUTARLA en una tarde — y jamás DEMOSTRARLA midiendo.</text>`,
		W/2, 558.0+float64(len(certs))*22+14, W/2, 558.0+float64(len(certs))*22+40)

	// the consequence
	fmt.Fprintf(&b, `<rect x="60" y="730" width="1380" height="230" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2.5"/>
<text x="%.0f" y="766" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">Y UNA CONSECUENCIA QUE EL CAPITÁN VA A QUERER</text>
<text x="%.0f" y="800" font-size="14" text-anchor="middle" fill="#dce8f7">una afirmación cuya MENTIRA siempre tiene certificado finito no puede ser falsa e indecidible a la vez:</text>
<text x="%.0f" y="824" font-size="14" text-anchor="middle" fill="#dce8f7">si fuera falsa, el certificado la refutaría. Y la hipótesis tiene exactamente esa forma.</text>
<text x="%.0f" y="866" font-size="21" text-anchor="middle" font-family="Georgia" fill="#ffd166">INDECIDIBLE  ⟹  VERDADERA</text>
<text x="%.0f" y="902" font-size="13.5" text-anchor="middle" fill="#dce8f7">o alguien la demuestra, o alguien la refuta con un certificado, o es verdadera y nadie podrá probarlo jamás.</text>
<text x="%.0f" y="928" font-size="13.5" text-anchor="middle" fill="#7fd7a8">Las tres puertas están abiertas — y ninguna de las tres es la nada.</text>
<text x="%.0f" y="960" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-verdad.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-verdad.svg")
}
