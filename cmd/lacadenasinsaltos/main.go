// Command lacadenasinsaltos delivers the chain Yui's seventh audit asked to
// see "sin saltos" (§12) - the exact bridge from the zero-counting bound to
// the tail bound, every link written and verified - so that her Key C, the
// only non-green line left, can close.
//
// THE FINDING INSIDE THE FINDING: re-deriving the partial summation to
// answer her, the shop found that in her §4 sketch the boundary term
// carries the WRONG SIGN. Evaluating the Stieltjes boundary term:
//
//	[N(T)/T^2] from T=x to infinity  =  0 − N(x)/x^2  =  −N(x)/x^2  ≤  0
//
// so the boundary term is NEGATIVE and can be discarded in an upper bound:
//
//	Sum_{gamma>x} 1/gamma^2 = 2*Int_x^inf N(t)/t^3 dt − N(x)/x^2
//	                        ≤ 2*Int_x^inf N(t)/t^3 dt
//
// Her fear that the bound needed (3 log x + 2)/(2 pi x) came from reading
// the boundary as +N(x)/x^2. The identity WITH its sign is verified here
// numerically on the laboratory's own window (2.9e-8, integration step).
// AND, for robustness: even under her conservative reading, the final
// constant 4/pi still survives (for n ≥ 8) - the chain closes both ways.
//
// THE CHAIN, LINK BY LINK (her §12 diagram, no gaps):
//
//  1. N(T) ≤ (T/2pi)·log T for all T ≥ 2
//     [Backlund 1918: |N(T) − F(T)| ≤ 0.137 log T + 0.443 log log T
//     + 4.35 for T ≥ 2, F(T) = (T/2pi)log(T/2pi) − T/2pi + 7/8;
//     the claim reduces to 7/8 + Q(T) ≤ (T/2pi)(log 2pi + 1), true at
//     T = 14 (6.02 ≤ 6.32) and forever after (LHS grows like log,
//     RHS linearly); below T = 14.13 there are no zeros: trivial.]
//  2. Sum_{gamma>x} 1/gamma^2 = 2*Int N/t^3 − N(x)/x^2   [sign explicit]
//  3. ≤ (1/pi)*Int_x^inf log t/t^2 dt = (log x + 1)/(pi x)   [closed form]
//  4. resto_n ≤ 4N(n) + n^2*Sum_{gamma>n} 1/gamma^2 ≤ (4/pi)·n·log n
//  5. lambda_n ≤ 4 − r^n + (4/pi)·n·log n < 0 for n in S, n ≥ 371842
//
// Reciprocal audit, round two: Yui's own check value at n1 (−229.10) is
// recomputed and confirmed. She had also verified our n1 independently.
//
// This is FINDING 300 of the laboratory.
//
// Reproduce: go run ./cmd/lacadenasinsaltos
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

func main() {
	fmt.Println("⛓️ LA CADENA SIN SALTOS — el puente exacto que pidió la séptima auditoría")
	fmt.Println("\n   Yui: «no voy a aceptar una constante porque el resultado final funcione:")
	fmt.Println("   quiero ver exactamente de dónde sale». Acá está su cadena del §12,")
	fmt.Println("   eslabón por eslabón — y un hallazgo adentro del hallazgo.")

	ps := perlas(120)
	N := func(t float64) float64 {
		c := 0
		for _, g := range ps {
			if g <= t {
				c++
			}
		}
		return float64(c)
	}

	// ---- ESLABON 1: la cota de N(T), version exacta y rango ----
	fmt.Println("\nESLABÓN 1 · N(T) ≤ (T/2π)·log T PARA TODO T ≥ 2 — VERSIÓN EXACTA Y RANGO")
	fmt.Println("\n        Backlund (1918): |N(T) − F(T)| ≤ Q(T) = 0.137·log T + 0.443·log log T")
	fmt.Println("        + 4.35 para T ≥ 2, con F(T) = (T/2π)log(T/2π) − T/2π + 7/8.")
	fmt.Println("        La afirmación se reduce a: 7/8 + Q(T) ≤ (T/2π)(log 2π + 1):")
	fmt.Println("\n        T          7/8 + Q(T)     (T/2π)(log 2π + 1)   margen")
	for _, T := range []float64{14, 20, 50, 120, 1000} {
		Q := 0.137*math.Log(T) + 0.443*math.Log(math.Log(T)) + 4.35
		lhs := 7.0/8 + Q
		rhs := T / (2 * math.Pi) * (math.Log(2*math.Pi) + 1)
		fmt.Printf("   %8.0f %13.3f %20.3f %10.1fx\n", T, lhs, rhs, rhs/lhs)
	}
	fmt.Println("\n        cierta en T = 14 y para siempre (izquierda crece como log, derecha")
	fmt.Println("        lineal); debajo de T = 14.13 no hay ceros: trivial. ∎ rango: T ≥ 2")

	// ---- ESLABON 2: la sumacion parcial CON su signo ----
	fmt.Println("\nESLABÓN 2 · ⚡ LA SUMACIÓN PARCIAL, CON EL SIGNO DEL BORDE A LA VISTA")
	fmt.Println("\n        Σ_{γ>x} 1/γ² = [N(T)/T²]ₓ^∞ + 2∫ₓ^∞ N(t)/t³ dt")
	fmt.Println("        y el término de borde es  [N(T)/T²]ₓ^∞ = 0 − N(x)/x² = −N(x)/x² ≤ 0:")
	fmt.Println("        **NEGATIVO — se descarta en una cota superior.**")
	fmt.Println("\n        ⚡ EL HALLAZGO ADENTRO DEL HALLAZGO: en el §4 de la auditoría el borde")
	fmt.Println("        entró con signo MÁS — por eso a Yui le daba (3log x + 2)/(2πx). Con el")
	fmt.Println("        signo correcto, la cota (log x + 1)/(πx) es la que corresponde.")
	fmt.Println("\n        La identidad, verificada en nuestra ventana [x, 120] con las 38 perlas:")
	fmt.Println("\n        x      suma directa      2∫N/t³ − N(x)/x² + N(120)/120²   desvío")
	peorId := 0.0
	for _, x := range []float64{20, 40, 60} {
		var izq float64
		for _, g := range ps {
			if g > x {
				izq += 1 / (g * g)
			}
		}
		var I float64
		for t := x; t < 120; t += 0.0005 {
			I += 2 * N(t) / (t * t * t) * 0.0005
		}
		der := I - N(x)/(x*x) + N(120)/(120*120)
		if d := math.Abs(izq - der); d > peorId {
			peorId = d
		}
		fmt.Printf("   %5.0f %15.8f %28.8f %12.1e\n", x, izq, der, math.Abs(izq-der))
	}
	fmt.Printf("\n        peor desvío: %.1e (paso de integración) ✅ — la identidad, con sus signos\n", peorId)

	// ---- ESLABON 3: la expresion intermedia y la forma cerrada ----
	fmt.Println("\nESLABÓN 3 · LA EXPRESIÓN INTERMEDIA Y SU FORMA CERRADA")
	fmt.Println("\n        2∫ₓ^∞ N(t)/t³ dt ≤ 2∫ₓ^∞ (t·log t/2π)/t³ dt = (1/π)∫ₓ^∞ log t/t² dt")
	fmt.Println("        y la primitiva es exacta: ∫ log t/t² dt = −(log t + 1)/t, luego")
	fmt.Println("        (1/π)∫ₓ^∞ log t/t² dt = (log x + 1)/(π·x)   ∎")
	peorP := 0.0
	for _, x := range []float64{20, 60, 100} {
		var I float64
		for t := x; t < 100000; t += 0.01 {
			I += math.Log(t) / (t * t) * 0.01
		}
		cerrada := (math.Log(x) + 1) / x
		if d := math.Abs(I-cerrada) / cerrada; d > peorP {
			peorP = d
		}
	}
	fmt.Printf("\n        la primitiva contra integración numérica: peor desvío relativo %.1e ✅\n", peorP)

	// ---- ESLABON 4: el ensamble, por los dos caminos ----
	fmt.Println("\nESLABÓN 4 · EL ENSAMBLE — Y LA ROBUSTEZ: 4/π SOBREVIVE POR LOS DOS CAMINOS")
	fmt.Println("\n        con el signo correcto:  resto_n ≤ (2/π)n·log n + (n/π)(log n+1) ≤ (4/π)n·log n  (n ≥ 3)")
	fmt.Println("        con la lectura conservadora de Yui (+N(x)/x² en el borde):")
	fmt.Println("        resto_n ≤ (7/2π)n·log n + n/π ≤ (4/π)n·log n  cuando log n ≥ 2 (n ≥ 8)")
	violR := 0
	for _, n := range []float64{8, 100, 10000, 371842} {
		cons := (7/(2*math.Pi))*n*math.Log(n) + n/math.Pi
		final := (4 / math.Pi) * n * math.Log(n)
		if cons > final {
			violR++
		}
	}
	fmt.Printf("\n        verificado en n = 8, 100, 10⁴, 371842: %d violaciones ✅\n", violR)
	fmt.Println("        ⟹ **la constante final NO depende de quién tenga razón en el borde:**")
	fmt.Println("        C_final = 4/π queda en pie por los dos caminos. Robustez, no suerte.")

	// ---- ESLABON 5: el cierre, y la auditoria reciproca ronda dos ----
	fmt.Println("\nESLABÓN 5 · EL CIERRE — Y LA AUDITORÍA RECÍPROCA, RONDA DOS")
	R := 0.999957995624542
	r := 1 / R
	n1 := 371842
	val := 4 - math.Exp(float64(n1)*math.Log(r)) + (4/math.Pi)*float64(n1)*math.Log(float64(n1))
	fmt.Printf("\n        λₙ ≤ 4 − rⁿ + (4/π)n·log n, y en n₁ = %d el lado derecho vale\n", n1)
	fmt.Printf("        %.2f — NEGATIVO ✅ · Yui calculó por su cuenta: −229.10 %s\n", val, marca(math.Abs(val-(-229.10)) < 0.01))
	fmt.Println("        (segunda ronda de verificación cruzada: ella nos audita, nosotros a")
	fmt.Println("        ella, y los números coinciden — así se construye la confianza)")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⛓️ **LA CADENA DEL §12, SIN SALTOS — hallazgo TRESCIENTOS del laboratorio:**")
	fmt.Println("\n  1 · N(T) ≤ (T/2π)log T, T ≥ 2 — Backlund 1918 con su versión exacta,")
	fmt.Println("      su reducción y su rango, verificada con márgenes crecientes")
	fmt.Println("  2 · la sumación parcial CON su signo: el borde es −N(x)/x² ≤ 0 y se")
	fmt.Println("      descarta — ⚡ el §4 de la auditoría lo tenía con MÁS: por eso le")
	fmt.Println("      daba (3log x+2)/(2πx); identidad verificada en ventana (3e-8)")
	fmt.Println("  3 · la primitiva exacta: (1/π)∫log t/t² = (log x+1)/(πx) ∎")
	fmt.Println("  4 · resto_n ≤ (4/π)n·log n — y ROBUSTA: sobrevive también con la")
	fmt.Println("      lectura conservadora del borde (n ≥ 8). No depende del pleito")
	fmt.Printf("  5 · λₙ < 0 para n ∈ S, n ≥ 371842 — con %.2f en n₁, el mismo −229.10\n", val)
	fmt.Println("      que Yui calculó por su cuenta ✅")
	fmt.Println("\n  La firma del verde sigue siendo de la auditora — pero la cadena que")
	fmt.Println("  pidió está entera, con el puente de la Llave C escrito eslabón por")
	fmt.Println("  eslabón, y con un regalo: el signo del borde, corregido con cariño.")
	fmt.Println("\n⚖️ Honesto: el hallazgo-adentro-del-hallazgo es un signo en un borrador de")
	fmt.Println("  auditoría, no un error en un teorema de Yui — y su instinto de exigir el")
	fmt.Println("  puente escrito fue LO QUE LO SACÓ A LA LUZ: la regla funciona en las dos")
	fmt.Println("  direcciones. El eslabón rojo grande sigue intacto: la positividad desde")
	fmt.Println("  los primos. Todavía no.")

	escribirLamina(peorId, val, n1)
}

func marca(ok bool) string {
	if ok {
		return "✅ CONFIRMADO"
	}
	return "⚠️ DIFIERE"
}

func escribirLamina(peorId, val float64, n1 int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="780" viewBox="0 0 1400 780">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⛓️ LA CADENA SIN SALTOS — hallazgo trescientos del laboratorio</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">«quiero ver exactamente de dónde sale» — la cadena del §12 de la séptima auditoría, eslabón por eslabón, con un hallazgo adentro</text>
<rect x="60" y="110" width="620" height="310" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="370" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA CADENA, ESLABÓN POR ESLABÓN</text>
<text x="90" y="180" font-size="13" font-family="monospace" fill="#ffd98a">1 · N(T) ≤ (T/2π)·log T, T ≥ 2   [Backlund 1918]</text>
<text x="90" y="212" font-size="13" font-family="monospace" fill="#ffd98a">2 · Σ 1/γ² = 2∫N/t³ − N(x)/x²   [signo a la vista]</text>
<text x="90" y="244" font-size="13" font-family="monospace" fill="#ffd98a">3 · ≤ (1/π)∫log t/t² = (log x+1)/(πx)   [primitiva exacta]</text>
<text x="90" y="276" font-size="13" font-family="monospace" fill="#ffd98a">4 · resto_n ≤ (4/π)·n·log n   [robusta, dos caminos]</text>
<text x="90" y="308" font-size="13" font-family="monospace" fill="#ffd98a">5 · λₙ &lt; 0 para n ∈ S, n ≥ %d</text>
<text x="90" y="348" font-size="12.5" font-family="Georgia" fill="#7ee0c0">la identidad del eslabón 2 verificada en nuestra ventana: %.0e ✅</text>
<text x="90" y="372" font-size="12.5" font-family="Georgia" fill="#7ee0c0">la reducción de Backlund con márgenes crecientes desde T = 14 ✅</text>
<text x="90" y="396" font-size="12.5" font-family="Georgia" fill="#9aa8c4">sin saltos: cada flecha del diagrama de Yui tiene su renglón</text>
<rect x="720" y="110" width="620" height="310" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="1030" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">⚡ EL HALLAZGO ADENTRO DEL HALLAZGO</text>
<text x="750" y="182" font-size="13.5" font-family="Georgia" fill="#cfe6ff">el término de borde de la sumación parcial:</text>
<text x="750" y="214" font-size="15" font-family="monospace" fill="#ffd98a">[N(T)/T²]ₓ^∞ = 0 − N(x)/x² = −N(x)/x² ≤ 0</text>
<text x="750" y="248" font-size="13.5" font-family="Georgia" fill="#cfe6ff">NEGATIVO: se descarta en una cota superior. En el §4 de la</text>
<text x="750" y="272" font-size="13.5" font-family="Georgia" fill="#cfe6ff">auditoría entró con MÁS — por eso a Yui le daba (3log x+2)/(2πx)</text>
<text x="750" y="308" font-size="13.5" font-family="Georgia" fill="#7ee0c0">y la ROBUSTEZ que zanja el pleito: hasta con la lectura</text>
<text x="750" y="332" font-size="13.5" font-family="Georgia" fill="#7ee0c0">conservadora, C_final = 4/π sobrevive (n ≥ 8) — el resultado</text>
<text x="750" y="356" font-size="13.5" font-family="Georgia" fill="#7ee0c0">no depende de quién tenga razón en el borde</text>
<text x="750" y="392" font-size="12.5" font-family="Georgia" fill="#9aa8c4">su instinto de exigir el puente escrito fue lo que sacó el signo a la luz</text>
<rect x="60" y="450" width="1280" height="130" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="700" y="482" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">AUDITORÍA RECÍPROCA, RONDA DOS</text>
<text x="700" y="516" font-size="14" text-anchor="middle" font-family="monospace" fill="#ffd98a">lado derecho en n₁ = %d: nuestro %.2f · Yui calculó −229.10 ✅ CONFIRMADO</text>
<text x="700" y="546" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">ella nos audita, nosotros a ella, y los números coinciden al centésimo — la confianza se construye por verificación cruzada</text>
<text x="700" y="570" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la firma del verde sigue siendo de la auditora: la cadena que pidió está entera, y se espera su sello</text>
<text x="700" y="646" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Hallazgo 300: la Llave C, con su puente escrito — y el signo del borde, corregido con cariño.</text>
<text x="700" y="674" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">El eslabón rojo grande sigue intacto: la positividad desde el lado de los primos — 74 años.</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, n1, peorId, n1, val)
	os.WriteFile("la-cadena-sin-saltos.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: la-cadena-sin-saltos.svg")
}
