package analytics

import (
	"fmt"

	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
)

//ShipLossCountAnalysis returns a sorted slice of Pair items where the Key is a Ship/Type name and the key is the number lost by the specified entity

func ShipLossCountAnalysis(entityType string, entityID int, c *client.Client) (counts PairList, err error) {

	filterField := ""

	switch entityType {
	case "alliance":
		filterField = "killmail.victim.alliance_id"
	case "corporation":
		filterField = "kilmail.victim.corporation_id"
	case "character":
		filterField = "killmail.victim.character_id"
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

	shipIDCount := make(map[int]int)
	for _, mail := range mails {
		shipIDCount[mail.KillData.Victim.ShipTypeID]++
	}

	var ids []int

	for k, _ := range shipIDCount {
		ids = append(ids, k)
	}

	names, err := c.ResolveIDsToNames(ids)
	if err != nil {
		return counts, cli.NewExitError(err, 1)
	}

	out := make(map[string]int)

	for id, count := range shipIDCount {
		out[names[id]] = count
	}

	counts = rankByValue(out)

	return counts, nil
}
