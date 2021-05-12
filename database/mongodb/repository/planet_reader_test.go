package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/database/mongodb"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanetReadById(t *testing.T) {
	mongodb.NewConnection()
	reader := NewReader()
	pd := createPlanet()
	defer reader.collection.Drop(mgm.Ctx())

	got, err := reader.ReadByPlanetId(pd.PlanetID)
	assert.Nil(t, err)
	assert.Equal(t, pd.PlanetID, got.PlanetID)
	assert.Equal(t, pd.Name, got.Name)
	assert.Equal(t, pd.Climate, got.Climate)
	assert.Equal(t, pd.Terrain, got.Terrain)
}

func TestPlanetReadAll(t *testing.T) {
	mongodb.NewConnection()
	reader := NewReader()
	_ = createPlanet()
	defer reader.collection.Drop(mgm.Ctx())

	type in struct {
		pageReq planet.PageFilterRequest
	}

	cases := []struct {
		name string
		in   in
		want int
	}{
		{
			name: "success",
			in: in{
				pageReq: planet.PageFilterRequest{
					Page:    1,
					PerPage: 5,
				},
			},
			want: 1,
		},
		{
			name: "success find by name",
			in: in{
				pageReq: planet.PageFilterRequest{
					Page:    1,
					PerPage: 5,
					Name:    "Alderaan",
				},
			},
			want: 1,
		},
		{
			name: "not found planet by name",
			in: in{
				pageReq: planet.PageFilterRequest{
					Page:    1,
					PerPage: 5,
					Name:    "not found",
				},
			},
			want: 0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := reader.ReadAll(c.in.pageReq)
			assert.Nil(t, err)
			assert.Len(t, got, c.want)
		})
	}
}

func TestPlanetCount(t *testing.T) {
	mongodb.NewConnection()
	reader := NewReader()
	_ = createPlanet()
	defer reader.collection.Drop(mgm.Ctx())

	got, err := reader.Count()
	assert.Nil(t, err)
	assert.Equal(t, 1, got)
}

func createPlanet() planet.PlanetDocument {
	writer := NewWriter()
	command := planet.CreatePlanetCommand{
		Name:    "Alderaan",
		Climate: "temperate",
		Terrain: "grasslands",
	}

	pd, _ := writer.Write(command)
	return pd
}
