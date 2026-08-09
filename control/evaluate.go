package control

import "math"

// significanceThreshold is the number of standard deviations a measurement must
// clear before it counts. Five sigma is the particle-physics convention: strict
// enough that running many detectors over the same data does not manufacture
// discoveries out of noise.
const significanceThreshold = 5.0

// Result is the verdict on one measurement.
type Result struct {
	// Observed is what the detector found in the real sequence.
	Observed int
	// Mean and StdDev describe what the decoys produced.
	Mean   float64
	StdDev float64
	// Ratio is Observed / Mean. Above 1 is an excess, below 1 a deficit.
	Ratio float64
	// ZScore is how many decoy standard deviations separate Observed from Mean.
	ZScore float64
	// Trials is how many decoys were drawn.
	Trials int
}

// Significant reports whether the measurement clears the threshold in either
// direction. A deficit is as much a finding as an excess: both mean the real
// arrangement differs from a shuffled one.
//
// Zero spread never counts. If every decoy returned the same number the control
// had no resolution, and no separation can be claimed from it. A result built
// with no trials at all carries a zero spread too, so it is covered here.
func (r Result) Significant() bool {
	if r.StdDev == 0 {
		return false
	}
	return math.Abs(r.ZScore) > significanceThreshold
}

// Evaluate scores observed against trials draws of sample.
//
// sample is called exactly trials times and must return the detector's count on
// a freshly built decoy. With no trials the result is null: nothing was
// compared, so nothing can be claimed.
func Evaluate(observed, trials int, sample func() int) Result {
	if trials < 1 {
		return Result{Observed: observed}
	}

	draws := make([]float64, trials)
	sum := 0.0
	for i := range draws {
		v := float64(sample())
		draws[i] = v
		sum += v
	}
	mean := sum / float64(trials)

	variance := 0.0
	for _, v := range draws {
		variance += (v - mean) * (v - mean)
	}
	stdDev := math.Sqrt(variance / float64(trials))

	result := Result{
		Observed: observed,
		Mean:     mean,
		StdDev:   stdDev,
		Trials:   trials,
	}
	if mean != 0 {
		result.Ratio = float64(observed) / mean
	}
	if stdDev != 0 {
		result.ZScore = (float64(observed) - mean) / stdDev
	}

	return result
}
