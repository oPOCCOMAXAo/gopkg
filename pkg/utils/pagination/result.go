package pagination

type Page[T any] struct {
	Items []T
	Total int64
}
