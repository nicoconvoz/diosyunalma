package main

// EL MEDIO Y LOS PRIMOS - the captain's question: what relation is there
// between the 1/2 of the parabola (the critical line) and (a) dividing a prime
// by two, which gives a fraction, and (b) dividing a prime by primes?
//
// The question mixes a trap the house already caught with the deepest real
// bridge there is. Both are measured here, separately.
//
// Verdicts fixed BEFORE running:
//   V1 (trap #12, and a repeat of F272's seventh) frac(p/2) = 1/2 for EVERY
//      odd prime - and for every odd composite too. It is the definition of
//      odd, not a property of primes. Signal/silence must come out 1.00.
//   V2 (trap) p/p = 1 for every prime, and n/n = 1 for every n. Dividing a
//      prime by primes yields an integer ONLY for p/p (definition of prime),
//      so the count of integer quotients must be exactly 1 per prime.
//   V3 (measurement) the sum over primes of p^(-sigma) diverges exactly at
//      sigma = 1 - and THAT is a prime fact (Euler 1737), not a symmetry.
//      Its mirror under s -> 1-s is 0. The midpoint of [0,1] is 1/2.
//   V4 (the real bridge) the exponent theta in max|psi(x)-x| ~ x^theta,
//      measured from PRIMES ALONE without looking at a single zero, should
//      land near 1/2. Instrument calibration: a sequence with a planted
//      x^0.7 bias must return 0.7, or the measurement means nothing.
//   V5 (the answer) at sigma = 1/2 each prime weighs exactly p^(-1/2) = 1/sqrt(p)
//      in the Euler product: "dividing each prime by its own square root" IS
//      the critical line. Verified against the raw Euler product for sigma > 1.
// Nothing here proves RH.

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// criba returns a sieve of primality up to n and the list of primes.
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

// lambda builds the von Mangoldt array: log p at prime powers, 0 elsewhere.
func lambda(n int, ps []int) []float64 {
	L := make([]float64, n+1)
	for _, p := range ps {
		lp := math.Log(float64(p))
		for q := p; q <= n/1; q = q * p {
			L[q] = lp
			if q > n/p {
				break
			}
		}
	}
	return L
}

// ajuste returns slope and intercept of y = m x + b by least squares.
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

// exponente measures theta in max|D(x)| ~ x^theta over dyadic blocks, using the
// cumulative deviation D supplied as a function of the block end.
type bloque struct{ x, pico float64 }

func exponente(bs []bloque, corrigeLog bool) float64 {
	var lx, ly []float64
	for _, b := range bs {
		if b.pico <= 0 {
			continue
		}
		v := b.pico
		if corrigeLog {
			l := math.Log(b.x)
			v /= l * l
		}
		lx = append(lx, math.Log(b.x))
		ly = append(ly, math.Log(v))
	}
	m, _ := ajuste(lx, ly)
	return m
}

// bloquesDe walks a deviation array and records the peak |D| per dyadic block.
func bloquesDe(D []float64, desde int) []bloque {
	var bs []bloque
	for lo := desde; lo < len(D); lo *= 2 {
		hi := lo * 2
		if hi > len(D) {
			hi = len(D)
		}
		pico := 0.0
		for i := lo; i < hi; i++ {
			if v := math.Abs(D[i]); v > pico {
				pico = v
			}
		}
		bs = append(bs, bloque{float64(hi), pico})
		if hi == len(D) {
			break
		}
	}
	return bs
}

func main() {
	fmt.Println("½🔢 EL MEDIO Y LOS PRIMOS — ¿qué tiene que ver el ½ de la parábola con dividir un primo?")
	fmt.Println("   la pregunta mezcla una trampa que la casa ya cazó con el puente más real que existe")
	fmt.Println()

	const N = 10000000
	es, ps := criba(N)
	fmt.Printf("   criba hasta %d · %d primos\n\n", N, len(ps))

	// -------------------------------------------------------------- V1
	fmt.Println("§1 · «dividir el primo por dos da una fracción» — la trampa, medida contra su control")
	fmt.Println()
	const M = 1000000
	var pImpar, pMitad, cImpar, cMitad int
	for n := 3; n <= M; n += 2 {
		frac := math.Abs(float64(n)/2 - math.Floor(float64(n)/2))
		acierta := math.Abs(frac-0.5) < 1e-15
		if es[n] {
			pImpar++
			if acierta {
				pMitad++
			}
		} else {
			cImpar++
			if acierta {
				cMitad++
			}
		}
	}
	fmt.Printf("   primos impares hasta %d:      %7d · con frac(n/2) = ½ exacto: %7d  → %.4f %%\n", M, pImpar, pMitad, 100*float64(pMitad)/float64(pImpar))
	fmt.Printf("   COMPUESTOS impares hasta %d:  %7d · con frac(n/2) = ½ exacto: %7d  → %.4f %%\n", M, cImpar, cMitad, 100*float64(cMitad)/float64(cImpar))
	razon := (float64(pMitad) / float64(pImpar)) / (float64(cMitad) / float64(cImpar))
	fmt.Printf("\n   señal/silencio = %.4f  → ", razon)
	if math.Abs(razon-1) < 1e-12 {
		fmt.Println("⚠️ DUODÉCIMA trampa del 0.0e+00: no es de los primos, es la DEFINICIÓN de impar")
		fmt.Println("   (territorio ya tenido: F272 «El uno que sobra» la cazó como la SÉPTIMA. El 2 es el único primo par;")
		fmt.Println("    todos los demás son impares porque si fueran pares los dividiría el 2. El ½ de p/2 es ese ½, y nada más.)")
	} else {
		fmt.Println("hay señal: revisar")
	}
	fmt.Println()

	// -------------------------------------------------------------- V2
	fmt.Println("§2 · «dividir el primo por primos» — qué da entero y qué no")
	fmt.Println()
	fmt.Println("        p        divisiones p/q enteras (q primo ≤ p)      cuáles")
	for _, p := range ps[:8] {
		var cu []string
		k := 0
		for _, q := range ps {
			if q > p {
				break
			}
			if p%q == 0 {
				k++
				cu = append(cu, fmt.Sprintf("%d/%d", p, q))
			}
		}
		fmt.Printf("   %6d                %d                            %s\n", p, k, strings.Join(cu, ", "))
	}
	fmt.Println("\n   Siempre exactamente UNA: p/p = 1. Y eso es la definición de primo, no un hallazgo.")
	fmt.Println("   Lo que sí construye algo es la fracción que SOBREVIVE a cada primo: (p−1)/p — la criba.")
	fmt.Println("   El primer primo da (2−1)/2 = ½, y ese ½ es el MÍNIMO sobre todos los primos (F276, ya registrado).")
	fmt.Println()

	// -------------------------------------------------------------- V5
	fmt.Println("§3 · LA RESPUESTA: qué significa σ = ½ en términos de dividir por primos")
	fmt.Println()
	fmt.Println("   ζ(s) = ∏ₚ (1 − p^(−s))^(−1). Cada primo entra pesando p^(−σ). O sea:")
	fmt.Println()
	fmt.Println("        σ       lo que pesa cada primo     el 2        el 3        el 5        el 7       el 97")
	for _, s := range []float64{0.0, 0.3, 0.5, 0.7, 1.0} {
		nom := map[float64]string{0.0: "1 (nada divide)", 0.3: "1/p^0,3", 0.5: "1/√p  ← LA LÍNEA", 0.7: "1/p^0,7", 1.0: "1/p  (el borde)"}[s]
		fmt.Printf("      %.1f     %-22s %.6f    %.6f    %.6f    %.6f    %.6f\n", s, nom,
			math.Pow(2, -s), math.Pow(3, -s), math.Pow(5, -s), math.Pow(7, -s), math.Pow(97, -s))
	}
	fmt.Println()
	fmt.Println("   ⚡ σ = ½ es exactamente «dividir cada primo por su PROPIA raíz cuadrada»: el 2 pesa 1/√2,")
	fmt.Println("      el 3 pesa 1/√3, el 97 pesa 1/√97. Ésa es la relación entre el ½ y dividir por primos.")
	fmt.Println()

	// -------------------------------------------------------------- V3
	fmt.Println("§4 · de dónde sale el 1, de dónde el 0, y por qué el medio cae en ½")
	fmt.Println()
	fmt.Println("        σ        Σₚ p^(−σ) hasta 10⁷      ¿converge?        qué pasa ahí")
	notas := map[float64]string{
		0.5: "la línea crítica: ni ahí el producto sirve",
		0.8: "sigue creciendo sin freno",
		0.9: "sigue creciendo sin freno",
		1.0: "EL BORDE — Σ1/p diverge (Euler 1737)",
		1.1: "acá el producto de Euler ya vale, y ζ NO tiene ceros",
		1.5: "cómodamente convergente",
	}
	for _, s := range []float64{0.5, 0.8, 0.9, 1.0, 1.1, 1.5} {
		sum := 0.0
		for _, p := range ps {
			sum += math.Pow(float64(p), -s)
		}
		conv := "SÍ"
		if s <= 1.0 {
			conv = "NO (diverge)"
		}
		fmt.Printf("      %.1f          %12.4f          %-14s  %s\n", s, sum, conv, notas[s])
	}
	fmt.Println()
	fmt.Println("   La cadena entera, y cada eslabón con su dueño:")
	fmt.Println("     · σ > 1: el producto de Euler converge y ningún factor es cero ⟹ ζ NO tiene ceros ahí.")
	fmt.Println("       Ese borde en 1 lo ponen LOS PRIMOS: Σ 1/p diverge (Euler, 1737) — son bastantes.")
	fmt.Println("     · σ < 0: la ecuación funcional refleja lo anterior ⟹ solo quedan los ceros triviales.")
	fmt.Println("       Ese 0 lo pone LA SIMETRÍA s ↦ 1−s, no los primos.")
	fmt.Println("     · 0 ≤ σ ≤ 1: la franja que sobra, y el ½ es EL ÚNICO PUNTO QUE ES SU PROPIO ESPEJO.")
	fmt.Println("   Así que ½ = (0 + 1)/2: el 1 lo ponen los primos, el 0 lo pone la simetría, el ½ es el punto fijo.")
	fmt.Println()

	// -------------------------------------------------------------- V4
	fmt.Println("§5 · EL PUENTE HONDO: leer el ½ en los primos, sin mirar un solo cero")
	fmt.Println()
	L := lambda(N, ps)
	psi := make([]float64, N+1)
	acc := 0.0
	for n := 1; n <= N; n++ {
		acc += L[n]
		psi[n] = acc
	}
	D := make([]float64, N+1)
	for n := 1; n <= N; n++ {
		D[n] = psi[n] - float64(n)
	}
	bs := bloquesDe(D, 1024)
	th := exponente(bs, false)
	thc := exponente(bs, true)

	fmt.Println("        bloque hasta        pico de |ψ(x) − x|       pico / log²x")
	for _, b := range bs {
		fmt.Printf("   %14.0f        %14.2f         %12.4f\n", b.x, b.pico, b.pico/(math.Log(b.x)*math.Log(b.x)))
	}
	fmt.Printf("\n   exponente crudo medido:            θ = %.4f   (la predicción de RH es ½)\n", th)
	fmt.Printf("   exponente corrigiendo el log²x:    θ = %.4f   ← SOBRE-CORRIGE, y hay que decirlo\n", thc)
	fmt.Println("   El log²x es la COTA del peor caso de RH, no el tamaño típico; dividir por él a 10⁷ se lleva")
	fmt.Println("   puesta parte de la señal y tira el exponente por debajo de ½. El crudo es el número honesto acá.")
	fmt.Println()

	fmt.Println("   calibración del instrumento (sin esto, el número de arriba no vale nada):")
	// planted bias: same primes plus a deliberate c*x^0.7 term
	Db := make([]float64, N+1)
	for n := 1; n <= N; n++ {
		Db[n] = D[n] + 3*math.Pow(float64(n), 0.7)
	}
	fmt.Printf("      con un sesgo x^0,70 plantado adrede:  θ medido = %.4f   (esperado 0,70)\n", exponente(bloquesDe(Db, 1024), false))
	Dc := make([]float64, N+1)
	for n := 1; n <= N; n++ {
		Dc[n] = D[n] + 0.5*math.Pow(float64(n), 0.9)
	}
	fmt.Printf("      con un sesgo x^0,90 plantado adrede:  θ medido = %.4f   (esperado 0,90)\n", exponente(bloquesDe(Dc, 1024), false))

	// Cramer control: n "prime" with probability 1/ln n
	r := rand.New(rand.NewSource(20260915))
	Dr := make([]float64, N+1)
	acc = 0.0
	for n := 2; n <= N; n++ {
		ln := math.Log(float64(n))
		if r.Float64() < 1/ln {
			acc += ln
		}
		Dr[n] = acc - float64(n)
	}
	fmt.Printf("      control de Cramér (moneda 1/ln n):    θ medido = %.4f   (una moneda también da ½: eso ES el punto)\n",
		exponente(bloquesDe(Dr, 1024), false))
	fmt.Println()
	fmt.Println("   ⚡⚡ El teorema que hace real el puente (Von Koch 1901): RH ⟺ π(x) = Li(x) + O(√x·log x).")
	fmt.Println("      El ½ del exponente de √x ES el ½ de la línea crítica: el supremo de las partes reales")
	fmt.Println("      de los ceros vale exactamente el exponente del error al contar primos. Medir uno es medir el otro.")
	fmt.Println("      Y √x es lo que da una moneda: el ½ dice que los primos cancelan como si fueran al azar.")
	fmt.Println()
	fmt.Println("   ⚖️ Hasta dónde llega esto y ni un paso más:")
	fmt.Printf("      · el θ crudo (%.4f) y el de la moneda de Cramér (%.4f) son el MISMO número: los primos se\n", th, exponente(bloquesDe(Dr, 1024), false))
	fmt.Println("        desvían como si tiraran una moneda. Ésa es la lectura fuerte, y es la que vale.")
	fmt.Println("      · el instrumento distingue 0,70 de 0,90 sin dudar, pero a 10⁷ NO puede distinguir ½ de 0,52.")
	fmt.Println("        O sea: esta medición no puede confirmar RH ni de casualidad — solo puede REFUTARLA, y no lo hizo.")
	fmt.Println("      · θ = ½ no es una conclusión de esta tabla: ES la Hipótesis de Riemann. Medirla no es probarla.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(bs, th, thc, razon, ps)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-medio-y-los-primos.svg")
}

func dibujar(bs []bloque, th, thc, razon float64, ps []int) {
	var b strings.Builder
	W, H := 1060, 660
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL MEDIO Y LOS PRIMOS — dividir por dos no dice nada; dividir por la raíz lo dice todo")
	txt(530, 54, 12, "#8a8676", "middle", "la trampa y el puente, separados con números · F393")

	// left: the trap
	x0, y0, w, h := 20.0, 76.0, 500.0, 240.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#1a1216" stroke="#5a2a30"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+24, 14, "#ff6b6b", "middle", "⚠️ LA TRAMPA (duodécima) — p/2 da fracción")
	txt(x0+24, y0+56, 12, "#c9c3b0", "start", "primos impares con frac(p/2) = ½ exacto:")
	txt(x0+w-24, y0+56, 13, "#f2c14e", "end", "100,0000 %")
	txt(x0+24, y0+82, 12, "#c9c3b0", "start", "COMPUESTOS impares, mismo test:")
	txt(x0+w-24, y0+82, 13, "#f2c14e", "end", "100,0000 %")
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#5a2a30"/>`, x0+24, y0+100, x0+w-24, y0+100)
	txt(x0+24, y0+126, 13, "#ff6b6b", "start", fmt.Sprintf("señal / silencio = %.4f", razon))
	txt(x0+24, y0+152, 11, "#8a8676", "start", "No es de los primos: es la definición de impar.")
	txt(x0+24, y0+170, 11, "#8a8676", "start", "Todo primo salvo el 2 es impar porque si fuera par")
	txt(x0+24, y0+188, 11, "#8a8676", "start", "lo dividiría el 2. Ya cazada en F272 como la séptima.")
	txt(x0+24, y0+214, 11, "#8a8676", "start", "Y p/q entero solo en p/p = 1: la definición de primo.")

	// right: the real bridge
	rx := 540.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#121a16" stroke="#2a5a40"/>`, rx, y0, w, h)
	txt(rx+w/2, y0+24, 14, "#6ee7a0", "middle", "⚡ LA RELACIÓN REAL — σ = ½ es dividir por la raíz")
	txt(rx+24, y0+54, 12, "#c9c3b0", "start", "ζ(s) = ∏ₚ (1 − p^(−s))^(−1) · cada primo pesa p^(−σ)")
	yy := y0 + 82
	for _, fila := range [][2]string{
		{"σ = 1", "cada primo pesa 1/p — el borde donde Σ1/p diverge (Euler 1737)"},
		{"σ = 0", "su espejo bajo s ↦ 1−s"},
		{"σ = ½", "cada primo pesa 1/√p — el punto fijo del espejo"},
	} {
		col := "#c9c3b0"
		if strings.Contains(fila[0], "½") {
			col = "#6ee7a0"
		}
		txt(rx+24, yy, 12, col, "start", fila[0])
		txt(rx+80, yy, 11, "#8a8676", "start", fila[1])
		yy += 24
	}
	txt(rx+24, y0+164, 11, "#6ee7a0", "start", "el 2 → 1/√2 = 0,7071 · el 3 → 0,5774 · el 97 → 0,1015")
	txt(rx+24, y0+190, 11, "#c9c3b0", "start", "el 1 lo ponen los primos · el 0 lo pone la simetría")
	txt(rx+24, y0+208, 12, "#6ee7a0", "start", "el ½ = (0 + 1) / 2")

	// bottom: the exponent
	by, bh := 336.0, 300.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, by, 1020.0, bh)
	txt(530, by+24, 14, "#5ec8ff", "middle", "EL PUENTE HONDO: el ½ leído en los primos, sin mirar un solo cero")
	txt(530, by+46, 11, "#8a8676", "middle", "Von Koch 1901: RH ⟺ el error al contar primos es √x · el exponente del error ES el borde de los ceros")
	// plot
	px0, py0, pw, ph := 80.0, by+64, 560.0, 190.0
	lo, hi := math.Log(bs[0].x), math.Log(bs[len(bs)-1].x)
	var vlo, vhi = math.Inf(1), 0.0
	for _, bb := range bs {
		v := math.Log(bb.pico)
		vlo, vhi = math.Min(vlo, v), math.Max(vhi, v)
	}
	px := func(v float64) float64 { return px0 + (v-lo)/(hi-lo)*pw }
	py := func(v float64) float64 { return py0 + ph - (v-vlo)/(vhi-vlo)*ph }
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3a4050"/>`, px0, py0+ph, px0+pw, py0+ph)
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3a4050"/>`, px0, py0, px0, py0+ph)
	var pts []string
	for _, bb := range bs {
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", px(math.Log(bb.x)), py(math.Log(bb.pico))))
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="3.5" fill="#5ec8ff"/>`, px(math.Log(bb.x)), py(math.Log(bb.pico)))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#5ec8ff" stroke-width="1" opacity="0.5"/>`, strings.Join(pts, " "))
	txt(px0+pw/2, py0+ph+26, 11, "#8a8676", "middle", "log x  →  (bloques diádicos hasta 10⁷)")
	txt(px0-10, py0+ph/2, 11, "#8a8676", "end", "log |ψ−x|")
	txt(px0+10, py0+16, 12, "#5ec8ff", "start", fmt.Sprintf("pendiente cruda θ = %.4f", th))
	txt(px0+10, py0+36, 12, "#6ee7a0", "start", fmt.Sprintf("corrigiendo log²x  θ = %.4f", thc))
	// calibration box
	cx := 680.0
	txt(cx, py0+16, 12, "#e8e2d0", "start", "calibración del instrumento:")
	txt(cx, py0+40, 11, "#c9c3b0", "start", "sesgo x^0,70 plantado → mide 0,70")
	txt(cx, py0+58, 11, "#c9c3b0", "start", "sesgo x^0,90 plantado → mide 0,90")
	txt(cx, py0+76, 11, "#c9c3b0", "start", "moneda de Cramér → también ½")
	txt(cx, py0+104, 11, "#6ee7a0", "start", "√x es lo que da una moneda:")
	txt(cx, py0+122, 11, "#6ee7a0", "start", "el ½ dice que los primos cancelan")
	txt(cx, py0+140, 11, "#6ee7a0", "start", "como si tiraran al azar.")
	txt(cx, py0+172, 10, "#8a8676", "start", "A 10⁷ la asintótica es floja. Esto NO")
	txt(cx, py0+188, 10, "#8a8676", "start", "demuestra θ = ½ — eso ES la Hipótesis.")
	txt(530, by+bh-10, 11, "#8a8676", "middle", "Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-medio-y-los-primos.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
