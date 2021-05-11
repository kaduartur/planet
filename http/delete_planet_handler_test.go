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

func TestDeletePlanetHandle(t *testing.T) {
	type in struct {
		deleter  planet.Deleter
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
				deleter:  mock.PlanetRepository{},
				planetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
			},
			out: out{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name: "error to delete planet",
			in: in{
				deleter: mock.PlanetRepository{
					DeleteErr: planet.ErrPlanetIDEmpty,
				},
				planetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
			},
			out: out{
				statusCode: http.StatusUnprocessableEntity,
				resBody:    "{\"code\":\"PNT-003\",\"message\":\"The PlanetID must not be empty.\"}",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := setupDeletePlanetHandle(c.in.deleter)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/v1/planets/"+string(c.in.planetID), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, c.out.statusCode, w.Code)
			assert.Equal(t, c.out.resBody, w.Body.String())
		})
	}
}

func setupDeletePlanetHandle(
	deleter planet.Deleter,
) *gin.Engine {
	r := gin.Default()
	deletePlanet := NewDeletePlanetHandler(deleter)

	v1 := r.Group("/v1")
	{
		v1.DELETE("/planets/:planet_id", deletePlanet.Handle)
	}

	return r
}
