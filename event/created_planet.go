package event

import (
	"encoding/json"
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/swapi"
	"log"
	"path"
	"strings"
	"time"
)

const layoutISO = "2006-01-02"

type CreatedPlanet struct {
	planet          planet.ReadUpdater
	api             swapi.Finder
	kafka           planet.KafkaWriter
	kafkaTopic      string
	maxRetryProcess int
}

func NewCreatePlanetProcess(
	planet planet.ReadUpdater,
	api swapi.Finder,
	kafka planet.KafkaWriter,
	kafkaTopic string,
	retryProcess int,
) CreatedPlanet {
	return CreatedPlanet{
		planet:          planet,
		api:             api,
		kafka:           kafka,
		kafkaTopic:      kafkaTopic,
		maxRetryProcess: retryProcess,
	}
}

func (c CreatedPlanet) Process(data interface{}) error {
	var event planet.CreatePlanetEvent
	_ = json.Unmarshal(data.([]byte), &event)

	pd, err := c.planet.ReadByPlanetId(event.PlanetID)
	if err != nil {
		return err
	}

	if pd.Status == planet.Processed.String() {
		log.Printf("The planet %q has already been processed [PlanetID: %s]\n", pd.Name, pd.PlanetID)
		return planet.ErrPlanetAlreadyProcessed
	}

	planetRes, err := c.api.FindPlanetByName(pd.Name)
	if err != nil {
		log.Printf("error to find planet into swapi [planetId: %s] - [error: %s]", pd.PlanetID, err)
		c.retry(pd, event)
		return err
	}

	var films planet.Films
	for _, url := range planetRes.Films {
		filmId := path.Base(url)
		film, err := c.api.FindFilmById(filmId)
		if err != nil {
			log.Println(err)
			c.retry(pd, event)
			return err
		}

		release, err := time.Parse(layoutISO, film.ReleaseData)
		if err != nil {
			log.Println(err)
			return err
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
	pd.Status = planet.Processed.String()

	if err := c.update(pd); err != nil {
		c.retry(pd, event)
		return err
	}

	return nil
}

func (c CreatedPlanet) update(pd planet.PlanetDocument) error {
	if err := c.planet.Update(pd); err != nil {
		log.Printf("Error to update planet [planet: %+v] - [error: %s]\n", pd, err)
		return err
	}

	return nil
}

func (c CreatedPlanet) retry(pd planet.PlanetDocument, event planet.CreatePlanetEvent) {
	if event.RetryCount >= c.maxRetryProcess {
		log.Printf("all retries have already been carried out\n")

		pd.Status = planet.Failed.String()
		_ = c.update(pd)
		return
	}

	event.RetryCount++
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("error to retry event [event: %+v]\n", event)
		return
	}

	if err := c.kafka.Write(c.kafkaTopic, planet.CreatedEvent, data)
		err != nil {
		log.Printf("error to write message into kafka")
		return
	}
}
