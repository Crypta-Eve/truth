package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

func (client *Client) ResolveIDsToNames(ids []int) (names map[int]string, err error) {

	names = make(map[int]string)

	numIDs := len(ids)

	switch {
	case numIDs == 0:
		return names, nil
	case numIDs < 1000:
		names, err = client.resolveIDsToNamesESI(ids)
	case numIDs > 1000:
		chunkSize := 900
		chunks := split(ids, chunkSize)
		for _, chunk := range chunks {
			n, err := client.resolveIDsToNamesESI(chunk)
			if err != nil {
				return names, errors.Wrap(err, "Error performing esi name lookup on split data")
			}
			for k, v := range n {
				names[k] = v
			}
		}
	}

	return names, err
}

//split slice into smaller shucks acceptable for esi
//sourced from https://gist.github.com/xlab/6e204ef96b4433a697b3
func split(ids []int, lim int) [][]int {
	var chunk []int
	chunks := make([][]int, 0, len(ids)/lim+1)
	for len(ids) >= lim {
		chunk, ids = ids[:lim], ids[lim:]
		chunks = append(chunks, chunk)
	}

	if len(ids) > 0 {
		chunks = append(chunks, ids)
	}

	return chunks
}

func (client *Client) resolveIDsToNamesESI(ids []int) (names map[int]string, err error) {
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

	//Could not resolve, so just return the ids
	if res.StatusCode != 200 {
		for id := range ids {
			names[id] = strconv.Itoa(id)
		}
		return names, nil
	}

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
