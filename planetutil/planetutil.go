package planetutil

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/kaduartur/planet"
	"strings"
)

func GeneratePlanetID(name string) planet.ID {
	hasher := sha256.New()
	hasher.Write([]byte(strings.ToTitle(name)))
	planetID := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return planet.ID(planetID)
}
