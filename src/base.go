package evodevo3

type Setting struct {
	basename       string
	seed           int
	with_cue       bool
	max_pop        int
	num_cell_x     int
	num_cell_y     int
	num_env        int
	num_layers     int
	num_dev        int
	num_components []int
	topology       Spmat
}
