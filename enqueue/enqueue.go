package enqueue

import (
	"fmt"
	"time"

	"github.com/Crypta-Eve/truth/client"
	"github.com/Crypta-Eve/truth/store"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// ScrapeCharacterJob will queue a scrape of a particular character into the db
func ScrapeCharacterJob(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	if argnum == 0 {
		client.Log.Fatal("No char id given")
	} else if argnum == 1 {

		charid := c.Args().First()

		job := store.Queue{
			ID:          store.JobScrapeCharacter,
			Attempts:    0,
			CreatedAt:   time.Now(),
			AvailableAt: time.Now(),
			Args:        charid,
			Complete:    false,
		}

		err := client.Store.InsertQueueJob(job)
		if err != nil {
			err = errors.Wrap(err, "Failed to insert job")
			return cli.NewExitError(err, 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	return nil
}

// ScrapeCorporationJob will queue a scrape of a particular character into the db
func ScrapeCorporationJob(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	if argnum == 0 {
		client.Log.Fatal("No corp id given")
	} else if argnum == 1 {

		corpid := c.Args().First()

		job := store.Queue{
			ID:          store.JobScrapeCorporation,
			Attempts:    0,
			CreatedAt:   time.Now(),
			AvailableAt: time.Now(),
			Args:        corpid,
			Complete:    false,
		}

		err := client.Store.InsertQueueJob(job)
		if err != nil {
			err = errors.Wrap(err, "Failed to insert job")
			return cli.NewExitError(err, 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	return nil
}

// ScrapeAllianceJob will queue a scrape of a particular character into the db
func ScrapeAllianceJob(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	if argnum == 0 {
		client.Log.Fatal("No alliance id given")
	} else if argnum == 1 {

		allid := c.Args().First()

		job := store.Queue{
			ID:          store.JobScrapeAlliance,
			Attempts:    0,
			CreatedAt:   time.Now(),
			AvailableAt: time.Now(),
			Args:        allid,
			Complete:    false,
		}

		err := client.Store.InsertQueueJob(job)
		if err != nil {
			err = errors.Wrap(err, "Failed to insert job")
			return cli.NewExitError(err, 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	return nil
}

// ScrapeDailyJob will queue a scrape of a set of days
func ScrapeHistoryJob(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	if argnum == 2 {

		start := c.Args().Get(0)
		end := c.Args().Get(1)

		layout := "20060102"
		_, err := time.Parse(layout, start)
		if err != nil {
			return errors.New("Invalid start date, please use format YYYYMMDD")
		}

		_, err = time.Parse(layout, end)
		if err != nil {
			return errors.New("Invalid end date, please use format YYYYMMDD")
		}

		job := store.Queue{
			ID:          store.JobScrapeDaily,
			Attempts:    0,
			CreatedAt:   time.Now(),
			AvailableAt: time.Now(),
			Args:        fmt.Sprintf("%s|%s", start, end),
			Complete:    false,
		}

		err = client.Store.InsertQueueJob(job)
		if err != nil {
			err = errors.Wrap(err, "Failed to insert job")
			return cli.NewExitError(err, 1)
		}
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	return nil
}
