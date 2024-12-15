package multicell

import (
	"gonum.org/v1/gonum/stat/distuv"
	"log"
	"math/rand/v2"
)

// index for sparse matrices etc.
type IntPair struct {
	I, J int
}

// sparse matrix
type SpMat struct {
	Nrow int
	Ncol int
	M    map[IntPair]float64
}

func (sp SpMat) Do(f func(i, j int, v float64)) {
	for ij, v := range sp.M {
		f(ij.I, ij.J, v)
	}
}

func (sp0 *SpMat) Equal(sp1 SpMat) bool {
	if sp0.Nrow != sp1.Nrow || sp0.Ncol != sp1.Ncol {
		return false
	}
	for ij, v := range sp0.M {
		if v != sp1.M[ij] {
			return false
		}
	}
	for ij, v := range sp1.M {
		if v != sp1.M[ij] {
			return false
		}
	}
	return true
}

// Create a new sparse matrix
func NewSpMat(nrow, ncol int) SpMat {
	mat := make(map[IntPair]float64)
	return SpMat{
		Nrow: nrow,
		Ncol: ncol,
		M:    mat}

}

func (sp *SpMat) At(i, j int) float64 {
	return sp.M[IntPair{i, j}]
}

func (sp *SpMat) Set(i, j int, v float64) {
	sp.M[IntPair{i, j}] = v
}

// copy a sparse matrix
func (sp *SpMat) Copy() SpMat {
	nsp := NewSpMat(sp.Nrow, sp.Ncol)
	for ij, v := range sp.M {
		nsp.M[ij] = v
	}
	return nsp
}

// multiply a sparse matrix to a vector
func (sp *SpMat) MultVec(vin, vout Vec) {
	vout.SetAll(0.0)
	sp.Do(func(i, j int, x float64) {
		vout[i] += x * vin[j]
	})
}

func (sp *SpMat) ToVec() Vec {
	var vec Vec

	for i := range sp.Nrow {
		for j := range sp.Ncol {
			vec = append(vec, sp.M[IntPair{i, j}])
		}
	}

	return vec
}

// random matrix
func (sp *SpMat) Randomize(density float64) {
	dist := distuv.Poisson{Lambda: density * float64(sp.Nrow*sp.Ncol)}
	n := int(dist.Rand())
	for range n {
		i := rand.IntN(sp.Nrow)
		j := rand.IntN(sp.Ncol)
		if rand.IntN(2) == 1 {
			sp.M[IntPair{i, j}] = 1
		} else {
			sp.M[IntPair{i, j}] = -1
		}
	}
}

func (sp SpMat) Mutate(density float64) {
	ij := IntPair{I: rand.IntN(sp.Nrow), J: rand.IntN(sp.Ncol)}
	if rand.Float64() >= density {
		delete(sp.M, ij)
	} else if rand.IntN(2) == 1 {
		sp.M[ij] = 1.0
	} else {
		sp.M[ij] = -1.0
	}
}

func MateSpMats(mat0, mat1 SpMat) (SpMat, SpMat) {
	if mat0.Nrow != mat1.Nrow || mat0.Ncol != mat1.Ncol {
		log.Fatal("MateSpMats: incompatible matrices")
	}

	nmat0 := NewSpMat(mat0.Nrow, mat0.Ncol)
	nmat1 := NewSpMat(mat0.Nrow, mat0.Ncol)
	ind0 := make([][]IntPair, mat0.Nrow)
	ind1 := make([][]IntPair, mat1.Nrow)

	for ij := range mat0.M {
		ind0[ij.I] = append(ind0[ij.I], ij)
	}
	for ij := range mat1.M {
		ind1[ij.I] = append(ind1[ij.I], ij)
	}
	for i := range mat0.Nrow {
		if rand.IntN(2) == 1 {
			for _, ij := range ind0[i] {
				nmat0.M[ij] = mat0.M[ij]
			}
			for _, ij := range ind1[i] {
				nmat1.M[ij] = mat1.M[ij]
			}
		} else {
			for _, ij := range ind0[i] {
				nmat1.M[ij] = mat0.M[ij]
			}
			for _, ij := range ind1[i] {
				nmat0.M[ij] = mat1.M[ij]
			}
		}
	}
	return nmat0, nmat1
}
