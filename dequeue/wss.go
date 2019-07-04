package dequeue

import (
	"strings"

	"github.com/Crypta-Eve/truth/client"
	"github.com/Crypta-Eve/truth/store"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	// The zkillboard wss url
	address          string = "wss://zkillboard.com:2096"
	subscribeMessage string = "{\"action\":\"sub\", \"channel\":\"killstream\"}"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func WholeWSS(c *cli.Context) error {

	client, err := client.New()

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return cli.NewExitError(err, 1)
	}

	client.Log.Println("Starting WSS Scraper")

	var dialer *websocket.Dialer

	ws, resp, err := dialer.Dial(address, nil)
	if err != nil {
		if resp != nil {
			return cli.NewExitError(errors.Wrap(err, "Failed to dial the wss"), 1)
		} else {
			return cli.NewExitError(errors.Wrap(err, "Failed to dial the wss"), 1)
		}
	}

	ws.WriteMessage(websocket.TextMessage, []byte(subscribeMessage))
	for {

		var message store.WSSKillmail
		err = ws.ReadJSON(&message)
		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Error retreiving wss message"), 1)
		}
		idhash, km := convertKillmail(message)
		client.Log.Printf("Received wss killmail - %#v", km.KillID)

		err = client.Store.InsertKillIDHash(idhash)
		if err != nil {
			if strings.Contains(err.Error(), "dup key") {
				client.Log.Printf("Duplicate idhash ignored - %v", idhash.ID)
				continue
			}
			return cli.NewExitError(errors.Wrap(err, "Error attempting to insert idhash"), 1)
		}

		err = client.Store.InsertKillmail(km)
		if err != nil {
			if strings.Contains(err.Error(), "dup key") {
				client.Log.Printf("Duplicate killmail ignored - %v", idhash.ID)
				continue
			}
			return cli.NewExitError(errors.Wrap(err, "Error attempting to insert killmail"), 1)
		}
	}

}

func convertKillmail(w store.WSSKillmail) (idhash store.ScrapeQueue, km store.KillmailData) {

	out := store.KillmailData{}
	esik := store.ESIKillmail{}

	out.KillID = w.KillmailID
	esik.Attackers = w.Attackers
	esik.KillmailID = w.KillmailID
	esik.KillmailTime = w.KillmailTime
	esik.SolarSystemID = w.SolarSystemID
	esik.Victim = w.Victim
	out.KillData = esik
	out.ZKBData = w.ZKB

	idhash = store.ScrapeQueue{ID: out.KillID, Hash: out.ZKBData.Hash}

	return idhash, out
}
