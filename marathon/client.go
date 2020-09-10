package marathon

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Client interface {
	ReadMarathonFile() (string, error)
	CreateOrUpdateApplication() error
}

type client struct {
	params *Marathon
}

func NewClient(params *Marathon) Client {
	return &client{params}
}

func (c *client) ReadMarathonFile() (string, error) {
	if c.params.Debug {
		log.Printf("Params: %v", c.params)
	}
	data, err := ioutil.ReadFile(c.params.MarathonFile)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("marathon").Parse(string(data))
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, *c.params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *client) CreateOrUpdateApplication() error {
	groupDefinition, err := c.ReadMarathonFile()
	if err != nil {
		return err
	}

	if c.params.Debug {
		fmt.Println("Request", groupDefinition)
	}

	if c.params.GroupName == "" {
		return errors.New("Group name is required")
	}

	status, body, err := c.sendToServer("PUT", "/v2/groups/"+c.params.GroupName+"?force=true", []byte(groupDefinition))
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
