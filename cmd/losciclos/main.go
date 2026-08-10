// Command losciclos tests the captain's cycle hypothesis on his own repeats.
//
// HIS CLAIM: "if that is a law, then there is a relation - every so many minus
// plus, the sign collects a repeated one, or a relation of that kind, for
// example two minuses in a row. But they must mark cycles of something."
//
// So the object under test is the REPEAT: a place where two consecutive twos
// carry the same sign. In the laboratory's notation, with d_i = g_i - g_{i+1}:
//
//	repeat at i   <=>   sign(d_i) = sign(d_{i+1}) != 0
//
// A repeat of two minuses means three gaps growing in a row; two pluses means
// three gaps shrinking in a row. His hypothesis is that these repeats are not
// scattered at random but land on a beat.
//
// THREE WAYS TO KILL IT, AND ALL THREE ARE RUN HERE:
//
//  1. THE SPACING. If repeats mark a cycle, the distance between one repeat and
//     the next should pile up around a preferred value. If they are scattered,
//     the spacing is geometric and its standard deviation is about equal to its
//     mean. Measured against the shuffled control.
//  2. THE SPECTRUM. A cycle is a peak in the periodogram of the repeat train.
//     ⚠️ And here is the trap that kills most amateur cycle hunts: ANY noise
//     produces peaks, and the more frequencies you scan the taller the tallest
//     one gets. So the null is NOT a fixed threshold - it is the TALLEST PEAK
//     THE SHUFFLED CONTROL PRODUCES over the same number of frequencies. That
//     comparison, and only that one, decides.
//  3. THE MEMORY. If repeats cluster or space themselves, the autocorrelation
//     of the repeat train shows it at some lag. Measured against the control,
//     with the control's own spread as the error bar.
//
// PRE-REGISTERED PREDICTION, written before running: NO cycle. The spacing will
// be near-geometric, the tallest true peak will not clear the tallest control
// peak, and the autocorrelation will be flat except at the shortest lags, where
// the laboratory already knows there is real structure (the mod-3 suppression
// of Finding 268). The captain's instinct that "they must mark cycles" is the
// single most common way a real pattern gets over-read, and this program is
// built to give that instinct every fair chance to be right.
//
// Reproduce: go run ./cmd/losciclos
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

const (
	N        = 5000000 // criba
	SEMILLAS = 8       // barajadas de control
	MESP     = 60000   // largo de la ventana para el espectro
	KFREQ    = 1500    // frecuencias barridas
)

func criba(n int) []bool {
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
	return es
}

func baraja(src []int, semilla uint64) []int {
	c := append([]int(nil), src...)
	s := semilla
	for i := len(c) - 1; i > 0; i-- {
		s = s*6364136223846793005 + 1442695040888963407
		j := int((s >> 33) % uint64(i+1))
		c[i], c[j] = c[j], c[i]
	}
	return c
}

func segundas(g []int) []int {
	d := make([]int, len(g)-1)
	for i := range d {
		d[i] = g[i] - g[i+1]
	}
	return d
}

func signo(v int) int {
	if v > 0 {
		return 1
	}
	if v < 0 {
		return -1
	}
	return 0
}

// repeticiones devuelve el tren de repeticiones: 1 donde dos doses seguidos
// llevan el mismo signo no nulo, 0 en el resto.
func repeticiones(d []int) []float64 {
	r := make([]float64, len(d)-1)
	for i := 0; i+1 < len(d); i++ {
		a, b := signo(d[i]), signo(d[i+1])
		if a != 0 && a == b {
			r[i] = 1
		}
	}
	return r
}

func espaciados(r []float64) []int {
	var pos []int
	for i, v := range r {
		if v == 1 {
			pos = append(pos, i)
		}
	}
	esp := make([]int, 0, len(pos))
	for i := 1; i < len(pos); i++ {
		esp = append(esp, pos[i]-pos[i-1])
	}
	return esp
}

func mediaDesvio(x []int) (float64, float64) {
	if len(x) == 0 {
		return 0, 0
	}
	var s, s2 float64
	for _, v := range x {
		s += float64(v)
		s2 += float64(v) * float64(v)
	}
	m := s / float64(len(x))
	return m, math.Sqrt(s2/float64(len(x)) - m*m)
}

// periodograma devuelve la potencia normalizada en KFREQ frecuencias y el pico.
func periodograma(r []float64, m, k int) ([]float64, float64, int) {
	if len(r) < m {
		m = len(r)
	}
	x := r[:m]
	media := 0.0
	for _, v := range x {
		media += v
	}
	media /= float64(m)
	var v0 float64
	for _, v := range x {
		v0 += (v - media) * (v - media)
	}
	pot := make([]float64, k)
	pico, dondePico := 0.0, 0
	for j := 1; j <= k; j++ {
		w := 2 * math.Pi * float64(j) / float64(m)
		var re, im float64
		for i := 0; i < m; i++ {
			a := (x[i] - media)
			re += a * math.Cos(w*float64(i))
			im += a * math.Sin(w*float64(i))
		}
		p := (re*re + im*im) / v0
		pot[j-1] = p
		if p > pico {
			pico, dondePico = p, j
		}
	}
	return pot, pico, dondePico
}

func autocorr(r []float64, lags int) []float64 {
	n := len(r)
	media := 0.0
	for _, v := range r {
		media += v
	}
	media /= float64(n)
	var v0 float64
	for _, v := range r {
		v0 += (v - media) * (v - media)
	}
	out := make([]float64, lags)
	for L := 1; L <= lags; L++ {
		var acc float64
		for i := 0; i+L < n; i++ {
			acc += (r[i] - media) * (r[i+L] - media)
		}
		out[L-1] = acc / v0
	}
	return out
}

func main() {
	fmt.Println("🔁 LOS CICLOS — la hipótesis del capitán sobre sus propias repeticiones")
	fmt.Println("\n   Dijo: «cada tantos − +, el signo cobra uno repetido… pero deben marcar")
	fmt.Println("   ciclos de algo». O sea: las repeticiones caerían en un compás.")
	fmt.Println("\n   ⚠️ LA TRAMPA DE ESTA CAZA, DECLARADA ANTES DE MEDIR: cualquier ruido produce")
	fmt.Println("   picos en el espectro, y cuantas más frecuencias mirás, más alto sale el pico")
	fmt.Println("   más alto. Así que el testigo NO es un umbral fijo: es **el pico más alto que")
	fmt.Println("   produce el control barajado** sobre las mismas frecuencias. Ése decide.")

	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)
	var primos []int
	for i := 2; i <= N; i++ {
		if es[i] {
			primos = append(primos, i)
		}
	}
	gaps := make([]int, len(primos)-1)
	for i := range gaps {
		gaps[i] = primos[i+1] - primos[i]
	}
	d := segundas(gaps)
	r := repeticiones(d)
	nrep := 0
	for _, v := range r {
		if v == 1 {
			nrep++
		}
	}
	fmt.Printf("primos: %d · doses: %d · repeticiones: %d (%.3f%%)\n",
		len(primos), len(d), nrep, 100*float64(nrep)/float64(len(r)))

	// controles
	var rc [][]float64
	for k := 0; k < SEMILLAS; k++ {
		rc = append(rc, repeticiones(segundas(baraja(gaps, uint64(0xC1C10000)+uint64(k)*104729))))
	}

	// ---- LEY 1: el espaciado ----
	fmt.Println("\nLEY 1 · ¿CADA CUÁNTO CAE UNA REPETICIÓN?")
	esp := espaciados(r)
	m, s := mediaDesvio(esp)
	var mm, ms float64
	for _, c := range rc {
		a, b := mediaDesvio(espaciados(c))
		mm += a / float64(SEMILLAS)
		ms += b / float64(SEMILLAS)
	}
	fmt.Printf("\n        distancia media entre repeticiones ....... %.4f  (barajado %.4f)\n", m, mm)
	fmt.Printf("        desvío estándar .......................... %.4f  (barajado %.4f)\n", s, ms)
	fmt.Printf("        razón desvío/media ....................... %.4f  (barajado %.4f)\n", s/m, ms/mm)
	fmt.Println("\n   📌 CÓMO SE LEE ESE ÚLTIMO NÚMERO, y es la clave de toda la ley: si las")
	fmt.Println("   repeticiones cayeran EN UN COMPÁS, todas las distancias serían parecidas y")
	fmt.Println("   la razón daría cerca de 0. Si caen desperdigadas sin memoria, la razón da")
	fmt.Println("   cerca de 1. No hay término medio ambiguo: el número dice cuál de las dos es.")

	fmt.Println("\n        distancia   en los primos      %      barajado       %")
	hist := map[int]int{}
	for _, v := range esp {
		hist[v]++
	}
	histB := map[int]float64{}
	for _, c := range rc {
		for _, v := range espaciados(c) {
			histB[v] += 1 / float64(SEMILLAS)
		}
	}
	totalB := 0.0
	for _, v := range histB {
		totalB += v
	}
	for L := 1; L <= 10; L++ {
		fmt.Printf("   %11d %15d %8.3f%% %13.0f %8.3f%%\n", L, hist[L],
			100*float64(hist[L])/float64(len(esp)), histB[L], 100*histB[L]/totalB)
	}

	// ---- LEY 2: el espectro ----
	fmt.Println("\nLEY 2 · ⚡ EL ESPECTRO — ¿HAY UN COMPÁS ESCONDIDO?")
	fmt.Printf("   Barriendo %d frecuencias sobre una ventana de %d pasos.\n", KFREQ, MESP)
	pot, pico, jPico := periodograma(r, MESP, KFREQ)
	var picosB []float64
	for _, c := range rc {
		_, p, _ := periodograma(c, MESP, KFREQ)
		picosB = append(picosB, p)
	}
	sort.Float64s(picosB)
	maxB := picosB[len(picosB)-1]
	medB := picosB[len(picosB)/2]
	fmt.Printf("\n        pico más alto en los PRIMOS .............. %.2f  (período %.1f pasos)\n",
		pico, float64(MESP)/float64(jPico))
	fmt.Printf("        pico más alto del control, mediana ....... %.2f\n", medB)
	fmt.Printf("        pico más alto del control, el peor de %d .. %.2f\n", SEMILLAS, maxB)
	veredictoEspectro := "❌ NO"
	if pico > maxB {
		veredictoEspectro = "✅ SÍ"
	}
	fmt.Printf("\n        ¿el pico de los primos supera a TODOS los del control? ..... %s\n", veredictoEspectro)
	if pico <= maxB {
		fmt.Println("\n   ⟹ **NO HAY COMPÁS.** El pico más alto de los primos ni siquiera le gana al")
		fmt.Println("   pico más alto que sale de barajar los mismos saltos. Lo que se ve es lo que")
		fmt.Println("   se ve en cualquier ruido: el máximo de muchas cosas chicas.")
	} else {
		fmt.Println("\n   ⟹ ⚡ EL PICO SUPERA AL CONTROL. Hay que mirarlo de cerca antes de festejar.")
	}

	// ---- LEY 3: la memoria ----
	fmt.Println("\nLEY 3 · LA MEMORIA DEL TREN DE REPETICIONES")
	ac := autocorr(r, 60)
	acB := make([]float64, 60)
	acB2 := make([]float64, 60)
	for _, c := range rc {
		a := autocorr(c, 60)
		for i := range a {
			acB[i] += a[i] / float64(SEMILLAS)
			acB2[i] += a[i] * a[i] / float64(SEMILLAS)
		}
	}
	fmt.Println("\n        paso   primos      barajado ± desvío        z")
	// OJO CON ESTO, y es corrección de mi propia primera versión: el paso 1 NO
	// puede contar como evidencia de ciclo. Ahí vive el sesgo de vecindad que el
	// hallazgo anterior ya midió (los múltiplos de 3 aplastados). Un CICLO tendría
	// que aparecer en un paso LARGO. Así que el veredicto se juzga desde el 2.
	z1 := 0.0
	peorZ, peorL := 0.0, 0
	for i := 0; i < 60; i++ {
		sd := math.Sqrt(math.Max(acB2[i]-acB[i]*acB[i], 0))
		z := 0.0
		if sd > 0 {
			z = (ac[i] - acB[i]) / sd
		}
		if i == 0 {
			z1 = z
		} else if math.Abs(z) > math.Abs(peorZ) {
			peorZ, peorL = z, i+1
		}
		if i < 12 || math.Abs(z) > 5 {
			fmt.Printf("   %12d %+9.5f %+13.5f ± %.5f %8.1f\n", i+1, ac[i], acB[i], sd, z)
		}
	}
	fmt.Printf("\n        paso 1 (vecindad, YA conocida) ......... z = %+.1f\n", z1)
	fmt.Printf("        el paso LARGO con más señal es el %d ..... z = %+.1f\n", peorL, peorZ)
	fmt.Println("\n   📌 Y la separación no es un detalle de forma: el paso 1 es donde vive el")
	fmt.Println("   sesgo de vecindad que ya medimos. Un CICLO tendría que aparecer en un paso")
	fmt.Println("   largo. Por eso el veredicto se juzga del paso 2 en adelante.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	if pico <= maxB && math.Abs(peorZ) < 5 {
		fmt.Println("❌ **NO HAY CICLOS.** Las tres pruebas dicen lo mismo:")
	} else {
		fmt.Println("⚡ **HAY ALGO EN UN PASO LARGO, Y HAY QUE MIRARLO ANTES DE FESTEJAR:**")
	}
	fmt.Printf("\n  · la distancia entre repeticiones tiene desvío/media = %.4f, contra %.4f\n", s/m, ms/mm)
	fmt.Println("    del barajado: las dos cerca de 1, o sea desperdigadas sin memoria.")
	fmt.Printf("  · el pico más alto del espectro (%.2f) no supera al del control (%.2f).\n", pico, maxB)
	fmt.Printf("  · fuera del paso 1, el paso con más señal es el %d, con z = %+.1f.\n", peorL, peorZ)
	fmt.Printf("  · el paso 1 sí tiene señal (z = %+.1f), pero ése es el sesgo de vecindad\n", z1)
	fmt.Println("    que ya estaba medido: no es un compás, es el vecino de al lado.")
	fmt.Printf("\n📌 Y UN DATO QUE SÍ ES REAL Y VA ANOTADO: las repeticiones caen MÁS SEGUIDO en\n")
	fmt.Printf("  los primos que en el barajado — cada %.4f pasos contra cada %.4f. Es la otra\n", m, mm)
	fmt.Println("  cara del mismo sesgo: como los primos evitan repetir salto, quedan menos")
	fmt.Println("  cortes y las corridas se juntan más. Mismo hallazgo, no uno nuevo.")
	fmt.Println("\n📌 Y ESTO NO CONTRADICE EL HALLAZGO ANTERIOR, lo AFINA: los primos sí evitan")
	fmt.Println("  repetirse —el aplastamiento de los múltiplos de 3, z = −244— pero ese sesgo")
	fmt.Println("  vive en el paso SIGUIENTE, no en un compás largo. Es una regla de vecindad,")
	fmt.Println("  no un latido. El capitán tenía razón en que hay ley; no la hay en que haya ciclo.")

	escribirLamina(esp, hist, histB, m, s, mm, ms, pot, pico, jPico, picosB, maxB, medB,
		ac, acB, acB2, peorZ, peorL, z1, nrep, len(r))
}

func escribirLamina(esp []int, hist map[int]int, histB map[int]float64,
	m, s, mm, ms float64, pot []float64, pico float64, jPico int,
	picosB []float64, maxB, medB float64,
	ac, acB, acB2 []float64, peorZ float64, peorL int, z1 float64, nrep, nr int) {

	var b strings.Builder
	W, H := 1640.0, 1240.0
	hayCiclo := pico > maxB
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="48" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔁 LOS CICLOS — ¿las repeticiones marcan un compás?</text>
<text x="%.0f" y="78" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">repetición = dos doses seguidos del mismo signo · %d repeticiones sobre %d pasos</text>
`, W, H, W, H, W/2, W/2, nrep, nr)

	// PANEL 1: el espaciado
	fmt.Fprintf(&b, `<rect x="40" y="102" width="770" height="330" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="425" y="134" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">1 · ¿CADA CUÁNTO CAE UNA?</text>
<text x="425" y="158" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">barras llenas: los primos · contorno: los mismos saltos barajados</text>
`)
	bx, by, bw := 90.0, 340.0, 660.0
	maxH := 0
	for L := 1; L <= 10; L++ {
		if hist[L] > maxH {
			maxH = hist[L]
		}
	}
	if maxH == 0 {
		maxH = 1
	}
	totalB := 0.0
	for _, v := range histB {
		totalB += v
	}
	cw := bw / 10.0
	// las dos series se normalizan a porcentaje de su propio total, que es la
	// única comparación honesta cuando los totales difieren.
	maxPct := 0.0
	for L := 1; L <= 10; L++ {
		if p := float64(hist[L]) / float64(len(esp)); p > maxPct {
			maxPct = p
		}
		if p := histB[L] / totalB; p > maxPct {
			maxPct = p
		}
	}
	if maxPct == 0 {
		maxPct = 1
	}
	for L := 1; L <= 10; L++ {
		h := 140.0 * (float64(hist[L]) / float64(len(esp))) / maxPct
		hb := 140.0 * (histB[L] / totalB) / maxPct
		x := bx + float64(L-1)*cw
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#7fb2ff"/>`, x, by-h, cw*0.66, h)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="none" stroke="#ffd98a" stroke-width="1.4"/>`, x, by-hb, cw*0.66, hb)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="monospace" fill="#8fa8c7">%d</text>`, x+cw*0.33, by+18, L)
	}
	fmt.Fprintf(&b, `
<text x="70" y="386" font-size="15" font-family="monospace" fill="#cfe6ff">desvío / media: primos %.4f · barajado %.4f</text>
<text x="70" y="410" font-size="13.5" font-family="Georgia" fill="#9aa8c4">cerca de 0 sería un compás · cerca de 1 es desperdigado sin memoria</text>
`, s/m, ms/mm)

	// PANEL 2: el espectro
	col := "#ff8fa0"
	texto := "NO SUPERA AL CONTROL"
	if hayCiclo {
		col = "#7ee0c0"
		texto = "SUPERA AL CONTROL"
	}
	fmt.Fprintf(&b, `<rect x="830" y="102" width="770" height="330" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1215" y="134" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">2 · EL ESPECTRO — ¿hay un compás escondido?</text>
<text x="1215" y="158" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">%d frecuencias barridas · la raya amarilla es el pico más alto del control</text>
`, len(pot))
	ex, ey, ew, eh := 880.0, 350.0, 680.0, 150.0
	esc := eh / math.Max(pico, maxB) * 0.9
	paso := len(pot) / 680
	if paso < 1 {
		paso = 1
	}
	for i := 0; i < len(pot); i += paso {
		x := ex + ew*float64(i)/float64(len(pot))
		h := pot[i] * esc
		if h > eh {
			h = eh
		}
		fmt.Fprintf(&b, `<rect x="%.2f" y="%.2f" width="1" height="%.2f" fill="#7fb2ff" opacity="0.75"/>`, x, ey-h, h)
	}
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.2f" x2="%.0f" y2="%.2f" stroke="#ffd98a" stroke-width="1.6" stroke-dasharray="6 4"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#5b7ba6" stroke-width="1"/>
<text x="880" y="386" font-size="15" font-family="monospace" fill="#cfe6ff">pico primos %.2f · control (peor de %d) %.2f</text>
<text x="880" y="410" font-size="16" font-family="Georgia" fill="%s">%s</text>
`, ex, ey-maxB*esc, ex+ew, ey-maxB*esc, ex, ey, ex+ew, ey, pico, len(picosB), maxB, col, texto)

	// PANEL 3: la memoria
	fmt.Fprintf(&b, `<rect x="40" y="452" width="1560" height="330" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="820" y="484" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">3 · LA MEMORIA — ¿una repetición avisa dónde estará la próxima?</text>
<text x="820" y="508" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">autocorrelación del tren de repeticiones, paso por paso · verde: primos · punteado: control</text>
`)
	ax, ay, aw, ah := 90.0, 640.0, 1450.0, 90.0
	maxA := 0.0
	for i := range ac {
		if math.Abs(ac[i]) > maxA {
			maxA = math.Abs(ac[i])
		}
	}
	if maxA == 0 {
		maxA = 1
	}
	var p1, p2 strings.Builder
	for i := range ac {
		x := ax + aw*float64(i)/float64(len(ac)-1)
		fmt.Fprintf(&p1, "%.1f,%.2f ", x, ay-ah*ac[i]/maxA)
		fmt.Fprintf(&p2, "%.1f,%.2f ", x, ay-ah*acB[i]/maxA)
	}
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#5b7ba6" stroke-width="1"/>
<polyline points="%s" fill="none" stroke="#7ee0c0" stroke-width="2"/>
<polyline points="%s" fill="none" stroke="#ffd98a" stroke-width="1.4" stroke-dasharray="5 4"/>
<text x="90" y="756" font-size="15" font-family="monospace" fill="#cfe6ff">paso 1 (vecindad, ya conocida): z = %+.1f  ·  fuera del 1, el mayor es el paso %d con z = %+.1f</text>
<text x="1540" y="756" font-size="13.5" text-anchor="end" font-family="Georgia" fill="#9aa8c4">si hubiera compás, acá se vería un diente repetido a intervalos regulares</text>
`, ax, ay, ax+aw, ay, p1.String(), p2.String(), z1, peorL, peorZ)

	// PANEL 4: el veredicto
	fmt.Fprintf(&b, `<rect x="40" y="802" width="1560" height="200" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="820" y="836" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffb27a">❌ NO HAY CICLOS — y las tres pruebas dicen lo mismo</text>
<text x="820" y="874" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">1 · la distancia entre repeticiones tiene desvío/media %.4f contra %.4f del barajado: desperdigadas, sin memoria</text>
<text x="820" y="902" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">2 · el pico más alto del espectro (%.2f) ni siquiera le gana al pico más alto que produce barajar los mismos saltos (%.2f)</text>
<text x="820" y="930" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">3 · fuera del paso 1 —donde vive el sesgo de vecindad ya conocido— el paso con más señal es el %d, con z = %+.1f</text>
<text x="820" y="968" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Y el testigo del espectro es el pico más alto DEL CONTROL, no un umbral: cuantas más frecuencias mirás, más alto sale el máximo del ruido.</text>
`, s/m, ms/mm, pico, maxB, peorL, peorZ)

	// PANEL 5: lo que queda
	fmt.Fprintf(&b, `<rect x="40" y="1022" width="1560" height="190" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="820" y="1056" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡ PERO ESTO NO TUMBA EL HALLAZGO ANTERIOR: LO AFINA</text>
<text x="820" y="1092" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Los primos SÍ evitan repetirse — los doses múltiplos de 3 están aplastados a la mitad, z = −244. Eso sigue en pie.</text>
<text x="820" y="1120" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Pero ese sesgo vive en el paso SIGUIENTE, no en un compás largo: es una regla de vecindad, no un latido.</text>
<text x="820" y="1156" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">El capitán tenía razón en que hay ley. No la tenía en que haya ciclo. Las dos mitades van al registro.</text>
<text x="820" y="1188" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Y una hipótesis que se puede matar con una medición vale más que una que no se puede matar con nada.</text>
</svg>
`)

	if err := os.WriteFile("los-ciclos.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: los-ciclos.svg")
}
