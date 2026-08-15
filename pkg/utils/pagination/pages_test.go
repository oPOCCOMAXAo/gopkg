package pagination

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPages(t *testing.T) {
	tests := []struct {
		name   string
		params PagesParams
		want   []PageData
	}{
		{
			name:   "zero total",
			params: PagesParams{Total: 0, Current: 1, FirstN: 2, LastN: 2, NearN: 1},
			want:   nil,
		},
		{
			name:   "single page",
			params: PagesParams{Total: 1, Current: 1, FirstN: 2, LastN: 2, NearN: 1},
			want: []PageData{
				{Number: 1, IsCurrent: true},
			},
		},
		{
			name:   "5 pages current=3 overlaps first/last/near",
			params: PagesParams{Total: 5, Current: 3, FirstN: 2, LastN: 2, NearN: 1},
			want: []PageData{
				{Number: 1, IsCurrent: false},
				{Number: 2, IsCurrent: false},
				{Number: 3, IsCurrent: true},
				{Number: 4, IsCurrent: false},
				{Number: 5, IsCurrent: false},
			},
		},
		{
			name:   "10 pages current=5 with gaps",
			params: PagesParams{Total: 10, Current: 5, FirstN: 2, LastN: 2, NearN: 1},
			want: []PageData{
				{Number: 1, IsCurrent: false, PrevSkipped: false},
				{Number: 2, IsCurrent: false, PrevSkipped: false},
				{Number: 4, IsCurrent: false, PrevSkipped: true},
				{Number: 5, IsCurrent: true, PrevSkipped: false},
				{Number: 6, IsCurrent: false, PrevSkipped: false},
				{Number: 9, IsCurrent: false, PrevSkipped: true},
				{Number: 10, IsCurrent: false, PrevSkipped: false},
			},
		},
		{
			name:   "10 pages current=1",
			params: PagesParams{Total: 10, Current: 1, FirstN: 2, LastN: 2, NearN: 1},
			want: []PageData{
				{Number: 1, IsCurrent: true, PrevSkipped: false},
				{Number: 2, IsCurrent: false, PrevSkipped: false},
				{Number: 9, IsCurrent: false, PrevSkipped: true},
				{Number: 10, IsCurrent: false, PrevSkipped: false},
			},
		},
		{
			name:   "10 pages current=10",
			params: PagesParams{Total: 10, Current: 10, FirstN: 2, LastN: 2, NearN: 1},
			want: []PageData{
				{Number: 1, IsCurrent: false, PrevSkipped: false},
				{Number: 2, IsCurrent: false, PrevSkipped: false},
				{Number: 9, IsCurrent: false, PrevSkipped: true},
				{Number: 10, IsCurrent: true, PrevSkipped: false},
			},
		},
		{
			name:   "current clamped below",
			params: PagesParams{Total: 5, Current: 0, FirstN: 2, LastN: 2, NearN: 1},
			want: []PageData{
				{Number: 1, IsCurrent: true, PrevSkipped: false},
				{Number: 2, IsCurrent: false, PrevSkipped: false},
				{Number: 4, IsCurrent: false, PrevSkipped: true},
				{Number: 5, IsCurrent: false, PrevSkipped: false},
			},
		},
		{
			name:   "current clamped above",
			params: PagesParams{Total: 5, Current: 99, FirstN: 2, LastN: 2, NearN: 1},
			want: []PageData{
				{Number: 1, IsCurrent: false, PrevSkipped: false},
				{Number: 2, IsCurrent: false, PrevSkipped: false},
				{Number: 4, IsCurrent: false, PrevSkipped: true},
				{Number: 5, IsCurrent: true, PrevSkipped: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Pages(tt.params))

			require.Equal(t, tt.want, got)
		})
	}
}
