package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/arkinjo/evodevo3/multicell"
	//	"log"
	"os"
	//	"strings"
	//	"math"
)

func dumpjson(js []byte) {
	fmt.Println(len(js))
	os.Stdout.Write(js)
}

func main() {
	s := multicell.GetDefaultSetting("hoge", 5, 15)
	s.SetOmega()
	s.Max_pop = 500
	s.ProductionRun = true
	s.Outdir = "traj"
	s.FullModel()
	s.SaveEnvs("envs.json", 50)
	envs := s.LoadEnvs("envs.json")
	env := envs[0]
	pop := s.NewPopulation(env)
	fmt.Println("#With_cue= ", s.With_cue)
	pop.RunEpochs(s, 1, 10, env)
	ofilename := pop.Dump(s)
	pop = s.LoadPopulation(ofilename, env)
	fmt.Println(s.Selecting_env(env))
	fmt.Println(pop.Indivs[0].Selected_pheno())
	s.Dump("mysetting.json")
	s2 := multicell.LoadSetting("mysetting.json")
	s2.Dump("mysetting2.json")
}
