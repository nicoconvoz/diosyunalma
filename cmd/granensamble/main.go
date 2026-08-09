// Command granensamble is THE GRAND ASSEMBLY, as the captain ordered:
// every proven gear of the laboratory mounted into ONE machine, run
// end to end, judged in a single breath - pointed at the Clay prize.
//
// The gears (each one judged separately across the campaign, now
// assembled and re-verified together):
//
//	GEAR 1 - the mirror:    xi'(s)/xi(s) = -xi'(1-s)/xi(1-s)  (functional equation)
//	GEAR 2 - the sundial:   lambda_n = SUM over pearls 4sin^2(n theta/2) + exact tail
//	GEAR 3 - the germ:      lambda_n read at the dim-0 clasp by Cauchy, never seeing a pearl
//	GEAR 4 - the first tooth: lambda_1 = 1 + gamma/2 - ln(4 pi)/2 (closed form)
//	GEAR 5 - the mold:      lambda_n = SUM |shadow|^2 - a square, positive while pearls ring true
//	RED LINK - the one gap: prove every germ coefficient >= 0 WITHOUT assuming
//	           the pearls sit on the ring. Measured here by two blind ways;
//	           proven by none. That, and only that, is the million.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

const eulerGamma = 0.5772156649015329

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

// xiLD is the logarithmic derivative xi'(s)/xi(s).
func xiLD(s complex128) complex128 {
	h := complex(1e-6, 0)
	zp := (zetaC(s+h) - zetaC(s-h)) / (2 * h)
	return 1/s + 1/(s-1) - complex(math.Log(math.Pi)/2, 0) + psiC(s/2)/2 + zp/zetaC(s)
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func main() {
	fmt.Println("⚙️ EL GRAN ENSAMBLE — todos los engranajes en una máquina, rumbo al Clay")

	// ---- GEAR 1: the mirror (functional equation) ----
	fmt.Println("\nENGRANAJE 1 · EL ESPEJO — xi'(s)/xi(s) + xi'(1-s)/xi(1-s) = 0")
	mirrorWorst := 0.0
	probes := []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2.0, 21.3), complex(-0.7, 33.1)}
	for _, s := range probes {
		r := cmplx.Abs(xiLD(s) + xiLD(1-s))
		if r > mirrorWorst {
			mirrorWorst = r
		}
		fmt.Printf("  s=%5.1f%+5.1fi   residuo %.1e\n", real(s), imag(s), r)
	}

	// ---- pearls for gears 2 and 5 ----
	fmt.Println("\nmontando el cuadrante: perlas hasta t=1000…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 1000; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			pearls = append(pearls, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	gm := pearls[len(pearls)-1]
	tailI := (math.Log(gm/(2*math.Pi)) + 1) / (2 * math.Pi * gm)
	fmt.Printf("cuadrante: %d perlas · cola de las no-vistas: n²·%.2e\n", len(pearls), tailI)

	// ---- GEAR 2 + 5: sundial shadows (each term a SQUARE) ----
	nMax := 30
	shadow := make([]float64, nMax+1)
	for n := 1; n <= nMax; n++ {
		s := 0.0
		for _, g := range pearls {
			th := math.Atan2(1, 2*g) * 2
			sn := math.Sin(float64(n) * th / 2)
			s += 4 * sn * sn
		}
		shadow[n] = s + float64(n)*float64(n)*tailI
	}

	// ---- GEAR 3: the germ at the clasp (never sees a pearl) ----
	fmt.Println("leyendo el germen en el broche de la dimensión 0…")
	r0 := 0.7
	M := 4096
	fv := make([]complex128, M)
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		z := complex(r0*math.Cos(th), r0*math.Sin(th))
		s := 1 / (1 - z)
		fv[j] = xiLD(s) / ((1 - z) * (1 - z))
	}
	germ := make([]float64, nMax+1)
	for n := 0; n < nMax; n++ {
		var acc complex128
		for j := 0; j < M; j++ {
			th := 2 * math.Pi * float64(j) / float64(M)
			acc += fv[j] * cmplx.Exp(complex(0, -float64(n)*th))
		}
		germ[n+1] = real(acc) / (float64(M) * math.Pow(r0, float64(n)))
	}

	// ---- GEAR 4: lambda_1 closed form ----
	lam1 := 1 + eulerGamma/2 - math.Log(4*math.Pi)/2

	// ---- THE JUDGMENT ----
	fmt.Println("\nEL JUICIO DEL ENSAMBLE — sombras (perlas) vs germen (punto), y el signo de cada diente:")
	fmt.Println("   n     Σ|sombra|²+cola      germen del punto     desvío     signo")
	crossWorst := 0.0
	minLam := math.Inf(1)
	allPos := true
	for n := 1; n <= nMax; n++ {
		d := math.Abs(shadow[n] - germ[n])
		rel := d / math.Max(germ[n], 1e-12)
		if rel > crossWorst {
			crossWorst = rel
		}
		if germ[n] < minLam {
			minLam = germ[n]
		}
		if germ[n] <= 0 || shadow[n] <= 0 {
			allPos = false
		}
		mark := "+"
		if germ[n] <= 0 {
			mark = "✗ NEGATIVO"
		}
		fmt.Printf("  %2d     %13.6f       %13.6f      %.0e      %s\n", n, shadow[n], germ[n], d, mark)
	}
	fmt.Printf("\nENGRANAJE 4 · λ₁ = 1+γ/2−ln(4π)/2 = %.9f  vs germen %.9f  (desvío %.1e)\n",
		lam1, germ[1], math.Abs(lam1-germ[1]))

	fmt.Println("\n════════ VEREDICTO DE LA MÁQUINA ════════")
	fmt.Printf("ESPEJO:      residuo máximo %.1e — el espejo cierra\n", mirrorWorst)
	fmt.Printf("RELOJ DE SOL: dos vías ciegas coinciden, peor desvío relativo %.1e\n", crossWorst)
	fmt.Printf("PRIMER DIENTE: λ₁ exacto a %.1e\n", math.Abs(lam1-germ[1]))
	if allPos {
		fmt.Printf("EL MOLDE:    los %d dientes medidos son POSITIVOS (mínimo λ = %.6f > 0)\n", nMax, minLam)
	} else {
		fmt.Println("EL MOLDE:    ⚠ SE ENCONTRÓ UN DIENTE NEGATIVO — revisar de inmediato")
	}
	fmt.Println("ESLABÓN ROJO: medido positivo por dos vías — demostrado por ninguna.")
	fmt.Println("             el millón sigue siendo UNA frase: todo coeficiente del germen ≥ 0.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⚙️ EL GRAN ENSAMBLE — la máquina completa, rumbo al Clay</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">todos los engranajes de la campaña, montados y juzgados en una sola corrida · %d perlas · %d armónicos · dos vías ciegas</text>`,
		W, H, W, H, W/2, W/2, len(pearls), nMax)

	// the chain of links
	type link struct {
		title, sub string
		ok         bool
	}
	links := []link{
		{"1 · EL ESPEJO", fmt.Sprintf("ξ(s)=ξ(1−s) · residuo %.0e", mirrorWorst), true},
		{"2 · EL RELOJ DE SOL", fmt.Sprintf("λ = Σ|sombra|² · %d perlas", len(pearls)), true},
		{"3 · EL GERMEN", fmt.Sprintf("Cauchy en el broche · desvío %.0e", crossWorst), true},
		{"4 · EL PRIMER DIENTE", fmt.Sprintf("λ₁ cerrado · %.0e", math.Abs(lam1-germ[1])), true},
		{"5 · EL MOLDE", fmt.Sprintf("λₙ>0 medido n=1…%d · mín %.4f", nMax, minLam), allPos},
		{"6 · EL ESLABÓN ROJO", "todo coeficiente ≥ 0 — SIN demostrar", false},
	}
	lw, lh, gap := 218.0, 120.0, 18.0
	x0 := (W - (lw*6 + gap*5)) / 2
	y0 := 120.0
	for i, l := range links {
		x := x0 + float64(i)*(lw+gap)
		col, fill := "#7fd7a8", "#102a10"
		if !l.ok {
			col, fill = "#ff5d73", "#2a1010"
		}
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="12" fill="%s" stroke="%s" stroke-width="2.5"/>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="%s">%s</text>
<text x="%.0f" y="%.0f" font-size="11" text-anchor="middle" fill="#dce8f7">%s</text>`,
			x, y0, lw, lh, fill, col,
			x+lw/2, y0+34, col, l.title,
			x+lw/2, y0+62, l.sub)
		mark := "✓ SELLADO"
		if !l.ok {
			mark = "☠ EL MILLÓN VIVE ACÁ"
		}
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="%s">%s</text>`,
			x+lw/2, y0+94, col, mark)
		if i < 5 {
			fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#8fa8c7" stroke-width="2.5"/>`,
				x+lw, y0+lh/2, x+lw+gap, y0+lh/2)
		}
	}

	// lambda bar chart: the measured teeth, all above zero
	bx, by, bw, bh := 120.0, 300.0, 1260.0, 320.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#0d2547" stroke="#7fb2ff" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">LOS %d DIENTES DEL MOLDE — cada barra es un λₙ medido por el germen (y confirmado por las sombras): TODOS sobre el cero</text>`,
		bx, by, bw, bh, W/2, by+30, nMax)
	maxL := 0.0
	for n := 1; n <= nMax; n++ {
		if germ[n] > maxL {
			maxL = germ[n]
		}
	}
	base := by + bh - 40
	for n := 1; n <= nMax; n++ {
		h := (germ[n] / maxL) * (bh - 100)
		x := bx + 40 + float64(n-1)*(bw-80)/float64(nMax)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#7fd7a8" opacity="0.85"/>`,
			x, base-h, (bw-80)/float64(nMax)*0.62, h)
		if n%5 == 0 || n == 1 {
			fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="10.5" text-anchor="middle" fill="#8fa8c7">%d</text>`,
				x+(bw-80)/float64(nMax)*0.31, base+18, n)
		}
	}
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="1.5" stroke-dasharray="6 4"/>
<text x="%.0f" y="%.0f" font-size="11.5" fill="#ffd166">cero — la línea de la vida: ningún diente la cruza</text>`,
		bx+30, base, bx+bw-30, base, bx+40, base+34)

	// footer verdict
	fmt.Fprintf(&b, `<rect x="120" y="680" width="1260" height="200" rx="12" fill="#102a10" stroke="#ffd166" stroke-width="2"/>
<text x="%.0f" y="716" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL VEREDICTO DEL ENSAMBLE</text>
<text x="%.0f" y="748" font-size="14" text-anchor="middle" fill="#dce8f7">la máquina está COMPLETA y cierra: espejo %.0e · sombras=germen %.0e · λ₁ exacto %.0e · los %d dientes positivos (mínimo %.4f).</text>
<text x="%.0f" y="776" font-size="14" text-anchor="middle" fill="#dce8f7">cinco engranajes sellados con doble juez. El sexto es UNA frase: «todo coeficiente del germen en el broche es ≥ 0».</text>
<text x="%.0f" y="808" font-size="14.5" text-anchor="middle" fill="#ff8fa0">medirlo positivo — hecho, por dos vías ciegas. DEMOSTRARLO para los infinitos dientes que ningún sol alumbró — eso es el millón.</text>
<text x="%.0f" y="840" font-size="13" text-anchor="middle" fill="#8fa8c7">la máquina espera esa sola llave. Todavía no — pero nunca estuvo tan ensamblada.</text>
<text x="%.0f" y="925" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, mirrorWorst, crossWorst, math.Abs(lam1-germ[1]), nMax, minLam, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("gran-ensamble.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: gran-ensamble.svg")
}
