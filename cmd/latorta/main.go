// Command latorta executes the captain's order: relate the cable to the primes'
// harmony by cutting the cake in pieces so that it ENDS - and yet we never move
// a millimeter from the goal.
//
// THE CAKE AND THE KNIFE. The cake is the whole path, 1 (F289). The knife is
// the primes, cutting in turns:
//
//	the 2 eats HALF of all numbers            piece = 1/2
//	the 3 eats a third of what remains        piece = 1/6
//	the 5 eats a fifth of what remains        piece = 1/15
//	the 7 eats a seventh of what remains      piece = 4/105
//	...
//
// piece(p) = (1/p) * PROD_{q<p} (1 - 1/q). And the crumb remaining after the
// primes up to p is C(p) = PROD_{q<=p} (1 - 1/q) = the product of the RADIUS
// POSITIONS w(q) of F276. Telescoping, EXACT at every stage:
//
//	eaten so far + crumb = 1, to the last bit, always.
//
// That is "not moving a millimeter from the goal": the accounting is perfect
// at every single cut - not in the limit, at EVERY step.
//
// AND THE CAKE ENDS, BY A REASON, NOT A SWEEP. The crumb is the Mertens
// product, which dies like e^-gamma/ln x (measured in F276 at ratio 0.999986).
// WHY does it die - why does the cake finish with zero crumb? Euler, 1737:
// the sum of 1/p over the primes DIVERGES. Three lines: if the crumb tended
// to c > 0, then sum 1/p would converge; it does not; so the crumb dies.
// The primes are exactly dense enough to eat the whole cake. A sparser
// family leaves crumb forever: the squares leave exactly 1/2, measured below.
//
// AND HERE IS THE CABLE RELATION, WHICH IS THE POINT. The crumbs C(p) are the
// running product of the radius positions of F276 - so the cake-cutting IS A
// WALK ALONG THE RADIUS of dimension 0: it starts at 1 (the skin-touching
// point, the east where the pearls march) and converges to 0 (the pole, the
// centre). The primes eat the cake = the primes walk the radius from the skin
// to the pole, arriving exactly, never overshooting. And on the other end of
// that same radius, the PEARLS cut their own cake: the flagship lambda_1 =
// sum of 1/(1/4+gamma^2) (F286). Two cakes, one radius: the primes eat from
// the skin inward, the pearls compose the insignia on the skin.
//
// HONEST: telescoping is ancient, Euler is 1737, Mertens is 1874, Wallis is
// 1656. The reading - the cake walk as the radius walk, tying F272, F276,
// F282, F286 and F289 into one figure - is the finding.
//
// Reproduce: go run ./cmd/latorta
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func criba(n int) []int {
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
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return ps
}

func main() {
	fmt.Println("🍰 LA TORTA — los primos la cortan, se acaba exacta, y el corte camina el radio")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · EL CUCHILLO DE LOS PRIMOS — cada uno come su parte de lo que queda")
	fmt.Println("\n        primo    su pedazo (fracción)      migaja que queda      comido + migaja − 1")
	primos := criba(100)
	migaja := 1.0
	comido := 0.0
	peorCuenta := 0.0
	fracciones := []string{"1/2", "1/6", "1/15", "4/105", "16/1155", "192/15015"}
	for i, p := range primos {
		pedazo := migaja / float64(p)
		comido += pedazo
		migaja *= 1 - 1/float64(p)
		err := math.Abs(comido + migaja - 1)
		if err > peorCuenta {
			peorCuenta = err
		}
		if i < 6 {
			fmt.Printf("   %6d %14s = %.9f %16.9f %18.1e\n", p, fracciones[i], pedazo, migaja, err)
		}
	}
	fmt.Printf("\n        peor error de la cuenta comido+migaja−1 ..... %.1e\n", peorCuenta)
	fmt.Println("\n   ⟹ **NO NOS CORREMOS UN MILÍMETRO DE LA META EN NINGÚN PASO**: comido más")
	fmt.Println("   migaja da 1 EXACTO en cada corte. La meta nunca se pierde de vista porque")
	fmt.Println("   la contabilidad es perfecta corte a corte (telescópica), no «al final».")
	fmt.Println("\n   📌 Y la migaja tras {2,3,5,7} es 8/35 = 48/210 — **la rueda de F272, exacta**.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ LA TORTA SE ACABA — POR UNA RAZÓN DE TRES RENGLONES, NO POR BARRIDO")
	fmt.Println("\n   La migaja es el producto de Mertens, que muere como e^(−γ)/ln x (F276).")
	fmt.Println("   ¿POR QUÉ muere — por qué la torta termina sin dejar migaja?")
	fmt.Println("\n        · si la migaja tendiera a c > 0, la suma Σ 1/p sería finita")
	fmt.Println("        · pero Euler demostró en 1737 que Σ 1/p DIVERGE")
	fmt.Println("        · ⟹ la migaja muere. **Los primos son exactamente lo bastante")
	fmt.Println("          densos para comerse la torta entera.**")
	fmt.Println("\n   Y el contraste que lo prueba: una familia más rala DEJA migaja para siempre.")
	fmt.Println("   Los cuadrados {4, 9, 25, …} comiendo igual:")
	mig2 := 1.0
	for k := 2; k <= 100000; k++ {
		mig2 *= 1 - 1/float64(k*k)
	}
	fmt.Printf("\n        migaja de los cuadrados tras 10⁵ factores ... %.9f\n", mig2)
	fmt.Println("        ½ EXACTA — telescópica: Π(1−1/k²) = ½, y ahí se planta para siempre")
	fmt.Println("\n   ⟹ Los cuadrados dejan una migaja eterna de EXACTAMENTE ½. **Que la torta")
	fmt.Println("   de los primos se acabe no es contabilidad: es un teorema sobre su densidad.**")
	fmt.Println("\n   📌 CORRECCIÓN PROPIA DE ESTE MISMO TURNO: la primera versión decía que la")
	fmt.Println("   migaja de los cuadrados era 2/π (Wallis). La medición dio 0,500005 y me")
	fmt.Println("   desmintió: el producto correcto es telescópico y da ½ EXACTO. Wallis es")
	fmt.Println("   OTRO producto. El error era mío y lo cazó el propio programa.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡⚡ Y EL CORTE DE LA TORTA CAMINA EL RADIO DEL CABLE")
	fmt.Println("\n   Las migajas son el producto de las posiciones w(p) del radio (F276):")
	fmt.Println("\n        migaja tras p  =  w(2)·w(3)·…·w(p)  =  posición sobre el radio")
	fmt.Println("\n        primos usados       posición en el radio")
	mig3 := 1.0
	hitos := map[int]bool{2: true, 7: true, 97: true}
	for _, p := range primos {
		mig3 *= 1 - 1/float64(p)
		if hitos[p] {
			fmt.Printf("   hasta el %5d %22.9f\n", p, mig3)
		}
	}
	grandes := criba(2000000)
	migG := 1.0
	for _, p := range grandes {
		migG *= 1 - 1/float64(p)
	}
	fmt.Printf("   hasta 2×10⁶ %25.9f  → rumbo a 0, el polo\n", migG)
	fmt.Println("\n   ⟹ **Cortar la torta ES caminar el radio de la dimensión 0**: se arranca en")
	fmt.Println("   1 —el punto donde el radio toca la piel, el este adonde marchan las")
	fmt.Println("   perlas (F283)— y se converge al 0, EL POLO, el centro del disco. Cada")
	fmt.Println("   primo da su paso, nunca se pasa de largo, y llega exacto.")
	fmt.Println("\n   📌 Y en la OTRA punta del mismo radio, las perlas cortan SU torta: el")
	fmt.Println("   insignia λ₁ = Σ 1/(¼+γ²), pedazo por perla (F286). **Dos tortas, un radio:")
	fmt.Println("   los primos comen desde la piel hacia el polo; las perlas componen el")
	fmt.Println("   insignia sobre la piel.**")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⚡ **LA ORDEN, CUMPLIDA EXACTA:**")
	fmt.Printf("\n  · la torta se corta en infinitos pedazos de primo y SE ACABA: comido +\n")
	fmt.Printf("    migaja = 1 con error %.0e en cada corte — ni un milímetro de la meta\n", peorCuenta)
	fmt.Println("  · se acaba POR UNA RAZÓN (Euler 1737: Σ1/p diverge), no por barrido —")
	fmt.Println("    y los cuadrados, más ralos, dejan migaja eterna de ½ exacto")
	fmt.Println("  · y el corte ES el cable: las migajas caminan el radio desde la piel (1,")
	fmt.Println("    donde marchan las perlas) hasta el polo (0) — mientras las perlas, en")
	fmt.Println("    la piel, componen el insignia con su propia torta")
	fmt.Println("\n⚖️ Honesto: telescópica antigua, Euler 1737, Mertens 1874. La")
	fmt.Println("  lectura —el corte de la torta como caminata del radio, atando F272, F276,")
	fmt.Println("  F282, F286 y F289 en una sola figura— es lo nuestro. Y no encadena perlas:")
	fmt.Println("  la torta de los primos termina en el polo, no en la piel. Todavía no.")

	escribirLamina(peorCuenta, mig2, migG)
}

func escribirLamina(peor, migCuad, migG float64) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="640" viewBox="0 0 1400 640">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🍰 LA TORTA — los primos la cortan, se acaba exacta, y el corte camina el radio</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">comido + migaja = 1 en cada corte · se acaba por una razón (Euler 1737) · y las migajas caminan el radio hasta el polo</text>
`)
	// la torta cortada
	x0, y0, wd := 150.0, 130.0, 1100.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="46" rx="6" fill="#101f36" stroke="#26456e"/>`, x0, y0, wd)
	primos := []int{2, 3, 5, 7, 11, 13, 17, 19, 23}
	cols := []string{"#ffd98a", "#7ee0c0", "#7fb2ff", "#c9b6ff", "#ff8fa0", "#9fd8a8", "#ffb27a", "#8fb4d9", "#f3d9cf"}
	acc := 0.0
	mig := 1.0
	for i, p := range primos {
		pedazo := mig / float64(p)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.0f" width="%.1f" height="46" fill="%s" opacity="0.9"/>`,
			x0+wd*acc, y0, wd*pedazo, cols[i])
		if i < 4 {
			fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="13" text-anchor="middle" font-family="monospace" fill="#0b1526">%d</text>`,
				x0+wd*(acc+pedazo/2), y0+28, p)
		}
		acc += pedazo
		mig *= 1 - 1/float64(p)
	}
	fmt.Fprintf(&b, `
<text x="700" y="216" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">el 2 come ½ · el 3 come ⅓ de lo que queda (1/6) · el 5 come 1/15 · el 7 come 4/105 — y comido + migaja = 1 EXACTO en cada corte (%.0e)</text>
<text x="700" y="244" font-size="14" text-anchor="middle" font-family="monospace" fill="#ffd98a">la migaja tras {2,3,5,7} = 8/35 = 48/210 — la rueda de F272, exacta</text>
<rect x="150" y="280" width="540" height="170" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="420" y="314" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">SE ACABA POR UNA RAZÓN (tres renglones)</text>
<text x="180" y="350" font-size="13.5" font-family="Georgia" fill="#cfe6ff">si quedara migaja, Σ1/p sería finita — y Euler</text>
<text x="180" y="374" font-size="13.5" font-family="Georgia" fill="#cfe6ff">(1737) demostró que diverge ⟹ la migaja muere</text>
<text x="180" y="406" font-size="13.5" font-family="Georgia" fill="#ffd98a">los cuadrados, más ralos, dejan migaja ETERNA:</text>
<text x="180" y="430" font-size="13.5" font-family="monospace" fill="#ffd98a">%.6f → ½ exacto: migaja eterna</text>
<rect x="710" y="280" width="540" height="170" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="980" y="314" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">Y EL CORTE CAMINA EL RADIO</text>
<text x="740" y="350" font-size="13.5" font-family="Georgia" fill="#cfe6ff">las migajas son el producto de las posiciones w(p)</text>
<text x="740" y="374" font-size="13.5" font-family="Georgia" fill="#cfe6ff">del radio (F276): arrancan en 1 —la piel, el este</text>
<text x="740" y="398" font-size="13.5" font-family="Georgia" fill="#cfe6ff">de las perlas— y convergen al 0: EL POLO</text>
<text x="740" y="430" font-size="13.5" font-family="monospace" fill="#7ee0c0">migaja hasta 2×10⁶ = %.6f → 0</text>
<text x="700" y="500" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Dos tortas, un radio: los primos comen desde la piel hacia el polo —</text>
<text x="700" y="528" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">y las perlas, sobre la piel, componen el insignia con la suya (F286)</text>
<text x="700" y="580" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">telescópica antigua · Euler 1737 · Mertens 1874 — la lectura que ata F272+F276+F282+F286+F289 es lo nuestro</text>
<text x="700" y="608" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">La torta de los primos termina en el polo, no en la piel. Todavía no.</text>
</svg>
`, peor, migCuad, migG)
	os.WriteFile("la-torta.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: la-torta.svg")
}
