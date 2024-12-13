package multicell_test

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

var ENVSFILE = "envs.json"
var MODELS = []string{"Full", "NoCue", "NoHie", "NoDev", "Null",
	"NullCue", "NullHie", "NullDev"}

func cleanup() {
	os.Remove(ENVSFILE)

	files, err := filepath.Glob("traj/*")
	if err != nil {
		panic("error in traj/*")
	}
	for _, f := range files {
		os.Remove(f)
	}
}

func myrecover(t *testing.T, msg string) {
	if r := recover(); r != nil {
		t.Errorf(msg)
	}
}

func TestMain(m *testing.M) {
	m.Run()
	cleanup()
}

func TestSetting(t *testing.T) {
	s := multicell.GetDefaultSetting()
	for _, model := range MODELS {
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
	s.Outdir = "traj"
	s.SetOmega()

	envs := s.SaveEnvs(ENVSFILE, 50)
	pop := s.NewPopulation(envs[0])
	if len(pop.Indivs) != s.MaxPopulation {
		t.Errorf("len(pop.Indivs)=%d; want %d", len(pop.Indivs), s.MaxPopulation)
	}
}

func TestModels(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.LoadEnvs(ENVSFILE)
	for _, model := range MODELS {
		fmt.Println("### Testing model: ", model)
		s.SetModel(model)
		pop := s.NewPopulation(envs[0])
		pop.Evolve(s, envs[0])
		defer myrecover(t, "Panic in Evolve, model: "+model)

	}

}

func TestProjection(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.SetModel("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 30
	envs := s.SaveEnvs(ENVSFILE, 50)

	env := envs[0]
	pop := s.NewPopulation(env)
	pop = pop.Evolve(s, env)

	s.ProductionRun = true
	env = envs[1]
	pop.Iepoch = 1
	pop = pop.Evolve(s, envs[1])

	file00 := s.TrajectoryFilename(1, 0, "traj.gz")
	pop0 := s.LoadPopulation(file00, env)
	file10 := s.TrajectoryFilename(1, s.MaxGeneration, "traj.gz")
	pop1 := s.LoadPopulation(file10, env)
	g0 := multicell.AverageVecs(pop0.GenomeVecs(s))
	p0 := envs[0].SelectingEnv(s)
	gaxis := s.GetGenomeAxis(&pop0, &pop1)
	paxis := s.GetPhenoAxis(envs[0], env)

	for igen := range s.MaxGeneration {
		file := s.TrajectoryFilename(1, igen, "traj.gz")
		pop := s.LoadPopulation(file, env)
		ofile := s.TrajectoryFilename(1, igen, "gpplot")
		pop.ProjectGenoPheno(s, ofile, g0, gaxis, p0, paxis)
	}
}
