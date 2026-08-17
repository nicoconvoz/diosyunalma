// planos_e.go - desk E: the twelve-fold compass with its bow (plate 14) and
// the two postal services of the DeLorean (plate 15). Every number on these
// two sheets is measured in the run that draws them.
package main

import (
	"fmt"
	"math"
	"math/big"
)

// ===========================================================================
// PLATE 14 - LA BRÚJULA DE DOCE Y LA PROA
//
// cmd/starship compassTwelve smells |sPred| over a five-spacing window, keeps
// the loudest twelfth, cuts THAT twelfth into twelve, and repeats until the
// sector is finer than a thousandth of the spacing. At t = 1e24 a float64
// cannot hold t·ln p mod 2π at all, so the anchor phase is folded here in
// big.Float exactly the way the ship folds it in double-double.
// ===========================================================================

const p14bits = 260

func p14nf(x float64) *big.Float { return new(big.Float).SetPrec(p14bits).SetFloat64(x) }

// p14atanh sums z + z³/3 + z⁵/5 + … at p14bits.
func p14atanh(z *big.Float) *big.Float {
	sum := new(big.Float).SetPrec(p14bits).Set(z)
	z2 := new(big.Float).SetPrec(p14bits).Mul(z, z)
	term := new(big.Float).SetPrec(p14bits).Set(z)
	t := new(big.Float).SetPrec(p14bits)
	for k := 1; k < 600; k++ {
		term.Mul(term, z2)
		t.Quo(term, p14nf(float64(2*k+1)))
		sum.Add(sum, t)
		if t.Sign() == 0 || t.MantExp(nil) < sum.MantExp(nil)-p14bits-8 {
			break
		}
	}
	return sum
}

// p14ln is ln(x) for a positive integer x: halve into [1/√2, √2] (exact in
// binary) and use 2·atanh((y-1)/(y+1)) plus m·ln 2.
func p14ln(x float64) *big.Float {
	ln2 := new(big.Float).SetPrec(p14bits).Mul(p14nf(2),
		p14atanh(new(big.Float).SetPrec(p14bits).Quo(p14nf(1), p14nf(3))))
	m, y := 0, x
	for y > 1.4142135623730951 {
		y /= 2
		m++
	}
	for y < 0.7071067811865476 {
		y *= 2
		m--
	}
	by := p14nf(y)
	num := new(big.Float).SetPrec(p14bits).Sub(by, p14nf(1))
	den := new(big.Float).SetPrec(p14bits).Add(by, p14nf(1))
	r := new(big.Float).SetPrec(p14bits).Mul(p14nf(2),
		p14atanh(new(big.Float).SetPrec(p14bits).Quo(num, den)))
	return r.Add(r, new(big.Float).SetPrec(p14bits).Mul(p14nf(float64(m)), ln2))
}

func p14atanInv(n float64) *big.Float {
	N := p14nf(n)
	N2 := new(big.Float).SetPrec(p14bits).Mul(N, N)
	term := new(big.Float).SetPrec(p14bits).Quo(p14nf(1), N)
	sum := new(big.Float).SetPrec(p14bits).Set(term)
	t := new(big.Float).SetPrec(p14bits)
	for k := 1; k < 900; k++ {
		term.Quo(term, N2)
		t.Quo(term, p14nf(float64(2*k+1)))
		if k%2 == 1 {
			sum.Sub(sum, t)
		} else {
			sum.Add(sum, t)
		}
		if t.Sign() == 0 || t.MantExp(nil) < sum.MantExp(nil)-p14bits-8 {
			break
		}
	}
	return sum
}

// p14dosPi is 2π by Machin: π = 16·atan(1/5) - 4·atan(1/239).
func p14dosPi() *big.Float {
	a := new(big.Float).SetPrec(p14bits).Mul(p14nf(4), p14atanInv(5))
	a.Sub(a, p14atanInv(239))
	return a.Mul(a, p14nf(8))
}

var p14primos = []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43,
	47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101}

// p14cuerdas returns each string's anchor phase at t0 (folded mod 2π with no
// loss), its ln p and its amplitude 1/√p - the split coordinate of the ship.
func p14cuerdas(t0 float64) (anc, lnp, amp []float64) {
	dp, bt := p14dosPi(), p14nf(t0)
	for _, p := range p14primos {
		lp := p14ln(p)
		q := new(big.Float).SetPrec(p14bits).Mul(lp, bt)
		q.Quo(q, dp)
		qi, _ := q.Int(nil)
		q.Sub(q, new(big.Float).SetPrec(p14bits).SetInt(qi))
		q.Mul(q, dp)
		ph, _ := q.Float64()
		lv, _ := lp.Float64()
		anc, lnp, amp = append(anc, ph), append(lnp, lv), append(amp, 1/math.Sqrt(p))
	}
	return
}

// p14s is S(t0 + u) heard through the 26 strings.
func p14s(anc, lnp, amp []float64, u float64) float64 {
	s := 0.0
	for i := range anc {
		s -= math.Sin(anc[i]+u*lnp[i]) * amp[i]
	}
	return s / math.Pi
}

// p14paso records one cut: its twelve tastings and the twelfth kept.
type p14paso struct {
	v      [12]float64
	pick   int
	lo, hi float64
}

// p14ventana picks the stretch of sea the compass is aimed at. A window whose
// treasure sits on the boundary teaches nothing - the compass would climb to
// the edge and every base-twelve digit would come out the same. So the plate
// aims at the richest stretch whose peak lies well inside it, which is what
// aiming a compass MEANS.
func p14ventana(anc, lnp, amp []float64, spc, span float64) float64 {
	mejor, lo0 := -1.0, 0.0
	for o := 0.0; o < 24*spc; o += spc / 8 {
		pico, uPico := -1.0, 0.0
		for i := 0; i <= 720; i++ {
			u := o + span*float64(i)/720
			if v := math.Abs(p14s(anc, lnp, amp, u)); v > pico {
				pico, uPico = v, u
			}
		}
		if rel := (uPico - o) / span; rel > 0.2 && rel < 0.8 && pico > mejor {
			mejor, lo0 = pico, o
		}
	}
	return lo0
}

func p14compas(anc, lnp, amp []float64, spc, lo0, span float64) ([]int, []p14paso, float64, float64) {
	var digs []int
	var pasos []p14paso
	lo, hi := lo0, lo0+span
	for hi-lo > spc/1000 {
		w := (hi - lo) / 12
		st := p14paso{lo: lo, hi: hi}
		best := -1.0
		for j := 0; j < 12; j++ {
			st.v[j] = math.Abs(p14s(anc, lnp, amp, lo+w*(float64(j)+0.5)))
			if st.v[j] > best {
				best, st.pick = st.v[j], j
			}
		}
		digs, pasos = append(digs, st.pick+1), append(pasos, st)
		lo, hi = lo+w*float64(st.pick), lo+w*float64(st.pick+1)
	}
	return digs, pasos, (lo+hi)/2 - (lo0 + span/2), hi - lo
}

// p14cofre draws a small treasure chest.
func p14cofre(l *Lienzo, x, y, s float64) {
	l.Rect(x-s, y-s*0.42, 2*s, s*1.05, cGold, cBg, 1)
	l.Line(x-s, y-s*0.42, x+s, y-s*0.42, cBg, 1.6, "")
	l.Rect(x-s*0.16, y-s*0.52, s*0.32, s*0.42, cBg, "none", 1)
}

func plano14() {
	const t0 = 1e24
	anc, lnp, amp := p14cuerdas(t0)
	spc := espaciado(t0)
	span := 5 * spc
	lo0 := p14ventana(anc, lnp, amp, spc, span)
	digs, pasos, proa, ancho := p14compas(anc, lnp, amp, spc, lo0, span)
	addr := ""
	for i, d := range digs {
		if i > 0 {
			addr += "."
		}
		addr += fmt.Sprintf("%d", d)
	}
	uEst, razon := lo0+span/2+proa, spc/ancho
	// the wide band the tide panel draws, scanned once so the peak is measured
	uLo, uHi := -3*span, 4*span
	sMax, uPico, vPico := 0.0, 0.0, 0.0
	vals := make([]float64, 900)
	for i := range vals {
		u := uLo + (uHi-uLo)*float64(i)/899
		vals[i] = p14s(anc, lnp, amp, u)
		if math.Abs(vals[i]) > sMax {
			sMax, uPico, vPico = math.Abs(vals[i]), u, vals[i]
		}
	}
	fmt.Printf("SS plano14 t=1e24 espaciado=%.6f ventana=%.6f cortes=%d direccion=%s ancho_final=%.4e razon=espaciado/%.0f proa=%+.6f=%+.4f_esp u*=%.6f |S(u*)|=%.5f margen_min=%.4f_esp pico_banda |S|=%.5f en u=%+.4f_esp\n",
		spc, span, len(digs), addr, ancho, razon, proa, proa/spc, uEst,
		math.Abs(p14s(anc, lnp, amp, uEst)), math.Min(uEst-lo0, lo0+span-uEst)/spc, sMax, uPico/spc)

	l := NewLienzo(1400, 950)
	l.Titulo("⑭ LA BRÚJULA DE DOCE Y LA PROA",
		"cmd/starship compassTwelve — cortar en doce, oler, cortar otra vez; y después correr el marco para que el tesoro quede al medio")

	// --- A: the compass rose, one ring per measured cut ---------------------
	l.Panel(45, 100, 420, 390, cGold, "LA BRÚJULA — un anillo por corte")
	cx, cy, K := 255.0, 312.0, len(digs)
	rad := func(k int) float64 {
		if K < 2 {
			return 148
		}
		return 148 - float64(k)*(102/float64(K-1))
	}
	ang := func(j float64) float64 { return (-90 + j*30) * math.Pi / 180 }
	for j := 0; j < 12; j++ {
		a := ang(float64(j))
		l.Line(cx+30*math.Cos(a), cy+30*math.Sin(a), cx+160*math.Cos(a), cy+160*math.Sin(a), cGrid, 1, "")
		b := ang(float64(j) + 0.5)
		l.Mono(cx+170*math.Cos(b), cy+170*math.Sin(b)+4, 10, cDim, "middle", fmt.Sprintf("%d", j+1))
	}
	ruta := make([][2]float64, 0, K)
	for k := 0; k < K; k++ {
		a := ang(float64(digs[k]) - 0.5)
		ruta = append(ruta, [2]float64{cx + rad(k)*math.Cos(a), cy + rad(k)*math.Sin(a)})
	}
	l.Camino(ruta, cGold, 3, 0.75)
	for k := 0; k < K; k++ {
		r := rad(k)
		l.Circ(cx, cy, r, "none", cGrid, 1.1)
		for j := 0; j < 12; j++ {
			b := ang(float64(j) + 0.5)
			bx, by := cx+r*math.Cos(b), cy+r*math.Sin(b)
			if j == digs[k]-1 {
				l.Circ(bx, by, 11.5, cGold, "", 0)
				l.Txt(bx, by+4, 11, cBg, "middle", fmt.Sprintf("%d", j+1))
			} else {
				l.Circ(bx, by, 4.5, cPanel, cDim, 1)
			}
		}
	}
	aL := ang(float64(digs[K-1]) - 0.5)
	l.Flecha(cx+(rad(K-1)-13)*math.Cos(aL), cy+(rad(K-1)-13)*math.Sin(aL),
		cx+19*math.Cos(aL), cy+19*math.Sin(aL), cGold, 2.2)
	p14cofre(l, cx, cy, 13)
	l.Mono(455, 124, 11, cGold, "end", "dir "+addr)

	// --- B: the smell, twelve tastings per cut ------------------------------
	l.Panel(480, 100, 430, 390, cBlue, "EL OLFATO — |S| en las doce catas (altura relativa al corte)")
	rowH := 344.0 / float64(K)
	for k, st := range pasos {
		top := 138 + float64(k)*rowH
		base, vmax, vmin := top+rowH-24, -1.0, math.Inf(1)
		for _, v := range st.v {
			if v > vmax {
				vmax = v
			}
			if v < vmin {
				vmin = v
			}
		}
		rango := vmax - vmin
		if rango <= 0 {
			rango = 1
		}
		l.Mono(492, top+11, 10, cDim, "start",
			fmt.Sprintf("corte %d  ancho %.2e  |S| de %.4f a %.4f  elige %d", k+1, st.hi-st.lo, vmin, vmax, digs[k]))
		for j := 0; j < 12; j++ {
			x := 500 + float64(j)*33.33
			h := 5 + (st.v[j]-vmin)/rango*(rowH-58)
			if j == st.pick {
				l.Rect(x, base-h, 25, h, cGold, "none", 0.95)
			} else {
				l.Rect(x, base-h, 25, h, cBlue, "none", 0.42)
			}
		}
		l.Line(500, base, 900, base, cGrid, 1.2, "")
		if k+1 < K {
			px := 500 + float64(st.pick)*33.33
			l.Line(px, base+2, 500, top+rowH+4, cGold, 0.8, "4 4")
			l.Line(px+25, base+2, 900, top+rowH+4, cGold, 0.8, "4 4")
		}
	}

	// --- C: the bow, the same treasure framed twice -------------------------
	l.Panel(925, 100, 430, 390, cGreen, "LA PROA — el mismo tesoro, dos encuadres")
	vmax := 0.0
	for i := 0; i <= 260; i++ {
		for _, off := range []float64{0, proa} {
			if v := math.Abs(p14s(anc, lnp, amp, off+span*float64(i)/260)); v > vmax {
				vmax = v
			}
		}
	}
	marco := func(fy, off float64, etiqueta string, cen bool) {
		l.Rect(945, fy, 380, 104, cBg, cDim, 1)
		l.Txt(950, fy-6, 11.5, cDim, "start", etiqueta)
		pts := make([][2]float64, 0, 261)
		for i := 0; i <= 260; i++ {
			u := span * float64(i) / 260
			v := math.Abs(p14s(anc, lnp, amp, off+u))
			pts = append(pts, [2]float64{945 + 380*float64(i)/260, fy + 96 - v/vmax*84})
		}
		l.Camino(pts, cBlue, 1.8, 0.9)
		l.Line(1135, fy, 1135, fy+104, cRose, 1, "5 4")
		tx := 1135.0
		if !cen {
			tx = 945 + 380*uEst/span
		}
		l.Line(tx, fy+4, tx, fy+100, cGold, 1.4, "3 3")
		p14cofre(l, tx, fy+88, 9)
		// margins, measured in spacings, drawn to scale
		mi := uEst
		if cen {
			mi = span / 2
		}
		l.Rect(945, fy+112, 380*mi/span, 9, cRose, "none", 0.85)
		l.Rect(945+380*mi/span, fy+112, 380*(span-mi)/span, 9, cGreen, "none", 0.85)
		l.Mono(948, fy+136, 10.5, cRose, "start", fmt.Sprintf("margen izq %.3f esp", mi/spc))
		l.Mono(1322, fy+136, 10.5, cGreen, "end", fmt.Sprintf("margen der %.3f esp", (span-mi)/spc))
		mn := math.Min(mi, span-mi) / spc
		vt, vc := fmt.Sprintf("margen mínimo %.3f esp — entra entero", mn), cGreen
		if mn < 0.5 {
			vt, vc = fmt.Sprintf("¡CORTADO! margen mínimo %.3f esp", mn), cRose
		}
		l.Mono(1135, fy+72, 10.5, vc, "middle", vt)
	}
	marco(150, 0, "SIN proa — el marco arranca en cero", false)
	l.Flecha(945+380*uEst/span, 294, 1135, 294, cGold, 2.4)
	l.Mono(1135, 310, 11, cGold, "middle", fmt.Sprintf("la proa corre el marco %+.4f espaciados", proa/spc))
	marco(336, proa, "CON proa — el tesoro queda al medio", true)

	// --- D: the tide as a sum of prime strings ------------------------------
	l.Panel(45, 505, 1310, 195, cRose, "LA MAREA — S(u) es la suma de 26 cuerdas primas")
	xa := func(u float64) float64 { return escala(u, uLo, uHi, 70, 1340) }
	l.Mono(1340, 529, 10, cBlue, "end", "en azul las sumas parciales de 1, 3 y 8 cuerdas — en dorado las 26")
	l.Rect(xa(0), 546, xa(span)-xa(0), 100, cGold, "none", 0.10)
	l.Line(xa(0), 546, xa(0), 654, cGold, 1, "4 4")
	l.Line(xa(span), 546, xa(span), 654, cGold, 1, "4 4")
	l.Line(70, 596, 1340, 596, cGrid, 1, "")
	for _, nc := range []int{1, 3, 8} {
		vp := make([][2]float64, 0, 460)
		for i := 0; i <= 459; i++ {
			u := uLo + (uHi-uLo)*float64(i)/459
			s := 0.0
			for q := 0; q < nc; q++ {
				s -= math.Sin(anc[q]+u*lnp[q]) * amp[q]
			}
			vp = append(vp, [2]float64{xa(u), 596 - (s/math.Pi)/sMax*44})
		}
		l.Camino(vp, cBlue, 1.3, 0.55)
	}
	pts := make([][2]float64, 0, 900)
	for i, v := range vals {
		pts = append(pts, [2]float64{xa(uLo + (uHi-uLo)*float64(i)/899), 596 - v/sMax*44})
	}
	l.Camino(pts, cGold, 2.2, 0.95)
	for k := 0; k <= 5; k++ {
		l.Line(xa(float64(k)*spc), 646, xa(float64(k)*spc), 654, cDim, 1, "")
		l.Mono(xa(float64(k)*spc), 666, 9.5, cDim, "middle", fmt.Sprintf("%d", k))
	}
	l.Line(xa(uEst), 546, xa(uEst), 654, cGreen, 1.4, "2 3")
	l.Mono(xa(uEst)+6, 558, 10.5, cGreen, "start", "u* del compás")
	l.Line(xa(uPico), 546, xa(uPico), 654, cRose, 1.4, "6 4")
	l.Circ(xa(uPico), 596-vPico/sMax*44, 4, cRose, "", 0)
	l.Mono(78, 684, 10.5, cGold, "start", "la ventana del compás mide 5 espaciados")
	l.Mono(1340, 684, 10.5, cRose, "end",
		fmt.Sprintf("el pico de |S| cae en u = %+.3f esp — FUERA de la ventana, por eso el compás se va al borde", uPico/spc))

	l.Formula(45, 712, 1310,
		"S(t) ~ -(1/pi)·sum_{p<=101} sin(t·ln p)/sqrt(p)      espaciado = 2pi/ln(t/2pi)")
	l.Nota(45, 758, 1310, cGold, []string{
		"EN CRIOLLO: la ventana donde se pesca se corta en doce pedazos y el aparato huele cuál de los doce hace más ruido. Ese doceavo se queda.",
		fmt.Sprintf("Después corta ESE doceavo en doce otra vez, y otra vez: en %d cortes el pedacito quedó %.0f veces más fino que el paso entre ceros (dirección %s).", K, razon, addr),
		"El olor es la marea: una sola cuerda prima casi no mueve el agua, con tres ya hay ola, con las 26 el mar habla (las curvas azules de abajo).",
		fmt.Sprintf("La proa es el último ajuste: corre el marco %+.4f espaciados, y el tesoro pasa de %.3f espaciados de margen contra el borde a %.3f.",
			proa/spc, math.Min(uEst-lo0, lo0+span-uEst)/spc, 2.5),
	})
	l.Pie(fmt.Sprintf("medido en vivo a t = 1e24 — dirección base doce %s — sector final %.3e = espaciado/%.0f — proa %+.4f espaciados", addr, ancho, razon, proa/spc))
	l.Guardar("laminas/plano-14-la-brujula.svg")
}

// ===========================================================================
// PLATE 15 - LOS DOS CORREOS
// ===========================================================================

func p15criba(n int) []bool {
	c := make([]bool, n+1)
	c[0], c[1] = true, true
	for i := 2; i*i <= n; i++ {
		if !c[i] {
			for j := i * i; j <= n; j += i {
				c[j] = true
			}
		}
	}
	return c
}

// p15inv inverts the smooth count: the t where N(t) = objetivo.
func p15inv(objetivo float64) float64 {
	lo, hi := 10.0, 6000.0
	for i := 0; i < 160; i++ {
		m := (lo + hi) / 2
		if cuentaCeros(m) < objetivo {
			lo = m
		} else {
			hi = m
		}
	}
	return (lo + hi) / 2
}

func p15sep(n int) string {
	s, out := fmt.Sprintf("%d", n), ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out += "."
		}
		out += string(c)
	}
	return out
}

func p15puerta(l *Lienzo, x, base float64, col, etiqueta string) {
	l.Rect(x-9, base-28, 18, 28, cPanel, col, 1)
	l.Circ(x+4, base-14, 2, col, "", 0)
	l.Txt(x, base+16, 10.5, col, "middle", etiqueta)
}

func p15globo(l *Lienzo, cx, cy, r float64, col string) {
	l.Circ(cx, cy, r, cPanel, col, 1.8)
	l.Line(cx-r, cy, cx+r, cy, col, 1, "")
	l.Curva(cx, cy-r, cx, cy+r, r*0.55, col, 1, "")
	l.Curva(cx, cy-r, cx, cy+r, -r*0.55, col, 1, "")
	l.Line(cx-r*0.87, cy-r*0.5, cx+r*0.87, cy-r*0.5, col, 0.7, "3 3")
	l.Line(cx-r*0.87, cy+r*0.5, cx+r*0.87, cy+r*0.5, col, 0.7, "3 3")
}

func p15sobre(l *Lienzo, x, y, w, h float64, col, num string) {
	l.Rect(x, y, w, h, cPanel, col, 1)
	l.Line(x, y, x+w/2, y+h*0.6, col, 1.2, "")
	l.Line(x+w, y, x+w/2, y+h*0.6, col, 1.2, "")
	l.Mono(x+w/2, y+h-7, 11, col, "middle", num)
}

func plano15() {
	const g100, g1000 = 236.524, 1419.422
	n100, n1000 := cuentaCeros(g100), cuentaCeros(g1000)
	th := map[int]float64{}
	for n := 96; n <= 104; n++ {
		th[n] = p15inv(float64(n) - 0.5)
	}
	brecha, spcL := g100-th[100], espaciado(g100)

	brujula := func(n float64) float64 { return n * (math.Log(n) + math.Log(math.Log(n)) - 1) }
	x6, x8 := brujula(1e6), brujula(1e8)
	dev6, dev8 := (x6-15485863)/15485863*100, (x8-2038074743)/2038074743*100
	lim, xi := 15486000, int(x6)
	comp := p15criba(lim)
	pix, cnt, nth := 0, 0, 0
	for i := 2; i <= lim; i++ {
		if !comp[i] {
			cnt++
			if i <= xi {
				pix = cnt
			}
			if cnt == 1000000 {
				nth = i
			}
		}
	}
	W := 120
	spf, sobre := make([]int, W), 0
	for i := 0; i < W; i++ {
		v := xi + i
		if !comp[v] {
			sobre++
			continue
		}
		spf[i] = -1
		for _, q := range []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47} {
			if v%q == 0 {
				spf[i] = q
				break
			}
		}
	}
	fmt.Printf("SS plano15 N(g100=%.3f)=%.5f dev=%+.5f N(g1000=%.3f)=%.5f dev=%+.5f t^100=%.5f brecha=%+.5f=%+.4f_esp x0(1e6)=%.0f primo1e6_medido=%d dev=%+.5f%% pi(x0)=%d faltan=%d x0(1e8)=%.0f dev=%+.5f%% criba %d enteros -> %d sobreviven\n",
		g100, n100, n100-100, g1000, n1000, n1000-1000, th[100], brecha, brecha/spcL,
		x6, nth, dev6, pix, 1000000-pix, x8, dev8, W, sobre)

	l := NewLienzo(1400, 950)
	l.Titulo("⑮ LOS DOS CORREOS — saltar por dirección, sin recorrer el camino",
		"cmd/starship: al cero número N se llega por la escalera θ(T)/π+1 — cmd/hipersalto: al primo número n se llega por brújula, cancelación y criba")

	// --- LEFT GLOBE: the zeros ---------------------------------------------
	l.Panel(45, 100, 645, 590, cGreen, "GLOBO DE LOS CEROS — el listín de direcciones N(T)")
	p15globo(l, 95, 178, 28, cGreen)
	p15sobre(l, 148, 160, 56, 36, cGold, "N=1000")
	ax, ay, aw, ah := 262.0, 150.0, 393.0, 250.0
	xt := func(t float64) float64 { return escala(t, 10, 1500, ax, ax+aw) }
	yn := func(n float64) float64 { return escala(n, 0, 1060, ay+ah, ay) }
	l.Ejes(ax, ay, aw, ah, "", "N(T)")
	cp := make([][2]float64, 0, 400)
	for i := 0; i <= 399; i++ {
		t := 10 + 1490*float64(i)/399
		cp = append(cp, [2]float64{xt(t), yn(cuentaCeros(t))})
	}
	l.Camino(cp, cBlue, 2.2, 0.95)
	l.Flecha(206, 178, ax-6, yn(1000), cGold, 2)
	for _, d := range [][2]float64{{100, g100}, {1000, g1000}} {
		l.Line(ax, yn(d[0]), xt(d[1]), yn(d[0]), cGreen, 1.2, "5 4")
		l.Line(xt(d[1]), yn(d[0]), xt(d[1]), ay+ah, cGreen, 1.2, "5 4")
		l.Circ(xt(d[1]), yn(d[0]), 4, cGreen, "", 0)
		l.Mono(ax+4, yn(d[0])-6, 10.5, cGreen, "start", fmt.Sprintf("N = %.0f", d[0]))
		p15puerta(l, xt(d[1]), ay+ah, cGreen, fmt.Sprintf("%.3f", d[1]))
	}
	for _, tk := range []float64{0, 500, 1000} {
		l.Mono(xt(tk), ay+ah+34, 9.5, cDim, "middle", fmt.Sprintf("%.0f", tk))
	}
	l.Mono(ax+aw, ay+ah+34, 9.5, cDim, "end", "T")
	// staircase zoom
	l.Txt(80, 452, 12, cGold, "start", "ZOOM — la escalera real de casas cerca del número 100")
	bx, bw, by2, bh := 90.0, 565.0, 470.0, 175.0
	tl, tr := th[97]-0.6, th[103]+0.6
	xz := func(t float64) float64 { return escala(t, tl, tr, bx, bx+bw) }
	yz := func(n float64) float64 { return escala(n, 96.2, 103.8, by2+bh, by2) }
	l.Line(bx, by2+bh, bx+bw, by2+bh, cGrid, 1.4, "")
	for n := 97; n <= 103; n++ {
		l.Line(xz(th[n]), yz(float64(n-1)), xz(th[n]), yz(float64(n)), cBlue, 2, "")
		xr := bx + bw
		if n < 103 {
			xr = xz(th[n+1])
		}
		l.Line(xz(th[n]), yz(float64(n)), xr, yz(float64(n)), cBlue, 2, "")
		l.Mono(xz(th[n])+4, yz(float64(n))-5, 9.5, cDim, "start", fmt.Sprintf("%d", n))
	}
	l.Line(xz(th[100]), by2, xz(th[100]), by2+bh, cGold, 1.4, "4 4")
	l.Mono(xz(th[100]), by2-4, 10, cGold, "middle", fmt.Sprintf("dirección %.4f", th[100]))
	l.Line(xz(g100), by2+18, xz(g100), by2+bh, cGreen, 2, "")
	p15puerta(l, xz(g100), by2+bh, cGreen, fmt.Sprintf("puerta real %.3f", g100))
	l.Flecha(xz(th[100]), by2+34, xz(g100), by2+34, cRose, 1.8)
	l.Mono((xz(th[100])+xz(g100))/2, by2+26, 10, cRose, "middle",
		fmt.Sprintf("%+.4f = %+.3f espaciados", brecha, brecha/spcL))

	// --- RIGHT GLOBE: the primes -------------------------------------------
	l.Panel(710, 100, 645, 590, cGold, "GLOBO DE LOS PRIMOS — brújula, cancelación y criba")
	p15globo(l, 760, 178, 28, cGold)
	p15sobre(l, 812, 160, 60, 36, cGold, "n=10^6")
	l.Rect(900, 168, 440, 14, cPanel, cDim, 1)
	l.Mono(900, 160, 9.5, cDim, "start", "0")
	l.Mono(1340, 160, 9.5, cDim, "end", p15sep(15485863))
	zx := 900 + 440*x6/15485863
	l.Rect(zx-3, 166, 8, 18, cGold, "none", 0.95)
	l.Line(zx-3, 184, 762, 232, cGold, 0.8, "4 4")
	l.Line(zx+5, 184, 1340, 232, cGold, 0.8, "4 4")
	lo6, hi6 := x6-6000, 15485863.0+6000
	xp := func(v float64) float64 { return escala(v, lo6, hi6, 762, 1340) }
	l.Line(762, 268, 1340, 268, cDim, 2, "")
	l.Line(xp(x6), 240, xp(x6), 268, cGold, 2, "")
	l.Circ(xp(x6), 240, 5, cGold, "", 0)
	l.Mono(xp(x6), 232, 10.5, cGold, "start", fmt.Sprintf("x0 = %s", p15sep(int(x6))))
	p15puerta(l, xp(15485863), 268, cGreen, fmt.Sprintf("p = %s", p15sep(nth)))
	l.Flecha(xp(x6)+4, 258, xp(15485863)-12, 258, cRose, 1.8)
	l.Mono((xp(x6)+xp(15485863))/2, 250, 10.5, cRose, "middle",
		fmt.Sprintf("%s enteros = %.3f %% — faltan %s primos por contar", p15sep(15485863-int(x6)), math.Abs(dev6), p15sep(1000000-pix)))
	l.Mono(762, 300, 10.5, cDim, "start", fmt.Sprintf("π(x0) por cancelación Meissel/P2 = %s (medido con criba completa)", p15sep(pix)))
	// deviation bars
	l.Txt(720, 334, 12, cGold, "start", "la brújula afina con la altura — |desvío| de x0 respecto del primo verdadero")
	tope := math.Max(math.Abs(dev6), math.Abs(dev8)) * 1.28
	for i, d := range [][2]float64{{6, math.Abs(dev6)}, {8, math.Abs(dev8)}} {
		y := 350 + float64(i)*30
		l.Mono(790, y+12, 10.5, cDim, "end", fmt.Sprintf("n = 10^%.0f", d[0]))
		l.Rect(800, y, 460*d[1]/tope, 16, cBlue, "none", 0.75)
		l.Mono(806+460*d[1]/tope, y+12, 10.5, cGold, "start", fmt.Sprintf("%.4f %%", d[1]))
	}
	// the comb
	l.Txt(720, 434, 12, cGold, "start", fmt.Sprintf("LA CRIBA LOCAL — %d enteros desde x0, peinados primo por primo", W))
	col := map[int]string{2: cRose, 3: cBlue, 5: cGold, 7: cGreen, 11: cDim, 13: cInk}
	xc := func(i int) float64 { return 764 + float64(i)*4.5 }
	for r, q := range []int{2, 3, 5, 7, 11, 13} {
		y := 452 + float64(r)*15
		l.Mono(758, y+4, 10, col[q], "end", fmt.Sprintf("x%d", q))
		l.Line(764, y, 1300, y, cGrid, 0.7, "")
		for i := 0; i < W; i++ {
			if (xi+i)%q == 0 {
				l.Circ(xc(i), y, 2.4, col[q], "", 0)
			}
		}
	}
	base, alterna := 620.0, 0
	l.Line(760, base, 1304, base, cGrid, 1.4, "")
	for i := 0; i < W; i++ {
		x := xc(i)
		if spf[i] == 0 { // survivor: stands up
			l.Line(x, base, x, base-56, cGreen, 2.2, "")
			l.Circ(x, base-56, 3.6, cGreen, "", 0)
			l.Mono(x, base-66-float64(alterna%2)*11, 8.5, cGreen, "middle", fmt.Sprintf("+%d", i))
			alterna++
			continue
		}
		c, ok := col[spf[i]] // struck: falls below the line, in the colour that killed it
		if !ok {
			c = cDim
		}
		l.Line(x, base, x, base+14, c, 1.8, "")
	}
	l.Line(760, base+7, 1304, base+7, cRose, 1.2, "")
	l.Mono(1304, base+30, 9.5, cRose, "end", "abajo de la raya, los tachados — el color es el primo que los mató")
	l.Mono(760, base+52, 11, cGold, "start",
		fmt.Sprintf("de %d enteros sobreviven %d — la casa está entre ellos", W, sobre))

	l.Formula(45, 702, 1310,
		"N(T) = theta(T)/pi + 1      |      x0 = n·(ln n + ln ln n - 1)  ->  Meissel/P2  ->  criba local")
	l.Nota(45, 748, 1310, cGold, []string{
		"EN CRIOLLO: acá no se camina el camino, se toca el timbre. Al cero número mil no se llega contando mil ceros: la escalera N(T) te da la altura y ahí caés.",
		fmt.Sprintf("La escalera es una predicción honesta: para la casa 100 dice %.4f y la puerta real está en %.3f — se erró por %.3f de un espaciado, nada.", th[100], g100, math.Abs(brecha/spcL)),
		fmt.Sprintf("Con los primos igual: la brújula x0 cae a %.3f %% de la casa del primo un millón, la cancelación cuenta cuántos hay hasta ahí sin listarlos,", math.Abs(dev6)),
		fmt.Sprintf("y el último tramo lo camina la criba: tacha múltiplos de 2, de 3, de 5… y de %d enteros quedan %d en pie. Uno de esos es la casa.", W, sobre),
	})
	l.Pie(fmt.Sprintf("medido en vivo — N(γ100) = %.4f (desvío %+.4f) · N(γ1000) = %.4f (desvío %+.4f) · x0(10⁶) desvía %+.4f %% · x0(10⁸) desvía %+.4f %%",
		n100, n100-100, n1000, n1000-1000, dev6, dev8))
	l.Guardar("laminas/plano-15-los-dos-correos.svg")
}
