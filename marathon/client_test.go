package marathon

import (
	"strings"
	"testing"
	"unicode"
	"net/http"
	"net/http/httptest"
)

func TestGenerateSimple(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.ID = "foobar"
	m.Cpus = 0.10
	m.Mem = 128.1
	m.Instances = 10
    m.Cmd = "cmd foo bar"
	m.BackoffFactor = 1.54
	m.BackoffSeconds = 2
	m.MaxLaunchDelaySeconds = 600
	m.Disk = 1.9

	var tests = []struct {
		expected string
		name     string
	}{
		{`"id":"foobar"`, "ID"},
		{`"cpus":0.1`, "Cpus"},
		{`"mem":128.1`, "Mem"},
		{`"instances":10`, "Instances"},
        {`"cmd":"cmdfoobar"`, "Cmd"},
		{`"backoffFactor":1.54`, "BackoffFactor"},
		{`"backoffSeconds":2`, "BackoffSeconds"},
		{`"maxLaunchDelaySeconds":600`, "MaxLaunchDelaySeconds"},
		{`"disk":1.9`, "Disk"},
	}

	conf := generateAppDef(client, t)

	for _, test := range tests {
		if !strings.Contains(conf, test.expected) {
			t.Errorf("%s not found %s", test.name, conf)
		}
	}
}

func TestNotGenerateOnEmptyParams(t *testing.T) {
	m := new(Marathon)

	mustNotContain(m, `"args":`, t)
	mustNotContain(m, `"uris":`, t)
	mustNotContain(m, `"healthChecks":`, t)
	mustNotContain(m, `"constraints":`, t)
	mustNotContain(m, `"acceptedResourceRoles":`, t)
	mustNotContain(m, `"portMappings":`, t)
	mustNotContain(m, `"volumes":`, t)
	mustNotContain(m, `"env":{`, t)
	mustNotContain(m, `"parameters":[`, t)
	mustNotContain(m, `"dependencies":[`, t)
	mustNotContain(m, `"fetch":`, t)
	mustNotContain(m, `"container": {`, t)
}

func TestGenerateArgs(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.Args = []string{"foo", "bar"}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"args":["foo","bar"]`) {
		t.Errorf("args not found %s", conf)
	}
}

func TestGenerateUris(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.Uris = []string{"foo", "bar"}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"uris":["foo","bar"]`) {
		t.Errorf("uris not found %s", conf)
	}
}

func TestGenerateFetch(t *testing.T)  {
	m := new(Marathon)
	client := NewClient(m)
	m.Fetchs = []Fetch{
		{Uri: "https://foo.com/archive.zip",
		Executable: false, Extract: true, Cache: false},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"fetch":[{"uri":"https://foo.com/archive.zip","executable":false,"extract":true,"cache":false}]`) {
		t.Errorf("fetch not found %s", conf)
	}
}

func TestGenerateUpgradeStrategy(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.MinHealthCapacity = 0.5
	m.MaxOverCapacity = 0.8

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"upgradeStrategy":{"minimumHealthCapacity":0.5,"maximumOverCapacity":0.8}`) {
		t.Errorf("upgradeStrategy not found %s", conf)
	}
}

func TestGenerateHealthCheckHttp(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.HealthChecks = []HealthCheck{
		{Protocol: "HTTP", Path: "/health", GracePeriodSeconds: 3,
			IntervalSeconds: 10, PortIndex: 1, TimeoutSeconds: 10,
			MaxConsecutiveFailures: 3},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"healthChecks":[{"protocol":"HTTP","path":"/health","portIndex":1,"gracePeriodSeconds":3,"intervalSeconds":10,"timeoutSeconds":10,"maxConsecutiveFailures":3}]`) {
		t.Errorf("healthChecks http not found %s", conf)
	}
}

func TestGenerateHealthCheckCommand(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.HealthChecks = []HealthCheck{
		{Protocol: "COMMAND", Command: "curl http://foo.bar", GracePeriodSeconds: 3,
			IntervalSeconds: 10, TimeoutSeconds: 10,
			MaxConsecutiveFailures: 3},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"healthChecks":[{"protocol":"COMMAND","command":{"value":"curlhttp://foo.bar"},"gracePeriodSeconds":3,"intervalSeconds":10,"timeoutSeconds":10,"maxConsecutiveFailures":3}]`) {
		t.Errorf("healthChecks command not found %s", conf)
	}
}

func TestNormalizeHealthCheck(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.HealthChecks = []HealthCheck{
		{Protocol: "HTTP", Path: "/health"},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"healthChecks":[{"protocol":"HTTP","path":"/health","gracePeriodSeconds":0,"intervalSeconds":0,"timeoutSeconds":0,"maxConsecutiveFailures":0}]`) {
		t.Errorf("healthChecks not found %s", conf)
	}
}

func TestGenerateConstraints(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.Constraints = []Constraint{
		{Field: "hostname", Operator: "CLUSTER", Value: "foobar"},
		{Field: "hostname", Operator: "UNIQUE"},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"constraints":[["hostname","CLUSTER","foobar"],["hostname","UNIQUE"]]`) {
		t.Errorf("constraints not found %s", conf)
	}
}

func TestGenerateAcceptedResourceRoles(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.AcceptedResourceRoles = []string{"foo", "bar"}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"acceptedResourceRoles":["foo","bar"]`) {
		t.Errorf("acceptedResourceRoles not found %s", conf)
	}
}

func TestGenerateDependencies(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.Dependencies = []string{"/product/db/mongo", "/product/db", "../../db"}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"dependencies":["/product/db/mongo","/product/db","../../db"],`) {
		t.Errorf("dependencies not found %s", conf)
	}
}

func TestGenerateProcessEnv(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.ProcessEnv = make(map[string]string)
	m.ProcessEnv["foo"] = "bar"
	m.ProcessEnv["STATUS"] = "EMPTY"

	conf := generateAppDef(client, t)

	if !(strings.Contains(conf, `"env":{"STATUS":"EMPTY","foo":"bar"}`) ||
		strings.Contains(conf, `"env":{"foo":"bar","STATUS":"EMPTY"}`)) {
		t.Errorf("environment not found %s", conf)
	}
}

func TestGenerateDocker(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.DockerImage = "registry.local/test-server:1.0"
	m.DockerForcePull = true
	m.DockerPrivileged = true
	m.DockerNetwork = "BRIDGE"

	var tests = []struct {
		expected string
		name     string
	}{
		{`"image":"registry.local/test-server:1.0"`, "DockerImage"},
		{`"forcePullImage":true`, "DockerForcePull"},
		{`"privileged":true`, "DockerPrivileged"},
		{`"network":"BRIDGE"`, "DockerNetwork"},
	}

	conf := generateAppDef(client, t)

	for _, test := range tests {
		if !strings.Contains(conf, test.expected) {
			t.Errorf("%s not found %s", test.name, conf)
		}
	}
}

func TestGenerateDockerPortMap(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.DockerImage = "someimg"
	m.DockerPortMappings = []DockerPortMapping{
		{ContainerPort: 3000, HostPort: 0, ServicePort: 0, Protocol: "udp"},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"portMappings":[{"containerPort":3000,"hostPort":0,"servicePort":0,"protocol":"udp"}]`) {
		t.Errorf("portMappings not found %s", conf)
	}
}

func TestNormalizeGenerateDockerPortMap(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.DockerImage = "someimg"
	m.DockerPortMappings = []DockerPortMapping{
		{ContainerPort: 0, HostPort: 0, ServicePort: 0, Protocol: ""},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"portMappings":[{"containerPort":0,"hostPort":0,"servicePort":0,"protocol":"tcp"}]`) {
		t.Errorf("portMappings normalization failed %s", conf)
	}
}

func TestGenerateDockerVolumes(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.DockerImage = "someimg"
	m.DockerVolumes = []DockerVolume{
		{ContainerPath: "/foo", HostPath: "/bar", Mode: "RW"},
	}

	conf := generateAppDef(client, t)

	if !strings.Contains(conf, `"volumes":[{"containerPath":"/foo","hostPath":"/bar","mode":"RW"}]`) {
		t.Errorf("volumes not found %s", conf)
	}
}

func TestGenerateDockerParams(t *testing.T) {
	m := new(Marathon)
	client := NewClient(m)
	m.DockerImage = "someimg"	
	m.DockerParams = make(map[string]string)
	m.DockerParams["foo"] = "bar"
	m.DockerParams["STATUS"] = "EMPTY"

	conf := generateAppDef(client, t)

	if !(strings.Contains(conf, `"parameters":[{"key":"STATUS","value":"EMPTY"},{"key":"foo","value":"bar"}]`) ||
		strings.Contains(conf, `"parameters":[{"key":"foo","value":"bar"},{"key":"STATUS","value":"EMPTY"}]`)) {
		t.Errorf("parameters not found %s", conf)
	}
}

func TestMarathonStatusCodes(t *testing.T)  {
	NoError := false
	Error := true
	checkRequest(http.StatusCreated, NoError, t)
	checkRequest(http.StatusOK, NoError, t)
	checkRequest(http.StatusBadRequest, Error, t)
	checkRequest(http.StatusUnauthorized, Error, t)
	checkRequest(http.StatusForbidden, Error, t)
	checkRequest(422, Error, t)
}

func checkRequest(statusCode int, expectError bool, t *testing.T)  {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
	defer ts.Close()
	m := new(Marathon)
	m.Server = ts.URL
	client := NewClient(m)
	err := client.CreateOrUpdateApplication()
	if (expectError && err == nil) || (!expectError && err != nil) {
		t.Error(err)
	}
}

func generateAppDef(c Client, t *testing.T) string {
	conf, err := c.GenerateApplicationDefinition()
	if err != nil {
		t.Error(err)
	}
	conf = removeSpaces(conf)
	return conf
}

func removeSpaces(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func mustNotContain(m *Marathon, stringToSearch string, t *testing.T) {
	client := NewClient(m)

	conf := generateAppDef(client, t)

	if strings.Contains(conf, stringToSearch) {
		t.Errorf("%s found in %s", stringToSearch, conf)
	}
}
