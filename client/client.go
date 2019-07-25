package client

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Crypta-Eve/truth/store"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	Client struct {
		HTTP         *http.Client
		Store        store.Store
		Log          *log.Logger
		UserAgent    string
		ESIRateLimit *safeCounter
		ZKBRateLimit *safeCounter
		RetryLimit   int
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

	rateLimESI := &safeCounter{}
	rateLimZkill := &safeCounter{}

	go func() {
		for {
			time.Sleep(time.Second)
			if rateLimESI.Value() > 0 {
				rateLimESI.Dec()
			}
			if rateLimZkill.Value() > 0 {
				rateLimZkill.Dec()
			}

		}
	}()

	return &Client{
		HTTP: &http.Client{
			Timeout: time.Second * 40,
			Transport: &http.Transport{
				MaxConnsPerHost:     10,
				MaxIdleConnsPerHost: 2,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
		},
		Store:        store,
		Log:          logger,
		UserAgent:    viper.GetString("user_agent"),
		ESIRateLimit: rateLimESI,
		ZKBRateLimit: rateLimZkill,
		RetryLimit:   25,
	}, nil

}

func (c *Client) makeRawHTTPGet(url string) ([]byte, int, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to buid http request")
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTP.Do(req)

	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to make request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to read response from request")
	}

	return body, res.StatusCode, nil
}

func (c *Client) MakeESIGet(url string) (out []byte, err error) {

	retriesRemain := c.RetryLimit
	for retriesRemain > 1 {
		retriesRemain--

		for c.ESIRateLimit.Value() > 10 {
			time.Sleep(500 * time.Millisecond)
		}

		body, status, err := c.makeRawHTTPGet(url)
		if err != nil {
			if strings.Contains(err.Error(), "too many open files") {
				// This is not going to hurt to keep retrying
				retriesRemain++
			} else {
				// fmt.Printf("ESI GET ERROR - %v\n", err)
			}
			continue
		}
		if !(status == 200) {
			// Increment the counter :(
			c.ESIRateLimit.Inc()
			// fmt.Printf("ESI GET RESPONSE ERROR - %v - %v - %v\n", status, url, string(body))
			time.Sleep(250 * time.Millisecond)
			continue
		}

		return body, err
	}

	return nil, fmt.Errorf("Max retries exceeded for url: ; err: %v", url, err)
}

func (c *Client) MakeZKBGet(url string) ([]byte, error) {

	retriesRemain := c.RetryLimit
	for retriesRemain > 1 {
		retriesRemain--

		for c.ZKBRateLimit.Value() > 10 {
			time.Sleep(500 * time.Millisecond)
		}

		body, status, err := c.makeRawHTTPGet(url)
		if err != nil {
			continue
		}
		if !(status == 200) {
			// Increment the counter :(
			c.ZKBRateLimit.Inc()
			time.Sleep(250 * time.Millisecond)
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
			time.Sleep(500 * time.Millisecond)
		}

		body, status, err := c.makeRawHTTPGet(url)
		if err != nil {
			continue
		}
		if !(status >= 200 && status < 300) {
			time.Sleep(250 * time.Millisecond)
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
	if c.cnt > 0 {
		c.cnt--
	}
	c.mux.Unlock()
}

// Value returns the current value of the counter.
func (c *safeCounter) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.cnt
}
