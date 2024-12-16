package multicell

import (
	"fmt"
	"log"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"
)

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

func MeanVecs(vecs []Vec) Vec {
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

func GetAxis(v0, v1 Vec) Vec {
	dv := make(Vec, len(v0))
	dv.Diff(v1, v0)
	mag2 := DotVecs(dv, dv)
	dv.ScaleBy(1 / mag2)
	return dv
}

func (s *Setting) GetPhenoAxis(env0, env1 Environment) (Vec, Vec) {
	senv0 := env0.SelectingEnv(s)
	senv1 := env1.SelectingEnv(s)
	return senv0, GetAxis(senv0, senv1)
}

func (s *Setting) GetGenomeAxis(pop0, pop1 Population) (Vec, Vec) {
	g0 := MeanVecs(pop0.GenomeVecs(s))
	g1 := MeanVecs(pop1.GenomeVecs(s))
	return g0, GetAxis(g0, g1)
}

func (s *Setting) GetCueAxis(env0, env1 Environment) (Vec, Vec) {
	return env1, GetAxis(env0, env1)
}

func ProjectOnAxis(vecs []Vec, v0 Vec, axis Vec) Vec {
	vt := make(Vec, len(v0))
	var ps Vec
	for _, v := range vecs {
		vt.Diff(v, v0)
		ps = append(ps, DotVecs(vt, axis))
	}
	return ps
}

func (pop *Population) Project(s *Setting, filename string,
	p0, paxis, g0, gaxis, c0, caxis Vec) {
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()

	// for alignment calculation
	punit := paxis.Copy()
	punit.Normalize()

	gvecs := pop.GenomeVecs(s)
	pvecs := pop.PhenoVecs(s)
	// Geno-Pheno Plot
	gs := ProjectOnAxis(gvecs, g0, gaxis)
	ps := ProjectOnAxis(pvecs, p0, paxis)

	mp := MeanVecs(pvecs)
	mg := MeanVecs(gvecs)

	// Pheno-Geno Cross-Covariance
	svGeno, uGeno, vGeno := XPCA(pvecs, mp, paxis, gvecs, mg, gaxis)
	rsvGeno0 := svGeno[0] / svGeno.Norm2()

	aliG := DotVecs(uGeno, punit)

	Pgs := ProjectOnAxis(pvecs, mp, uGeno)
	Ggs := ProjectOnAxis(gvecs, mg, vGeno)
	if DotVecs(Ggs, Pgs) < 0 {
		Ggs.ScaleBy(-1)
	}
	// Pheno-Cue Cross-Covariance
	cvecs := pop.CueVecs(s)
	mc := c0 //MeanVecs(cvecs)
	svCue, uCue, vCue := XPCA(pvecs, mp, paxis, cvecs, mc, caxis)

	rsvCue0 := svCue[0] / svCue.Norm2()
	aliC := DotVecs(uCue, punit)

	Pcs := ProjectOnAxis(pvecs, mp, uCue)
	Ccs := ProjectOnAxis(cvecs, mc, vCue)
	if DotVecs(Ccs, Pcs) < 0 {
		Ccs.ScaleBy(-1)
	}

	var ns, ms Vec
	for _, indiv := range pop.Indivs {
		ns = append(ns, float64(indiv.Ndev))
		ms = append(ms, indiv.Mismatch)
	}
	ga, gv := avesd(gs)
	pa, pv := avesd(ps)
	na, nv := avesd(ns)
	ma, mv := avesd(ms)
	fmt.Fprintf(fout, "P\t%d\t%d", pop.Igen, len(gs))
	fmt.Fprintf(fout, "\t%f\t%f", ga, pa)
	fmt.Fprintf(fout, "\t%f\t%f", gv, pv)
	fmt.Fprintf(fout, "\t%f\t%f", na, nv)
	fmt.Fprintf(fout, "\t%f\t%f", ma, mv)
	fmt.Fprintf(fout, "\t%f\t%f\t%f", svGeno[0], rsvGeno0, aliG)
	fmt.Fprintf(fout, "\t%f\t%f\t%f", svCue[0], rsvCue0, aliC)
	fmt.Fprintf(fout, "\n")

	for i := range pop.Indivs {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f", i, gs[i], ps[i])
		fmt.Fprintf(fout, "\t%f\t%f", Ggs[i], Pgs[i])
		fmt.Fprintf(fout, "\t%f\t%f", Ccs[i], Pcs[i])
		fmt.Fprintf(fout, "\n")
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
				cov.Set(i, j, cov.At(i, j)+xt[i]*yt[j])
			}
		}
	}

	fac := 1.0 / float64(len(xs))
	cov.Scale(fac, cov)

	return cov
}

// Get singular values, the first left and right singular vectors.
func XPCA(xs []Vec, x0, xaxis Vec, ys []Vec, y0, yaxis Vec) (Vec, Vec, Vec) {
	ccov := CovarianceMatrix(xs, x0, ys, y0)
	var svd mat.SVD
	ok := svd.Factorize(ccov, mat.SVDThin)
	if !ok {
		log.Fatal("SVD failed")
	}
	var u, v mat.Dense
	sv := svd.Values(nil)
	svd.UTo(&u)
	svd.VTo(&v)
	u0 := make(Vec, len(x0))
	for i := range len(x0) {
		u0[i] = u.At(i, 0)
	}
	if DotVecs(u0, xaxis) < 0 {
		u0.ScaleBy(-1)
	}
	v0 := make(Vec, len(y0))
	for i := range len(y0) {
		v0[i] = v.At(i, 0)
	}
	if DotVecs(v0, yaxis) < 0 {
		v0.ScaleBy(-1)
	}
	return sv, u0, v0
}
