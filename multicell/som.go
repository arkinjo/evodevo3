package multicell

// sparse matrix of anything.
type SliceOfMaps[T any] struct {
	M []map[int]T
}

func NewSliceOfMaps[T any](n int) SliceOfMaps[T] {
	t := make([]map[int]T, n)
	for i := range n {
		t[i] = make(map[int]T)
	}
	return SliceOfMaps[T]{t}
}

func (sm SliceOfMaps[T]) At(i, j int) T {
	return sm.M[i][j]
}

func (sm SliceOfMaps[T]) Set(i, j int, v T) {
	sm.M[i][j] = v
}

func (sm SliceOfMaps[T]) Do(f func(i, j int, v T)) {
	for i, mi := range sm.M {
		for j, v := range mi {
			f(i, j, v)
		}
	}
}

func (sm SliceOfMaps[T]) EachRow(f func(i int, mi map[int]T)) {
	for i, mi := range sm.M {
		f(i, mi)
	}
}
