package multicell

import (
	"gonum.org/v1/gonum/stat/distuv"
	"log"
	"maps"
	"math/rand/v2"
)

// sparse matrix
type SpMat struct {
	Nrow int
	Ncol int
	M    SliceOfMaps[float64]
}

func (sp SpMat) Do(f func(i, j int, v float64)) {
	sp.M.Do(f)
}

func (sp0 SpMat) Equal(sp1 SpMat) bool {
	if sp0.Nrow != sp1.Nrow || sp0.Ncol != sp1.Ncol {
		return false
	}
	for i, mi := range sp0.M {
		if !maps.Equal(mi, sp1.M[i]) {
			return false
		}
	}
	return true
}

// Create a new sparse matrix
func NewSpMat(nrow, ncol int) SpMat {
	mat := NewSliceOfMaps[float64](nrow)
	return SpMat{
		Nrow: nrow,
		Ncol: ncol,
		M:    mat}

}

func (sp *SpMat) At(i, j int) float64 {
	return sp.M[i][j]
}

func (sp *SpMat) Set(i, j int, v float64) {
	sp.M[i][j] = v
}

// copy a sparse matrix
func (sp *SpMat) Clone() SpMat {
	nsp := NewSpMat(sp.Nrow, sp.Ncol)
	sp.Do(func(i, j int, v float64) {
		nsp.M[i][j] = v
	})
	return nsp
}

// multiply a sparse matrix to a vector. vout is NOT initialized!!
func (vout Vec) MultSpMatVec(sp SpMat, vin Vec) {
	sp.Do(func(i, j int, x float64) {
		vout[i] += x * vin[j]
	})
}

func (sp *SpMat) ToVec() Vec {
	var vec Vec

	for i := range sp.Nrow {
		for j := range sp.Ncol {
			vec = append(vec, sp.M[i][j])
		}
	}

	return vec
}

// random matrix
func (sp SpMat) Randomize(density float64) {
	dist := distuv.Poisson{Lambda: density * float64(sp.Nrow*sp.Ncol)}
	n := int(dist.Rand())
	for range n {
		i := rand.IntN(sp.Nrow)
		j := rand.IntN(sp.Ncol)
		if rand.IntN(2) == 1 {
			sp.M[i][j] = 1
		} else {
			sp.M[i][j] = -1
		}
	}
}

func (sp SpMat) Mutate(density float64) {
	i := rand.IntN(sp.Nrow)
	j := rand.IntN(sp.Ncol)
	if rand.Float64() >= density {
		delete(sp.M[i], j)
	} else if rand.IntN(2) == 1 {
		sp.M[i][j] = 1.0
	} else {
		sp.M[i][j] = -1.0
	}
}

func (mat0 SpMat) MateWith(mat1 SpMat) (SpMat, SpMat) {
	if mat0.Nrow != mat1.Nrow || mat0.Ncol != mat1.Ncol {
		log.Fatal("MateSpMats: incompatible matrices")
	}

	nmat0 := NewSpMat(mat0.Nrow, mat0.Ncol)
	nmat1 := NewSpMat(mat0.Nrow, mat0.Ncol)

	for i := range mat0.Nrow {
		if rand.IntN(2) == 1 {
			nmat0.M[i] = maps.Clone(mat0.M[i])
			nmat1.M[i] = maps.Clone(mat1.M[i])
		} else {
			nmat1.M[i] = maps.Clone(mat0.M[i])
			nmat0.M[i] = maps.Clone(mat1.M[i])
		}
	}
	return nmat0, nmat1
}
