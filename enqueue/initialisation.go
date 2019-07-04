package enqueue

import (
	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func PerformSetup(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	// Check if the DB exists, if not doesnt really matter but I like giving feedback
	names, err := client.Store.ListDatabaseNames()
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Error getting database names"), 1)
	}

	dbExists := false
	for _, db := range names {
		if db == "truth" {
			dbExists = true
		}
	}

	if dbExists {
		client.Log.Println("Database 'truth' exists")
	} else {
		client.Log.Println("Database 'truth' does not exist, it will be created")
	}

	// Check for all of our collections

	err = client.Store.SeedDB()
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to seed the database"), 1)
	}

	client.Log.Println("Database has been successfully seeded")

	err = client.Store.AddIndexes()
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to index the database"), 1)
	}

	client.Log.Println("Database has been successfully indexed")

	return nil
}
