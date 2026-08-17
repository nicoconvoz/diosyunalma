// planos_d.go - plates 11..13 of the museum tour.
//
//	11: the two machines share one law and differ only in one constant.
//	12: the quartic hull whose needles are re-nailed every 4096 clicks.
//	13: one Fresnel fold per hunting window - only the carrier turns.
//
// Every number on these sheets is measured in this run: the block budgets come
// from the bench (bandL / blockFracL / eta), and plates 12-13 carry their own
// arbitrary-precision reference, standing in for the ships' double-double.
package main

import (
	"fmt"
	"math"
	"math/big"
)

// ---------------------------------------------------------------------------
// arbitrary-precision bench (shared by plates 12 and 13)
// ---------------------------------------------------------------------------

const p12prec = 240

var p12pi = p12str("3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211707")

func p12str(s string) *big.Float {
	f, _, _ := big.ParseFloat(s, 10, p12prec, big.ToNearestEven)
	return f
}
func p12f(v float64) *big.Float         { return new(big.Float).SetPrec(p12prec).SetFloat64(v) }
func p12mul(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(p12prec).Mul(a, b) }
func p12add(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(p12prec).Add(a, b) }
func p12div(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(p12prec).Quo(a, b) }
func p12sub(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(p12prec).Sub(a, b) }

// p12mod is the floored remainder x mod m.
func p12mod(x, m *big.Float) *big.Float {
	qi, _ := p12div(x, m).Int(nil)
	r := p12sub(x, p12mul(m, new(big.Float).SetPrec(p12prec).SetInt(qi)))
	if r.Sign() < 0 {
		r = p12add(r, m)
	}
	return r
}

// p12m2p mirrors the ship's mod2pi: reduce, round to float64, and if the
// rounding lands on 2pi itself, fold it back to zero exactly as the ship does.
func p12m2p(x *big.Float) float64 {
	v, _ := p12mod(x, p12mul(p12pi, p12f(2))).Float64()
	if v >= 2*math.Pi || v < 0 {
		v = math.Mod(v, 2*math.Pi)
		if v < 0 {
			v += 2 * math.Pi
		}
	}
	return v
}

func p12frac1(x *big.Float) float64 {
	v, _ := p12mod(x, p12f(1)).Float64()
	return frac(v)
}

// p12coef carves the facet's four screws (cmd/starship build): the Taylor
// coefficients of -t*ln(k0+j) in j, up to the quartic.
func p12coef(t, k0 float64) [4]*big.Float {
	T, iv := p12f(t), p12div(p12f(1), p12f(k0))
	i2 := p12mul(iv, iv)
	return [4]*big.Float{
		p12mul(T, iv),
		p12div(p12mul(T, i2), p12f(-2)),
		p12div(p12mul(T, p12mul(i2, iv)), p12f(3)),
		p12div(p12mul(T, p12mul(i2, i2)), p12f(-4)),
	}
}

func p12phi(c [4]*big.Float, j int64) float64 {
	jb := new(big.Float).SetPrec(p12prec).SetInt64(j)
	p, s := p12f(1), p12f(0)
	for i := 0; i < 4; i++ {
		p = p12mul(p, jb)
		s = p12add(s, p12mul(c[i], p))
	}
	return p12m2p(s)
}

// p12refresh is the ship's facet.refresh: the hand that re-nails phase and the
// three needles at step j, straight from the high-precision chain.
func p12refresh(c [4]*big.Float, j int64) (ph, d1, d2, d3 float64) {
	jf := float64(j)
	j2, j3 := jf*jf, jf*jf*jf
	ph = p12phi(c, j)
	d1 = p12m2p(p12add(c[0], p12add(p12mul(c[1], p12f(2*jf+1)),
		p12add(p12mul(c[2], p12f(3*j2+3*jf+1)), p12mul(c[3], p12f(4*j3+6*j2+4*jf+1))))))
	d2 = p12m2p(p12add(p12mul(c[1], p12f(2)),
		p12add(p12mul(c[2], p12f(6*jf+6)), p12mul(c[3], p12f(12*j2+24*jf+14)))))
	d3 = p12m2p(p12add(p12mul(c[2], p12f(6)), p12mul(c[3], p12f(24*jf+36))))
	return
}

// p12paseo walks the forward-difference chain in float64 exactly as the facet
// tier does. cada>0 re-nails the needles every `cada` steps; cada=0 never does.
func p12paseo(c [4]*big.Float, n, cada, paso int64, ref []float64) (errs []float64, peor, ultimo float64, perdida int64) {
	const tp = 2 * math.Pi
	ph, d1, d2, d3 := p12refresh(c, 0)
	d4 := p12m2p(p12mul(c[3], p12f(24)))
	perdida = -1
	k := 0
	for j := int64(0); j <= n; j++ {
		if cada > 0 && j > 0 && j%cada == 0 {
			ph, d1, d2, d3 = p12refresh(c, j)
		}
		if j%paso == 0 && k < len(ref) {
			e := math.Abs(math.Remainder(ph-ref[k], tp))
			errs = append(errs, e)
			if e > peor {
				peor = e
			}
			if perdida < 0 && e > 1e-3 {
				perdida = j
			}
			ultimo, k = e, k+1
		}
		ph += d1
		if ph >= tp {
			ph -= tp
		}
		d1 += d2
		if d1 >= tp {
			d1 -= tp
		}
		d2 += d3
		if d2 >= tp {
			d2 -= tp
		}
		d3 += d4
		if d3 >= tp {
			d3 -= tp
		}
	}
	return
}

// ---------------------------------------------------------------------------
// drawn parts
// ---------------------------------------------------------------------------

// p11loco draws a schematic locomotive standing on its rail at (x, y).
func p11loco(l *Lienzo, x, y float64, col string) {
	l.Rect(x+2, y-10, 186, 12, "none", col, 0.9)
	l.Rect(x+8, y-70, 50, 60, "none", col, 0.9)
	l.Rect(x+58, y-52, 116, 42, "none", col, 0.9)
	l.Rect(x+150, y-78, 20, 26, "none", col, 0.9)
	l.Rect(x+172, y-38, 14, 28, "none", col, 0.9)
	l.Circ(x+96, y-52, 11, "none", col, 0.9)
	l.Rect(x+16, y-62, 34, 20, "none", cDim, 0.7)
	for _, w := range [][3]float64{{34, 6, 20}, {88, 2, 14}, {126, 2, 14}, {160, 2, 14}} {
		l.Circ(x+w[0], y+w[1], w[2], "none", col, 1.6)
		l.Circ(x+w[0], y+w[1], 3.2, col, "", 0)
	}
	for i := 0; i < 3; i++ {
		f := float64(i)
		l.Circ(x+164+f*18, y-92-f*14, 7+f*2.5, "none", cDim, 0.8)
	}
	l.Line(x-14, y+22, x+204, y+22, cGrid, 3, "")
}

// p11auto draws a schematic gullwing wedge (the DeLorean) at (x, y).
func p11auto(l *Lienzo, x, y float64, col string) {
	body := [][2]float64{{0, -12}, {12, -26}, {50, -32}, {76, -54}, {126, -56},
		{152, -34}, {178, -28}, {182, -14}, {176, -4}, {6, -4}, {0, -12}}
	pts := make([][2]float64, len(body))
	for i, p := range body {
		pts[i] = [2]float64{x + p[0], y + p[1]}
	}
	l.Camino(pts, col, 1.9, 1)
	l.Camino([][2]float64{{x + 84, y - 52}, {x + 122, y - 54}, {x + 118, y - 38}, {x + 88, y - 36}, {x + 84, y - 52}}, cDim, 1.2, 0.9)
	l.Camino([][2]float64{{x + 100, y - 56}, {x + 108, y - 92}, {x + 148, y - 78}, {x + 138, y - 52}}, col, 1.4, 0.75)
	for _, w := range [][2]float64{{44, -4}, {142, -4}} {
		l.Circ(x+w[0], y+w[1], 14, "none", col, 1.7)
		l.Circ(x+w[0], y+w[1], 3.2, col, "", 0)
	}
	l.Line(x-14, y+16, x+204, y+16, cGrid, 3, "")
}

// p11engranaje draws a toothed wheel of the given radius.
func p11engranaje(l *Lienzo, cx, cy, r float64, dientes int, col string) {
	l.Circ(cx, cy, r, "none", col, 1.8)
	l.Circ(cx, cy, r*0.32, "none", col, 1.2)
	for i := 0; i < dientes; i++ {
		a := 2 * math.Pi * float64(i) / float64(dientes)
		s, c := math.Sincos(a)
		l.Line(cx+c*r, cy+s*r, cx+c*(r+6.5), cy+s*(r+6.5), col, 3.2, "")
	}
}

// p11viga draws a beam of the given pixel length sagging by etaR radians.
func p11viga(l *Lienzo, x, y, largo, etaR, pxrad float64, col string, sw float64) {
	l.Line(x, y, x+largo, y, cGrid, 1.2, "5 5")
	pts := make([][2]float64, 0, 61)
	for i := 0; i <= 60; i++ {
		u := float64(i) / 60
		pts = append(pts, [2]float64{x + u*largo, y + pxrad*etaR*u*u*u})
	}
	l.Camino(pts, col, sw, 1)
	l.Line(x, y-9, x, y+9, col, 2, "")
	l.Line(x+largo, y-9, x+largo, y+pxrad*etaR+9, col, 2, "")
}

// p12dial draws a gauge whose needle sits at the measured angle (radians).
func p12dial(l *Lienzo, cx, cy, r, ang float64, col, nombre, valor string) {
	l.Circ(cx, cy, r, cPanel, col, 1.7)
	for i := 0; i < 12; i++ {
		a := 2 * math.Pi * float64(i) / 12
		s, c := math.Sincos(a)
		l.Line(cx+c*(r-7), cy+s*(r-7), cx+c*r, cy+s*r, cGrid, 1.6, "")
	}
	s, c := math.Sincos(ang - math.Pi/2)
	l.Line(cx, cy, cx+c*(r-10), cy+s*(r-10), col, 2.6, "")
	l.Circ(cx, cy, 4, col, "", 0)
	l.Txt(cx, cy+r+20, 14, col, "middle", nombre)
	l.Mono(cx, cy+r+36, 11.5, cDim, "middle", valor)
}

// p12mano draws the schematic hand that comes down to re-nail a needle.
func p12mano(l *Lienzo, x, y float64, col string) {
	l.Rect(x-26, y-46, 52, 34, "none", col, 0.95)
	for i := 0; i < 4; i++ {
		f := -18 + float64(i)*12
		l.Rect(x+f, y-12, 8, 16, "none", col, 0.8)
	}
	l.Rect(x-38, y-40, 12, 20, "none", col, 0.8)
}

// ---------------------------------------------------------------------------

func plano11() {
	aguas := []float64{1e12, 1e18, 1e24, 1e30, 1e36, 1e42, 1e48}
	type m11 struct{ t, k0, ld, lt, r, ed, et float64 }
	var ms []m11
	rLo, rHi := math.Inf(1), math.Inf(-1)
	edLo, edHi := math.Inf(1), math.Inf(-1)
	etLo, etHi := math.Inf(1), math.Inf(-1)
	for _, t := range aguas {
		k0 := math.Sqrt(t/(2*math.Pi)) / 2
		ld, lt := blockFracL(t, k0), bandL(t, k0)
		x := m11{t, k0, ld, lt, lt / ld, eta(t, k0, ld), eta(t, k0, lt)}
		ms = append(ms, x)
		rLo, rHi = math.Min(rLo, x.r), math.Max(rHi, x.r)
		edLo, edHi = math.Min(edLo, x.ed), math.Max(edHi, x.ed)
		etLo, etHi = math.Min(etLo, x.et), math.Max(etHi, x.et)
	}
	c50 := math.Cbrt(50)
	ej := ms[3]
	fmt.Printf("SS plano11 aguas=%d cociente=[%.9f..%.9f] cbrt50=%.9f etaDeLorean=[%.12f..%.12f] etaTren=[%.12f..%.12f] ej t=%.0e k0=%.4e Ld=%.4e Lt=%.4e\n",
		len(ms), rLo, rHi, c50, edLo, edHi, etLo, etHi, ej.t, ej.k0, ej.ld, ej.lt)

	l := NewLienzo(1400, 950)
	l.Titulo("⑪ LA MISMA LEY, DOS PRESUPUESTOS — el tren y el DeLorean comparten geometría",
		"una sola regla para el largo del bloque; cada máquina elige cuánta panza le tolera a la viga")

	// the shared law, hanging over both machines
	l.Rect(468, 98, 464, 46, cPanel, cGold, 0.55)
	l.Mono(700, 128, 20, cGold, "middle", "L = k · (C / t)^(1/3)")
	l.Txt(700, 158, 12.5, cDim, "middle", "LA MISMA LEY")
	l.Flecha(560, 146, 372, 184, cGold, 1.6)
	l.Flecha(840, 146, 1030, 184, cGold, 1.6)

	l.Panel(40, 176, 600, 206, cBlue, "EL DELOREAN — presupuesto C = 0.009")
	p11auto(l, 62, 312, cBlue)
	l.Panel(760, 176, 600, 206, cRose, "EL TREN — presupuesto C = 0.45")
	p11loco(l, 782, 312, cRose)
	// the same two gears carved in both boxes
	for i, gx := range []float64{330, 1050} {
		col := cBlue
		if i == 1 {
			col = cRose
		}
		l.Rect(gx-16, 214, 292, 158, cBg, cGrid, 0.75)
		p11engranaje(l, gx+56, 280, 38, 12, col)
		p11engranaje(l, gx+166, 280, 38, 12, col)
		l.Txt(gx+56, 344, 12.5, cDim, "middle", "espejo  b → 1−b")
		l.Txt(gx+166, 344, 12.5, cDim, "middle", "cizalla b → b−½")
		l.Txt(gx+111, 364, 11.5, cGold, "middle", "las mismas dos ruedas")
	}

	l.Panel(40, 396, 940, 304, cGold, "LAS DOS VIGAS A ESCALA — misma agua, mismo k")
	const px = 250.0
	ldPx := 150.0
	ltPx := ldPx * ej.r
	p11viga(l, 96, 468, ldPx, ej.ed, px, cBlue, 3)
	l.Txt(96, 446, 13, cBlue, "start", "DeLorean")
	l.Mono(96+ldPx+14, 472, 12, cDim, "start", fmt.Sprintf("η = %.4f rad", ej.ed))
	p11viga(l, 96, 580, ltPx, ej.et, px, cRose, 3)
	l.Txt(96, 558, 13, cRose, "start", "tren")
	l.Mono(96+ltPx+14, 584, 12, cDim, "start", fmt.Sprintf("η = %.4f rad", ej.et))
	l.Line(96, 648, 96+ltPx, 648, cGold, 1.2, "4 4")
	l.Line(96, 642, 96, 654, cGold, 1.6, "")
	l.Line(96+ltPx, 642, 96+ltPx, 654, cGold, 1.6, "")
	l.Txt(96+ltPx/2, 668, 12.5, cGold, "middle", fmt.Sprintf("el bloque del tren es %.4f veces más largo", ej.r))
	// magnifier over the DeLorean's beam
	l.Circ(866, 500, 66, cBg, cGold, 1.6)
	aum := ej.et / ej.ed
	pts := make([][2]float64, 0, 41)
	for i := 0; i <= 40; i++ {
		u := float64(i) / 40
		pts = append(pts, [2]float64{806 + u*120, 480 + px*ej.ed*aum*u*u*u})
	}
	l.Line(806, 480, 926, 480, cGrid, 1.1, "4 4")
	l.Camino(pts, cBlue, 2.4, 1)
	l.Txt(866, 552, 12.5, cGold, "middle", fmt.Sprintf("×%.0f", aum))

	l.Panel(1000, 396, 360, 304, cGreen, "EL COCIENTE EN SIETE AGUAS")
	l.Ejes(1042, 452, 288, 190, "", "L_tren / L_DeLorean")
	l.Line(1042, escala(c50, 3.60, 3.76, 642, 452), 1330, escala(c50, 3.60, 3.76, 642, 452), cGold, 1.6, "6 5")
	for _, x := range ms {
		cx := escala(math.Log10(x.t), 12, 48, 1050, 1322)
		cy := escala(x.r, 3.60, 3.76, 642, 452)
		l.Circ(cx, cy, 5.5, cGreen, "", 0)
		l.Mono(cx, 662, 10, cDim, "middle", fmt.Sprintf("%.0f", math.Log10(x.t)))
	}
	l.Mono(1326, escala(c50, 3.60, 3.76, 642, 452)-13, 11, cGold, "end", fmt.Sprintf("∛50 = %.6f", c50))
	l.Txt(1186, 682, 11.5, cDim, "middle", "log₁₀ t — siete aguas, un solo cociente")

	l.Formula(40, 712, 1320, fmt.Sprintf("L = k·(C/t)^(1/3)  ⇒  η = t·L³/(3k³) = C/3   |   0.009 → %.4f rad ;  0.45 → %.4f rad ;  cociente = ∛50 = %.6f", ej.ed, ej.et, ej.r))
	l.Nota(40, 760, 1320, cGold, []string{
		"EN CRIOLLO: las dos máquinas cortan el mar en bloques y usan la MISMA regla para decidir el largo de cada bloque.",
		"Lo único distinto es cuánta panza le permiten a la viga: el DeLorean tolera " + fmt.Sprintf("%.4f", ej.ed) + " radianes de curva sobrante, el tren " + fmt.Sprintf("%.4f", ej.et) + ".",
		fmt.Sprintf("Cincuenta veces más tolerancia le da al tren bloques %.4f veces más largos — la raíz cúbica de cincuenta, en cualquier agua.", ej.r),
		"Y adentro de las dos cajas están las mismas dos ruedas dentadas: el espejo (b→1−b) y la cizalla (b→b−½).",
	})
	l.Pie("Una sola geometría; dos presupuestos de curvatura. El resto es la misma máquina.")
	l.Guardar("laminas/plano-11-la-misma-ley.svg")
}

func plano12() {
	t := 1e24
	ub := math.Pow(0.1/t, 0.2)
	k0 := 6.4e10
	n := int64(ub * k0)
	const cada, paso = 4096, 500
	c := p12coef(t, k0)
	ref := make([]float64, 0, n/paso+1)
	for j := int64(0); j <= n; j += paso {
		ref = append(ref, p12phi(c, j))
	}
	d4 := p12m2p(p12mul(c[3], p12f(24)))
	eLib, peorL, ultL, perdL := p12paseo(c, n, 0, paso, ref)
	eCla, peorC, ultC, _ := p12paseo(c, n, cada, paso, ref)
	_, dd1, dd2, dd3 := p12refresh(c, n/2)
	cv := make([]float64, 4)
	for i := 0; i < 4; i++ {
		cv[i], _ = c[i].Float64()
	}
	fmt.Printf("SS plano12 t=%.0e k0=%.4e n=%d cada=%d cD=%.4e d4=%.3e libre[peor=%.4f rad, final=%.4f rad, >1e-3 en j=%d] clavado[peor=%.3e rad, final=%.3e rad] cuartica_soltada=%.4e rad\n",
		t, k0, n, cada, cv[3], d4, peorL, ultL, perdL, peorC, ultC, math.Abs(cv[3])*math.Pow(float64(cada), 4))

	l := NewLienzo(1400, 950)
	l.Titulo("⑫ EL CASCO QUE SE RE-ANCLA — la faceta cuártica",
		fmt.Sprintf("t = %.0e, k₀ = %.3g, faceta de %d pasos: la mano baja cada %d clics y vuelve a clavar las agujas", t, k0, n, cada))

	l.Panel(40, 96, 640, 604, cBlue, "LA FACETA: una regla curva con cuatro tornillos")
	// the curved ruler IS the facet's measured sag: the phase left over once
	// the straight line cA*j is taken away, normalized to 70 px.
	nf := float64(n)
	comba := func(u float64) float64 {
		j := u * nf
		return cv[1]*j*j + cv[2]*j*j*j + cv[3]*j*j*j*j
	}
	cMax := math.Abs(comba(1))
	rul := make([][2]float64, 0, 121)
	for i := 0; i <= 120; i++ {
		u := float64(i) / 120
		rul = append(rul, [2]float64{80 + u*560, 168 + 70*math.Abs(comba(u))/cMax})
	}
	l.Camino(rul, cBlue, 3.4, 1)
	l.Line(80, 168, 640, 168, cGrid, 1.2, "5 5")
	nom := []string{"cA", "cB", "cC", "cD"}
	for i := 0; i < 4; i++ {
		u := 0.16 + float64(i)*0.24
		sx, sy := 80+u*560, 168+70*math.Abs(comba(u))/cMax
		l.Circ(sx, sy, 9, cBg, cGold, 2)
		l.Line(sx-6, sy-6, sx+6, sy+6, cGold, 1.6, "")
		l.Line(sx-6, sy+6, sx+6, sy-6, cGold, 1.6, "")
		l.Txt(sx, sy-18, 14, cGold, "middle", nom[i])
		l.Mono(sx, sy+26, 10.5, cDim, "middle", fmt.Sprintf("%.3e", cv[i]))
	}
	l.Txt(360, 282, 12.5, cDim, "middle", "los cuatro tornillos de la fase: φ(j) = cA·j + cB·j² + cC·j³ + cD·j⁴")

	// the click track and the descending hand
	l.Line(80, 410, 640, 410, cDim, 2, "")
	marcas := int(n / cada)
	for i := 0; i <= marcas; i++ {
		x := 80 + 560*float64(i)/float64(marcas)
		l.Line(x, 404, x, 416, cGrid, 1.2, "")
	}
	for i := 0; i <= marcas; i += 16 {
		x := 80 + 560*float64(i)/float64(marcas)
		l.Line(x, 398, x, 422, cGreen, 1.8, "")
	}
	p12mano(l, 432, 366, cGreen)
	l.Flecha(432, 372, 432, 400, cGreen, 2)
	l.Txt(432, 306, 13, cGreen, "middle", fmt.Sprintf("la mano baja %d veces", marcas))
	l.Mono(80, 442, 11, cDim, "start", "j = 0")
	l.Mono(640, 442, 11, cDim, "end", fmt.Sprintf("j = %d", n))
	l.Txt(360, 462, 12, cGold, "middle", fmt.Sprintf("cada %d clics las agujas se re-clavan en el valor exacto (cadena de doble-doble)", cada))

	// the four needles, at their measured mid-walk values
	muerta := d4 == 0
	for i, v := range []float64{dd1, dd2, dd3, d4} {
		col := cGreen
		et := fmt.Sprintf("%.6f rad", v)
		if i == 3 && muerta {
			col, et = cRose, "0 — bajo la resolución"
		}
		p12dial(l, 145+float64(i)*140, 552, 48, v, col, []string{"d1", "d2", "d3", "d4"}[i], et)
	}
	pie12 := "las cuatro agujas avanzan solas; la mano las vuelve a clavar en el valor exacto"
	if muerta {
		pie12 = "el cuarto tornillo es más fino que la aguja: sin la mano, la cuártica se pierde entera"
	}
	l.Txt(360, 668, 12, cDim, "middle", pie12)

	// the measured error curves
	l.Panel(700, 96, 660, 604, cRose, "EL ERROR MEDIDO — con la mano y sin la mano")
	px0, px1, py0, py1 := 776.0, 1336.0, 150.0, 616.0
	l.Ejes(px0, py0, px1-px0, py1-py0, "j (pasos de la faceta)", "log₁₀ |error de fase| (rad)")
	for _, g := range []float64{-16, -12, -8, -4, 0} {
		y := escala(g, -16, 1, py1, py0)
		l.Line(px0, y, px1, y, cGrid, 1, "3 6")
		l.Mono(px0-8, y+4, 10.5, cDim, "end", fmt.Sprintf("%.0f", g))
	}
	yPi := escala(math.Log10(math.Pi), -16, 1, py1, py0)
	l.Line(px0, yPi, px1, yPi, cGold, 1.4, "7 5")
	l.Mono(px1-4, yPi-8, 11, cGold, "end", "π — la fase está perdida")
	curva := func(e []float64, col string, sw float64) {
		pts := make([][2]float64, 0, len(e))
		for i, v := range e {
			if v < 1e-16 {
				v = 1e-16
			}
			pts = append(pts, [2]float64{
				escala(float64(i)*paso, 0, float64(n), px0, px1),
				escala(math.Log10(v), -16, 1, py1, py0)})
		}
		l.Camino(pts, col, sw, 0.95)
	}
	curva(eLib, cRose, 2.2)
	curva(eCla, cGreen, 1.5)
	if perdL > 0 {
		xp := escala(float64(perdL), 0, float64(n), px0, px1)
		l.Line(xp, py0, xp, py1, cRose, 1.2, "4 5")
		l.Mono(xp+6, py0+16, 10.5, cRose, "start", fmt.Sprintf("j=%d: ya erró 10⁻³ rad", perdL))
	}
	l.Circ(800, 646, 6, cRose, "", 0)
	l.Mono(814, 650, 11.5, cInk, "start", fmt.Sprintf("sin re-anclaje: peor %.3f rad", peorL))
	l.Circ(1076, 646, 6, cGreen, "", 0)
	l.Mono(1090, 650, 11.5, cInk, "start", fmt.Sprintf("cada %d pasos: peor %.2e rad", cada, peorC))
	l.Txt(1068, 676, 11.5, cDim, "middle", "el diente de sierra verde es cada visita de la mano")

	l.Formula(40, 712, 1320, fmt.Sprintf("φ(j) = cA·j + cB·j² + cC·j³ + cD·j⁴  —  diferencias constantes d1..d4  |  re-anclaje en doble-doble cada %d pasos", cada))
	nota2 := fmt.Sprintf("El cuarto tornillo (cD = %.2e) le llega a la aguja d4 como %.3e rad por paso.", cv[3], d4)
	if muerta {
		nota2 = fmt.Sprintf("Pero el cuarto tornillo (cD = %.2e) es tan fino que la aguja d4 lo lee como CERO: sola, la cadena se va de curso.", cv[3])
	}
	nota3 := fmt.Sprintf("Sin la mano el error termina en %.3e rad al cabo del recorrido.", ultL)
	if perdL > 0 {
		nota3 = fmt.Sprintf("Sin la mano, a los %d pasos el error ya pasó la milésima de radián y termina en %.2f rad: la fase quedó perdida.", perdL, ultL)
	}
	l.Nota(40, 760, 1320, cGold, []string{
		"EN CRIOLLO: el casco avanza sumando cuatro agujas, una tras otra, sin recalcular nada: es baratísimo y por eso la nave vuela.",
		nota2,
		nota3,
		fmt.Sprintf("Con la mano bajando cada %d clics, el error nunca supera %.1e rad: se re-clava antes de que la deriva se note.", cada, peorC),
	})
	l.Pie("El truco no es calcular mejor: es volver a clavar la aguja antes de que la deriva crezca.")
	l.Guardar("laminas/plano-12-el-casco.svg")
}

func plano13() {
	t := 1e30
	nTop := int64(math.Sqrt(t / (2 * math.Pi)))
	cFrac := blockFracL(t, 1)
	foldFrom := int64(math.Ceil(256 / cFrac))
	spacing := espaciado(t)
	span := 5 * spacing
	h := 2 * math.Pi / math.Log(float64(nTop)) / 3
	guard := 8 * h
	S := int((span+2*guard)/h) + 2

	// pick a real fold-tier block whose chirp shows a handful of loops
	k0, L, bb := 2.5e13, 0, 0.0
	for i := 0; i < 400000; i++ {
		L = int(blockFracL(t, k0))
		bb = frac(-t / (4 * math.Pi * k0 * k0))
		lo := 2 * math.Min(bb, 1-bb) * float64(L)
		if lo >= 8 && lo <= 16 {
			break
		}
		k0 += 1e8
	}
	kc := k0 + float64(L)/2
	paseo := func(a, b float64) [][2]float64 {
		pts := make([][2]float64, 0, L+1)
		x, y, ph := 0.0, 0.0, 0.0
		pts = append(pts, [2]float64{0, 0})
		for j := 0; j < L; j++ {
			s, co := math.Sincos(2 * math.Pi * ph)
			x, y = x+co, y+s
			pts = append(pts, [2]float64{x, y})
			ph = frac(ph + frac(a+b*float64(2*j+1)))
		}
		return pts
	}
	// the S positions of the hunting window: only the carrier should move
	walks := make([][][2]float64, S)
	angs := make([]float64, S)
	derivaMax, sepMax := 0.0, 0.0
	for s := 0; s < S; s++ {
		d := float64(s)*h - guard
		ts := p12add(p12f(t), p12f(d))
		a := p12frac1(p12div(ts, p12mul(p12mul(p12pi, p12f(2)), p12f(k0))))
		b := frac(-(t + d) / (4 * math.Pi * k0 * k0))
		walks[s] = paseo(a, b)
		angs[s] = math.Mod(-d*math.Log(kc), 2*math.Pi)
		dv := math.Abs(d)*float64(L)/k0 + math.Abs(d)*float64(L)*float64(L)/(2*k0*k0)
		if dv > derivaMax {
			derivaMax = dv
		}
		for j := range walks[s] {
			e := math.Hypot(walks[s][j][0]-walks[0][j][0], walks[s][j][1]-walks[0][j][1])
			if e > sepMax {
				sepMax = e
			}
		}
	}
	terminos := nTop - foldFrom + 1
	bloques := math.Log(float64(nTop)/float64(foldFrom)) / math.Log1p(cFrac)
	sinPlegar := float64(terminos) * float64(S)
	giro := math.Mod(h*math.Log(kc), 2*math.Pi)
	fmt.Printf("SS plano13 t=%.0e nTop=%d foldFrom=%d c=%.4e ventana=%d puntos (h=%.5f, span=%.4f) k0=%.4e L=%d b=%.6f deriva_forma=%.3e rad sep_dibujada=%.3e giro/paso=%.4f rad terminos=%.4e bloques=%.4e ahorro=%.1fx\n",
		t, nTop, foldFrom, cFrac, S, h, span, k0, L, bb, derivaMax, sepMax, giro, float64(terminos), bloques, sinPlegar/bloques)

	l := NewLienzo(1400, 950)
	l.Titulo("⑬ UN PLIEGUE POR VENTANA — el engranaje de Fresnel",
		fmt.Sprintf("t = %.0e, bloque k₀ = %.4g con L = %d términos: en los %d puntos de la ventana sólo gira la portadora", t, k0, L, S))

	l.Panel(40, 96, 830, 604, cGreen, "LA VENTANA DE CAZA — la misma forma dibujada en cada posición")
	paso13 := 762.0
	if S > 1 {
		paso13 = 762 / float64(S-1)
	}
	for s := 0; s < S; s++ {
		cx := 74 + float64(s)*paso13
		sn, cs := math.Sincos(angs[s] - math.Pi/2)
		l.Circ(cx, 152, 13, cBg, cGrid, 1.2)
		l.Line(cx, 152, cx+cs*11, 152+sn*11, cGold, 2, "")
	}
	l.Txt(455, 188, 12.5, cGold, "middle", fmt.Sprintf("la portadora e^(−i t ln k_c): gira %.3f rad por paso — lo ÚNICO que cambia", giro))
	// the overlapped block shapes: one drawn curve, S times
	xLo, xHi, yLo, yHi := math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)
	for _, p := range walks[0] {
		xLo, xHi = math.Min(xLo, p[0]), math.Max(xHi, p[0])
		yLo, yHi = math.Min(yLo, p[1]), math.Max(yHi, p[1])
	}
	esc := math.Min(700/(xHi-xLo+1e-9), 400/(yHi-yLo+1e-9))
	ox, oy := 455-esc*(xLo+xHi)/2, 424-esc*(yLo+yHi)/2
	for s := 0; s < S; s++ {
		pts := make([][2]float64, len(walks[s]))
		for i, p := range walks[s] {
			pts[i] = [2]float64{ox + esc*p[0], oy + esc*p[1]}
		}
		col := cGreen
		if s%2 == 1 {
			col = cBlue
		}
		l.Camino(pts, col, 1.1, 0.55)
	}
	l.Txt(455, 646, 12.5, cInk, "middle", fmt.Sprintf("%d paseos dibujados uno sobre otro — se ve UNA sola curva", S))
	l.Mono(455, 668, 11.5, cDim, "middle", fmt.Sprintf("separación máxima medida = %.2e unidades del paseo (deriva física de la forma: %.2e rad)", sepMax, derivaMax))

	l.Panel(890, 96, 470, 300, cGold, "EL PLIEGUE: un bloque, un solo trámite")
	for i := 0; i < 26; i++ {
		y := 148 + float64(i)*7.2
		l.Flecha(916, y, 1046, 216+float64(i)*0.6, cDim, 0.9)
	}
	l.Txt(916, 138, 11.5, cDim, "start", fmt.Sprintf("%d términos", L))
	l.Rect(1050, 186, 96, 74, cPanel, cGold, 0.95)
	l.Txt(1098, 218, 14, cGold, "middle", "FRESNEL")
	l.Txt(1098, 238, 11, cDim, "middle", "un pliegue")
	l.Flecha(1152, 222, 1252, 222, cGreen, 5)
	l.Circ(1300, 236, 34, "none", cGreen, 2)
	l.Curva(1266, 236, 1334, 236, 16, cGreen, 2, "")
	l.Txt(1300, 292, 12, cGreen, "middle", "el balde de luz")
	l.Mono(916, 330, 11.5, cInk, "start", fmt.Sprintf("términos × puntos : %.3e", sinPlegar))
	l.Mono(916, 350, 11.5, cInk, "start", fmt.Sprintf("bloques plegados  : %.3e", bloques))
	l.Mono(916, 370, 12.5, cGold, "start", fmt.Sprintf("ahorro            : %.0f×", sinPlegar/bloques))

	l.Panel(890, 412, 470, 288, cRose, "EL BLOQUE SIN CURVATURA — un trazo y listo")
	l.Ejes(922, 470, 410, 176, "a", "D(a)/L")
	dk := make([][2]float64, 0, 401)
	for i := 0; i <= 400; i++ {
		a := -3.5/float64(L) + 7.0/float64(L)*float64(i)/400
		v := nucleoDirichlet(a, float64(L)) / float64(L)
		dk = append(dk, [2]float64{escala(float64(i), 0, 400, 922, 1332), escala(v, -0.35, 1.05, 646, 470)})
	}
	l.Line(922, escala(0, -0.35, 1.05, 646, 470), 1332, escala(0, -0.35, 1.05, 646, 470), cGrid, 1.2, "")
	l.Camino(dk, cRose, 2.2, 1)
	l.Mono(1127, 462, 11, cDim, "middle", fmt.Sprintf("D(a) = sin(πaL)/sin(πa),  L = %d", L))
	l.Txt(1127, 676, 11.5, cDim, "middle", "si b·L < 1e-9 el bloque es una progresión: núcleo de Dirichlet")

	l.Formula(40, 712, 1320, fmt.Sprintf("bloque(t) = e^(−i t ln k_c) · F(forma) ,  con la forma constante a menos de %.1e rad en toda la ventana", derivaMax))
	l.Nota(40, 760, 1320, cGold, []string{
		fmt.Sprintf("EN CRIOLLO: un bloque de %d términos se pliega UNA sola vez por ventana; después, en cada uno de los %d puntos, sólo se lo hace girar.", L, S),
		fmt.Sprintf("La forma interna del bloque no cambia: en toda la ventana se mueve %.1e radianes, o sea nada — por eso los %d dibujos son uno solo.", derivaMax, S),
		fmt.Sprintf("Eso cambia %.2e evaluaciones por %.2e pliegues: %.0f veces menos trabajo, y el bloque entra al balde como un solo término gordo.", sinPlegar, bloques, sinPlegar/bloques),
		"Y si el bloque no tiene curvatura, ni pliegue hace falta: es una progresión geométrica y se resuelve de un trazo.",
	})
	l.Pie("Plegar una vez y girar muchas: la ventana entera al precio de un solo bloque.")
	l.Guardar("laminas/plano-13-un-pliegue.svg")
}
