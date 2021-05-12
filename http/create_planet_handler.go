package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/planetutil"
	"net/http"
)

type CreatePlanetHandler struct {
	planetW    planet.Writer
	planetR    planet.Reader
	kafka      planet.KafkaWriter
	kafkaTopic string
}

func NewCreatePlanetHandler(
	repoW planet.Writer,
	repoR planet.Reader,
	producer planet.KafkaWriter,
	kafkaTopic string,
) CreatePlanetHandler {
	return CreatePlanetHandler{
		planetW:    repoW,
		planetR:    repoR,
		kafka:      producer,
		kafkaTopic: kafkaTopic,
	}
}

func (cp CreatePlanetHandler) Handle(c *gin.Context) {
	var createPlanet planet.CreatePlanetCommand
	if err := c.BindJSON(&createPlanet); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, planet.ErrDecodeRequest)
		return
	}

	errors := createPlanet.Validate()
	if len(errors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"fields": errors})
		return
	}

	pd, err := cp.planetR.ReadByPlanetId(planetutil.GeneratePlanetID(createPlanet.Name))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if pd.PlanetID != "" {
		cp.success(c, pd)
		return
	}

	pd, err = cp.planetW.Write(createPlanet)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	event := planet.CreatePlanetEvent{
		PlanetID: pd.PlanetID,
	}

	data, err := json.Marshal(event)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, planet.ErrUnknown)
		return
	}

	if err := cp.kafka.Write(cp.kafkaTopic, planet.CreatedEvent, data); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
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
