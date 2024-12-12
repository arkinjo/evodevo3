package multicell_test

import (
	"math"
	"os"
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

var ENVSFILE = "envs.json"

func cleanup() {
	os.Remove(ENVSFILE)
	os.Remove("traj/Full_00_010.traj.gz")
}

func TestMain(m *testing.M) {
	m.Run()
	cleanup()
}

func TestSetting(t *testing.T) {
	s := multicell.GetDefaultSetting()
	models := []string{"Full", "NoCue", "NoHie", "NoDev", "Null",
		"NullCue", "NullHie", "NullDev"}
	for _, model := range models {
		s.SetModel(model)
		got := len(s.Topology)
		if got != s.NumLayers {
			t.Errorf("got len(Topology) = %d; want %d", got, s.NumLayers)
		}
	}
}

func TestEnvironment(t *testing.T) {
	s := multicell.GetDefaultSetting()
	envs := s.SaveEnvs(ENVSFILE, 50)
	if len(envs) != 50 {
		t.Errorf("len(envs)= %d; want 50", len(envs))
	}
	for i, env := range envs {
		if env.Len() != s.LenFace*4 {
			t.Errorf("env[%d].Len()= %d; want %d", i, env.Len(), s.LenFace*4)
		}
	}
}

func TestCell(t *testing.T) {
	s := multicell.GetDefaultSetting()
	cell := s.NewCell(0)
	if len(cell.E) != multicell.NumFaces {
		t.Errorf("len(cell.E)= %d; want %d", len(cell.E), multicell.NumFaces)
	}
	for l, state := range cell.S {
		if len(state) != s.LenLayer[l] {
			t.Errorf("len(cell.S[%d]=%d; want %d",
				l, len(state), s.LenLayer[l])
		}
	}
	lp := s.LenLayer[s.NumLayers-1]
	if len(cell.Pave) != lp {
		t.Errorf("len(cell.Pave)=%d; want %d",
			len(cell.Pave), lp)
	}
	if len(cell.Pvar) != lp {
		t.Errorf("len(cell.Pvar)=%d; want %d",
			len(cell.Pvar), lp)
	}
}

func TestIndividual(t *testing.T) {
	s := multicell.GetDefaultSetting()
	envs := s.SaveEnvs(ENVSFILE, 50)
	indiv := s.NewIndividual(113, envs[0])
	if indiv.Id != 113 {
		t.Errorf("indiv.Id=%d; want 113", indiv.Id)
	}

	ncells := s.NumCellX * s.NumCellY
	if indiv.NumCells() != ncells {
		t.Errorf("indiv.NumCells()=%d; want %d", indiv.NumCells(), ncells)
	}
}

func TestGenome(t *testing.T) {
	s := multicell.GetDefaultSetting()
	g := s.NewGenome()
	if len(g.E) != multicell.NumFaces {
		t.Errorf("len(g.E)=%d; want %d", len(g.E), multicell.NumFaces)
	}
	if len(g.M) != s.NumLayers {
		t.Errorf("len(g.M)=%d; want %d", len(g.M), s.NumLayers)
	}

}

func TestActivation(t *testing.T) {
	lcatan := multicell.LCatan(1.0)
	v0 := lcatan(1.0)
	v1 := 6.0 * math.Atan(1.0/multicell.Sqrt3) / math.Pi
	if v0 != v1 {
		t.Errorf("lcatan(1.0)(1.0)=%f; want %f", v0, v1)
	}

	tanh := multicell.Tanh(2.0)
	if tanh(1.0) != math.Tanh(2.0*1.0) {
		t.Errorf("tanh(1.0) = %f; want %f", tanh(1.0), math.Tanh(2.0*1.0))
	}
}

func TestPopulation(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.SetOmega()
	s.Outdir = "traj"
	envs := s.SaveEnvs(ENVSFILE, 50)
	pop := s.NewPopulation(envs[0])
	if len(pop.Indivs) != s.MaxPopulation {
		t.Errorf("len(pop.Indivs)=%d; want %d", len(pop.Indivs), s.MaxPopulation)
	}
	s.MaxGeneration = 10
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic in Evolve")
		}
	}()

	pop.Evolve(s, envs[0])

}
