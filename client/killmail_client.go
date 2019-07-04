package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Crypta-Eve/truth/store"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (client *Client) FetchAndInsertKillmail(id int, hash string) error {

	url := fmt.Sprintf("https://esi.evetech.net/v1/killmails/%v/%v/", id, hash)

	body, err := client.MakeESIGet(url)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to query esi killmail id: %v | %v", id, hash))
	}
	km := store.ESIKillmail{}
	err = json.Unmarshal(body, &km)
	if err != nil {
		return errors.Wrap(err, "Error morphing esi killmail into struct")
	}

	ins := store.KillmailData{KillID: id, KillData: km}

	err = client.Store.InsertKillmail(ins)

	return err
}

func (client *Client) FetchAndInsertZKB(id int) error {

	url := fmt.Sprintf("https://zkillboard.com/api/killID/%d/", id)

	body, err := client.MakeGetRequest(url)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to query zkb for killid: %v", id))
	}

	km := []store.ZKBResponse{}
	err = json.Unmarshal(body, &km)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error morphing zkb killmail into struct data: \"%v\"", string(body)))
	}

	// Select the record to update with the zkb field
	filter := bson.M{
		"_id": id,
	}

	// Create the update
	update := bson.M{
		"$set": bson.M{
			"zkb": km[0].ZKB,
		},
	}

	err = client.Store.UpdateKillmail(filter, update)

	return err
}
