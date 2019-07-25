package client

import (
	"bytes"
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

	body, err := client.MakeZKBGet(url)
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

func (client *Client) FetchAndInsertAxiom(mail store.KillmailData) error {

	data, err := json.Marshal(mail.KillData)
	if err != nil {
		client.Log.Printf("Error processing mail:\n %#v\nerr:%s\n", mail, err)
		return errors.Wrap(err, "Failed to marshal killmail...")
	}

	req, err := http.NewRequest("POST", "http://localhost:3005/killmail", bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Failed to build axiom request")
	}

	res, err := client.HTTP.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to contact axiom")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Error reading axiom body")
	}

	attr := store.FittingAttributes{}
	err = json.Unmarshal(body, &attr)
	if err != nil {
		return errors.Wrap(err, "Failed to decode axiom data")
	}

	filter := bson.M{
		"_id": mail.KillID,
	}

	update := bson.M{
		"$set": bson.M{
			"axiom": attr,
		},
	}

	err = client.Store.UpdateKillmail(filter, update)
	if err != nil {
		return errors.Wrap(err, "Failed to insert axiom data into mongo")
	}
	return nil
}
