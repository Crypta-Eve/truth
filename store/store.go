package store

import "time"

type (
	Store interface {
		InsertKillmail(id int, data string) error

		InsertKillIDHash(idhash ScrapeQueue) error

		InsertQueueJob(job Queue) error
		PopQueueJob() (job Queue, err error)
		MaintainJobReservation(job Queue, time time.Time) error
		MarkJobComplete(job Queue) error
	}

	Queue struct {
		ID          int       `json:"id" bson:"id"`
		Args        string    `json:"args" bson:"args"`
		Attempts    int       `json:"attempts"  bson:"attempts"`
		ReservedAt  time.Time `json:"reserved_at" bson:"reservedAt"`
		AvailableAt time.Time `json:"available_at" bson:"availableAt"`
		CreatedAt   time.Time `json:"created_at" bson:"createdAt"`
		Complete    bool      `json:"complete" bson:"complete"`
	}

	ScrapeQueue struct {
		ID   int    `json:"_id" bson:"_id"`
		Hash string `json:"hash" bson:"hash"`
	}

	KillmailData struct {
		KillID   int    `json:"_id"`
		KillData string `json:"killmail"`
	}
)

//List of all possible jobs
const (
	JobScrapeCharacter   = 1
	JobScrapeCorporation = 2
	JobScrapeAlliance    = 3
)
