package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kamva/mgm/v3"
	"log"
)

var _ planet.Updater = PlanetUpdate{}

type PlanetUpdate struct {
	collection *mgm.Collection
}

func NewUpdater() PlanetUpdate {
	return PlanetUpdate{
		collection: mgm.CollectionByName("planets"),
	}
}

func (p PlanetUpdate) Update(document planet.PlanetDocument) error {
	if err := p.collection.Update(&document); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
