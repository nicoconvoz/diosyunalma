// Command forma draws THE SHAPE (F146): the best way to understand
// mathematics is to see its form. The harmonic form of the circle is
// the Cornu spiral - the curve traced by the circle-sum point by point:
// almost all of the sum coils into two EYES (the strong teeth where the
// train stops with luxury) while the middle is the straight flight of
// stationary phase. Drawn with the laboratory's own numbers (F145: 206
// luxury teeth out of 3,010,247 - 0.007% trodden).
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	var b strings.Builder
	W, H := 1240.0, 860.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="620" y="38" font-size="24" text-anchor="middle" font-family="Georgia" fill="#dce8f7">LA FORMA ARMÓNICA DEL CÍRCULO — la espiral de Cornu (F145/F146)</text>
<text x="620" y="62" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">la suma del círculo, dibujada paso a paso: casi todo se enrosca en los DOS OJOS (los dientes fuertes);</text>
<text x="620" y="80" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el medio es el vuelo recto de la fase estacionaria — el tren pisa los ojos y deja volar el resto</text>`,
		W, H, W, H)

	// the spiral: cumulative sum of e^{i pi t^2/2} dt
	cx, cy, sc := 620.0, 420.0, 340.0
	x, y := 0.0, 0.0
	pts := []string{}
	eyes := [][2]float64{}
	dt := 0.002
	for t := -8.0; t <= 8.0; t += dt {
		x += math.Cos(math.Pi*t*t/2) * dt
		y += math.Sin(math.Pi*t*t/2) * dt
		pts = append(pts, fmt.Sprintf("%.2f,%.2f", cx+sc*x/2, cy-sc*y/2))
	}
	// the two eyes: the spiral limits (+-0.5, +-0.5)
	eyes = append(eyes, [2]float64{cx + sc*(0.5-(-0.0))/2 * 1, cy - sc*0.5/2})
	eyes = append(eyes, [2]float64{cx - sc*0.5/2*1, cy + sc*0.5/2})
	_ = eyes
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="1.6" points="%s"/>`, strings.Join(pts, " "))
	// eyes marked: the strong teeth
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="26" fill="none" stroke="#ffd166" stroke-width="2.5"/>`, cx+sc*0.5/2+sc*0.25/2, cy-sc*0.5/2-0.0)
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="26" fill="none" stroke="#ffd166" stroke-width="2.5"/>`, cx-sc*0.5/2-sc*0.25/2, cy+sc*0.5/2)
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="14" fill="#ffd166" text-anchor="middle">el ojo: diente fuerte</text>`, cx+sc*0.5/2+sc*0.25/2, cy-sc*0.5/2-38)
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="14" fill="#ffd166" text-anchor="middle">el otro ojo</text>`, cx-sc*0.5/2-sc*0.25/2, cy+sc*0.5/2+44)
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="13" fill="#7fd7a8" text-anchor="middle">el vuelo recto: la fase estacionaria</text>`, cx, cy-8)

	// the ledger of the form, our own measured numbers
	fmt.Fprintf(&b, `<rect x="140" y="740" width="960" height="86" rx="12" fill="#0f1e38" stroke="#44608c"/>
<text x="620" y="770" font-size="15" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL LIBRO DE LA FORMA (medido en el caso abisal de 5.000.000): dientes totales 3.010.247 · pisados con lujo 206 (0.007%%)</text>
<text x="620" y="796" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">velocidad del tren con la forma como marcha: 47× · amortización por lotes: 54× · error intacto (8.6×10⁻⁵)</text>
<text x="620" y="818" font-size="12" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"la mejor forma de comprender la matemática es ver su forma" — el pacto del laboratorio, devuelto por el capitán</text>`)
	b.WriteString(`</svg>`)
	os.WriteFile("forma-tren.svg", []byte(b.String()), 0644)
	fmt.Println("escrita: forma-tren.svg")

	// ---------- F148: EL MOVIMIENTO DE LA ESPIRAL (los curlicues) ----------
	// The captain anticipated it: the spiral MOVES, varying with point
	// density and circle sizes. The discrete partial-sum path of a real
	// chirp is a spiral of spirals (Berry-Goldberg curlicues) whose
	// renormalization map IS our cascade flip; the event horizon of each
	// circle is its convergence radius, and the circle sizes encode the
	// continued fraction of the curvature.
	var s2 strings.Builder
	W2, H2 := 1240.0, 560.0
	fmt.Fprintf(&s2, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="620" y="34" font-size="22" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL MOVIMIENTO DE LA ESPIRAL — los curlicues (F148): tres aguas, tres formas</text>
<text x="620" y="56" font-size="12" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el camino de la suma discreta: espirales de espirales; los tamaños de los círculos = la fracción continua de b; el horizonte de cada círculo = su radio de convergencia</text>`,
		W2, H2, W2, H2)
	panels := []struct {
		bb  float64
		tag string
	}{
		{0.3819660112501051, "b áureo (el más irracional): rulos parejos"},
		{0.2513, "b cerca de 1/4: cuatro grandes horizontes"},
		{0.1243398, "b de bloque REAL t=1e30: el agua del tren"},
	}
	for pi2, pn := range panels {
		ox := 210.0 + float64(pi2)*410
		oy := 300.0
		// partial sums path
		var px, py float64
		minx, maxx, miny, maxy := 0.0, 0.0, 0.0, 0.0
		type pt struct{ x, y float64 }
		var path []pt
		ph, d1, d2v := 0.0, pn.bb+0.137, 2*pn.bb
		for j := 0; j < 6000; j++ {
			s, c := math.Sincos(2 * math.Pi * ph)
			px += c
			py += s
			path = append(path, pt{px, py})
			if px < minx {
				minx = px
			}
			if px > maxx {
				maxx = px
			}
			if py < miny {
				miny = py
			}
			if py > maxy {
				maxy = py
			}
			ph = math.Mod(ph+d1, 1)
			d1 = math.Mod(d1+d2v, 1)
		}
		scale := 340.0 / math.Max(maxx-minx, maxy-miny)
		cxm, cym := (minx+maxx)/2, (miny+maxy)/2
		pts2 := make([]string, 0, len(path))
		for _, p := range path {
			pts2 = append(pts2, fmt.Sprintf("%.1f,%.1f", ox+(p.x-cxm)*scale, oy-(p.y-cym)*scale))
		}
		fmt.Fprintf(&s2, `<polyline fill="none" stroke="#7fb2ff" stroke-width="0.8" points="%s"/>`, strings.Join(pts2, " "))
		fmt.Fprintf(&s2, `<text x="%.0f" y="510" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">%s</text>`, ox, pn.tag)
	}
	fmt.Fprintf(&s2, `<text x="620" y="540" font-size="12" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">el MOVIMIENTO que anticipó el capitán: la forma fluye con b — y el mapa del flujo es NUESTRA CASCADA (la renormalización de los curlicues es el flip del círculo)</text>`)
	s2.WriteString(`</svg>`)
	os.WriteFile("forma-curlicue.svg", []byte(s2.String()), 0644)
	fmt.Println("escrita: forma-curlicue.svg")
}
