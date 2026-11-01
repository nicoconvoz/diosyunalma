package main

// EL PATRON DEL CAPITAN - he wrote the four bricks himself, each as p plus the
// root of p repeated some number of times:
//
//   (2 + sqrt(2))(2 - sqrt(2))                         1 two   -> sqrt(2)
//   (3 + sqrt(3+3))(3 - sqrt(3+3))                     2 threes -> sqrt(6)
//   (5 + sqrt(5+5+5+5))(5 - sqrt(5+5+5+5))             4 fives -> sqrt(20)
//   (7 + sqrt(7*8))(7 - sqrt(7*8))                     8 sevens -> sqrt(56)
//
// and asked whether he is wrong. He is NOT: all four are correct. But the
// pattern he saw in the counts - 1, 2, 4, 8, the powers of two - is a
// COINCIDENCE of the first four primes, and the real rule underneath is an
// identity that works for every number, prime or not.
//
// Verdicts fixed BEFORE running:
//   V1 all four of the captain's expressions must verify exactly (up to sign
//      for the seven, which used p+1 copies instead of p-1).
//   V2 the powers-of-two reading must BREAK at the fifth prime: 2^4 = 16 is
//      neither 11-1 nor 11+1, so 11 with sixteen copies must fail.
//   V3 the real rule is m = p-1 (or p+1) copies, which is the algebraic
//      identity n^2 - n(n-1) = n. Verified for the primes.
//   V4 (trap #20, the consequence) that identity holds for EVERY integer, not
//      only primes. Measured against composites: if 4, 6, 9 also work, the
//      pattern carries no information about primality.
// Nothing new; the captain found a true identity and the house measures its
// reach, which is the whole job.

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

func main() {
	fmt.Println("🧮 EL PATRÓN DEL CAPITÁN — «decime si estoy mal»")
	fmt.Println()
	fmt.Println("   No está mal. Las cuatro son correctas. Lo que sigue es medir POR QUÉ funcionan,")
	fmt.Println("   porque la razón no es la que parece a simple vista.")
	fmt.Println()

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · las cuatro del capitán, verificadas una por una")
	fmt.Println()
	fmt.Println("        lo que escribió                                cuántas copias    bajo la raíz    resultado")
	casos := []struct {
		p, copias int
		texto     string
	}{
		{2, 1, "(2 + √2)(2 − √2)"},
		{3, 2, "(3 + √(3+3))(3 − √(3+3))"},
		{5, 4, "(5 + √(5+5+5+5))(5 − √(5+5+5+5))"},
		{7, 8, "(7 + √(7+7+7+7+7+7+7+7))(7 − √(…))"},
	}
	for _, c := range casos {
		bajo := c.p * c.copias
		v := c.p*c.p - bajo
		et := "✅ correcto"
		if v == -c.p {
			et = "✅ correcto (da −" + fmt.Sprint(c.p) + ", el mismo en valor absoluto)"
		}
		fmt.Printf("        %-44s      %d            %3d          %+3d   %s\n", c.texto, c.copias, bajo, v, et)
	}
	fmt.Println()
	fmt.Println("   ⚡ Las cuatro dan el número de vuelta. El 7 sale con signo menos porque usó OCHO copias;")
	fmt.Println("      con SEIS habría dado +7. Los dos sirven, y el capitán no se equivocó en ninguna.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · ahora la parte que hay que probar: ¿por qué 1, 2, 4, 8?")
	fmt.Println()
	fmt.Println("   Hay dos lecturas posibles de esa secuencia, y solo una sobrevive:")
	fmt.Println()
	fmt.Println("        lectura A: son POTENCIAS DE DOS — 2⁰, 2¹, 2², 2³, una por cada primo en orden")
	fmt.Println("        lectura B: son p − 1 (o p + 1) — porque 2−1=1, 3−1=2, 5−1=4, 7+1=8")
	fmt.Println()
	fmt.Println("   Para los cuatro primeros primos las dos lecturas COINCIDEN. Se separan en el quinto:")
	fmt.Println()
	fmt.Println("        primo     lectura A (2^k)     ¿funciona?      lectura B (p∓1)     ¿funciona?")
	_, ps := criba(200)
	for k, p := range ps {
		if k >= 8 {
			break
		}
		a := 1 << uint(k)
		va := p*p - p*a
		oka := "NO"
		if va == p || va == -p {
			oka = "sí"
		}
		bmen, bmas := p-1, p+1
		vb := p*p - p*bmen
		okb := "NO"
		if vb == p || vb == -p {
			okb = "sí"
		}
		_ = bmas
		marca := ""
		if oka == "NO" && okb == "sí" {
			marca = "   ← acá se separan"
		}
		fmt.Printf("         %3d         %6d            %-4s           %6d            %-4s%s\n", p, a, oka, bmen, okb, marca)
	}
	fmt.Println()
	fmt.Println("   ⛔ La lectura A (potencias de dos) se rompe en el 11: 2⁴ = 16, y 11² − 11·16 = −55, que no es ±11.")
	fmt.Println("   ✅ La lectura B (p−1 copias) aguanta siempre. Ésa es la regla verdadera.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · la regla verdadera, y es de una línea")
	fmt.Println()
	fmt.Println("        n² − n·(n−1) = n² − n² + n = n")
	fmt.Println()
	fmt.Println("   O sea:  n = (n + √(n·(n−1)))(n − √(n·(n−1)))  — con n−1 copias de n bajo la raíz.")
	fmt.Println()
	fmt.Println("        primo     n−1 copias    bajo la raíz    los dos factores                  producto")
	for _, p := range ps[:6] {
		bajo := p * (p - 1)
		r := math.Sqrt(float64(bajo))
		fmt.Printf("         %3d          %3d           %5d      %12.9f × %12.9f  =  %.9f\n",
			p, p-1, bajo, float64(p)+r, float64(p)-r, (float64(p)+r)*(float64(p)-r))
	}
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · LA TRAMPA (vigésima): ¿le pasa solo a los primos?")
	fmt.Println()
	fmt.Println("   La regla se descubrió mirando primos. Pero la cuenta n² − n(n−1) = n no dice «primo» en")
	fmt.Println("   ninguna parte. Así que hay que probarla en compuestos, que es lo que el control exige:")
	fmt.Println()
	fmt.Println("        número    ¿primo?    n−1 copias    los dos factores                  producto")
	for _, n := range []int{4, 6, 8, 9, 10, 12, 100} {
		bajo := n * (n - 1)
		r := math.Sqrt(float64(bajo))
		fmt.Printf("         %4d       no          %4d      %12.6f × %12.6f  =  %.6f\n",
			n, n-1, float64(n)+r, float64(n)-r, (float64(n)+r)*(float64(n)-r))
	}
	es, _ := criba(100000)
	okP, okC, totP, totC := 0, 0, 0, 0
	for n := 2; n <= 100000; n++ {
		v := n*n - n*(n-1)
		if es[n] {
			totP++
			if v == n {
				okP++
			}
		} else {
			totC++
			if v == n {
				okC++
			}
		}
	}
	fmt.Println()
	fmt.Printf("        primos hasta 100.000 que cumplen:      %6d / %-6d   →  %.4f %%\n", okP, totP, 100*float64(okP)/float64(totP))
	fmt.Printf("        COMPUESTOS hasta 100.000 que cumplen:  %6d / %-6d   →  %.4f %%\n", okC, totC, 100*float64(okC)/float64(totC))
	fmt.Printf("\n        señal / silencio = %.4f  → ", (float64(okP)/float64(totP))/(float64(okC)/float64(totC)))
	fmt.Println("⚠️ VIGÉSIMA trampa del 0.0e+00")
	fmt.Println()
	fmt.Println("   El 100 % contra el 100 % otra vez. La identidad es CIERTA y es LINDA, pero vale para todo")
	fmt.Println("   número entero — el 4, el 6, el 100 —, así que no distingue un primo de un compuesto.")
	fmt.Println()

	fmt.Println("§5 · EL VEREDICTO")
	fmt.Println()
	fmt.Println("   ✅ El capitán NO está mal: las cuatro expresiones son correctas, verificadas.")
	fmt.Println("   ✅ Y encontró una identidad de verdad: n = (n + √(n(n−1)))(n − √(n(n−1))), exacta siempre.")
	fmt.Println("   ⚠️ Pero el patrón 1, 2, 4, 8 que vio en las copias es una COINCIDENCIA de los cuatro")
	fmt.Println("      primeros primos: se rompe en el 11, donde 2⁴ = 16 ya no sirve y hay que usar 10.")
	fmt.Println("   ⚠️ Y la identidad no dice nada de primalidad: le pasa al 4, al 6 y al 100 exactamente igual.")
	fmt.Println()
	fmt.Println("   La diferencia con F398 está en el molde: allá el número bajo la raíz era FIJO (siempre 2) y")
	fmt.Println("   por eso la mitad de los primos no entraba — eso era información. Acá el número bajo la raíz")
	fmt.Println("   CAMBIA con cada primo, y un molde que se adapta a cada pieza entra siempre, y no mide nada.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(ps)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-patron-del-capitan.svg")
}

func dibujar(ps []int) {
	var b strings.Builder
	W, H := 1060, 660
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL PATRÓN DEL CAPITÁN — las cuatro están bien, y la regla es otra")
	txt(530, 54, 12, "#8a8676", "middle", "una identidad verdadera, un patrón que era coincidencia, y una trampa · F403")

	// the four, verified
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="170" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 102, 13, "#6ee7a0", "middle", "✅ LAS CUATRO SON CORRECTAS")
	yy := 132.0
	for _, c := range [][3]string{
		{"(2 + √2)(2 − √2)", "4 − 2", "= 2"},
		{"(3 + √(3+3))(3 − √(3+3))", "9 − 6", "= 3"},
		{"(5 + √(5+5+5+5))(5 − √(5+5+5+5))", "25 − 20", "= 5"},
		{"(7 + √(7·8))(7 − √(7·8))", "49 − 56", "= −7"},
	} {
		txt(70, yy, 13, "#e8e2d0", "start", c[0])
		txt(560, yy, 13, "#c9c3b0", "start", c[1])
		txt(700, yy, 14, "#6ee7a0", "start", c[2])
		yy += 26
	}
	txt(830, 158, 11, "#8a8676", "start", "el 7 usó 8 copias y da −7;")
	txt(830, 176, 11, "#8a8676", "start", "con 6 copias daría +7.")
	txt(830, 194, 11, "#8a8676", "start", "Los dos sirven.")

	// the two readings
	fmt.Fprintf(&b, `<rect x="20" y="262" width="500" height="200" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(270, 288, 13, "#ff6b6b", "middle", "⛔ lectura A: potencias de dos")
	txt(44, 320, 12, "#c9c3b0", "start", "1, 2, 4, 8 = 2⁰, 2¹, 2², 2³")
	txt(44, 346, 11, "#8a8676", "start", "funciona para 2, 3, 5 y 7…")
	txt(44, 376, 12, "#ff6b6b", "start", "…y SE ROMPE en el 11:")
	txt(44, 402, 12, "#e8e2d0", "start", "2⁴ = 16  ⟹  121 − 176 = −55")
	txt(44, 428, 12, "#ff6b6b", "start", "que no es ±11")
	txt(270, 452, 10, "#8a8676", "middle", "era coincidencia de los cuatro primeros")

	fmt.Fprintf(&b, `<rect x="540" y="262" width="500" height="200" fill="#121a16" stroke="#2a5a40"/>`)
	txt(790, 288, 13, "#6ee7a0", "middle", "✅ lectura B: p − 1 copias")
	txt(564, 322, 12, "#c9c3b0", "start", "2−1=1 · 3−1=2 · 5−1=4 · 7+1=8")
	txt(564, 352, 11, "#8a8676", "start", "y en el 11: 11−1 = 10 ⟹ 121 − 110 = 11 ✓")
	txt(790, 394, 15, "#6ee7a0", "middle", "n² − n(n−1) = n")
	txt(790, 422, 11, "#c9c3b0", "middle", "una identidad de una sola línea, exacta siempre")
	txt(790, 452, 10, "#8a8676", "middle", "ésta es la regla verdadera que el capitán encontró")

	// the trap
	fmt.Fprintf(&b, `<rect x="20" y="478" width="1020" height="164" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 504, 14, "#ff6b6b", "middle", "⚠️ LA TRAMPA VIGÉSIMA: la identidad es cierta, pero no es de los primos")
	txt(120, 540, 12, "#c9c3b0", "start", "primos hasta 100.000 que cumplen:")
	txt(560, 540, 13, "#f2c14e", "start", "9.592 / 9.592  =  100,0000 %")
	txt(120, 568, 12, "#c9c3b0", "start", "COMPUESTOS hasta 100.000 que cumplen:")
	txt(560, 568, 13, "#f2c14e", "start", "90.406 / 90.406  =  100,0000 %")
	fmt.Fprintf(&b, `<line x1="120" y1="584" x2="940" y2="584" stroke="#5a2a30"/>`)
	txt(530, 608, 13, "#ff6b6b", "middle", "señal / silencio = 1,0000 — al 4, al 6 y al 100 les pasa exactamente igual")
	txt(530, 632, 11, "#8a8676", "middle", "un molde que se adapta a cada pieza entra siempre, y por eso no mide nada. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-patron-del-capitan.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
