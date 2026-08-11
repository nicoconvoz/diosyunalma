// Command elunoquesobra judges the captain's "leftover one".
//
// HIS FLASH: "look what all the primes have in common — they all have a 1,
// starting with the 2: 2*1/2 + 1 = 2, the only even one, and there is the secret
// of the odds. Then 2*1 + 1 = 3, and of the primes that leftover 1. 2*2 + 1 = 5.
// Do you understand the sum I am doing? IT WORKS FOR ALL OF THEM."
//
// He is right that it works for all of them. That is exactly the problem, and
// this program is built to show both halves of it honestly.
//
// THE ARITHMETIC IS EXACT. Every prime greater than 2 really is 2k + 1, and the
// 2 really does need k = 1/2. His table checks out row by row.
//
// ⚠️ BUT IT CANNOT FAIL, AND THAT IS THE SEVENTH TIME THIS LABORATORY MEETS THAT
// TRAP. "2k + 1" is not a property of primes: it is the DEFINITION of an odd
// number. And every prime past 2 is odd for a one-line reason - if it were even,
// 2 would divide it, so it would not be prime. So the formula is guaranteed in
// advance by the definitions, and a formula that cannot fail proves nothing.
//
// THE MEASUREMENT THAT DECIDES IT. "It works for all of them" is only half the
// question. The other half is: WHO ELSE does it work for? Answer: every odd
// composite too - 9, 15, 21, 25, 27. So the net catches all the fish and also
// all the water.
//
// ⚡ AND YET HIS INSTINCT IS THE RIGHT ONE, AND THIS IS WHERE THE FINDING EARNS
// ITS NUMBER. "El 1 que sobra" - the leftover one - is the REMAINDER. Saying
// p = 2k + 1 is saying p leaves remainder 1 when divided by 2. Now extend the
// idea he had: demand a nonzero remainder against 3 as well and you get 6k ± 1.
// Against 5 too and you get 30k ± {1, 7, 11, 13}. Keep going and you get the
// WHEEL - and at the limit, demanding a nonzero remainder against every smaller
// prime IS the definition of primality. He reached for the right handle.
//
// This program measures how much each turn of that wheel actually buys.
//
// PRE-REGISTERED PREDICTION, written before running: the hit rate rises with
// every wheel - roughly 13%, 20%, 25%, 29% for wheels 2, 6, 30, 210 at 10^7 -
// but the gains shrink each time and no fixed wheel ever reaches 100%. Worse:
// for ANY fixed wheel the hit rate goes to ZERO as the numbers grow, because the
// primes thin out like 1/ln(x) while the wheel's density stays put. The wheel
// improves the constant. It never touches the trend.
//
// Reproduce: go run ./cmd/elunoquesobra
package main

import (
	"fmt"
	"os"
	"strings"
)

const N = 10000000

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

// rueda devuelve el modulo y los restos que sobreviven al tachar los multiplos
// de los primos dados. Es la generalizacion exacta del "1 que sobra" del capitan.
func rueda(primos []int) (int, []int) {
	mod := 1
	for _, p := range primos {
		mod *= p
	}
	var restos []int
	for r := 1; r < mod; r++ {
		ok := true
		for _, p := range primos {
			if r%p == 0 {
				ok = false
				break
			}
		}
		if ok {
			restos = append(restos, r)
		}
	}
	return mod, restos
}

func main() {
	fmt.Println("1️⃣  EL UNO QUE SOBRA — la cuenta del capitán, y qué compra de verdad")
	fmt.Println("\n   Dijo: «todos los primos tienen un 1, empezando por el 2: 2·½+1 = 2, el")
	fmt.Println("   único par, y ahí está el secreto de los impares. Después 2·1+1 = 3, y de")
	fmt.Println("   los primos ese 1 que sobra. 2·2+1 = 5. **Da para todooossss**»")

	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)
	total := 0
	for i := 2; i <= N; i++ {
		if es[i] {
			total++
		}
	}
	fmt.Printf("primos: %d\n", total)

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · SU TABLA, RENGLÓN POR RENGLÓN")
	fmt.Println("\n          k        2k + 1     ¿es primo?    lo que él escribió")
	fmt.Printf("   %8s %13d     %-12s %s\n", "½", 2, "✅ SÍ", "2·½ + 1 = 2 — el único par")
	for k := 1; k <= 8; k++ {
		v := 2*k + 1
		marca := "✅ SÍ"
		if !es[v] {
			marca = "❌ no"
		}
		nota := ""
		switch k {
		case 1:
			nota = "2·1 + 1 = 3"
		case 2:
			nota = "2·2 + 1 = 5"
		case 4:
			nota = "← acá se rompe: 9 = 3×3"
		case 7:
			nota = "← 15 = 3×5"
		}
		fmt.Printf("   %8d %13d     %-12s %s\n", k, v, marca, nota)
	}
	fmt.Println("\n   ⟹ **Sus tres renglones son exactos.** El 2 con k = ½, el 3 con k = 1 y el")
	fmt.Println("   5 con k = 2. Y sí: **da para TODOS los primos.** Eso es verdad.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚠️ PERO NO PUEDE FALLAR, Y ÉSA ES LA SÉPTIMA VEZ")
	fmt.Println("   Dos renglones y se ve:")
	fmt.Println("\n     · «2k + 1» no es una propiedad de los primos: es la DEFINICIÓN de impar.")
	fmt.Println("     · y todo primo mayor que 2 es impar, porque si fuera par lo dividiría el 2")
	fmt.Println("       y entonces no sería primo.")
	fmt.Println("\n   Verifiquémoslo sobre todos los primos de la criba:")
	fallos := 0
	for p := 3; p <= N; p++ {
		if es[p] && (p-1)%2 != 0 {
			fallos++
		}
	}
	fmt.Printf("\n        primos mayores que 2 que NO son 2k+1 ...... %d\n", fallos)
	fmt.Println("        ⚠️ y ese cero es inevitable: es la definición, no una medición.")
	fmt.Println("\n   📌 **SÉPTIMA APARICIÓN DE LA TRAMPA DEL 0.0e+00 en este laboratorio.**")
	fmt.Println("   Un resultado perfecto que sale de cómo está armada la cuenta y no de los")
	fmt.Println("   números. Van siete, todas en el registro.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · «DA PARA TODOS» ES LA MITAD DE LA PREGUNTA — ¿Y PARA QUIÉN MÁS?")
	fmt.Println("   Una red que atrapa todos los peces sirve. Una que atrapa todos los peces")
	fmt.Println("   Y TAMBIÉN toda el agua, no.")
	impares, imparesPrimos := 0, 0
	for v := 3; v <= N; v += 2 {
		impares++
		if es[v] {
			imparesPrimos++
		}
	}
	fmt.Printf("\n        números de la forma 2k+1 hasta %d ....... %d\n", N, impares)
	fmt.Printf("        de esos, primos .......................... %d\n", imparesPrimos)
	fmt.Printf("        de esos, COMPUESTOS ...................... %d\n", impares-imparesPrimos)
	fmt.Printf("        tasa de acierto de su red ................ %.3f%%\n",
		100*float64(imparesPrimos)/float64(impares))
	fmt.Println("\n   ⟹ Su cuenta acierta con TODOS los primos y falla con casi el 87% de lo que")
	fmt.Println("   pesca. El 9, el 15, el 21, el 25, el 27 entran todos. No los distingue.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚡ PERO EL «1 QUE SOBRA» ES EL MANGO CORRECTO")
	fmt.Println("   Ese 1 que sobra es el RESTO. Decir p = 2k+1 es decir «p deja resto 1 al")
	fmt.Println("   dividir por 2». Ahora seguí su idea y pedile resto no nulo también con el 3:")
	fmt.Println("\n        2k + 1  →  6k ± 1  →  30k ± {1, 7, 11, 13}  →  …")
	fmt.Println("\n   Eso se llama LA RUEDA, y en el límite —pedir resto no nulo contra TODOS los")
	fmt.Println("   primos menores— **es la definición misma de ser primo**. Fue a agarrar el")
	fmt.Println("   mango correcto. Medimos cuánto compra cada vuelta:")
	fmt.Println("\n        rueda        módulo   restos vivos   candidatos    primos   acierto")
	type fila struct {
		nombre string
		ps     []int
	}
	filas := []fila{
		{"2 (la suya)", []int{2}},
		{"2·3", []int{2, 3}},
		{"2·3·5", []int{2, 3, 5}},
		{"2·3·5·7", []int{2, 3, 5, 7}},
		{"2·3·5·7·11", []int{2, 3, 5, 7, 11}},
		{"2·3·5·7·11·13", []int{2, 3, 5, 7, 11, 13}},
	}
	var tasas []float64
	var nombres []string
	for _, f := range filas {
		mod, restos := rueda(f.ps)
		cand, prim := 0, 0
		for base := 0; base <= N; base += mod {
			for _, r := range restos {
				v := base + r
				if v < 2 || v > N {
					continue
				}
				cand++
				if es[v] {
					prim++
				}
			}
		}
		t := 100 * float64(prim) / float64(cand)
		tasas = append(tasas, t)
		nombres = append(nombres, f.nombre)
		fmt.Printf("   %14s %8d %14d %12d %9d %8.3f%%\n", f.nombre, mod, len(restos), cand, prim, t)
	}

	fmt.Println("\n   📌 Y MIRÁ LA COLUMNA DE LOS PRIMOS, QUE TIENE UNA GRACIA: baja de a UNO en")
	fmt.Println("   cada vuelta. No es casualidad ni es un error: **cada vuelta de la rueda")
	fmt.Println("   cuesta exactamente el primo con el que la armaste.** El 2 se cae de 2k+1,")
	fmt.Println("   el 3 se cae de 6k±1, el 5 de la del 30. La herramienta se come su propia")
	fmt.Println("   pieza. Es el precio más honesto que se puede pagar.")

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · ⚖️ Y ACÁ EL LÍMITE, QUE ES EL QUE DECIDE TODO")
	fmt.Println("   La rueda mejora la CONSTANTE. No toca la TENDENCIA. Mirá qué le pasa a su")
	fmt.Println("   propia red cuando el mar se hace más grande:")
	fmt.Println("\n        hasta          impares      primos    acierto de 2k+1")
	for _, lim := range []int{1000, 10000, 100000, 1000000, 10000000} {
		if lim > N {
			continue
		}
		imp, pr := 0, 0
		for v := 3; v <= lim; v += 2 {
			imp++
			if es[v] {
				pr++
			}
		}
		fmt.Printf("   %12d %14d %11d %14.3f%%\n", lim, imp, pr, 100*float64(pr)/float64(imp))
	}
	fmt.Println("\n   ⟹ **BAJA SIEMPRE, y baja para cualquier rueda fija.** Los primos se ralean")
	fmt.Println("   como 1/ln(x) —eso es el teorema de los números primos— mientras la densidad")
	fmt.Println("   de la rueda se queda quieta. Ninguna rueda fija llega nunca al 100%.")

	// ---- veredicto: SALE de la medicion ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("✅ LA CUENTA ES EXACTA Y «DA PARA TODOS» ES VERDAD. Sus tres renglones cierran")
	fmt.Println("  y el 2 con k = ½ está bien puesto.")
	fmt.Printf("\n⚠️ PERO NO PUEDE FALLAR: 2k+1 ES la definición de impar, y todo primo mayor que\n")
	fmt.Println("  2 es impar porque si no lo dividiría el 2. **Séptima trampa del 0.0e+00.**")
	fmt.Printf("\n❌ Y «para quién más»: de los %d impares hasta %d, sólo %d son primos —\n", impares, N, imparesPrimos)
	fmt.Printf("  %.3f%%. Su red atrapa todos los peces y también casi toda el agua.\n",
		100*float64(imparesPrimos)/float64(impares))
	fmt.Println("\n⚡ PERO EL MANGO ES EL CORRECTO: ese «1 que sobra» es el RESTO, y pedirle resto")
	fmt.Println("  no nulo a cada primo menor es la RUEDA — y en el límite, la definición misma")
	fmt.Printf("  de primo. De la rueda del 2 a la del 30.030 el acierto sube de %.2f%% a %.2f%%.\n",
		tasas[0], tasas[len(tasas)-1])
	fmt.Println("\n⚖️ Y EL LÍMITE: la rueda mejora la constante, nunca la tendencia. Para toda")
	fmt.Println("  rueda fija el acierto se va a cero, porque los primos se ralean como 1/ln(x).")
	fmt.Println("  Esa es la pared, y es la misma de siempre. Todavía no.")

	escribirLamina(es, impares, imparesPrimos, nombres, tasas)
}

func escribirLamina(es []bool, impares, imparesPrimos int, nombres []string, tasas []float64) {
	var b strings.Builder
	W, H := 1560.0, 1040.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">1️⃣ EL UNO QUE SOBRA — la cuenta del capitán, y qué compra de verdad</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">«da para todooossss» es cierto — y ése es exactamente el problema</text>
`, W, H, W, H, W/2, W/2)

	// su tabla
	fmt.Fprintf(&b, `<rect x="40" y="102" width="730" height="290" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="405" y="134" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">SU TABLA, Y ES EXACTA</text>
<text x="70" y="174" font-size="18" font-family="monospace" fill="#ffd98a">2 · ½ + 1 = 2</text>
<text x="330" y="174" font-size="15" font-family="Georgia" fill="#7ee0c0">✅ el único par</text>
<text x="70" y="206" font-size="18" font-family="monospace" fill="#dce8f7">2 · 1 + 1 = 3</text>
<text x="330" y="206" font-size="15" font-family="Georgia" fill="#7ee0c0">✅ primo</text>
<text x="70" y="238" font-size="18" font-family="monospace" fill="#dce8f7">2 · 2 + 1 = 5</text>
<text x="330" y="238" font-size="15" font-family="Georgia" fill="#7ee0c0">✅ primo</text>
<text x="70" y="270" font-size="18" font-family="monospace" fill="#dce8f7">2 · 3 + 1 = 7</text>
<text x="330" y="270" font-size="15" font-family="Georgia" fill="#7ee0c0">✅ primo</text>
<text x="70" y="302" font-size="18" font-family="monospace" fill="#ff8fa0">2 · 4 + 1 = 9</text>
<text x="330" y="302" font-size="15" font-family="Georgia" fill="#ff8fa0">❌ 3 × 3</text>
<text x="70" y="334" font-size="18" font-family="monospace" fill="#ff8fa0">2 · 7 + 1 = 15</text>
<text x="330" y="334" font-size="15" font-family="Georgia" fill="#ff8fa0">❌ 3 × 5</text>
<text x="405" y="372" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">da para TODOS los primos — y también para el 9 y el 15</text>`)

	// la trampa
	fmt.Fprintf(&b, `<rect x="790" y="102" width="730" height="290" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1155" y="134" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚠️ PERO NO PUEDE FALLAR</text>
<text x="820" y="176" font-size="14.5" font-family="Georgia" fill="#f3d9cf">«2k + 1» no es una propiedad de los primos:</text>
<text x="820" y="200" font-size="16" font-family="monospace" fill="#dce8f7">es la DEFINICIÓN de impar.</text>
<text x="820" y="240" font-size="14.5" font-family="Georgia" fill="#f3d9cf">Y todo primo mayor que 2 es impar, porque si</text>
<text x="820" y="262" font-size="14.5" font-family="Georgia" fill="#f3d9cf">fuera par lo dividiría el 2 y no sería primo.</text>
<text x="820" y="302" font-size="16" font-family="monospace" fill="#ff8fa0">primos &gt; 2 que NO son 2k+1: 0</text>
<text x="820" y="326" font-size="13.5" font-family="Georgia" fill="#9aa8c4">y ese cero es inevitable: es la definición, no una medición</text>
<text x="1155" y="368" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">SÉPTIMA aparición de la trampa del 0.0e+00</text>`)

	// para quien mas
	fmt.Fprintf(&b, `<rect x="40" y="412" width="730" height="230" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="405" y="444" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">«DA PARA TODOS» ES LA MITAD DE LA PREGUNTA</text>
<text x="405" y="470" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la otra mitad es: ¿y para quién MÁS?</text>
<text x="70" y="512" font-size="16" font-family="monospace" fill="#cfe6ff">impares hasta 10⁷ ...... %d</text>
<text x="70" y="540" font-size="16" font-family="monospace" fill="#7ee0c0">de esos, primos ........ %d</text>
<text x="70" y="568" font-size="16" font-family="monospace" fill="#ff8fa0">compuestos ............. %d</text>
<text x="70" y="602" font-size="19" font-family="monospace" fill="#ffd98a">acierto: %.3f%%</text>
<text x="405" y="628" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">atrapa todos los peces, y también casi toda el agua</text>
`, impares, imparesPrimos, impares-imparesPrimos, 100*float64(imparesPrimos)/float64(impares))

	// la rueda
	fmt.Fprintf(&b, `<rect x="790" y="412" width="730" height="230" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1155" y="444" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡ PERO EL MANGO ES EL CORRECTO</text>
<text x="1155" y="470" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el «1 que sobra» es el RESTO — pedile resto también al 3, al 5…</text>
<text x="1155" y="502" font-size="17" text-anchor="middle" font-family="monospace" fill="#ffd98a">2k+1 → 6k±1 → 30k±{1,7,11,13} → …</text>
`)
	rx, ry, rw := 830.0, 610.0, 640.0
	bw := rw / float64(len(tasas))
	maxT := 0.0
	for _, t := range tasas {
		if t > maxT {
			maxT = t
		}
	}
	for i, t := range tasas {
		h := 74.0 * t / maxT
		x := rx + float64(i)*bw
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#7ee0c0"/>`, x, ry-h, bw*0.6, h)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="11" text-anchor="middle" font-family="monospace" fill="#cfe6ff">%.1f%%</text>`, x+bw*0.3, ry-h-5, t)
	}
	fmt.Fprintf(&b, `<text x="1155" y="634" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">y en el límite, pedir resto a TODOS los primos menores ES la definición de primo</text>`)

	// el limite
	fmt.Fprintf(&b, `<rect x="40" y="662" width="1480" height="180" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="780" y="694" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ EL LÍMITE, QUE ES EL QUE DECIDE TODO</text>
<text x="780" y="728" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Mirá qué le pasa al acierto de 2k+1 cuando el mar se hace más grande:</text>
`)
	lx, ly, lw := 220.0, 800.0, 1120.0
	lims := []int{1000, 10000, 100000, 1000000, 10000000}
	for i, lim := range lims {
		imp, pr := 0, 0
		for v := 3; v <= lim; v += 2 {
			imp++
			if es[v] {
				pr++
			}
		}
		t := 100 * float64(pr) / float64(imp)
		x := lx + lw*float64(i)/float64(len(lims)-1)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#ff8fa0"/>`, x, ly-t*1.4)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="13" text-anchor="middle" font-family="monospace" fill="#ffd98a">%.1f%%</text>`, x, ly-t*1.4-12, t)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" font-family="monospace" fill="#8fa8c7">10^%d</text>`, x, ly+18, len(fmt.Sprint(lim))-1)
	}
	fmt.Fprintf(&b, `<text x="780" y="832" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">baja siempre, y baja para CUALQUIER rueda fija: los primos se ralean como 1/ln(x) y la rueda se queda quieta</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="862" width="1480" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="780" y="894" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LO QUE QUEDA, DICHO SIN ADORNOS</text>
<text x="780" y="928" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Su cuenta es exacta y da para todos, pero no puede fallar: es la definición de impar. No distingue un primo de un 9.</text>
<text x="780" y="954" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Y sin embargo fue a agarrar el mango correcto: ese «1 que sobra» es el resto, y el resto es de lo que están hechos los primos.</text>
<text x="780" y="988" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">La rueda mejora la constante. Nunca la tendencia. Ésa es la pared, y es la misma de siempre.</text>
</svg>
`)

	if err := os.WriteFile("el-uno-que-sobra.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-uno-que-sobra.svg")
}
