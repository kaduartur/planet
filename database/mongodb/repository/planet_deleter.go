package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kamva/mgm/v3"
)

var _ planet.Deleter = PlanetDelete{}

type PlanetDelete struct {
	planet     planet.Reader
	collection *mgm.Collection
}

func NewDeleter(r planet.Reader) PlanetDelete {
	return PlanetDelete{
		planet:     r,
		collection: mgm.CollectionByName("planets"),
	}
}

func (p PlanetDelete) Delete(id planet.ID) error {
	if id == "" {
		return planet.ErrPlanetIDEmpty
	}

	pd, err := p.planet.ReadByPlanetId(id)
	if err != nil {
		return err
	}

	if err := p.collection.Delete(&pd); err != nil {
		return err
	}

	return nil
}
