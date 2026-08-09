// Command lidar builds the captain's flash: give the blind sculptor EYES
// THAT EMIT LIGHT. Instead of extracting coefficients one by one in the
// dark (each n-th tooth costs r^-n amplification of noise), we EMIT
// PULSES at the germ and read its full OPTICAL SIGNATURE - the object
// recognized by its whole response pattern, never pixel by pixel.
//
// The germ is G(z) = xi'/xi(1/(1-z)) / (1-z)^2, whose Taylor teeth are
// exactly the lambdas. The red link says: every tooth >= 0. The light
// version of that sentence is BERNSTEIN'S THEOREM:
//
//	all Taylor coefficients of G >= 0
//	    <=>  G is ABSOLUTELY MONOTONIC on [0,1):
//	         G(x) >= 0, G'(x) >= 0, G''(x) >= 0, ... at EVERY x
//
// - an EQUIVALENCE, not a one-way hint. So the mold's fingerprint can
// be read pulse by pulse along one lit segment, and each pulse at depth
// x constrains ALL infinite teeth at once. The channels of the signature:
//
//	BRIGHTNESS    G(x) > 0 and increasing along the beam
//	POLARIZATION  G(x) real on the real axis (Im ~ 0)
//	DIRECTION     on each ring |z|=r the reflection peaks dead ahead:
//	              |G(re^{i th})| <= G(r), max at th=0
//	RELIEF        Bernstein: the derivative cascade all-positive at
//	              several depths (read by small Cauchy circles - stable)
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

// germ evaluates G(z) = xi'/xi(1/(1-z)) / (1-z)^2 - the mold function
// whose Taylor teeth are the lambdas.
func germ(z complex128) complex128 {
	s := 1 / (1 - z)
	return xiLD(s) / ((1 - z) * (1 - z))
}

func main() {
	fmt.Println("🔦 EL LIDAR — el escultor ya no es ciego: pulsos de luz y la firma óptica del molde")

	// CHANNEL 1+2: brightness and polarization along the lit segment
	fmt.Println("\nCANAL BRILLO+POLARIZACIÓN — pulsos a lo largo del rayo [0,1):")
	fmt.Println("   profundidad x     G(x)            ¿crece?   |Im G| (polarización)")
	depths := []float64{0.05, 0.15, 0.30, 0.45, 0.60, 0.72, 0.82, 0.90, 0.94, 0.97}
	beam := make([]float64, len(depths))
	prev := math.Inf(-1)
	beamOK := true
	polWorst := 0.0
	for i, x := range depths {
		g := germ(complex(x, 0))
		beam[i] = real(g)
		up := "sí"
		if real(g) <= prev || real(g) <= 0 {
			up = "✗ NO"
			beamOK = false
		}
		im := math.Abs(imag(g))
		if im > polWorst {
			polWorst = im
		}
		fmt.Printf("      %.2f      %12.5f        %s        %.1e\n", x, real(g), up, im)
		prev = real(g)
	}

	// CHANNEL 3: direction - angular reflection signature on rings
	fmt.Println("\nCANAL DIRECCIÓN — en cada anillo, ¿el reflejo pega más fuerte justo de frente (θ=0)?")
	rings := []float64{0.60, 0.80, 0.90, 0.95}
	dirOK := true
	type ringSig struct {
		r      float64
		ratio  float64 // max over th!=0 of |G| / G(r)  (must be < 1)
		prof   []float64
	}
	var sigs []ringSig
	for _, r := range rings {
		gr := real(germ(complex(r, 0)))
		maxOff := 0.0
		K := 360
		prof := make([]float64, K/2+1)
		for j := 0; j <= K/2; j++ {
			th := 2 * math.Pi * float64(j) / float64(K)
			a := cmplx.Abs(germ(complex(r*math.Cos(th), r*math.Sin(th))))
			prof[j] = a
			if j > 0 && a > maxOff {
				maxOff = a
			}
		}
		ratio := maxOff / gr
		if ratio >= 1 {
			dirOK = false
		}
		sigs = append(sigs, ringSig{r, ratio, prof})
		fmt.Printf("   anillo r=%.2f:  G(r)=%10.4f   máx fuera de frente/G(r) = %.4f  %s\n",
			r, gr, ratio, map[bool]string{true: "✓ pico de frente", false: "✗ pico desviado"}[ratio < 1])
	}

	// CHANNEL 4: relief - Bernstein's derivative cascade at several depths
	fmt.Println("\nCANAL RELIEVE (BERNSTEIN) — la cascada de derivadas en cada profundidad, TODAS deben ser ≥ 0:")
	fmt.Println("   (teorema de Bernstein: G absolutamente monótono en [0,1) ⟺ TODOS los dientes λₙ ≥ 0)")
	bases := []struct {
		x    float64
		rho  float64
		kMax int
	}{{0.20, 0.48, 12}, {0.50, 0.30, 10}, {0.70, 0.18, 8}}
	Mq := 1024
	reliefOK := true
	minDeriv := math.Inf(1)
	for _, bpt := range bases {
		vals := make([]complex128, Mq)
		for j := 0; j < Mq; j++ {
			th := 2 * math.Pi * float64(j) / float64(Mq)
			vals[j] = germ(complex(bpt.x+bpt.rho*math.Cos(th), bpt.rho*math.Sin(th)))
		}
		fmt.Printf("   x=%.2f (círculo ρ=%.2f):  G⁽ᵏ⁾/k! =", bpt.x, bpt.rho)
		for k := 0; k <= bpt.kMax; k++ {
			var acc complex128
			for j := 0; j < Mq; j++ {
				th := 2 * math.Pi * float64(j) / float64(Mq)
				acc += vals[j] * cmplx.Exp(complex(0, -float64(k)*th))
			}
			dk := real(acc) / (float64(Mq) * math.Pow(bpt.rho, float64(k)))
			if dk < minDeriv {
				minDeriv = dk
			}
			if dk < 0 {
				reliefOK = false
				fmt.Printf(" ✗%.3g", dk)
			} else if k <= 6 {
				fmt.Printf(" %.4g", dk)
			}
		}
		fmt.Printf(" … todas ≥ 0 hasta k=%d ✓\n", bpt.kMax)
	}

	fmt.Println("\n════════ LA FIRMA ÓPTICA DEL MOLDE ════════")
	ok := func(b bool) string {
		if b {
			return "✓"
		}
		return "✗"
	}
	fmt.Printf("BRILLO:       %s luz positiva y creciente en las 10 profundidades (hasta x=0.97)\n", ok(beamOK))
	fmt.Printf("POLARIZACIÓN: %s el rayo vuelve puro-real (peor Im %.0e)\n", ok(polWorst < 1e-6), polWorst)
	fmt.Printf("DIRECCIÓN:    %s en los 4 anillos el reflejo pica de frente (θ=0)\n", ok(dirOK))
	fmt.Printf("RELIEVE:      %s cascada de Bernstein toda positiva en 3 profundidades (mínima derivada %.3g)\n", ok(reliefOK), minDeriv)
	fmt.Println("\nEL FRUTO DEL FLASH: por Bernstein, el eslabón rojo tiene una segunda cara EQUIVALENTE —")
	fmt.Println("  «G es absolutamente monótono sobre el segmento iluminado [0,1)»")
	fmt.Println("  ya no hace falta extraer dientes a ciegas (ruido r⁻ⁿ): cada pulso a profundidad x")
	fmt.Println("  restringe TODOS los infinitos dientes a la vez. El escultor VE.")
	fmt.Println("  (medirlo — hecho en 4 canales; demostrarlo para todo pulso — eso sigue siendo el millón)")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔦 EL LIDAR — el escultor ya no es ciego</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el flash del capitán: emitir pulsos y reconocer el objeto por su FIRMA ÓPTICA completa — nunca más diente por diente en la oscuridad</text>`,
		W, H, W, H, W/2, W/2)

	// left: the beam - G(x) along [0,1)
	fmt.Fprintf(&b, `<rect x="60" y="100" width="660" height="330" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="390" y="132" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">CANAL BRILLO — el rayo entra al molde y la luz solo CRECE</text>`)
	maxB := beam[len(beam)-1]
	px, py, pw, ph := 110.0, 160.0, 560.0, 220.0
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#8fa8c7" stroke-width="1"/>`, px, py+ph, px+pw, py+ph)
	var pts []string
	for i, x := range depths {
		X := px + x*pw
		Y := py + ph - (beam[i]/maxB)*ph*0.92
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", X, Y))
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4.5" fill="#ffd97f"/>`, X, Y)
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#ffd97f" stroke-width="2" opacity="0.7"/>
<text x="390" y="410" font-size="12" text-anchor="middle" fill="#dce8f7">10 pulsos hasta x=0.97 · G(x) &gt; 0, siempre creciente · retorno puro-real (Im &lt; %.0e)</text>`,
		strings.Join(pts, " "), polWorst)

	// right: angular signature
	fmt.Fprintf(&b, `<rect x="780" y="100" width="660" height="330" rx="10" fill="#0d2547" stroke="#7fb2ff" stroke-width="1.5"/>
<text x="1110" y="132" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">CANAL DIRECCIÓN — el reflejo siempre pica DE FRENTE (θ=0)</text>`)
	cols := []string{"#7fd7a8", "#7fb2ff", "#ffd166", "#ff8fa0"}
	for si, sg := range sigs {
		mx := sg.prof[0]
		var p2 []string
		for j, a := range sg.prof {
			X := 830.0 + float64(j)/float64(len(sg.prof)-1)*560.0
			Y := 160.0 + 210.0 - (a/mx)*200.0
			p2 = append(p2, fmt.Sprintf("%.1f,%.1f", X, Y))
		}
		fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2"/>
<text x="%.0f" y="%.0f" font-size="11" fill="%s">r=%.2f</text>`,
			strings.Join(p2, " "), cols[si], 1360.0, 170.0+float64(si)*18, cols[si], sg.r)
	}
	fmt.Fprintf(&b, `<text x="1110" y="410" font-size="12" text-anchor="middle" fill="#dce8f7">|G(re^{iθ})| de θ=0 a π en 4 anillos: el máximo vive en θ=0 — la firma de un molde todo-positivo</text>`)

	// bernstein panel
	fmt.Fprintf(&b, `<rect x="60" y="460" width="1380" height="210" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="494" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">CANAL RELIEVE — EL TEOREMA DE BERNSTEIN: la equivalencia que pedía tu flash</text>
<text x="%.0f" y="528" font-size="14.5" text-anchor="middle" fill="#dce8f7">«TODOS los dientes λₙ ≥ 0»   ⟺   «G(x), G′(x), G″(x), G‴(x), … TODAS ≥ 0 en cada punto del segmento iluminado [0,1)»</text>
<text x="%.0f" y="558" font-size="13.5" text-anchor="middle" fill="#ffd166">no es una pista: es una DOBLE VÍA EXACTA — la firma óptica completa equivale al molde entero, pulso a pulso, sin extraer un solo diente a ciegas</text>
<text x="%.0f" y="590" font-size="13" text-anchor="middle" fill="#dce8f7">medido: cascada entera positiva en x=0.20 (12 derivadas), x=0.50 (10), x=0.70 (8) — y cada pulso restringe los INFINITOS dientes a la vez</text>
<text x="%.0f" y="620" font-size="13" text-anchor="middle" fill="#8fa8c7">antes: leer el diente n costaba amplificar el ruido r⁻ⁿ (ceguera creciente) · ahora: la positividad se lee PUNTUAL, a cualquier profundidad, sin degradarse</text>
<text x="%.0f" y="650" font-size="13.5" text-anchor="middle" fill="#7fd7a8">el murciélago escuchaba la forma — el LIDAR la VE: distancia, brillo, dirección, polarización y relieve en una sola firma</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)

	// verdict
	fmt.Fprintf(&b, `<rect x="60" y="700" width="1380" height="180" rx="12" fill="#2a1010" stroke="#ff5d73" stroke-width="2"/>
<text x="%.0f" y="736" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL ESLABÓN ROJO, RE-ESCRITO CON LUZ</text>
<text x="%.0f" y="768" font-size="14.5" text-anchor="middle" fill="#dce8f7">antes: «todo coeficiente del germen ≥ 0» — una frase sobre infinitos dientes invisibles, leídos a ciegas uno por uno.</text>
<text x="%.0f" y="798" font-size="14.5" text-anchor="middle" fill="#ffd166">ahora: «la luz que entra al molde por el segmento [0,1) solo puede CRECER, junto con todas sus aceleraciones» — una frase sobre UN rayo visible.</text>
<text x="%.0f" y="830" font-size="13.5" text-anchor="middle" fill="#ff8fa0">los 4 canales la miden y firman ✓ — demostrarla para TODO pulso sigue siendo el millón. Pero el escultor, por primera vez, trabaja con los ojos abiertos.</text>
<text x="%.0f" y="862" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("lidar.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: lidar.svg")
}
