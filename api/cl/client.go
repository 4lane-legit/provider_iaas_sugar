package cl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iaas_sugar/api/sr"
	"io"
	"net/http"
)

// Client holds all of the information required to connect to a server
type Client struct {
	hostname   string
	port       int
	authToken  string
	httpClient *http.Client
}

// NewClient returns a new client configured to communicate on a server with the
// given hostname and port and to send an Authorization Header with the value of
// token
func NewClient(hostname string, port int, token string) *Client {
	return &Client{
		hostname:   hostname,
		port:       port,
		authToken:  token,
		httpClient: &http.Client{},
	}
}

// GetAll Retrieves all of the Minions from the server
func (c *Client) GetAll() (*map[string]sr.Minion, error) {
	body, err := c.httpRequest("minions", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	minionss := map[string]sr.Minion{}
	err = json.NewDecoder(body).Decode(&minionss)
	if err != nil {
		return nil, err
	}
	return &minionss, nil
}

// GetMinion gets an minions with a specific name from the server
func (c *Client) GetMinion(name string) (*sr.Minion, error) {
	body, err := c.httpRequest(fmt.Sprintf("minions/%v", name), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	minions := &sr.Minion{}
	err = json.NewDecoder(body).Decode(minions)
	if err != nil {
		return nil, err
	}
	return minions, nil
}

// NewMinion creates a new Minion
func (c *Client) NewMinion(minions *sr.Minion) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(minions)
	if err != nil {
		return err
	}
	_, err = c.httpRequest("minions", "POST", buf)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMinion updates the values of an minions
func (c *Client) UpdateMinion(minions *sr.Minion) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(minions)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("minions/%s", minions.Name), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMinion removes an minions from the server
func (c *Client) DeleteMinion(minionsName string) error {
	_, err := c.httpRequest(fmt.Sprintf("minions/%s", minionsName), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) httpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authToken)
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%s:%v/%s", c.hostname, c.port, path)
}
