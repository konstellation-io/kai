package domain

type Version struct {
	ID            string
	Product       string
	KeyValueStore string

	Workflows []*Workflow
}

type Workflow struct {
	ID            string
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
