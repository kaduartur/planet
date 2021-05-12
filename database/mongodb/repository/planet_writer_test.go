package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/database/mongodb"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanetWrite(t *testing.T) {
	mongodb.NewConnection()
	writer := NewWriter()
	defer writer.collection.Drop(mgm.Ctx())

	command := planet.CreatePlanetCommand{
		Name:    "Alderaan",
		Climate: "temperate",
		Terrain: "grasslands",
	}

	pd, err := writer.Write(command)

	assert.Nil(t, err)
	assert.NotEmpty(t, pd.PlanetID)
	assert.Equal(t, pd.Name, command.Name)
	assert.Equal(t, pd.Climate[0], command.Climate)
	assert.Equal(t, pd.Terrain[0], command.Terrain)
}
