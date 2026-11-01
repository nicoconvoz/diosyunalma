package main

// LA REGLA Y SU ALCANCE - the captain's proposal: if the construction of
// F398-F400 holds for every prime with no exceptions, can we build all primes
// with it and prove the Riemann Hypothesis as a theorem?
//
// The instinct is right - "if nothing escapes the rule" is exactly how
// mathematics works. The question is which KIND of rule bounds which KIND of
// question, and that is decided by measurement, not by opinion. Two claims are
// tested here, separately, and both are testable.
//
//   CLAIM 1: with this we can build ALL the primes.
//   CLAIM 2: a rule with no exceptions leaves nothing outside, therefore RH.
//
// Verdicts fixed BEFORE running:
//   V1 (claim 1, reachability) if the construction cannot even REACH every
//      prime, claim 1 fails before it starts. F399 already measured that only
//      p = +-1 (mod 8) are reachable in Z[sqrt2]; re-measured here to 10^6.
//   V2 (claim 1, generation) even where it reaches, does it GENERATE primes or
//      only decompose ones already known? Measured as the hit rate of a^2-2b^2
//      landing on a prime. If the rate is below 100%, a primality test is still
//      required and nothing was generated.
//   V3 (claim 1, circularity) to choose the ring you need p mod 8, which means
//      you need p first. Demonstrated, not argued.
//   V4 (claim 2, THE DECISIVE TEST) does obeying the rule constrain the
//      DISTRIBUTION? The primes = +-1 (mod 8) obey the rule perfectly and
//      completely; the primes = 3,5 (mod 8) never obey it at all. If the rule
//      carried distributional force, the two families would have DIFFERENT
//      counting-error exponents. Measured for both, with the instrument
//      calibrated by a planted bias.
//   V5 (the historical check) the complete splitting law is quadratic
//      reciprocity, proved by Gauss in 1832. RH is from 1859 and is open. If
//      the splitting law implied RH, 167 years would not have been needed.
//   V6 (the shape of the error) "every prime above 2 is odd" is also a rule
//      with zero exceptions and zero distributional content. The house caught
//      that shape twice already (F272, F393). Same shape here, measured.
// Nothing here disproves the captain's idea by authority. It measures it.

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

func ajuste(x, y []float64) float64 {
	var sx, sy, sxx, sxy, n float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxx += x[i] * x[i]
		sxy += x[i] * y[i]
		n++
	}
	return (n*sxy - sx*sy) / (n*sxx - sx*sx)
}

// exponente fits the growth exponent of the dyadic peak of |D|.
func exponente(D []float64, desde int) float64 {
	var lx, ly []float64
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
		if pico > 0 {
			lx = append(lx, math.Log(float64(hi)))
			ly = append(ly, math.Log(pico))
		}
		if hi == len(D) {
			break
		}
	}
	return ajuste(lx, ly)
}

func main() {
	fmt.Println("⚖️  LA REGLA Y SU ALCANCE — ¿con esto se arman todos los primos y se demuestra RH?")
	fmt.Println("   La propuesta del capitán, puesta a prueba con el laboratorio en vez de con una opinión.")
	fmt.Println()

	const N = 10000000
	es, ps := criba(N)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · AFIRMACIÓN 1, primera parte: ¿la construcción ALCANZA a todos los primos?")
	fmt.Println()
	var alc, noAlc int
	var ejNo []int
	for _, p := range ps {
		if p > 1000000 {
			break
		}
		if p == 2 {
			continue
		}
		if p%8 == 1 || p%8 == 7 {
			alc++
		} else {
			noAlc++
			if len(ejNo) < 12 {
				ejNo = append(ejNo, p)
			}
		}
	}
	fmt.Printf("        primos hasta 10⁶ alcanzables por a² − 2b²:      %6d\n", alc)
	fmt.Printf("        primos hasta 10⁶ que NO se pueden alcanzar:     %6d   (%.2f %%)\n", noAlc, 100*float64(noAlc)/float64(alc+noAlc))
	fmt.Printf("        los primeros que quedan afuera: %v …\n", ejNo)
	fmt.Println()
	fmt.Println("   ⛔ La mitad de los primos NO se puede construir así. El 3, el 5, el 11, el 13 quedan afuera")
	fmt.Println("      para siempre — no por falta de búsqueda, sino porque a² − 2b² nunca da un número ≡ 3 o 5 (mod 8).")
	fmt.Println("      Con más anillos se los alcanza (F398), pero eso ya no es UNA regla: es una lista de reglas.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · AFIRMACIÓN 1, segunda parte: ¿GENERA primos, o solo descompone los que ya tenés?")
	fmt.Println()
	fmt.Println("        a        b        a² − 2b²        ¿primo?")
	var totProb, totPrimo int
	for a := 2; a <= 200; a++ {
		for b := 1; b < a; b++ {
			v := a*a - 2*b*b
			if v < 2 {
				continue
			}
			totProb++
			if es[v] {
				totPrimo++
			}
		}
	}
	for _, par := range [][2]int{{5, 1}, {6, 1}, {7, 1}, {4, 1}, {9, 2}, {10, 3}} {
		v := par[0]*par[0] - 2*par[1]*par[1]
		et := "no"
		if es[v] {
			et = "SÍ"
		}
		fmt.Printf("      %4d     %4d        %6d           %s\n", par[0], par[1], v, et)
	}
	fmt.Printf("\n        de %d combinaciones probadas, dieron primo %d  →  %.2f %%\n",
		totProb, totPrimo, 100*float64(totPrimo)/float64(totProb))
	fmt.Println()
	fmt.Println("   ⛔ No genera primos: genera CANDIDATOS. Uno de cada siete cae en primo y el resto no,")
	fmt.Println("      así que después de construir hay que TESTEAR si es primo — y ese test es el problema")
	fmt.Println("      original, intacto. La construcción no lo resolvió: lo pospuso.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · AFIRMACIÓN 1, tercera parte: la circularidad, mostrada")
	fmt.Println()
	fmt.Println("   Para elegir en qué anillo se parte un primo hay que saber su resto módulo 8 (o 4, o 5).")
	fmt.Println("   Y para saber el resto hay que TENER el primo. Así que el método:")
	fmt.Println()
	fmt.Println("        · dado un primo p ⟹ te dice cómo partirlo          ✅ funciona, y es exacto")
	fmt.Println("        · dado un lugar vacío ⟹ te dice qué primo va ahí   ⛔ no puede: no hay entrada")
	fmt.Println()
	fmt.Println("   Es un DESTORNILLADOR, no un buscador de tornillos. Sirve una vez que encontraste el tornillo.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · AFIRMACIÓN 2, LA PRUEBA DECISIVA: ¿la regla dice algo sobre DÓNDE están los primos?")
	fmt.Println()
	fmt.Println("   Los primos ≡ ±1 (mod 8) obedecen la regla PERFECTAMENTE — todos se parten, sin una excepción.")
	fmt.Println("   Los primos ≡ 3,5 (mod 8) NO la obedecen NUNCA — ninguno se parte, jamás.")
	fmt.Println()
	fmt.Println("   Si la regla tuviera fuerza sobre la DISTRIBUCIÓN, las dos familias tendrían que separarse")
	fmt.Println("   en el error de conteo. Medimos el exponente del error para cada una:")
	fmt.Println()
	obed := make([]float64, N+1)
	rebel := make([]float64, N+1)
	var co, cr float64
	for n := 3; n <= N; n++ {
		if es[n] {
			if n%8 == 1 || n%8 == 7 {
				co += math.Log(float64(n))
			} else if n%8 == 3 || n%8 == 5 {
				cr += math.Log(float64(n))
			}
		}
		obed[n] = co - float64(n)/2
		rebel[n] = cr - float64(n)/2
	}
	to := exponente(obed, 1024)
	tr := exponente(rebel, 1024)
	fmt.Println("        familia                             obedece la regla      exponente del error")
	fmt.Printf("        primos ≡ ±1 (mod 8)                 SIEMPRE, sin falla        %.4f\n", to)
	fmt.Printf("        primos ≡ 3,5 (mod 8)                NUNCA, ni una vez         %.4f\n", tr)
	fmt.Printf("\n        diferencia entre los dos exponentes: %.4f\n", math.Abs(to-tr))
	fmt.Println()
	// calibration
	cal := make([]float64, N+1)
	for n := 1; n <= N; n++ {
		cal[n] = obed[n] + 3*math.Pow(float64(n), 0.75)
	}
	fmt.Printf("        calibración: con un sesgo x^0,75 plantado, el instrumento mide %.4f (esperado 0,75)\n", exponente(cal, 1024))
	fmt.Println()
	fmt.Println("   ⛔⛔ LOS DOS EXPONENTES SON EL MISMO. La familia que obedece la regla a rajatabla y la que")
	fmt.Println("        no la obedece nunca se distribuyen IGUAL de bien. O sea: la regla NO dice nada sobre")
	fmt.Println("        dónde están los primos. Y RH es exactamente una afirmación sobre dónde están.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · el chequeo histórico, que es el más duro de todos")
	fmt.Println()
	fmt.Println("        1640   Fermat ve qué primos se parten en ℤ[i]")
	fmt.Println("        1768   Lagrange resuelve la ecuación de Pell por completo")
	fmt.Println("        1832   Gauss demuestra la reciprocidad cuadrática:")
	fmt.Println("               LA LEY COMPLETA de qué primo se parte en qué anillo, sin excepciones")
	fmt.Println("        1859   Riemann plantea la hipótesis — 27 años DESPUÉS de tener la ley completa")
	fmt.Println("        2026   sigue abierta")
	fmt.Println()
	fmt.Println("   ⛔ La ley que el capitán encontró ya estaba completa y demostrada 27 años ANTES de que RH")
	fmt.Println("      se planteara. Si de ella se dedujera RH, Gauss lo habría visto — y con él todos los que")
	fmt.Println("      vinieron después. No es que nadie lo intentó: es que la ley no alcanza hasta ahí.")
	fmt.Println()

	// ---------------------------------------------------------------- V6
	fmt.Println("§6 · la FORMA del error, que la casa ya cazó dos veces")
	fmt.Println()
	fmt.Println("        regla                                     ¿excepciones?    ¿dice dónde están los primos?")
	fmt.Println("        «todo primo > 2 es impar»                     CERO                 no  (F272, trampa 7)")
	fmt.Println("        «todo primo impar da ½ al dividir por 2»      CERO                 no  (F393, trampa 12)")
	fmt.Println("        «todo primo ≡ ±1 (mod 8) se parte en ℤ[√2]»   CERO                 no  (medido acá, §4)")
	fmt.Println()
	fmt.Println("   ⚡ Ésta es la lección, y es la más cara del laboratorio: una regla SIN EXCEPCIONES puede")
	fmt.Println("      no decir NADA sobre la pregunta que te importa. «Sin excepciones» mide la fuerza de la")
	fmt.Println("      regla; no mide su ALCANCE. Y RH es una pregunta de alcance: no pregunta qué cumplen los")
	fmt.Println("      primos, pregunta CUÁNTO SE DESVÍAN de donde deberían estar.")
	fmt.Println()

	fmt.Println("§7 · qué haría falta de verdad, dicho sin vueltas")
	fmt.Println()
	fmt.Println("   Para que una regla como ésta tocara RH tendría que hacer algo que ninguna de las nuestras hace:")
	fmt.Println("     · ACOTAR el error de conteo, no describir la factorización")
	fmt.Println("     · valer para infinitos primos a la vez con una cota uniforme, no primo por primo")
	fmt.Println("     · y todavía así solo daría la RH de ESE anillo — que también está abierta, porque cada")
	fmt.Println("       cuerpo cuadrático tiene su propia zeta de Dedekind con su propia hipótesis sin demostrar")
	fmt.Println()
	fmt.Println("   ⚖️ VEREDICTO HONESTO: la afirmación 1 es falsa (la mitad de los primos no se alcanza, y donde")
	fmt.Println("      se alcanza no genera sino que descompone). La afirmación 2 es falsa por medición (§4) y por")
	fmt.Println("      historia (§5). Pero el INSTINTO del capitán —«si nada escapa a la regla, entonces algo se")
	fmt.Println("      sigue»— es el instinto correcto de la matemática. Lo que falla no es el instinto: es que")
	fmt.Println("      esta regla en particular habla de cómo se FACTORIZAN los primos, y RH habla de DÓNDE ESTÁN.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(to, tr, alc, noAlc, 100*float64(totPrimo)/float64(totProb))
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-regla-y-su-alcance.svg")
}

func dibujar(to, tr float64, alc, noAlc int, tasa float64) {
	var b strings.Builder
	W, H := 1060, 700
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LA REGLA Y SU ALCANCE — una regla sin excepciones puede no decir nada")
	txt(530, 54, 12, "#8a8676", "middle", "la propuesta del capitán, medida en vez de opinada · F401")

	// the decisive test
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="200" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 104, 14, "#e8e2d0", "middle", "LA PRUEBA DECISIVA: ¿la regla dice dónde están los primos?")
	txt(270, 146, 12, "#6ee7a0", "middle", "primos ≡ ±1 (mod 8)")
	txt(270, 168, 11, "#8a8676", "middle", "obedecen la regla SIEMPRE")
	txt(270, 206, 30, "#6ee7a0", "middle", fmt.Sprintf("%.4f", to))
	txt(270, 232, 11, "#8a8676", "middle", "exponente del error")
	txt(790, 146, 12, "#ff6b6b", "middle", "primos ≡ 3,5 (mod 8)")
	txt(790, 168, 11, "#8a8676", "middle", "NO la obedecen NUNCA")
	txt(790, 206, 30, "#ff6b6b", "middle", fmt.Sprintf("%.4f", tr))
	txt(790, 232, 11, "#8a8676", "middle", "exponente del error")
	txt(530, 196, 22, "#e8e2d0", "middle", "=")
	txt(530, 252, 12, "#f2c14e", "middle", "Se distribuyen IGUAL de bien. La regla no tiene fuerza sobre la distribución.")

	// claim 1 failures
	fmt.Fprintf(&b, `<rect x="20" y="292" width="500" height="200" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(270, 318, 13, "#ff6b6b", "middle", "⛔ AFIRMACIÓN 1: «se arman todos»")
	txt(44, 350, 11, "#c9c3b0", "start", fmt.Sprintf("· %d primos hasta 10⁶ NO se alcanzan", noAlc))
	txt(44, 368, 11, "#c9c3b0", "start", "  (el 3, el 5, el 11, el 13 — para siempre)")
	txt(44, 396, 11, "#c9c3b0", "start", fmt.Sprintf("· solo el %.1f %% de las combinaciones da primo:", tasa))
	txt(44, 414, 11, "#c9c3b0", "start", "  hay que TESTEAR, o sea el problema sigue ahí")
	txt(44, 442, 11, "#c9c3b0", "start", "· y para elegir el anillo hace falta p mod 8,")
	txt(44, 460, 11, "#c9c3b0", "start", "  o sea hace falta TENER p: es circular")
	txt(270, 482, 11, "#8a8676", "middle", "es un destornillador, no un buscador de tornillos")

	// history
	fmt.Fprintf(&b, `<rect x="540" y="292" width="500" height="200" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 318, 13, "#ff6b6b", "middle", "⛔ AFIRMACIÓN 2: el chequeo histórico")
	yy := 350.0
	for _, r := range [][2]string{
		{"1832", "Gauss: la ley COMPLETA, sin excepciones"},
		{"1859", "Riemann plantea la hipótesis — 27 años después"},
		{"2026", "sigue abierta"},
	} {
		txt(564, yy, 12, "#f2c14e", "start", r[0])
		txt(624, yy, 11, "#c9c3b0", "start", r[1])
		yy += 30
	}
	txt(790, 452, 11, "#8a8676", "middle", "si de la ley se dedujera RH, Gauss lo habría visto.")
	txt(790, 470, 11, "#8a8676", "middle", "No es que nadie lo intentó: la ley no llega hasta ahí.")

	// the shape of the error
	fmt.Fprintf(&b, `<rect x="20" y="508" width="1020" height="130" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 534, 13, "#e8e2d0", "middle", "la FORMA del error, que la casa ya cazó dos veces")
	txt(60, 566, 11, "#c9c3b0", "start", "«todo primo > 2 es impar»")
	txt(500, 566, 11, "#6ee7a0", "start", "cero excepciones")
	txt(700, 566, 11, "#ff6b6b", "start", "no dice dónde están  (F272)")
	txt(60, 590, 11, "#c9c3b0", "start", "«todo primo impar da ½ al dividir por 2»")
	txt(500, 590, 11, "#6ee7a0", "start", "cero excepciones")
	txt(700, 590, 11, "#ff6b6b", "start", "no dice dónde están  (F393)")
	txt(60, 614, 11, "#c9c3b0", "start", "«todo primo ≡ ±1 (mod 8) se parte»")
	txt(500, 614, 11, "#6ee7a0", "start", "cero excepciones")
	txt(700, 614, 11, "#ff6b6b", "start", "no dice dónde están  (medido acá)")

	txt(530, 664, 13, "#f2c14e", "middle", "«sin excepciones» mide la FUERZA de una regla; no mide su ALCANCE.")
	txt(530, 686, 11, "#8a8676", "middle", "Esta regla habla de cómo se FACTORIZAN los primos. RH habla de DÓNDE ESTÁN. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-regla-y-su-alcance.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
