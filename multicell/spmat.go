package multicell

import (
	// "gonum.org/v1/gonum/stat/distuv"
	"math/rand"
)

// Create a vector with initial values of "v".
func NewVec(n int, v float64) Vec {
	vec := make([]float64, n)
	for i := range n {
		vec[i] = v
	}
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
func MultSpMatVec(sp SpMat, v, vout Vec) {
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

func ApplyFVec(f func(float64) float64, omega float64, vin, vout Vec) {
	for i, v := range vin {
		vout[i] = f(v * omega)
	}
}

func AddVecs(v0, v1, vout Vec) {
	for i, v := range v0 {
		vout[i] = v + v1[i]
	}
}

func DiffVecs(v0, v1, vout Vec) {
	for i, v := range v0 {
		vout[i] = v - v1[i]
	}
}
