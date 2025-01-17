package multicell

import (
	"gonum.org/v1/gonum/stat/distuv"
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
	return slices.Clone(vec)
}

// sum all elements.
func (vec Vec) Sum() float64 {
	s := 0.0
	for _, v := range vec {
		s += v
	}
	return s
}

func (vec Vec) Mean() float64 {
	return vec.Sum() / float64(len(vec))
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

func (vec Vec) ScaleBy(f float64) Vec {
	for i, v := range vec {
		vec[i] = f * v
	}
	return vec
}

func (v Vec) Normalize() Vec {
	mag := v.Norm2()
	v.ScaleBy(1 / mag)
	return v
}

func (vout Vec) Add(v0, v1 Vec) Vec {
	for i, v := range v0 {
		vout[i] = v + v1[i]
	}
	return vout
}

// Accumulate
func (vout Vec) Acc(vin Vec) Vec {
	for i, v := range vin {
		vout[i] += v
	}
	return vout
}

func (vout Vec) Diff(v0, v1 Vec) Vec {
	if len(vout) == 0 {
		vout = make(Vec, len(v0))
	}
	for i, v := range v0 {
		vout[i] = v - v1[i]
	}
	return vout
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

func (v Vec) NormInf() float64 {
	d := 0.0
	for _, x := range v {
		d = max(d, math.Abs(x))
	}
	return d
}

func NoisyVec(n int, noise float64) Vec {
	m := int(noise * float64(n))
	v := make(Vec, n)
	for _, p := range rand.Perm(n) {
		if p < m {
			v[p] = -1
		} else {
			v[p] = 1
		}
	}
	return v
}

func DiffMats(vs1, vs0 []Vec) []Vec {
	dvs := make([]Vec, len(vs0))
	for i, v1 := range vs1 {
		dvs[i] = make(Vec, len(v1))
		dvs[i].Diff(v1, vs0[i])
	}
	return dvs
}

func CorrVecs(vs0, vs1 Vec) (float64, float64) {
	var m0, m1 float64
	for i, v := range vs0 {
		m0 += v
		m1 += vs1[i]
	}
	f := float64(len(vs0))
	m0 /= f
	m1 /= f
	var v0, v1, corr float64
	for i, v := range vs0 {
		d0 := v - m0
		d1 := vs1[i] - m1
		v0 += d0 * d0
		v1 += d1 * d1
		corr += d0 * d1
	}
	r := corr / math.Sqrt(v0*v1)
	tstat := r * math.Sqrt((f-2)/(1-r*r))
	dist := distuv.StudentsT{0, 1, f - 2, nil}
	pval := 2 * dist.CDF(-math.Abs(tstat))
	return r, pval
}
