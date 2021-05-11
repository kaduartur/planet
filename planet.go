package planet

import (
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"time"
)

type HttpHandler interface {
	Handle(c *gin.Context)
}

type KafkaWriter interface {
	Write(topic string, eventType EventType, msg []byte) error
}

type KafkaRead interface {
	Read()
}

type EventProcessor interface {
	Process(data interface{}) error
}

type Writer interface {
	Write(createPlanet CreatePlanetCommand) (PlanetDocument, error)
}

type Deleter interface {
	Delete(id ID) error
}

type Reader interface {
	ReadByPlanetId(id ID) (PlanetDocument, error)
	ReadAll(filter PageFilterRequest) ([]PlanetDocument, error)
	Count() (int, error)
}

type Updater interface {
	Update(planet PlanetDocument) error
}

type ReadUpdater interface {
	Reader
	Updater
}

type Status int

const (
	Processing Status = iota
	Processed
	Failed
)

func (p Status) String() string {
	return [...]string{"PROCESSING", "PROCESSED", "FAIL"}[p]
}

type EventType int

const CreatedEvent EventType = iota

func (e EventType) String() string {
	return [...]string{"CREATED"}[e]
}

type EventsProcessor map[string]EventProcessor

type CreatePlanetEvent struct {
	PlanetID   ID  `json:"planet_id"`
	RetryCount int `json:"retry_count"`
}

type CreatePlanetCommand struct {
	Name    string `json:"name"`
	Climate string `json:"climate"`
	Terrain string `json:"terrain"`
}

func (c *CreatePlanetCommand) Validate() map[string]string {
	validation := make(map[string]string)
	if c.Name == "" {
		validation["name"] = "This field cannot be empty."
	}

	if c.Climate == "" {
		validation["climate"] = "This field cannot be empty."
	}

	if c.Terrain == "" {
		validation["terrain"] = "This field cannot be empty."
	}

	return validation
}

type StateResponse struct {
	PlanetID ID     `json:"planet_id"`
	Status   string `json:"status"`
}

type PageFilterRequest struct {
	Name    string `form:"name"`
	Page    int    `form:"page"`
	PerPage int    `form:"per_page"`
}

func (p *PageFilterRequest) Validate() map[string]string {
	validation := make(map[string]string)
	if p.Page <= 0 {
		validation["page"] = "the page query param cannot be zero or negative"
	}

	if p.PerPage <= 0 {
		validation["per_page"] = "the per_page query param cannot be zero or negative"
	}

	return validation
}

type Metadata struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}

type Pagination struct {
	Meta    Metadata    `json:"_metadata"`
	Results interface{} `json:"results"`
}

type PlanetDocument struct {
	PlanetID       ID       `bson:"planetID" json:"planet_id"`
	Name           string   `bson:"name" json:"name"`
	Climate        []string `bson:"climate" json:"climate"`
	Terrain        []string `bson:"terrain" json:"terrain"`
	Films          Films    `bson:"films" json:"films"`
	Status         string   `bson:"status" json:"status"`
	mgm.IDField    `bson:" ,inline" json:"-"`
	mgm.DateFields `bson:" ,inline"`
}

type ID string

type Films []Film

type Film struct {
	Title       string    `bson:"title" json:"title"`
	ReleaseDate time.Time `bson:"release_date" json:"release_date"`
}
