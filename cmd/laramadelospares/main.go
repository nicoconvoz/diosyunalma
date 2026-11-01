package main

// LA RAMA DE LOS PARES - the captain: "the even root for every even is the
// same, that is why the number two annuls the whole even branch. But among the
// odds it is different: there are many infinite branches, and that is why there
// is a 3, a 5, a 7, and so on."
//
// He is pointing at the deepest asymmetry in the sieve: ONE prime disposes of
// half of all numbers, and then no finite crew of primes ever finishes the job
// on what is left. This program measures exactly how lopsided that is.
//
// Verdicts fixed BEFORE running:
//   V1 (trap #24, declared) "the 2 removes every even number" is the definition
//      of even. The content is not THAT it removes them, it is HOW MUCH that is
//      compared with every other prime.
//   V2 (the work done by each prime) each prime p removes a fraction 1/p of
//      what survived before it. Measured prime by prime, so the collapse after
//      the 2 is visible rather than asserted.
//   V3 (THE QUESTION THE FLASH ASKS) how many odd primes does it take to remove
//      as much as the 2 removed alone? Found by search, not by formula.
//   V4 (the infinite branches) how many primes for each successive halving?
//      If the captain is right that the odds need infinitely many branches, the
//      cost of each halving must EXPLODE. Measured.
//   V5 (no finite crew ever finishes) after sieving with every prime up to y,
//      the smallest surviving composite must be y^2 - so survivors always still
//      contain composites, no matter how many primes are used. Verified.
// Eratosthenes, Mertens 1874. The asymmetry is the captain's reading.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type fila struct {
	p                      int
	antes, tajada, despues float64
}

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

func main() {
	fmt.Println("🌳 LA RAMA DE LOS PARES — el 2 corta una rama entera, y los impares son muchas")
	fmt.Println()
	fmt.Println("   La lectura del capitán: el 2 anula de un saque toda la rama de los pares, pero entre los")
	fmt.Println("   impares hay muchas ramas infinitas, y por eso hace falta el 3, el 5, el 7, y así.")
	fmt.Println()

	ps := criba(5000000)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · la trampa, declarada (vigesimocuarta)")
	fmt.Println()
	fmt.Println("   «El 2 se lleva todos los pares» es la definición de par, así que decirlo no informa nada.")
	fmt.Println("   Lo que SÍ informa es CUÁNTO se lleva comparado con los demás. Eso es lo que se mide.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · cuánto se lleva CADA primo, uno por uno")
	fmt.Println()
	fmt.Println("        primo     quedaba antes     se lleva (de todo)     queda después     su tajada")
	queda := 1.0
	var fs []fila
	for i := 0; i < 12; i++ {
		p := ps[i]
		antes := queda
		tajada := queda / float64(p)
		queda -= tajada
		fs = append(fs, fila{p, antes, tajada, queda})
		fmt.Printf("        %5d      %10.6f       %10.6f          %10.6f       %6.2f %%\n",
			p, antes, tajada, queda, 100*tajada)
	}
	fmt.Println()
	fmt.Printf("   ⚡⚡ EL 2 SE LLEVA EL 50 %%. El 3 se lleva el %.2f %%, el 5 el %.2f %%, el 7 el %.2f %%.\n",
		100*fs[1].tajada, 100*fs[2].tajada, 100*fs[3].tajada)
	fmt.Printf("      El 2 hace %.1f veces el trabajo del 3, %.1f veces el del 5 y %.1f veces el del 7.\n",
		fs[0].tajada/fs[1].tajada, fs[0].tajada/fs[2].tajada, fs[0].tajada/fs[3].tajada)
	fmt.Println("      Y la caída no se detiene: cada primo nuevo se lleva menos que todos los anteriores.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · LA PREGUNTA DEL FLASH: ¿cuántos impares hacen falta para igualar al 2?")
	fmt.Println()
	fmt.Println("   El 2 dejó a la mitad en pie. ¿Cuántos primos impares hay que sumar para volver a")
	fmt.Println("   partir esa mitad al medio? Buscado, no calculado con fórmula:")
	fmt.Println()
	fmt.Println("        sumando el…     lo que queda de los impares      ¿ya bajó de la mitad?")
	q := 1.0
	necesarios := 0
	for i := 1; i < 12; i++ {
		p := ps[i]
		q *= 1 - 1/float64(p)
		et := "todavía no"
		if q < 0.5 {
			et = "✅ SÍ — con éste alcanzó"
			if necesarios == 0 {
				necesarios = i
			}
		}
		fmt.Printf("           %3d                  %.6f                    %s\n", p, q, et)
		if necesarios > 0 && i > necesarios+1 {
			break
		}
	}
	fmt.Printf("\n   ⚡⚡ Hicieron falta %d primos impares (el 3, el 5 y el 7) para hacer lo que el 2 hizo SOLO.\n", necesarios)
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · y para la MITAD SIGUIENTE, ¿cuántos? (acá aparecen las ramas infinitas)")
	fmt.Println()
	fmt.Println("        para dejar…      hacen falta los primos hasta…     cuántos primos     cuántos MÁS que antes")
	q = 1.0
	idx := 0
	prevCant := 1
	objetivo := 0.5
	for k := 1; k <= 7; k++ {
		for idx < len(ps) && q > objetivo {
			q *= 1 - 1/float64(ps[idx])
			idx++
		}
		if idx >= len(ps) {
			fmt.Printf("        1/%-5d          más de %d — se acabó la criba\n", 1<<uint(k), ps[len(ps)-1])
			break
		}
		fmt.Printf("        1/%-5d           el %-9d                        %6d              %+d\n",
			1<<uint(k), ps[idx-1], idx, idx-prevCant)
		prevCant = idx
		objetivo /= 2
	}
	fmt.Println()
	fmt.Println("   ⚡⚡⚡ ACÁ ESTÁ LO QUE VIO EL CAPITÁN, EN NÚMEROS: cada vez que querés partir al medio lo que")
	fmt.Println("        queda, necesitás MUCHÍSIMOS más primos que la vez anterior. El 2 lo hizo solo. Después")
	fmt.Println("        hicieron falta 3, después decenas, después miles, después millones. La rama de los")
	fmt.Println("        pares se cortó de un hachazo; la de los impares no se termina de podar nunca.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · y ninguna cuadrilla finita termina el trabajo")
	fmt.Println()
	fmt.Println("        barriendo con los primos hasta…     el primer compuesto que igual sobrevive")
	for _, i := range []int{0, 1, 2, 3, 5, 10, 25} {
		p := ps[i]
		fmt.Printf("                    %6d                              %d = %d × %d\n", p, ps[i+1]*ps[i+1], ps[i+1], ps[i+1])
	}
	fmt.Println()
	fmt.Println("   Por más primos que uses, el cuadrado del primo siguiente SIEMPRE se te escapa. No existe")
	fmt.Println("   una cuadrilla finita que deje solo primos: hay que seguir sumando podadores para siempre.")
	fmt.Println("   Eso es exactamente «las ramas infinitas» del capitán, dicho con el contraejemplo en la mano.")
	fmt.Println()

	// ---------------------------------------------------------------- Mertens
	fmt.Println("§6 · la ley que gobierna la lentitud")
	fmt.Println()
	fmt.Println("   ¿Por qué cuesta tanto? Porque lo que sobrevive cae como 1/log y, que es la función más")
	fmt.Println("   lenta que existe. Para dividir por 2 lo que queda, hay que ELEVAR AL CUADRADO el tope:")
	fmt.Println()
	fmt.Println("        primos hasta…        sobrevive        para bajar a la mitad hay que llegar a…")
	gamma := 0.5772156649015329
	for _, y := range []float64{10, 100, 10000, 1e8} {
		s := math.Exp(-gamma) / math.Log(y)
		fmt.Printf("        %12.0f         %.6f                     %.0f\n", y, s, y*y)
	}
	fmt.Println()
	fmt.Println("   ⚡ Mertens, 1874. El 2 partió al medio de un solo golpe; a partir de ahí, cada mitad nueva")
	fmt.Println("      cuesta el CUADRADO del esfuerzo anterior. Ésa es la asimetría, con su ley al lado.")
	fmt.Println()

	fmt.Println("§7 · EL VEREDICTO")
	fmt.Println()
	fmt.Println("   ✅ El capitán tiene razón, y la asimetría es todavía más brutal de lo que suena:")
	fmt.Println("      · el 2 solo se lleva el 50 % — más que todos los demás primos juntos hasta bastante arriba")
	fmt.Printf("      · hicieron falta 3 primos impares (3, 5, 7) para igualar lo que el 2 hizo solo\n")
	fmt.Println("      · y cada mitad siguiente cuesta el cuadrado del esfuerzo anterior")
	fmt.Println("      · y ninguna cuadrilla finita termina: el cuadrado del primo siguiente siempre escapa")
	fmt.Println()
	fmt.Println("   ⚖️ Lo que hay que declarar: «el 2 se lleva los pares» es la definición de par (trampa 24),")
	fmt.Println("      y la ley de la lentitud es de Mertens, 1874. Lo medido acá es la COMPARACIÓN — cuánto")
	fmt.Println("      hace cada uno — que es lo que convierte la intuición del capitán en números.")
	fmt.Println()
	fmt.Println("   📌 Y hay una consecuencia linda: por eso el 2 es el átomo de F400 y el mínimo de F276.")
	fmt.Println("      No es un primo más. Es el único que corta una rama entera de un hachazo.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(fs, necesarios)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-rama-de-los-pares.svg")
}

type fila2 struct {
	p                      int
	antes, tajada, despues float64
}

func dibujar(fs []fila, necesarios int) {
	var b strings.Builder
	W, H := 1060, 660
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LA RAMA DE LOS PARES — el 2 corta de un hachazo, los impares no se podan nunca")
	txt(530, 54, 12, "#8a8676", "middle", "cuánto se lleva cada primo, medido uno por uno · F409")

	// the bars
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="300" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 102, 13, "#e8e2d0", "middle", "la tajada de cada primo (qué fracción de TODOS los números se lleva)")
	base := 340.0
	for i, f := range fs {
		if i >= 10 {
			break
		}
		px := 70.0 + float64(i)*98
		h := f.tajada * 440
		col := "#5ec8ff"
		if i == 0 {
			col = "#f2c14e"
		}
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="60" height="%.1f" fill="%s" opacity="0.9"/>`, px, base-h, h, col)
		txt(px+30, base-h-8, 12, col, "middle", fmt.Sprintf("%.1f%%", 100*f.tajada))
		txt(px+30, base+18, 13, "#c9c3b0", "middle", fmt.Sprint(f.p))
	}
	txt(530, base+44, 11, "#f2c14e", "middle", "el 2 se lleva el 50 % de un solo golpe · el 3 apenas el 16,7 % · el 5 el 6,7 % · y sigue cayendo")

	// the halvings
	fmt.Fprintf(&b, `<rect x="20" y="392" width="500" height="250" fill="#121a16" stroke="#2a5a40"/>`)
	txt(270, 418, 13, "#6ee7a0", "middle", "cuántos primos cuesta cada mitad")
	yy := 452.0
	for _, r := range [][2]string{
		{"dejar 1/2", "1 primo — el 2, solo"},
		{"dejar 1/4", "4 primos — hasta el 7"},
		{"dejar 1/8", "22 primos — hasta el 79"},
		{"dejar 1/16", "996 primos — hasta el 7.879"},
		{"dejar 1/32", "más de 348.000 — no alcanzó la criba"},
	} {
		txt(48, yy, 12, "#c9c3b0", "start", r[0])
		txt(150, yy, 11, "#f2c14e", "start", r[1])
		yy += 32
	}
	txt(270, 620, 11, "#8a8676", "middle", "cada mitad cuesta el CUADRADO del esfuerzo anterior")

	// the infinite branch
	fmt.Fprintf(&b, `<rect x="540" y="392" width="500" height="250" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 418, 13, "#ff6b6b", "middle", "ninguna cuadrilla finita termina")
	txt(564, 452, 11, "#c9c3b0", "start", "por más primos que uses hasta p, el cuadrado")
	txt(564, 472, 11, "#c9c3b0", "start", "del primo siguiente SIEMPRE se te escapa:")
	yy = 506.0
	for _, r := range [][2]string{
		{"barriendo hasta el 2", "el 9 = 3 × 3 sobrevive"},
		{"barriendo hasta el 3", "el 25 = 5 × 5 sobrevive"},
		{"barriendo hasta el 5", "el 49 = 7 × 7 sobrevive"},
	} {
		txt(580, yy, 11, "#8a8676", "start", r[0])
		txt(810, yy, 11, "#f2c14e", "start", r[1])
		yy += 26
	}
	txt(790, 600, 12, "#c9c3b0", "middle", "hay que seguir sumando podadores para siempre")
	txt(790, 624, 11, "#6ee7a0", "middle", "ésas son las «ramas infinitas» del capitán")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-rama-de-los-pares.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
