package multicell_test

import (
	"fmt"
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

func TestGenome(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	g := s.NewGenome()
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

func TestGenomeClone(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	g0 := s.NewGenome()
	g1 := g0.Clone()
	if !g0.Equal(g1) {
		t.Errorf("Genome cloning failed.")
	}
	g1.M[1][0].Randomize(0.1)
	if g0.M[1][0].Equal(g1.M[1][0]) {
		t.Errorf("Genome randomization failed (1).")
	}
	if !g0.M[2][1].Equal(g1.M[2][1]) {
		t.Errorf("Genome randomization failed (2).")
	}
}

func TestGenomeMutate(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.MutRate = 0.0005
	g0 := s.NewGenome()
	g1 := g0.Clone()
	if !g0.Equal(g1) {
		t.Errorf("Genome cloning failed")
	}
	g1.Mutate(s)
	v0 := g0.ToVec(s)
	v1 := g1.ToVec(s)
	dv := make(multicell.Vec, len(v1))
	dv.Diff(v0, v1)
	for i, d := range dv {
		if d != 0.0 {
			fmt.Printf("GenomeMutate: %d %f\n", i, d)
		}
	}
}

func TestGenomeEqual(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)

	pop0 := s.NewPopulation(envs[0])
	pop0, dumpfile := pop0.Evolve(s, envs[0])
	pop0.Sort()

	pop1 := s.LoadPopulation(dumpfile)
	pop1.Initialize(s, envs[1])
	pop1.Develop(s, envs[1])
	pop1.Sort()
	for i, indiv := range pop0.Indivs {
		if !indiv.Genome.Equal(pop1.Indivs[i].Genome) {
			t.Errorf("Genomes are not equal")
		}
	}
}

func TestGenomeVecs(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)

	pop0 := s.NewPopulation(envs[0])
	pop0, dumpfile := pop0.Evolve(s, envs[0])
	pop0.Sort()

	pop1 := s.LoadPopulation(dumpfile)
	pop1.Initialize(s, envs[1])
	pop1.Develop(s, envs[1])
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
			fmt.Printf("genome equality: %d %t\n", i, g0.Equal(g1))
			for k, x := range v0 {
				if x != v1[k] {
					fmt.Printf("%d\t%d\t%f\t%f\n", i, k, x, v1[k])
				}
			}
		}
	}
}
