package multicell

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// faces
const (
	Left = iota
	Top
	Right
	Bottom
	NumFaces // Numbef of faces per cell
)

const (
	Sqrt3              = 1.7320508075688772
	default_num_layers = 4
	default_len_face   = 50
	default_num_cell_x = 1
	default_num_cell_y = 1
	default_density    = 0.02 // genome matrix density
	default_W_scale    = 1.1
)

// sparse matrix of anything.
type SliceOfMaps[T any] []map[int]T

func (sm SliceOfMaps[T]) Do(f func(i, j int, v T)) {
	for i, mi := range sm {
		for j, v := range mi {
			f(i, j, v)
		}
	}
}

func NewSliceOfMaps[T any](n int) SliceOfMaps[T] {
	t := make([]map[int]T, n)
	for i := range n {
		t[i] = make(map[int]T)
	}
	return t
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
	WScale        float64 // scaling factor for genome.W mutation

	WithCue    bool                 // with cue or not
	MaxDevelop int                  // maximum number of developmental steps
	Alpha      float64              // weight for exponential moving average
	NumLayers  int                  // number of middle layers
	LenLayer   []int                // Length of each state vector
	DensityEM  float64              // input -> middle layer genome density
	Topology   SliceOfMaps[float64] // densities of genome matrices
	Omega      Vec                  // initial scaling factors of activation functions
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
		WScale:        default_W_scale,

		// parameters to be determined in SetModel are:
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
