package multicell

import (
	"fmt"
	"log"
	"maps"
	"math/rand/v2"

	"gonum.org/v1/gonum/stat/distuv"
)

// BLOCK-DIAGONAL sparse matrix
type SpMat struct {
	Ncol int // number of columns
	Nblk int // number of diagonal blocks
	SliceOfMaps[float64]
}

func (sp SpMat) Nrows() int {
	return len(sp.M)
}

func (sp SpMat) Ncols() int {
	return sp.Ncol
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
func NewSpMat(nrow, ncol, nblock int) SpMat {
	mat := NewSliceOfMaps[float64](nrow)
	return SpMat{ncol, nblock, mat}
}

// copy a sparse matrix
func (sp *SpMat) Clone() SpMat {
	nsp := NewSpMat(sp.Nrows(), sp.Ncols(), sp.Nblk)
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

	for _, mi := range sp.M {
		for j := range sp.Ncol {
			vec = append(vec, mi[j])
		}
	}

	return vec
}

func (sp SpMat) Density() float64 {
	nonz := 0.0
	sp.Do(func(_, _ int, _ float64) {
		nonz += 1.0
	})
	return nonz / float64(sp.Nrows()*sp.Ncols())
}

func (sp SpMat) PickRandomElements(n int) SliceOfMaps[float64] {
	nr := sp.Nrows() / sp.Nblk // elements per block in one row.
	nc := sp.Ncols() / sp.Nblk // elements per block in one column
	npb := min(n/sp.Nblk, nr*nc)
	ps := NewSliceOfMaps[float64](sp.Nrows())
	for ib := range sp.Nblk {
		dist := distuv.Poisson{Lambda: float64(npb)}
		nb := min(int(dist.Rand()), nr*nc)

		for _, p := range rand.Perm(nr * nc)[:nb] {
			i := p/nc + ib*nr
			j := p%nc + ib*nc
			//	log.Println("i,j", ib, p, i, j)
			ps.M[i][j] = rand.Float64()
		}
	}

	return ps
}

// random matrix
func (sp SpMat) Randomize(density float64) {
	if density*float64(sp.Nblk) > 1.0 {
		err := fmt.Errorf("SpMat.Randomize: density %f is too big for %d blocks", density, sp.Nblk)
		log.Fatal(err)
	}
	nr := sp.Nrows()
	nc := sp.Ncols()
	dist := distuv.Poisson{Lambda: density * float64(nr*nc)}
	n := int(dist.Rand())
	sp.PickRandomElements(n).Do(func(i, j int, r float64) {
		if r < 0.5 {
			sp.M[i][j] = 1
		} else {
			sp.M[i][j] = -1
		}
	})
}

func (sp SpMat) Mutate(rate float64, density float64) {
	d2 := density / 2.0
	nr := sp.Nrows()
	nc := sp.Ncols()
	dist := distuv.Poisson{Lambda: rate * float64(nr*nc)}
	n := int(dist.Rand())
	sp.PickRandomElements(n).Do(func(i, j int, r float64) {
		if r < density {
			if v, ok := sp.M[i][j]; ok {
				sp.M[i][j] = -v
			} else if r < d2 {
				sp.M[i][j] = 1.0
			} else {
				sp.M[i][j] = -1.0
			}
		} else {
			delete(sp.M[i], j)
		}
	})
}

func (mat0 SpMat) MateWith(mat1 SpMat) (SpMat, SpMat) {
	if mat0.Nrows() != mat1.Nrows() || mat0.Ncols() != mat1.Ncols() {
		log.Fatal("MateSpMats: incompatible matrices")
	}

	nmat0 := NewSpMat(mat0.Nrows(), mat0.Ncols(), mat0.Nblk)
	nmat1 := NewSpMat(mat0.Nrows(), mat0.Ncols(), mat0.Nblk)

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
