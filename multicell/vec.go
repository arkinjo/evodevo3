package multicell

import (
	"math"
	"math/rand/v2"
	"slices"
)

// vector
type Vec []float64

/*
This is different from:	type Vec = []float64

	Note the "=" here! For the latter, methods cannot be defined.
*/

func (vec Vec) SetAll(v float64) {
	for i := range vec {
		vec[i] = v
	}
}

func (vec Vec) Clone() Vec {
	vout := make(Vec, len(vec))
	copy(vout, vec)
	return vout
}

// sum all elements.
func (vec Vec) Sum() float64 {
	s := 0.0
	for _, v := range vec {
		s += v
	}
	return s
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

func (v0 Vec) MateWith(v1 Vec) (Vec, Vec) {
	nv0 := slices.Clone(v0)
	nv1 := slices.Clone(v1)

	for i, v := range nv0 {
		if rand.IntN(2) == 1 {
			nv0[i] = nv1[i]
			nv1[i] = v
		}
	}

	return nv0, nv1
}
