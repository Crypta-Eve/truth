package client

import (
	"log"
	"time"

	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/Crypta-Eve/truth/store"
	"github.com/go-redis/redis"
)

type (
	Client struct {
		HTTP      *http.Client
		Store     store.Store
		Log       *log.Logger
		UserAgent string
	}
)

func New() (*Client, error) {
	logger := log.New(os.Stdout, "CLIENT:", log.Lshortfile|log.Ldate|log.Ltime)

	envRedis := viper.GetStringMapString("redis")

	rclient := redis.NewClient(&redis.Options{
		Addr:     envRedis["host"] + ":" + envRedis["port"],
		Password: envRedis["password"],
		DB:       0,
	})

	_, err := rclient.Ping().Result()
	if err != nil {
		return nil, err
	}

	// now check we have access to mongo

	envDB := viper.GetStringMapString("db")
	store, err := store.SetupStore(envDB)

	if err != nil {
		return nil, err
	}

	return &Client{
		HTTP: &http.Client{
			Timeout: time.Second * 30,
		},
		Store:     store,
		Log:       logger,
		UserAgent: viper.GetString("user_agent"),
	}, nil

}
