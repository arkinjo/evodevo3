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
	nenvsP := flag.Int("num_envs", 50, "number of environments")
	envsfileP := flag.String("envs", "", "input environments JSON file")
	outfileP := flag.String("o", "", "Output environment JSON file")
	denvP := flag.Float64("denv", 0.5, "degree of environmental change")
	appendP := flag.Bool("append", false, "append new environments after existing ones")
	seedP := flag.Uint64("seed", 13, "random seed for environments")
	flag.Parse()

	s := multicell.GetDefaultSetting("Full")

	s.Seed = *seedP
	s.Denv = *denvP

	log.Printf("Seed=%d; Denv=%f\n", s.Seed, s.Denv)

	if *outfileP == "" {
		panic("Specify output file!")
	}
	if *appendP {
		if *envsfileP == "" {
			panic("Provide environment file with -envs!")
		}
		envs := s.LoadEnvs(*envsfileP)
		env0 := envs[len(envs)-1]
		aenvs := env0.GenerateEnvs(s, *nenvsP)
		s.DumpEnvs(*outfileP, aenvs)
	} else {
		env0 := s.NewEnvironment()
		envs := env0.GenerateEnvs(s, *nenvsP)
		s.DumpEnvs(*outfileP, envs)
	}
	log.Printf("Brand new environments saved in: %s\n", *outfileP)

}
