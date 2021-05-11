package event

import (
	"encoding/json"
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/internal/mock"
	"github.com/kaduartur/planet/swapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcess(t *testing.T) {
	type in struct {
		repo     planet.ReadUpdater
		swapi    swapi.Finder
		producer planet.KafkaWriter
		retry    int
		event    planet.CreatePlanetEvent
	}

	cases := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				event: planet.CreatePlanetEvent{
					PlanetID:   "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
					RetryCount: 0,
				},
				repo: mock.PlanetRepository{
					DocReadById: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Status:   planet.Processing.String(),
					},
				},
				swapi: mock.Swapi{
					Planet: swapi.Planet{
						Name: "Tatooine",
						Films: []string{
							"https://swapi.dev/api/films/1",
						},
						Climate: "arid",
						Terrain: "desert",
					},
					Film: swapi.Film{
						Title:       "test",
						ReleaseData: "2021-06-25",
					},
				},
				producer: mock.KafkaWrite{},
				retry:    1,
			},
		},
		{
			name: "error to read planet by id",
			in: in{
				repo: mock.PlanetRepository{
					ReadByIdErr: planet.ErrUnknown,
				},
				swapi:    mock.Swapi{},
				producer: mock.KafkaWrite{},
				retry:    1,
			},
			want: planet.ErrUnknown,
		},
		{
			name: "error planet already processed",
			in: in{
				repo: mock.PlanetRepository{
					DocReadById: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Status:   planet.Processed.String(),
					},
				},
				swapi:    mock.Swapi{},
				producer: mock.KafkaWrite{},
				retry:    1,
			},
			want: planet.ErrPlanetAlreadyProcessed,
		},
		{
			name: "error to request planet to swapi",
			in: in{
				event: planet.CreatePlanetEvent{
					PlanetID:   "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
					RetryCount: 0,
				},
				repo: mock.PlanetRepository{
					DocReadById: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Status:   planet.Processing.String(),
					},
				},
				swapi: mock.Swapi{
					PlanetErr: swapi.ErrPlanetNotFound,
				},
				producer: mock.KafkaWrite{},
				retry:    1,
			},
			want: swapi.ErrPlanetNotFound,
		},
		{
			name: "error to update planet",
			in: in{
				event: planet.CreatePlanetEvent{
					PlanetID:   "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
					RetryCount: 0,
				},
				repo: mock.PlanetRepository{
					DocReadById: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Status:   planet.Processing.String(),
					},
					UpdateErr: planet.ErrUnknown,
				},
				swapi:    mock.Swapi{},
				producer: mock.KafkaWrite{},
				retry:    1,
			},
			want: planet.ErrUnknown,
		},
		{
			name: "update planet with failed status after exceeding all retries",
			in: in{
				event: planet.CreatePlanetEvent{
					PlanetID:   "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
					RetryCount: 1,
				},
				repo: mock.PlanetRepository{
					DocReadById: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Status:   planet.Processing.String(),
					},
				},
				swapi:    mock.Swapi{
					PlanetErr: swapi.ErrPlanetNotFound,
				},
				producer: mock.KafkaWrite{},
				retry:    1,
			},
			want: swapi.ErrPlanetNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			process := NewCreatePlanetProcess(c.in.repo, c.in.swapi, c.in.producer, c.in.retry)

			b, _ := json.Marshal(c.in.event)

			got := process.Process(b)

			assert.Equal(t, c.want, got)
		})
	}
}
