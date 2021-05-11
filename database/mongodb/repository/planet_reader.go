package repository

import (
	"github.com/kaduartur/planet"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ planet.Reader = PlanetRead{}

type PlanetRead struct {
	collection *mgm.Collection
}

func NewReader() PlanetRead {
	return PlanetRead{
		collection: mgm.CollectionByName("planets"),
	}
}

func (p PlanetRead) ReadByPlanetId(id planet.ID) (planet.PlanetDocument, error) {
	one := p.collection.FindOne(mgm.Ctx(), bson.M{"planetID": id})
	var pd planet.PlanetDocument
	if err := one.Decode(&pd); err != nil && err != mongo.ErrNoDocuments {
		return planet.PlanetDocument{}, err
	}

	return pd, nil
}

func (p PlanetRead) ReadByName(name string) (planet.PlanetDocument, error) {
	one := p.collection.FindOne(mgm.Ctx(), bson.M{"name": name})
	var pd planet.PlanetDocument
	if err := one.Decode(&pd); err != nil && err != mongo.ErrNoDocuments {
		return planet.PlanetDocument{}, err
	}

	return pd, nil
}

func (p PlanetRead) ReadAll(f planet.PageFilterRequest) ([]planet.PlanetDocument, error) {
	result := make([]planet.PlanetDocument, 0)
	offset := int64((f.Page - 1) * f.PerPage)
	limit := int64(f.PerPage)

	var filter bson.M
	if f.Name != "" {
		filter = bson.M{"name": f.Name}
	}


	opts := options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
	}

	if err := p.collection.SimpleFind(&result, filter, &opts); err != nil {
		return nil, err
	}

	return result, nil
}

func (p PlanetRead) Count() (int, error) {
	total, err := p.collection.CountDocuments(mgm.Ctx(), bson.M{})
	if err != nil {
		return 0, err
	}

	return int(total), err
}
