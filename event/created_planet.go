package event

import (
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/swapi"
	"log"
	"path"
	"strings"
	"time"
)

const layoutISO = "2006-01-02"

type CreatedPlanet struct {
	planet planet.ReadUpdater
	api    swapi.Finder
}

func NewCreatePlanetProcess(planet planet.ReadUpdater, api swapi.Finder) CreatedPlanet {
	return CreatedPlanet{
		planet: planet,
		api:    api,
	}
}

func (c CreatedPlanet) Process(id planet.ID) {
	pd, err := c.planet.ReadByPlanetId(id)
	if err != nil {
		log.Println(err)
		return
	}

	if pd.Status == planet.Processed.String() {
		log.Printf("The planet %q has already been processed [PlanetID: %s]\n", pd.Name, pd.PlanetID)
		return
	}

	pd.Status = planet.Processed.String()
	planetRes, err := c.api.FindPlanetByName(pd.Name)
	if err != nil {
		log.Println(err)
		c.update(pd)
		return
	}

	var films planet.Films
	for _, url := range planetRes.Films {
		filmId := path.Base(url)
		film, err := c.api.FindFilmById(filmId)
		if err != nil {
			log.Println(err)
			continue
		}

		release, err := time.Parse(layoutISO, film.ReleaseData)
		if err != nil {
			log.Println(err)
			continue
		}

		f := planet.Film{
			Title:       film.Title,
			ReleaseDate: release,
		}

		films = append(films, f)
	}

	pd.Films = films
	pd.Climate = strings.Split(planetRes.Climate, ", ")
	pd.Terrain = strings.Split(planetRes.Terrain, ", ")

	c.update(pd)
}

func (c CreatedPlanet) update(pd planet.PlanetDocument)  {
	if err := c.planet.Update(pd); err != nil {
		log.Println(err)
	}
}
