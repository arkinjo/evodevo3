package multicell

import (
	// "gonum.org/v1/gonum/stat/distuv"
	"math"
	"math/rand"
)

// vector
type Vec = []float64

// sparse matrix
type SpMat = [](map[int]float64)

func SetVec(vec Vec, v float64) {
	for i := range vec {
		vec[i] = v
	}
}

// Create a vector with initial values of "v".
func NewVec(n int, v float64) Vec {
	vec := make([]float64, n)
	SetVec(vec, v)
	return vec
}

// Create a new sparse matrix
func NewSpMat(nrow int) SpMat {
	mat := make(SpMat, nrow)
	for i := range mat {
		mat[i] = make(map[int]float64)
	}
	return mat
}

// copy a sparse matrix
func CopySpMat(sp SpMat) SpMat {
	nsp := NewSpMat(len(sp))
	for i, m := range sp {
		for j, d := range m {
			nsp[i][j] = d
		}
	}
	return nsp
}

// multiply a sparse matrix to a vector
func MultSpMatVec(vout Vec, sp SpMat, v Vec) {
	for i, ri := range sp {
		vout[i] = 0.0
		for j, a := range ri {
			vout[i] += a * v[j]
		}
	}
}

// random matrix
func RandomizeSpMat(sp SpMat, ncol int, density float64) {
	d2 := density / 2
	for _, ri := range sp {
		for j := 0; j < ncol; j++ {
			r := rand.Float64()
			if r < d2 {
				ri[j] = 1
			} else if r < density {
				ri[j] = -1
			}
		}
	}
}

func ApplyFVec(vout Vec, f func(float64) float64, vin Vec) {
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

func MultVecSca(vout, vin Vec, f float64) {
	for i, v := range vin {
		vout[i] = f * v
	}
}

func NormalizeVec(v Vec) {
	mag := VecNorm2(v)
	MultVecSca(v, v, 1/mag)
}

func AddVecs(vout, v0, v1 Vec) {
	for i, v := range v0 {
		vout[i] = v + v1[i]
	}
}

func DiffVecs(vout, v0, v1 Vec) {
	for i, v := range v0 {
		vout[i] = v - v1[i]
	}
}

func VecNorm1(v Vec) float64 {
	d := 0.0
	for _, x := range v {
		d += math.Abs(x)
	}
	return d
}

func VecNorm2(v Vec) float64 {
	d := 0.0
	for _, x := range v {
		d += x * x
	}
	return math.Sqrt(d)
}
