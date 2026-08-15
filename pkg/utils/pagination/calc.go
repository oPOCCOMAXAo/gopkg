package pagination

// TotalPages returns the number of pages needed to fit count items
// with the given perPage limit.
func TotalPages(count int64, perPage int64) int64 {
	if perPage <= 0 {
		return 0
	}

	totalPages := count / perPage
	if count%perPage > 0 {
		totalPages++
	}

	return totalPages
}

// Offset returns the zero-based offset for the given 1-based page number
// with the given perPage limit.
func Offset(page int64, perPage int64) int64 {
	if page < 1 {
		return 0
	}

	return (page - 1) * perPage
}
