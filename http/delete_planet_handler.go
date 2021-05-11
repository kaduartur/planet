package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"log"
	"net/http"
)

var _ planet.HttpHandler = DeletePlanetHandler{}

type DeletePlanetHandler struct {
	planet planet.Deleter
}

func NewDeletePlanetHandler(planet planet.Deleter) DeletePlanetHandler {
	return DeletePlanetHandler{planet: planet}
}

func (d DeletePlanetHandler) Handle(c *gin.Context) {
	planetID := c.Param("planet_id")
	log.Printf("Deleting planet [planetId=%s]", planetID)
	if err := d.planet.Delete(planet.ID(planetID)); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	c.Status(http.StatusNoContent)
}
