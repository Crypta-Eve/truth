package dequeue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
		case store.JobSeedDatabase:
			//Really this never should be hit... Mark it complete so it shouldnt get hit again.
			client.Store.MarkJobComplete(job)
			continue

		case store.JobScrapeCharacter, store.JobScrapeCorporation, store.JobScrapeAlliance:
			err := scrapePlayer(client, job)
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("Error running zkill scrape for job - %+v", job))
				client.Log.Println(err)
				continue
			}

		case store.JobScrapeDaily:
			err := scrapeHistory(client, job)
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("Error running zkill scrape for daily job - %+v", job))
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

	for {

		missingMails, err := client.Store.ListMissingKillmails()
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Failed to get kills that are not in list"), 1)
		}

		numMissing := len(missingMails)
		if numMissing == 0 {
			client.Log.Println("No missing killmails, sleeping...")
			time.Sleep(5 * time.Minute)
			continue
		}

		client.Log.Printf("Have %v killmails to fetch", numMissing)

		time.Sleep(10 * time.Second)

		const numThreads = 1000

		if numMissing < numThreads {

			for i, mail := range missingMails {
				client.Log.Printf("Processing mail %d/%d - %d", i+1, numMissing, mail.ID)
				err := client.FetchAndInsertKillmail(mail.ID, mail.Hash)

				if err != nil {
					if strings.Contains(err.Error(), "dup key") {
						client.Log.Printf("Duplicate killmail ignored - %v", mail.ID)
						continue
					}
					return cli.NewExitError(errors.Wrap(err, "Error attempting to insert killmail"), 1)
				}
			}
			return nil
		}

		batchSize := (numMissing / numThreads) + 1

		client.Log.Printf("Batch Size is %v", batchSize)

		var batches [][]store.ScrapeQueue

		for batchSize < len(missingMails) {
			missingMails, batches = missingMails[batchSize:], append(batches, missingMails[0:batchSize:batchSize])
		}

		batches = append(batches, missingMails)

		// for _, batch := range batches {
		// 	fmt.Println(batch)
		// }

		var waitgroup sync.WaitGroup

		for i, batch := range batches {
			waitgroup.Add(1)

			btch := batch
			i1 := i
			go func() {
				client.Log.Printf("Starting sub process %v", i1)
				// client.Log.Printf("Sub Process Batch -  %v", btch)
				amount := len(btch)
				for j, mail := range btch {
					client.Log.Printf("%d - Processing mail %d/%d - %d", i1, j+1, amount, mail.ID)
					err := client.FetchAndInsertKillmail(mail.ID, mail.Hash)
					if err != nil {
						if strings.Contains(err.Error(), "dup key") {
							client.Log.Printf("%d - Duplicate killmail ignored - %v", i1, mail.ID)
							continue
						}
						client.Log.Printf("err: %v", err)
					}
				}

				waitgroup.Done()
				client.Log.Printf("Subprocess %v done!", i1)
			}()
		}

		waitgroup.Wait()

		client.Log.Println("Job Complete")

	}

}

func ProcessMissingZKB(c *cli.Context) error {
	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	for {

		ids, err := client.Store.GetKillsMissingZKB()

		client.Log.Printf("Found %v killmails that need updating", len(ids))

		if err != nil {
			err = errors.Wrap(err, "Failed to get list of killmails that are short zkb")
			return cli.NewExitError(err, 1)

		}

		numMissing := len(ids)
		if numMissing == 0 {
			client.Log.Println("No missing zkb, sleeping for 5s")
			time.Sleep(1 * time.Minute)
			continue
		}

		// Lets batch these out, for the case where we have bulk missing zkbs....... Sorry Squizz for hammer.....

		numThreads := 50

		if numMissing < numThreads {

			for i, mail := range ids {
				client.Log.Printf("Processing mail %d/%d - %d", i, numMissing, mail.ID)
				err := client.FetchAndInsertZKB(mail.ID)
				if err != nil {
					client.Log.Println(errors.Wrap(err, "Error trying to create new killmail"))
				}
			}

			return nil
		}

		var batches [][]store.ScrapeQueue

		batchSize := (numMissing / numThreads) + 1

		for batchSize < len(ids) {
			ids, batches = ids[batchSize:], append(batches, ids[0:batchSize:batchSize])
		}

		batches = append(batches, ids)

		var waitgroup sync.WaitGroup

		for i, batch := range batches {
			waitgroup.Add(1)

			btch := batch
			i1 := i
			go func() {
				client.Log.Printf("Starting sub process %v", i1)
				// client.Log.Printf("Sub Process Batch -  %v", btch)
				amount := len(btch)
				for j, mail := range btch {
					client.Log.Printf("Processing zkb %d/%d - %d", j, amount, mail.ID)
					err := client.FetchAndInsertZKB(mail.ID)
					if err != nil {
						client.Log.Println(errors.Wrap(err, "Error trying to create new zkb"))
					}
				}

				waitgroup.Done()
				client.Log.Printf("Subprocess %v done!", i1)
			}()
		}

		waitgroup.Wait()

		client.Log.Println("Job Complete")
	}
}

func scrapePlayer(c *client.Client, job store.Queue) error {

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

		err := c.Store.InsertKillIDHash(mailIDHash)
		if err != nil {
			if strings.Contains(err.Error(), "dup key") {
				c.Log.Printf("Duplicate killmail hash ignored - %v", mailIDHash.ID)
				continue
			} else {
				return cli.NewExitError(errors.Wrap(err, "Error attempting to insert killmail"), 1)
			}
		}

	}

	c.Store.MarkJobComplete(job)

	return nil

}

func scrapeHistory(c *client.Client, job store.Queue) error {

	layout := "20060102"

	start, _ := time.Parse(layout, (strings.Split(job.Args, "|")[0]))
	end, _ := time.Parse(layout, (strings.Split(job.Args, "|")[1]))

	days := int(end.Sub(start).Hours()/24) + 1
	c.Log.Printf("History Scrape across %v days", days)

	const maxprocs = 10

	switch {
	case days >= 0 && days < 10:
		for i := 0; i < days; i++ {
			dt := start.Add(time.Duration(i*24) * time.Hour)
			err := scrapeDaily(c, dt.Format("20060102"))
			if err != nil {
				c.Log.Println(err)
				continue
			}
		}
		c.Store.MarkJobComplete(job)

	case days > 10:

		dates := make([]string, days)

		for i := 0; i < days; i++ {
			dt := start.Add(time.Duration(i*24) * time.Hour)
			dates[i] = dt.Format("20060102")
		}

		var batches [][]string

		for maxprocs < len(dates) {
			dates, batches = dates[maxprocs:], append(batches, dates[0:maxprocs:maxprocs])
		}

		batches = append(batches, dates)

		// Now we can spawn multiple routines to scrape it all!

		var waitgroup sync.WaitGroup

		for i, batch := range batches {
			waitgroup.Add(1)

			btch := batch
			i1 := i
			go func() {
				c.Log.Printf("Starting sub process %v", i1)
				c.Log.Printf("Sub Process Batch -  %v", btch)
				for _, day := range btch {
					err := scrapeDaily(c, day)
					if err != nil {
						c.Log.Println(errors.Wrap(err, fmt.Sprintf("Failed to scrape a historical day %v. In multi process.", day)))
					}
				}

				waitgroup.Done()
				c.Log.Printf("Subprocess %v done!", i1)
			}()
		}

		waitgroup.Wait()

		c.Log.Println("Job Complete")

		c.Store.MarkJobComplete(job)
	}

	return nil
}

func scrapeDaily(c *client.Client, date string) error {
	apiURL := "https://zkillboard.com/api/history/%v/"

	urlToHit := fmt.Sprintf(apiURL, date)

	req, err := http.NewRequest(http.MethodGet, urlToHit, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to build daily zkill request")
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTP.Do(req)

	if err != nil {
		return errors.Wrap(err, "Failed to complete daily zkill request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read daily response")
	}

	defer res.Body.Close()

	pageResponse := make(map[string]string, 0)

	err = json.Unmarshal(body, &pageResponse)

	if err != nil {
		return err
	}

	for k, v := range pageResponse {
		if k == " day" {
			// there is a a day key in there.....
			continue
		}

		// Now to sanity check, the killID should be an int, the hash should be 40 characters long...
		// This will fail semi silently....
		killID, err := strconv.Atoi(k)
		if err != nil || len(v) != 40 {
			c.Log.Printf("Invalid response from zkill history api. %v, %v.\n", k, v)
			continue
		}

		ins := store.ScrapeQueue{ID: killID, Hash: v}
		err = c.Store.InsertKillIDHash(ins)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue
			}
			c.Log.Println(err)
			continue
		}
	}

	return nil

}
