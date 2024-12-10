package multicell

import (
	"encoding/json"
	"log"
	"math"
	"os"
)

const (
	sqrt3              = 1.732050807568877
	default_num_cell_x = 1
	default_num_cell_y = 1
	default_len_face   = 50
	default_density    = 0.02 // genome matrix density
)

// vector
type Vec = []float64

// sparse matrix
type SpMat = [](map[int]float64)

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
	Outdir        string  // output directory for trajectory
	WithCue       bool    // with cue or not
	MaxPopulation int     // maximum population size
	NumCellX      int     // number of cells in the x-axis
	NumCellY      int     // number of cells in the y-axis
	LenFace       int     // face length
	NumLayers     int     // number of middle layers
	MaxDevelop    int     // maximum number of developmental steps
	LenLayer      []int   // Length of each state vector
	DensityEM     float64 // input -> middle layer genome density
	DensityMP     float64 // middle -> output layer genome density
	Topology      SpMat   // densities of genome matrices
	Omega         Vec     // scaling factors of activation function inputs
	OmegaP        float64 // scaling factor of the output layer
	EnvNoise      float64 // noise level
	MutRate       float64 // mutation rate
	ConvDevelop   float64 // convergence limit
	Denv          float64 // size of an environmental change
	Selstrength   float64 // selection strength
	Alpha         float64 // weight for exponential moving average
	ProductionRun bool    // true if production run (i.e. "test" phase)
}

// LeCun-inspired arctan function
func LCatan(x float64) float64 {
	return 6.0 * math.Atan(x/sqrt3) / math.Pi
}

func GetDefaultSetting() *Setting {
	num_layers := 3
	num_components := make([]int, num_layers)
	omega := make(Vec, num_layers)
	topology := NewSpMat(num_layers)
	ncx := default_num_cell_x
	ncy := default_num_cell_y

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
		Basename:      "Full",
		Outdir:        ".",
		WithCue:       true,
		MaxPopulation: 500,
		NumCellX:      ncx,
		NumCellY:      ncy,
		LenFace:       default_len_face,
		NumLayers:     num_layers,
		MaxDevelop:    200,
		LenLayer:      num_components,
		DensityEM:     default_density,
		DensityMP:     default_density,
		Topology:      topology,
		Omega:         omega,
		OmegaP:        1.0,
		EnvNoise:      0.05,
		MutRate:       0.001,
		ConvDevelop:   1e-5,
		Denv:          0.5,
		Selstrength:   20.0,
		Alpha:         1.0 / 3.0,
		ProductionRun: false}
}

func (s *Setting) SetOmega() {
	s.Omega = make(Vec, s.NumLayers)
	for l, tl := range s.Topology {
		omega := 0.0
		if l == 0 {
			omega += s.DensityEM * float64(s.LenFace*NumFaces)
		}
		for k, d := range tl {
			omega += d * float64(s.LenLayer[k])
		}
		if omega > 0 {
			s.Omega[l] = 1.0 / math.Sqrt(omega)
		}
	}
	s.OmegaP = 1.0 /
		math.Sqrt(s.DensityMP*float64(s.LenLayer[s.NumLayers-1]))
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
	s := GetDefaultSetting()
	err = json.Unmarshal(buffer, s)
	JustFail(err)

	return s
}
