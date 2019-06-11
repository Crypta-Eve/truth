package store

import (
	"time"
)

type (
	Store interface {
		// Insert new data
		InsertKillmail(kill KillmailData) error
		InsertKillIDHash(idhash ScrapeQueue) error

		//Update data
		UpdateKillmail(filter interface{}, update interface{}) error

		//RawFindRequest
		GetData(filter interface{}) (results []KillmailData, err error)

		// Queue management
		InsertQueueJob(job Queue) error
		PopQueueJob() (job Queue, err error)
		ListAllExistingIDs() (ids []int, err error)
		GetKillsNotInList(existing []int) (hashes []ScrapeQueue, err error)
		GetKillsMissingZKB() (hashes []ScrapeQueue, err error)
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
		KillID   int         `json:"_id" bson:"_id"`
		KillData ESIKillmail `json:"killmail" bson:"killmail"`
		ZKBData  ZKBData     `json:"zkb" bson:"zkb"`
	}

	ESIKillmail struct {
		Attackers     []ESIAttacker `json:"attackers" bson:"attackers"`
		KillmailID    int           `json:"killmail_id" bson:"killmail_id"`
		KillmailTime  time.Time     `json:"killmail_time" bson:"killmail_time"`
		SolarSystemID int           `json:"solar_system_id" bson:"solar_system_id"`
		Victim        ESIVictim     `json:"victim" bson:"victim"`
	}

	ESIAttacker struct {
		AllianceID     int     `json:"alliance_id,omitempty" bson:"alliance_id,omitempty"`
		CorporationID  int     `json:"corporation_id" bson:"corporation_id"`
		CharacterID    int     `json:"character_id" bson:"character_id"`
		DamageDone     int     `json:"damage_done" bson:"damage_done"`
		FinalBlow      bool    `json:"final_blow" bson:"final_blow"`
		SecurityStatus float32 `json:"security_status" bson:"security_status"`
		ShipTypeID     int     `json:"ship_type_id" bson:"ship_type_id"`
		WeaponTypeID   int     `json:"weapon_type_id" bson:"weapon_type_id"`
	}

	ESIVictim struct {
		AllianceID    int         `json:"alliance_id,omitempty" bson:"alliance_id,omitempty"`
		CorporationID int         `json:"corporation_id" bson:"corporation_id"`
		CharacterID   int         `json:"character_id" bson:"character_id"`
		DamageTaken   int         `json:"damage_taken" bson:"damage_taken"`
		Items         []ESIItem   `json:"items" bson:"items"`
		Position      ESIPosition `json:"position" bson:"position"`
		ShipTypeID    int         `json:"ship_type_id" bson:"ship_type_id"`
	}

	ESIItem struct {
		Flag              int `json:"flag" bson:"flag"`
		ItemTypeID        int `json:"item_type_id" bson:"item_type_id"`
		QuantityDropped   int `json:"quantity_dropped,omitempty" bson:"quantity_dropped,omitempty"`
		QuantityDestroyed int `json:"quantity_destroyed,omitempty" bson:"quantity_destroyed,omitempty"`
		Singleton         int `json:"singleton" bson:"singleton"`
	}

	ESIPosition struct {
		X float64 `json:"x" bson:"x"`
		Y float64 `json:"y" bson:"y"`
		Z float64 `json:"z" bson:"z"`
	}

	ZKBResponse struct {
		KillmailID int     `json:"killmail_id"`
		ZKB        ZKBData `json:"zkb"`
	}

	ZKBData struct {
		LocationID  int     `json:"locationID" bson:"location_id"`
		Hash        string  `json:"hash" bson:"hash"`
		FittedValue float64 `json:"fittedValue" bson:"fitted_value"`
		TotalValue  float64 `json:"totalValue" bson:"total_value"`
		Points      int     `json:"points" bson:"points"`
		NPC         bool    `json:"npc" bson:"npc"`
		Solo        bool    `json:"solo" bson:"solo"`
		Awox        bool    `json:"awox" bson:"awox"`
	}
)

//List of all possible jobs
const (
	JobScrapeCharacter   = 1
	JobScrapeCorporation = 2
	JobScrapeAlliance    = 3
)
