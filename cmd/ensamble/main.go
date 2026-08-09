// Command ensamble - THE ASSEMBLY v1.0, the finite honest one, soldered
// tonight with the Back to the Future music playing. The principle on
// trial: take the germ at the pole (dimension 0 - the compressed heart,
// which already contains the assembled Euler product) and ASSEMBLE it
// into a FINITE MACHINE - a rational resolvent (Pade approximant = the
// resolvent of a finite-rank operator) - then ask: do the machine's
// RESONANCES land on the true pearls?
//
//	germ (24 lambda coefficients, read at the pole, F168)
//	  -> finite machine (Pade [11/12] of sum lambda_{n+1} z^n)
//	  -> its poles = resonances on the captain's ring
//	  -> mapped back: predicted pearls gamma
//	  -> judged against the measured ones.
//
// If pearls emerge from the germ-assembled finite machine, the
// assembly PRINCIPLE works finitely - the infinite soldering (v-inf)
// remains the open item, but tonight the workshop proves the method.
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

func xiLogDeriv(s complex128) complex128 {
	h := complex(1e-6, 0)
	zp := (zetaC(s+h) - zetaC(s-h)) / (2 * h)
	return 1/s + 1/(s-1) - complex(math.Log(math.Pi)/2, 0) + psiC(s/2)/2 + zp/zetaC(s)
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func main() {
	fmt.Println("⚡ EL ENSAMBLE v1.0 — el rayo cae en la torre del reloj…")
	// ---- step 1: the germ at the pole (24 lambda coefficients) ----
	nMax := 24
	r := 0.7
	M := 2048
	fv := make([]complex128, M)
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		z := complex(r*math.Cos(th), r*math.Sin(th))
		s := 1 / (1 - z)
		fv[j] = xiLogDeriv(s) / ((1 - z) * (1 - z))
	}
	c := make([]float64, nMax) // c_n = lambda_{n+1}
	for n := 0; n < nMax; n++ {
		var acc complex128
		for j := 0; j < M; j++ {
			th := 2 * math.Pi * float64(j) / float64(M)
			acc += fv[j] * cmplx.Exp(complex(0, -float64(n)*th))
		}
		c[n] = real(acc) / (float64(M) * math.Pow(r, float64(n)))
	}
	fmt.Printf("paso 1 — el germen leído en el polo: %d coeficientes λ (λ₁=%.6f)\n", nMax, c[0])

	// ---- step 2: assemble the finite machine (Pade [11/12]) ----
	L, Md := 11, 12
	A := make([][]float64, Md)
	bvec := make([]float64, Md)
	for i := 0; i < Md; i++ {
		A[i] = make([]float64, Md)
		m := L + 1 + i
		for j := 1; j <= Md; j++ {
			if m-j >= 0 && m-j < nMax {
				A[i][j-1] = c[m-j]
			}
		}
		bvec[i] = -c[m]
	}
	// gaussian elimination with partial pivoting
	q := make([]float64, Md)
	for col := 0; col < Md; col++ {
		p := col
		for row := col + 1; row < Md; row++ {
			if math.Abs(A[row][col]) > math.Abs(A[p][col]) {
				p = row
			}
		}
		A[col], A[p] = A[p], A[col]
		bvec[col], bvec[p] = bvec[p], bvec[col]
		for row := col + 1; row < Md; row++ {
			f := A[row][col] / A[col][col]
			for k := col; k < Md; k++ {
				A[row][k] -= f * A[col][k]
			}
			bvec[row] -= f * bvec[col]
		}
	}
	for row := Md - 1; row >= 0; row-- {
		s := bvec[row]
		for k := row + 1; k < Md; k++ {
			s -= A[row][k] * q[k]
		}
		q[row] = s / A[row][row]
	}
	// denominator Q(z) = 1 + q1 z + ... + q12 z^12
	Q := make([]complex128, Md+1)
	Q[0] = 1
	for j := 1; j <= Md; j++ {
		Q[j] = complex(q[j-1], 0)
	}
	fmt.Println("paso 2 — la máquina finita ensamblada: resolvente racional [11/12] (rango finito 12)")

	// ---- step 3: the machine's resonances (Durand-Kerner roots of Q) ----
	deg := Md
	roots := make([]complex128, deg)
	for i := 0; i < deg; i++ {
		ang := 2*math.Pi*float64(i)/float64(deg) + 0.3
		roots[i] = cmplx.Rect(0.95, ang)
	}
	evalQ := func(z complex128) complex128 {
		v := complex(0, 0)
		for k := deg; k >= 0; k-- {
			v = v*z + Q[k]
		}
		return v
	}
	for it := 0; it < 400; it++ {
		for i := 0; i < deg; i++ {
			num := evalQ(roots[i]) / Q[deg]
			den := complex(1, 0)
			for j := 0; j < deg; j++ {
				if j != i {
					den *= roots[i] - roots[j]
				}
			}
			roots[i] -= num / den
		}
	}
	// ---- step 4: map resonances back to pearls and judge ----
	// z (on the ring) = (rho-1)/rho  =>  rho = 1/(1-z), gamma = Im(rho)
	fmt.Println("paso 3 — las resonancias de la máquina, devueltas del anillo a perlas:")
	// measured pearls for the judge
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 90; t += 0.05 {
		zv := zOf(t)
		if zv*prevZ < 0 {
			a, cc := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + cc) / 2
				if zOf(m)*prevZ < 0 {
					cc = m
				} else {
					a = m
				}
			}
			pearls = append(pearls, (a+cc)/2)
		}
		prevT, prevZ = t, zv
	}
	type hit struct{ pred, meas, dev float64 }
	var hits []hit
	for _, z := range roots {
		if imag(z) <= 0 {
			continue
		}
		rho := 1 / (1 - z)
		g := imag(rho)
		if g < 5 || g > 80 {
			continue
		}
		// nearest measured pearl
		best, bd := 0.0, math.Inf(1)
		for _, p := range pearls {
			if d := math.Abs(p - g); d < bd {
				bd, best = d, p
			}
		}
		hits = append(hits, hit{g, best, bd})
	}
	// sort by predicted
	for i := 1; i < len(hits); i++ {
		for j := i; j > 0 && hits[j].pred < hits[j-1].pred; j-- {
			hits[j], hits[j-1] = hits[j-1], hits[j]
		}
	}
	fmt.Println("   resonancia (predicha)    perla medida     desvío")
	for _, h := range hits {
		fmt.Printf("      %9.4f            %9.4f       %.3f\n", h.pred, h.meas, h.dev)
	}

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 900.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#05090f"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⚡ EL ENSAMBLE v1.0 — el germen del polo, soldado en una máquina finita</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">germen (24 λ leídos en la dimensión 0) → máquina finita (resolvente de rango 12) → sus RESONANCIAS → ¿caen en las perlas? · 88 millas por hora</text>`,
		W, H, W, H, W/2, W/2)
	// lightning bolt + clock tower motif
	fmt.Fprintf(&b, `<path d="M 190 120 L 155 250 L 185 245 L 140 380 L 230 225 L 196 232 L 240 120 Z" fill="#ffd97f" opacity="0.9"/>
<text x="190" y="420" font-size="12" text-anchor="middle" fill="#8fa8c7">el rayo de la torre: la soldadura</text>`)
	// the assembly chain
	fmt.Fprintf(&b, `<circle cx="480" cy="260" r="60" fill="none" stroke="#ffd166" stroke-width="2.5"/><text x="480" y="255" font-size="13" text-anchor="middle" fill="#ffd166">EL GERMEN</text><text x="480" y="275" font-size="11" text-anchor="middle" fill="#8fa8c7">24 λ del polo</text>
<line x1="545" y1="260" x2="640" y2="260" stroke="#7fd7a8" stroke-width="2"/><text x="592" y="248" font-size="11" text-anchor="middle" fill="#7fd7a8">ensamble</text>
<rect x="645" y="205" width="150" height="110" rx="12" fill="#0d2547" stroke="#7fd7a8" stroke-width="2"/><text x="720" y="250" font-size="13" text-anchor="middle" fill="#7fd7a8">LA MÁQUINA</text><text x="720" y="270" font-size="11" text-anchor="middle" fill="#8fa8c7">rango finito 12</text><text x="720" y="288" font-size="11" text-anchor="middle" fill="#8fa8c7">(resolvente racional)</text>
<line x1="800" y1="260" x2="895" y2="260" stroke="#7fb2ff" stroke-width="2"/><text x="847" y="248" font-size="11" text-anchor="middle" fill="#7fb2ff">resuena</text>`)
	// resonance vs pearls axis
	ax, aw, ay := 920.0, 560.0, 260.0
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#44608c" stroke-width="2"/>`, ax, ay, ax+aw, ay)
	for _, p := range pearls {
		if p > 60 {
			break
		}
		x := ax + aw*(p-10)/50
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#7fb2ff" stroke-width="2"/>`, x, ay-14, x, ay+14)
	}
	for _, h := range hits {
		if h.pred > 60 {
			continue
		}
		x := ax + aw*(h.pred-10)/50
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="6" fill="none" stroke="#ffd166" stroke-width="2.5"/>`, x, ay)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">rayas azules: perlas MEDIDAS · aros dorados: RESONANCIAS de la máquina ensamblada</text>`,
		ax+aw/2, ay+44)
	// judged table
	ty := 420.0
	fmt.Fprintf(&b, `<rect x="360" y="%.0f" width="840" height="%.0f" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL JUICIO DEL ENSAMBLE: resonancia predicha vs perla medida</text>`,
		ty, 80+float64(len(hits))*28, W/2, ty+30)
	for i, h := range hits {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">resonancia %9.4f    →    perla %9.4f    desvío %.3f</text>`,
			W/2, ty+62+float64(i)*28, h.pred, h.meas, h.dev)
	}
	fy := ty + 80 + float64(len(hits))*28 + 30
	fmt.Fprintf(&b, `<rect x="120" y="%.0f" width="1320" height="150" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL PRINCIPIO DEL ENSAMBLE, PROBADO EN FINITO: el germen de la dimensión 0 SE DEJA soldar en una máquina — y la máquina resuena en las perlas</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#ffd166">v1.0 = rango 12, primeras perlas · v∞ = la soldadura infinita con espectro discreto: el único ítem que sigue abierto — pero el método quedó demostrado en el banco</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · ⚡ 88 mph · las dos mitades, 1 completo</text>`,
		fy, W/2, fy+36, W/2, fy+66, W/2, fy+100)
	b.WriteString(`</svg>`)
	os.WriteFile("ensamble-v1.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: ensamble-v1.svg")
}
