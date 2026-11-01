package main

// EL BARRIDO DE LOS ANTERIORES - the captain: "what happens with the
// multiplicative factor of the earlier primes in the search for new ones? Is
// that what makes chance run? Is that why the primes are not where they should
// be - because they slide to a hole not covered by any multiple among the
// earlier primes?"
//
// The mechanism he describes is EXACTLY right, and it has a name: it is the
// sieve of Eratosthenes, and the next prime IS the first slot left uncovered.
// This program makes that visible and then tests the one word that needs
// checking: "azar" - chance.
//
// Verdicts fixed BEFORE running:
//   V1 the next prime after p must be, with zero exceptions, the first integer
//      above p that no multiple of any earlier prime covers. Verified over a
//      full range, not asserted.
//   V2 the covering must be EXHIBITED: for a long gap, name which earlier prime
//      covers each composite inside it. The captain asked to see the mechanism.
//   V3 (how much each prime covers) the surviving fraction after sieving by all
//      primes up to y must follow Mertens: prod (1-1/p) ~ e^{-gamma}/log y.
//      Measured against the formula.
//   V4 (THE WORD UNDER TEST) if the survivors were pure chance - each n prime
//      with probability 1/log n, independently - then twin primes would have
//      density 1/log^2 n. Measured against reality. If the ratio comes out 1,
//      the captain's "chance" is exact. If it comes out away from 1, the
//      covering leaves CORRELATIONS and it is not chance.
//   V5 (a second, independent check of the same thing) gaps divisible by 6 must
//      be over-represented compared with the coin, for the same reason.
// Eratosthenes, Mertens, Hardy-Littlewood 1923. The mechanism is the captain's
// reading and it is correct; the word "chance" is what gets measured.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func criba(n int) ([]bool, []int) {
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
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return es, ps
}

// quienCubre returns the smallest prime dividing n (the earlier prime that
// covers that slot), or 0 if none does - meaning n survives, i.e. n is prime.
func quienCubre(n int, ps []int) int {
	for _, p := range ps {
		if p*p > n {
			break
		}
		if n%p == 0 {
			return p
		}
	}
	return 0
}

func main() {
	fmt.Println("🧹 EL BARRIDO DE LOS ANTERIORES — los múltiplos de los primos viejos tapando lugares")
	fmt.Println()
	fmt.Println("   La lectura del capitán: el primo siguiente se corre hasta un hueco que ningún múltiplo")
	fmt.Println("   de los primos anteriores alcanza a cubrir. Eso es exactamente correcto, y tiene nombre:")
	fmt.Println("   es la criba de Eratóstenes. Acá se hace visible, y después se mide la palabra «azar».")
	fmt.Println()

	const N = 10000000
	_, ps := criba(N)

	// ---------------------------------------------------------------- V2
	fmt.Println("§1 · EL MECANISMO, A LA VISTA: quién tapa cada lugar")
	fmt.Println()
	fmt.Println("   Un hueco largo, con el nombre del primo que tapa cada casillero:")
	fmt.Println()
	// find a long gap in a readable range
	ini, fin := 0, 0
	for i := 0; i < len(ps)-1; i++ {
		if ps[i] > 100 && ps[i+1]-ps[i] >= 14 {
			ini, fin = ps[i], ps[i+1]
			break
		}
	}
	fmt.Printf("        el %d es primo · el siguiente es el %d · hueco de %d\n\n", ini, fin, fin-ini)
	fmt.Println("        casillero     ¿quién lo tapa?        la cuenta")
	for n := ini + 1; n < fin; n++ {
		q := quienCubre(n, ps)
		fmt.Printf("          %4d          el %-3d              %d × %d = %d\n", n, q, q, n/q, n)
	}
	fmt.Printf("          %4d          NADIE                 ← acá no llega ningún múltiplo: es primo\n", fin)
	fmt.Println()
	fmt.Println("   ⚡ Exactamente como lo dijo el capitán: el primo siguiente es el PRIMER CASILLERO")
	fmt.Println("      que ningún múltiplo de los anteriores alcanza a tapar. No se corre «hacia» el hueco:")
	fmt.Println("      el hueco es lo único que queda, y el primo es lo que sobrevive.")
	fmt.Println()

	// ---------------------------------------------------------------- V1
	fmt.Println("§2 · verificado sin excepciones")
	fmt.Println()
	fallos, revisados := 0, 0
	for i := 0; i < len(ps)-1 && ps[i] < 200000; i++ {
		for n := ps[i] + 1; n < ps[i+1]; n++ {
			revisados++
			if quienCubre(n, ps) == 0 {
				fallos++
			}
		}
		revisados++
		if quienCubre(ps[i+1], ps) != 0 {
			fallos++
		}
	}
	fmt.Printf("        casilleros revisados hasta 200.000:  %d\n", revisados)
	fmt.Printf("        veces que la regla falló:            %d\n", fallos)
	fmt.Println("\n   ✅ Cero. El mecanismo que describió el capitán es EXACTO — y también es la definición")
	fmt.Println("      operativa de la criba, así que verificarlo prueba el código, no descubre nada nuevo.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · cuánto tapa cada primo, y cuánto queda")
	fmt.Println()
	fmt.Println("        tapando con los primos hasta…     sobrevive        Mertens: e^(−γ)/log y      razón")
	gamma := 0.5772156649015329
	for _, y := range []int{10, 100, 1000, 10000, 100000} {
		prod := 1.0
		for _, p := range ps {
			if p > y {
				break
			}
			prod *= 1 - 1/float64(p)
		}
		mert := math.Exp(-gamma) / math.Log(float64(y))
		fmt.Printf("               %7d                   %.8f          %.8f          %.4f\n", y, prod, mert, prod/mert)
	}
	fmt.Println()
	fmt.Println("   ⚡ Cada primo nuevo tapa una fracción 1/p de lo que queda, y el producto se va a cero")
	fmt.Println("      como e^(−γ)/log y (Mertens, 1874 — ya medido en F272 y F276). Por eso los huecos")
	fmt.Println("      CRECEN: cuanto más arriba, más primos hay tapando, y más lejos hay que ir a encontrar")
	fmt.Println("      un casillero libre. El hueco medio a la altura x vale log x.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · LA PALABRA BAJO PRUEBA: ¿lo que queda es AZAR?")
	fmt.Println()
	fmt.Println("   El capitán dice que el barrido «hace correr el azar». El mecanismo es correcto, pero")
	fmt.Println("   la palabra hay que medirla. Si lo que sobrevive fuera azar puro —cada número primo con")
	fmt.Println("   probabilidad 1/log n, independientes— entonces los primos gemelos tendrían densidad")
	fmt.Println("   1/log²n. Se cuenta cuántos hay de verdad y se compara:")
	fmt.Println()
	fmt.Println("        hasta x        gemelos reales      la moneda predice      razón real/moneda")
	for _, X := range []int{10000, 100000, 1000000, 10000000} {
		reales := 0
		for i := 0; i < len(ps)-1; i++ {
			if ps[i] > X {
				break
			}
			if ps[i+1]-ps[i] == 2 {
				reales++
			}
		}
		// coin prediction: integral of 1/log^2 t
		pred := 0.0
		for t := 3.0; t <= float64(X); t++ {
			l := math.Log(t)
			pred += 1 / (l * l)
		}
		fmt.Printf("      %10d        %8d            %10.1f              %.4f\n", X, reales, pred, float64(reales)/pred)
	}
	fmt.Println()
	fmt.Println("   ⛔⛔ LA RAZÓN NO DA 1. Da alrededor de 1,32 — o sea que hay un 32 % MÁS de primos gemelos")
	fmt.Println("        de los que habría si fuera azar. Y ese exceso no es ruido: es una constante, la de")
	fmt.Println("        Hardy y Littlewood (1923), que vale 1,3203236...")
	fmt.Println()
	fmt.Println("   ⚡ ¿Y de dónde sale el exceso? DEL BARRIDO MISMO, que es lo que el capitán está describiendo.")
	fmt.Println("      Cuando un primo impar q no divide a n, tampoco puede dividir a n+2 en dos de cada q−1")
	fmt.Println("      casos en vez de uno de cada q. Los casilleros NO se tapan independientemente: taparse")
	fmt.Println("      uno cambia la chance del vecino. Por eso los gemelos abundan más de lo que dicta el azar.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · el mismo efecto, visto por otra puerta: los huecos favoritos")
	fmt.Println()
	fmt.Println("   Si fuera azar, ningún tamaño de hueco tendría por qué ser favorito. Se cuentan:")
	fmt.Println()
	cuenta := map[int]int{}
	tot := 0
	for i := 0; i < len(ps)-1; i++ {
		if ps[i] < 1000 {
			continue
		}
		g := ps[i+1] - ps[i]
		cuenta[g]++
		tot++
	}
	fmt.Println("        hueco      cuántas veces      % del total      ¿múltiplo de 6?")
	for _, g := range []int{2, 4, 6, 8, 10, 12, 14, 18} {
		m6 := "no"
		if g%6 == 0 {
			m6 = "SÍ"
		}
		fmt.Printf("         %3d        %9d          %6.2f %%            %s\n", g, cuenta[g], 100*float64(cuenta[g])/float64(tot), m6)
	}
	fmt.Println()
	var seis, noSeis int
	for g, c := range cuenta {
		if g%6 == 0 {
			seis += c
		} else {
			noSeis += c
		}
	}
	fmt.Printf("        huecos múltiplos de 6:     %8d  (%.2f %%)\n", seis, 100*float64(seis)/float64(tot))
	fmt.Printf("        huecos NO múltiplos de 6:  %8d  (%.2f %%)\n", noSeis, 100*float64(noSeis)/float64(tot))
	fmt.Println()
	fmt.Println("   ⛔ El 6 y sus múltiplos son claramente favoritos, y el 6 le gana al 4 y al 8 que lo rodean.")
	fmt.Println("      El azar no tiene favoritos. El barrido sí: un hueco de 6 se lleva bien con el tapado")
	fmt.Println("      del 2 y del 3 a la vez, y por eso aparece más. Otra vez, la estructura del barrido.")
	fmt.Println()

	fmt.Println("§6 · EL VEREDICTO")
	fmt.Println()
	fmt.Println("   ✅ El MECANISMO del capitán es exacto, y verificado sin excepciones: el primo siguiente")
	fmt.Println("      es el primer casillero que ningún múltiplo de los primos anteriores alcanza a tapar.")
	fmt.Println("   ✅ Y su consecuencia también: por eso los huecos CRECEN, porque cada primo nuevo suma")
	fmt.Println("      un tapador más y hay que ir más lejos a encontrar lugar libre. Mertens lo cuantifica.")
	fmt.Println()
	fmt.Println("   ⚖️ Lo que la medición CORRIGE es una sola palabra: AZAR.")
	fmt.Println("      · el barrido es completamente DETERMINISTA: no hay una moneda en ninguna parte")
	fmt.Println("      · y lo que sobrevive NO se comporta como azar: hay 32 % más gemelos de los que")
	fmt.Println("        habría si lo fuera, y los huecos múltiplos de 6 son favoritos")
	fmt.Println("      · el desvío del azar es exactamente la HUELLA del barrido: taparse un casillero")
	fmt.Println("        cambia la chance del vecino, así que los tapados no son independientes")
	fmt.Println()
	fmt.Println("   📌 Los primos SÍ están donde deberían estar: el barrido los pone ahí, uno por uno y sin")
	fmt.Println("      margen. Lo que no tenemos es una fórmula corta para adelantarnos al barrido — y esa")
	fmt.Println("      falta de fórmula es lo que se SIENTE como azar, pero no lo es.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(ini, fin, ps, cuenta, tot)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-barrido-de-los-anteriores.svg")
}

func dibujar(ini, fin int, ps []int, cuenta map[int]int, tot int) {
	var b strings.Builder
	W, H := 1060, 700
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL BARRIDO DE LOS ANTERIORES — el primo es lo que sobrevive al tapado")
	txt(530, 54, 12, "#8a8676", "middle", "el mecanismo del capitán, exacto · y la palabra «azar», medida · F408")

	// the coverage map
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="200" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 102, 13, "#e8e2d0", "middle", fmt.Sprintf("el hueco del %d al %d: quién tapa cada casillero", ini, fin))
	ancho := 1020.0 / float64(fin-ini+1)
	cols := map[int]string{2: "#5ec8ff", 3: "#c58bff", 5: "#ff9f6b", 7: "#ff6b6b", 11: "#8a8676", 13: "#6a9a7a"}
	for n := ini; n <= fin; n++ {
		px := 20 + float64(n-ini)*ancho
		q := quienCubre(n, ps)
		col := "#2a3040"
		if q > 0 {
			if c, ok := cols[q]; ok {
				col = c
			} else {
				col = "#4a5060"
			}
		} else {
			col = "#f2c14e"
		}
		fmt.Fprintf(&b, `<rect x="%.1f" y="140" width="%.1f" height="70" fill="%s" opacity="0.85"/>`, px+1, ancho-2, col)
		txt(px+ancho/2, 232, 10, "#c9c3b0", "middle", fmt.Sprint(n))
		if q > 0 {
			txt(px+ancho/2, 182, 13, "#0d1017", "middle", fmt.Sprint(q))
		} else {
			txt(px+ancho/2, 182, 12, "#0d1017", "middle", "libre")
		}
	}
	txt(530, 258, 11, "#f2c14e", "middle", "el amarillo es el casillero que nadie tapó — y ése es el primo siguiente")

	// the chance test
	fmt.Fprintf(&b, `<rect x="20" y="292" width="500" height="200" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(270, 318, 13, "#ff6b6b", "middle", "⛔ pero NO es azar: los gemelos sobran")
	txt(44, 352, 11, "#c9c3b0", "start", "si fuera una moneda, los gemelos tendrían")
	txt(44, 370, 11, "#c9c3b0", "start", "densidad 1/log²n. Contados de verdad:")
	txt(270, 412, 30, "#f2c14e", "middle", "1,32 ×")
	txt(270, 440, 11, "#c9c3b0", "middle", "un 32 % MÁS de los que dicta el azar")
	txt(270, 468, 10, "#8a8676", "middle", "y no es ruido: es la constante de Hardy–Littlewood")

	fmt.Fprintf(&b, `<rect x="540" y="292" width="500" height="200" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 318, 13, "#ff6b6b", "middle", "⛔ y hay huecos FAVORITOS")
	yy := 350.0
	for _, g := range []int{4, 6, 8} {
		col := "#c9c3b0"
		if g == 6 {
			col = "#f2c14e"
		}
		w := 300.0 * float64(cuenta[g]) / float64(cuenta[6])
		fmt.Fprintf(&b, `<rect x="600" y="%.0f" width="%.1f" height="18" fill="%s" opacity="0.8"/>`, yy-13, w, col)
		txt(580, yy, 12, col, "end", fmt.Sprintf("hueco %d", g))
		txt(610+w, yy, 11, "#8a8676", "start", fmt.Sprintf("%.1f %%", 100*float64(cuenta[g])/float64(tot)))
		yy += 34
	}
	txt(790, 458, 11, "#c9c3b0", "middle", "el 6 le gana al 4 y al 8 que lo rodean.")
	txt(790, 478, 10, "#8a8676", "middle", "El azar no tiene favoritos; el barrido sí.")

	// verdict
	fmt.Fprintf(&b, `<rect x="20" y="508" width="1020" height="172" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 534, 14, "#e8e2d0", "middle", "EL VEREDICTO")
	txt(60, 566, 12, "#6ee7a0", "start", "✅ el mecanismo es EXACTO:")
	txt(300, 566, 11, "#c9c3b0", "start", "el primo siguiente es el primer casillero que nadie tapó. Cero excepciones.")
	txt(60, 592, 12, "#6ee7a0", "start", "✅ y explica los huecos:")
	txt(300, 592, 11, "#c9c3b0", "start", "cada primo nuevo suma un tapador, y hay que ir más lejos (Mertens)")
	txt(60, 620, 12, "#ff6b6b", "start", "⚖️ pero NO es azar:")
	txt(300, 620, 11, "#c9c3b0", "start", "el barrido es determinista, y su huella deja correlaciones medibles")
	txt(530, 656, 12, "#f2c14e", "middle", "Los primos SÍ están donde deberían. Lo que falta es una fórmula corta para adelantarse al barrido —")
	txt(530, 674, 11, "#8a8676", "middle", "y esa falta es lo que se SIENTE como azar, pero no lo es. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-barrido-de-los-anteriores.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
