package store

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) ListDatabaseNames() ([]string, error) {

	ctx := context.Background()

	existingDB, err := db.Database.ListDatabaseNames(ctx, bson.M{})

	return existingDB, err

}

func (db *DB) SeedDB() error {

	// Need a better method of storing this
	idhash := ScrapeQueue{ID: 30297766, Hash: "687ae7b8619cfec0f8ceb10d641ac28a228e429a"}
	killmail := KillmailData{
		KillID: 30304916,
		KillData: ESIKillmail{
			Attackers: []ESIAttacker{
				{
					CorporationID:  166334732,
					CharacterID:    546867805,
					DamageDone:     609,
					FinalBlow:      true,
					SecurityStatus: -10,
					ShipTypeID:     29986,
					WeaponTypeID:   3520,
				},
			},
			KillmailID:    30304916,
			KillmailTime:  time.Unix(1367785740, 0),
			SolarSystemID: 30004979,
			Victim: ESIVictim{
				AllianceID:    434243723,
				CorporationID: 109299958,
				CharacterID:   92168555,
				DamageTaken:   609,
				Items: []ESIItem{
					{
						Flag:              89,
						ItemTypeID:        3216,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        3195,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        3213,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        19540,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        24343,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        24344,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        24347,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        19551,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        24345,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        19555,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        24346,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        19556,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        19553,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
					{
						Flag:              89,
						ItemTypeID:        19554,
						QuantityDestroyed: 1,
						Singleton:         0,
					},
				},
				Position: ESIPosition{
					X: 0,
					Y: 0,
					Z: 0,
				},
				ShipTypeID: 670,
			},
		},
		ZKBData: ZKBData{
			Hash:        "687ae7b8619cfec0f8ceb10d641ac28a228e429a",
			FittedValue: 4371230780.14,
			TotalValue:  4371230780.14,
			Points:      1,
			NPC:         false,
			Solo:        false,
			Awox:        false,
		},
	}
	job := Queue{
		ID:          JobSeedDatabase,
		Args:        "",
		Attempts:    1,
		CreatedAt:   time.Now(),
		ReservedAt:  time.Now(),
		AvailableAt: time.Now(),
		Complete:    true,
	}

	// Hashes Table
	err := db.InsertKillIDHash(idhash)
	if err != nil {
		if !strings.Contains(err.Error(), "dup key") {
			return errors.Wrap(err, "Failed to seed the hashes collection")
		}
	}

	err = db.InsertKillmail(killmail)
	if err != nil {
		if !strings.Contains(err.Error(), "dup key") {
			return errors.Wrap(err, "Failed to seed the killmail collection")
		}
	}

	err = db.InsertQueueJob(job)
	if err != nil {
		if !strings.Contains(err.Error(), "dup key") {
			return errors.Wrap(err, "Failed to seed the jobqueue collection")
		}
	}
	return nil
}

func (db *DB) AddIndexes() error {

	jobqueueIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"complete", 1},
				{"availableAt", 1},
			},
			Options: options.Index().
				SetName("_jobqueue_avail_jobs"),
		},
	}

	// Commented out until I ad these later......
	// killmailsIndexes := []mongo.IndexModel{}

	_, err := db.Database.Database("truth").Collection("jobqueue").Indexes().CreateMany(context.Background(), jobqueueIndexes)
	if err != nil {
		return errors.Wrap(err, "Failed to index jobqueue collection")
	}

	// Commented out until I ad these later......
	// _, err = db.Database.Database("truth").Collection("killmails").Indexes().CreateMany(context.Background(), killmailsIndexes)
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to index killmails collection")
	// }

	return nil
}
