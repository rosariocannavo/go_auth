package db

import (
	"context"
	"log"

	"github.com/rosariocannavo/go_auth/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() error {
	clientOptions := options.Client().ApplyURI(config.MongoUri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	Client = client
	return nil
}

func CloseDB() {
	if Client != nil {
		if err := Client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
}
