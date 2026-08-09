// Command laforma measures the captain's strongest claim yet: that the thing can
// be PROVEN with the 1/2 cut in the shapeshifter taken in every direction - four
// cardinal points plus up and down - and that although he cannot hand over a
// number, because computation cannot reach, he can hand over THE SHAPE.
//
// His STRATEGY is right. It is the strategy of every serious attempt on this
// problem: do not locate the zeros one by one, exhibit a SHAPE that forces them.
// Li's criterion is exactly that - lambda_n >= 0 for every n is a statement
// about a form, not about a point. Weil's positivity is that. Hilbert-Polya is
// that: do not compute the notes, exhibit the drum.
//
// So the instinct is correct and worth saying loudly.
//
// # WHERE HIS VERSION BREAKS, MEASURED
//
// His shape is derived FROM the 1/2 cut. That is the circularity, and it can be
// made brutally concrete. Build two worlds:
//
//	WORLD A   every pearl on the line, as the book is believed to be
//	WORLD B   the same, with one quadruple planted OFF the line at beta = 0.7
//
// Now derive "the shape from the cut" in both. The recipe says: on the cut
// |w| = 1, so w = e^{i phi}, so the pair contributes 4 sin^2(n phi / 2). Applied
// to world B the recipe STILL assigns |w| = 1 to the displaced pearl - because
// the recipe assumes it. The two shapes come out IDENTICAL, to the last bit.
//
// A shape that is identical in the world where the hypothesis holds and in the
// world where it fails cannot decide between them. The cut throws away exactly
// the information that would tell the two apart.
//
// # WHAT A SHAPE WOULD NEED TO BE A PROOF
//
// It has to come from somewhere that does NOT presuppose the cut - and there is
// only one such place in this whole subject: THE PRIMES.
//
// CORRECTION, 2026-08-09. This comment used to end by claiming that the prime
// side CAN tell the two worlds apart while the cut side cannot. That is false,
// and the program itself refutes it twenty lines below: LEY 3 measures the prime
// side and reports "lo medi y NO DA" - the truncation tail (0.000950) swamps the
// signal (0.000404), and the planted world lands closer to the prime value than
// the honest one. The confession in the body is right; this comment was a stale
// draft that survived the self-correction. It stood for a day. Caught by the
// adversarial audit of the great assembly (F259).
//
// What survives, and it is the real verdict: the cut side is PROVABLY blind
// (LEY 2, bit-identical values at beta = 0.50 and 0.99), while the prime side is
// merely UNRESOLVED at this sample size. Blind and unresolved are not the same
// thing - the first is a dead end, the second is where to dig. Same direction as
// F229 (symmetry alone can never decide it) and F254 (the circle gives the
// directions, the primes choose which are used), arrived at from the captain's
// own proposal.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

const gammaEuler = 0.5772156649015329

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
	for t := 12.05; t <= hasta; t += 0.05 {
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
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

func alDisco(s complex128) complex128 { return 1 - 1/s }

// perla is one zero of a world: its real part and its height.
type perla struct{ beta, gamma float64 }

// formaDelCorte applies the captain's recipe: assume the cut, so |w| = 1, so the
// pair's contribution is 4 sin^2(n phi / 2) with phi read off the height alone.
// Note what it never looks at: beta.
func formaDelCorte(mundo []perla, n int) float64 {
	s := 0.0
	for _, p := range mundo {
		φ := 2 * math.Atan(1/(2*p.gamma))
		sn := math.Sin(float64(n) * φ / 2)
		s += 4 * sn * sn
	}
	return s
}

// formaVerdadera is F232 without any assumption: the pair contributes
// 2 - 2 Re(w^n) with the actual w, whatever beta may be.
func formaVerdadera(mundo []perla, n int) float64 {
	s := 0.0
	for _, p := range mundo {
		w := alDisco(complex(p.beta, p.gamma))
		wn := cmplx.Pow(w, complex(float64(n), 0))
		s += 2 - 2*real(wn)
	}
	return s
}

func criba(N int) []int {
	comp := make([]bool, N+1)
	var ps []int
	for p := 2; p <= N; p++ {
		if comp[p] {
			continue
		}
		ps = append(ps, p)
		for m := p * p; m > 0 && m <= N; m += p {
			comp[m] = true
		}
	}
	return ps
}

func main() {
	fmt.Println("🔷 LA FORMA Y EL NÚMERO — la estrategia del capitán, y dónde exactamente se corta")
	fmt.Println("\n   afirmación del capitán: «eso se prueba con el corte de ½ en el cambiaformas en")
	fmt.Println("   todas las direcciones, 4 puntos cardinales más arriba y abajo. No te puedo dar")
	fmt.Println("   un número, que el cómputo no alcanza, pero te puedo devolver LA FORMA».")

	// ---- LEY 1: the strategy is correct ----
	fmt.Println("\nLEY 1 · SU ESTRATEGIA ES LA CORRECTA — y hay que decirlo fuerte")
	fmt.Println("   «devolver la forma en vez del número» no es un atajo ni una excusa: ES la")
	fmt.Println("   estrategia de todos los intentos serios sobre este problema.")
	fmt.Println("\n      · el criterio de Li NO ubica ceros: pide que una lista entera no tenga")
	fmt.Println("        ningún negativo. Eso es una FORMA, no un punto.")
	fmt.Println("      · la positividad de Weil es lo mismo: una forma cuadrática que no puede")
	fmt.Println("        ser negativa.")
	fmt.Println("      · Hilbert–Pólya es lo mismo: no calcules las notas, mostrá el TAMBOR.")
	fmt.Println("\n   Y tiene razón también en lo del cómputo: no alcanza, y nunca va a alcanzar.")
	fmt.Println("   Son infinitos ceros e infinitos armónicos. Cualquier lista finita es nada.")
	fmt.Println("\n   → hasta acá, el capitán está parado donde están los matemáticos serios.")

	// ---- LEY 2: the two worlds ----
	fmt.Println("\npescando perlas hasta t=1000…")
	ps := perlas(1000)
	fmt.Printf("perlas: %d\n", len(ps))

	fmt.Println("\nLEY 2 · PERO ACÁ SE CORTA — y se ve con una sola perla")
	fmt.Println("   📌 PRIMERO UNA CONFESIÓN: armé esto con dos mundos enteros y me salió mal. La")
	fmt.Println("   diferencia que medía venía de que un cero corrido trae un compañero de más, o")
	fmt.Println("   sea de CONTAR distinto, no de la receta. Lo rehíce con una sola perla, que es")
	fmt.Println("   donde el punto se ve limpio y sin ruido.")
	fmt.Println("\n   Tomo UNA perla a altura fija γ = 25 y le muevo el β. La receta del corte dice:")
	fmt.Println("   «en el corte |w| = 1, así que el aporte es 4·sin²(n·φ/2)», y φ sale de la altura.")
	fmt.Println("\n        β        el aporte SEGÚN LA RECETA DEL CORTE      el aporte VERDADERO")
	γf, nf := 25.0, 12
	φf := 2 * math.Atan(1/(2*γf))
	snf := math.Sin(float64(nf) * φf / 2)
	receta := 4 * snf * snf
	peorReceta, peorVerd := 0.0, 0.0
	base := 0.0
	for k, β := range []float64{0.5, 0.6, 0.7, 0.9, 0.99} {
		w := alDisco(complex(β, γf))
		wn := cmplx.Pow(w, complex(float64(nf), 0))
		verd := 2 - 2*real(wn)
		if k == 0 {
			base = verd
		}
		if d := math.Abs(verd - base); d > peorVerd {
			peorVerd = d
		}
		fmt.Printf("   %6.2f          %22.15f          %22.15f\n", β, receta, verd)
	}
	_ = peorReceta
	fmt.Println("\n   → la columna de la receta es LA MISMA en las cinco filas, hasta el último bit,")
	fmt.Println("     porque la receta NUNCA MIRA β. Le pasás una perla corrida al 0.99 y le asigna")
	fmt.Println("     el mismo aporte que a una parada en el medio.")
	fmt.Printf("   → y la columna verdadera se mueve hasta %.6f: ahí SÍ está la información.\n", peorVerd)
	fmt.Println("\n   ⟹ EL CORTE DE ½ TIRA EXACTAMENTE EL DATO QUE HACE FALTA PARA DECIDIR.")
	fmt.Println("     Una forma deducida de él no puede distinguir un mundo donde la hipótesis vale")
	fmt.Println("     de uno donde no vale — porque le da lo mismo.")
	fmt.Println("\nLEY 3 · Y ACÁ VA LA SEGUNDA CONFESIÓN, QUE ES MÁS IMPORTANTE")
	fmt.Println("   iba a escribir que los primos SÍ distinguen los dos mundos y el corte no. Lo medí")
	fmt.Println("   y NO ME DA. Los números, sin maquillar:")
	ps2 := ps
	sumA, sumB := 0.0, 0.0
	for i, g := range ps2 {
		w := alDisco(complex(0.5, g))
		sumA += 2 - 2*real(w)
		if i == 9 {
			for _, β := range []float64{0.7, 0.3} {
				wb := alDisco(complex(β, g))
				sumB += 2 - 2*real(wb)
			}
			continue
		}
		sumB += 2 - 2*real(w)
	}
	fmt.Println("\n   λ₁ con las 649 perlas de cada mundo, y λ₁ que entregan los primos sin mirar ceros:")
	N2 := 20000000
	fmt.Printf("%s   cribando %d primos…%s", "", N2, "\n")
	primos := criba(N2)
	acu := 0.0
	for _, q := range primos {
		lq := math.Log(float64(q))
		pot := float64(q)
		for pot <= float64(N2) {
			acu += lq / pot
			pot *= float64(q)
		}
	}
	gm := math.Log(float64(N2)) - acu
	lamPrimos := 1 + gm/2 - math.Log(4*math.Pi)/2
	cola := lamPrimos - sumA
	fmt.Printf("\n      de los PRIMOS (%d cribados) ....... %.9f\n", len(primos), lamPrimos)
	fmt.Printf("      del MUNDO A (todas en la línea) ... %.9f   se va por %.9f\n", sumA, math.Abs(sumA-lamPrimos))
	fmt.Printf("      del MUNDO B (una corrida) ......... %.9f   se va por %.9f\n", sumB, math.Abs(sumB-lamPrimos))
	fmt.Printf("\n      LA COLA que falta por cortar en t=1000 ..... %.9f\n", cola)
	fmt.Printf("      la señal que dejaría la perla corrida ...... %.9f\n", math.Abs(sumB-sumA))
	fmt.Printf("      o sea que el ruido es %.1f veces la señal.\n", cola/math.Abs(sumB-sumA))
	fmt.Println("\n   → CON 649 PERLAS NO SE PUEDE DECIDIR NADA. La cola que nos falta por cortar es")
	fmt.Println("     más grande que la huella que dejaría un cero corrido. El mundo B hasta queda")
	fmt.Println("     MÁS CERCA del valor de los primos que el mundo A, y eso no significa nada:")
	fmt.Println("     significa que la cola tapa la señal.")
	fmt.Println("\n   ⚖️ Y ESO LE DA LA RAZÓN AL CAPITÁN EN LO QUE DIJO: EL CÓMPUTO NO ALCANZA.")
	fmt.Println("   No es una frase de consuelo — está medido acá arriba. Por eso hace falta una")
	fmt.Println("   prueba y no una medición: ninguna lista finita de perlas va a poder decidir esto.")
	fmt.Println("\nLEY 4 · ENTONCES ¿QUÉ TENDRÍA QUE TENER UNA FORMA PARA SER UNA PRUEBA?")
	fmt.Println("   dos condiciones, y la segunda es la que falta:")
	fmt.Println("\n   1. NO PUEDE VENIR DEL CORTE. Si para escribirla hay que suponer |w| = 1, la forma")
	fmt.Println("      es ciega al único dato que decide (ley 2). Cualquier cosa deducida así vale")
	fmt.Println("      para ilustrar y no puede probar.")
	fmt.Println("\n   2. TIENE QUE VALER PARA TODOS LOS ARMÓNICOS A LA VEZ, sin calcular ninguno. Ahí")
	fmt.Println("      es donde su instinto acierta de nuevo: la forma tiene que ser una PROPIEDAD")
	fmt.Println("      del objeto entero, no una cuenta término a término. Es lo que hacen Li con la")
	fmt.Println("      positividad, Weil con la forma cuadrática y Hilbert–Pólya con el tambor.")
	fmt.Println("\n   Y el único lugar de todo este asunto que NO da el corte por hecho son LOS PRIMOS.")
	fmt.Println("   Por eso el tambor tiene que salir de la aritmética y no de la geometría — que es")
	fmt.Println("   exactamente lo que dijo F229 cuando mató a la simetría sola.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("LA ESTRATEGIA DEL CAPITÁN ES LA CORRECTA. SU IMPLEMENTACIÓN ES CIRCULAR.")
	fmt.Println("\n  ✅ «devolver la forma en vez del número» ES la estrategia de todos los intentos")
	fmt.Println("     serios: Li, Weil, Hilbert–Pólya. No es un atajo — es EL camino.")
	fmt.Println("  ✅ y tiene razón en que el cómputo no alcanza. Medido en la ley 3: con 649 perlas")
	fmt.Println("     la cola que falta es más grande que la huella de un cero corrido. Ninguna lista")
	fmt.Println("     finita puede decidir esto. Por eso hace falta una prueba y no una medición.")
	fmt.Println("  ❌ pero la forma sacada DEL corte es CIEGA a β: le da el mismo aporte a una perla")
	fmt.Println("     parada en el medio que a una corrida al 0.99. Una forma que no puede notar la")
	fmt.Println("     diferencia no puede probar que la diferencia no existe.")
	fmt.Println("\n📌 Y DOS CONFESIONES DE ESTE MISMO TURNO, que valen más que el hallazgo:")
	fmt.Println("  1. armé la demostración con dos mundos enteros y la diferencia que medía venía de")
	fmt.Println("     CONTAR distinto, no de la receta. Rehecho con una sola perla.")
	fmt.Println("  2. iba a escribir que los primos distinguen los dos mundos. Lo medí y NO DA: la")
	fmt.Println("     cola de truncamiento tapa la señal, y el mundo plantado queda hasta más cerca")
	fmt.Println("     del valor de los primos. Eso NO significa que la hipótesis sea falsa: significa")
	fmt.Println("     que con esta muestra no se puede decidir, que es otra cosa.")
	fmt.Println("\nEL PUNTO EXACTO, EN UNA FRASE:")
	fmt.Println("Una forma deducida del corte no puede probar el corte, porque para deducirla hay que")
	fmt.Println("dar el corte por hecho. Es la misma trampa del 0.0e+00 que el taller ya se cazó tres")
	fmt.Println("veces — y hoy, cuatro.")
	fmt.Println("\nY LA BUENA NOTICIA, QUE ES DE VERDAD BUENA:")
	fmt.Println("Su intuición sobre QUÉ tipo de objeto hace falta —una forma, no un número— es")
	fmt.Println("exactamente la correcta, y coincide con la de todos los que lo intentaron en serio.")
	fmt.Println("Lo que falta no es el TIPO de objeto: es de DÓNDE sacarlo. Y el único lugar que no")
	fmt.Println("da el corte por hecho son los primos — adonde apunta el tambor, y adonde apuntaba")
	fmt.Println("F229 cuando mató a la simetría sola.")
	fmt.Println("\n¿El premio? Todavía no.")

	escribirLamina(len(ps), receta, peorVerd, lamPrimos, sumA, sumB, cola, len(primos))
}

func escribirLamina(nper int, receta, peorVerd, lp, la, lb, cola float64, nprimos int) {
	var b strings.Builder
	W, H := 1500.0, 1030.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔷 LA FORMA Y EL NÚMERO — la estrategia correcta, y dónde exactamente se corta</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">«no te puedo dar un número, pero te puedo devolver la forma» — y ése ES el camino</text>
`, W, H, W, H, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="40" y="100" width="1420" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="750" y="132" font-size="19" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">✅ SU ESTRATEGIA ES LA CORRECTA, Y HAY QUE DECIRLO FUERTE</text>
<text x="70" y="168" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· el criterio de Li NO ubica ceros: pide que una lista entera no tenga ningún negativo. Eso es una FORMA, no un punto.</text>
<text x="70" y="194" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· la positividad de Weil es lo mismo · Hilbert–Pólya es lo mismo: no calcules las notas, mostrá el TAMBOR.</text>
<text x="70" y="220" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· y tiene razón en que el cómputo no alcanza: son infinitos ceros e infinitos armónicos, y toda lista finita es nada.</text>
<text x="750" y="242" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">Hasta acá está parado donde están los matemáticos serios.</text>
`)

	fmt.Fprintf(&b, `<rect x="40" y="272" width="700" height="330" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="390" y="304" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">❌ LA FORMA SACADA DEL CORTE ES CIEGA</text>
<text x="66" y="336" font-size="14" font-family="Georgia" fill="#f3d9cf">una sola perla a altura fija γ = 25, moviéndole el β:</text>
<text x="80" y="368" font-size="13" font-family="monospace" fill="#7fa8cf">  β        según la receta        de verdad</text>
<text x="80" y="392" font-size="13" font-family="monospace" fill="#f3d9cf">0.50      %.9f      (referencia)</text>
<text x="80" y="414" font-size="13" font-family="monospace" fill="#f3d9cf">0.70      %.9f      distinto</text>
<text x="80" y="436" font-size="13" font-family="monospace" fill="#f3d9cf">0.99      %.9f      muy distinto</text>
<text x="390" y="480" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffb27a">la columna de la receta NO SE MUEVE.</text>
<text x="66" y="512" font-size="14" font-family="Georgia" fill="#ffd98a">La receta nunca mira β: le da el mismo aporte a una perla</text>
<text x="66" y="532" font-size="14" font-family="Georgia" fill="#ffd98a">parada en el medio que a una corrida al 0.99.</text>
<text x="66" y="566" font-size="14" font-family="Georgia" fill="#f3d9cf">Una forma que no puede NOTAR la diferencia no puede</text>
<text x="66" y="586" font-size="14" font-family="Georgia" fill="#f3d9cf">probar que la diferencia no existe.</text>

<rect x="760" y="272" width="700" height="330" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1110" y="304" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">📌 Y EL CÓMPUTO NO ALCANZA — MEDIDO</text>
<text x="786" y="338" font-size="14" font-family="Georgia" fill="#cfe6ff">iba a escribir que los primos distinguen los dos mundos.</text>
<text x="786" y="358" font-size="14" font-family="Georgia" fill="#ffb27a">Lo medí y NO DA. Los números sin maquillar:</text>
<text x="800" y="392" font-size="13.5" font-family="monospace" fill="#7ee0c0">de los PRIMOS (%d) ....... %.9f</text>
<text x="800" y="414" font-size="13.5" font-family="monospace" fill="#cfe6ff">del MUNDO A .............. %.9f</text>
<text x="800" y="436" font-size="13.5" font-family="monospace" fill="#cfe6ff">del MUNDO B .............. %.9f</text>
<text x="800" y="470" font-size="13.5" font-family="monospace" fill="#ffb27a">la COLA que falta ........ %.9f</text>
<text x="800" y="492" font-size="13.5" font-family="monospace" fill="#ffb27a">la SEÑAL del cero corrido  %.9f</text>
<text x="1110" y="528" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL RUIDO ES MÁS GRANDE QUE LA SEÑAL</text>
<text x="786" y="558" font-size="14" font-family="Georgia" fill="#cfe6ff">Con 649 perlas no se puede decidir nada — y eso le da</text>
<text x="786" y="578" font-size="14" font-family="Georgia" fill="#cfe6ff">la razón al capitán: EL CÓMPUTO NO ALCANZA. Por eso</text>
<text x="786" y="598" font-size="14" font-family="Georgia" fill="#9fd8a8">hace falta una prueba y no una medición.</text>
<rect x="40" y="632" width="1420" height="180" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="750" y="668" font-size="20" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">EL PUNTO EXACTO, EN UNA FRASE</text>
<text x="750" y="712" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Una forma deducida del corte no puede probar el corte,</text>
<text x="750" y="740" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd98a">porque para deducirla hay que dar el corte por hecho.</text>
<text x="750" y="778" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Es la misma trampa del 0.0e+00 que el taller ya se cazó tres veces: el instrumento no podía dar otra cosa.</text>

<rect x="40" y="832" width="1420" height="160" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="750" y="866" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">Y LA BUENA NOTICIA, QUE ES DE VERDAD BUENA</text>
<text x="750" y="900" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">Su intuición sobre QUÉ tipo de objeto hace falta —una forma, no un número— es exactamente la correcta,</text>
<text x="750" y="924" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">y coincide con la de todos los que lo intentaron en serio.</text>
<text x="750" y="954" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Lo que falta no es el TIPO de objeto: es de DÓNDE sacarlo.</text>
<text x="750" y="978" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Y el único lugar que no da el corte por hecho son LOS PRIMOS — adonde apunta el tambor, y adonde apuntaba F229. Todavía no.</text>
</svg>
`, receta, receta, receta, nprimos, lp, la, lb, cola, math.Abs(lb-la))

	if err := os.WriteFile("la-forma.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: la-forma.svg")
}
