package http

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/internal/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePlanetHandle(t *testing.T) {
	type in struct {
		writer      planet.Writer
		reader      planet.Reader
		kafkaWriter planet.KafkaWriter
		command     planet.CreatePlanetCommand
	}

	type out struct {
		statusCode int
		resBody    string
	}

	cases := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success",
			in: in{
				writer: mock.PlanetRepository{
					DocWrite: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Status:   planet.Processing.String(),
					},
				},
				reader:      mock.PlanetRepository{},
				kafkaWriter: mock.KafkaWrite{},
				command: planet.CreatePlanetCommand{
					Name:    "Tatooine",
					Climate: "arid",
					Terrain: "desert",
				},
			},
			out: out{
				statusCode: http.StatusOK,
				resBody:    "{\"planet_id\":\"PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=\",\"status\":\"PROCESSING\"}",
			},
		},
		{
			name: "success idempotent",
			in: in{
				reader: mock.PlanetRepository{
					DocReadById: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Status:   planet.Processing.String(),
					},
				},
				command: planet.CreatePlanetCommand{
					Name:    "Tatooine",
					Climate: "arid",
					Terrain: "desert",
				},
			},
			out: out{
				statusCode: http.StatusOK,
				resBody:    "{\"planet_id\":\"PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=\",\"status\":\"PROCESSING\"}",
			},
		},
		{
			name: "invalid request body",
			in: in{
				command: planet.CreatePlanetCommand{},
			},
			out: out{
				statusCode: http.StatusBadRequest,
				resBody:    "{\"fields\":{\"climate\":\"This field cannot be empty.\",\"name\":\"This field cannot be empty.\",\"terrain\":\"This field cannot be empty.\"}}",
			},
		},
		{
			name: "error to read planet by id",
			in: in{
				reader: mock.PlanetRepository{
					ReadByIdErr: planet.ErrUnknown,
				},
				command: planet.CreatePlanetCommand{
					Name:    "Tatooine",
					Climate: "arid",
					Terrain: "desert",
				},
			},
			out: out{
				statusCode: http.StatusBadRequest,
				resBody:    "{\"code\":\"PNT-001\",\"message\":\"There was an unknown error processing your request.\"}",
			},
		},
		{
			name: "error to write planet",
			in: in{
				reader: mock.PlanetRepository{},
				writer: mock.PlanetRepository{
					WriteErr: planet.ErrUnknown,
				},
				command: planet.CreatePlanetCommand{
					Name:    "Tatooine",
					Climate: "arid",
					Terrain: "desert",
				},
			},
			out: out{
				statusCode: http.StatusUnprocessableEntity,
				resBody:    "{\"code\":\"PNT-001\",\"message\":\"There was an unknown error processing your request.\"}",
			},
		},
		{
			name: "error to write planet",
			in: in{
				reader: mock.PlanetRepository{},
				writer: mock.PlanetRepository{
					WriteErr: planet.ErrUnknown,
				},
				command: planet.CreatePlanetCommand{
					Name:    "Tatooine",
					Climate: "arid",
					Terrain: "desert",
				},
			},
			out: out{
				statusCode: http.StatusUnprocessableEntity,
				resBody:    "{\"code\":\"PNT-001\",\"message\":\"There was an unknown error processing your request.\"}",
			},
		},
		{
			name: "error to publish planet message to kafka",
			in: in{
				reader: mock.PlanetRepository{},
				writer: mock.PlanetRepository{},
				kafkaWriter: mock.KafkaWrite{
					Err: planet.ErrUnknown,
				},
				command: planet.CreatePlanetCommand{
					Name:    "Tatooine",
					Climate: "arid",
					Terrain: "desert",
				},
			},
			out: out{
				statusCode: http.StatusUnprocessableEntity,
				resBody:    "{\"code\":\"PNT-001\",\"message\":\"There was an unknown error processing your request.\"}",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := setupCreatePlanetHandle(c.in.writer, c.in.reader, c.in.kafkaWriter)
			bodyReq, _ := json.Marshal(c.in.command)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/v1/planets", bytes.NewBuffer(bodyReq))
			router.ServeHTTP(w, req)

			assert.Equal(t, c.out.statusCode, w.Code)
			assert.Equal(t, c.out.resBody, w.Body.String())
		})
	}
}

func setupCreatePlanetHandle(
	writer planet.Writer,
	reader planet.Reader,
	kafkaWriter planet.KafkaWriter,
) *gin.Engine {
	r := gin.Default()
	createPlanet := NewCreatePlanetHandler(writer, reader, kafkaWriter)

	v1 := r.Group("/v1")
	{
		v1.POST("/planets", createPlanet.Handle)
	}

	return r
}
