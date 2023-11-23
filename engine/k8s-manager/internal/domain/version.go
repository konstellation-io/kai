package domain

type Version struct {
	Tag                  string
	Product              string
	GlobalKeyValueStore  string
	VersionKeyValueStore string

	Workflows          []*Workflow
	MinioConfiguration MinioConfiguration
	ServiceAccount     ServiceAccount
}

type Workflow struct {
	Name          string
	Stream        string
	KeyValueStore string
	Processes     []*Process
}

type MinioConfiguration struct {
	Bucket string
}

type ServiceAccount struct {
	Username string
	Password string
}
