package multicell

import (
	//	"log"
	"slices"
)

/*
		Genome is an array of maps of sparse matrices.
	        g Genome
                g.E[iface] is the matrix connecting cell.E[iface] and cell.S[0]
	        g.M[l][k] is the matrix connecting from cell.S[k] to cell.S[l].
	        Feedforward if k < l.
	        Feedback if k > l.
	        Self-loop if l == k.
*/

type Genome struct {
	M SliceOfMaps[SpMat]
	W Vec // weight of activation function
}

// Expected variance of a random genome.
func (s *Setting) RandomGenomeVariance() float64 {
	v := 0.0
	s.Topology.Do(func(l, k int, density float64) {
		v += float64(s.LenLayer[l]*s.LenLayer[k]) * density
	})

	return v
}

func (s *Setting) NewGenome() Genome {
	M := NewSliceOfMaps[SpMat](s.NumLayers)
	s.Topology.Do(func(l, k int, density float64) {
		M[l][k] = NewSpMat(s.LenLayer[l], s.LenLayer[k])
		M[l][k].Randomize(density)
	})

	W := NewVec(s.NumLayers, 1.0)
	return Genome{M: M, W: W}
}

func (genome *Genome) Clone() Genome {
	M := NewSliceOfMaps[SpMat](len(genome.M))
	genome.M.Do(func(l, k int, mat SpMat) {
		M[l][k] = mat.Clone()
	})

	W := slices.Clone(genome.W)

	return Genome{M: M, W: W}
}

func (genome *Genome) Mutate(s *Setting) {
	s.Topology.Do(func(l, k int, density float64) {
		genome.M[l][k].Mutate(s.MutRate, density)
	})

	genome.W.Mutate(s.MutRate, s.WScale)
}

func (g0 *Genome) MateWith(g1 Genome) (Genome, Genome) {
	M0 := NewSliceOfMaps[SpMat](len(g0.M))
	M1 := NewSliceOfMaps[SpMat](len(g1.M))
	g0.M.Do(func(l, k int, m0 SpMat) {
		M0[l][k], M1[l][k] = m0.MateWith(g1.M[l][k])
	})

	W0, W1 := g0.W.MateWith(g1.W)

	kid0 := Genome{M: M0, W: W0}
	kid1 := Genome{M: M1, W: W1}
	return kid0, kid1
}

func (g *Genome) ToVec(s *Setting) Vec {
	var vec Vec
	// Go's map is UNORDERED (random order for every "range").
	for l := range s.NumLayers {
		for k := range s.NumLayers {
			if mat, ok := g.M[l][k]; ok {
				vec = append(vec, mat.ToVec()...)
			}
		}
	}

	vec = append(vec, g.W...)

	return vec
}

func (g0 Genome) Equal(g1 Genome) bool {
	for l, ml := range g0.M {
		for k, m := range ml {
			if !m.Equal(g1.M[l][k]) {
				return false
			}
		}
	}

	return slices.Equal(g0.W, g1.W)
}
