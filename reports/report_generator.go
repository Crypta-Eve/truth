package reports

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/Crypta-Eve/truth/analytics"
	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type ChartData struct {
	Labels []string
	Values []int
}

type AllianceReport struct {
	CorpKills  ChartData
	CorpLosses ChartData

	PilotKills  ChartData
	PilotLosses ChartData

	ShipKillsWith ChartData
	ShipLosses    ChartData

	LocationKills  ChartData
	LocationLosses ChartData

	TZKills  ChartData
	TZLosses ChartData
}

func AllianceReportServer(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	argnum := len(c.Args())

	allID := 0

	if argnum == 0 {
		client.Log.Fatal("No alliance id given")
	} else if argnum == 1 {
		arg := c.Args().First()
		allID, err = strconv.Atoi(arg)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Alliance ID not an integer"), 1)
		}

	} else if argnum == 3 {
		client.Log.Print("Dates NYI")
	} else {
		client.Log.Fatal("invalid number of arguments")
	}

	client.Log.Printf("Starting ReportGen for Alliance %v", allID)

	startTime := time.Now()

	data := AllianceReport{}

	// Corp Kills
	corpKills, err := analytics.AggregateKillsCountAnalysis("corporation", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.CorpKills = trimToTop30(corpKills)

	// Corp Losses
	corpLosses, err := analytics.AggregateLossCountAnalysis("corporation", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.CorpLosses = trimToTop30(corpLosses)

	// Pilot Kills
	pilotKills, err := analytics.AggregateKillsCountAnalysis("character", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.PilotKills = trimToTop30(pilotKills)

	// Pilot Losses
	pilotLosses, err := analytics.AggregateLossCountAnalysis("character", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.PilotLosses = trimToTop30(pilotLosses)

	//TODO: UP TO HERE

	// Ship Kills
	shipKills, err := analytics.AggregateKillsCountAnalysis("ship", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.ShipKillsWith = trimToTop30(shipKills)

	// Ship Losses
	shipLosses, err := analytics.AggregateLossCountAnalysis("ship", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.ShipLosses = trimToTop30(shipLosses)

	// Location Kills
	locKills, err := analytics.AggregateKillsCountAnalysis("system", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.LocationKills = trimToTop30(locKills)

	// Location Losses
	locLosses, err := analytics.AggregateLossCountAnalysis("system", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.LocationLosses = trimToTop30(locLosses)

	// TZ Kills
	tzKills, err := analytics.AggregateKillsCountAnalysis("hour", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.TZKills = trimToTop30(tzKills)

	// TZ Losses
	tzLosses, err := analytics.AggregateLossCountAnalysis("hour", "alliance", allID, client)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	data.TZLosses = trimToTop30(tzLosses)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("reports/main.html", "reports/report.html", "reports/plotly.html")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		t.ExecuteTemplate(w, "report", data)
	})

	endTime := time.Now()

	client.Log.Printf("Time taken to process - %v", endTime.Sub(startTime))

	client.Log.Print("Report ready on port 1080")
	http.ListenAndServe(":1080", nil)

	return nil
}

func trimToTop30(data analytics.PairList) (out ChartData) {

	length := len(data)

	if length > 30 {

		labels := make([]string, 30)
		values := make([]int, 30)

		for i := 0; i < 30; i++ {
			pair := data[i]
			labels[i] = pair.Key
			values[i] = pair.Value
		}

		return ChartData{Labels: labels, Values: values}
	} else {

		labels := make([]string, length)
		values := make([]int, length)

		for i := 0; i < length; i++ {
			pair := data[i]
			labels[i] = pair.Key
			values[i] = pair.Value
		}

		return ChartData{Labels: labels, Values: values}
	}
}
