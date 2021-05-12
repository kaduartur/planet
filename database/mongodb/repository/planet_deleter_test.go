package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/database/mongodb"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanetDelete(t *testing.T) {
	mongodb.NewConnection()
	reader := NewReader()
	deleter := NewDeleter(reader)
	pd := createPlanet()
	defer deleter.collection.Drop(mgm.Ctx())

	cases := []struct {
		name string
		in   planet.ID
		want error
	}{
		{
			name: "success",
			in:   pd.PlanetID,
		},
		{
			name: "error empty planet id",
			in:   "",
			want: planet.ErrPlanetIDEmpty,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := deleter.Delete(c.in)
			assert.Equal(t, c.want, got)
		})
	}
}
