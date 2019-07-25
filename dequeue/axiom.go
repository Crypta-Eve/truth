package dequeue

import (
	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

//Axiom should be available on the localhost at :3005/killmail
// Posting an esi formatted killmail there will return an AxiomAttributes response

func ProcessMissingAttributes(c *cli.Context) error {
	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	mails, err := client.Store.GetKillsMissingAxiom()
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to get killmails missing axiom data"), 1)
	}

	for _, mail := range mails {
		err = client.FetchAndInsertAxiom(mail)
		if err != nil {
			client.Log.Println(err)
		}
	}

	return nil
}
