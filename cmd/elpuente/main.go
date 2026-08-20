package main

// EL PUENTE - referee verification for Auditoria 64 §1. The bridge error
// E_n = T_n - T_n^M has EXACT leading terms by pure algebra: the model's
// closed form n^{-s}/(e^{i theta}-1) resums exactly the (it)-parts of the
// Euler-Maclaurin series of T_n, so the difference starts with
//   D1 = n^{1-s} [ 1/(s-1) - 1/(it) ]        (drift mismatch)
//   D2 = sigma n^{-s-1} / 12                 (sigma part of the B2 term)
// and the remaining tail is the series of sigma-corrections to higher
// Bernoulli terms (O(t^2 n^{-sigma-3}) leading). This program measures E_n
// directly and compares with D1+D2 in the direct zone theta <= pi.

import (
	"fmt"
	"math"
	"math/cmplx"
)

func main() {
	fmt.Println("🌉 EL PUENTE — E_n medido contra sus términos exactos D₁+D₂")
	fmt.Println()
	fmt.Println("   σ     θ      n      |E| medido    |D₁+D₂|     residuo |E−(D₁+D₂)|/|E|   |E|/|T^M|")
	t := 1000.0
	for _, sigma := range []float64{0.3, 0.5, 0.7} {
		s := complex(sigma, t)
		// center: EM with drift + half + B2 term at X = 6t (higher precision
		// than the experiments, so the bridge is not polluted by E_C)
		X := int(6 * t)
		var S complex128
		var Sn complex128
		ns := func(n float64) complex128 { // n^{-s}
			return cmplx.Exp(-s * complex(math.Log(n), 0))
		}
		for m := 1; m <= X; m++ {
			S += ns(float64(m))
		}
		fX := float64(X)
		C := S + ns(fX)*complex(fX, 0)/(s-1) - ns(fX)/2 + s*ns(fX)/complex(12*fX, 0)

		for _, th := range []float64{0.3, 1.0, 2.0, 3.0} {
			n := t / th
			nn := int(n)
			fn := float64(nn)
			thn := t / fn
			// partial sum up to nn
			Sn = 0
			for m := 1; m <= nn; m++ {
				Sn += ns(float64(m))
			}
			T := C - Sn
			TM := ns(fn) / (cmplx.Exp(complex(0, thn)) - 1)
			E := T - TM
			D1 := ns(fn) * complex(fn, 0) * (1/(s-1) - 1/complex(0, t))
			D2 := complex(sigma/12, 0) * ns(fn) / complex(fn, 0)
			pred := D1 + D2
			fmt.Printf("   %.1f  %4.1f  %5d    %.3e    %.3e         %.3f                %.5f\n",
				sigma, th, nn, cmplx.Abs(E), cmplx.Abs(pred),
				cmplx.Abs(E-pred)/cmplx.Abs(E), cmplx.Abs(E)/cmplx.Abs(TM))
		}
	}
	fmt.Println()
	fmt.Println("   lectura honesta: D₁+D₂ capturan el puente al 0,2–2,6% para θ ≤ 1 (n ≳ t);")
	fmt.Println("   hacia θ = 2 el residuo crece (24–66%) y en θ = 3 la cola de la serie")
	fmt.Println("   domina. CASO B: primer orden demostrado y verificado en su zona; la cola")
	fmt.Println("   (correcciones σ a los términos de Bernoulli superiores) queda abierta.")
}
