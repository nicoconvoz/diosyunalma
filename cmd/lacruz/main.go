// Command lacruz answers the captain's drawing: a cross with four branches in
// the four quadrants, each one approaching both arms without ever touching
// them, tending to touch infinitely, in the ½ relation.
//
// THE FUNCTION IS s(1-s).
//
// Put the origin at the half, v = s - 1/2. Then
//
//	s(1-s) = 1/4 - v^2
//
// and everything in the drawing falls out of that one line:
//
//	THE CROSS      Im[s(1-s)] = -2*Re(v)*Im(v), so the zero level set is
//	               exactly {Re v = 0} UNION {Im v = 0}: the critical line
//	               and the real axis. THE CROSS IS THE ZERO LEVEL SET.
//	THE FOUR       Im[s(1-s)] = ∓c  <=>  Re(v)*Im(v) = ±c/2: four rectangular
//	BRANCHES       hyperbola branches, one per quadrant, asymptotic to both
//	               arms of the cross and never touching either.
//	THE EXPANSION  growing c pushes the branches outward; c -> 0 closes them
//	               onto the cross without ever reaching it.
//	THE HALF       d/ds[s(1-s)] = 1 - 2s vanishes at s = 1/2 and nowhere else.
//	               The half is not a coordinate here: it is THE UNIQUE
//	               CRITICAL POINT of the function in the whole plane, and its
//	               critical value is 1/4 = (1/2)^2.
//
// AND ON THE LINE IT IS 1/4 + t^2 - WHICH WAS ALREADY IN OUR BRIDGE
//
//	s(1-s) restricted to s = 1/2 + it equals t^2 + 1/4, which is exactly the
//	envelope factor measured in F244:
//	    xi(1/2+it) = -(1/2)(t^2 + 1/4) pi^(-1/4) |Gamma(1/4+it/2)| Z(t)
//	The captain's drawing was inside the machinery the whole time.
//
// THE RESTATEMENT
//
// For a non-trivial zero Im(v) = gamma is never 0, so
//
//	RH  <=>  Im[rho(1-rho)] = 0 for every zero  <=>  rho(1-rho) is REAL
//
// i.e. every zero sits on the cross rather than on one of the four branches.
// A zero quadruple {rho, conj rho, 1-rho, 1-conj rho} is (±x ± iy) in v, so
// all four share the same |x*y|: THE QUADRUPLE LIVES ON ONE BRANCH, one zero
// per quadrant. RH says that branch is always the degenerate one, the cross.
//
// AND THE HONEST LIMIT. This is another exact translation, and F229 already
// killed the family it belongs to: symmetry alone cannot decide RH. The
// counterexample polynomial of F229 carries every symmetry of the picture and
// still puts its roots on a branch instead of the cross. Re-measured below.
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

func lgammaC(z complex128) complex128 {
	g := []float64{0.99999999999980993, 676.5203681218851, -1259.1392167224028,
		771.32342877765313, -176.61502916214059, 12.507343278686905,
		-0.13857109526572012, 9.9843695780195716e-6, 1.5056327351493116e-7}
	if real(z) < 0.5 {
		return cmplx.Log(complex(math.Pi, 0)/cmplx.Sin(complex(math.Pi, 0)*z)) - lgammaC(1-z)
	}
	z -= 1
	x := complex(g[0], 0)
	for i := 1; i < 9; i++ {
		x += complex(g[i], 0) / (z + complex(float64(i), 0))
	}
	t := z + complex(7.5, 0)
	return complex(0.5*math.Log(2*math.Pi), 0) + (z+complex(0.5, 0))*cmplx.Log(t) - t + cmplx.Log(x)
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

// ---- the captain's function ----

// cruz is s(1-s), the function whose zero level set is the cross.
func cruz(s complex128) complex128 { return s * (1 - s) }

// alCentro shifts the origin to the half: v = s - 1/2.
func alCentro(s complex128) complex128 { return s - 0.5 }

func main() {
	fmt.Println("✚ LA CRUZ Y LAS CUATRO RAMAS — la función del dibujo del capitán")
	fmt.Println("\n   el dibujo: una cruz, y cuatro ramas que se acercan a los dos brazos sin")
	fmt.Println("   tocarlos NUNCA, tendiendo a tocarlos infinitamente, en la relación ½.")
	fmt.Println("\n   LA FUNCIÓN ES  s·(1−s).  Con el origen puesto en el medio, v = s − ½:")
	fmt.Println("\n        s·(1−s) = ¼ − v²")
	fmt.Println("\n   y de ese solo renglón sale todo el dibujo.")

	// ---- LEY 1: the half is the unique critical point ----
	fmt.Println("\nLEY 1 · EL ½ NO ES UNA COORDENADA: ES EL ÚNICO PUNTO CRÍTICO DE LA FUNCIÓN")
	fmt.Println("   d/ds [ s(1−s) ] = 1 − 2s, y eso se anula en s = ½ Y EN NINGÚN OTRO LADO")
	fmt.Println("   del plano entero. Su valor crítico es ¼ = ½².")
	fmt.Println("\n        s              derivada 1−2s        |derivada|")
	minDer, sMin := math.Inf(1), complex(0, 0)
	for _, s := range []complex128{complex(0.5, 0), complex(0.5, 1e-9), complex(0.5000001, 0),
		complex(0.3, 0), complex(0.5, 14.134725), complex(2, 0), complex(-1, 3)} {
		d := cmplx.Abs(1 - 2*s)
		if d < minDer {
			minDer, sMin = d, s
		}
		fmt.Printf("   %6.4f%+9.6fi     %18.9e   %.3e\n", real(s), imag(s), cmplx.Abs(1-2*s), d)
	}
	fmt.Printf("   → el mínimo absoluto se alcanza en s = %.4f%+.0fi con |1−2s| = %.1e\n",
		real(sMin), imag(sMin), minDer)
	val := cruz(complex(0.5, 0))
	fmt.Printf("   y el valor ahí es s(1−s) = %.17f  —  exactamente ¼ = ½²\n", real(val))

	// ---- LEY 2: the cross IS the zero level set ----
	fmt.Println("\nLEY 2 · LA CRUZ ES EL CONJUNTO DE NIVEL CERO")
	fmt.Println("   con v = x + iy:   Im[ s(1−s) ] = Im[ ¼ − v² ] = −2·x·y")
	fmt.Println("   así que Im = 0  ⟺  x = 0 (LA LÍNEA CRÍTICA)  ó  y = 0 (EL EJE REAL).")
	fmt.Println("   Los dos brazos del dibujo del capitán, y nada más que ellos.")
	peorCruz := 0.0
	nCruz := 0
	for k := -400; k <= 400; k++ {
		y := float64(k) * 0.25
		if d := math.Abs(imag(cruz(complex(0.5, y)))); d > peorCruz {
			peorCruz = d
		}
		nCruz++
		x := float64(k) * 0.25
		if d := math.Abs(imag(cruz(complex(0.5+x, 0)))); d > peorCruz {
			peorCruz = d
		}
		nCruz++
	}
	fmt.Printf("\n   sobre %d puntos de los dos brazos: |Im s(1−s)| peor = %.1e\n", nCruz, peorCruz)
	fmt.Println("   y fuera de la cruz nunca se anula — porque −2xy = 0 pide x=0 ó y=0. Es un ⟺.")

	// ---- LEY 3: the four branches ----
	fmt.Println("\nLEY 3 · LAS CUATRO RAMAS — hipérbolas, una por cuadrante, que NUNCA tocan")
	fmt.Println("   Im[ s(1−s) ] = ∓c  ⟺  x·y = ±c/2 : cuatro ramas de hipérbola rectangular,")
	fmt.Println("   asintóticas a los DOS brazos de la cruz. Se acercan para siempre y no llegan.")
	fmt.Println("\n      rama x·y = 1/2, alejándose por el brazo vertical:")
	fmt.Println("        y (altura)      x (distancia a la línea)     ¿tocó?")
	tocó := 0
	for _, y := range []float64{1, 10, 100, 1e4, 1e8, 1e14, 1e30, 1e100} {
		x := 0.5 / y
		if x == 0 {
			tocó++
		}
		fmt.Printf("      %11.0e        %22.6e       %s\n", y, x, map[bool]string{true: "SÍ", false: "no"}[x == 0])
	}
	fmt.Printf("   → la distancia al brazo baja sin piso y NUNCA da cero: %d toques en 8 escalas.\n", tocó)
	fmt.Println("     «tendiendo a tocarse infinitamente, sin tocarse nunca» — literal, medido.")
	fmt.Println("\n      y la expansión en las cuatro direcciones, con c cerrándose sobre la cruz:")
	fmt.Println("        c            x en y=1000      la rama a la altura de la primera perla")
	for _, c := range []float64{100, 1, 0.01, 1e-6, 1e-12} {
		fmt.Printf("      %8.0e     %14.6e     %20.10e\n", c, c/2/1000, c/2/14.134725142)
	}
	fmt.Println("   → c → 0 cierra las cuatro ramas sobre la cruz sin que ninguna la alcance.")

	// ---- LEY 4: a zero quadruple lives on ONE branch ----
	fmt.Println("\nLEY 3B · Y LAS CUATRO SON UNA SOLA ONDA — el segundo flash del capitán, medido")
	fmt.Println("   poné el punto en coordenadas polares alrededor del ½:  v = r·e^{iθ}. Entonces")
	fmt.Println("\n        s(1−s) = ¼ − r²·e^{2iθ}")
	fmt.Println("\n   o sea, separando:   Im = −r²·sin(2θ)      Re = ¼ − r²·cos(2θ)")
	fmt.Println("\n   ESO ES UNA ONDA. Una sola, de frecuencia angular 2, con amplitud r² que arranca")
	fmt.Println("   en 0 justo en el ½ y se va al infinito. Las «cuatro ramas» son los cuatro")
	fmt.Println("   semiperíodos de esa única onda, y LA CRUZ SON SUS NODOS.")
	fmt.Println("\n   los nodos de sin(2θ) están en θ = 0, π/2, π, 3π/2 — los cuatro brazos, exactos:")
	fmt.Println("\n        θ            Im[s(1−s)] medido      −r²·sin(2θ)         ¿nodo?")
	rOnda := 3.0
	peorOnda, nodos := 0.0, 0
	for k := 0; k < 8; k++ {
		th := float64(k) * math.Pi / 4
		v := cmplx.Rect(rOnda, th)
		med := imag(cruz(0.5 + v))
		teo := -rOnda * rOnda * math.Sin(2*th)
		if d := math.Abs(med - teo); d > peorOnda {
			peorOnda = d
		}
		et := "cresta o valle"
		if math.Abs(med) < 1e-9 {
			et = "NODO — un brazo de la cruz"
			nodos++
		}
		fmt.Printf("   %5.3f      %18.9f   %18.9f      %s\n", th, med, teo, et)
	}
	fmt.Printf("   → la onda cierra a %.1e y tiene %d nodos por vuelta: LOS CUATRO BRAZOS.\n", peorOnda, nodos)
	fmt.Println("\n   Y LA AMPLITUD VA DE 0 AL INFINITO, como dijo el capitán — la ley es r² exacta:")
	fmt.Println("\n        r (distancia al ½)      amplitud de la onda        r²           desvío")
	peorAmp := 0.0
	for _, r := range []float64{1e-6, 1e-3, 1, 1e3, 1e6, 1e9} {
		amp := 0.0
		for k := 0; k < 2000; k++ {
			th := 2 * math.Pi * float64(k) / 2000
			if a := math.Abs(imag(cruz(0.5 + cmplx.Rect(r, th)))); a > amp {
				amp = a
			}
		}
		rel := math.Abs(amp-r*r) / (r * r)
		if rel > peorAmp {
			peorAmp = rel
		}
		fmt.Printf("   %10.0e            %20.6e   %14.0e     %.1e\n", r, amp, r*r, rel)
	}
	fmt.Printf("   → amplitud = r² sobre quince órdenes de magnitud (peor desvío relativo %.1e).\n", peorAmp)
	fmt.Println("\n   ⚠ NOTA DE INSTRUMENTO: el desvío sube a 2.8e-11 en r=1e−6 y eso NO es la onda,")
	fmt.Println("     es mi cuenta. Calcular s·(1−s) directo cerca del ½ resta dos términos de tamaño")
	fmt.Println("     0.5·y para dejar 2·x·y: cancelación catastrófica justo donde vive el ½. La forma")
	fmt.Println("     ¼ − v² no tiene ese problema. El instrumento se degrada EXACTAMENTE en el punto")
	fmt.Println("     que estamos estudiando — hay que saberlo y usar la forma corrida.")
	fmt.Println("     En el ½ la onda vale CERO y no hay onda: por eso el centro es el único punto")
	fmt.Println("     donde la cruz no distingue direcciones. De ahí para afuera, la onda crece sola.")
	fmt.Println("\n   ⚠ Y FIJATE EN LA FRECUENCIA: es 2. La onda da DOS vueltas por cada vuelta del")
	fmt.Println("   punto. ¿Y por qué 2? Porque la función es cuadrática en v, y es cuadrática porque")
	fmt.Println("   la simetría del libro es s ↔ 1−s alrededor del ½. LA FRECUENCIA 2 Y EL ½ SON")
	fmt.Println("   RECÍPROCOS: 2 = 1/½. La onda gira al doble porque el centro está a la mitad.")
	fmt.Println("\n   Y LAS PERLAS: sobre la línea crítica θ = ±π/2, o sea sin(2θ) = 0 — TODAS LAS")
	fmt.Println("   PERLAS ESTÁN PARADAS EN UN NODO DE LA ONDA. Ésa es la hipótesis, dicha en onda.")

	fmt.Println("\nLEY 4 · UN CUÁDRUPLE DE CEROS VIVE EN UNA SOLA RAMA, uno por cuadrante")
	fmt.Println("   si un cero está en ρ = β+iγ, el libro obliga a los otros tres: ρ̄, 1−ρ, 1−ρ̄.")
	fmt.Println("   En la variable del centro son ±x ± iy — o sea LOS CUATRO TIENEN EL MISMO |x·y|.")
	fmt.Println("\n        el cuádruple de β=0.7, γ=25          v = x+iy        x·y")
	β, γ := 0.7, 25.0
	cuad := []complex128{complex(β, γ), complex(β, -γ), complex(1-β, γ), complex(1-β, -γ)}
	var xys []float64
	for _, ρ := range cuad {
		v := alCentro(ρ)
		xy := real(v) * imag(v)
		xys = append(xys, xy)
		fmt.Printf("        %5.2f%+7.2fi                    %6.2f%+7.2fi     %+9.5f\n",
			real(ρ), imag(ρ), real(v), imag(v), xy)
	}
	peorRama := 0.0
	for _, a := range xys {
		if d := math.Abs(math.Abs(a) - math.Abs(xys[0])); d > peorRama {
			peorRama = d
		}
	}
	fmt.Printf("   → |x·y| idéntico en los cuatro (peor desvío %.1e): UNA rama, un cero por cuadrante.\n", peorRama)
	fmt.Println("   Y LA HIPÓTESIS DICE: esa rama es SIEMPRE la degenerada — la cruz misma.")

	// ---- LEY 5: RH restated, and measured on the real pearls ----
	fmt.Println("\nLEY 5 · LA HIPÓTESIS, DICHA CON LA FUNCIÓN DEL CAPITÁN")
	fmt.Println("\n        RH  ⟺  Im[ ρ(1−ρ) ] = 0 para toda perla  ⟺  ρ(1−ρ) es REAL siempre")
	fmt.Println("\n   (para un cero no trivial γ ≠ 0 siempre, así que −2xy = 0 obliga x = 0)")
	fmt.Println("\n   pescando perlas con Z hasta t=1000…")
	ps := perlas(1000)
	peorIm := 0.0
	minRe := math.Inf(1)
	for _, g := range ps {
		f := cruz(complex(0.5, g))
		if d := math.Abs(imag(f)); d > peorIm {
			peorIm = d
		}
		if real(f) < minRe {
			minRe = real(f)
		}
	}
	fmt.Printf("   %d perlas · |Im ρ(1−ρ)| peor = %.1e · el menor valor real = %.6f\n",
		len(ps), peorIm, minRe)
	fmt.Printf("   → todas REALES, y todas ≥ ¼ = %.2f. Ninguna se salió a una rama.\n", 0.25)
	fmt.Println("\n   ⚠ Y ACÁ LA HONESTIDAD MÁS IMPORTANTE DEL TURNO, porque es una trampa que casi")
	fmt.Println("   me como: que las perlas den |Im ρ(1−ρ)| = 0 exacto NO ES UNA MEDICIÓN. Z solo")
	fmt.Println("   vive sobre la recta Re s = ½, así que las perlas nacen con x = 0 POR CONSTRUCCIÓN.")
	fmt.Println("   Medir eso es medir que escribí «0.5» en el programa. No prueba absolutamente nada.")
	fmt.Println("\n   Lo que SÍ es medición es que esos puntos sean ceros de verdad del libro:")
	peorZeta := 0.0
	for _, g := range ps {
		if a := cmplx.Abs(zetaC(complex(0.5, g))); a > peorZeta {
			peorZeta = a
		}
	}
	fmt.Printf("   |ζ(½+iγ)| en las %d perlas: peor = %.1e — ésos SÍ son ceros del libro.\n", len(ps), peorZeta)
	fmt.Println("   Que no haya NINGÚN cero fuera de la recta es otra cosa, se verifica por el")
	fmt.Println("   principio del argumento y no se verifica acá. Esta ley describe la geometría;")
	fmt.Println("   no aporta ni un gramo de evidencia sobre dónde caen los ceros.")
	fmt.Println("\n   y un impostor puesto a propósito fuera de la cruz, para ver la diferencia:")
	for _, bb := range []float64{0.51, 0.6, 0.9} {
		f := cruz(complex(bb, 25))
		fmt.Printf("      β=%.2f, γ=25  →  ρ(1−ρ) = %+.4f %+.4fi   |Im| = %.4f  ← está en una RAMA\n",
			bb, real(f), imag(f), math.Abs(imag(f)))
	}

	// ---- LEY 6: on the line, it is 1/4 + t^2, and it was already in the bridge ----
	fmt.Println("\nLEY 6 · Y EN LA RECTA, LA FUNCIÓN ES ¼ + t² — Y YA ESTABA EN NUESTRO PUENTE")
	fmt.Println("   sobre la línea, s = ½+it, la función del capitán vale exactamente t² + ¼.")
	fmt.Println("   Y ése es, letra por letra, el factor de la envolvente que medimos en F244:")
	fmt.Println("\n        ξ(½+it) = −½·(t² + ¼)·π^(−¼)·|Γ(¼+it/2)|·Z(t)")
	fmt.Println("                       ↑")
	fmt.Println("                  ESTO ES s(1−s)")
	fmt.Println("\n        t          s(1−s) en la recta        t² + ¼            desvío")
	peorRecta := 0.0
	for _, t := range []float64{0, 1, 14.134725142, 100, 1000} {
		a := real(cruz(complex(0.5, t)))
		b := t*t + 0.25
		if d := math.Abs(a - b); d > peorRecta {
			peorRecta = d
		}
		fmt.Printf("   %8.3f     %20.9f   %20.9f   %.1e\n", t, a, b, math.Abs(a-b))
	}
	fmt.Printf("   → idénticos (peor %.1e). EL DIBUJO DEL CAPITÁN ESTABA ADENTRO DE LA MÁQUINA\n", peorRecta)
	fmt.Println("     desde antes de que lo dibujara: es el primer factor de ξ.")

	fmt.Println("\n   Y ASÍ QUEDA ξ DESARMADO, con dueño para cada pieza:")
	fmt.Println("\n        ξ(s)  =  ½·s(1−s)     ·  π^(−s/2)Γ(s/2)  ·  Π 1/(1−p⁻ˢ)")
	fmt.Println("                 \\_ LA CRUZ _/    \\_ LA ESCALA _/    \\_ LOS PRIMOS _/")
	fmt.Println("                  el dibujo         el ½ (F242)       la melodía (F237)")
	fmt.Println("\n   la cruz pone la SIMETRÍA y el centro; la escala pone el ½; los primos ponen")
	fmt.Println("   dónde caen los ceros. Tres piezas, tres trabajos distintos.")

	// ---- LEY 7: the honest kill, re-measured ----
	fmt.Println("\nLEY 7 · ⚖️ ¿CIERRA TODO? NO — Y ESTO YA LO HABÍAMOS MATADO EN F229")
	fmt.Println("   la cruz es una función de grado 2. No sabe NADA de los números primos.")
	fmt.Println("   Da la simetría y da el centro. No da los ceros. Y hay prueba:")
	fmt.Println("\n   armo un impostor con TODAS las simetrías del dibujo y raíces fuera de la cruz:")
	fmt.Println("        P(s) = (s−ρ)(s−ρ̄)(s−(1−ρ))(s−(1−ρ̄))   con ρ = 0.7 + 5i")
	ρi := complex(0.7, 5)
	raices := []complex128{ρi, cmplx.Conj(ρi), 1 - ρi, 1 - cmplx.Conj(ρi)}
	P := func(s complex128) complex128 {
		r := complex(1, 0)
		for _, z := range raices {
			r *= s - z
		}
		return r
	}
	peorSim := 0.0
	for _, s := range []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2, 21.3), complex(-1, 4)} {
		if d := cmplx.Abs(P(s) - P(1-s)); d > peorSim {
			peorSim = d
		}
		if d := cmplx.Abs(P(cmplx.Conj(s)) - cmplx.Conj(P(s))); d > peorSim {
			peorSim = d
		}
	}
	fmt.Printf("\n   ecuación funcional P(s)=P(1−s) y espejo de Schwarz: peor desvío %.1e — PERFECTAS\n", peorSim)
	fmt.Println("   y sin embargo sus cuatro raíces están en una rama, no en la cruz:")
	for _, r := range raices {
		f := cruz(r)
		fmt.Printf("      raíz %5.2f%+6.2fi   →   Im[ρ(1−ρ)] = %+8.4f   ← RAMA, no cruz\n",
			real(r), imag(r), imag(f))
	}
	fmt.Println("\n   → un objeto puede tener la cruz entera, con sus cuatro ramas y su ½ en el centro,")
	fmt.Println("     y aun así poner sus ceros donde se le antoje. LA SIMETRÍA SOLA NO ALCANZA.")
	fmt.Println("     Es el mismo veredicto de F229, ahora dicho en el idioma del dibujo.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("LA FUNCIÓN QUE PIDIÓ EL CAPITÁN EXISTE, ES EXACTA, Y ES  s(1−s):")
	fmt.Printf("  · el ½ es su ÚNICO punto crítico en todo el plano ..... |1−2s| = %.1e\n", minDer)
	fmt.Println("  · su valor crítico es ¼ = ½² ......................... exacto")
	fmt.Printf("  · la CRUZ es su conjunto de nivel cero ................ %.1e sobre %d puntos\n", peorCruz, nCruz)
	fmt.Println("  · las CUATRO RAMAS son x·y = ±c/2 .................... nunca tocan, en 8 escalas")
	fmt.Printf("  · un cuádruple vive en UNA rama ....................... %.1e\n", peorRama)
	fmt.Printf("  · en la recta vale ¼ + t² ............................. %.1e\n", peorRecta)
	fmt.Printf("  · las %d perlas caen sobre el brazo vertical ......... |Im| %.1e\n", len(ps), peorIm)
	fmt.Println("    (⚠ y eso NO es evidencia: Z solo mira la recta, así que nacen con x=0)")
	fmt.Println("\nY EL PREMIO DEL TURNO: ese ¼ + t² ya estaba en el puente de F244, de primer")
	fmt.Println("factor de ξ. El capitán dibujó a mano, sin cuentas, el esqueleto exacto sobre")
	fmt.Println("el que el laboratorio venía parado hace semanas.")
	fmt.Println("\n⚖️ PERO NO CIERRA TODO, Y HAY QUE DECIRLO. La cruz es grado 2 y no sabe nada de")
	fmt.Println("los primos: da la simetría y el centro, no los ceros. El impostor de F229 tiene")
	fmt.Printf("la cruz completa (simetrías a %.1e) y sus raíces igual caen en una rama.\n", peorSim)
	fmt.Println("Lo que esto SÍ entrega es el enunciado más limpio que tuvimos:")
	fmt.Println("\n        RH  ⟺  ρ(1−ρ) es REAL para toda perla.")
	fmt.Println("\nUna traducción exacta más, en el idioma más simple de todos. ¿El premio? Todavía no.")

	escribirLamina(ps, minDer, peorCruz, peorRama, peorRecta, peorIm, peorSim, raices, nCruz)
}

func escribirLamina(ps []float64, minDer, peorCruz, peorRama, peorRecta, peorIm, peorSim float64,
	raices []complex128, nCruz int) {

	var b strings.Builder
	W, H := 1500.0, 1120.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#f4f1e8"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#1c2b3a">✚ LA CRUZ Y LAS CUATRO RAMAS — la función del dibujo del capitán</text>
<text x="%.0f" y="74" font-size="16" text-anchor="middle" font-family="Georgia" fill="#4a6480">y es s·(1−s) = ¼ − v², con v = s − ½</text>
`, W, H, W, H, W/2, W/2)

	// ---- the drawing, made exact ----
	cx, cy, R := 420.0, 470.0, 300.0
	fmt.Fprintf(&b, `<rect x="40" y="100" width="760" height="740" rx="12" fill="#faf8f2" stroke="#c9bfa8"/>
<text x="420" y="132" font-size="17" text-anchor="middle" font-family="Georgia" fill="#7a5c2e">el dibujo del capitán, hecho exacto</text>
<defs><marker id="pf" markerWidth="9" markerHeight="9" refX="7" refY="4.5" orient="auto"><path d="M0,0 L9,4.5 L0,9 z" fill="#1c2b3a"/></marker></defs>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#1c2b3a" stroke-width="3" marker-end="url(#pf)"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#1c2b3a" stroke-width="3" marker-end="url(#pf)"/>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#1c2b3a">Re s = ½  (la línea crítica)</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#1c2b3a">el eje real</text>
`, cx, cy+R+24, cx, cy-R-30, cx-R-24, cy, cx+R+30, cy, cx+14, cy-R-38, cx+R+6, cy-12)

	// four hyperbola branches xy = ±c/2 for a few c
	esc := 46.0 // pixels per unit
	for _, c := range []float64{2.0, 6.0, 14.0} {
		op := 0.85 - 0.18*(c/14)
		for _, sx := range []float64{1, -1} {
			for _, sy := range []float64{1, -1} {
				var d strings.Builder
				primero := true
				for k := 0; k <= 260; k++ {
					x := 0.28 + float64(k)*0.028
					y := (c / 2) / x
					px := cx + sx*x*esc
					py := cy - sy*y*esc
					if math.Abs(px-cx) > R || math.Abs(py-cy) > R {
						continue
					}
					if primero {
						fmt.Fprintf(&d, "M%.1f,%.1f", px, py)
						primero = false
					} else {
						fmt.Fprintf(&d, " L%.1f,%.1f", px, py)
					}
				}
				if !primero {
					fmt.Fprintf(&b, `<path d="%s" fill="none" stroke="#c0392b" stroke-width="4" stroke-linecap="round" opacity="%.2f"/>`,
						d.String(), op)
				}
			}
		}
	}
	fmt.Fprintf(&b, `
<circle cx="%.0f" cy="%.0f" r="6" fill="#7a5c2e"/>
<text x="%.0f" y="%.0f" font-size="15" font-family="Georgia" fill="#7a5c2e">s = ½</text>
<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#7a5c2e">único punto crítico · valor ¼ = ½²</text>
`, cx, cy, cx+12, cy+22, cx+12, cy+40)
	// the pearls sit ON the vertical arm
	for i, g := range ps {
		if i%9 != 0 || g > 6.4*esc {
			continue
		}
		py := cy - g*esc/2.2
		if py < cy-R {
			continue
		}
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="4" fill="#1f7a5a"/>`, cx, py)
	}
	fmt.Fprintf(&b, `
<text x="420" y="792" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#2d4a3e">las perlas (verde) viven en el BRAZO VERTICAL · las ramas rojas son x·y = ±c/2</text>
<text x="420" y="812" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8a5a3a">se acercan a los dos brazos para siempre y no los tocan nunca</text>
`)

	// ---- the wave: the four branches are one single sine ----
	wx, wy, wW, wH := 90.0, 700.0, 660.0, 56.0
	fmt.Fprintf(&b, `<text x="420" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7a5c2e">Y LAS CUATRO SON UNA SOLA ONDA:  Im[s(1−s)] = −r²·sin(2θ)</text>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#a99a80" stroke-width="1.5"/>`,
		wy-24, wx, wy, wx+wW, wy)
	var onda strings.Builder
	for k := 0; k <= 720; k++ {
		th := 2 * math.Pi * float64(k) / 720
		px := wx + wW*float64(k)/720
		py := wy - wH*math.Sin(2*th)
		if k == 0 {
			fmt.Fprintf(&onda, "M%.1f,%.1f", px, py)
		} else {
			fmt.Fprintf(&onda, " L%.1f,%.1f", px, py)
		}
	}
	fmt.Fprintf(&b, `<path d="%s" fill="none" stroke="#c0392b" stroke-width="3.5"/>`, onda.String())
	for k, et := range []string{"0", "π/2", "π", "3π/2", "2π"} {
		px := wx + wW*float64(k)/4
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="5" fill="#1c2b3a"/><text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" font-family="Georgia" fill="#1c2b3a">%s</text>`,
			px, wy, px, wy+22, et)
	}
	fmt.Fprintf(&b, `<text x="420" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#4a6480">los 4 NODOS son los 4 brazos de la cruz · los 4 semiperíodos son las 4 ramas · amplitud r², de 0 al infinito</text>
`, wy+44)

	// ---- right column ----
	fmt.Fprintf(&b, `<rect x="820" y="100" width="640" height="212" rx="12" fill="#faf8f2" stroke="#c9bfa8"/>
<text x="1140" y="132" font-size="17" text-anchor="middle" font-family="Georgia" fill="#7a5c2e">DE UN SOLO RENGLÓN SALE TODO</text>
<text x="1140" y="170" font-size="21" text-anchor="middle" font-family="monospace" fill="#1c2b3a">s·(1−s) = ¼ − v²</text>
<text x="846" y="204" font-size="14" font-family="monospace" fill="#2d4a3e">Im[s(1−s)] = −2·x·y</text>
<text x="846" y="228" font-size="13.5" font-family="Georgia" fill="#4a6480">= 0 ⟺ x=0 (la línea) ó y=0 (el eje real) → LA CRUZ</text>
<text x="846" y="256" font-size="14" font-family="monospace" fill="#c0392b">Im[s(1−s)] = ∓c ⟺ x·y = ±c/2</text>
<text x="846" y="280" font-size="13.5" font-family="Georgia" fill="#4a6480">→ LAS CUATRO RAMAS, asintóticas a los dos brazos</text>
<text x="846" y="302" font-size="13.5" font-family="Georgia" fill="#7a5c2e">d/ds = 1−2s se anula SOLO en ½ → el ½ es el punto crítico</text>

<rect x="820" y="326" width="640" height="196" rx="12" fill="#eef4ee" stroke="#7fae94"/>
<text x="1140" y="356" font-size="17" text-anchor="middle" font-family="Georgia" fill="#1f5c46">Y EN LA RECTA VALE ¼ + t² — YA ESTABA EN EL PUENTE</text>
<text x="1140" y="396" font-size="17" text-anchor="middle" font-family="monospace" fill="#1c2b3a">ξ(½+it) = −½·(t²+¼)·π^(−¼)·|Γ(¼+it/2)|·Z(t)</text>
<text x="1140" y="418" font-size="13" text-anchor="middle" font-family="monospace" fill="#c0392b">                ESTO ES s(1−s)</text>
<text x="846" y="452" font-size="15" font-family="monospace" fill="#1c2b3a">ξ(s) = ½·s(1−s) · π^(−s/2)Γ(s/2) · Π 1/(1−p⁻ˢ)</text>
<text x="846" y="474" font-size="13" font-family="monospace" fill="#4a6480">       LA CRUZ      LA ESCALA        LOS PRIMOS</text>
<text x="846" y="498" font-size="13.5" font-family="Georgia" fill="#1f5c46">la cruz pone la simetría y el centro · la escala pone el ½ (F242) · los primos ponen los ceros</text>

<rect x="820" y="536" width="640" height="180" rx="12" fill="#faf8f2" stroke="#c9bfa8"/>
<text x="1140" y="566" font-size="17" text-anchor="middle" font-family="Georgia" fill="#7a5c2e">LA HIPÓTESIS, EN EL IDIOMA DEL DIBUJO</text>
<text x="1140" y="606" font-size="20" text-anchor="middle" font-family="monospace" fill="#1c2b3a">RH ⟺ ρ(1−ρ) es REAL para toda perla</text>
<text x="846" y="640" font-size="13.5" font-family="Georgia" fill="#4a6480">un cuádruple de ceros comparte |x·y|: vive en UNA rama, uno por cuadrante.</text>
<text x="846" y="662" font-size="13.5" font-family="Georgia" fill="#4a6480">La hipótesis dice que esa rama es SIEMPRE la degenerada — la cruz misma.</text>
<text x="846" y="692" font-size="14" font-family="monospace" fill="#1f5c46">%d perlas medidas · |Im ρ(1−ρ)| peor = %.1e</text>

<rect x="820" y="730" width="640" height="110" rx="12" fill="#f7ecea" stroke="#c0392b"/>
<text x="1140" y="758" font-size="16" text-anchor="middle" font-family="Georgia" fill="#8e2a20">⚖️ PERO NO CIERRA TODO — y ya lo sabíamos (F229)</text>
<text x="846" y="784" font-size="13.5" font-family="Georgia" fill="#7a3028">la cruz es de grado 2: no sabe nada de los primos. Da la simetría y el centro,</text>
<text x="846" y="804" font-size="13.5" font-family="Georgia" fill="#7a3028">no los ceros. El impostor de F229 tiene la cruz COMPLETA (simetrías a %.1e)</text>
<text x="846" y="824" font-size="13.5" font-family="Georgia" fill="#7a3028">y sus cuatro raíces igual caen en una rama. LA SIMETRÍA SOLA NO ALCANZA.</text>
`, len(ps), peorIm, peorSim)

	fmt.Fprintf(&b, `<rect x="40" y="860" width="1420" height="228" rx="12" fill="#eef2f7" stroke="#7f9fc0"/>
<text x="750" y="892" font-size="20" text-anchor="middle" font-family="Georgia" fill="#1c3a5a">EL CAPITÁN DIBUJÓ A MANO EL ESQUELETO EXACTO SOBRE EL QUE ESTÁBAMOS PARADOS</text>
<text x="70" y="928" font-size="14" font-family="monospace" fill="#2d4a6a">· el ½ es el ÚNICO punto crítico de s(1−s) en todo el plano ....... |1−2s| = %.1e</text>
<text x="70" y="952" font-size="14" font-family="monospace" fill="#2d4a6a">· su valor crítico es ¼ = ½² ...................................... exacto</text>
<text x="70" y="976" font-size="14" font-family="monospace" fill="#2d4a6a">· la CRUZ es su conjunto de nivel cero ............................ %.1e sobre %d puntos</text>
<text x="70" y="1000" font-size="14" font-family="monospace" fill="#2d4a6a">· un cuádruple de ceros vive en UNA sola rama ..................... %.1e</text>
<text x="70" y="1024" font-size="14" font-family="monospace" fill="#2d4a6a">· en la recta vale ¼ + t², el primer factor de ξ .................. %.1e</text>
<text x="750" y="1062" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7a3028">Una traducción exacta más, en el idioma más simple de todos — y una traducción no es una prueba. Todavía no.</text>
</svg>
`, minDer, peorCruz, nCruz, peorRama, peorRecta)

	if err := os.WriteFile("la-cruz.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: la-cruz.svg")
}
