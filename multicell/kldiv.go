package multicell

import (
	"math"
)

type GenomeDist struct {
	Nelems int
	P      [3]SliceOfMaps[SpMat] // 0, +1, -1
}

func (s *Setting) NewGenomeDist() GenomeDist {
	var gd GenomeDist
	s.Topology.Do(func(l, k int, _ float64) {
		gd.Nelems += s.LenLayer[l] * s.LenLayer[k]
	})
	for i := range gd.P {
		gd.P[i] = NewSliceOfMaps[SpMat](s.NumLayers)
		s.Topology.Do(func(l, k int, _ float64) {
			gd.P[i].M[l][k] = NewSpMat(s.LenLayer[l], s.LenLayer[k])

		})
	}
	return gd
}

func (pop *Population) GetGenomeDist(s *Setting) GenomeDist {
	gd := s.NewGenomeDist()
	for _, indiv := range pop.Indivs {
		indiv.Genome.Do(func(l, k int, M SpMat) {
			M.Do(func(i, j int, v float64) {
				p := int(3+v) % 3
				gd.P[p].M[l][k].M[i][j] += 1.0

			})
		})
	}
	n := 1.0 / float64(len(pop.Indivs))
	for i, P := range gd.P {
		P.Do(func(l, k int, sp SpMat) {
			gd.P[i].M[l][k] = sp.ScaleBy(n)
		})
	}
	return gd
}

func (gd1 GenomeDist) KLDivergenceFrom(gd0 GenomeDist) float64 {
	kldiv := 0.0
	for ind, P := range gd1.P {
		P.Do(func(l, k int, sp SpMat) {
			Q0 := gd0.P[ind].M[l][k]
			sp.Do(func(i, j int, p float64) {
				q := Q0.M[i][j]
				if p > 0 && q > 0 {
					kldiv += p * math.Log(p/q)
				}
			})
		})
	}
	return kldiv / float64(gd1.Nelems)
}
