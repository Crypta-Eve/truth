package client

import (
	"log"
	"time"

	"net/http"
	"os"
	"sync"

	"github.com/Crypta-Eve/truth/store"
	"github.com/spf13/viper"
)

type (
	Client struct {
		HTTP      *http.Client
		Store     store.Store
		Log       *log.Logger
		UserAgent string
		ESIRateLimit *safeCounter
	}

	safeCounter struct {
		cnt int
		mux sync.Mutex
	}
)

func New() (*Client, error) {
	logger := log.New(os.Stdout, "CLIENT:", log.Lshortfile|log.Ldate|log.Ltime)

	// now check we have access to mongo

	envDB := viper.GetStringMapString("db")
	store, err := store.SetupStore(envDB)

	if err != nil {
		return nil, err
	}

	rateLim := &make(safeCounter)

	go func() {
		for {
			time.Sleep(time.Second)
			if rateLim.Value() > 0 {
				rateLim.Dec()
			}
		}
	}()

	return &Client{
		HTTP: &http.Client{
			Timeout: time.Second * 40,
		},
		Store:        store,
		Log:          logger,
		UserAgent:    viper.GetString("user_agent"),
		ESIRateLimit: rateLim,
		RetryLimit:   10,
	}, nil

}

func (c *Client) makeRawHTTPGet(url string) ([]byte, int, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to buid http request")
	}

	req.Header.Set("User-Agent", client.UserAgent)

	res, err := client.HTTP.Do(req)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to read response from request")
	}

	return body, res.StatusCode, nil
}

func (c *Client) MakeESIGet(url string) ([]byte, error) {

	retriesRemain := c.RetryLimit
	for retriesRemain > 1 {
		retriesRemain--

		for c.ESIRateLimit.Value() > 10 {
			time.sleep(500 * time.Millisecond)
		}

		body, status, err := makeRawHTTPGet(url)
		if err != nil {
			continue
		}
		if !(status >= 200 && status < 300) {
			// Increment the counter :(
			c.ESIRateLimit.Inc()
			continue
		}

		return body, err
	}

	return nil, fmt.Errorf("Max retries exceeded for url: ", url)
}

func (c *Client) MakeGetRequest(url string) ([]byte, error) {
	retriesRemain := c.RetryLimit
	for retriesRemain > 1 {
		retriesRemain--

		for c.ESIRateLimit.Value() > 10 {
			time.sleep(500 * time.Millisecond)
		}

		body, status, err := makeRawHTTPGet(url)
		if err != nil {
			continue
		}
		if !(status >= 200 && status < 300) {
			// Increment the counter :(
			c.ESIRateLimit.Inc()
			continue
		}

		return body, err
	}

	return nil, fmt.Errorf("Max retries exceeded for url: ", url)
}

// Inc increments the counter.
func (c *safeCounter) Inc() {
	c.mux.Lock()
	c.cnt++
	c.mux.Unlock()
}

// Dec decrements the counter.
func (c *safeCounter) Dec() {
	c.mux.Lock()
	c.cnt++
	c.mux.Unlock()
}

// Value returns the current value of the counter.
func (c *safeCounter) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.cnt
}
