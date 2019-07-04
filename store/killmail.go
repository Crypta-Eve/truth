package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	collection := db.Database.Database("truth").Collection("hashes")

	ctx := context.Background()

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	rand := fmt.Sprintf("temp_%s", string(b))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$lookup", bson.D{
				{"from", "killmails"},
				{"localField", "_id"},
				{"foreignField", "_id"},
				{"as", "mail"},
			}}},
		bson.D{
			{"$unwind", bson.D{
				{"path", "$mail"},
				{"preserveNullAndEmptyArrays", true},
			}}},
		bson.D{
			{"$match", bson.D{
				{"mail", bson.D{
					{"$exists", false},
				}},
			}}},
		bson.D{
			{"$out", rand},
		},
	}

	startA := time.Now()
	// options.Aggregate().SetBatchSize(500000)
	c, err := collection.Aggregate(ctx, pipeline, options.Aggregate().SetBypassDocumentValidation(true))
	if err != nil {
		return mails, errors.Wrap(err, "error retrieving existing hashes")
	}

	fmt.Printf("Aggregate took %vs\n", time.Now().Sub(startA).Seconds())

	defer c.Close(ctx)

	// start := time.Now()

	// -----------------------------------------------------------------------

	c2 := db.Database.Database("truth").Collection(rand)

	ctx2 := context.Background()

	c3, err := c2.Find(ctx2, bson.M{})
	if err != nil {
		return mails, errors.Wrap(err, "error retrieving existing hashes")
	}

	defer c.Close(ctx)

	start := time.Now()

	for c3.Next(ctx) {

		var mail = ScrapeQueue{}

		err := c3.Decode(&mail)
		if err != nil {
			return mails, errors.Wrap(err, "Failed to morp killmail into struct")
		}

		mails = append(mails, mail)
	}

	fmt.Printf("Cursor took %vs\n", time.Now().Sub(start).Seconds())
	return mails, nil

	//------------------------------------------------------------------------

	// outputChan := make(chan ScrapeQueue, 50)
	// var waitgroup sync.WaitGroup
	// hashes := []ScrapeQueue{}

	// for i := 0; i < 1; i++ {
	// 	//Start up the go routines
	// 	waitgroup.Add(1)
	// 	i1 := i
	// 	go func() {
	// 		for j := 0; j < 10; j++ {
	// 			for c.Next(context.Background()) {
	// 				hsh := ScrapeQueue{}
	// 				err := c.Decode(&hsh)
	// 				if err == nil {
	// 					outputChan <- hsh
	// 				}
	// 				fmt.Print(i1)
	// 				// time.Sleep(500 * time.Nanosecond)
	// 			}
	// 			time.Sleep(15 * time.Second)
	// 		}
	// 		waitgroup.Done()
	// 		fmt.Printf("%v done\n", i1)
	// 	}()
	// }

	// go func() {
	// 	for hash := range outputChan {
	// 		hashes = append(hashes, hash)
	// 		// fmt.Print(".")
	// 	}
	// }()

	// waitgroup.Wait()

	// close(outputChan)

	// // for c.Next(ctx) {

	// // 	var hash = ScrapeQueue{}

	// // 	err := c.Decode(&hash)
	// // 	if err != nil {
	// // 		return mails, errors.Wrap(err, "Failed to morp idhash into struct")
	// // 	}

	// // 	mails = append(mails, hash)
	// // }

	// fmt.Printf("Cursor took %vs\n", time.Now().Sub(start).Seconds())
	// return hashes, nil
}

func (db *DB) GetKillsMissingZKB() (hashes []ScrapeQueue, err error) {
	collection := db.Database.Database("truth").Collection("killmails")

	filter := bson.M{
		"$or": []bson.M{
			{
				"zkb": bson.M{
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
