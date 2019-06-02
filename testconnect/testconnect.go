package testconnect

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestExternalAppConnections will ensure that we have a connection to all required local services
func TestExternalAppConnections(c *cli.Context) error {

	envRedis := viper.GetStringMapString("redis")

	rclient := redis.NewClient(&redis.Options{
		Addr:     envRedis["host"] + ":" + envRedis["port"],
		Password: envRedis["password"],
		DB:       0,
	})

	_, err := rclient.Ping().Result()
	if err != nil {
		return err
	}

	fmt.Println("Connected to Redis!")

	// now check we have access to mongo

	envDB := viper.GetStringMapString("db")
	clientOptions := options.Client().ApplyURI("mongodb://" + envDB["host"] + ":" + envDB["port"])

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB!")

	err = client.Disconnect(context.TODO())

	if err != nil {
		return err
	}
	fmt.Println("Connection to MongoDB closed.")

	return nil

}
