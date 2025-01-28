package multicell

import (
	"log"
	"math"
)

type model_t struct {
	cue     bool
	develop bool
	nLayers int
}

var models = map[string]model_t{
	"Full":  {true, true, 3},
	"NoCue": {false, true, 3},
	"NoDev": {true, false, 3},
	"NoHie": {true, true, 0},

	"Hie2": {true, true, 2},
	"Hie1": {true, true, 1},

	"NullHie": {false, false, 3},
	"NullCue": {true, false, 0},
	"NullDev": {false, true, 0},
	"Null":    {false, false, 0},
}

// Self-loop for hidden layers
func (s *Setting) SetLayerSL(n int) {
	s.NumLayers = n + 2
	s.LenLayer = make([]int, s.NumLayers)
	slen := s.LenFace * NumFaces
	s.LenLayer[0] = slen
	s.LenLayer[s.NumLayers-1] = slen

	s.Topology = NewSliceOfMaps[float64](s.NumLayers)

	switch n {
	case 3:
		s.LenLayer[1] = slen
		s.LenLayer[2] = slen
		s.LenLayer[3] = slen

		// feedforward
		s.Topology.Set(1, 0, default_density)
		s.Topology.Set(2, 1, default_density)
		s.Topology.Set(3, 2, default_density)
		s.Topology.Set(4, 3, default_density)
		// feedback
		s.Topology.Set(1, 1, default_density)
		s.Topology.Set(2, 2, default_density)
		s.Topology.Set(3, 3, default_density)
	case 2:
		s.LenLayer[1] = slen * 3 / 2
		s.LenLayer[2] = slen * 3 / 2
		// feedforward
		s.Topology.Set(1, 0, default_density*2.0/3.0)
		s.Topology.Set(2, 1, default_density*8.0/9.0)
		s.Topology.Set(3, 2, default_density*2.0/3.0)
		// feedback
		s.Topology.Set(1, 1, default_density*2.0/3.0)
		s.Topology.Set(2, 2, default_density*2.0/3.0)
	case 1:
		s.LenLayer[1] = slen * 3
		// feedforward
		s.Topology.Set(1, 0, default_density*2.0/3.0)
		s.Topology.Set(2, 1, default_density*2.0/3.0)
		// feedback
		s.Topology.Set(1, 1, default_density/3.0)
	case 0:
		s.Topology.Set(1, 0, default_density*7.0)
	default:
		log.Printf("SetLayer: unknown number of layers: %d\n", n)
		panic("SetLayer")
	}
}

// Feedback loop to the previous layer
func (s *Setting) SetLayerM1(n int) {
	s.NumLayers = n + 2
	s.LenLayer = make([]int, s.NumLayers)
	slen := s.LenFace * NumFaces
	s.LenLayer[0] = slen
	s.LenLayer[s.NumLayers-1] = slen

	s.Topology = NewSliceOfMaps[float64](s.NumLayers)

	switch n {
	case 3:
		s.LenLayer[1] = slen
		s.LenLayer[2] = slen
		s.LenLayer[3] = slen

		// feedforward
		s.Topology.Set(1, 0, default_density)
		s.Topology.Set(2, 1, default_density)
		s.Topology.Set(3, 2, default_density)
		s.Topology.Set(4, 3, default_density)
		// feedback
		s.Topology.Set(1, 2, default_density)
		s.Topology.Set(2, 3, default_density)

	case 2:
		s.LenLayer[1] = slen * 3 / 2
		s.LenLayer[2] = slen * 3 / 2
		//feedforward
		s.Topology.Set(1, 0, default_density*2.0/3.0)
		s.Topology.Set(2, 1, default_density*8.0/9.0)
		s.Topology.Set(3, 2, default_density*2.0/3.0)
		// feedback
		s.Topology.Set(1, 2, default_density*8.0/9.0)

	case 1:
		s.LenLayer[1] = slen * 3
		//feedforward
		s.Topology.Set(1, 0, default_density*2.0/3.0)
		s.Topology.Set(2, 1, default_density*2.0/3.0)
		// feedback
		s.Topology.Set(1, 1, default_density*2.0/9.0)
	case 0:
		s.Topology.Set(1, 0, default_density*6.0)
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

func (s *Setting) SetOmega() {
	s.Omega = make(Vec, s.NumLayers)

	if s.WithCue {
		s.Omega[0] = 2.0
	} else {
		s.Omega[0] = 1.0
	}
	s.Topology.Do(func(l, k int, density float64) {
		s.Omega[l] += density * float64(s.LenLayer[k])
	})

	for l, omega := range s.Omega {
		if omega > 0 {
			s.Omega[l] = 1.0 / math.Sqrt(omega)
		}
	}
}

func (s *Setting) SetModel(basename string) {
	m, ok := models[basename]
	if ok {
		s.Basename = basename
		s.SetLayerSL(m.nLayers)
		s.SetDevelop(m.develop)
		s.WithCue = m.cue
		s.SetOmega()
	} else {
		log.Fatal("Unknown model: " + basename)
	}
}
