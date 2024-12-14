package multicell

import (
	// "gonum.org/v1/gonum/stat/distuv"
	"math"
	"math/rand/v2"
)

// vector
type Vec = []float64

// sparse matrix
type SpMat struct {
	Nrow int
	Ncol int
	M    []map[int]float64
}

func (sp0 *SpMat) Equal(sp1 SpMat) bool {
	if sp0.Nrow != sp1.Nrow || sp0.Ncol != sp1.Ncol {
		return false
	}
	for i, vi := range sp0.M {
		for j, v := range vi {
			if v != sp1.M[i][j] {
				return false
			}
		}
	}
	for i, vi := range sp1.M {
		for j, v := range vi {
			if v != sp1.M[i][j] {
				return false
			}
		}
	}
	return true
}

func VecSet(vec Vec, v float64) {
	for i := range vec {
		vec[i] = v
	}
}

// Create a vector with initial values of "v".
func NewVec(n int, v float64) Vec {
	vec := make([]float64, n)
	VecSet(vec, v)
	return vec
}

// Create a new sparse matrix
func NewSpMat(nrow, ncol int) SpMat {
	mat := make([]map[int]float64, nrow)
	for i := range nrow {
		mat[i] = make(map[int]float64)
	}
	return SpMat{
		Nrow: nrow,
		Ncol: ncol,
		M:    mat}

}

func (sp *SpMat) At(i, j int) float64 {
	return sp.M[i][j]
}

func (sp *SpMat) SetAt(i, j int, v float64) {
	sp.M[i][j] = v
}

// copy a sparse matrix
func (sp *SpMat) Copy() SpMat {
	nsp := NewSpMat(sp.Nrow, sp.Ncol)
	for i, vi := range sp.M {
		for j, v := range vi {
			nsp.M[i][j] = v
		}
	}
	return nsp
}

// multiply a sparse matrix to a vector
func (sp *SpMat) MultVec(vin, vout Vec) {
	for i, xi := range sp.M {
		vout[i] = 0.0
		for j, x := range xi {
			vout[i] += x * vin[j]
		}
	}
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
func (sp *SpMat) Randomize(density float64) {
	d2 := density / 2
	for i := range sp.Nrow {
		sp.M[i] = make(map[int]float64)
		for j := range sp.Ncol {
			r := rand.Float64()
			if r < d2 {
				sp.M[i][j] = 1
			} else if r < density {
				sp.M[i][j] = -1
			}
		}
	}
}

func ApplyFVec(f func(float64) float64, vin, vout Vec) {
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

func VecScale(vin Vec, f float64) {
	for i, v := range vin {
		vin[i] = f * v
	}
}

func NormalizeVec(v Vec) {
	mag := VecNorm2(v)
	VecScale(v, 1/mag)
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

func MateSpMats(mat0, mat1 SpMat) (SpMat, SpMat) {
	nmat0 := NewSpMat(mat0.Nrow, mat0.Ncol)
	nmat1 := NewSpMat(mat0.Nrow, mat0.Ncol)
	for i, ri := range mat0.M {
		if rand.IntN(2) == 1 {
			for j, v := range ri {
				nmat0.M[i][j] = v
			}
			for j, v := range mat1.M[i] {
				nmat1.M[i][j] = v
			}
		} else {
			for j, v := range ri {
				nmat1.M[i][j] = v
			}
			for j, v := range mat1.M[i] {
				nmat0.M[i][j] = v
			}
		}
	}
	return nmat0, nmat1
}

func (sp SpMat) Do(f func(i, j int, v float64)) {
	for i, vi := range sp.M {
		for j, v := range vi {
			f(i, j, v)
		}
	}
}
