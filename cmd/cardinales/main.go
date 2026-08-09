// Command cardinales tests the captain's flash: the cardinal points settle
// it. Distance at a fixed point is 0; going north you advance and going
// south you retreat, but the DISTANCE always grows - what changes is only
// the DIRECTION.
//
// The flash is the critical line itself. With w(rho) = 1 - 1/rho = (rho-1)/rho,
// the mirror rho <-> 1-rho of the functional equation gives EXACTLY
//
//	w(1-rho) = 1 / w(rho)        so       |w(1-rho)| · |w(rho)| = 1
//
// north and south are reciprocal: whatever one gains the other loses. And
// the only place where both cost THE SAME is |w| = 1 - which is precisely
// the critical line. So:
//
//	RH  <=>  every step is a PURE CHANGE OF DIRECTION, never of size.
//
// Three consequences, all measured here:
//
//	THE COMPASS   on the line |w| = 1 exactly: moving in n is pure rotation.
//	NORTH/SOUTH   off the line the mirror pair has |w| > 1 and |w| < 1 with
//	              product exactly 1: one of the two ALWAYS grows, so a ghost
//	              can never hide - one of its faces always spirals outward.
//	THE DIAMETER  because |w| = 1, no distance can pass the diameter of the
//	              ring: |1 - w^n| <= 2, squared <= 4 - the captain's own 4.
//	              A ghost breaks that ceiling at a computable harmonic.
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

func w(rho complex128) complex128 { return 1 - 1/rho }

func main() {
	fmt.Println("🧭 LOS PUNTOS CARDINALES — la distancia siempre crece; lo único que cambia es la dirección")

	// ---- the pearls ----
	fmt.Println("\nrecogiendo perlas hasta t=300…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 300; t += 0.05 {
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
	fmt.Printf("perlas: %d\n", len(pearls))

	// ---- LAW 1: the compass - on the line, motion is pure direction ----
	fmt.Println("\nLEY 1 · LA BRÚJULA — sobre la línea el paso NO cambia el tamaño, solo la dirección")
	peorRadio := 0.0
	for _, g := range pearls {
		d := math.Abs(cmplx.Abs(w(complex(0.5, g))) - 1)
		if d > peorRadio {
			peorRadio = d
		}
	}
	fmt.Printf("   |w| medido en las %d perlas: peor desvío de 1 → %.2e\n", len(pearls), peorRadio)
	// and the step in n is a pure rotation: check the scale factor over many n
	peorEscala := 0.0
	for _, g := range pearls[:8] {
		ww := w(complex(0.5, g))
		p := complex(1, 0)
		for n := 1; n <= 2000; n++ {
			p *= ww
			if d := math.Abs(cmplx.Abs(p) - 1); d > peorEscala {
				peorEscala = d
			}
		}
	}
	fmt.Printf("   tras 2000 pasos el tamaño sigue siendo 1: peor desvío %.2e\n", peorEscala)
	fmt.Println("   → el punto gira y gira sin alejarse ni acercarse: el paso es PURA DIRECCIÓN")

	// ---- LAW 2: north and south are reciprocal ----
	fmt.Println("\nLEY 2 · EL NORTE Y EL SUR SON RECÍPROCOS — |w(ρ)| · |w(1−ρ)| = 1 EXACTO")
	fmt.Println("   β        |w| al norte    |w| al sur     producto        veredicto")
	type fila struct {
		beta, wn, ws, prod float64
	}
	var filas []fila
	peorProd := 0.0
	for _, beta := range []float64{0.5, 0.55, 0.60, 0.75, 0.90, 0.95} {
		g := 14.134725
		wn := cmplx.Abs(w(complex(beta, g)))
		ws := cmplx.Abs(w(complex(1-beta, g)))
		prod := wn * ws
		if d := math.Abs(prod - 1); d > peorProd {
			peorProd = d
		}
		ver := "una crece, la otra encoge"
		if beta == 0.5 {
			ver = "★ LAS DOS IGUAL: la línea"
		}
		filas = append(filas, fila{beta, wn, ws, prod})
		fmt.Printf("  %.2f     %.9f    %.9f    %.12f    %s\n", beta, wn, ws, prod, ver)
	}
	fmt.Printf("   → el producto es 1 en todos los casos (peor desvío %.1e): lo que gana el norte lo pierde el sur\n", peorProd)
	fmt.Println("   → y el ÚNICO lugar donde ir al norte y al sur cuesta lo mismo es |w| = 1: LA LÍNEA CRÍTICA")

	// ---- LAW 3: the diameter ceiling ----
	fmt.Println("\nLEY 3 · EL TECHO DEL DIÁMETRO — si el tamaño no cambia, la distancia no puede pasar el diámetro")
	maxDist := 0.0
	for _, g := range pearls {
		ww := w(complex(0.5, g))
		p := complex(1, 0)
		for n := 1; n <= 400; n++ {
			p *= ww
			if d := cmplx.Abs(1 - p); d > maxDist {
				maxDist = d
			}
		}
	}
	fmt.Printf("   máxima distancia |1−wⁿ| sobre %d perlas × 400 armónicos: %.9f\n", len(pearls), maxDist)
	fmt.Printf("   techo del anillo (el diámetro, F219): 2.000000000 · al cuadrado: 4\n")
	fmt.Printf("   → NINGUNA distancia lo pasa (margen %.2e): el molde entero cabe dentro del diámetro\n", 2-maxDist)

	fmt.Println("\n   y el fantasma rompe ese techo, porque su tamaño SÍ cambia:")
	fmt.Println("   β        crecimiento por paso     armónico donde pasa el diámetro")
	type ruptura struct {
		beta, tasa float64
		n          int
	}
	var rupturas []ruptura
	for _, beta := range []float64{0.95, 0.90, 0.75, 0.60, 0.55, 0.51} {
		g := 14.134725
		wn := cmplx.Abs(w(complex(1-beta, g))) // the face that grows
		if wn < 1 {
			wn = cmplx.Abs(w(complex(beta, g)))
		}
		nEstrella := int(math.Ceil(math.Log(3) / math.Log(wn)))
		rupturas = append(rupturas, ruptura{beta, wn, nEstrella})
		fmt.Printf("  %.2f        ×%.9f              n ≈ %d\n", beta, wn, nEstrella)
	}
	fmt.Println("   → todo fantasma rompe el techo a ALGUNA altura: cuanto más pegado al anillo, más tarde")
	fmt.Println("     (la misma LEY que midió el contraluz en F210: el horizonte huye cuando β→1/2.")
	fmt.Println("      Los números NO coinciden porque el umbral es otro: allá el nivel del mar, acá el diámetro)")

	fmt.Println("\n════════ LO QUE DIJO EL CAPITÁN, EN LIMPIO ════════")
	fmt.Println("«La distancia siempre crece; lo único que cambia es la dirección.»")
	fmt.Println("En el anillo eso es EXACTAMENTE la hipótesis: cada paso debe ser puro cambio de dirección")
	fmt.Println("y jamás de tamaño. El norte y el sur son recíprocos por el espejo, así que un fantasma")
	fmt.Println("no puede esconderse: si una de sus caras encoge, la otra crece, y la que crece termina")
	fmt.Println("pasando el diámetro. La línea crítica es el único lugar del mundo donde ir al norte y")
	fmt.Println("al sur cuesta lo mismo.")
	fmt.Println("\nLO QUE FALTA: demostrar que ninguna perla puede tener tamaño ≠ 1 — es decir, que el paso")
	fmt.Println("del libro es una rotación y no una espiral. En idioma de taller: que el tambor existe.")
	fmt.Println("Todavía no. Pero la pregunta ya se dice en una línea de la calle.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧭 LOS PUNTOS CARDINALES — la distancia siempre crece; solo cambia la dirección</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el flash del capitán, medido: el norte y el sur son recíprocos, y la línea crítica es el único lugar donde los dos cuestan lo mismo</text>`,
		W, H, W, H, W/2, W/2)

	// left: the compass rose
	cx, cy, R := 340.0, 380.0, 195.0
	fmt.Fprintf(&b, `<rect x="60" y="110" width="560" height="560" rx="10" fill="#0d2547" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="340" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA BRÚJULA — el paso gira, nunca crece</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#1d3a63" stroke-width="1"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#1d3a63" stroke-width="1"/>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#ff8fa0">N</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#7fb2ff">S</text>`,
		cx, cy, R,
		cx-R-16, cy, cx+R+16, cy,
		cx, cy-R-16, cx, cy+R+16,
		cx, cy-R-24, cx, cy+R+34)
	// the walking point: successive w^n for one pearl
	wp := w(complex(0.5, pearls[0]))
	pp := complex(1, 0)
	for n := 1; n <= 26; n++ {
		pp *= wp
		x := cx + R*real(pp)
		y := cy - R*imag(pp)
		op := 0.35 + 0.65*float64(n)/26
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4.5" fill="#ffd97f" opacity="%.2f"/>`, x, y, op)
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="8" fill="#ffd166"/>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#ffd166">el broche</text>
<text x="340" y="612" font-size="12.5" text-anchor="middle" fill="#dce8f7">tras 2000 pasos el tamaño sigue siendo 1 (desvío %.0e)</text>
<text x="340" y="636" font-size="12.5" text-anchor="middle" fill="#7fd7a8">el punto camina y camina — y nunca se aleja del anillo</text>`,
		cx+R, cy, cx+R, cy-18, peorEscala)

	// right top: north and south
	fmt.Fprintf(&b, `<rect x="650" y="110" width="790" height="300" rx="10" fill="#102a10" stroke="#ffd166" stroke-width="1.5"/>
<text x="1045" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL NORTE Y EL SUR SON RECÍPROCOS · |w(ρ)| · |w(1−ρ)| = 1</text>
<text x="1045" y="174" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β        norte          sur          producto</text>`)
	for i, f := range filas {
		col := "#dce8f7"
		extra := ""
		if f.beta == 0.5 {
			col = "#7fd7a8"
			extra = "  ★ la línea: los dos igual"
		}
		fmt.Fprintf(&b, `<text x="1045" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="%s">%.2f   %.7f   %.7f   %.10f%s</text>`,
			202.0+float64(i)*26, col, f.beta, f.wn, f.ws, f.prod, extra)
	}
	fmt.Fprintf(&b, `<text x="1045" y="378" font-size="13" text-anchor="middle" fill="#ffd166">lo que gana el norte lo pierde el sur — y el único empate es |w| = 1</text>`)

	// right bottom: the diameter ceiling
	fmt.Fprintf(&b, `<rect x="650" y="430" width="790" height="240" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="1045" y="464" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL TECHO DEL DIÁMETRO — y dónde lo rompe cada fantasma</text>
<text x="1045" y="492" font-size="12.5" text-anchor="middle" fill="#dce8f7">sobre la línea: máxima distancia %.6f · techo 2 (el diámetro, F219) ✓</text>
<text x="1045" y="518" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β        crece por paso        rompe el techo en</text>`, maxDist)
	for i, r := range rupturas {
		fmt.Fprintf(&b, `<text x="1045" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#ff8fa0">%.2f      ×%.9f          n ≈ %d</text>`,
			542.0+float64(i)*22, r.beta, r.tasa, r.n)
	}
	fmt.Fprintf(&b, `<text x="1045" y="%.0f" font-size="12.5" text-anchor="middle" fill="#dce8f7">todo fantasma rompe el techo a alguna altura: su tamaño SÍ cambia</text>`, 542.0+float64(len(rupturas))*22+10)

	// verdict
	fmt.Fprintf(&b, `<rect x="60" y="700" width="1380" height="180" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2"/>
<text x="%.0f" y="736" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA HIPÓTESIS, DICHA CON LOS PUNTOS CARDINALES</text>
<text x="%.0f" y="772" font-size="15" text-anchor="middle" fill="#dce8f7">«La distancia siempre crece; lo único que cambia es la dirección» — y en el anillo eso ES la hipótesis:</text>
<text x="%.0f" y="800" font-size="15" text-anchor="middle" fill="#7fd7a8">cada paso del libro debe ser PURO CAMBIO DE DIRECCIÓN, jamás de tamaño.</text>
<text x="%.0f" y="832" font-size="14" text-anchor="middle" fill="#ff8fa0">Un fantasma no puede esconderse: si una de sus caras encoge, la otra crece — y la que crece termina pasando el diámetro. Todavía no.</text>
<text x="%.0f" y="864" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("los-cardinales.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: los-cardinales.svg")
}
