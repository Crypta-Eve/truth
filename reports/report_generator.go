package reports

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Crypta-Eve/truth/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("reports/main.html", "reports/report.html", "reports/plotly.html")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		t.ExecuteTemplate(w, "plotly", fmt.Sprintf("Hello World! - %v", allID))
	})

	http.ListenAndServe(":1080", nil)

	return nil
}
