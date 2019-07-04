package setup

import (
	"github.com/Crypta-Eve/truth/client"
	"github.com/Crypta-Eve/truth/store"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func HashTableMaint(c *cli.Context) error {
	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	mails, err := client.Store.ListAllExistingKillmails()
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to get existing killmails"), 1)
	}

	// Morph to a scrape queue list
	hashList := make([]store.ScrapeQueue, len(mails))
	for i, mail := range mails {
		hashList[i] = store.ScrapeQueue{ID: mail.KillID, Hash: mail.ZKBData.Hash}
	}

	err = client.Store.UpdateManyHashes(hashList)

	return err

}
