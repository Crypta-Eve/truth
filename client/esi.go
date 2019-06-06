package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func (client *Client) ResolveIDsToNames(ids []int) (names map[int]string, err error) {

	names = make(map[int]string)

	const url = "https://esi.evetech.net/latest/universe/names/?datasource=tranquility"

	bd, err := json.Marshal(ids)

	if err != nil {
		return names, errors.Wrap(err, "Failed to build body for esi name request")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bd))
	if err != nil {
		return names, errors.Wrap(err, "Failed to build ID resolving request")
	}

	req.Header.Set("User-Agent", client.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.HTTP.Do(req)
	if err != nil {
		return names, errors.Wrap(err, "Failed to request names from ESI")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return names, errors.Wrap(err, "Error in decoding names response from esi")
	}

	type ESINameResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var response []ESINameResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return names, errors.Wrap(err, "Failed to unwrap data from esi names")
	}

	for _, item := range response {
		names[item.ID] = item.Name
	}

	return names, nil

}
