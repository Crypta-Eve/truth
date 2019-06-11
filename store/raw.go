package store

import (
	"context"

	"github.com/pkg/errors"
)

func (db *DB) GetData(filter interface{}) (results []KillmailData, err error) {
	collection := db.Database.Database("truth").Collection("killmails")

	ctx := context.TODO()

	c, err := collection.Find(ctx, filter)
	defer c.Close(ctx)

	if err != nil {
		return results, errors.Wrap(err, "Failed to query from db")
	}

	for c.Next(ctx) {

		var data KillmailData

		err := c.Decode(&data)
		if err != nil {
			return results, errors.Wrap(err, "Error decoding killmail from db")
		}

		results = append(results, data)
	}

	return results, nil
}
