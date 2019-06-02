package store

import "time"

type (
	Store interface {
		InsertKillmail(id int, data string) error

		InsertQueueJob(job Queue) error
		PopQueueJob() (job Queue, err error)
	}

	Queue struct {
		ID          int       `json:"id"`
		Args        string    `json:"args"`
		Attempts    int       `json:"attempts"`
		ReservedAt  time.Time `json:"reserved_at"`
		AvailableAt time.Time `json:"available_at"`
		CreatedAt   time.Time `json:"created_at"`
		Complete    bool      `json:"complete"`
	}

	ScrapeQueue struct {
		ID   int    `json:"_id"`
		Hash string `json:"hash"`
	}

	KillmailData struct {
		KillID   int    `json:"_id"`
		KillData string `json:"killmail"`
	}
)

//List of all
const (
	JobScrapeCharacter   = 1
	JobScrapeCorporation = 2
	JobScrapeAlliance    = 3
)
