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
	Sqrt3                 = 1.7320508075688772
	default_num_layers    = 4
	default_len_face      = 50
	default_num_cell_x    = 1
	default_num_cell_y    = 1
	default_density       = 0.02 // genome matrix density
	default_mutation_rate = 0.001
	default_W_scale       = 1.1
)

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
	Topology   SliceOfMaps[float64] // densities of genome matrices
	Omega      Vec                  // initial scaling factors of activation functions
}

func GetDefaultSetting(modelname string) *Setting {
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
		MutRate:       default_mutation_rate,
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
		//Topology
		//Omega
	}

	s.SetModel(modelname)
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
	var s Setting
	err = json.Unmarshal(buffer, &s)
	JustFail(err)

	return &s
}
