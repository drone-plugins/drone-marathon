package marathon

type Marathon struct {
	Server       string `json:"server"`
	MarathonFile string `json:"marathonfile"`
	GroupName    string `json:"group_name"`
  Branch       string `json:"branch"`
  BuildNumber  string `json:"build_number"`
	Debug        bool   `json:"debug"`
}
