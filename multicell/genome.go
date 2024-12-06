package multicell

import (
	"gonum.org/v1/gonum/stat/distuv"
	"math/rand/v2"
)

/*
		Genome is an array of maps of sparse matrices.
	        g Genome
	        g[l][k] is the matrix connecting the k-th layer to the l-th layer.
	        Feedforward if k < l.
	        Feedbackward if k > l.
	        Self-loop if l == k.
*/

type Genome [](map[int]SpMat)

func (s *Setting) NewGenome() Genome {
	genome := make([](map[int]SpMat), s.Num_layers)
	for l, rl := range s.Topology {
		genome[l] = make(map[int]SpMat)
		for k, density := range rl {
			m := NewSpMat(s.Num_components[l])
			RandomizeSpMat(m, s.Num_components[k], density)
			genome[l][k] = m
		}
	}
	return genome
}

func MutateSpMat(sp SpMat, nrow, ncol int, density float64) {
	i := rand.IntN(nrow)
	j := rand.IntN(ncol)

	if rand.Float64() > density {
		sp[i][j] = 0.0
	} else if rand.IntN(2) == 1 {
		sp[i][j] = 1.0
	} else {
		sp[i][j] = -1.0
	}
}

func (s *Setting) MutateGenome(genome Genome) {
	for l, rl := range s.Topology {
		nl := s.Num_components[l]
		for k, density := range rl {
			nk := s.Num_components[k]
			lambda := s.Mut_rate * float64(nl*nk)
			dist := distuv.Poisson{Lambda: lambda}
			nmut := int(dist.Rand())
			for n := 0; n < nmut; n++ {
				MutateSpMat(genome[l][k], nl, nk, density)
			}
		}
	}
}

func MateSpMats(mat0, mat1 SpMat) (SpMat, SpMat) {
	nrow := len(mat0)
	nmat0 := NewSpMat(nrow)
	nmat1 := NewSpMat(nrow)
	for i, ri := range mat0 {
		if rand.IntN(2) == 1 {
			for j, v := range ri {
				nmat0[i][j] = v
			}
			for j, v := range mat1[i] {
				nmat1[i][j] = v
			}
		} else {
			for j, v := range ri {
				nmat1[i][j] = v
			}
			for j, v := range mat1[i] {
				nmat0[i][j] = v
			}
		}
	}

	return nmat0, nmat1
}

func (s *Setting) MateGenomes(gen0, gen1 Genome) (Genome, Genome) {
	kid0 := make([](map[int]SpMat), len(gen0))
	kid1 := make([](map[int]SpMat), len(gen1))

	for l, rl := range gen0 {
		kid0[l] = make(map[int]SpMat)
		kid1[l] = make(map[int]SpMat)
		for k, sp := range rl {
			nmat0, nmat1 := MateSpMats(sp, gen1[l][k])
			kid0[l][k] = nmat0
			kid1[l][k] = nmat1
		}
	}

	return kid0, kid1
}
