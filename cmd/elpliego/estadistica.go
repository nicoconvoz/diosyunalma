package main

import (
	"math"
	"sort"
)

// estadistica.go - the statistics the adversarial review demanded: a bin-free
// distance (the canto exam is discontinuous on a degenerate spectrum, so a
// picket fence's verdict depends on where the bin edges fall), a two-sample
// test on the means (the exam's own per-realisation spread answers a different
// question), and the number variance - the two-level statistic that is NOT
// blind where the one-spacing exam is.

// wignerCDF integrates the GUE surmise to s on a fine grid.
func wignerCDF(s float64) float64 {
	if s <= 0 {
		return 0
	}
	const n = 4000
	h := s / n
	acc := 0.0
	for i := 0; i < n; i++ {
		acc += (wigner(float64(i)*h) + wigner(float64(i+1)*h)) / 2 * h
	}
	return acc
}

// ks is the Kolmogorov-Smirnov distance from a spacing sample to the Wigner
// law. Unlike canto it has no bins, so a degenerate spectrum gets one answer.
func ks(sp []float64) float64 {
	v := append([]float64(nil), sp...)
	sort.Float64s(v)
	n := float64(len(v))
	d := 0.0
	for i, s := range v {
		f := wignerCDF(s)
		if x := math.Abs(float64(i+1)/n - f); x > d {
			d = x
		}
		if x := math.Abs(f - float64(i)/n); x > d {
			d = x
		}
	}
	return d
}

// tDeMedias is the two-sample t on the means - the right question when asking
// whether two candidates differ, rather than whether one draw can be told apart.
func tDeMedias(a, b []float64) float64 {
	ma, mb := media(a), media(b)
	sa, sb := desvio(a), desvio(b)
	na, nb := float64(len(a)), float64(len(b))
	return math.Abs(ma-mb) / math.Sqrt(sa*sa/na+sb*sb/nb)
}

// varianzaNumero is Sigma^2(L): the variance of how many unfolded levels fall
// in a window of length L. It sees CORRELATIONS between levels, which no
// one-spacing statistic can. Rigid spectra give small values; a memoryless
// sequence with the same spacing law gives large ones.
func varianzaNumero(sp []float64, L float64, pasos int) float64 {
	x := make([]float64, len(sp)+1)
	for i, s := range sp {
		x[i+1] = x[i] + s
	}
	total := x[len(x)-1]
	if total <= L*1.2 {
		return math.NaN()
	}
	var ns []float64
	paso := (total - L) / float64(pasos)
	for k := 0; k < pasos; k++ {
		a := float64(k) * paso
		lo := sort.SearchFloat64s(x, a)
		hi := sort.SearchFloat64s(x, a+L)
		ns = append(ns, float64(hi-lo))
	}
	m := media(ns)
	v := 0.0
	for _, y := range ns {
		v += (y - m) * (y - m)
	}
	return v / float64(len(ns))
}

// normalizar divides a sample by its own mean - the preprocessing the GUE arm
// already had, applied to every arm so the comparison is matched.
func normalizar(sp []float64) []float64 {
	m := media(sp)
	out := make([]float64, len(sp))
	for i, s := range sp {
		out[i] = s / m
	}
	return out
}

// maxi is the largest of a sample.
func maxi(v []float64) float64 {
	m := v[0]
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}
