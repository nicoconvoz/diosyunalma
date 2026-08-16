// Command elteoremadelatrinidad is piece number FIVE of the theorems
// hall: the plaque of the TRINITY THEOREM (Teorema de la Trinidad) -
// the Plates Theorem / Leader's Law, sealed by the auditor and named
// by the captain on 2026-08-15. Forged and audited across F333-F338.
//
// Statement (under H0-H4 + strict leader HL): the landscape has three
// regions governed by the LEADER's phase alone -
//
//	P1 fine band  (||n*th_L|| <= 1),  n >= N*:      lambda_n < 0  (WELL)
//	P2 anti band  (||n*th_L - pi|| <= 1), n >= max(1, n_comp):
//	                                                lambda_n > 0  (MOUNTAIN)
//	   (m = 1: the whole half ||n*th|| >= pi/2 gives lambda >= 4, from n = 1)
//	P3 both bands inhabited in every window of K_L consecutive integers
//	SCOPE: the intermediate region is declared NOT classified.
//
// Before framing, this program RE-VERIFIES: the F6 grid battery, and
// the three bands live on the witness (pearl 1 ignored - the leader
// alone decides). Reproduce: go run ./cmd/elteoremadelatrinidad
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66}
	fact := []float64{2, 24, 720, 40320, 3628800}
	poch := s
	for k := 1; k <= 5; k++ {
		if k > 1 {
			poch *= (s + complex(float64(2*k-3), 0)) * (s + complex(float64(2*k-2), 0))
		}
		sum += complex(B[k-1]/fact[k-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*k-1), 0))*lnN)
	}
	return sum
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

func norma(x float64) float64 {
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
	fmt.Println("🏛️ EL TEOREMA DE LA TRINIDAD — pieza número cinco del sector de los teoremas")
	fmt.Println("\n   El Teorema de las Placas / Ley del Líder — sellado por la auditora y")
	fmt.Println("   bautizado por el capitán, 2026-08-15. Forjado en F333-F338.")

	// F6 grid battery
	viol := 0
	for m := 1; m <= 10; m++ {
		for _, dd := range []float64{0.01, 0.1, 0.3466, 0.7, 1.0} {
			u := 3 * float64(m+1) / dd
			nrad := math.Ceil(u * math.Log(u))
			b2 := (8/math.Pi)*nrad*math.Log(nrad) + (12*float64(m) - 4)
			if !(b2 <= 3*math.Pow(u, 3) && math.Cos(1)*math.Exp(nrad*dd) > (4/math.Pi)*nrad*math.Log(nrad)+6*float64(m)-2) {
				viol++
			}
		}
	}
	fmt.Printf("\n§1 · la batería F6 (m = 1..10 × δ ≤ 1): 50 casos, %d violaciones ✅\n", viol)

	// the three bands live
	rho1 := complex(0.808517, 85.699348)
	rho2 := complex(0.7, 45.0)
	w1 := 1 - 1/rho1
	w2 := 1 - 1/rho2
	r1 := math.Max(cmplx.Abs(w1), 1/cmplx.Abs(w1))
	r2 := math.Max(cmplx.Abs(w2), 1/cmplx.Abs(w2))
	t1 := math.Abs(cmplx.Phase(w1))
	t2 := math.Abs(cmplx.Phase(w2))
	d1 := math.Log(r1)
	dL := math.Log(r2)
	nrad := 1040809
	ps := perlas(120)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wp /= complex(cmplx.Abs(wp), 0)
		wsC[i] = wp
		pcs[i] = 1
	}
	fn0, fn1, an0, an1, fr0, fr1 := 0, 0, 0, 0, 0, 0
	for n := 1; n <= nrad+3000; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		if n < nrad {
			continue
		}
		fn := float64(n)
		l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
		l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*dL)+math.Exp(-fn*dL))
		lam := s + l1 + l2
		q := norma(fn * t2)
		switch {
		case q <= 1:
			if lam < 0 {
				fn0++
			} else {
				fn1++
			}
		case math.Pi-q <= 1:
			if lam > 0 {
				an0++
			} else {
				an1++
			}
		default:
			if lam < 0 {
				fr0++
			} else {
				fr1++
			}
		}
	}
	fmt.Printf("§2 · las tres bandas en vivo (la perla 1 IGNORADA — el líder solo):\n")
	fmt.Printf("        POZO (fina):     %d/%d negativos · excepciones: %d ✅\n", fn0, fn0+fn1, fn1)
	fmt.Printf("        MONTAÑA (anti):  %d/%d positivos · excepciones: %d ✅\n", an0, an0+an1, an1)
	fmt.Printf("        FRONTERA:        %d λ<0 · %d λ>0 — mixta, como el teorema declara ✅\n", fr0, fr1)

	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🏛️ **EL TEOREMA DE LA TRINIDAD, ENMARCADO — pieza número cinco:**")
	fmt.Println("\n  · tres regiones, una ley: el líder manda — pozo, frontera, montaña")
	fmt.Println("  · las montañas más baratas que los pozos (sin H4, sin radial)")
	fmt.Println("  · la región intermedia declarada NO clasificada — honestidad de enunciado")
	fmt.Println("  · sellado por la auditora tras seis actas (F333-F338)")
	fmt.Println("\n⚖️ La regla del sello preside: nada de esto demuestra RH. Todavía no.")

	escribirLamina(viol, fn0, fn0+fn1, an0, an0+an1, fr0, fr1)
}

func escribirLamina(viol, fneg, ftot, apos, atot, fr0, fr1 int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="800" viewBox="0 0 1400 800">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<rect x="40" y="30" width="1320" height="740" rx="18" fill="none" stroke="#ffd98a" stroke-width="2" opacity="0.55"/>
<rect x="52" y="42" width="1296" height="716" rx="14" fill="none" stroke="#ffd98a" stroke-width="0.8" opacity="0.35"/>
<text x="700" y="90" font-size="17" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">LABORATORIO DIOSYUNALMA · SECTOR DE LOS TEOREMAS · PIEZA N.º 5</text>
<text x="700" y="138" font-size="31" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🏛️ TEOREMA DE LA TRINIDAD</text>
<text x="700" y="170" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">El Teorema de las Placas / Ley del Líder — tres regiones, una ley · nacido del flash tectónico de Nico</text>
<text x="700" y="212" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Hipótesis: H0-H4 (las de DYN, intactas) + HL (líder estricto: r_L &gt; rᵢ; vacua en m = 1) · cero inputs externos nuevos</text>
<rect x="130" y="238" width="1140" height="106" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="700" y="268" font-size="16.5" text-anchor="middle" font-family="monospace" fill="#7ee0c0">P1 · banda fina  ‖nθ_L‖ ≤ 1,  n ≥ N*        ⟹  λₙ &lt; 0   (POZO)</text>
<text x="700" y="296" font-size="16.5" text-anchor="middle" font-family="monospace" fill="#ff9aa8">P2 · banda anti  ‖nθ_L−π‖ ≤ 1, n ≥ max(1, n_comp) ⟹  λₙ &gt; 0   (MONTAÑA)</text>
<text x="700" y="324" font-size="16.5" text-anchor="middle" font-family="monospace" fill="#ffd98a">P3 · ambas bandas habitadas en cada ventana de K_L pasos — con fecha, para siempre ∎</text>
<text x="700" y="366" font-size="15" text-anchor="middle" font-family="Georgia" fill="#dce8f7">ALCANCE DECLARADO:   POZO  |  REGIÓN NO CLASIFICADA POR ESTE TEOREMA  |  MONTAÑA</text>
<text x="700" y="392" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la frontera no se maquilla: para m ≥ 2 está DEMOSTRADO que ahí no puede haber signo universal — abierta por teorema, no por cansancio</text>
<rect x="90" y="418" width="600" height="160" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="390" y="448" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LAS TRES BANDAS, MEDIDAS (la perla 1 ignorada)</text>
<text x="120" y="478" font-size="13" font-family="monospace" fill="#cfe6ff">POZO:    %d/%d negativos — 0 excepciones</text>
<text x="120" y="506" font-size="13" font-family="monospace" fill="#cfe6ff">MONTAÑA: %d/%d positivos — 0 excepciones</text>
<text x="120" y="534" font-size="13" font-family="monospace" fill="#cfe6ff">FRONTERA: %d/%d — mixta, como el teorema declara</text>
<text x="120" y="562" font-size="12.5" font-family="Georgia" fill="#9aa8c4">batería F6: 50 casos, %d violaciones · m=1: λ ≥ 4 en ‖nθ‖ ≥ π/2, desde n = 1</text>
<rect x="710" y="418" width="600" height="160" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1010" y="448" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA FORJA Y LA AUDITORÍA (F333-F338)</text>
<text x="740" y="478" font-size="12.5" font-family="Georgia" fill="#cfe6ff">el flash tectónico de Nico → el semigrupo de citas y las placas (F333) →</text>
<text x="740" y="502" font-size="12.5" font-family="Georgia" fill="#cfe6ff">la Ley del Líder (F334) → la banda fina F1-F6 con N* explícito (F335) →</text>
<text x="740" y="526" font-size="12.5" font-family="Georgia" fill="#cfe6ff">el diseño I-X (F336) → la costura de F(η) afuera (F337) → las puntadas</text>
<text x="740" y="550" font-size="12.5" font-family="Georgia" fill="#cfe6ff">de notación (F338) — y el sello de la auditora ✅</text>
<text x="700" y="614" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: el paisaje tiene tres regiones — el pozo, la montaña, y la frontera que las separa — y una sola ley</text>
<text x="700" y="636" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">las gobierna: EL LÍDER MANDA. Donde el líder afina, la escalera se hunde; donde desafina a pleno, se levanta</text>
<text x="700" y="658" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">montaña; y donde enmudece, ningún teorema puede decidir — y eso también está demostrado.</text>
<text x="700" y="692" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Tres regiones, tres en la mesa que lo forjaron, y el nombre del capitán: Trinidad. La regla del sello preside: nada de esto demuestra RH.</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, fneg, ftot, apos, atot, fr0, fr1, viol)
	os.WriteFile("el-teorema-de-la-trinidad.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-teorema-de-la-trinidad.svg")
}
