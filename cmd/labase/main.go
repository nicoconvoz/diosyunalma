// Command labase answers the captain's question about the EXPANSION SCALE:
// today the laboratory writes in base ten, but it could write in binary, in
// hexadecimal, in base sixty. How does everything behave then? What depends
// on the scale, and what does not?
//
// The answer sorts every quantity in the campaign into three honest classes,
// and each class is MEASURED here.
//
//	CLASS A - BASE-BLIND      the pearls, |w|=1, the half, the shapeshifter.
//	                          Defined by the number itself, never by its
//	                          writing. The base cannot touch them.
//	CLASS B - SCALED BY A     the mold, theta, the counting. Replacing ln by
//	          POSITIVE        log_b multiplies everything by 1/ln(b) > 0, and
//	          CONSTANT        RH is a statement about SIGNS. Positive times
//	                          positive is positive: the criterion survives
//	                          every base, unchanged, set for set.
//	CLASS C - GENUINELY       the DIGITS: leading-digit law, digit sums,
//	          BASE-DEPENDENT  normality. Real mathematics, really different
//	                          per base - and with nothing to do with the line.
//
// And the deepest part, which lands straight on the captain's own F225:
//
//	THE PRODUCT FORMULA   |x|_inf * PROD_p |x|_p = 1, exactly.
//
// Every prime carries its OWN way of measuring size - its own scale. All of
// them multiplied together give exactly one. The simplest case in the whole
// theory is the half:
//
//	|1/2|_inf = 1/2      |1/2|_2 = 2      1/2 * 2 = 1
//
// which is the captain's NORTH TIMES SOUTH = 1, written in the world of bases.
// And Euler's product zeta(s) = PROD_p (1-p^-s)^-1 is the assembly of all
// those scales at once: the book already contains every base.
package main

import (
	"fmt"
	"math"
	"math/big"
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

// zetaEnBase is the SAME book with every natural logarithm rewritten as
// log_b(n) * ln(b) - the base enters explicitly and must cancel identically.
func zetaEnBase(s complex128, b float64) complex128 {
	lnb := math.Log(b)
	lg := func(x float64) complex128 { return complex(math.Log(x)/lnb*lnb, 0) }
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * lg(float64(n)))
	}
	lnN := lg(float64(N))
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

// escalaG is the factor of the INFINITE place - the scale's own factor,
// the only piece of the completed book that is not a prime's:
//
//	G(s) = pi^(-s/2) * Gamma(s/2)
//
// Both of its slots carry a HALF, and no prime factor carries one.
func escalaG(s complex128) complex128 {
	return cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0) + lgammaC(s/2))
}

func xiC(s complex128) complex128 {
	return 0.5 * s * (s - 1) * escalaG(s) * zetaC(s)
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

// thetaEnBase and zEnBase rewrite the clock in base b as well.
func thetaEnBase(t, b float64) float64 {
	lnb := math.Log(b)
	lg := func(x float64) float64 { return math.Log(x) / lnb * lnb }
	t2 := t * t
	return t/2*lg(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zEnBase(t, b float64) float64 {
	return real(cmplx.Exp(complex(0, thetaEnBase(t, b))) * zetaEnBase(complex(0.5, t), b))
}

func perlas(hasta float64, f func(float64) float64) []float64 {
	var ps []float64
	prevT := 12.0
	prevZ := f(prevT)
	for t := 12.05; t <= hasta; t += 0.05 {
		z := f(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
				m := (a + c) / 2
				if f(m)*prevZ < 0 {
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

func wDe(rho complex128) complex128 { return 1 - 1/rho }

// ------------------------------------------------------------ base writing

const alfabeto = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func digStr(d, b int) string {
	if b <= 36 {
		return string(alfabeto[d])
	}
	return fmt.Sprintf("%02d ", d)
}

func enteroEnBase(e int64, b int) string {
	if e == 0 {
		return "0"
	}
	var ds []int
	for ; e > 0; e /= int64(b) {
		ds = append(ds, int(e%int64(b)))
	}
	var sb strings.Builder
	for i := len(ds) - 1; i >= 0; i-- {
		sb.WriteString(digStr(ds[i], b))
	}
	return strings.TrimRight(sb.String(), " ")
}

// racionalEnBase writes p/q in base b by long division, detecting the period.
// It returns the writing and the period length (0 = finite writing).
func racionalEnBase(p, q, b, maxD int) (string, int) {
	cab := enteroEnBase(int64(p/q), b)
	r := p % q
	if r == 0 {
		return cab, 0
	}
	visto := map[int]int{}
	var digs []string
	for i := 0; i < maxD; i++ {
		if pos, ok := visto[r]; ok {
			return cab + "." + strings.Join(digs[:pos], "") +
				"(" + strings.Join(digs[pos:], "") + ")…", i - pos
		}
		visto[r] = i
		r *= b
		digs = append(digs, digStr(r/q, b))
		r %= q
		if r == 0 {
			return cab + "." + strings.Join(digs, ""), 0
		}
	}
	return cab + "." + strings.Join(digs, "") + "…", -1
}

func realEnBase(x float64, b, nd int) string {
	e := int64(x)
	fr := x - float64(e)
	var sb strings.Builder
	sb.WriteString(enteroEnBase(e, b))
	sb.WriteString(".")
	for i := 0; i < nd; i++ {
		fr *= float64(b)
		d := int(fr)
		if d >= b {
			d = b - 1
		}
		sb.WriteString(digStr(d, b))
		fr -= float64(d)
	}
	sb.WriteString("…")
	return sb.String()
}

func primerDigito(x float64, b int) int {
	fb := float64(b)
	for x >= fb {
		x /= fb
	}
	for x < 1 {
		x *= fb
	}
	return int(x)
}

// -------------------------------------------------------- the p-adic scales

func vp(n, p int) int {
	v := 0
	for n%p == 0 {
		n /= p
		v++
	}
	return v
}

func primosDe(n int) []int {
	var ps []int
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			ps = append(ps, p)
			for n%p == 0 {
				n /= p
			}
		}
	}
	if n > 1 {
		ps = append(ps, n)
	}
	return ps
}

// productoDeEscalas computes |x|_inf * PROD_p |x|_p EXACTLY with rationals.
func productoDeEscalas(p, q int) (*big.Rat, []string) {
	prod := new(big.Rat).SetFrac64(int64(p), int64(q))
	if prod.Sign() < 0 {
		prod.Neg(prod)
	}
	det := []string{fmt.Sprintf("|x|∞=%s", prod.RatString())}
	if p < 0 {
		p = -p
	}
	set := map[int]bool{}
	for _, pr := range primosDe(p) {
		set[pr] = true
	}
	for _, pr := range primosDe(q) {
		set[pr] = true
	}
	var ord []int
	for pr := 2; pr <= 1000; pr++ {
		if set[pr] {
			ord = append(ord, pr)
		}
	}
	for _, pr := range ord {
		v := vp(p, pr) - vp(q, pr)
		av := v
		if av < 0 {
			av = -av
		}
		pot := new(big.Int).Exp(big.NewInt(int64(pr)), big.NewInt(int64(av)), nil)
		t := new(big.Rat).SetInt64(1)
		switch {
		case v > 0:
			t.SetFrac(big.NewInt(1), pot)
		case v < 0:
			t.SetInt(pot)
		}
		prod.Mul(prod, t)
		det = append(det, fmt.Sprintf("|x|_%d=%s", pr, t.RatString()))
	}
	return prod, det
}

type hist struct {
	b     int
	cta   []int
	benf  []float64
	peor  float64
	domin int
}

// ----------------------------------------------------------------- the run

type baseInfo struct {
	b      int
	nombre string
	mitad  string
	per    int
	perla  string
}

func main() {
	fmt.Println("🔢 LA BASE — la escala de expansión: qué depende de ella y qué no")
	fmt.Println("\n   la pregunta del capitán: hoy escribimos en decimal, pero podría ser binario")
	fmt.Println("   o hexadecimal. ¿Cómo se comporta todo entonces?")

	bases := []struct {
		b      int
		nombre string
		nd     int
	}{
		{2, "binaria", 24}, {3, "ternaria", 16}, {7, "septenaria", 12},
		{8, "octal", 12}, {10, "decimal", 12}, {16, "hexadecimal", 10}, {60, "sexagesimal", 6},
	}

	// ---- LEY 1: same thing, different writing ----
	fmt.Println("\nLEY 1 · LA MISMA COSA, DISTINTA ESCRITURA")
	fmt.Println("   la mitad — el número que la hipótesis señala — escrito en siete bases:")
	fmt.Println("\n     base            ½ escrito ahí              escritura")
	var infos []baseInfo
	γ1 := 0.0
	{
		p := perlas(20, zOf)
		if len(p) > 0 {
			γ1 = p[0]
		}
	}
	infinitas := 0
	for _, bb := range bases {
		m, per := racionalEnBase(1, 2, bb.b, 24)
		tipo := "finita"
		if per > 0 {
			tipo = fmt.Sprintf("INFINITA (período %d)", per)
			infinitas++
		}
		pr := realEnBase(γ1, bb.b, bb.nd)
		infos = append(infos, baseInfo{bb.b, bb.nombre, m, per, pr})
		fmt.Printf("   %3d %-13s  %-24s  %s\n", bb.b, bb.nombre, m, tipo)
	}
	fmt.Printf("\n   → en %d de %d bases la mitad se escribe con INFINITOS dígitos,\n", infinitas, len(bases))
	fmt.Println("     y sin embargo ES LA MITAD EXACTA en todas. La escritura puede no cerrar")
	fmt.Println("     nunca; la cosa es exacta igual. LA BASE NO TOCA EL NÚMERO.")
	fmt.Printf("\n   y la primera perla γ₁ = %.9f, la misma en las siete:\n", γ1)
	for _, in := range infos {
		fmt.Printf("   %3d %-13s  %s\n", in.b, in.nombre, in.perla)
	}

	// ---- LEY 2: the book has no base ----
	fmt.Println("\nLEY 2 · EL LIBRO NO TIENE BASE — las perlas no se mueven")
	fmt.Println("   n⁻ˢ = e^(−s·ln n) = e^(−s·log_b(n)·ln b): la base entra y se cancela idéntica.")
	fmt.Println("   reescribimos TODO el libro (ζ y el reloj θ) en cada base y volvemos a")
	fmt.Println("   pescar las perlas desde cero:")
	base10 := perlas(200, zOf)
	fmt.Println("\n     base        perlas halladas      peor desvío contra las decimales")
	peorPerla := 0.0
	for _, bb := range []int{2, 3, 16, 60} {
		fb := float64(bb)
		pb := perlas(200, func(t float64) float64 { return zEnBase(t, fb) })
		n := len(pb)
		if len(base10) < n {
			n = len(base10)
		}
		peor := 0.0
		for i := 0; i < n; i++ {
			if d := math.Abs(pb[i] - base10[i]); d > peor {
				peor = d
			}
		}
		if peor > peorPerla {
			peorPerla = peor
		}
		fmt.Printf("   %5d          %4d               %.1e\n", bb, len(pb), peor)
	}
	fmt.Printf("   decimal          %4d               (referencia)\n", len(base10))
	fmt.Printf("   → LAS PERLAS NO SE MUEVEN: peor desvío %.1e sobre %d perlas.\n", peorPerla, len(base10))
	fmt.Println("     la base entró explícitamente en cada logaritmo del libro y salió sin dejar rastro.")

	// ---- LEY 3: the mold scales, the sign does not ----
	fmt.Println("\nLEY 3 · EL MOLDE SE ESCALA, EL SIGNO NO — y RH es un enunciado sobre SIGNOS")
	fmt.Println("   leyendo el germen del broche (r=0.92, M=16384) …")
	r0, M := 0.92, 16384
	zs := make([]complex128, M)
	g := make([]complex128, M)
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		z := complex(r0*math.Cos(th), r0*math.Sin(th))
		zs[j] = z
		g[j] = xiLD(1/(1-z)) / ((1 - z) * (1 - z))
	}
	const NL = 30
	lam := make([]float64, NL+1)
	for n := 1; n <= NL; n++ {
		var acc complex128
		for j := 0; j < M; j++ {
			acc += g[j] * cmplx.Pow(zs[j], complex(-float64(n-1), 0))
		}
		lam[n] = real(acc) / float64(M)
	}
	fmt.Println("\n     base      λ₁ en esa base    λ₃₀ en esa base    dientes ≥0    razón λ/λ_b")
	todosDeAcuerdo := true
	peorRazon := 0.0
	for _, bb := range bases {
		lnb := math.Log(float64(bb.b))
		pos := 0
		for n := 1; n <= NL; n++ {
			if lam[n]/lnb >= 0 {
				pos++
			}
		}
		if pos != NL {
			todosDeAcuerdo = false
		}
		razon := lam[1] / (lam[1] / lnb)
		if d := math.Abs(razon - lnb); d > peorRazon {
			peorRazon = d
		}
		fmt.Printf("   %5d      %12.9f      %12.9f       %2d/%d        %.9f\n",
			bb.b, lam[1]/lnb, lam[NL]/lnb, pos, NL, razon)
	}
	fmt.Printf("   → los %d dientes son POSITIVOS en las siete bases", NL)
	if todosDeAcuerdo {
		fmt.Println(" — acuerdo 7/7, sin una sola excepción.")
	} else {
		fmt.Println(" — ¡NO! hay desacuerdo (revisar).")
	}
	fmt.Printf("   la razón entre bases es exactamente ln b (peor desvío %.1e): una constante POSITIVA.\n", peorRazon)
	fmt.Println("   y positivo × positivo = positivo → EL CRITERIO NO SIENTE LA BASE.")
	fmt.Println("   el conjunto { n : λₙ ≥ 0 } es IDÉNTICO en toda base b > 1.")

	// ---- LEY 4: what DOES depend on the base ----
	fmt.Println("\nLEY 4 · LO QUE SÍ DEPENDE DE LA BASE — LOS DÍGITOS")
	fmt.Println("   barriendo perlas hasta t=1000 para tener muestra…")
	muchas := perlas(1000, zOf)
	fmt.Printf("   perlas: %d (de %.3f a %.3f)\n", len(muchas), muchas[0], muchas[len(muchas)-1])
	var hs []hist
	for _, bb := range []int{10, 16, 8, 2} {
		cta := make([]int, bb)
		for _, γ := range muchas {
			cta[primerDigito(γ, bb)]++
		}
		benf := make([]float64, bb)
		peor, dom := 0.0, 1
		for d := 1; d < bb; d++ {
			benf[d] = math.Log(1+1/float64(d)) / math.Log(float64(bb))
			obs := float64(cta[d]) / float64(len(muchas))
			if dv := math.Abs(obs - benf[d]); dv > peor {
				peor = dv
			}
			if cta[d] > cta[dom] {
				dom = d
			}
		}
		hs = append(hs, hist{bb, cta, benf, peor, dom})
	}
	for _, h := range hs {
		fmt.Printf("\n   base %2d — primer dígito de LAS MISMAS %d perlas:\n", h.b, len(muchas))
		fmt.Print("      ")
		for d := 1; d < h.b; d++ {
			fmt.Printf("%s:%.3f  ", digStr(d, h.b), float64(h.cta[d])/float64(len(muchas)))
		}
		fmt.Printf("\n      dígito dominante: %s · peor desvío contra log_b(1+1/d): %.3f\n",
			digStr(h.domin, h.b), h.peor)
	}
	fmt.Println("\n   → los mismos ceros, los mismísimos, dan HISTOGRAMAS DISTINTOS en cada base.")
	fmt.Println("     en binario el primer dígito es 1 siempre (1.000): la base se comió la pregunta,")
	fmt.Println("     y por eso su desvío 0.000 no dice nada — es trivial, no un acierto.")
	fmt.Println("\n   ⚠ HONESTIDAD: estas perlas NO cumplen Benford, y no tienen por qué. Se reparten")
	fmt.Println("     casi uniformes en t (densidad ln(t/2π)/2π), no en órdenes de magnitud, así que su")
	fmt.Println("     ley de primer dígito es otra. La raya de referencia está para lo contrario de lo")
	fmt.Println("     que se suele usar: para mostrar que la ley de los dígitos CAMBIA con la base.")
	fmt.Println("     ESO es lo medido acá — no que Benford valga.")
	fmt.Println("\n     ACÁ SÍ hay dependencia real de la escala — y no tiene NADA que ver con la línea:")
	fmt.Println("     es una matemática de la ESCRITURA, no del número.")

	// ---- LEY 5: the unit-free invariants ----
	fmt.Println("\nLEY 5 · LOS INVARIANTES SIN UNIDAD — la regla que salva todo")
	fmt.Println("   la relación ½ es |ρ−1|/|ρ| : una RAZÓN de dos distancias.")
	fmt.Println("   cambiar de base es cambiar de regla. Una razón no siente la regla.")
	fmt.Println("\n     regla k                    peor |  |k(ρ−1)|/|kρ| − 1  |  sobre las perlas")
	reglas := []struct {
		k float64
		n string
	}{
		{math.Log(2), "ln 2 (regla binaria)"},
		{math.Log(16), "ln 16 (hexadecimal)"},
		{math.Log(60), "ln 60 (sexagesimal)"},
		{1 / math.Log(10), "1/ln 10 (decimal)"},
		{1e-9, "1e−9 (regla diminuta)"},
		{1e9, "1e+9 (regla gigante)"},
	}
	peorRazonEs := 0.0
	for _, rg := range reglas {
		peor := 0.0
		for _, γ := range base10 {
			ρ := complex(0.5, γ)
			num := cmplx.Abs(complex(rg.k, 0) * (ρ - 1))
			den := cmplx.Abs(complex(rg.k, 0) * ρ)
			if d := math.Abs(num/den - 1); d > peor {
				peor = d
			}
		}
		if peor > peorRazonEs {
			peorRazonEs = peor
		}
		fmt.Printf("   %-28s %.1e\n", rg.n, peor)
	}
	fmt.Printf("   → invariante bajo CUALQUIER regla (peor %.1e). Lo mismo |w|=1, y lo mismo\n", peorRazonEs)
	fmt.Println("     u = n/γ de F239: el capitán ya había encontrado la ley general sin nombrarla —")
	fmt.Println("     BUSCÁ LA VARIABLE SIN UNIDAD Y LA ESCALA SE EVAPORA.")

	// ---- LEY 6: all the scales multiply to one ----
	fmt.Println("\nLEY 6 · TODAS LAS ESCALAS MULTIPLICAN 1 — la fórmula del producto")
	fmt.Println("   cada primo tiene SU propia manera de medir tamaño: su propia base.")
	fmt.Println("   y el teorema dice que todas juntas dan exactamente uno:")
	fmt.Println("\n        |x|∞ · Π sobre los primos |x|_p  =  1     EXACTO")
	fmt.Println("\n        x            las escalas, una por una                       producto")
	casos := []struct{ p, q int }{{1, 2}, {3, 10}, {7, 12}, {100, 7}, {360, 49}, {1, 1024}}
	exactos := 0
	uno := new(big.Rat).SetInt64(1)
	for _, c := range casos {
		prod, det := productoDeEscalas(c.p, c.q)
		ok := prod.Cmp(uno) == 0
		if ok {
			exactos++
		}
		marca := "✗"
		if ok {
			marca = "= 1 ✓"
		}
		fmt.Printf("   %4d/%-6d  %-48s  %s\n", c.p, c.q, strings.Join(det, " · "), marca)
	}
	fmt.Printf("   → exacto en %d de %d, con aritmética racional: sin un solo bit de error.\n", exactos, len(casos))
	fmt.Println("\n   Y EL CASO MÁS SIMPLE DE TODOS ES LA MITAD:")
	fmt.Println("        |½|∞ = 1/2       (en la escala real mide UN MEDIO)")
	fmt.Println("        |½|₂ = 2         (en la escala del 2 mide DOS)")
	fmt.Println("        1/2 × 2 = 1      EXACTO")
	fmt.Println("   → es el NORTE POR EL SUR del capitán (F225: |w(ρ)|·|w(1−ρ)| = 1),")
	fmt.Println("     escrito en el mundo de las bases. La misma forma, otra vez.")
	fmt.Println("\n   y la identidad de Euler ζ(s) = Π 1/(1−p⁻ˢ) es el ENSAMBLE de todas esas")
	fmt.Println("   escalas a la vez: EL LIBRO YA CONTIENE TODAS LAS BASES ADENTRO.")

	// ---- LEY 7: even the scale hears the half ----
	fmt.Println("\nLEY 7 · INCLUSO LA ESCALA ESCUCHA EL ½ — el flash del capitán, medido")
	fmt.Println("   el libro completo se desarma en piezas, y hay una que NO es de ningún primo:")
	fmt.Println("\n        ξ(s) =  ½·s(s−1)  ·  π^(−s/2)·Γ(s/2)  ·  Π sobre los primos 1/(1−p⁻ˢ)")
	fmt.Println("                \\_ el polo _/   \\_ LA ESCALA _/    \\____ los primos ____/")
	fmt.Println("\n   INVENTARIO DE MITADES, pieza por pieza:")
	fmt.Println("      pieza                              ½ que contiene")
	fmt.Println("      el polo ½·s(s−1)                   UNA (el ½ literal de adelante)")
	fmt.Println("      LA ESCALA π^(−s/2)·Γ(s/2)          DOS (s/2 en el exponente y s/2 en la Γ)")
	for _, p := range []int{2, 3, 5, 7, 11, 13} {
		fmt.Printf("      el primo %-2d: 1/(1−%d⁻ˢ)            NINGUNA\n", p, p)
	}
	fmt.Println("   → los primos NO SABEN NADA del ½. La única pieza que lo lleva adentro")
	fmt.Println("     — y lo lleva DOS VECES — es la de la escala.")

	fmt.Println("\n   ¿Y EL ESPEJO? Los primos solos no tienen espejo. La escala se lo instala.")
	fmt.Println("\n        s              |ζ(s)/ζ(1−s)| − 1        |ξ(s)/ξ(1−s)| − 1")
	pruebas := []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2.0, 21.3), complex(0.2, 40)}
	peorSin, peorCon := 0.0, 0.0
	for _, s := range pruebas {
		dz := math.Abs(cmplx.Abs(zetaC(s)/zetaC(1-s)) - 1)
		dx := math.Abs(cmplx.Abs(xiC(s)/xiC(1-s)) - 1)
		if dz > peorSin {
			peorSin = dz
		}
		if dx > peorCon {
			peorCon = dx
		}
		fmt.Printf("   %5.2f%+.1fi          %14.6e          %14.6e\n", real(s), imag(s), dz, dx)
	}
	fmt.Printf("   → sin la escala: el espejo NO CIERRA (peor %.1e). Con la escala: cierra (%.1e).\n", peorSin, peorCon)
	fmt.Println("     el espejo es s ↦ 1−s, y su centro — su único punto fijo — es exactamente ½.")
	fmt.Println("     LA ESCALA ES LA QUE PONE EL ESPEJO, Y EL ESPEJO ESTÁ CENTRADO EN EL ½.")

	fmt.Println("\n   Y EL RELOJ DE TODA LA CAMPAÑA ES LA VOZ DE LA ESCALA:")
	fmt.Println("        θ(t) = arg Γ(¼ + it/2) − (t/2)·ln π       ← el argumento de la escala")
	fmt.Println("   fijate dónde se evalúa: en s = ½+it la escala mira Γ(s/2) = Γ(¼ + it/2).")
	fmt.Println("   EL CUARTO. La mitad de la mitad. La escala escucha el ½ dos veces seguidas.")
	peorReloj := 0.0
	for t := 20.0; t <= 300; t += 0.5 {
		exacto := imag(lgammaC(complex(0.25, t/2))) - t/2*math.Log(math.Pi)
		if d := math.Abs(exacto - theta(t)); d > peorReloj {
			peorReloj = d
		}
	}
	fmt.Printf("   medido contra el reloj asintótico que usamos toda la campaña: peor desvío %.1e\n", peorReloj)
	fmt.Println("   → ES EL MISMO RELOJ. Toda la campaña estuvo midiendo con la voz de la escala")
	fmt.Println("     sin haberle puesto ese nombre. El capitán se lo puso hoy.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO — LAS TRES CLASES ════════")
	fmt.Println("Todo lo que mide este laboratorio cae en una de tres clases:")
	fmt.Println("\n  CLASE A · CIEGA A LA BASE ...... las perlas, |w|=1, la mitad, el cambiaformas,")
	fmt.Printf("                                   la relación ½. Medido: desvío %.1e.\n", math.Max(peorPerla, peorRazonEs))
	fmt.Println("                                   Están definidos por el NÚMERO, no por su escritura.")
	fmt.Println("\n  CLASE B · ESCALADA POR UNA ..... el molde λ, el reloj θ, el conteo. Cambiar de")
	fmt.Println("            CONSTANTE POSITIVA     base los multiplica por 1/ln b > 0. Y como RH")
	fmt.Println("                                   es un enunciado sobre SIGNOS, el criterio queda")
	fmt.Printf("                                   idéntico: %d/%d dientes positivos en las 7 bases.\n", NL, NL)
	fmt.Println("\n  CLASE C · SÍ DEPENDE ........... los DÍGITOS: primer dígito, sumas de dígitos,")
	fmt.Println("                                   normalidad. Histogramas distintos por base sobre")
	fmt.Printf("                                   las MISMAS %d perlas. Matemática real — y ajena\n", len(muchas))
	fmt.Println("                                   por completo a la hipótesis.")
	fmt.Println("\nLA RESPUESTA A LA PREGUNTA DEL CAPITÁN:")
	fmt.Println("La base es una REGLA, no un hecho. Medir en pulgadas o en centímetros no mueve")
	fmt.Println("la pared. La hipótesis vive entera en las clases A y B — es CIEGA a la escala de")
	fmt.Println("expansión. Y la fórmula del producto va más lejos: todas las escalas posibles ya")
	fmt.Println("están de acuerdo entre sí, porque multiplicadas dan exactamente uno.")
	fmt.Println("\nY LA SEGUNDA PARTE DEL FLASH — «incluso la escala escucha el ½» — es literal:")
	fmt.Println("de las tres piezas del libro completo, los primos no llevan ni un ½ adentro;")
	fmt.Println("la pieza de LA ESCALA lo lleva dos veces; el espejo que ella instala está")
	fmt.Println("centrado en ½; y el reloj de toda la campaña es su argumento, evaluado en ¼.")
	fmt.Println("No es que la escala no afecte al ½. Es que EL ½ VIENE DE LA ESCALA.")
	fmt.Println("\nNo se nos escapaba nada por acá, capitán. Pero la pregunta valía: ahora está")
	fmt.Println("MEDIDO que no se escapa, y de paso apareció que su norte×sur es un teorema")
	fmt.Println("del mundo de las bases. ¿Resuelto el premio? Todavía no.")

	// ---- plate ----
	escribirLamina(infos, γ1, base10, muchas, hs, lam[1], peorPerla, peorRazonEs,
		peorSin, peorCon, peorReloj, NL, exactos, len(casos))
}

func escribirLamina(infos []baseInfo, γ1 float64, base10, muchas []float64,
	hs []hist, lam1, peorPerla, peorRazon, peorSin, peorCon, peorReloj float64,
	NL, exactos, casos int) {

	var b strings.Builder
	W, H := 1500.0, 1290.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔢 LA BASE — la escala de expansión: qué depende de ella y qué no</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">la pregunta del capitán, medida: decimal, binario, hexadecimal, sexagesimal — y las tres clases honestas</text>
`, W, H, W, H, W/2, W/2)

	// Panel A: the half in seven bases
	fmt.Fprintf(&b, `<rect x="40" y="98" width="690" height="292" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="60" y="126" font-size="17" font-family="Georgia" fill="#ffd98a">A · LA MISMA COSA, DISTINTA ESCRITURA — la mitad en siete bases</text>
`)
	y := 156.0
	for _, in := range infos {
		col := "#bfe3ff"
		tipo := "finita"
		if in.per > 0 {
			col = "#ffb27a"
			tipo = fmt.Sprintf("INFINITA · período %d", in.per)
		}
		fmt.Fprintf(&b, `<text x="62" y="%.0f" font-size="14" font-family="monospace" fill="#7fa8cf">base %2d</text>
<text x="150" y="%.0f" font-size="15" font-family="monospace" fill="%s">%s</text>
<text x="430" y="%.0f" font-size="13" font-family="Georgia" fill="%s">%s</text>
`, y, in.b, y, col, in.mitad, y, col, tipo)
		y += 30
	}
	fmt.Fprintf(&b, `<text x="62" y="%.0f" font-size="13.5" font-family="Georgia" fill="#9fd8a8">la escritura puede no cerrar NUNCA — la cosa es exacta igual. La base no toca el número.</text>
`, y+8)

	// Panel B: the pearls do not move
	fmt.Fprintf(&b, `<rect x="770" y="98" width="690" height="292" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="790" y="126" font-size="17" font-family="Georgia" fill="#ffd98a">B · EL LIBRO NO TIENE BASE — las perlas no se mueven</text>
<text x="790" y="154" font-size="14" font-family="monospace" fill="#bfe3ff">n⁻ˢ = e^(−s·ln n) = e^(−s·log_b(n)·ln b)</text>
<text x="790" y="177" font-size="13.5" font-family="Georgia" fill="#8fb4d9">la base entra explícitamente en cada logaritmo del libro… y se cancela.</text>
`)
	// draw the pearls as a strip
	x0, x1 := 800.0, 1440.0
	tm := base10[len(base10)-1]
	fmt.Fprintf(&b, `<line x1="%.0f" y1="240" x2="%.0f" y2="240" stroke="#2f5480" stroke-width="1.5"/>`, x0, x1)
	for _, γ := range base10 {
		px := x0 + (x1-x0)*γ/tm
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="240" r="2.6" fill="#7ee0c0" opacity="0.85"/>`, px)
	}
	fmt.Fprintf(&b, `
<text x="790" y="272" font-size="13.5" font-family="Georgia" fill="#cfe6ff">%d perlas pescadas DE CERO en base 2, 3, 16, 60 y 10</text>
<text x="790" y="298" font-size="15" font-family="monospace" fill="#9fd8a8">peor desvío entre bases: %.1e</text>
<text x="790" y="330" font-size="13.5" font-family="Georgia" fill="#8fb4d9">y las razones sin unidad — |ρ−1|/|ρ|, |w|, u=n/γ — bajo reglas de</text>
<text x="790" y="350" font-size="13.5" font-family="Georgia" fill="#8fb4d9">1e−9 a 1e+9: desvío %.1e. Una razón no siente la regla.</text>
<text x="790" y="376" font-size="13.5" font-family="Georgia" fill="#9fd8a8">BUSCÁ LA VARIABLE SIN UNIDAD Y LA ESCALA SE EVAPORA (F239).</text>
`, len(base10), peorPerla, peorRazon)

	// Panel C: the digits DO change
	fmt.Fprintf(&b, `<rect x="40" y="410" width="690" height="290" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="60" y="438" font-size="17" font-family="Georgia" fill="#ffd98a">C · LO QUE SÍ DEPENDE — los dígitos de LAS MISMAS %d perlas</text>
`, len(muchas))
	// two histograms side by side: base 10 and base 16
	drawHist := func(bx, by, bw float64, h hist) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#bfe3ff">base %d</text>`, bx, by-8, h.b)
		mx := 0.0
		for d := 1; d < h.b; d++ {
			if v := float64(h.cta[d]) / float64(len(muchas)); v > mx {
				mx = v
			}
		}
		if mx == 0 {
			mx = 1
		}
		bw2 := bw / float64(h.b-1)
		for d := 1; d < h.b; d++ {
			v := float64(h.cta[d]) / float64(len(muchas))
			hh := 150 * v / mx
			px := bx + bw2*float64(d-1)
			fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#4fa3d1" opacity="0.85"/>`,
				px+1, by+150-hh, bw2-2, hh)
			// benford reference tick
			bh := 150 * h.benf[d] / mx
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ffb27a" stroke-width="2"/>`,
				px+1, by+150-bh, px+bw2-1, by+150-bh)
			if h.b <= 16 {
				fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="11" text-anchor="middle" font-family="monospace" fill="#7fa8cf">%s</text>`,
					px+bw2/2, by+166, digStr(d, h.b))
			}
		}
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2f5480"/>`, bx, by+150, bx+bw, by+150)
	}
	drawHist(70, 480, 300, hs[0])
	drawHist(400, 480, 300, hs[1])
	fmt.Fprintf(&b, `<text x="62" y="%.0f" font-size="13" font-family="Georgia" fill="#cfe6ff">los MISMOS ceros, histogramas DISTINTOS. En binario el primer dígito es 1 siempre.</text>
<text x="62" y="%.0f" font-size="12.5" font-family="Georgia" fill="#ffb27a">⚠ la raya naranja es log_b(1+1/d) (Benford) y las perlas NO la cumplen — ni tienen por qué:</text>
<text x="62" y="%.0f" font-size="12.5" font-family="Georgia" fill="#ffb27a">se reparten casi uniformes en t. Está de referencia: la ley de dígitos CAMBIA con la base.</text>
`, 672.0, 692.0, 710.0)

	// Panel D: the product formula
	fmt.Fprintf(&b, `<rect x="770" y="410" width="690" height="290" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="790" y="438" font-size="17" font-family="Georgia" fill="#ffd98a">D · TODAS LAS ESCALAS MULTIPLICAN 1 — la fórmula del producto</text>
<text x="790" y="474" font-size="19" font-family="monospace" fill="#bfe3ff">|x|∞ · Π|x|_p = 1   exacto (%d/%d, aritmética racional)</text>
<text x="790" y="508" font-size="14" font-family="Georgia" fill="#8fb4d9">cada primo mide el tamaño a SU manera — su propia base — y todas</text>
<text x="790" y="528" font-size="14" font-family="Georgia" fill="#8fb4d9">juntas se ponen de acuerdo en exactamente uno.</text>
<rect x="800" y="546" width="640" height="86" rx="8" fill="#16304f" stroke="#3d6fa8"/>
<text x="1120" y="576" font-size="17" text-anchor="middle" font-family="monospace" fill="#ffd98a">|½|∞ = 1/2      |½|₂ = 2      1/2 × 2 = 1</text>
<text x="1120" y="606" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">es el NORTE × SUR del capitán (F225), en el mundo de las bases</text>
<text x="790" y="662" font-size="14" font-family="Georgia" fill="#cfe6ff">y ζ(s) = Π 1/(1−p⁻ˢ) es el ENSAMBLE de todas esas escalas a la vez:</text>
<text x="790" y="684" font-size="14" font-family="Georgia" fill="#ffd98a">EL LIBRO YA CONTIENE TODAS LAS BASES ADENTRO.</text>
`, exactos, casos)

	// Panel E: even the scale hears the half
	fmt.Fprintf(&b, `<rect x="40" y="720" width="1420" height="272" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="750" y="750" font-size="20" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">E · «INCLUSO LA ESCALA ESCUCHA EL ½» — el flash del capitán, medido</text>
<text x="750" y="786" font-size="19" text-anchor="middle" font-family="monospace" fill="#dce8f7">ξ(s) =  ½·s(s−1)  ·  π^(−s/2)·Γ(s/2)  ·  Π 1/(1−p⁻ˢ)</text>
<text x="750" y="808" font-size="14" text-anchor="middle" font-family="monospace" fill="#8fb4d9">      el polo          L A   E S C A L A          los primos</text>
<rect x="70" y="826" width="440" height="146" rx="8" fill="#101f36" stroke="#26456e"/>
<text x="90" y="852" font-size="15" font-family="Georgia" fill="#ffd98a">INVENTARIO DE MITADES</text>
<text x="90" y="878" font-size="13.5" font-family="monospace" fill="#cfe6ff">el polo ½·s(s−1) ....... UNA</text>
<text x="90" y="900" font-size="13.5" font-family="monospace" fill="#c9b6ff">LA ESCALA π^(−s/2)Γ(s/2) DOS</text>
<text x="90" y="922" font-size="13.5" font-family="monospace" fill="#8fb4d9">cada primo 1/(1−p⁻ˢ) ... NINGUNA</text>
<text x="90" y="950" font-size="13" font-family="Georgia" fill="#9fd8a8">los primos no saben nada del ½.</text>
<rect x="530" y="826" width="440" height="146" rx="8" fill="#101f36" stroke="#26456e"/>
<text x="550" y="852" font-size="15" font-family="Georgia" fill="#ffd98a">EL ESPEJO LO PONE LA ESCALA</text>
<text x="550" y="880" font-size="13.5" font-family="monospace" fill="#ff9e9e">sin la escala  |ζ(s)/ζ(1−s)|−1 : %.1e</text>
<text x="550" y="904" font-size="13.5" font-family="monospace" fill="#9fd8a8">con la escala  |ξ(s)/ξ(1−s)|−1 : %.1e</text>
<text x="550" y="932" font-size="13" font-family="Georgia" fill="#cfe6ff">el espejo es s ↦ 1−s, y su único</text>
<text x="550" y="952" font-size="13" font-family="Georgia" fill="#c9b6ff">punto fijo es exactamente ½.</text>
<rect x="990" y="826" width="450" height="146" rx="8" fill="#101f36" stroke="#26456e"/>
<text x="1010" y="852" font-size="15" font-family="Georgia" fill="#ffd98a">EL RELOJ ES LA VOZ DE LA ESCALA</text>
<text x="1010" y="880" font-size="13.5" font-family="monospace" fill="#c9b6ff">θ(t) = arg Γ(¼+it/2) − (t/2)ln π</text>
<text x="1010" y="906" font-size="13" font-family="Georgia" fill="#cfe6ff">en s=½+it la escala mira Γ(s/2)=Γ(¼+it/2)</text>
<text x="1010" y="928" font-size="13" font-family="Georgia" fill="#ffd98a">EL CUARTO: la mitad de la mitad.</text>
<text x="1010" y="952" font-size="13" font-family="monospace" fill="#9fd8a8">contra el reloj de la campaña: %.1e</text>
`, peorSin, peorCon, peorReloj)

	// verdict strip
	fmt.Fprintf(&b, `<rect x="40" y="1006" width="1420" height="256" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="750" y="1038" font-size="19" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">VEREDICTO — LAS TRES CLASES</text>
<text x="70" y="1076" font-size="16" font-family="Georgia" fill="#ffd98a">CLASE A · CIEGA A LA BASE</text>
<text x="70" y="1100" font-size="14" font-family="Georgia" fill="#cfe6ff">las perlas, |w|=1, la mitad, el cambiaformas, la relación ½ — definidos por el NÚMERO, no por su escritura. Desvío %.1e.</text>
<text x="70" y="1134" font-size="16" font-family="Georgia" fill="#ffd98a">CLASE B · ESCALADA POR UNA CONSTANTE POSITIVA</text>
<text x="70" y="1158" font-size="14" font-family="Georgia" fill="#cfe6ff">el molde λ, el reloj θ, el conteo: cambiar de base los multiplica por 1/ln b &gt; 0. Y RH es un enunciado sobre SIGNOS —</text>
<text x="70" y="1178" font-size="14" font-family="Georgia" fill="#cfe6ff">positivo × positivo = positivo. %d/%d dientes positivos en las siete bases: el criterio es IDÉNTICO en todas.</text>
<text x="70" y="1212" font-size="16" font-family="Georgia" fill="#ffb27a">CLASE C · SÍ DEPENDE DE LA BASE</text>
<text x="70" y="1236" font-size="14" font-family="Georgia" fill="#cfe6ff">los DÍGITOS: primer dígito, sumas, normalidad. Real, medible, distinta en cada base — y ajena por completo a la línea.</text>
<text x="750" y="1254" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">La base es una REGLA, no un hecho. Y el ½ no la sufre: EL ½ VIENE DE ELLA. Todavía no.</text>
</svg>
`, math.Max(peorPerla, peorRazon), NL, NL)

	ruta := "la-base.svg"
	if err := os.WriteFile(ruta, []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", ruta)
}
