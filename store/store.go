package store

import (
	"time"
)

type (
	Store interface {
		// Insert new data
		InsertKillmail(kill KillmailData) error
		InsertKillIDHash(idhash ScrapeQueue) error

		// Get specific data
		GetKillmail(id int) (mail KillmailData, err error)

		//Update data
		UpdateKillmail(filter interface{}, update interface{}) error

		//RawFindRequest
		GetData(filter interface{}) (results []KillmailData, err error)

		// Queue management
		InsertQueueJob(job Queue) error
		PopQueueJob() (job Queue, err error)
		MaintainJobReservation(job Queue, time time.Time) error
		MarkJobComplete(job Queue) error

		// Killmail Maint
		ListAllExistingIDs() (ids []int, err error)
		ListAllExistingKillmails() (mails []KillmailData, err error)
		ListMissingKillmails() (mails []ScrapeQueue, err error)
		GetKillsMissingZKB() (hashes []ScrapeQueue, err error)
		GetKillsMissingAxiom() (mails []KillmailData, err error)
		UpdateManyHashes(hashes []ScrapeQueue) error

		// SetupTasks
		ListDatabaseNames() ([]string, error)
		SeedDB() error
		AddIndexes() error

		// StaticData
		DeleteStaticData() error
		InsertRegion(region ESIRegion) error
		InsertConstellation(cons ESIConstellation) error
		InsertSystem(system ESISystem) error
		InsertStar(star ESIStar) error
		InsertPlanet(planet ESIPlanet) error
		InsertMoon(moon ESIMoon) error
		InsertAsteroidBelt(belt ESIAsteroidBelt) error
		InsertStargate(gate ESIStargate) error
		InsertStation(station ESIStation) error
		InsertType(typeESI ESIType) error
		InsertGroup(group ESIGroup) error
		InsertCategory(category ESICategory) error

		GetSystems() (systems []ESISystem, err error)
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
		KillID         int               `json:"_id" bson:"_id"`
		KillData       ESIKillmail       `json:"killmail" bson:"killmail"`
		ZKBData        ZKBData           `json:"zkb,omitempty" bson:"zkb,omitempty"`
		AxiomAttribute FittingAttributes `json:"axiom,omitempty" bson:"axiom,omitempty"`
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

	WSSKillmail struct {
		Attackers     []ESIAttacker `json:"attackers"`
		KillmailID    int           `json:"killmail_id"`
		KillmailTime  time.Time     `json:"killmail_time"`
		SolarSystemID int           `json:"solar_system_id"`
		Victim        ESIVictim     `json:"victim"`
		ZKB           ZKBData       `json:"zkb"`
	}

	ESIRegion struct {
		Constellations []int  `json:"constellations" bson:"constellations"`
		Description    string `json:"description" bson:"description"`
		Name           string `json:"name" bson:"name"`
		RegionID       int    `json:"region_id" bson:"_id"`
	}

	ESIConstellation struct {
		ConstellationID int         `json:"constellation_id" bson:"_id"`
		Name            string      `json:"name" bson:"name"`
		Systems         []int       `json:"systems" bson:"systems"`
		Postion         ESIPosition `json:"position" bson:"position"`
		RegionID        int         `json:"region_id" bson:"region_id"`
	}

	ESISystem struct {
		ConstellationID int                `json:"constellation_id" bson:"constellation_id"`
		Name            string             `json:"name" bson:"name"`
		Planets         []ESISystemPlanets `json:"planets" bson:"planets"`
		Position        ESIPosition        `json:"position" bson:"position"`
		SecurityClass   string             `json:"security_class" bson:"security_class"`
		SecurityStatus  float64            `json:"security_status" bson:"security_status"`
		StarID          int                `json:"star_id" bson:"star_id"`
		Stargates       []int              `json:"stargates" bson:"stargates"`
		Stations        []int              `json:"stations" bson:"stations"`
		SystemID        int                `json:"system_id" bson:"_id"`
	}

	ESISystemPlanets struct {
		PlanetID      int   `json:"planet_id" bson:"planet_id"`
		Moons         []int `json:"moons" bson:"moons,omitempty"`
		AsteroidBelts []int `json:"asteroid_belts" bson:"asteroid_belts,omitempty"`
	}

	ESIStar struct {
		Age           int64   `json:"age" bson:"age"`
		Luminosity    float64 `json:"luminosity" bson:"luminosity"`
		Name          string  `json:"name" bson:"name"`
		Radius        int64   `json:"radius" bson:"radius"`
		SolarSystemID int     `json:"solar_system_id" bson:"solar_system_id"`
		SpectralClass string  `json:"spectral_class" bson:"spectral_class"`
		Temperature   int     `json:"temperature" bson:"temperature"`
		TypeID        int     `json:"type_id" bson:"type_id"`
		StarID        int     `json:"star_id,omitempty" bson:"_id"`
	}

	ESIPlanet struct {
		Name     string      `json:"name" bson:"name"`
		PlanetID int32       `json:"planet_id" bson:"_id"`
		Position ESIPosition `json:"position" bson:"position"`
		SystemID int32       `json:"system_id" bson:"system_id"`
		TypeID   int32       `json:"type_id" bson:"type_id"`
	}

	ESIMoon struct {
		MoonID   int32       `json:"moon_id" bson:"_id"`
		Name     string      `json:"name" bson:"name"`
		Position ESIPosition `json:"position" bson:"position"`
		SystemID int32       `json:"system_id" bson:"system_id"`
	}

	ESIAsteroidBelt struct {
		BeltID   int32       `json:"belt_id,omitempty" bson:"_id"`
		Name     string      `json:"name" bson:"name"`
		Position ESIPosition `json:"position" bson:"position"`
		SystemID int32       `json:"system_id" bson:"system_id"`
	}

	ESIStargate struct {
		Destination ESIStargateDestination `json:"destination" bson:"destination"`
		Name        string                 `json:"name" bson:"name"`
		Position    ESIPosition            `json:"position" bson:"position"`
		StargateID  int32                  `json:"stargate_id" bson:"_id"`
		SystemID    int32                  `json:"system_id" bson:"system_id"`
		TypeID      int32                  `json:"type_id" bson:"type_id"`
	}

	ESIStargateDestination struct {
		StargateID int32 `json:"stargate_id" bson:"stargate_id"`
		SystemID   int32 `json:"system_id" bson:"system_id"`
	}

	ESIStation struct {
		MaxDockableShipVolume  float64     `json:"max_dockable_ship_volume" bson:"max_dockable_ship_volume"`
		Name                   string      `json:"name" bson:"name"`
		OfficeRentalCost       float64     `json:"office_rental_cost" bson:"office_rental_cost"`
		Owner                  int32       `json:"owner" bson:"owner"`
		Position               ESIPosition `json:"position" bson:"position"`
		RaceID                 int32       `json:"race_id" bson:"race_id"`
		ReprocessingEfficiency float32     `json:"reprocessing_efficiency" bson:"reprocessing_efficiency"`
		Services               []string    `json:"services" bson:"services"`
		StationID              int32       `json:"station_id" bson:"_id"`
		SystemID               int32       `json:"system_id" bson:"system_id"`
		TypeID                 int32       `json:"type_id" bson:"type_id"`
	}

	ESIType struct {
		Capacity        float64              `json:"capacity,omitempty" bson:"capacity,omitempty"`
		Description     string               `json:"description" bson:"description"`
		DogmaAttributes []TypeDogmaAttribute `json:"dogma_attributes,omitempty" bson:"dogma_attributes,omitempty"`
		DogmaEffects    []TypeDogmaEffect    `json:"dogma_effects,omitempty" bson:"dogma_effects,omitempty"`
		GraphicID       int32                `json:"graphic_id,omitempty" bson:"graphic_id,omitempty"`
		GroupID         int32                `json:"group_id" bson:"group_id"`
		IconID          int32                `json:"icon_id,omitempty" bson:"icon_id,omitempty"`
		MarketGroupID   int32                `json:"market_group_id,omitempty" bson:"market_group_id,omitempty"`
		Mass            float64              `json:"mass,omitempty" bson:"mass,omitempty"`
		Name            string               `json:"name" bson:"name"`
		PackagedVolume  float64              `json:"packaged_volume,omitempty" bson:"packaged_volume,omitempty"`
		PortionSize     int32                `json:"portion_size,omitempty" bson:"portion_size,omitempty"`
		Published       bool                 `json:"published" bson:"published"`
		Radius          float64              `json:"radius,omitempty" bson:"radius,omitempty"`
		TypeID          int32                `json:"type_id" bson:"_id"`
		Volume          float64              `json:"volume,omitempty" bson:"volume,omitempty"`
	}

	TypeDogmaAttribute struct {
		AttributeID int32   `json:"attribute_id" bson:"attribute_id"`
		Value       float64 `json:"value" bson:"value"`
	}

	TypeDogmaEffect struct {
		EffectID  int32 `json:"effect_id" bson:"effect_id"`
		IsDefault bool  `json:"is_default" bson:"is_default"`
	}

	ESIGroup struct {
		CategoryID int32   `json:"category_id" bson:"category_id"`
		GroupID    int32   `json:"group_id" bson:"_id"`
		Name       string  `json:"name" bson:"name"`
		Published  bool    `json:"published" bson:"published"`
		Types      []int32 `json:"types" bson:"types"`
	}

	ESICategory struct {
		CategoryID int32   `json:"category_id" bson:"_id"`
		Groups     []int32 `json:"groups" bson:"groups"`
		Name       string  `json:"name" bson:"name"`
		Published  bool    `json:"published" bson:"published"`
	}

	FittingAttributes struct {
		Ship   map[string]float64   `json:"ship,omitempty" bson:"ship,omitempty"`
		Drones []map[string]float64 `json:"drones,omitempty" bson:"drones,omitempty"`
	}
)

//List of all possible jobs
const (
	JobSeedDatabase      = 13270
	JobScrapeCharacter   = 1
	JobScrapeCorporation = 2
	JobScrapeAlliance    = 3
	JobScrapeDaily       = 4
)
