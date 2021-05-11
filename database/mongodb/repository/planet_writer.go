package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/planetutil"
	"github.com/kamva/mgm/v3"
	"log"
)

var _ planet.Writer = PlanetWrite{}

type PlanetWrite struct {
	collection *mgm.Collection
}

func NewWriter() PlanetWrite {
	return PlanetWrite{
		collection: mgm.CollectionByName("planets"),
	}
}

func (p PlanetWrite) Write(cmd planet.CreatePlanetCommand) (planet.PlanetDocument, error) {
	planetID := planetutil.GeneratePlanetID(cmd.Name)
	document := &planet.PlanetDocument{
		PlanetID: planetID,
		Name:     cmd.Name,
		Climate:  []string{cmd.Climate},
		Terrain:  []string{cmd.Terrain},
		Films:    make(planet.Films, 0),
		Status:   planet.Processing.String(),
	}

	if err := p.collection.Create(document); err != nil {
		log.Printf("Error to write planet [Planet: %+v] - [Error: %s]\n", cmd, err)
		return planet.PlanetDocument{}, planet.ErrUnknown
	}

	return *document, nil
}
