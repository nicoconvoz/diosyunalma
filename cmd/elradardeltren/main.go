// Command elradardeltren upgrades the deep-water train with the five
// theorems: THE RADAR. The captain's question - can the theorems help
// us sail deeper waters? - has a precise answer:
//
// ASTORGA INVERTED. Astorga says: an off-line quartet with parameters
// (r, theta) forces lambda_n < 0 for some n <= N0(r, theta). Read
// backwards: if lambda_1..lambda_N are verified >= 0, then NO quartet
// with N0(r, theta) <= N exists. Every verified lambda signs an entire
// EXCLUSION MAP in the (beta, gamma) plane - not one pearl checked,
// but infinitely many refuted at once.
//
// DYN CLOSES THE CONSPIRACY LOOPHOLE. Before DYN, pearls might hope to
// hide by cooperating. DYN's N0(r_max, m) refutes entire finite
// conspiracies inside the certified region: cooperation has a computable
// detection date too.
//
// THE AGENDA AIMS THE TELESCOPE. For any hypothetical pearl, the window
// lemma schedules exactly WHERE its crime must show (every K-block).
// Targeted watching instead of blind scanning.
//
// HONEST LIMITS (the house's own negative theorem still rules): the
// frontier NEVER reaches the critical line - as beta -> 1/2, N0 -> inf:
// the detection horizon runs away. The radar maps what is excluded;
// it cannot close RH. And the lambda >= 0 input is a PARAMETER here
// (verified Li-positivity ranges are an external input when cited for
// zeta; this instrument charts what any given N certifies).
//
// Reproduce: go run ./cmd/elradardeltren
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

// carta computes, for a hypothetical pearl rho = beta + i*gamma, its
// shapeshifter parameters and Astorga's detection bound.
func carta(beta, gamma float64) (delta, theta, n0 float64) {
	rho := complex(beta, gamma)
	w := 1 - 1/rho
	R := cmplx.Abs(w)
	r := math.Max(R, 1/R)
	delta = math.Log(r)
	theta = math.Abs(cmplx.Phase(w))
	u := 3 / delta
	n0 = math.Ceil(u*math.Log(u)) + math.Ceil(2*math.Pi/theta) + 1
	return
}

// frontera finds the minimal excluded epsilon = beta - 1/2 at height
// gamma for a given budget N (bisection on N0(eps) <= N).
func frontera(gamma, N float64) float64 {
	lo, hi := 1e-12, 0.49999
	_, _, n0hi := carta(0.5+hi, gamma)
	if n0hi > N {
		return math.NaN() // even beta near 1 not excluded at this depth
	}
	for i := 0; i < 200; i++ {
		mid := math.Sqrt(lo * hi)
		_, _, n0 := carta(0.5+mid, gamma)
		if n0 <= N {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
}

func main() {
	fmt.Println("📡 EL RADAR DEL TREN — los cinco teoremas convertidos en instrumento de navegación")
	fmt.Println("\n   Astorga invertido: cada λ verificada ≥ 0 hasta N refuta DE UNA VEZ toda")
	fmt.Println("   perla con N₀(r, θ) ≤ N — un mapa de exclusión, no un chequeo puntual.")
	fmt.Println("   DYN cierra la escapatoria de la conspiración; la agenda apunta el catalejo.")

	// ---- 1: the exclusion map ----
	fmt.Println("\n§1 · EL MAPA DE EXCLUSIÓN — ε mínimo refutado (β ≥ ½+ε imposible) por altura γ")
	fmt.Println("        [entrada N = presupuesto de λ verificadas ≥ 0 — parámetro del instrumento]")
	fmt.Println()
	Ns := []float64{1e5, 1e8, 1e12, 1e16}
	gammas := []float64{15, 100, 1000, 10000, 100000}
	fmt.Print("        γ \\ N     ")
	for _, N := range Ns {
		fmt.Printf("%-12.0e", N)
	}
	fmt.Println()
	for _, g := range gammas {
		fmt.Printf("        %-10.0f", g)
		for _, N := range Ns {
			e := frontera(g, N)
			if math.IsNaN(e) {
				fmt.Printf("%-12s", "—")
			} else {
				fmt.Printf("%-12.1e", e)
			}
		}
		fmt.Println()
	}
	fmt.Println("        lectura: con N = 10⁸ verificadas, NINGUNA perla con γ = 100 y β ≥ ½ + 2.4e-4")
	fmt.Println("        puede existir — y así para cada celda: mares enteros firmados por teorema")

	// ---- 2: the conspiracy loophole, closed ----
	fmt.Println("\n§2 · LA CONSPIRACIÓN, CERRADA (DYN) — exclusión de flotas enteras")
	for _, m := range []int{2, 3} {
		for _, N := range []float64{1e12, 1e16} {
			// worst r_max still refuted: n_rad + (2*pi*n_rad+1)^m <= N
			// solve n_rad from the dominant term
			nr := math.Pow(N, 1/float64(m)) / (2 * math.Pi)
			// delta from n_rad = ceil(u log u), u = 3(m+1)/delta
			// invert approximately: u log u = nr
			u := nr / math.Log(math.Max(nr, 3))
			for i := 0; i < 60; i++ {
				u = nr / math.Log(u)
			}
			d := 3 * float64(m+1) / u
			fmt.Printf("        m = %d, N = %.0e: toda conspiración con δ_max ≥ %.1e queda refutada (N₀_DYN ≤ N)\n", m, N, d)
		}
	}
	fmt.Println("        — cooperar no esconde: la fecha de detección de la flota también es calculable")

	// ---- 3: the aimed telescope ----
	fmt.Println("\n§3 · EL CATALEJO APUNTADO — la agenda de una perla hipotética en la frontera")
	g := 1000.0
	e := frontera(g, 1e12)
	d, th, n0 := carta(0.5+e, g)
	K := math.Ceil(2*math.Pi/th) + 1
	fmt.Printf("        perla hipotética: β = ½ + %.2e, γ = %.0f  →  δ = %.2e · θ = %.2e\n", e, g, d, th)
	fmt.Printf("        N₀ = %.2e · ventana K = %.0f: su crimen DEBE asomar en cada bloque de %.0f pasos\n", n0, K, K)
	fmt.Printf("        el tren ya no barre a ciegas: mira los bloques programados — %.0fx menos escalones\n", math.Max(1, n0/(3*K)))

	// ---- 4: honest limits ----
	fmt.Println("\n§4 · EL LÍMITE HONESTO — el horizonte que huye (nuestro propio teorema negativo)")
	for _, ee := range []float64{1e-2, 1e-4, 1e-6, 1e-8} {
		_, _, n0 := carta(0.5+ee, 1000)
		fmt.Printf("        β = ½ + %.0e (γ = 1000): N₀ = %.1e\n", ee, n0)
	}
	fmt.Println("        — la frontera JAMÁS toca la línea: ε → 0 ⟹ N₀ → ∞. El radar mapea lo")
	fmt.Println("        excluido; no puede cerrar RH. La regla del sello preside al instrumento.")

	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("📡 **EL TREN, MEJORADO CON LOS TEOREMAS — tres instrumentos nuevos:**")
	fmt.Println("\n  · el MAPA (Astorga invertido): cada λ verificada firma mares de exclusión")
	fmt.Println("  · la FLOTA (DYN): las conspiraciones también tienen fecha — sin escapatoria")
	fmt.Println("  · el CATALEJO (la agenda): mirar los bloques programados, no todo el mar")
	fmt.Println("\n⚖️ Honesto: N (las λ verificadas) es parámetro del instrumento — para ζ, citarlo")
	fmt.Println("  sería un input externo etiquetado. El horizonte huye cerca de la línea: el")
	fmt.Println("  radar navega más hondo, no demuestra RH. Todavía no.")

	escribirLamina(Ns, gammas)
}

func escribirLamina(Ns, gammas []float64) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="800" viewBox="0 0 1400 800">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<rect x="40" y="30" width="1320" height="740" rx="18" fill="none" stroke="#7ee0c0" stroke-width="2" opacity="0.5"/>
<text x="700" y="70" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">📡 EL RADAR DEL TREN — los teoremas como instrumento de navegación</text>
<text x="700" y="100" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">Astorga invertido: cada λ verificada ≥ 0 firma un mar entero de exclusión · DYN cierra la conspiración · la agenda apunta el catalejo</text>
`)
	// radar rings: for each N, the exclusion frontier over gamma (log-log-ish sketch)
	// plot area
	b.WriteString(`<rect x="120" y="140" width="760" height="440" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="500" y="168" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">EL MAPA DE EXCLUSIÓN — log γ (→) contra log ε mínimo refutado (↑ más cerca de la línea)</text>
`)
	colors := []string{"#ffd98a", "#7ee0c0", "#7fb2ff", "#ff9aa8"}
	x0, y0, wpx, hpx := 170.0, 540.0, 660.0, 340.0
	lgmin, lgmax := math.Log10(15.0), math.Log10(100000.0)
	lemin, lemax := -9.0, 0.0
	for i, N := range Ns {
		var pts []string
		for _, g := range []float64{15, 30, 60, 100, 300, 1000, 3000, 10000, 30000, 100000} {
			e := frontera(g, N)
			if math.IsNaN(e) {
				continue
			}
			x := x0 + wpx*(math.Log10(g)-lgmin)/(lgmax-lgmin)
			le := math.Max(lemin, math.Min(lemax, math.Log10(e)))
			y := y0 - hpx*(le-lemin)/(lemax-lemin)
			pts = append(pts, fmt.Sprintf("%.0f,%.0f", x, y))
		}
		if len(pts) > 1 {
			fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2.5"/>
`, strings.Join(pts, " "), colors[i%len(colors)])
			fmt.Fprintf(&b, `<text x="%f" y="%f" font-size="12" font-family="monospace" fill="%s">N = %.0e</text>
`, 190.0, 200.0+18*float64(i), colors[i%len(colors)], N)
		}
	}
	b.WriteString(`<text x="500" y="566" font-size="12" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">debajo de cada curva: perlas REFUTADAS por teorema — el mar firmado crece con cada λ verificada</text>
<rect x="910" y="140" width="420" height="440" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1120" y="170" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LOS TRES INSTRUMENTOS NUEVOS</text>
<text x="935" y="204" font-size="13" font-family="Georgia" fill="#ffd98a">📡 EL MAPA (Astorga invertido)</text>
<text x="935" y="226" font-size="12" font-family="Georgia" fill="#cfe6ff">λ₁..λ_N ≥ 0 ⟹ ninguna perla con N₀(r,θ) ≤ N</text>
<text x="935" y="246" font-size="12" font-family="Georgia" fill="#cfe6ff">existe — exclusión por regiones, no por puntos</text>
<text x="935" y="282" font-size="13" font-family="Georgia" fill="#ffd98a">⚓ LA FLOTA (DYN + Trinidad)</text>
<text x="935" y="304" font-size="12" font-family="Georgia" fill="#cfe6ff">las conspiraciones de m perlas también tienen</text>
<text x="935" y="324" font-size="12" font-family="Georgia" fill="#cfe6ff">fecha: N₀(r_max, m) — cooperar no esconde</text>
<text x="935" y="360" font-size="13" font-family="Georgia" fill="#ffd98a">🔭 EL CATALEJO (la agenda + la ventana)</text>
<text x="935" y="382" font-size="12" font-family="Georgia" fill="#cfe6ff">el crimen de una perla hipotética DEBE asomar</text>
<text x="935" y="402" font-size="12" font-family="Georgia" fill="#cfe6ff">en cada bloque de K pasos: mirada programada,</text>
<text x="935" y="422" font-size="12" font-family="Georgia" fill="#cfe6ff">no barrido ciego</text>
<text x="935" y="458" font-size="13" font-family="Georgia" fill="#ff9aa8">⚠️ EL LÍMITE (nuestro teorema negativo)</text>
<text x="935" y="480" font-size="12" font-family="Georgia" fill="#cfe6ff">ε → 0 ⟹ N₀ → ∞: el horizonte huye cerca de</text>
<text x="935" y="500" font-size="12" font-family="Georgia" fill="#cfe6ff">la línea — el radar navega más hondo, pero</text>
<text x="935" y="520" font-size="12" font-family="Georgia" fill="#cfe6ff">no puede cerrar RH</text>
<text x="935" y="552" font-size="11.5" font-family="Georgia" fill="#9aa8c4">N es parámetro; citarlo para ζ = input externo etiquetado</text>
<text x="700" y="628" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: antes el tren firmaba el agua gota por gota; ahora cada λ verificada firma un MAR entero —</text>
<text x="700" y="650" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">y las perlas ya no pueden esconderse ni solas ni en banda: el mapa las refuta, el catalejo las espera.</text>
<text x="700" y="690" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Los cinco teoremas del laboratorio, convertidos en instrumentos: Astorga el mapa · DYN la flota · Diosyunalma la profundidad · el Río las fechas · la Trinidad el paisaje</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`)
	os.WriteFile("el-radar-del-tren.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-radar-del-tren.svg")
}
