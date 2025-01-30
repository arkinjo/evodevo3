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

func (s *Setting) GetPhenoAxis(pop0, pop1 Population, env0, env1 Environment) (Vec, Vec) {

	pop0.Initialize(s, env0)
	pop0.Develop(s, env0)
	p0 := MeanVecs(pop0.PhenoVecs(s))
	p1 := MeanVecs(pop1.PhenoVecs(s))
	return p0, GetAxis(p0, p1)
}

func (s *Setting) GetSelectedPhenoAxis(pop0, pop1 Population, env0, env1 Environment) (Vec, Vec) {

	pop0.Initialize(s, env0)
	pop0.Develop(s, env0)
	p0 := MeanVecs(pop0.SelectedPhenoVecs(s))
	p1 := MeanVecs(pop1.SelectedPhenoVecs(s))
	//	fmt.Println("p0: ", p0)
	//	fmt.Println("p1: ", p1)
	return p0, GetAxis(p0, p1)
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
	ali1 := 0.0
	ali2 := 0.0
	if axis != nil {
		ali1 = math.Abs(DotVecs(u[0], axis))
		ali2 = math.Abs(DotVecs(u[1], axis))
	}

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

	corr0, pval0 := CorrVecs(px, ps)
	corr1, pval1 := CorrVecs(py, ps)
	fmt.Fprintf(fout, "%s\t%d\t%f\t%f\t%f\t%f\t%f\t%e\t%f\t%e\n",
		label, pop.Igen, sv[0], sv[0]/sv.Norm2(), ali1, ali2,
		corr0, pval0, corr1, pval1)

	filvec := s.TrajectoryFilename(pop.Iepoch, pop.Igen, label)
	fvec, err := os.OpenFile(filvec, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fvec.Close()

	for i, ui := range u[0] {
		fmt.Fprintf(fvec, "U\t%d\t%e\t%e\n", i, ui, u[1][i])
	}

	return px, py
}

func (pop *Population) GetProjected2(s *Setting, fout *os.File, label string, xs []Vec, x0 Vec, ys []Vec, y0 Vec, uaxis, vaxis, ps, gs Vec) (Vec, Vec) {
	sv, u, v := XPCA(xs, x0, ys, y0)
	uali := 0.0
	if uaxis != nil {
		uali = DotVecs(u[0], uaxis)
	}
	vali := 0.0
	if vaxis != nil {
		vali = DotVecs(v[0], vaxis)
	}

	px := ProjectOnAxis(xs, x0, u[0])
	py := ProjectOnAxis(ys, y0, v[0])

	if uali < 0 {
		uali *= -1
		vali *= -1
		px.ScaleBy(-1)
		py.ScaleBy(-1)
		u[0].ScaleBy(-1)
		v[0].ScaleBy(-1)
	}
	corr, pval := CorrVecs(px, py)
	corrp, pvalp := CorrVecs(ps, px)
	corrg, pvalg := CorrVecs(gs, py)
	svtot := sv.Norm2()
	fmt.Fprintf(fout, "%s\t%d\t%f\t%f\t%f\t%f\t%f\t%e\t%f\t%e\t%f\t%e\n",
		label, pop.Igen, svtot, sv[0]/svtot, uali, vali,
		corrp, pvalp, corrg, pvalg, corr, pval)

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

func (pop *Population) PrintPopStats(fout *os.File, gs, ps, ali0 Vec) {
	ga, gv := avesd(gs)
	pa, pv := avesd(ps)
	fmt.Fprintf(fout, "Proj\t%d\t%d\t%f\t%f\t%f\t%f\n",
		pop.Igen, len(gs), ga, pa, gv, pv)

	var ns, ms, fs Vec
	for _, indiv := range pop.Indivs {
		ns = append(ns, float64(indiv.Ndev))
		ms = append(ms, indiv.Align)
		fs = append(fs, indiv.Fitness)
	}
	na, nv := avesd(ns)
	ma, mv := avesd(ms)
	a0a, a0v := avesd(ali0)
	fa, fv := avesd(fs)
	fmt.Fprintf(fout, "Ndev\t%d\t%f\t%f\n", pop.Igen, na, nv)
	fmt.Fprintf(fout, "AliNov\t%d\t%f\t%f\n", pop.Igen, ma, mv)
	fmt.Fprintf(fout, "AliAnc\t%d\t%f\t%f\n", pop.Igen, a0a, a0v)
	fmt.Fprintf(fout, "Fit\t%d\t%e\t%e\n", pop.Igen, fa, fv)
}

func (pop *Population) GenoPhenoPlot(s *Setting, p0, paxis, g0, gaxis, env0 Vec) {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen, "gpplot")
	fout, err := os.Create(filename)
	JustFail(err)
	defer fout.Close()

	selenv0 := env0.SelectingEnv(s)
	// Geno-Pheno Projection Plot
	gvecs := pop.GenomeVecs(s)
	gs := ProjectOnAxis(gvecs, g0, gaxis)
	pvecs := pop.SelectedPhenoVecs(s)
	ps := ProjectOnAxis(pvecs, p0, paxis)

	ali0 := make(Vec, len(pvecs))
	for i, p := range pvecs {
		ali0[i] = DotVecs(selenv0, p) / float64(len(selenv0))
	}

	pop.PrintPopStats(fout, gs, ps, ali0)

	mg := MeanVecs(gvecs)
	mp := MeanVecs(pvecs)
	gvar := MatTotVar(gvecs, mg)
	pvar := MatTotVar(pvecs, mp)
	rvar := s.RandomGenomeVariance()
	ali := make(Vec, len(pop.Indivs))
	for i, indiv := range pop.Indivs {
		ali[i] = indiv.Align
	}
	pacor, pvala := CorrVecs(ps, ali)
	fmt.Fprintf(fout, "PAcorr\t%d\t%f\t%e\n", pop.Igen, pacor, pvala)
	fmt.Fprintf(fout, "GPvar\t%d\t%f\t%f\t%f\n", pop.Igen, gvar, pvar, rvar)
	fmt.Fprintf(fout, "#\t%3s\t%8s\t%8s\t%8s\t%8s", "id", "AliNov", "AliAnc", "g", "p")
	fmt.Fprintf(fout, "\n")
	for i, indiv := range pop.Indivs {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f\t%f\t%f",
			i, indiv.Align, ali0[i], gs[i], ps[i])
		fmt.Fprintf(fout, "\n")
	}
	log.Printf("Projection saved in: %s", filename)
}

func (pop *Population) PGCov(s *Setting, p0, paxis, g0, gaxis, env0, env1 Vec) {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen, "pgcov")
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()

	selenv0 := env0.SelectingEnv(s)
	selenv0.Normalize()
	selenv1 := env1.SelectingEnv(s)
	selenv1.Normalize()

	// for alignment calculation
	punit := paxis.Clone()
	punit.Normalize()
	gunit := gaxis.Clone()
	gunit.Normalize()

	gvecs := pop.GenomeVecs(s)
	pvecs := pop.SelectedPhenoVecs(s)
	mg := MeanVecs(gvecs)
	mp := MeanVecs(pvecs)

	// Geno-Pheno Projection Plot
	gs := ProjectOnAxis(gvecs, g0, gaxis)
	ps := ProjectOnAxis(pvecs, p0, paxis)

	sv, u, v := XPCA(pvecs, mp, gvecs, mg)
	svtot := sv.Norm2()
	var pks, gks []Vec
	fmt.Fprintf(fout, "Tot\t%f\n", svtot)
	for k, sk := range sv {
		acc := sv[:k+1].Norm2() / svtot
		fmt.Fprintf(fout, "SV\t%d\t%f\t%f\t%f\n", k, sk, sk/svtot, acc)
	}
	for k := range u {
		pk := ProjectOnAxis(pvecs, mp, u[k])
		gk := ProjectOnAxis(gvecs, mg, v[k])
		if DotVecs(pk, ps) < 0 {
			pk.ScaleBy(-1)
			gk.ScaleBy(-1)
			u[k].ScaleBy(-1)
			v[k].ScaleBy(-1)
		}
		pks = append(pks, pk)
		gks = append(gks, gk)

		ali0 := DotVecs(selenv0, u[k])
		ali1 := DotVecs(selenv1, u[k])
		fmt.Fprintf(fout, "AliNov\t%d\t%f\n", k, ali1)
		fmt.Fprintf(fout, "AliAnc\t%d\t%f\n", k, ali0)
		corr, pval := CorrVecs(pk, gk)
		corrp, pvalp := CorrVecs(pk, ps)
		corrg, pvalg := CorrVecs(gk, gs)
		fmt.Fprintf(fout, "PGcorr\t%d\t%f\t%e\n", k, corr, pval)
		fmt.Fprintf(fout, "PUcorr\t%d\t%f\t%e\n", k, corrp, pvalp)
		fmt.Fprintf(fout, "GVcorr\t%d\t%f\t%e\n", k, corrg, pvalg)
	}

	for i, gi := range gs {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f", i, gi, ps[i])
		for k, gk := range gks {
			fmt.Fprintf(fout, "\t%f\t%f", gk[i], pks[k][i])
		}
		fmt.Fprintf(fout, "\n")
	}
}

func (pop *Population) SVDProject(s *Setting, p0, paxis, g0, gaxis, c0, caxis Vec) {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen, "xpca")
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()

	// for alignment calculation
	punit := paxis.Clone()
	punit.Normalize()
	gunit := gaxis.Clone()
	gunit.Normalize()
	cunit := caxis.Clone()
	cunit.Normalize()

	// Geno-Pheno Projection Plot
	gvecs := pop.GenomeVecs(s)
	pvecs := pop.PhenoVecs(s)

	mg := MeanVecs(gvecs)
	mp := MeanVecs(pvecs)
	gs := ProjectOnAxis(gvecs, g0, gaxis)
	ps := ProjectOnAxis(pvecs, p0, paxis)

	//Pheno-Pheno variance-covariance
	Pp0, Pp1 := pop.GetProjected1(s, fout, "PPcov", pvecs, mp, punit, ps)

	// Pheno-Cue Cross-Covariance
	cvecs := pop.CueVecs(s)
	mc := c0 //MeanVecs(cvecs)
	cs := ProjectOnAxis(cvecs, c0, caxis)
	Ppc, Cpc := pop.GetProjected2(s, fout, "PCcov", pvecs, mp, cvecs, mc, punit, cunit, ps, cs)

	// Pheno-Geno Cross-Covariance
	Ppg, Gpg := pop.GetProjected2(s, fout, "PGcov", pvecs, mp, gvecs, mg, punit, gunit, ps, gs)

	// State variance-covariance
	var Ss0, Ss1 Vec
	if s.NumLayers > 2 {
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
	if s.NumLayers > 2 {
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
		if s.NumLayers > 2 {
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
	npca := 3
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
	if npca > len(sv) {
		npca = len(sv)
	}
	u0 := make([]Vec, npca)
	for j := range len(u0) {
		u0[j] = make(Vec, len(x0))
		for i := range len(x0) {
			u0[j][i] = u.At(i, j)
		}
	}

	v0 := make([]Vec, npca)
	for j := range len(v0) {
		v0[j] = make(Vec, len(y0))
		for i := range len(y0) {
			v0[j][i] = v.At(i, j)
		}
	}

	return sv, u0, v0
}

// Analyze adaptive plastic responses to various environmental changes.
func (pop *Population) AnalyzeVarEnvs(s *Setting, env0 Environment, n int) {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen, "varenv")
	log.Printf("AnalyzeVarEnvs output to %s\n", filename)
	fout, err := os.Create(filename)
	JustFail(err)
	defer fout.Close()

	gvecs0 := pop.GenomeVecs(s)
	mg0 := MeanVecs(gvecs0)
	pvecs0 := pop.PhenoVecs(s)
	mp0 := MeanVecs(pvecs0)
	sv0, u0, v0 := XPCA(pvecs0, mp0, gvecs0, mg0)
	fmt.Fprintf(fout, "SV\t%d\t%d\t%e\t%e\n", pop.Iepoch, 0, sv0[0], sv0[0]/sv0.Norm2())

	ps0 := ProjectOnAxis(pvecs0, mp0, u0[0])
	gs0 := ProjectOnAxis(gvecs0, mg0, v0[0])

	var us, vs, envs []Vec
	envs = append(envs, env0)
	us = append(us, u0[0])
	vs = append(vs, v0[0])

	var pss0, pss, gss []Vec
	pss0 = append(pss0, ps0)
	pss = append(pss, ps0)
	gss = append(gss, gs0)
	for i := range n {
		log.Printf("Environment %d\n", i+1)
		env := env0.ChangeEnv(s)

		envs = append(envs, env)
		pop.Initialize(s, env)
		pop.Develop(s, env)
		pvecs := pop.PhenoVecs(s)
		mp := MeanVecs(pvecs)
		sv, u, v := XPCA(pvecs, mp, gvecs0, mg0)

		// Evolve for 200 generations
		pop1, _ := pop.Evolve(s, env)
		pvecs1 := pop1.PhenoVecs(s)
		mp1 := MeanVecs(pvecs1)
		dp := make(Vec, len(mp1))
		dp.Diff(mp1, mp0)
		dp.Normalize()

		gvecs1 := pop1.GenomeVecs(s)
		mg1 := MeanVecs(gvecs1)
		dg := make(Vec, len(mg0))
		dg.Diff(mg1, mg0)
		dg.Normalize()

		pali := DotVecs(dp, u[0])
		gali := DotVecs(dg, v[0])
		if pali < 0 {
			pali *= -1
			gali *= -1
			u[0].ScaleBy(-1)
			v[0].ScaleBy(-1)
		}
		us = append(us, u[0])
		vs = append(vs, v[0])
		psOn0 := ProjectOnAxis(pvecs, mp0, u0[0])
		pss0 = append(pss0, psOn0)

		ps := ProjectOnAxis(pvecs, mp, u[0])
		pss = append(pss, ps)

		gs := ProjectOnAxis(gvecs0, mg0, v[0])
		gss = append(gss, gs)

		ts := sv.Norm2()
		fmt.Fprintf(fout, "SV\t%d\t%d\t%e\t%e\t%e\t%e\n", pop.Iepoch, i+1, sv[0], sv[0]/ts, pali, gali)
	}
	for i, ps := range pss {
		cpg, ppg := CorrVecs(ps, gss[i])
		fmt.Fprintf(fout, "CorrPG\t%d\t%d\t%e\t%e\n", pop.Iepoch, i, cpg, ppg)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n+1; j++ {
			fmt.Fprintf(fout, "CorrE\t%d\t%d", i, j)

			corre, pe := CorrVecs(envs[i], envs[j])
			fmt.Fprintf(fout, "\t%f\t%e\n", corre, pe)
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n+1; j++ {
			fmt.Fprintf(fout, "CorrU\t%d\t%d", i, j)
			corru, pu := CorrVecs(us[i], us[j])
			dot := DotVecs(us[i], us[j])
			fmt.Fprintf(fout, "\t%f\t%e\t%f\n", corru, pu, dot)
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n+1; j++ {
			fmt.Fprintf(fout, "CorrV\t%d\t%d", i, j)
			corrv, pv := CorrVecs(vs[i], vs[j])
			dot := DotVecs(vs[i], vs[j])
			fmt.Fprintf(fout, "\t%f\t%e\t%f\n", corrv, pv, dot)
		}
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n+1; j++ {
			fmt.Fprintf(fout, "CorrP0\t%d\t%d", i, j)
			corrv, pv := CorrVecs(pss0[i], pss0[j])
			fmt.Fprintf(fout, "\t%f\t%e\n", corrv, pv)
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n+1; j++ {
			fmt.Fprintf(fout, "CorrG\t%d\t%d", i, j)
			corrv, pv := CorrVecs(gss[i], gss[j])
			fmt.Fprintf(fout, "\t%f\t%e\n", corrv, pv)
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n+1; j++ {
			fmt.Fprintf(fout, "CorrPS\t%d\t%d", i, j)
			corrv, pv := CorrVecs(pss[i], pss[j])
			fmt.Fprintf(fout, "\t%f\t%e\n", corrv, pv)
		}
	}

	for i := range len(env0) {
		fmt.Fprintf(fout, "Env\t%d", i)
		for k := range envs {
			fmt.Fprintf(fout, "\t%f", envs[k][i])
		}
		fmt.Fprintf(fout, "\n")
	}

	fmt.Fprintf(fout, "#Ip0\tepoch\tind\t%8s\t%8s\tp1...\n",
		"v0(g)", "u0(p)")
	for i := range len(gs0) {
		fmt.Fprintf(fout, "Ip0\t%d\t%d\t%e\t%e", pop.Iepoch, i,
			gs0[i], ps0[i])
		for k := range pss0 {
			fmt.Fprintf(fout, "\t%e", pss0[k][i])
		}
		fmt.Fprintf(fout, "\n")
	}

	fmt.Fprintf(fout, "#Igp\tepoch\tind\t%8s\t%8s...\n", "v(g)", "u(p)")
	for i := range len(gs0) {
		fmt.Fprintf(fout, "Igp\t%d\t%d", pop.Iepoch, i)
		for k := range gss {
			fmt.Fprintf(fout, "\t%e\t%e", gss[k][i], pss[k][i])
		}
		fmt.Fprintf(fout, "\n")
	}

	fmt.Fprintf(fout, "#U\tepoch\tind\t%8s...\n", "u1")
	for i := range len(us[0]) {
		fmt.Fprintf(fout, "U\t%d\t%d", pop.Iepoch, i)
		for k := range us {
			fmt.Fprintf(fout, "\t%e", us[k][i])
		}
		fmt.Fprintf(fout, "\n")
	}

	fmt.Fprintf(fout, "#V\tepoch\tind\t%8s...\n", "v1")
	for i := range len(vs[0]) {
		fmt.Fprintf(fout, "V\t%d\t%d", pop.Iepoch, i)
		for k := range vs {
			fmt.Fprintf(fout, "\t%e", vs[k][i])
		}
		fmt.Fprintf(fout, "\n")
	}

}

func ConservedGenomeSites(mg1, vg1 Vec, gvecs []Vec) (map[int]int, []int) {
	conserved := make(map[int]int)
	for k, g := range mg1 {
		if g != 0.0 && vg1[k] == 0.0 {
			conserved[k] += 0
		}
	}
	count := make([]int, len(gvecs))
	for i, g := range gvecs {
		for k := range conserved {
			if math.Abs(mg1[k]-g[k]) < 1e-10 {
				conserved[k] += 1
				count[i] += 1
			}
		}
	}
	return conserved, count
}

// Comparing adaptive plastic response in env0 to evolutionary adaptation to env1
func (s *Setting) AnalyzeAPRGeno(env0, env1 Environment, pop0, pop1 Population) {
	// Generation 1 in novel environment.
	gvecs0 := pop0.GenomeVecs(s)
	mg0 := MeanVecs(gvecs0)
	vg0 := VarVecs(gvecs0, mg0)
	pvecs0N := pop0.PhenoVecs(s)
	var denv Vec
	denv.Diff(env1, env0)
	p0, paxis := s.GetPhenoAxis(pop0, pop1, env0, env1)
	punit := paxis.Clone().Normalize()

	pproj0 := ProjectOnAxis(pvecs0N, p0, paxis)

	// develop in ancestral environment.
	pop0.Initialize(s, env0)
	pop0.Develop(s, env0)
	pvecs0A := pop0.PhenoVecs(s)

	dpvecs0 := DiffMats(pvecs0N, pvecs0A)
	mp0 := MeanVecs(dpvecs0)
	sv, us, vs := XPCA(dpvecs0, mp0, gvecs0, mg0)

	// Generation 200(?) adapted to novel environment.
	gvecs1 := pop1.GenomeVecs(s)
	mg1 := MeanVecs(gvecs1)
	vg1 := VarVecs(gvecs1, mg1)
	dg1 := make(Vec, len(mg1))
	dg1.Diff(mg1, mg0)
	gaxis := GetAxis(mg0, mg1)
	gunit := gaxis.Clone().Normalize()
	gproj0 := ProjectOnAxis(gvecs0, mg0, gaxis)

	conserved1, count1 := ConservedGenomeSites(mg1, vg1, gvecs0)
	conserved0, count0 := ConservedGenomeSites(mg0, vg0, gvecs1)
	shared := make(map[int]bool)
	for k := range conserved0 {
		if _, ok := conserved1[k]; ok {
			shared[k] = true
		}
	}

	if DotVecs(punit, us[0]) < 0 {
		vs[0].ScaleBy(-1)
		us[0].ScaleBy(-1)
	}
	vproj0 := ProjectOnAxis(gvecs0, mg0, vs[0])
	uproj0 := ProjectOnAxis(dpvecs0, mp0, us[0])

	filename := s.TrajectoryFilename(pop0.Iepoch, pop0.Igen, "aprgeno")
	log.Printf("AnalyzeAPRGeno output to %s\n", filename)
	fout, err := os.Create(filename)
	JustFail(err)
	defer fout.Close()

	fmt.Fprintf(fout, "SV\t%e\t%e\t%e\t%e\n",
		sv[0], sv[0]/sv.Norm2(),
		math.Abs(DotVecs(punit, us[0])),
		math.Abs(DotVecs(gunit, vs[0])))
	fmt.Fprintf(fout, "Cons\t%d\t%d\t%d\n", len(conserved0), len(conserved1), len(shared))
	fmt.Fprintf(fout, "#\tind\t%8s\t%8s\t%8s\t%8s\t%4s\t%4s\n",
		"Gproj", "Pproj", "Vproj", "Uproj", "Cons0", "Cons1")
	for i, pp := range pproj0 {
		fmt.Fprintf(fout, "I\t%d\t%f\t%f", i, gproj0[i], pp)
		fmt.Fprintf(fout, "\t%f\t%f", vproj0[i], uproj0[i])
		fmt.Fprintf(fout, "\t%d\t%d\n", count0[i], count1[i])
	}
	fmt.Fprintf(fout, "#\tind\t%8s\t%8s\n", "U1", "denv")
	for i, u := range us[0] {
		fmt.Fprintf(fout, "P\t%d\t%e\t%e\n", i, u, denv[i])
	}
	for i, v := range vs[0] {
		fmt.Fprintf(fout, "G\t%d\t%e\t%e\t%d\t%d",
			i, v, dg1[i], conserved0[i], conserved1[i])
		fmt.Fprintf(fout, "\t%e\t%e", mg0[i], vg0[i])
		fmt.Fprintf(fout, "\t%e\t%e", mg1[i], vg1[i])
		fmt.Fprintf(fout, "\n")
	}
}
