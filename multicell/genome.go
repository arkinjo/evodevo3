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
	B []Vec
	SliceOfMaps[SpMat]
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
	B := make([]Vec, s.NumLayers)
	for l, nl := range s.LenLayer {
		B[l] = NewVec(nl, 0.0)
	}
	G := NewSliceOfMaps[SpMat](s.NumLayers)
	s.Topology.Do(func(l, k int, density float64) {
		G.M[l][k] = NewSpMat(s.LenLayer[l], s.LenLayer[k], s.NumBlocks)
		G.M[l][k].Randomize(density)
	})

	return Genome{B, G}
}

func (genome Genome) Clone() Genome {
	B := make([]Vec, len(genome.B))
	if with_bias {
		for l, bl := range genome.B {
			B[l] = bl.Clone()
		}
	}
	G := NewSliceOfMaps[SpMat](len(genome.M))
	genome.Do(func(l, k int, mat SpMat) {
		G.M[l][k] = mat.Clone()
	})

	return Genome{B, G}
}

func (genome Genome) Mutate(s *Setting) {
	if with_bias {
		for l := range genome.B {
			genome.B[l].Mutate(s.MutRate)
		}
	}
	s.Topology.Do(func(l, k int, density float64) {
		genome.M[l][k].Mutate(s.MutRate, density)
	})
}

func (g0 Genome) MateWith(g1 Genome) (Genome, Genome) {
	B0 := make([]Vec, len(g0.B))
	B1 := make([]Vec, len(g1.B))
	if with_bias {
		for l, b0 := range g0.B {
			B0[l], B1[l] = b0.MateWith(g1.B[l])
		}
	}
	M0 := NewSliceOfMaps[SpMat](len(g0.M))
	M1 := NewSliceOfMaps[SpMat](len(g1.M))
	g0.Do(func(l, k int, m0 SpMat) {
		M0.M[l][k], M1.M[l][k] = m0.MateWith(g1.M[l][k])
	})

	return Genome{B0, M0}, Genome{B1, M1}
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
		if with_bias {
			vec = append(vec, g.B[l]...)
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
