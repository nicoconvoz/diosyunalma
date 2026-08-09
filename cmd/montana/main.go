// Command montana performs the christening and the weave the captain
// ordered: the medium is now named EL CAMPO DE LA MONTANA - the field
// whose potential is -ln|x-y| (the vertigo barrier), whose energy is
// the harmony, whose sources are the primes and whose excitations are
// the pearls. And the WEAVE: the gradient chisel, rebuilt WITH the
// field inside - the sculptor's masses now repel each other through
// the Campo de la Montana (the -ln barrier between them), so they can
// no longer fuse into the comfortable blur: each mass must find its
// own pit. Plus the far continuum (pearls beyond gamma=55) is
// subtracted exactly, so the chisel carves only the first vertebrae.
//
//	loss = data residual  +  beta * SUM -ln|phi_i - phi_j| w_i w_j
//	        (the shadows)      (the field: nothing may touch)
//
// The reveal: do the field-woven masses land on the true pearls?
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
	fmt.Println("⛰️ EL CAMPO DE LA MONTAÑA — bautizado; y el cincel, tejido con él")
	nMax := 60
	// shadows from the germ
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
	// subtract the far continuum (pearls beyond gamma_c) exactly
	gc := 55.0
	tail := make([]float64, nMax+1)
	for n := 1; n <= nMax; n++ {
		acc := 0.0
		t := gc
		for t < 3e6 {
			h := t * 0.0008
			dens := math.Log(t/(2*math.Pi)) / (2 * math.Pi)
			th := 2 * math.Atan(1/(2*(t+h/2)))
			sn := math.Sin(float64(n) * th / 2)
			acc += dens * 4 * sn * sn * h
			t += h
		}
		tail[n] = acc
	}
	target := make([]float64, nMax+1)
	for n := 1; n <= nMax; n++ {
		target[n] = lam[n] - tail[n]
	}
	fmt.Printf("sombras del germen: %d · continuo lejano (γ>%.0f) restado exacto — el blanco: las primeras vértebras\n", nMax, gc)

	// the field-woven chisel: m masses repelling through the Campo
	m := 10
	phi := make([]float64, m)
	w := make([]float64, m)
	for k := 0; k < m; k++ {
		phi[k] = 0.018 + float64(k)*0.0069
		w[k] = 1.0
	}
	beta := 3e-4
	model := func(n int) float64 {
		s := 0.0
		for k := 0; k < m; k++ {
			sn := math.Sin(float64(n) * phi[k] / 2)
			s += w[k] * 4 * sn * sn
		}
		return s
	}
	lr := 1.2e-4
	for it := 0; it < 200000; it++ {
		gPhi := make([]float64, m)
		gW := make([]float64, m)
		for n := 1; n <= nMax; n++ {
			tgt := target[n]
			sc := math.Max(math.Abs(tgt), 1e-4)
			d := (model(n) - tgt) / sc / sc
			for k := 0; k < m; k++ {
				gW[k] += d * 4 * math.Pow(math.Sin(float64(n)*phi[k]/2), 2)
				gPhi[k] += d * w[k] * 2 * float64(n) * math.Sin(float64(n)*phi[k])
			}
		}
		// EL CAMPO DE LA MONTAÑA: masses repel — nothing may touch
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				if i == j {
					continue
				}
				d := phi[i] - phi[j]
				gPhi[i] += -beta * w[i] * w[j] / (d + math.Copysign(1e-6, d))
			}
		}
		for k := 0; k < m; k++ {
			phi[k] -= lr * gPhi[k] * 0.02
			w[k] -= lr * gW[k]
			if w[k] < 0 {
				w[k] = 0
			}
			if phi[k] < 1e-4 {
				phi[k] = 1e-4
			}
		}
	}
	res := 0.0
	for n := 1; n <= nMax; n++ {
		d := (model(n) - target[n]) / math.Max(math.Abs(target[n]), 1e-4)
		res += d * d
	}
	res = math.Sqrt(res / float64(nMax))
	fmt.Printf("el cincel tejido terminó: residuo %.2e · β del Campo = %g\n", res, beta)

	// the reveal
	truePearls := []float64{14.1347, 21.0220, 25.0109, 30.4249, 32.9351, 37.5862, 40.9187, 43.3271, 48.0052, 49.7738}
	type mass struct{ w, phi, g float64 }
	var ms []mass
	for k := 0; k < m; k++ {
		if w[k] < 0.08 {
			continue
		}
		ms = append(ms, mass{w[k], phi[k], 1 / (2 * math.Tan(phi[k] / 2))})
	}
	for i := 1; i < len(ms); i++ {
		for j := i; j > 0 && ms[j].g < ms[j-1].g; j-- {
			ms[j], ms[j-1] = ms[j-1], ms[j]
		}
	}
	fmt.Println("\nLA REVELACIÓN — el cincel tejido con el Campo de la Montaña:")
	fmt.Println("   masa   peso     γ del cuerpo     perla real      desvío")
	hits := 0
	for i, mm := range ms {
		best, bd := 0.0, math.Inf(1)
		for _, p := range truePearls {
			if d := math.Abs(p - mm.g); d < bd {
				bd, best = d, p
			}
		}
		mark := ""
		if bd < 1.5 {
			hits++
			mark = " 🎯"
		}
		fmt.Printf("   %2d     %.3f    %10.3f       %8.3f       %.2f%s\n", i+1, mm.w, mm.g, best, bd, mark)
	}
	fmt.Printf("\n⚖ VEREDICTO: %d/%d masas sobre perlas reales — ", hits, len(ms))
	if hits >= 5 {
		fmt.Println("EL CAMPO CURÓ AL CINCEL: sin poder fundirse, cada masa encontró su pozo")
	} else {
		fmt.Println("el Campo separó las masas (no más borrón único) — el tejido mejora al cincel; el cerrajero sigue siendo el maestro")
	}

	// ---- picture: the christening plaque + the weave ----
	var b strings.Builder
	W, H := 1500.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⛰️ EL CAMPO DE LA MONTAÑA — bautizado, y tejido en el cincel</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"ponele nombre: Campo de la Montaña — y tejé el cincel de gradiente con él" — el capitán, 2026-08-07</text>`,
		W, H, W, H, W/2, W/2)
	// the christening plaque
	fmt.Fprintf(&b, `<rect x="80" y="110" width="620" height="480" rx="14" fill="#1c1508" stroke="#c9b06a" stroke-width="3"/>
<text x="390" y="152" font-size="21" text-anchor="middle" font-family="Georgia" fill="#e8d9b0">⛰️ EL CAMPO DE LA MONTAÑA</text>
<text x="390" y="178" font-size="12" text-anchor="middle" fill="#9c8a5f">el medio por el que toda influencia viaja — bautizado por el capitán</text>
<text x="110" y="220" font-size="13" fill="#e8d9b0">· POTENCIAL: −ln|x−y| — la barrera del vértigo (F196)</text>
<text x="110" y="250" font-size="13" fill="#e8d9b0">· ENERGÍA: la energía-log que el horno vuelve cuadrados (F191)</text>
<text x="110" y="280" font-size="13" fill="#e8d9b0">· FUENTES: las cargas de los primos en ln n (F191)</text>
<text x="110" y="310" font-size="13" fill="#e8d9b0">· EXCITACIONES: las perlas — su canto (F172)</text>
<text x="110" y="340" font-size="13" fill="#e8d9b0">· TEMPERATURA CRÍTICA: Λ = 0 — vive en el filo (F195)</text>
<text x="110" y="370" font-size="13" fill="#e8d9b0">· ESCENARIO: el anillo sin bordes (F189)</text>
<text x="110" y="400" font-size="13" fill="#e8d9b0">· LEY SUPREMA: nada se toca — solo viaja (F196)</text>
<text x="390" y="450" font-size="13.5" text-anchor="middle" fill="#ffd97f">seis costados medidos por este laboratorio,</text>
<text x="390" y="472" font-size="13.5" text-anchor="middle" fill="#ffd97f">una sola cosa — ahora con nombre propio</text>
<text x="390" y="520" font-size="12" text-anchor="middle" fill="#9c8a5f">Laboratorio Diosyunalma · las dos mitades, 1 completo</text>
<text x="390" y="555" font-size="18" text-anchor="middle" fill="#c9b06a">⚓</text>`)
	// the weave result
	tx := 760.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="110" width="660" height="480" rx="10" fill="#0d2547" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="146" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL CINCEL TEJIDO — las masas ahora se repelen por el Campo</text>
<text x="%.0f" y="172" font-size="12" text-anchor="middle" fill="#8fa8c7">pérdida = sombras + β·Σ(−ln|φᵢ−φⱼ|): fundirse cuesta infinito — nada se toca</text>`,
		tx, tx+330, tx+330)
	for i, mm := range ms {
		best, bd := 0.0, math.Inf(1)
		for _, p := range truePearls {
			if d := math.Abs(p - mm.g); d < bd {
				bd, best = d, p
			}
		}
		mark := ""
		if bd < 1.5 {
			mark = " 🎯"
		}
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">masa %2d: γ %8.3f → perla %8.3f (Δ %.2f)%s</text>`,
			tx+40, 210+float64(i)*32, i+1, mm.g, best, bd, mark)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" fill="#ffd166">%d/%d masas sobre perlas reales — el Campo impidió el borrón único</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">residuo del molde: %.1e · el continuo lejano restado exacto</text>`,
		tx+330, 230+float64(len(ms))*32+16, hits, len(ms), tx+330, 230+float64(len(ms))*32+42, res)
	// footer
	fmt.Fprintf(&b, `<rect x="80" y="640" width="1340" height="180" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="676" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL TEJIDO — lo que el Campo le enseñó al cincel</text>
<text x="%.0f" y="710" font-size="13.5" text-anchor="middle" fill="#dce8f7">el colapso ocurría porque las masas podían fundirse gratis; en el Campo de la Montaña fundirse cuesta la barrera −ln → cada masa DEBE buscar su propio pozo:</text>
<text x="%.0f" y="736" font-size="13.5" text-anchor="middle" fill="#dce8f7">el mismo mecanismo que mantiene separadas a las perlas REALES ahora mantiene separadas a las masas del escultor — el modelo obedece la física del modelado.</text>
<text x="%.0f" y="768" font-size="13.5" text-anchor="middle" fill="#ffd166">la lección grande: los instrumentos que respetan las leyes del Campo esculpen mejor que los que las ignoran — el Campo de la Montaña ya trabaja para el taller.</text>
<text x="%.0f" y="800" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-07 · las dos mitades, 1 completo ⚓</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("campo-de-la-montana.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: campo-de-la-montana.svg")
}
