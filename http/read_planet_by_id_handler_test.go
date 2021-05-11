package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/internal/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindPlanetByIdHandle(t *testing.T) {
	type in struct {
		reader   planet.Reader
		planetID planet.ID
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
				reader: mock.PlanetRepository{
					DocReadById: planet.PlanetDocument{
						PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
						Name:     "Toosda",
						Climate:  []string{"arid"},
						Terrain:  []string{"desert"},
						Films:    planet.Films{},
						Status:   planet.Processed.String(),
					},
				},
				planetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
			},
			out: out{
				statusCode: http.StatusOK,
				resBody:    "{\"planet_id\":\"PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=\",\"name\":\"Toosda\",\"climate\":[\"arid\"],\"terrain\":[\"desert\"],\"films\":[],\"status\":\"PROCESSED\",\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}",
			},
		},
		{
			name: "error to find planet by id",
			in: in{
				reader: mock.PlanetRepository{
					ReadByIdErr: planet.ErrUnknown,
				},
				planetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
			},
			out: out{
				statusCode: http.StatusUnprocessableEntity,
				resBody:    "{\"code\":\"PNT-001\",\"message\":\"There was an unknown error processing your request.\"}",
			},
		},
		{
			name: "error planet not found",
			in: in{
				reader:   mock.PlanetRepository{},
				planetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
			},
			out: out{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := setupReadPlanetByIdHandle(c.in.reader)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v1/planets/"+string(c.in.planetID), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, c.out.statusCode, w.Code)
			assert.Equal(t, c.out.resBody, w.Body.String())
		})
	}
}

func setupReadPlanetByIdHandle(reader planet.Reader) *gin.Engine {
	r := gin.Default()
	planetById := NewFindPlanetByIdHandler(reader)

	v1 := r.Group("/v1")
	{
		v1.GET("/planets/:planet_id", planetById.Handle)
	}

	return r
}
