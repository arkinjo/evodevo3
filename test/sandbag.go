package main

import (
	"encoding/json"
	"fmt"
	"github.com/arkinjo/evodevo3/multicell"
	"log"
	"os"
	//	"strings"
	//	"math"
)

func dumpjson(js []byte) {
	fmt.Println(len(js))
	os.Stdout.Write(js)
}

func main() {
	s := multicell.Get_default_setting("hoge", 5, 15)
	s.Set_omega()
	s.Max_pop = 100

	env := s.NewEnvironment()
	pop := s.NewPopulation(env)

	js, err := json.Marshal(pop)
	if err != nil {
		log.Fatal(err)
	}
	dumpjson(js)

	pop.Evolve(s, env, 200)
	fmt.Println(s.Selecting_env(env))
	fmt.Println(pop.Indivs[0].Selected_pheno())

}
