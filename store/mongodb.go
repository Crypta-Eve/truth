package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Database *mongo.Client
}

func SetupStore(envDB map[string]string) (s Store, err error) {

	clientOptions := options.Client().ApplyURI("mongodb://" + envDB["host"] + ":" + envDB["port"])

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return &DB{Database: client}, nil
}
