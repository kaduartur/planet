package swapi

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestFindPlanetByName(t *testing.T) {
	type in struct {
		url        string
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
				url:        "http://127.0.0.1:8882",
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
				url:        "http://127.0.0.1:8882",
				planetName: "not found planet",
			},
			out: out{
				err: ErrPlanetNotFound,
			},
		},
		{
			name: "empty result",
			in: in{
				url:        "http://127.0.0.1:8882",
				planetName: "empty_result",
			},
			out: out{
				err: ErrPlanetNotFound,
			},
		},
		{
			name: "invalid json response",
			in: in{
				url:        "http://127.0.0.1:8882",
				planetName: "invalid_json",
			},
			out: out{
				err: io.EOF,
			},
		},
		{
			name: "server error",
			in: in{
				url:        "http://127.0.0.1:8882",
				planetName: "server_error",
			},
			out: out{
				err: errors.New("an unexpected error has occurred. STATUS=500 Server Error"),
			},
		},
		{
			name: "invalid url",
			in: in{
				url:        "invalid",
				planetName: "tatooine",
			},
			out: out{
				err: &url.Error{
					Op:  "Get",
					URL: "invalid/api/planets?search=tatooine",
					Err: errors.New("unsupported protocol scheme \"\""),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			swapi := NewClient(c.in.url, http.DefaultClient)
			got, err := swapi.FindPlanetByName(c.in.planetName)

			assert.Equal(t, c.out.err, err)
			assert.Equal(t,  c.out.want, got)
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
