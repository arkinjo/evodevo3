package multicell_test

import (
	"fmt"
	"github.com/arkinjo/evodevo3/multicell"

	"testing"
	"time"
)

func TestSVDTrans(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.MaxGeneration = 10
	envs := s.SaveEnvs(ENVSFILE, 50)
	s.SetModel("Full")
	pop := s.NewPopulation(envs[0])
	pop.Evolve(s, envs[0])

	gvecs := pop.GenomeVecs(s)
	pvecs := pop.PhenoVecs(s, true)
	mg := multicell.MeanVecs(gvecs)
	mp := multicell.MeanVecs(pvecs)
	t0 := time.Now()
	multicell.XPCA(pvecs, mp, gvecs, mg)
	fmt.Println("XPCA(p,g): ", time.Since(t0))

	//	t1 := time.Now()
	//	multicell.XPCA(gvecs, mg, pvecs, mp)
	//	fmt.Println("XPCA(p,g): ", time.Since(t1))
}
