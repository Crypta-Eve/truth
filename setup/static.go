package setup

import (
	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Eve ID Ranges - https://gist.github.com/a-tal/5ff5199fdbeb745b77cb633b7f4400bb

func PopulateStaticData(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	populateUniverse(client)

	return nil

}

func populateUniverse(client *client.Client) error {
	// First Step is to populate the region list

	// const urlRegion = "https://esi.evetech.net/latest/universe/regions/?datasource=tranquility"
	// regionBody, err := client.MakeESIGet(url)
	// var regions []int
	// err = json.Unmarshal(body, &regions)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func populateTypes(client *client.Client) error {
	return nil
}

func populateGroups(client *client.Client) error {
	return nil
}

func populateCategories(client *client.Client) error {
	return nil
}
