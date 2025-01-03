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
	ave.ScaleBy(1.0 / n)
	return ave
}

func VarVecs(vecs []Vec, mv Vec) Vec {
	vvec := make(Vec, len(mv))
	for _, vec := range vecs {
		for i, v := range vec {
			d := v - mv[i]
			vvec[i] += d * d
		}
	}
	n := float64(len(vecs))
	vvec.ScaleBy(1.0 / n)

	return vvec
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

func (pop *Population) GetProjected1(s *Setting, fout *os.File, label string, xs []Vec, x0 Vec, axis, ps Vec) (Vec, Vec) {
	sv, u, _ := XPCA(xs, x0, xs, x0)
	ali := 0.0
	if axis != nil {
		ali = math.Abs(DotVecs(u[0], axis))
	}

	fmt.Fprintf(fout, "%s\t%d\t%f\t%f\t%f\n",
		label, pop.Igen, sv.Norm2(), sv[0], ali)

	px := ProjectOnAxis(xs, x0, u[0])
	py := ProjectOnAxis(xs, x0, u[1])

	if DotVecs(px, ps) < 0 {
		px.ScaleBy(-1)
		u[0].ScaleBy(-1)
	}

	if DotVecs(py, ps) < 0 {
		py.ScaleBy(-1)
		u[1].ScaleBy(-1)
	}

	filvec := s.TrajectoryFilename(pop.Iepoch, pop.Igen, label)
	fvec, err := os.OpenFile(filvec, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fvec.Close()
	for i, ui := range u[0] {
		fmt.Fprintf(fvec, "%d\t%e\t%e\n", i, ui, u[1][i])
	}

	return px, py
}

func (pop *Population) GetProjected2(s *Setting, fout *os.File, label string, xs []Vec, x0 Vec, ys []Vec, y0 Vec, axis, ps Vec) (Vec, Vec) {
	sv, u, v := XPCA(xs, x0, ys, y0)
	ali := 0.0
	if axis != nil {
		ali = math.Abs(DotVecs(u[0], axis))
	}
	fmt.Fprintf(fout, "%s\t%d\t%f\t%f\t%f\n",
		label, pop.Igen, sv.Norm2(), sv[0], ali)

	px := ProjectOnAxis(xs, x0, u[0])
	py := ProjectOnAxis(ys, y0, v[0])

	if DotVecs(px, ps) < 0 {
		px.ScaleBy(-1)
		py.ScaleBy(-1)
		u[0].ScaleBy(-1)
		v[0].ScaleBy(-1)
	}

	filvec := s.TrajectoryFilename(pop.Iepoch, pop.Igen, label)
	fvec, err := os.OpenFile(filvec, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fvec.Close()
	for i, ui := range u[0] {
		fmt.Fprintf(fvec, "U\t%d\t%e\n", i, ui)
	}
	for i, vi := range v[0] {
		fmt.Fprintf(fvec, "V\t%d\t%e\n", i, vi)
	}

	return px, py
}

func (pop *Population) PrintPopStats(fout *os.File, gs, ps Vec) {
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
	fmt.Fprintf(fout, "Fit\t%d\t%e\t%e\n", pop.Igen, fa, fv)
}

func (pop *Population) GenoPhenoPlot(s *Setting, p0, paxis, g0, gaxis Vec) {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen, "gpplot")
	fout, err := os.Create(filename)
	JustFail(err)
	defer fout.Close()
	// Geno-Pheno Projection Plot
	gvecs := pop.GenomeVecs(s)
	gs := ProjectOnAxis(gvecs, g0, gaxis)

	pvecs := pop.PhenoVecs(s)
	ps := ProjectOnAxis(pvecs, p0, paxis)

	pop.PrintPopStats(fout, gs, ps)

	mg := MeanVecs(gvecs)
	mp := MeanVecs(pvecs)
	gvar := MatTotVar(gvecs, mg)
	pvar := MatTotVar(pvecs, mp)
	rvar := s.RandomGenomeVariance()
	fmt.Fprintf(fout, "GPvar\t%d\t%f\t%f\t%f\n", pop.Igen, gvar, pvar, rvar)
	fmt.Fprintf(fout, "#\t%3s\t%8s\t%8s", "gen", "g", "p")
	fmt.Fprintf(fout, "\n")
	for i := range pop.Indivs {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f", i, gs[i], ps[i])
		fmt.Fprintf(fout, "\n")
	}
	log.Printf("Projection saved in: %s", filename)
}

func (pop *Population) Project(s *Setting, p0, paxis, g0, gaxis, c0, caxis Vec) {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen, "gpplot")
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()

	// for alignment calculation
	punit := paxis.Clone()
	punit.Normalize()

	// Geno-Pheno Projection Plot
	gvecs := pop.GenomeVecs(s)
	pvecs := pop.PhenoVecs(s)
	mg := MeanVecs(gvecs)
	mp := MeanVecs(pvecs)
	gs := ProjectOnAxis(gvecs, g0, paxis)
	ps := ProjectOnAxis(pvecs, p0, paxis)

	//Pheno-Pheno variance-covariance
	Pp0, Pp1 := pop.GetProjected1(s, fout, "PPcov", pvecs, mp, punit, ps)

	// Pheno-Cue Cross-Covariance
	cvecs := pop.CueVecs(s)
	mc := c0 //MeanVecs(cvecs)
	Ppc, Cpc := pop.GetProjected2(s, fout, "PCcov", pvecs, mp, cvecs, mc, punit, ps)

	// Pheno-Geno Cross-Covariance
	Ppg, Gpg := pop.GetProjected2(s, fout, "PGcov", pvecs, mp, gvecs, mg, punit, ps)

	// State variance-covariance
	var Ss0, Ss1 Vec
	if s.NumLayers > 1 {
		svecs := pop.StateVecs()
		ms := MeanVecs(svecs)
		Ss0, Ss1 = pop.GetProjected1(s, fout, "SScov", svecs, ms, nil, ps)
		// State-Cue cross-covariance
		//Ssc, Csc := pop.GetProjected2(s, fout, "SCcov", svecs, ms, cvecs, mc, nil, ps)

		// State-Genome cross-covariance (very slow)
		//	Ssg, Gsg := pop.GetProjected2(s, fout, "SGcov", svecs, ms, gvecs, mg, nil, ps)
	}

	fmt.Fprintf(fout, "#\t%3s\t%8s\t%8s", "gen", "g", "p")
	fmt.Fprintf(fout, "\t%8s\t%8s", "Ppheno0", "Ppheno1")
	fmt.Fprintf(fout, "\t%8s\t%8s", "Ccue", "Pcue")
	fmt.Fprintf(fout, "\t%8s\t%8s", "Ggeno", "Pgeno")
	if s.NumLayers > 1 {
		fmt.Fprintf(fout, "\t%8s\t%8s", "SS0", "SS1")
		//	fmt.Fprintf(fout, "\t%8s\t%8s", "Gsc", "Ssc")
		//	fmt.Fprintf(fout, "\t%8s\t%8s", "Gsg", "Ssg")
	}
	fmt.Fprintf(fout, "\n")
	for i := range pop.Indivs {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f", i, gs[i], ps[i])
		fmt.Fprintf(fout, "\t%f\t%f", Pp0[i], Pp1[i])

		fmt.Fprintf(fout, "\t%f\t%f", Cpc[i], Ppc[i])
		fmt.Fprintf(fout, "\t%f\t%f", Gpg[i], Ppg[i])
		if s.NumLayers > 1 {
			fmt.Fprintf(fout, "\t%f\t%f", Ss0[i], Ss1[i])
			//	fmt.Fprintf(fout, "\t%f\t%f", Csc[i], Ssc[i])
			//	fmt.Fprintf(fout, "\t%f\t%f", Gsg[i], Ssg[i])
		}
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
		for j, yj := range yt {
			for i, xi := range xt {
				cov.Set(i, j, cov.At(i, j)+xi*yj)
			}
		}
	}

	fac := 1.0 / float64(len(xs))
	cov.Scale(fac, cov)

	return cov
}

// Get singular values, the first left and right singular vectors.
func XPCA(xs []Vec, x0 Vec, ys []Vec, y0 Vec) (Vec, []Vec, []Vec) {
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
	u0 := make([]Vec, 2)
	for j := range len(u0) {
		u0[j] = make(Vec, len(x0))
		for i := range len(x0) {
			u0[j][i] = u.At(i, j)
		}
	}

	v0 := make([]Vec, 2)
	for j := range len(v0) {
		v0[j] = make(Vec, len(y0))
		for i := range len(y0) {
			v0[j][i] = v.At(i, j)
		}
	}

	return sv, u0, v0
}
