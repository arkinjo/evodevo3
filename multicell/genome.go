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

type Genome struct {
	E []SpMat
	M [](map[int]SpMat)
	P []SpMat
}

func (s *Setting) NewGenome() Genome {
	E := make([]SpMat, NumFaces)
	P := make([]SpMat, NumFaces)
	M := make([](map[int]SpMat), s.NumLayers)

	for i := range NumFaces {
		E[i] = NewSpMat(s.LenLayer[0])
		RandomizeSpMat(E[i], s.LenFace, s.DensityEM)

		P[i] = NewSpMat(s.LenFace)
		RandomizeSpMat(P[i], s.LenLayer[s.NumLayers-1], s.DensityMP)
	}

	for l, rl := range s.Topology {
		M[l] = make(map[int]SpMat)
		for k, density := range rl {
			m := NewSpMat(s.LenLayer[l])
			RandomizeSpMat(m, s.LenLayer[k], density)
			M[l][k] = m
		}
	}
	return Genome{E: E, M: M, P: P}
}

func MutateSpMat(sp SpMat, ncol int, density float64) {
	i := rand.IntN(len(sp))
	j := rand.IntN(ncol)

	if rand.Float64() >= density {
		delete(sp[i], j)
	} else if rand.IntN(2) == 1 {
		sp[i][j] = 1.0
	} else {
		sp[i][j] = -1.0
	}
}

func (genome *Genome) MutateGenome(s *Setting) {
	var nk, nmut int
	dist := distuv.Poisson{Lambda: 1.0}

	for i := range NumFaces {
		nk = s.LenLayer[0]
		dist.Lambda = s.MutRate * float64(s.LenFace*nk)
		nmut = int(dist.Rand())
		for n := 0; n < nmut; n++ {
			MutateSpMat(genome.E[i], s.LenFace, s.DensityEM)
		}

		nk = s.LenLayer[s.NumLayers-1]
		dist.Lambda = s.MutRate * float64(s.LenFace*nk)
		nmut = int(dist.Rand())
		for n := 0; n < nmut; n++ {
			MutateSpMat(genome.P[i], nk, s.DensityMP)
		}
	}

	for l, rl := range s.Topology {
		nl := s.LenLayer[l]
		for k, density := range rl {
			nk := s.LenLayer[k]
			lambda := s.MutRate * float64(nl*nk)
			dist := distuv.Poisson{Lambda: lambda}
			nmut := int(dist.Rand())
			for n := 0; n < nmut; n++ {
				MutateSpMat(genome.M[l][k], nk, density)
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

func (s *Setting) MateGenomes(g0, g1 Genome) (Genome, Genome) {
	E0 := make([]SpMat, NumFaces)
	E1 := make([]SpMat, NumFaces)
	P0 := make([]SpMat, NumFaces)
	P1 := make([]SpMat, NumFaces)
	M0 := make([](map[int]SpMat), s.NumLayers)
	M1 := make([](map[int]SpMat), s.NumLayers)

	for i := range NumFaces {
		E0[i], E1[i] = MateSpMats(g0.E[i], g1.E[i])
		P0[i], P1[i] = MateSpMats(g0.P[i], g1.P[i])
	}
	for l, rl := range g0.M {
		M0[l] = make(map[int]SpMat)
		M1[l] = make(map[int]SpMat)
		for k, sp := range rl {
			nmat0, nmat1 := MateSpMats(sp, g1.M[l][k])
			M0[l][k] = nmat0
			M1[l][k] = nmat1
		}
	}
	kid0 := Genome{E: E0, M: M0, P: P0}
	kid1 := Genome{E: E1, M: M1, P: P1}
	kid0.MutateGenome(s)
	kid1.MutateGenome(s)
	return kid0, kid1
}

func (g *Genome) ToVec(s *Setting) Vec {
	var vec Vec
	for _, e := range g.E {
		for _, ei := range e { // i in [0:s.LenLayer[0]]
			for j := range s.LenFace {
				v, ok := ei[j]
				if !ok {
					v = 0.0
				}
				vec = append(vec, v)
			}
		}
	}
	for _, ml := range g.M {
		for k, m := range ml {
			for _, mi := range m {
				for j := range s.LenLayer[k] {
					v, ok := mi[j]
					if !ok {
						v = 0.0
					}
					vec = append(vec, v)
				}
			}
		}
	}
	for _, p := range g.P {
		for _, pi := range p {
			for j := range s.LenLayer[s.NumLayers-1] {
				v, ok := pi[j]
				if !ok {
					v = 0.0
				}
				vec = append(vec, v)
			}
		}
	}

	return vec
}
