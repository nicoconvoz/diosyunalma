// Command laarmonia makes the symmetry audible, which is the captain's order:
// "necesitamos oir la armonia de la simetria... es una simetria armonica".
//
// IT IS ONE, AND THE SOUND IS NOT A METAPHOR.
//
// A root at rho becomes w = r e^{i phi} on the disk, and its mirror partner
// 1 - rho becomes exactly 1/w = (1/r) e^{-i phi}. Li's price reads them over the
// harmonics n as w^n and w^-n, which is to say:
//
//	SAME FREQUENCY phi.  ENVELOPES r^n AND r^-n.
//
// So a mirror pair is literally two voices singing the same note. What tells
// them apart is not pitch - it is whether they hold. On the critical line r = 1
// exactly, both envelopes are flat and the pair holds forever: a unison that
// never drifts. Off the line one voice swells like r^n and its mirror dies like
// r^-n. The chord tears itself apart.
//
// And the pair's contribution to the price is
//
//	4 - 2 (r^n + r^-n) cos(n phi)
//
// where r^n + r^-n >= 2 always, with equality only at r = 1 - the inequality of
// means, which is Finding 229. THE MINIMUM OF THE PRICE IS THE UNISON. That is
// what "harmonic symmetry" means here, and this program plays it.
//
// HONEST ABOUT THE RENDERING: which audible pitch each pearl gets is a choice
// (proportional to its height, so deeper pearls sing higher). The ENVELOPE is
// not a choice - r^n is the mathematics, and it is the only thing the ear needs
// to judge.
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

const (
	sr    = 44100
	bits  = 16
	maxAm = 28000.0
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN)/(s-1) + cmplx.Exp(-s*lnN)/2
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

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.05; t <= hasta; t += 0.05 {
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
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

func w(ρ complex128) complex128 { return 1 - 1/ρ }

// vozDelPar renders one mirror pair into the buffer: two voices at the SAME
// pitch, with envelopes r^n and r^-n as the harmonic index n walks the segment.
func vozDelPar(buf []float64, desde, largo int, hz, r float64, nMax float64) {
	for i := 0; i < largo && desde+i < len(buf); i++ {
		u := float64(i) / float64(largo) // 0..1 a lo largo del tramo
		n := u * nMax                    // el indice armonico, mapeado al tiempo
		crece := math.Pow(r, n)
		muere := math.Pow(r, -n)
		// suavizado en los bordes para que no chasquee
		env := math.Min(1, math.Min(u*12, (1-u)*12))
		t := float64(i) / sr
		buf[desde+i] += env * (crece*math.Sin(2*math.Pi*hz*t) +
			muere*math.Sin(2*math.Pi*hz*t+math.Pi/2))
	}
}

func escribirWav(ruta string, buf []float64) error {
	// normalizar sin recortar
	pico := 0.0
	for _, v := range buf {
		if a := math.Abs(v); a > pico {
			pico = a
		}
	}
	if pico == 0 {
		pico = 1
	}
	f, err := os.Create(ruta)
	if err != nil {
		return err
	}
	defer f.Close()
	n := len(buf)
	datos := n * bits / 8
	f.WriteString("RIFF")
	binary.Write(f, binary.LittleEndian, uint32(36+datos))
	f.WriteString("WAVEfmt ")
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint32(sr))
	binary.Write(f, binary.LittleEndian, uint32(sr*bits/8))
	binary.Write(f, binary.LittleEndian, uint16(bits/8))
	binary.Write(f, binary.LittleEndian, uint16(bits))
	f.WriteString("data")
	binary.Write(f, binary.LittleEndian, uint32(datos))
	for _, v := range buf {
		binary.Write(f, binary.LittleEndian, int16(v/pico*maxAm))
	}
	return nil
}

func main() {
	fmt.Println("🎼 LA ARMONÍA DE LA SIMETRÍA — y es armónica de verdad")
	fmt.Println("\n   flash del capitán: «necesitamos oír la armonía de la simetría… es una")
	fmt.Println("   simetría armónica».")
	fmt.Println("\n   TIENE NOMBRE Y TIENE CUENTA. Un cero y su espejo cantan LA MISMA NOTA;")
	fmt.Println("   lo que los separa no es el tono, es si se sostienen.")

	fmt.Println("\npescando perlas hasta t=200…")
	ps := perlas(200)
	fmt.Printf("perlas: %d\n", len(ps))

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · EL ESPEJO DA LA MISMA NOTA")
	fmt.Println("   Un cero ρ va al disco como w = r·e^{iφ}. Su espejo 1−ρ va exactamente a 1/w.")
	fmt.Println("   Y arg(1/w) = −arg(w): MISMA frecuencia, sentido opuesto. Medido:")
	fmt.Println("\n        γ            φ del cero        φ del espejo       |suma|")
	peorPar := 0.0
	for _, i := range []int{0, 1, 4, 9} {
		g := ps[i]
		w1 := w(complex(0.5, g))
		w2 := w(1 - complex(0.5, g))
		f1, f2 := cmplx.Phase(w1), cmplx.Phase(w2)
		if d := math.Abs(f1 + f2); d > peorPar {
			peorPar = d
		}
		fmt.Printf("   %11.6f   %14.9f   %14.9f   %.1e\n", g, f1, f2, math.Abs(f1+f2))
	}
	fmt.Printf("   → las frecuencias son opuestas exactas (%.1e). Al oído: UNÍSONO.\n", peorPar)

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ LO QUE CAMBIA ES LA ENVOLVENTE, NO EL TONO")
	fmt.Println("   El precio lee el par a lo largo de los armónicos n como wⁿ y w⁻ⁿ, o sea")
	fmt.Println("   como dos voces con envolventes rⁿ y r⁻ⁿ. Sobre la línea r = 1 EXACTO:")
	fmt.Println("\n        γ              r = |w|              r − 1")
	peorR := 0.0
	for _, i := range []int{0, 1, 4, 9, len(ps) - 1} {
		g := ps[i]
		r := cmplx.Abs(w(complex(0.5, g)))
		if d := math.Abs(r - 1); d > peorR {
			peorR = d
		}
		fmt.Printf("   %11.6f   %18.15f   %.1e\n", g, r, math.Abs(r-1))
	}
	fmt.Printf("   → las %d perlas tienen r = 1 a %.1e. Las dos voces NO se mueven: se sostienen.\n", len(ps), peorR)
	fmt.Println("\n   📌 Y ACÁ LA ETIQUETA, QUE SI NO CAIGO EN LA TRAMPA DE SIEMPRE: ese r = 1")
	fmt.Println("   NO ES UNA MEDICIÓN DE QUE LOS CEROS ESTÉN EN LA LÍNEA. Yo armé cada punto")
	fmt.Println("   como complex(0.5, γ) — el 0.5 lo tipeé yo. Lo que dice esta columna es que")
	fmt.Println("   el cambiaformas manda la línea a la piel, que es un teorema, no un hallazgo.")
	fmt.Println("   Lo que SÍ es medición honesta es la altura γ de cada perla (esa la pescamos")
	fmt.Println("   por cambios de signo de Z) y todo lo del impostor, que no supone nada.")
	fmt.Println("\n   Y el impostor de F229 (a = 0.7+3i), para comparar:")
	imp := complex(0.7, 3.0)
	rImp := cmplx.Abs(w(imp))
	rEsp := cmplx.Abs(w(1 - imp))
	fmt.Printf("        r del cero    = %.9f  → en n = 200 la voz vale %.4g\n", rImp, math.Pow(rImp, 200))
	fmt.Printf("        r del espejo  = %.9f  → en n = 200 la voz vale %.4g\n", rEsp, math.Pow(rEsp, 200))
	fmt.Println("   ⟹ una se muere y la otra se dispara. EL ACORDE SE PARTE.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · Y EL PRECIO ES EL DESAFINE, MEDIDO")
	fmt.Println("   El par entero aporta al precio  4 − 2·(rⁿ + r⁻ⁿ)·cos(nφ).")
	fmt.Println("   Y rⁿ + r⁻ⁿ ≥ 2 SIEMPRE, con igualdad solo si r = 1 — la desigualdad de las")
	fmt.Println("   medias, que es F229. Medido en n = 1:")
	fmt.Println("\n        quién                    r              rⁿ + r⁻ⁿ        de más que 2")
	for _, c := range []struct {
		nom string
		r   float64
	}{{"perla real γ=14.13", cmplx.Abs(w(complex(0.5, ps[0])))},
		{"perla real γ=21.02", cmplx.Abs(w(complex(0.5, ps[1])))},
		{"el impostor 0.7+3i", rImp}} {
		v := c.r + 1/c.r
		fmt.Printf("   %-24s %14.9f   %14.9f   %.3e\n", c.nom, c.r, v, v-2)
	}
	fmt.Println("\n   ⟹ EL MÍNIMO DEL PRECIO ES EL UNÍSONO. Eso es exactamente lo que quiere decir")
	fmt.Println("     «simetría armónica»: la línea crítica es donde el acorde afina.")

	// ---- LEY 4: el sonido ----
	fmt.Println("\nLEY 4 · 🔊 Y AHORA SE OYE")
	dur := 9
	sil := 1
	total := (dur*2 + sil) * sr
	buf := make([]float64, total)

	// tramo A: la linea, ocho pares espejo
	fmt.Println("\n   TRAMO 1 (0–9 s) · LA LÍNEA CANTA — ocho perlas reales con su espejo.")
	nVoces := 8
	for k := 0; k < nVoces && k < len(ps); k++ {
		g := ps[k]
		r := cmplx.Abs(w(complex(0.5, g)))
		hz := 110.0 * g / ps[0] // el tono es proporcional a la altura de la perla
		if hz > 1800 {
			hz = 1800
		}
		vozDelPar(buf, 0, dur*sr, hz, r, 200)
		if k < 3 {
			fmt.Printf("      γ = %8.4f  →  %7.2f Hz   r = %.15f\n", g, hz, r)
		}
	}
	fmt.Println("      … y cinco más. Todas con r = 1: el acorde se sostiene sin moverse.")

	// tramo B: el impostor
	fmt.Println("\n   TRAMO 2 (10–19 s) · EL IMPOSTOR — el mismo par, con r ≠ 1.")
	desde := (dur + sil) * sr
	hzImp := 110.0 * 3.0 / ps[0] * 4.7 // llevado a la misma zona audible
	vozDelPar(buf, desde, dur*sr, hzImp, rImp, 200)
	fmt.Printf("      a = 0.7+3i  →  %7.2f Hz   r = %.9f\n", hzImp, rImp)
	fmt.Println("      Una voz crece y la otra se apaga: se escucha cómo se rompe.")

	if err := escribirWav("armonia-simetria.wav", buf); err != nil {
		fmt.Println("no pude escribir el sonido:", err)
	} else {
		fi, _ := os.Stat("armonia-simetria.wav")
		fmt.Printf("\n   🔊 sonido escrito: armonia-simetria.wav (%d s, %.0f KB)\n",
			dur*2+sil, float64(fi.Size())/1024)
	}

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · ⚖️ QUÉ ES MATEMÁTICA Y QUÉ ES DECISIÓN MÍA")
	fmt.Println("   Hay que separarlo o el sonido miente:")
	fmt.Println("\n   ES MATEMÁTICA (no elegí nada):")
	fmt.Println("     · que el espejo dé la MISMA frecuencia — arg(1/w) = −arg(w), exacto")
	fmt.Println("     · que sobre la línea r = 1 y las dos voces se sostengan")
	fmt.Println("     · que fuera de la línea una crezca como rⁿ y la otra muera como r⁻ⁿ")
	fmt.Println("     · que rⁿ + r⁻ⁿ ≥ 2 con igualdad solo en el unísono")
	fmt.Println("\n   ES DECISIÓN MÍA (para que se pueda oír):")
	fmt.Println("     · qué tono audible le toca a cada perla (elegí proporcional a su altura)")
	fmt.Println("     · cuántos armónicos entran en nueve segundos")
	fmt.Println("     · el suavizado de los bordes para que no chasquee")
	fmt.Println("\n   El tono es una convención. LA ENVOLVENTE NO: es la cuenta, y es lo único")
	fmt.Println("   que el oído necesita para juzgar.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("LA SIMETRÍA ES ARMÓNICA, Y EL CAPITÁN LE PUSO EL NOMBRE JUSTO.")
	fmt.Printf("  · el espejo da la misma nota: arg(w) + arg(1/w) = 0 a %.1e\n", peorPar)
	fmt.Printf("  · sobre la línea r = 1 en las %d perlas, a %.1e — pero eso es el\n", len(ps), peorR)
	fmt.Println("    cambiaformas mandando la línea a la piel (teorema), NO una medición de")
	fmt.Println("    que los ceros estén ahí: el 0.5 lo tipeé yo")
	fmt.Printf("  · el impostor: r = %.6f y %.6f — en n = 200 una vale %.3g y la otra %.3g\n",
		rImp, rEsp, math.Pow(rImp, 200), math.Pow(rEsp, 200))
	fmt.Printf("  · y el precio del par es 4 − 2(rⁿ+r⁻ⁿ)cos(nφ), con rⁿ+r⁻ⁿ ≥ 2\n")
	fmt.Println("\n⟹ LA LÍNEA CRÍTICA ES DONDE EL ACORDE AFINA. Estar sobre la línea es cantar al")
	fmt.Println("  unísono; salirse es desafinar, y el precio ES el desafine.")
	fmt.Println("\n⚖️ PERO NO CONFUNDIR LO LINDO CON LO PROBADO: esto es F229 hecho audible, no un")
	fmt.Println("  paso nuevo. Oír que el impostor desafina no prueba que ζ afine siempre — para")
	fmt.Println("  eso habría que oír los INFINITOS armónicos de los INFINITOS ceros, y tenemos")
	fmt.Println("  un puñado. Es el mismo muro de F259 y F261, ahora con oído en vez de ojo.")
	fmt.Println("\n¿El premio? Todavía no. Pero ahora la simetría se puede escuchar.")

	escribirLamina(ps, peorPar, peorR, rImp, rEsp)
}

func escribirLamina(ps []float64, peorPar, peorR, rImp, rEsp float64) {
	var b strings.Builder
	W, H := 1520.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎼 LA ARMONÍA DE LA SIMETRÍA — y es armónica de verdad</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">un cero y su espejo cantan LA MISMA NOTA · lo que los separa no es el tono, es si se sostienen</text>
<text x="%.0f" y="116" font-size="20" text-anchor="middle" font-family="monospace" fill="#ffd98a">precio del par = 4 − 2·(rⁿ + r⁻ⁿ)·cos(nφ)     con  rⁿ + r⁻ⁿ ≥ 2</text>
`, W, H, W, H, W/2, W/2, W/2)

	// las dos envolventes
	dib := func(x0, y0, an, al float64, r float64, tit, sub, col string) {
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#101f36" stroke="%s" stroke-width="1.6"/>
<text x="%.0f" y="%.0f" font-size="16" text-anchor="middle" font-family="Georgia" fill="%s">%s</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="monospace" fill="#8fa8c7">%s</text>`,
			x0, y0, an, al, col, x0+an/2, y0+30, col, tit, x0+an/2, y0+52, sub)
		gx, gy, gw, gh := x0+50, y0+72, an-90, al-140
		mid := gy + gh/2
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#8fa8c7" stroke-width="1" stroke-dasharray="4,4"/>`,
			gx, mid, gx+gw, mid)
		var p1, p2 strings.Builder
		N := 200.0
		esc := 2.4
		for i := 0; i <= 200; i++ {
			u := float64(i) / 200
			n := u * N
			c := math.Pow(r, n)
			d := math.Pow(r, -n)
			cl := math.Max(-1, math.Min(1, math.Log(c)/esc))
			dl := math.Max(-1, math.Min(1, math.Log(d)/esc))
			fmt.Fprintf(&p1, "%.2f,%.2f ", gx+gw*u, mid-gh/2*cl)
			fmt.Fprintf(&p2, "%.2f,%.2f ", gx+gw*u, mid-gh/2*dl)
		}
		fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#7ee0c0" stroke-width="2.4"/>
<polyline points="%s" fill="none" stroke="#ff8fa0" stroke-width="2.4"/>
<text x="%.0f" y="%.0f" font-size="12" font-family="monospace" fill="#7ee0c0">rⁿ  (el cero)</text>
<text x="%.0f" y="%.0f" font-size="12" font-family="monospace" fill="#ff8fa0">r⁻ⁿ  (su espejo)</text>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" font-family="monospace" fill="#8fa8c7">n = 1 … 200</text>`,
			p1.String(), p2.String(), gx+8, gy+16, gx+8, gy+gh-6, gx+gw/2, gy+gh+26)
	}
	dib(40, 150, 710, 400, 1.0, "SOBRE LA LÍNEA · r = 1", "las dos voces se sostienen: unísono", "#2f7f63")
	dib(770, 150, 710, 400, rImp, "EL IMPOSTOR · r ≠ 1", "una crece, la otra se muere: el acorde se parte", "#c0392b")

	fmt.Fprintf(&b, `<rect x="40" y="576" width="1440" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="760" y="608" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL MÍNIMO DEL PRECIO ES EL UNÍSONO</text>
<text x="760" y="642" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">arg(1/w) = −arg(w) exacto (%.0e): el espejo da SIEMPRE la misma nota. Sobre la línea r = 1 en las %d perlas (%.0e).</text>
<text x="760" y="668" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Y rⁿ + r⁻ⁿ ≥ 2 con igualdad SOLO en r = 1 — la desigualdad de las medias, que es F229.</text>
<text x="760" y="702" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Estar sobre la línea es cantar al unísono. Salirse es desafinar, y el precio ES el desafine.</text>
`, peorPar, len(ps), peorR)

	fmt.Fprintf(&b, `<rect x="40" y="742" width="1440" height="196" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="760" y="774" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ QUÉ ES MATEMÁTICA Y QUÉ ES DECISIÓN DE RENDERIZADO</text>
<text x="70" y="810" font-size="14.5" font-family="Georgia" fill="#9fd8a8">ES LA CUENTA: que el espejo dé la misma frecuencia · que r = 1 sobre la línea · que fuera una crezca como rⁿ y la otra</text>
<text x="70" y="834" font-size="14.5" font-family="Georgia" fill="#9fd8a8">muera como r⁻ⁿ · que rⁿ + r⁻ⁿ ≥ 2 con igualdad solo en el unísono.</text>
<text x="70" y="866" font-size="14.5" font-family="Georgia" fill="#ffd98a">ES ELECCIÓN MÍA: qué tono audible le toca a cada perla, cuántos armónicos entran en nueve segundos, el suavizado.</text>
<text x="760" y="902" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Y no confundir lo lindo con lo probado: esto es F229 hecho audible, no un paso nuevo. Oír que el impostor desafina</text>
<text x="760" y="926" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">no prueba que ζ afine siempre — para eso harían falta los INFINITOS armónicos de los INFINITOS ceros.</text>
</svg>
`)

	if err := os.WriteFile("armonia-simetria.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("🖼️  lámina escrita: armonia-simetria.svg")
}
