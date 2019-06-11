package analytics

import (
	"fmt"
	"strconv"

	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func ShipLosses(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	allID := 0
	entityType := ""

	entityOptions := []string{"character", "corporation", "alliance"}

	if argnum == 2 {

		entityType = c.Args().First()
		valid := false

		for _, v := range entityOptions {
			if v == entityType {
				valid = true
			}
		}

		if valid != true {
			return cli.NewExitError(fmt.Sprintf("Invalid entity type. Valid options are %v", entityOptions), 1)
		}

		arg := c.Args().Get(1)
		allID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Entity ID not an integer"), 1)
		}

	} else if argnum == 4 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	output, err := ShipLossCountAnalysis(entityType, allID, client)
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to perform ship loss analysis"), 1)
	}

	for _, loss := range output {
		client.Log.Printf("%v", loss)
	}

	return nil

}

func PilotLosses(c *cli.Context) error {

	// client, err := client.New()

	// if err != nil {
	// 	err = errors.Wrap(err, "failed to create client")
	// 	return cli.NewExitError(err, 1)
	// }

	// argnum := len(c.Args())

	// corpID := 0

	// if argnum == 0 {
	// 	client.Log.Fatal("No alliance id given")
	// } else if argnum == 1 {
	// 	arg := c.Args().First()
	// 	corpID, err = strconv.Atoi(arg)
	// 	if err != nil {
	// 		return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
	// 	}

	// } else if argnum == 3 {
	// 	client.Log.Print("Dates NYI")
	// } else {
	// 	client.Log.Fatal("invalid number of arguments")
	// }

	// mails, err := client.Store.GetAllianceLosses(corpID)

	// if err != nil {
	// 	err = errors.Wrap(err, "Failed to get ship losses for alliance")
	// 	return cli.NewExitError(err, 1)
	// }

	// pilotIDCount := make(map[int]int)
	// for _, mail := range mails {
	// 	pilotIDCount[mail.Victim.CharacterID]++
	// }

	// var ids []int

	// for k, _ := range pilotIDCount {
	// 	ids = append(ids, k)
	// }

	// names, err := client.ResolveIDsToNames(ids)
	// if err != nil {
	// 	return cli.NewExitError(err, 1)
	// }

	// out := make(map[string]int)

	// for id, count := range pilotIDCount {
	// 	out[names[id]] = count
	// }

	// output := rankByValue(out)

	// for _, loss := range output {
	// 	client.Log.Printf("%v", loss)
	// }

	// err = GenerateSingleValuePlot("Pilot Losses", output)

	// if err != nil {
	// 	return cli.NewExitError(errors.Wrap(err, "Failed to generate plot"), 1)
	// }

	return nil

}

func CorpLosses(c *cli.Context) error {

	// client, err := client.New()

	// if err != nil {
	// 	err = errors.Wrap(err, "failed to create client")
	// 	return cli.NewExitError(err, 1)
	// }

	// argnum := len(c.Args())

	// corpID := 0

	// if argnum == 0 {
	// 	client.Log.Fatal("No alliance id given")
	// } else if argnum == 1 {
	// 	arg := c.Args().First()
	// 	corpID, err = strconv.Atoi(arg)
	// 	if err != nil {
	// 		return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
	// 	}

	// } else if argnum == 3 {
	// 	client.Log.Print("Dates NYI")
	// } else {
	// 	client.Log.Fatal("invalid number of arguments")
	// }

	// mails, err := client.Store.GetAllianceLosses(corpID)

	// if err != nil {
	// 	err = errors.Wrap(err, "Failed to get ship losses for alliance")
	// 	return cli.NewExitError(err, 1)
	// }

	// corpIDCount := make(map[int]int)
	// for _, mail := range mails {
	// 	corpIDCount[mail.Victim.CorporationID]++
	// }

	// var ids []int

	// for k, _ := range corpIDCount {
	// 	ids = append(ids, k)
	// }

	// names, err := client.ResolveIDsToNames(ids)
	// if err != nil {
	// 	return cli.NewExitError(err, 1)
	// }

	// out := make(map[string]int)

	// for id, count := range corpIDCount {
	// 	out[names[id]] = count
	// }

	// output := rankByValue(out)

	// for _, loss := range output {
	// 	client.Log.Printf("%v", loss)
	// }

	return nil

}

func LocationLosses(c *cli.Context) error {

	// client, err := client.New()

	// if err != nil {
	// 	err = errors.Wrap(err, "failed to create client")
	// 	return cli.NewExitError(err, 1)
	// }

	// argnum := len(c.Args())

	// corpID := 0

	// if argnum == 0 {
	// 	client.Log.Fatal("No alliance id given")
	// } else if argnum == 1 {
	// 	arg := c.Args().First()
	// 	corpID, err = strconv.Atoi(arg)
	// 	if err != nil {
	// 		return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
	// 	}

	// } else if argnum == 3 {
	// 	client.Log.Print("Dates NYI")
	// } else {
	// 	client.Log.Fatal("invalid number of arguments")
	// }

	// mails, err := client.Store.GetAllianceLosses(corpID)

	// if err != nil {
	// 	err = errors.Wrap(err, "Failed to get ship losses for alliance")
	// 	return cli.NewExitError(err, 1)
	// }

	// locationIDCount := make(map[int]int)
	// for _, mail := range mails {
	// 	locationIDCount[mail.SolarSystemID]++
	// }

	// var ids []int

	// for k, _ := range locationIDCount {
	// 	ids = append(ids, k)
	// }

	// names, err := client.ResolveIDsToNames(ids)
	// if err != nil {
	// 	return cli.NewExitError(err, 1)
	// }

	// out := make(map[string]int)

	// for id, count := range locationIDCount {
	// 	out[names[id]] = count
	// }

	// output := rankByValue(out)

	// for _, loss := range output {
	// 	client.Log.Printf("%v", loss)
	// }

	return nil

}

func TZLosses(c *cli.Context) error {

	// client, err := client.New()

	// if err != nil {
	// 	err = errors.Wrap(err, "failed to create client")
	// 	return cli.NewExitError(err, 1)
	// }

	// argnum := len(c.Args())

	// corpID := 0

	// if argnum == 0 {
	// 	client.Log.Fatal("No alliance id given")
	// } else if argnum == 1 {
	// 	arg := c.Args().First()
	// 	corpID, err = strconv.Atoi(arg)
	// 	if err != nil {
	// 		return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
	// 	}

	// } else if argnum == 3 {
	// 	client.Log.Print("Dates NYI")
	// } else {
	// 	client.Log.Fatal("invalid number of arguments")
	// }

	// mails, err := client.Store.GetAllianceLosses(corpID)

	// if err != nil {
	// 	err = errors.Wrap(err, "Failed to get ship losses for alliance")
	// 	return cli.NewExitError(err, 1)
	// }

	// tzCount := make(map[string]int)
	// for _, mail := range mails {
	// 	tzCount[strconv.Itoa(mail.KillmailTime.Hour())]++
	// }

	// output := rankByValue(tzCount)

	// for _, loss := range output {
	// 	client.Log.Printf("%v", loss)
	// }

	return nil

}
