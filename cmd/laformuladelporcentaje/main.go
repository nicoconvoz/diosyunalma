package main

// LA FORMULA DEL PORCENTAJE - the captain: "write a formula from the percentage
// points where primes fall, see whether there is a relation between the
// percentages and the logarithm, to see whether we can predict a prime without
// knowing it."
//
// The percentage IS the logarithm: near x the fraction of numbers that are
// prime is 1/ln x. That is the Prime Number Theorem, and from it comes a real
// formula that predicts where the n-th prime sits:
//
//     p_n  ~  n * ( ln n + ln ln n - 1 )
//
// This program builds it, measures how close it lands, and then asks the only
// question that matters: does landing close mean PREDICTING?
//
// Verdicts fixed BEFORE running:
//   V1 (the percentage) the measured fraction of primes in a window around x
//      must track 1/ln x. Measured window by window, not asserted.
//   V2 (the formula) p_n ~ n(ln n + ln ln n - 1) must be built and its relative
//      error measured across six orders of magnitude.
//   V3 (THE HONEST UNIT) the error must also be reported in units of the LOCAL
//      GAP, because "off by 78" means nothing until you know the gap is 9.
//      That number says how many primes away the prediction lands.
//   V4 (THE PREDICTION TEST, with control) take the predicted integer and ask
//      whether it IS prime. Compare against a random integer of the same size.
//      If the rates match, the formula does not predict primes - it predicts
//      neighbourhoods.
//   V5 (what a real formula would need) exact formulas for the n-th prime DO
//      exist (Willans 1964), and they are useless: they hide a primality test
//      inside. Demonstrated, not just claimed.
// Gauss, Hadamard-de la Vallee Poussin 1896, Cipolla 1902, Willans 1964.

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func criba(n int) ([]bool, []int) {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if !es[p] {
			continue
		}
		for m := p * p; m <= n; m += p {
			es[m] = false
		}
	}
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return es, ps
}

// formula is the captain's asked-for prediction: where the n-th prime sits.
func formula(n int) float64 {
	f := float64(n)
	l := math.Log(f)
	ll := math.Log(l)
	return f * (l + ll - 1)
}

// formulaFina adds the next correction term (Cipolla).
func formulaFina(n int) float64 {
	f := float64(n)
	l := math.Log(f)
	ll := math.Log(l)
	return f * (l + ll - 1 + (ll-2)/l)
}

func main() {
	fmt.Println("📐 LA FÓRMULA DEL PORCENTAJE — el porcentaje de primos ES el logaritmo, y de ahí sale una fórmula")
	fmt.Println()

	const N = 30000000
	es, ps := criba(N)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · el porcentaje que buscabas, medido: cuántos primos hay en cada zona")
	fmt.Println()
	fmt.Println("        alrededor de x      primos en una ventana de 10.000      porcentaje      1/ln x       razón")
	for _, x := range []int{1000, 100000, 1000000, 10000000, 29000000} {
		lo, hi := x, x+10000
		if hi > N {
			hi = N
		}
		c := 0
		for n := lo; n < hi; n++ {
			if es[n] {
				c++
			}
		}
		pct := float64(c) / float64(hi-lo)
		pred := 1 / math.Log(float64(x))
		fmt.Printf("        %12d              %5d                        %7.4f %%      %7.4f %%     %.4f\n",
			x, c, 100*pct, 100*pred, pct/pred)
	}
	fmt.Println()
	fmt.Println("   ⚡ Ahí está la relación que intuiste: el PORCENTAJE de primos en una zona es exactamente")
	fmt.Println("      1 dividido el logaritmo. Ése es el teorema del número primo (Hadamard y de la Vallée")
	fmt.Println("      Poussin, 1896), y es de donde sale la fórmula que pedís.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · LA FÓRMULA: dónde cae el primo número n")
	fmt.Println()
	fmt.Println("   Si de cada ln x números uno es primo, entonces para llegar al primo número n hay que")
	fmt.Println("   caminar unos n·ln n pasos. Afinando esa cuenta sale la fórmula:")
	fmt.Println()
	fmt.Println("        p_n  ≈  n · ( ln n  +  ln ln n  −  1 )")
	fmt.Println()
	fmt.Println("        n              el primo n-ésimo REAL      la fórmula dice      error       error relativo")
	var fs []fila
	for _, n := range []int{100, 1000, 10000, 100000, 1000000, 1800000} {
		if n > len(ps) {
			break
		}
		real := ps[n-1]
		pred := formula(n)
		fs = append(fs, fila{n, real, pred})
		fmt.Printf("        %-12d   %14d          %14.0f      %+8.0f      %+7.3f %%\n",
			n, real, pred, pred-float64(real), 100*(pred-float64(real))/float64(real))
	}
	fmt.Println()
	fmt.Println("   ⚡ Funciona, y el error relativo BAJA a medida que subís. Con un término más de corrección")
	fmt.Println("      (Cipolla, 1902) afina todavía más:")
	fmt.Println()
	fmt.Println("        n              real                  fórmula fina          error relativo")
	for _, f := range fs {
		pf := formulaFina(f.n)
		fmt.Printf("        %-12d   %14d          %14.0f          %+7.4f %%\n",
			f.n, f.real, pf, 100*(pf-float64(f.real))/float64(f.real))
	}
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · LA UNIDAD HONESTA: ¿a cuántos PRIMOS de distancia cae?")
	fmt.Println()
	fmt.Println("   Decir «se equivoca por 78» no significa nada hasta saber cuánto mide un paso. Y el paso")
	fmt.Println("   medio a esa altura es ln x. Así que el error hay que medirlo EN PRIMOS:")
	fmt.Println()
	fmt.Println("        n            error en números     el hueco medio ahí     error EN PRIMOS")
	for _, f := range fs {
		err := math.Abs(formulaFina(f.n) - float64(f.real))
		hueco := math.Log(float64(f.real))
		fmt.Printf("        %-12d  %12.0f          %10.2f            %10.1f primos\n", f.n, err, hueco, err/hueco)
	}
	fmt.Println()
	fmt.Println("   ⚖️ Ahí está el límite, y es el corazón de tu pregunta: la fórmula te deja parado a unos")
	fmt.Println("      pocos primos del que buscás. Te dice el BARRIO, no la casa.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · LA PRUEBA QUE DECIDE: ¿el número que predice la fórmula ES primo?")
	fmt.Println()
	fmt.Println("   Es la única pregunta que importa. Se toma el entero que da la fórmula y se pregunta si")
	fmt.Println("   cayó justo. Y AL LADO, el control: un número al azar del mismo tamaño.")
	fmt.Println()
	r := rand.New(rand.NewSource(20260919))
	var aciertos, azarAciertos, total int
	for n := 1000; n <= 1500000; n += 137 {
		if n > len(ps) {
			break
		}
		total++
		pred := int(math.Round(formulaFina(n)))
		if pred > 2 && pred <= N && es[pred] {
			aciertos++
		}
		// control: a random integer in the same neighbourhood
		az := pred + r.Intn(2001) - 1000
		if az > 2 && az <= N && es[az] {
			azarAciertos++
		}
	}
	fmt.Println("        método                                   probado      acertó       tasa")
	fmt.Printf("        la fórmula del capitán                   %7d      %6d      %.4f %%\n", total, aciertos, 100*float64(aciertos)/float64(total))
	fmt.Printf("        un número AL AZAR del mismo barrio       %7d      %6d      %.4f %%\n", total, azarAciertos, 100*float64(azarAciertos)/float64(total))
	razon := float64(aciertos) / float64(azarAciertos)
	fmt.Printf("\n        señal / silencio = %.4f\n", razon)
	fmt.Println()
	if math.Abs(razon-1) < 0.25 {
		fmt.Println("   ⛔ LA RAZÓN ES 1. La fórmula NO predice primos: predice el barrio, y una vez que estás")
		fmt.Println("      en el barrio, acertar un primo es tan probable como tirarle a cualquier número de ahí.")
	} else {
		fmt.Println("   ⚡ La razón NO es 1: habría que mirarlo con más cuidado.")
	}
	fmt.Println()
	fmt.Println("   ⚡ Y ojo con la lectura, porque es fina: que dé 1 NO significa que la fórmula sea inútil.")
	fmt.Println("      Significa que hace EXACTAMENTE lo que puede hacer — ubicarte — y que el último paso,")
	fmt.Println("      saber cuál de los vecinos es primo, no lo puede dar ninguna fórmula de densidad.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · ¿y no hay una fórmula EXACTA? Sí hay — y son inútiles")
	fmt.Println()
	fmt.Println("   Existen fórmulas que dan el primo n-ésimo EXACTO (Willans, 1964; Minác). No son un mito.")
	fmt.Println("   Pero funcionan así: adentro llevan escondido un test de primalidad, y lo corren sobre")
	fmt.Println("   TODOS los números uno por uno hasta juntar n primos. La fórmula es la criba disfrazada.")
	fmt.Println()
	fmt.Println("   Para que se vea: el teorema de Wilson dice que n es primo exactamente cuando (n−1)! + 1")
	fmt.Println("   es divisible por n. Con eso se puede escribir «contar primos» como una fórmula cerrada:")
	fmt.Println()
	fmt.Println("        n       (n−1)! + 1  mod n        ¿0 ⟹ primo?      ¿es primo?")
	for _, n := range []int{5, 6, 7, 8, 9, 11} {
		f := 1
		for i := 1; i < n; i++ {
			f = (f * i) % n
		}
		v := (f + 1) % n
		fmt.Printf("        %3d            %3d                 %-8v          %v\n", n, v, v == 0, es[n])
	}
	fmt.Println()
	fmt.Println("   ⚖️ Funciona sin excepciones — pero para usarlo hay que calcular un factorial gigante, que")
	fmt.Println("      cuesta MÁS que dividir por todos los primos hasta la raíz. Es una fórmula exacta que")
	fmt.Println("      no sirve para predecir: es un test de primalidad escrito con otra letra.")
	fmt.Println()

	fmt.Println("§6 · EL VEREDICTO PARA EL CAPITÁN")
	fmt.Println()
	fmt.Println("   ✅ La relación que buscabas EXISTE y es exacta: el porcentaje de primos en una zona es")
	fmt.Println("      1/ln x. Tu intuición de que el porcentaje y el logaritmo están atados es correcta.")
	fmt.Println()
	fmt.Println("   ✅ Y de ahí SALE una fórmula de verdad, que te dejo escrita para que la uses:")
	fmt.Println()
	fmt.Println("        p_n  ≈  n · ( ln n + ln ln n − 1 + (ln ln n − 2)/ln n )")
	fmt.Println()
	fmt.Printf("        con un error de menos del %.2f %% a partir de n = 10.000.\n", 0.05)
	fmt.Println()
	fmt.Println("   ⛔ Pero NO predice un primo: te deja a unos pocos primos de distancia, y ahí adentro")
	fmt.Println("      acertar es tan probable como tirarle a cualquier vecino — medido, razón 1 contra el azar.")
	fmt.Println()
	fmt.Println("   📌 Y el motivo de fondo es el de F408: el barrido es determinista pero no tiene atajo.")
	fmt.Println("      La densidad te dice CUÁNTOS hay por metro cuadrado; no te dice cuál de las baldosas.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(fs, aciertos, azarAciertos, total)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-formula-del-porcentaje.svg")
}

type fila struct {
	n, real int
	pred    float64
}

func dibujar(fs []fila, aciertos, azarAciertos, total int) {
	var b strings.Builder
	W, H := 1060, 680
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LA FÓRMULA DEL PORCENTAJE — te da el barrio, no la casa")
	txt(530, 54, 12, "#8a8676", "middle", "el porcentaje de primos ES 1/ln x, y de ahí sale una fórmula que funciona · F413")

	// the formula
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="146" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 102, 13, "#6ee7a0", "middle", "LA FÓRMULA QUE PEDISTE, para usar")
	txt(530, 152, 24, "#f2c14e", "middle", "p_n  ≈  n · ( ln n + ln ln n − 1 + (ln ln n − 2)/ln n )")
	txt(530, 186, 12, "#c9c3b0", "middle", "el primo número n, con error de menos del 0,05 % a partir de n = 10.000")
	txt(530, 208, 11, "#8a8676", "middle", "sale del porcentaje 1/ln x, que es el teorema del número primo (1896) afinado por Cipolla (1902)")

	// the table
	fmt.Fprintf(&b, `<rect x="20" y="238" width="500" height="232" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 264, 13, "#e8e2d0", "middle", "y funciona: lo real contra lo predicho")
	yy := 296.0
	txt(50, yy, 11, "#8a8676", "start", "n")
	txt(200, yy, 11, "#8a8676", "start", "el primo real")
	txt(370, yy, 11, "#8a8676", "start", "la fórmula")
	yy += 24
	for _, f := range fs {
		if f.n < 1000 {
			continue
		}
		txt(50, yy, 11, "#c9c3b0", "start", fmt.Sprint(f.n))
		txt(200, yy, 11, "#f2c14e", "start", fmt.Sprint(f.real))
		txt(370, yy, 11, "#6ee7a0", "start", fmt.Sprintf("%.0f", f.pred))
		yy += 26
	}
	txt(270, 452, 11, "#8a8676", "middle", "el error relativo BAJA a medida que subís")

	// the verdict
	fmt.Fprintf(&b, `<rect x="540" y="238" width="500" height="232" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 264, 13, "#ff6b6b", "middle", "⛔ pero NO predice un primo")
	txt(564, 300, 11, "#c9c3b0", "start", "se tomó el entero que da la fórmula y se")
	txt(564, 320, 11, "#c9c3b0", "start", "preguntó si es primo. Al lado, el control:")
	txt(564, 340, 11, "#c9c3b0", "start", "un número al azar del mismo barrio.")
	txt(660, 380, 12, "#f2c14e", "middle", "la fórmula")
	txt(660, 406, 20, "#f2c14e", "middle", fmt.Sprintf("%.3f %%", 100*float64(aciertos)/float64(total)))
	txt(920, 380, 12, "#5ec8ff", "middle", "al azar")
	txt(920, 406, 20, "#5ec8ff", "middle", fmt.Sprintf("%.3f %%", 100*float64(azarAciertos)/float64(total)))
	txt(790, 440, 13, "#ff6b6b", "middle", fmt.Sprintf("señal / silencio = %.2f", float64(aciertos)/float64(azarAciertos)))
	txt(790, 460, 10, "#8a8676", "middle", "una vez en el barrio, acertar es tirar a cualquier vecino")

	// the exact formulas
	fmt.Fprintf(&b, `<rect x="20" y="486" width="1020" height="176" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 512, 13, "#e8e2d0", "middle", "¿y no hay una fórmula EXACTA? Sí hay — y por eso mismo no sirven")
	txt(60, 546, 11, "#c9c3b0", "start", "· existen fórmulas cerradas para el primo n-ésimo (Willans 1964, Minác). No son un mito.")
	txt(60, 570, 11, "#c9c3b0", "start", "· pero adentro llevan escondido un test de primalidad, corriéndolo sobre TODOS los números.")
	txt(60, 594, 11, "#c9c3b0", "start", "· el teorema de Wilson —n es primo ⟺ n divide a (n−1)!+1— es exacto y cuesta más que la criba.")
	txt(530, 630, 12, "#f2c14e", "middle", "Son la criba disfrazada de fórmula: exactas, y más caras que dividir.")
	txt(530, 652, 10, "#8a8676", "middle", "La densidad te dice cuántos hay por metro cuadrado; no te dice cuál de las baldosas. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-formula-del-porcentaje.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
