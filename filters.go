package dsp

import "math"

// NewLowPassFilter creates a new low-pass filter
func NewLowPassFilter(fC, fS float64) *Filter {
	wcT := 2 * math.Pi * fC / fS
	K := math.Tan(wcT / 2)
	K2 := K * K

	// all coeff denoms are the same
	denom := (1 + math.Sqrt2*K + K2)

	b0 := 1.0
	b1 := (2 * (K2 - 1)) / denom
	b2 := (1 - math.Sqrt2*K + K2) / denom
	a0 := K2 / denom
	a1 := 2 * K2 / denom
	a2 := K2 / denom

	A := []float64{a0, a1, a2}
	B := []float64{b0, b1, b2}

	return &Filter{B, A}
}

// NewHighPassFilter creates a new high-pass filter
func NewHighPassFilter(fC, fS float64) *Filter {
	wcT := 2 * math.Pi * fC / fS
	K := math.Tan(wcT / 2)
	K2 := K * K

	// all coeff denoms are the same
	denom := (1 + math.Sqrt2*K + K2)

	b0 := 1.0
	b1 := (2 * (K2 - 1)) / denom
	b2 := (1 - math.Sqrt2*K + K2) / denom
	a0 := 1 / denom
	a1 := -2 / denom
	a2 := 1 / denom

	A := []float64{a0, a1, a2}
	B := []float64{b0, b1, b2}

	return &Filter{B, A}
}

// NewBandPassFilter creates a new Band-pass filter
func NewBandPassFilter(fC, bw, fS float64) *Filter {
	Q := fS / bw
	wcT := 2 * math.Pi * fC / fS
	K := math.Tan(wcT / 2)
	K2 := K * K

	// all coeff denoms are the same
	denom := (1 + (1/Q)*K + K2)

	b0 := 1.0
	b1 := (2 * (K2 - 1)) / denom
	b2 := (1 - (1/Q)*K + K2) / denom
	a0 := (1 / Q) * K / denom
	a1 := 0.0
	a2 := (1 / Q) * K / denom

	A := []float64{a0, a1, a2}
	B := []float64{b0, b1, b2}

	return &Filter{B, A}
}

// Filter contains the coefficients for a filter.
type Filter struct {
	B, A []float64
}

// Filter executes the filter on the given data.
func (f Filter) Filter(X []float64) []float64 {
	n := len(f.A)
	z := make([]float64, n)
	Y := make([]float64, len(X))

	for m := 0; m < len(Y); m++ {
		Y[m] = f.A[0]*X[m] + z[0]

		for i := 1; i < n; i++ {
			z[i-1] = f.A[i]*X[m] + z[i] - f.B[i]*Y[m]
		}
	}
	return Y
}
