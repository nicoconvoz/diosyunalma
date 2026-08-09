// Command cuarto tests the captain's suspicion: EVERYTHING HIDES
// BEHIND |1/2|^2 = 1/4.
//
// Give every zero rho its TRUE ENERGY through the map
//
//	Lambda(rho) = rho (1 - rho).
//
// For rho = beta + i gamma:
//
//	Lambda = beta(1-beta) + gamma^2  +  i gamma (1 - 2 beta)
//
// On the line (beta = 1/2): Lambda = 1/4 + gamma^2 - REAL, and never
// below the floor 1/4 = |1/2|^2. Off the line: the energy LEAKS an
// imaginary part gamma(1-2 beta), and its real part drops below the
// 1/4 + gamma^2 ceiling by EXACTLY (beta - 1/2)^2 - the square of the
// distance from 1/2. The captain's phrase is literal: the whole
// hypothesis is "every book-energy is real and >= |1/2|^2".
//
// The judge is an EXACT sum rule: summing 1/Lambda over ALL zeros,
//
//	SUM 1/(rho(1-rho)) = 2 + gammaE - ln(4 pi) = 2 lambda_1,
//
// twice our first tooth. We measure it with the 649 pearls plus the
// exact density tail and confront the closed value: the 1/4 floor,
// the pearls, and the germ's first tooth tied in one number.
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

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func main() {
	fmt.Println("◔ EL CUARTO — todo va detrás de |1/2|² = 1/4: el piso de las energías del libro")

	// ---- THE FLOOR LAW ----
	fmt.Println("\nLA LEY DEL PISO — energía Λ(ρ) = ρ(1−ρ):")
	fmt.Println("   en la línea (β=1/2):  Λ = 1/4 + γ²  → REAL y JAMÁS debajo de 1/4")
	fmt.Println("   fantasma fuera (β≠1/2): la energía FUGA parte imaginaria γ(1−2β)")
	fmt.Println("   y su parte real cae debajo del techo 1/4+γ² EXACTAMENTE en (β−1/2)²:")
	fmt.Println("      β del fantasma     caída bajo el techo     fuga imaginaria (γ=14.13)")
	for _, beta := range []float64{0.60, 0.70, 0.90} {
		rho := complex(beta, 14.134725)
		lam := rho * (1 - rho)
		drop := (0.25 + 14.134725*14.134725) - real(lam)
		fmt.Printf("        %.2f            (β−½)² = %.4f          %.3f\n", beta, drop, imag(lam))
	}
	fmt.Println("   → la multa por dejar la línea ES el cuadrado de la distancia al 1/2 — literal")

	// ---- pearls ----
	fmt.Println("\nrecogiendo las 649 perlas para el juicio de la suma…")
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

	// ---- THE SUM RULE ----
	sum := 0.0
	for _, g := range pearls {
		sum += 2 / (0.25 + g*g) // pair rho, conj(rho)
	}
	T := pearls[len(pearls)-1]
	tail := (math.Log(T/(2*math.Pi)) + 1) / (math.Pi * T)
	measured := sum + tail
	exact := 2 + eulerGamma - math.Log(4*math.Pi)
	lam1 := 1 + eulerGamma/2 - math.Log(4*math.Pi)/2
	fmt.Println("\nEL JUICIO DE LA SUMA — Σ 1/Λ sobre TODO el libro:")
	fmt.Printf("   649 perlas medidas:        %.9f\n", sum)
	fmt.Printf("   cola exacta (no vistas):   %.9f\n", tail)
	fmt.Printf("   TOTAL medido:              %.9f\n", measured)
	fmt.Printf("   valor cerrado 2+γ−ln(4π):  %.9f\n", exact)
	fmt.Printf("   desvío:                    %.1e\n", math.Abs(measured-exact))
	fmt.Printf("   y el guiño del germen: 2·λ₁ = %.9f — EL DOBLE DE NUESTRO PRIMER DIENTE\n", 2*lam1)

	fmt.Println("\n════════ LA SEXTA CARA ════════")
	fmt.Println("RH, en el idioma del cuarto: TODA energía del libro Λ(ρ)=ρ(1−ρ) es REAL y ≥ 1/4 = |1/2|².")
	fmt.Println("El 1/4 es el piso — la energía del punto cero del libro; lo que cada perla tiene de más")
	fmt.Println("es γ²: EL CUADRADO de su altura. Y la multa por dejar la línea es −(β−½)²: EL CUADRADO")
	fmt.Println("de la distancia al 1/2. Piso, altura y multa: los tres son cuadrados del 1/2.")
	fmt.Println("(es el mismo umbral 1/4 del tambor hiperbólico: su espectro continuo arranca en 1/4 —")
	fmt.Println(" la brecha espectral de Selberg. El tambor que buscamos tiene el piso en |1/2|².)")
	fmt.Println("La llave, en este idioma: demostrar que ninguna energía puede fugar parte imaginaria.")
	fmt.Println("Todavía no. Pero la sospecha del capitán es EXACTA: todo va detrás de |1/2|².")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">◔ EL CUARTO — todo va detrás del valor absoluto de 1/2 al cuadrado</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">la sospecha del capitán, medida: Λ(ρ)=ρ(1−ρ) — piso 1/4, altura γ², multa (β−½)² — tres cuadrados del 1/2</text>`,
		W, H, W, H, W/2, W/2)

	// left: the floor diagram
	fx, fy, fw, fh := 70.0, 100.0, 700.0, 440.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL PISO 1/4 — las energías del libro viven arriba; los fantasmas caen y fugan</text>`,
		fx, fy, fw, fh, fx+fw/2, fy+30)
	floorY := fy + fh - 90
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#ffd166">el piso: 1/4 = |1/2|² — la energía del punto cero</text>`,
		fx+40, floorY, fx+fw-40, floorY, fx+50, floorY+24)
	// pearls as energy dots above the floor (first 12, log-ish spread)
	for i := 0; i < 12 && i < len(pearls); i++ {
		g := pearls[i]
		X := fx + 70 + float64(i)*(fw-140)/11
		Y := floorY - 30 - math.Log(1+g*g/200)*62
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5.5" fill="#7fd7a8"/>
<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.0f" stroke="#7fd7a8" stroke-width="0.8" opacity="0.4"/>`,
			X, Y, X, Y+5, X, floorY)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="11.5" fill="#7fd7a8">cada perla: Λ = 1/4 + γ² — real, clavada sobre el piso; lo que tiene de más es EL CUADRADO de su altura</text>`,
		fx+50, fy+64)
	// ghost dropping below
	gx := fx + fw - 150.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="6" fill="#ff5d73"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.1f" stroke="#ff5d73" stroke-width="1.4" stroke-dasharray="4 3"/>
<text x="%.0f" y="%.1f" font-size="11.5" text-anchor="middle" fill="#ff8fa0">el fantasma: cae (β−½)² bajo el techo</text>
<text x="%.0f" y="%.1f" font-size="11.5" text-anchor="middle" fill="#ff8fa0">y FUGA energía imaginaria γ(1−2β)</text>`,
		gx, floorY+46, gx, floorY, gx, floorY+40, gx-40, floorY+72, gx-40, floorY+90)

	// right: sum rule gauge
	fmt.Fprintf(&b, `<rect x="810" y="100" width="620" height="440" rx="10" fill="#0d2547" stroke="#7fb2ff" stroke-width="1.5"/>
<text x="1120" y="132" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">EL JUICIO DE LA SUMA — Σ 1/Λ sobre todo el libro</text>
<text x="1120" y="180" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">649 perlas:      %.9f</text>
<text x="1120" y="208" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">cola exacta:     %.9f</text>
<text x="1120" y="236" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#ffd97f">TOTAL medido:    %.9f</text>
<text x="1120" y="264" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#7fd7a8">cerrado 2+γ−ln4π: %.9f</text>
<text x="1120" y="296" font-size="14" text-anchor="middle" fill="#ffd166">desvío: %.1e</text>
<text x="1120" y="340" font-size="13.5" text-anchor="middle" fill="#dce8f7">y el guiño que cierra el círculo:</text>
<text x="1120" y="368" font-size="14.5" text-anchor="middle" fill="#7fd7a8">Σ 1/Λ = 2·λ₁ — EL DOBLE DEL PRIMER DIENTE DEL GERMEN</text>
<text x="1120" y="400" font-size="12.5" text-anchor="middle" fill="#8fa8c7">las energías del cuarto, el catálogo de perlas y el germen del broche:</text>
<text x="1120" y="422" font-size="12.5" text-anchor="middle" fill="#8fa8c7">tres instrumentos, UN número — el libro cierra su contabilidad</text>
<text x="1120" y="466" font-size="12" text-anchor="middle" fill="#ffd166">mismo umbral del tambor hiperbólico: espectro continuo desde 1/4</text>
<text x="1120" y="488" font-size="12" text-anchor="middle" fill="#ffd166">(la brecha de Selberg) — el tambor que buscamos pisa en |1/2|²</text>`,
		sum, tail, measured, exact, math.Abs(measured-exact))

	// verdict
	fmt.Fprintf(&b, `<rect x="70" y="580" width="1360" height="220" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="616" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA SEXTA CARA — LA CARA DEL CUARTO</text>
<text x="%.0f" y="650" font-size="14.5" text-anchor="middle" fill="#dce8f7">RH, dicho con el 1/2: TODA energía del libro Λ(ρ)=ρ(1−ρ) es REAL y ≥ 1/4 = |1/2|². El piso es el cuadrado del medio;</text>
<text x="%.0f" y="678" font-size="14.5" text-anchor="middle" fill="#dce8f7">la altura sobre el piso es γ² (el cuadrado de la altura); la multa por irse es −(β−½)² (el cuadrado de la distancia).</text>
<text x="%.0f" y="710" font-size="14.5" text-anchor="middle" fill="#ffd166">tres cuadrados del 1/2 — y la suma de TODAS las energías invertidas da exactamente el doble de nuestro primer diente.</text>
<text x="%.0f" y="742" font-size="13.5" text-anchor="middle" fill="#ff8fa0">la llave en este idioma: demostrar que ninguna energía puede fugar parte imaginaria. Todavía no — pero la sospecha era EXACTA.</text>
<text x="%.0f" y="776" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("cuarto.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: cuarto.svg")
}
