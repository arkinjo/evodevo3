package multicell

import (
	//	"log"
	"math"
	"slices"
)

type Individual struct {
	Id       int
	MomId    int
	DadId    int
	Genome   Genome
	Cells    []Cell
	Ndev     int
	Mismatch float64
	Fitness  float64
}

func (indiv *Individual) NumCells() int {
	return len(indiv.Cells)
}

func (s *Setting) SetCellEnv(cells []Cell, env Environment) {
	nenv := env.AddNoise(s) // noisy environment.
	for i, c := range cells {
		for iface, iop := range c.Facing {
			if iop < 0 {
				cells[i].E[iface] = nenv.Face(s, iface)
			} else {
				cells[i].E[iface] = cells[iop].OppositeFace(s, iface)
			}
		}
	}
}

func (indiv *Individual) CueVec(s *Setting) Vec {
	var vec Vec

	for _, c := range indiv.Cells {
		for iface, iopp := range c.Facing {
			if iopp < 0 {
				vec = append(vec, c.E[iface]...)
			}
		}
	}
	return vec
}

func (s *Setting) CellId(i, j int) int {
	return i*s.NumCellY + j
}

func (s *Setting) NewIndividual(id int, env Environment) Individual {
	cells := make([]Cell, s.NumCellX*s.NumCellY)
	for i := range s.NumCellX {
		for j := range s.NumCellY {
			id := s.CellId(i, j)
			cells[id] = s.NewCell(id)
			if i > 0 {
				cells[id].Facing[Left] = s.CellId(i-1, j)
			}
			if i < s.NumCellX-1 {
				cells[id].Facing[Right] = s.CellId(i+1, j)
			}
			if j > 0 {
				cells[id].Facing[Bottom] = s.CellId(i, j-1)
			}
			if j < s.NumCellY-1 {
				cells[id].Facing[Top] = s.CellId(i, j+1)
			}
		}
	}

	s.SetCellEnv(cells, env)

	return Individual{
		Id:    id,
		MomId: -1,
		DadId: -1,
		//Genome:   s.NewGenome(), given later
		Cells:    cells,
		Ndev:     0,
		Mismatch: 100000.0,
		Fitness:  0}
}

func (indiv *Individual) SelectedPhenotype(s *Setting) []Vec {
	var p []Vec
	for _, c := range indiv.Cells {
		if c.Facing[Left] < 0 {
			p = append(p, c.Left(s))
		}
	}

	return p
}

func (indiv *Individual) SelectedPhenotypeVec(s *Setting) Vec {
	return slices.Concat(indiv.SelectedPhenotype(s)...)
}

func (indiv *Individual) Initialize(s *Setting, env Environment) {
	for i := range indiv.Cells {
		indiv.Cells[i].Initialize(s)
	}
	indiv.Ndev = 0
	s.SetCellEnv(indiv.Cells, env)
}

func (indiv *Individual) GetMismatch(s *Setting, selenv Vec) float64 {
	selphen := indiv.SelectedPhenotype(s)
	dv := make(Vec, len(selenv))
	dev := 0.0
	for _, p := range selphen {
		dv.Diff(p, selenv)
		dev += dv.Norm1()
	}

	return dev / float64(len(selenv)*len(selphen))
}

func (indiv *Individual) Develop(s *Setting, selenv Vec) Individual {
	istep := 0
	for istep = range s.MaxDevelop {
		dev := 0.0
		for i := range indiv.Cells {
			dev += indiv.Cells[i].DevStep(s, indiv.Genome, istep)
		}
		if dev < s.ConvDevelop {
			break
		}
	}

	indiv.Ndev = istep + 1
	indiv.Mismatch = indiv.GetMismatch(s, selenv)

	if istep < s.MaxDevelop {
		indiv.Fitness = math.Exp(-indiv.Mismatch * s.SelStrength)
	} else {
		indiv.Fitness = 0
	}
	return *indiv
}

func (s *Setting) MateIndividuals(indiv0, indiv1 Individual, env Environment) (Individual, Individual) {
	g0, g1 := indiv0.Genome.MateWith(indiv1.Genome)
	kid0 := s.NewIndividual(-1, env)
	kid1 := s.NewIndividual(-2, env)

	kid0.MomId = indiv0.Id
	kid0.DadId = indiv1.Id
	kid0.Genome = g0
	kid0.Genome.Mutate(s)

	kid0.MomId = indiv1.Id
	kid0.DadId = indiv0.Id
	kid1.Genome = g1
	kid1.Genome.Mutate(s)

	return kid0, kid1
}
