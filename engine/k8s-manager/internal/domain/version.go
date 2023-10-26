package domain

type Version struct {
	Tag                  string
	Product              string
	GlobalKeyValueStore  string
	VersionKeyValueStore string

	Workflows []*Workflow
}

type Workflow struct {
	Name          string
	Stream        string
	KeyValueStore string
	Processes     []*Process
}
