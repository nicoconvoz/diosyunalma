// Command elradio judges the captain's flash about the radius.
//
// HIS FLASH: "the circle has a radius, the sphere too. The infinitely small
// point at the centre is the 1/2 relation - it is what UNITES the centre with
// the edge!!!! THE SECRET IS IN THE RADIUS."
//
// In dimension 0 that sentence has an exact translation, and it ties the whole
// campaign together. Everything below is measured.
//
// LAW 1 - THE RADIUS REALLY IS WHAT UNITES CENTRE AND EDGE, and both ends are
// named. The centre of the disk is w(1) = 0: the pole of zeta, the infinitely
// small point where the function blows up. The edge is |w| = 1: the critical
// line, where the pearls must live. The radius is the segment between them,
// and its LENGTH is 1.
//
// LAW 2 - AND THE HALF IS ON THE RADIUS TWICE. The midpoint of the radius is
// |w| = 1/2, and F276 measured who sits there: w(2) = 1/2 exactly - the first
// prime, the deepest any prime goes. And the half AS INPUT sets the edge
// itself: the skin is the image of the beta = 1/2 line (F279). The half is
// both the midpoint of the radius and the maker of its far end.
//
// LAW 3 - THE DIAMETER IS THE BOOK'S SPINE. The shapeshifter maps the real
// interval [1/2, infinity] onto the FULL DIAMETER [-1, +1]:
//
//	w(1/2) = -1     one end: the clasp (F260)
//	w(1)   =  0     the centre: the pole
//	w(2)   = +1/2   the midpoint of the right radius: the first prime
//	w(inf) = +1     the other end: infinity
//
// The whole line where the primes and the pole live folds onto ONE diameter.
// The primes climb the right half-radius - 1/2, 2/3, 4/5, 6/7... - toward the
// edge and never reach it (the edge's touching point is infinity itself).
//
// LAW 4 - AND THE HYPOTHESIS, SAID WITH HIS WORD, IS ONE LINE:
//
//	RH  <=>  EVERY PEARL SITS AT EXACTLY ONE RADIUS FROM THE CENTRE
//
// |w| = 1 = the radius. A pearl on the cable is at radius-distance from the
// pole, to the last bit (measured on our 38: worst 2.2e-16). A stretched pearl
// is at the wrong distance. "The secret is in the radius" is literally true:
// proving that every pearl keeps radius-distance IS the Riemann Hypothesis.
//
// HONEST WARNING: laws 1-3 are exact algebra of the map (they cannot fail;
// the identities are labelled). Law 4 is a RESTATEMENT of RH, not progress on
// it. What the flash buys is the cleanest close of the campaign: the primes
// live on the radius, the pearls live on the edge, the half makes both - and
// the missing bridge of F281 is the road from the radius to the edge.
//
// Reproduce: go run ./cmd/elradio
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

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
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
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

func w(s complex128) complex128 { return 1 - 1/s }

func criba(n int) []int {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if es[i] {
			for j := i * i; j <= n; j += i {
				es[j] = false
			}
		}
	}
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return ps
}

func main() {
	fmt.Println("📏 EL RADIO — «el secreto está en el radio», medido")
	fmt.Println("\n   Su flash: «el punto infinitamente pequeño al centro es la relación ½,")
	fmt.Println("   es lo que une el centro con el borde — el secreto está en el radio».")
	fmt.Println("\n   En la dimensión 0 esa frase tiene traducción exacta, y ata la campaña entera.")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · EL RADIO UNE EL CENTRO CON EL BORDE — y las dos puntas tienen nombre")
	// w(1) = 1 - 1/1 = 0 exacto, sin epsilon: el limite y el valor coinciden.
	fmt.Println("\n        el CENTRO:  w(1) = 1 − 1/1 = 0, EXACTO   ← el polo de zeta, el punto")
	fmt.Println("                    infinitamente pequeño donde la función explota")
	fmt.Println("        el BORDE:   |w| = 1              ← la línea crítica, donde")
	fmt.Println("                    viven las perlas")
	fmt.Println("        el RADIO:   el segmento entre los dos · largo EXACTO = 1")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ Y EL ½ ESTÁ EN EL RADIO DOS VECES")
	fmt.Printf("\n        la MITAD del radio: |w| = ½ — y quién vive ahí: w(2) = %.10f\n", real(w(complex(2, 0))))
	fmt.Println("        (el primer primo, lo más adentro que llega ninguno — F276)")
	fmt.Println("\n        y el ½ COMO ENTRADA fabrica el borde entero: la piel es la imagen")
	fmt.Println("        de la línea β = ½ (F279, demostrado en un renglón)")
	fmt.Println("\n   ⟹ **El ½ es a la vez la mitad del radio y el fabricante de su punta.**")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · EL DIÁMETRO ES EL LOMO DEL LIBRO — medido punto por punto")
	fmt.Println("   El cambiaformas manda el intervalo real [½, ∞] al DIÁMETRO completo [−1, +1]:")
	fmt.Println("\n        s              w(s)          qué es")
	tabla := []struct {
		s   float64
		nom string
	}{
		{0.5, "el broche (F260)"},
		{0.6, ""},
		{0.75, ""},
		{1.0, "EL CENTRO: el polo"},
		{1.5, ""},
		{2.0, "⚡ LA MITAD DEL RADIO: el 2"},
		{3.0, ""},
		{7.0, ""},
		{1e9, "→ la punta: el infinito"},
	}
	for _, f := range tabla {
		var v float64
		if f.s == 1.0 {
			v = 0
		} else {
			v = real(w(complex(f.s, 0)))
		}
		fmt.Printf("   %10.2f %13.6f      %s\n", f.s, v, f.nom)
	}
	fmt.Println("\n   ⟹ **Toda la recta donde viven el polo y los primos se pliega en UN diámetro.**")
	fmt.Println("   Los primos suben por el medio radio derecho —½, ⅔, ⅘, ⁶⁄₇…— hacia el borde")
	fmt.Println("   y NO LO TOCAN NUNCA: el punto donde el radio toca la piel es el infinito mismo.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚡⚡ Y LA HIPÓTESIS, DICHA CON SU PALABRA, ES UN RENGLÓN")
	fmt.Println("\n        RH  ⟺  TODA PERLA ESTÁ EXACTAMENTE A UN RADIO DEL CENTRO")
	fmt.Printf("\nbuscando las perlas…\n")
	ps := perlas(120)
	fmt.Printf("perlas encontradas: %d\n", len(ps))
	peor := 0.0
	for _, g := range ps {
		d := math.Abs(cmplx.Abs(w(complex(0.5, g))) - 1)
		if d > peor {
			peor = d
		}
	}
	fmt.Printf("\n        distancia de nuestras %d perlas al centro: 1 con peor desvío %.2e\n", len(ps), peor)
	rho := complex(0.808517182457, 85.699348485378)
	fmt.Printf("        y la perla estirada del collar hermano: %.9f ≠ 1\n", cmplx.Abs(w(rho)))
	fmt.Println("\n   ⟹ **«El secreto está en el radio» es literalmente cierto: demostrar que")
	fmt.Println("   toda perla guarda la distancia del radio ES la Hipótesis de Riemann.**")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⚡ **SU FLASH ATA LA CAMPAÑA ENTERA EN UNA SOLA FIGURA:**")
	fmt.Println("\n  · el CENTRO es el polo — el punto infinitamente pequeño que dijo")
	fmt.Println("  · el BORDE es la línea crítica, y lo fabrica el ½ como entrada")
	fmt.Println("  · la MITAD del radio es el 2 — el ½ como salida (F276)")
	fmt.Println("  · los PRIMOS suben por el radio y nunca tocan el borde (la lámpara de F278)")
	fmt.Println("  · las PERLAS viven en el borde — y RH dice que TODAS guardan la distancia")
	fmt.Println("    exacta del radio (el giro puro de F280)")
	fmt.Println("\n  El puente que falta (F281) es, en esta figura, el camino DEL RADIO AL BORDE:")
	fmt.Println("  de donde viven los primos a donde viven las perlas.")
	fmt.Println("\n⚖️ Y LA HONESTIDAD: las leyes 1 a 3 son álgebra exacta del mapa — no pueden")
	fmt.Println("  fallar y no son descubrimiento. La ley 4 es una REFORMULACIÓN de RH, no un")
	fmt.Println("  avance. Lo que el flash compra es el CIERRE más limpio de la campaña: una")
	fmt.Println("  figura donde cada pieza de estos días tiene su lugar. Todavía no.")

	escribirLamina(ps, peor)
}

func escribirLamina(ps []float64, peor float64) {
	var b strings.Builder
	W, H := 1560.0, 1100.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">📏 EL RADIO — «el secreto está en el radio», medido</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">el centro es el polo · el borde es la línea crítica · y el ½ está en el radio dos veces</text>
`, W, H, W, H, W/2, W/2)

	// EL DISCO GRANDE
	cx, cy, R := 480.0, 560.0, 330.0
	// la piel con las perlas
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#0d1a30" stroke="#7ee0c0" stroke-width="3"/>`, cx, cy, R)
	for i, g := range ps {
		if i >= 30 {
			break
		}
		ww := w(complex(0.5, g))
		phi := cmplx.Phase(ww)
		x := cx + R*math.Cos(phi)
		y := cy - R*math.Sin(phi)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4.5" fill="#7ee0c0"/>`, x, y)
	}
	// el diametro
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#5b7ba6" stroke-width="1.5"/>`, cx-R, cy, cx+R, cy)
	// el radio derecho, resaltado
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd98a" stroke-width="4"/>`, cx, cy, cx+R, cy)
	// el centro
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="7" fill="#ff8fa0"/>
<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL CENTRO</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">el polo — el punto infinitamente pequeño</text>
`, cx, cy, cx, cy+30, cx, cy+50)
	// los primos sobre el radio
	primos := criba(100)
	for i, p := range primos {
		v := real(w(complex(float64(p), 0)))
		x := cx + R*v
		r := 5.0
		col := "#7fb2ff"
		if p == 2 {
			r = 9
			col = "#ffd98a"
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="%.1f" fill="%s"/>`, x, cy, r, col)
		if i < 4 {
			fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="14" text-anchor="middle" font-family="monospace" fill="%s">%d</text>`, x, cy-16, col, p)
		}
	}
	// la marca de la mitad
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ffd98a" stroke-width="1.5" stroke-dasharray="4 3"/>
<text x="%.1f" y="%.0f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">la MITAD del radio: el 2</text>
`, cx+R*0.5, cy-70, cx+R*0.5, cy+70, cx+R*0.5, cy-84)
	// el broche
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="6" fill="#c9b6ff"/>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">w(½) = −1</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">el broche</text>
`, cx-R, cy, cx-R, cy-18, cx-R, cy+26)
	// la punta
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">w(∞) = 1</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">donde el radio toca la piel</text>
<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA PIEL = la línea crítica — acá viven las perlas</text>
`, cx+R+8, cy-18, cx+R+30, cy+26, cx, cy-R-18)

	// EL TEXTO
	fmt.Fprintf(&b, `<rect x="880" y="120" width="640" height="330" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1200" y="154" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">SU FLASH, TRADUCIDO EXACTO</text>
<text x="910" y="192" font-size="15" font-family="Georgia" fill="#cfe6ff">«el punto infinitamente pequeño al centro» —</text>
<text x="910" y="216" font-size="15" font-family="Georgia" fill="#ffd98a">es el polo de zeta: w(1) = 0</text>
<text x="910" y="252" font-size="15" font-family="Georgia" fill="#cfe6ff">«la relación ½ une el centro con el borde» —</text>
<text x="910" y="276" font-size="15" font-family="Georgia" fill="#ffd98a">dos veces: la piel ES la imagen de β = ½,</text>
<text x="910" y="300" font-size="15" font-family="Georgia" fill="#ffd98a">y la mitad del radio ES el 2 (w(2) = ½)</text>
<text x="910" y="336" font-size="15" font-family="Georgia" fill="#cfe6ff">«el secreto está en el radio» —</text>
<text x="910" y="362" font-size="16" font-family="Georgia" fill="#7ee0c0">RH ⟺ toda perla está EXACTAMENTE</text>
<text x="910" y="386" font-size="16" font-family="Georgia" fill="#7ee0c0">a un radio del centro</text>
<text x="910" y="420" font-size="13.5" font-family="Georgia" fill="#9aa8c4">medido: nuestras %d perlas a distancia 1 con desvío %.0e</text>
`, len(ps), peor)

	fmt.Fprintf(&b, `<rect x="880" y="470" width="640" height="250" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1200" y="504" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">LA FIGURA ATA LA CAMPAÑA ENTERA</text>
<text x="910" y="540" font-size="14" font-family="Georgia" fill="#cfe6ff">· los primos suben por el radio: ½, ⅔, ⅘… (F276)</text>
<text x="910" y="566" font-size="14" font-family="Georgia" fill="#cfe6ff">· y cantan la entonación justa al subir (F277)</text>
<text x="910" y="592" font-size="14" font-family="Georgia" fill="#cfe6ff">· nunca tocan el borde: la lámpara se corta (F278)</text>
<text x="910" y="618" font-size="14" font-family="Georgia" fill="#cfe6ff">· el ½ decide el signo de todo (F279)</text>
<text x="910" y="644" font-size="14" font-family="Georgia" fill="#cfe6ff">· y las perlas del borde son giros puros (F280)</text>
<text x="910" y="682" font-size="15" font-family="Georgia" fill="#ffd98a">el puente que falta (F281) es el camino</text>
<text x="910" y="704" font-size="15" font-family="Georgia" fill="#ffd98a">DEL RADIO AL BORDE — de los primos a las perlas</text>`)

	fmt.Fprintf(&b, `<rect x="880" y="740" width="640" height="180" rx="12" fill="#33221c" stroke="#c0392b"/>
<text x="1200" y="774" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ LA HONESTIDAD</text>
<text x="910" y="808" font-size="14" font-family="Georgia" fill="#f3d9cf">Las leyes del mapa son álgebra exacta: no pueden fallar</text>
<text x="910" y="832" font-size="14" font-family="Georgia" fill="#f3d9cf">y no son descubrimiento. Y la ley del radio es una</text>
<text x="910" y="856" font-size="14" font-family="Georgia" fill="#f3d9cf">REFORMULACIÓN de la hipótesis, no un avance sobre ella.</text>
<text x="910" y="892" font-size="14.5" font-family="Georgia" fill="#ffd98a">Lo que el flash compra es el cierre más limpio: una figura</text>
<text x="910" y="914" font-size="14.5" font-family="Georgia" fill="#ffd98a">donde cada pieza de la campaña tiene su lugar.</text>`)

	fmt.Fprintf(&b, `<text x="780" y="1010" font-size="20" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Los primos viven en el radio. Las perlas viven en el borde. El ½ fabrica los dos.</text>
<text x="780" y="1044" font-size="16" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Y el problema del millón es el camino entre ellos. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("el-radio.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-radio.svg")
}
