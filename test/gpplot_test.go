package multicell_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

func TestGPPlot(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	s.Dump()
	envs := s.SaveEnvs(ENVSFILE, 50)

	env0 := envs[0]
	pop := s.NewPopulation(env0)
	pop, _ = pop.Evolve(s, env0)

	s.ProductionRun = true
	env1 := envs[1]
	pop.Iepoch = 1
	pop, _ = pop.Evolve(s, env1)

	file00 := s.TrajectoryFilename(1, 0, "traj.gz")
	pop0 := s.LoadPopulation(file00)
	file10 := s.TrajectoryFilename(1, s.MaxGeneration, "traj.gz")
	pop1 := s.LoadPopulation(file10)
	g0, gaxis := s.GetGenomeAxis(pop0, pop1)
	p0, paxis := s.GetPhenoAxis(pop0, pop1, env0, env1)
	for igen := range s.MaxGeneration {
		file := s.TrajectoryFilename(1, igen, "traj.gz")
		t0 := time.Now()
		pop := s.LoadPopulation(file)
		fmt.Println("GPPlot(load): ", time.Since(t0))
		t0 = time.Now()
		pop.GenoPhenoPlot(s, p0, paxis, g0, gaxis, env0)
		fmt.Println("GPPlot(project): ", time.Since(t0))
	}
}

func TestSVDProjection(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)

	env0 := envs[0]
	pop := s.NewPopulation(env0)
	pop, _ = pop.Evolve(s, env0)

	s.ProductionRun = true
	env1 := envs[1]
	pop.Iepoch = 1
	pop, _ = pop.Evolve(s, env1)

	file00 := s.TrajectoryFilename(1, 0, "traj.gz")
	pop0 := s.LoadPopulation(file00)
	file10 := s.TrajectoryFilename(1, s.MaxGeneration, "traj.gz")
	pop1 := s.LoadPopulation(file10)
	g0, gaxis := s.GetGenomeAxis(pop0, pop1)
	p0, paxis := s.GetPhenoAxis(pop0, pop1, env0, env1)
	c0, caxis := s.GetCueAxis(env1, env0)
	for igen := range s.MaxGeneration {
		file := s.TrajectoryFilename(1, igen, "traj.gz")
		pop := s.LoadPopulation(file)
		pop.SVDProject(s, p0, paxis, g0, gaxis, c0, caxis)
	}
}
