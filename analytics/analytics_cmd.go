package analytics

import (
	"fmt"
	"strconv"

	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func AggregateLosses(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	allID := 0
	entityType := ""
	aggregateType := ""

	entityOptions := []string{"character", "corporation", "alliance"}
	aggregateOptions := []string{"corporation", "character", "ship", "system", "hour"}

	if argnum == 3 {

		aggregateType = c.Args().Get(0)
		validAggregate := false

		for _, v := range aggregateOptions {
			if v == aggregateType {
				validAggregate = true
			}
		}

		if validAggregate != true {
			return cli.NewExitError(fmt.Sprintf("Invalid aggregate type. Valid options are %v", aggregateOptions), 1)
		}

		entityType = c.Args().Get(1)
		validEntity := false

		for _, v := range entityOptions {
			if v == entityType {
				validEntity = true
			}
		}

		if validEntity != true {
			return cli.NewExitError(fmt.Sprintf("Invalid entity type. Valid options are %v", entityOptions), 1)
		}

		arg := c.Args().Get(2)
		allID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Entity ID not an integer"), 1)
		}

	} else if argnum == 5 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	output, err := AggregateLossCountAnalysis(aggregateType, entityType, allID, client)
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to perform ship loss analysis"), 1)
	}

	for _, loss := range output {
		client.Log.Printf("%v", loss)
	}

	return nil

}

func AggregateKills(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	allID := 0
	entityType := ""
	aggregateType := ""

	entityOptions := []string{"character", "corporation", "alliance"}
	aggregateOptions := []string{"corporation", "character", "ship", "system", "hour"}

	if argnum == 3 {

		aggregateType = c.Args().Get(0)
		validAggregate := false

		for _, v := range aggregateOptions {
			if v == aggregateType {
				validAggregate = true
			}
		}

		if validAggregate != true {
			return cli.NewExitError(fmt.Sprintf("Invalid aggregate type. Valid options are %v", aggregateOptions), 1)
		}

		entityType = c.Args().Get(1)
		validEntity := false

		for _, v := range entityOptions {
			if v == entityType {
				validEntity = true
			}
		}

		if validEntity != true {
			return cli.NewExitError(fmt.Sprintf("Invalid entity type. Valid options are %v", entityOptions), 1)
		}

		arg := c.Args().Get(2)
		allID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Entity ID not an integer"), 1)
		}

	} else if argnum == 5 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	output, err := AggregateKillCountAnalysis(aggregateType, entityType, allID, client)
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to perform ship loss analysis"), 1)
	}

	for _, loss := range output {
		client.Log.Printf("%v", loss)
	}

	return nil

}
