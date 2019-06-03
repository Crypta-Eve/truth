package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Crypta-Eve/truth/dequeue"
	"github.com/Crypta-Eve/truth/enqueue"
	"github.com/Crypta-Eve/truth/testconnect"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var app *cli.App

var defSleep = cli.IntFlag{
	Name:  "sleep",
	Usage: "Custom Sleep Duration",
}

func init() {
	app = cli.NewApp()
	app.Name = "Eve Truth CLI"
	app.Version = "0.1.1"
	app.Authors = []cli.Author{
		{
			Name:  "Crypta Electrica",
			Email: "crypta@crypta.tech",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "queue",
			Category:  "Queue",
			Usage:     "Handles queueing scrapes",
			UsageText: "queue [subcommand]",
			Subcommands: []cli.Command{
				cli.Command{
					Name:        "character",
					Category:    "Queue",
					Action:      enqueue.ScrapeCharacterJob,
					Usage:       "Handles queuing of character scrapes",
					UsageText:   "character [charid] (optional [startdate] [enddate])",
					Description: "Queue a character scrape",
				},

				cli.Command{
					Name:        "corporation",
					Category:    "Queue",
					Action:      enqueue.ScrapeCorporationJob,
					Usage:       "Handles queuing of corp scrapes",
					UsageText:   "character [corpid] (optional [startdate] [enddate])",
					Description: "Queue a corp scrape",
				},

				cli.Command{
					Name:        "alliance",
					Category:    "Queue",
					Action:      enqueue.ScrapeAllianceJob,
					Usage:       "Handles queuing of alliance scrapes",
					UsageText:   "character [allianceid] (optional [startdate] [enddate])",
					Description: "Queue an alliance scrape",
				},
			},
		},
		cli.Command{
			Name:        "test",
			Category:    "Test",
			Usage:       "Test connections",
			UsageText:   "test",
			Action:      testconnect.TestExternalAppConnections,
			Description: "Test connectivity to external appliances",
		},
		cli.Command{
			Name:      "process",
			Category:  "Process",
			Usage:     "Used to process pending jobs",
			UsageText: "process [subcommand]",
			Subcommands: []cli.Command{
				cli.Command{
					Name:        "jobs",
					Category:    "Process",
					Action:      dequeue.ProcessJobQueue,
					Usage:       "Works on jobs that are in the job queue",
					UsageText:   "jobs",
					Description: "Work on items in the job queue. Will check for new jobs every 5s and loop until interrupted",
				},
			},
		},
	}

	viper.SetConfigType("json")
	viper.SetConfigFile("env.json")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "unable to read in configuration file: ", err)
		os.Exit(1)
	}

}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
