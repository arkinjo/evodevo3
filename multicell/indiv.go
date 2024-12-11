package multicell

import (
	//	"fmt"
	"math"
)

type Individual struct {
	Id       int
	MomId    int
	DadId    int
	Genome   Genome
	Cells    [][]Cell
	Ndev     int
	Mismatch float64
	Fitness  float64
}

func (s *Setting) SetCellEnv(cells [][]Cell, env Environment) {
	envs := s.NewCellEnvs(env)

	var left, right, top, bottom Vec
	for i, cs := range cells {
		for j := range cs {
			if i == 0 {
				left = envs.Lefts[j]
			} else {
				left = cells[i-1][j].Pave[Right]
			}

			if j == s.NumCellY-1 {
				top = envs.Tops[i]
			} else {
				top = cells[i][j+1].Pave[Bottom]
			}

			if i == s.NumCellX-1 {
				right = envs.Rights[j]
			} else {
				right = cells[i+1][j].Pave[Left]
			}

			if j == 0 {
				bottom = envs.Bottoms[i]
			} else {
				bottom = cells[i][j-1].Pave[Top]
			}

			cells[i][j].E[Left] = left
			cells[i][j].E[Top] = top
			cells[i][j].E[Right] = right
			cells[i][j].E[Bottom] = bottom
		}
	}
}

func (indiv *Individual) CueVec(s *Setting) Vec {
	var vec Vec

	for j := range s.NumCellY {
		vec = append(vec, indiv.Cells[0][j].E[Left]...)
	}
	for i := range s.NumCellX {
		vec = append(vec, indiv.Cells[i][s.NumCellY-1].E[Top]...)
	}
	for j := range s.NumCellY {
		vec = append(vec, indiv.Cells[s.NumCellX-1][j].E[Right]...)
	}
	for i := range s.NumCellX {
		vec = append(vec, indiv.Cells[i][0].E[Bottom]...)
	}

	return vec
}

func (s *Setting) NewIndividual(id int, env Environment) Individual {
	cells := make([][]Cell, s.NumCellX)
	for i := range s.NumCellX {
		cells[i] = make([]Cell, s.NumCellY)
		for j := range s.NumCellY {
			cells[i][j] = s.NewCell()
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

func (indiv *Individual) SelectedPhenotype() Vec {
	var p Vec
	for _, cell := range indiv.Cells[0] {
		p = append(p, cell.Pave[Left]...)
	}

	return p
}

func (indiv *Individual) Initialize(s *Setting, env Environment) {
	for i, cs := range indiv.Cells {
		for j := range cs {
			indiv.Cells[i][j].Initialize(s)
		}
	}
	indiv.Ndev = 0
	s.SetCellEnv(indiv.Cells, env)
}

func (indiv *Individual) GetMismatch(s *Setting, selenv Environment) float64 {
	selphen := indiv.SelectedPhenotype()
	dev := 0.0
	for i, e := range selphen {
		dev += math.Abs(e - selenv[i])
	}

	return dev / float64(len(selenv))
}

func (indiv *Individual) Develop(s *Setting, selenv Environment) Individual {
	istep := 0
	for istep = range s.MaxDevelop {
		dev := 0.0
		for i := range s.NumCellX {
			for j := range s.NumCellY {
				dev += indiv.Cells[i][j].DevStep(s, indiv.Genome, istep)
			}
		}
		if dev < s.ConvDevelop {
			break
		}
	}

	indiv.Ndev = istep + 1
	indiv.Mismatch = indiv.GetMismatch(s, selenv)

	if istep < s.MaxDevelop {
		indiv.Fitness = math.Exp(-indiv.Mismatch * s.Selstrength)
	} else {
		indiv.Fitness = 0
	}
	return *indiv
}

func (s *Setting) MateIndividuals(indiv0, indiv1 Individual, env Environment) (Individual, Individual) {
	g0, g1 := s.MateGenomes(indiv0.Genome, indiv1.Genome)
	kid0 := s.NewIndividual(-1, env)
	kid1 := s.NewIndividual(-2, env)

	kid0.Genome = g0
	kid0.MomId = indiv0.Id
	kid0.DadId = indiv1.Id

	kid1.Genome = g1
	kid0.MomId = indiv1.Id
	kid0.DadId = indiv0.Id

	return kid0, kid1
}
