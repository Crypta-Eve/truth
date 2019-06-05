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
		"availableAt": bson.M{
			"$lte": time.Now(),
		},
	}

	update := bson.M{
		"$inc": bson.D{
			{"attempts", 1},
		},
		"$set": bson.M{
			"reservedAt":  time.Now(),
			"availableAt": time.Now().Add(60 * time.Second),
		},
	}

	result := collection.FindOneAndUpdate(context.TODO(), filter, update)

	err = result.Decode(&job)

	if err != nil {
		return Queue{}, err
	}

	return job, nil
}


func (db *DB) MaintainJobReservation(job Queue, t time.Time) error {

	collection := db.Database.Database("truth").Collection("jobqueue")

	filter := bson.M{
		"id":   job.ID,
		"args": job.Args,
		// "attempts": job.Attempts, (cant guarentee this one)
		"complete":   job.Complete,
		"reservedAt": job.ReservedAt,
		"createdAt":  job.CreatedAt,
	}

	update := bson.M{
		"$set": bson.M{
			"availableat": t,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func (db *DB) MarkJobComplete(job Queue) error {

	collection := db.Database.Database("truth").Collection("jobqueue")

	filter := bson.M{
		"id":        job.ID,
		"args":      job.Args,
		"complete":  false,
		"createdAt": job.CreatedAt,
	}

	update := bson.M{
		"$set": bson.M{
			"availableAt": time.Now(),
			"complete":    true,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err

}
