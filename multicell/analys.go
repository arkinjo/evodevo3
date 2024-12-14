package multicell

import (
	"fmt"
	"os"
)

func (s *Setting) GetPhenoAxis(env0, env1 Environment) Vec {
	senv0 := env0.SelectingEnv(s)
	senv1 := env1.SelectingEnv(s)
	dv := make(Vec, len(senv0))
	DiffVecs(senv1, senv0, dv)
	mag2 := DotVecs(dv, dv)
	VecScale(dv, 1/mag2)

	return dv
}

func AverageVecs(vecs []Vec) Vec {
	ave := make(Vec, len(vecs[0]))
	for _, vec := range vecs {
		for i, v := range vec {
			ave[i] += v
		}
	}
	n := float64(len(vecs))
	VecScale(ave, 1/n)
	return ave
}

func (s *Setting) GetGenomeAxis(pop0, pop1 *Population) Vec {
	g0 := AverageVecs(pop0.GenomeVecs(s))
	g1 := AverageVecs(pop1.GenomeVecs(s))
	dg := make(Vec, len(g0))
	DiffVecs(g0, g1, dg)
	mag2 := DotVecs(dg, dg)
	VecScale(dg, 1/mag2)
	return dg
}

func (pop *Population) ProjectGenoPheno(s *Setting, filename string,
	g0, gaxis, p0, paxis Vec) {
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()

	for i, indiv := range pop.Indivs {
		gt := indiv.Genome.ToVec(s)
		pt := indiv.SelectedPhenotypeVec(s)
		DiffVecs(gt, g0, gt)
		DiffVecs(pt, p0, pt)
		gc := DotVecs(gt, gaxis)
		pc := DotVecs(pt, paxis)
		fmt.Fprintf(fout, "%d\t%f\t%f\n", i, gc, pc)
	}
}
