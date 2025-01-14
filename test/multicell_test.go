package multicell_test

import (
	"fmt"
	//	"gonum.org/v1/gonum/blas/blas64"
	"math"
	"os"
	"path/filepath"
	//	"reflect"
	//	"slices"
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

var ENVSFILE = "envs.json"
var MODELS = []string{"Full", "NoCue", "NoHie", "Hie1", "Hie2", "NoDev", "Null",
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

func TestSpMatMutate(t *testing.T) {
	m0 := multicell.NewSpMat(200, 200)
	m0.Randomize(0.02)
	m1 := m0.Clone()
	if !m0.Equal(m1) {
		t.Errorf("Clone failed.")
	}

	m0.Mutate(0.001, 0.02)

	if m0.Equal(m1) {
		t.Errorf("Mutation failed.")
	}

	fmt.Printf("Densities %f %f\n", m0.Density(), m1.Density())
}

func TestSetting(t *testing.T) {
	for _, model := range MODELS {
		s := multicell.GetDefaultSetting(model)
		got := len(s.Topology.M)
		if got != s.NumLayers {
			t.Errorf("got len(Topology) = %d; want %d", got, s.NumLayers)
		}
	}
}

func TestSettingDump(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.Dump()
}

func TestCell(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	cell := s.NewCell(0)
	if len(cell.Cue) != multicell.NumFaces {
		t.Errorf("len(cell.Cue)= %d; want %d", len(cell.Cue), multicell.NumFaces)
	}
	for l, state := range cell.S {
		if len(state) != s.LenLayer[l] {
			t.Errorf("len(cell.S[%d]=%d; want %d",
				l, len(state), s.LenLayer[l])
		}
	}
	lp := s.LenLayer[s.NumLayers-1]
	if len(cell.S[s.NumLayers-1]) != lp {
		t.Errorf("len(cell.S[%d])=%d; want %d",
			s.NumLayers-1, len(cell.S[s.NumLayers-1]), lp)
	}
}

func TestIndividual(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
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
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.SetOmega()

	envs := s.SaveEnvs(ENVSFILE, 50)
	pop := s.NewPopulation(envs[0])
	if len(pop.Indivs) != s.MaxPopulation {
		t.Errorf("len(pop.Indivs)=%d; want %d", len(pop.Indivs), s.MaxPopulation)
	}
}

func TestModels(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)
	for _, model := range MODELS {
		fmt.Println("### Testing model: ", model)
		s.SetModel(model)
		pop := s.NewPopulation(envs[0])
		pop.Evolve(s, envs[0])
		defer myrecover(t, "Panic in Evolve, model: "+model)

	}

}

func TestPopulationDump(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)
	pop := s.NewPopulation(envs[0])
	pop.Evolve(s, envs[0])
	pop.Dump(s)
}

func TestPopulationDumpJSON(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)
	pop := s.NewPopulation(envs[0])
	pop.Evolve(s, envs[0])
	pop.DumpJSON(s)
}
