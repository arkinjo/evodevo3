package multicell

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
)

const (
	Sqrt3              = 1.7320508075688772
	default_num_layers = 4
	default_len_state  = 200 // multiple of 4.
	default_num_cell_x = 1
	default_num_cell_y = 1
	default_len_face   = 50
	default_density    = 0.02 // genome matrix density
)

// faces
const (
	Left = iota
	Top
	Right
	Bottom
	NumFaces // Numbef of faces per cell
)

type Topology_t []map[int]float64

func NewTopology(n int) Topology_t {
	t := make([]map[int]float64, n)
	for i := range n {
		t[i] = make(map[int]float64)
	}
	return t
}

func (top Topology_t) Do(f func(l, k int, v float64)) {
	for l, tl := range top {
		for k, v := range tl {
			f(l, k, v)
		}
	}
}

// various set-ups
type Setting struct {
	Basename      string  // name of the model
	Seed          uint64  // random seed
	Outdir        string  // output directory for trajectory
	MaxPopulation int     // maximum population size
	MaxGeneration int     // maximum number of generations per epoch
	NumCellX      int     // number of cells in the x-axis
	NumCellY      int     // number of cells in the y-axis
	LenFace       int     // face length
	ProductionRun bool    // true if production run (i.e. "test" phase)
	EnvNoise      float64 // noise level
	MutRate       float64 // mutation rate
	ConvDevelop   float64 // convergence limit
	Denv          float64 // size of an environmental change
	SelStrength   float64 // selection strength

	WithCue    bool       // with cue or not
	MaxDevelop int        // maximum number of developmental steps
	Alpha      float64    // weight for exponential moving average
	NumLayers  int        // number of middle layers
	LenLayer   []int      // Length of each state vector
	DensityEM  float64    // input -> middle layer genome density
	Topology   Topology_t // densities of genome matrices
	Omega      Vec        // initial scaling factors of activation functions
}

// LeCun-inspired arctan function
func LCatan(omega float64) func(float64) float64 {
	b := omega / Sqrt3
	return func(x float64) float64 {
		return 6.0 * math.Atan(b*x) / math.Pi
	}
}

// LeCun tanh (without multiplying by Sqrt3).
func LCtanh(omega float64) func(float64) float64 {
	b := math.Atanh(1.0/Sqrt3) * omega
	return func(x float64) float64 {
		return math.Tanh(b * x)
	}
}

// Scaled tanh
func SCtanh(omega float64) func(float64) float64 {
	b := math.Atanh(0.99) * omega
	return func(x float64) float64 {
		return math.Tanh(b * x)
	}
}

// Simple arctan
func Atan(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		return math.Atan(omega * x)
	}
}

func Tanh(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		return math.Tanh(omega * x)
	}
}

func GetDefaultSetting() *Setting {
	s := Setting{
		Seed:          13,
		Outdir:        ".",
		MaxPopulation: 500,
		MaxGeneration: 200,
		NumCellX:      default_num_cell_x,
		NumCellY:      default_num_cell_y,
		LenFace:       default_len_face,
		ProductionRun: false,
		EnvNoise:      0.05,
		MutRate:       0.001,
		ConvDevelop:   1e-5,
		Denv:          0.5,
		SelStrength:   20.0,
		// undefined parameters are:
		//WithCue
		//MaxDevelop
		//Alpha
		//NumLayers
		//LenLayer
		//DensityEM
		//Topology
		//Omega
	}

	s.SetModel("Full")
	return &s
}

func (s *Setting) SetOmega() {
	s.Omega = make(Vec, s.NumLayers)

	s.Omega[0] = s.DensityEM * float64(s.LenFace*NumFaces)

	s.Topology.Do(func(l, k int, density float64) {
		s.Omega[l] += density * float64(s.LenLayer[k])
	})

	for l, omega := range s.Omega {
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

func (s *Setting) Dump() {
	filename := fmt.Sprintf("%s/Setting_%s.json", s.Outdir, s.Basename)
	log.Printf("Setting file saved in: %s\n", filename)
	json, err := json.MarshalIndent(s, "", "    ")
	JustFail(err)
	os.WriteFile(filename, json, 0644)
}

func LoadSetting(filename string) *Setting {
	buffer, err := os.ReadFile(filename)
	JustFail(err)
	s := GetDefaultSetting()
	err = json.Unmarshal(buffer, s)
	JustFail(err)

	return s
}
