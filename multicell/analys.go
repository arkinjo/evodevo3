package multicell

import (
	"fmt"
	"log"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"
)

func (s *Setting) GetPhenoAxis(env0, env1 Environment) Vec {
	senv0 := env0.SelectingEnv(s)
	senv1 := env1.SelectingEnv(s)
	dv := make(Vec, len(senv0))
	dv.Diff(senv1, senv0)
	mag2 := DotVecs(dv, dv)
	dv.ScaleBy(1 / mag2)

	return dv
}

func avesd(vec Vec) (float64, float64) {
	a := 0.0
	n := float64(len(vec))
	for _, x := range vec {
		a += x
	}

	a /= n
	v := 0.0
	for _, x := range vec {
		v += (x - a) * (x - a)
	}
	v = math.Sqrt(v / n)

	return a, v
}

func AverageVecs(vecs []Vec) Vec {
	ave := make(Vec, len(vecs[0]))
	for _, vec := range vecs {
		for i, v := range vec {
			ave[i] += v
		}
	}
	n := float64(len(vecs))
	ave.ScaleBy(1 / n)
	return ave
}

func (s *Setting) GetGenomeAxis(pop0, pop1 *Population) Vec {
	g0 := AverageVecs(pop0.GenomeVecs(s))
	g1 := AverageVecs(pop1.GenomeVecs(s))
	dg := make(Vec, len(g0))
	dg.Diff(g1, g0)
	mag2 := DotVecs(dg, dg)
	dg.ScaleBy(1 / mag2)
	return dg
}

func (pop *Population) ProjectGenoPheno(s *Setting, filename string,
	g0, gaxis, p0, paxis Vec) {
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()

	var gs, ps, ns, ms Vec

	for _, indiv := range pop.Indivs {
		gt := indiv.Genome.ToVec(s)
		pt := indiv.SelectedPhenotypeVec(s)
		gt.Diff(gt, g0)
		pt.Diff(pt, p0)
		gs = append(gs, DotVecs(gt, gaxis))
		ps = append(ps, DotVecs(pt, paxis))
		ns = append(ns, float64(indiv.Ndev))
		ms = append(ms, indiv.Mismatch)
	}
	ga, gv := avesd(gs)
	pa, pv := avesd(ps)
	na, nv := avesd(ns)
	ma, mv := avesd(ms)
	fmt.Fprintf(fout, "P\t%d\t%d\t%f\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n",
		pop.Igen, len(gs), ga, pa, gv, pv, na, nv, ma, mv)
	for i, indiv := range pop.Indivs {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f\t%d\t%f\n",
			i, gs[i], ps[i], indiv.Ndev, ms[i])
	}
}

func CovarianceMatrix(xs []Vec, x0 Vec, ys []Vec, y0 Vec) *mat.Dense {
	if len(xs) != len(ys) {
		log.Printf("CovarianceMatrix: size mismatch %d != %d\n",
			len(xs), len(ys))
		panic("CovarianceMatrix")
	}
	nx := len(x0)
	ny := len(y0)
	xt := make(Vec, nx)
	yt := make(Vec, ny)
	cov := mat.NewDense(nx, ny, nil)
	for n, x := range xs {
		xt.Diff(x, x0)
		yt.Diff(ys[n], y0)
		for i := range nx {
			for j := range ny {
				cov.Set(i, j, xt[i]*yt[j])
			}
		}
	}

	return cov
}
