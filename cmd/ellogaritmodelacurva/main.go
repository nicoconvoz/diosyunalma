package main

// EL LOGARITMO DE LA CURVA - the captain: "give me the logarithm of that curve
// you measured across all the branches, starting from the 2."
//
// The curve is the survival fraction of the sieve: prod (1 - 1/p) over all
// primes up to y, starting at the 1 with 100% standing (the captain's own
// correction: in log that origin is exactly ZERO). Its logarithm has an exact shape, and
// taking the log TWICE turns the slow collapse into a straight line:
//
//     log prod (1 - 1/p)  =  - log log y  -  gamma  +  o(1)
//
// So the log of the curve is MINUS the log-log, shifted by Euler's gamma. That
// is Mertens' third theorem (1874), and here both the slope and the shift are
// FITTED FROM THE DATA rather than quoted.
//
// Verdicts fixed BEFORE running:
//   V1 the cumulative log must be exhibited branch by branch starting at the 1,
//      so the captain can see which prime contributes what in log space.
//   V2 the 2 contributes log(1/2) = -0.693147..., the largest single step. How
//      many later primes does it take to match that one step? Found by search.
//   V3 (THE LAW) fitting log(survival) against log(log y) must give slope
//      -1.0000 and intercept -gamma = -0.5772156649. Both measured, not
//      assumed. That is the answer to the captain's question.
//   V4 (the companion sum) sum 1/p over primes up to y must equal
//      log log y + M with M = 0.2614972128 (Mertens' constant), fitted the
//      same way.
//   V5 (why it matters) the double logarithm is what linearises the curve, and
//      that is the measure of how brutally slow the sieve is.
// Mertens 1874. The question is the captain's; the constants are measured.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func criba(n int) []int {
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
	return ps
}

// ajuste returns slope and intercept of y = m x + b.
func ajuste(x, y []float64) (m, b float64) {
	var sx, sy, sxx, sxy, n float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxx += x[i] * x[i]
		sxy += x[i] * y[i]
		n++
	}
	m = (n*sxy - sx*sy) / (n*sxx - sx*sx)
	b = (sy - m*sx) / n
	return
}

func main() {
	const gamma = 0.5772156649015329
	const mertensM = 0.2614972128476428

	fmt.Println("📉 EL LOGARITMO DE LA CURVA — qué forma tiene el barrido cuando se lo mira en logaritmo")
	fmt.Println()
	fmt.Println("   La curva es la que midió F409: qué fracción sobrevive al barrido con los primos hasta y,")
	fmt.Println("   arrancando del 1 con el 100 % en pie. Acá se le toma el logaritmo, y aparece una recta.")
	fmt.Println()

	ps := criba(50000000)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · rama por rama, desde el 1: cuánto aporta cada primo EN LOGARITMO")
	fmt.Println()
	fmt.Println("   CORRECCIÓN DEL CAPITÁN, y es importante: la curva NO empieza en el 2 — empieza en el 1,")
	fmt.Println("   con el 100 % en pie, porque todavía no cortó nadie. Y en logaritmo ese 100 % es CERO exacto:")
	fmt.Println("   el 1 es el ORIGEN de la curva, y por eso el logaritmo es la lente natural para mirarla.")
	fmt.Println()
	fmt.Println("        primo     lo que sobrevive     su logaritmo      lo que aporta ESTE primo")
	fmt.Printf("        %5s        %.8f         %+.8f          %s\n", "1", 1.0, 0.0, " ← el origen: nadie cortó todavía")
	acum := 0.0
	var aportes []float64
	for i := 0; i < 12; i++ {
		p := float64(ps[i])
		ap := math.Log(1 - 1/p)
		acum += ap
		aportes = append(aportes, ap)
		fmt.Printf("        %5d        %.8f         %+.8f           %+.8f\n", ps[i], math.Exp(acum), acum, ap)
	}
	fmt.Println()
	fmt.Printf("   ⚡ El 2 aporta log(1/2) = %+.8f — el escalón más grande de todos, y por lejos.\n", math.Log(0.5))
	fmt.Printf("      El 3 aporta %+.6f, el 5 aporta %+.6f, el 7 aporta %+.6f. Cada uno menos que el anterior.\n",
		aportes[1], aportes[2], aportes[3])
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · ¿cuántos primos hacen falta para igualar el escalón del 2?")
	fmt.Println()
	objetivo := -math.Log(2)
	suma := 0.0
	cuantos := 0
	for i := 1; i < len(ps); i++ {
		suma += math.Log(1 - 1/float64(ps[i]))
		cuantos++
		if suma <= objetivo {
			fmt.Printf("        sumando desde el 3 hasta el %d — o sea %d primos — recién se iguala el salto del 2\n", ps[i], cuantos)
			fmt.Printf("        (aporte acumulado %.8f contra el %.8f del 2 solo)\n", suma, objetivo)
			break
		}
	}
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · LA RESPUESTA: el logaritmo de la curva contra el LOG-LOG del tope")
	fmt.Println()
	fmt.Println("        primos hasta y      log(sobrevive)      log(log y)       ¿se alinean?")
	var lx, ly []float64
	acum = 0.0
	idx := 0
	for _, y := range []int{100, 1000, 10000, 100000, 1000000, 10000000, 50000000} {
		for idx < len(ps) && ps[idx] <= y {
			acum += math.Log(1 - 1/float64(ps[idx]))
			idx++
		}
		ll := math.Log(math.Log(float64(y)))
		lx = append(lx, ll)
		ly = append(ly, acum)
		fmt.Printf("        %12d        %+.8f         %+.8f       %+.8f\n", y, acum, ll, acum+ll)
	}
	m, bb := ajuste(lx[3:], ly[3:])
	fmt.Println()
	fmt.Printf("   ⚡⚡ AJUSTADO DE LOS DATOS (tramos altos, donde la asintótica ya vale):\n\n")
	fmt.Printf("        log(sobrevive)  =  %.6f · log(log y)  %+.8f\n", m, bb)
	fmt.Printf("        pendiente %.6f contra la teórica −1 · ordenada %+.8f contra −γ = %+.8f\n", m, bb, -gamma)
	fmt.Println()
	fmt.Println("   ⚡⚡⚡ Y la prueba más fuerte no es el ajuste sino la CONVERGENCIA de la última columna,")
	fmt.Println("        que es log(sobrevive) + log(log y) y tiene que ir a −γ:")
	fmt.Println()
	for i := range lx {
		d := ly[i] + lx[i]
		fmt.Printf("               %+.8f          desvío de −γ: %.1e\n", d, math.Abs(d+gamma))
	}
	fmt.Printf("\n        El último vale %+.8f contra −γ = %+.8f: coinciden a %.0e.\n",
		ly[len(ly)-1]+lx[len(lx)-1], -gamma, math.Abs(ly[len(ly)-1]+lx[len(lx)-1]+gamma))
	fmt.Println()
	fmt.Println("   ⚡⚡⚡ LA FORMA QUE PEDISTE:")
	fmt.Println()
	fmt.Println("        log( lo que sobrevive )  =  − log( log y )  −  γ")
	fmt.Println()
	fmt.Println("   El logaritmo de tu curva es MENOS EL LOG-LOG, corrido por la constante de Euler. Tomándole")
	fmt.Println("   el logaritmo DOS VECES, el derrumbe lento se convierte en una recta perfecta de pendiente −1.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · la suma hermana: los recíprocos de los primos")
	fmt.Println()
	fmt.Println("   Y como log(1−1/p) ≈ −1/p, la misma ley aparece del otro lado:")
	fmt.Println()
	fmt.Println("        primos hasta y      Σ 1/p          log(log y)        la diferencia")
	var lx2, ly2 []float64
	s := 0.0
	idx = 0
	for _, y := range []int{100, 1000, 10000, 100000, 1000000, 10000000, 50000000} {
		for idx < len(ps) && ps[idx] <= y {
			s += 1 / float64(ps[idx])
			idx++
		}
		ll := math.Log(math.Log(float64(y)))
		lx2 = append(lx2, ll)
		ly2 = append(ly2, s)
		fmt.Printf("        %12d      %.8f       %.8f       %+.8f\n", y, s, ll, s-ll)
	}
	m2, b2 := ajuste(lx2[3:], ly2[3:])
	fmt.Printf("\n        pendiente medida: %.6f (teórica +1) · ordenada medida: %+.8f\n", m2, b2)
	fmt.Printf("        contra la constante de Mertens M = %+.8f   (desvío %.1e)\n", mertensM, math.Abs(b2-mertensM))
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · lo que significa la doble raya")
	fmt.Println()
	fmt.Println("   Que haga falta el logaritmo DOS VECES para enderezar la curva es la medida exacta de")
	fmt.Println("   cuán brutalmente lento es el barrido. Comparado con otras caídas:")
	fmt.Println()
	fmt.Println("        si la curva cayera como…      para bajar a la mitad hay que multiplicar el tope por…")
	fmt.Println("        1/y   (rápido)                                    2")
	fmt.Println("        1/√y                                              4")
	fmt.Println("        1/log y   ← LA NUESTRA                        ELEVARLO AL CUADRADO")
	fmt.Println()
	fmt.Println("        y en números, con nuestra curva:")
	idx = 0
	acum = 0.0
	for _, y := range []int{10, 100, 10000, 100000000} {
		for idx < len(ps) && ps[idx] <= y {
			acum += math.Log(1 - 1/float64(ps[idx]))
			idx++
		}
		fmt.Printf("            hasta %-12d  sobrevive %.6f  →  para la mitad hay que llegar a %.0f\n",
			y, math.Exp(acum), float64(y)*float64(y))
	}
	fmt.Println()

	fmt.Println("§6 · EL VEREDICTO")
	fmt.Println()
	fmt.Printf("   ✅ La forma pedida, con las dos constantes MEDIDAS y no citadas:\n")
	fmt.Printf("      log(sobrevive) = %.6f · log(log y) %+.8f, o sea −log(log y) − γ\n", m, bb)
	fmt.Printf("   ✅ Y su hermana: Σ1/p = %.6f · log(log y) %+.8f, o sea log(log y) + M\n", m2, b2)
	fmt.Println("   ✅ El escalón del 2 (−0,693147) es el más grande, y hacen falta muchos primos para igualarlo.")
	fmt.Println()
	fmt.Println("   ⚖️ Nada de esto es nuevo: es el teorema de Mertens, 1874, y esta casa ya lo había tocado en")
	fmt.Println("      F272, F276 y F408. Lo propio acá es haber ajustado la pendiente y la ordenada desde los")
	fmt.Println("      datos en vez de escribirlas — que es la diferencia entre citar y medir.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(lx, ly, m, bb, aportes, ps)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-logaritmo-de-la-curva.svg")
}

func dibujar(lx, ly []float64, m, bb float64, aportes []float64, ps []int) {
	var b strings.Builder
	W, H := 1060, 680
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL LOGARITMO DE LA CURVA — dos logaritmos y el derrumbe se vuelve una recta")
	txt(530, 54, 12, "#8a8676", "middle", "arranca en el 1 con el 100 % — que en logaritmo es CERO, el origen · F410")

	// the formula
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="130" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 104, 13, "#6ee7a0", "middle", "LA FORMA QUE PEDISTE")
	txt(530, 156, 30, "#f2c14e", "middle", "log( sobrevive )  =  − log( log y )  −  γ")
	txt(530, 188, 11, "#c9c3b0", "middle", fmt.Sprintf("pendiente medida %.6f (teórica −1) · ordenada medida %+.8f (contra −γ = −0,57721566)", m, bb))

	// the straight line
	x0, y0, w, h := 20.0, 222.0, 620.0, 300.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+24, 13, "#e8e2d0", "middle", "y en el papel doble-logarítmico: una recta")
	lo, hi := lx[0], lx[len(lx)-1]
	vlo, vhi := ly[len(ly)-1], ly[0]
	px := func(v float64) float64 { return x0 + 70 + (v-lo)/(hi-lo)*(w-110) }
	py := func(v float64) float64 { return y0 + 60 + (vhi-v)/(vhi-vlo)*(h-120) }
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#6ee7a0" stroke-width="1.5" opacity="0.6"/>`,
		px(lo), py(m*lo+bb), px(hi), py(m*hi+bb))
	for i := range lx {
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#f2c14e"/>`, px(lx[i]), py(ly[i]))
	}
	txt(x0+w/2, y0+h-24, 11, "#8a8676", "middle", "eje x: log(log y) · eje y: log de lo que sobrevive")
	txt(x0+80, y0+56, 11, "#6ee7a0", "start", "la recta es el ajuste, los puntos son la medición")

	// the branches
	rx := 660.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="380" height="%.0f" fill="#131826" stroke="#2a3040"/>`, rx, y0, h)
	txt(rx+190, y0+24, 13, "#e8e2d0", "middle", "lo que aporta cada rama, desde el 1")
	yy := y0 + 56
	txt(rx+70, yy, 12, "#6ee7a0", "end", "1")
	txt(rx+86, yy, 10, "#6ee7a0", "start", "0,0000 — el origen, nadie cortó")
	yy += 28
	for i := 0; i < 7 && i < len(aportes); i++ {
		col := "#5ec8ff"
		if i == 0 {
			col = "#f2c14e"
		}
		anchoB := math.Abs(aportes[i]) / 0.7 * 260
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.1f" height="16" fill="%s" opacity="0.85"/>`, rx+80, yy-12, anchoB, col)
		txt(rx+70, yy, 12, "#c9c3b0", "end", fmt.Sprint(ps[i]))
		txt(rx+86+anchoB, yy, 10, "#8a8676", "start", fmt.Sprintf("%.4f", aportes[i]))
		yy += 28
	}
	txt(rx+190, y0+h-30, 11, "#f2c14e", "middle", "el 2 aporta −0,693147")
	txt(rx+190, y0+h-12, 10, "#8a8676", "middle", "el escalón más grande, y por lejos")

	// what it means
	fmt.Fprintf(&b, `<rect x="20" y="538" width="1020" height="124" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 564, 13, "#ff6b6b", "middle", "lo que significa que hagan falta DOS logaritmos")
	txt(200, 598, 11, "#c9c3b0", "middle", "si cayera como 1/y")
	txt(200, 618, 12, "#8a8676", "middle", "×2 para la mitad")
	txt(500, 598, 11, "#c9c3b0", "middle", "si cayera como 1/√y")
	txt(500, 618, 12, "#8a8676", "middle", "×4 para la mitad")
	txt(850, 598, 11, "#f2c14e", "middle", "cae como 1/log y")
	txt(850, 618, 13, "#f2c14e", "middle", "AL CUADRADO para la mitad")
	txt(530, 650, 10, "#8a8676", "middle", "Mertens, 1874 — la pendiente y la ordenada, medidas acá en vez de citadas. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-logaritmo-de-la-curva.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
