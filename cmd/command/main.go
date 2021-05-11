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
	kafkaProducer := kafka.NewProducer()

	planetRead := repository.NewReader()
	planetWrite := repository.NewWriter()
	planetDelete := repository.NewDeleter(planetRead)

	r := gin.Default()
	createPlanet := http.NewCreatePlanetHandler(planetWrite, planetRead, kafkaProducer)
	deletePlanet := http.NewDeletePlanetHandler(planetDelete)

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
