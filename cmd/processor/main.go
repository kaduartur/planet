package main

import (
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/database/mongodb"
	"github.com/kaduartur/planet/database/mongodb/repository"
	"github.com/kaduartur/planet/event"
	"github.com/kaduartur/planet/kafka"
	"github.com/kaduartur/planet/swapi"
	"net/http"
	"os"
	"time"
)

func main() {
	mongodb.NewConnection()
	reader := repository.NewReader()
	updater := repository.NewUpdater()
	readUpdater := repository.NewReadUpdater(reader, updater)

	httpClient := &http.Client{Timeout: time.Second * 5}
	swapiClient := swapi.NewClient(os.Getenv("SWAPI_URL"), httpClient)

	createPlanetProcessor := event.NewCreatePlanetProcess(readUpdater, swapiClient)
	eventsProcessor := planet.EventsProcessor{
		planet.CreatedEvent.String(): createPlanetProcessor,
	}

	consumer := kafka.NewConsumer(eventsProcessor)
	consumer.Read()

}
