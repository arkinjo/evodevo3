package multicell_test

import (
	"fmt"
	//	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
	"math"
	"os"
	"path/filepath"
	//	"reflect"
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

var ENVSFILE = "envs.json"
var MODELS = []string{"Full", "NoCue", "Hie0", "Hie1", "Hie2", "NoDev", "Null",
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
	m0 := multicell.NewSpMat(10, 10)
	m0.Randomize(0.1)
	m1 := m0.Clone()
	if !m0.Equal(m1) {
		t.Errorf("Clone failed.")
	}

	for range 100 {
		m0.Mutate(0.1)
	}
	if m0.Equal(m1) {
		t.Errorf("Mutation failed.")
	}
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

func TestSettingDump(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.Outdir = "traj"
	s.Dump()
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
	if len(cell.S[s.NumLayers-1]) != lp {
		t.Errorf("len(cell.S[%d])=%d; want %d",
			s.NumLayers-1, len(cell.S[s.NumLayers-1]), lp)
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
	maxL := 0
	for l := range g.M {
		if maxL < l {
			maxL = l
		}
	}
	if maxL != s.NumLayers-1 {
		t.Errorf("maxL=%d; want %d", maxL, s.NumLayers-1)
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
	s := multicell.GetDefaultSetting()
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)
	s.SetModel("Full")
	pop := s.NewPopulation(envs[0])
	pop.Evolve(s, envs[0])
	pop.Dump(s)
}

func TestPopulationDumpJSON(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)
	s.SetModel("Full")
	pop := s.NewPopulation(envs[0])
	pop.Evolve(s, envs[0])
	pop.DumpJSON(s)
}

func TestProjection(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.SetModel("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)

	env := envs[0]
	pop := s.NewPopulation(env)
	pop, _ = pop.Evolve(s, env)

	s.ProductionRun = true
	env = envs[1]
	pop.Iepoch = 1
	pop, _ = pop.Evolve(s, envs[1])

	file00 := s.TrajectoryFilename(1, 0, "traj.gz")
	pop0 := s.LoadPopulation(file00)
	file10 := s.TrajectoryFilename(1, s.MaxGeneration, "traj.gz")
	pop1 := s.LoadPopulation(file10)
	g0, gaxis := s.GetGenomeAxis(pop0, pop1)
	sel0 := envs[0].SelectingEnv(s)
	sel1 := env.SelectingEnv(s)
	p0, paxis := s.GetPhenoAxis(sel0, sel1)
	c0, caxis := s.GetCueAxis(env, envs[0])
	for igen := range s.MaxGeneration {
		file := s.TrajectoryFilename(1, igen, "traj.gz")
		pop := s.LoadPopulation(file)
		pop.Project(s, p0, paxis, g0, gaxis, c0, caxis)
	}
}

func TestGenomeEqual(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.SetModel("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)

	pop0 := s.NewPopulation(envs[0])
	pop0, dumpfile := pop0.Evolve(s, envs[0])
	pop0.Sort()

	pop1 := s.LoadPopulation(dumpfile)
	pop1.Initialize(s, envs[1])
	selenv := envs[1].SelectingEnv(s)
	pop1.Develop(s, selenv)
	pop1.Sort()
	for i, indiv := range pop0.Indivs {
		if !indiv.Genome.Equal(&pop1.Indivs[i].Genome) {
			t.Errorf("Genomes are not equal")
		}
	}
}

func TestGenomeVecs(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.SetModel("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)

	pop0 := s.NewPopulation(envs[0])
	pop0, dumpfile := pop0.Evolve(s, envs[0])
	pop0.Sort()

	pop1 := s.LoadPopulation(dumpfile)
	pop1.Initialize(s, envs[1])
	selenv := envs[1].SelectingEnv(s)
	pop1.Develop(s, selenv)
	pop1.Sort()

	vecs0 := pop0.GenomeVecs(s)
	vecs1 := pop1.GenomeVecs(s)
	dvec := make(multicell.Vec, len(vecs0[0]))
	for i, v0 := range vecs0 {
		g0 := pop0.Indivs[i].Genome
		g1 := pop1.Indivs[i].Genome
		v1 := vecs1[i]
		v1.Diff(dvec, v0)
		del := dvec.Norm1()
		if del > 0 {
			t.Errorf("genome vecs %d differ by %f", i, del)
			fmt.Printf("genome equality: %d %t\n", i, g0.Equal(&g1))
			for k, x := range v0 {
				if x != v1[k] {
					fmt.Printf("%d\t%d\t%f\t%f\n", i, k, x, v1[k])
				}
			}
		}
	}
}

func TestSVD(t *testing.T) {
	M := mat.NewDense(2, 3, []float64{1, 2, 3, 4, 5, 6})
	var svd mat.SVD
	ok := svd.Factorize(M, mat.SVDThin)
	if !ok {
		t.Errorf("SVD failed")
	}
	var u, v mat.Dense
	sv := svd.Values(nil)
	svd.UTo(&u)
	svd.VTo(&v)
	u0 := u.ColView(0).(*mat.VecDense).RawVector().Data
	v0 := v.RawRowView(0)
	fmt.Println(len(sv), u0, v0, v)
}

func TestGenomeW(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.Outdir = "traj"
	s.SelStrength = 20
	s.SetModel("Full")
	envs := s.SaveEnvs(ENVSFILE, 50)
	pop := s.NewPopulation(envs[0])
	pop, _ = pop.Evolve(s, envs[0])
	for i, indiv := range pop.Indivs {
		fmt.Println(i, "\t", indiv.Genome.W)
	}

}
