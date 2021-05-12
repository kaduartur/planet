package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet/database/mongodb"
	"github.com/kaduartur/planet/database/mongodb/repository"
	"github.com/kaduartur/planet/http"
	"github.com/kaduartur/planet/kafka"
	"os"
)

func main() {
	mongodb.NewConnection()
	planetRead := repository.NewReader()
	planetWrite := repository.NewWriter()
	planetDelete := repository.NewDeleter(planetRead)

	kafkaProducer := kafka.NewProducer()

	topic := os.Getenv("KAFKA_PLANET_TOPIC")
	createPlanet := http.NewCreatePlanetHandler(planetWrite, planetRead, kafkaProducer, topic)
	deletePlanet := http.NewDeletePlanetHandler(planetDelete)

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/planets", createPlanet.Handle)
		v1.DELETE("/planets/:planet_id", deletePlanet.Handle)
	}

	gin.Logger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := r.Run(":" + port); err != nil {
		os.Exit(1)
	}
}
