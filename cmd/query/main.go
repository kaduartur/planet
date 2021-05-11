package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet/database/mongodb"
	"github.com/kaduartur/planet/database/mongodb/repository"
	"github.com/kaduartur/planet/http"
	"os"
)

func main() {
	mongodb.NewConnection()
	planetRead := repository.NewReader()

	planetById := http.NewFindPlanetByIdHandler(planetRead)
	planetList := http.NewListPlanetHandler(planetRead)

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/planets", planetList.Handle)
		v1.GET("/planets/:planet_id", planetById.Handle)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	if err := router.Run(":" + port); err != nil {
		os.Exit(1)
	}
}
