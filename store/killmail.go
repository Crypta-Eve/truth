package store

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *DB) InsertKillIDHash(idhash ScrapeQueue) error {

	collection := db.Database.Database("truth").Collection("hashes")

	_, err := collection.InsertOne(context.TODO(), idhash)
	if err != nil {
		return errors.Wrap(err, "failed to insert kill id hash")
	}

	return nil
}

func (db *DB) InsertKillmail(kill KillmailData) error {

	collection := db.Database.Database("truth").Collection("killmails")

	_, err := collection.InsertOne(context.TODO(), kill)
	if err != nil {
		return errors.Wrap(err, "Failed to insert killmail to db")
	}

	return nil
}

//
func (db *DB) ListAllExistingIDs() (ids []int, err error) {

	collection := db.Database.Database("truth").Collection("killmails")

	ctx := context.Background()

	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return ids, errors.Wrap(err, "error retrieving existing hashes")
	}

	defer c.Close(ctx)

	for c.Next(ctx) {

		type fetchid struct {
			Id       int      `bson:"_id"`
			Killmail bson.Raw `bson:"killmail"`
		}

		var id fetchid

		err := c.Decode(&id)
		if err != nil {
			return ids, errors.Wrap(err, "Failed to morp killmail into struct")
		}

		ids = append(ids, id.Id)
	}

	return ids, nil
}

func (db *DB) GetKillsNotInList(existing []int) (hashes []ScrapeQueue, err error) {

	// {"_id": { $nin : [77029755, 77029736....]}}

	collection := db.Database.Database("truth").Collection("hashes")

	filter := bson.M{
		"_id": bson.M{
			"$nin": existing,
		},
	}

	c, err := collection.Find(context.TODO(), filter)

	if err != nil {
		return hashes, errors.Wrap(err, "error retrieving missing hashes")
	}

	defer c.Close(context.TODO())

	for c.Next(context.TODO()) {

		var hsh ScrapeQueue
		err := c.Decode(&hsh)
		if err != nil {
			return hashes, err
		}

		hashes = append(hashes, hsh)
	}

	return hashes, nil

}

func (db *DB) GetKillsMissingZKB() (hashes []ScrapeQueue, err error) {
	collection := db.Database.Database("truth").Collection("killmails")

	filter := bson.M{
		"zkb": bson.M{
			"$exists": false,
		},
	}

	c, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return hashes, errors.Wrap(err, "error retrieving killmails missing zkb element")
	}

	defer c.Close(context.TODO())

	for c.Next(context.TODO()) {
		var hsh ScrapeQueue
		err = c.Decode(&hsh)
		if err != nil {
			return hashes, err
		}

		hashes = append(hashes, hsh)
	}

	return hashes, nil
}

func (db *DB) UpdateKillmail(filter interface{}, update interface{}) error {
	collection := db.Database.Database("truth").Collection("killmails")

	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err
}
