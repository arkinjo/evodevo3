package main

import (
	"flag"
	//	"fmt"
	"log"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Envs    []multicell.Environment
}

func main() {
	nenvsP := flag.Int("n", 50, "number of environments")
	envsfileP := flag.String("envs", "", "input environments JSON file")
	outfileP := flag.String("o", "", "Output environment JSON file")
	denvP := flag.Float64("denv", 0.5, "degree of environmental change")
	replaceP := flag.Int("replace", 0, "replace new environments after the epoch. should be >0.")
	seedP := flag.Uint64("seed", 13, "random seed for environments")
	flag.Parse()

	s := multicell.GetDefaultSetting("Full")

	s.Seed = *seedP
	s.Denv = *denvP

	log.Printf("Seed=%d; Denv=%f\n", s.Seed, s.Denv)

	var env0 multicell.Environment
	var envs multicell.EnvironmentS
	if *outfileP == "" {
		panic("Specify output file!")
	}
	if *replaceP <= 0 {
		env0 = s.NewEnvironment()
		envs = env0.GenerateEnvs(s, *nenvsP)
	} else {
		if *envsfileP == "" {
			panic("Provide environment file with -envs!")
		}
		aenvs := s.LoadEnvs(*envsfileP)
		env0 = aenvs[*replaceP-1]
		nenvs := env0.GenerateEnvs(s, *nenvsP+1)
		envs = append(aenvs[0:*replaceP], nenvs[1:]...)
	}

	envs.DumpEnvs(*outfileP)
	log.Printf("Brand new environments saved in: %s\n", *outfileP)

}
