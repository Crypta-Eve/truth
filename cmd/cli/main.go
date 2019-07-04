package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Crypta-Eve/truth/analytics"
	"github.com/Crypta-Eve/truth/dequeue"
	"github.com/Crypta-Eve/truth/enqueue"
	"github.com/Crypta-Eve/truth/reports"
	"github.com/Crypta-Eve/truth/testconnect"
	"github.com/Crypta-Eve/truth/wss"

	// "github.com/pkg/profile"
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
	app.Description = "CLI interface for Eve Truth Engine"
	app.Authors = []cli.Author{
		{
			Name:  "Crypta Electrica",
			Email: "crypta@crypta.tech",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "queue",
			Category:  "Queue",
			Usage:     "Handles queueing scrapes",
			UsageText: "queue [subcommand]",
			Subcommands: []cli.Command{
				{
					Name:        "character",
					Category:    "Queue",
					Action:      enqueue.ScrapeCharacterJob,
					Usage:       "Handles queuing of character scrapes",
					UsageText:   "character [charid] (optional [startdate] [enddate])",
					Description: "Queue a character scrape",
				},
				{
					Name:        "corporation",
					Category:    "Queue",
					Action:      enqueue.ScrapeCorporationJob,
					Usage:       "Handles queuing of corp scrapes",
					UsageText:   "character [corpid] (optional [startdate] [enddate])",
					Description: "Queue a corp scrape",
				},
				{
					Name:        "alliance",
					Category:    "Queue",
					Action:      enqueue.ScrapeAllianceJob,
					Usage:       "Handles queuing of alliance scrapes",
					UsageText:   "character [allianceid] (optional [startdate] [enddate])",
					Description: "Queue an alliance scrape",
				},
				{
					Name:        "history",
					Category:    "Queue",
					Action:      enqueue.ScrapeHistoryJob,
					Usage:       "Handles queuing of daily scrapes",
					UsageText:   "histroy startdate enddate",
					Description: "Queue a historical scrape. date format it YYYYMMDD",
				},
			},
		},
		{
			Name:        "test",
			Category:    "Test",
			Usage:       "Test connections",
			UsageText:   "test",
			Action:      testconnect.TestExternalAppConnections,
			Description: "Test connectivity to external appliances",
		},
		{
			Name:      "process",
			Category:  "Process",
			Usage:     "Used to process pending jobs",
			UsageText: "process [subcommand]",
			Subcommands: []cli.Command{
				{
					Name:        "jobs",
					Category:    "Process",
					Action:      dequeue.ProcessJobQueue,
					Usage:       "Works on jobs that are in the job queue",
					UsageText:   "jobs",
					Description: "Work on items in the job queue. Will check for new jobs every 5s and loop until interrupted",
				},
				{
					Name:        "missing",
					Category:    "Process",
					Action:      dequeue.ProcessMissingKillmails,
					Usage:       "Works on making sure we have all the killmails we have ids for",
					UsageText:   "jobs",
					Description: "",
				},
				{
					Name:        "zkb",
					Category:    "Process",
					Action:      dequeue.ProcessMissingZKB,
					Usage:       "Make sure all kills have a zkb field",
					UsageText:   "zkb",
					Description: "",
				},
				{
					Name:        "wss",
					Category:    "Process",
					Action:      wss.WholeWSS,
					Usage:       "Monitor WSS for real time killmails",
					UsageText:   "zkb",
					Description: "",
				},
				{
					Name:        "hashmaint",
					Category:    "Process",
					Action:      dequeue.HashTableMaint,
					Usage:       "Maintain Hashes",
					UsageText:   "hashmaint",
					Description: "",
				},
			},
		},
		{
			Name:      "analyse",
			Category:  "Analysis",
			Usage:     "get info from the db",
			UsageText: "analyse [subcommand]",
			Subcommands: []cli.Command{
				{
					Name:        "losses",
					Category:    "Analysis",
					Usage:       "count the sum of each shiptype lost",
					UsageText:   "losses [aggregate] [entityType] [entityID] (optional [startdate] [enddate])",
					Description: "",
					Action:      analytics.AggregateLosses,
				},
				{
					Name:        "killed",
					Category:    "Analysis",
					Usage:       "group by aggregate those killed by the entity",
					UsageText:   "killed [aggregate] [entityType] [entityID] (optional [startdate] [enddate])",
					Description: "",
					Action:      analytics.AggregateKilled,
				},
				{
					Name:        "kills",
					Category:    "Analysis",
					Usage:       "group by aggregate those killing from within the entity",
					UsageText:   "kills [aggregate] [entityType] [entityID] (optional [startdate] [enddate])",
					Description: "",
					Action:      analytics.AggregateKills,
				},
			},
		},
		{
			Name:      "report",
			Category:  "Report",
			Usage:     "Serve a html report of the given entity",
			UsageText: "report [subcommand]",
			Subcommands: []cli.Command{
				{
					Name:        "alliance",
					Category:    "Report",
					Usage:       "alliance [allianceid]",
					UsageText:   "",
					Description: "Serve a html report of an alliance",
					Action:      reports.AllianceReportServer,
				},
			},
		},
		{
			Name:        "setup",
			Category:    "Setup",
			Usage:       "Setup",
			UsageText:   "Setup the systems",
			Action:      enqueue.PerformSetup,
			Description: "This will ensure everything is set up correctly",
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

	// defer profile.Start().Stop()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
