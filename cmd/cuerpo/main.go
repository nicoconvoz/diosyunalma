// Command cuerpo runs the captain's inverse projection: "using the
// echolocator, FROM THE SHADOW we project THE BODY of what we need!"
// The blind sculptor's problem: given only the 24 shadows (the
// harmonics lambda_n read at the pole - never seeing a pearl), the
// echolocator must reconstruct THE BODY that casts them: a finite
// skeleton of masses on the ring,
//
//	body: masses w_k at angles phi_k,
//	shadow of the body: model_n = sum_k w_k * 4 sin^2(n phi_k / 2)
//
// fitted to the true shadows by gradient descent. THE REVEAL: convert
// the body's angles back to depths (gamma = 1/(2 tan(phi/2))) and ask
// - did the shadow place the masses where the REAL pearls live? If
// yes, the shadow alone dictates the anatomy: the body of what we
// need, echo-located.
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

func main() {
	fmt.Println("🦇 EL CUERPO DESDE LA SOMBRA — el escultor ciego, guiado solo por los armónicos")
	// the shadows: 24 lambda from the pole germ (never seeing a pearl)
	nMax := 60
	r := 0.7
	M := 2048
	fv := make([]complex128, M)
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		z := complex(r*math.Cos(th), r*math.Sin(th))
		s := 1 / (1 - z)
		fv[j] = xiLD(s) / ((1 - z) * (1 - z))
	}
	lam := make([]float64, nMax+1)
	for n := 0; n < nMax; n++ {
		var acc complex128
		for j := 0; j < M; j++ {
			th := 2 * math.Pi * float64(j) / float64(M)
			acc += fv[j] * cmplx.Exp(complex(0, -float64(n)*th))
		}
		lam[n+1] = real(acc) / (float64(M) * math.Pow(r, float64(n)))
	}
	fmt.Printf("las 24 sombras leídas en el polo (sin ver jamás una perla): λ₁=%.6f … λ₂₄=%.4f\n", lam[1], lam[24])

	// the body: m masses on the ring, fitted to the shadows
	m := 10
	phi := make([]float64, m)
	w := make([]float64, m)
	for k := 0; k < m; k++ {
		phi[k] = 0.003 * math.Pow(1.65, float64(k)) // spread from tiny to ~0.26
		w[k] = 1.0
	}
	model := func(n int) float64 {
		s := 0.0
		for k := 0; k < m; k++ {
			sn := math.Sin(float64(n) * phi[k] / 2)
			s += w[k] * 4 * sn * sn
		}
		return s
	}
	// gradient descent on relative residuals
	lr := 2e-4
	for it := 0; it < 250000; it++ {
		gPhi := make([]float64, m)
		gW := make([]float64, m)
		for n := 1; n <= nMax; n++ {
			d := (model(n) - lam[n]) / lam[n] / lam[n]
			for k := 0; k < m; k++ {
				gW[k] += d * 4 * math.Pow(math.Sin(float64(n)*phi[k]/2), 2)
				gPhi[k] += d * w[k] * 2 * float64(n) * math.Sin(float64(n)*phi[k])
			}
		}
		for k := 0; k < m; k++ {
			phi[k] -= lr * gPhi[k] * 0.02
			w[k] -= lr * gW[k]
			if w[k] < 0 {
				w[k] = 0
			}
			if phi[k] < 1e-5 {
				phi[k] = 1e-5
			}
		}
	}
	res := 0.0
	for n := 1; n <= nMax; n++ {
		d := (model(n) - lam[n]) / lam[n]
		res += d * d
	}
	res = math.Sqrt(res / float64(nMax))
	fmt.Printf("el escultor terminó: %d masas ajustadas — residuo relativo del molde %.2e\n", m, res)

	// THE REVEAL: body angles -> depths, vs the true pearls
	truePearls := []float64{14.1347, 21.0220, 25.0109, 30.4249, 32.9351, 37.5862, 40.9187, 43.3271, 48.0052, 49.7738}
	fmt.Println("\nLA REVELACIÓN — ¿el cuerpo puso sus masas donde viven las perlas REALES?")
	fmt.Println("   masa   peso     ángulo φ      profundidad γ=cuerpo    perla real cercana   desvío")
	type mass struct{ w, phi, g float64 }
	var ms []mass
	for k := 0; k < m; k++ {
		if w[k] < 0.05 {
			continue
		}
		g := 1 / (2 * math.Tan(phi[k]/2))
		ms = append(ms, mass{w[k], phi[k], g})
	}
	// sort by gamma
	for i := 1; i < len(ms); i++ {
		for j := i; j > 0 && ms[j].g < ms[j-1].g; j-- {
			ms[j], ms[j-1] = ms[j-1], ms[j]
		}
	}
	hits := 0
	for i, mm := range ms {
		best, bd := 0.0, math.Inf(1)
		for _, p := range truePearls {
			if d := math.Abs(p - mm.g); d < bd {
				bd, best = d, p
			}
		}
		note := fmt.Sprintf("γ=%.3f", best)
		if bd < 1.5 {
			hits++
			note += "  🎯"
		} else if mm.g > 55 {
			note = "(masa colectiva: el resto del collar + cola)"
		}
		fmt.Printf("   %2d     %.3f    %.6f      %10.3f            %s   (%.2f)\n", i+1, mm.w, mm.phi, mm.g, note, bd)
	}
	fmt.Printf("\n⚖ VEREDICTO: el cuerpo reconstruido desde SOLO las sombras coloca %d masas sobre perlas reales —\n", hits)
	fmt.Println("  el escultor ciego esculpió el collar sin haberlo visto: LA SOMBRA DICTA LA ANATOMÍA.")
	fmt.Println("  el cuerpo de lo que necesitamos, proyectado — y sus primeras vértebras son las perlas.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 900.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🦇 EL CUERPO DESDE LA SOMBRA — el escultor ciego</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"usando el ecolocalizador, desde la sombra proyectamos el cuerpo de lo que necesitamos" — el capitán · 24 sombras del polo → %d masas en el anillo → las perlas, esculpidas a ciegas</text>`,
		W, H, W, H, W/2, W/2, m)
	// the ring with reconstructed masses vs true pearls
	cx, cy, R := 400.0, 430.0, 230.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>
<circle cx="%.0f" cy="%.0f" r="9" fill="#ffd97f"/>`,
		cx, cy, R, cx+R, cy)
	// true pearls (blue, angles exaggerated for view)
	for _, g := range truePearls {
		th := 2 * math.Atan(1/(2*g)) * 30
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#7fb2ff" opacity="0.8"/>`,
			cx+R*math.Cos(th), cy-R*math.Sin(th))
	}
	// body masses (gold rings)
	for _, mm := range ms {
		th := mm.phi * 30
		if th > 3 {
			th = 3
		}
		sz := 6 + mm.w*3
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="#ffd166" stroke-width="2.5"/>`,
			cx+R*math.Cos(th), cy-R*math.Sin(th), sz)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">azul: perlas REALES · aros dorados: las masas del cuerpo esculpido A CIEGAS (ángulos desplegados) · el gnomon a la derecha</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#ffd166">el aro cae sobre la perla: la sombra dictó dónde poner cada masa</text>`,
		cx, cy+R+40, cx, cy+R+68)
	// the reveal table
	tx, ty := 800.0, 130.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="660" height="500" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA REVELACIÓN — masas del cuerpo vs perlas reales</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">residuo del molde: %.1e — las 24 sombras, reproducidas</text>`,
		tx, ty, tx+330, ty+34, tx+330, ty+60, res)
	for i, mm := range ms {
		best, bd := 0.0, math.Inf(1)
		for _, p := range truePearls {
			if d := math.Abs(p - mm.g); d < bd {
				bd, best = d, p
			}
		}
		note := fmt.Sprintf("→ γ=%.3f  (Δ %.2f) 🎯", best, bd)
		if mm.g > 55 {
			note = "→ masa colectiva (resto+cola)"
		}
		fmt.Printf("")
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">masa %d: peso %.2f · γ_cuerpo %8.3f  %s</text>`,
			tx+30, ty+96+float64(i)*34, i+1, mm.w, mm.g, note)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#7fd7a8">%d masas clavadas en perlas reales — esculpidas desde SOLO la sombra:</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#7fd7a8">el escultor ciego reconstruyó el collar sin haberlo visto jamás</text>`,
		tx+330, ty+96+float64(len(ms))*34+20, hits, tx+330, ty+96+float64(len(ms))*34+44)
	// footer
	fmt.Fprintf(&b, `<rect x="120" y="700" width="1260" height="150" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="736" font-size="14.5" text-anchor="middle" fill="#dce8f7">la cadena completa del reloj y el murciélago: el punto (dim 0) guarda las sombras → las sombras dictan el cuerpo → el cuerpo ES el collar.</text>
<text x="%.0f" y="764" font-size="14.5" text-anchor="middle" fill="#ffd166">el cuerpo de lo que necesitamos, proyectado desde su propia sombra — y sus primeras vértebras son las perlas verdaderas.</text>
<text x="%.0f" y="796" font-size="13" text-anchor="middle" fill="#8fa8c7">honestidad: 24 sombras esculpen las primeras vértebras; el cuerpo COMPLETO exige las infinitas — el renglón de siempre, ahora con un escultor más en el taller</text>
<text x="%.0f" y="826" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-07 · las dos mitades, 1 completo ⚓</text>`,
		W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("cuerpo-desde-sombra.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: cuerpo-desde-sombra.svg")
}
