package domain

type Version struct {
	Name          string
	Product       string
	KeyValueStore string

	Workflows []*Workflow
}

type Workflow struct {
	Name          string
	Stream        string
	KeyValueStore string
	Processes     []*Process
}
