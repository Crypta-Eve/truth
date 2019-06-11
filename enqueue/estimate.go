package enqueue

// import (
// 	"strconv"
// 	"time"

// 	"github.com/Crypta-Eve/truth/client"
// 	"github.com/pkg/errors"
// 	"github.com/urfave/cli"
// )

//Come back to this, I want to be able to estimate the number of killmails

// func EstimateAllianceScrapeCount(c *cli.Context) error {
// 	client, err := client.New()

// 	if err != nil {
// 		err = errors.Wrap(err, "failed to create client")
// 		return cli.NewExitError(err, 1)
// 	}

// 	argnum := len(c.Args())

// 	allianceID := 0
// 	startDate := time.Now()
// 	endDate := time.Now()
// 	dates := false

// 	if argnum == 0 {
// 		client.Log.Fatal("No alliance id given")
// 	} else if argnum == 1 {

// 		id := c.Args().First()
// 		allianceID, err = strconv.Atoi(id)

// 		if err != nil {
// 			return cli.NewExitError(errors.Wrap(err, "Alliance ID must be an int"), 1)
// 		}

// 	} else if argnum == 3 {
// 		id := c.Args().First()
// 		allianceID, err = strconv.Atoi(id)

// 		if err != nil {
// 			return cli.NewExitError(errors.Wrap(err, "Alliance ID must be an int"), 1)
// 		}

// 		// c.Args()[1] //start
// 		// c.Args()[2] //End
// 	} else {
// 		client.Log.Fatal("invalid number of arguments")
// 	}

// 	return nil
// }
