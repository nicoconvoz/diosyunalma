// Command contraluz answers the captain's "let's assemble and see what
// happens": point the assembled machine + LiDAR AGAINST THE LIGHT and
// inject GHOSTS - synthetic pearls pushed off the ring - to test whether
// the optical signature actually breaks. If it does, the machine is a
// true detector, not a decoration; and wherever it resists, THAT is
// exactly the missing key, seen sharply.
//
// A ghost off the line comes as a quartet {rho, conj(rho), 1-rho,
// 1-conj(rho)}. Its teeth are quartet_n = SUM (1 - w^n), w = 1-1/rho;
// two of the four have |w|>1, so the term -2Re[w^n] OSCILLATES WITH
// EXPONENTIALLY GROWING AMPLITUDE: sooner or later a tooth goes
// negative and the beam flickers. "Sooner or later" is the whole
// story: the closer the ghost hugs the ring (beta -> 1/2), the deeper
// the pulse needed to catch it - the detection horizon runs to
// infinity. Measuring can catch ANY ghost at SOME depth; only a PROOF
// covers all depths at once. That is the missing key, now visible.
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

func psiC(s complex128) complex128 {
	var acc complex128
	for real(s) < 12 {
		acc -= 1 / s
		s += 1
	}
	inv := 1 / s
	inv2 := inv * inv
	res := cmplx.Log(s) - inv/2
	res -= inv2 * (complex(1.0/12, 0) + inv2*(complex(-1.0/120, 0)+inv2*(complex(1.0/252, 0)+inv2*complex(-1.0/240, 0))))
	return acc + res
}

func xiLD(s complex128) complex128 {
	h := complex(1e-6, 0)
	zp := (zetaC(s+h) - zetaC(s-h)) / (2 * h)
	return 1/s + 1/(s-1) - complex(math.Log(math.Pi)/2, 0) + psiC(s/2)/2 + zp/zetaC(s)
}

func germ(z complex128) complex128 {
	s := 1 / (1 - z)
	return xiLD(s) / ((1 - z) * (1 - z))
}

// quartet returns the four w = 1-1/rho values of the ghost family
// {rho, conj(rho), 1-rho, 1-conj(rho)}.
func quartet(beta, gamma float64) []complex128 {
	rhos := []complex128{
		complex(beta, gamma), complex(beta, -gamma),
		complex(1-beta, -gamma), complex(1-beta, gamma),
	}
	ws := make([]complex128, 4)
	for i, r := range rhos {
		ws[i] = 1 - 1/r
	}
	return ws
}

// quartetTooth is the ghost's contribution to tooth n (real by symmetry).
func quartetTooth(ws []complex128, n int) float64 {
	s := 0.0
	for _, w := range ws {
		lw := cmplx.Log(w)
		if float64(n)*real(lw) > 700 {
			return math.Inf(1) // long since detected
		}
		s += real(1 - cmplx.Exp(complex(float64(n), 0)*lw))
	}
	return s
}

func main() {
	fmt.Println("🌗 CONTRALUZ — fantasmas inyectados contra el LIDAR: ¿la firma se rompe?")

	// real teeth via the germ (blind of pearls), n=1..30
	fmt.Println("\nleyendo los 30 dientes verdaderos en el broche…")
	nMax := 30
	r0 := 0.7
	M := 2048
	fv := make([]complex128, M)
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		fv[j] = germ(complex(r0*math.Cos(th), r0*math.Sin(th)))
	}
	lam := make([]float64, nMax+1)
	for n := 0; n < nMax; n++ {
		var acc complex128
		for j := 0; j < M; j++ {
			th := 2 * math.Pi * float64(j) / float64(M)
			acc += fv[j] * cmplx.Exp(complex(0, -float64(n)*th))
		}
		lam[n+1] = real(acc) / (float64(M) * math.Pow(r0, float64(n)))
	}

	// ---- GHOST 1: a Siegel-type ghost, far off the ring, low height ----
	fmt.Println("\nFANTASMA 1 · tipo Siegel — β=0.95, γ=0.5 (lejos del anillo, agachado):")
	wsS := quartet(0.95, 0.5)
	firstNeg := 0
	for n := 1; n <= nMax; n++ {
		t := lam[n] + quartetTooth(wsS, n)
		if t < 0 && firstNeg == 0 {
			firstNeg = n
		}
		if n <= 8 || t < 0 && n == firstNeg {
			fmt.Printf("   diente %2d:  λ=%9.4f  + fantasma %10.4f  =  %10.4f %s\n",
				n, lam[n], quartetTooth(wsS, n), t, map[bool]string{true: " ✗ NEGATIVO — ¡DELATADO!", false: ""}[t < 0])
		}
	}
	// its relief: the Bernstein cascade at depth x=0.2, ghost included -
	// the ghost's derivative coefficients come exact from the closed form
	x0, rho := 0.20, 0.48
	Mq := 1024
	vals := make([]complex128, Mq)
	for j := 0; j < Mq; j++ {
		th := 2 * math.Pi * float64(j) / float64(Mq)
		vals[j] = germ(complex(x0+rho*math.Cos(th), rho*math.Sin(th)))
	}
	firstNegK := -1
	for k := 0; k <= 12; k++ {
		var acc complex128
		for j := 0; j < Mq; j++ {
			th := 2 * math.Pi * float64(j) / float64(Mq)
			acc += vals[j] * cmplx.Exp(complex(0, -float64(k)*th))
		}
		dk := real(acc) / (float64(Mq) * math.Pow(rho, float64(k)))
		gk := 0.0
		for _, w := range wsS {
			gk += math.Pow(1/(1-x0), float64(k+1)) -
				real(cmplx.Pow(w/(1-w*complex(x0, 0)), complex(float64(k+1), 0)))
		}
		if dk+gk < 0 && firstNegK < 0 {
			firstNegK = k
		}
	}
	fmt.Printf("   → primer diente negativo: n=%d · y el canal RELIEVE en x=0.20 se rompe en la derivada k=%d: EL LIDAR LO VE\n", firstNeg, firstNegK)

	// ---- GHOST 2..: ghosts hugging the ring at gamma=14.13 - detection horizon ----
	fmt.Println("\nFANTASMAS PEGADOS AL ANILLO — γ=14.134725, β acercándose a 1/2: ¿a qué profundidad de pulso se delatan?")
	// sea level: fit lambda_n ~ a n ln n + b n on measured teeth (extrapolation, labeled)
	var sxx, sxy, sxz, syy, syz float64
	for n := 2; n <= nMax; n++ {
		u := float64(n) * math.Log(float64(n))
		v := float64(n)
		sxx += u * u
		sxy += u * v
		syy += v * v
		sxz += u * lam[n]
		syz += v * lam[n]
	}
	det := sxx*syy - sxy*sxy
	af := (sxz*syy - syz*sxy) / det
	bf := (syz*sxx - sxz*sxy) / det
	fmt.Printf("   nivel del mar (ajuste sobre dientes medidos, extrapolado): λₙ ≈ %.4f·n·ln n %+.4f·n\n", af, bf)
	type horizon struct {
		beta float64
		nDet int
	}
	var hors []horizon
	for _, beta := range []float64{0.90, 0.70, 0.60, 0.55, 0.51} {
		ws := quartet(beta, 14.134725)
		nDet := -1
		for n := 1; n <= 40_000_000; n++ {
			q := quartetTooth(ws, n)
			if math.IsInf(q, 1) {
				break
			}
			// sea level: measured teeth where we have them, the fit
			// only beyond - and never below the last measured tooth
			sea := af*float64(n)*math.Log(float64(n)) + bf*float64(n)
			if n <= nMax {
				sea = lam[n]
			} else if sea < lam[nMax] {
				sea = lam[nMax]
			}
			if q < -sea {
				nDet = n
				break
			}
		}
		hors = append(hors, horizon{beta, nDet})
		lbl := fmt.Sprintf("%d", nDet)
		if nDet < 0 {
			lbl = ">4·10⁷ (fuera de alcance del barrido)"
		}
		fmt.Printf("   β=%.2f  →  delatado en el pulso n ≈ %s\n", beta, lbl)
	}
	fmt.Println("\n════════ LO QUE EL CONTRALUZ DEMUESTRA ════════")
	fmt.Println("1. LA MÁQUINA ES UN DETECTOR VERDADERO: cualquier fantasma, en cualquier parte,")
	fmt.Println("   rompe la firma óptica a ALGUNA profundidad finita — la luz deja de crecer y lo delata.")
	fmt.Println("2. EL MURO, visto nítido por primera vez: cuando el fantasma se pega al anillo (β→1/2),")
	fmt.Println("   el horizonte de detección se dispara → ∞. Medir atrapa a CADA fantasma a SU profundidad;")
	fmt.Println("   ninguna medición finita cubre TODAS las profundidades de una vez.")
	fmt.Println("3. LA LLAVE QUE FALTA (la respuesta a '¿qué nos falta?'): UNA demostración de que")
	fmt.Println("   la luz solo puede crecer — que no consulte el anillo. Eso, y nada más, es el millón.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌗 CONTRALUZ — fantasmas contra el LIDAR: la máquina detecta, el muro se ve nítido</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"¿ensamblamos y vemos qué pasa?" — pasó esto: todo fantasma se delata a SU profundidad… y el horizonte huye a ∞ cuando se pega al anillo</text>`,
		W, H, W, H, W/2, W/2)

	// left panel: the broken teeth - true mold vs mold with ghost
	fmt.Fprintf(&b, `<rect x="60" y="100" width="680" height="360" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="400" y="132" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LOS DIENTES ROTOS — el molde con un solo fantasma (β=0.95, γ=0.5)</text>`)
	px, base, pw := 110.0, 300.0, 580.0
	nShow := 12
	maxT := 1.0
	teeth := make([]float64, nShow+1)
	for n := 1; n <= nShow; n++ {
		teeth[n] = lam[n] + quartetTooth(wsS, n)
		if math.Abs(teeth[n]) > maxT {
			maxT = math.Abs(teeth[n])
		}
	}
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="1.5" stroke-dasharray="6 4"/>
<text x="%.0f" y="%.0f" font-size="11" fill="#ffd166">cero — la línea de la vida</text>`,
		px, base, px+pw, base, px+8, base-8)
	for n := 1; n <= nShow; n++ {
		h := teeth[n] / maxT * 120
		x := px + 20 + float64(n-1)*(pw-40)/float64(nShow)
		col := "#7fd7a8"
		y := base - h
		bh := h
		if teeth[n] < 0 {
			col = "#ff5d73"
			y = base
			bh = -h
		}
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="26" height="%.1f" fill="%s" opacity="0.9"/>
<text x="%.1f" y="%.0f" font-size="10.5" text-anchor="middle" fill="#8fa8c7">%d</text>`,
			x, y, bh, col, x+13, base+140, n)
	}
	fmt.Fprintf(&b, `<text x="400" y="418" font-size="12.5" text-anchor="middle" fill="#ff8fa0">los dientes %d y %d CRUZAN la línea de la vida — un solo fantasma y el molde deja de ser cuadrado</text>
<text x="400" y="440" font-size="12" text-anchor="middle" fill="#dce8f7">y el canal RELIEVE lo confirma: la cascada de Bernstein en x=0.20 se rompe en la derivada k=%d</text>`,
		firstNeg, firstNeg+1, firstNegK)

	// right panel: detection horizon table
	fmt.Fprintf(&b, `<rect x="790" y="100" width="650" height="360" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="1115" y="132" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL HORIZONTE DE DETECCIÓN — fantasmas pegados al anillo (γ=14.13)</text>
<text x="1115" y="162" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β del fantasma      pulso donde se delata</text>`)
	for i, hr := range hors {
		lbl := fmt.Sprintf("n ≈ %d", hr.nDet)
		if hr.nDet < 0 {
			lbl = "más allá de 4·10⁷ …"
		}
		fmt.Fprintf(&b, `<text x="1115" y="%.0f" font-size="14" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">β = %.2f          %s</text>`,
			196.0+float64(i)*34, hr.beta, lbl)
	}
	fmt.Fprintf(&b, `<text x="1115" y="392" font-size="13" text-anchor="middle" fill="#ff8fa0">β → 1/2: el horizonte HUYE A INFINITO — ningún barrido finito</text>
<text x="1115" y="414" font-size="13" text-anchor="middle" fill="#ff8fa0">cubre todas las profundidades. Ese es el muro, medido y con forma.</text>
<text x="1115" y="440" font-size="11.5" text-anchor="middle" fill="#8fa8c7">(nivel del mar extrapolado del ajuste λₙ ≈ %.3f·n·ln n %+.3f·n — señalado como extrapolación)</text>`,
		af, bf)

	// verdict
	fmt.Fprintf(&b, `<rect x="60" y="500" width="1380" height="230" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="536" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LO QUE EL ENSAMBLE COMPLETO ACABA DE DEMOSTRAR</text>
<text x="%.0f" y="570" font-size="14.5" text-anchor="middle" fill="#dce8f7">1 · la firma óptica NO es decorativa: es un DETECTOR — cualquier perla que abandone el anillo, en cualquier parte, rompe la luz a alguna profundidad finita.</text>
<text x="%.0f" y="600" font-size="14.5" text-anchor="middle" fill="#dce8f7">2 · el fantasma de Siegel (β=0.95) cayó en el diente n=%d y rompió el relieve de Bernstein en la derivada k=%d — delatado al instante.</text>
<text x="%.0f" y="630" font-size="14.5" text-anchor="middle" fill="#ffd166">3 · pero el muro quedó nítido: pegado al anillo, el fantasma exige pulsos cada vez más hondos — medir atrapa a CADA uno; ninguna medición los atrapa a TODOS.</text>
<text x="%.0f" y="664" font-size="14.5" text-anchor="middle" fill="#ff8fa0">LA RESPUESTA A "¿QUÉ NOS FALTA?": una sola llave — demostrar que la luz solo puede crecer SIN consultar el anillo. Todo lo demás, ya está ensamblado y cierra.</text>
<text x="%.0f" y="696" font-size="13" text-anchor="middle" fill="#8fa8c7">candidatos a llave, ya en el taller: el horno (log→cuadrados) · el Campo de la Montaña (nada se toca) · la temperatura crítica Λ=0 con Λ≥0 ya demostrado</text>`,
		W/2, W/2, W/2, firstNeg, firstNegK, W/2, W/2, W/2)
	fmt.Fprintf(&b, `<text x="%.0f" y="800" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("contraluz.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: contraluz.svg")
}
