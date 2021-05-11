package mock

import "github.com/kaduartur/planet"

var _ planet.Updater = PlanetRepository{}

type PlanetRepository struct {
	DocWrite    planet.PlanetDocument
	DocReadById planet.PlanetDocument
	DocsReadAll []planet.PlanetDocument
	CountValue  int
	WriteErr    error
	ReadByIdErr error
	ReadAllErr  error
	CountErr    error
	DeleteErr   error
	UpdateErr   error
}

func (m PlanetRepository) Write(createPlanet planet.CreatePlanetCommand) (planet.PlanetDocument, error) {
	return m.DocWrite, m.WriteErr
}

func (m PlanetRepository) ReadByPlanetId(id planet.ID) (planet.PlanetDocument, error) {
	return m.DocReadById, m.ReadByIdErr
}

func (m PlanetRepository) ReadAll(filter planet.PageFilterRequest) ([]planet.PlanetDocument, error) {
	return m.DocsReadAll, m.ReadAllErr
}

func (m PlanetRepository) Count() (int, error) {
	return m.CountValue, m.CountErr
}

func (m PlanetRepository) Delete(id planet.ID) error {
	return m.DeleteErr
}

func (m PlanetRepository) Update(planet planet.PlanetDocument) error {
	return m.UpdateErr
}
