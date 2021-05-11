package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"github.com/kaduartur/planet/internal/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestListPlanetHandle(t *testing.T) {
	type in struct {
		reader  planet.Reader
		page    string
		perPage string
		name    string
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
				page:    "1",
				perPage: "10",
				reader: mock.PlanetRepository{
					DocsReadAll: []planet.PlanetDocument{
						{
							PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
							Name:     "Tatooine",
							Climate:  []string{"arid"},
							Terrain:  []string{"desert"},
							Films:    planet.Films{},
							Status:   planet.Processed.String(),
						},
						{
							PlanetID: "dfhtHhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
							Name:     "Yavin IV",
							Climate:  []string{"temperate", "tropical"},
							Terrain:  []string{"jungle", "rainforests"},
							Films:    planet.Films{},
							Status:   planet.Processed.String(),
						},
					},
					CountValue: 2,
				},
			},
			out: out{
				statusCode: http.StatusOK,
				resBody: "{\"_metadata\":{\"page\":1,\"per_page\":10,\"page_count\":2,\"total_count\":2}," +
					"\"results\":[{\"planet_id\":\"PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=\",\"name\":" +
					"\"Tatooine\",\"climate\":[\"arid\"],\"terrain\":[\"desert\"],\"films\":[],\"status\":" +
					"\"PROCESSED\",\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}" +
					",{\"planet_id\":\"dfhtHhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=\",\"name\":\"Yavin IV\",\"climate\"" +
					":[\"temperate\",\"tropical\"],\"terrain\":[\"jungle\",\"rainforests\"],\"films\":[],\"status\"" +
					":\"PROCESSED\",\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}]}",
			},
		},
		{
			name: "success find by name",
			in: in{
				page:    "1",
				perPage: "10",
				name:    "Tatooine",
				reader: mock.PlanetRepository{
					DocsReadAll: []planet.PlanetDocument{
						{
							PlanetID: "PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=",
							Name:     "Tatooine",
							Climate:  []string{"arid"},
							Terrain:  []string{"desert"},
							Films:    planet.Films{},
							Status:   planet.Processed.String(),
						},
					},
					CountValue: 1,
				},
			},
			out: out{
				statusCode: http.StatusOK,
				resBody:    "{\"_metadata\":{\"page\":1,\"per_page\":10,\"page_count\":1,\"total_count\":1},\"results\":[{\"planet_id\":\"PmGs6HhdFVyz01FZgKHoiz_SKxz6B23UH2HZwu5L7BY=\",\"name\":\"Tatooine\",\"climate\":[\"arid\"],\"terrain\":[\"desert\"],\"films\":[],\"status\":\"PROCESSED\",\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}]}",
			},
		},
		{
			name: "success without data",
			in: in{
				page:    "1",
				perPage: "10",
				reader: mock.PlanetRepository{
					DocsReadAll: []planet.PlanetDocument{},
					CountValue:  0,
				},
			},
			out: out{
				statusCode: http.StatusOK,
				resBody:    "{\"_metadata\":{\"page\":1,\"per_page\":10,\"page_count\":0,\"total_count\":0},\"results\":[]}",
			},
		},
		{
			name: "error request without page and perPage query params",
			in: in{
				reader: mock.PlanetRepository{
					DocsReadAll: []planet.PlanetDocument{},
					CountValue:  0,
				},
			},
			out: out{
				statusCode: http.StatusBadRequest,
				resBody:    "{\"params\":{\"page\":\"the page query param cannot be zero or negative\",\"per_page\":\"the per_page query param cannot be zero or negative\"}}",
			},
		},
		{
			name: "error to read all planets",
			in: in{
				page:    "1",
				perPage: "10",
				reader: mock.PlanetRepository{
					ReadAllErr: planet.ErrUnknown,
				},
			},
			out: out{
				statusCode: http.StatusUnprocessableEntity,
				resBody:    "{\"code\":\"PNT-001\",\"message\":\"There was an unknown error processing your request.\"}",
			},
		},
		{
			name: "error to count planets",
			in: in{
				page:    "1",
				perPage: "10",
				reader: mock.PlanetRepository{
					DocsReadAll: []planet.PlanetDocument{},
					CountErr:    planet.ErrUnknown,
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
			router := setupListPlanetHandle(c.in.reader)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v1/planets", nil)
			params := url.Values{}
			params.Set("page", c.in.page)
			params.Set("per_page", c.in.perPage)
			params.Set("name", c.in.name)
			req.URL.RawQuery = params.Encode()

			router.ServeHTTP(w, req)

			assert.Equal(t, c.out.statusCode, w.Code)
			assert.Equal(t, c.out.resBody, w.Body.String())
		})
	}
}

func setupListPlanetHandle(reader planet.Reader) *gin.Engine {
	r := gin.Default()
	planetById := NewListPlanetHandler(reader)

	v1 := r.Group("/v1")
	{
		v1.GET("/planets", planetById.Handle)
	}

	return r
}
