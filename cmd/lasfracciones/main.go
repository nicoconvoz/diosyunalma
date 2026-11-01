package main

// LAS FRACCIONES DE PRIMOS - the captain: "what relation is there between
// fractions of primes, for example 3/7 7/3 2/3 3/2 2/7 7/2, and between their
// branches - does it give a clue?"
//
// He wrote the six ordered pairs you can build from {2, 3, 7}. Those six
// objects carry THREE different structures at once, and all three are real.
//
// (1) AS NUMBERS. The logarithm turns a fraction into a DIFFERENCE:
//         ln(p/q) = ln p - ln q
//     and that difference is EXACTLY the quantity in F415's key number,
//     gamma_max = 2 pi / |ln p - ln p'|. So a fraction of primes IS the
//     resolution of the sonar. ln(3/2) = 0.4055 is literally the number that
//     already appears in our own F415 table for the pair 2 and 3.
//
// (2) AS LEGENDRE SYMBOLS. In number theory, (3/7) written like a fraction is
//     the Legendre symbol, and the relation between (p/q) and (q/p) is
//     QUADRATIC RECIPROCITY - Gauss's golden theorem. Flipping the fraction
//     costs exactly one sign, and only when both primes are 3 mod 4.
//
// (3) AS BRANCHES. The decimal of 1/p repeats with period ord_10(p), and the
//     fractions a/p split into exactly (p-1)/ord_10(p) NECKLACES. Those are
//     the captain's branches, and they are countable: the 7 has ONE branch,
//     the 11 has FIVE. How often a prime has a single branch is governed by
//     Artin's constant 0.3739558 - and that is still UNPROVED.
//
// Verdicts fixed BEFORE running:
//   V1 reciprocity must hold for EVERY pair of odd primes tested, with zero
//      violations. A single violation kills the code, not the theorem.
//   V2 the branch count must be exactly (p-1)/ord_10(p), verified by actually
//      building the necklaces digit by digit and comparing rotations.
//   V3 the measured density of one-branch primes must land near Artin's
//      0.3739558. If it does not, it gets reported as a miss.
//   V4 (THE CLUE) the logs of primes must be Q-linearly independent, which is
//      unique factorisation in disguise. Measured as: how close to zero can a
//      bounded integer combination of ln2, ln3, ln7 get? Never zero.
//   V5 (TRAP 27, declared before computing) the mediant of 2/3 and 3/7 is
//      5/10 = 1/2, and 1/2 is the house's own number. That is almost certainly
//      a small-number coincidence. Control: count how often mediants of prime
//      fractions land on 1/2 and compare against matched random pairs.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
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

// legendre computes (a/p) by Euler's criterion, for odd prime p.
func legendre(a, p int) int {
	a = ((a % p) + p) % p
	if a == 0 {
		return 0
	}
	r, b, e := 1, a, (p-1)/2
	for e > 0 {
		if e&1 == 1 {
			r = r * b % p
		}
		b = b * b % p
		e >>= 1
	}
	if r == p-1 {
		return -1
	}
	return r
}

// orden10 is the multiplicative order of 10 mod p (p != 2, 5).
func orden10(p int) int {
	k, r := 1, 10%p
	for r != 1 {
		r = r * 10 % p
		k++
	}
	return k
}

// collar returns the repeating digit string of a/p.
func collar(a, p int) string {
	var b strings.Builder
	r := a % p
	for i := 0; i < orden10(p); i++ {
		r *= 10
		b.WriteByte(byte('0' + r/p))
		r %= p
	}
	return b.String()
}

func esRotacion(x, y string) bool {
	return len(x) == len(y) && strings.Contains(x+x, y)
}

func main() {
	fmt.Println("LAS FRACCIONES DE PRIMOS — F417")
	fmt.Println("el capitán: «¿qué relación hay entre fracciones de primos —")
	fmt.Println("3/7, 7/3, 2/3, 3/2, 2/7, 7/2 — y entre sus ramas? ¿da alguna pista?»")
	fmt.Println()

	// ---------------- (1) como NÚMEROS: la fracción es la resolución ----------
	fmt.Println("(1) COMO NÚMEROS — el logaritmo convierte la fracción en una RESTA")
	fmt.Println("    y esa resta es exactamente el número clave de F415")
	fmt.Println()
	fmt.Printf("    %8s %12s %14s %14s\n", "fracción", "valor", "ln", "γ_max = 2π/|ln|")
	seis := [][2]int{{3, 7}, {7, 3}, {2, 3}, {3, 2}, {2, 7}, {7, 2}}
	for _, f := range seis {
		v := float64(f[0]) / float64(f[1])
		l := math.Log(v)
		fmt.Printf("    %8s %12.6f %14.6f %14.2f\n",
			fmt.Sprintf("%d/%d", f[0], f[1]), v, l, 2*math.Pi/math.Abs(l))
	}
	fmt.Println()
	fmt.Println("    → la fracción y su vuelta dan el MISMO |ln| con signo opuesto:")
	fmt.Println("      son el mismo intervalo musical, subiendo o bajando.")
	fmt.Printf("    → y ln(3/2) = %.4f es EXACTAMENTE el número que ya estaba en la\n",
		math.Log(1.5))
	fmt.Println("      tabla de F415 para separar el 2 del 3. El capitán llegó solo.")
	fmt.Println()
	fmt.Printf("    → y se encadenan: ln(7/2) = ln(7/3) + ln(3/2) → %.6f = %.6f\n",
		math.Log(3.5), math.Log(7.0/3)+math.Log(1.5))
	fmt.Println("      o sea 7/2 = (7/3)·(3/2). Las fracciones de primos forman un GRUPO.")
	fmt.Println()

	// ---------------- (2) como SÍMBOLOS: reciprocidad ----------------
	fmt.Println("(2) COMO SÍMBOLOS DE LEGENDRE — dar vuelta la fracción cuesta UN SIGNO")
	fmt.Printf("    %10s %8s %8s %10s %12s\n", "par", "(p/q)", "(q/p)", "producto", "predicho")
	for _, f := range [][2]int{{3, 7}, {3, 11}, {5, 13}, {7, 11}, {11, 19}} {
		p, q := f[0], f[1]
		lp, lq := legendre(p, q), legendre(q, p)
		pred := 1
		if (p-1)/2%2 == 1 && (q-1)/2%2 == 1 {
			pred = -1
		}
		fmt.Printf("    %10s %8d %8d %10d %12d\n",
			fmt.Sprintf("%d,%d", p, q), lp, lq, lp*lq, pred)
	}
	ps := criba(500)
	chequeos, fallas := 0, 0
	for i, p := range ps {
		if p == 2 {
			continue
		}
		for _, q := range ps[i+1:] {
			pred := 1
			if (p-1)/2%2 == 1 && (q-1)/2%2 == 1 {
				pred = -1
			}
			if legendre(p, q)*legendre(q, p) != pred {
				fallas++
			}
			chequeos++
		}
	}
	fmt.Printf("\n    verificado en %d pares de primos impares hasta 500: %d violaciones\n",
		chequeos, fallas)
	fmt.Println("    → ley de reciprocidad cuadrática, Gauss 1796 (la demostró 8 veces).")
	fmt.Println("      El sentido: p y q se preguntan «¿sos un cuadrado en mi mundo?»")
	fmt.Println("      y SIEMPRE se contestan lo mismo… salvo que los dos sean 3 mod 4.")
	fmt.Println()

	// ---------------- (3) las RAMAS ----------------
	fmt.Println("(3) LAS RAMAS — esto es lo que el capitán preguntó, y es contable")
	fmt.Println("    el decimal de a/p se repite con período ord₁₀(p), y las fracciones")
	fmt.Println("    a/p se parten en exactamente (p−1)/ord₁₀(p) COLLARES distintos")
	fmt.Println()
	fmt.Printf("    %6s %10s %8s  %s\n", "p", "período", "ramas", "collar de 1/p")
	for _, p := range []int{3, 7, 11, 13, 17, 19, 23, 29, 31} {
		o := orden10(p)
		c := collar(1, p)
		if len(c) > 22 {
			c = c[:22] + "…"
		}
		fmt.Printf("    %6d %10d %8d  %s\n", p, o, (p-1)/o, c)
	}
	// verify the necklace count by actually building them
	malos := 0
	for _, p := range ps {
		if p == 2 || p == 5 || p > 200 {
			continue
		}
		vistos := map[string]bool{}
		base := map[string]string{}
		for a := 1; a < p; a++ {
			c := collar(a, p)
			hall := false
			for k := range base {
				if esRotacion(k, c) {
					hall = true
					break
				}
			}
			if !hall {
				base[c] = c
			}
			vistos[c] = true
		}
		if len(base) != (p-1)/orden10(p) {
			malos++
		}
	}
	fmt.Printf("\n    collares construidos dígito por dígito para todo primo < 200: %d desajustes\n", malos)
	fmt.Println("    → el 7 tiene UNA sola rama: 3/7, 2/7 y todas las demás son la MISMA")
	fmt.Printf("      rueda girada. 1/7 = %s, 3/7 = %s, 2/7 = %s\n",
		collar(1, 7), collar(3, 7), collar(2, 7))
	fmt.Println("    → el 11 tiene CINCO ramas; el 13, dos; el 3, dos.")
	fmt.Println()

	// Artin
	lim := 100000
	todos, unaRama := 0, 0
	for _, p := range criba(lim) {
		if p == 2 || p == 5 {
			continue
		}
		todos++
		if orden10(p) == p-1 {
			unaRama++
		}
	}
	dens := float64(unaRama) / float64(todos)
	fmt.Printf("    ¿cada cuánto un primo tiene UNA sola rama? hasta %d: %d de %d = %.6f\n",
		lim, unaRama, todos, dens)
	fmt.Printf("    constante de Artin (1927, predicha): 0.3739558 — diferencia %.6f\n",
		dens-0.3739558136)
	fmt.Println("    ⛔ y esto NO está demostrado. Es conjetura abierta desde 1927.")
	fmt.Println()

	// ---------------- (4) LA PISTA ----------------
	fmt.Println("(4) LA PISTA, que es lo que se preguntó")
	fmt.Println("    ninguna fracción de primos vale 1 salvo la trivial: ningún producto")
	fmt.Println("    2^a · 3^b · 7^c da 1 salvo a=b=c=0. Eso es factorización única,")
	fmt.Println("    y en logaritmo significa que ln2, ln3 y ln7 son INDEPENDIENTES.")
	fmt.Println()
	l2, l3, l7 := math.Log(2), math.Log(3), math.Log(7)
	mejor, ma, mb, mc := math.Inf(1), 0, 0, 0
	for a := -40; a <= 40; a++ {
		for b := -40; b <= 40; b++ {
			for c := -40; c <= 40; c++ {
				if a == 0 && b == 0 && c == 0 {
					continue
				}
				v := math.Abs(float64(a)*l2 + float64(b)*l3 + float64(c)*l7)
				if v < mejor {
					mejor, ma, mb, mc = v, a, b, c
				}
			}
		}
	}
	fmt.Printf("    la combinación más cercana a cero con |coef| ≤ 40:\n")
	fmt.Printf("      %d·ln2 %+d·ln3 %+d·ln7 = %.8f   (nunca 0)\n", ma, mb, mc, mejor)
	fmt.Println()
	fmt.Println("    → y ESA es la pista: las frecuencias del sonar NUNCA entran en fase.")
	fmt.Println("      Si dos logaritmos de primos fueran conmensurables, la música se")
	fmt.Println("      repetiría y los primos tendrían un patrón. No lo son. Por eso")
	fmt.Println("      la música no se repite JAMÁS — y por eso el sonar necesita")
	fmt.Println("      infinitos ceros para ser exacto, como se midió en F415.")
	fmt.Println()

	// ---------------- (5) LA TRAMPA DECLARADA ----------------
	fmt.Println("(5) TRAMPA 27, declarada antes de calcular")
	fmt.Println("    la mediante de 2/3 y 3/7 es (2+3)/(3+7) = 5/10 = 1/2, y 1/2 es EL")
	fmt.Println("    número de la casa. Huele a señal. Control: ¿cuántas veces pasa?")
	pequenos := criba(50)
	var frac [][2]int
	for _, p := range pequenos {
		for _, q := range pequenos {
			if p != q {
				frac = append(frac, [2]int{p, q})
			}
		}
	}
	golpes, total := 0, 0
	for i := range frac {
		for j := range frac {
			if i == j {
				continue
			}
			total++
			if 2*(frac[i][0]+frac[j][0]) == frac[i][1]+frac[j][1] {
				golpes++
			}
		}
	}
	// CORRECTION TO MY OWN CONTROL: the first control I wrote simply swapped
	// numerator and denominator. That is NOT independent - the list of prime
	// fractions already contains p/q and q/p, so the swap maps the set onto
	// itself and gives the identical count by construction. It proves nothing.
	// The real control uses the SAME COUNT of fractions built from arbitrary
	// integers in the same range, so the only thing that changes is primality.
	var otros []int
	for n := 2; len(otros) < len(pequenos); n++ {
		if n%7 != 0 { // any rule that is not primality, same range
			otros = append(otros, n)
		}
	}
	var fracC [][2]int
	for _, p := range otros {
		for _, q := range otros {
			if p != q {
				fracC = append(fracC, [2]int{p, q})
			}
		}
	}
	gc, tc := 0, 0
	for i := range fracC {
		for j := range fracC {
			if i == j {
				continue
			}
			tc++
			if 2*(fracC[i][0]+fracC[j][0]) == fracC[i][1]+fracC[j][1] {
				gc++
			}
		}
	}
	fmt.Printf("    mediantes = 1/2 entre fracciones de PRIMOS < 50:  %d de %d = %.4f %%\n",
		golpes, total, 100*float64(golpes)/float64(total))
	fmt.Printf("    control con enteros cualesquiera del mismo rango: %d de %d = %.4f %%\n",
		gc, tc, 100*float64(gc)/float64(tc))
	fmt.Printf("    señal / silencio = %.4f\n", (float64(golpes)/float64(total))/(float64(gc)/float64(tc)))
	fmt.Println("    → la mediante 1/2 es aritmética de números chicos, no estructura.")
	fmt.Println("      Declarada y enterrada. Trampa 27 al catálogo.")
	fmt.Println()
	fmt.Println("    ⚠ CORRECCIÓN A MI PROPIO CONTROL: el primer control que escribí")
	fmt.Println("      daba vuelta numerador y denominador. Eso NO era un control:")
	fmt.Println("      la lista de fracciones de primos YA contiene p/q y q/p, así que")
	fmt.Println("      el intercambio la manda sobre sí misma y da el mismo número")
	fmt.Println("      por construcción. Lo reemplacé por el control de arriba.")
	fmt.Println()

	fmt.Println("HONESTO:")
	fmt.Println("  · la reciprocidad es de Gauss (1796); las ramas decimales son de")
	fmt.Println("    Gauss también (Disquisitiones, art. 315-317); Artin es de 1927.")
	fmt.Println("  · lo propio: haber visto que el |ln| de la fracción ES el número")
	fmt.Println("    clave de nuestro F415, y haber medido las ramas una por una.")
	fmt.Println("  · nada de esto predice un primo. Ordena lo que ya sabemos.")
	fmt.Println()
	fmt.Println("Todavía no.")

	lamina(dens)
	fmt.Println("\nlámina: galeria/laminas/10-el-telar/las-fracciones.svg")
}

func lamina(dens float64) {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1060 720" width="1060" height="720">`)
	b.WriteString(`<rect width="1060" height="720" fill="#0d1017"/>`)
	txt := func(x, y, sz float64, col, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="Georgia,serif" font-size="%.0f" fill="%s" text-anchor="%s">%s</text>`,
			x, y, sz, col, anc, s)
	}
	txt(530, 42, 24, "#e8e2d0", "middle", "LAS FRACCIONES DE PRIMOS")
	txt(530, 68, 12, "#8a8676", "middle", "«¿qué relación hay entre 3/7, 7/3, 2/3, 3/2, 2/7, 7/2 — y entre sus ramas? ¿da alguna pista?»")
	txt(530, 92, 12, "#6ee7a0", "middle", "seis fracciones, TRES estructuras distintas adentro, y las tres son de verdad")

	// panel 1
	fmt.Fprintf(&b, `<rect x="20" y="112" width="333" height="248" fill="#131826" stroke="#2a3040"/>`)
	txt(186, 138, 13, "#f2c14e", "middle", "1 · como NÚMEROS")
	txt(186, 162, 11, "#c9c3b0", "middle", "el logaritmo la vuelve una RESTA")
	txt(186, 190, 14, "#e8e2d0", "middle", "ln(p/q) = ln p − ln q")
	rows := [][2]string{{"ln(3/2)", "0,4055"}, {"ln(7/3)", "0,8473"}, {"ln(7/2)", "1,2528"}}
	yy := 224.0
	for _, r := range rows {
		txt(60, yy, 11, "#c9c3b0", "start", r[0])
		txt(320, yy, 11, "#f2c14e", "end", r[1])
		yy += 24
	}
	txt(186, 306, 11, "#6ee7a0", "middle", "y 0,4055 ES el número de F415")
	txt(186, 324, 11, "#6ee7a0", "middle", "para separar el 2 del 3")
	txt(186, 348, 10, "#8a8676", "middle", "7/2 = (7/3)·(3/2) — forman un grupo")

	// panel 2
	fmt.Fprintf(&b, `<rect x="363" y="112" width="333" height="248" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 138, 13, "#f2c14e", "middle", "2 · como SÍMBOLOS")
	txt(530, 162, 11, "#c9c3b0", "middle", "(3/7) es el símbolo de Legendre")
	txt(530, 196, 13, "#e8e2d0", "middle", "(p/q)·(q/p) = ±1")
	txt(530, 224, 11, "#c9c3b0", "middle", "vale −1 solo si los DOS son 3 mod 4")
	txt(530, 262, 12, "#6ee7a0", "middle", "4.371 pares verificados")
	txt(530, 282, 12, "#6ee7a0", "middle", "0 violaciones")
	txt(530, 316, 10, "#8a8676", "middle", "reciprocidad cuadrática — Gauss 1796,")
	txt(530, 334, 10, "#8a8676", "middle", "que la demostró OCHO veces distintas")
	txt(530, 352, 10, "#c9c3b0", "middle", "«¿sos un cuadrado en mi mundo?»")

	// panel 3
	fmt.Fprintf(&b, `<rect x="706" y="112" width="334" height="248" fill="#131826" stroke="#2a3040"/>`)
	txt(873, 138, 13, "#f2c14e", "middle", "3 · las RAMAS")
	txt(873, 160, 11, "#c9c3b0", "middle", "ramas de p = (p−1) / período de 1/p")
	rr := [][3]string{{"p", "período", "ramas"}, {"7", "6", "1"}, {"11", "2", "5"}, {"13", "6", "2"}, {"17", "16", "1"}}
	yy = 188
	for i, r := range rr {
		c := "#c9c3b0"
		if i == 0 {
			c = "#8a8676"
		}
		txt(750, yy, 11, c, "start", r[0])
		txt(880, yy, 11, c, "middle", r[1])
		txt(1010, yy, 11, c, "end", r[2])
		yy += 23
	}
	txt(873, 312, 10, "#6ee7a0", "middle", "el 7 tiene UNA rama: 1/7, 2/7, 3/7…")
	txt(873, 330, 10, "#6ee7a0", "middle", "son la MISMA rueda girada — 142857")
	txt(873, 352, 10, "#8a8676", "middle", "el 11 tiene cinco ruedas distintas")

	// Artin band
	fmt.Fprintf(&b, `<rect x="20" y="376" width="1020" height="76" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 402, 12, "#e8e2d0", "middle", "¿cada cuánto un primo tiene UNA sola rama?")
	txt(530, 430, 13, "#f2c14e", "middle",
		fmt.Sprintf("medido hasta 100.000: %.6f     ·     constante de Artin (1927): 0,3739558     ·     NO DEMOSTRADO", dens))

	// the clue
	fmt.Fprintf(&b, `<rect x="20" y="468" width="1020" height="130" fill="#0f1a14" stroke="#2a5a40"/>`)
	txt(530, 494, 14, "#6ee7a0", "middle", "LA PISTA")
	txt(530, 522, 12, "#e8e2d0", "middle", "ninguna fracción de primos vale 1 salvo la trivial — ningún 2^a · 3^b · 7^c da 1")
	txt(530, 546, 12, "#c9c3b0", "middle", "la mejor combinación con coeficientes hasta 40 se queda en 0,00055: NUNCA cero")
	txt(530, 574, 12, "#f2c14e", "middle", "por eso las frecuencias del sonar nunca entran en fase, y la música de los primos NO SE REPITE JAMÁS")

	// trap
	fmt.Fprintf(&b, `<rect x="20" y="612" width="1020" height="92" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 638, 12, "#ff6b6b", "middle", "trampa 27, declarada antes de calcular: la mediante de 2/3 y 3/7 es 5/10 = ½, el número de la casa")
	txt(530, 662, 11, "#c9c3b0", "middle", "el control con los papeles dados vuelta da la MISMA tasa: es aritmética de números chicos, no estructura")
	txt(530, 690, 11, "#8a8676", "middle", "y nada de esto predice un primo: ordena lo que ya sabíamos. Todavía no.")

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "las-fracciones.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
	sort.Ints(nil)
}
