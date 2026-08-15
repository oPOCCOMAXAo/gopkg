package pagination

import (
	"iter"
	"strconv"
)

type PagesParams struct {
	// Total number of pages.
	Total int64

	// Current page number.
	Current int64

	// FirstN number of pages to show from start.
	FirstN int64

	// LastN number of pages to show from end.
	LastN int64

	// NearN number of pages to show around current page.
	NearN int64
}

type PageData struct {
	// Page number.
	Number int64

	// IsCurrent is current selected/opened page.
	IsCurrent bool

	// PrevSkipped is true if there are skipped pages before this page.
	PrevSkipped bool
}

func (p *PageData) String() string {
	return strconv.FormatInt(p.Number, 10)
}

// Pages returns iterator for pagination render.
func Pages(params PagesParams) iter.Seq[PageData] {
	return func(yield func(PageData) bool) {
		if params.Total <= 0 {
			return
		}

		params.Current = min(max(params.Current, 1), params.Total)
		params.FirstN = max(params.FirstN, 0)
		params.LastN = max(params.LastN, 0)
		params.NearN = max(params.NearN, 0)

		ranges := [][2]int64{
			{1, min(params.FirstN, params.Total)},
			{max(params.Current-params.NearN, 1), min(params.Current+params.NearN, params.Total)},
			{max(params.Total-params.LastN+1, 1), params.Total},
		}

		var prev int64 = 0

		for _, r := range ranges {
			lo, hi := r[0], r[1]

			for page := lo; page <= hi; page++ {
				if page <= prev {
					continue
				}

				if !yield(PageData{
					Number:      page,
					IsCurrent:   page == params.Current,
					PrevSkipped: prev > 0 && page > prev+1,
				}) {
					return
				}

				prev = page
			}
		}
	}
}
