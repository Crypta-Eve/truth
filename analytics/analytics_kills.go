package analytics

import (
	"fmt"

	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
)

//AggregateLossCountAnalysis returns a sorted slice of Pair items where the Key is a Entity name and the key is the number killed for the entity type
func AggregateKillCountAnalysis(aggregate string, entityType string, entityID int, c *client.Client) (counts PairList, err error) {

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

	var ids []int

	for k, _ := range idCount {
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
