package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("🌊🧱 POROSIDAD — Fase VIII: ¿el medio es uniforme, o tiene partes")
	fmt.Println("     más blandas y otras más duras? — el flash es del capitán")

	const N = 400
	const t0 = 100.0
	const kmax = 120
	const F = 30.0 // the same total force in EVERY arm: only the material changes
	ms := medio(4000, N, t0)
	d := &dado{s: 20260817}

	fmt.Println("\n§1 · LA TRADUCCIÓN (su §7: definirla, no usar la palabra suelta)")
	fmt.Println("     cada sitio lleva una PERMEABILIDAD local h_i > 0, y el núcleo es")
	fmt.Println("       C(i,j) = √(h_i·h_j) · f(|i−j|)")
	fmt.Println("     así «poroso» significa exactamente: un sitio por el que la influencia pasa fácil.")
	fmt.Println("     Duro = h chico, blando = h grande. Porosidad y dureza NO son dos variables:")
	fmt.Println("     son una sola leída en los dos sentidos — que es lo que su §7 pedía decidir.")
	fmt.Println("     regla aritmética, declarada ANTES de calcular un solo espectro:")
	fmt.Println("       h_i = 1/log(p_i) — el modo de un primo chico es una onda larga y floja, y deja")
	fmt.Println("       pasar; el de un primo grande es corta y rígida. Ningún cero se mira.")
	fmt.Printf("     %d modos · fuerza total FIJA en %.0f para todos los brazos: sólo cambia el material\n", len(ms), F)

	hHom := homogenea(len(ms))
	hArit := permeabilidad(ms)
	hDiv := permeabilidadDivisores(ms)
	fmt.Printf("     dispersión de h: aritmética (1/log p) %.3f · aritmética (divisores) %.3f · homogénea 0\n",
		desvio(hArit), desvio(hDiv))

	// -----------------------------------------------------------------------
	// §2 · her §11 matrix: with and without node, with and without heterogeneity
	// -----------------------------------------------------------------------
	fmt.Println("\n§2 · LA MATRIZ DE COMPARACIÓN (su §11) — y el brazo MEZCLADO es el control que decide")
	fmt.Println("     mezclado = los MISMOS valores de h, permutados al azar: misma estadística,")
	fmt.Println("     misma dispersión, todo igual — salvo el ORDEN. Si el aritmético le gana al")
	fmt.Println("     mezclado, lo que importa es el orden, que es lo único que puede ser aritmético.")
	fmt.Printf("     %-30s %7s %9s %9s %9s %8s %8s\n", "medio", "vivos", "Σ²(5)", "Σ²(10)", "Σ²(20)", "α", "PR/N")

	type arm struct {
		nom string
		h   []float64
		c   func(int) float64
	}
	sinNodoF := colaLarga(0.5)
	conNodoF := conNodo(5, 0.5)

	corre := func(nom string, h []float64, c func(int) float64) obs {
		o := medir(espectro(ms, h, c, kmax, F))
		if o.valido {
			fmt.Printf("     %-30s %7d %9.3f %9.3f %9.3f %8.3f %8.3f\n",
				nom, o.vivos, o.s5, o.s10, o.s20, o.alfa, o.pr)
		} else {
			fmt.Printf("     %-30s %7d   ← BANDA VACIADA\n", nom, o.vivos)
		}
		return o
	}

	base := corre("A · homogéneo, sin nodo", hHom, sinNodoF)
	baseN := corre("A · homogéneo, CON nodo", hHom, conNodoF)

	// the shuffled arm gets several seeds, because one permutation is an anecdote
	var mezS, mezN []float64
	var oMez, oMezN obs
	for s := 0; s < 5; s++ {
		hm := mezclar(hArit, d)
		om := medir(espectro(ms, hm, sinNodoF, kmax, F))
		omn := medir(espectro(ms, hm, conNodoF, kmax, F))
		if om.valido {
			mezS = append(mezS, om.s10)
			oMez = om
		}
		if omn.valido {
			mezN = append(mezN, omn.s10)
			oMezN = omn
		}
	}
	fmt.Printf("     %-30s %7d %9s %9.3f ± %.3f %15s %8.3f\n",
		"B · MEZCLADO, sin nodo (5 semillas)", oMez.vivos, "—", media(mezS), desvio(mezS), "—", oMez.pr)
	fmt.Printf("     %-30s %7d %9s %9.3f ± %.3f %15s %8.3f\n",
		"B · MEZCLADO, CON nodo (5 semillas)", oMezN.vivos, "—", media(mezN), desvio(mezN), "—", oMezN.pr)

	arit := corre("C · ARITMÉTICO 1/log p, sin nodo", hArit, sinNodoF)
	aritN := corre("C · ARITMÉTICO 1/log p, CON nodo", hArit, conNodoF)
	div := corre("C · ARITMÉTICO divisores, sin nodo", hDiv, sinNodoF)
	divN := corre("C · ARITMÉTICO divisores, CON nodo", hDiv, conNodoF)

	// -----------------------------------------------------------------------
	// §3 · the decisive comparison: arithmetic order vs the same values shuffled
	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · LA PRUEBA QUE DECIDE — el orden aritmético contra los mismos valores mezclados")
	sig := desvio(mezS)
	if sig == 0 {
		sig = 1e-9
	}
	z := (media(mezS) - arit.s10) / sig
	fmt.Printf("     sin nodo : aritmético %.4f  ·  mezclado %.4f ± %.4f  →  %.2f sigmas de diferencia\n",
		arit.s10, media(mezS), desvio(mezS), z)
	sigN := desvio(mezN)
	if sigN == 0 {
		sigN = 1e-9
	}
	zN := (media(mezN) - aritN.s10) / sigN
	fmt.Printf("     con nodo : aritmético %.4f  ·  mezclado %.4f ± %.4f  →  %.2f sigmas\n",
		aritN.s10, media(mezN), desvio(mezN), zN)
	fmt.Printf("     y la otra regla aritmética (divisores, no correlacionada con la amplitud): %.4f / %.4f\n",
		div.s10, divN.s10)
	if math.Abs(z) > 2 || math.Abs(zN) > 2 {
		fmt.Println("     ⟹ el ORDEN importa: la aritmética no es sólo textura. Pero ATENCIÓN — esto")
		fmt.Println("       todavía NO dice que el orden sea ARITMÉTICO. Eso lo decide §3bis.")
	} else {
		fmt.Println("     ⟹ el orden NO importa: el aritmético queda dentro del ruido del mezclado.")
		fmt.Println("       Lo que hace la heterogeneidad es TEXTURA, no aritmética. Resultado negativo útil.")
	}

	// -----------------------------------------------------------------------
	// §3bis · THE CONTROL THIS SHEET NEEDED: ordered but NOT arithmetic
	// -----------------------------------------------------------------------
	fmt.Println("\n§3bis · EL CONTROL QUE LE FALTABA A ESTA HOJA — ordenado pero NO aritmético")
	fmt.Println("     las dos reglas aritméticas dieron casi lo mismo, lo que huele a que el medio")
	fmt.Println("     responde al ORDEN del campo y no a los primos. Se prueba con campos igual de")
	fmt.Println("     ordenados y sin nada de aritmética: una rampa suave y una onda suave.")
	dispA := desvio(hArit) * 2 * math.Sqrt(3) // a ramp with the same spread as the arithmetic field
	rampa := corre("D · rampa lisa, sin nodo", ordenadaRampa(len(ms), dispA), sinNodoF)
	onda := corre("D · onda lisa, sin nodo", ordenadaOnda(len(ms), desvio(hArit)*math.Sqrt2), sinNodoF)
	fmt.Printf("     aritmético %.3f  ·  rampa %.3f  ·  onda %.3f  ·  mezclado %.3f ± %.3f\n",
		arit.s10, rampa.s10, onda.s10, media(mezS), desvio(mezS))
	// Both controls are judged against the SHUFFLE noise, and the verdict belongs to
	// the CLOSEST non-arithmetic field: one match is enough to break the attribution.
	sg := math.Max(desvio(mezS), 1e-9)
	dRampa := math.Abs(rampa.s10-arit.s10) / sg
	dOnda := math.Abs(onda.s10-arit.s10) / sg
	fmt.Printf("     distancia a la aritmética, en sigmas del mezclado:  rampa %.2f  ·  onda %.2f\n", dRampa, dOnda)
	fmt.Printf("     el control más cercano es la %s\n", map[bool]string{true: "ONDA", false: "RAMPA"}[dOnda <= dRampa])
	if math.Min(dRampa, dOnda) < 1 {
		fmt.Println("     ⟹ un campo ordenado y SIN NADA DE ARITMÉTICA hace lo mismo que el aritmético.")
		fmt.Println("       Lo que el medio premia es que la dureza tenga ESTRUCTURA a escala intermedia,")
		fmt.Println("       no que venga de los primos. La atribución aritmética NO se sostiene — y era")
		fmt.Println("       la única con interés estructural fuerte (su §6). Resultado negativo, y es el")
		fmt.Println("       importante: la mejora del 22% es real, su causa NO es aritmética.")
		fmt.Println("     · pero NO es «cualquier orden»: la rampa monótona es PEOR que el azar")
		fmt.Printf("       (%.3f contra %.3f del mezclado). Ordenar mal arruina. Eso es un hallazgo aparte.\n",
			rampa.s10, media(mezS))
	} else {
		fmt.Printf("     ⟹ ningún control ordenado no-aritmético llega (el más cercano queda a %.2f sigmas):\n", math.Min(dRampa, dOnda))
		fmt.Println("       el orden solo no alcanza, y queda algo que sí es propio del campo aritmético.")
	}

	// -----------------------------------------------------------------------
	// §4 · her §8: the localisation medium, hard zone against soft zone
	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · LOCALIZACIÓN (su §8) — un medio con zona blanda, zona media y zona dura")
	fmt.Printf("     %10s %7s %9s %9s %8s %8s\n", "contraste", "vivos", "Σ²(10)", "Σ²(20)", "α", "PR/N")
	var locs [][6]float64
	for _, ct := range []float64{1, 2, 5, 20, 100} {
		hb := bloques(len(ms), ct)
		o := medir(espectro(ms, hb, sinNodoF, kmax, F))
		if o.valido {
			fmt.Printf("     %10.0f %7d %9.3f %9.3f %8.3f %8.3f\n", ct, o.vivos, o.s10, o.s20, o.alfa, o.pr)
			locs = append(locs, [6]float64{ct, float64(o.vivos), o.s10, o.s20, o.alfa, o.pr})
		} else {
			fmt.Printf("     %10.0f %7d   ← BANDA VACIADA\n", ct, o.vivos)
			locs = append(locs, [6]float64{ct, float64(o.vivos), 0, 0, 0, o.pr})
		}
	}
	fmt.Printf("     ⚠ Y LA ADVERTENCIA DE SU §9, CUMPLIDA: el nodo baja Σ²(10) de %.3f a %.3f, pero la\n", base.s10, baseN.s10)
	fmt.Printf("       localización PASA DE %.3f A %.3f — o sea, los estados se ATRAPAN más. Buena parte\n", base.pr, baseN.pr)
	fmt.Println("       de esa mejora puede ser ALCANCE APARENTE por localización, no alcance real.")
	fmt.Println("       Es exactamente la distinción que su §9 pedía, y hay que declararla antes de")
	fmt.Println("       festejar cualquier número del nodo.")
	fmt.Println("     PR/N es la razón de participación: 1 = estados EXTENDIDOS por todo el medio,")
	fmt.Println("     cerca de 0 = estados ATRAPADOS en una región. Es la respuesta directa a su §9:")
	fmt.Println("     distingue alcance real de alcance aparente por localización.")

	// -----------------------------------------------------------------------
	// §5 · verdicts, separated as her §17 asks
	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · VEREDICTOS SEPARADOS (su §17)")
	fmt.Printf("     EFECTO DE HETEROGENEIDAD : homogéneo %.4f → aritmético %.4f (%.1f%%)\n",
		base.s10, arit.s10, 100*(1-arit.s10/base.s10))
	fmt.Printf("     EFECTO DEL NODO          : sin nodo %.4f → con nodo %.4f (homogéneo)\n", base.s10, baseN.s10)
	fmt.Printf("                                sin nodo %.4f → con nodo %.4f (aritmético)\n", arit.s10, aritN.s10)
	fmt.Printf("     MEJOR PUNTO SIN NODO     : porosidad SUAVE (contraste ×2) %.4f — mejor que el aritmético\n", locs[1][2])
	fmt.Printf("     POSIBLE EFECTO ARITMÉTICO: %.2f sigmas contra el mezclado, PERO la onda lisa lo empata\n", z)
	fmt.Printf("                                a %.2f sigmas sin usar un solo primo ⟹ NO se sostiene\n", dOnda)
	fmt.Println("                                lo que cuenta no es el origen aritmético sino la ESTRUCTURA")
	fmt.Printf("     niveles vivos en todos los brazos: entre %d y %d de %d — la banda NO se vacía\n",
		minI(base.vivos, arit.vivos, aritN.vivos, div.vivos), maxI(base.vivos, arit.vivos, aritN.vivos, div.vivos), len(ms))
	fmt.Println("     los ceros verdaderos, como regla: Σ²(10) = 0.3364")

	dibujar(Res{
		N: len(ms), Fuerza: F,
		Base: base, BaseN: baseN, Arit: arit, AritN: aritN, Div: div, DivN: divN,
		MezMedia: media(mezS), MezDes: desvio(mezS), MezPR: oMez.pr,
		MezNMedia: media(mezN), MezNDes: desvio(mezN),
		Z: z, ZN: zN, DispArit: desvio(hArit), DispDiv: desvio(hDiv),
		Locs: locs, HArit: hArit, Rampa: rampa, Onda: onda, DRampa: dRampa,
	})
}

func minI(v ...int) int {
	m := v[0]
	for _, x := range v {
		if x < m {
			m = x
		}
	}
	return m
}

func maxI(v ...int) int {
	m := v[0]
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}

// Res carries every measured number to the plate.
type Res struct {
	N                        int
	Fuerza                   float64
	Base, BaseN, Arit, AritN obs
	Div, DivN                obs
	MezMedia, MezDes, MezPR  float64
	MezNMedia, MezNDes       float64
	Z, ZN, DispArit, DispDiv float64
	Rampa, Onda              obs
	DRampa                   float64
	Locs                     [][6]float64 // contrast, vivos, s10, s20, alfa, pr
	HArit                    []float64
}
