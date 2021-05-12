package sliceutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMerge(t *testing.T) {
	type in struct {
		aa []string
		bb []string
	}
	cases := []struct {
		name string
		in   in
		want []string
	}{
		{
			name: "success",
			in: in{
				aa: []string{"a", "b", "c", "x"},
				bb: []string{"a", "b", "c", "d"},
			},
			want: []string{"a", "b", "c", "d", "x"},
		},
		{
			name: "success and remove empty values",
			in: in{
				aa: []string{"a", "", "c", "x"},
				bb: []string{"a", "c", "d"},
			},
			want: []string{"a", "c", "d", "x"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Merge(c.in.aa, c.in.bb...)
			assert.Equal(t, c.want, got)
		})
	}
}
