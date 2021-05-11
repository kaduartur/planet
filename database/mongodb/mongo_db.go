package mongodb

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

const uriPattern = "mongodb://%s:%s@%s:%s"

var (
	dbHost = os.Getenv("DATABASE_HOST")
	dbPort = os.Getenv("DATABASE_PORT")
	dbUser = os.Getenv("DATABASE_USER")
	dbPass = os.Getenv("DATABASE_PASS")
	dbName = os.Getenv("DATABASE_NAME")
)

func NewConnection() {
	opts := options.Client().ApplyURI(uri())
	err := mgm.SetDefaultConfig(nil, dbName, opts)
	if err != nil {
		log.Fatal(err)
	}
}

func uri() string {
	return fmt.Sprintf(uriPattern, dbUser, dbPass, dbHost, dbPort)
}
