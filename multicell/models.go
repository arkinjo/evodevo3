package multicell

import (
	"log"
)

func (s *Setting) FullModel() {
	s.Basename = "Full"
}

func (s *Setting) NoCueModel() {
	s.Basename = "NoCue"
	s.With_cue = false
}

func (s *Setting) NoDevModel() {
	s.Basename = "NoDev"
	s.Max_dev = 1
}

func (s *Setting) NoHieModel() {
	s.Basename = "NoHie"
	s.Num_layers = 3
	s.Num_components = []int{200, 600, 200}
	topology := NewSpMat(3)

	topology[1][0] = default_density * 2.0 / 3.0 //(1,0) and (2,1) in Full
	topology[1][1] = default_density / 3.0       // (1,1), (2,2), (3,3) in Full
	topology[2][1] = default_density * 2.0 / 3.0 //(3,2), (4,3) in Full
	s.Topology = topology
	s.Omega = make(Vec, 3)
	s.SetOmega()
}

func (s *Setting) NullModel() {
	s.NoHieModel()
	s.With_cue = false
	s.Max_dev = 1
	s.Basename = "Null"
	s.SetOmega()
}

func (s *Setting) NullCueModel() {
	s.NullModel()
	s.With_cue = true
	s.Basename = "NullCue"
}

func (s *Setting) NullDevModel() {
	s.NullModel()
	s.Max_dev = 200
	s.Basename = "NullDev"
}

func (s *Setting) NullHieModel() {
	s.With_cue = false
	s.Max_dev = 1
	s.Basename = "NullHie"
}

func (s *Setting) SetModel(model string) {
	switch model {
	case "Full":
		s.FullModel()
	case "NoHie":
		s.NoHieModel()
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
	default:
		log.Println("SetModel: invalid model name. Must be one of Full, NoCue, NoHie, NoDev, Null, NullCue, NullHie, NullDev\n")
		log.Fatal("Invalid model")

	}
}
