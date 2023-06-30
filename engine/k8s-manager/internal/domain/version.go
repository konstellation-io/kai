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

func (v Version) GetAmountOfProcesses() int {
	amount := 0
	for _, w := range v.Workflows {
		amount += len(w.Processes)
	}

	return amount
}
