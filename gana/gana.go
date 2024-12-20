package main

// Genome Analysis

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Files   []string // trajectory files
}

func GetSetting() Simulation {
	settingP := flag.String("setting", "", "saved settings file")
	flag.Parse()

	if *settingP == "" {
		log.Fatal("specify a settings file with -setting")
	}
	s := multicell.LoadSetting(*settingP)
	s.ProductionRun = true

	return Simulation{
		Setting: s,
		Files:   flag.Args()}

}

func PrintMeanVarVecs(filename string, mv, vv multicell.Vec) {
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	multicell.JustFail(err)
	defer fout.Close()

	for i, m := range mv {
		fmt.Fprintf(fout, "%d\t%f\t%f\n", i, m, vv[i])
	}
}

func PrintVarMatrix(filename string, vv multicell.Vec) {
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	multicell.JustFail(err)
	defer fout.Close()

	for i, v := range vv {
		fmt.Fprintf(fout, "%f", v*100)
		if (i+1)%200 == 0 {
			fmt.Fprintf(fout, "\n")
		} else {
			fmt.Fprintf(fout, " ")
		}
	}
	fmt.Fprintf(fout, "\n")
}

func main() {
	t0 := time.Now()
	sim := GetSetting()

	for _, traj := range sim.Files {
		pop := sim.Setting.LoadPopulation(traj)
		fname := sim.Setting.TrajectoryFilename(pop.Iepoch, pop.Igen, "genomvar.tsv")
		gvecs := pop.GenomeVecs(sim.Setting)
		gmean := multicell.MeanVecs(gvecs)
		gvar := multicell.VarVecs(gvecs, gmean)
		PrintMeanVarVecs(fname, gmean, gvar)
		//		PrintVarMatrix(fname, gvar)
	}

	log.Println("Time: ", time.Since(t0))
}
