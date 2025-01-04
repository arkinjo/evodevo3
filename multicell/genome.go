package multicell

/*
		Genome is an array of maps of sparse matrices.
	        g Genome
	        g.G[l][k] is the matrix connecting from cell.S[k] to cell.S[l].
	        Feedforward if k < l.
	        Feedback if k >= l.
*/

// This better be
// type Genome SliceOfMaps[SpMat]
// but it doesn't work...methods of SliceOfMaps[T] are lost as Genome is a completely different type.
// type Genome = SliceOfMaps[SpMat]
// doesn't work either... new methods can't be defined for Genome.

type Genome struct {
	M SliceOfMaps[SpMat]
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
	G := NewSliceOfMaps[SpMat](s.NumLayers)
	s.Topology.Do(func(l, k int, density float64) {
		G[l][k] = NewSpMat(s.LenLayer[l], s.LenLayer[k])
		G[l][k].Randomize(density)
	})

	return Genome{G}
}

func (genome Genome) Clone() Genome {
	G := NewSliceOfMaps[SpMat](len(genome.M))
	genome.M.Do(func(l, k int, mat SpMat) {
		G[l][k] = mat.Clone()
	})

	return Genome{G}
}

func (genome Genome) Mutate(s *Setting) {
	s.Topology.Do(func(l, k int, density float64) {
		genome.M[l][k].Mutate(s.MutRate, density)
	})
}

func (g0 Genome) MateWith(g1 Genome) (Genome, Genome) {
	M0 := NewSliceOfMaps[SpMat](len(g0.M))
	M1 := NewSliceOfMaps[SpMat](len(g1.M))
	g0.M.Do(func(l, k int, m0 SpMat) {
		M0[l][k], M1[l][k] = m0.MateWith(g1.M[l][k])
	})

	return Genome{M0}, Genome{M1}
}

func (g Genome) ToVec(s *Setting) Vec {
	var vec Vec
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

func (g0 Genome) Equal(g1 Genome) bool {
	for l, ml := range g0.M {
		for k, m := range ml {
			if !m.Equal(g1.M[l][k]) {
				return false
			}
		}
	}

	return true
}
