package marathon

import (
	"bytes"
	"fmt"
  "log"
	"io/ioutil"
	"net/http"
)

type Client interface {
	ReadMarathonFile() ([]byte, error)
	CreateOrUpdateApplication() error
}

type client struct {
	params *Marathon
}

func NewClient(params *Marathon) Client {
	return &client{params}
}

func (c *client) ReadMarathonFile() ([]byte, error) {
  log.Printf("Params: %v", c.params)
  data, err := ioutil.ReadFile(c.params.MarathonFile)
  if err != nil {
    return nil, err
  }
	return data, nil
}

func (c *client) CreateOrUpdateApplication() error {
	groupDefinition, err := c.ReadMarathonFile()
	if err != nil {
		return err
	}

	if c.params.Debug {
		fmt.Println("Request", string(groupDefinition))
	}

  status, body, err := c.sendToServer("PUT", "/v2/groups/", groupDefinition)
  if err != nil {
    return err
  }

	if c.params.Debug {
		if err == nil {
			fmt.Printf("Response (%v): %s\n", status, body)
		} else {
			fmt.Println("Error:", err)
		}
	}

	if err == nil && status != http.StatusCreated && status != http.StatusOK {
		err = fmt.Errorf("Response (%v): %s\n", status, body)
	}

	return err
}

func (c *client) sendToServer(action, path string, definition []byte) (int, string, error) {
	req, err := http.NewRequest(action, c.params.Server+path, bytes.NewBuffer(definition))
	req.Header.Set("Content-Type", "application/json")
	httpclient := &http.Client{}
	resp, err := httpclient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}
