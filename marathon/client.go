package marathon

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

type Client interface {
	GenerateApplicationDefinition() (string, error)
	CreateOrUpdateApplication() error
}

type client struct {
	params *Marathon
}

func NewClient(params *Marathon) Client {
	return &client{params}
}

func (c *client) GenerateApplicationDefinition() (string, error) {
	c.extractProcessEnvKeys()
	c.extractDockerParamsKeys()
	c.normalizeDockerPortMappings()

	tmpl, err := template.New("marathon").Parse(applicationDefinition)
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

func (c *client) extractProcessEnvKeys() {
	c.params.ProcessEnvKeys = []string{}
	for k := range c.params.ProcessEnv {
		c.params.ProcessEnvKeys = append(c.params.ProcessEnvKeys, k)
	}
}

func (c *client) extractDockerParamsKeys() {
	c.params.DockerParamsKeys = []string{}
	for k := range c.params.DockerParams {
		c.params.DockerParamsKeys = append(c.params.DockerParamsKeys, k)
	}
}

func (c *client) normalizeDockerPortMappings() {
	for i, each := range c.params.DockerPortMappings {
		if each.Protocol == "" {
			c.params.DockerPortMappings[i].Protocol = "tcp"
		}
	}
}

func (c *client) CreateOrUpdateApplication() error {
	appDefinition, err := c.GenerateApplicationDefinition()
	if err != nil {
		return err
	}

	if c.params.Debug {
		fmt.Println("Request", appDefinition)
	}

	payload := []byte(appDefinition)

	status, body, err := c.sendToServer("POST", "/v2/apps", payload)
	if err != nil {
		return err
	}

	if status == http.StatusConflict {
		status, body, err = c.sendToServer("PUT", "/v2/apps/"+c.params.ID, payload)
		if err != nil {
			return err
		}
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
	req, _ := http.NewRequest(action, c.params.Server+path, bytes.NewBuffer(definition))
	if c.params.Username != "" && c.params.Password != "" {
		req.SetBasicAuth(c.params.Username, c.params.Password)
	}
	req.Header.Set("Content-Type", "application/json")
	httpclient := &http.Client{}
	resp, err := httpclient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}
