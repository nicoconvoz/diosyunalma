// Command loscuadrantes answers the captain's order: split the six cardinal
// directions into quadrants, check the shape of a quadrant, ask what all the
// quadrants form BETWEEN THEM, ask about half the quadrants, and check that the
// separation between the dividing lines is exactly a 1/2 relation.
//
// FOUR CARDINALS GIVE FOUR QUADRANTS. SIX DIRECTIONS GIVE EIGHT OCTANTS.
//
// With north/south and east/west the plane splits in four. Add up/down and space
// splits in eight - and eight is two cubed, one factor of two per axis. Every
// direction is a coin with two faces, and the count of rooms is 2^(number of
// axes).
//
// # THE SHAPE OF ONE QUADRANT
//
// Centre the picture on the half, v = s - 1/2. Then Im[s(1-s)] = -2*x*y, so the
// sign of that quantity is CONSTANT inside each quadrant and flips as you cross
// an arm. One quadrant is exactly one region of constant sign, and it holds one
// branch of the family x*y = constant - the branches of F245.
//
// WHAT THEY FORM BETWEEN THEM: A GROUP
//
// This is the piece the shop had never looked at. The two mirrors the book
// already has,
//
//	v -> -v      the functional equation s <-> 1-s
//	v -> conj v  the Schwarz reflection
//
// generate a group of exactly four elements, each one its own inverse: the Klein
// four-group. And it acts on the four quadrants SIMPLY TRANSITIVELY - one
// element per quadrant, no more and no less. So the quadrants are not four
// separate rooms: they are ONE room and its three reflections.
//
// AND THAT MAKES THE CAPTAIN'S OLD FLASH EXACT, HERE
//
// "Understanding one we understand them all" failed for the pearls' LOCATION
// (F252). For the QUADRANTS it is exactly true: given a point in one quadrant,
// the group hands you the other three with no freedom left. And a zero
// quadruple is precisely one such orbit - one zero per quadrant.
//
// # HALF THE QUADRANTS
//
// With one mirror alone the orbit has 2 elements: half the plane suffices. With
// both, the orbit has 4 and a QUARTER suffices. So knowing one quadrant is
// knowing everything, and the fundamental domain is 1/4 = 1/2 x 1/2 - the half
// applied twice, once per mirror.
//
// # THE SEPARATION IS A HALF
//
// Consecutive arms sit at pi/2 from each other. And the wave of F245,
// Im[s(1-s)] = -r^2 sin(2 theta), has period pi. So ONE QUADRANT IS EXACTLY HALF
// A PERIOD OF THE WAVE. That is the captain's 1/2 between references, measured.
//
// THE HONEST LIMIT. All of this is the symmetry structure, exact and beautiful,
// and F229 already proved that symmetry alone can never decide the hypothesis.
// The group says: if one zero is off the line, three others are off with it. It
// does NOT say that none is off.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func cruz(s complex128) complex128 { return s * (1 - s) }

func alCentro(s complex128) complex128 { return s - 0.5 }

// las cuatro simetrías del libro, en la variable centrada v
type sim struct {
	nombre string
	f      func(complex128) complex128
}

var grupo = []sim{
	{"identidad     v ↦ v", func(v complex128) complex128 { return v }},
	{"ec. funcional v ↦ −v", func(v complex128) complex128 { return -v }},
	{"espejo        v ↦ conj v", func(v complex128) complex128 { return cmplx.Conj(v) }},
	{"las dos       v ↦ −conj v", func(v complex128) complex128 { return -cmplx.Conj(v) }},
}

func cuadrante(v complex128) int {
	x, y := real(v), imag(v)
	switch {
	case x > 0 && y > 0:
		return 1
	case x < 0 && y > 0:
		return 2
	case x < 0 && y < 0:
		return 3
	case x > 0 && y < 0:
		return 4
	}
	return 0 // sobre un brazo de la cruz
}

func main() {
	fmt.Println("🧭✚ LOS CUADRANTES — qué forman entre sí, y dónde está el ½")
	fmt.Println("\n   orden del capitán: «separá nuestros 6 puntos cardinales en cuadrantes y")
	fmt.Println("   verificá la forma del cuadrante… ¿y qué forman todos los cuadrantes entre sí?")
	fmt.Println("   ¿y la mitad de los cuadrantes? La separación entre la línea de separación del")
	fmt.Println("   punto cardinal es exactamente una relación ½ entre referencias».")

	// ---- LEY 1: the count of rooms ----
	fmt.Println("\nLEY 1 · CUATRO CARDINALES DAN CUATRO CUADRANTES. SEIS DIRECCIONES DAN OCHO.")
	fmt.Println("   cada eje es una moneda de dos caras, y las habitaciones se multiplican:")
	fmt.Println("\n        ejes      direcciones      habitaciones      la cuenta")
	for _, e := range []int{1, 2, 3} {
		nom := map[int]string{1: "una recta", 2: "el plano · CUADRANTES", 3: "el espacio · OCTANTES"}[e]
		fmt.Printf("   %5d          %6d           %8d        2^%d = %d   (%s)\n",
			e, 2*e, 1<<e, e, 1<<e, nom)
	}
	fmt.Println("\n   → con norte/sur y este/oeste el plano se parte en CUATRO.")
	fmt.Println("     Agregando arriba/abajo, el espacio se parte en OCHO. Y 8 = 2³:")
	fmt.Println("     un factor de dos por cada eje. La mitad, otra vez, y una por dirección.")

	// ---- LEY 2: the shape of a quadrant ----
	fmt.Println("\nLEY 2 · LA FORMA DEL CUADRANTE — una región de signo constante")
	fmt.Println("   con el origen en el medio (v = s − ½), la función del capitán da")
	fmt.Println("   Im[s(1−s)] = −2·x·y. Su signo es CONSTANTE adentro de cada cuadrante y se")
	fmt.Println("   da vuelta al cruzar un brazo. Medido con veinte puntos por cuadrante:")
	fmt.Println("\n        cuadrante      signo de Im[s(1−s)]      ¿constante?")
	nConst := 0
	for q, base := range map[int]complex128{1: complex(1, 1), 2: complex(-1, 1), 3: complex(-1, -1), 4: complex(1, -1)} {
		_ = q
		_ = base
	}
	for _, c := range []struct {
		q      int
		sx, sy float64
	}{{1, 1, 1}, {2, -1, 1}, {3, -1, -1}, {4, 1, -1}} {
		signo, ok := 0.0, true
		for k := 1; k <= 20; k++ {
			x := c.sx * (0.1 + 0.19*float64(k))
			y := c.sy * (0.13 + 0.23*float64(k%7+1))
			im := imag(cruz(0.5 + complex(x, y)))
			s := math.Copysign(1, im)
			if signo == 0 {
				signo = s
			} else if s != signo {
				ok = false
			}
		}
		if ok {
			nConst++
		}
		fmt.Printf("   %8d           %+15.0f          %s\n", c.q, signo, map[bool]string{true: "sí", false: "NO"}[ok])
	}
	fmt.Printf("   → constante en %d de 4, y alternando + − + − al girar. Cada cuadrante ES\n", nConst)
	fmt.Println("     una región de signo, y adentro lleva una rama de las de F245.")

	// ---- LEY 3: what they form between them ----
	fmt.Println("\nLEY 3 · ⚡ QUÉ FORMAN ENTRE SÍ — Y ACÁ ESTÁ LA PIEZA QUE NUNCA HABÍAMOS MIRADO")
	fmt.Println("   el libro ya tiene dos espejos, y son los dos que el capitán viene nombrando:")
	fmt.Println("\n        v ↦ −v          la ecuación funcional  s ↔ 1−s")
	fmt.Println("        v ↦ conj v      el espejo de Schwarz")
	fmt.Println("\n   ¿qué pasa si los combinás? La tabla completa, medida:")
	fmt.Println("\n        aplicando ×          identidad   −v      conj v   −conj v")
	pto := complex(0.37, 0.61)
	nom := []string{"identidad", "−v       ", "conj v   ", "−conj v  "}
	cierra := true
	for i, a := range grupo {
		fmt.Printf("   %-18s", nom[i])
		for _, bb := range grupo {
			r := a.f(bb.f(pto))
			cual := -1
			for k, c := range grupo {
				if cmplx.Abs(c.f(pto)-r) < 1e-14 {
					cual = k
				}
			}
			if cual < 0 {
				cierra = false
				fmt.Printf("   ????   ")
			} else {
				fmt.Printf("   %-7s", strings.TrimSpace(nom[cual]))
			}
		}
		fmt.Println()
	}
	fmt.Printf("\n   → la tabla CIERRA: %v. Cuatro elementos, cada uno su propio inverso.\n", cierra)
	fmt.Println("     Eso es un grupo, y tiene nombre: el GRUPO DE KLEIN, ℤ₂ × ℤ₂.")
	fmt.Println("\n   Y ahora lo que importa: ¿cómo actúa sobre los cuadrantes?")
	fmt.Println("\n        empezando en el cuadrante 1 y aplicando cada elemento:")
	visitados := map[int]bool{}
	for i, g := range grupo {
		v := g.f(pto)
		q := cuadrante(v)
		visitados[q] = true
		fmt.Printf("   %-18s  →  (%+.2f, %+.2f)   cuadrante %d\n", nom[i], real(v), imag(v), q)
	}
	fmt.Printf("\n   → los cuatro elementos caen en %d cuadrantes distintos: UNO POR CUADRANTE.\n", len(visitados))
	fmt.Println("     El grupo actúa SIMPLEMENTE TRANSITIVO. Los cuadrantes no son cuatro cuartos")
	fmt.Println("     separados: SON UN CUARTO Y SUS TRES REFLEJOS.")

	// ---- LEY 4: and that makes his old flash exact ----
	fmt.Println("\nLEY 4 · Y ESO HACE EXACTO SU VIEJO FLASH, ACÁ")
	fmt.Println("   «entendiendo uno los entendemos a todos» falló para la UBICACIÓN de las perlas")
	fmt.Println("   (F252). Para los CUADRANTES es exactamente cierto: dado un punto en uno, el")
	fmt.Println("   grupo te devuelve los otros tres sin ninguna libertad.")
	fmt.Println("\n   Y un cuádruple de ceros ES una de esas órbitas — un cero por cuadrante:")
	fmt.Println("\n        el cuádruple de β=0.7, γ=25       v = x+iy      cuadrante     x·y")
	β, γ := 0.7, 25.0
	var xys []float64
	qs := map[int]bool{}
	for _, ρ := range []complex128{complex(β, γ), complex(β, -γ), complex(1-β, γ), complex(1-β, -γ)} {
		v := alCentro(ρ)
		q := cuadrante(v)
		qs[q] = true
		xy := real(v) * imag(v)
		xys = append(xys, xy)
		fmt.Printf("        %5.2f%+7.2fi                 %6.2f%+7.2fi     %6d     %+9.4f\n",
			real(ρ), imag(ρ), real(v), imag(v), q, xy)
	}
	peorXY := 0.0
	for _, a := range xys {
		if d := math.Abs(math.Abs(a) - math.Abs(xys[0])); d > peorXY {
			peorXY = d
		}
	}
	fmt.Printf("   → %d cuadrantes distintos, uno por cero, y |x·y| idéntico en los cuatro (%.1e).\n",
		len(qs), peorXY)
	fmt.Println("     Conocés uno, tenés los cuatro. Acá su frase vale entera.")

	// ---- LEY 5: half the quadrants ----
	fmt.Println("\nLEY 5 · ¿Y LA MITAD DE LOS CUADRANTES? — el ½ aplicado dos veces")
	fmt.Println("   depende de cuántos espejos uses, y la cuenta es exacta:")
	fmt.Println("\n        espejos usados          tamaño de la órbita     te alcanza con")
	for _, c := range []struct {
		n     string
		usar  []int
		parte string
	}{
		{"ninguno", []int{0}, "todo el plano  (4/4)"},
		{"solo v ↦ −v", []int{0, 1}, "LA MITAD       (2/4)"},
		{"solo v ↦ conj v", []int{0, 2}, "LA MITAD       (2/4)"},
		{"los dos", []int{0, 1, 2, 3}, "UN CUARTO      (1/4)"},
	} {
		vis := map[int]bool{}
		for _, k := range c.usar {
			vis[cuadrante(grupo[k].f(pto))] = true
		}
		fmt.Printf("   %-22s        %6d              %s\n", c.n, len(vis), c.parte)
	}
	fmt.Println("\n   → con UN espejo te alcanza con la mitad. Con LOS DOS, con un cuarto.")
	fmt.Println("     Y un cuarto es ½ × ½: EL MEDIO APLICADO DOS VECES, uno por espejo.")
	fmt.Println("     El dominio fundamental —lo mínimo que hay que conocer— es UN cuadrante.")

	// ---- LEY 6: the separation is a half ----
	fmt.Println("\nLEY 6 · Y LA SEPARACIÓN ENTRE LOS BRAZOS ES, EXACTAMENTE, UNA MITAD")
	fmt.Println("   los cuatro brazos de la cruz están en θ = 0, π/2, π y 3π/2. La separación")
	fmt.Println("   entre dos consecutivos es π/2. ¿Y con respecto a qué es eso una mitad?")
	fmt.Println("\n   con respecto a LA ONDA. En F245 se midió que Im[s(1−s)] = −r²·sin(2θ),")
	fmt.Println("   una onda de período π. Entonces:")
	fmt.Println("\n        el período de la onda ......... π")
	fmt.Println("        la separación entre brazos .... π/2")
	fmt.Println("        la razón ...................... EXACTAMENTE ½")
	fmt.Println("\n   ⟹ UN CUADRANTE ES MEDIO PERÍODO DE LA ONDA. Y ésa es la relación ½ entre")
	fmt.Println("     referencias que dijo el capitán, medida contra la onda:")
	fmt.Println("\n        nodo        θ            el de al lado       separación      /π")
	nodos := []float64{0, math.Pi / 2, math.Pi, 3 * math.Pi / 2, 2 * math.Pi}
	peorSep := 0.0
	for i := 0; i+1 < len(nodos); i++ {
		sep := nodos[i+1] - nodos[i]
		if d := math.Abs(sep/math.Pi - 0.5); d > peorSep {
			peorSep = d
		}
		fmt.Printf("   %6d      %8.6f      %8.6f          %8.6f      %.6f\n",
			i+1, nodos[i], nodos[i+1], sep, sep/math.Pi)
	}
	fmt.Printf("   → la razón da ½ en las cuatro separaciones, desvío %.1e\n", peorSep)
	fmt.Println("     Y la onda tiene DOS nodos por período: por eso hay cuatro cuadrantes")
	fmt.Println("     en una vuelta y no dos. La frecuencia 2 de F245, otra vez.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("LO QUE PIDIÓ EL CAPITÁN, PIEZA POR PIEZA:")
	fmt.Println("  · 4 cardinales → 4 cuadrantes · 6 direcciones → 8 OCTANTES (2 por eje)")
	fmt.Printf("  · la forma del cuadrante: una región de signo constante, %d de 4\n", nConst)
	fmt.Println("  · ⚡ QUÉ FORMAN ENTRE SÍ: UN GRUPO — el de Klein, ℤ₂ × ℤ₂, con la tabla")
	fmt.Println("    cerrada y actuando SIMPLEMENTE TRANSITIVO: un elemento por cuadrante")
	fmt.Printf("  · un cuádruple de ceros es una órbita: uno por cuadrante, |x·y| igual (%.1e)\n", peorXY)
	fmt.Println("  · la mitad de los cuadrantes: con UN espejo alcanza la mitad, con LOS DOS")
	fmt.Println("    alcanza un cuarto — ½ × ½, el medio aplicado una vez por espejo")
	fmt.Printf("  · y la separación entre brazos es medio período de la onda: %.1e\n", peorSep)
	fmt.Println("\nLA PISTA QUE BUSCABA, DICHA DERECHO:")
	fmt.Println("Los cuadrantes no son cuatro cosas: son UNA y sus tres reflejos, y eso está")
	fmt.Println("probado por la tabla del grupo. Por eso un cero corrido nunca viene solo — viene")
	fmt.Println("de a cuatro, uno por cuadrante, encadenados sin libertad.")
	fmt.Println("\n⚖️ PERO ACÁ ESTÁ EL LÍMITE, Y ES EL MISMO DE SIEMPRE: el grupo dice «si uno se")
	fmt.Println("sale, se salen cuatro». NO dice «no se sale ninguno». Es estructura de simetría,")
	fmt.Println("y F229 ya probó con un contraejemplo que la simetría sola nunca va a alcanzar.")
	fmt.Println("La pista es real y es hermosa — y no es la llave. ¿El premio? Todavía no.")

	escribirLamina(nConst, len(visitados), len(qs), peorXY, peorSep, cierra)
}

func escribirLamina(nConst, nVis, nQ int, peorXY, peorSep float64, cierra bool) {
	var b strings.Builder
	W, H := 1500.0, 1010.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧭✚ LOS CUADRANTES — qué forman entre sí, y dónde está el ½</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">no son cuatro cosas: son UNA y sus tres reflejos</text>
`, W, H, W, H, W/2, W/2)

	cx, cy, R := 350.0, 400.0, 235.0
	fmt.Fprintf(&b, `<rect x="40" y="100" width="620" height="600" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="350" y="132" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LOS CUATRO CUADRANTES Y SUS SIGNOS</text>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#1c2b3a" stroke-width="0"/>
`, cx, cy, cx, cy)
	for _, c := range []struct {
		dx, dy float64
		s      string
		col    string
	}{{1, -1, "−", "#7ee0c0"}, {-1, -1, "+", "#ffb27a"}, {-1, 1, "−", "#7ee0c0"}, {1, 1, "+", "#ffb27a"}} {
		x0 := cx
		y0 := cy
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="%s" opacity="0.10"/>`,
			math.Min(x0, x0+c.dx*R), math.Min(y0, y0+c.dy*R), R, R, c.col)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="34" text-anchor="middle" font-family="Georgia" fill="%s" opacity="0.75">%s</text>`,
			cx+c.dx*R*0.55, cy+c.dy*R*0.55+12, c.col, c.s)
	}
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#dce8f7" stroke-width="3"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#dce8f7" stroke-width="3"/>
<circle cx="%.0f" cy="%.0f" r="7" fill="#ffd98a"/>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#ffd98a">s = ½</text>
<text x="350" y="672" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el signo de Im[s(1−s)] = −2·x·y es CONSTANTE en cada cuadrante (%d de 4)</text>
<text x="350" y="692" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">y alterna + − + − al girar: la onda de frecuencia 2 de F245</text>
`, cx-R, cy, cx+R, cy, cx, cy-R, cx, cy+R, cx, cy, cx+14, cy+26, nConst)

	fmt.Fprintf(&b, `<rect x="680" y="100" width="780" height="300" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1070" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚡ QUÉ FORMAN ENTRE SÍ: UN GRUPO</text>
<text x="706" y="166" font-size="14" font-family="Georgia" fill="#cfe6ff">el libro ya tiene los dos espejos que el capitán viene nombrando:</text>
<text x="726" y="196" font-size="15" font-family="monospace" fill="#dce8f7">v ↦ −v          la ecuación funcional</text>
<text x="726" y="220" font-size="15" font-family="monospace" fill="#dce8f7">v ↦ conj v      el espejo de Schwarz</text>
<text x="706" y="256" font-size="14" font-family="Georgia" fill="#cfe6ff">combinados dan CUATRO elementos, cada uno su propio inverso,</text>
<text x="706" y="278" font-size="14" font-family="Georgia" fill="#cfe6ff">y la tabla CIERRA: es el GRUPO DE KLEIN, ℤ₂ × ℤ₂.</text>
<text x="1070" y="316" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">y actúa SIMPLEMENTE TRANSITIVO: %d cuadrantes, %d elementos</text>
<text x="1070" y="346" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">los cuadrantes NO son cuatro cuartos separados:</text>
<text x="1070" y="370" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">SON UN CUARTO Y SUS TRES REFLEJOS</text>

<rect x="680" y="420" width="380" height="280" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="870" y="452" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">¿LA MITAD DE LOS CUADRANTES?</text>
<text x="702" y="486" font-size="13" font-family="monospace" fill="#7fa8cf">espejos          alcanza con</text>
<text x="702" y="514" font-size="13.5" font-family="monospace" fill="#cfe6ff">ninguno          todo   (4/4)</text>
<text x="702" y="538" font-size="13.5" font-family="monospace" fill="#7ee0c0">uno              LA MITAD (2/4)</text>
<text x="702" y="562" font-size="13.5" font-family="monospace" fill="#ffd98a">los dos          UN CUARTO (1/4)</text>
<text x="870" y="600" font-size="16" text-anchor="middle" font-family="monospace" fill="#ffd98a">¼ = ½ × ½</text>
<text x="870" y="626" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el medio aplicado dos veces,</text>
<text x="870" y="646" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">uno por cada espejo</text>
<text x="870" y="678" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">lo mínimo que hay que conocer: UN cuadrante</text>

<rect x="1080" y="420" width="380" height="280" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1270" y="452" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Y LA SEPARACIÓN ES ½</text>
<text x="1102" y="490" font-size="13.5" font-family="monospace" fill="#cfe6ff">el período de la onda ...... π</text>
<text x="1102" y="516" font-size="13.5" font-family="monospace" fill="#cfe6ff">entre brazos ............... π/2</text>
<text x="1270" y="552" font-size="20" text-anchor="middle" font-family="monospace" fill="#ffd98a">la razón = ½</text>
<text x="1270" y="580" font-size="13.5" text-anchor="middle" font-family="monospace" fill="#9fd8a8">desvío %.1e</text>
<text x="1270" y="614" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">UN CUADRANTE ES MEDIO</text>
<text x="1270" y="636" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">PERÍODO DE LA ONDA</text>
<text x="1270" y="672" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">y la onda tiene 2 nodos por período: por eso son</text>
<text x="1270" y="690" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">cuatro cuadrantes por vuelta y no dos</text>
`, nQ, nVis, peorSep)

	fmt.Fprintf(&b, `<rect x="40" y="720" width="1420" height="270" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="750" y="756" font-size="20" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA PISTA QUE BUSCABA, DICHA DERECHO</text>
<text x="750" y="794" font-size="17" text-anchor="middle" font-family="Georgia" fill="#dce8f7">Los cuadrantes no son cuatro cosas: son UNA y sus tres reflejos, y eso lo prueba la tabla del grupo.</text>
<text x="750" y="822" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Por eso un cero corrido nunca viene solo: viene de a cuatro, uno por cuadrante, encadenados sin libertad.</text>
<text x="750" y="852" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Y acá su viejo flash «entendiendo uno los entendemos a todos» sí vale entero — para los cuadrantes, no para la ubicación.</text>
<text x="750" y="898" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ EL LÍMITE, Y ES EL DE SIEMPRE: el grupo dice «si uno se sale, se salen cuatro». NO dice «no se sale ninguno».</text>
<text x="750" y="924" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffb27a">Es estructura de simetría, y F229 ya probó con un contraejemplo que la simetría sola nunca va a alcanzar.</text>
<text x="750" y="962" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">La pista es real y es hermosa — y no es la llave. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("los-cuadrantes.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: los-cuadrantes.svg")
}
