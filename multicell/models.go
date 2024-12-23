package multicell

import (
	"log"
)

func (s *Setting) FullModel() {
	s.Basename = "Full"
	s.NumLayers = 4
	s.LenLayer = make([]int, s.NumLayers)
	topology := NewTopology(s.NumLayers)
	s.DensityEM = default_density
	for l := range s.NumLayers {
		s.LenLayer[l] = default_len_state

		// feedforward
		if l > 0 {
			topology[l][l-1] = default_density
		}
		// feedback
		topology[l][l] = default_density
	}

	s.Topology = topology
	s.SetOmega()
}

func (s *Setting) NoCueModel() {
	s.Basename = "NoCue"
	s.WithCue = false
}

func (s *Setting) NoDevModel() {
	s.Basename = "NoDev"
	s.MaxDevelop = 1
	s.Alpha = 1.0
}

func (s *Setting) OriginalModel() {
	s.Basename = "Original"
	s.NumLayers = 4
	topology := NewTopology(s.NumLayers)

	// feedforward
	topology[1][0] = default_density
	topology[2][1] = default_density
	topology[3][2] = default_density

	// feedback
	topology[0][1] = default_density
	topology[2][2] = default_density
	s.Topology = topology
	s.SetOmega()
}

func (s *Setting) NoHieModel() {
	s.Basename = "NoHie"
	s.NumLayers = 1
	s.LenLayer = []int{default_len_state}
	topology := NewTopology(s.NumLayers)

	//feedforward
	s.DensityEM = default_density * 4.0

	// feedback
	topology[0][0] = default_density * 4.0

	s.Topology = topology
	s.SetOmega()
}

func (s *Setting) Hie1Model() {
	s.Basename = "Hie1"
	s.NumLayers = 2
	s.LenLayer = []int{600, 200}
	topology := NewTopology(s.NumLayers)

	//feedforward
	s.DensityEM = default_density * 2.0 / 3.0
	topology[1][0] = default_density * 2.0 / 3.0
	// feedback
	topology[0][0] = default_density / 3.0
	topology[1][1] = default_density

	s.Topology = topology
	s.SetOmega()
}

func (s *Setting) Hie2Model() {
	s.Basename = "Hie2"
	s.NumLayers = 3
	s.LenLayer = []int{300, 300, 200}
	topology := NewTopology(s.NumLayers)

	//feedforward
	s.DensityEM = default_density * 2.0 / 3.0
	topology[1][0] = default_density * 8.0 / 9.0
	topology[2][1] = default_density * 2.0 / 3.0
	// feedback
	topology[0][0] = default_density * 2.0 / 3.0
	topology[1][1] = default_density * 2.0 / 3.0
	topology[2][2] = default_density

	s.Topology = topology
	s.SetOmega()
}

func (s *Setting) NullModel() {
	s.NoHieModel()
	s.WithCue = false
	s.MaxDevelop = 1
	s.Alpha = 1.0
	s.Basename = "Null"
	s.SetOmega()
}

func (s *Setting) NullCueModel() {
	s.NullModel()
	s.WithCue = true
	s.Basename = "NullCue"
}

func (s *Setting) NullDevModel() {
	s.NullModel()
	s.MaxDevelop = 200
	s.Basename = "NullDev"
}

func (s *Setting) NullHieModel() {
	s.WithCue = false
	s.MaxDevelop = 1
	s.Basename = "NullHie"
}

func (s *Setting) SetModel(model string) {
	switch model {
	case "Full":
		s.FullModel()
	case "NoHie":
		s.NoHieModel()
	case "Hie1":
		s.Hie1Model()
	case "Hie2":
		s.Hie2Model()
	case "NoCue":
		s.NoCueModel()
	case "NoDev":
		s.NoDevModel()
	case "Null":
		s.NullModel()
	case "NullCue":
		s.NullCueModel()
	case "NullHie":
		s.NullHieModel()
	case "NullDev":
		s.NullDevModel()
	case "Original":
		s.OriginalModel()
	default:
		log.Println("SetModel: invalid model name. Must be one of Full, NoCue, NoHie, NoDev, Null, NullCue, NullHie, NullDev\n")
		log.Fatal("Invalid model: " + model)
	}

	s.SetOmega()
}
