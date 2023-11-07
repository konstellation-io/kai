package compensator

import "errors"

type Compensation func() error

type Compensator struct {
	compensations []Compensation
}

func New() *Compensator {
	return &Compensator{
		compensations: []Compensation{},
	}
}

func (c *Compensator) AddCompensation(compensation Compensation) {
	c.compensations = append(c.compensations, compensation)
}

func (c *Compensator) Execute() error {
	var errs error

	for i := len(c.compensations) - 1; i >= 0; i-- {
		err := c.compensations[i]()
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
