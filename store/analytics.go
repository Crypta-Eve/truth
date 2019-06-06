package store

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *DB) GetAllianceLosses(allianceID int) (results []ESIKillmail, err error) {

	collection := db.Database.Database("truth").Collection("killmails")

	filter := bson.M{
		"killmail.killmail_time": bson.M{
			"$gte": time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		"killmail.victim.alliance_id": allianceID,
	}

	ctx := context.TODO()

	c, err := collection.Find(ctx, filter)
	defer c.Close(ctx)

	if err != nil {
		return results, errors.Wrap(err, "Failed to query alliance losses from db")
	}

	for c.Next(ctx) {

		var data KillmailData

		err := c.Decode(&data)
		if err != nil {
			return results, errors.Wrap(err, "Failed to decode the ship counts")
		}

		results = append(results, data.KillData)
	}

	return results, nil
}
