package multicell

import (
	"gonum.org/v1/gonum/stat/distuv"
	"math/rand/v2"
)

/*
		Genome is an array of maps of sparse matrices.
	        g Genome
                g.E[iface] is the matrix connecting cell.E[iface] and cell.S[0]
	        g.M[IntPair{l,k}] is the matrix connecting cell.S[k] to cell.S[l].
	        Feedforward if k < l.
	        Feedbackward if k > l.
	        Self-loop if l == k.
*/

type Genome struct {
	E [NumFaces]SpMat   // input to middle layers
	M map[IntPair]SpMat // within middle layers
	W Vec               // weight of activation function
}

// Expected variance of a random genome.
func (s *Setting) RandomGenomeVariance() float64 {
	v := float64(NumFaces*s.LenLayer[0]*s.LenFace) * s.DensityEM
	s.Topology.Do(func(l, k int, density float64) {
		v += float64(s.LenLayer[l]*s.LenLayer[k]) * density
	})

	return v
}

func (s *Setting) NewGenome() Genome {
	var E [NumFaces]SpMat
	for i := range NumFaces {
		E[i] = NewSpMat(s.LenLayer[0], s.LenFace)
		E[i].Randomize(s.DensityEM)
	}

	M := make(map[IntPair]SpMat)
	s.Topology.Do(func(l, k int, density float64) {
		m := NewSpMat(s.LenLayer[l], s.LenLayer[k])
		m.Randomize(density)
		M[IntPair{l, k}] = m
	})

	W := make(Vec, s.NumLayers)
	copy(W, s.Omega)
	return Genome{E: E, M: M, W: W}
}

func (genome *Genome) Copy() Genome {
	var E [NumFaces]SpMat

	for i, mat := range genome.E {
		E[i] = mat.Copy()
	}
	M := make(map[IntPair]SpMat)
	for lk, mat := range genome.M {
		M[lk] = mat.Copy()
	}

	W := make(Vec, len(genome.W))
	copy(W, genome.W)

	return Genome{E: E, M: M, W: W}
}

func (genome *Genome) Mutate(s *Setting) {
	dist := distuv.Poisson{Lambda: 1.0}

	for i := range NumFaces {
		dist.Lambda = s.MutRate * float64(s.LenLayer[0]*s.LenFace)
		nmut := int(dist.Rand())
		for n := 0; n < nmut; n++ {
			genome.E[i].Mutate(s.DensityEM)
		}
	}

	s.Topology.Do(func(l, k int, density float64) {
		lambda := s.MutRate * float64(s.LenLayer[l]*s.LenLayer[k])
		dist := distuv.Poisson{Lambda: lambda}
		nmut := int(dist.Rand())
		for range nmut {
			genome.M[IntPair{l, k}].Mutate(density)
		}
	})

	for l := range genome.W {
		if rand.Float64() < s.MutRate {
			if rand.IntN(2) == 1 {
				genome.W[l] *= 1.1
			} else {
				genome.W[l] /= 1.1
			}
		}
	}
}

func (g0 *Genome) MateWith(g1 Genome) (Genome, Genome) {
	var E0, E1 [NumFaces]SpMat
	for i, e0 := range g0.E {
		E0[i], E1[i] = MateSpMats(e0, g1.E[i])
	}

	M0 := make(map[IntPair]SpMat)
	M1 := make(map[IntPair]SpMat)
	for lk, m0 := range g0.M {
		nmat0, nmat1 := MateSpMats(m0, g1.M[lk])
		M0[lk] = nmat0
		M1[lk] = nmat1
	}
	W0 := make(Vec, len(g0.W))
	W1 := make(Vec, len(g0.W))

	for l, w := range g0.W {
		if rand.IntN(2) == 1 {
			W0[l] = w
			W1[l] = g1.W[l]
		} else {
			W1[l] = w
			W0[l] = g1.W[l]
		}
	}
	kid0 := Genome{E: E0, M: M0, W: W0}
	kid1 := Genome{E: E1, M: M1, W: W1}
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
			if mat, ok := g.M[IntPair{l, k}]; ok {
				vec = append(vec, mat.ToVec()...)
			}
		}
	}

	vec = append(vec, g.W...)

	return vec
}

func (g0 *Genome) Equal(g1 *Genome) bool {
	for iface, e := range g0.E {
		if !e.Equal(g1.E[iface]) {
			return false
		}
	}
	for lk, m := range g0.M {
		if !m.Equal(g1.M[lk]) {
			return false
		}
	}

	for l, w := range g0.W {
		if w != g1.W[l] {
			return false
		}
	}
	return true
}
