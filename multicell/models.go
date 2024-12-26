package multicell

import (
	"log"
)

type model_t struct {
	cue     bool
	develop bool
	nLayers int
}

var models = map[string]model_t{
	"Full":    {true, true, 3},
	"Hie0":    {true, true, 0},
	"Hie1":    {true, true, 1},
	"Hie2":    {true, true, 2},
	"Null":    {false, false, 0},
	"NoCue":   {false, true, 3},
	"NoDev":   {true, false, 3},
	"NullHie": {false, false, 3},
	"NullCue": {true, false, 0},
	"NullDev": {true, true, 0},
}

func (s *Setting) SetLayer(n int) {
	s.NumLayers = n + 1
	s.LenLayer = make([]int, s.NumLayers)
	slen := s.LenFace * NumFaces
	s.LenLayer[s.NumLayers-1] = slen

	s.Topology = NewSliceOfMaps[float64](s.NumLayers)

	switch n {
	case 3:
		s.LenLayer[0] = slen
		s.LenLayer[1] = slen
		s.LenLayer[2] = slen
		s.DensityEM = default_density
		for l := range s.NumLayers {
			// feedforward
			if l > 0 {
				s.Topology[l][l-1] = default_density
			}
			// feedback (no feedback for the phenotype layer)
			if l < s.NumLayers-1 {
				s.Topology[l][l] = default_density
			}
		}
	case 2:
		s.LenLayer[0] = slen * 3 / 2
		s.LenLayer[1] = slen * 3 / 2
		//feedforward
		s.DensityEM = default_density * 2.0 / 3.0
		s.Topology[1][0] = default_density * 8.0 / 9.0
		s.Topology[2][1] = default_density * 2.0 / 3.0
		// feedback
		s.Topology[0][0] = default_density * 2.0 / 3.0
		s.Topology[1][1] = default_density * 2.0 / 3.0
	case 1:
		s.LenLayer[0] = slen * 3
		//feedforward
		s.DensityEM = default_density * 2.0 / 3.0
		s.Topology[1][0] = default_density * 2.0 / 3.0
		// feedback
		s.Topology[0][0] = default_density / 3.0
	case 0:
		s.DensityEM = default_density * 7.0
	default:
		log.Printf("SetLayer: unknown number of layers: %d\n", n)
		panic("SetLayer")
	}
}

func (s *Setting) SetDevelop(flag bool) {
	if flag {
		s.MaxDevelop = 200
		s.Alpha = 1.0 / 3.0
	} else {
		s.MaxDevelop = 1
		s.Alpha = 1.0
	}
}

func (s *Setting) SetModel(basename string) {
	m, ok := models[basename]
	if ok {
		s.Basename = basename
		s.SetLayer(m.nLayers)
		s.SetDevelop(m.develop)
		s.WithCue = m.cue
		s.SetOmega()
	} else {
		log.Fatal("Unknown model: " + basename)
	}
}
