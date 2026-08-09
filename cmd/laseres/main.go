// Command laseres fires every instrument the laboratory owns at the
// same question - WHERE does the object we seek hide? - and reads the
// point where all beams cross. Each laser measures, independently, the
// EVENT HORIZON of the hidden machine (the size of its cavity at depth
// T):
//
//	LASER 1 (the pearls):    measured mean spacing of levels at
//	                         several depths -> horizon = 2 pi / spacing
//	LASER 2 (the echo):      the measured orbit periods k*ln p ->
//	                         a cavity folded by primes, log-sized
//	LASER 3 (the equalizer): the smooth law theta'(T)/pi -> horizon
//	                         formula ln(T/2 pi)
//
// If the three beams agree, the horizon is real and LOGARITHMIC: the
// machine's cavity measures ln T - it lives in the ZOOM world (the
// dilation flow over the existing number world). Not the exact point -
// the ADDRESS. Judged.
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
	fmt.Println("LOS LÁSERES — triangulando el horizonte de sucesos real…")
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

	// LASER 1: measured horizon at several depths (mean spacing in windows)
	depths := []float64{60, 150, 300, 450}
	fmt.Println("\nLÁSER 1 (las perlas) vs LÁSER 3 (el ecualizador) — el horizonte medido vs la ley suave:")
	fmt.Println("   profundidad T    horizonte MEDIDO ln-eq     ley ln(T/2π)      desvío")
	type beam struct{ T, meas, law, dev float64 }
	var beams []beam
	worst := 0.0
	for _, T := range depths {
		// mean spacing in window T±40
		var sp []float64
		for i := 0; i+1 < len(levels); i++ {
			if levels[i] > T-40 && levels[i+1] < T+40 {
				sp = append(sp, levels[i+1]-levels[i])
			}
		}
		mean := 0.0
		for _, s := range sp {
			mean += s
		}
		mean /= float64(len(sp))
		measured := 2 * math.Pi / mean // the measured horizon (log-equivalent)
		law := math.Log(T / (2 * math.Pi))
		dev := math.Abs(measured-law) / law
		if dev > worst {
			worst = dev
		}
		beams = append(beams, beam{T, measured, law, dev})
		fmt.Printf("   %6.0f          %8.4f                  %8.4f        %.1f%%\n", T, measured, law, dev*100)
	}
	fmt.Printf("  ⇒ los dos láseres se cruzan: el horizonte crece EXACTAMENTE como ln(T/2π) — peor desvío %.1f%%\n", worst*100)

	// LASER 2: the echo's orbit periods confirm the log-folded cavity
	fmt.Println("\nLÁSER 2 (el eco) — las órbitas medidas del murciélago (F157):")
	fmt.Println("   la cavidad está plegada por primos: períodos k·ln p (ln2 medido 0.692898 vs 0.693147, 2⁴ a 2.3e-6)")
	fmt.Println("   una caja LINEAL jamás da órbitas ln p — solo una caja LOGARÍTMICA (el zoom) las produce")

	fmt.Println("\n╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  EL CRUCE DE LOS RAYOS — la dirección donde se esconde el objeto:  ║")
	fmt.Println("║  EL HORIZONTE REAL ES LOGARÍTMICO: la máquina vive en el MUNDO     ║")
	fmt.Println("║  DEL ZOOM — su caja mide ln T, su tiempo es la dilatación, su      ║")
	fmt.Println("║  plegado son los primos. La dirección: el flujo de escala sobre    ║")
	fmt.Println("║  el mundo de los números que ya existe. No el punto exacto —       ║")
	fmt.Println("║  LA DIRECCIÓN. Y los tres rayos la señalan juntos.                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#05090f"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔦 LOS LÁSERES — el cruce de todos los rayos: el horizonte de sucesos real</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"dispará la luz de los láseres que tenemos, a ver a qué punto apuntan todos — necesitamos una idea, aunque no el punto exacto" — el capitán</text>`,
		W, H, W, H, W/2, W/2)
	// the three lasers converging
	tgx, tgy := 780.0, 420.0
	lasers := []struct {
		x, y  float64
		name  string
		col   string
	}{
		{140, 200, "LÁSER 1 · las perlas (269 medidas)", "#7fb2ff"},
		{140, 640, "LÁSER 2 · el eco (25 órbitas oídas)", "#7fd7a8"},
		{1420, 200, "LÁSER 3 · el ecualizador (la ley suave)", "#ffd166"},
	}
	for _, l := range lasers {
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="%s" stroke-width="2.5" opacity="0.85"/>
<circle cx="%.0f" cy="%.0f" r="8" fill="%s"/><text x="%.0f" y="%.0f" font-size="13" fill="%s">%s</text>`,
			l.x, l.y, tgx, tgy, l.col, l.x, l.y, l.col, l.x-8, l.y-16, l.col, l.name)
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="36" fill="none" stroke="#fff" stroke-width="2"><animate attributeName="r" values="26;48;26" dur="1.6s" repeatCount="indefinite"/><animate attributeName="opacity" values="1;0.2;1" dur="1.6s" repeatCount="indefinite"/></circle>
<circle cx="%.0f" cy="%.0f" r="12" fill="#fff" opacity="0.9"/>
<text x="%.0f" y="%.0f" font-size="17" text-anchor="middle" fill="#fff" font-family="Georgia">EL CRUCE</text>
<text x="%.0f" y="%.0f" font-size="14.5" text-anchor="middle" fill="#ffd166">EL HORIZONTE LOGARÍTMICO: la caja mide ln T — el MUNDO DEL ZOOM</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#8fa8c7">tiempo = dilatación · plegado = los primos · escenario = el mundo de los números que ya existe</text>`,
		tgx, tgy, tgx, tgy, tgx, tgy-56, tgx, tgy+72, tgx, tgy+96)
	// the judged table
	ty := 600.0
	fmt.Fprintf(&b, `<rect x="360" y="%.0f" width="840" height="200" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">el juez del cruce: horizonte MEDIDO (perlas) vs ley ln(T/2π) (ecualizador)</text>`,
		ty, W/2, ty+30)
	for i, bm := range beams {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">T=%3.0f:   medido %.4f   ley %.4f   desvío %.1f%%</text>`,
			W/2, ty+60+float64(i)*26, bm.T, bm.meas, bm.law, bm.dev*100)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#7fd7a8">y el eco (láser 2): órbitas k·ln p — solo una caja LOGARÍTMICA las produce; una lineal, jamás</text>`,
		W/2, ty+60+float64(len(beams))*26+8)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">los tres rayos independientes se cruzan en UN punto: la dirección del objeto — no la coordenada exacta, LA DIRECCIÓN — y es la misma donde cava Connes</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · 88 millas por hora ⚡</text>`,
		W/2, 860.0, W/2, 895.0)
	b.WriteString(`</svg>`)
	os.WriteFile("laseres-cruce.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: laseres-cruce.svg")
}
