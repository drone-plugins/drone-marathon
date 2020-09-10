package marathon

type Marathon struct {
	Branch            string `json:"branch"`
	BuildNumber       string `json:"build_number"`
	CommitSha         string `json:"commit_sha"`
	CommitAuthor      string `json:"commit_author"`
	CommitAuthorEmail string `json:"commit_author_email"`
	CommitBranch      string `json:"commit_branch"`
	CommitLink        string `json:"commit_link"`
	Debug             bool   `json:"debug"`
  DeployTo          string `json:"deploy_to"`
	GroupName         string `json:"group_name"`
	MarathonFile      string `json:"marathonfile"`
	Server            string `json:"server"`
	Tag               string `json:"tag"`
}
