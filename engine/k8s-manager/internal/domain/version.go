package domain

type Version struct {
	Tag                  string
	Product              string
	GlobalKeyValueStore  string
	VersionKeyValueStore string

	Workflows          []*Workflow
	MinioConfiguration MinioConfiguration
}

type Workflow struct {
	Name          string
	Stream        string
	KeyValueStore string
	Processes     []*Process
}

type MinioConfiguration struct {
	User     string
	Password string
	Bucket   string
}
