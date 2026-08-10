// Command elanterior judges the captain's tightening of his own formula.
//
// After the panel showed that "his formula IS Goldbach" claimed more than the
// algebra licenses - with X a whole number, 4X reaches only multiples of 4 - he
// did not defend the wide claim. He narrowed it:
//
//	"For the formula to CLOSE, X must be PRIME, and Y is not 2Y, it is Y,
//	 and it must be the PREVIOUS PRIME."
//
// That is a different statement, and a much better one: it is falsifiable in a
// single line. It is no longer about every even number. It says that a prime
// GENERATES two more primes from the prime that came before it:
//
//	X prime, Y = the prime immediately before X
//	P = 2X - Y   and   Q = 2X + Y   are both prime
//
// PRE-REGISTERED PREDICTION, written before the sweep was run: the strict rule
// will fail, and it will fail early - his own third row (X = 3, Y = 2) already
// gives 4 and 8. What is genuinely open, and what this program actually
// measures, is whether the WEAKER form survives: is the smallest offset that
// works usually a prime? If it is, his instinct has real content even though
// the strict rule does not.
//
// Reproduce: go run ./cmd/elanterior
package main

import (
	"fmt"
	"os"
	"strings"
)

const N = 20000000 // el índice más grande que se lee es 2X+Y < 3X ≤ 3·10^7... acotado abajo

func criba(n int) []bool {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if es[i] {
			for j := i * i; j <= n; j += i {
				es[j] = false
			}
		}
	}
	return es
}

func main() {
	fmt.Println("🔙 EL ANTERIOR — el capitán le pone dos candados a su propia fórmula")
	fmt.Println("\n   Dijo: «para que la fórmula CIERRE, X debe ser primo, e Y no es 2Y sino Y,")
	fmt.Println("   y debe ser el PRIMO ANTERIOR».")
	fmt.Println("\n   Eso ya NO es Goldbach: es un generador. Un primo, con el primo que lo")
	fmt.Println("   precede, tendría que fabricar otros dos primos. Y eso se mide en un renglón.")

	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)

	// lista de primos hasta la mitad, para que 2X+Y siempre entre en la criba
	tope := N / 3
	var primos []int
	for i := 2; i <= tope; i++ {
		if es[i] {
			primos = append(primos, i)
		}
	}
	fmt.Printf("primos disponibles como X: %d (hasta %d)\n", len(primos), tope)

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · SU PROPIA TABLA, LEÍDA CON LA REGLA NUEVA")
	fmt.Println("   Si Y tiene que ser el primo anterior a X, sus tres renglones dicen:")
	fmt.Println("\n        X    primo anterior   2X − Y    2X + Y    ¿cierra?")
	type fila struct {
		x, yUsada int
	}
	for _, f := range []fila{{2, 1}, {3, 1}, {5, 3}} {
		ant := 0
		for _, p := range primos {
			if p < f.x {
				ant = p
			} else {
				break
			}
		}
		antTxt := fmt.Sprint(ant)
		if ant == 0 {
			antTxt = "no hay"
		}
		veredicto, p, q := "—", 0, 0
		if ant > 0 {
			p, q = 2*f.x-ant, 2*f.x+ant
			if es[p] && es[q] {
				veredicto = "✅ SÍ"
			} else {
				veredicto = "❌ NO"
			}
		}
		if ant > 0 {
			fmt.Printf("   %5d %14s %9d %9d    %s   (él usó Y = %d)\n", f.x, antTxt, p, q, veredicto, f.yUsada)
		} else {
			fmt.Printf("   %5d %14s %9s %9s    %s   (él usó Y = %d)\n", f.x, antTxt, "—", "—", veredicto, f.yUsada)
		}
	}
	fmt.Println("\n   📌 ACÁ YA HAY ALGO QUE DECIR, Y ES DE SU PROPIA TABLA: en el renglón del 3")
	fmt.Println("   el primo anterior es el 2, y 3+3−2 = 4 y 3+3+2 = 8 no son primos. Él usó")
	fmt.Println("   Y = 1, que no es primo. Así que la regla estricta se rompe en el segundo")
	fmt.Println("   renglón que él mismo escribió. El único que cierra con el anterior es el 5.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · LA REGLA ESTRICTA, BARRIDA SOBRE TODOS LOS PRIMOS")
	fmt.Println("   X primo · Y = el primo inmediatamente anterior · ¿P y Q primos los dos?")
	aciertos, probados := 0, 0
	primerosFallos := []string{}
	primerosAciertos := []string{}
	for i := 1; i < len(primos); i++ {
		x, y := primos[i], primos[i-1]
		p, q := 2*x-y, 2*x+y
		probados++
		if es[p] && es[q] {
			aciertos++
			if len(primerosAciertos) < 8 {
				primerosAciertos = append(primerosAciertos, fmt.Sprintf("%d(y=%d)→%d,%d", x, y, p, q))
			}
		} else if len(primerosFallos) < 8 {
			primerosFallos = append(primerosFallos, fmt.Sprintf("%d(y=%d)→%d,%d", x, y, p, q))
		}
	}
	fmt.Printf("\n        pares (X, anterior) probados ......... %d\n", probados)
	fmt.Printf("        cierran los dos lados ................ %d\n", aciertos)
	fmt.Printf("        tasa de acierto ...................... %.4f%%\n", 100*float64(aciertos)/float64(probados))
	fmt.Println("\n        primeros que CIERRAN:  " + strings.Join(primerosAciertos, "  "))
	fmt.Println("        primeros que FALLAN:   " + strings.Join(primerosFallos, "  "))

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ¿ES EL «PRIMO ANTERIOR» O ALCANZA CON QUE SEA PRIMO?")
	fmt.Println("   Para cada X primo, buscamos el MENOR desvío Y que hace primos los dos lados.")
	fmt.Println("   Y después preguntamos qué es ese Y ganador.")
	var conY0, yPrimo, yUno, yCompuesto, yEsAnterior, sinY int
	muestra := 0
	if len(primos) > 200000 {
		muestra = 200000
	} else {
		muestra = len(primos)
	}
	for i := 1; i < muestra; i++ {
		x, ant := primos[i], primos[i-1]
		hallada := -1
		for y := 0; 2*x+y <= N && y < 2*x; y++ {
			if es[2*x-y] && es[2*x+y] {
				hallada = y
				break
			}
		}
		switch {
		case hallada < 0:
			sinY++
		case hallada == 0:
			conY0++
		case hallada == 1:
			yUno++
		case es[hallada]:
			yPrimo++
			if hallada == ant {
				yEsAnterior++
			}
		default:
			yCompuesto++
		}
	}
	total := muestra - 1
	fmt.Printf("\n        primos X examinados .................. %d\n", total)
	fmt.Printf("        el menor Y que sirve es 0 ............ %d  (%.2f%%)  ← 2X ya era primo, imposible: 2X es par\n", conY0, 100*float64(conY0)/float64(total))
	fmt.Printf("        el menor Y que sirve es 1 ............ %d  (%.2f%%)\n", yUno, 100*float64(yUno)/float64(total))
	fmt.Printf("        el menor Y que sirve es PRIMO ........ %d  (%.2f%%)\n", yPrimo, 100*float64(yPrimo)/float64(total))
	fmt.Printf("           …y de esos, es el ANTERIOR a X .... %d  (%.2f%% del total)\n", yEsAnterior, 100*float64(yEsAnterior)/float64(total))
	fmt.Printf("        el menor Y que sirve es COMPUESTO .... %d  (%.2f%%)\n", yCompuesto, 100*float64(yCompuesto)/float64(total))
	fmt.Printf("        sin ningún Y ......................... %d\n", sinY)

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · EL CONTROL: ¿HACE FALTA QUE X SEA PRIMO?")
	fmt.Println("   Misma pregunta con X COMPUESTO y Y = el primo inmediatamente anterior a X.")
	fmt.Println("   Si la tasa es parecida, entonces «X primo» no está haciendo ningún trabajo.")
	// PREDICCIÓN ANOTADA ANTES DE MIRAR: los compuestos que cierran tienen que ser
	// casi todos múltiplos de 3, porque ésa es la única puerta que la LEY 5 deja
	// abierta y a los primos les está cerrada.
	acC, prC, acC3 := 0, 0, 0
	var excepciones []string
	ant := 0
	for x := 3; x <= tope && prC < probados; x++ {
		if es[x] {
			ant = x
			continue
		}
		if ant == 0 {
			continue
		}
		p, q := 2*x-ant, 2*x+ant
		prC++
		if p >= 0 && q <= N && es[p] && es[q] {
			acC++
			if x%3 == 0 {
				acC3++
			} else {
				excepciones = append(excepciones, fmt.Sprintf("X=%d (anterior %d) → %d, %d", x, ant, p, q))
			}
		}
	}
	fmt.Printf("\n        X compuestos probados ................ %d\n", prC)
	fmt.Printf("        cierran los dos lados ................ %d\n", acC)
	fmt.Printf("        tasa de acierto ...................... %.4f%%\n", 100*float64(acC)/float64(prC))
	fmt.Printf("        de los que cierran, múltiplos de 3 ... %d de %d\n", acC3, acC)
	if len(excepciones) == 0 {
		fmt.Println("        excepciones (cierran sin ser múltiplo de 3): NINGUNA")
	} else {
		fmt.Printf("        excepciones (cierran sin ser múltiplo de 3): %d\n", len(excepciones))
		for _, e := range excepciones {
			fmt.Println("           · " + e)
		}
		fmt.Println("        📌 y ojo con esa lista: la única puerta que queda es que el primo")
		fmt.Println("        anterior sea EL 3, el único primo múltiplo de 3. La misma puerta")
		fmt.Println("        por la que se cuela el (5, 3) del capitán — la ley no tiene fugas.")
	}

	tasaP := 100 * float64(aciertos) / float64(probados)
	tasaC := 100 * float64(acC) / float64(prC)

	// ---- LEY 5: el POR QUÉ, y es demostrable en tres renglones ----
	fmt.Println("\nLEY 5 · ⚡ Y ACÁ ESTÁ EL POR QUÉ — Y ES UNA DEMOSTRACIÓN, NO UN CONTEO")
	fmt.Println("   Tomá X e Y primos, los dos mayores que 3. Ninguno es múltiplo de 3, así que")
	fmt.Println("   cada uno deja resto 1 o 2 al dividir por 3. Cuatro casos, y nada más:")
	fmt.Println("\n        X mod 3   Y mod 3   2X − Y mod 3   2X + Y mod 3   ¿quién muere?")
	for _, xr := range []int{1, 2} {
		for _, yr := range []int{1, 2} {
			pr, qr := ((2*xr-yr)%3+3)%3, (2*xr+yr)%3
			quien := "—"
			if pr == 0 {
				quien = "muere 2X − Y (múltiplo de 3)"
			} else if qr == 0 {
				quien = "muere 2X + Y (múltiplo de 3)"
			}
			fmt.Printf("   %7d %9d %14d %14d   %s\n", xr, yr, pr, qr, quien)
		}
	}
	fmt.Println("\n   ⟹ EN LOS CUATRO CASOS UNO DE LOS DOS LADOS ES MÚLTIPLO DE 3. Y un múltiplo")
	fmt.Println("   de 3 mayor que 3 no puede ser primo. **Así que la regla es IMPOSIBLE para")
	fmt.Println("   todo X primo mayor que 5.** No falla a veces: no puede cerrar nunca.")
	fmt.Println("\n   📌 Y LAS DOS ÚNICAS ESCAPATORIAS SON LOS DOS RENGLONES DE SU TABLA:")
	fmt.Println("     · X = 3 — el propio X es múltiplo de 3, el argumento no aplica. Da 4 y 8:")
	fmt.Println("       no cierra igual, pero por otro motivo.")
	fmt.Println("     · X = 5, Y = 3 — el 3 es EL ÚNICO primo múltiplo de 3, y ahí el argumento")
	fmt.Println("       se rompe. Da 7 y 13. **CIERRA.**")
	fmt.Println("\n   ⟹ **(5, 3) ES EL ÚNICO PAR DE PRIMOS CONSECUTIVOS DEL UNIVERSO QUE CIERRA**,")
	fmt.Println("     y usted lo escribió. No encontró una ley: encontró la única excepción que hay.")
	fmt.Println("\n   Y el control lo confirma por el otro lado: los X compuestos aciertan MÁS")
	fmt.Println("   justamente porque pueden ser múltiplos de 3, que es la puerta que a los")
	fmt.Println("   primos les está cerrada. Medido arriba: de los que cierran, casi todos lo son.")

	// ---- veredicto: SALE de la medición, no está tipeado ----
	fmt.Println("\n════════ VEREDICTO ════════")
	if tasaP > 99.9 {
		fmt.Printf("✅ LA REGLA ESTRICTA CIERRA: %.4f%% de acierto sobre %d primos.\n", tasaP, probados)
	} else {
		fmt.Printf("❌ LA REGLA ESTRICTA NO CIERRA: %d acierto(s) en %d primos (%.4f%%).\n", aciertos, probados, tasaP)
		fmt.Println("   Y no es que falle casi siempre: la LEY 5 demuestra que **NO PUEDE cerrar**")
		fmt.Println("   para ningún X primo mayor que 5. Uno de los dos lados siempre es múltiplo")
		fmt.Println("   de 3. El único acierto de toda la recta es el (5, 3) — su propio renglón.")
	}
	razon := tasaP / tasaC
	fmt.Printf("\n   · X primo:     %.4f%%\n", tasaP)
	fmt.Printf("   · X compuesto: %.4f%%   (razón %.2f×)\n", tasaC, razon)
	if razon > 1.5 {
		fmt.Println("   ⚡ EL CANDADO «X PRIMO» SÍ HACE TRABAJO: la tasa cambia de verdad.")
	} else if razon < 0.67 {
		fmt.Println("   ⚡ EL CANDADO «X PRIMO» HACE TRABAJO, PERO AL REVÉS: con X compuesto acierta MÁS.")
	} else {
		fmt.Println("   📌 EL CANDADO «X PRIMO» NO HACE TRABAJO: la tasa es prácticamente la misma.")
	}
	fmt.Println("\n⚖️ Y LO QUE SÍ QUEDA EN PIE, QUE NO ES POCO: el capitán dejó de defender el")
	fmt.Println("  enunciado ancho y lo APRETÓ hasta volverlo falsable en un renglón. Eso es")
	fmt.Println("  método, no suerte. Una regla que se puede matar con un solo contraejemplo")
	fmt.Println("  vale más que una frase que no se puede matar con nada.")

	escribirLamina(probados, aciertos, tasaP, tasaC, total, yPrimo, yUno, yCompuesto, yEsAnterior)
}

func escribirLamina(probados, aciertos int, tasaP, tasaC float64, total, yPrimo, yUno, yCompuesto, yEsAnterior int) {
	var b strings.Builder
	W, H := 1520.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔙 EL ANTERIOR — el capitán le pone dos candados a su fórmula</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">X debe ser primo · Y debe ser el primo anterior · y eso se puede matar con un solo contraejemplo</text>
`, W, H, W, H, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="40" y="106" width="1440" height="180" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="760" y="140" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">EL ENUNCIADO APRETADO</text>
<text x="760" y="180" font-size="19" text-anchor="middle" font-family="monospace" fill="#dce8f7">X primo   ·   Y = el primo inmediatamente anterior a X</text>
<text x="760" y="212" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">P = 2X − Y      Q = 2X + Y      ¿los dos primos?</text>
<text x="760" y="250" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">Ya no habla de todos los pares: dice que un primo, con el que lo precede, FABRICA otros dos.</text>
<text x="760" y="274" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Eso es falsable en un renglón — y por eso vale más que la versión ancha.</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="306" width="710" height="240" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="395" y="338" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">SU PROPIA TABLA LO ROMPE</text>
<text x="70" y="374" font-size="15" font-family="monospace" fill="#cfe6ff">X = 2   anterior: no hay      él usó Y = 1</text>
<text x="70" y="402" font-size="15" font-family="monospace" fill="#ff8fa0">X = 3   anterior: 2   →  4 y 8   ❌ ninguno primo</text>
<text x="70" y="430" font-size="15" font-family="monospace" fill="#7ee0c0">X = 5   anterior: 3   →  7 y 13  ✅ los dos primos</text>
<text x="70" y="470" font-size="14.5" font-family="Georgia" fill="#f3d9cf">En el renglón del 3 él no usó el primo anterior: usó el 1,</text>
<text x="70" y="492" font-size="14.5" font-family="Georgia" fill="#f3d9cf">que no es primo. La regla estricta se rompe adentro de</text>
<text x="70" y="514" font-size="14.5" font-family="Georgia" fill="#f3d9cf">la misma tabla que la propone. Cierra uno de tres.</text>`)

	fmt.Fprintf(&b, `<rect x="770" y="306" width="710" height="240" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1125" y="338" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL BARRIDO, Y EL CONTROL</text>
<text x="800" y="376" font-size="15" font-family="monospace" fill="#cfe6ff">pares (X, anterior) probados .... %d</text>
<text x="800" y="402" font-size="15" font-family="monospace" fill="#cfe6ff">cierran los dos lados ........... %d</text>
<text x="800" y="434" font-size="17" font-family="monospace" fill="#ff8fa0">X primo:     %.4f%%</text>
<text x="800" y="462" font-size="17" font-family="monospace" fill="#8fb4d9">X compuesto: %.4f%%</text>
<text x="800" y="500" font-size="14.5" font-family="Georgia" fill="#cfe6ff">Si las dos tasas son parecidas, el candado «X primo»</text>
<text x="800" y="522" font-size="14.5" font-family="Georgia" fill="#cfe6ff">no está haciendo ningún trabajo. El control lo dice.</text>
`, probados, aciertos, tasaP, tasaC)

	fmt.Fprintf(&b, `<rect x="40" y="566" width="1440" height="180" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="760" y="598" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">Y AHORA LA PREGUNTA BUENA: ¿QUÉ ES EL Y QUE SÍ SIRVE?</text>
<text x="760" y="628" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Para cada primo X buscamos el MENOR desvío que hace primos los dos lados, y le preguntamos qué es.</text>
<text x="200" y="668" font-size="15" font-family="monospace" fill="#7ee0c0">Y = 1 ......... %d</text>
<text x="200" y="694" font-size="15" font-family="monospace" fill="#7ee0c0">Y primo ....... %d</text>
<text x="760" y="668" font-size="15" font-family="monospace" fill="#ffd98a">…y de esos, el ANTERIOR a X ..... %d</text>
<text x="760" y="694" font-size="15" font-family="monospace" fill="#ff8fa0">Y compuesto ..................... %d</text>
<text x="760" y="728" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">sobre %d primos examinados</text>
`, yUno, yPrimo, yEsAnterior, yCompuesto, total)

	fmt.Fprintf(&b, `<rect x="40" y="766" width="1440" height="140" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="760" y="798" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ LO QUE QUEDA EN PIE</text>
<text x="760" y="832" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">El capitán no defendió el enunciado ancho: lo APRETÓ hasta volverlo falsable en un renglón.</text>
<text x="760" y="858" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Una regla que se puede matar con un solo contraejemplo vale más que una frase que no se puede matar con nada.</text>
<text x="760" y="888" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Eso es método. Y el método es lo único que después se transfiere.</text>
</svg>
`)

	if err := os.WriteFile("el-anterior.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-anterior.svg")
}
