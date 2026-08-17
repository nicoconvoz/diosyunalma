package main

import (
	"fmt"
	"math"
	"strings"
)

// p05mil groups an integer with Spanish thousand dots.
func p05mil(n int64) string {
	s := fmt.Sprintf("%d", n)
	out := make([]byte, 0, len(s)+len(s)/3)
	for i := 0; i < len(s); i++ {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, s[i])
	}
	return string(out)
}

// p05coma turns a decimal point into the Rioplatense comma (floats only).
func p05coma(s string) string { return strings.ReplaceAll(s, ".", ",") }

// p05entrada is the curvature the cascade sees at the door: frac, then the
// mirror b>1/2 -> 1-b. The shear acts AFTER this, so this is what gets cut.
func p05entrada(b float64) float64 {
	b = frac(b)
	if b > 0.5 {
		b = 1 - b
	}
	return b
}

// p05escalera turns a measured ladder into staircase pixels: horizontal =
// how much n shrank (log), vertical = where the curvature jumped to.
func p05escalera(bs []float64, ns []int64, fx func(int64) float64, fy func(float64) float64) [][2]float64 {
	pts := [][2]float64{{fx(ns[0]), fy(bs[0])}}
	for i := 0; i+1 < len(bs); i++ {
		pts = append(pts, [2]float64{fx(ns[i+1]), fy(bs[i])})
		pts = append(pts, [2]float64{fx(ns[i+1]), fy(bs[i+1])})
	}
	return pts
}

// plano05 draws the same case taking two destinies: without the shear it drowns
// in the fixed point b=1/2; with it, one cut and it is home.
func plano05() {
	const nIni, nCom int64 = 1000000, 1000
	type caso struct{ b, c float64 }
	var elegido caso
	var sinB, conB []float64
	var sinN, conN []int64
	mejor, res := -1.0, ""
	for _, k := range []caso{{0.4999, 0.437}, {0.082150, 0.437}} {
		sb, sn := cascada(k.b, k.c, nIni, nCom, false)
		cb, cn := cascada(k.b, k.c, nIni, nCom, true)
		res += fmt.Sprintf("b=%.6f sin %d vueltas (n→%d) / con %d (n→%d) · ", k.b, len(sb), sn[len(sn)-1], len(cb), cn[len(cn)-1])
		if r := float64(len(sb)) / float64(len(cb)); r > mejor {
			mejor, elegido, sinB, sinN, conB, conN = r, k, sb, sn, cb, cn
		}
	}
	ent := p05entrada(elegido.b)
	fmt.Printf("SS plano05 · %selegido b=%.6f: %d vs %d vueltas (×%.0f) · entra b=%.6f → cizallada b=%.6g\n",
		res, elegido.b, len(sinB), len(conB), mejor, ent, conB[0])

	l := NewLienzo(1400, 950)
	l.Titulo("⑤ EL CORTE DE CIZALLA — el mismo caso, dos destinos",
		"la misma curvatura entra por dos puertas: una la manda al remolino del ½, la otra la deja en el fondo de un solo tijeretazo")

	dibuja := func(px, py float64, color, tit string, bs []float64, ns []int64, cizalla bool) {
		l.Panel(px, py, 640, 500, color, tit)
		x0, x1, y0, y1 := px+66, px+620, py+85, py+385
		fx := func(n int64) float64 { return escala(math.Log10(float64(n)), 6.06, 2.0, x0, x1) }
		fy := func(b float64) float64 { return escala(b, 0.55, 0.0, y0, y1) }
		l.Ejes(x0, y0, x1-x0, y1-y0, "", "b")
		for _, tk := range []float64{0, 0.25, 0.5} {
			l.Mono(x0-8, fy(tk)+4, 10.5, cDim, "end", p05coma(fmt.Sprintf("%.2f", tk)))
		}
		for k, e := range []int64{1000000, 100000, 10000, 1000, 100} {
			l.Mono(fx(e), y1+19, 10.5, cDim, "middle", []string{"10⁶", "10⁵", "10⁴", "10³", "10²"}[k])
		}
		l.Txt((x0+x1)/2, y1+42, 11.5, cDim, "middle", "n — cuántos términos quedan por sumar (escala log, cae hacia la derecha)")
		l.Line(x0, fy(0.25), x1, fy(0.25), cGold, 1.7, "8 6")
		l.Txt(x1-4, fy(0.25)+18, 11.5, cGold, "end", "b = ¼ — la puerta de la cizalla")
		if !cizalla {
			l.Txt((x0+x1)/2, fy(0.25)+90, 12.5, cDim, "middle", "esta mitad del tablero queda vacía: la escalera nunca baja hasta acá")
		}
		l.Line(x0, fy(0.5), x1, fy(0.5), cRose, 1.2, "3 5")
		l.Line(fx(nCom), y0, fx(nCom), y1, cGreen, 1.2, "4 5")
		l.Txt(fx(nCom), y0-9, 11, cGreen, "middle", "puerto: n ≤ 1000")
		ex, ey := fx(ns[0]), fy(ent)
		if cizalla {
			l.Flecha(ex, ey, ex, fy(bs[0])-7, cRose, 3.4) // the cut the ladder never records
			sx, sy := ex+70, fy(0.25)
			l.Line(sx-44, sy-27, sx+36, sy+23, cGold, 3.4, "")
			l.Line(sx-44, sy+27, sx+36, sy-23, cGold, 3.4, "")
			l.Circ(sx-52, sy-33, 11, "none", cGold, 2.6)
			l.Circ(sx-52, sy+33, 11, "none", cGold, 2.6)
			l.Circ(sx-4, sy, 3.4, cGold, "", 0)
			l.Txt(sx+46, sy+5, 12.5, cGold, "start", "la cizalla: b − ½")
			wx, wy := px+230, py+140 // the half turn the scissors take away
			for i := 0; i < 2; i++ {
				var w [][2]float64
				for k := 0; k <= 60; k++ {
					a := float64(k) / 60 * math.Pi
					w = append(w, [2]float64{wx + float64(i)*110 + a/math.Pi*110, wy - 26*math.Sin(a)*float64(1-2*i)})
				}
				if i == 0 {
					l.Camino(w, cGreen, 2.4, 0.95)
				} else {
					l.Camino(w, cRose, 2.4, 0.95)
				}
			}
			l.Line(wx+110, wy-36, wx+110, wy+36, cRose, 1.6, "4 4")
			l.Txt(wx+110, wy+56, 12, cRose, "middle", "le saca media vuelta a cada término")
			l.Txt(wx+110, wy+74, 12, cDim, "middle", "y la suma no se entera: j(j−1) es par")
		}
		l.Circ(ex, ey, 5.5, cInk, "", 0)
		if cizalla {
			l.Txt(ex+10, ey-9, 11.5, cInk, "start", p05coma(fmt.Sprintf("entra b = %.4f", ent)))
		}
		pts := p05escalera(bs, ns, fx, fy)
		l.Camino(pts, color, 3, 0.95)
		fin := pts[len(pts)-1]
		if cizalla {
			l.Circ(pts[0][0], pts[0][1], 5, cGreen, "", 0)
			l.Txt(pts[0][0]+12, pts[0][1]-10, 11.5, cGreen, "start", p05coma(fmt.Sprintf("cae a b = %.4f", bs[0])))
			l.Circ(fin[0], fin[1], 8, "none", cGreen, 2.6)
			l.Txt(fin[0]-12, fin[1]-14, 12, cGreen, "end", fmt.Sprintf("llega a n = %s", p05mil(ns[len(ns)-1])))
		} else { // the whirlpool: the ladder never leaves the fixed point
			cxw, cyw := (pts[0][0]+fin[0])/2+30, (pts[0][1]+fin[1])/2
			var esp [][2]float64
			for k := 0; k <= 300; k++ {
				a, r := float64(k)*0.075, 56*(1-float64(k)/300)
				esp = append(esp, [2]float64{cxw + r*math.Cos(a), cyw + r*math.Sin(a)*0.72})
			}
			l.Camino(esp, cRose, 1.7, 0.9)
			l.Line(fin[0]-6, fin[1]-6, fin[0]+6, fin[1]+6, cRose, 3, "")
			l.Line(fin[0]+6, fin[1]-6, fin[0]-6, fin[1]+6, cRose, 3, "")
			l.Txt(cxw, cyw-58, 12.5, cRose, "middle", "punto fijo: 1/(4·½) = ½")
			l.Txt(cxw-6, cyw+62, 12, cInk, "middle", p05coma(fmt.Sprintf("entra b = %.4f y ahí queda", ent)))
			l.Txt(cxw-6, cyw+80, 12, cRose, "middle", fmt.Sprintf("se planta en n = %s", p05mil(ns[len(ns)-1])))
			// la lupa: the first steps, magnified until the treads show
			k := len(bs)
			if k > 24 {
				k = 24
			}
			lb, ln := bs[:k], ns[:k]
			l.Line(cxw+52, cyw-22, 258, 240, cGold, 1.5, "5 4")
			l.Panel(250, 196, 410, 150, cGold, fmt.Sprintf("CON LUPA — los primeros %d de los %d escalones", k, len(bs)))
			bmin, bmax := lb[0], lb[0]
			for _, v := range lb {
				bmin, bmax = math.Min(bmin, v), math.Max(bmax, v)
			}
			gx := func(n int64) float64 {
				return escala(float64(n), float64(ln[0]), float64(ln[k-1]), 268, 640)
			}
			gy := func(b float64) float64 { return escala(b, bmax, bmin, 238, 300) }
			l.Camino(p05escalera(lb, ln, gx, gy), cRose, 1.6, 0.95)
			l.Mono(268, 322, 10.5, cDim, "start", p05coma(fmt.Sprintf("b: %.6f → %.6f — se movió %.1e en %d vueltas", bs[0], bs[len(bs)-1], bs[0]-bs[len(bs)-1], len(bs))))
		}
		l.Mono(px+18, py+448, 12, cInk, "start", fmt.Sprintf("vueltas: %d      b final: %.6g", len(bs), bs[len(bs)-1]))
		l.Mono(px+18, py+470, 12, color, "start", fmt.Sprintf("n: %s → %s", p05mil(ns[0]), p05mil(ns[len(ns)-1])))
	}
	dibuja(46, 100, cRose, "SIN CIZALLA — la escalera se enrosca", sinB, sinN, false)
	dibuja(714, 100, cGreen, "CON CIZALLA — un tijeretazo y a puerto", conB, conN, true)

	l.Formula(46, 640, 1308, "b ∈ (¼,½] → (b−½, c+½): la suma NO cambia — j(j−1) es par ⟹ e^{−iπ j(j−1)} = 1   ·   espejo b > ½ → (1−b, 1−c) y se conjuga")
	l.Nota(46, 694, 1308, cGold, []string{
		"EN CRIOLLO: la máquina baja de un millón de términos a mil cambiando la escalera larga por otra más corta, una y otra vez. Pero hay una trampa:",
		"si la curvatura queda pegada a un medio, cada escalón acorta casi nada y la escalera se enrosca en un remolino que no llega nunca a puerto.",
		"La cizalla es un tijeretazo: le saca media vuelta a cada término. Media vuelta por el cuadrado de un entero da siempre vuelta entera, así que",
		"la cuenta da EXACTAMENTE lo mismo — y la curvatura, que estaba pegada al medio, cae al fondo. Mismo mar, mismo resultado, dos escalones.",
	})
	l.Mono(700, 828, 13, cGold, "middle", fmt.Sprintf("medido en esta corrida: %d vueltas sin cizalla contra %d con cizalla — %.0f veces menos escalones para el mismo número",
		len(sinB), len(conB), mejor))
	l.Pie("go run ./cmd/losplanos · lámina ⑤ del recorrido de las máquinas · riel 5 del tren · Todavía no.")
	l.Guardar("laminas/plano-05-la-cizalla.svg")
}

// p06nivel is one storey of the descent with its edge accounting.
type p06nivel struct {
	b, c     float64
	n        int64
	izq, der int64
	dientes  int64
}

// p06niveles walks cmd/circulo's cascadeDD skeleton in float64, recording at
// each storey the dual teeth crossed and the two Fresnel edge bands
// (margin = MFh/sqrt(2b), MFh capped at 24 by the adaptive horizon F148).
func p06niveles(b, c float64, n int64, nCom int64, MF float64) []p06nivel {
	var out []p06nivel
	for k := 0; k < 400; k++ {
		for {
			b, c = frac(b), frac(c)
			if b > 0.5 {
				b, c = 1-b, frac(1-c)
			}
			if b > 0.25 {
				b, c = b-0.5, frac(c+0.5)
				continue
			}
			break
		}
		if n <= nCom || b < 1e-9 {
			return out
		}
		s2b := math.Sqrt(2 * b)
		MFh := MF
		if need := 1 / (2 * math.Pi * (1e-7 * math.Sqrt(float64(n))) * s2b); need > MFh {
			MFh = math.Min(need, 24)
		}
		marg := MFh / s2b
		xLo, xHi := -0.5, float64(n)-0.5
		mLo := int64(math.Ceil(c + 2*b*(xLo-marg)))
		mHi := int64(math.Floor(c + 2*b*(xHi+marg)))
		mA := int64(math.Ceil(c + 2*b*(xLo+marg)))
		mB := int64(math.Floor(c + 2*b*(xHi-marg)))
		out = append(out, p06nivel{b, c, n, mA - mLo, mHi - mB, mHi - mLo + 1})
		if mB < mA {
			return out
		}
		b, c, n = 1/(4*b), (float64(mA)-c)/(2*b), mB-mA+1
	}
	return out
}

// p06cornu integrates F(u) = int_0^u e^{i pi tau^2} dtau by composite Simpson,
// returning the whole measured curve as (u, Re, Im).
func p06cornu(uMax float64, steps int) [][3]float64 {
	h := uMax / float64(steps)
	pts := make([][3]float64, 0, steps+1)
	var re, im float64
	pts = append(pts, [3]float64{0, 0, 0})
	for k := 0; k < steps; k++ {
		u := float64(k) * h
		s0, c0 := math.Sincos(math.Pi * u * u)
		sm, cm := math.Sincos(math.Pi * (u + h/2) * (u + h/2))
		s1, c1 := math.Sincos(math.Pi * (u + h) * (u + h))
		re += h / 6 * (c0 + 4*cm + c1)
		im += h / 6 * (s0 + 4*sm + s1)
		pts = append(pts, [3]float64{u + h, re, im})
	}
	return pts
}

// plano06 draws the dual window and the two thin edges that cost the luxury.
func plano06() {
	nv := p06niveles(0.199689, 0.437, 5000000, 1000, 8)
	var evals, dientes int64
	for _, v := range nv {
		evals += 2 * (v.izq + v.der)
		dientes += v.dientes
	}
	pct := 100 * float64(evals) / float64(dientes)
	// the Cornu curve, and its eye measured by averaging one whole winding turn
	cur := p06cornu(5.2, 52000)
	var er, ei float64
	var cnt int
	for _, p := range cur {
		if p[0] >= 5.0 {
			er, ei, cnt = er+p[1], ei+p[2], cnt+1
		}
	}
	er, ei = er/float64(cnt), ei/float64(cnt)
	diag := 2 * math.Hypot(er, ei)
	ang := math.Atan2(ei, er) * 180 / math.Pi
	fmt.Printf("SS plano06 · %d niveles · dientes duales %s · evaluaciones Fresnel %d = %.4f%% · orillas nivel 0: %d+%d · Cornu medido: diagonal %.6f a %.3f° (e^{iπ/4} = 1 a 45°), media cuerda %.6f\n",
		len(nv), p05mil(dientes), evals, pct, nv[0].izq, nv[0].der, diag, ang, diag/2)

	l := NewLienzo(1400, 950)
	l.Titulo(p05coma(fmt.Sprintf("⑥ LA VENTANA Y SUS ORILLAS — el %.3f%% del círculo", pct)),
		"adentro todos los dientes llevan el mismo sombrero de 45°; sólo los de la orilla piden la integral exacta")

	// ---- the window and its teeth ----
	l.Panel(46, 98, 838, 412, cBlue, "LA VENTANA DUAL — un palito por cada m, todos de la misma altura 1/√(2b)")
	wx0, wx1, ty, by := 100.0, 856.0, 192.0, 418.0
	l.Rect(wx0, 176, wx1-wx0, 256, "none", cGrid, 0.9)
	const nPal = 97
	for i := 0; i < nPal; i++ {
		sx := wx0 + (wx1-wx0)*float64(i)/(nPal-1)
		if i <= 7 || i >= nPal-8 { // the two edge bands
			l.Line(sx, by, sx, ty, cRose, 2.4, "")
			l.Curva(sx-6, ty+2, sx+6, ty+2, 7, cRose, 1.6, "")
			continue
		}
		l.Line(sx, by, sx, ty, cGreen, 1.5, "")
		l.Line(sx-7, ty+7, sx+7, ty-7, cGold, 2.2, "") // el sombrero de 45°
	}
	l.Line(wx0, by, wx1, by, cDim, 1.6, "")
	for _, lc := range []float64{128, 828} { // the golden magnifiers
		l.Circ(lc, 305, 54, "none", cGold, 3.2)
		l.Circ(lc, 305, 48, "none", cGold, 1.1)
		l.Line(lc+38, 343, lc+62, 367, cGold, 5, "")
	}
	l.Mono(wx0, 438, 11.5, cDim, "start", "x = −½")
	l.Mono(wx1, 438, 11.5, cDim, "end", "x = N − ½")
	l.Txt(150, 462, 12.5, cGold, "middle", "Fresnel exacto")
	l.Txt(806, 462, 12.5, cGold, "middle", "Fresnel exacto")
	l.Txt(478, 462, 13, cGreen, "middle", "A ESTOS LOS LLEVA LA ARMONÍA — cada uno aporta el mismo e^{iπ/4}")
	l.Mono(465, 490, 12, cInk, "middle", fmt.Sprintf("nivel 0 medido: %s dientes duales, y sólo %d+%d en las orillas (±%d alrededor de cada borde)",
		p05mil(nv[0].dientes), nv[0].izq, nv[0].der, nv[0].izq/2))

	// ---- the measured Cornu spiral ----
	l.Panel(900, 98, 454, 412, cGold, "EL ESPIRAL DE CORNU — integrado acá, punto por punto")
	mx := func(v float64) float64 { return escala(v, -0.58, 0.58, 995, 1259) }
	my := func(v float64) float64 { return escala(v, -0.58, 0.58, 452, 188) }
	l.Line(mx(-0.58), my(0), mx(0.58), my(0), cGrid, 1.2, "")
	l.Line(mx(0), my(-0.58), mx(0), my(0.58), cGrid, 1.2, "")
	var esp [][2]float64
	for i := len(cur) - 1; i >= 0; i -= 52 { // negative branch: F(-u) = -F(u)
		if cur[i][0] <= 4.6 {
			esp = append(esp, [2]float64{mx(-cur[i][1]), my(-cur[i][2])})
		}
	}
	for i := 0; i < len(cur); i += 52 {
		if cur[i][0] <= 4.6 {
			esp = append(esp, [2]float64{mx(cur[i][1]), my(cur[i][2])})
		}
	}
	l.Camino(esp, cBlue, 1.5, 0.95)
	l.Flecha(mx(-er), my(-ei), mx(er), my(ei), cGold, 3.2)
	const off = 0.075 // the half chord, drawn parallel so both stay visible
	l.Flecha(mx(off), my(-off), mx(er+off), my(ei-off), cRose, 3.2)
	l.Circ(mx(er), my(ei), 4.5, cGold, "", 0)
	l.Circ(mx(-er), my(-ei), 4.5, cGold, "", 0)
	l.Circ(mx(0), my(0), 4.5, cInk, "", 0)
	l.Txt(mx(0)-14, my(0)-8, 11.5, cInk, "end", "F(0) = 0")
	l.Mono(1127, 150, 12, cGold, "middle", p05coma(fmt.Sprintf("ojo a ojo = %.6f a %.2f°", diag, ang)))
	l.Txt(1127, 170, 12, cGold, "middle", "el diente de adentro: la diagonal entera")
	l.Txt(1127, 478, 12, cRose, "middle", p05coma(fmt.Sprintf("el diente justo en el borde: la MITAD (%.4f)", diag/2)))

	// ---- the ledger of the descent ----
	l.Panel(46, 524, 1308, 176, cGreen, "LA CUENTA DEL DESCENSO — dientes recorridos (barra) contra dientes de lujo (los topes rosados)")
	for i, v := range nv {
		y := 566 + float64(i)*20
		l.Mono(62, y+4, 11, cDim, "start", fmt.Sprintf("nivel %d   b = %s   n = %s", i, p05coma(fmt.Sprintf("%.9f", v.b)), p05mil(v.n)))
		w := escala(math.Log10(float64(v.dientes)), 2, math.Log10(float64(nv[0].dientes)), 10, 380)
		l.Rect(420, y-4, w, 9, cBlue, "none", 0.85)
		l.Rect(417, y-8, 6, 17, cRose, "none", 1)
		l.Rect(420+w-3, y-8, 6, 17, cRose, "none", 1)
		l.Mono(816, y+4, 11, cInk, "start", fmt.Sprintf("dientes %s · orillas %d+%d · lujo %d evaluaciones",
			p05mil(v.dientes), v.izq, v.der, 2*(v.izq+v.der)))
	}
	l.Mono(700, 690, 12.5, cGold, "middle", fmt.Sprintf("TOTAL: %s dientes duales pisados y %d evaluaciones exactas = %s%% del círculo",
		p05mil(dientes), evals, p05coma(fmt.Sprintf("%.4f", pct))))

	l.Formula(46, 712, 1308, "F(u) = ∫₀^u e^{iπτ²} dτ   ·   F(+∞) − F(−∞) = e^{iπ/4}   ·   medio extremo:  Σ_{a≤j≤b} f(j) = ∫_a^b f(x) dx + ½[f(a) + f(b)] + …")
	l.Nota(46, 762, 1308, cBlue, []string{
		"EN CRIOLLO: cuando la máquina da vuelta el círculo le queda una fila larguísima de dientes. Los del medio son todos iguales: cada uno aporta",
		"el mismo giro de 45 grados, y por eso se los lleva puestos una fórmula sola, sin mirarlos de a uno. Los únicos distintos son los quince de",
		"cada punta, donde la ventana se corta al medio y el diente aporta apenas la mitad. A esos —y sólo a esos— se les hace la integral exacta.",
		"Es la diferencia entre pesar el vagón entero y pesar los dos tornillos de las puntas: por eso el lujo sale una milésima de lo que parecía.",
	})
	l.Pie("go run ./cmd/losplanos · lámina ⑥ del recorrido de las máquinas · riel 4 del tren, la forma armónica F145 · Todavía no.")
	l.Guardar("laminas/plano-06-las-orillas.svg")
}

// p07agua is one water with its measured block budget.
type p07agua struct {
	exp, t, k0, L, eta float64
	piezas             int64
	parte              bool
}

// plano07 draws the block law: the beam grows without bound, its sag never does.
func plano07() {
	var aguas []p07agua
	for i := 0; i < 8; i++ {
		e := 21 + float64(i)*27/7
		t := math.Pow(10, e)
		k0 := math.Sqrt(t / (2 * math.Pi))
		L := bandL(t, k0)
		et := eta(t, k0, L)
		aguas = append(aguas, p07agua{e, t, k0, L, et,
			int64(math.Ceil(math.Cbrt(math.Abs(et) / 0.04))), math.Abs(et) > 0.05 && L > 2000})
	}
	eMin, eMax := aguas[0].eta, aguas[0].eta
	for _, a := range aguas {
		eMin, eMax = math.Min(eMin, a.eta), math.Max(eMax, a.eta)
	}
	et, pz := aguas[0].eta, aguas[7].piezas
	razL := aguas[7].L / aguas[0].L
	fmt.Printf("SS plano07 · 8 aguas 10^21..10^48 · L de %.2f a %.5g (×%.0f) · eta = %.12f en TODAS (dispersión %.1e) · piezas = ceil((eta/0.04)^(1/3)) = %d · en 10^21 L=%.0f<2000 → 1 sola pieza\n",
		aguas[0].L, aguas[7].L, razL, et, eMax-eMin, pz, aguas[0].L)

	l := NewLienzo(1400, 950)
	l.Titulo("⑦ LA LEY DEL BLOQUE — la viga que se comba y se parte en dos",
		fmt.Sprintf("el bloque se estira ×%s entre 10²¹ y 10⁴⁸, y la comba que deja es siempre la misma: %s radianes",
			p05mil(int64(math.Round(razL))), p05coma(fmt.Sprintf("%.2f", et))))

	// ---- the beam and its sag ----
	l.Panel(46, 98, 1308, 392, cGold, "LA VIGA Y SU COMBA — la fase cúbica que el primer orden no explica")
	bx0, bx1, byy := 150.0, 1180.0, 168.0
	prof := func(f float64) float64 { return escala(f, 0, 0.20, 0, 150) } // radians -> pixels of sag
	comba := func(x0, x1, top, e float64, col string) {
		l.Line(x0, top, x1, top, cDim, 1.4, "6 5")
		var c [][2]float64
		for k := 0; k <= 120; k++ {
			u := float64(k) / 120
			c = append(c, [2]float64{x0 + (x1-x0)*u, top + prof(e*u*u*u)})
			if k%12 == 6 { // the hangers, so it reads as a beam under load
				l.Line(x0+(x1-x0)*u, top, x0+(x1-x0)*u, top+prof(e*u*u*u), cGrid, 1.2, "")
			}
		}
		l.Camino(c, col, 3, 0.98)
		l.Line(x0-6, top-8, x0+6, top+8, cDim, 2, "")
		l.Line(x0+6, top-8, x0-6, top+8, cDim, 2, "")
	}
	l.Txt(bx0-14, byy-10, 12, cDim, "end", "la viga ideal")
	comba(bx0, bx1, byy, et, cRose)
	xc := bx0 + (bx1-bx0)*math.Cbrt(0.05/et)
	l.Rect(xc, byy+prof(0.05), bx1-xc, prof(et)-prof(0.05), cRose, "none", 0.16)
	l.Line(bx0, byy+prof(0.05), bx1, byy+prof(0.05), cGold, 1.8, "8 6")
	l.Txt(bx0+8, byy+prof(0.05)-8, 12, cGold, "start", "presupuesto: 0,05 rad")
	l.Mono(bx1+8, byy+prof(et)+5, 12.5, cRose, "start", p05coma(fmt.Sprintf("η = %.6f", et)))
	l.Txt(xc+12, byy+prof(et)+20, 12, cRose, "start", "de acá para allá la viga ya se pasó de comba")
	l.Mono(bx0, byy+prof(et)+42, 12, cDim, "start", p05coma(fmt.Sprintf("una sola pieza de largo L — comba η = t·L³/(3k₀³) = %.6f rad", et)))
	l.Flecha(760, 306, 760, 344, cGold, 3)
	l.Txt(776, 332, 13, cGold, "start", fmt.Sprintf("se parte en %d", pz))

	// the broken beam: each half sags eta / pieces^3
	e2, top2 := et/math.Pow(float64(pz), 3), 372.0
	comba(bx0, 655, top2, e2, cGreen)
	comba(675, bx1, top2, e2, cGreen)
	for k := 0; k < 5; k++ { // the dd stitch
		l.Line(657, top2-12+float64(k)*8, 673, top2-6+float64(k)*8, cGold, 2, "")
	}
	l.Txt(665, top2-24, 12, cGold, "middle", "costura dd")
	l.Mono(bx1+8, top2+prof(e2)+5, 12, cGreen, "start", p05coma(fmt.Sprintf("η/%d³ = %.6f", pz, e2)))
	l.Txt(bx0, top2+prof(e2)+28, 12, cGreen, "start", "cada mitad, casi recta: bajo el presupuesto, y el tren se la come de un saque")
	l.Mono(665, 458, 11.5, cGold, "middle", "la costura gira −T·[ln(k₀+off) − ln k₀], la fase del arranque de cada mitad, llevada en doble-doble")

	// ---- the ruler of waters ----
	l.Panel(46, 504, 1308, 196, cBlue, "EL DURMIENTE EN CADA AGUA — ocho mares, ocho vigas, una sola comba")
	lx := func(e float64) float64 { return escala(e, 21, 48, 142, 1266) }
	for _, a := range aguas { // sleepers first, rails on top
		w := escala(math.Log10(a.L), math.Log10(aguas[0].L), math.Log10(aguas[7].L), 26, 120)
		col := cGrid
		if !a.parte {
			col = "#4a2740"
		}
		l.Rect(lx(a.exp)-w/2, 548, w, 22, col, cDim, 0.95)
		l.Mono(lx(a.exp), 590, 10.5, cInk, "middle", fmt.Sprintf("L=%.0f", a.L))
		l.Mono(lx(a.exp), 608, 11, cDim, "middle", fmt.Sprintf("t=10^%.0f", a.exp))
	}
	l.Line(120, 554, 1302, 554, cGold, 2, "")
	l.Line(120, 564, 1302, 564, cGold, 2, "")
	l.Txt(lx(21), 632, 11, cRose, "middle", fmt.Sprintf("L = %.0f < 2000: 1 pieza", aguas[0].L))
	l.Txt(1302, 632, 11, cDim, "end", "durmiente dibujado en escala logarítmica (el número de al lado es el real)")
	l.Line(120, 664, 1302, 664, cGold, 2.2, "")
	for _, a := range aguas {
		l.Circ(lx(a.exp), 664, 5.5, cGold, "", 0)
	}
	l.Txt(120, 654, 12, cGold, "start", "la comba η")
	l.Mono(1302, 654, 12, cGold, "end", p05coma(fmt.Sprintf("η = %.6f rad — PLANA, agua por agua", et)))
	l.Mono(700, 688, 12, cGreen, "middle", fmt.Sprintf("L se multiplica por %s · η se multiplica por %s (dispersión medida entre las 8 aguas: %s)",
		p05mil(int64(math.Round(razL))), p05coma(fmt.Sprintf("%.6f", aguas[7].eta/aguas[0].eta)), p05coma(fmt.Sprintf("%.1e", eMax-eMin))))

	l.Formula(46, 712, 1308, p05coma(fmt.Sprintf("L = k₀·(0,45/t)^{1/3}  ⟹  η = t·L³/(3k₀³) = 0,45/3 = %.2f rad, en TODA agua   ·   piezas = ⌈(η/0,04)^{1/3}⌉ = %d", et, pz)))
	l.Nota(46, 762, 1308, cGreen, []string{
		"EN CRIOLLO: el tren no puede sumar el mar entero de una, así que lo corta en bloques. ¿De qué largo? El largo se elige para que la viga se combe",
		"siempre lo mismo: no importa si el agua tiene veintiún ceros o cuarenta y ocho, el bloque se estira treinta mil veces pero la panza que hace es",
		"idéntica. Esa es la ley: la regla no mide metros, mide comba. Y como esa comba pasa el presupuesto, la viga se parte al medio; cada mitad comba",
		"ocho veces menos y queda derecha. La costura que las vuelve a unir se lleva en doble-doble, para que no se note el corte por ningún lado.",
	})
	l.Pie("go run ./cmd/losplanos · lámina ⑦ del recorrido de las máquinas · riel 7 del tren · Todavía no.")
	l.Guardar("laminas/plano-07-la-ley-del-bloque.svg")
}
