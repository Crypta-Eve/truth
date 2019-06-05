package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Crypta-Eve/truth/store"
	"github.com/pkg/errors"
)

func (client *Client) FetchAndInsertKillmail(id int, hash string) error {

	url := fmt.Sprintf("https://esi.evetech.net/v1/killmails/%v/%v/", id, hash)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return errors.Wrap(err, "Failed to buid esi killmail request")
	}

	req.Header.Set("User-Agent", client.UserAgent)

	res, err := client.HTTP.Do(req)

	if err != nil {
		return errors.Wrap(err, "Failed to make esi killmail request")
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return errors.Wrap(err, "Failed to read response from esi killmail request")
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
