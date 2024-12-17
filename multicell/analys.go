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

// Total variance of a matrix
func MatTotVar(vecs []Vec, mv Vec) float64 {
	sd2 := 0.0
	vt := make(Vec, len(mv))
	for _, vec := range vecs {
		vt.Diff(vec, mv)
		sd2 += DotVecs(vt, vt)
	}
	return sd2 / float64(len(vecs))
}

func (pop *Population) Project(s *Setting, p0, paxis, g0, gaxis, c0, caxis Vec) {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen, "gpplot")
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()

	// for alignment calculation
	punit := paxis.Copy()
	punit.Normalize()

	// Geno-Pheno Plot
	gvecs := pop.GenomeVecs(s)
	gs := ProjectOnAxis(gvecs, g0, gaxis)

	pvecs := pop.PhenoVecs(s)
	ps := ProjectOnAxis(pvecs, p0, paxis)
	ga, gv := avesd(gs)
	pa, pv := avesd(ps)
	fmt.Fprintf(fout, "Proj\t%d\t%d\t%f\t%f\t%f\t%f\n",
		pop.Igen, len(gs), ga, pa, gv, pv)

	var ns, ms, fs Vec
	for _, indiv := range pop.Indivs {
		ns = append(ns, float64(indiv.Ndev))
		ms = append(ms, indiv.Mismatch)
		fs = append(fs, indiv.Fitness)
	}
	na, nv := avesd(ns)
	ma, mv := avesd(ms)
	fa, fv := avesd(fs)
	fmt.Fprintf(fout, "Ndev\t%d\t%f\t%f\n", pop.Igen, na, nv)
	fmt.Fprintf(fout, "Mis\t%d\t%f\t%f\n", pop.Igen, ma, mv)
	fmt.Fprintf(fout, "Fit\t%d\t%f\t%f\n", pop.Igen, fa, fv)

	mg := MeanVecs(gvecs)
	mp := MeanVecs(pvecs)

	gvar := MatTotVar(gvecs, mg)
	pvar := MatTotVar(pvecs, mp)
	fmt.Fprintf(fout, "GPvar\t%d\t%f\t%f\n", pop.Igen, gvar, pvar)

	//Pheno-Pheno variance-covariance
	svPheno, uPheno, _ := XPCA(pvecs, mp, paxis, pvecs, mp, paxis)
	rsvPheno0 := svPheno[0] / svPheno.Norm2()
	aliP := DotVecs(uPheno, punit)
	fmt.Fprintf(fout, "PPcov\t%d\t%f\t%f\t%f\n",
		pop.Igen, svPheno[0], rsvPheno0, aliP)

	Pps := ProjectOnAxis(pvecs, mp, uPheno)

	// Pheno-Geno Cross-Covariance
	svGeno, uGeno, vGeno := XPCA(pvecs, mp, paxis, gvecs, mg, gaxis)
	rsvGeno0 := svGeno[0] / svGeno.Norm2()
	aliG := DotVecs(uGeno, punit)
	fmt.Fprintf(fout, "PGcov\t%d\t%f\t%f\t%f\n",
		pop.Igen, svGeno[0], rsvGeno0, aliG)

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
	fmt.Fprintf(fout, "PCcov\t%d\t%f\t%f\t%f\n",
		pop.Igen, svCue[0], rsvCue0, aliC)

	Pcs := ProjectOnAxis(pvecs, mp, uCue)
	Ccs := ProjectOnAxis(cvecs, mc, vCue)
	if DotVecs(Ccs, Pcs) < 0 {
		Ccs.ScaleBy(-1)
	}

	fmt.Fprintf(fout, "#\t%3s\t%8s\t%8s", "gen", "g", "p")
	fmt.Fprintf(fout, "\t%8s", "Ppheno")
	fmt.Fprintf(fout, "\t%8s\t%8s", "Ggeno", "Pgeno")
	fmt.Fprintf(fout, "\t%8s\t%8s", "Gcue", "Pcue")
	fmt.Fprintf(fout, "\n")
	for i := range pop.Indivs {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f", i, gs[i], ps[i])
		fmt.Fprintf(fout, "\t%f", Pps[i])
		fmt.Fprintf(fout, "\t%f\t%f", Ggs[i], Pgs[i])
		fmt.Fprintf(fout, "\t%f\t%f", Ccs[i], Pcs[i])
		fmt.Fprintf(fout, "\n")
	}
	log.Printf("Projection saved in: %s", filename)
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
