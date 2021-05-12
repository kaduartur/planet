package mongodb

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
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
	conf := mgm.Config{CtxTimeout: time.Second * 5}
	if err := mgm.SetDefaultConfig(&conf, dbName, opts); err != nil {
		log.Fatalln(err)
	}

	_, client, _, err := mgm.DefaultConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	if err := client.Ping(mgm.Ctx(), readpref.Primary()); err != nil {
		log.Fatalln(err)
	}
}

func uri() string {
	return fmt.Sprintf(uriPattern, dbUser, dbPass, dbHost, dbPort)
}
