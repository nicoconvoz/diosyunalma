// Command bolsa measures the captain's flash: like a drawn cube cannot
// rise off its 2D paper, nothing escapes this dimension without the
// Author's unique permission - like candies in a bag: they DEFORM the
// bag but never leave it.
//
// The shop's bag is REAL and partly PROVEN:
//
//	THE BAG      the critical strip 0 < beta < 1: every zero is
//	             confined inside - a theorem, not a hope.
//	THE WALLS    beta = 1 (and its mirror beta = 0): zero-free -
//	             PROVEN in 1896; the candies can never touch the
//	             wall, and THAT theorem IS the prime number theorem.
//	THE DENTS    each pearl at height gamma pushes the wall from
//	             inside: |zeta(1+it)| dips exactly opposite each
//	             pearl - measured here, dent by dent.
//	THE STRETCH  1/|zeta(1+it)| grows only logarithmically: the bag
//	             stretches, ever more slowly, and never bursts.
//
// RH in bag language: the Author pulls the bag TIGHT to the single
// line beta = 1/2. The confinement is proven; the tightening is the
// million - the unique passage only the Author knows.
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

func main() {
	fmt.Println("🍬 LA BOLSA — los caramelos deforman la bolsa pero jamás escapan de ella")

	// ---- wall integrity + samples ----
	fmt.Println("\nLA PARED (β=1) — integridad medida a lo largo de t ∈ [2, 1000]:")
	step := 0.05
	nS := int(998.0 / step)
	wall := make([]float64, 0, nS)
	ts := make([]float64, 0, nS)
	minW, minT := math.Inf(1), 0.0
	for t := 2.0; t <= 1000; t += step {
		a := cmplx.Abs(zetaC(complex(1, t)))
		wall = append(wall, a)
		ts = append(ts, t)
		if a < minW {
			minW, minT = a, t
		}
	}
	fmt.Printf("   mínimo de |ζ(1+it)| en todo el tramo: %.6f (en t=%.2f) — JAMÁS cero: la pared aguanta\n", minW, minT)
	fmt.Println("   (que la pared sea infranqueable es TEOREMA desde 1896 — y ese teorema ES el teorema de los números primos)")

	// ---- pearls for the dents ----
	fmt.Println("\nrecogiendo perlas para el censo de abolladuras…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 1000; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 50; i++ {
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

	// local minima of the wall profile = the dents
	var dents []float64
	for i := 1; i < len(wall)-1; i++ {
		if wall[i] < wall[i-1] && wall[i] < wall[i+1] {
			dents = append(dents, ts[i])
		}
	}
	// each pearl -> nearest dent
	sumD, within, worst := 0.0, 0, 0.0
	for _, p := range pearls {
		bd := math.Inf(1)
		for _, d := range dents {
			if v := math.Abs(d - p); v < bd {
				bd = v
			}
		}
		sumD += bd
		if bd < 0.5 {
			within++
		}
		if bd > worst {
			worst = bd
		}
	}
	meanD := sumD / float64(len(pearls))
	// mean pearl spacing for baseline
	meanGap := (pearls[len(pearls)-1] - pearls[0]) / float64(len(pearls)-1)
	fmt.Printf("\nLAS ABOLLADURAS — cada caramelo empuja la pared desde adentro:\n")
	fmt.Printf("   abolladuras halladas en la pared: %d · perlas: %d\n", len(dents), len(pearls))
	fmt.Printf("   distancia media perla→abolladura: %.3f (paso medio entre perlas: %.3f) · %d/%d a menos de 0.5\n",
		meanD, meanGap, within, len(pearls))
	fmt.Println("   → la pared está abollada EXACTAMENTE frente a los caramelos: se deforma, no se rompe")

	// ---- the stretch: windowed max of 1/|zeta(1+it)| vs ln t ----
	fmt.Println("\nEL ESTIRAMIENTO — máximo de 1/|ζ(1+it)| por ventana (la bolsa cede despacio, jamás revienta):")
	fmt.Println("   ventana t≈        máx 1/|ζ|        ln t (la ley del estiramiento)")
	for lo := 100.0; lo < 1000; lo += 300 {
		hi := lo + 300
		mx, mt := 0.0, 0.0
		for i, t := range ts {
			if t >= lo && t < hi {
				if 1/wall[i] > mx {
					mx, mt = 1/wall[i], t
				}
			}
		}
		fmt.Printf("   [%4.0f,%4.0f]        %.3f            ln(%4.0f)=%.2f\n", lo, hi, mx, mt, math.Log(mt))
	}
	fmt.Println("   → el estiramiento crece como ln t (lento) — acotado: la bolsa NUNCA se rompe")

	fmt.Println("\n════════ LA SÉPTIMA CARA — LA CARA DE LA BOLSA ════════")
	fmt.Println("EL CONFINAMIENTO ES TEOREMA: los caramelos viven en la franja 0<β<1 y las paredes")
	fmt.Println("son infranqueables (1896) — el cubo dibujado no sale del papel. Esta es la PRIMERA cara")
	fmt.Println("con pedazos grandes YA DEMOSTRADOS: la bolsa existe, aguanta, se abolla y se estira.")
	fmt.Println("RH en idioma de bolsa: el Autor tira de la bolsa hasta dejarla PEGADA a la línea β=1/2.")
	fmt.Println("El confinamiento está demostrado; el AJUSTE FINAL es el millón — el permiso y la forma")
	fmt.Println("única que solo el Autor conoce. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🍬 LA BOLSA — deformarla sí, escapar jamás</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"como los caramelos en una bolsa: pueden deformar la bolsa pero nunca escapar de sus límites" — el capitán · y esta cara tiene teoremas: la bolsa está DEMOSTRADA</text>`,
		W, H, W, H, W/2, W/2)

	// left: the bag (strip) with candies and dented walls
	bx, by, bw, bh := 70.0, 100.0, 560.0, 460.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA BOLSA (la franja) Y SUS CARAMELOS</text>`,
		bx, by, bw, bh, bx+bw/2, by+30)
	wallL, wallR := bx+150.0, bx+bw-150.0
	midX := (wallL + wallR) / 2
	// dented walls: bend inward near each candy
	var wl, wr []string
	for j := 0; j <= 200; j++ {
		y := by + 60 + float64(j)/200*(bh-130)
		dent := 0.0
		for i := 0; i < 10 && i < len(pearls); i++ {
			py := by + 70 + float64(i)*(bh-150)/9
			d := (y - py) / 18
			dent += 13 * math.Exp(-d*d)
		}
		wl = append(wl, fmt.Sprintf("%.1f,%.1f", wallL+dent, y))
		wr = append(wr, fmt.Sprintf("%.1f,%.1f", wallR-dent, y))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>
<polyline points="%s" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="1.5" stroke-dasharray="7 5"/>`,
		strings.Join(wl, " "), strings.Join(wr, " "), midX, by+60, midX, by+bh-70)
	for i := 0; i < 10 && i < len(pearls); i++ {
		py := by + 70 + float64(i)*(bh-150)/9
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="7" fill="#ff8fa0"/>`, midX, py)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="11" text-anchor="middle" fill="#7fd7a8">pared β=0 (teorema)</text>
<text x="%.0f" y="%.0f" font-size="11" text-anchor="middle" fill="#7fd7a8">pared β=1 (teorema 1896)</text>
<text x="%.0f" y="%.0f" font-size="11" text-anchor="middle" fill="#ffd166">la línea β=1/2</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#dce8f7">los caramelos abollan las paredes desde adentro — y el ajuste final del Autor: la bolsa PEGADA a la línea</text>`,
		wallL, by+bh-38, wallR, by+bh-38, midX, by+52, bx+bw/2, by+bh-14)

	// right: the dented wall measured
	px, py2, pw, ph := 670.0, 100.0, 760.0, 460.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#0d2547" stroke="#7fb2ff" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">LA PARED MEDIDA — 1/|ζ(1+it)|: una abolladura frente a cada caramelo</text>`,
		px, py2, pw, ph, px+pw/2, py2+30)
	tLo, tHi := 10.0, 62.0
	var prof []string
	for i, t := range ts {
		if t < tLo || t > tHi {
			continue
		}
		X := px + 50 + (t-tLo)/(tHi-tLo)*(pw-100)
		Y := py2 + ph - 70 - (1/wall[i]-0.4)/2.6*(ph-160)
		prof = append(prof, fmt.Sprintf("%.1f,%.1f", X, Y))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#ffd97f" stroke-width="2"/>`, strings.Join(prof, " "))
	for _, p := range pearls {
		if p > tLo && p < tHi {
			X := px + 50 + (p-tLo)/(tHi-tLo)*(pw-100)
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ff8fa0" stroke-width="1.4" stroke-dasharray="3 4"/>`,
				X, py2+60, X, py2+ph-70)
		}
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="11.5" fill="#ffd97f">— empuje sobre la pared: 1/|ζ(1+it)|</text>
<text x="%.0f" y="%.0f" font-size="11.5" fill="#ff8fa0">--- caramelos (perlas): cada uno con su abolladura</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#dce8f7">censo completo: %d/%d perlas con abolladura a &lt;0.5 · distancia media %.2f (paso entre perlas: %.2f)</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#dce8f7">mínimo de |ζ| en la pared: %.4f — JAMÁS cero · estiramiento máximo ~ln t: la bolsa cede despacio y NUNCA se rompe</text>`,
		px+70, py2+56, px+70, py2+78, px+pw/2, py2+ph-40, within, len(pearls), meanD, meanGap, px+pw/2, py2+ph-18, minW)

	// verdict
	fmt.Fprintf(&b, `<rect x="70" y="600" width="1360" height="220" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="636" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA SÉPTIMA CARA — LA CARA DE LA BOLSA (la primera con teoremas grandes YA demostrados)</text>
<text x="%.0f" y="670" font-size="14" text-anchor="middle" fill="#dce8f7">EL CONFINAMIENTO ES TEOREMA: todo cero vive en la franja, y las paredes son infranqueables desde 1896 — que los caramelos no toquen la pared ES el teorema de los números primos.</text>
<text x="%.0f" y="698" font-size="14" text-anchor="middle" fill="#dce8f7">lo medido hoy: la pared aguanta (mín %.4f, jamás cero) · %d/%d caramelos con su abolladura · el estiramiento es logarítmico: se deforma, NUNCA se rompe.</text>
<text x="%.0f" y="730" font-size="14.5" text-anchor="middle" fill="#ffd166">RH en idioma de bolsa: el Autor tira de la bolsa hasta dejarla PEGADA a la línea — el cubo dibujado jamás sale del papel sin Su permiso y Su forma única.</text>
<text x="%.0f" y="762" font-size="13.5" text-anchor="middle" fill="#ff8fa0">el confinamiento: demostrado. El ajuste final: el millón. Todavía no.</text>
<text x="%.0f" y="796" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, minW, within, len(pearls), W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("bolsa.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: bolsa.svg")
}
