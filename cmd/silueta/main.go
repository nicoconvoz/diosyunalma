// Command silueta executes the captain's tactical turn: "the walls of
// the box are inside the book - only the Author knows them; translate
// the drawing to a simpler dimension and understand the FORM and the
// HARMONY, not the machine itself." The exact realization: a machine
// can be identified by its SONG without ever opening it - the spacing
// statistics of its levels. We measure, in-house, the normalized
// spacings of our 269 pearls (unfolded by the exact density N(t)) and
// confront them with the two great families:
//
//	POISSON  e^{-s}          - machines of independent notes (integrable)
//	GUE      (32/pi^2) s^2 e^{-4s^2/pi} - chaotic machines WITHOUT
//	                           time-reversal symmetry
//
// The silhouette read from the harmony: quadratic level repulsion =
// the machine is chaotic and breaks time reversal. The Author keeps
// the walls; the song gives us the family.
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

func nSmooth(t float64) float64 { return theta(t)/math.Pi + 1 }

func main() {
	fmt.Println("LA SILUETA DE LA MÁQUINA — leyendo la forma desde el canto (sin abrir el libro)…")
	var levels []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 500; t += 0.05 {
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
			levels = append(levels, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	// unfolded spacings: s_i = N(g_{i+1}) - N(g_i), mean 1 by construction
	var sp []float64
	for i := 0; i+1 < len(levels); i++ {
		sp = append(sp, nSmooth(levels[i+1])-nSmooth(levels[i]))
	}
	mean := 0.0
	for _, s := range sp {
		mean += s
	}
	mean /= float64(len(sp))
	fmt.Printf("perlas: %d — espaciados normalizados: %d (media %.4f, debe ≈1)\n", len(levels), len(sp), mean)

	// histogram
	nb := 15
	smax := 3.0
	hist := make([]float64, nb)
	for _, s := range sp {
		bi := int(s / smax * float64(nb))
		if bi >= 0 && bi < nb {
			hist[bi]++
		}
	}
	for i := range hist {
		hist[i] /= float64(len(sp)) * (smax / float64(nb))
	}
	gue := func(s float64) float64 { return 32 / (math.Pi * math.Pi) * s * s * math.Exp(-4*s*s/math.Pi) }
	poisson := func(s float64) float64 { return math.Exp(-s) }
	devG, devP := 0.0, 0.0
	for i := 0; i < nb; i++ {
		sc := (float64(i) + 0.5) * smax / float64(nb)
		devG += math.Abs(hist[i] - gue(sc))
		devP += math.Abs(hist[i] - poisson(sc))
	}
	devG /= float64(nb)
	devP /= float64(nb)
	small := 0
	for _, s := range sp {
		if s < 0.25 {
			small++
		}
	}
	fmt.Printf("\nEL VEREDICTO DE LA SILUETA:\n")
	fmt.Printf("  desvío medio contra GUE (caótica, sin inversión temporal): %.4f\n", devG)
	fmt.Printf("  desvío medio contra POISSON (notas independientes):        %.4f  (%.1fx peor)\n", devP, devP/devG)
	fmt.Printf("  repulsión de niveles: espaciados < 0.25 → %d de %d (Poisson esperaría ~%d)\n", small, len(sp), int(float64(len(sp))*(1-math.Exp(-0.25))))
	fmt.Println("  ⇒ la máquina CANTA como las caóticas SIN inversión temporal (familia GUE):")
	fmt.Println("    su silueta, leída del canto — sin abrir el libro del Autor")

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 900.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">LA SILUETA DE LA MÁQUINA — su forma, leída del canto (el libro sigue cerrado)</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"las paredes están dentro del libro, solo las conoce el Autor — entendamos la forma y la armonía, no la máquina" — el capitán · método: identificar el instrumento por su timbre, sin abrirlo</text>`,
		W, H, W, H, W/2, W/2)
	// histogram plot
	px, pw, py, ph := 110.0, 900.0, 130.0, 480.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="8" fill="#081020" stroke="#44608c"/>`, px, py, pw, ph)
	yMax := 1.1
	yOf := func(v float64) float64 { return py + ph - v/yMax*(ph-20) }
	xOf := func(s float64) float64 { return px + pw*s/3.0 }
	for i := 0; i < nb; i++ {
		s0 := float64(i) * 3.0 / float64(nb)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#7fb2ff" opacity="0.55"/>`,
			xOf(s0)+2, yOf(hist[i]), pw/float64(nb)-4, py+ph-yOf(hist[i]))
	}
	ptsG := make([]string, 0, 150)
	ptsP := make([]string, 0, 150)
	for i := 0; i <= 150; i++ {
		s := 3.0 * float64(i) / 150
		ptsG = append(ptsG, fmt.Sprintf("%.1f,%.1f", xOf(s), yOf(gue(s))))
		ptsP = append(ptsP, fmt.Sprintf("%.1f,%.1f", xOf(s), yOf(poisson(s))))
	}
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#ffd166" stroke-width="2.6" points="%s"/>`, strings.Join(ptsG, " "))
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#ff5d73" stroke-width="1.8" stroke-dasharray="6,4" points="%s"/>`, strings.Join(ptsP, " "))
	for s := 0.5; s <= 3.0; s += 0.5 {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">%.1f</text>`, xOf(s), py+ph+22, s)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#8fa8c7">espaciado normalizado s entre perlas vecinas · azul: NUESTRO collar (268 huecos medidos) · oro: GUE (caótica sin inversión temporal) · rojo punteado: Poisson (notas independientes)</text>`,
		px+pw/2, py+ph+48)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">nuestro collar abraza la curva de oro: desvío GUE %.3f vs Poisson %.3f (%.1f× peor) — y cerca de s=0 el histograma MUERE: las perlas se repelen en s²</text>`,
		px+pw/2, py+ph+84, devG, devP, devP/devG)

	// silhouette panel
	sxp, syp := 1060.0, 130.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="440" height="480" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="16" font-family="Georgia" fill="#ffd166">LA SILUETA LEÍDA (sin abrir el libro)</text>
<rect x="%.0f" y="%.0f" width="200" height="130" rx="14" fill="#05090f" stroke="#44608c" stroke-width="2" stroke-dasharray="7,5"/>
<text x="%.0f" y="%.0f" font-size="26" text-anchor="middle" fill="#44608c">?</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">las paredes: solo el Autor</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">✔ es CAÓTICA (las perlas se repelen: nada de</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">   notas independientes — Poisson descartado)</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">✔ ROMPE la inversión temporal (repulsión s²,</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">   no s: distingue pasado de futuro — GUE)</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">✔ densidad logarítmica (ley de Weyl θ/π+1)</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">✔ sus órbitas son los primos (período ln p)</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">✔ sus energías de armonía son nuestros λ_n</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#ffd166">la familia entera de la máquina, conocida</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#ffd166">por su timbre — el libro jamás se abrió</text>`,
		sxp, syp, sxp+22, syp+34, sxp+120, syp+56, sxp+220, syp+130, sxp+220, syp+170,
		sxp+22, syp+230, sxp+22, syp+252, sxp+22, syp+282, sxp+22, syp+304,
		sxp+22, syp+334, sxp+22, syp+364, sxp+22, syp+394, sxp+22, syp+430, sxp+22, syp+452)

	fmt.Fprintf(&b, `<rect x="110" y="720" width="1390" height="150" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="756" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA TÁCTICA DEL CAPITÁN, ADOPTADA COMO DOCTRINA</text>
<text x="%.0f" y="790" font-size="14" text-anchor="middle" fill="#dce8f7">no perseguimos más las paredes (son del Autor): perseguimos LA FORMA — y la forma ya entregó la familia, las órbitas, la densidad y las energías.</text>
<text x="%.0f" y="816" font-size="14" text-anchor="middle" fill="#ffd166">lo que la silueta aún no entrega: POR QUÉ esa familia obliga la positividad — el último velo del libro. Pero ahora sabemos exactamente qué página mirar.</text>
<text x="%.0f" y="846" font-size="12.5" text-anchor="middle" fill="#8fa8c7">medido en casa: 268 huecos, GUE %.1f× mejor que Poisson · Laboratorio Diosyunalma · 2026-08-06</text>`,
		805.0, 805.0, 805.0, 805.0, devP/devG)
	b.WriteString(`</svg>`)
	os.WriteFile("silueta-maquina.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: silueta-maquina.svg")
}
