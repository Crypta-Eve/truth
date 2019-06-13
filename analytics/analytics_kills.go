package analytics

import (
	"fmt"
	"strconv"

	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
)

//AggregateLossCountAnalysis returns a sorted slice of Pair items where the Key is a Entity name and the key is the number killed for the entity type
func AggregateKilledCountAnalysis(aggregate string, entityType string, entityID int, c *client.Client) (counts PairList, err error) {

	filterField := ""

	switch entityType {
	case "alliance":
		filterField = "killmail.attackers.alliance_id"
	case "corporation":
		filterField = "killmail.attackers.corporation_id"
	case "character":
		filterField = "killmail.attackers.character_id"
	default:
		return counts, errors.New("Invalid entityType got through.... This shouldnt happen")
	}

	filter := bson.M{
		filterField: entityID,
	}

	mails, err := c.Store.GetData(filter)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to query database for filter - %v", filter))
		return counts, err
	}

	c.Log.Printf("Found %v killmails to aggregate", len(mails))

	idCount := make(map[int]int)

	if aggregate == "hour" {
		for i := 0; i < 24; i++ {
			idCount[i] = 0
		}
	}

	for _, mail := range mails {
		switch aggregate {
		case "corporation":
			idCount[mail.KillData.Victim.CorporationID]++
		case "character":
			idCount[mail.KillData.Victim.CharacterID]++
		case "ship":
			idCount[mail.KillData.Victim.ShipTypeID]++
		case "system":
			idCount[mail.KillData.SolarSystemID]++
		case "hour":
			idCount[mail.KillData.KillmailTime.Hour()]++
		default:
			return counts, errors.New("Invalid aggregate (shouldnt have got here), options are - corporation, character, ship, system")
		}

	}

	if aggregate == "hour" {
		counts = make(PairList, len(idCount))

		for k, v := range idCount{
			counts[k] = Pair{Key:strconv.Itoa(k), Value:v}
		}

		return counts, nil
	}

	var ids []int

	for k := range idCount {
		ids = append(ids, k)
	}

	names, err := c.ResolveIDsToNames(ids)
	if err != nil {
		return counts, cli.NewExitError(err, 1)
	}

	out := make(map[string]int)

	for id, count := range idCount {
		out[names[id]] = count
	}

	counts = rankByValue(out)

	return counts, nil
}

func AggregateKillsCountAnalysis(aggregate string, entityType string, entityID int, c *client.Client) (counts PairList, err error) {

	filterField := ""

	switch entityType {
	case "alliance":
		filterField = "killmail.attackers.alliance_id"
	case "corporation":
		filterField = "killmail.attackers.corporation_id"
	case "character":
		filterField = "killmail.attackers.character_id"
	default:
		return counts, errors.New("Invalid entityType got through.... This shouldnt happen")
	}

	filter := bson.M{
		filterField: entityID,
	}

	mails, err := c.Store.GetData(filter)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to query database for filter - %v", filter))
		return counts, err
	}

	c.Log.Printf("Found %v killmails to aggregate", len(mails))

	//I am not happy about what follows, but am undecided on a better course of action

	type void struct{}
	var setter void

	idCount := make(map[int]int)

	if aggregate == "hour" {
		for i := 0; i < 24; i++ {
			idCount[i] = 0
		}
	}

	for _, mail := range mails {
		switch aggregate {
		case "corporation":
			temp := make(map[int]void)

			for _, attacker := range mail.KillData.Attackers {
				switch entityType {
				case "alliance":
					if attacker.AllianceID == entityID {
						temp[attacker.CorporationID] = setter
					}
				case "corporation":
					if attacker.CorporationID == entityID {
						temp[attacker.CorporationID] = setter
					}
				case "character":
					if attacker.CharacterID == entityID {
						temp[attacker.CorporationID] = setter
					}
				}
			}

			for id := range temp {
				idCount[id]++
			}

		case "character":
			temp := make(map[int]void)

			for _, attacker := range mail.KillData.Attackers {
				switch entityType {
				case "alliance":
					if attacker.AllianceID == entityID {
						temp[attacker.CharacterID] = setter
					}
				case "corporation":
					if attacker.CorporationID == entityID {
						temp[attacker.CharacterID] = setter
					}
				case "character":
					if attacker.CharacterID == entityID {
						temp[attacker.CharacterID] = setter
					}
				}
			}

			for id := range temp {
				idCount[id]++
			}
		case "ship":
			for _, attacker := range mail.KillData.Attackers {
				switch entityType {
				case "alliance":
					if attacker.AllianceID == entityID {
						idCount[attacker.ShipTypeID]++
					}
				case "corporation":
					if attacker.CorporationID == entityID {
						idCount[attacker.ShipTypeID]++
					}
				case "character":
					if attacker.CharacterID == entityID {
						idCount[attacker.ShipTypeID]++
					}
				}
			}
		case "system":
			idCount[mail.KillData.SolarSystemID]++
		case "hour":
			idCount[mail.KillData.KillmailTime.Hour()]++
		default:
			return counts, errors.New("Invalid aggregate (shouldnt have got here), options are - corporation, character, ship, system")
		}

	}

	if aggregate == "hour" {
		counts = make(PairList, len(idCount))

		for k, v := range idCount{
			counts[k] = Pair{Key:strconv.Itoa(k), Value:v}
		}

		return counts, nil
	}

	var ids []int

	for k := range idCount {
		ids = append(ids, k)
	}

	names, err := c.ResolveIDsToNames(ids)
	if err != nil {
		return counts, cli.NewExitError(err, 1)
	}

	out := make(map[string]int)

	for id, count := range idCount {
		out[names[id]] = count
	}

	counts = rankByValue(out)

	return counts, nil
}
