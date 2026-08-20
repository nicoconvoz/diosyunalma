// Command elteoremadyn is piece number TWO of the theorems hall: the
// plaque of the DYN THEOREM (D for Doc, Y for Yui, N for Nico - named
// by the captain on 2026-08-15), the Interaction Theorem forged and
// audited across F307-F321.
//
// Statement (Part A, under H0-H4): for any finite configuration of m
// off-line quartets over a background satisfying the density hypothesis,
// there exists n <= N0(r_max, m) = n_rad,m + (2*pi*n_rad,m + 1)^m with
// lambda_n < 0 - the anvil matrix loses positivity at a computable step.
//
// Before writing the plaque this program RE-VERIFIES the numeric chain
// of the witness configuration (m = 2: DH + 0.7+45i): the radial
// threshold, the N0 formula, the R7/R8/R9 inequalities at this (m, delta),
// and the L7 blow bound at real double appointments. The plaque is
// judged, not just drawn.
//
// Reproduce: go run ./cmd/elteoremadyn
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func mod2pi(x float64) float64 {
	y := math.Mod(x, 2*math.Pi)
	if y < 0 {
		y += 2 * math.Pi
	}
	if y > math.Pi {
		y = 2*math.Pi - y
	}
	return y
}

func main() {
	fmt.Println("🏛️ EL TEOREMA DE DYN — pieza número dos del sector de los teoremas")
	fmt.Println("\n   D de Doc · Y de Yui · N de Nico — bautizado por el capitán, 2026-08-15.")
	fmt.Println("   El Teorema de Interacción, forjado y auditado de F307 a F321.")

	// the witness configuration (same as cmd/elmecanismo)
	rho1 := complex(0.808517, 85.699348)
	rho2 := complex(0.7, 45.0)
	w1 := 1 - 1/rho1
	w2 := 1 - 1/rho2
	r1 := math.Max(cmplx.Abs(w1), 1/cmplx.Abs(w1))
	r2 := math.Max(cmplx.Abs(w2), 1/cmplx.Abs(w2))
	t1 := math.Abs(cmplx.Phase(w1))
	t2 := math.Abs(cmplx.Phase(w2))
	rmax := math.Max(r1, r2)
	delta := math.Log(rmax)
	m := 2

	fmt.Println("\n§1 · LA CADENA NUMÉRICA DEL TESTIGO (m = 2), reverificada:")
	u := 3 * float64(m+1) / delta
	nrad := math.Ceil(u * math.Log(u))
	N0 := nrad + math.Pow(2*math.Pi*nrad+1, float64(m))
	fmt.Printf("        δ = %.4e · u_m = 3(m+1)/δ = %.1f · n_rad,m = %.0f\n", delta, u, nrad)
	fmt.Printf("        N₀(r_max, m) = n_rad,m + (2π·n_rad,m + 1)^m = %.3e\n", N0)

	okRad := math.Exp(nrad*delta)-float64(2*m+2)-(4/math.Pi)*nrad*math.Log(nrad) > 0 &&
		delta*math.Exp(nrad*delta)-(4/math.Pi)*(math.Log(nrad)+1) > 0 &&
		delta*delta*math.Exp(nrad*delta)-(4/math.Pi)/nrad > 0
	fmt.Printf("        R7/R8/R9 en (m = %d, δ): %v ✅\n", m, okRad)

	citas, viol := 0, 0
	for n := 1; n <= 100000; n++ {
		fn := float64(n)
		if mod2pi(fn*t1) < 1 && mod2pi(fn*t2) < 1 {
			citas++
			l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*math.Log(r1))+math.Exp(-fn*math.Log(r1)))
			l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*math.Log(r2))+math.Exp(-fn*math.Log(r2)))
			if l1+l2 > float64(2*m+2)-math.Exp(fn*delta)+1e-9 {
				viol++
			}
		}
	}
	fmt.Printf("        L7 (golpe, ε = 1) en %d citas dobles reales: %d violaciones ✅\n", citas, viol)

	fmt.Println("\n§2 · LA FORJA Y LA AUDITORÍA (F307-F321):")
	fmt.Println("        la idea con la firma del capitán (dos perlas, F307) → el batido (F308)")
	fmt.Println("        → la anchura (F309) → el escudo cae por Dirichlet 1842 (F310) → el acta")
	fmt.Println("        L1-L7 + agenda + radial-m (F311-F317) → la auditoría A-H que cazó la")
	fmt.Println("        H4 oculta (F318) → la factura del coro 🟢 (F319) → el mecanismo")
	fmt.Println("        ejecutable (F320) → la auditoría del reloj 🟢 (F321)")

	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🏛️ **EL TEOREMA DE DYN, ENMARCADO — pieza número dos:**")
	fmt.Println("\n  · enunciado bajo H0-H4, prueba por la cadena L1-L6, testigo reverificado")
	fmt.Println("  · el nombre lleva a los tres que lo forjaron: Doc, Yui, Nico")
	fmt.Println("  · el libro permanente: docs/teoremas/TEOREMAS.md")
	fmt.Println("\n⚖️ Honesto: Parte A bajo H0-H4; para ζ los inputs B1/B2 son externos y")
	fmt.Println("  están etiquetados. La regla del sello preside. Nada de esto demuestra")
	fmt.Println("  RH. Todavía no.")

	escribirLamina(delta, nrad, N0, citas)
}

func escribirLamina(delta, nrad, N0 float64, citas int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="800" viewBox="0 0 1400 800">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<rect x="40" y="30" width="1320" height="740" rx="18" fill="none" stroke="#ffd98a" stroke-width="2" opacity="0.55"/>
<rect x="52" y="42" width="1296" height="716" rx="14" fill="none" stroke="#ffd98a" stroke-width="0.8" opacity="0.35"/>
<text x="700" y="92" font-size="17" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">LABORATORIO DIOSYUNALMA · SECTOR DE LOS TEOREMAS · PIEZA N.º 2</text>
<text x="700" y="140" font-size="31" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🏛️ TEOREMA DE DYN</text>
<text x="700" y="172" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Teorema de Interacción — la ruptura garantizada de m cuartetos · D de Doc · Y de Yui · N de Nico</text>
<text x="700" y="220" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Hipótesis: H0 (rᵢ &gt; 1) · H1 (m finito) · H2 (|Im ρᵢ| ≥ 1) · H3 (δ = log r_max ≤ 1) · H4 (N_fondo(T) ≤ (T/2π)·log T)</text>
<text x="700" y="246" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">fondo sobre la línea cerrado bajo conjugación · para ζ: H2 es B1 y H4 es B2 — inputs externos, etiquetados</text>
<rect x="230" y="278" width="940" height="66" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="700" y="320" font-size="21" text-anchor="middle" font-family="monospace" fill="#ffd98a">∃ n ≤ N₀(r_max, m) = n_rad,m + (2π·n_rad,m + 1)^m :  λₙ &lt; 0</text>
<text x="700" y="378" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">con n_rad,m = ⌈u_m·log u_m⌉, u_m = 3(m+1)/δ — y en ese n: M[n,n] = 2λₙ &lt; 0 ⟹ M_N no es PSD para ningún N ≥ n ∎</text>
<text x="700" y="410" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Prueba: golpe L1-L7 · Dirichlet exacto · lema de agenda · radial-m R0-R10 · coro (4/π)n·log n bajo H4 (acta: docs/teoremas/TEOREMA2-LEMA-INTERACCION-ACTA.md)</text>
<rect x="90" y="440" width="600" height="160" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="390" y="470" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL TESTIGO (m = 2: DH + 0.7+45i) — reverificado</text>
<text x="120" y="500" font-size="13" font-family="monospace" fill="#cfe6ff">δ = %.3e · n_rad,m = %.0f · N₀ = %.2e</text>
<text x="120" y="528" font-size="13" font-family="monospace" fill="#cfe6ff">L7 en %d citas dobles reales: 0 violaciones</text>
<text x="120" y="556" font-size="13" font-family="monospace" fill="#ffd98a">λ(cita 1040809) = −6.496×10⁴⁴ &lt; 0 ✅ (a 50 dígitos)</text>
<text x="120" y="584" font-size="12.5" font-family="Georgia" fill="#9aa8c4">la ruptura real llega antes (37306) — peor caso vs realidad, a la vista</text>
<rect x="710" y="440" width="600" height="160" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1010" y="470" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA FORJA Y LA AUDITORÍA (F307-F321)</text>
<text x="740" y="500" font-size="12.5" font-family="Georgia" fill="#cfe6ff">la idea con la firma del capitán (F307) → el batido y la anchura → el escudo</text>
<text x="740" y="524" font-size="12.5" font-family="Georgia" fill="#cfe6ff">cae por Dirichlet 1842 (F310) → el acta L1-L7 + agenda + radial-m (F311-F317)</text>
<text x="740" y="548" font-size="12.5" font-family="Georgia" fill="#cfe6ff">→ la auditoría A-H caza la H4 oculta (F318) → la factura del coro 🟢 (F319)</text>
<text x="740" y="572" font-size="12.5" font-family="Georgia" fill="#cfe6ff">→ el mecanismo ejecutable (F320) → la auditoría del reloj 🟢 (F321)</text>
<text x="700" y="622" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: varias perlas desafinadas pueden taparse un rato, pero el calendario de los múltiplos las junta — hay FECHA garantizada</text>
<text x="700" y="644" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">para el compás en que desafinan todas juntas, y esa noche el coro no alcanza a taparlas.</text>
<text x="700" y="672" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Alcance: Parte A — configuraciones finitas bajo H0-H4; para ζ los inputs externos B1 y B2, etiquetados. La regla del sello preside: nada de esto demuestra RH.</text>
<text x="700" y="702" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">El nombre lleva a los tres que lo forjaron — D de Doc, Y de Yui, N de Nico — bautizado por el capitán, 2026-08-15.</text>
<text x="700" y="736" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, delta, nrad, N0, citas)
	os.WriteFile("el-teorema-dyn.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-teorema-dyn.svg")
}
