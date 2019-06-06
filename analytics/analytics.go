package analytics

import (
	"sort"
	"strconv"

	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type (
	Pair struct {
		Key   string
		Value int
	}

	PairList []Pair
)

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//rankByValue adapted from post by Andrew Gerrand - https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
func rankByValue(mip map[string]int) PairList {
	pl := make(PairList, len(mip))
	i := 0
	for k, v := range mip {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func ShipLosses(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	allID := 0

	if argnum == 0 {
		client.Log.Fatal("No alliance id given")
	} else if argnum == 1 {
		arg := c.Args().First()
		allID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	mails, err := client.Store.GetAllianceLosses(allID)

	if err != nil {
		err = errors.Wrap(err, "Failed to get ship losses for alliance")
		return cli.NewExitError(err, 1)
	}

	shipIDCount := make(map[int]int)
	for _, mail := range mails {
		shipIDCount[mail.Victim.ShipTypeID]++
	}

	var ids []int

	for k, _ := range shipIDCount {
		ids = append(ids, k)
	}

	names, err := client.ResolveIDsToNames(ids)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	out := make(map[string]int)

	for id, count := range shipIDCount {
		out[names[id]] = count
	}

	output := rankByValue(out)

	for _, loss := range output {
		client.Log.Printf("%v", loss)
	}

	return nil

}

func PilotLosses(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	corpID := 0

	if argnum == 0 {
		client.Log.Fatal("No alliance id given")
	} else if argnum == 1 {
		arg := c.Args().First()
		corpID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	mails, err := client.Store.GetAllianceLosses(corpID)

	if err != nil {
		err = errors.Wrap(err, "Failed to get ship losses for alliance")
		return cli.NewExitError(err, 1)
	}

	pilotIDCount := make(map[int]int)
	for _, mail := range mails {
		pilotIDCount[mail.Victim.CharacterID]++
	}

	var ids []int

	for k, _ := range pilotIDCount {
		ids = append(ids, k)
	}

	names, err := client.ResolveIDsToNames(ids)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	out := make(map[string]int)

	for id, count := range pilotIDCount {
		out[names[id]] = count
	}

	output := rankByValue(out)

	for _, loss := range output {
		client.Log.Printf("%v", loss)
	}

	return nil

}

func CorpLosses(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	corpID := 0

	if argnum == 0 {
		client.Log.Fatal("No alliance id given")
	} else if argnum == 1 {
		arg := c.Args().First()
		corpID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	mails, err := client.Store.GetAllianceLosses(corpID)

	if err != nil {
		err = errors.Wrap(err, "Failed to get ship losses for alliance")
		return cli.NewExitError(err, 1)
	}

	corpIDCount := make(map[int]int)
	for _, mail := range mails {
		corpIDCount[mail.Victim.CorporationID]++
	}

	var ids []int

	for k, _ := range corpIDCount {
		ids = append(ids, k)
	}

	names, err := client.ResolveIDsToNames(ids)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	out := make(map[string]int)

	for id, count := range corpIDCount {
		out[names[id]] = count
	}

	output := rankByValue(out)

	for _, loss := range output {
		client.Log.Printf("%v", loss)
	}

	return nil

}

func LocationLosses(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	corpID := 0

	if argnum == 0 {
		client.Log.Fatal("No alliance id given")
	} else if argnum == 1 {
		arg := c.Args().First()
		corpID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	mails, err := client.Store.GetAllianceLosses(corpID)

	if err != nil {
		err = errors.Wrap(err, "Failed to get ship losses for alliance")
		return cli.NewExitError(err, 1)
	}

	locationIDCount := make(map[int]int)
	for _, mail := range mails {
		locationIDCount[mail.SolarSystemID]++
	}

	var ids []int

	for k, _ := range locationIDCount {
		ids = append(ids, k)
	}

	names, err := client.ResolveIDsToNames(ids)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	out := make(map[string]int)

	for id, count := range locationIDCount {
		out[names[id]] = count
	}

	output := rankByValue(out)

	for _, loss := range output {
		client.Log.Printf("%v", loss)
	}

	return nil

}
