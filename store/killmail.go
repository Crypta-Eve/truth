package store

import (
	"context"

	"github.com/pkg/errors"
)

func (db *DB) InsertKillIDHash(idhash ScrapeQueue) error {

	collection := db.Database.Database("truth").Collection("hashes")

	_, err := collection.InsertOne(context.TODO(), idhash)
	if err != nil {
		return errors.Wrap(err, "failed to insert scrape")
	}

	return nil
}

func (db *DB) InsertKillmail(id int, data string) error {

	return errors.New("NYI")
}
