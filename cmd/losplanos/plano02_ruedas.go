package main

import (
	"fmt"
	"math"
	"time"
)

// p02caminata returns the running partial sums of sum_{n<q} exp(i*pi*p*n^2/q),
// mirroring the train's rail-1 gaussSum: n^2 is never formed, it is carried
// modulo 2q by differences, so the phase argument is an EXACT integer over q.
func p02caminata(p, q, sign int64) [][2]float64 {
	m := 2 * q
	pts := make([][2]float64, 0, q+1)
	var re, im float64
	var n2, dn int64 = 0, 1
	pts = append(pts, [2]float64{0, 0})
	for n := int64(0); n < q; n++ {
		ph := math.Pi * float64((p*n2)%m) / float64(q) * float64(sign)
		s, c := math.Sincos(ph)
		re += c
		im += s
		pts = append(pts, [2]float64{re, im})
		n2 = (n2 + dn) % m
		dn = (dn + 2) % m
	}
	return pts
}

func p02suma(p, q, sign int64) (float64, float64) {
	w := p02caminata(p, q, sign)
	return w[len(w)-1][0], w[len(w)-1][1]
}

// plano02 draws rail 1: the exact reciprocity of 1893, as two walks that land
// on the same point - one of q steps and one of p steps.
func plano02() {
	const p, q = 3, 2000
	izq := p02caminata(p, q, +1)
	der := p02caminata(q, p, -1)
	f := math.Sqrt(float64(q) / float64(p))
	c4, s4 := math.Sqrt2/2, math.Sqrt2/2
	dx, dy := der[len(der)-1][0], der[len(der)-1][1]
	predX := f * (c4*dx - s4*dy)
	predY := f * (c4*dy + s4*dx)
	realX, realY := izq[len(izq)-1][0], izq[len(izq)-1][1]
	err := math.Hypot(realX-predX, realY-predY) / math.Hypot(realX, realY)

	// the same identity at the scale the machine actually uses: a wheel of a
	// hundred million teeth answered by a wheel of seven.
	t0 := time.Now()
	gr, gi := p02suma(7, 1e8, +1)
	dtGrande := time.Since(t0)
	t1 := time.Now()
	hr, hi := p02suma(1e8, 7, -1)
	dtChica := time.Since(t1)
	gf := math.Sqrt(1e8 / 7)
	gpx := gf * (c4*hr - s4*hi)
	gpy := gf * (c4*hi + s4*hr)
	gErr := math.Hypot(gr-gpx, gi-gpy) / math.Hypot(gr, gi)
	modGrande := math.Hypot(gr, gi)
	modChica := math.Hypot(hr, hi)

	fmt.Printf("SS plano02 · %d pasos vs %d: err %.1e · rueda 1e8 vs 7: |S| = %.6f (√q = %.0f), err %.1e, %.3fs vs %.6fs\n",
		q, p, err, modGrande, math.Sqrt(1e8), gErr, dtGrande.Seconds(), dtChica.Seconds())

	l := NewLienzo(1400, 950)
	l.Titulo("② LAS DOS RUEDAS — riel 1 del tren: la reciprocidad EXACTA de 1893",
		"dos caminatas de largo abismalmente distinto que terminan, al último bit, en el mismo punto")

	// ---- left: the long walk (the curlicue) ----
	l.Panel(46, 104, 620, 556, cGold, "LA CAMINATA LARGA — un paso por término")
	minX, maxX, minY, maxY := 0.0, 0.0, 0.0, 0.0
	for _, pt := range izq {
		minX, maxX = math.Min(minX, pt[0]), math.Max(maxX, pt[0])
		minY, maxY = math.Min(minY, pt[1]), math.Max(maxY, pt[1])
	}
	pad := 0.10 * math.Max(maxX-minX, maxY-minY)
	cx, cy := (minX+maxX)/2, (minY+maxY)/2
	half := math.Max(maxX-minX, maxY-minY)/2 + pad
	mapa := func(pt [2]float64) [2]float64 {
		return [2]float64{escala(pt[0], cx-half, cx+half, 76, 636), escala(pt[1], cy-half, cy+half, 630, 150)}
	}
	scr := make([][2]float64, len(izq))
	for i, pt := range izq {
		scr[i] = mapa(pt)
	}
	l.Camino(scr, cGreen, 1.25, 0.92)
	l.Circ(scr[0][0], scr[0][1], 5, cInk, "", 0)
	l.Txt(scr[0][0]+11, scr[0][1]+17, 11.5, cDim, "start", "salida (n = 0)")
	fin := scr[len(scr)-1]
	l.Circ(fin[0], fin[1], 9, "none", cRose, 2.6)
	l.Circ(fin[0], fin[1], 3.5, cRose, "", 0)
	l.Txt(fin[0]+13, fin[1]-11, 12.5, cRose, "start", "llegada")
	l.Txt(356, 636, 12, cDim, "middle", fmt.Sprintf("%d pasos, cada uno girando más rápido que el anterior (n² acelera el giro)", q))

	// ---- right: the short walk, and the two operations of the identity ----
	l.Panel(690, 104, 664, 556, cGreen, "LA CAMINATA CORTA — y la MISMA llegada")
	sminX, smaxX, sminY, smaxY := 0.0, 0.0, 0.0, 0.0
	for _, pt := range der {
		sminX, smaxX = math.Min(sminX, pt[0]), math.Max(smaxX, pt[0])
		sminY, smaxY = math.Min(sminY, pt[1]), math.Max(smaxY, pt[1])
	}
	spad := 0.5 * math.Max(smaxX-sminX, smaxY-sminY)
	scx, scy := (sminX+smaxX)/2, (sminY+smaxY)/2
	shalf := math.Max(smaxX-sminX, smaxY-sminY)/2 + spad
	smapa := func(pt [2]float64) [2]float64 {
		return [2]float64{escala(pt[0], scx-shalf, scx+shalf, 726, 986), escala(pt[1], scy-shalf, scy+shalf, 430, 170)}
	}
	sscr := make([][2]float64, len(der))
	for i, pt := range der {
		sscr[i] = smapa(pt)
	}
	for i := 0; i+1 < len(sscr); i++ {
		l.Flecha(sscr[i][0], sscr[i][1], sscr[i+1][0], sscr[i+1][1], cBlue, 2.6)
	}
	l.Circ(sscr[0][0], sscr[0][1], 5, cInk, "", 0)
	l.Circ(sscr[len(sscr)-1][0], sscr[len(sscr)-1][1], 8, "none", cBlue, 2.4)
	l.Txt(856, 462, 12, cDim, "middle", fmt.Sprintf("%d pasos. Eso es TODO lo que hay que sumar.", p))

	// the gearing: rotate 45 degrees, then stretch
	gx, gy := 1180.0, 268.0
	l.Circ(gx, gy, 52, "none", cGold, 2.2)
	for i := 0; i < 12; i++ {
		a := float64(i) / 12 * 2 * math.Pi
		l.Line(gx+48*math.Cos(a), gy+48*math.Sin(a), gx+58*math.Cos(a), gy+58*math.Sin(a), cGold, 2.4, "")
	}
	l.Flecha(gx, gy, gx+40*math.Cos(-math.Pi/4), gy+40*math.Sin(-math.Pi/4), cGold, 2.6)
	l.Curva(gx+46, gy-46, gx-46, gy-46, 26, cGold, 2, "5 4")
	l.Txt(gx, gy-72, 13, cGold, "middle", "① girar 45°")
	l.Mono(gx, gy+80, 12.5, cDim, "middle", "e^{iπ/4}")
	l.Txt(gx, gy+118, 13, cGold, "middle", "② estirar")
	l.Flecha(gx-58, gy+140, gx+58, gy+140, cGold, 2.4)
	l.Mono(gx, gy+168, 12.5, cDim, "middle", fmt.Sprintf("×√(q/p) = %.2f", f))

	l.Rect(700, 500, 644, 148, cPanel, cRose, 0.45)
	l.Txt(1022, 526, 13.5, cRose, "middle", "LAS DOS LLEGADAS, UNA AL LADO DE LA OTRA")
	l.Mono(1022, 552, 12.5, cInk, "middle", fmt.Sprintf("larga (%d pasos)          → (%+.7f, %+.7f)", q, realX, realY))
	l.Mono(1022, 574, 12.5, cInk, "middle", fmt.Sprintf("corta, girada y estirada → (%+.7f, %+.7f)", predX, predY))
	l.Mono(1022, 600, 13, cGreen, "middle", fmt.Sprintf("diferencia relativa: %.1e — el ruido del coseno, nada más", err))
	l.Mono(1022, 628, 12, cDim, "middle", "las fases son enteras exactas: p·n² mod 2q, arrastrado por diferencias")

	// ---- the wheel measured at machine scale ----
	l.Panel(46, 676, 1308, 92, cBlue, "Y A ESCALA DE MÁQUINA — la rueda de cien millones de dientes contra la rueda de siete")
	wy := 730.0
	l.Circ(150, wy, 26, "none", cBlue, 2)
	for i := 0; i < 90; i++ { // a rim so fine it reads as smooth
		a := float64(i) / 90 * 2 * math.Pi
		l.Line(150+25*math.Cos(a), wy+25*math.Sin(a), 150+29*math.Cos(a), wy+29*math.Sin(a), cBlue, 0.8, "")
	}
	l.Mono(150, wy+46, 11, cDim, "middle", "10⁸ dientes")
	l.Circ(300, wy, 20, "none", cGreen, 2)
	for i := 0; i < 7; i++ {
		a := float64(i) / 7 * 2 * math.Pi
		l.Line(300+18*math.Cos(a), wy+18*math.Sin(a), 300+30*math.Cos(a), wy+30*math.Sin(a), cGreen, 3, "")
	}
	l.Mono(300, wy+46, 11, cDim, "middle", "7 dientes")
	l.Mono(370, wy-8, 12.5, cInk, "start", fmt.Sprintf("|suma de 10⁸ términos| = %.6f   y   √q = %.0f", modGrande, math.Sqrt(1e8)))
	l.Mono(370, wy+14, 12.5, cInk, "start", fmt.Sprintf("|suma de 7 términos| = %.6f = √7      diferencia relativa entre las dos: %.1e", modChica, gErr))
	l.Mono(370, wy+36, 12.5, cGold, "start", fmt.Sprintf("tiempo: %.3f s la rueda grande · %.6f s la chica — %.0f millones de veces más rápido",
		dtGrande.Seconds(), dtChica.Seconds(), dtGrande.Seconds()/math.Max(dtChica.Seconds(), 1e-9)/1e6))

	l.Formula(46, 782, 1308, "Σ_{n<q} e^{iπpn²/q}  =  √(q/p) · e^{iπ/4} · Σ_{n<p} e^{−iπqn²/p}      (Landsberg–Schaar, 1893, con pq par)")

	l.Nota(46, 826, 1308, cGreen, []string{
		"EN CRIOLLO: caminás dos mil pasos girando cada vez un poco más y terminás en un punto. Después caminás TRES pasos, girás la hoja 45 grados,",
		"estirás con la regla — y estás parado exactamente en el mismo punto. No cerca: el mismo. Por eso el tren no rema el mar: cambia la caminata larga",
		"por la corta cuando la corta dice lo mismo. Nació de un flash del capitán (F137): «un numerito chiquito arriba apuntando a un número gigante abajo».",
	})
	l.Pie("go run ./cmd/losplanos · lámina ② del recorrido de las máquinas · nació en F137-F138 · Todavía no.")
	l.Guardar("laminas/plano-02-las-dos-ruedas.svg")
}
