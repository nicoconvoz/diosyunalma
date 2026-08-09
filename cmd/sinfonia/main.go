// Command sinfonia is the flash of a flash of a flash: the whole
// laboratory turned into sound around the carrier's counting. Three
// staves, all real data:
//
//	THE SEA    - a drone whose breath is |Z(t)|: it swells between
//	             pearls and DIES exactly at each zero;
//	THE COUNT  - one bell per pearl (the carrier's winding, audible):
//	             269 chimes accelerating as the pearls tighten; each
//	             bell's pitch is the local spacing vs the expected
//	             mean (tight pair = high tense bell, wide gap = low);
//	THE EQUALIZER - the prime orchestra underneath: each prime p a
//	             band pulsing with its own explicit-formula wave
//	             cos(t ln p), amplitude ln p / sqrt(p).
//
// Listen to the carrier count. Plus the equalizer sheet (SVG score).
package main

import (
	"encoding/binary"
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
	// ---- compute the sea and find the pearls ----
	fmt.Println("LA SINFONÍA DEL CONTEO — componiendo con el mar, las perlas y los primos…")
	t0, t1 := 5.0, 506.0
	dt := 0.02
	nG := int((t1 - t0) / dt)
	sea := make([]float64, nG)
	for i := 0; i < nG; i++ {
		sea[i] = zOf(t0 + float64(i)*dt)
	}
	var pearls []float64
	for i := 1; i < nG; i++ {
		if sea[i]*sea[i-1] < 0 {
			pearls = append(pearls, t0+(float64(i)-0.5)*dt)
		}
	}
	fmt.Printf("mar calculado (%d muestras) · perlas: %d campanadas\n", nG, len(pearls))

	// ---- the audio ----
	const sr = 44100
	dur := 76.0
	rate := (t1 - t0) / (dur - 4)
	buf := make([]float64, int(dur*sr))
	tOf := func(tau float64) float64 { return t0 + tau*rate }
	// THE SEA: drone amplitude = |Z| (smoothed), dying at each pearl
	maxZ := 0.0
	for _, v := range sea {
		if math.Abs(v) > maxZ {
			maxZ = math.Abs(v)
		}
	}
	for i := range buf {
		tau := float64(i) / sr
		tt := tOf(tau)
		if tt >= t1 {
			break
		}
		gi := int((tt - t0) / dt)
		if gi >= nG {
			gi = nG - 1
		}
		amp := math.Abs(sea[gi]) / maxZ
		buf[i] += 0.30 * amp * (math.Sin(2*math.Pi*220*tau) + 0.5*math.Sin(2*math.Pi*110*tau))
	}
	// THE COUNT: one bell per pearl; pitch = local spacing vs mean
	bell := func(tau0, f, amp float64) {
		i0 := int(tau0 * sr)
		n := int(1.1 * sr)
		for i := 0; i < n && i0+i < len(buf); i++ {
			tt := float64(i) / sr
			env := math.Exp(-3.2*tt) * math.Min(1, tt/0.004)
			buf[i0+i] += amp * env * (math.Sin(2*math.Pi*f*tt) + 0.35*math.Sin(2*math.Pi*2.02*f*tt))
		}
	}
	for k, g := range pearls {
		var gap float64
		if k+1 < len(pearls) {
			gap = pearls[k+1] - g
		} else {
			gap = pearls[k] - pearls[k-1]
		}
		mean := 2 * math.Pi / math.Log(g/(2*math.Pi))
		ratio := mean / gap
		f := 440 * math.Pow(2, math.Log2(ratio))
		if f < 220 {
			f = 220
		}
		if f > 1760 {
			f = 1760
		}
		bell((g-t0)/rate, f, 0.34)
	}
	// THE EQUALIZER: prime bands pulsing their explicit-formula waves
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for i := range buf {
		tau := float64(i) / sr
		tt := tOf(tau)
		if tt >= t1 {
			break
		}
		var v float64
		for _, p := range primes {
			lp := math.Log(float64(p))
			w := lp / math.Sqrt(float64(p))
			v += w * (0.5 + 0.5*math.Cos(tt*lp)) * math.Sin(2*math.Pi*(72*lp)*tau)
		}
		buf[i] += 0.10 * v
	}
	// normalize + write
	peak := 1e-9
	for _, v := range buf {
		if math.Abs(v) > peak {
			peak = math.Abs(v)
		}
	}
	pcm := make([]int16, len(buf))
	for i, v := range buf {
		pcm[i] = int16(v / peak * 0.9 * 32767)
	}
	wf, _ := os.Create("sinfonia-conteo.wav")
	dataLen := uint32(len(pcm) * 2)
	wf.WriteString("RIFF")
	binary.Write(wf, binary.LittleEndian, uint32(36+dataLen))
	wf.WriteString("WAVEfmt ")
	binary.Write(wf, binary.LittleEndian, uint32(16))
	binary.Write(wf, binary.LittleEndian, uint16(1))
	binary.Write(wf, binary.LittleEndian, uint16(1))
	binary.Write(wf, binary.LittleEndian, uint32(sr))
	binary.Write(wf, binary.LittleEndian, uint32(sr*2))
	binary.Write(wf, binary.LittleEndian, uint16(2))
	binary.Write(wf, binary.LittleEndian, uint16(16))
	wf.WriteString("data")
	binary.Write(wf, binary.LittleEndian, dataLen)
	binary.Write(wf, binary.LittleEndian, pcm)
	wf.Close()
	fmt.Printf("escrita: sinfonia-conteo.wav (%.0f s · %d campanadas del portador · ecualizador de %d primos)\n", dur, len(pearls), len(primes))

	// ---- the score / equalizer sheet ----
	var b strings.Builder
	W, H := 1560.0, 860.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="44" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎼 LA SINFONÍA DEL CONTEO — la partitura (76 s · %d campanadas · ecualizador de %d primos)</text>
<text x="%.0f" y="70" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el mar respira |Z| y MUERE en cada perla · cada muerte, una campana del portador (aguda = par apretado, grave = hueco ancho) · debajo, cada primo pulsa su onda cos(t·ln p)</text>`,
		W, H, W, H, W/2, len(pearls), len(primes), W/2)
	// the sea envelope with bells
	px, pw, py, ph := 80.0, 1400.0, 110.0, 300.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="8" fill="#081020" stroke="#44608c"/>`, px, py, pw, ph)
	pts := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		gi := i * nG / 1000
		v := math.Abs(sea[gi]) / maxZ
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", px+pw*float64(i)/1000, py+ph-8-v*(ph-30)))
	}
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="1.2" points="%s"/>`, strings.Join(pts, " "))
	for _, g := range pearls {
		x := px + pw*(g-t0)/(t1-t0)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ffd166" stroke-width="0.8" opacity="0.75"/>`, x, py+ph-8, x, py+ph-26)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">el aliento del mar |Z(t)| (azul) — cada marca dorada: una campanada del portador; escuchá cómo se APURAN: las perlas se aprietan al bajar</text>`,
		W/2, py+ph+26)
	// the equalizer
	ey := py + ph + 60
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="15" font-family="Georgia" fill="#ffd166">EL ECUALIZADOR DE LOS PRIMOS — cada banda pulsa su onda de la fórmula explícita</text>`, px, ey)
	for i, p := range primes {
		lp := math.Log(float64(p))
		w := lp / math.Sqrt(float64(p))
		bw := w / 0.635 * 320
		y := ey + 26 + float64(i)*30
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Consolas,monospace" fill="#8fa8c7">p=%d</text>
<rect x="%.0f" y="%.0f" width="%.1f" height="18" rx="4" fill="#7fd7a8" opacity="0.75"/>
<text x="%.0f" y="%.0f" font-size="11.5" fill="#8fa8c7">tono %.0f Hz · pulso cos(t·ln %d) · peso ln p/√p = %.3f</text>`,
			px, y+14, p, px+64, y, bw, px+80+bw+12, y+14, 72*lp, p, w)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">tres voces, todas de datos reales: el mar (nuestra Z), el conteo (las perlas del portador), la orquesta (los primos con sus pesos exactos) — el laboratorio entero, sonando</text>`,
		W/2, ey+26+float64(len(primes))*30+24)
	b.WriteString(`</svg>`)
	os.WriteFile("sinfonia-conteo.svg", []byte(b.String()), 0644)
	fmt.Println("escrita: sinfonia-conteo.svg (la partitura)")
}
