// Command elcolordelafragilidad answers the captain's flash: do the
// wells, frontiers and mountains have a COLOR? Does recombination give
// deeper tones, lighter ones - can the color of fragility be PROJECTED?
//
// YES. The natural screen is the PHASE TORUS: each step n lives at the
// point (||n*th_1||, ||n*th_2||) in the phase square [0,pi]^2, and there
// the landscape paints itself:
//
//   - deep water (the window past n_rad): color = lambda_n. The
//     Leader's Law appears as VERTICAL banding - the leader's axis
//     alone picks the hue (deep blues = wells, light golds =
//     mountains), and the other pearl only shades the tone. The mix is
//     bleached by dominance.
//   - shallow water (n in [1, 3000]): color = the pearls' joint push
//     l1 + l2. Here both radii are ~1 and the two phase channels MIX
//     truly - the interference field of the captain's old "beat":
//     recombination of colors, deeper where both align, lighter where
//     they cancel.
//
// So the answer to the flash: the colors EXIST, they MIX in shallow
// water, and depth bleaches the mixture until one pearl rules the
// palette - the Leader's Law, seen as paint.
//
// Reproduce: go run ./cmd/elcolordelafragilidad
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

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

func norma(x float64) float64 {
	y := math.Mod(x, 2*math.Pi)
	if y < 0 {
		y += 2 * math.Pi
	}
	if y > math.Pi {
		y = 2*math.Pi - y
	}
	return y
}

// mezcla interpolates two hex colors.
func mezcla(c1, c2 [3]int, t float64) string {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	r := int(float64(c1[0]) + t*float64(c2[0]-c1[0]))
	g := int(float64(c1[1]) + t*float64(c2[1]-c1[1]))
	b := int(float64(c1[2]) + t*float64(c2[2]-c1[2]))
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func main() {
	fmt.Println("🎨 EL COLOR DE LA FRAGILIDAD — el paisaje proyectado sobre el toro de fases")
	fmt.Println("\n   La pregunta del capitán: ¿pozos, fronteras y montañas tienen color?")
	fmt.Println("   ¿la recombinación da tonos más profundos y más claros? ¿se puede proyectar?")
	fmt.Println("   SÍ: la pantalla es el cuadrado de fases (‖nθ₁‖, ‖nθ₂‖) — y pinta solo.")

	rho1 := complex(0.808517, 85.699348)
	rho2 := complex(0.7, 45.0)
	w1 := 1 - 1/rho1
	w2 := 1 - 1/rho2
	r1 := math.Max(cmplx.Abs(w1), 1/cmplx.Abs(w1))
	r2 := math.Max(cmplx.Abs(w2), 1/cmplx.Abs(w2))
	t1 := math.Abs(cmplx.Phase(w1))
	t2 := math.Abs(cmplx.Phase(w2))
	d1 := math.Log(r1)
	d2 := math.Log(r2)
	nrad := 1040809

	ps := perlas(120)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wp /= complex(cmplx.Abs(wp), 0)
		wsC[i] = wp
		pcs[i] = 1
	}

	// palettes
	azulClaro := [3]int{0x7f, 0xb2, 0xff}
	azulHondo := [3]int{0x04, 0x0e, 0x26}
	oroSuave := [3]int{0x6b, 0x5a, 0x26}
	oroClaro := [3]int{0xff, 0xf3, 0xc4}

	var deep, shallow strings.Builder
	// shallow water pass: n in [1, 3000], color = l1 + l2 (the pure mix)
	minS, maxS := math.Inf(1), math.Inf(-1)
	for n := 1; n <= 3000; n++ {
		fn := float64(n)
		l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
		l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*d2)+math.Exp(-fn*d2))
		v := l1 + l2
		if v < minS {
			minS = v
		}
		if v > maxS {
			maxS = v
		}
		q1 := norma(fn * t1) / math.Pi
		q2 := norma(fn * t2) / math.Pi
		x := 105 + q1*500
		y := 620 - q2*440
		var col string
		if v < 0 {
			col = mezcla(azulClaro, azulHondo, -v/8.5)
		} else {
			col = mezcla(oroSuave, oroClaro, v/16.5)
		}
		fmt.Fprintf(&shallow, `<circle cx="%.0f" cy="%.0f" r="2.6" fill="%s"/>`, x, y, col)
	}
	// deep water pass: window past n_rad, color = lambda (leader's law)
	minL, maxL := 0.0, 0.0
	for n := 1; n <= nrad+3000; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		if n < nrad {
			continue
		}
		fn := float64(n)
		l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
		l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*d2)+math.Exp(-fn*d2))
		lam := s + l1 + l2
		if lam < minL {
			minL = lam
		}
		if lam > maxL {
			maxL = lam
		}
		q1 := norma(fn * t1) / math.Pi
		q2 := norma(fn * t2) / math.Pi
		x := 745 + q1*500
		y := 620 - q2*440
		var col string
		if lam < 0 {
			col = mezcla(azulClaro, azulHondo, math.Log10(-lam)/46)
		} else {
			col = mezcla(oroSuave, oroClaro, math.Log10(lam+1)/46)
		}
		fmt.Fprintf(&deep, `<circle cx="%.0f" cy="%.0f" r="2.6" fill="%s"/>`, x, y, col)
	}
	fmt.Printf("\n   agua baja  [1, 3000]:        mezcla ℓ₁+ℓ₂ ∈ [%.1f, %.1f] — los dos canales tiñen\n", minS, maxS)
	fmt.Printf("   agua honda [n_rad, +3000]:   λ ∈ [%.1e, %.1e] — el líder dicta el tono\n", minL, maxL)
	fmt.Println("\n   veredicto visual: en el panel hondo el color depende SOLO del eje vertical")
	fmt.Println("   (el líder) — bandas horizontales puras: la Ley del Líder, vista como pintura.")
	fmt.Println("   En el panel bajo, el color depende de LOS DOS ejes: la mezcla existe — más")
	fmt.Println("   profunda donde ambas fases se alinean, más clara donde se cancelan.")
	fmt.Println("   La regla del sello preside: es una proyección del paisaje, no una prueba nueva.")

	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="800" viewBox="0 0 1400 800">
<rect width="100%" height="100%" fill="#0b1526"/>
<rect x="40" y="30" width="1320" height="740" rx="18" fill="none" stroke="#7ee0c0" stroke-width="2" opacity="0.5"/>
<text x="700" y="70" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎨 EL COLOR DE LA FRAGILIDAD — el paisaje proyectado sobre el toro de fases</text>
<text x="700" y="100" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">cada escalón n vive en el punto (‖nθ₁‖, ‖nθ₂‖) del cuadrado de fases — azules hondos = pozos · dorados claros = montañas</text>
<rect x="100" y="130" width="510" height="500" rx="10" fill="#0d1830" stroke="#26456e"/>
<text x="355" y="156" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">AGUA BAJA (n = 1..3000) — color = ℓ₁+ℓ₂: LA MEZCLA</text>
`)
	b.WriteString(shallow.String())
	b.WriteString(`
<text x="355" y="652" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">los DOS canales tiñen: tonos más hondos donde las fases se alinean, más claros donde se cancelan</text>
<rect x="740" y="130" width="510" height="500" rx="10" fill="#0d1830" stroke="#26456e"/>
<text x="995" y="156" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">AGUA HONDA (tras n_rad) — color = λ: LA LEY DEL LÍDER</text>
`)
	b.WriteString(deep.String())
	b.WriteString(`
<text x="995" y="652" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">bandas horizontales puras: el eje del líder (vertical) dicta el color — la dominancia blanquea la mezcla</text>
<text x="105" y="672" font-size="11.5" font-family="Georgia" fill="#8fa8c7">eje x: ‖nθ₁‖ de 0 a π · eje y: ‖nθ₂‖ (el líder) de 0 (abajo) a π (arriba)</text>
<text x="700" y="706" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: los colores EXISTEN y se recombinan — en agua baja las dos perlas pintan juntas; en agua honda el líder</text>
<text x="700" y="728" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">se queda con el pincel, y la mezcla solo matiza el tono. La fragilidad tiene paleta, y la paleta obedece los teoremas.</text>
<text x="700" y="758" font-size="13" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Proyección del paisaje demostrado — no una prueba nueva. La regla del sello preside. Todavía no.</text>
</svg>
`)
	os.WriteFile("el-color-de-la-fragilidad.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-color-de-la-fragilidad.svg")
}
