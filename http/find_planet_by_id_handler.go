package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"log"
	"net/http"
)

var _ planet.HttpHandler = FindPlanetByIdHandler{}

type FindPlanetByIdHandler struct {
	planet planet.Reader
}

func NewFindPlanetByIdHandler(planet planet.Reader) FindPlanetByIdHandler {
	return FindPlanetByIdHandler{planet: planet}
}

func (d FindPlanetByIdHandler) Handle(c *gin.Context) {
	planetID := c.Param("planet_id")
	log.Printf("Deleting planet [planetId=%s]", planetID)
	pd, err := d.planet.ReadByPlanetId(planet.ID(planetID))
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if pd.PlanetID == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, pd)
}
