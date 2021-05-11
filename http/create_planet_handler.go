package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/planetutil"
	"net/http"
)

type CreatePlanetHandler struct {
	planetW planet.Writer
	planetR planet.Reader
	kafka   planet.KafkaWriter
}

func NewCreatePlanetHandler(repoW planet.Writer, repoR planet.Reader, producer planet.KafkaWriter) CreatePlanetHandler {
	return CreatePlanetHandler{planetW: repoW, planetR: repoR, kafka: producer}
}

func (cp CreatePlanetHandler) Handle(c *gin.Context) {
	var createPlanet planet.CreatePlanetCommand
	if err := c.BindJSON(&createPlanet); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errors := createPlanet.Validate()
	if len(errors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"fields": errors})
		return
	}

	pd, err := cp.planetR.ReadByPlanetId(planetutil.GeneratePlanetID(createPlanet.Name))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if pd.PlanetID != "" {
		cp.success(c, pd)
		return
	}

	pd, err = cp.planetW.Write(createPlanet)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if err := cp.kafka.Write("planet-processor", planet.CreatedEvent, []byte(pd.PlanetID)); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	cp.success(c, pd)
}

func (cp CreatePlanetHandler) success(c *gin.Context, pd planet.PlanetDocument) {
	res := planet.StateResponse{
		PlanetID: pd.PlanetID,
		Status:   pd.Status,
	}

	c.JSON(http.StatusOK, res)
}
