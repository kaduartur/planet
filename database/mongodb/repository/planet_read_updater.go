package repository

import "github.com/kaduartur/planet"

type PlanetReadUpdater struct {
	planet.Reader
	planet.Updater
}

func NewReadUpdater(r planet.Reader, u planet.Updater) PlanetReadUpdater {
	return PlanetReadUpdater{r, u}
}
