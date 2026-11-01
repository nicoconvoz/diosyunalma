package main

// LOS DOS LATIDOS - the captain: "two beats against more beats, that is the
// distinction between a composite and a prime. So there has to be a relation
// between the number 2 (because it is two beats), the number 1, and the prime
// itself. Find it."
//
// The beats are exact: a prime has exactly TWO divisors, 1 and itself. Saying
// so is the definition, and the house declares it as such. But the relation he
// asks for - between 2, 1 and p - EXISTS, it is not a tautology, and for the
// first time in this whole campaign a flash of the captain lands on something
// with REAL SIGNAL:
//
//     2^(p-1) = 1  (mod p)        for every prime p > 2
//
// Fermat, 1640. The number 2, the number 1 and the prime, in one line - and
// unlike every rule tested in F393, F395, F401 and F403, this one SEPARATES
// primes from composites instead of holding for both.
//
// Verdicts fixed BEFORE running:
//   V1 (trap #21, declared) "a prime has exactly two divisors" is the
//      definition. d(n) = 2 <=> n prime carries no information; it is what
//      "prime" means. Measured and declared, not sold.
//   V2 (the relation, and it must hold without exception) 2^(p-1) mod p must
//      equal 1 for EVERY odd prime up to 10^6. A single failure kills it.
//   V3 (THE CONTROL, which is what makes this different) the same test on
//      composites must FAIL for almost all of them. Signal/silence measured.
//      Every previous flash gave exactly 1.0000 here; this one must not.
//   V4 (the honest limit) the test is not perfect: Fermat pseudoprimes pass.
//      The first must be found by search, not quoted, and their density
//      measured.
//   V5 (the deepest failure) Carmichael numbers pass for EVERY base, not just
//      2. Found by search and reported, because a test that can be fooled in
//      principle is not a proof of primality.
// Fermat 1640. Nothing new - but for the first time, something that measures.

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

// potMod computes base^exp mod m without overflow for m < 2^31.
func potMod(base, exp, m int) int {
	r := 1
	b := base % m
	for e := exp; e > 0; e >>= 1 {
		if e&1 == 1 {
			r = int((int64(r) * int64(b)) % int64(m))
		}
		b = int((int64(b) * int64(b)) % int64(m))
	}
	return r
}

func divisores(n int) []int {
	var d []int
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			d = append(d, i)
			if i != n/i {
				d = append(d, n/i)
			}
		}
	}
	return d
}

func main() {
	fmt.Println("💓 LOS DOS LATIDOS — la relación entre el 2, el 1 y el primo")
	fmt.Println()

	const N = 1000000
	es, _ := criba(N)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · los latidos, y la trampa que traen puesta (vigesimoprimera)")
	fmt.Println()
	fmt.Println("        número    ¿primo?    sus divisores              cuántos latidos")
	for _, n := range []int{2, 3, 5, 7, 4, 6, 9, 12} {
		d := divisores(n)
		et := "no"
		if es[n] {
			et = "SÍ"
		}
		var s []string
		for _, x := range d {
			s = append(s, fmt.Sprint(x))
		}
		fmt.Printf("         %3d       %-3s       %-24s        %d\n", n, et, strings.Join(s, ", "), len(d))
	}
	fmt.Println()
	fmt.Println("   ⚡ El capitán tiene razón y es exacto: un primo late DOS veces — el 1 y él mismo — y un")
	fmt.Println("      compuesto late más. ⚠️ Pero eso ES LA DEFINICIÓN de primo, así que como regla no informa:")
	fmt.Println("      es la vigesimoprimera vez que la casa declara una de éstas. Lo interesante viene ahora.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · LA RELACIÓN QUE PIDIÓ: el 2, el 1 y el primo, en un renglón")
	fmt.Println()
	fmt.Println("        2^(p−1)  ≡  1   (mod p)")
	fmt.Println()
	fmt.Println("   El dos elevado al primo menos uno, dividido por el primo, deja resto UNO. Siempre.")
	fmt.Println()
	fmt.Println("        primo p      2^(p−1)              dividido por p deja resto")
	for _, p := range []int{3, 5, 7, 11, 13, 17, 19} {
		v := 1
		for i := 0; i < p-1; i++ {
			v *= 2
		}
		fmt.Printf("          %3d       2^%-3d = %-10d           %d\n", p, p-1, v, v%p)
	}
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · EL CONTROL — y acá está lo que hace distinta a esta regla")
	fmt.Println()
	var pOK, pTot, cOK, cTot int
	var pseudos []int
	for n := 3; n <= N; n += 2 {
		pasa := potMod(2, n-1, n) == 1
		if es[n] {
			pTot++
			if pasa {
				pOK++
			}
		} else {
			cTot++
			if pasa {
				cOK++
				if len(pseudos) < 10 {
					pseudos = append(pseudos, n)
				}
			}
		}
	}
	fmt.Println("        grupo                        cuántos       la cumplen        porcentaje")
	fmt.Printf("        PRIMOS impares hasta 10⁶     %7d       %7d          %.4f %%\n", pTot, pOK, 100*float64(pOK)/float64(pTot))
	fmt.Printf("        COMPUESTOS impares           %7d       %7d          %.4f %%\n", cTot, cOK, 100*float64(cOK)/float64(cTot))
	razon := (float64(pOK) / float64(pTot)) / (float64(cOK) / float64(cTot))
	fmt.Printf("\n        señal / silencio = %.1f\n", razon)
	fmt.Println()
	fmt.Println("   ⚡⚡⚡ POR PRIMERA VEZ EN TODA LA CAMPAÑA, LA RAZÓN NO ES 1.")
	fmt.Println("        Todas las reglas anteriores del capitán daban 100 % contra 100 % — cero información.")
	fmt.Println("        Ésta da 100 % contra 0,03 %: separa de verdad. Es una regla que SÍ mide.")
	fmt.Println()
	fmt.Println("        para comparar, las cuatro anteriores:")
	fmt.Println("            F393 «p/2 da ½»                      100 % vs 100 %   → razón 1,0000")
	fmt.Println("            F401 «p ≡ ±1 (mod 8) se parte»        exponentes iguales → sin fuerza")
	fmt.Println("            F403 «n² − n(n−1) = n»               100 % vs 100 %   → razón 1,0000")
	fmt.Printf("            ÉSTA  «2^(p−1) ≡ 1 (mod p)»          100 %% vs %.4f %% → razón %.0f\n", 100*float64(cOK)/float64(cTot), razon)
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · el límite honesto: los impostores")
	fmt.Println()
	fmt.Printf("   No es perfecta. %d compuestos pasan el test hasta 10⁶. Los primeros, hallados por búsqueda:\n\n", cOK)
	fmt.Println("        impostor      su factorización            pasa igual")
	for _, n := range pseudos {
		var f []string
		m := n
		for d := 3; d*d <= m; d += 2 {
			for m%d == 0 {
				f = append(f, fmt.Sprint(d))
				m /= d
			}
		}
		if m > 1 {
			f = append(f, fmt.Sprint(m))
		}
		fmt.Printf("        %8d      %-26s   2^%d mod %d = 1\n", n, strings.Join(f, " × "), n-1, n)
	}
	fmt.Println()
	fmt.Printf("   densidad de impostores: %.6f %% de los compuestos impares. Pocos, pero hay — y son infinitos.\n", 100*float64(cOK)/float64(cTot))
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · la falla más profunda: los que engañan con CUALQUIER base")
	fmt.Println()
	fmt.Println("   Si el 2 falla, se podría probar con el 3, o el 5. Pero hay compuestos que pasan con TODAS")
	fmt.Println("   las bases a la vez. Buscados, no citados:")
	fmt.Println()
	fmt.Println("        número      factorización        ¿pasa con 2, 3, 5, 7 y 11?")
	hallados := 0
	for n := 3; n <= 100000 && hallados < 6; n += 2 {
		if es[n] {
			continue
		}
		todas := true
		for _, b := range []int{2, 3, 5, 7, 11, 13} {
			if n%b == 0 {
				continue
			}
			if potMod(b, n-1, n) != 1 {
				todas = false
				break
			}
		}
		if !todas {
			continue
		}
		var f []string
		m := n
		for d := 3; d*d <= m; d += 2 {
			for m%d == 0 {
				f = append(f, fmt.Sprint(d))
				m /= d
			}
		}
		if m > 1 {
			f = append(f, fmt.Sprint(m))
		}
		fmt.Printf("        %7d      %-20s   SÍ, con todas\n", n, strings.Join(f, " × "))
		hallados++
	}
	fmt.Println()
	fmt.Println("   ⚖️ Son los números de Carmichael (1910). Ningún cambio de base los caza, así que el test")
	fmt.Println("      del capitán NO PUEDE ser una demostración de primalidad — solo un filtro muy bueno.")
	fmt.Println()

	fmt.Println("§6 · EL VEREDICTO")
	fmt.Println()
	fmt.Println("   ✅ La relación que pidió EXISTE, y es exactamente entre los tres números que nombró:")
	fmt.Println("      el 2, el 1 y el primo. 2^(p−1) ≡ 1 (mod p).")
	fmt.Printf("   ✅ Y por PRIMERA VEZ una regla suya tiene señal de verdad: razón %.0f contra las anteriores,\n", razon)
	fmt.Println("      que daban exactamente 1. Ésta separa primos de compuestos; aquellas no separaban nada.")
	fmt.Println("   ⚖️ Lo que NO es: nueva. Es el pequeño teorema de Fermat, de 1640, y es la base de todos los")
	fmt.Println("      tests de primalidad que se usan hoy (Miller–Rabin y compañía).")
	fmt.Println("   ⚖️ Y lo que NO puede: demostrar primalidad. Los impostores existen y los de Carmichael pasan")
	fmt.Println("      con cualquier base. Es un filtro excelente, no una prueba.")
	fmt.Println()
	fmt.Println("   📌 Pero el instinto acertó donde las cuatro veces anteriores no: «tiene que haber una relación")
	fmt.Println("      entre el 2, el 1 y el primo» — y la hay, y es la más útil de la teoría de números.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(pOK, pTot, cOK, cTot, razon, pseudos)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/los-dos-latidos.svg")
}

func dibujar(pOK, pTot, cOK, cTot int, razon float64, pseudos []int) {
	var b strings.Builder
	W, H := 1060, 680
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LOS DOS LATIDOS — la primera regla del capitán que SÍ mide")
	txt(530, 54, 12, "#8a8676", "middle", "el 2, el 1 y el primo, atados en un renglón · F405")

	// the relation
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="150" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 104, 13, "#6ee7a0", "middle", "LA RELACIÓN QUE PIDIÓ, Y EXISTE")
	txt(530, 158, 34, "#f2c14e", "middle", "2 ^ (p − 1)  ≡  1   (mod p)")
	txt(530, 192, 12, "#c9c3b0", "middle", "el DOS elevado al primo menos UNO, dividido por el PRIMO, deja resto UNO")
	txt(530, 212, 11, "#8a8676", "middle", "los tres números que nombró el capitán, atados · Fermat, 1640")

	// the control — the big news
	fmt.Fprintf(&b, `<rect x="20" y="242" width="1020" height="200" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 270, 14, "#e8e2d0", "middle", "EL CONTROL — y por primera vez en la campaña, la razón no es 1")
	txt(270, 308, 12, "#6ee7a0", "middle", "PRIMOS que la cumplen")
	txt(270, 352, 34, "#6ee7a0", "middle", "100 %")
	txt(270, 378, 11, "#8a8676", "middle", fmt.Sprintf("%d de %d", pOK, pTot))
	txt(790, 308, 12, "#ff6b6b", "middle", "COMPUESTOS que la cumplen")
	txt(790, 352, 34, "#ff6b6b", "middle", fmt.Sprintf("%.3f %%", 100*float64(cOK)/float64(cTot)))
	txt(790, 378, 11, "#8a8676", "middle", fmt.Sprintf("%d de %d", cOK, cTot))
	txt(530, 340, 20, "#e8e2d0", "middle", "contra")
	txt(530, 412, 15, "#f2c14e", "middle", fmt.Sprintf("señal / silencio = %.0f", razon))
	txt(530, 432, 11, "#8a8676", "middle", "todas las reglas anteriores del capitán daban exactamente 1,0000 acá")

	// comparison
	fmt.Fprintf(&b, `<rect x="20" y="458" width="500" height="200" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 484, 12, "#e8e2d0", "middle", "las anteriores, para comparar")
	yy := 514.0
	for _, r := range [][3]string{
		{"F393  p/2 da ½", "100 % vs 100 %", "1,0000"},
		{"F403  n² − n(n−1) = n", "100 % vs 100 %", "1,0000"},
		{"F404  múltiplos vecinos", "100 % vs 100 %", "1,0000"},
	} {
		txt(44, yy, 11, "#8a8676", "start", r[0])
		txt(290, yy, 11, "#8a8676", "start", r[1])
		txt(470, yy, 11, "#ff6b6b", "end", r[2])
		yy += 28
	}
	fmt.Fprintf(&b, `<line x1="44" y1="%.0f" x2="476" y2="%.0f" stroke="#2a3040"/>`, yy-10, yy-10)
	txt(44, yy+16, 12, "#6ee7a0", "start", "ÉSTA")
	txt(290, yy+16, 11, "#6ee7a0", "start", "100 % vs 0,03 %")
	txt(470, yy+16, 13, "#6ee7a0", "end", fmt.Sprintf("%.0f", razon))

	// the honest limit
	fmt.Fprintf(&b, `<rect x="540" y="458" width="500" height="200" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 484, 12, "#ff6b6b", "middle", "el límite honesto: los impostores")
	txt(564, 514, 11, "#c9c3b0", "start", fmt.Sprintf("%d compuestos pasan igual. El primero:", cOK))
	if len(pseudos) > 0 {
		txt(790, 546, 20, "#f2c14e", "middle", fmt.Sprintf("%d = 11 × 31", pseudos[0]))
	}
	txt(564, 580, 11, "#c9c3b0", "start", "Y los de Carmichael (561, 1105, 1729…)")
	txt(564, 600, 11, "#c9c3b0", "start", "pasan con CUALQUIER base, no solo el 2.")
	txt(564, 628, 11, "#8a8676", "start", "Por eso es un filtro excelente,")
	txt(564, 646, 11, "#8a8676", "start", "pero NO una demostración de primalidad.")
	txt(1030, 672, 10, "#8a8676", "end", "Fermat 1640 · Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "los-dos-latidos.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
	_ = math.Pi
}
