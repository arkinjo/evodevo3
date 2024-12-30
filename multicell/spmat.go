package multicell

import (
	"gonum.org/v1/gonum/stat/distuv"
	"log"
	"maps"
	"math/rand/v2"
)

// sparse matrix
type SpMat struct {
	Ncol int
	M    SliceOfMaps[float64]
}

func (sp SpMat) Nrows() int {
	return len(sp.M)
}

func (sp SpMat) Ncols() int {
	return sp.Ncol
}

func (sp SpMat) Do(f func(i, j int, v float64)) {
	sp.M.Do(f)
}

func (sp0 SpMat) Equal(sp1 SpMat) bool {
	if sp0.Nrows() != sp1.Nrows() || sp0.Ncols() != sp1.Ncols() {
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
	nsp := NewSpMat(sp.Nrows(), sp.Ncols())
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

	for i := range sp.Nrows() {
		for j := range sp.Ncol {
			vec = append(vec, sp.M[i][j])
		}
	}

	return vec
}

func (sp SpMat) Density() float64 {
	nonz := 0.0
	sp.M.Do(func(i, j int, _ float64) {
		nonz += 1.0
	})
	return nonz / float64(sp.Nrows()*sp.Ncols())
}

// random matrix
func (sp SpMat) Randomize(density float64) {
	nr := sp.Nrows()
	nc := sp.Ncols()
	dist := distuv.Poisson{Lambda: density * float64(nr*nc)}
	n := int(dist.Rand())
	for range n {
		i := rand.IntN(nr)
		j := rand.IntN(nc)
		if rand.IntN(2) == 1 {
			sp.M[i][j] = 1
		} else {
			sp.M[i][j] = -1
		}
	}
}

func (sp SpMat) Mutate(rate float64, density float64) {
	nr := sp.Nrows()
	nc := sp.Ncols()
	dist := distuv.Poisson{Lambda: rate * float64(nr*nc)}
	n := int(dist.Rand())
	d2 := density * 0.5
	for range n {
		i := rand.IntN(nr)
		j := rand.IntN(nc)
		delete(sp.M[i], j)
		r := rand.Float64()
		if r < d2 {
			sp.M[i][j] = 1.0
		} else if r < density {
			sp.M[i][j] = -1.0
		}
	}
}

func (mat0 SpMat) MateWith(mat1 SpMat) (SpMat, SpMat) {
	if mat0.Nrows() != mat1.Nrows() || mat0.Ncols() != mat1.Ncols() {
		log.Fatal("MateSpMats: incompatible matrices")
	}

	nmat0 := NewSpMat(mat0.Nrows(), mat0.Ncols())
	nmat1 := NewSpMat(mat0.Nrows(), mat0.Ncols())

	for i := range mat0.Nrows() {
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
