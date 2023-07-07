package version

var (
	BuildDate  string
	Version    string
	Commit     string
	CommitDate string
	Branch     string
	Tag        string
	Dirty      string
)

type VersionInfo struct {
	BuildDate  string `json:"build_date"`
	Version    string `json:"build_version"`
	Commit     string `json:"build_commit"`
	CommitDate string `json:"build_commit_date"`
	Branch     string `json:"build_commit_branch"`
	Tag        string `json:"build_tag"`
	Dirty      string `json:"build_dirty"`
}

func New() *VersionInfo {
	return &VersionInfo{
		BuildDate:  BuildDate,
		Version:    Version,
		Commit:     Commit,
		CommitDate: CommitDate,
		Branch:     Branch,
		Tag:        Tag,
		Dirty:      Dirty,
	}
}
