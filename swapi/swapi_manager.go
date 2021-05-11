package swapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var ErrPlanetNotFound = errors.New("the planet was not found")

type Manager struct {
	url    string
	client *http.Client
}

func NewClient(url string, client *http.Client) Manager {
	return Manager{url: url, client: client}
}

func (s Manager) FindPlanetByName(name string) (Planet, error) {
	planetsURL := fmt.Sprintf("%s/api/planets", s.url)
	req, err := http.NewRequest(http.MethodGet, planetsURL, nil)
	if err != nil {
		return Planet{}, err
	}

	params := url.Values{}
	params.Set("search", name)
	req.URL.RawQuery = params.Encode()

	res, err := s.client.Do(req)
	if err != nil {
		return Planet{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		return Planet{}, ErrPlanetNotFound
	}

	if res.StatusCode != http.StatusOK {
		return Planet{}, err
	}

	defer res.Body.Close()
	var body PlanetResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return Planet{}, err
	}

	if len(body.Planets) <= 0 {
		return Planet{}, ErrPlanetNotFound
	}

	return body.Planets[0], err
}

func (s Manager) FindFilmById(id string) (Film, error) {
	filmsURL := fmt.Sprintf("%s/api/films/%s", s.url, id)

	req, err := http.NewRequest(http.MethodGet, filmsURL, nil)
	if err != nil {
		return Film{}, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return Film{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		return Film{}, ErrPlanetNotFound
	}

	if res.StatusCode != http.StatusOK {
		return Film{}, err
	}

	defer res.Body.Close()
	var body Film
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return Film{}, err
	}

	return body, nil
}
