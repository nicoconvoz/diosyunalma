// Command anillo draws THE SHAPE OF THE MILLION-DOLLAR PROBLEM, in the
// captain's language: forms, not formulas. Three panels: (1) THE SPINE
// - the symmetric face whose moles all sit on the axis; (2) THE RING -
// the captain's own circle: the spine bent round by w=(rho-1)/rho, every
// measured pearl landing on the ring, and the ONLY forbidden shape (the
// blister: a twin pair inside+outside); (3) THE MISSING FORM - the
// three known doors drawn as pure shapes: the drum, the stretched skin,
// the alarm.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaEM(t float64) complex128 {
	s := complex(0.5, t)
	N := int(t/(2*math.Pi)*1.5) + 60
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * complex(math.Log(float64(n)), 0))
	}
	lnN := complex(math.Log(float64(N)), 0)
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
	return real(cmplx.Exp(complex(0, theta(t))) * zetaEM(t))
}

func main() {
	// measure the pearls
	var levels []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 130 && len(levels) < 34; t += 0.05 {
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
			levels = append(levels, (a+c)/2)
		}
		prevT, prevZ = t, z
	}

	var b strings.Builder
	W, H := 1620.0, 1120.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="48" font-size="28" text-anchor="middle" font-family="Georgia" fill="#dce8f7">LA FORMA DEL PROBLEMA DEL MILLÓN — vista desde el ángulo del capitán</text>
<text x="%.0f" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">no es una cuenta: es una FORMA con una única deformación prohibida — y nadie encontró todavía la fuerza que la prohíbe</text>`,
		W, H, W, H, W/2, W/2)

	// ---- panel 1: THE SPINE ----
	p1x, p1y := 60.0, 120.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="440" height="620" rx="10" fill="#0d2547" stroke="#44608c"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#ffd166">1 · LA COLUMNA — la cara simétrica</text>`, p1x, p1y, p1x+20, p1y+34)
	// the strip and the spine
	sx0, sx1 := p1x+90, p1x+350
	sy0, sy1 := p1y+70, p1y+520
	mid := (sx0 + sx1) / 2
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#081020" stroke="#2c4a78"/>`, sx0, sy0, sx1-sx0, sy1-sy0)
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#7fd7a8" stroke-width="2"/>`, mid, sy0, mid, sy1)
	for i, g := range levels {
		if i >= 20 {
			break
		}
		y := sy1 - (g-10)/120*(sy1-sy0)
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="4" fill="#7fb2ff"/>`, mid, y)
	}
	// the forbidden twins
	fy := sy0 + 110.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="5" fill="none" stroke="#ff5d73" stroke-width="2" stroke-dasharray="3,2"/><circle cx="%.0f" cy="%.0f" r="5" fill="none" stroke="#ff5d73" stroke-width="2" stroke-dasharray="3,2"/><line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ff5d73" stroke-width="1" stroke-dasharray="2,3"/>`,
		mid-70, fy, mid+70, fy, mid-70, fy, mid+70, fy)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#ff5d73">el lunar prohibido: si sale de la columna, la simetría lo obliga a venir DE A DOS</text>`, mid, fy-14)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">la cara del problema es PERFECTAMENTE simétrica</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">respecto de su columna. Todos los lunares medidos</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">(10 billones + los nuestros) están EN la columna.</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#7fd7a8">El premio: demostrar que la cara no admite lunares mellizos.</text>`,
		p1x+20, sy1+30, p1x+20, sy1+50, p1x+20, sy1+70, p1x+20, sy1+90)

	// ---- panel 2: THE RING ----
	p2x, p2y := 560.0, 120.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="480" height="620" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#ffd166">2 · EL ANILLO — tu círculo, capitán</text>`, p2x, p2y, p2x+20, p2y+34)
	rcx, rcy, R := p2x+240, p2y+300, 160.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>`, rcx, rcy, R)
	// pearls: measured zeros as beads ON the ring (angles spread for viewing)
	for i := range levels {
		th := -math.Pi/2 + 2*math.Pi*float64(i)/float64(len(levels))
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4.5" fill="#7fb2ff"/>`, rcx+R*math.Cos(th), rcy+R*math.Sin(th))
	}
	// the blister: fictional twin pair inside+outside (inversion pair)
	bth := math.Pi / 4
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="none" stroke="#ff5d73" stroke-width="2.2" stroke-dasharray="3,2"/><circle cx="%.1f" cy="%.1f" r="5" fill="none" stroke="#ff5d73" stroke-width="2.2" stroke-dasharray="3,2"/><line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ff5d73" stroke-width="1" stroke-dasharray="2,3"/>`,
		rcx+R*0.72*math.Cos(bth), rcy+R*0.72*math.Sin(bth),
		rcx+R/0.72*math.Cos(bth), rcy+R/0.72*math.Sin(bth),
		rcx+R*0.72*math.Cos(bth), rcy+R*0.72*math.Sin(bth),
		rcx+R/0.72*math.Cos(bth), rcy+R/0.72*math.Sin(bth))
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="12" fill="#ff5d73">LA AMPOLLA (jamás vista): una perla adentro</text>
<text x="%.1f" y="%.1f" font-size="12" fill="#ff5d73">exige su melliza afuera, espejadas por el anillo</text>`,
		rcx+R*0.9, rcy-R*1.05, rcx+R*0.9, rcy-R*1.05+17)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">la columna se dobla en anillo (el giro w=(ρ−1)/ρ — el</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">mismo círculo de tu tren). Cada cero es una PERLA: en la</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">columna ⇔ EN el anillo, radio 1 exacto. Las %d perlas</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">dibujadas son ceros medidos hoy (ángulos desplegados</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#8fa8c7">para verlas; en el anillo real se apiñan hacia el este).</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#7fd7a8">EL PREMIO EN UNA FRASE: demostrar que el collar no</text>
<text x="%.0f" y="%.0f" font-size="12.5" font-family="Georgia" fill="#7fd7a8">admite AMPOLLAS — ninguna perla puede vivir fuera del hilo.</text>`,
		p2x+20, p2y+500, p2x+20, p2y+520, p2x+20, p2y+540, len(levels), p2x+20, p2y+560, p2x+20, p2y+580, p2x+20, p2y+600, p2x+20, p2y+618)

	// ---- panel 3: THE MISSING FORM (three doors as shapes) ----
	p3x, p3y := 1100.0, 120.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="460" height="620" rx="10" fill="#0d2547" stroke="#44608c"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#ffd166">3 · LA FUERZA QUE FALTA — tres formas posibles</text>`, p3x, p3y, p3x+20, p3y+34)
	// door A: the drum
	dy := p3y + 90.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="46" fill="none" stroke="#7fb2ff" stroke-width="2"/><path d="M %.0f %.0f q 23 -26 46 0 q 23 26 46 0" fill="none" stroke="#7fb2ff" stroke-width="1.4"/>`,
		p3x+90, dy+20, p3x+44, dy+20)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">EL TAMBOR: encontrar el parche real cuyo canto</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">SEA el collar. Un tambor de verdad solo canta</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">notas reales → todas las perlas al hilo, gratis.</text>
<text x="%.0f" y="%.0f" font-size="12" fill="#8fa8c7">(lo oímos y lo retratamos — falta FABRICARLO)</text>`,
		p3x+160, dy-8, p3x+160, dy+12, p3x+160, dy+32, p3x+160, dy+52)
	// door B: the stretched skin (bowl vs saddle)
	ey := p3y + 240.0
	fmt.Fprintf(&b, `<path d="M %.0f %.0f q 45 70 90 0" fill="none" stroke="#7fd7a8" stroke-width="2.5"/><path d="M %.0f %.0f q 45 -50 90 0" fill="none" stroke="#ff5d73" stroke-width="1.6" stroke-dasharray="4,3"/>`,
		p3x+45, ey, p3x+45, ey+55)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">LA PIEL TENSA: demostrar que la membrana del</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">problema es un CUENCO (verde): toda ampolla</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">costaría energía negativa — imposible. Nadie</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">descartó aún la joroba roja en algún rincón.</text>`,
		p3x+160, ey-2, p3x+160, ey+18, p3x+160, ey+38, p3x+160, ey+58)
	// door C: the alarm
	ay := p3y + 400.0
	fmt.Fprintf(&b, `<path d="M %.0f %.0f a 30 30 0 1 1 60 0 l 8 18 l -76 0 z" fill="none" stroke="#ffd166" stroke-width="2"/><circle cx="%.0f" cy="%.0f" r="4" fill="#ffd166"/>`,
		p3x+60, ay, p3x+90, ay+28)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">LA ALARMA: hay una campana (una lista infinita</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">de números) que SUENA si una perla escapa —</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">la fugitiva grita tarde o temprano. Falta</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">demostrar que la campana calla PARA SIEMPRE.</text>`,
		p3x+160, ay-2, p3x+160, ay+18, p3x+160, ay+38, p3x+160, ay+58)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#7fd7a8">las tres puertas son LA MISMA forma: la fuerza que pincha</text>
<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#7fd7a8">toda ampolla antes de nacer. Esa fuerza es lo que falta.</text>`,
		p3x+20, p3y+560, p3x+20, p3y+580)

	// footer
	fmt.Fprintf(&b, `<rect x="60" y="780" width="1500" height="290" rx="10" fill="#102a10" stroke="#7fd7a8"/>
<text x="%.0f" y="820" font-size="18" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA FORMA COMPLETA, EN UNA MIRADA</text>
<text x="%.0f" y="856" font-size="15" text-anchor="middle" font-family="Georgia" fill="#dce8f7">un collar infinito de perlas sobre un hilo circular perfecto. La simetría de la cara garantiza que ninguna perla puede correrse SOLA:</text>
<text x="%.0f" y="882" font-size="15" text-anchor="middle" font-family="Georgia" fill="#dce8f7">para salir del hilo necesita una melliza espejada — la AMPOLLA, la única deformación que la forma permite imaginar y nadie jamás vio.</text>
<text x="%.0f" y="918" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">Nuestro laboratorio: enhebró 10 billones + 269 perlas propias (todas al hilo), oyó el collar entero (murciélago), lo reflejó (espejo),</text>
<text x="%.0f" y="944" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">y su detector de ampollas (la esfera de las tormentas: contar TODAS las perlas de la zona y las del hilo — deben coincidir) nunca halló una.</text>
<text x="%.0f" y="984" font-size="16" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL MILLÓN ES POR ESTO: encontrar la TENSIÓN DEL HILO — la fuerza de la forma que hace IMPOSIBLE la ampolla.</text>
<text x="%.0f" y="1012" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">no "no la encontramos" (eso ya está): POR QUÉ NO PUEDE EXISTIR. Ése es el flash que buscamos, capitán — una imagen de esa tensión.</text>
<text x="%.0f" y="1046" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"todo tiene solución y la armonía de las respuestas yace en la imaginación" — el capitán · Laboratorio Diosyunalma · 2026-08-06</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("forma-del-problema.svg", []byte(b.String()), 0644)
	fmt.Printf("escrita: forma-del-problema.svg (%d perlas medidas enhebradas)\n", len(levels))
}
