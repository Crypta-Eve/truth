package dequeue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Crypta-Eve/truth/client"
	"github.com/Crypta-Eve/truth/store"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type (
	zkillmail struct {
		Killmail_id int `json:"killmail_id"`
		Zkb         zkb `json:"zkb"`
	}

	zkb struct {
		Hash string `json:"hash"`
	}
)

func ProcessJobQueue(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	for {
		//Now we need to grab a queue from the job queue

		job, err := client.Store.PopQueueJob()

		if err != nil {
			// No jobs available
			client.Log.Println("No jobs available, sleeping for 5s")
			time.Sleep(5 * time.Second)
			continue
		}

		// Job available, lets go to work

		switch job.ID {
		case store.JobScrapeCharacter, store.JobScrapeCorporation, store.JobScrapeAlliance:
			err := scrapeCharacter(client, job)
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("Error running zkill scrape for character job - %+v", job))
				client.Log.Println(err)
				continue
			}

		default:
			err := fmt.Errorf("Invalid job description from job - %d", job.ID)
			client.Log.Println(err)
			continue

		}
	}
}

func ProcessMissingKillmails(c *cli.Context) error {
	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	ids, err := client.Store.ListAllExistingIDs()

	client.Log.Printf("Found %v existing killmail ids", len(ids))

	if err != nil {
		err = errors.Wrap(err, "Failed to get list of all existing killmails")
		return cli.NewExitError(err, 1)

	}

	// Now to build the massive request that will tell us what we are missing
	var IDList []int
	for _, v := range ids {
		IDList = append(IDList, v)
	}

	missingMails, err := client.Store.GetKillsNotInList(IDList)
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Failed to get kills that are not in list"), 1)
	}

	numMissing := len(missingMails)
	client.Log.Printf("Have %v killmails to fetch", numMissing)

	for i, mail := range missingMails {
		client.Log.Printf("Processing mail %d/%d - %d", i, numMissing, mail.ID)
		err := client.FetchAndInsertKillmail(mail.ID, mail.Hash)
		if err != nil {
			client.Log.Println(errors.Wrap(err, "Error trying to create new killmail"))
		}
	}

	return nil
}

func ProcessMissingZKB(c *cli.Context) error {
	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	ids, err := client.Store.GetKillsMissingZKB()

	client.Log.Printf("Found %v killmails that need updating", len(ids))

	if err != nil {
		err = errors.Wrap(err, "Failed to get list of killmails that are short zkb")
		return cli.NewExitError(err, 1)

	}

	numMissing := len(ids)

	for i, mail := range ids {
		client.Log.Printf("Processing mail %d/%d - %d", i, numMissing, mail.ID)
		err := client.FetchAndInsertZKB(mail.ID)
		if err != nil {
			client.Log.Println(errors.Wrap(err, "Error trying to create new killmail"))
		}
	}

	return nil
}

func scrapeCharacter(c *client.Client, job store.Queue) error {

	//Grab the charID from args. Atm we do not support date ranges
	charID, err := strconv.Atoi(strings.Split(job.Args, "|")[0])

	if err != nil {
		return errors.Wrap(err, "CharID in job was not an integer")
	}

	var entityType string

	switch job.ID {
	case store.JobScrapeCharacter:
		entityType = "character"
	case store.JobScrapeCorporation:
		entityType = "corporation"
	case store.JobScrapeAlliance:
		entityType = "alliance"
	}

	apiURL := fmt.Sprintf("https://zkillboard.com/api/%vID/%v/page/", entityType, charID)

	pagenum := 0
	var pages []zkillmail

	for {
		pagenum++
		urlToHit := fmt.Sprintf(apiURL+"%d/", pagenum)

		req, err := http.NewRequest(http.MethodGet, urlToHit, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to build character zkill request")
		}

		req.Header.Set("User-Agent", c.UserAgent)

		res, err := c.HTTP.Do(req)

		if err != nil {
			return errors.Wrap(err, "Failed to complete character zkill request")
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "Failed to read zkill response")
		}

		defer res.Body.Close()

		fmt.Print(">")

		if len(body) < 4 {
			break
		}

		pg := []zkillmail{}
		err = json.Unmarshal(body, &pg)

		if err != nil {
			return errors.Wrap(err, "Trying to decode zkill character response")
		}

		c.Store.MaintainJobReservation(job, time.Now().Add(60*time.Second))

		pages = append(pages, pg...)
	}

	// fmt.Printf("%+v\n", pages)

	for _, mail := range pages {

		mailIDHash := store.ScrapeQueue{ID: mail.Killmail_id, Hash: mail.Zkb.Hash}

		c.Store.InsertKillIDHash(mailIDHash)

	}

	c.Store.MarkJobComplete(job)

	return nil

}
