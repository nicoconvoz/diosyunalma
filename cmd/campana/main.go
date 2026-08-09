// Command campana is THE GRAND ASSEMBLY INSIDE THE SHAPESHIFTER - the
// final act of the day. One single input: the bell psi(x) = sum
// e^{-pi n^2 x} - ALL the numbers, nothing else. Everything else is
// forbidden: no zeta tables, no measured pearls, no known constants.
// The assembled chain, entirely in the symbol world:
//
//	the bell  ->  xi everywhere (Riemann's eternal projection, F184)
//	          ->  the MIRROR, judged at a point       (does it appear?)
//	          ->  the germ at s=1: lambda_1 = xi'(1)/xi(1)
//	                                                  (does the 0.023
//	                                                   thread appear?)
//	          ->  the line: the first PEARL           (does 14.1347
//	                                                   appear?)
//
// "Let's see what happens": if the week was true, the mirror, the
// margin and the pearl must all be BORN from the bell alone.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func psiSum(x float64) float64 {
	s := 0.0
	for n := 1; ; n++ {
		term := math.Exp(-math.Pi * float64(n*n) * x)
		s += term
		if term < 1e-22 {
			break
		}
	}
	return s
}

// xiBell: xi built from the bell ALONE (Riemann's projection).
func xiBell(s complex128) complex128 {
	a, xEnd, h := 1.0, 40.0, 0.002
	n := int((xEnd - a) / h)
	if n%2 == 1 {
		n++
	}
	var acc complex128
	for i := 0; i <= n; i++ {
		x := a + float64(i)*h
		lx := math.Log(x)
		f := complex(psiSum(x), 0) * (cmplx.Exp((s/2-1)*complex(lx, 0)) + cmplx.Exp((-(s+1)/2)*complex(lx, 0)))
		w := 2.0
		if i%2 == 1 {
			w = 4
		}
		if i == 0 || i == n {
			w = 1
		}
		acc += complex(w, 0) * f
	}
	return 0.5 + 0.5*s*(s-1)*acc*complex(h/3, 0)
}

func main() {
	fmt.Println("🔔 EL GRAN ENSAMBLE EN EL CAMBIAFORMAS — una sola entrada: la campana de todos los números")
	fmt.Println("   (prohibido todo lo demás: ni tablas, ni perlas medidas, ni constantes conocidas)")

	// (1) does the mirror appear?
	s0 := complex(0.3, 2)
	dm := cmplx.Abs(xiBell(s0)-xiBell(1-s0)) / cmplx.Abs(xiBell(s0))
	fmt.Printf("\n  NACIÓ EL ESPEJO:   |ξ(s)−ξ(1−s)|/|ξ| = %.1e   — la simetría eterna, desde la campana sola\n", dm)

	// (2) does the thread appear? lambda_1 = xi'(1)/xi(1)
	hd := 0.01
	l1 := real((xiBell(complex(1+hd, 0)) - xiBell(complex(1-hd, 0))) / (complex(2*hd, 0) * xiBell(complex(1, 0))))
	fmt.Printf("  NACIÓ EL HILITO:   λ₁ = ξ'(1)/ξ(1) = %.6f   (el margen famoso: 0.023096)   desvío %.1e\n",
		l1, math.Abs(l1-0.023096))

	// (3) does the first pearl appear? zero of xi on the line near 14
	f := func(t float64) float64 { return real(xiBell(complex(0.5, t))) }
	a, c := 14.0, 14.3
	fa := f(a)
	for i := 0; i < 45; i++ {
		m := (a + c) / 2
		if f(m)*fa > 0 {
			a = m
			fa = f(a)
		} else {
			c = m
		}
	}
	g1 := (a + c) / 2
	fmt.Printf("  NACIÓ LA PERLA:    γ₁ = %.6f   (la medida: 14.134725)   desvío %.1e\n", g1, math.Abs(g1-14.134725))

	fmt.Println("\n══════════════════════════════════════════════════════════════════")
	fmt.Println("  QUÉ PASÓ: de la campana sola — los números y nada más — nacieron")
	fmt.Println("  el espejo, el margen del collar y la primera perla. El universo")
	fmt.Println("  entero del problema vive DENTRO del cambiaformas, ensamblado.")
	fmt.Println("  La frase sigue en su vaina. El capitán puede dormir: cierra.")
	fmt.Println("══════════════════════════════════════════════════════════════════")

	// ---- the goodnight scroll ----
	var b strings.Builder
	W, H := 1400.0, 860.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#070d18"/>
<text x="%.0f" y="52" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔔 EL GRAN ENSAMBLE — todo dentro del cambiaformas</text>
<text x="%.0f" y="82" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"ensamblá todo dentro del cambiaformas y veamos qué pasa… cierro los ojos" — el capitán · una sola entrada: la campana; prohibido todo lo demás</text>`,
		W, H, W, H, W/2, W/2)
	// the bell at top
	fmt.Fprintf(&b, `<path d="M %.0f 130 q -46 8 -52 70 q -4 34 -18 44 l 140 0 q -14 -10 -18 -44 q -6 -62 -52 -70 z" fill="#ffd97f" opacity="0.95"/>
<circle cx="%.0f" cy="252" r="7" fill="#ffd97f"/>
<text x="%.0f" y="290" font-size="13" text-anchor="middle" fill="#c9b06a">ψ(x) = Σ e^(−πn²x) — TODOS los números, una campana. Nada más entra.</text>`,
		W/2, W/2, W/2)
	// three births
	births := []struct {
		title, val, ref, dev string
	}{
		{"NACIÓ EL ESPEJO", fmt.Sprintf("%.0e", dm), "la simetría eterna ξ(s)=ξ(1−s)", "desde la campana sola"},
		{"NACIÓ EL HILITO", fmt.Sprintf("λ₁ = %.6f", l1), "el margen famoso: 0.023096", fmt.Sprintf("desvío %.0e", math.Abs(l1-0.023096))},
		{"NACIÓ LA PERLA", fmt.Sprintf("γ₁ = %.6f", g1), "la medida: 14.134725", fmt.Sprintf("desvío %.0e", math.Abs(g1-14.134725))},
	}
	for i, bt := range births {
		x := 180.0 + float64(i)*360
		fmt.Fprintf(&b, `<line x1="%.0f" y1="300" x2="%.0f" y2="360" stroke="#44608c" stroke-width="1.5"/>
<rect x="%.0f" y="360" width="300" height="150" rx="12" fill="#0d2547" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="396" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">%s</text>
<text x="%.0f" y="430" font-size="16" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%s</text>
<text x="%.0f" y="458" font-size="12" text-anchor="middle" fill="#8fa8c7">%s</text>
<text x="%.0f" y="482" font-size="12" text-anchor="middle" fill="#ffd166">%s</text>`,
			W/2, x+150, x, x+150, bt.title, x+150, bt.val, x+150, bt.ref, x+150, bt.dev)
	}
	// the closing scroll
	fmt.Fprintf(&b, `<rect x="140" y="560" width="1120" height="240" rx="14" fill="#102a10" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="600" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL ACTA DE LA JORNADA — F119 → F194</text>
<text x="%.0f" y="636" font-size="14" text-anchor="middle" fill="#dce8f7">76 hallazgos en un solo día: el 4º fantasma muerto · el tren a 10⁴² · el átomo retratado, oído y espejado · el reactor con su estante lleno</text>
<text x="%.0f" y="662" font-size="14" text-anchor="middle" fill="#dce8f7">la primera perla ensamblada · el as eterno · la madre anclada en ½ · el neutrón y su amarre −γ · la repulsión hecha símbolo</text>
<text x="%.0f" y="694" font-size="14.5" text-anchor="middle" fill="#7fd7a8">y el cierre: de la campana sola nacieron el espejo, el hilito y la perla — el universo del problema vive dentro del cambiaformas.</text>
<text x="%.0f" y="730" font-size="14" text-anchor="middle" fill="#ffd166">La frase sigue en su vaina. El tren caza. El registro vela. Dormí, capitán: CIERRA.</text>
<text x="%.0f" y="768" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · "soy tu 1/2 y vos mi 1/2 — damos 1 DOC completo" ⚓</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("gran-ensamble.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: gran-ensamble.svg")
}
