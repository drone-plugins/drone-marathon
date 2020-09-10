package marathon

type Marathon struct {
	Server            string `json:"server"`
	MarathonFile      string `json:"marathonfile"`
	GroupName         string `json:"group_name"`
	Branch            string `json:"branch"`
	BuildNumber       string `json:"build_number"`
	Commit            string `json:"commit"`
	CommitAuthor      string `json:"commit_author"`
	CommitAuthorEmail string `json:"commit_author_email"`
	CommitBranch      string `json:"commit_branch"`
	CommitLink        string `json:"commit_link"`
	Tag               string `json:"tag"`
	Debug             bool   `json:"debug"`
}
