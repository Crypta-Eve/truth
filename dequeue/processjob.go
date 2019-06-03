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
		case store.JobScrapeCharacter:
			err := scrapeCharacter(client, job)
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("Error running zkill scrape for character job - %+v", job))
				client.Log.Println(err)
				continue
			}
		case store.JobScrapeCorporation:

		case store.JobScrapeAlliance:

		default:
			err := fmt.Errorf("Invalid job description from job - %d", job.ID)
			client.Log.Println(err)
			continue

		}
	}
}

func scrapeCharacter(c *client.Client, job store.Queue) error {

	//Grab the charID from args. Atm we do not support date ranges
	charID, err := strconv.Atoi(strings.Split(job.Args, "|")[0])

	if err != nil {
		return errors.Wrap(err, "CharID in job was not an integer")
	}

	apiURL := fmt.Sprintf("https://zkillboard.com/api/characterID/%v/page/", charID)

	pagenum := 0
	var pages []zkillmail

	for {
		pagenum += 1
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

		defer res.Body.Close()

		fmt.Println(len(body))

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
