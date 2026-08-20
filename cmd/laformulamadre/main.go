package main

// LA FORMULA MADRE - the captain's flash: put the Midpoint Theorem inside the
// mother formula (the zeros<->primes duality) and send the bat to find it.
// Three probes, smallest possible, on 649 real zeros (Riemann-Siegel):
//   A) single echo at the twin CENTERS m (12, 18, 30, ...): the mother formula
//      predicts SILENCE - every center is composite (m = 0 mod 6), Lambda = 0.
//   B) echo at the PAIR period log(p*q) = log(m^2-1): also composite, silence.
//   C) the PAIR layer: form factor K(T) = |sum e^{i gamma T}|^2 / M - does the
//      two-zero layer hear the centers where the one-zero layer is deaf?
// Exploration mode: nothing is registered unless something significant appears.

import (
	"fmt"
	"math"
	"os"
)

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zetaZ(t float64) float64 {
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sg := 1.0
	if N%2 == 0 {
		sg = -1
	}
	return s + sg*math.Pow(2*math.Pi/t, 0.25)*c0
}

func ceros(t1 float64) []float64 {
	var g []float64
	a, za := 10.0, zetaZ(10.0)
	for b := 10.05; b <= t1; b += 0.05 {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 50; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zetaZ(m) < 0) {
					hi = m
				} else {
					lo = m
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

type dado struct{ s uint64 }

func (d *dado) u() float64 {
	d.s ^= d.s << 13
	d.s ^= d.s >> 7
	d.s ^= d.s << 17
	return float64(d.s>>11) / float64(uint64(1)<<53)
}

func media(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func desvio(v []float64) float64 {
	m := media(v)
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return math.Sqrt(s / float64(len(v)-1))
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "audit20" {
		audit20()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fase20" {
		fase20()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fase19" {
		fase19()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fase18" {
		fase18()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fase17" {
		fase17()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fase16" {
		fase16()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fase15" {
		fase15()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fase14" {
		fase14()
		return
	}
	fmt.Println("🦇🧬 LA FÓRMULA MADRE — ¿los ceros escuchan los CENTROS de los gemelos?")
	g := ceros(1000)
	M := float64(len(g))
	fmt.Printf("   %d ceros hasta γ=1000 · exploración: nada se registra sin hallazgo\n", len(g))
	fmt.Println("   predicción (la ley Λ de F360): los centros son compuestos ⟹ la capa de UN")
	fmt.Println("   cero debe callar; la pregunta abierta es la capa de a DOS (factor de forma).")

	eco := func(T float64) float64 {
		s := 0.0
		for _, x := range g {
			s += math.Cos(x * T)
		}
		return 2 * s / M
	}
	ff := func(T float64) float64 { // form factor |sum e^{i g T}|^2 / M
		c, s := 0.0, 0.0
		for _, x := range g {
			c += math.Cos(x * T)
			s += math.Sin(x * T)
		}
		return (c*c + s*s) / M
	}

	centros := []int{12, 18, 30, 42, 60, 72, 102, 108, 138, 150, 180, 192, 198}
	primos := []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}
	d := &dado{s: 20260820}
	tMin, tMax := math.Log(5.0), math.Log(200.0)

	z := func(vals []float64, f func(float64) float64) (float64, float64) {
		var res []float64
		for r := 0; r < 2000; r++ {
			res = append(res, f(tMin+(tMax-tMin)*d.u()))
		}
		return (media(vals) - media(res)) / math.Max(desvio(res)/math.Sqrt(float64(len(vals))), 1e-12), media(res)
	}

	// A) single-layer echo at centers / primes / pair-products
	var eC, eP, eQ, aC, aP []float64
	for _, m := range centros {
		eC = append(eC, -eco(math.Log(float64(m))))
		eQ = append(eQ, -eco(math.Log(float64(m*m-1))))
		aC = append(aC, ff(math.Log(float64(m))))
	}
	for _, p := range primos {
		eP = append(eP, -eco(math.Log(float64(p))))
		aP = append(aP, ff(math.Log(float64(p))))
	}
	zC, _ := z(eC, func(T float64) float64 { return -eco(T) })
	zP, _ := z(eP, func(T float64) float64 { return -eco(T) })
	zQ, _ := z(eQ, func(T float64) float64 { return -eco(T) })
	fmt.Println("\n§A · LA CAPA DE UN CERO — eco -E(T), z de la MEDIA contra períodos al azar")
	fmt.Printf("   primos p (control positivo)      : media %+.4f → %+6.1fσ\n", media(eP), zP)
	fmt.Printf("   CENTROS gemelos m                : media %+.4f → %+6.1fσ\n", media(eC), zC)
	fmt.Printf("   período de pareja log(m²−1)=log(p·q): media %+.4f → %+6.1fσ\n", media(eQ), zQ)

	// C) pair layer: form factor at centers vs primes vs random
	zFC, base := z(aC, ff)
	zFP, _ := z(aP, ff)
	fmt.Println("\n§C · LA CAPA DE A DOS — factor de forma K(T) = |Σe^{iγT}|²/M")
	fmt.Printf("   base al azar %.3f · primos %+6.1fσ · CENTROS %+6.1fσ\n", base, zFP, zFC)

	// the per-center detail, because 12 sits 0.08 from log 11 and log 13
	fmt.Println("\n§D · centro por centro (contagio vigilado: el 12 vive a 0,08 de log 11/13)")
	for _, m := range centros[:6] {
		fmt.Printf("   m=%3d  -E(log m) = %+7.4f   K(log m) = %7.3f\n",
			m, -eco(math.Log(float64(m))), ff(math.Log(float64(m))))
	}

	fmt.Println("\n§ LECTURA")
	switch {
	case math.Abs(zC) < 3 && math.Abs(zFC) < 3:
		fmt.Println("   los centros CALLAN en las dos capas: la fórmula madre no los conoce ni de a")
		fmt.Println("   uno ni de a dos a esta profundidad. El Punto Medio vive en una capa que los")
		fmt.Println("   ceros no escuchan directo — si entra, entra por correlaciones más finas o")
		fmt.Println("   más ceros. Nulo limpio, coherente con la ley Λ: se registra sólo la lectura.")
	case math.Abs(zFC) >= 3 && math.Abs(zC) < 3:
		fmt.Println("   ⚡ LA CAPA DE A DOS ESCUCHA LO QUE LA DE A UNO NO: los centros suenan en el")
		fmt.Println("   factor de forma. ESO sería hallazgo — replicar con más ceros antes de nombrar.")
	default:
		fmt.Println("   ⚡ hay señal en la capa de un cero donde la ley Λ predice silencio — o es")
		fmt.Println("   contagio de los primos vecinos o es hallazgo: separar antes de decir nada.")
	}

	espejo(g, d)

	// replication of the mirror's section C with twice the depth
	fmt.Println("\n🪞 RÉPLICA — lo mismo con ceros hasta γ = 2000")
	g2 := ceros(2000)
	fmt.Printf("   %d ceros\n", len(g2))
	espejoC(g2, d)
}
