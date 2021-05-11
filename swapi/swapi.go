package swapi

type Finder interface {
	FindPlanetByName(name string) (Planet, error)
	FindFilmById(id string) (Film, error)
}

type PlanetResponse struct {
	Planets Planets `json:"results"`
}

type Planets []Planet

type Planet struct {
	Name    string   `json:"name"`
	Films   []string `json:"films"`
	Climate string   `json:"climate"`
	Terrain string   `json:"terrain"`
}

type Film struct {
	Title       string `json:"title"`
	ReleaseData string `json:"release_date"`
}
