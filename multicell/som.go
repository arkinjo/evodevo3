package multicell

// sparse matrix of anything.
type SliceOfMaps[T any] []map[int]T

func NewSliceOfMaps[T any](n int) SliceOfMaps[T] {
	t := make([]map[int]T, n)
	for i := range n {
		t[i] = make(map[int]T)
	}
	return t
}

func (sm SliceOfMaps[T]) Do(f func(i, j int, v T)) {
	for i, mi := range sm {
		for j, v := range mi {
			f(i, j, v)
		}
	}
}
