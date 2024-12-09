package multicell

import (
	//	"fmt"
	"math"
)

type Individual struct {
	Id       int
	Mom_id   int
	Dad_id   int
	Genome   Genome
	Cells    [][]Cell
	Ndev     int
	Mismatch float64
	Fitness  float64
}

func (s *Setting) Set_cell_env(cells [][]Cell, env Environment) {
	envs := s.NewCell_envs(env)

	var left, right, top, bottom Vec
	for i, cs := range cells {
		for j := range cs {
			if i == 0 {
				left = envs.Lefts[j]
			} else {
				left = cells[i-1][j].Pr
			}
			if i == s.Num_cell_x-1 {
				right = envs.Rights[j]
			} else {
				right = cells[i+1][j].Pl
			}

			if j == 0 {
				bottom = envs.Bottoms[i]
			} else {
				bottom = cells[i][j-1].Pt
			}
			if j == s.Num_cell_y-1 {
				top = envs.Tops[i]
			} else {
				top = cells[i][j+1].Pb
			}

			cells[i][j].El = left
			cells[i][j].Et = top
			cells[i][j].Er = right
			cells[i][j].Eb = bottom
		}
	}
}

func (s *Setting) NewIndividual(id int, env Environment) Individual {
	cells := make([][]Cell, s.Num_cell_x)
	for i := range s.Num_cell_x {
		cells[i] = make([]Cell, s.Num_cell_y)
		for j := range s.Num_cell_y {
			cells[i][j] = s.NewCell()
		}
	}

	s.Set_cell_env(cells, env)

	return Individual{
		Id:       id,
		Mom_id:   -1,
		Dad_id:   -1,
		Genome:   s.NewGenome(),
		Cells:    cells,
		Ndev:     0,
		Mismatch: 100000.0,
		Fitness:  0}
}

func (indiv *Individual) Selected_pheno() Vec {
	var p Vec
	for _, cell := range indiv.Cells[0] {
		p = append(p, cell.Pl...)
	}

	return p
}

func (indiv *Individual) Initialize(s *Setting, env Environment) {
	for i, cs := range indiv.Cells {
		for j := range cs {
			indiv.Cells[i][j].Initialize(s)
		}
	}
	s.Set_cell_env(indiv.Cells, env)
}

func (indiv *Individual) Get_mismatch(s *Setting, selenv Environment) float64 {
	selphen := indiv.Selected_pheno()
	dev := 0.0
	for i, e := range selphen {
		dev += math.Abs(e - selenv[i])
	}

	return dev / float64(len(selenv))
}

func (indiv *Individual) Develop(s *Setting, selenv Environment) Individual {
	istep := 0
	for istep = range s.Max_dev {
		dev := 0.0
		for i := range s.Num_cell_x {
			for j := range s.Num_cell_y {
				dev += indiv.Cells[i][j].Dev_step(s, indiv.Genome, istep)
			}
		}
		if dev < s.Conv_dev {
			break
		}
	}

	indiv.Ndev = istep + 1
	indiv.Mismatch = indiv.Get_mismatch(s, selenv)

	if istep < s.Max_dev {
		indiv.Fitness = math.Exp(-indiv.Mismatch * s.Selstrength)
	} else {
		indiv.Fitness = 0
	}
	return *indiv
}

func (s *Setting) MateIndividuals(indiv0, indiv1 Individual, env Environment) (Individual, Individual) {
	geno0, geno1 := s.MateGenomes(indiv0.Genome, indiv1.Genome)
	kid0 := s.NewIndividual(-1, env)
	kid1 := s.NewIndividual(-2, env)

	kid0.Genome = geno0
	kid0.Mom_id = indiv0.Id
	kid0.Dad_id = indiv1.Id

	kid1.Genome = geno1
	kid0.Mom_id = indiv1.Id
	kid0.Dad_id = indiv0.Id

	return kid0, kid1
}
