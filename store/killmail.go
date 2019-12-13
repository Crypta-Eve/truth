package store

import (
	"context"
	"fmt"

	// "math/rand"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) InsertKillIDHash(idhash ScrapeQueue) error {

	collection := db.Database.Database("truth").Collection("killmails")

	filter := bson.M{"_id": idhash.ID}
	update := bson.M{"$set": bson.M{"hash": idhash.Hash}}

	_, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return errors.Wrap(err, "failed to insert kill id hash")
	}

	return nil
}

func (db *DB) InsertKillmail(kill KillmailData) error {

	collection := db.Database.Database("truth").Collection("killmails")
	filter := bson.M{"_id": kill.KillID}
	update := bson.M{"$set": bson.M{"killmail": kill.KillData}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, "Failed to insert killmail to db")
	}

	return nil
}

func (db *DB) GetKillmail(id int) (mail KillmailData, err error) {
	collection := db.Database.Database("truth").Collection("killmails")

	filter := bson.M{
		"_id": id,
	}

	ctx := context.TODO()

	c := collection.FindOne(ctx, filter)

	err = c.Decode(&mail)

	return mail, err
}

func (db *DB) ListAllExistingIDs() (ids []int, err error) {

	collection := db.Database.Database("truth").Collection("killmails")

	ctx := context.Background()

	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return ids, errors.Wrap(err, "error retrieving existing hashes")
	}

	defer c.Close(ctx)

	start := time.Now()

	for c.Next(ctx) {

		type fetchid struct {
			ID int `bson:"_id"`
			// Killmail bson.Raw `bson:"killmail"`
		}

		var id fetchid

		err := c.Decode(&id)
		if err != nil {
			return ids, errors.Wrap(err, "Failed to morp killmail into struct")
		}

		ids = append(ids, id.ID)
	}

	fmt.Printf("Cursor took %vs\n", time.Now().Sub(start).Seconds())
	return ids, nil
}

func (db *DB) ListAllExistingKillmails() (mails []KillmailData, err error) {

	collection := db.Database.Database("truth").Collection("killmails")

	ctx := context.Background()

	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return mails, errors.Wrap(err, "error retrieving existing hashes")
	}

	defer c.Close(ctx)

	start := time.Now()

	for c.Next(ctx) {

		var mail = KillmailData{}

		err := c.Decode(&mail)
		if err != nil {
			return mails, errors.Wrap(err, "Failed to morp killmail into struct")
		}

		mails = append(mails, mail)
	}

	fmt.Printf("Cursor took %vs\n", time.Now().Sub(start).Seconds())
	return mails, nil
}

func (db *DB) ListMissingKillmails() (mails []ScrapeQueue, err error) {

	collection := db.Database.Database("truth").Collection("killmails")

	ctx := context.Background()
	c, err := collection.Find(ctx, bson.M{"killmail.killmail_id": bson.M{"$exists": false}})
	if err != nil {
		return mails, errors.Wrap(err, "error retrieving missing killmails")
	}

	defer c.Close(ctx)

	start := time.Now()

	for c.Next(ctx) {

		var mail = ScrapeQueue{}

		err := c.Decode(&mail)
		if err != nil {
			return mails, errors.Wrap(err, "Failed to morp killmail into struct")
		}

		mails = append(mails, mail)
	}

	fmt.Printf("Cursor took %vs\n", time.Now().Sub(start).Seconds())
	return mails, nil

}

func (db *DB) GetKillsMissingZKB() (hashes []ScrapeQueue, err error) {
	collection := db.Database.Database("truth").Collection("killmails")

	filter := bson.M{
		"$or": []bson.M{
			{
				"zkb.hash": bson.M{
					"$exists": false,
				},
			},
			{
				"zkb.hash": "",
			},
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

func (db *DB) GetKillsMissingAxiom() (mails []KillmailData, err error) {
	collection := db.Database.Database("truth").Collection("killmails")

	filter := bson.M{
		"$or": []bson.M{
			{
				"axiom": bson.M{
					"$exists": false,
				},
			},
			{
				"axiom.hp": bson.M{
					"$exists": false,
				},
			},
		},
	}

	c, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return mails, errors.Wrap(err, "error retrieving killmails missing zkb element")
	}

	defer c.Close(context.TODO())

	for c.Next(context.TODO()) {
		var hsh KillmailData
		err = c.Decode(&hsh)
		if err != nil {
			return mails, err
		}

		mails = append(mails, hsh)
	}

	return mails, nil
}

func (db *DB) UpdateKillmail(filter interface{}, update interface{}) error {
	collection := db.Database.Database("truth").Collection("killmails")

	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func (db *DB) UpdateManyHashes(hashes []ScrapeQueue) error {

	documents := make([]interface{}, len(hashes))
	for i, hash := range hashes {
		documents[i] = hash
	}
	collection := db.Database.Database("truth").Collection("hashes")

	_, err := collection.InsertMany(context.TODO(), documents, options.InsertMany().SetOrdered(false))

	return err
}
