package swapi

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFindPlanetByName(t *testing.T) {
	type in struct {
		planetName string
	}

	type out struct {
		err  error
		want Planet
	}

	cases := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success",
			in: in{
				planetName: "tatooine",
			},
			out: out{
				want: Planet{
					Name:    "Tatooine",
					Climate: "arid",
					Terrain: "desert",
					Films:   []string{"http://swapi.dev/api/films/1/"},
				},
			},
		},
		{
			name: "not found",
			in: in{
				planetName: "not found planet",
			},
			out: out{
				err: ErrPlanetNotFound,
			},
		},
		{
			name: "server error",
			in: in{
				planetName: "server_error",
			},
			out: out{
				err: errors.New("an unexpected error has occurred. STATUS=500 Server Error"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			swapi := NewClient("http://127.0.0.1:8882", http.DefaultClient)
			got, err := swapi.FindPlanetByName(c.in.planetName)

			assert.Equal(t, err, c.out.err)
			assert.Equal(t, got, c.out.want)
		})
	}
}

func TestFindFilmById(t *testing.T) {
	type in struct {
		id string
	}

	type out struct {
		err  error
		want Film
	}

	cases := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success",
			in: in{
				id: "1",
			},
			out: out{
				want: Film{
					Title:       "A New Hope",
					ReleaseData: "1977-05-25",
				},
			},
		},
		{
			name: "not found",
			in: in{
				id: "3",
			},
			out: out{
				err: ErrFilmNotFound,
			},
		},
		{
			name: "server error",
			in: in{
				id: "5",
			},
			out: out{
				err: errors.New("an unexpected error has occurred. STATUS=500 Server Error"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			swapi := NewClient("http://127.0.0.1:8882", http.DefaultClient)
			got, err := swapi.FindFilmById(c.in.id)

			assert.Equal(t, err, c.out.err)
			assert.Equal(t, got, c.out.want)
		})
	}

}
