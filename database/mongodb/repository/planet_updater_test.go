package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/database/mongodb"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanetUpdate(t *testing.T) {
	mongodb.NewConnection()
	reader := NewReader()
	updater := NewUpdater()
	readUpdater := NewReadUpdater(reader, updater)
	pd := createPlanet()
	defer updater.collection.Drop(mgm.Ctx())

	pd.Status = planet.Processing.String()
	err := readUpdater.Update(pd)
	assert.Nil(t, err)

	pd, _ = readUpdater.ReadByPlanetId(pd.PlanetID)
	assert.Equal(t, pd.Status, planet.Processing.String())
}
