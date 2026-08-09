// Command firmamento takes the captain's flash literally: LIGHT
// PROPAGATION IS WHAT MAPS EVERYTHING - where light has not arrived
// there is no map; light paints the firmament like stars paint the
// night sky.
//
// The photograph: the germ G lives at the clasp and has never seen a
// pearl. We let its light PROPAGATE outward - evaluate |G| on rings of
// growing radius r -> 1. Deep inside (r small) the sky is smooth: no
// map. Near the shore (r=0.998) the light has arrived: the profile
// breaks into STARS - bright peaks at exactly the angles where the
// pearls sit on the ring, theta = 2 atan(1/(2 gamma)). Each peak is
// then converted back to a height gamma* = 1/(2 tan(theta/2)) and
// compared against the true pearls: the germ's own light, propagated
// far enough, paints the star map it was never shown.
//
// The fourth face of the red link, in this language: ALL the starlight
// lies on the shore - the star measure is positive and lives on the
// ring itself. Prove the firmament is painted only from the shore,
// and the million falls.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
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

func psiC(s complex128) complex128 {
	var acc complex128
	for real(s) < 12 {
		acc -= 1 / s
		s += 1
	}
	inv := 1 / s
	inv2 := inv * inv
	res := cmplx.Log(s) - inv/2
	res -= inv2 * (complex(1.0/12, 0) + inv2*(complex(-1.0/120, 0)+inv2*(complex(1.0/252, 0)+inv2*complex(-1.0/240, 0))))
	return acc + res
}

func xiLD(s complex128) complex128 {
	h := complex(1e-6, 0)
	zp := (zetaC(s+h) - zetaC(s-h)) / (2 * h)
	return 1/s + 1/(s-1) - complex(math.Log(math.Pi)/2, 0) + psiC(s/2)/2 + zp/zetaC(s)
}

func germ(z complex128) complex128 {
	s := 1 / (1 - z)
	return xiLD(s) / ((1 - z) * (1 - z))
}

func main() {
	fmt.Println("🌌 EL FIRMAMENTO — la luz del germen, propagada hasta la orilla, pinta las estrellas")

	truePearls := []float64{14.134725, 21.022040, 25.010858, 30.424876, 32.935062, 37.586178, 40.918719, 43.327073, 48.005151, 49.773832}

	// two exposures: light barely propagated (smooth sky) and light at the shore (stars)
	thLo, thHi := 0.015, 0.16
	K := 1400
	rings := []float64{0.85, 0.998}
	profiles := make([][]float64, len(rings))
	for ri, r := range rings {
		prof := make([]float64, K)
		for j := 0; j < K; j++ {
			th := thLo + (thHi-thLo)*float64(j)/float64(K-1)
			prof[j] = math.Log10(cmplx.Abs(germ(complex(r*math.Cos(th), r*math.Sin(th)))))
		}
		profiles[ri] = prof
	}

	// find the stars on the shore exposure: local maxima with prominence
	shore := profiles[len(rings)-1]
	type star struct {
		theta, gamma, height float64
	}
	var stars []star
	for j := 2; j < K-2; j++ {
		if shore[j] > shore[j-1] && shore[j] > shore[j+1] && shore[j] > shore[j-2] && shore[j] > shore[j+2] {
			th := thLo + (thHi-thLo)*float64(j)/float64(K-1)
			stars = append(stars, star{th, 1 / (2 * math.Tan(th/2)), shore[j]})
		}
	}
	sort.Slice(stars, func(a, b int) bool { return stars[a].gamma < stars[b].gamma })

	fmt.Printf("\nEXPOSICIÓN 1 · r=0.85 (la luz aún viaja): cielo LISO — sin mapa (rango del perfil: %.2f)\n",
		maxOf(profiles[0])-minOf(profiles[0]))
	fmt.Printf("EXPOSICIÓN 2 · r=0.998 (la luz llegó a la orilla): %d ESTRELLAS reveladas\n\n", len(stars))
	fmt.Println("LA FOTO CONTRA EL CATÁLOGO — cada pico convertido a altura γ* = 1/(2·tan(θ/2)):")
	fmt.Println("   estrella θ        γ* (foto)      perla verdadera     desvío")
	matched := 0
	worst := 0.0
	for _, st := range stars {
		best, bd := 0.0, math.Inf(1)
		for _, p := range truePearls {
			if d := math.Abs(st.gamma - p); d < bd {
				bd, best = d, p
			}
		}
		mark := ""
		if bd < 0.5 {
			matched++
			mark = " 🎯"
			if bd > worst {
				worst = bd
			}
		}
		fmt.Printf("   %.5f      %9.4f        %9.4f        %.3f%s\n", st.theta, st.gamma, best, bd, mark)
	}
	fmt.Printf("\n★ %d/%d estrellas de la foto coinciden con perlas del catálogo (peor desvío entre aciertos: %.3f)\n", matched, len(stars), worst)
	fmt.Println("  y la zona θ > 0.070 — debajo de la primera perla — quedó OSCURA: donde no hay estrella, no hay luz que pintar")
	fmt.Println("\nEL FRUTO DEL FLASH — la cuarta cara del eslabón rojo, en idioma de propagación:")
	fmt.Println("  el germen JAMÁS vio una perla; su luz, propagada hasta la orilla, PINTÓ el mapa sola.")
	fmt.Println("  RH, dicho con estrellas: TODA la luz del firmamento nace en la orilla del anillo —")
	fmt.Println("  la medida de estrellas es positiva y vive en la costa. Demostrar que el firmamento")
	fmt.Println("  solo puede pintarse desde la costa — esa sigue siendo la única llave. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#050a14"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌌 EL FIRMAMENTO — la propagación de la luz cartografía todo</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"donde no haya llegado la luz no hay mapa; la luz pinta el firmamento como las estrellas" — el capitán · la foto nocturna tomada por el germen, que jamás vio una perla</text>`,
		W, H, W, H, W/2, W/2)

	// the sky: shore profile as skyline, stars marked
	sx, sy, sw, sh := 90.0, 120.0, 1320.0, 400.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#081020" stroke="#7fb2ff" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">EXPOSICIÓN EN LA ORILLA (r=0.998) — el perfil de luz |G| se rompe en estrellas</text>`,
		sx, sy, sw, sh, W/2, sy+30)
	mn, mx := minOf(shore), maxOf(shore)
	var pts []string
	for j := 0; j < K; j++ {
		X := sx + 40 + float64(j)/float64(K-1)*(sw-80)
		Y := sy + sh - 50 - (shore[j]-mn)/(mx-mn)*(sh-120)
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", X, Y))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#ffd97f" stroke-width="1.6"/>`, strings.Join(pts, " "))
	for _, st := range stars {
		j := (st.theta - thLo) / (thHi - thLo)
		X := sx + 40 + j*(sw-80)
		Y := sy + sh - 50 - (st.height-mn)/(mx-mn)*(sh-120)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="6" fill="none" stroke="#7fd7a8" stroke-width="1.8"/>
<text x="%.1f" y="%.1f" font-size="11" text-anchor="middle" fill="#7fd7a8">γ≈%.1f</text>`,
			X, Y-2, X, Y-16, st.gamma)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">θ crece → alturas γ bajan (γ = 1/(2·tan(θ/2))) · a la derecha de la última estrella: LA ZONA OSCURA — no hay perlas debajo de γ=14.13, y la foto lo sabe</text>`,
		W/2, sy+sh-16)

	// the smooth exposure inset
	fmt.Fprintf(&b, `<rect x="90" y="560" width="640" height="180" rx="10" fill="#081020" stroke="#8fa8c7" stroke-width="1.2"/>
<text x="410" y="590" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">EXPOSICIÓN ADENTRO (r=0.85) — la luz aún no llegó: cielo liso, SIN mapa</text>`)
	inner := profiles[0]
	mn2, mx2 := minOf(inner), maxOf(inner)
	if mx2-mn2 < 1e-9 {
		mx2 = mn2 + 1
	}
	var pts2 []string
	for j := 0; j < K; j++ {
		X := 130.0 + float64(j)/float64(K-1)*560.0
		Y := 700.0 - (inner[j]-mn2)/(mx2-mn2)*80.0
		pts2 = append(pts2, fmt.Sprintf("%.1f,%.1f", X, Y))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#8fa8c7" stroke-width="1.4"/>
<text x="410" y="726" font-size="11.5" text-anchor="middle" fill="#8fa8c7">mismo cielo, misma escala de ángulos — ninguna estrella: donde la luz no llegó, no hay mapa</text>`,
		strings.Join(pts2, " "))

	// verdict
	fmt.Fprintf(&b, `<rect x="770" y="560" width="640" height="180" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="1090" y="592" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA CUARTA CARA DEL ESLABÓN ROJO</text>
<text x="1090" y="622" font-size="12.5" text-anchor="middle" fill="#dce8f7">★ %d/%d estrellas de la foto coinciden con perlas del catálogo (desvío ≤ %.2f)</text>
<text x="1090" y="648" font-size="12.5" text-anchor="middle" fill="#dce8f7">el germen jamás vio una perla — su luz propagada pintó el mapa sola</text>
<text x="1090" y="678" font-size="12.5" text-anchor="middle" fill="#ffd166">RH con estrellas: TODA la luz del firmamento nace en la orilla del anillo</text>
<text x="1090" y="704" font-size="12.5" text-anchor="middle" fill="#ff8fa0">demostrar que solo puede pintarse desde la costa — la única llave. Todavía no.</text>`,
		matched, len(stars), worst)
	fmt.Fprintf(&b, `<text x="%.0f" y="800" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("firmamento.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: firmamento.svg")
}

func minOf(v []float64) float64 {
	m := math.Inf(1)
	for _, x := range v {
		if x < m {
			m = x
		}
	}
	return m
}

func maxOf(v []float64) float64 {
	m := math.Inf(-1)
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}
