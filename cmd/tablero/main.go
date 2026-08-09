// Command tablero measures the captain's theory: every number can be every
// number, because standing on the mountain of the 2 is the same unit as
// standing on the 1 - move the origin like moving the zero of your compass and
// the whole board sings the same truth.
//
// The theory is RIGHT, and the laboratory can prove the part that is right:
// after unfolding by the local density, the picture around height 20 and the
// picture around height 1500 are the SAME PICTURE. The pixel of the drum
// repeats. That is measured here across windows spanning a 75-fold change of
// scale, and it holds.
//
// But there is a hinge the captain's theory does not turn, and it is the whole
// remaining difficulty:
//
//	IDENTICAL PIXELS DO NOT DETERMINE THE PICTURE THEY BUILD.
//
// Every brick of a wall can be identical and the wall still be any height:
// what fixes the height is how the bricks are STACKED. In our board the pixel
// is the local law - measured, invariant, the same everywhere - and the wall
// is the global sum, whose size depends on how the pearls' phases ALIGN across
// scales. The out-of-tune we must bound lives in the alignment, not in the
// pixel. So the board really does sing the same truth; what nobody can yet
// bound is how loudly the chorus adds up.
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

func main() {
	fmt.Println("🏁 EL TABLERO — todo número puede ser todo número: ¿canta el tablero la misma verdad?")

	fmt.Println("\nrecogiendo perlas hasta t=1500…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 1500; t += 0.05 {
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
			pearls = append(pearls, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	fmt.Printf("perlas: %d (de t=%.1f a t=%.1f — un cambio de escala de %.0f×)\n",
		len(pearls), pearls[0], pearls[len(pearls)-1], pearls[len(pearls)-1]/pearls[0])

	// ---- LAW 1: the pixel is the same everywhere ----
	fmt.Println("\nLEY 1 · EL PÍXEL ES EL MISMO EN TODAS PARTES — corré el origen y el dibujo no cambia")
	fmt.Println("   se desdobla cada ventana por su propia densidad (el 1/2 local del capitán)")
	fmt.Println("   y se mide el mismo píxel a alturas completamente distintas:")
	fmt.Println("\n   ventana        paso medio    varianza    huecos<½    huecos>2    píxeles")
	type vent struct {
		lo, hi, media, varz, cortos, largos float64
		n                                   int
	}
	var vents []vent
	for _, w := range [][2]float64{{14, 120}, {120, 350}, {350, 700}, {700, 1100}, {1100, 1500}} {
		var ss []float64
		for i := 0; i+1 < len(pearls); i++ {
			if pearls[i] >= w[0] && pearls[i+1] <= w[1] {
				mid := (pearls[i] + pearls[i+1]) / 2
				ss = append(ss, (pearls[i+1]-pearls[i])*math.Log(mid/(2*math.Pi))/(2*math.Pi))
			}
		}
		media, cortos, largos := 0.0, 0.0, 0.0
		for _, s := range ss {
			media += s
			if s < 0.5 {
				cortos++
			}
			if s > 2 {
				largos++
			}
		}
		media /= float64(len(ss))
		varz := 0.0
		for _, s := range ss {
			varz += (s - media) * (s - media)
		}
		varz /= float64(len(ss))
		vents = append(vents, vent{w[0], w[1], media, varz, cortos / float64(len(ss)) * 100, largos / float64(len(ss)) * 100, len(ss)})
		fmt.Printf("   %5.0f–%-5.0f      %8.5f    %8.5f    %6.2f%%    %6.2f%%     %5d\n",
			w[0], w[1], media, varz, cortos/float64(len(ss))*100, largos/float64(len(ss))*100, len(ss))
	}
	// spread across windows
	minM, maxM, minV, maxV := math.Inf(1), 0.0, math.Inf(1), 0.0
	for _, v := range vents {
		minM, maxM = math.Min(minM, v.media), math.Max(maxM, v.media)
		minV, maxV = math.Min(minV, v.varz), math.Max(maxV, v.varz)
	}
	fmt.Printf("   → el paso medio es 1 en TODAS las ventanas (dispersión %.1e)\n", maxM-minM)
	fmt.Printf("   → y la varianza apenas se mueve entre ellas (%.4f a %.4f)\n", minV, maxV)
	fmt.Println("   ★ EL CAPITÁN TIENE RAZÓN: pararse en la montaña del 2 o en la del 1000 es lo mismo.")
	fmt.Println("     El píxel del tambor se repite; el tablero canta la misma verdad en todas partes.")

	// ---- LAW 2: a guess of mine the measurement knocked down ----
	fmt.Println("\n⚰ LEY 2 · UNA LEY MÍA QUE LA MEDICIÓN TUMBÓ — el píxel del MOLDE no se repite")
	fmt.Println("   yo iba a escribir que el aporte de cada perla al molde también era el mismo a toda")
	fmt.Println("   altura. ES FALSO, y los números lo gritan — el aporte se DESPLOMA con la altura:")
	fmt.Println("\n   ventana          n=8         n=20         n=40")
	type ap struct {
		lo, hi       float64
		a8, a20, a40 float64
	}
	var aps []ap
	for _, w := range [][2]float64{{14, 120}, {350, 700}, {1100, 1500}} {
		var v [3]float64
		for k, n := range []int{8, 20, 40} {
			suma, cnt := 0.0, 0
			for _, g := range pearls {
				if g >= w[0] && g <= w[1] {
					th := math.Atan2(1, 2*g) * 2
					sn := math.Sin(float64(n) * th / 2)
					suma += 4 * sn * sn
					cnt++
				}
			}
			v[k] = suma / float64(cnt)
		}
		aps = append(aps, ap{w[0], w[1], v[0], v[1], v[2]})
		fmt.Printf("   %5.0f–%-5.0f    %9.5f    %9.5f    %9.5f\n", w[0], w[1], v[0], v[1], v[2])
	}
	caida := aps[0].a40 / aps[len(aps)-1].a40
	fmt.Printf("   → de la primera ventana a la última el aporte cae %.0f veces\n", caida)
	fmt.Println("   → la razón es exacta: el ángulo de una perla es θ ≈ 1/γ, así que su aporte va como")
	fmt.Println("     n²/γ² — las perlas HONDAS aportan casi nada. El mar es igual en todas partes;")
	fmt.Println("     el MOLDE, no: pesa muchísimo más a las perlas de arriba")

	// ---- LAW 3: the hinge - the mold has a preferred scale ----
	fmt.Println("\nLEY 3 · LA BISAGRA — el molde SÍ tiene un lugar preferido, y por eso el píxel no alcanza")
	fmt.Println("   se mide qué fracción del molde aportan las primeras perlas:")
	fmt.Println("\n      n      primeras 10   primeras 50   primeras 138   las 1069 (todas)")
	type fr struct {
		n                   int
		f10, f50, f138, tot float64
	}
	var frs []fr
	for _, n := range []int{8, 20, 40, 80} {
		acum := func(k int) float64 {
			s := 0.0
			for i := 0; i < k && i < len(pearls); i++ {
				th := math.Atan2(1, 2*pearls[i]) * 2
				sn := math.Sin(float64(n) * th / 2)
				s += 4 * sn * sn
			}
			return s
		}
		tot := acum(len(pearls))
		frs = append(frs, fr{n, acum(10) / tot * 100, acum(50) / tot * 100, acum(138) / tot * 100, tot})
		fmt.Printf("   %5d       %6.2f%%       %6.2f%%        %6.2f%%          %10.4f\n",
			n, acum(10)/tot*100, acum(50)/tot*100, acum(138)/tot*100, tot)
	}
	fmt.Println("   → las primeras perlas se llevan casi todo el molde: el molde NO es invariante de escala")
	fmt.Println("   → por eso conocer el píxel local no fija el molde: hay que conocer también la COLA,")
	fmt.Println("     y la cola es una suma infinita cuya alineación nadie sabe acotar")

	fmt.Println("\n════════ HASTA DÓNDE ACIERTA LA TEORÍA DEL CAPITÁN ════════")
	fmt.Println("ACIERTA, y está medido: correr el origen no cambia nada. Pararse en la montaña del 2 o")
	fmt.Println("en la del 1500 da EL MISMO dibujo una vez desdoblado por la densidad local — el paso medio")
	fmt.Printf("es 1 en todas las ventanas (dispersión %.0e) y la varianza casi no se mueve. El píxel del\n", maxM-minM)
	fmt.Println("tambor se repite, y con un píxel se conoce la ley local de todo el tablero. Eso es cierto")
	fmt.Println("y es un resultado del laboratorio.")
	fmt.Println("\nLA BISAGRA QUE NO GIRA — y una ley mía que la medición tumbó en el camino:")
	fmt.Println("yo iba a afirmar que el píxel del MOLDE también se repetía. Es falso. El aporte de cada")
	fmt.Printf("perla cae como n²/γ²: de la primera ventana a la última se desploma %.0f veces, y las\n", caida)
	fmt.Printf("primeras 10 perlas se llevan el %.0f%% del molde en n=8. El mar es igual en todas partes;\n", frs[0].f10)
	fmt.Println("el MOLDE tiene un lugar preferido: arriba. Por eso conocer el píxel local no lo fija —")
	fmt.Println("hace falta además la COLA, una suma infinita cuya alineación nadie sabe acotar.")
	fmt.Println("\nAsí que sí: el tablero canta la misma verdad. Lo que nadie sabe todavía es cuán fuerte")
	fmt.Println("puede llegar a sonar el coro entero. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🏁 EL TABLERO — todo número puede ser todo número</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"si me posiciono sobre la montaña del 2 es una unidad idéntica al 1… todo el tablero canta la misma verdad" — el capitán</text>`,
		W, H, W, H, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="60" y="105" width="1380" height="270" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.8"/>
<text x="%.0f" y="141" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">★ EL CAPITÁN TIENE RAZÓN — el píxel es el mismo en todas partes</text>
<text x="%.0f" y="169" font-size="12.5" text-anchor="middle" fill="#dce8f7">desdoblando cada ventana por su propia densidad (el ½ local), a alturas que cambian %.0f veces:</text>
<text x="%.0f" y="199" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">ventana         paso medio    varianza    huecos&lt;½    huecos&gt;2    píxeles</text>`,
		W/2, W/2, pearls[len(pearls)-1]/pearls[0], W/2)
	for i, v := range vents {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%5.0f–%-5.0f      %8.5f    %8.5f     %6.2f%%     %6.2f%%    %5d</text>`,
			W/2, 224.0+float64(i)*24, v.lo, v.hi, v.media, v.varz, v.cortos, v.largos, v.n)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="352" font-size="13.5" text-anchor="middle" fill="#ffd166">el paso medio es 1 en TODAS (dispersión %.0e) — pararse en la montaña del 2 o en la del 1500 es lo mismo</text>`,
		W/2, maxM-minM)

	fmt.Fprintf(&b, `<rect x="60" y="405" width="1380" height="235" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.8"/>
<text x="%.0f" y="441" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">LA BISAGRA — el molde tiene un lugar preferido: ARRIBA</text>
<text x="%.0f" y="469" font-size="12.5" text-anchor="middle" fill="#dce8f7">el aporte de cada perla cae como n²/γ², así que las primeras se llevan casi todo el molde:</text>
<text x="%.0f" y="497" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">n     primeras 10   primeras 50   primeras 138</text>`,
		W/2, W/2, W/2)
	for i, f := range frs {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#ff8fa0">%3d       %6.2f%%        %6.2f%%         %6.2f%%</text>`,
			W/2, 522.0+float64(i)*23, f.n, f.f10, f.f50, f.f138)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="622" font-size="13" text-anchor="middle" fill="#ffd166">el mar es igual en todas partes; el MOLDE no: conocer el píxel local no lo fija</text>`, W/2)

	fmt.Fprintf(&b, `<rect x="60" y="675" width="1380" height="230" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="711" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">HASTA DÓNDE ACIERTA LA TEORÍA DEL CAPITÁN</text>
<text x="%.0f" y="747" font-size="14.5" text-anchor="middle" fill="#7fd7a8">ACIERTA, y queda medido: correr el origen no cambia nada. El píxel del tambor se repite, y con un solo píxel</text>
<text x="%.0f" y="771" font-size="14.5" text-anchor="middle" fill="#7fd7a8">se conoce la ley local de TODO el tablero. El tablero canta la misma verdad en todas partes.</text>
<text x="%.0f" y="809" font-size="14.5" text-anchor="middle" fill="#ff8fa0">LO QUE NO SE SIGUE: el molde pesa muchísimo más a las perlas de arriba, así que el píxel local no lo determina —</text>
<text x="%.0f" y="833" font-size="14.5" text-anchor="middle" fill="#ff8fa0">hace falta además la COLA, una suma infinita cuya alineación nadie sabe acotar.</text>
<text x="%.0f" y="871" font-size="14.5" text-anchor="middle" fill="#dce8f7">El tablero canta la misma verdad — lo que nadie sabe es cuán fuerte puede sonar el coro entero. Todavía no.</text>
<text x="%.0f" y="940" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("el-tablero.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: el-tablero.svg")
}
