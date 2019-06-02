package store

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *DB) InsertQueueJob(job Queue) error {

	collection := db.Database.Database("truth").Collection("jobqueue")

	_, err := collection.InsertOne(context.TODO(), job)
	if err != nil {
		return errors.Wrap(err, "failed to insert job")
	}

	return nil
}

func (db *DB) PopQueueJob() (job Queue, err error) {

	collection := db.Database.Database("truth").Collection("jobqueue")

	//Build our valid job filter
	//mongo -> {complete:false, availableat: {$lte: new Date() }}

	filter := bson.M{
		"complete": false,
		"availableat": bson.M{
			"$lte": time.Now(),
		},
	}

	update := bson.M{
		"$inc": bson.D{
			{"attempts", 1},
		},
		"$set": bson.M{
			"reservedat":  time.Now(),
			"availableat": time.Now().Add(60 * time.Second),
		},
	}

	result := collection.FindOneAndUpdate(context.TODO(), filter, update)

	err = result.Decode(&job)

	if err != nil {
		return Queue{}, err
	}

	return job, nil
}
