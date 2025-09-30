package rickandmorty

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func New(baseURL string) *Client {
	return &Client{
		cli:     http.DefaultClient,
		baseURL: baseURL,
	}
}

type Client struct {
	cli     *http.Client
	baseURL string
}

func (c *Client) GetCaracters() (*ResponseCaracters, error) {
	var payload ResponseCaracters
	err := c.do(http.MethodGet, "/api/character", &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func (c *Client) do(method, uri string, payload interface{}) error {
	r, err := http.NewRequest(method, c.baseURL+uri, nil)
	if err != nil {
		return err
	}

	resp, err := c.cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("sdk rick and morty: uri: %v , status not valid: %v", uri, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, payload)
	if err != nil {
		return err
	}

	return nil
}
