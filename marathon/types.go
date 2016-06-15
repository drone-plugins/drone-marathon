package marathon

type Marathon struct {
  Server                string      `json:"server"`
  Username              string      `json:"username"`
  Password              string      `json:"password"`
  ID                    string      `json:"id"`
  Instances             int         `json:"instances"`
  Cpus                  float64     `json:"cpus"`
  Mem                   float64     `json:"mem"` 
  Cmd                   string      `json:"cmd"`
  Args                  []string    `json:"args"`
  Uris                  []string    `json:"uris"`
  Fetchs                []Fetch     `json:"fetch"` 
  MinHealthCapacity     float64     `json:"min_health_capacity"`
  MaxOverCapacity       float64     `json:"max_over_capacity"`
  HealthChecks          []HealthCheck       `json:"health_checks"`
  Constraints           []Constraint        `json:"constraints"`
  AcceptedResourceRoles []string            `json:"accepted_resource_roles"`
  BackoffFactor         float64             `json:"backoff_factor"`
  BackoffSeconds        int                 `json:"backoff_seconds"`
  MaxLaunchDelaySeconds int                 `json:"max_launch_delay_seconds"`
  Dependencies          []string            `json:"dependencies"`
  Disk                  float64             `json:"disk"`
  ProcessEnv            map[string]string   `json:"process_environment"`
  ProcessEnvKeys        []string
  DockerImage           string              `json:"docker_image"`
  DockerNetwork         string              `json:"docker_network"`
  DockerForcePull       bool                `json:"docker_force_pull"`
  DockerPrivileged      bool                `json:"docker_privileged"`
  DockerPortMappings    []DockerPortMapping `json:"docker_port_mappings"`
  DockerVolumes         []DockerVolume      `json:"docker_volumes"`
  DockerParams          map[string]string   `json:"docker_parameters"`
  DockerParamsKeys      []string
  Debug                 bool                `json:"debug"`
}

type HealthCheck struct {
  Protocol               string `json:"protocol"`
  Path                   string `json:"path"`
  Command                string `json:"command"`
  GracePeriodSeconds     int    `json:"grace_period_seconds"`
  IntervalSeconds        int    `json:"interval_seconds"`
  Port                   int    `json:"port"`
  PortIndex              int    `json:"port_index"`
  TimeoutSeconds         int    `json:"timeout_seconds"`
  MaxConsecutiveFailures int    `json:"max_consecutive_failures"`
}

type Constraint struct {
  Field    string `json:"field"`
  Operator string `json:"operator"`
  Value    string `json:"value"`
}

type Fetch struct {
  Uri        string `json:"uri"`
  Executable bool   `json:"executable"`
  Extract    bool   `json:"extract"`
  Cache      bool   `json:"cache"`
}

type DockerPortMapping struct {
  ContainerPort int    `json:"container_port"`
  HostPort      int    `json:"host_port"`
  ServicePort   int    `json:"service_port"`
  Protocol      string `json:"protocol"`
}

type DockerVolume struct {
  ContainerPath string `json:"container_path"`
  HostPath      string `json:"host_path"`
  Mode          string `json:"mode"`
}

const applicationDefinition = `
{
  "id": "{{.ID}}",
  "cpus": {{.Cpus}},
  "mem": {{.Mem}},
  "instances": {{.Instances}},
  {{if .Cmd}}
  "cmd": "{{.Cmd}}",
  {{end}}
  {{if .Args}}
  "args": [
    {{range $index, $element := .Args}}
        {{if $index}},{{end}}
        "{{$element}}"
    {{end}}      
  ],
  {{end}}
  {{if .DockerImage}}
  "container": {
    "type": "DOCKER",
    "docker": {
      "image": "{{.DockerImage}}",
      "network": "{{.DockerNetwork}}",
      {{if .DockerPortMappings}}
      "portMappings": [
        {{range $index, $element := .DockerPortMappings}}
          {{if $index}},{{end}}
          {
            "containerPort": {{$element.ContainerPort}},
            "hostPort": {{$element.HostPort}},
            "servicePort": {{$element.ServicePort}},
            "protocol": "{{$element.Protocol}}"
          }
        {{end}}
      ],
      {{end}}
      {{if .DockerVolumes}}
      "volumes": [
        {{range $index, $element := .DockerVolumes}}
          {{if $index}},{{end}}
          {
            "containerPath": "{{$element.ContainerPath}}",              
            "hostPath": "{{$element.HostPath}}",
            "mode": "{{$element.Mode}}"
          }
        {{end}}
      ],
      {{end}}
      {{if .DockerParamsKeys}}
      "parameters": [
        {{range $index, $key := .DockerParamsKeys}}
          {{if $index}},{{end}}
          { "key": "{{$key}}", "value": "{{index $.DockerParams $key}}" }  
        {{end}}         
      ],
      {{end}}
      "forcePullImage": {{.DockerForcePull}},
      "privileged": {{.DockerPrivileged}}
    }
  },
  {{end}}  
  {{if .ProcessEnvKeys}}
  "env": {
    {{range $index, $key := .ProcessEnvKeys}}
        {{if $index}},{{end}}  
        "{{$key}}": "{{index $.ProcessEnv $key}}"
    {{end}}
  },
  {{end}} 
  {{if .HealthChecks}}
  "healthChecks": [
    {{range $index, $element := .HealthChecks}}
      {{if $index}},{{end}}
      {
        "protocol": "{{$element.Protocol}}",
        {{if $element.Path}}"path": "{{$element.Path}}",{{end}}
        {{if $element.PortIndex}}"portIndex": {{$element.PortIndex}},{{end}}
        {{if $element.Port}}"port": {{$element.Port}},{{end}}
        {{if $element.Command}}"command": { "value": "{{$element.Command}}" },{{end}}
        "gracePeriodSeconds": {{$element.GracePeriodSeconds}},
        "intervalSeconds": {{$element.IntervalSeconds}},
        "timeoutSeconds": {{$element.TimeoutSeconds}},
        "maxConsecutiveFailures": {{$element.MaxConsecutiveFailures}}
      }
    {{end}}
  ],
  {{end}}
  {{if .Uris}}
  "uris": [
    {{range $index, $element := .Uris}}
        {{if $index}},{{end}}
        "{{$element}}"
    {{end}}      
  ],
  {{end}}
  {{if .Fetchs}}
  "fetch": [
    {{range $index, $element := .Fetchs}}
      {{if $index}},{{end}}
      {
        "uri": "{{$element.Uri}}",              
        "executable": {{$element.Executable}},
        "extract": {{$element.Extract}},
        "cache": {{$element.Cache}}
      }
    {{end}}
  ],
  {{end}}
  {{if .Constraints}}
  "constraints": [
    {{range $index, $element := .Constraints}}
      {{if $index}},{{end}}
      ["{{$element.Field}}","{{$element.Operator}}"{{if $element.Value}},"{{$element.Value}}"{{end}}]
    {{end}}      
  ],  
  {{end}}
  {{if .AcceptedResourceRoles}}
  "acceptedResourceRoles": [
    {{range $index, $element := .AcceptedResourceRoles}}
        {{if $index}},{{end}}
        "{{$element}}"
    {{end}}      
  ],
  {{end}}
  {{if .Dependencies}}
  "dependencies": [
    {{range $index, $element := .Dependencies}}
        {{if $index}},{{end}}
        "{{$element}}"
    {{end}}      
  ],
  {{end}}
  "disk": {{.Disk}},
  "backoffFactor": {{.BackoffFactor}},
  "backoffSeconds": {{.BackoffSeconds}},
  "maxLaunchDelaySeconds": {{.MaxLaunchDelaySeconds}},
  "upgradeStrategy": {
    "minimumHealthCapacity": {{.MinHealthCapacity}},
    "maximumOverCapacity": {{.MaxOverCapacity}}
  }
}
`
