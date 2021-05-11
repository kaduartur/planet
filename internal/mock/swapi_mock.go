package mock

import "github.com/kaduartur/planet/swapi"

type Swapi struct {
	Planet    swapi.Planet
	Film      swapi.Film
	PlanetErr error
	FilmErr   error
}

func (s Swapi) FindPlanetByName(name string) (swapi.Planet, error) {
	return s.Planet, s.PlanetErr
}

func (s Swapi) FindFilmById(id string) (swapi.Film, error) {
	return s.Film, s.FilmErr
}
