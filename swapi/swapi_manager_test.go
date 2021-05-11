package swapi

import (
	"net/http"
	"reflect"
	"testing"
)

func TestFindPlanetByName(t *testing.T) {
	type in struct {
		planetName string
	}

	type out struct {
		err  error
		want []Planet
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
				want: []Planet{
					{
						Name: "Tatooine",
						Films: []string{
							"http://swapi.dev/api/films/1/",
							"http://swapi.dev/api/films/3/",
							"http://swapi.dev/api/films/4/",
							"http://swapi.dev/api/films/5/",
							"http://swapi.dev/api/films/6/"},
					},
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			swapi := NewClient("http://127.0.0.1:8882", http.DefaultClient)
			got, err := swapi.FindPlanetByName(c.in.planetName)

			if err != c.out.err {
				t.Errorf("FindPlanetByName(%q) == %q, want %q", c.name, err, c.out.err)
			}

			if !reflect.DeepEqual(got, c.out.want) {
				t.Errorf("FindPlanetByName(%q) == %q, want %q", c.name, got, c.out.want)
			}
		})
	}

}
