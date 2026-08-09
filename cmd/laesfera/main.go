// Command laesfera answers the captain's flash: "the hourglass in every
// position, which form a sphere - and infinite spheres touching at infinite
// points of the hourglasses."
//
// There is a real, exact geometry underneath, and it has a name.
//
// # THE PLANE IS A SPHERE
//
// Add one point at infinity to the plane and it closes into a sphere - the
// Riemann sphere, 1857. Under stereographic projection:
//
//	the disk of dimension 0  ->  the southern cap
//	THE SKIN, |z| = 1        ->  THE EQUATOR
//	the outside, |z| > 1     ->  the northern cap
//
// So the critical line is not a line at all: it is a CIRCLE, the equator of a
// sphere, and it only looks like a line because we are standing on it.
//
// # THE HOURGLASS IN EVERY POSITION
//
// On the sphere the pearls all sit on the equator, crowding toward one point of
// it - the clasp - like the waist of an hourglass. And a sphere can be turned:
// rotate it and the crowding point goes anywhere. That is exactly the captain's
// "the hourglass in every position", and the positions really do sweep out the
// whole sphere.
//
// # INFINITE SPHERES TOUCHING AT INFINITE POINTS
//
// That picture already exists in mathematics and it is called the FORD CIRCLES.
// For every fraction p/q there is a circle sitting on the real line at p/q with
// diameter 1/q^2. They never cross. They TOUCH - and two of them touch exactly
// when |p*s - q*r| = 1, which is the Farey-neighbour condition. Infinitely many
// circles, infinitely many contact points, no overlap anywhere.
//
// # AND WHERE IT LEADS, HONESTLY
//
// The Ford circles are the horocycles of the modular group, which is the group
// that produces modular forms - and the shop built one in F246, Ramanujan's
// Delta, whose L-function has its own centre at 6. So this really is the house
// the second book lives in. What it does NOT do is say where any zero sits.
// It is the floor plan, not the address.
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

func alDisco(s complex128) complex128 { return 1 - 1/s }

// aLaEsfera is stereographic projection: it wraps the plane onto a sphere of
// radius 1, sending the unit circle to the equator and infinity to the north
// pole.
func aLaEsfera(z complex128) [3]float64 {
	x, y := real(z), imag(z)
	d := 1 + x*x + y*y
	return [3]float64{2 * x / d, 2 * y / d, (x*x + y*y - 1) / d}
}

// alPlano undoes it.
func alPlano(p [3]float64) complex128 {
	return complex(p[0]/(1-p[2]), p[1]/(1-p[2]))
}

// girar turns the sphere about the x axis by the angle a.
func girar(p [3]float64, a float64) [3]float64 {
	c, s := math.Cos(a), math.Sin(a)
	return [3]float64{p[0], c*p[1] - s*p[2], s*p[1] + c*p[2]}
}

func norma(p [3]float64) float64 {
	return math.Sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])
}

func mcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

type frac struct{ p, q int }

func main() {
	fmt.Println("🌐 LA ESFERA — el reloj de arena en todas las posiciones, y las que se tocan")
	fmt.Println("\n   flash del capitán: «el reloj de arena en todas las posiciones que forman una")
	fmt.Println("   esfera, e infinitas esferas tocándose en infinitos puntos de los relojes».")
	fmt.Println("\n   Hay una geometría exacta abajo de eso, y tiene nombre.")

	fmt.Println("\npescando perlas hasta t=1000…")
	ps := perlas(1000)
	fmt.Printf("perlas: %d\n", len(ps))

	// ---- LEY 1: the plane is a sphere ----
	fmt.Println("\nLEY 1 · EL PLANO ES UNA ESFERA — y la línea crítica es su ECUADOR")
	fmt.Println("   agregale al plano un solo punto, el del infinito, y el plano se cierra: queda")
	fmt.Println("   una esfera. Es la esfera de Riemann, de 1857. Y bajo esa envoltura:")
	fmt.Println("\n        el disco de la dimensión 0  →  el casquete SUR")
	fmt.Println("        LA PIEL, |z| = 1            →  EL ECUADOR")
	fmt.Println("        el afuera, |z| > 1          →  el casquete NORTE")
	fmt.Println("\n   comprobado con las perlas: cada una vive en la piel, así que cada una tiene")
	fmt.Println("   que caer justo en el ecuador — o sea con altura CERO sobre la esfera:")
	fmt.Println("\n        γ (la altura)      altura en la esfera      ¿en el ecuador?")
	peorEc := 0.0
	for _, i := range []int{0, 4, 49, 299, len(ps) - 1} {
		g := ps[i]
		p := aLaEsfera(alDisco(complex(0.5, g)))
		if d := math.Abs(p[2]); d > peorEc {
			peorEc = d
		}
		fmt.Printf("   %14.6f     %20.16f      %s\n", g, p[2], map[bool]string{true: "sí", false: "NO"}[math.Abs(p[2]) < 1e-12])
	}
	for _, g := range ps {
		if d := math.Abs(aLaEsfera(alDisco(complex(0.5, g)))[2]); d > peorEc {
			peorEc = d
		}
	}
	fmt.Printf("   → las %d perlas en el ecuador, peor desvío %.1e\n", len(ps), peorEc)
	fmt.Println("   ⟹ LA LÍNEA CRÍTICA SE VUELVE UN CÍRCULO — pero con tres precisiones obligatorias,")
	fmt.Println("     que me trajeron los refutadores y sin las cuales la frase es poesía y no geometría:")
	fmt.Println("\n   ⚠ (1) LA PROYECCIÓN SE APLICA AL PLANO w, NO AL PLANO ρ. Si proyectaras la línea")
	fmt.Println("       crítica directo desde el plano del libro, saldría un círculo que pasa por el")
	fmt.Println("       POLO NORTE, no el ecuador. El ecuador aparece SOLO después del cambiaformas.")
	fmt.Println("       Y «parece recta porque estamos parados encima» es metáfora: la razón de verdad")
	fmt.Println("       es que miramos la coordenada ρ en vez de la w.")
	fmt.Println("\n   ⚠ (2) EL BROCHE w=1 NO ES IMAGEN DE NINGÚN ρ FINITO. Pedirlo daría −1 = 0. Así que")
	fmt.Println("       la imagen de la línea crítica de verdad es el ecuador MENOS UN PUNTO: un círculo")
	fmt.Println("       PINCHADO. w=1 es la imagen de ρ = ∞. Para que «la línea es un círculo» sea")
	fmt.Println("       literalmente cierto hay que agregarle el punto del infinito. Sin eso es falso")
	fmt.Println("       por un punto — y un punto alcanza.")
	peorBro := math.Inf(1)
	for _, g := range []float64{1e3, 1e6, 1e12, 1e30} {
		w := alDisco(complex(0.5, g))
		if d := cmplx.Abs(w - 1); d < peorBro {
			peorBro = d
		}
	}
	fmt.Printf("       medido: hasta γ=1e30 la perla queda a %.1e del broche, y nunca llega.\n", peorBro)
	fmt.Println("\n   ✅ (3) Y UNA QUE NO TENÍA, Y ES LA MEJOR: el ecuador no es un círculo cualquiera")
	fmt.Println("       elegido a dedo. El espejo del libro, ρ ↦ 1−ρ, se convierte en el disco en")
	fmt.Println("       w ↦ 1/w — y el conjunto que ese espejo deja quieto es EXACTAMENTE |w| = 1.")
	fmt.Println("       O sea: el ecuador es el conjunto invariante del espejo. Eso le da estatus.")
	fmt.Println("\n        ρ                w(1−ρ)                    1/w(ρ)              desvío")
	peorInv := 0.0
	for _, r := range []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2, 4.1), complex(0.5, 14.134725)} {
		a := alDisco(1 - r)
		c := 1 / alDisco(r)
		d := cmplx.Abs(a - c)
		if d > peorInv {
			peorInv = d
		}
		fmt.Printf("   %5.2f%+7.2fi   %9.6f%+9.6fi   %9.6f%+9.6fi    %.1e\n",
			real(r), imag(r), real(a), imag(a), real(c), imag(c), d)
	}
	fmt.Printf("   → el espejo del libro ES la inversión del disco (%.1e), y su conjunto quieto\n", peorInv)
	fmt.Println("     es el ecuador. Por eso el ecuador es EL círculo y no UN círculo.")
	fmt.Println("\n   ⚖️ Y LA ADVERTENCIA DE FONDO, que hay que decir fuerte: |w| = 1 ⟺ Re ρ = ½ es una")
	fmt.Println("   reformulación TRIVIAL — es la equidistancia de 0 y 1, nada más. «Todas las perlas")
	fmt.Println("   tienen |w| = 1» y la Hipótesis de Riemann son LA MISMA FRASE ESCRITA DOS VECES.")
	fmt.Println("   La esfera es un escenario hermoso. No es un teorema.")

	fmt.Println("\nLEY 2 · EL RELOJ DE ARENA EN TODAS LAS POSICIONES — y sí forman la esfera")
	fmt.Println("   sobre el ecuador las perlas se apretujan contra UN punto: el broche. Ésa es")
	fmt.Println("   la cintura del reloj de arena. Medida sobre la esfera:")
	fmt.Println("\n        γ          distancia de arco al broche       ≈ 1/γ")
	broche := aLaEsfera(complex(1, 0))
	for _, g := range []float64{14.134725, 100, 1000, 1e6, 1e12} {
		p := aLaEsfera(alDisco(complex(0.5, g)))
		cosang := p[0]*broche[0] + p[1]*broche[1] + p[2]*broche[2]
		if cosang > 1 {
			cosang = 1
		}
		fmt.Printf("   %10.0e     %22.12e     %14.6e\n", g, math.Acos(cosang), 1/g)
	}
	fmt.Println("   → cuanto más alta la perla, más pegada al broche, y nunca lo toca.")
	fmt.Println("\n   Y ACÁ ESTÁ LA PARTE DEL CAPITÁN: una esfera se puede GIRAR. Girala y la")
	fmt.Println("   cintura del reloj se va a cualquier lado. Las posiciones del reloj barren")
	fmt.Println("   la esfera entera, y el dibujo no se deforma: los giros son rígidos.")
	fmt.Println("\n        giro       ¿siguen sobre la esfera?     ¿se conservan las distancias?")
	peorRig, peorNor := 0.0, 0.0
	base := make([][3]float64, 0, len(ps))
	for _, g := range ps {
		base = append(base, aLaEsfera(alDisco(complex(0.5, g))))
	}
	for _, a := range []float64{0.3, 1.0, math.Pi / 2, 2.5} {
		for i := range base {
			q := girar(base[i], a)
			if d := math.Abs(norma(q) - 1); d > peorNor {
				peorNor = d
			}
			if i > 0 {
				p0, p1 := girar(base[i-1], a), q
				d0 := base[i][0]*base[i-1][0] + base[i][1]*base[i-1][1] + base[i][2]*base[i-1][2]
				d1 := p0[0]*p1[0] + p0[1]*p1[1] + p0[2]*p1[2]
				if d := math.Abs(d0 - d1); d > peorRig {
					peorRig = d
				}
			}
		}
		fmt.Printf("   %7.3f          sí (%.1e)                  sí (%.1e)\n", a, peorNor, peorRig)
	}
	fmt.Printf("   → los %d puntos giran juntos sin deformarse: normas %.1e, distancias %.1e\n",
		len(ps), peorNor, peorRig)
	fmt.Println("     EL RELOJ DE ARENA EN TODAS SUS POSICIONES ES, LITERALMENTE, LA ESFERA.")

	// ---- LEY 3: infinite spheres touching at infinite points ----
	fmt.Println("\nLEY 3 · «INFINITAS ESFERAS TOCÁNDOSE EN INFINITOS PUNTOS» — ya existe y tiene nombre")
	fmt.Println("   se llaman CÍRCULOS DE FORD. Por cada fracción p/q hay un círculo apoyado en la")
	fmt.Println("   recta justo en p/q, de diámetro 1/q². Y pasa esto:")
	fmt.Println("\n        NUNCA se cruzan. O están separados, o SE TOCAN en un punto.")
	fmt.Println("        Y se tocan exactamente cuando |p·s − q·r| = 1.")
	var fr []frac
	for q := 1; q <= 14; q++ {
		for p := 0; p <= q; p++ {
			if mcd(p, q) == 1 {
				fr = append(fr, frac{p, q})
			}
		}
	}
	cruces, tangentes, sueltos, pares := 0, 0, 0, 0
	peorTg := 0.0
	for i := 0; i < len(fr); i++ {
		for j := i + 1; j < len(fr); j++ {
			a, b := fr[i], fr[j]
			ra, rb := 1/(2*float64(a.q*a.q)), 1/(2*float64(b.q*b.q))
			dx := float64(a.p)/float64(a.q) - float64(b.p)/float64(b.q)
			dy := ra - rb
			d := math.Hypot(dx, dy)
			suma := ra + rb
			det := a.p*b.q - a.q*b.p
			if det < 0 {
				det = -det
			}
			pares++
			switch {
			case d < suma-1e-12:
				cruces++
			case math.Abs(d-suma) <= 1e-12:
				tangentes++
				if det != 1 {
					fmt.Printf("   ⚠ tangente con |ps−qr| = %d (esperaba 1): %d/%d y %d/%d\n", det, a.p, a.q, b.p, b.q)
				}
				if e := math.Abs(d - suma); e > peorTg {
					peorTg = e
				}
			default:
				sueltos++
				if det == 1 {
					fmt.Printf("   ⚠ |ps−qr| = 1 pero NO se tocan: %d/%d y %d/%d\n", a.p, a.q, b.p, b.q)
				}
			}
		}
	}
	fmt.Printf("\n   con las fracciones de denominador hasta 14 → %d círculos, %d pares probados:\n", len(fr), pares)
	fmt.Printf("      se CRUZAN ......... %d\n", cruces)
	fmt.Printf("      se TOCAN .......... %d   (todos con |p·s − q·r| = 1, desvío %.1e)\n", tangentes, peorTg)
	fmt.Printf("      quedan sueltos .... %d\n", sueltos)
	fmt.Println("   → CERO cruces. Infinitos círculos, infinitos puntos de contacto, y ninguno")
	fmt.Println("     invade a otro. Es exactamente el dibujo que describió el capitán.")
	fmt.Println("\n   Y LOS REFUTADORES ME TRAJERON LA CUENTA EXACTA, que es mejor que la medición:")
	fmt.Println("\n        d² − (r₁+r₂)²  =  [ (p·s − q·r)² − 1 ] / (q²·s²)")
	fmt.Println("\n   el divisor es positivo siempre, así que el signo lo decide SOLO el entero")
	fmt.Println("   D = |p·s − q·r|, y como es entero no hay valores intermedios:")
	fmt.Println("\n        D = 0   → negativo → es la MISMA fracción (círculos anidados)")
	fmt.Println("        D = 1   → CERO exacto → se tocan en un punto")
	fmt.Println("        D ≥ 2   → positivo → separados, con hueco")
	fmt.Println("\n   verificada la identidad contra la medición sobre los mismos pares:")
	peorId := 0.0
	for i := 0; i < len(fr); i++ {
		for j := i + 1; j < len(fr); j++ {
			a, bb := fr[i], fr[j]
			ra, rb := 1/(2*float64(a.q*a.q)), 1/(2*float64(bb.q*bb.q))
			dx := float64(a.p)/float64(a.q) - float64(bb.p)/float64(bb.q)
			d2 := dx*dx + (ra-rb)*(ra-rb)
			izq := d2 - (ra+rb)*(ra+rb)
			det := a.p*bb.q - a.q*bb.p
			der := (float64(det*det) - 1) / float64(a.q*a.q*bb.q*bb.q)
			if e := math.Abs(izq - der); e > peorId {
				peorId = e
			}
		}
	}
	fmt.Printf("      %d pares · peor desvío entre la fórmula y la medición: %.1e\n", pares, peorId)
	fmt.Println("\n   ⚠ Y DOS CONDICIONES QUE ME FALTABAN, y no son cosméticas:")
	fmt.Println("   · las fracciones tienen que estar en TÉRMINOS MÍNIMOS **y ser DISTINTAS**. Si no,")
	fmt.Println("     aparece D = 0: 1/2 contra 2/4 da círculos ANIDADOS, tangentes por adentro. Se")
	fmt.Println("     tocan en un punto, sí, pero uno está metido adentro del otro: rompe el enunciado.")
	fmt.Println("   · la tangencia entre fracciones distintas y reducidas es SIEMPRE EXTERNA, porque")
	fmt.Println("     d² − (r₁−r₂)² = [(p·s−q·r)² + 1]/(q²s²) es positivo SIEMPRE.")
	fmt.Println("   · y de yapa: la condición se autoprotege — si |p·s−q·r| = 1 entonces las fracciones")
	fmt.Println("     ya están reducidas solas, porque el máximo común divisor divide a p·s−q·r.")

	fmt.Println("\nLEY 4 · ¿Y QUÉ EXPLICA? — acá los refutadores me borraron una frase entera")
	fmt.Println("\n   📌 LO QUE HABÍA ESCRITO Y ESTABA MAL: «los círculos de Ford son LOS horociclos del")
	fmt.Println("   grupo modular, que es el grupo que FABRICA las formas modulares, ASÍ QUE esta")
	fmt.Println("   geometría es LA CASA donde vive el segundo libro». Tres errores en un renglón.")
	fmt.Println("\n   ✅ LO CORRECTO, dicho con el artículo en su lugar:")
	fmt.Println("   · los círculos de Ford son UNA ÓRBITA de horociclos, no «los» horociclos: son la")
	fmt.Println("     órbita del horociclo de altura 1 en la cúspide infinito, bajo el grupo modular.")
	fmt.Println("     Horociclos hay uno por cada punto del borde y cada altura; Im z = 2 es uno y no")
	fmt.Println("     es de Ford. Además el horociclo del infinito es una RECTA, no un círculo.")
	fmt.Println("   · el grupo modular no FABRICA formas modulares: les IMPONE la condición. El peso,")
	fmt.Println("     el factor de automorfia, no vive en la familia de horociclos.")
	fmt.Println("\n   ❌ Y LA FRASE QUE HAY QUE BORRAR ES «LA CASA DONDE VIVE EL SEGUNDO LIBRO», porque")
	fmt.Println("   convierte una identificación geométrica correcta en una insinuación sobre ceros")
	fmt.Println("   que no se sostiene. Entre los círculos de Ford y la función L de la Δ hay SEIS")
	fmt.Println("   pasos escondidos, y los decisivos son invisibles desde esta geometría:")
	fmt.Println("\n      1. Ford es la órbita del grupo sobre el semiplano — sin peso ni multiplicador")
	fmt.Println("      2. una forma de peso 12 lleva un factor de automorfia que no está en los horociclos")
	fmt.Println("      3. que Δ sea cuspidal de peso 12 no se lee en ningún círculo de Ford")
	fmt.Println("      4. los τ(n) son la serie de Fourier a lo largo de las rectas Im z = const —")
	fmt.Println("         ÉSE es el único contacto genuino, y es con el horociclo del infinito (la recta)")
	fmt.Println("      5. L(s,Δ) = Σ τ(n)/nˢ necesita teoría de Hecke: dato aritmético, no geométrico")
	fmt.Println("      6. la ecuación funcional sale de z ↦ −1/z aplicado a una integral sobre el eje")
	fmt.Println("         imaginario — que es una GEODÉSICA, no un horociclo")
	fmt.Println("\n   ⟹ LA GEOMETRÍA DE FORD NO DICE NADA SOBRE LOS CEROS DE L(s,Δ). Es contabilidad de")
	fmt.Println("   Farey y de fracciones continuas, y el contorno clásico del método del círculo de")
	fmt.Println("   Rademacher. Útil, concreta, y sin relación con dónde caen los ceros.")
	fmt.Println("\n   ✅ LO QUE SÍ EXPLICA, sin adornos:")
	fmt.Println("   · que la línea crítica no tenga puntas: bajo el cambiaformas es un círculo")
	fmt.Println("     (pinchado en el broche, que es la imagen del infinito)")
	fmt.Println("   · que el espejo del libro sea la inversión del disco, y que el ecuador sea su")
	fmt.Println("     conjunto invariante — o sea EL círculo, no UN círculo")
	fmt.Println("   · que existe un empaquetamiento infinito de círculos tangentes sin cruzarse, con")
	fmt.Println("     la regla exacta |p·s − q·r| = 1, tal cual lo describió el capitán")
	fmt.Println("\n   ⚖️ Y NADA MÁS. Es el plano de una casa, no la dirección de nadie.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL FLASH DEL CAPITÁN DESCRIBE UNA GEOMETRÍA QUE EXISTE Y ES EXACTA:")
	fmt.Printf("  · bajo el cambiaformas la piel se vuelve el ECUADOR ..... %.1e sobre %d perlas\n", peorEc, len(ps))
	fmt.Printf("  · el espejo del libro ES la inversión del disco .......... %.1e\n", peorInv)
	fmt.Printf("  · el reloj en todas las posiciones = la esfera ........... giros rígidos a %.1e\n", peorRig)
	fmt.Printf("  · infinitas esferas tocándose sin cruzarse ............... %d contactos, %d cruces\n", tangentes, cruces)
	fmt.Printf("  · con la identidad exacta d²−(r₁+r₂)² = [(ps−qr)²−1]/(q²s²)  %.1e\n", peorId)
	fmt.Println("\n📌 Y TRES COSAS QUE LOS REFUTADORES ME CORRIGIERON, que valen más que lo de arriba:")
	fmt.Println("  1. la línea no ES el ecuador: SE VUELVE el ecuador bajo el cambiaformas, y es un")
	fmt.Println("     círculo PINCHADO — el broche es la imagen del infinito, no de ningún ρ finito.")
	fmt.Println("  2. faltaban condiciones en Ford: fracciones DISTINTAS y en términos mínimos, o")
	fmt.Println("     aparecen círculos anidados. Y la tangencia es siempre externa.")
	fmt.Println("  3. me borraron entera la frase «esta geometría es la casa donde vive el segundo")
	fmt.Println("     libro»: entre Ford y la función L de la Δ hay seis pasos, y los decisivos no se")
	fmt.Println("     ven desde acá. La geometría de Ford NO dice nada sobre ceros.")
	fmt.Println("\n⚖️ Y LA ADVERTENCIA MÁS IMPORTANTE DE TODAS: |w| = 1 ⟺ Re ρ = ½ es una")
	fmt.Println("reformulación TRIVIAL. «Todas las perlas tienen |w| = 1» y la Hipótesis de Riemann")
	fmt.Println("son la misma frase escrita dos veces. Que las perlas caigan en el ecuador no se")
	fmt.Println("midió: se supuso, porque las pescamos sobre la línea y la línea ES el ecuador.")
	fmt.Println("\nEs el PLANO DE LA CASA, no la dirección de nadie. ¿El premio? Todavía no.")

	escribirLamina(ps, peorEc, peorRig, peorInv, peorId, tangentes, cruces, sueltos, len(fr), fr)
}

func escribirLamina(ps []float64, peorEc, peorRig, peorInv, peorId float64,
	tangentes, cruces, sueltos, nfr int, fr []frac) {
	var b strings.Builder
	W, H := 1500.0, 1020.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌐 LA ESFERA — el reloj de arena en todas las posiciones, y las que se tocan</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">la línea crítica no es una línea: es el ECUADOR de una esfera</text>
`, W, H, W, H, W/2, W/2)

	// the sphere with the pearls on the equator
	cx, cy, R := 350.0, 380.0, 215.0
	fmt.Fprintf(&b, `<rect x="40" y="100" width="620" height="560" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="350" y="132" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA ESFERA DE RIEMANN</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#0d1c31" stroke="#3d6fa8" stroke-width="2"/>
<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="%.0f" fill="none" stroke="#7ee0c0" stroke-width="2.2"/>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">EL ECUADOR = la línea crítica</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">casquete norte · el afuera</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">casquete sur · el disco de la dimensión 0</text>
`, cx, cy, R, cx, cy, R, R*0.30, cx, cy-R-14, cx, cy-R*0.55, cx, cy+R*0.72)
	for i, g := range ps {
		if i%4 != 0 {
			continue
		}
		φ := 2 * math.Atan(1/(2*g))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2.6" fill="#7ee0c0" opacity="0.9"/>`,
			cx+R*math.Cos(φ), cy-R*0.30*math.Sin(φ))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2.6" fill="#7ee0c0" opacity="0.9"/>`,
			cx+R*math.Cos(-φ), cy-R*0.30*math.Sin(-φ))
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="6" fill="#ffb27a"/>
<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#ffb27a">el broche · la cintura del reloj</text>
<text x="350" y="624" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">%d perlas en el ecuador, altura %.1e · girá la esfera y la cintura va a cualquier lado</text>
<text x="350" y="646" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">los giros son RÍGIDOS: el dibujo no se deforma (%.1e)</text>
`, cx+R, cy, cx+R-160, cy-16, len(ps), peorEc, peorRig)

	// the Ford circles
	fx0, fy0, fw := 700.0, 560.0, 740.0
	fmt.Fprintf(&b, `<rect x="680" y="100" width="780" height="560" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1070" y="132" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">INFINITAS ESFERAS TOCÁNDOSE — los círculos de Ford</text>
<text x="1070" y="156" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">uno por cada fracción p/q, de diámetro 1/q², apoyado en p/q</text>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#3d6fa8" stroke-width="2"/>
`, fx0, fy0, fx0+fw, fy0)
	for _, f := range fr {
		r := fw / 2 / float64(f.q*f.q)
		x := fx0 + fw*float64(f.p)/float64(f.q)
		op := 0.95 - 0.05*float64(f.q)
		if op < 0.28 {
			op = 0.28
		}
		col := "#7ee0c0"
		if f.q > 3 {
			col = "#4fa3d1"
		}
		if f.q > 7 {
			col = "#5a4fa8"
		}
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="%.2f" fill="none" stroke="%s" stroke-width="1.6" opacity="%.2f"/>`,
			x, fy0-r, r, col, op)
		if f.q <= 3 {
			fmt.Fprintf(&b, `<text x="%.2f" y="%.0f" font-size="11" text-anchor="middle" font-family="monospace" fill="#7fa8cf">%d/%d</text>`,
				x, fy0+18, f.p, f.q)
		}
	}
	fmt.Fprintf(&b, `
<text x="1070" y="600" font-size="14" text-anchor="middle" font-family="monospace" fill="#9fd8a8">%d círculos · %d contactos · %d cruces · %d sueltos</text>
<text x="1070" y="624" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">se tocan EXACTAMENTE cuando |p·s − q·r| = 1 (vecinos de Farey)</text>
<text x="1070" y="646" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">infinitos círculos, infinitos contactos, ninguno invade a otro</text>
`, nfr, tangentes, cruces, sueltos)

	fmt.Fprintf(&b, `<rect x="40" y="680" width="700" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="390" y="710" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">✅ LO QUE EXPLICA</text>
<text x="66" y="742" font-size="14" font-family="Georgia" fill="#cfe6ff">· por qué la línea crítica no tiene puntas: ES UN CÍRCULO</text>
<text x="66" y="768" font-size="14" font-family="Georgia" fill="#cfe6ff">· por qué el cambiaformas no deforma: es un movimiento rígido</text>
<text x="66" y="794" font-size="14" font-family="Georgia" fill="#cfe6ff">· que existe un empaquetamiento infinito de círculos tangentes</text>
<text x="80" y="814" font-size="14" font-family="Georgia" fill="#cfe6ff">sin cruzarse, con la regla exacta |p·s − q·r| = 1</text>

<rect x="760" y="680" width="700" height="150" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1110" y="710" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ LO QUE NO EXPLICA</text>
<text x="786" y="742" font-size="14" font-family="Georgia" fill="#f3d9cf">· NO dice dónde está ninguna perla. Ni una.</text>
<text x="786" y="768" font-size="14" font-family="Georgia" fill="#f3d9cf">· que las perlas caigan en el ecuador NO se midió: se supuso,</text>
<text x="800" y="788" font-size="14" font-family="Georgia" fill="#f3d9cf">porque las pescamos sobre la línea y la línea ES el ecuador</text>
<text x="786" y="814" font-size="14" font-family="Georgia" fill="#f3d9cf">· los refutadores me borraron «la casa donde vive el segundo libro»</text>

<rect x="40" y="850" width="1420" height="130" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="750" y="886" font-size="21" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">ES EL PLANO DE LA CASA, NO LA DIRECCIÓN DE NADIE</text>
<text x="750" y="920" font-size="15" text-anchor="middle" font-family="Georgia" fill="#dce8f7">el capitán describió, sin saberlo, la esfera de Riemann (1857) y los círculos de Ford (1938). Pero ojo:</text>
<text x="750" y="942" font-size="15" text-anchor="middle" font-family="Georgia" fill="#dce8f7">|w|=1 ⟺ Re ρ=½ es una reformulación TRIVIAL: «todas en el ecuador» y la Hipótesis son la misma frase dos veces.</text>
<text x="750" y="968" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">Ninguna de las dos toca la hipótesis. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("la-esfera.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: la-esfera.svg")
}
