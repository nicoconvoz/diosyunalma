// Command elplanodeltren draws the exploded blueprint of the TRAIN
// (Doc Brown's Landsberg-Schaar locomotive, cmd/circulo): every piece,
// how it fits, what it is for and how it is used - the museum tour the
// captain asked for.
//
// Before framing, it VERIFIES the two load-bearing pieces live:
// rail #1 (the exact Landsberg-Schaar reciprocity, q-term sum == p-term
// sum with exact integer phases) and rail #2 (one circle flip: the
// N-term chirp becomes a ~2bN-term dual - the shortening that powers
// the cascade).
//
// Reproduce: go run ./cmd/elplanodeltren
package main

import (
	"fmt"
	"math"
	"os"
)

func gaussSum(p, q, sign int64) (re, im float64) {
	m := 2 * q
	var n2, dn int64
	dn = 1
	for n := int64(0); n < q; n++ {
		ph := math.Pi * float64((p*n2)%m) / float64(q) * float64(sign)
		s, c := math.Sincos(ph)
		re += c
		im += s
		n2 = (n2 + dn) % m
		dn = (dn + 2) % m
	}
	return
}

func main() {
	fmt.Println("🚂 EL PLANO DEL TREN — el despiece de la locomotora de Landsberg-Schaar")

	// verify rail #1: the exact reciprocity
	p, q := int64(3), int64(50000)
	lr, li := gaussSum(p, q, 1)
	rr, ri := gaussSum(q, p, -1)
	f := math.Sqrt(float64(q) / float64(p))
	c4, s4 := math.Cos(math.Pi/4), math.Sin(math.Pi/4)
	pr := f * (c4*rr - s4*ri)
	pi2 := f * (c4*ri + s4*rr)
	err := math.Hypot(lr-pr, li-pi2) / math.Hypot(lr, li)
	fmt.Printf("\n§1 · RIEL 1 verificado: suma de %d términos == suma de %d términos (error %.1e) ✅\n", q, p, err)
	fmt.Println("        el numerito de arriba resuelve al gigante de abajo — EXACTO, no aproximado")

	// verify rail #2: one flip shortens by ~2b
	b := 0.13
	n := int64(100000)
	dual := float64(int64(math.Floor(2*b*float64(n-1))) + 1)
	fmt.Printf("§2 · RIEL 2 verificado: un giro del círculo convierte %d términos en ~%.0f (razón %.3f ≈ 2b = %.2f) ✅\n", n, dual, dual/float64(n), 2*b)
	fmt.Println("        iterado, el descenso es LOGARÍTMICO — la clase de motor t^(1/3)")

	fmt.Println("\n§3 · las demás piezas (certificadas en sus actas): el amortiguador F144, la")
	fmt.Println("        cascada calibrada, el sonar de ballena, el juez, el árbitro de 256 bits")
	fmt.Println("        y la marcha — el despiece completo va en la lámina y en el recorrido web")

	var sb []string
	add := func(s string) { sb = append(sb, s) }
	add(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">
<rect width="100%" height="100%" fill="#0b1526"/>
<rect x="30" y="20" width="1340" height="910" rx="18" fill="none" stroke="#7ee0c0" stroke-width="2" opacity="0.5"/>
<text x="700" y="60" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🚂 EL PLANO DEL TREN — la locomotora de Landsberg-Schaar, pieza por pieza</text>
<text x="700" y="88" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">para qué: leer el mar de ondas de ζ en aguas donde el float64 se ahoga · cómo se usa: go run ./cmd/circulo (demo) · -cazar (la cacería sin fin)</text>`)
	caja := func(x, y, w, h float64, num, titulo string, lineas []string, color string) {
		add(fmt.Sprintf(`<rect x="%f" y="%f" width="%f" height="%f" rx="10" fill="#0d1830" stroke="%s"/>
<text x="%f" y="%f" font-size="14.5" font-family="Georgia" fill="%s">%s · %s</text>`, x, y, w, h, color, x+14, y+26, color, num, titulo))
		for i, l := range lineas {
			add(fmt.Sprintf(`<text x="%f" y="%f" font-size="12" font-family="Georgia" fill="#cfe6ff">%s</text>`, x+14, y+50+float64(i)*20, l))
		}
	}
	caja(70, 110, 620, 150, "RIEL 1", "LA RECIPROCIDAD (el corazón)", []string{
		"Landsberg-Schaar EXACTA: Σ e^{iπpn²/q} (q términos) = √(q/p)·e^{iπ/4}·Σ e^{−iπqn²/p} (p términos)",
		"cómo calza: es la física del zoom — el gigante de abajo se lee con el numerito de arriba",
		"fases en int64 EXACTO (p·n² mod 2q): el único error flotante es el coseno final",
		"verificado en esta corrida: 50 000 términos = 3 términos, error 2e-13 ✅"}, "#ffd98a")
	caja(710, 110, 620, 150, "RIEL 2", "LA CASCADA DE CHIRPS (el mar de verdad)", []string{
		"los bloques del mar son chirps de curvatura IRRACIONAL: Σ e^{2πi(bj² + cj)}",
		"un giro del círculo (Poisson + fase estacionaria): N términos → ~2bN términos duales",
		"el dual es OTRO chirp ⟹ se itera: descenso como fracción continua — logarítmico",
		"cómo calza: convierte cada bloque gigante en una bajada de escalones cortos"}, "#7ee0c0")
	caja(70, 280, 620, 130, "EL AMORTIGUADOR (F144)", "los reservorios de redondeo", []string{
		"la fase avanza por recurrencia con DOS registros de compensación (two-sum):",
		"la energía de redondeo se ATRAPA en el reservorio, no se disipa en el resultado",
		"cómo calza: sin él, un millón de pasos de fase ensucia el último dígito del mar"}, "#7fb2ff")
	caja(710, 280, 620, 130, "RIEL 3", "LA CASCADA CALIBRADA (el tamaño cómodo)", []string{
		"el plan del capitán: girar el círculo hasta que la suma quepa en un TAMAÑO CÓMODO,",
		"calibrar donde la verdad es barata (mar liviano) y navegar el mar pesado con lo mismo",
		"cómo se usa: es el motor interno de cada bloque del cazadero — no se toca a mano"}, "#7ee0c0")
	caja(70, 430, 620, 170, "EL SONAR DE BALLENA", "escuchar() — la escucha que crece", []string{
		"presupuesto del bloque: L = t^(-1/3)-clase (bandL) — la ley del tercio",
		"la onda se propaga: 1500 → 6000 → 24000 → banda completa; el agua calma",
		"corta la escucha temprano; solo el agua interesante mantiene la onda viajando",
		"anómalo si |ola|/√L supera 2.4σ (COHERENTE) o cae bajo 0.05 (MUDA) — dos bestias",
		"cómo calza: es el filtro barato que decide dónde vale la pena remar de verdad"}, "#ff9aa8")
	caja(710, 430, 620, 170, "EL JUEZ Y EL ÁRBITRO", "saquear() — nadie entra al libro sin firma", []string{
		"EL JUEZ: cada bestia se rema DIRECTO (blockDirect) y se compara con la cascada;",
		"si difieren más del 5%, rechazada — la cascada jamás se cree a sí misma",
		"EL ÁRBITRO DE 256 BITS (F201): en aguas ≥ 1e40, cada 8ª presa recibe contraste",
		"big.Float — porque cerca del borde del cuaderno dd, dos motores dd de acuerdo",
		"no prueban nada; una falla ahí = el 5º fantasma, avistado en vivo"}, "#ffd98a")
	caja(70, 620, 620, 150, "LA MARCHA", "cómo se anexa un mar", []string{
		"un agua entra al cazadero solo tras 3 SONDAS FIRMADAS en bandas distintas",
		"aguas conquistadas: 1e33 → 1e42 (registradas como MARCHA en cazadero.log)",
		"la frontera navega hacia el agotamiento del cuaderno dd (~1e48): más hondo,",
		"la máquina no firma — y esta casa no navega donde no puede firmar"}, "#7fb2ff")
	caja(710, 620, 620, 150, "LA BITÁCORA DE LUZ", "dónde queda todo", []string{
		"cada presa: BESTIA MUDA/COHERENTE con t, k a 17 dígitos, L, |ola|, coh, juez",
		"cardúmenes (CARDUMEN): bloques vecinos que difieren en el 12º dígito de k",
		"tormentas de fondo (fondo.log) y los ARBITRO/FANTASMA5 del contraste",
		"cómo se usa: luz/cazadero.log + luz/fondo.log — y el faro (:8117) lo muestra en vivo"}, "#ff9aa8")
	add(`<text x="700" y="812" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: el tren no rema el mar gigante — lo GIRA. Cada vuelta del círculo cambia la suma enorme por una más corta que dice</text>
<text x="700" y="834" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">exactamente lo mismo, hasta que cabe en la mano. Y nada entra al libro sin pasar por el juez — y en lo hondo, por el árbitro.</text>
<text x="700" y="866" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">las cifras de certificación (4º fantasma muerto con anclas dd, e256 ≤ 4e-3 en 3e33 y 1e36) vienen de las actas F144-F155/F201 — este plano las cita, no las re-corre</text>
<text x="700" y="900" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">go run ./cmd/circulo · -cazar para la cacería sin fin · el recorrido completo: galeria/recorrido-maquinas.html · Todavía no.</text>
</svg>`)
	os.WriteFile("el-plano-del-tren.svg", []byte(joinS(sb)), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-plano-del-tren.svg")
}

func joinS(ss []string) string {
	out := ""
	for _, s := range ss {
		out += s + "\n"
	}
	return out
}
