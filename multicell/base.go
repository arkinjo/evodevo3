package multicell

import (
	"encoding/json"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

const (
	sqrt3              = 1.732050807568877
	default_num_cell_x = 1
	default_num_cell_y = 1
	default_density    = 0.02 // genome matrix density
)

// vector
type Vec = []float64

// sparse matrix
type SpMat = [](map[int]float64)

// various set-ups
type Setting struct {
	Basename       string  // name of the model
	Outdir         string  // output directory for trajectory
	Seed           int     // random seed
	With_cue       bool    // with cue or not
	Max_pop        int     // maximum population size
	Num_cell_x     int     // number of cells in the x-axis
	Num_cell_y     int     // number of cells in the y-axis
	Num_env        int     // number of environmental factors
	Num_layers     int     // number of layers
	Max_dev        int     // maximum number of developmental steps
	Num_components []int   // number of components of each state vector
	Topology       SpMat   // densities of genome matrices
	Omega          Vec     // scaling factors of activation function inputs
	Env_noise      float64 // noise level
	Mut_rate       float64 // mutation rate
	Conv_dev       float64 // convergence limit
	Denv           float64 // size of an environmental change
	Selstrength    float64 // selection strength
	Alpha          float64 // weight for exponential moving average
	ProductionRun  bool    // true if production run (i.e. "test" phase)
}

// LeCun-inspired arctan function
func LCatan(x float64) float64 {
	return 6.0 * math.Atan(x/sqrt3) / math.Pi
}

func GetDefaultSetting(basename string, num_layers int, seed int) *Setting {
	for range seed {
		rand.IntN(2)
	}
	num_components := make([]int, num_layers)
	omega := make(Vec, num_layers)
	topology := NewSpMat(num_layers)
	ncx := default_num_cell_x
	ncy := default_num_cell_y
	num_env := (200 / 4) * (ncx + ncy) * 2

	for i := 0; i < num_layers; i++ {
		num_components[i] = 200
		if i > 0 {
			topology[i][i-1] = default_density
			if i < num_layers-1 {
				topology[i][i] = default_density
			}
		}
	}

	return &Setting{
		Basename:       basename,
		Outdir:         ".",
		Seed:           seed,
		With_cue:       true,
		Max_pop:        500,
		Num_cell_x:     ncx,
		Num_cell_y:     ncy,
		Num_env:        num_env,
		Num_layers:     num_layers,
		Max_dev:        200,
		Num_components: num_components,
		Topology:       topology,
		Omega:          omega,
		Env_noise:      0.05,
		Mut_rate:       0.001,
		Conv_dev:       1e-5,
		Denv:           0.5,
		Selstrength:    20.0,
		Alpha:          1.0 / 3.0,
		ProductionRun:  false}
}

func (s *Setting) SetOmega() {
	for l, tl := range s.Topology {
		omega := 0.0
		for k, d := range tl {
			omega += d * float64(s.Num_components[k])
		}
		if omega > 0 {
			s.Omega[l] = 1.0 / math.Sqrt(omega)
		}
	}
}

func JustFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Setting) Dump(filename string) {
	json, err := json.MarshalIndent(s, "", "    ")
	JustFail(err)
	os.WriteFile(filename, json, 0644)
}

func LoadSetting(filename string) *Setting {
	buffer, err := os.ReadFile(filename)
	JustFail(err)
	s := GetDefaultSetting("", 1, 1)
	err = json.Unmarshal(buffer, s)
	JustFail(err)

	return s
}
