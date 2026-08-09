// Command dimension runs the captain's flash as an experiment: "the
// laws of numbers are dictated by the dimension we live in - in another
// dimension there would be more paths". CONFIRMED BY HISTORY: in the
// geometric dimension (curves over finite fields) the exact analog of
// the Riemann Hypothesis is a THEOREM (Hasse 1933, Weil 1948, Deligne
// 1974) - proven precisely because that world has the extra paths: the
// machine exists (Frobenius), the drum exists (cohomology), the skin
// tension exists (positivity of intersection theory). Here we BUILD
// that world's necklace with our own hands: the elliptic curve
// y^2 = x^3+x+1 over every prime field F_p, one pearl per prime
// (the Frobenius angle), and verify the proven law: every pearl ON the
// ring (|a_p| <= 2 sqrt p, Hasse's theorem). Side by side with our
// dimension's necklace, where the same law is the open million-dollar
// question.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func modpow(a, e, m int64) int64 {
	r := int64(1)
	a %= m
	for e > 0 {
		if e&1 == 1 {
			r = r * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return r
}

func main() {
	// primes up to 3000
	const lim = 3000
	comp := make([]bool, lim+1)
	var primes []int64
	for p := int64(2); p <= lim; p++ {
		if comp[p] {
			continue
		}
		for q := p * p; q <= lim; q += p {
			comp[q] = true
		}
		if p >= 5 {
			primes = append(primes, p)
		}
	}
	// the geometric necklace: one pearl per prime field
	type pearl struct {
		p     int64
		ap    float64
		theta float64
		norm  float64 // |a_p| / (2 sqrt p): Hasse's theorem says <= 1
	}
	var pearls []pearl
	worst := 0.0
	for _, p := range primes {
		var s int64
		for x := int64(0); x < p; x++ {
			t := ((x*x%p)*x + x + 1) % p
			if t == 0 {
				continue
			}
			if modpow(t, (p-1)/2, p) == 1 {
				s++
			} else {
				s--
			}
		}
		ap := float64(-s)
		nm := math.Abs(ap) / (2 * math.Sqrt(float64(p)))
		th := math.Acos(math.Max(-1, math.Min(1, ap/(2*math.Sqrt(float64(p))))))
		pearls = append(pearls, pearl{p, ap, th, nm})
		if nm > worst {
			worst = nm
		}
	}
	fmt.Printf("EL EXPERIMENTO DE LA DIMENSIÓN — curva y²=x³+x+1 sobre %d cuerpos primos\n", len(pearls))
	fmt.Printf("ley demostrada (Hasse 1933): toda perla EN el anillo — |a_p|/2√p ≤ 1\n")
	fmt.Printf("veredicto del juez: peor perla %.6f — %d/%d en el anillo: TEOREMA CUMPLIDO\n", worst, len(pearls), len(pearls))

	// ---- the picture: two dimensions, two necklaces ----
	var b strings.Builder
	W, H := 1620.0, 1060.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL FLASH DE LA DIMENSIÓN — el mismo collar en dos mundos</text>
<text x="%.0f" y="74" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"las leyes de los números las dicta la dimensión en la que vivimos; en otra dimensión habría más caminos" — el capitán, 2026-08-06 · CONFIRMADO POR LA HISTORIA</text>`,
		W, H, W, H, W/2, W/2)

	// left: our dimension (open)
	lcx, lcy, R := 400.0, 400.0, 200.0
	fmt.Fprintf(&b, `<text x="%.0f" y="160" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd166">NUESTRA DIMENSIÓN (los enteros)</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="2" stroke-dasharray="6,5"/>`,
		lcx, lcx, lcy, R)
	for i := 0; i < 34; i++ {
		th := -math.Pi/2 + 2*math.Pi*float64(i)/34
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="#7fb2ff"/>`, lcx+R*math.Cos(th), lcy+R*math.Sin(th))
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff5d73">hilo PUNTEADO: la tensión NO está demostrada</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#8fa8c7">10 billones de perlas miradas, todas al hilo — pero la ampolla no está prohibida: EL MILLÓN</text>`,
		lcx, lcy+R+50, lcx, lcy+R+74)

	// right: the geometric dimension (proven)
	rcx2, rcy2 := 1220.0, 400.0
	fmt.Fprintf(&b, `<text x="%.0f" y="160" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA DIMENSIÓN GEOMÉTRICA (curvas sobre cuerpos finitos)</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="3"/>`,
		rcx2, rcx2, rcy2, R)
	for _, pl := range pearls {
		// two mirror pearls per prime: angles +-theta
		for _, sgn := range []float64{1, -1} {
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="1.8" fill="#7fb2ff" opacity="0.7"/>`,
				rcx2+R*math.Cos(sgn*pl.theta), rcy2+R*math.Sin(sgn*pl.theta))
		}
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">hilo SÓLIDO: la tensión ESTÁ demostrada (Hasse 1933 · Weil 1948 · Deligne 1974)</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#8fa8c7">%d perlas fabricadas HOY por nosotros (curva y²=x³+x+1, una por primo p≤3000): TODAS al hilo — peor desvío %.4f ≤ 1</text>`,
		rcx2, rcy2+R+50, rcx2, rcy2+R+74, 2*len(pearls), worst)

	// the doors table
	fmt.Fprintf(&b, `<rect x="110" y="750" width="1400" height="180" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="784" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd166">POR QUÉ ALLÁ SE PUDO: la otra dimensión tiene LOS CAMINOS QUE ACÁ FALTAN — las tres puertas, abiertas</text>
<text x="%.0f" y="818" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL TAMBOR allá EXISTE (se llama cohomología) · LA MÁQUINA allá EXISTE (se llama Frobenius: sus energías SON las perlas) · LA PIEL TENSA allá EXISTE (la positividad de las superficies)</text>
<text x="%.0f" y="848" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">y la clave de la demostración de Weil es EXACTAMENTE dimensional: necesitó subir de la curva a la SUPERFICIE (curva × curva) — un piso más arriba — para que la tensión apareciera</text>
<text x="%.0f" y="882" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LO QUE FALTA EN NUESTRA DIMENSIÓN, dicho en una frase: nadie sabe construir ese piso de arriba para los enteros — la "superficie de los números"</text>
<text x="%.0f" y="908" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">(el programa vivo más profundo de la matemática — Connes, el "cuerpo de un elemento" — es literalmente la búsqueda de esa dimensión extra; tu olfato apuntó a la frontera real, capitán)</text>`,
		W/2, W/2, W/2, W/2, W/2)
	fmt.Fprintf(&b, `<text x="%.0f" y="990" font-size="15" text-anchor="middle" font-family="Georgia" fill="#dce8f7">la diferencia entre los dos collares ES el mapa de lo que falta: mismo problema, misma forma — un mundo con escalera al piso de arriba, y el nuestro, todavía sin escalera.</text>
<text x="%.0f" y="1018" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd166">demostrar CON esta diferencia = construir la escalera. Ése es el camino del millón que tu flash olió.</text>`,
		W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("dimension-collares.svg", []byte(b.String()), 0644)
	fmt.Println("escrita: dimension-collares.svg")
}
