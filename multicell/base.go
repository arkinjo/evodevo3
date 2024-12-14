package multicell

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
)

const (
	Sqrt3              = 1.732050807568877
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

// various set-ups
type Setting struct {
	Basename      string  // name of the model
	Seed          uint64  // random seed
	Outdir        string  // output directory for trajectory
	MaxPopulation int     // maximum population size
	MaxGeneration int     // maximum number of generations per epoch
	WithCue       bool    // with cue or not
	NumCellX      int     // number of cells in the x-axis
	NumCellY      int     // number of cells in the y-axis
	LenFace       int     // face length
	NumLayers     int     // number of middle layers
	MaxDevelop    int     // maximum number of developmental steps
	LenLayer      []int   // Length of each state vector
	DensityEM     float64 // input -> middle layer genome density
	Topology      SpMat   // densities of genome matrices
	Omega         Vec     // scaling factors of activation function inputs
	EnvNoise      float64 // noise level
	MutRate       float64 // mutation rate
	ConvDevelop   float64 // convergence limit
	Denv          float64 // size of an environmental change
	SelStrength   float64 // selection strength
	Alpha         float64 // weight for exponential moving average
	ProductionRun bool    // true if production run (i.e. "test" phase)
}

// LeCun-inspired arctan function
func LCatan(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		return 6.0 * math.Atan(omega*x/Sqrt3) / math.Pi
	}
}

func Tanh(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		return math.Tanh(omega * x)
	}
}

func GetDefaultSetting() *Setting {
	num_layers := default_num_layers
	num_components := make([]int, num_layers)
	topology := NewSpMat(num_layers, num_layers)
	ncx := default_num_cell_x
	ncy := default_num_cell_y

	for i := 0; i < num_layers; i++ {
		num_components[i] = default_len_state
		if i < num_layers-1 { // no loop for the last layer.
			topology.SetAt(i, i, default_density)
		}
		if i > 0 {
			topology.SetAt(i, i-1, default_density)
		}
	}

	return &Setting{
		Basename:      "Full",
		Seed:          13,
		Outdir:        ".",
		MaxGeneration: 200,
		WithCue:       true,
		MaxPopulation: 500,
		NumCellX:      ncx,
		NumCellY:      ncy,
		LenFace:       default_len_face,
		NumLayers:     num_layers,
		MaxDevelop:    200,
		LenLayer:      num_components,
		DensityEM:     default_density,
		Topology:      topology,
		EnvNoise:      0.05,
		MutRate:       0.001,
		ConvDevelop:   1e-5,
		Denv:          0.5,
		SelStrength:   20.0,
		Alpha:         1.0 / 3.0,
		ProductionRun: false}
	// Omega is set in SetOmega().
}

func (s *Setting) SetOmega() {
	s.Omega = make(Vec, s.NumLayers)
	for l := range s.NumLayers {
		omega := 0.0
		if l == 0 {
			omega += s.DensityEM * float64(s.LenFace*NumFaces)
		}
		for k := range s.NumLayers {
			omega += s.Topology.At(l, k) * float64(s.LenLayer[k])
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
