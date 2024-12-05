package evodevo3

import (
	"gonum.org/v1/gonum/stat/distuv"
	"math/rand"
)

type Spmat struct {
	Ncol int // number of columns
	Mat  [](map[int]float64)
}

func NewSpmat(nrow, ncol int) Spmat {
	mat := make([](map[int]float64), nrow)
	for i := range mat {
		mate[i] = make(map[int]float64)
	}
	return Spmat{ncol, mat}
}

func (sp *Spmat) Copy() Spmat {
	nsp := NewSpmat(len(sp.Mat), sp.Ncol)
	for i, m := range sp.Mat {
		for j, d := range m {
			nsp.Mat[i][j] = d
		}
	}
	return nsp
}

// Randomize entries of sparse matrix
func (sp *Spmat) Randomize(density float64) {
	if density == 0 {
		return
	}

	density2 := density / 2
	for i := range sp.Mat {
		for j := 0; j < sp.Ncol; j++ {
			r := rand.Float64()
			if r < density2 {
				sp.Mat[i][j] = 1
			} else if r < density {
				sp.Mat[i][j] = -1
			}
		}
	}
}

// Apply an sparse matrix to a vector.
func MultMatVec(vout Vec, mat Spmat, vin Vec) { //Matrix multiplication
	for i := range vout {
		vout[i] = 0.0
	}

	for i, m := range mat.Mat {
		for j, d := range m {
			vout[i] += d * vin[j]
		}
	}
	return
}

// mutating a sparse matrix
// Note: This implementation has non-zero probability of choosing same element to be mutated twice.
func (mat *Spmat) mutateSpmat(density, mutrate float64) {
	if density == 0.0 {
		return
	}

	nrow := len(mat.Mat)
	lambda := mutrate * float64(nrow*mat.Ncol)
	dist := distuv.Poisson{Lambda: lambda}
	nmut := int(dist.Rand())
	density2 := density * 0.5
	for n := 0; n < nmut; n++ {
		i := rand.Intn(nrow)
		j := rand.Intn(mat.Ncol)
		r := rand.Float64()
		delete(mat.Mat[i], j)
		if r < density2 {
			mat.Mat[i][j] = 1.0
		} else if r < density {
			mat.Mat[i][j] = -1.0
		}
	}

	return
}

func CrossoverSpmats(mat0, mat1 Spmat) {
	for i, ri := range mat0.Mat {
		r := rand.Float64()
		if r < 0.5 {
			mat0.Mat[i] = mat1.Mat[i]
			mat1.Mat[i] = ri
		}
	}

}
