package mothistoryapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultHost = "https://beta.check-mot.service.gov.uk"
)

// Client is an API Client for MOT History API
type Client struct {
	apiKey   string
	baseHost string
	client   *http.Client
}

// NewClient returns a new MOT History API Client
func NewClient(apiKey, host string) *Client {
	if host == "" {
		host = defaultHost
	}

	return &Client{
		apiKey:   apiKey,
		baseHost: host,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetVehicleHistory fetches the MOT History for a specific vehicle
func (c *Client) GetVehicleHistory(registrationNumber string) (*Vehicle, error) {
	requestURL := c.baseHost + "/trade/vehicles/mot-tests"

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("registration", registrationNumber)
	req.URL.RawQuery = q.Encode()

	res := Vehicles{}
	err = c.sendRequest(req, &res)

	if err != nil {
		return nil, err
	}

	return &res[0], nil
}

func (c *Client) sendRequest(req *http.Request, vehicles *Vehicles) error {
	req.Header.Set("Accept", "application/json+v6")
	req.Header.Set("X-Api-Key", c.apiKey)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, body)
	}

	return json.NewDecoder(res.Body).Decode(&vehicles)
}
