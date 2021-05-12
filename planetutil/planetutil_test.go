package planetutil

import (
	"github.com/kaduartur/planet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePlanetID(t *testing.T) {
	const (
		name     = "Earth"
		expectID = planet.ID("vP7iWouvaAj85f9OY88hyNEUhTyn6s3MPCENc8WNq2Y=")
	)

	planetId := GeneratePlanetID(name)
	assert.Equal(t, expectID, planetId)
}
