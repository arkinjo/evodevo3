package multicell

import (
	"gonum.org/v1/gonum/stat/distuv"
)

/*
		Genome is an array of maps of sparse matrices.
	        g Genome
                g.E[iface] is the matrix connecting cell.E[iface] and cell.S[0]
	        g.M[l][k] is the matrix connecting cell.S[k] to cell.S[l].
	        Feedforward if k < l.
	        Feedbackward if k > l.
	        Self-loop if l == k.
*/

type Genome struct {
	E [NumFaces]SpMat
	M []map[int]SpMat
}

func (s *Setting) NewGenome() Genome {
	var E [NumFaces]SpMat
	M := make([](map[int]SpMat), s.NumLayers)

	for i := range NumFaces {
		E[i] = NewSpMat(s.LenLayer[0], s.LenFace)
		E[i].Randomize(s.DensityEM)
	}

	for l := range s.NumLayers {
		M[l] = make(map[int]SpMat)
		for k := range s.NumLayers {
			if density := s.Topology.At(l, k); density > 0 {
				m := NewSpMat(s.LenLayer[l], s.LenLayer[k])
				m.Randomize(density)
				M[l][k] = m
			}
		}
	}
	return Genome{E: E, M: M}
}

func (genome *Genome) Mutate(s *Setting) {
	var nk, nmut int
	dist := distuv.Poisson{Lambda: 1.0}

	for i := range NumFaces {
		nk = s.LenLayer[0]
		dist.Lambda = s.MutRate * float64(s.LenFace*nk)
		nmut = int(dist.Rand())
		for n := 0; n < nmut; n++ {
			genome.E[i].Mutate(s.DensityEM)
		}
	}

	s.Topology.Do(func(l, k int, density float64) {
		lambda := s.MutRate * float64(s.LenLayer[l]*s.LenLayer[k])
		dist := distuv.Poisson{Lambda: lambda}
		nmut := int(dist.Rand())
		for n := 0; n < nmut; n++ {
			genome.M[l][k].Mutate(density)
		}
	})
}

func (s *Setting) MateGenomes(g0, g1 Genome) (Genome, Genome) {
	var E0, E1 [NumFaces]SpMat
	M0 := make([](map[int]SpMat), s.NumLayers)
	M1 := make([](map[int]SpMat), s.NumLayers)

	for i := range NumFaces {
		E0[i], E1[i] = MateSpMats(g0.E[i], g1.E[i])
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
	kid0 := Genome{E: E0, M: M0}
	kid1 := Genome{E: E1, M: M1}
	kid0.Mutate(s)
	kid1.Mutate(s)
	return kid0, kid1
}

func (g *Genome) ToVec(s *Setting) Vec {
	var vec Vec
	for _, e := range g.E {
		vec = append(vec, e.ToVec()...)
	}
	// Go's map is UNORDERED (random order for every "range").
	for l := range s.NumLayers {
		for k := range s.NumLayers {
			if mat, ok := g.M[l][k]; ok {
				vec = append(vec, mat.ToVec()...)
			}
		}
	}

	return vec
}

func (g0 *Genome) Equal(g1 *Genome) bool {
	for iface, e := range g0.E {
		if !e.Equal(g1.E[iface]) {
			return false
		}
	}
	for l, ml := range g0.M {
		for k, m := range ml {
			if !m.Equal(g1.M[l][k]) {
				return false
			}
		}
	}
	return true
}
