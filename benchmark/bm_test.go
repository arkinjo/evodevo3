package multicell_test

import (
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

func BenchmarkFullTrain(b *testing.B) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxPopulation = 500
	s.MaxGeneration = 200
	envs := s.SaveEnvs("envs.json", 10)
	pop := s.NewPopulation(envs[0])
	s.ProductionRun = false
	pop.Evolve(s, envs[0])
	s.ProductionRun = true
	pop.Evolve(s, envs[1])

}
