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
	swapiURL := os.Getenv("SWAPI_URL")
	topic := os.Getenv("KAFKA_PLANET_TOPIC")

	mongodb.NewConnection()
	reader := repository.NewReader()
	updater := repository.NewUpdater()
	readUpdater := repository.NewReadUpdater(reader, updater)

	httpClient := &http.Client{Timeout: time.Second * 5}
	swapiClient := swapi.NewClient(swapiURL, httpClient)

	producer := kafka.NewProducer()
	createPlanetProcessor := event.NewCreatePlanetProcess(readUpdater, swapiClient, producer, topic, 5)
	eventsProcessor := planet.EventsProcessor{
		planet.CreatedEvent.String(): createPlanetProcessor,
	}

	consumer := kafka.NewConsumer(eventsProcessor)
	consumer.Read()
}
