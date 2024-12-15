package multicell

import (
	"math"
)

// vector
type Vec []float64

func (vec Vec) At(i int) float64 {
	return vec[i]
}

func (vec Vec) SetAll(v float64) {
	for i := range vec {
		vec[i] = v
	}
}

// Create a vector with initial values of "v".
func NewVec(n int, v float64) Vec {
	vec := make(Vec, n)
	vec.SetAll(v)
	return vec
}

func (vout Vec) ApplyFVec(f func(float64) float64, vin Vec) {
	for i, v := range vin {
		vout[i] = f(v)
	}
}

func DotVecs(v0, v1 Vec) float64 {
	dot := 0.0
	for i, v := range v0 {
		dot += v * v1[i]
	}
	return dot
}

func (vin Vec) ScaleBy(f float64) {
	for i, v := range vin {
		vin[i] = f * v
	}
}

func (v Vec) Normalize() {
	mag := v.Norm2()
	v.ScaleBy(1 / mag)
}

func (vout Vec) Add(v0, v1 Vec) {
	for i, v := range v0 {
		vout[i] = v + v1[i]
	}
}

// Accumulate
func (vout Vec) Acc(vin Vec) {
	for i, v := range vin {
		vout[i] += v
	}
}

func (vout Vec) Diff(v0, v1 Vec) {
	for i, v := range v0 {
		vout[i] = v - v1[i]
	}
}

func (v Vec) Norm1() float64 {
	d := 0.0
	for _, x := range v {
		d += math.Abs(x)
	}
	return d
}

func (v Vec) Norm2() float64 {
	d := 0.0
	for _, x := range v {
		d += x * x
	}
	return math.Sqrt(d)
}
