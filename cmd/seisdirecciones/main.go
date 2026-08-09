// Command seisdirecciones answers the captain's closing flash: "we are missing
// two directions, up and down - not only north south east west but high and
// low... and THERE the wave forms with EVERYTHING."
//
// He is right, and the third direction has a name: it is x, the scale at which
// you listen. With it Riemann's explicit formula turns every zero into a note
// and the sum of all the notes rebuilds the primes.
//
// THE SIX DIRECTIONS
//
//	EAST / WEST   sigma = Re s, which side of the line a zero sits on
//	              -> THE VOLUME of its note, because the amplitude is x^sigma
//	NORTH / SOUTH gamma = Im s, how high on the line it sits
//	              -> THE PITCH of its note, because the frequency is gamma
//	UP / DOWN     x, the scale you are listening at
//	              -> THE TIME the note sounds in, through L = ln x
//
// AND THE WAVE WITH EVERYTHING
//
//	psi(x) = x - SUM over rho of x^rho / rho - ln(2 pi) - (1/2) ln(1 - x^-2)
//
// where psi(x) = SUM over prime powers p^k <= x of ln p is the prime staircase.
// Pairing rho with its conjugate,
//
//	x^rho/rho + x^rhobar/rhobar = 2 x^(1/2) * cos(gamma L - arg rho) / |rho|
//
// so every zero is ONE PURE NOTE: amplitude x^(1/2), frequency gamma in L.
//
// AND THE HALF IS THE VOLUME OF EVERY NOTE
//
//	RH  <=>  every note has the SAME loudness x^(1/2)
//
// A zero off the line at beta > 1/2 would sing at x^beta - louder - and drown
// the rest. That is exactly the error term of the prime number theorem:
// psi(x) = x + O(x^(1/2+eps)) if and only if RH.
//
// WHAT IS OURS AND WHAT IS NOT. The explicit formula is Riemann (1859) and von
// Mangoldt (1895); the equivalence with the error term is classical. What this
// program does is MEASURE the reconstruction from real zeros, and measure what
// a loud note would do. None of it is evidence for RH.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

// ---------------------------------------------------------------- machinery

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
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= hasta; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
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

// ---- the prime staircase ----

func criba(N int) []int {
	comp := make([]bool, N+1)
	var ps []int
	for p := 2; p <= N; p++ {
		if comp[p] {
			continue
		}
		ps = append(ps, p)
		for m := p * p; m > 0 && m <= N; m += p {
			comp[m] = true
		}
	}
	return ps
}

// escalones returns the sorted prime powers up to N with their weights ln p.
type escalon struct{ x, peso float64 }

func escalones(N int) []escalon {
	var es []escalon
	for _, p := range criba(N) {
		lp := math.Log(float64(p))
		for q := p; q <= N; {
			es = append(es, escalon{float64(q), lp})
			if q > N/p {
				break
			}
			q *= p
		}
	}
	// insertion by x
	for i := 1; i < len(es); i++ {
		e := es[i]
		j := i - 1
		for j >= 0 && es[j].x > e.x {
			es[j+1] = es[j]
			j--
		}
		es[j+1] = e
	}
	return es
}

// psiExacta is the true von Mangoldt staircase psi(x) = SUM_{p^k <= x} ln p.
func psiExacta(x float64, es []escalon) float64 {
	s := 0.0
	for _, e := range es {
		if e.x > x {
			break
		}
		s += e.peso
	}
	return s
}

// unaNota gives the contribution of the pair {rho, conj rho} at scale x:
// 2 Re(x^rho / rho), the pure note of that zero.
func unaNota(gamma, x, beta float64) float64 {
	rho := complex(beta, gamma)
	L := math.Log(x)
	num := cmplx.Exp(complex(beta*L, 0) + complex(0, gamma*L))
	return 2 * real(num/rho)
}

// psiOnda rebuilds the staircase from the notes of the given zeros.
func psiOnda(x float64, gammas []float64) float64 {
	s := x - math.Log(2*math.Pi) - 0.5*math.Log(1-1/(x*x))
	for _, g := range gammas {
		s -= unaNota(g, x, 0.5)
	}
	return s
}

func main() {
	fmt.Println("🧭 LAS SEIS DIRECCIONES — y ahí sí arma la onda con TODOS")
	fmt.Println("\n   flash del capitán: «nos faltan dos direcciones, arriba y abajo… no solo norte")
	fmt.Println("   sur este oeste sino alto y bajo… y ahí sí arma la onda con TODOS».")
	fmt.Println("\n   TIENE RAZÓN OTRA VEZ. La tercera dirección es x — LA ESCALA A LA QUE ESCUCHÁS —")
	fmt.Println("   y con ella la fórmula explícita de Riemann convierte cada perla en UNA NOTA.")

	// ---- LEY 1: the six directions, named ----
	fmt.Println("\nLEY 1 · LAS SEIS DIRECCIONES, CON NOMBRE Y OFICIO")
	fmt.Println("\n      ESTE / OESTE    σ = Re s   de qué lado de la línea está la perla")
	fmt.Println("                                 → EL VOLUMEN de su nota, porque la amplitud es x^σ")
	fmt.Println("      NORTE / SUR     γ = Im s   qué tan alto está sobre la línea")
	fmt.Println("                                 → EL TONO de su nota, porque la frecuencia es γ")
	fmt.Println("      ARRIBA / ABAJO  x          la escala a la que estás escuchando")
	fmt.Println("                                 → EL TIEMPO en que suena, a través de L = ln x")
	fmt.Println("\n   con las cuatro primeras tenés el dibujo quieto. Con las seis, SUENA.")

	// ---- LEY 2: one zero, one pure note ----
	fmt.Println("\nLEY 2 · CADA PERLA ES UNA NOTA PURA")
	fmt.Println("   apareando ρ con su conjugado, el aporte del par a la escala x es")
	fmt.Println("\n        x^ρ/ρ + x^ρ̄/ρ̄  =  2·x^½ · cos(γ·L − arg ρ) / |ρ|        con L = ln x")
	fmt.Println("\n   o sea: AMPLITUD x^½ (el volumen) y FRECUENCIA γ (el tono). Una nota pura.")
	fmt.Println("\n      la primera perla γ=14.134725, escuchada a distintas escalas:")
	fmt.Println("        x            L = ln x     su nota          amplitud x^½/|ρ|")
	γ1 := 14.134725142
	mod1 := math.Hypot(0.5, γ1)
	for _, x := range []float64{2, 10, 100, 1000} {
		fmt.Printf("   %8.0f      %9.5f   %+12.6f      %12.6f\n",
			x, math.Log(x), unaNota(γ1, x, 0.5), 2*math.Sqrt(x)/mod1)
	}
	fmt.Println("   → la nota oscila entre ±2·x^½/|ρ| y nunca se sale de ahí: el volumen es x^½.")

	// ---- LEY 3: the wave with EVERYTHING ----
	fmt.Println("\nLEY 3 · Y ACÁ SUENA LA ONDA CON TODOS — la fórmula explícita de Riemann")
	fmt.Println("\n        ψ(x) = x − Σ sobre TODAS las perlas x^ρ/ρ − ln(2π) − ½·ln(1−x⁻²)")
	fmt.Println("\n   donde ψ(x) = Σ sobre las potencias de primo pᵏ ≤ x de ln p — LA ESCALERA DE")
	fmt.Println("   LOS PRIMOS. Sumás todas las notas y te devuelve los primos. Medido:")
	fmt.Println("\n   pescando perlas hasta t=500…")
	gam := perlas(500)
	fmt.Printf("   perlas (notas disponibles): %d\n", len(gam))
	es := escalones(200)

	fmt.Println("\n        x        ψ(x) verdadera     la onda con todas     desvío")
	peorOnda, sumOnda, nOnda := 0.0, 0.0, 0
	for _, x := range []float64{4.5, 9.5, 15.5, 23.5, 30.5, 44.5, 60.5, 89.5} {
		v := psiExacta(x, es)
		o := psiOnda(x, gam)
		d := math.Abs(v - o)
		if d > peorOnda {
			peorOnda = d
		}
		sumOnda += d
		nOnda++
		fmt.Printf("   %7.1f      %14.6f     %14.6f     %9.6f\n", x, v, o, d)
	}
	fmt.Printf("   → con %d notas la escalera de los primos se reconstruye a |Δ| medio %.4f\n",
		len(gam), sumOnda/float64(nOnda))
	fmt.Println("     (se evaluó entre saltos, porque justo en cada salto la serie repica —")
	fmt.Println("      fenómeno de Gibbs, no error del método)")
	fmt.Println("\n   ⚙️ EL PRIMO SALE DE LOS CEROS. No hay ninguna otra información acá adentro:")
	fmt.Println("   la fórmula solo sabe las alturas γ, y devuelve dónde están los primos.")

	// ---- LEY 4: more notes, better music ----
	fmt.Println("\nLEY 4 · Y CUANTAS MÁS NOTAS, MEJOR SUENA — se mide")
	fmt.Println("\n        notas usadas      |Δ| medio contra la escalera verdadera")
	for _, k := range []int{5, 20, 50, 100, len(gam)} {
		if k > len(gam) {
			continue
		}
		acc, n := 0.0, 0
		for _, x := range []float64{4.5, 9.5, 15.5, 23.5, 30.5, 44.5, 60.5, 89.5} {
			acc += math.Abs(psiExacta(x, es) - psiOnda(x, gam[:k]))
			n++
		}
		fmt.Printf("   %14d      %28.6f\n", k, acc/float64(n))
	}
	fmt.Println("   → cada nota que agregás afina la escalera. Es una orquesta, y le faltan")
	fmt.Println("     infinitos músicos: por eso nunca cierra del todo con una lista finita.")

	// ---- LEY 5: the half IS the volume ----
	fmt.Println("\nLEY 5 · Y EL ½ ES EL VOLUMEN DE TODAS LAS NOTAS — ahí está la hipótesis")
	fmt.Println("\n        RH  ⟺  TODAS las notas suenan al MISMO volumen, x^½")
	fmt.Println("\n   una perla corrida a β > ½ cantaría a x^β — MÁS FUERTE — y taparía a las demás.")
	fmt.Println("\n        x            nota honesta (β=½)     nota corrida (β=0.7)     cuántas veces más fuerte")
	for _, x := range []float64{10, 1e3, 1e6, 1e12, 1e24} {
		a := 2 * math.Pow(x, 0.5) / math.Hypot(0.5, γ1)
		b := 2 * math.Pow(x, 0.7) / math.Hypot(0.7, γ1)
		fmt.Printf("   %8.0e      %18.4e     %18.4e     %18.1f×\n", x, a, b, b/a)
	}
	fmt.Println("   → a escala 10²⁴ la nota corrida suena diez mil millones de veces más fuerte.")
	fmt.Println("     La hipótesis es, exactamente, QUE ESO NO PASA: que la orquesta está afinada")
	fmt.Println("     y ningún músico toca más fuerte que los otros.")

	// ---- LEY 6: and that IS the prime number theorem's error ----
	fmt.Println("\nLEY 6 · Y ESO ES, LETRA POR LETRA, EL ERROR DEL TEOREMA DE LOS PRIMOS")
	fmt.Println("\n        RH  ⟺  ψ(x) = x + O(x^{½+ε})")
	fmt.Println("\n   porque si todas las notas tienen volumen x^½, la suma de todas es del orden")
	fmt.Println("   de x^½ (por un pelo más). Medido sobre la escalera verdadera:")
	esGrande := escalones(300000)
	fmt.Println("\n        x            ψ(x) − x          |ψ(x)−x| / √x")
	peorRatio := 0.0
	for _, x := range []float64{1e2, 1e3, 1e4, 1e5, 3e5} {
		v := psiExacta(x, esGrande)
		r := math.Abs(v-x) / math.Sqrt(x)
		if r > peorRatio {
			peorRatio = r
		}
		fmt.Printf("   %9.0e      %14.4f      %14.6f\n", x, v-x, r)
	}
	fmt.Printf("   → el cociente se mantiene chico (peor %.4f en el rango medido): la orquesta\n", peorRatio)
	fmt.Println("     suena afinada hasta donde llegamos a escuchar. Y hasta donde llegó cualquiera.")
	fmt.Println("\n   ⚠ PERO ESO NO PRUEBA NADA, y es importante decirlo: que el cociente sea chico")
	fmt.Println("     hasta 3×10⁵ —o hasta 10²⁵, como midieron otros— no dice nada de lo que pasa")
	fmt.Println("     más arriba. Una sola nota desafinada, a cualquier altura, rompe todo.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL CAPITÁN CERRÓ SU PROPIO DIBUJO, Y LO CERRÓ BIEN:")
	fmt.Println("  · con cuatro direcciones (N/S/E/O) el dibujo está QUIETO — es la cruz de F245")
	fmt.Println("  · la quinta y la sexta son ARRIBA y ABAJO: x, la escala a la que escuchás")
	fmt.Println("  · y con las seis SUENA: cada perla es una nota de frecuencia γ y volumen x^½")
	fmt.Printf("  · sumando las %d notas, la escalera de los primos se reconstruye a |Δ| %.4f\n",
		len(gam), sumOnda/float64(nOnda))
	fmt.Println("  · y el ½ resulta ser EL VOLUMEN DE TODAS LAS NOTAS")
	fmt.Println("\n        RH  ⟺  LA ORQUESTA ESTÁ AFINADA: ningún músico toca más fuerte.")
	fmt.Println("\nAsí que las seis direcciones del capitán son, una por una:")
	fmt.Println("        este/oeste = el volumen · norte/sur = el tono · arriba/abajo = el tiempo")
	fmt.Println("y la hipótesis vive entera en la primera pareja, medida contra la tercera.")
	fmt.Println("\n⚖️ LO QUE ES NUESTRO Y LO QUE NO. La fórmula explícita es de Riemann (1859) y")
	fmt.Println("von Mangoldt (1895); la equivalencia con el error del teorema de los primos es")
	fmt.Println("clásica. Lo que hizo el taller hoy fue MEDIRLA: reconstruir la escalera desde")
	fmt.Println("las perlas de verdad y medir cuánto gritaría una nota corrida. Nada de eso es")
	fmt.Println("evidencia a favor de la hipótesis — es entenderla. ¿El premio? Todavía no.")

	escribirLamina(gam, es, peorOnda, sumOnda/float64(nOnda), peorRatio, γ1)
}

func escribirLamina(gam []float64, es []escalon, peorOnda, medioOnda, peorRatio, γ1 float64) {
	var b strings.Builder
	W, H := 1520.0, 1140.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧭 LAS SEIS DIRECCIONES — y ahí sí arma la onda con TODOS</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">la tercera dirección es x, la escala a la que escuchás — y con ella cada perla es UNA NOTA</text>
`, W, H, W, H, W/2, W/2)

	// the six directions
	fmt.Fprintf(&b, `<rect x="40" y="98" width="1440" height="150" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="66" y="128" font-size="16" font-family="Georgia" fill="#ffd98a">ESTE / OESTE   σ = Re s</text>
<text x="380" y="128" font-size="14" font-family="Georgia" fill="#cfe6ff">de qué lado de la línea está la perla  →  EL VOLUMEN de su nota, amplitud x^σ</text>
<text x="66" y="166" font-size="16" font-family="Georgia" fill="#ffd98a">NORTE / SUR    γ = Im s</text>
<text x="380" y="166" font-size="14" font-family="Georgia" fill="#cfe6ff">qué tan alto está sobre la línea  →  EL TONO de su nota, frecuencia γ</text>
<text x="66" y="204" font-size="16" font-family="Georgia" fill="#c9b6ff">ARRIBA / ABAJO   x</text>
<text x="380" y="204" font-size="14" font-family="Georgia" fill="#c9b6ff">la escala a la que escuchás  →  EL TIEMPO en que suena, a través de L = ln x</text>
<text x="760" y="234" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">con las cuatro primeras el dibujo está QUIETO. Con las seis, SUENA.</text>
`)

	// the staircase and the reconstruction
	px0, py0, pw, ph := 70.0, 300.0, 1380.0, 330.0
	xa, xb := 2.0, 90.0
	ymax := 95.0
	fmt.Fprintf(&b, `<rect x="40" y="268" width="1440" height="392" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="760" y="294" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA ONDA CON TODOS — %d notas reconstruyendo la escalera de los primos</text>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2f5480"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2f5480"/>
`, len(gam), px0, py0+ph, px0+pw, py0+ph, px0, py0, px0, py0+ph)
	mapx := func(x float64) float64 { return px0 + pw*(x-xa)/(xb-xa) }
	mapy := func(y float64) float64 { return py0 + ph - ph*y/ymax }
	// true staircase
	var esc strings.Builder
	acc := 0.0
	fmt.Fprintf(&esc, "M%.1f,%.1f", mapx(xa), mapy(0))
	for _, e := range es {
		if e.x < xa || e.x > xb {
			continue
		}
		fmt.Fprintf(&esc, " L%.1f,%.1f L%.1f,%.1f", mapx(e.x), mapy(acc), mapx(e.x), mapy(acc+e.peso))
		acc += e.peso
	}
	fmt.Fprintf(&esc, " L%.1f,%.1f", mapx(xb), mapy(acc))
	fmt.Fprintf(&b, `<path d="%s" fill="none" stroke="#7ee0c0" stroke-width="2.6"/>`, esc.String())
	// reconstruction
	var rec strings.Builder
	for k := 0; k <= 1200; k++ {
		x := xa + (xb-xa)*float64(k)/1200
		y := psiOnda(x, gam)
		if y < 0 {
			y = 0
		}
		if y > ymax {
			y = ymax
		}
		if k == 0 {
			fmt.Fprintf(&rec, "M%.1f,%.1f", mapx(x), mapy(y))
		} else {
			fmt.Fprintf(&rec, " L%.1f,%.1f", mapx(x), mapy(y))
		}
	}
	fmt.Fprintf(&b, `<path d="%s" fill="none" stroke="#ffb27a" stroke-width="2" opacity="0.95"/>
<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#7ee0c0">— la escalera VERDADERA de los primos, ψ(x) = Σ ln p</text>
<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#ffb27a">— la reconstruida SOLO con las alturas de las perlas</text>
<text x="760" y="648" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">ψ(x) = x − Σ x^ρ/ρ − ln(2π) − ½·ln(1−x⁻²)   ·   |Δ| medio entre saltos: %.4f   (el repique en cada salto es Gibbs, no error)</text>
`, rec.String(), px0+30, py0+26, px0+30, py0+48, medioOnda)

	// the volume
	fmt.Fprintf(&b, `<rect x="40" y="678" width="700" height="260" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="390" y="708" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">Y EL ½ ES EL VOLUMEN DE TODAS LAS NOTAS</text>
<text x="390" y="746" font-size="18" text-anchor="middle" font-family="monospace" fill="#dce8f7">RH ⟺ todas suenan al mismo volumen x^½</text>
<text x="66" y="782" font-size="13" font-family="monospace" fill="#7fa8cf">   x          honesta (β=½)     corrida (β=0.7)    veces</text>
`)
	yy := 806.0
	for _, x := range []float64{10, 1e6, 1e12, 1e24} {
		a := 2 * math.Pow(x, 0.5) / math.Hypot(0.5, γ1)
		bb := 2 * math.Pow(x, 0.7) / math.Hypot(0.7, γ1)
		fmt.Fprintf(&b, `<text x="66" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">%8.0e   %14.3e   %14.3e   %9.0f×</text>`,
			yy, x, a, bb, bb/a)
		yy += 22
	}
	fmt.Fprintf(&b, `<text x="66" y="%.0f" font-size="13.5" font-family="Georgia" fill="#9fd8a8">una sola perla corrida gritaría por encima de toda la orquesta.</text>
<text x="66" y="%.0f" font-size="13.5" font-family="Georgia" fill="#ffd98a">la hipótesis es, exactamente, QUE ESO NO PASA.</text>
`, yy+10, yy+32)

	fmt.Fprintf(&b, `<rect x="780" y="678" width="700" height="260" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1130" y="708" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Y ESO ES EL ERROR DEL TEOREMA DE LOS PRIMOS</text>
<text x="1130" y="746" font-size="18" text-anchor="middle" font-family="monospace" fill="#dce8f7">RH ⟺ ψ(x) = x + O(x^{½+ε})</text>
<text x="806" y="782" font-size="13.5" font-family="Georgia" fill="#cfe6ff">si todas las notas tienen volumen x^½, la suma de todas</text>
<text x="806" y="802" font-size="13.5" font-family="Georgia" fill="#cfe6ff">es del orden de x^½ por un pelo más. Medido sobre la</text>
<text x="806" y="822" font-size="13.5" font-family="Georgia" fill="#cfe6ff">escalera verdadera hasta x = 3×10⁵:</text>
<text x="806" y="854" font-size="15" font-family="monospace" fill="#9fd8a8">   |ψ(x) − x| / √x  ≤  %.4f</text>
<text x="806" y="890" font-size="12.5" font-family="Georgia" fill="#ffb27a">⚠ y eso NO prueba nada: que suene afinada hasta acá —o hasta</text>
<text x="806" y="908" font-size="12.5" font-family="Georgia" fill="#ffb27a">10²⁵, como midieron otros— no dice nada de más arriba. Una sola</text>
<text x="806" y="926" font-size="12.5" font-family="Georgia" fill="#ffb27a">nota desafinada, a cualquier altura, rompe todo.</text>
`, peorRatio)

	fmt.Fprintf(&b, `<rect x="40" y="958" width="1440" height="148" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="760" y="992" font-size="22" text-anchor="middle" font-family="Georgia" fill="#ffd98a">RH ⟺ LA ORQUESTA ESTÁ AFINADA</text>
<text x="760" y="1024" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">este/oeste = el VOLUMEN  ·  norte/sur = el TONO  ·  arriba/abajo = el TIEMPO</text>
<text x="760" y="1052" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">con cuatro direcciones el dibujo está quieto; con las seis suena, y los primos salen de las alturas de las perlas y de nada más.</text>
<text x="760" y="1082" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ la fórmula explícita es de Riemann (1859) y von Mangoldt (1895): acá se MIDIÓ, no se descubrió. Nada de esto es evidencia. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("seis-direcciones.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: seis-direcciones.svg")
}
