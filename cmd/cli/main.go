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
				cli.Command{
					Name:        "missing",
					Category:    "Process",
					Action:      dequeue.ProcessMissingKillmails,
					Usage:       "Works on making sure we have all the killmails we have ids for",
					UsageText:   "jobs",
					Description: "",
				},
				cli.Command{
					Name:        "zkb",
					Category:    "Process",
					Action:      dequeue.ProcessMissingZKB,
					Usage:       "Make sure all kills have a zkb field",
					UsageText:   "zkb",
					Description: "",
				},
			},
		},
		cli.Command{
			Name:      "analyse",
			Category:  "Analysis",
			Usage:     "get info from the db",
			UsageText: "analyse [subcommand]",
			Subcommands: []cli.Command{
				cli.Command{
					Name:        "shiplosses",
					Category:    "Analysis",
					Usage:       "count the sum of each shiptype lost",
					UsageText:   "shiplosses [entityType] [entityID] (optional) [startdate] [enddate]",
					Description: "",
					Action:      analytics.ShipLosses,
				},
				//Need to work through and fix all the ones below here to follow the new model
				cli.Command{
					Name:        "alliance",
					Category:    "Analysis",
					Usage:       "",
					UsageText:   "alliance [subcommand]",
					Description: "",
					Subcommands: []cli.Command{

						cli.Command{
							Name:        "pilotlosses",
							Category:    "Analysis",
							Usage:       "pilotlosses",
							UsageText:   "pilotlosses [allianceid]",
							Description: "",
							Action:      analytics.PilotLosses,
						},
						cli.Command{
							Name:        "corplosses",
							Category:    "Analysis",
							Usage:       "corplosses",
							UsageText:   "corplosses [allianceid]",
							Description: "",
							Action:      analytics.CorpLosses,
						},
						cli.Command{
							Name:        "loclosses",
							Category:    "Analysis",
							Usage:       "loclosses",
							UsageText:   "loclosses [allianceid]",
							Description: "",
							Action:      analytics.LocationLosses,
						},
						cli.Command{
							Name:        "tzlosses",
							Category:    "Analysis",
							Usage:       "tzlosses",
							UsageText:   "tzlosses [allianceid]",
							Description: "",
							Action:      analytics.TZLosses,
						},
					},
				},
			},
		},
		cli.Command{
			Name:      "report",
			Category:  "Report",
			Usage:     "Serve a html report of the given entity",
			UsageText: "report [subcommand]",
			Subcommands: []cli.Command{
				cli.Command{
					Name:        "alliance",
					Category:    "Report",
					Usage:       "alliance [allianceid]",
					UsageText:   "",
					Description: "Serve a html report of an alliance",
					Action:      reports.AllianceReportServer,
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
