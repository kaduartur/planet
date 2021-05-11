package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kaduartur/planet"
	"net/http"
)

var _ planet.HttpHandler = ListPlanetHandler{}

type ListPlanetHandler struct {
	planet planet.Reader
}

func NewListPlanetHandler(planet planet.Reader) ListPlanetHandler {
	return ListPlanetHandler{planet: planet}
}

func (d ListPlanetHandler) Handle(c *gin.Context) {
	var pageReq planet.PageFilterRequest
	if err := c.BindQuery(&pageReq); err != nil {
		return
	}

	validate := pageReq.Validate()
	if len(validate) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"params": validate})
		return
	}

	pd, err := d.planet.ReadAll(pageReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	total, err := d.planet.Count()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	pg := planet.Pagination{
		Meta: planet.Metadata{
			Page:       pageReq.Page,
			PerPage:    pageReq.PerPage,
			PageCount:  len(pd),
			TotalCount: total,
		},
		Results: pd,
	}

	c.JSON(http.StatusOK, pg)
}
