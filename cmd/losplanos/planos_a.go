// planos_a.go - desk A of the museum tour: the bridge (wave -> chirp), the
// turn of the circle (stationary phase) and the descent (Euclid's cascade).
// Every number printed or drawn on these three plates is measured here.
package main

import (
	"fmt"
	"math"
)

// ---- shared drawing helpers for this desk (p01/p03 prefixes per house rule) ----

// p01caja fits a measured point cloud into a box, preserving aspect and
// flipping y so mathematical up is visual up.
func p01caja(pts [][2]float64, x, y, w, h float64) [][2]float64 {
	if len(pts) == 0 {
		return pts
	}
	x0, x1, y0, y1 := pts[0][0], pts[0][0], pts[0][1], pts[0][1]
	for _, p := range pts {
		x0, x1 = math.Min(x0, p[0]), math.Max(x1, p[0])
		y0, y1 = math.Min(y0, p[1]), math.Max(y1, p[1])
	}
	s := math.Min(w/math.Max(x1-x0, 1e-12), h/math.Max(y1-y0, 1e-12))
	ox, oy := x+(w-(x1-x0)*s)/2, y+(h-(y1-y0)*s)/2
	out := make([][2]float64, len(pts))
	for i, p := range pts {
		out[i] = [2]float64{ox + (p[0]-x0)*s, oy + (y1-p[1])*s}
	}
	return out
}

// p01dial draws a real gauge: 270-degree sweep, ticks and a needle standing at
// the measured value.
func p01dial(l *Lienzo, cx, cy, r, lo, hi, v float64, col, nombre, valor, unidad string) {
	l.Circ(cx, cy, r, cPanel, col, 2)
	for i := 0; i <= 12; i++ {
		a := (135 + 270*float64(i)/12) * math.Pi / 180
		q := 0.86
		if i%3 == 0 {
			q = 0.74
		}
		l.Line(cx+r*q*math.Cos(a), cy+r*q*math.Sin(a), cx+r*0.97*math.Cos(a), cy+r*0.97*math.Sin(a), cDim, 1.4, "")
	}
	f := math.Max(0, math.Min(1, (v-lo)/(hi-lo)))
	a := (135 + 270*f) * math.Pi / 180
	l.Line(cx, cy, cx+(r-12)*math.Cos(a), cy+(r-12)*math.Sin(a), cGold, 3.2, "")
	l.Circ(cx, cy, 5, cGold, "", 0)
	if nombre != "" {
		l.Txt(cx, cy-r-13, 15, col, "middle", nombre)
	}
	if valor != "" {
		l.Mono(cx, cy+r+21, 13.5, cInk, "middle", valor)
	}
	if unidad != "" {
		l.Txt(cx, cy+r+39, 12, cDim, "middle", unidad)
	}
}

// p04mil groups thousands with dots, the way the museum's labels read.
func p04mil(n int64) string {
	s, out := fmt.Sprintf("%d", n), ""
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out += "."
		}
		out += string(r)
	}
	return out
}

// p03dientes draws a row of teeth of the given count spanning width w.
func p03dientes(l *Lienzo, x, y, w, h float64, cuantos int, col string, sw float64) {
	l.Line(x, y+h, x+w, y+h, col, 1.6, "")
	for i := 0; i < cuantos; i++ {
		px := x + (float64(i)+0.5)/float64(cuantos)*w
		l.Line(px, y, px, y+h, col, sw, "")
	}
}

// ---------------------------------------------------------------------------
// PLATE 01 - EL PUENTE
// ---------------------------------------------------------------------------

func plano01() {
	const t, k0 = 1e24, 2.5e11
	c, b, g := puente(t, k0)
	L := bandL(t, k0)
	n := int(L)
	T := t / (2 * math.Pi)
	ulp := math.Nextafter(T/k0, 2*T/k0) - T/k0

	// the same block walked twice: the pure chirp and the honest zeta terms.
	// the difference between them is exactly the cubic remainder of ln(1+u).
	andar := make([][2]float64, 2*n)
	sag := make([][2]float64, n)
	var ar, ai, mr, mi float64
	maxd, maxj := 0.0, 0
	for j := 0; j < n; j++ {
		fj := float64(j)
		ph := math.Mod(b*fj*fj+c*fj, 1)
		s, co := math.Sincos(2 * math.Pi * ph)
		ar, ai = ar+co, ai+s
		u := fj / k0
		rem := -T * (u*u*u/3 - u*u*u*u/4)
		s2, co2 := math.Sincos(2 * math.Pi * math.Mod(ph+rem, 1))
		amp := 1 / math.Sqrt(1+u)
		mr, mi = mr+amp*co2, mi+amp*s2
		andar[j] = [2]float64{ar, ai}
		andar[n+j] = [2]float64{mr, mi}
		sag[j] = [2]float64{fj, rem}
		if d := math.Hypot(mr-ar, mi-ai); d > maxd {
			maxd, maxj = d, j
		}
	}
	fmt.Printf("SS plano01 t=1e24 k0=2.5e11 c=%.9f b=%.9f g=%.4e L=%d eta=%.4frad gL3=%.6fciclos sep_max=%.4f@j=%d sqrtN=%.2f ulp(T/k0)=%.2e\n",
		c, b, g, n, eta(t, k0, L), g*L*L*L, maxd, maxj, math.Sqrt(float64(n)), ulp)

	l := NewLienzo(1400, 950)
	l.Titulo("① EL PUENTE — de ola a chirp: por qué el mar de zeta ES un peine",
		"un pedazo de la suma de zeta, mirado con lupa, resulta ser una sola cuerda que sube de tono")

	// --- the sea: a real Z(t) window, measured ---
	l.Panel(36, 100, 420, 300, cBlue, "EL MAR — Z(t), la ola que hay que leer")
	var mar [][2]float64
	for i := 0; i <= 620; i++ {
		tt := 1000 + 24*float64(i)/620
		th := theta(tt)
		K := int(math.Sqrt(tt / (2 * math.Pi)))
		z := 0.0
		for kk := 1; kk <= K; kk++ {
			z += math.Cos(th-tt*math.Log(float64(kk))) / math.Sqrt(float64(kk))
		}
		mar = append(mar, [2]float64{tt, 2 * z})
	}
	zmax := 0.0
	for _, p := range mar {
		zmax = math.Max(zmax, math.Abs(p[1]))
	}
	px := func(tt float64) float64 { return escala(tt, 1000, 1024, 62, 434) }
	l.Line(62, 236, 434, 236, cGrid, 1.4, "")
	var mp [][2]float64
	for _, p := range mar {
		mp = append(mp, [2]float64{px(p[0]), 236 - p[1]/zmax*70})
	}
	l.Camino(mp, cBlue, 2, 0.95)
	ceros := 0
	for i := 1; i < len(mar); i++ {
		if mar[i-1][1]*mar[i][1] < 0 {
			l.Circ(px(mar[i][0]), 236, 3.6, cRose, "", 0)
			ceros++
		}
	}
	p03dientes(l, 100, 330, 320, 20, 26, cDim, 1.4)
	l.Line(196, 326, 284, 326, cGold, 2.4, "")
	l.Line(196, 326, 196, 336, cGold, 2.4, "")
	l.Line(284, 326, 284, 336, cGold, 2.4, "")
	l.Txt(240, 320, 12, cGold, "middle", "un BLOQUE")
	l.Txt(62, 372, 11, cRose, "start", fmt.Sprintf("● %d ceros en la ventana", ceros))
	l.Txt(430, 372, 11, cDim, "end", "los términos k = 1, 2, 3, …")
	l.Txt(62, 390, 11, cDim, "start", "ventana t ∈ [1000,1024] para que se vea; el tren lee la misma ola a t = 10²⁴")

	// --- the lens ---
	l.Flecha(462, 251, 524, 251, cGreen, 3)
	l.Curva(566, 172, 566, 330, 30, cGold, 2.6, "")
	l.Curva(566, 172, 566, 330, -30, cGold, 2.6, "")
	l.Line(566, 172, 566, 330, cGold, 1, "4 6")
	l.Flecha(608, 251, 692, 251, cGreen, 3)
	l.Txt(566, 158, 13.5, cGold, "middle", "LA LUPA")
	l.Mono(566, 356, 11.5, cInk, "middle", "ln(k0+j) = ln k0")
	l.Mono(566, 374, 11.5, cInk, "middle", "+ u − u²/2 + u³/3 …")
	l.Mono(566, 392, 11.5, cDim, "middle", "u = j/k0")

	// --- the three screws ---
	l.Panel(700, 100, 664, 300, cGold, "LOS TRES TORNILLOS que salen de la lupa")
	p01dial(l, 830, 232, 62, 0, 1, c, cGreen, "c — pendiente", fmt.Sprintf("%.6f", c), "vueltas por paso")
	p01dial(l, 1032, 232, 62, 0, 1, b, cBlue, "b — curvatura", fmt.Sprintf("%.6f", b), "vueltas por paso²")
	p01dial(l, 1234, 232, 62, 0, -5e-12, g, cRose, "g — torcedura", fmt.Sprintf("%.4e", g), "vueltas por paso³")
	l.Txt(1032, 384, 11.5, cDim, "middle", fmt.Sprintf("a esta hondura float64 ya sólo resuelve c a ±%.1e — por eso el DeLorean navega en doble-doble", ulp))

	// --- the same block, twice ---
	l.Panel(36, 420, 900, 300, cGreen, fmt.Sprintf("EL MISMO BLOQUE, DOS VECES — la caminata de sus %d términos", n))
	fit := p01caja(andar, 70, 446, 830, 236)
	l.Camino(fit[:n], cBlue, 4.8, 0.6)
	l.Camino(fit[n:], cRose, 1.4, 0.95)
	l.Circ(fit[0][0], fit[0][1], 5, cGold, "", 0)
	l.Line(fit[n-1][0], fit[n-1][1], fit[2*n-1][0], fit[2*n-1][1], cGold, 2, "3 3")
	l.Circ(fit[n-1][0], fit[n-1][1], 4.5, cBlue, "", 0)
	l.Circ(fit[2*n-1][0], fit[2*n-1][1], 4.5, cRose, "", 0)
	// zoom inset on the tail, where the third screw finally shows
	const M = 240
	zoom := append(append([][2]float64{}, andar[n-M:n]...), andar[2*n-M:]...)
	zf := p01caja(zoom, 100, 476, 212, 172)
	l.Rect(88, 462, 236, 200, cBg, cGrid, 0.9)
	l.Camino(zf[:M], cBlue, 4.4, 0.6)
	l.Camino(zf[M:], cRose, 1.6, 0.95)
	l.Line(zf[M-1][0], zf[M-1][1], zf[2*M-1][0], zf[2*M-1][1], cGold, 2, "3 3")
	l.Circ(zf[M-1][0], zf[M-1][1], 4.5, cBlue, "", 0)
	l.Circ(zf[2*M-1][0], zf[2*M-1][1], 4.5, cRose, "", 0)
	l.Txt(206, 678, 11, cDim, "middle", fmt.Sprintf("de cerca: los últimos %d pasos", M))
	l.Curva(330, 560, fit[2*n-1][0]-14, fit[2*n-1][1], -30, cGold, 1.4, "4 4")
	l.Txt(56, 704, 12, cBlue, "start", "azul grueso: el chirp e^(2πi(b j² + c j))")
	l.Txt(370, 704, 12, cRose, "start", "rosa fino: los términos de zeta")
	l.Txt(916, 704, 12, cGold, "end", fmt.Sprintf("se despegan al final: %.2f de √N = %.1f", maxd, math.Sqrt(float64(n))))

	// --- the sag ---
	l.Panel(960, 420, 404, 300, cRose, "LA TORCEDURA QUE SOBRA — g·j³")
	l.Ejes(1000, 466, 330, 204, "j", "")
	var sp [][2]float64
	for _, p := range sag {
		sp = append(sp, [2]float64{escala(p[0], 0, float64(n), 1000, 1330), escala(p[1], 0, sag[n-1][1], 466, 650)})
	}
	l.Camino(sp, cRose, 2.6, 0.95)
	l.Line(1000, 650, 1330, 650, cGold, 1.2, "5 4")
	l.Mono(1326, 644, 10.5, cGold, "end", fmt.Sprintf("presupuesto η = %.3f rad = %.5f vueltas", eta(t, k0, L), -g*L*L*L))
	l.Txt(1006, 600, 11.5, cDim, "start", "crece como j³ — por eso el bloque se corta")

	l.Formula(36, 736, 1328, "fase(j)/2π = −T·ln(k0+j) = cte + c·j + b·j² + g·j³ + …    con T = t/2π,  u = j/k0,  c = frac(−T/k0),  b = frac(T/2k0²),  g = −T/3k0³")
	l.Nota(36, 782, 1328, cGreen, []string{
		"EN CRIOLLO: el mar de zeta se arma sumando muchas olitas, una por cada número entero. Si agarrás un pedacito de esa suma y lo mirás con lupa, todas las olitas resultan ser LA MISMA",
		"cuerda, que va subiendo de tono de a poquito: eso es un chirp, el ruido del avión que despega. Y una cuerda que sube parejo se describe con tres tornillos nada más: con qué velocidad",
		"arranca (c), cuánto acelera (b) y cuánto se tuerce al final (g). Abajo está el mismo bloque caminado dos veces: en azul el chirp de dos tornillos, en rosa los términos verdaderos.",
		"Van pegados casi todo el viaje y recién al final se separan — esa separación ES el tercer tornillo. Por eso el bloque tiene tamaño máximo: se corta antes de que la torcedura se note.",
	})
	l.Pie(fmt.Sprintf("medido en vivo: t = 10²⁴, k0 = 2,5·10¹¹, bloque de %d términos — el presupuesto del tren", n))
	l.Guardar("laminas/plano-01-el-puente.svg")
}

// ---------------------------------------------------------------------------
// PLATE 03 - LA VUELTA
// ---------------------------------------------------------------------------

func plano03() {
	m1a, m2a, larA := dualRango(0.031, 0.29, 200000)
	c, b, _ := puente(1e24, 2.5e11)
	m1b, m2b, larB := dualRango(b, c, 120000)
	// drawing case: gentle curvature so the eye can follow every single arrow.
	const bd, cd = 0.02, -0.28
	xm := -cd / (2 * bd)
	R, ancho := 13.0, 1/math.Sqrt(2*bd)
	fmt.Printf("SS plano03 dualA(b=0.031,c=0.29,N=200000) m=[%d..%d] largo=%d vs 2bN=%.1f | dualB(b=%.9f,c=%.9f,N=120000) m=[%d..%d] largo=%d vs 2bN=%.2f | dibujo x_m=%.2f ancho=1/sqrt(2b)=%.2f\n",
		m1a, m2a, larA, 2*0.031*200000, b, c, m1b, m2b, larB, 2*b*120000, xm, ancho)

	l := NewLienzo(1400, 950)
	l.Titulo("③ LA VUELTA — dónde se alinean las flechas (fase estacionaria)",
		"de N flechas que casi todas se cancelan sobreviven sólo 2bN: ahí nace la vuelta del círculo")

	// --- parabola with spinning arrows ---
	l.Panel(36, 100, 720, 350, cBlue, "LA FASE ES UNA PARÁBOLA — y cada término lleva su flecha")
	fx := func(x float64) float64 { return escala(x, xm-R, xm+R, 74, 728) }
	fy := func(p float64) float64 { return escala(p, -1.6, 2.8, 424, 148) }
	l.Ejes(74, 148, 654, 276, "j", "fase (vueltas)")
	l.Rect(fx(xm-ancho/2), 150, fx(xm+ancho/2)-fx(xm-ancho/2), 272, cGold, cGold, 0.11)
	var par [][2]float64
	for i := 0; i <= 400; i++ {
		x := xm - R + 2*R*float64(i)/400
		par = append(par, [2]float64{fx(x), fy(bd*x*x + cd*x)})
	}
	l.Camino(par, cDim, 2, 0.75)
	l.Line(74, fy(bd*xm*xm+cd*xm), 728, fy(bd*xm*xm+cd*xm), cGreen, 2, "7 5")
	for i := 0; i < 40; i++ {
		x := xm - R + 2*R*float64(i)/39
		ph := 2 * math.Pi * frac(bd*x*x+cd*x)
		gx, gy := fx(x), fy(bd*x*x+cd*x)
		dx, dy := 15*math.Cos(ph), 15*math.Sin(ph)
		col := cRose
		if math.Abs(x-xm) <= ancho/2 {
			col = cGold
		}
		l.Flecha(gx-dx, gy-dy, gx+dx, gy+dy, col, 2)
	}
	l.Circ(fx(xm), fy(bd*xm*xm+cd*xm), 5, cGold, "", 0)
	l.Txt(fx(xm), 144, 13, cGold, "middle", "acá las flechas se alinean")
	l.Txt(96, fy(bd*xm*xm+cd*xm)-11, 12, cGreen, "start", "tangente de pendiente m")
	l.Mono(fx(xm), 442, 12, cGold, "middle", fmt.Sprintf("x_m = (m−c)/(2b) = %.1f     ancho ≈ 1/√(2b) = %.2f", xm, ancho))

	// --- the Cornu walk ---
	l.Panel(780, 100, 584, 350, cRose, "LA CAMINATA — lejos de x_m todo se enrosca y se anula")
	var cor [][2]float64
	var cr, ci float64
	for x := xm - R; x <= xm+R; x += 0.02 {
		s, co := math.Sincos(2 * math.Pi * (bd*x*x + cd*x))
		cr, ci = cr+co*0.02, ci+s*0.02
		cor = append(cor, [2]float64{cr, ci})
	}
	cf := p01caja(cor, 812, 146, 520, 268)
	l.Camino(cf, cRose, 2.2, 0.95)
	l.Circ(cf[0][0], cf[0][1], 5, cBlue, "", 0)
	l.Circ(cf[len(cf)-1][0], cf[len(cf)-1][1], 5, cGreen, "", 0)
	h := len(cf) / 2
	l.Flecha(cf[h-45][0], cf[h-45][1], cf[h+45][0], cf[h+45][1], cGold, 3)
	l.Txt(1072, 434, 12, cDim, "middle", "los dos rulos son el resto del bloque anulándose solo; el tramo derecho es todo el aporte")

	// --- N thin teeth -> 2bN fat teeth ---
	l.Panel(36, 470, 1328, 232, cGreen, "EL PEINE SE AFLOJA — de N dientes finos quedan 2bN dientes gruesos")
	finos := 124
	gruesos := int(math.Round(float64(finos) * float64(larA) / 200000))
	p03dientes(l, 90, 506, 1160, 36, finos, cBlue, 1.2)
	l.Mono(1262, 532, 13, cBlue, "start", "N = 200.000")
	l.Flecha(670, 552, 670, 590, cGold, 3)
	l.Mono(690, 578, 13, cGold, "start", fmt.Sprintf("×2b — el peine se afloja al %.1f%%", 100*float64(larA)/200000))
	p03dientes(l, 90, 594, 1160, 36, gruesos, cRose, 4.5)
	l.Mono(1262, 620, 13, cRose, "start", fmt.Sprintf("2bN = %d", larA))
	l.Mono(90, 660, 12, cDim, "start", fmt.Sprintf("medido A — b=0.031  c=0.29  N=200.000  →  m ∈ [%d .. %d], largo %d   (2bN = %.1f)", m1a, m2a, larA, 2*0.031*200000))
	l.Mono(90, 682, 12, cDim, "start", fmt.Sprintf("medido B — bloque real de t=10²⁴: b=%.6f  N=120.000  →  m ∈ [%d .. %d], largo %d   (2bN = %.2f)", b, m1b, m2b, larB, 2*b*120000))

	l.Formula(36, 716, 1328, "x_m = (m−c)/(2b)   ·   curvatura dual = −1/(4b)   ·   amplitud = 1/√(2b)   ·   giro fijo = e^(iπ/4)   ⇒   N flechas se vuelven 2bN flechas")
	l.Nota(36, 762, 1328, cGreen, []string{
		"EN CRIOLLO: sumar el bloque es sumar N flechitas, una por término, y cada flecha apunta a un ángulo que gira cada vez más rápido. Las flechas que giran rápido se comen entre ellas: se anulan.",
		"Lo único que sobrevive son los pocos lugares donde el giro se frena y varias flechas apuntan casi para el mismo lado; ahí sí se suman de verdad. Esos lugares están espaciados parejo y son 2bN,",
		"muchísimos menos que N. Así que en vez de sumar N cosas se arma una lista nueva, mucho más corta: el peine fino se convierte en peine grueso y el trabajo se derrumba. Eso es una vuelta del círculo.",
		"La caminata de la derecha lo muestra de un vistazo: los dos rulos de las puntas son términos que se cancelan solos, y el tramo derecho del medio es todo el aporte que queda.",
	})
	l.Pie("medido en vivo — la vuelta acorta el peine por el factor 2b, y se puede volver a dar")
	l.Guardar("laminas/plano-03-la-vuelta.svg")
}

// ---------------------------------------------------------------------------
// PLATE 04 - EL DESCENSO
// ---------------------------------------------------------------------------

func plano04() {
	const comodo = 1000
	bs, ns := cascada(0.199689, 0.437, 5000000, comodo, true)
	k := len(bs)
	fmt.Printf("SS plano04 b0=0.199689 c0=0.437 n0=5000000 comodo=%d vueltas=%d escalera:", comodo, k)
	for i := range bs {
		fmt.Printf(" (b=%.9f,n=%d)", bs[i], ns[i])
	}
	fmt.Printf(" recorte=%.1fx\n", float64(ns[0])/float64(ns[k-1]))

	l := NewLienzo(1400, 950)
	l.Titulo("④ EL DESCENSO Y SU TACÓMETRO — la cascada de Euclides",
		"cada vuelta del círculo es un escalón: el peine se afloja, el problema se achica, hasta que entra en el bote")

	// --- the staircase ---
	l.Panel(36, 100, 1000, 620, cBlue, "LA ESCALERA — un escalón por vuelta; al margen, el tacómetro de b")
	l.Txt(560, 142, 11.5, cBlue, "middle", "el peine de ese nivel — ancho ∝ log n")
	l.Mono(906, 142, 11.5, cDim, "start", "términos")
	lo := math.Log10(float64(ns[k-1])) - 0.35
	hi := math.Log10(float64(ns[0]))
	ancho := make([]float64, k)
	caida, salto := 0, 0.0
	for i := range bs {
		ancho[i] = escala(math.Log10(float64(ns[i])), lo, hi, 80, 640)
		if i > 0 {
			if d := math.Log10(float64(ns[i-1])) - math.Log10(float64(ns[i])); d > salto {
				caida, salto = i, d
			}
		}
	}
	for i := range bs {
		y := 176 + float64(i)*80
		col := cBlue
		if i == k-1 {
			col = cGreen
		}
		dientes := int(ancho[i] / 12)
		if dientes < 3 {
			dientes = 3
		}
		p03dientes(l, 250, y-22, ancho[i], 22, dientes, col, 1.3)
		p01dial(l, 110, y, 30, 0, 0.26, bs[i], cGold, "", "", "")
		l.Mono(148, y+4, 12, cGold, "start", fmt.Sprintf("%.6f", bs[i]))
		l.Mono(906, y+4, 12, cDim, "start", p04mil(ns[i]))
		if i < k-1 {
			l.Line(250+ancho[i], y, 250+ancho[i+1], y, cGold, 1.6, "4 4")
			l.Line(250+ancho[i+1], y, 250+ancho[i+1], y+80, cGold, 1.6, "4 4")
		}
	}
	if caida > 0 {
		l.Txt(250+ancho[caida]+16, 176+float64(caida)*80-34, 12, cRose, "start",
			fmt.Sprintf("aguja casi en cero (b=%.4f) ⇒ caída de %.0fx", bs[caida-1], math.Pow(10, salto)))
	}

	// --- the continued-fraction signature ---
	l.Panel(1056, 100, 308, 380, cRose, "LA FIRMA — los b de cada vuelta")
	l.Ejes(1090, 150, 250, 280, "vuelta", "b")
	var fp [][2]float64
	for i, v := range bs {
		fp = append(fp, [2]float64{escala(float64(i), 0, float64(k-1), 1090, 1340), escala(v, 0, 0.26, 430, 150)})
	}
	l.Camino(fp, cRose, 2.4, 0.95)
	for _, p := range fp {
		l.Circ(p[0], p[1], 4, cGold, "", 0)
	}
	l.Mono(1090, 466, 11, cDim, "start", "b ← 1/(4b) mod 1 — la firma de Euclides")

	// --- the comfortable size, with its rower ---
	l.Panel(1056, 500, 308, 220, cGreen, "EL BOTE")
	l.Txt(1210, 552, 15, cGold, "middle", fmt.Sprintf("TAMAÑO CÓMODO = %d", comodo))
	l.Mono(1210, 574, 12.5, cInk, "middle", fmt.Sprintf("n final = %d ≤ %d", ns[k-1], comodo))
	l.Txt(1210, 600, 11.5, cDim, "middle", "acá se rema derecho, término por término")
	l.Camino([][2]float64{{1150, 666}, {1162, 684}, {1258, 684}, {1270, 666}}, cGold, 2.6, 1)
	l.Line(1210, 666, 1210, 638, cInk, 2.4, "")
	l.Circ(1210, 630, 8, cInk, "", 0)
	l.Line(1210, 648, 1140, 690, cInk, 2, "")
	l.Line(1210, 648, 1280, 690, cInk, 2, "")
	l.Line(1132, 686, 1148, 696, cGold, 4, "")
	l.Line(1288, 686, 1272, 696, cGold, 4, "")
	l.Curva(1080, 694, 1340, 694, 7, cBlue, 2, "")
	l.Curva(1080, 706, 1340, 706, -7, cBlue, 2, "")

	l.Formula(36, 736, 1328, "b ← 1/(4b) (mod 1)   ·   c ← α/(2b) con α = ⌈c⌉−c   ·   n ← ⌊c+2b(n−1)⌋−⌈c⌉+1 ≈ 2bn   ·   multiplicador ← multiplicador · e^(iπ/4)/√(2b)")
	l.Nota(36, 782, 1328, cGreen, []string{
		fmt.Sprintf("EN CRIOLLO: la vuelta del círculo no se da una sola vez, se da de nuevo y de nuevo. Cada vuelta cambia el peine por uno más flojo y deja un problema más chico: %s términos, después %s,", p04mil(ns[0]), p04mil(ns[1])),
		fmt.Sprintf("después %s… hasta que quedan %d y ya conviene remar a mano. Son %d escalones para bajar de cinco millones a menos de mil: %.0f veces menos trabajo, y ni una cuenta de más.", p04mil(ns[2]), ns[k-1], k, float64(ns[0])/float64(ns[k-1])),
		"El tacómetro del margen marca el b de cada nivel, que es la aguja del motor: cuando la aguja queda casi en cero el escalón siguiente se desploma, porque el peine nuevo tiene poquísimos dientes.",
		"Esa lista de b es exactamente la misma cuenta que hacía Euclides para el máximo común divisor, dos mil trescientos años atrás. La máquina baja la escalera de Euclides y abajo la espera el bote.",
	})
	l.Pie("medido en vivo — cascada(b=0,199689  c=0,437  n=5.000.000) con la cura de cizalla puesta")
	l.Guardar("laminas/plano-04-el-descenso.svg")
}
